package v1

import (
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	routepb "github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
	"google.golang.org/grpc"
)

type V1 struct {
	routepb.UnimplementedRouteServer

	uc usecase.IProducerUseCase
	l  logger.Interface
	v  *validator.Validate
}

func (r *V1) RouteChat(stream grpc.BidiStreamingServer[routepb.OBUData, routepb.OBUData]) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		r.l.Info(fmt.Sprintf("recieved: %+v/n", in))
		r.uc.ProcessOBUData()
	}
}
