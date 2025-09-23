package consumer

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type (
	uc interface{ // usecase
		ProcessMessage(*kafka.Message)
	}
)
