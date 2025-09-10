package main

import (
	"context"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/seeder/internal/repo/grpcclient"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func genCord() float64 {
	return rand.Float64() * 100
}

func runDevice(ctx context.Context, client *grpcclient.Client, deviceID uuid.UUID) {
	// Start session for this device
	_, err := client.StartSession(ctx, &route.DeviceID{
		DeviceId: deviceID[:],
	})
	if err != nil {
		log.Printf("[Device %s] failed to start session: %v", deviceID, err)
		return
	}

	// Open streaming connection
	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Printf("[Device %s] failed to open RouteChat: %v", deviceID, err)
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
					DeviceId:  deviceID[:],
					Latitude:  genCord(),
					Longitude: genCord(),
					Timestamp: timestamppb.New(time.Now()),
				}
				if err := stream.Send(msg); err != nil {
					log.Printf("[Device %s] failed to send data: %v", deviceID, err)
					return
				}
				log.Printf("[Device %s] sent OBUData: %+v", deviceID, msg)
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
		deviceID := uuid.New()
		go runDevice(ctx, client, deviceID)
	}

	// Wait for interrupt signal
	<-sigCh
	log.Println("shutting down client...")
	cancel()
}
