package usecase

import (
	"context"

	"github.com/rauan06/realtime-map/producer/internal/domain"
)

type (
	IProducerUseCase interface {
		StartSession(context.Context, []byte) error
		EndSession(context.Context, []byte) error
		ProcessOBUData(context.Context, domain.OBUData) error
	}
)
