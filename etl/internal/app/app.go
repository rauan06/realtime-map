package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rauan06/realtime-map/etl/config"
	"github.com/rauan06/realtime-map/etl/internal/domain"
	"github.com/rauan06/realtime-map/etl/internal/extractor"
	"github.com/rauan06/realtime-map/etl/internal/extractor/flight"
	"github.com/rauan06/realtime-map/etl/internal/extractor/road"
	"github.com/rauan06/realtime-map/etl/internal/extractor/ship"
	"github.com/rauan06/realtime-map/etl/internal/extractor/transport"
	"github.com/rauan06/realtime-map/etl/internal/loader"
	chloader "github.com/rauan06/realtime-map/etl/internal/loader/clickhouse"
	kfkloader "github.com/rauan06/realtime-map/etl/internal/loader/kafka"
	multiloader "github.com/rauan06/realtime-map/etl/internal/loader/multi"
	"github.com/rauan06/realtime-map/etl/internal/pipeline"
	"github.com/rauan06/realtime-map/etl/internal/transformer/location"
	kafkaproducer "github.com/rauan06/realtime-map/go-commons/pkg/kafka/producer"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

var errNoEnabledPipelines = errors.New("no enabled source pipelines — set at least one *_ENABLED=true")

const (
	sourceFlight    = "flight"
	sourceShip      = "ship"
	sourceTransport = "transport"
	sourceRoad      = "road"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipelines, cleanup, err := buildPipelines(ctx, cfg, l)
	if err != nil {
		l.Fatal(fmt.Errorf("build pipelines: %w", err))
	}
	defer cleanup()

	if len(pipelines) == 0 {
		l.Fatal(errNoEnabledPipelines)
	}

	var wg sync.WaitGroup

	for _, p := range pipelines {
		wg.Add(1)

		go func(p *pipeline.Pipeline) {
			defer wg.Done()

			if err := p.Run(ctx); err != nil {
				l.Error("pipeline exited with error: %s", err)
			}
		}(p)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	l.Info("app - Run - signal: %s", s.String())
	cancel()
	wg.Wait()
}

func buildPipelines(ctx context.Context, cfg *config.Config, l *logger.Logger) ([]*pipeline.Pipeline, func(), error) {
	var (
		pipelines []*pipeline.Pipeline
		closers   []func()
	)

	cleanup := func() {
		for i := len(closers) - 1; i >= 0; i-- {
			closers[i]()
		}
	}

	specs := []struct {
		name   string
		source string
		cfg    config.SourceConfig
		make   func() extractor.Extractor
	}{
		{
			name:   sourceFlight,
			source: sourceFlight,
			cfg:    cfg.Sources.Flight,
			make:   func() extractor.Extractor { return flight.New(cfg.Sources.HTTPTimeout) },
		},
		{
			name:   sourceShip,
			source: sourceShip,
			cfg:    cfg.Sources.Ship,
			make:   func() extractor.Extractor { return ship.New(cfg.Sources.HTTPTimeout) },
		},
		{
			name:   sourceTransport,
			source: sourceTransport,
			cfg:    cfg.Sources.Transport,
			make:   func() extractor.Extractor { return transport.New(cfg.Sources.HTTPTimeout) },
		},
		{
			name:   sourceRoad,
			source: sourceRoad,
			cfg:    cfg.Sources.Road,
			make:   func() extractor.Extractor { return road.New(cfg.Sources.HTTPTimeout, cfg.Sources.Road.Endpoint) },
		},
	}

	for _, s := range specs {
		if !s.cfg.Enabled {
			l.Info("skipping disabled source: %s", s.name)

			continue
		}

		ld, err := buildLoader(ctx, cfg, s.source, s.cfg.Topic, &closers)
		if err != nil {
			cleanup()

			return nil, nil, fmt.Errorf("build %s loader: %w", s.name, err)
		}

		p := pipeline.New(pipeline.Config{
			Name:          s.name,
			Extractor:     s.make(),
			Transformer:   location.New(),
			Loader:        ld,
			Logger:        l,
			FetchInterval: s.cfg.FetchInterval,
			FlushInterval: cfg.Sources.FlushInterval,
			BatchSize:     cfg.Sources.BatchSize,
		})
		pipelines = append(pipelines, p)
	}

	return pipelines, cleanup, nil
}

func buildLoader(ctx context.Context, cfg *config.Config, source, topic string, closers *[]func()) (loader.Loader, error) {
	var children []loader.Loader

	if cfg.Kafka.Enabled {
		kp, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": cfg.Kafka.BootstrapServers,
		})
		if err != nil {
			return nil, fmt.Errorf("kafka producer: %w", err)
		}

		*closers = append(*closers, func() { kp.Close() })

		prod, err := kafkaproducer.New(kp, topic)
		if err != nil {
			return nil, fmt.Errorf("kafka producer wrapper: %w", err)
		}

		children = append(children, kfkloader.New(prod))
	}

	if cfg.ClickHouse.Enabled {
		ch, err := chloader.New(ctx, chloader.Options{
			Addr:     cfg.ClickHouse.Addr,
			Database: cfg.ClickHouse.Database,
			Username: cfg.ClickHouse.Username,
			Password: cfg.ClickHouse.Password,
			Source:   source,
		})
		if err != nil {
			return nil, fmt.Errorf("clickhouse loader: %w", err)
		}

		*closers = append(*closers, func() { _ = ch.Close() })
		children = append(children, ch)
	}

	if len(children) == 0 {
		return nil, fmt.Errorf("%w: %s", domain.ErrNoLoaders, source)
	}

	if len(children) == 1 {
		return children[0], nil
	}

	return multiloader.New(children...), nil
}
