// Package lawtotalprob visualizes the law of total probability: splitting a
// hard-to-compute overall probability into a probability-weighted
// combination of easier per-scenario pieces. The running example is a
// factory floor with three machines, each producing a different share of
// total output at a different defect rate -- what fraction of ALL widgets
// coming off the floor are defective?
package lawtotalprob

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "law-of-total-probability",
		Seq:   91,
		Title: "Law of total probability (combining per-scenario pieces)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "shareA", Label: "Machine A's share of output", Min: 0.5, Max: 0.96, Step: 0.01, Def: 0.90},
			{Key: "shareB", Label: "Machine B's share of output", Min: 0.02, Max: 0.40, Step: 0.01, Def: 0.08},
			{Key: "rateB", Label: "Machine B's defect rate", Min: 0.02, Max: 0.50, Step: 0.01, Def: 0.20},
		},
		Render: render,
	})
}

// Shares normalizes the two free machine shares (A and B) into a valid
// three-way partition (a, b, c) that sums to exactly 1, with the remainder
// going to machine C. If shareA+shareB would leave machine C with less
// than 1% of output, both are scaled down proportionally so c stays at
// least 0.01 -- every scenario in the partition keeps some nonzero weight.
func Shares(shareA, shareB float64) (a, b, c float64) {
	if shareA < 0 {
		shareA = 0
	}
	if shareB < 0 {
		shareB = 0
	}
	if shareA+shareB > 0.99 {
		scale := 0.99 / (shareA + shareB)
		shareA *= scale
		shareB *= scale
	}
	return shareA, shareB, 1 - shareA - shareB
}

// TotalProbability implements the law of total probability itself:
// P(event) = Sum_i P(event|scenario_i) * P(scenario_i). Given each
// scenario's share of the population (shares, which should sum to 1) and
// the event's probability within that scenario (rates), it returns the
// event's overall probability across the whole population. Mismatched
// slice lengths or empty input return 0.
func TotalProbability(shares, rates []float64) float64 {
	sum := 0.0
	n := len(shares)
	if len(rates) < n {
		n = len(rates)
	}
	for i := 0; i < n; i++ {
		sum += shares[i] * rates[i]
	}
	return sum
}

// NaiveAverage returns the plain (unweighted) arithmetic mean of the
// per-scenario rates, ignoring how common each scenario actually is. It
// exists to demonstrate the common mistake of treating every scenario as
// equally likely instead of weighting by its real share -- see the
// "What's the common mistake here?" section.
func NaiveAverage(rates []float64) float64 {
	if len(rates) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range rates {
		sum += r
	}
	return sum / float64(len(rates))
}

func render(p map[string]float64) string {
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
