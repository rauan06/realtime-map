package usecase

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type (
	IAnalyticsUseCase interface {
		ProcessMessage(*kafka.Message)
	}
)
