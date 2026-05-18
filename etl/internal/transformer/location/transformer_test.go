package location

import (
	"testing"
	"time"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

func TestTransform_StampsSourceAndTimestamp(t *testing.T) {
	tx := New()
	ts := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)

	events, err := tx.Transform([]domain.RawRecord{{
		SourceID:  "icao24-abc",
		Timestamp: ts,
		Fields:    map[string]interface{}{"lat": 51.0, "lng": 71.0},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	ev := events[0]
	if ev.Key != "icao24-abc" {
		t.Errorf("key: got %q want %q", ev.Key, "icao24-abc")
	}
	fields, ok := ev.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("data is %T, want map[string]interface{}", ev.Data)
	}
	if fields["source_id"] != "icao24-abc" {
		t.Errorf("source_id not stamped: %+v", fields["source_id"])
	}
	wantTS := ts.Format(time.RFC3339)
	if fields["timestamp"] != wantTS {
		t.Errorf("timestamp: got %v want %v", fields["timestamp"], wantTS)
	}
}

func TestTransform_EmptyInputProducesEmptyOutput(t *testing.T) {
	tx := New()
	events, err := tx.Transform(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}
