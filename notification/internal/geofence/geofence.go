// Package geofence implements the two geometric primitives used by the
// notification service: circle (haversine distance) and polygon (ray casting).
// Both Contains methods take WGS84 (lat, lng) degrees.
package geofence

import "math"

const (
	earthRadiusM    = 6371000.0
	degToRad        = math.Pi / 180.0
	minPolygonSides = 3
	half            = 2.0
)

// Fence is anything that can answer "does this lat/lng sit inside me?". Named
// fences are keyed by Name in the alert state map, so two fences sharing a name
// would collide — keep names unique in the config.
type Fence interface {
	GetName() string
	GetLayers() []string
	Contains(lat, lng float64) bool
}

// Circle is a haversine-bounded fence (center + meter radius). Good enough
// for the relatively small "airport / port / depot" use cases on the map.
type Circle struct {
	Name      string   `json:"name"`
	Layers    []string `json:"layers,omitempty"`
	Lat       float64  `json:"lat"`
	Lng       float64  `json:"lng"`
	RadiusM   float64  `json:"radius_m"`
	AlertType string   `json:"alert_type,omitempty"`
}

func (c Circle) GetName() string     { return c.Name }
func (c Circle) GetLayers() []string { return c.Layers }
func (c Circle) Contains(lat, lng float64) bool {
	return haversineMeters(c.Lat, c.Lng, lat, lng) <= c.RadiusM
}

// Polygon is a ring of lat/lng vertices (last need not equal first; we close
// implicitly). Ray-casting against great-circle edges is overkill here —
// linear interpolation is fine at the scales we display.
type Polygon struct {
	Name      string       `json:"name"`
	Layers    []string     `json:"layers,omitempty"`
	Vertices  [][2]float64 `json:"vertices"` // each [lat, lng]
	AlertType string       `json:"alert_type,omitempty"`
}

func (p Polygon) GetName() string     { return p.Name }
func (p Polygon) GetLayers() []string { return p.Layers }

func (p Polygon) Contains(lat, lng float64) bool {
	n := len(p.Vertices)
	if n < minPolygonSides {
		return false
	}

	inside := false
	j := n - 1

	for i := range n {
		yi, xi := p.Vertices[i][0], p.Vertices[i][1]
		yj, xj := p.Vertices[j][0], p.Vertices[j][1]
		intersect := ((yi > lat) != (yj > lat)) &&
			(lng < (xj-xi)*(lat-yi)/(yj-yi)+xi)

		if intersect {
			inside = !inside
		}

		j = i
	}

	return inside
}

// haversineMeters returns the great-circle distance between two points in
// meters using the standard mean Earth radius (6371 km).
func haversineMeters(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := (lat2 - lat1) * degToRad
	dLng := (lng2 - lng1) * degToRad
	a := math.Sin(dLat/half)*math.Sin(dLat/half) +
		math.Cos(lat1*degToRad)*math.Cos(lat2*degToRad)*math.Sin(dLng/half)*math.Sin(dLng/half)

	return half * earthRadiusM * math.Asin(math.Sqrt(a))
}
