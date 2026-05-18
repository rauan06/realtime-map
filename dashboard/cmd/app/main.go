package main

import (
	"io/fs"
	"log"

	"github.com/rauan06/realtime-map/dashboard/config"
	"github.com/rauan06/realtime-map/dashboard/internal/app"
	"github.com/rauan06/realtime-map/dashboard/web"
)

func main() {
	templates, err := fs.Sub(web.FS, "templates")
	if err != nil {
		log.Fatalf("template embed root: %v", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("loading config: ", err)
	}

	app.Run(cfg, templates)
}
