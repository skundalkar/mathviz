// Package lln visualizes the law of large numbers: a running average of
// repeated coin flips, jumping around wildly at small n and settling in on
// the true probability as n grows -- without ever guaranteeing that any
// single next flip "corrects" the average.
package lln

import (
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

func render(p map[string]float64) string {
	c := viz.New(680, 420, 0, 300, 0, 1)
	return c.String()
}
