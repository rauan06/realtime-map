package grpcclient

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	routepb "github.com/rauan06/realtime-map/go-commons/gen/proto/route"
)

type Client struct {
	routepb.RouteClient
}

func New() *Client {
	client, err := grpc.NewClient("localhost:5081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	routeClient := routepb.NewRouteClient(client)

	return &Client{routeClient}
}
