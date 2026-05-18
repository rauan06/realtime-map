package consumer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

const readTimeout = time.Second

// Handler is invoked for every message from one of the subscribed topics.
type Handler interface {
	Handle(topic string, payload []byte)
}

type MultiTopic struct {
	consumer *kafka.Consumer
	topics   []string
	handler  Handler
	l        logger.Interface
}

func New(c *kafka.Consumer, topics []string, h Handler, l logger.Interface) *MultiTopic {
	return &MultiTopic{consumer: c, topics: topics, handler: h, l: l}
}

func (m *MultiTopic) Run(ctx context.Context) error {
	if err := m.consumer.SubscribeTopics(m.topics, nil); err != nil {
		return fmt.Errorf("anomaly consumer subscribe: %w", err)
	}
	defer m.consumer.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := m.consumer.ReadMessage(readTimeout)
			if err == nil {
				m.handler.Handle(*msg.TopicPartition.Topic, msg.Value)

				continue
			}

			var kafkaErr kafka.Error
			if errors.As(err, &kafkaErr) && !kafkaErr.IsTimeout() {
				m.l.Error("anomaly consumer error: %v", err)
			}
		}
	}
}
