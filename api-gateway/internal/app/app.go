package app

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rauan06/realtime-map/api-gateway/config"
	"github.com/rauan06/realtime-map/api-gateway/internal/controller"
	routepb "github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"github.com/rauan06/realtime-map/go-commons/pkg/httpserver"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Metrics server
	if cfg.Metrics.Enabled {
		go func() {
			mMux := http.NewServeMux()
			mMux.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(":2112", mMux); err != nil && err != http.ErrServerClosed {
				l.Error(err)
			}
		}()
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal(err)
		}
	}()

	select {}
}
