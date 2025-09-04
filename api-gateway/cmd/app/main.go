package main

import (
	"context"
	"log"
	"net/http"

	receivepb "github.com/rauan06/realtime-map/go-commons/gen/proto/receiver"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mux := runtime.NewServeMux()
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := receivepb.RegisterReceiverServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("HTTP gateway listening on :8080")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
