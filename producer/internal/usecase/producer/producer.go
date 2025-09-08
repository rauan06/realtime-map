package producer

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rauan06/realtime-map/producer/internal/domain"
	"github.com/rauan06/realtime-map/producer/internal/repo"
	"github.com/rauan06/realtime-map/producer/internal/usecase"
)

var _ (usecase.IProducerUseCase) = &UseCase{}

type UseCase struct {
	eventbus repo.IEventBus
	cache    repo.ICache
}

func New(eb repo.IEventBus, c repo.ICache) usecase.IProducerUseCase {
	return &UseCase{eb, c}
}

func (uc *UseCase) StartSession(ctx context.Context, ID []byte) ([]byte, error) {
	session, _ := uc.cache.Get(ctx, string(ID))
	if len(session) != 0 {
		return session, nil
	}

	session = []byte(uuid.NewString())
	err := uc.cache.Set(ctx, string(ID), session, 2*time.Minute)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *UseCase) EndSession(ctx context.Context, ID string) error {
	_, err := uc.cache.Get(ctx, ID)
	if err != nil {
		return errors.New("no active sessions")
	}

	err = uc.cache.Delete(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) ProcessOBUData(ctx context.Context, data domain.OBUData) error {
	session, err := uc.cache.Get(ctx, string(data.ID))
	if err != nil {
		return errors.New("no session found")
	}

	err = uc.eventbus.ProduceEvent(ctx, domain.KafkaMessage{
		Session: string(session),
		Data:    data,
	})
	if err != nil {
		return err
	}

	return nil
}
