// Package geometric visualizes the geometric distribution: the probability
// that the first success in a sequence of independent Bernoulli(p) trials
// lands exactly on trial k. binomial-distribution answers "out of a fixed n
// trials, how many succeed"; this concept answers the open-ended question
// "how many trials until the first success happens at all" -- flipping a
// coin until the first heads is the running example.
package geometric

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "geometric-distribution",
		Seq:   90,
		Title: "Geometric distribution (trials until first success)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "p", Label: "Success probability (p)", Min: 0.05, Max: 0.95, Step: 0.05, Def: 0.3},
			{Key: "k", Label: "Highlighted trial (k)", Min: 1, Max: 25, Step: 1, Def: 3},
		},
		Render: render,
	})
}

// PMF returns P(X=k): the probability the first success lands exactly on
// trial k, for k=1,2,3,... The first k-1 trials must all fail (each with
// probability 1-p) and trial k must succeed (probability p): PMF(p,k) =
// (1-p)^(k-1) * p. Returns 0 for k<1 or p outside (0,1].
func PMF(p float64, k int) float64 {
	if k < 1 || p <= 0 || p > 1 {
		return 0
	}
	return math.Pow(1-p, float64(k-1)) * p
}

// CDF returns P(X<=k): the probability the first success happens by trial k
// (i.e. within the first k trials). Summing the geometric series
// Sum_{i=1..k} (1-p)^(i-1)*p has the closed form 1-(1-p)^k -- "it didn't
// take longer than k trials" is the complement of "every one of the first
// k trials failed".
func CDF(p float64, k int) float64 {
	if k < 1 {
		return 0
	}
	if p <= 0 || p > 1 {
		return 0
	}
	return 1 - math.Pow(1-p, float64(k))
}

// Mean returns the expected number of trials until the first success, 1/p.
// A fair coin (p=0.5) takes 2 flips on average to see the first heads.
func Mean(p float64) float64 {
	if p <= 0 {
		return math.Inf(1)
	}
	return 1 / p
}

// Variance returns the variance of the trial count until first success,
// (1-p)/p^2.
func Variance(p float64) float64 {
	if p <= 0 {
		return math.Inf(1)
	}
	return (1 - p) / (p * p)
}

func render(p map[string]float64) string {
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
