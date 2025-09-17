package main

import (
	"context"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/seeder/internal/repo/grpcclient"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func genCord() float64 {
	return rand.Float64() * 100
}

func runDevice(ctx context.Context, client *grpcclient.Client) {
	// Start session for this device
	session, err := client.StartSession(ctx, nil)
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return
	}

	// Open streaming connection
	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Printf("[Session %s] failed to open RouteChat: %v", session.SessionId, err)
		return
	}

	// Sender goroutine
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				msg := &route.OBUData{
					SessionId: session.SessionId,
					Latitude:  genCord(),
					Longitude: genCord(),
					Timestamp: timestamppb.New(time.Now()),
				}
				if err := stream.Send(msg); err != nil {
					log.Printf("[Session %s] failed to send data: %v", session.SessionId, err)
					return
				}
				log.Printf("[Session %s] sent OBUData: %+v", session.SessionId, msg)
			}
		}
	}()
}

func main() {
	// Create client
	client := grpcclient.New()

	// Use cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for i := 0; i < 10; i++ {
		go runDevice(ctx, client)
	}

	// Wait for interrupt signal
	<-sigCh
	log.Println("shutting down client...")
	cancel()
}
