package repo

import (
	"context"
)

type (
	IEventBus interface {
		ProduceEvent(ctx context.Context, key string, data any) error
		// GetHistory(context.Context) ([]domain.OBUData, error)
	}
)
