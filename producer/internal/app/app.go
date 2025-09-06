package app

import (
	"log"

	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/config"
	"github.com/rauan06/realtime-map/go-commons/pkg/grpcserver"
	"github.com/rauan06/realtime-map/producer/internal/repo/eventbus"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	eb, err := eventbus.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer eb.Close()

	grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
	
}