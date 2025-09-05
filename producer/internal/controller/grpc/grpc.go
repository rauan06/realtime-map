package grpcserver

import (
	producerpb "github.com/rauan06/realtime-map/go-commons/gen/proto/producer"
	"github.com/rauan06/realtime-map/producer/internal/usecase/producer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerAPI struct {
	producerpb.UnimplementedProducerServiceServer
	producer producer.UseCase
}

func New(uc producer.UseCase) *ServerAPI {
	return &ServerAPI{producer: uc}
}

func (s *ServerAPI) SendLocation(stream grpc.ClientStreamingServer[producerpb.OBUData, producerpb.ProducerResponse]) error {
	stream.Recv()

	return status.Errorf(codes.Unimplemented, "method SendLocation not implemented")
}
