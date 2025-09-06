package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rauan06/realtime-map/go-commons/pkg/grpcserver"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/config"
	"github.com/rauan06/realtime-map/producer/internal/controller/grpcrouter"
	"github.com/rauan06/realtime-map/producer/internal/repo/eventbus"
	"github.com/rauan06/realtime-map/producer/internal/usecase/producer"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	eb, err := eventbus.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer eb.Close()

	uc := producer.New(*eb)

	grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
	grpcrouter.NewRoutes(grpcServer.App, grpcrouter.RouteConfig{
		UseCase:           uc,
		Logger:            l,
		ReflectinoEnabled: cfg.GRPC.ReflectionEnabled,
	})

	grpcServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	err = grpcServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}
}