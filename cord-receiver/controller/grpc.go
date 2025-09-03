package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/rauan06/realtime-map/go-commons/gen/proto/cord-receiver/locationpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// Implement the service
type locationServer struct {
	locationpb.UnimplementedLocationServiceServer
}

func (s *locationServer) SendLocation(ctx context.Context, req *locationpb.GpsData) (*locationpb.LocationResponse, error) {
	fmt.Printf("Received location: lat=%.6f, lon=%.6f, ts=%d\n", req.Latitude, req.Longitude, req.Timestamp)
	return &locationpb.LocationResponse{Status: "OK"}, nil
}

func main() {
	// Start gRPC server
	grpcServer := grpc.NewServer()
	locationpb.RegisterLocationServiceServer(grpcServer, &locationServer{})
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start HTTP Gateway
	mux := runtime.NewServeMux()
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err = locationpb.RegisterLocationServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("HTTP gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
