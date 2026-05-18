package extractor

import (
	"context"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

// Extractor fetches raw records from an external source.
type Extractor interface {
	Extract(ctx context.Context) ([]domain.RawRecord, error)
}
