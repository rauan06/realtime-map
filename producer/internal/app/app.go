package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rauan06/realtime-map/go-commons/pkg/grpcserver"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/config"
	"github.com/rauan06/realtime-map/producer/internal/controller/grpcrouter"
	"github.com/rauan06/realtime-map/producer/internal/repo/cache"
	"github.com/rauan06/realtime-map/producer/internal/repo/eventbus"
	"github.com/rauan06/realtime-map/producer/internal/usecase/producer"
	"github.com/redis/go-redis/v9"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka.BootstrapServers})
	if err != nil {
		l.Fatal(err)
	}
	eb, err := eventbus.New(kafkaProducer, cfg)
	if err != nil {
		l.Fatal(err)
	}
	defer kafkaProducer.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.RedisURI,
		Password: cfg.Redis.RedisPassword,
		DB:       0,
	})
	redis, err := cache.New(redisClient, *cfg)
	if err != nil {
		l.Fatal(err)
	}

	uc := producer.New(eb, redis)

	grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
	grpcrouter.NewRoutes(grpcServer.App, grpcrouter.RouteConfig{
		UseCase:           uc,
		Logger:            l,
		ReflectinoEnabled: cfg.GRPC.ReflectionEnabled,
	})

	grpcServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	err = grpcServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	err = redisClient.Close()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - redisClient.Shutdown: %w", err))
	}
}
