package usecase

import (
	"context"

	"github.com/rauan06/realtime-map/producer/internal/domain"
)

type IProducerUseCase interface {
	StartSession(context.Context, string) (string, error)
	EndSession(context.Context, string) error
	ProcessOBUData(context.Context, domain.OBUData) error
}
