// Package geometric visualizes the geometric distribution: the probability
// that the first success in a sequence of independent Bernoulli(p) trials
// lands exactly on trial k. binomial-distribution answers "out of a fixed n
// trials, how many succeed"; this concept answers the open-ended question
// "how many trials until the first success happens at all" -- flipping a
// coin until the first heads is the running example.
package geometric

import (
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

func render(p map[string]float64) string {
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
