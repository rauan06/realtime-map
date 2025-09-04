package main

import (
	"context"
	"fmt"
	"log"
	"net"

	receivepb "github.com/rauan06/realtime-map/go-commons/gen/proto/receiver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Implement the service
type receiverServer struct {
	receivepb.UnimplementedReceiverServiceServer
}

func (s *locationServer) SendLocation(ctx context.Context, req *receivepb.OBUData) (*receivepb.ReceiverResponse, error) {
	fmt.Printf("Received location: lat=%.6f, lon=%.6f, ts=%d\n", req.Latitude, req.Longitude, req.Timestamp)
	return &receivepb.ReceiverResponse{Status: "OK"}, nil
}

func main() {
	// Start gRPC server
	grpcServer := grpc.NewServer()
	receivepb.RegisterReceiverServiceHandler(grpcServer, &receiverServer{})
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
