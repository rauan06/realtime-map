package v1

import (
	"github.com/go-playground/validator/v10"
	producerpb "github.com/rauan06/realtime-map/go-commons/gen/proto/producer"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
	"google.golang.org/grpc"
)

type V1 struct {
	producerpb.ProducerServiceServer

	t usecase.IProducerUseCase
	l logger.Interface
	v *validator.Validate
}

func (r *V1) SendLocation(
	stream grpc.ClientStreamingServer[producerpb.OBUData, producerpb.ProducerResponse]) error {
		return nil
}
