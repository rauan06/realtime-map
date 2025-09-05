package producer

import (
	"github.com/rauan06/realtime-map/producer/internal/repo/eventbus"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
)

var _ (usecase.IProducerUseCase) = &UseCase{}

type UseCase struct {
	eventbus eventbus.EventBus
}

func New(eb eventbus.EventBus) *UseCase {
	return &UseCase{eb}
}

func (uc *UseCase) StartTracking()  {}
func (uc *UseCase) ProcessOBUData() {}
