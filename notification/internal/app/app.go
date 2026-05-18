package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/notification/config"
	"github.com/rauan06/realtime-map/notification/internal/consumer"
	"github.com/rauan06/realtime-map/notification/internal/geofence"
	"github.com/rauan06/realtime-map/notification/internal/notifier"
	"github.com/rauan06/realtime-map/notification/internal/usecase"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	reg, err := geofence.LoadFromFile(cfg.Geofence.ConfigPath)
	if err != nil {
		l.Fatal(err)
	}

	l.Info("loaded %d geofences from %s", reg.Len(), cfg.Geofence.ConfigPath)

	n := notifier.New(cfg.Notifier.WebhookURL, l)
	uc := usecase.New(reg, n, *l)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.BootstrapServers,
		"group.id":          cfg.Kafka.GroupID,
		"auto.offset.reset": "latest",
	})
	if err != nil {
		l.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mt := consumer.New(c, cfg.Kafka.Topics, uc, l)
	go func() {
		if err := mt.Run(ctx); err != nil {
			l.Error("consumer: %v", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	l.Info("notification - signal: %s", s.String())
	cancel()
}
