// Package geometric visualizes the geometric distribution: the probability
// that the first success in a sequence of independent Bernoulli(p) trials
// lands exactly on trial k. binomial-distribution answers "out of a fixed n
// trials, how many succeed"; this concept answers the open-ended question
// "how many trials until the first success happens at all" -- flipping a
// coin until the first heads is the running example.
package geometric

import (
	"fmt"
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

// displayTrials is how many trials the bar chart shows -- fixed regardless
// of p so the picture's shape (a steadily shrinking tail) is comparable
// across slider settings, matching the k slider's own range.
const displayTrials = 25

func render(params map[string]float64) string {
	p := params["p"]
	if p < 0.01 {
		p = 0.01
	}
	if p > 1 {
		p = 1
	}
	k := int(params["k"] + 0.5)
	if k < 1 {
		k = 1
	}
	if k > displayTrials {
		k = displayTrials
	}

	pmf := make([]float64, displayTrials+1) // index 0 unused, trials are 1-based
	maxP := 0.0
	for i := 1; i <= displayTrials; i++ {
		pmf[i] = PMF(p, i)
		if pmf[i] > maxP {
			maxP = pmf[i]
		}
	}
	yMax := maxP * 1.2
	if yMax <= 0 {
		yMax = 1
	}

	c := viz.New(680, 420, 0.5, float64(displayTrials)+0.5, 0, yMax)
	c.Axes()
	for x := 1.0; x <= displayTrials; x += 5 {
		c.Tick(x, fmt.Sprintf("%.0f", x))
	}

	mean := Mean(p)
	if mean <= float64(displayTrials) {
		c.Path([][2]float64{{mean, 0}, {mean, yMax}}, viz.Muted, 1)
	}

	barHalfW := 0.4
	for i := 1; i <= displayTrials; i++ {
		color := viz.Accent
		if i == k {
			color = viz.Warm
		}
		x0, x1 := c.X(float64(i)-barHalfW), c.X(float64(i)+barHalfW)
		y0, y1 := c.Y(0), c.Y(pmf[i])
		c.Rect(x0, y1, x1-x0, y0-y1, color, 0.8)
	}

	c.Text(16, 24, fmt.Sprintf("p=%.2f    mean trials until first success=1/p=%.2f    variance=%.2f",
		p, mean, Variance(p)), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("P(X=%d)=%.4f    P(X<=%d)=%.4f    P(X>%d)=%.4f",
		k, PMF(p, k), k, CDF(p, k), k, 1-CDF(p, k)), 15, viz.Warm, "start")

	return c.String()
}
