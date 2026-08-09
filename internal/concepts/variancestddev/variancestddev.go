// Package variancestddev visualizes why we square deviations before averaging
// them, then take a square root to get back to the original units. Raw
// deviations from the mean always sum to zero (positive and negative cancel),
// so squaring is what lets them accumulate into a real measure of spread —
// variance. But variance is in squared units, so stddev = √variance converts
// it back to something comparable to the data itself.
package variancestddev

import (
	"fmt"
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
		Blurb: "Two dart throwers, 4 throws each, distance from bullseye in inches (left = " +
			"negative, right = positive). Thrower 1: -2,-1,+1,+2. Thrower 2: -8,-4,+4,+8 — " +
			"visibly 4x wilder. Try the obvious fix: average the raw distances. Thrower 1: " +
			"(-2-1+1+2)/4 = 0. Thrower 2: (-8-4+4+8)/4 = 0. BOTH average to exactly zero, even " +
			"though one is clearly way more scattered — deviations from the mean always cancel, " +
			"every time, for any data. Square each deviation first (erases the sign, and punishes " +
			"far misses harder) and average THOSE: Thrower 1 gets variance 2.5, Thrower 2 gets " +
			"variance 40 — finally distinguishable. But 'inches squared' isn't a real unit anyone " +
			"uses, so square-root back down: stddev 1.58 vs. 6.32, in real inches again — the same " +
			"4x gap the raw numbers showed. Drag spread to watch stddev scale linearly while " +
			"variance scales with the square; drag the outlier to see how much one extreme point " +
			"inflates both.",
		Params: []concept.ParamSpec{
			{Key: "spread", Label: "Spread", Min: 0.3, Max: 2.5, Step: 0.1, Def: 1},
			{Key: "outlier", Label: "Outlier push", Min: 0, Max: 10, Step: 0.5, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	spread, outlier := p["spread"], p["outlier"]
	if spread <= 0 {
		spread = 0.1
	}
	xs := Sample(spread, outlier)
	mean := Mean(xs)
	variance := Variance(xs)
	stddev := StdDev(xs)

	lo, hi := xs[0], xs[0]
	for _, x := range xs {
		if x < lo {
			lo = x
		}
		if x > hi {
			hi = x
		}
	}
	pad := (hi-lo)*0.15 + 1
	lo -= pad
	hi += pad

	sq := make([]float64, len(xs))
	maxSq := 0.0
	for i, x := range xs {
		d := x - mean
		sq[i] = d * d
		if sq[i] > maxSq {
			maxSq = sq[i]
		}
	}
	if maxSq == 0 {
		maxSq = 1
	}

	c := viz.New(680, 360, lo, hi, 0, maxSq*1.3)
	c.Axes()
	step := 2.0
	start := math.Ceil(lo/step) * step
	for x := start; x <= hi; x += step {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	// Squared-deviation bars: taller bar = point further from the mean. Bars
	// never go negative, unlike the raw signed deviations they replace.
	for i, x := range xs {
		px := c.X(x)
		top := c.Y(sq[i])
		bottom := c.Y(0)
		color := viz.Good
		if x < mean {
			color = viz.Bad
		}
		c.Rect(px-4, top, 8, bottom-top, color, 0.55)
	}

	c.VLine(mean, viz.Ink, false)

	// The data points themselves, as a row of markers above the tallest bar.
	markerY := maxSq * 1.16
	for _, x := range xs {
		px, py := c.X(x), c.Y(markerY)
		c.Rect(px-3, py-3, 6, 6, viz.Ink, 0.85)
	}

	c.Text(c.X(mean), c.PadT+12, "mean", 12, viz.Ink, "middle")
	c.Text(20, 24, fmt.Sprintf("mean = %.2f", mean), 14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("variance = %.2f (squared units — bar heights above)", variance),
		13, viz.Muted, "start")
	c.Text(20, 62, fmt.Sprintf("stddev = √variance = %.2f (back to the data's own units)", stddev),
		13, viz.Muted, "start")

	return c.String()
}
