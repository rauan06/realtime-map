package location

import (
	"time"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

type Transformer struct{}

func New() *Transformer {
	return &Transformer{}
}

func (t *Transformer) Transform(records []domain.RawRecord) ([]domain.KafkaEvent, error) {
	events := make([]domain.KafkaEvent, 0, len(records))

	for _, r := range records {
		r.Fields["source_id"] = r.SourceID
		r.Fields["timestamp"] = r.Timestamp.Format(time.RFC3339)

		events = append(events, domain.KafkaEvent{
			Key:  r.SourceID,
			Data: r.Fields,
		})
	}

	return events, nil
}
