// Package bandit visualizes the multi-armed bandit problem: repeatedly
// choosing among several options with unknown reward rates, where every pull
// spent finding out which option is best is a pull that couldn't go to the
// option you already suspect is best. Epsilon-greedy is the simplest
// strategy that faces that trade-off head-on: pull a uniformly random arm
// with probability epsilon (explore), otherwise pull whichever arm currently
// looks best (exploit).
package bandit

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "multi-armed-bandit",
		Seq:   102,
		Title: "Multi-armed bandit (explore vs. exploit)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder -- filled in once the math and picture exist.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "epsilon", Label: "Exploration rate (ε)", Min: 0, Max: 1, Step: 0.05, Def: 0.2},
			{Key: "step", Label: "Pulls so far (step)", Min: 0, Max: 60, Step: 1, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(700, 460, 0, 1, 0, 1).String()
}
