package app

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rauan06/realtime-map/dashboard/config"
	"github.com/rauan06/realtime-map/dashboard/internal/consumer"
	"github.com/rauan06/realtime-map/dashboard/internal/controller"
	"github.com/rauan06/realtime-map/dashboard/internal/hub"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

const (
	readTimeout     = 5 * time.Second
	shutdownTimeout = 5 * time.Second
	hubSubBuffer    = 128
)

// Run starts the dashboard service: kafka consumers fanning into a hub,
// and an HTTP server serving the embedded UI + /ws. The templates filesystem
// is passed in from main so the embed.FS lives next to the binary entry
// point and we don't need a package-level global.
func Run(cfg *config.Config, templates fs.FS) {
	l := logger.New(cfg.Log.Level)

	h := hub.New(hubSubBuffer)

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

	mt := consumer.New(c, cfg.Kafka.Topics, h, l)

	go func() {
		if err := mt.Run(ctx); err != nil {
			l.Error("consumer exited: %v", err)
		}
	}()

	mux := http.NewServeMux()
	controller.Register(mux, h, templates, l)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTP.Port,
		Handler:           mux,
		ReadHeaderTimeout: readTimeout,
	}

	go func() {
		l.Info("dashboard listening on :%s (topics=%v)", cfg.HTTP.Port, cfg.Kafka.Topics)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("http server: %v", err)
			cancel()
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("dashboard - signal: %s", s.String())
	case <-ctx.Done():
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		l.Error("dashboard shutdown: %v", err)
	}

	cancel()
}
