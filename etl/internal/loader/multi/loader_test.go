package multi

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

var (
	errFakeA = errors.New("fake loader A failed")
	errFakeB = errors.New("fake loader B failed")
)

type fakeLoader struct {
	mu       sync.Mutex
	events   []domain.KafkaEvent
	flushed  int
	flushErr error
}

func (f *fakeLoader) Add(ev domain.KafkaEvent) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.events = append(f.events, ev)
}

func (f *fakeLoader) Flush(_ context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.flushed++

	if f.flushErr != nil {
		return f.flushErr
	}

	f.events = f.events[:0]

	return nil
}

func (f *fakeLoader) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return len(f.events)
}

func TestLoader_AddFansOutToAllChildren(t *testing.T) {
	t.Parallel()

	a, b := &fakeLoader{}, &fakeLoader{}
	ml := New(a, b)

	ml.Add(domain.KafkaEvent{Key: "k"})

	if got := a.Len(); got != 1 {
		t.Errorf("loader A: got %d want 1", got)
	}

	if got := b.Len(); got != 1 {
		t.Errorf("loader B: got %d want 1", got)
	}
}

func TestLoader_LenReturnsMaxBacklog(t *testing.T) {
	t.Parallel()

	a, b := &fakeLoader{}, &fakeLoader{}
	a.events = make([]domain.KafkaEvent, 3)
	b.events = make([]domain.KafkaEvent, 7)
	ml := New(a, b)

	if got := ml.Len(); got != 7 {
		t.Errorf("Len: got %d want 7 (max across children)", got)
	}
}

func TestLoader_FlushJoinsErrors(t *testing.T) {
	t.Parallel()

	a := &fakeLoader{flushErr: errFakeA}
	b := &fakeLoader{flushErr: errFakeB}
	ml := New(a, b)

	err := ml.Flush(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, errFakeA) || !errors.Is(err, errFakeB) {
		t.Errorf("expected joined errors with %v + %v, got %v", errFakeA, errFakeB, err)
	}
}
