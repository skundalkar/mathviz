// Package variancestddev visualizes why we square deviations before averaging
// them, then take a square root to get back to the original units. Raw
// deviations from the mean always sum to zero (positive and negative cancel),
// so squaring is what lets them accumulate into a real measure of spread —
// variance. But variance is in squared units, so stddev = √variance converts
// it back to something comparable to the data itself.
package variancestddev

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

// baseSample is a small, fixed dataset (deviations from its own mean, roughly)
// that every rendered sample is derived from. Keeping it as a literal — rather
// than generating it — keeps the picture concrete: these are "quiz score"
// style deltas, not synthetic noise.
var baseSample = []float64{-4, -3, -1, 0, 1, 1, 2, 4, 6}

// Sample scales baseSample by spread (so its variance scales with spread²)
// and, if outlier is nonzero, pushes the last point further out — showing how
// much a single extreme value can inflate variance versus a robust measure.
// Pure function: same inputs always produce the same slice.
func Sample(spread, outlier float64) []float64 {
	xs := make([]float64, len(baseSample))
	for i, v := range baseSample {
		xs[i] = v * spread
	}
	if len(xs) > 0 {
		xs[len(xs)-1] += outlier
	}
	return xs
}

// Mean is the arithmetic mean of xs. Returns 0 for an empty slice.
func Mean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

// Variance is the population variance of xs: the mean of the squared
// deviations from the mean. Returns 0 for an empty slice.
func Variance(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	mean := Mean(xs)
	sum := 0.0
	for _, x := range xs {
		d := x - mean
		sum += d * d
	}
	return sum / float64(len(xs))
}

// StdDev is the population standard deviation of xs: √Variance(xs). Same
// units as the data, unlike Variance.
func StdDev(xs []float64) float64 {
	return math.Sqrt(Variance(xs))
}

func init() {
	concept.Register(concept.Concept{
		ID:    "variance-vs-stddev",
		Title: "Variance vs. standard deviation",
		Blurb: "Raw deviations from the mean always sum to zero — positives and negatives " +
			"cancel out, so they can't measure spread directly. Squaring each deviation makes " +
			"every term positive (and punishes outliers hardest, since they get squared too), " +
			"so the average squared deviation — the variance — actually grows with spread. The " +
			"catch: variance is in squared units. Taking its square root, the standard " +
			"deviation, brings the number back to the same units as the data. Drag spread to " +
			"see stddev scale linearly while variance scales with the square; drag the outlier " +
			"to see how much one extreme point can inflate both.",
		Params: []concept.ParamSpec{
			{Key: "spread", Label: "Spread", Min: 0.3, Max: 2.5, Step: 0.1, Def: 1},
			{Key: "outlier", Label: "Outlier push", Min: 0, Max: 10, Step: 0.5, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(680, 360, 0, 1, 0, 1).String()
}
