package geofence

import "testing"

func TestCircle_ContainsRespectsRadius(t *testing.T) {
	c := Circle{Name: "astana", Lat: 51.1694, Lng: 71.4491, RadiusM: 5000}

	// Center: definitely inside.
	if !c.Contains(51.1694, 71.4491) {
		t.Errorf("center should be inside circle")
	}
	// ~2 km north: still inside.
	if !c.Contains(51.1874, 71.4491) {
		t.Errorf("point 2km away should be inside 5km circle")
	}
	// ~10 km north: outside.
	if c.Contains(51.2594, 71.4491) {
		t.Errorf("point 10km away should be outside 5km circle")
	}
}

func TestPolygon_RayCastingMatchesGeometry(t *testing.T) {
	// Unit square covering (lat 0-1, lng 0-1).
	p := Polygon{Name: "unit", Vertices: [][2]float64{
		{0, 0}, {0, 1}, {1, 1}, {1, 0},
	}}

	cases := []struct {
		lat, lng float64
		want     bool
	}{
		{0.5, 0.5, true},
		{0.01, 0.01, true},
		{0.99, 0.99, true},
		{1.5, 0.5, false},
		{-0.5, 0.5, false},
		{0.5, 1.5, false},
		{0.5, -0.5, false},
	}
	for _, c := range cases {
		if got := p.Contains(c.lat, c.lng); got != c.want {
			t.Errorf("Contains(%.2f,%.2f) = %v, want %v", c.lat, c.lng, got, c.want)
		}
	}
}

func TestPolygon_RejectsDegenerate(t *testing.T) {
	p := Polygon{Vertices: [][2]float64{{0, 0}, {1, 1}}}
	if p.Contains(0.5, 0.5) {
		t.Errorf("2-vertex polygon should never contain a point")
	}
}

func TestRegistry_FiltersByLayer(t *testing.T) {
	r := NewRegistry(
		Circle{Name: "flight-only", Layers: []string{"flight"}, Lat: 0, Lng: 0, RadiusM: 1000000},
		Circle{Name: "any", Lat: 0, Lng: 0, RadiusM: 1000000},
	)

	m := r.Match("flight", 0, 0)
	if len(m) != 2 {
		t.Errorf("flight should match both fences, got %v", m)
	}
	m = r.Match("ship", 0, 0)
	if len(m) != 1 || m[0] != "any" {
		t.Errorf("ship should only match unfiltered fence, got %v", m)
	}
}
