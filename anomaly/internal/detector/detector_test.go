package detector

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/rauan06/realtime-map/anomaly/internal/iforest"
)

const (
	shipLayer    = "ship"
	keyMMSI      = "mmsi"
	keyLat       = "lat"
	keyLng       = "lng"
	keySOG       = "sog"
	keyCOG       = "cog"
	keyHeading   = "heading"
	keyNavStat   = "nav_stat"
	speedingShip = "speeding-tanker"
)

// TestObserve_FlagsOutlierAfterWarmup feeds 250 ships with normal SOG
// around 8 knots, then injects one ship with SOG=80 knots and expects an
// alert. (A real container ship maxes ~25 kn — 80 is implausible.)
func TestObserve_FlagsOutlierAfterWarmup(t *testing.T) {
	t.Parallel()

	d := New(Options{
		Layer:     shipLayer,
		Extract:   ShipFeatures,
		Warmup:    200,
		Threshold: 0.6,
		Forest:    iforest.Options{NumTrees: 50, SampleSize: 128, Seed: 7},
	})

	rng := rand.New(rand.NewPCG(11, 13)) //nolint:gosec // test data generator

	const warmup = 250

	for i := range warmup {
		payload := map[string]any{
			keyMMSI:    fmt.Sprintf("normal-%d", i),
			keyLat:     55.0 + rng.NormFloat64()*0.5,
			keyLng:     20.0 + rng.NormFloat64()*0.5,
			keySOG:     8.0 + rng.NormFloat64()*1.5,
			keyCOG:     90.0 + rng.NormFloat64()*15,
			keyHeading: 90.0 + rng.NormFloat64()*15,
			keyNavStat: 0.0,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal normal[%d]: %v", i, err)
		}

		if _, err := d.Observe(body); err != nil {
			t.Fatalf("observe %d: %v", i, err)
		}
	}

	outlier := map[string]any{
		keyMMSI:    speedingShip,
		keyLat:     55.0,
		keyLng:     20.0,
		keySOG:     80.0,
		keyCOG:     90.0,
		keyHeading: 270.0,
		keyNavStat: 7.0,
	}

	body, err := json.Marshal(outlier)
	if err != nil {
		t.Fatalf("marshal outlier: %v", err)
	}

	alert, err := d.Observe(body)
	if err != nil {
		t.Fatalf("observe outlier: %v", err)
	}

	if alert == nil {
		t.Fatal("expected alert for sog=80 ship, got nil")
	}

	if alert.Score < 0.6 {
		t.Errorf("alert score %.3f below threshold 0.6", alert.Score)
	}

	if alert.SourceID != speedingShip {
		t.Errorf("source_id = %q want %q", alert.SourceID, speedingShip)
	}

	if len(alert.Reasons) == 0 {
		t.Error("expected at least one driver feature in Reasons")
	}
}

func TestObserve_NoAlertDuringWarmup(t *testing.T) {
	t.Parallel()

	d := New(Options{
		Layer:   shipLayer,
		Extract: ShipFeatures,
		Warmup:  500,
		Forest:  iforest.Options{NumTrees: 20, SampleSize: 64, Seed: 1},
	})

	payload := map[string]any{keyMMSI: "x", keyLat: 0.0, keyLng: 0.0, keySOG: 1000.0}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	alert, err := d.Observe(body)
	if err != nil {
		t.Fatalf("observe: %v", err)
	}

	if alert != nil {
		t.Errorf("expected no alert during warmup, got %+v", alert)
	}
}

func TestFlightFeatures_RejectsMissingID(t *testing.T) {
	t.Parallel()

	vec, lat, lng, id, ok := FlightFeatures(map[string]any{keyLat: 1.0, keyLng: 2.0})
	if ok {
		t.Error("expected ok=false when icao24 missing")
	}

	if vec != nil || lat != 0 || lng != 0 || id != "" {
		t.Errorf("expected zero return values when ok=false, got vec=%v lat=%v lng=%v id=%q",
			vec, lat, lng, id)
	}
}
