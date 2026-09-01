// Package backprop visualizes backpropagation: running the chain rule
// backward through a small network — one input, two hidden neurons, one
// output — to compute every weight's gradient from a single forward pass
// followed by a single backward pass, instead of re-deriving each gradient
// from scratch.
package backprop

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "backpropagation",
		Seq:   87,
		Title: "Backpropagation (the chain rule through a small network)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "x", Label: "Input (x)", Min: -2, Max: 2, Step: 0.1, Def: 1},
			{Key: "lr", Label: "Learning rate", Min: 0, Max: 1.5, Step: 0.1, Def: 0.5},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(700, 560, 0, 1, 0, 1)
	return c.String()
}
