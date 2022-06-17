package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strconv"
)

type webWorker struct {
	addr string
}

func (w *webWorker) Start() error {
	r := gin.Default()
	api := r.Group("/api")
	api.GET("/", w.get)
	api.GET("/api/temp", w.getTemperature)
	api.GET("/api/echofront", w.echoFront)
	api.POST("/api/settemp", w.postTemperature)
	return r.Run(w.addr)
}

func (w *webWorker) get(c *gin.Context) {
	log.Println(req.Header)
	_, err := rw.Write([]byte("ok"))
	if err != nil {
		log.Println(err)
		return
	}
}

func (w *webWorker) getTemperature(c *gin.Context) {
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

func (w *webWorker) postTemperature(c *gin.Context) {
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
		for _, con := range connDevice {
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

func (w *webWorker) echoFront(c *gin.Context) {
	c, err := upgrader.Upgrade(rw, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	_, message, err := c.ReadMessage()
	connFront[string(message)] = c

	defer delete(connDevice, string(message))

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if len(message) > 0 {
			msg := string(message)
			fmt.Println("Message:" + msg)

		}
	}
}
