package transformer

import "github.com/rauan06/realtime-map/etl/internal/domain"

// Transformer converts raw records into Kafka-ready events.
type Transformer interface {
	Transform(records []domain.RawRecord) ([]domain.KafkaEvent, error)
}
