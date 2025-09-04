package repo

import (
	"context"

	"github.com/rauan06/realtime-map/receiver/internal/domain"
)

type (
	IEventBus interface {
		Store(context.Context, domain.OBUData) error
		GetHistory(context.Context) ([]domain.OBUData, error)		
	}
)