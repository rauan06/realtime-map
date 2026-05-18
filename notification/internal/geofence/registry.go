package geofence

import (
	"encoding/json"
	"fmt"
	"os"
)

// Registry holds the set of geofences the notifier checks against.
type Registry struct {
	fences []Fence
}

// fenceFile is the on-disk JSON shape. Circles and polygons live in separate
// arrays so the schema stays trivially Go-typed; a "geometry" discriminator
// would cost more than it saves at this scale.
type fenceFile struct {
	Circles  []Circle  `json:"circles"`
	Polygons []Polygon `json:"polygons"`
}

func LoadFromFile(path string) (*Registry, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read geofence file %s: %w", path, err)
	}
	var ff fenceFile
	if err := json.Unmarshal(body, &ff); err != nil {
		return nil, fmt.Errorf("parse geofence file: %w", err)
	}
	r := &Registry{}
	for _, c := range ff.Circles {
		r.fences = append(r.fences, c)
	}
	for _, p := range ff.Polygons {
		r.fences = append(r.fences, p)
	}
	return r, nil
}

func NewRegistry(fences ...Fence) *Registry {
	return &Registry{fences: append([]Fence(nil), fences...)}
}

// Match returns the names of every fence containing (lat,lng) for the given layer.
// A fence with no layers configured matches every layer.
func (r *Registry) Match(layer string, lat, lng float64) []string {
	var matches []string
	for _, f := range r.fences {
		layers := f.GetLayers()
		if len(layers) > 0 && !contains(layers, layer) {
			continue
		}
		if f.Contains(lat, lng) {
			matches = append(matches, f.GetName())
		}
	}
	return matches
}

func (r *Registry) Len() int { return len(r.fences) }

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
