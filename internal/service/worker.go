package service

import (
	"fmt"
	"gitee.com/forkpoons/farm/internal/database"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

var conn map[string]*websocket.Conn

func Start() {
	conn = make(map[string]*websocket.Conn)
	log.Println("222")
	http.HandleFunc("/api/temp", postTemperature)
	http.HandleFunc("/echo", echo)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func postTemperature(w http.ResponseWriter, req *http.Request) {
	database.PostTemperature(1)
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	w.Header().Set("Content-Type", "application/json")
	for _, qwe := range conn {
		err := qwe.WriteMessage(1, []byte("qwe"))
		if err != nil {
			return
		}
	}
	_, err = w.Write([]byte("ok"))
	if err != nil {
		log.Println(err)
		return
	}
}

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

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
