// Package calibration answers a question the eval-playbook concept doesn't:
// even once a model separates the classes well (good loss, good ROC-AUC,
// good precision/recall), does its output NUMBER mean what it claims? A
// model that says "90% confident" should be right about 90% of the time it
// says that — ranking ability (AUC) says nothing about whether that's true.
// Temperature scaling — the same knob the sigmoid-softmax concept exposes —
// is the standard, cheap fix, and it never touches the ranking at all.
package calibration

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "calibration",
		Seq:   21,
		Title: "Calibration, ECE & temperature",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"Placeholder."},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "sep", Label: "Class separation", Min: 1, Max: 5, Step: 0.1, Def: 3},
			{Key: "temp", Label: "Temperature", Min: 0.3, Max: 3, Step: 0.1, Def: 0.5},
			{Key: "markP", Label: "Inspect confidence", Min: 0.5, Max: 0.99, Step: 0.01, Def: 0.9},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
