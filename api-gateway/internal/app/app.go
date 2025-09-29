package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rauan06/realtime-map/api-gateway/config"
	"github.com/rauan06/realtime-map/go-commons/pkg/httpserver"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	if cfg.Metrics.Enabled {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}

	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))

	httpServer.Start()

	l.Fatal("123")
}
