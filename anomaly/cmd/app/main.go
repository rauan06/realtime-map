package main

import (
	"log"

	"github.com/rauan06/realtime-map/anomaly/config"
	"github.com/rauan06/realtime-map/anomaly/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("loading config: ", err)
	}

	app.Run(cfg)
}
