package v1

import (
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"

	"github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
)

func NewTranslationRoutes(sv *grpc.Server, uc usecase.IProducerUseCase, l logger.Interface) {
	r := &V1{uc: uc, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	route.RegisterRouteServer(sv, r)
}
