// Package pascalstriangle visualizes Pascal's triangle: each entry built
// as the sum of the two entries above it, and why that same number also
// answers "how many ways are there to choose k items from n" — the
// binomial coefficient "n choose k."
package pascalstriangle

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "pascals-triangle",
		Seq:   34,
		Title: "Pascal's triangle (binomial coefficients, row by row)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Row n", Min: 0, Max: 8, Step: 1, Def: 5},
			{Key: "k", Label: "Position k (choose k of n)", Min: 0, Max: 8, Step: 1, Def: 2},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(600, 440, 0, 1, 0, 1)
	return c.String()
}
