package producer

import (
	"context"

	"github.com/rauan06/realtime-map/producer/internal/domain"
	"github.com/rauan06/realtime-map/producer/internal/repo"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
)

type UseCase struct {
	eventbus repo.IEventBus
}

func New(eb repo.IEventBus) usecase.IProducerUseCase {
	return &UseCase{eb}
}

func (uc *UseCase) StartSession(ctx context.Context, ID []byte) error {
	err := uc.eventbus.ProduceEvent(ctx, "session_start", ID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) EndSession(ctx context.Context, ID []byte) error {
	err := uc.eventbus.ProduceEvent(ctx, "session_start", ID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) ProcessOBUData(ctx context.Context, data domain.OBUData) error {
	err := uc.eventbus.ProduceEvent(ctx, "obu_data", data)
	if err != nil {
		return err
	}

	return nil
}
