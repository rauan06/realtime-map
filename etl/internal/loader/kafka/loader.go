package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rauan06/realtime-map/etl/internal/domain"
	kafkaproducer "github.com/rauan06/realtime-map/go-commons/pkg/kafka/producer"
)

var errFlushPending = errors.New("kafka flush left messages in queue")

const flushTimeoutMs = 30000

type Loader struct {
	producer *kafkaproducer.KafkaProducer
	mu       sync.Mutex
	buffer   []domain.KafkaEvent
}

func New(producer *kafkaproducer.KafkaProducer) *Loader {
	return &Loader{
		producer: producer,
		buffer:   make([]domain.KafkaEvent, 0),
	}
}

func (l *Loader) Add(event domain.KafkaEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buffer = append(l.buffer, event)
}

func (l *Loader) Flush(ctx context.Context) error {
	l.mu.Lock()
	batch := make([]domain.KafkaEvent, len(l.buffer))
	copy(batch, l.buffer)
	l.buffer = l.buffer[:0]
	l.mu.Unlock()

	if len(batch) == 0 {
		return nil
	}

	for _, event := range batch {
		if err := l.producer.ProduceEvent(ctx, event.Key, event.Data); err != nil {
			return fmt.Errorf("kafka loader flush: %w", err)
		}
	}

	remaining := l.producer.Flush(flushTimeoutMs)
	if remaining > 0 {
		return fmt.Errorf("kafka loader flush: %w (remaining=%d)", errFlushPending, remaining)
	}

	return nil
}

func (l *Loader) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return len(l.buffer)
}
