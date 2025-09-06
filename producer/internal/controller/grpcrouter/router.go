package grpcrouter

import (
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RouteConfig struct {
	uc         usecase.IProducerUseCase
	l          logger.Interface
	enableRefl bool
}

func NewRoutes(sv *grpc.Server, cfg RouteConfig) {
	

	if cfg.enableRefl {
		reflection.Register(sv)
	}
}
