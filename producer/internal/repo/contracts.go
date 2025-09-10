package repo

import (
	"context"
	"time"
)

type (
	IEventBus interface {
		ProduceEvent(ctx context.Context, key string, data interface{}) error
		// GetHistory(context.Context) ([]domain.OBUData, error)
	}
	ICache interface {
		Get(context.Context, string) ([]byte, error)
		Set(context.Context, string, []byte, time.Duration) error
		Delete(context.Context, string) error
	}
)
