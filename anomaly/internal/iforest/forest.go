// Package iforest implements an Isolation Forest anomaly detector
// (Liu, Ting, Zhou, 2008). A forest is an ensemble of randomly-built
// trees that recursively partition the feature space. Points that get
// isolated by short paths are anomalies.
//
// The anomaly score is
//
//	s(x, n) = 2^( -E[h(x)] / c(n) )
//
// where E[h(x)] is the mean path length of x across trees and c(n) is the
// expected path length of an unsuccessful BST search over n samples. A
// score close to 1 is highly anomalous; close to 0.5 or below is normal.
//
// This implementation is intentionally compact (~250 LOC) and aims for
// correctness, not raw speed. It is suitable for streaming use cases with
// a few thousand samples per retrain.
package iforest

import (
	"errors"
	"math"
	"math/rand/v2"
)

const (
	defaultTrees      = 100
	defaultSampleSize = 256
)

var (
	// ErrUntrained is returned when Score is called before any Fit.
	ErrUntrained = errors.New("isolation forest: untrained — call Fit first")

	// ErrEmptyInput is returned when Fit receives no samples.
	ErrEmptyInput = errors.New("isolation forest: empty training set")

	// ErrDimMismatch is returned when Score is called on a vector whose
	// dimensionality differs from the training set.
	ErrDimMismatch = errors.New("isolation forest: input dimensionality mismatch")
)

// Forest is an ensemble of isolation trees. The Score method is safe for
// concurrent use after Fit returns; Fit itself is not concurrent-safe.
type Forest struct {
	trees      []*tree
	sampleSize int
	dims       int
	rng        *rand.Rand
}

// Options tunes the forest. Zero values mean "use defaults".
type Options struct {
	NumTrees   int
	SampleSize int
	Seed       uint64
}

func New(opts Options) *Forest {
	if opts.NumTrees <= 0 {
		opts.NumTrees = defaultTrees
	}

	if opts.SampleSize <= 0 {
		opts.SampleSize = defaultSampleSize
	}

	const (
		defaultSeed  uint64 = 0xDEADBEEF
		seedScramble uint64 = 0x5A5A
	)

	if opts.Seed == 0 {
		opts.Seed = defaultSeed
	}

	return &Forest{
		sampleSize: opts.SampleSize,
		trees:      make([]*tree, 0, opts.NumTrees),
		rng:        rand.New(rand.NewPCG(opts.Seed, opts.Seed^seedScramble)), //nolint:gosec // statistical model, not security
	}
}

// Fit builds the forest on samples. Each sample is a fixed-length feature
// vector; all samples must share the same length. The number of trees is
// taken from cap(f.trees) — set via Options.NumTrees on New.
func (f *Forest) Fit(samples [][]float64) error {
	if len(samples) == 0 {
		return ErrEmptyInput
	}

	dims := len(samples[0])
	for i, s := range samples {
		if len(s) != dims {
			return ErrDimMismatch
		}

		_ = i
	}

	f.dims = dims
	f.trees = f.trees[:0]

	limit := int(math.Ceil(math.Log2(float64(f.sampleSize))))
	want := cap(f.trees)

	for range want {
		sample := subsample(samples, f.sampleSize, f.rng)
		t := buildTree(sample, 0, limit, f.rng)
		f.trees = append(f.trees, t)
	}

	return nil
}

// Score returns the anomaly score in [0,1] for x. A score above ~0.6
// typically indicates an anomaly. Scores in [0.4, 0.6] are inconclusive.
func (f *Forest) Score(x []float64) (float64, error) {
	if len(f.trees) == 0 {
		return 0, ErrUntrained
	}

	if len(x) != f.dims {
		return 0, ErrDimMismatch
	}

	var sum float64
	for _, t := range f.trees {
		sum += pathLength(t, x, 0)
	}

	mean := sum / float64(len(f.trees))
	c := cFactor(f.sampleSize)

	if c == 0 {
		return 0, nil
	}

	const scoreBase = 2.0

	return math.Pow(scoreBase, -mean/c), nil
}

// Dims reports the trained feature dimensionality, or 0 if Fit has not run.
func (f *Forest) Dims() int { return f.dims }

// ── tree internals ───────────────────────────────────────────────────────

type tree struct {
	// internal node fields
	feature int
	split   float64
	left    *tree
	right   *tree
	// leaf field
	size int
}

func (t *tree) leaf() bool { return t.left == nil && t.right == nil }

func buildTree(samples [][]float64, depth, limit int, rng *rand.Rand) *tree {
	n := len(samples)
	if n <= 1 || depth >= limit {
		return &tree{size: n}
	}

	dims := len(samples[0])
	// Pick a random feature with non-zero range; if all features collapsed
	// (every sample identical), we're done.
	picked := -1

	var lo, hi float64

	for range dims {
		f := rng.IntN(dims)
		lo, hi = featureRange(samples, f)

		if hi > lo {
			picked = f

			break
		}
	}

	if picked < 0 {
		return &tree{size: n}
	}

	split := lo + rng.Float64()*(hi-lo)

	leftSet, rightSet := partition(samples, picked, split)

	return &tree{
		feature: picked,
		split:   split,
		left:    buildTree(leftSet, depth+1, limit, rng),
		right:   buildTree(rightSet, depth+1, limit, rng),
	}
}

func pathLength(t *tree, x []float64, depth int) float64 {
	if t.leaf() {
		return float64(depth) + cFactor(t.size)
	}

	if x[t.feature] < t.split {
		return pathLength(t.left, x, depth+1)
	}

	return pathLength(t.right, x, depth+1)
}

// cFactor is c(n) = 2 H(n-1) - 2(n-1)/n with H(i) ≈ ln(i) + Euler-Mascheroni.
// It approximates the average path length of an unsuccessful BST search
// over n samples — the normalising term in the IF anomaly score.
func cFactor(n int) float64 {
	const (
		euler   = 0.5772156649015329
		hScale  = 2.0
		bstTerm = 2.0
	)

	if n <= 1 {
		return 0
	}

	h := math.Log(float64(n-1)) + euler

	return hScale*h - bstTerm*float64(n-1)/float64(n)
}

func featureRange(samples [][]float64, f int) (lo, hi float64) {
	lo = samples[0][f]
	hi = lo

	for _, s := range samples[1:] {
		v := s[f]
		if v < lo {
			lo = v
		}

		if v > hi {
			hi = v
		}
	}

	return lo, hi
}

func partition(samples [][]float64, f int, split float64) (left, right [][]float64) {
	for _, s := range samples {
		if s[f] < split {
			left = append(left, s)
		} else {
			right = append(right, s)
		}
	}

	return left, right
}

// subsample returns up to k samples drawn without replacement.
func subsample(samples [][]float64, k int, rng *rand.Rand) [][]float64 {
	n := len(samples)
	if k >= n {
		out := make([][]float64, n)
		copy(out, samples)

		return out
	}

	// Fisher-Yates partial shuffle.
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}

	for i := range k {
		j := i + rng.IntN(n-i)
		idx[i], idx[j] = idx[j], idx[i]
	}

	out := make([][]float64, k)
	for i := range k {
		out[i] = samples[idx[i]]
	}

	return out
}
