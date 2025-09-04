package app

import (
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/receiver/config"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	l.Debug("123")
}
