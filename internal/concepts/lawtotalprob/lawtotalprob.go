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

func render(p map[string]float64) string {
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
