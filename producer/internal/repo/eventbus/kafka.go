package eventbus

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/producer/config"
	"github.com/rauan06/realtime-map/producer/internal/repo"
)

type EventBus struct {
	*kafka.Producer
	TopicPartition kafka.TopicPartition
}

func New(p *kafka.Producer, cfg *config.Config) (repo.IEventBus, error) {
	return &EventBus{
		Producer:  p,
		TopicPartition: kafka.TopicPartition{Topic: &cfg.Kafka.Topic, Partition: kafka.PartitionAny},
	}, nil
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

func (eb *EventBus) ProduceEvent(ctx context.Context, key string, data interface{}) error {
	parsedData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = eb.Produce(&kafka.Message{
		Key:            []byte(key),
		TopicPartition: eb.TopicPartition,
		Value:          parsedData,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
