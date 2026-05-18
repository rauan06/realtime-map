// Package transport synthesizes a Kazakhstan public-transport feed.
//
// There is no widely-available realtime public-transport API for Kazakhstan
// cities (Almaty's "Onay" and Astana's smart-city feeds require auth and are
// not openly documented). For the realtime-map dashboard, we generate a
// plausible bus fleet moving on great-circle segments between major KZ
// cities. The payload shape (vehicle_id, lat, lng, bearing, speed, status,
// label) matches what an MBTA-shaped extractor would return so the dashboard
// renderer is source-agnostic.
//
// If a real API becomes available, this file is the only thing that needs
// to change.
package transport

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

const (
	earthRadiusM    = 6371000.0
	degToRad        = math.Pi / 180.0
	radToDeg        = 180.0 / math.Pi
	defaultSpeedMS  = 18.0 // ~65 km/h intercity coach average
	cityLoopSpeedMS = 8.0  // ~30 km/h within-city bus

	half          = 2.0
	fullCircleDeg = 360.0
	halfCircleDeg = 180.0
)

type Extractor struct {
	timeout time.Duration

	mu     sync.Mutex
	last   time.Time
	routes []*route
}

func New(timeout time.Duration) *Extractor {
	return &Extractor{
		timeout: timeout,
		routes:  buildRoutes(),
	}
}

// Extract advances every bus by (now-last) seconds along its route and
// returns the resulting fleet snapshot.
func (e *Extractor) Extract(_ context.Context) ([]domain.RawRecord, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now().UTC()
	dt := 0.0

	if !e.last.IsZero() {
		dt = now.Sub(e.last).Seconds()
	}

	e.last = now

	var out []domain.RawRecord

	for _, r := range e.routes {
		distM := haversineMeters(r.start, r.end)
		if distM <= 0 {
			continue
		}

		for _, b := range r.buses {
			step := r.speedMS * dt / distM
			if b.forward {
				b.progress += step
				if b.progress >= 1 {
					b.progress = 1
					b.forward = false
				}
			} else {
				b.progress -= step
				if b.progress <= 0 {
					b.progress = 0
					b.forward = true
				}
			}

			lat, lng := interp(r.start, r.end, b.progress)
			bearing := bearingDeg(r.start, r.end, b.forward)

			out = append(out, domain.RawRecord{
				SourceID:  b.id,
				Timestamp: now,
				Fields: map[string]any{
					"vehicle_id": b.id,
					"label":      b.label,
					"lat":        lat,
					"lng":        lng,
					"bearing":    bearing,
					"speed":      r.speedMS,
					"status":     status(b),
					"route":      r.name,
				},
			})
		}
	}

	return out, nil
}

type route struct {
	name    string
	start   [2]float64 // [lat, lng]
	end     [2]float64
	speedMS float64
	buses   []*bus
}

type bus struct {
	id       string
	label    string
	progress float64 // 0..1 along route start→end
	forward  bool
}

