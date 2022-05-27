package iot

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type repository interface {
	WriteTemperature(name string, temp float64) error
	WriteAction(name string, temp bool) error
	ReadTemperature(name string) (float64, error)
}

type worker struct {
	addr string
	repo repository
	cons map[string]*websocket.Conn
}

func New(addr string, repo repository) (*worker, error) {
	return &worker{
		addr: addr,
		repo: repo,
		cons: make(map[string]*websocket.Conn),
	}, nil
}

func (w *worker) Start() error {
	//connDevice = make(map[string]*websocket.Conn)
	//connFront = make(map[string]*websocket.Conn)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	if err := r.Run(); err != nil {
		return err
	}

	http.HandleFunc("/", w.get)
	http.HandleFunc("/api/temp", w.getTemperature)
	http.HandleFunc("/api/settemp", w.postTemperature)
	http.HandleFunc("/api/echo", w.echo)
	http.HandleFunc("/api/echofront", w.echoFront)

	return http.ListenAndServe(w.addr, nil)
}
