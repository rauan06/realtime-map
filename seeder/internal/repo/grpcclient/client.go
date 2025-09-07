package grpcclient

import (
	"log"

	routepb "github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	routepb.RouteClient
}

func New() *Client {
	client, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	prodcuerClient := routepb.NewRouteClient(client)
	return &Client{prodcuerClient}
}
