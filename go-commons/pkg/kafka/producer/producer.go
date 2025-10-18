package producer

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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

func (kp *KafkaProducer) ProduceEvent(_ context.Context, key string, data interface{}) error {
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
