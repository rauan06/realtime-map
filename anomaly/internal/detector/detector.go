// Package detector wraps the Isolation Forest with per-layer feature
// extraction and a rolling sample buffer so the model retrains as new
// observations stream in from Kafka.
package detector

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/rauan06/realtime-map/anomaly/internal/iforest"
)

const (
	defaultWarmup    = 200
	defaultBufferCap = 2000
	defaultThreshold = 0.62
	retrainEvery     = 100
)

// Alert is the result of a positive detection. SourceID identifies the
// flight (icao24) or ship (mmsi); Score is the IF anomaly score; Reason
// is a human-readable description of which feature(s) drove the score.
type Alert struct {
	Layer    string    `json:"layer"`
	SourceID string    `json:"source_id"`
	Lat      float64   `json:"lat"`
	Lng      float64   `json:"lng"`
	Score    float64   `json:"score"`
	Reasons  []string  `json:"reasons"`
	At       time.Time `json:"at"`
}

// FeatureFn extracts an ordered feature vector and the (lat,lng) for a
// payload. ok=false signals "no usable observation; skip".
type FeatureFn func(payload map[string]any) (vec []float64, lat, lng float64, sourceID string, ok bool)

// Detector is the per-layer detector. Safe for concurrent Observe.
type Detector struct {
	layer     string
	extract   FeatureFn
	threshold float64
	warmup    int

	mu      sync.Mutex
	buf     [][]float64
	bufCap  int
	added   int
	lastFit int
	forest  *iforest.Forest
	// per-feature mean/std snapshot used to label which feature(s) drove
	// the anomaly score. Refreshed every retrain.
	means []float64
	stds  []float64
}

type Options struct {
	Layer     string
	Extract   FeatureFn
	Warmup    int
	BufferCap int
	Threshold float64
	Forest    iforest.Options
}

func New(opts Options) *Detector {
	if opts.Warmup <= 0 {
		opts.Warmup = defaultWarmup
	}

	if opts.BufferCap <= 0 {
		opts.BufferCap = defaultBufferCap
	}

	if opts.Threshold <= 0 {
		opts.Threshold = defaultThreshold
	}

	return &Detector{
		layer:     opts.Layer,
		extract:   opts.Extract,
		threshold: opts.Threshold,
		warmup:    opts.Warmup,
		bufCap:    opts.BufferCap,
		forest:    iforest.New(opts.Forest),
	}
}

// Observe ingests a kafka message payload. If the IF is trained and the
// observation scores above threshold, returns a non-nil Alert.
func (d *Detector) Observe(payload []byte) (*Alert, error) {
	var p map[string]any
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, fmt.Errorf("detector unmarshal: %w", err)
	}

	vec, lat, lng, sourceID, ok := d.extract(p)
	if !ok {
		return nil, nil //nolint:nilnil // "no usable observation" is the legitimate skip path
	}

	if !allFinite(vec) {
		return nil, nil //nolint:nilnil // skip NaN/Inf observations
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.appendLocked(vec)

	if !d.isTrainedLocked() {
		return nil, nil //nolint:nilnil // forest not yet trained
	}

	score, err := d.forest.Score(vec)
	if err != nil {
		if errors.Is(err, iforest.ErrDimMismatch) {
			return nil, nil //nolint:nilnil // observation has a different schema; skip
		}

		return nil, fmt.Errorf("detector score: %w", err)
	}

	if score < d.threshold {
		return nil, nil //nolint:nilnil // below alerting threshold is the normal path
	}

	return &Alert{
		Layer:    d.layer,
		SourceID: sourceID,
		Lat:      lat,
		Lng:      lng,
		Score:    score,
		Reasons:  d.driverFeaturesLocked(vec),
		At:       time.Now().UTC(),
	}, nil
}

// appendLocked adds a sample to the rolling buffer and refits the forest
// once a full warmup window has accumulated or every retrainEvery additions.
func (d *Detector) appendLocked(vec []float64) {
	cp := make([]float64, len(vec))
	copy(cp, vec)

	if len(d.buf) < d.bufCap {
		d.buf = append(d.buf, cp)
	} else {
		// Ring overwrite: drop the oldest. We pay an O(n) shift but at the
		// scale we care about (≤ a few thousand) this is dominated by the
		// trees-rebuild step on retrain.
		copy(d.buf, d.buf[1:])
		d.buf[len(d.buf)-1] = cp
	}

	d.added++

	if (d.forest.Dims() == 0 && len(d.buf) >= d.warmup) ||
		(d.forest.Dims() > 0 && d.added-d.lastFit >= retrainEvery) {
		// Inputs are pre-validated above (allFinite + dim consistency via
		// the per-layer feature extractor), so Fit's documented errors
		// cannot trigger here.
		_ = d.forest.Fit(d.buf) //nolint:errcheck // pre-validated inputs
		d.lastFit = d.added
		d.refreshStatsLocked()
	}
}

func (d *Detector) isTrainedLocked() bool { return d.forest.Dims() > 0 }

// refreshStatsLocked recomputes per-feature mean/std on the current buffer
// so driverFeatures can pick the biggest z-score(s).
func (d *Detector) refreshStatsLocked() {
	if len(d.buf) == 0 {
		return
	}

	dims := len(d.buf[0])
	d.means = make([]float64, dims)
	d.stds = make([]float64, dims)

	for _, v := range d.buf {
		for i, x := range v {
			d.means[i] += x
		}
	}

	for i := range d.means {
		d.means[i] /= float64(len(d.buf))
	}

	for _, v := range d.buf {
		for i, x := range v {
			diff := x - d.means[i]
			d.stds[i] += diff * diff
		}
	}

	for i := range d.stds {
		d.stds[i] = math.Sqrt(d.stds[i] / float64(len(d.buf)))
	}
}

// driverFeaturesLocked returns up to maxDrivers feature indices whose
// z-score from the buffer mean is largest, formatted as "f<i>_z=<value>".
// Only features with |z| above driverMinZ are reported.
func (d *Detector) driverFeaturesLocked(vec []float64) []string {
	const (
		maxDrivers = 2
		driverMinZ = 1.5
	)

	if len(d.means) == 0 {
		return nil
	}

	type zscore struct {
		idx int
		z   float64
	}

	zs := make([]zscore, 0, len(vec))

	for i, x := range vec {
		std := d.stds[i]
		if std == 0 {
			continue
		}

		z := math.Abs(x-d.means[i]) / std
		zs = append(zs, zscore{i, z})
	}

	// Take top 2 by |z|.
	for i := range zs {
		for j := i + 1; j < len(zs); j++ {
			if zs[j].z > zs[i].z {
				zs[i], zs[j] = zs[j], zs[i]
			}
		}
	}

	out := make([]string, 0, maxDrivers)

	for i, z := range zs {
		if i >= maxDrivers {
			break
		}

		if z.z < driverMinZ {
			break
		}

		out = append(out, fmt.Sprintf("f%d_z=%.2f", z.idx, z.z))
	}

	return out
}

func allFinite(v []float64) bool {
	for _, x := range v {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return false
		}
	}

	return true
}
