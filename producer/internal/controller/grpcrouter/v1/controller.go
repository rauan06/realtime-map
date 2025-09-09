package v1

import (
	"context"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	routepb "github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/internal/domain"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type V1 struct {
	routepb.UnimplementedRouteServer

	uc usecase.IProducerUseCase
	l  logger.Interface
	v  *validator.Validate
}

func (r *V1) StartSession(ctx context.Context, in *routepb.DeviceID) (*routepb.InitResponse, error) {
	session, err := r.uc.StartSession(ctx, in.DeviceId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &routepb.InitResponse{SessionId: session}, nil
}

func (r *V1) EndSession(context.Context, *routepb.DeviceID) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EndSession not implemented")
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

		if err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}

		err = r.uc.ProcessOBUData(context.Background(),
			domain.OBUData{
				ID:        in.DeviceId,
				Long:      in.Longitude,
				Lat:       in.Latitude,
				Timestamp: in.Timestamp.AsTime(),
			})
		if err != nil {
			r.l.Error(err)
			return status.Error(codes.Unauthenticated, err.Error())
		}
	}
}
