package v1

import (
	"github.com/go-playground/validator/v10"
	producerpb "github.com/rauan06/realtime-map/go-commons/gen/proto/producer"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	usecase "github.com/rauan06/realtime-map/producer/internal/usecase"
	"google.golang.org/grpc"
)

func NewTranslationRoutes(sv *grpc.Server, t usecase.IProducerUseCase, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	producerpb.RegisterProducerServiceServer(sv, r)
}
