package eventbus

import (
	"context"

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

func New(cfg *config.Config, l logger.Logger) *EventBus {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	return &EventBus{p}
}

func (eb *EventBus) Run() {
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

func (eb *EventBus) ProduceEvent(context.Context, domain.OBUData) error {
	eb.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(word),
	}, nil)
	return nil
}
