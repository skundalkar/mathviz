// Package lln visualizes the law of large numbers: a running average of
// repeated coin flips, jumping around wildly at small n and settling in on
// the true probability as n grows -- without ever guaranteeing that any
// single next flip "corrects" the average.
package lln

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "law-of-large-numbers",
		Seq:   37,
		Title: "Law of large numbers (the running average settles down)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "p", Label: "True probability p", Min: 0.05, Max: 0.95, Step: 0.05, Def: 0.5},
			{Key: "n", Label: "Number of flips (n)", Min: 1, Max: 300, Step: 1, Def: 50},
		},
		Render: render,
	})
}

// hash01 turns an integer seed into a deterministic pseudo-random number in
// [0, 1) via a fixed irrational-multiplier trick. It's not a real random
// number generator -- it's a fixed, reproducible sequence that merely looks
// patternless, which is exactly what Render needs (pure: same input, same
// output, no time, no math/rand).
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// Outcome returns the result of the `trial`-th coin flip of a coin that
// comes up heads with true probability p: 1 for heads, 0 for tails.
// Deterministic in trial and p -- the same trial number always reproduces
// the same flip, so a whole run of flips can be recomputed exactly.
func Outcome(trial int, p float64) float64 {
	if hash01(trial) < p {
		return 1
	}
	return 0
}

// RunningAverages returns the running average of Outcome(1,p), Outcome(2,p),
// ..., Outcome(maxN,p) after each flip: result[i] is the average of the
// first i+1 flips. maxN < 1 is treated as 1.
func RunningAverages(maxN int, p float64) []float64 {
	if maxN < 1 {
		maxN = 1
	}
	out := make([]float64, maxN)
	sum := 0.0
	for i := 0; i < maxN; i++ {
		sum += Outcome(i+1, p)
		out[i] = sum / float64(i+1)
	}
	return out
}

// RunningAverage returns the average of the first n flips of a coin with
// true probability p -- the single value RunningAverages(n,p)'s last entry
// holds, computed directly for callers that only need one point.
func RunningAverage(n int, p float64) float64 {
	if n < 1 {
		n = 1
	}
	sum := 0.0
	for i := 1; i <= n; i++ {
		sum += Outcome(i, p)
	}
	return sum / float64(n)
}

// maxFlips is the fixed right edge of the curve's x-axis, matching the
// "Number of flips" slider's own Max so the whole slider range stays
// visible on the plot.
const maxFlips = 300

func render(params map[string]float64) string {
	trueP := params["p"]
	n := int(params["n"])
	if n < 1 {
		n = 1
	}

	series := RunningAverages(maxFlips, trueP)

	c := viz.New(680, 420, 0, maxFlips, 0, 1)
	c.Axes()
	for x := 0.0; x <= maxFlips; x += 50 {
		c.Tick(x, fmt.Sprintf("%.0f", x))
	}

	// The true probability, as a dashed reference line the running
	// average is converging toward.
	c.Path([][2]float64{{0, trueP}, {maxFlips, trueP}}, viz.Muted, 1)

	// The running-average curve itself: wide, noisy swings at small n,
	// settling toward the flat line as n grows.
	curve := make([][2]float64, maxFlips)
	for i, avg := range series {
		curve[i] = [2]float64{float64(i + 1), avg}
	}
	c.Path(curve, viz.Accent, 2)

	// The current n, highlighted.
	avg := series[n-1]
	px, py := c.X(float64(n)), c.Y(avg)
	c.Path([][2]float64{{float64(n), 0}, {float64(n), avg}}, viz.Warm, 1)
	c.Rect(px-4, py-4, 8, 8, viz.Warm, 1)

	c.Text(16, 24, fmt.Sprintf("true probability p = %.2f", trueP), 14, viz.Muted, "start")
	c.Text(16, 44, fmt.Sprintf("after n=%d flips: running average = %.3f", n, avg), 15, viz.Warm, "start")
	c.Text(16, 64, fmt.Sprintf("error |average - p| = %.3f", math.Abs(avg-trueP)), 13, viz.Ink, "start")

	return c.String()
}