// buildRoutes seeds the synthetic fleet. Each route gets a handful of buses
// spaced along it with alternating directions, so the dashboard shows a
// realistic spread instead of all buses clumped at the depot.
func buildRoutes() []*route {
	// Major Kazakhstan cities.
	almaty := [2]float64{43.2389, 76.8897}
	astana := [2]float64{51.1694, 71.4491}
	shymkent := [2]float64{42.3417, 69.5901}
	karaganda := [2]float64{49.8047, 73.1094}
	pavlodar := [2]float64{52.2873, 76.9674}
	aktobe := [2]float64{50.2839, 57.2294}
	semey := [2]float64{50.4111, 80.2275}
	kostanay := [2]float64{53.2198, 63.6354}
	taraz := [2]float64{42.9000, 71.3667}
	atyrau := [2]float64{47.0945, 51.9238}

	specs := []struct {
		name        string
		start, end  [2]float64
		buses       int
		speedMS     float64
		labelPrefix string
	}{
		{"almaty-shymkent", almaty, shymkent, 4, defaultSpeedMS, "ALM-SHY"},
		{"almaty-taraz", almaty, taraz, 3, defaultSpeedMS, "ALM-TAR"},
		{"almaty-astana", almaty, astana, 5, defaultSpeedMS, "ALM-AST"},
		{"astana-karaganda", astana, karaganda, 4, defaultSpeedMS, "AST-KGD"},
		{"astana-pavlodar", astana, pavlodar, 3, defaultSpeedMS, "AST-PVL"},
		{"astana-kostanay", astana, kostanay, 3, defaultSpeedMS, "AST-KST"},
		{"astana-semey", astana, semey, 3, defaultSpeedMS, "AST-SEM"},
		{"aktobe-atyrau", aktobe, atyrau, 3, defaultSpeedMS, "AKT-ATY"},
		// Small city-loop "routes" simulating intra-city buses.
		{"almaty-loop", almaty, [2]float64{43.2700, 76.9300}, 4, cityLoopSpeedMS, "ALM"},
		{"astana-loop", astana, [2]float64{51.1300, 71.4900}, 4, cityLoopSpeedMS, "AST"},
	}

	out := make([]*route, 0, len(specs))

	const (
		seedA uint64 = 1
		seedB uint64 = 2
	)

	// Deterministic seed so repeated restarts show the same fleet. This is
	// a synthetic data generator — math/rand/v2 is fit-for-purpose; no
	// cryptographic randomness needed.
	rng := rand.New(rand.NewPCG(seedA, seedB)) //nolint:gosec // synthetic fleet seed

	for _, s := range specs {
		r := &route{name: s.name, start: s.start, end: s.end, speedMS: s.speedMS}

		const dirCoinflip = 2

		for i := range s.buses {
			r.buses = append(r.buses, &bus{
				id:       fmt.Sprintf("%s-%02d", s.labelPrefix, i+1),
				label:    fmt.Sprintf("%s #%d", s.labelPrefix, i+1),
				progress: rng.Float64(),
				forward:  rng.IntN(dirCoinflip) == 0,
			})
		}

		out = append(out, r)
	}

	return out
}

// haversineMeters returns the great-circle distance between two lat/lng
// points using the mean Earth radius.
func haversineMeters(a, b [2]float64) float64 {
	lat1, lng1 := a[0]*degToRad, a[1]*degToRad
	lat2, lng2 := b[0]*degToRad, b[1]*degToRad
	dLat := lat2 - lat1
	dLng := lng2 - lng1
	h := math.Sin(dLat/half)*math.Sin(dLat/half) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLng/half)*math.Sin(dLng/half)

	return half * earthRadiusM * math.Asin(math.Sqrt(h))
}

// interp returns the lat/lng at fractional progress along a great-circle
// segment from a to b. p=0 → a, p=1 → b.
func interp(a, b [2]float64, p float64) (lat, lng float64) {
	if p <= 0 {
		return a[0], a[1]
	}

	if p >= 1 {
		return b[0], b[1]
	}

	lat1 := a[0] * degToRad
	lng1 := a[1] * degToRad
	lat2 := b[0] * degToRad
	lng2 := b[1] * degToRad

	d := half * math.Asin(math.Sqrt(
		math.Pow(math.Sin((lat2-lat1)/half), half)+
			math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin((lng2-lng1)/half), half),
	))
	if d == 0 {
		return a[0], a[1]
	}

	A := math.Sin((1-p)*d) / math.Sin(d)
	B := math.Sin(p*d) / math.Sin(d)

	x := A*math.Cos(lat1)*math.Cos(lng1) + B*math.Cos(lat2)*math.Cos(lng2)
	y := A*math.Cos(lat1)*math.Sin(lng1) + B*math.Cos(lat2)*math.Sin(lng2)
	z := A*math.Sin(lat1) + B*math.Sin(lat2)

	lat = math.Atan2(z, math.Sqrt(x*x+y*y)) * radToDeg
	lng = math.Atan2(y, x) * radToDeg

	return lat, lng
}

// bearingDeg returns the initial great-circle bearing from a to b, flipped
// when the bus is traveling in reverse.
func bearingDeg(a, b [2]float64, forward bool) float64 {
	lat1 := a[0] * degToRad
	lat2 := b[0] * degToRad
	dLng := (b[1] - a[1]) * degToRad

	y := math.Sin(dLng) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLng)
	br := math.Atan2(y, x) * radToDeg
	br = math.Mod(br+fullCircleDeg, fullCircleDeg)

	if !forward {
		br = math.Mod(br+halfCircleDeg, fullCircleDeg)
	}

	return br
}

func status(b *bus) string {
	if b.progress >= 1 || b.progress <= 0 {
		return "STOPPED_AT_STATION"
	}

	return "IN_TRANSIT_TO"
}
