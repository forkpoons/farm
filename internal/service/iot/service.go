package iot

import (
	"github.com/forkpoons/farm/internal/service/iot/greenhouse"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

type iotWorker struct {
	addr        string
	db          *sqlx.DB
	greenhouses map[int]*greenhouse.Device
}

func New(addr string, db *sqlx.DB) (*iotWorker, error) {
	return &iotWorker{
		addr:        addr,
		db:          db,
		greenhouses: make(map[int]*greenhouse.Device),
	}, nil
}

func (w *iotWorker) Start() error {
	greenhouse.ConnDevice = make(map[string]*websocket.Conn)
	greenhouse.Devices = make(map[string]*greenhouse.Device)
	greenhouse.Data = w.db
	r := gin.Default()
	api := r.Group("/api")
	greenhouse.Handlers(api.Group("/greenhouse"))
	return r.Run()
}
