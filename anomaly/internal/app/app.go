package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rauan06/realtime-map/anomaly/config"
	"github.com/rauan06/realtime-map/anomaly/internal/consumer"
	"github.com/rauan06/realtime-map/anomaly/internal/detector"
	"github.com/rauan06/realtime-map/anomaly/internal/iforest"
	"github.com/rauan06/realtime-map/anomaly/internal/publisher"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

// router fans incoming Kafka messages to the right per-layer detector and
// publishes any returned Alert to the output topic.
type router struct {
	flightTopic string
	shipTopic   string
	flight      *detector.Detector
	ship        *detector.Detector
	pub         *publisher.Publisher
	l           logger.Interface
}

func (r *router) Handle(topic string, payload []byte) {
	var (
		alert *detector.Alert
		err   error
	)

	switch topic {
	case r.flightTopic:
		alert, err = r.flight.Observe(payload)
	case r.shipTopic:
		alert, err = r.ship.Observe(payload)
	default:
		return
	}

	if err != nil {
		r.l.Error("anomaly observe (%s): %v", topic, err)

		return
	}

	if alert == nil {
		return
	}

	if err := r.pub.Publish(context.Background(), *alert); err != nil {
		r.l.Error("anomaly publish: %v", err)

		return
	}

	r.l.Info("ALERT layer=%s id=%s score=%.3f reasons=%v",
		alert.Layer, alert.SourceID, alert.Score, alert.Reasons)
}

// Run wires the anomaly service together.
//
//nolint:funlen // straight-line composition root: detectors, producer, consumer
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	forestOpts := iforest.Options{
		NumTrees:   cfg.Detector.NumTrees,
		SampleSize: cfg.Detector.SampleSize,
	}

	flight := detector.New(detector.Options{
		Layer:     "flight",
		Extract:   detector.FlightFeatures,
		Warmup:    cfg.Detector.Warmup,
		BufferCap: cfg.Detector.BufferCap,
		Threshold: cfg.Detector.Threshold,
		Cooldown:  cfg.Detector.Cooldown,
		Forest:    forestOpts,
	})

	ship := detector.New(detector.Options{
		Layer:     "ship",
		Extract:   detector.ShipFeatures,
		Warmup:    cfg.Detector.Warmup,
		BufferCap: cfg.Detector.BufferCap,
		Threshold: cfg.Detector.Threshold,
		Cooldown:  cfg.Detector.Cooldown,
		Forest:    forestOpts,
	})

	const bootstrapServersKey = "bootstrap.servers"

	prod, err := kafka.NewProducer(&kafka.ConfigMap{
		bootstrapServersKey: cfg.Kafka.BootstrapServers,
	})
	if err != nil {
		l.Fatal(err)
	}

	pub, err := publisher.New(prod, cfg.Kafka.OutputTopic)
	if err != nil {
		l.Fatal(err)
	}

	const metadataRefreshMs = 10000

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		bootstrapServersKey:                  cfg.Kafka.BootstrapServers,
		"group.id":                           cfg.Kafka.GroupID,
		"auto.offset.reset":                  "latest",
		"topic.metadata.refresh.interval.ms": metadataRefreshMs,
	})
	if err != nil {
		l.Fatal(err)
	}

	rt := &router{
		flightTopic: cfg.Kafka.FlightTopic,
		shipTopic:   cfg.Kafka.ShipTopic,
		flight:      flight,
		ship:        ship,
		pub:         pub,
		l:           l,
	}

	topics := []string{cfg.Kafka.FlightTopic, cfg.Kafka.ShipTopic}
	mt := consumer.New(c, topics, rt, l)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := mt.Run(ctx); err != nil {
			l.Error("anomaly consumer: %v", err)
		}
	}()

	l.Info("anomaly service: subscribed to %v → publishing to %s", topics, cfg.Kafka.OutputTopic)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	l.Info("anomaly - signal: %s", s.String())

	cancel()

	if err := pub.Close(); err != nil {
		l.Error("anomaly publisher close: %v", err)
	}
}
