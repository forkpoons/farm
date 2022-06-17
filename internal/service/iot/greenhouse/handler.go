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
	lastTemp    time.Time
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
	fmt.Println(d)
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
	deviceName := string(message[5:])
	err = con.WriteMessage(1, []byte("##10|15#"))
	w.connections[deviceName] = con
	if v, ok := w.devices[deviceName]; ok {
		v.onLine = true
	} else {
		w.devices[deviceName] = &Device{onLine: true, work: false, Temperature: 99.0, Humidity: 99.0, lastOnline: time.Now(), lastTemp: time.Now()}
	}
	defer delete(w.connections, deviceName)
	defer w.devices[deviceName].offlineDevice()
	//fmt.Println(w.devices[deviceName])
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
			// ##typeMsg|value#
			if msg[0] == '#' && msg[1] == '#' {
				msg = msg[2:]
				for n, m := range msg {
					if m == '|' {
						msgValue = append(msgValue, msg[nn:n])
						nn = n + 1
					} else if m == '#' {
						msgValue = append(msgValue, msg[nn:n])
						break
					}
				}
				fmt.Println(msgValue)
				if len(msgValue) >= 2 {
					switch msgValue[0] {
					case "temp":
						valueTemp, err := strconv.ParseFloat(msgValue[1], 8)
						if err != nil {
							log.Println(err)
							break
						}
						w.devices[deviceName].Temperature = valueTemp
						w.devices[deviceName].lastTemp = time.Now()
						if err := w.repo.WriteTemperature(deviceName, valueTemp); err != nil {
							log.Println(err)
						}
					case "humidity":
						valueH, err := strconv.ParseFloat(msgValue[1], 8)
						if err != nil {
							log.Println(err)
							break
						}
						w.devices[deviceName].Humidity = valueH
					case "action":
						fmt.Println("action" + msgValue[1])
						w.devices[deviceName].work = msgValue[1] == "1"
						if err := w.repo.WriteAction(deviceName, msgValue[1] == "1"); err != nil {
							log.Println(err)
						}
					}
				}
				fmt.Println(w.devices[deviceName])
			}
		}
	}
}
