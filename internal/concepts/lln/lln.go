// Package lln visualizes the law of large numbers: a running average of
// repeated coin flips, jumping around wildly at small n and settling in on
// the true probability as n grows -- without ever guaranteeing that any
// single next flip "corrects" the average.
package lln

import (
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

func render(p map[string]float64) string {
	c := viz.New(680, 420, 0, 300, 0, 1)
	return c.String()
}
