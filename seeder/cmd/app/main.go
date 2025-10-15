package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/seeder/internal/repo/grpcclient"
)

const (
	tickerRange     = 100 * time.Millisecond
	floatMultiplier = 100
)

func genCord() float64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}

	// Convert bytes to uint64, then to float64 in range [0, 1)
	u := binary.BigEndian.Uint64(b[:])
	f := float64(u) / float64(math.MaxUint64)

	return f * floatMultiplier
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
		ticker := time.NewTicker(tickerRange)
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
