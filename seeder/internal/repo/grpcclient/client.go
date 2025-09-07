package grpcclient

import (
	"log"

	producerpb "github.com/rauan06/realtime-map/go-commons/gen/proto/producer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	producerpb.ProducerServiceClient
}

func New() *Client {
	client, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	prodcuerClient := producerpb.NewProducerServiceClient(client)
	return &Client{prodcuerClient}
}
