package producer

import (
	"github.com/rauan06/realtime-map/producer/internal/repo/cache"
	"github.com/rauan06/realtime-map/producer/internal/repo/eventbus"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
)

var _ (usecase.IProducerUseCase) = &UseCase{}

type UseCase struct {
	eventbus eventbus.EventBus
	cache    cache.Cache
}

func New(eb eventbus.EventBus, c cache.Cache) *UseCase {
	return &UseCase{eb, c}
}

func (uc *UseCase) StartTracking()  {}
func (uc *UseCase) ProcessOBUData() {}
