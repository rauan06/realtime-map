package eventbus

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/config"
	"github.com/rauan06/realtime-map/producer/internal/domain"
	"github.com/rauan06/realtime-map/producer/internal/repo"
)

var _ (repo.IEventBus) = &EventBus{}

type EventBus struct {
	*kafka.Producer
	Topic string
}

func New(p *kafka.Producer, cfg *config.Config) (repo.IEventBus, error) {
	return &EventBus{p, cfg.Kafka.Topic}, nil
}

func (eb *EventBus) Run(l *logger.Logger) {
	// Delivery report handler for produced messages
	go func() {
		for e := range eb.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					l.Error("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					l.Info("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
}

func (eb *EventBus) ProduceEvent(ctx context.Context, data domain.KafkaMessage) error {
	parsedData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	eb.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &eb.Topic, Partition: kafka.PartitionAny},
		Value:          parsedData,
	}, nil)
	return nil
}
