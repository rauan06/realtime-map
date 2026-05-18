package multi

import (
	"context"
	"errors"

	"github.com/rauan06/realtime-map/etl/internal/domain"
	"github.com/rauan06/realtime-map/etl/internal/loader"
)

// Loader fans out events to multiple inner loaders. Add() goes to every child;
// Flush() flushes every child and joins any errors so a slow downstream does
// not silently mask a stuck Kafka or ClickHouse sink.
type Loader struct {
	children []loader.Loader
}

func New(children ...loader.Loader) *Loader {
	return &Loader{children: children}
}

func (l *Loader) Add(event domain.KafkaEvent) {
	for _, c := range l.children {
		c.Add(event)
	}
}

func (l *Loader) Flush(ctx context.Context) error {
	var errs []error
	for _, c := range l.children {
		if err := c.Flush(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Len returns the largest backlog across children so the pipeline batches
// based on the slowest sink, not the fastest.
func (l *Loader) Len() int {
	var maxN int
	for _, c := range l.children {
		if n := c.Len(); n > maxN {
			maxN = n
		}
	}
	return maxN
}
