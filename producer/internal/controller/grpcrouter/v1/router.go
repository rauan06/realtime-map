package v1

import (
	"github.com/go-playground/validator/v10"
	routepb "github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	usecase "github.com/rauan06/realtime-map/producer/internal/usecase"
	"google.golang.org/grpc"
)

func NewTranslationRoutes(sv *grpc.Server, t usecase.IProducerUseCase, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	routepb.RegisterRouteServer(sv, r)
}
