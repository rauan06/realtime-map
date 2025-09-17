package repo

import (
	"context"
	"time"

	"github.com/rauan06/realtime-map/analytics/internal/domain"
)

type (
	IDatabase[T domain.Entity] interface {
		Create(entity *T) error
		GetByID(id string, preload ...string) (*T, error)
		Update(entity *T) error
		Delete(id string) error
	}

	ICache interface {
		Get(context.Context, string) ([]byte, error)
		Set(context.Context, string, []byte, time.Duration) error
		Delete(context.Context, string) error
	}
)
