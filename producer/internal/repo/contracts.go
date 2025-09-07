package repo

import (
	"context"
	"time"

	"github.com/rauan06/realtime-map/producer/internal/domain"
)

type (
	IEventBus interface {
		ProduceEvent(context.Context, domain.OBUData) error
		// GetHistory(context.Context) ([]domain.OBUData, error)
	}
	ICache interface {
		Get(context.Context, string) ([]byte, error)
		Set(context.Context, string, []byte, time.Duration) error
		Delete(context.Context, string) error
	}
)
