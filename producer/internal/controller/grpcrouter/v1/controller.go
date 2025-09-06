package v1

import (
	"io"
	"log"

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
	waitc := make(chan struct{})

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			close(waitc)
			return nil
		}
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
			return err
		}
		log.Printf("Got message %d at point(%f, %f)", in.Timestamp, in.Latitude, in.Longitude)
	}
}
