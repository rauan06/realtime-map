package domain

import "errors"

var (
	ErrConfigFileLoad = errors.New("error loading .env file")
	ErrConfigParse    = errors.New("error parsing config")
	ErrUnknownSource  = errors.New("unknown source")

	// ErrUpstreamStatus is returned by extractors when the upstream API
	// responds with a non-2xx status.
	ErrUpstreamStatus = errors.New("upstream non-OK status")

	// ErrNoLoaders is returned when a pipeline is built with no enabled sink.
	ErrNoLoaders = errors.New("no loaders enabled for source")

	// ErrKafkaFlushPending is returned by the Kafka loader when messages
	// remain in the producer queue after the flush timeout.
	ErrKafkaFlushPending = errors.New("kafka flush left messages in queue")
)
