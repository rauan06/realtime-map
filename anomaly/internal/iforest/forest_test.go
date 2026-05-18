package iforest

import (
	"errors"
	"math/rand/v2"
	"testing"
)

func TestForest_RejectsUntrainedScore(t *testing.T) {
	t.Parallel()

	f := New(Options{})
	if _, err := f.Score([]float64{0, 0}); !errors.Is(err, ErrUntrained) {
		t.Errorf("expected ErrUntrained, got %v", err)
	}
}

func TestForest_RejectsEmptyAndMismatchedFit(t *testing.T) {
	t.Parallel()

	f := New(Options{})
	if err := f.Fit(nil); !errors.Is(err, ErrEmptyInput) {
		t.Errorf("expected ErrEmptyInput, got %v", err)
	}

	mismatch := [][]float64{{1, 2, 3}, {1, 2}}
	if err := f.Fit(mismatch); !errors.Is(err, ErrDimMismatch) {
		t.Errorf("expected ErrDimMismatch, got %v", err)
	}
}

func TestForest_GivesHigherScoreToOutliers(t *testing.T) {
	t.Parallel()

	// Build a 2D blob centered at (0,0) with σ≈1 and a clear outlier
	// at (50, 50). The outlier should score noticeably higher than any
	// of the blob points.
	rng := rand.New(rand.NewPCG(7, 11)) //nolint:gosec // test data generator

	const blob = 400

	samples := make([][]float64, 0, blob+1)

	for range blob {
		samples = append(samples, []float64{rng.NormFloat64(), rng.NormFloat64()})
	}

	outlier := []float64{50, 50}
	samples = append(samples, outlier)

	f := New(Options{NumTrees: 100, SampleSize: 128, Seed: 42})
	if err := f.Fit(samples); err != nil {
		t.Fatalf("fit: %v", err)
	}

	outlierScore, err := f.Score(outlier)
	if err != nil {
		t.Fatalf("score outlier: %v", err)
	}

	// Mean score of 50 random blob points.
	var blobSum float64

	const blobChecks = 50

	for i := range blobChecks {
		s, err := f.Score(samples[i])
		if err != nil {
			t.Fatalf("score blob[%d]: %v", i, err)
		}

		blobSum += s
	}

	blobMean := blobSum / float64(blobChecks)
	if outlierScore <= blobMean+0.1 {
		t.Errorf("expected outlier score (%.3f) >> blob mean (%.3f) by ≥ 0.1", outlierScore, blobMean)
	}

	if outlierScore < 0.55 {
		t.Errorf("expected outlier score > 0.55, got %.3f", outlierScore)
	}
}

func TestForest_ScoreDimMismatch(t *testing.T) {
	t.Parallel()

	f := New(Options{NumTrees: 10, SampleSize: 32, Seed: 1})
	if err := f.Fit([][]float64{{0, 0}, {1, 1}}); err != nil {
		t.Fatalf("fit: %v", err)
	}

	if _, err := f.Score([]float64{0, 0, 0}); !errors.Is(err, ErrDimMismatch) {
		t.Errorf("expected ErrDimMismatch, got %v", err)
	}
}

func TestCFactor_BoundaryValues(t *testing.T) {
	t.Parallel()

	if got := cFactor(0); got != 0 {
		t.Errorf("cFactor(0) = %v, want 0", got)
	}

	if got := cFactor(1); got != 0 {
		t.Errorf("cFactor(1) = %v, want 0", got)
	}

	if got := cFactor(2); got <= 0 {
		t.Errorf("cFactor(2) = %v, want > 0", got)
	}
}
