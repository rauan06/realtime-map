package consumer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rauan06/realtime-map/dashboard/internal/hub"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

const readTimeout = time.Second

// MultiTopic subscribes to every topic in topics and republishes each Kafka
// message to the hub tagged with the topic-derived layer name.
type MultiTopic struct {
	consumer *kafka.Consumer
	topics   []string
	hub      *hub.Hub
	l        logger.Interface
}

func New(c *kafka.Consumer, topics []string, h *hub.Hub, l logger.Interface) *MultiTopic {
	return &MultiTopic{consumer: c, topics: topics, hub: h, l: l}
}

func (m *MultiTopic) Run(ctx context.Context) error {
	if err := m.consumer.SubscribeTopics(m.topics, nil); err != nil {
		return fmt.Errorf("dashboard consumer subscribe: %w", err)
	}
	defer m.consumer.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := m.consumer.ReadMessage(readTimeout)
			if err == nil {
				layer := topicToLayer(*msg.TopicPartition.Topic)
				m.hub.Publish(hub.Message{Layer: layer, Payload: msg.Value})

				continue
			}

			var kafkaErr kafka.Error
			if errors.As(err, &kafkaErr) && !kafkaErr.IsTimeout() {
				m.l.Error("dashboard consumer error: %v", err)
			}
		}
	}
}

func topicToLayer(topic string) string {
	switch topic {
	case "etl_flights":
		return "flight"
	case "etl_ships":
		return "ship"
	case "etl_transport":
		return "transport"
	case "etl_roads":
		return "road"
	case "obu_data":
		return "obu"
	default:
		return topic
	}
}
