package main

import (
	"log"

	"github.com/rauan06/realtime-map/receiver/config"
	"github.com/rauan06/realtime-map/receiver/internal/app"
)

func main() {
	cfg,err := config.NewConfig()
	if err != nil {
		log.Fatal("Loading config:", err)
	}

	app.Run(cfg)
}
