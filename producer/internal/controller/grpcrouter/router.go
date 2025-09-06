package grpcrouter

import (
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	v1 "github.com/rauan06/realtime-map/producer/internal/controller/grpcrouter/v1"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RouteConfig struct {
	UseCase          usecase.IProducerUseCase
	Logger           logger.Interface
	ReflectinoEnabled bool
}

func NewRoutes(sv *grpc.Server, cfg RouteConfig) {
	v1.NewTranslationRoutes(sv, cfg.UseCase, cfg.Logger)

	if cfg.ReflectinoEnabled{
		reflection.Register(sv)
	}
}
