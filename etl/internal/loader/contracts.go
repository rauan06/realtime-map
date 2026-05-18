package loader

import (
	"context"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

// Loader sends events to a downstream sink.
type Loader interface {
	Add(event domain.KafkaEvent)
	Flush(ctx context.Context) error
	Len() int
}
