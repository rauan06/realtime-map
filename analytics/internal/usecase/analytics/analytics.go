package analytics

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rauan06/realtime-map/analytics/internal/usecase"
)

type AnalyticsUseCase struct {
}

func New() *usecase.IAnalyticsUseCase {
	return nil
}

func (uc *AnalyticsUseCase) ProcessMessage(*kafka.Message) {
	
}
