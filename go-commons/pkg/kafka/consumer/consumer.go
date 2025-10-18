package consumer

import (
	"context"
	"errors"
	"runtime"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

const workersPerThread = 3

type KafkaConsumer struct {
	*kafka.Consumer
	TopicPartition kafka.TopicPartition
	Cancel         context.CancelFunc
	Errors         chan error

	l   logger.Logger
	uc  uc
	ctx context.Context
}

func New(c *kafka.Consumer, usecase uc, l logger.Logger, topic string) (*KafkaConsumer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &KafkaConsumer{
		Consumer:       c,
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Cancel:         cancel,
		Errors:         make(chan error),

		l:   l,
		uc:  usecase,
		ctx: ctx,
	}, nil
}

func (kc *KafkaConsumer) Run() error {
	workers := workersPerThread * runtime.GOMAXPROCS(-1)

	limit := make(chan struct{}, workers)

	err := kc.SubscribeTopics([]string{*kc.TopicPartition.Topic}, nil)
	if err != nil {
		return err
	}

	for {
		limit <- struct{}{}

		select {
		case <-kc.ctx.Done():
			return nil
		case err := <-kc.Errors:
			kc.l.Error(err)
		default:
			msg, err := kc.ReadMessage(time.Second)
			if err == nil {
				go func() {
					defer func() { <-limit }()

					kc.uc.ProcessMessage(msg)
				}()
			}

			var kafkaErr kafka.Error
			if errors.As(err, &kafkaErr) && !kafkaErr.IsTimeout() {
				kc.l.Error("consumer error: %v (%v)\n", err, msg)
			}

			<-limit
		}
	}
}
