package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

type repository interface {
	WriteTemperature(temp int) error
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
	http.HandleFunc("/api/temp", w.postTemperature)
	http.HandleFunc("/echo", w.echo)
	return http.ListenAndServe(w.addr, nil)
}

func (w *worker) postTemperature(rw http.ResponseWriter, req *http.Request) {
	//if err := w.repo.WriteTemperature(1); err != nil {
	//	вот так будешь вызывать потом
	//}
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	rw.Header().Set("Content-Type", "application/json")
	for _, qwe := range conn {
		err := qwe.WriteMessage(1, []byte("qwe"))
		if err != nil {
			return
		}
	}
	_, err = rw.Write([]byte("ok"))
	if err != nil {
		log.Println(err)
		return
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
	conn[string(message)] = c
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		msg := string(message)
		nn := 0
		for n, m := range msg {
			if m == '|' {
				nn = n + 1
				break
			}
		}

		log.Printf("recv: %s", message[nn:])
	}
}
