package road

import (
	"fmt"
	"math"
	"testing"
)

// wgs84ToWebMercator is the forward projection (mirrors the inverse in
// extractor.go) used only to build round-trip test inputs.
func wgs84ToWebMercator(lng, lat float64) (float64, float64) {
	const r = 6378137.0
	x := lng * r * math.Pi / 180.0
	y := math.Log(math.Tan((90.0+lat)*math.Pi/360.0)) * r
	return x, y
}

func TestParseMultiLineString_RoundTripsWebMercator(t *testing.T) {
	// (0,0) in Web Mercator → (0,0) WGS84
	got := parseMultiLineString("MULTILINESTRING((0 0, 0 0))")
	if len(got) != 2 {
		t.Fatalf("got %d points, want 2", len(got))
	}
	for _, p := range got {
		if math.Abs(p[0]) > 1e-6 || math.Abs(p[1]) > 1e-6 {
			t.Errorf("expected (0,0), got %+v", p)
		}
	}

	// Forward-project Astana, parse it back, verify drift is sub-meter.
	x, y := wgs84ToWebMercator(71.4491, 51.1694)
	wkt := fmt.Sprintf("MULTILINESTRING((%f %f))", x, y)
	got = parseMultiLineString(wkt)
	if len(got) != 1 {
		t.Fatalf("got %d points, want 1", len(got))
	}
	if math.Abs(got[0][0]-71.4491) > 1e-4 || math.Abs(got[0][1]-51.1694) > 1e-4 {
		t.Errorf("Astana roundtrip drift: got %+v want ~(71.4491, 51.1694)", got[0])
	}
}

func TestParseMultiLineString_HandlesEmptyAndJunk(t *testing.T) {
	if got := parseMultiLineString("MULTILINESTRING EMPTY"); got != nil {
		t.Errorf("expected nil for EMPTY, got %+v", got)
	}
	if got := parseMultiLineString("POINT(1 2)"); got != nil {
		t.Errorf("expected nil for wrong type, got %+v", got)
	}
	if got := parseMultiLineString(""); got != nil {
		t.Errorf("expected nil for empty string, got %+v", got)
	}
}

func TestParseMultiLineString_MultipleLines(t *testing.T) {
	// Two sub-lines, each with two points → 4 flattened points.
	wkt := "MULTILINESTRING((0 0, 0 1000),(1000 0, 1000 1000))"
	got := parseMultiLineString(wkt)
	if len(got) != 4 {
		t.Fatalf("got %d points, want 4: %+v", len(got), got)
	}
}
