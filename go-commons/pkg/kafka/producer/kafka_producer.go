package kafka_producer

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

type KafkaProducer struct {
	*kafka.Producer
	TopicPartition kafka.TopicPartition
}

func New(p *kafka.Producer, topic string) (*KafkaProducer, error) {
	return &KafkaProducer{
		Producer:       p,
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	}, nil
}

func (kp *KafkaProducer) Run(l *logger.Logger) {
	// Delivery report handler for produced messages
	go func() {
		for e := range kp.Events() {
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

func (kp *KafkaProducer) ProduceEvent(ctx context.Context, key string, data interface{}) error {
	parsedData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = kp.Produce(&kafka.Message{
		Key:            []byte(key),
		TopicPartition: kp.TopicPartition,
		Value:          parsedData,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
