package main

import (
	"flag"
	"fmt"
	"github.com/forkpoons/farm/internal/database"
	"github.com/forkpoons/farm/internal/service/iot"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type cnf struct {
	iotAddr      string
	webAddr      string
	databaseType string
	databaseConn string
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.yml", "path to config file")
	flag.Parse()
	var c cnf

	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles(configFile)
	if err != nil {
		panic(err)
	}
	c.iotAddr = config.String("iotAddr")
	c.webAddr = config.String("webAddr")
	c.databaseType = config.String("databaseType")
	c.databaseConn = config.String("databaseConn")
	db, err := database.New(c.databaseType, c.databaseConn)
	if err != nil {
		panic(err)
	}
	worker, err := iot.New(c.iotAddr, db.Db)
	if err != nil {
		panic(err)
	}
	fmt.Printf("error while listen http: %v", worker.Start())
}
