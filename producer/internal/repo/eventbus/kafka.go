package eventbus

import (
	"github.com/rauan06/realtime-map/go-commons/pkg/kafka/broker/producer"
	"github.com/rauan06/realtime-map/producer/internal/repo"
)

var _ (*repo.IEventBus) = &EventBus{}

type EventBus struct {
	*producer.Producer
}

func New(eb *producer.Producer) *EventBus {
	return &EventBus{eb}
}

func (eb *EventBus) 