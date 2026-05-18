package main

import (
	"log"

	"github.com/rauan06/realtime-map/notification/config"
	"github.com/rauan06/realtime-map/notification/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("loading config: ", err)
	}
	app.Run(cfg)
}
