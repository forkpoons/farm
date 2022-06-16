package greenhouse

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"time"
)

type Device struct {
	onLine      bool
	work        bool
	Temperature float64
	Humidity    float64
	lastOnline  time.Time
	lastData    time.Time
}

type worker struct {
	repo        repo
	devices     map[string]*Device
	connections map[string]*websocket.Conn
}

func New(db *sqlx.DB) *worker {
	return &worker{
		repo:        *newDB(context.Background(), db),
		devices:     make(map[string]*Device),
		connections: make(map[string]*websocket.Conn),
	}
}

func (w *worker) Handlers(gh *gin.RouterGroup) {
	gh.GET("/echo", w.echo)
}

var upgrader = websocket.Upgrader{}

func (d *Device) offlineDevice() {
	d.onLine = false
	d.lastOnline = time.Now()
}

func (w *worker) echo(c *gin.Context) {
	con, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer con.Close()

	_, message, err := con.ReadMessage()
	if string(message[0:5]) != "name:" {
		if err := con.WriteMessage(1, []byte("error")); err != nil {
			return
		}
		if err = con.Close(); err != nil {
			return
		}
	}
	err = con.WriteMessage(1, []byte("##10|15#"))
	w.connections[string(message)] = con
	if v, ok := w.devices[string(message)]; ok {
		v.onLine = true
	} else {
		w.devices[string(message)] = &Device{onLine: true, work: false, Temperature: 99.0, Humidity: 99.0, lastOnline: time.Now(), lastData: time.Now()}
	}
	defer delete(w.connections, string(message))
	defer w.devices[string(message)].offlineDevice()
	fmt.Println(w.devices[string(message)])
	for {
		_, message, err := con.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if len(message) > 0 {
			msg := string(message)
			msgValue := make([]string, 0)
			nn := 0
			fmt.Println("Message:" + msg)
			// ##id|typemsg|value#
			if msg[0] == '#' && msg[1] == '#' {
				msg = msg[2:]
				for n, m := range msg {
					//fmt.Println(nn, "/", n, "/", string(m))
					if m == '|' {
						msgValue = append(msgValue, msg[nn:n])
						nn = n + 1
					} else if m == '#' {
						msgValue = append(msgValue, msg[nn:n])
						break
					}
				}
				fmt.Println(msgValue)
				switch msgValue[1] {
				case "temp":
					if len(msgValue) >= 3 {
						valueTemp, err := strconv.ParseFloat(msgValue[2], 8)
						if err != nil {
							log.Println(err)
							break
						}
						if err := w.repo.WriteTemperature(msgValue[0], valueTemp); err != nil {
							log.Println(err)
						}
						w.devices[msgValue[0]].Temperature = valueTemp
						w.devices[msgValue[0]].lastData = time.Now()
						if len(msgValue) >= 4 {
							valueH, err := strconv.ParseFloat(msgValue[3], 8)
							if err != nil {
								log.Println(err)
								break
							}
							w.devices[msgValue[0]].Humidity = valueH
						}
					}
				case "action":
					fmt.Println("action")
					if err := w.repo.WriteAction(msgValue[0], msgValue[2] == "1"); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
