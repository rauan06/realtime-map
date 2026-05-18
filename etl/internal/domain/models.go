package domain

import "time"

// RawRecord represents unprocessed data from an external API source.
type RawRecord struct {
	SourceID  string
	Timestamp time.Time
	Fields    map[string]interface{}
}

// KafkaEvent is the Kafka-ready payload produced by the transformer.
type KafkaEvent struct {
	Key  string
	Data interface{}
}
