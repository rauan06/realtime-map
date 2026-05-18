package consumer

import (
	"context"
	"errors"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

// Processor matches the shape go-commons consumer.uc expects, but kept local
// so the notification service can drive multiple topics from one consumer.
type Processor interface {
	ProcessMessage(*kafka.Message)
}

type MultiTopic struct {
	consumer *kafka.Consumer
	topics   []string
	proc     Processor
	l        logger.Interface
}

func New(c *kafka.Consumer, topics []string, p Processor, l logger.Interface) *MultiTopic {
	return &MultiTopic{consumer: c, topics: topics, proc: p, l: l}
}

func (m *MultiTopic) Run(ctx context.Context) error {
	if err := m.consumer.SubscribeTopics(m.topics, nil); err != nil {
		return err
	}
	defer m.consumer.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := m.consumer.ReadMessage(time.Second)
			if err == nil {
				m.proc.ProcessMessage(msg)
				continue
			}
			var kafkaErr kafka.Error
			if errors.As(err, &kafkaErr) && !kafkaErr.IsTimeout() {
				m.l.Error("notification consumer error: %v", err)
			}
		}
	}
}
