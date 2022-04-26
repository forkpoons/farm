package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type repository interface {
	WriteTemperature(name string, temp int) error
	WriteAction(name string, temp bool) error
	ReadTemperature(name string) (float64, error)
}

type worker struct {
	addr string
	repo repository
	cons map[string]*websocket.Conn
}

var conn map[string]*websocket.Conn

func New(addr string, repo repository) (*worker, error) {
	return &worker{
		addr: addr,
		repo: repo,
		cons: make(map[string]*websocket.Conn),
	}, nil
}

func (w *worker) Start() error {
	conn = make(map[string]*websocket.Conn)
	http.HandleFunc("/", w.get)
	http.HandleFunc("/api/temp", w.getTemperature)
	http.HandleFunc("/api/settemp", w.postTemperature)
	http.HandleFunc("/api/echo", w.echo)
	return http.ListenAndServe(w.addr, nil)
}

func (w *worker) get(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.Header)
	_, err := rw.Write([]byte("ok"))
	if err != nil {
		log.Println(err)
		return
	}
}

func (w *worker) getTemperature(rw http.ResponseWriter, req *http.Request) {
	//if err := w.repo.WriteTemperature(1); err != nil {
	//	вот так будешь вызывать потом
	//}

	rw.Header().Set("Content-Type", "application/json")
	temp, err := w.repo.ReadTemperature("qwe")
	if err != nil {
		log.Println(err)
	}
	_, err = rw.Write([]byte(fmt.Sprintf("%f", temp)))
	if err != nil {
		log.Println(err)
		return
	}
}

func (w *worker) postTemperature(rw http.ResponseWriter, req *http.Request) {
	//if err := w.repo.WriteTemperature(1); err != nil {
	//	вот так будешь вызывать потом
	//}
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	nn := 0
	min := 0
	max := 0
	if body[0] == '#' && body[1] == '#' {
		body = body[2:]
		for n, m := range body {
			//fmt.Println(nn, "/", n, "/", string(m))
			if m == '|' {
				min, err = strconv.Atoi(string(body[nn:n]))
				nn = n + 1
			} else if m == '#' {
				max, err = strconv.Atoi(string(body[nn:n]))
				break
			}
		}
	}
	rw.Header().Set("Content-Type", "application/json")
	if err == nil {
		for _, con := range conn {
			err := con.WriteMessage(1, []byte("##"+strconv.Itoa(min)+"|"+strconv.Itoa(max)+"#"))
			if err != nil {
				log.Println(err)
			}
		}
		_, err = rw.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println(err)
		_, err = rw.Write([]byte("error"))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

var upgrader = websocket.Upgrader{} // use default options

func (w *worker) echo(rw http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(rw, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	_, message, err := c.ReadMessage()
	err = c.WriteMessage(1, []byte("##10|15#"))
	conn[string(message)] = c
	defer delete(conn, string(message))
	for {
		_, message, err := c.ReadMessage()
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
					value, err := strconv.ParseFloat(msgValue[2], 8)
					if err != nil {
						log.Println(err)
						break
					}
					if err := w.repo.WriteTemperature(msgValue[0], int(value)); err != nil {
						log.Println(err)
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
