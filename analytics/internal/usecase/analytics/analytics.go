package analytics

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rauan06/realtime-map/analytics/internal/usecase"
)

type AnalyticsUseCase struct {
}

func New() usecase.IAnalyticsUseCase {
	return &AnalyticsUseCase{}
}

func (uc *AnalyticsUseCase) ProcessMessage(msg *kafka.Message) {
	fmt.Printf("Recieved msg: %+v, key: %s, value: %s\n", msg, msg.Key, msg.Value)
}
