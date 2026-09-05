// Package expectedvalue visualizes expected value: the probability-weighted
// average of a random variable's outcomes. The running example is a $5
// carnival game — flip a weighted coin; heads pays a net $15, tails nets a
// $5 loss — and the question is whether that game is worth playing on
// average, not on any single play.
package expectedvalue

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "expected-value",
		Seq:   89,
		Title: "Expected value (the number a bet centers on)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "p", Label: "Win probability (p)", Min: 0.01, Max: 0.99, Step: 0.01, Def: 0.2},
			{Key: "winAmt", Label: "Net win amount", Min: 0, Max: 50, Step: 1, Def: 15},
			{Key: "loseAmt", Label: "Net lose amount", Min: -50, Max: 0, Step: 1, Def: -5},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(680, 400, 0, 1, 0, 1).String()
}
