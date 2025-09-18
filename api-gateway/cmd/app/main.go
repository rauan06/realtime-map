package main

import (
	"log"

	"github.com/rauan06/realtime-map/api-gateway/config"
	"github.com/rauan06/realtime-map/api-gateway/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
