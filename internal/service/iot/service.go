package iot

import (
	"github.com/forkpoons/farm/internal/service/iot/greenhouse"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type iotWorker struct {
	addr             string
	db               *sqlx.DB
	greenhouseWorker service
}

type service interface {
	Handlers(gh *gin.RouterGroup)
}

func New(addr string, db *sqlx.DB) (*iotWorker, error) {
	return &iotWorker{
		addr: addr,
		db:   db,
	}, nil
}

func (w *iotWorker) Start() error {
	w.greenhouseWorker = greenhouse.New(w.db)
	r := gin.Default()
	api := r.Group("/api")
	w.greenhouseWorker.Handlers(api.Group("/greenhouse"))
	return r.Run()
}
