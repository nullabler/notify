package main

import (
	"flag"
	"log"
	"notify/pkg/application"
	"notify/pkg/service"
)

var PathToConf = flag.String("path", "./config.yaml", "Path to config file")

func main() {
	flag.Parse()
	app := application.NewApp(*PathToConf)
	svc, err := service.NewConsumerSvc(app)
	if err != nil {
		log.Panicf("Error from init consumer service: %v", err)
	}
	defer svc.Close()

	svc.Invoke()
}
