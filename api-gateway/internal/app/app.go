package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rauan06/realtime-map/api-gateway/config"
	"github.com/rauan06/realtime-map/api-gateway/internal/controller"
	routepb "github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/go-commons/pkg/httpserver"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

const (
	readTimeout = 5 * time.Second
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	mux := http.NewServeMux()

	_ = httpserver.New()

	// gRPC route client to producer
	grpcConn, err := grpc.NewClient("localhost:"+cfg.GRPC.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Fatal(err)
	}

	routeClient := routepb.NewRouteClient(grpcConn)

	controller.Register(mux, routeClient)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTP.Port,
		Handler:           mux,
		ReadHeaderTimeout: readTimeout,
	}

	// Metrics server
	if cfg.Metrics.Enabled {
		go func() {
			mMux := http.NewServeMux()
			mMux.Handle("/metrics", promhttp.Handler())

			server := &http.Server{
				Addr:         ":2112",
				Handler:      mMux,
				ReadTimeout:  readTimeout,
				WriteTimeout: readTimeout,
				IdleTimeout:  readTimeout,
			}

			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				l.Error(err)
			}
		}()
	}

	go func() {
		if errors.Is(err, srv.ListenAndServe()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Fatal(err)
		}
	}()

	select {}
}
