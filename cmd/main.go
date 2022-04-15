package main

import (
	"flag"
	"gitee.com/forkpoons/farm/internal/service"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type cnf struct {
	serverAddr string `mapstructure:"serverAddr"`
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
	if err := config.BindStruct("", &c); err != nil {
		panic(err)
	}
	//database.PostTemperature(10)
	service.Start(c.serverAddr)
}
