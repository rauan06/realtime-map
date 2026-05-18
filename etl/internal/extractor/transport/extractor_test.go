package transport

import (
	"context"
	"testing"
	"time"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

func TestExtractor_BusesAdvanceAndStayOnRoute(t *testing.T) {
	t.Parallel()

	e := New(0)

	// First Extract initializes state with dt=0 — every bus stays at its
	// seeded progress.
	r1, err := e.Extract(context.Background())
	if err != nil {
		t.Fatalf("first extract: %v", err)
	}

	if len(r1) == 0 {
		t.Fatal("expected non-zero fleet on first extract")
	}

	// Simulate clock advancing 60s by rewinding e.last; the second call
	// should advance every bus.
	e.mu.Lock()
	e.last = e.last.Add(-60 * time.Second)
	e.mu.Unlock()

	r2, err := e.Extract(context.Background())
	if err != nil {
		t.Fatalf("second extract: %v", err)
	}

	if len(r2) != len(r1) {
		t.Fatalf("fleet size changed between extracts: %d → %d", len(r1), len(r2))
	}

	moved := countMoved(t, r1, r2)
	if moved == 0 {
		t.Error("no buses moved after 60s — extractor is not advancing state")
	}
}

// countMoved compares two fleet snapshots and asserts every bus stayed
// inside the Kazakhstan bounding box. Returns the number that moved.
func countMoved(t *testing.T, r1, r2 []domain.RawRecord) int {
	t.Helper()

	moved := 0

	for i, rec1 := range r1 {
		rec2 := r2[i]
		lat1, lng1 := mustFloatLatLng(t, rec1)
		lat2, lng2 := mustFloatLatLng(t, rec2)

		if lat1 != lat2 || lng1 != lng2 {
			moved++
		}

		if lat2 < 40 || lat2 > 56 || lng2 < 45 || lng2 > 90 {
			t.Errorf("bus %s left KZ bbox: lat=%v lng=%v", rec2.SourceID, lat2, lng2)
		}
	}

	return moved
}

func mustFloatLatLng(t *testing.T, r domain.RawRecord) (float64, float64) {
	t.Helper()

	lat, ok1 := r.Fields["lat"].(float64)
	lng, ok2 := r.Fields["lng"].(float64)

	if !ok1 || !ok2 {
		t.Fatalf("non-float coords in record: %+v", r.Fields)
	}

	return lat, lng
}

func TestHaversineMeters_KnownDistance(t *testing.T) {
	t.Parallel()

	// Astana (51.17, 71.45) to Almaty (43.24, 76.89) is ~970 km by road,
	// ~960 km great-circle. Allow ±20 km tolerance.
	astana := [2]float64{51.1694, 71.4491}
	almaty := [2]float64{43.2389, 76.8897}

	d := haversineMeters(astana, almaty)
	if d < 940_000 || d > 980_000 {
		t.Errorf("haversine(Astana, Almaty) = %.0f m, expected ~960 km", d)
	}
}
