package consumer

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

type KafkaConsumer struct {
	*kafka.Consumer
	TopicPartition kafka.TopicPartition
	Cancel         context.CancelFunc

	workers int
	l       logger.Logger
	uc      func(*kafka.Message) error
	ctx     context.Context
}

func New(c *kafka.Consumer, usecase func(kafka.Message), l logger.Logger, topic string) (*KafkaConsumer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaConsumer{
		Consumer:       c,
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Cancel:         cancel,

		l:   l,
		ctx: ctx,
	}, nil
}

func (kc *KafkaConsumer) Run(l *logger.Logger) error {
	limit := make(chan struct{}, kc.workers)

	err := kc.SubscribeTopics([]string{*kc.TopicPartition.Topic}, nil)
	if err != nil {
		return err
	}

	for {
		limit <-struct{}{}
		
		select {
		case <-kc.ctx.Done():
			return nil
		default:
			msg, err := kc.ReadMessage(time.Second)
			if err == nil {
				kc.uc(msg)
			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				kc.l.Error("consumer error: %v (%v)\n", err, msg)
			}
		}

		<- limit
	}
}
