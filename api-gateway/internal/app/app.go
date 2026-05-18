package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rauan06/realtime-map/api-gateway/config"
	"github.com/rauan06/realtime-map/api-gateway/internal/auth"
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

	// Optional auth: when enabled, wrap the whole mux so /ws + any /api/v1/*
	// endpoint require a valid bearer token. /auth/login is registered before
	// the wrap so it stays publicly reachable.
	var handler http.Handler = mux
	if cfg.Auth.Enabled {
		if cfg.Auth.JWTSecret == "" {
			l.Fatal("AUTH_ENABLED=true but AUTH_JWT_SECRET is empty")
		}
		issuer := auth.NewIssue(cfg.Auth.JWTSecret, cfg.Auth.TokenTTL)

		publicMux := http.NewServeMux()
		controller.RegisterAuth(publicMux, controller.LoginConfig{
			Issuer:       issuer,
			SharedSecret: cfg.Auth.SharedSecret,
		})

		protected := auth.Middleware([]byte(cfg.Auth.JWTSecret))(mux)

		root := http.NewServeMux()
		root.Handle("/auth/", publicMux)
		root.Handle("/", protected)
		handler = root
	}

	srv := &http.Server{
		Addr:              ":" + cfg.HTTP.Port,
		Handler:           handler,
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
