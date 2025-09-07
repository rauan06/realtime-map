package main

import (
	"context"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/seeder/internal/domain"
	"github.com/rauan06/realtime-map/seeder/internal/repo/grpcclient"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	stream, err := client.RouteChat(context.Background())
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
			err = stream.Send(&route.OBUData{
				DeviceId:  []byte("33f04b11-a6ac-4a43-bae3-3cdbd1d2dcd8"),
				Latitude:  genCord(),
				Longitude: genCord(),
				Timestamp: timestamppb.New(time.Now()),
			})
			if err != nil {
				if err != io.EOF {
					log.Fatal(err)
				}
			}
		}
	}
}
