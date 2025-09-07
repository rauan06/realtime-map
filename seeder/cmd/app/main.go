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
	producerpb "github.com/rauan06/realtime-map/go-commons/gen/proto/producer"
	"github.com/rauan06/realtime-map/seeder/internal/domain"
	"github.com/rauan06/realtime-map/seeder/internal/repo/grpcclient"
)

func genCord() float64 {
	return rand.Float64() * 100
}

func genOBUData() domain.OBUData {
	return domain.OBUData{
		ID:   uuid.New(),
		Long: genCord(),
		Lat:  genCord(),
	}
}

func main() {
	client := grpcclient.New()
	stream, err := client.SendLocation(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Handle Ctrl+C or kill signal to end streaming gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-sigCh:
			log.Println("Stopping stream...")
			break loop
		default:
			time.Sleep(500 * time.Millisecond)
			err = stream.Send(&producerpb.OBUData{
				Latitude:  genCord(),
				Longitude: genCord(),
				Timestamp: time.Now().Unix(),
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// End stream properly
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Failed to close stream:", err)
	}
	log.Printf("Stream finished. Server response: %+v\n", resp)
}
