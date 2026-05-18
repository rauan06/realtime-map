package main

import (
	"log"

	"github.com/rauan06/realtime-map/etl/config"
	"github.com/rauan06/realtime-map/etl/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Loading config: ", err)
	}

	app.Run(cfg)
}
