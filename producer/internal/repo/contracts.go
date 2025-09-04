package repo

import (
	"context"

	"github.com/rauan06/realtime-map/producer/internal/domain"
)

type (
	IEventBus interface {
		ProduceEvent(context.Context, domain.OBUData) error
		// GetHistory(context.Context) ([]domain.OBUData, error)		
	}
)