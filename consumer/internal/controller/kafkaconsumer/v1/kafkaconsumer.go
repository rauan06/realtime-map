package kafkaconsumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/rauan06/realtime-map/consumer/internal/usecase"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

const (
	PoolSize = 30
)

type kafkaConsumer struct {
	l logger.Logger
	cfg *config.Config
	v   *validator.Validate
	uc  *usecase.IConsumerUseCase
	// metrics *metrics.ReaderServiceMetrics
}

func New() {}

func (k *kafkaConsumer) ProcessMessages(ctx context.Context, r *kafka.Consumer, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.ReadMessage(time.Second)
		if err != nil {
			k.l.Warn(fmt.Sprintf("workerID: %v, err: %v", workerID, err))
			continue
		}

		k.logProcessMessage(m, workerID)

		switch m.Topic {
		case k.cfg.KafkaTopics.ProductCreated.TopicName:
			k.processProductCreated(ctx, r, m)
		case k.cfg.KafkaTopics.ProductUpdated.TopicName:
			k.processProductUpdated(ctx, r, m)
		case k.cfg.KafkaTopics.ProductDeleted.TopicName:
			k.processProductDeleted(ctx, r, m)
		}
	}
}
