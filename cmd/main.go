package main

import (
	"flag"
	"fmt"
	"github.com/forkpoons/farm/internal/database"
	"github.com/forkpoons/farm/internal/service"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type cnf struct {
	serverAddr   string
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
	c.serverAddr = config.String("serverAddr")
	c.databaseType = config.String("databaseType")
	c.databaseConn = config.String("databaseConn")
	db, err := database.New(c.databaseType, c.databaseConn)
	if err != nil {
		panic(err)
	}
	worker, err := service.New(c.serverAddr, db)
	if err != nil {
		panic(err)
	}
	fmt.Printf("error while listen http: %v", worker.Start())
}
