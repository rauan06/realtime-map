// Package publisher sends Alerts to the anomalies Kafka topic. It is a
// thin adapter on top of go-commons' kafka producer that fire-and-forgets
// individual messages and Flushes on Close.
package publisher

import (
	"context"
	"errors"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rauan06/realtime-map/anomaly/internal/detector"
	kafkaproducer "github.com/rauan06/realtime-map/go-commons/pkg/kafka/producer"
)

const flushTimeoutMs = 5000

var errFlushPending = errors.New("anomaly publisher flush left messages in queue")

type Publisher struct {
	inner *kafkaproducer.KafkaProducer
}

func New(p *kafka.Producer, topic string) (*Publisher, error) {
	inner, err := kafkaproducer.New(p, topic)
	if err != nil {
		return nil, fmt.Errorf("publisher init: %w", err)
	}

	return &Publisher{inner: inner}, nil
}

func (p *Publisher) Publish(ctx context.Context, a detector.Alert) error {
	if err := p.inner.ProduceEvent(ctx, a.SourceID, a); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

func (p *Publisher) Close() error {
	if remaining := p.inner.Flush(flushTimeoutMs); remaining > 0 {
		return fmt.Errorf("%w (remaining=%d)", errFlushPending, remaining)
	}

	p.inner.Close()

	return nil
}
