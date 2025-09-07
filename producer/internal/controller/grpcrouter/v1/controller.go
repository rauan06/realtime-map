package v1

import (
	"github.com/go-playground/validator/v10"
	routepb "github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
	"google.golang.org/grpc"
)

type V1 struct {
	routepb.UnimplementedRouteServer

	uc usecase.IProducerUseCase
	l logger.Interface
	v *validator.Validate
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

		
		// Note: this copy prevents blocking other clients while serving this one.
		// We don't need to do a deep copy, because elements in the slice are
		// insert-only and never modified.
		rn := make([]*pb.RouteNote, len(s.routeNotes[key]))
		copy(rn, s.routeNotes[key])
		s.mu.Unlock()
4
		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}
