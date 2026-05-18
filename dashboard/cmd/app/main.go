package main

import (
	"log"

	"github.com/rauan06/realtime-map/dashboard/config"
	"github.com/rauan06/realtime-map/dashboard/internal/app"
	"github.com/rauan06/realtime-map/dashboard/web"
)

func main() {
	app.SetTemplates(web.FS, "templates")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("loading config: ", err)
	}
	app.Run(cfg)
}
