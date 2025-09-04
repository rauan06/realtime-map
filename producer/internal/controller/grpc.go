package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	producerpb "github.com/rauan06/realtime-map/go-commons/gen/proto/producer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Implement the service
type ProducerServer struct {
	producerpb.UnimplementedProducerServiceServer
}

func New() {

}

func (s *ProducerServer) SendLocation(ctx context.Context, req *producerpb.OBUData) (*producerpb.ProducerResponse, error) {
	fmt.Printf("Received location: lat=%.6f, lon=%.6f, ts=%d\n", req.Latitude, req.Longitude, req.Timestamp)
	return &producerpb.ProducerResponse{Status: "OK"}, nil
}

func RunGRPCServer() {
	grpcServer := grpc.NewServer()
	producerpb.RegisterProducerServiceServer(grpcServer, &ProducerServer{})
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
