// Package knapsack visualizes the 0/1 knapsack problem solved by dynamic
// programming: instead of checking every possible subset of items (2^n of
// them), build up a table of the best value achievable for every smaller
// item-count/capacity combination, then reuse those already-solved
// subproblems to answer the full problem in one pass.
package knapsack

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "dynamic-programming-knapsack",
		Seq:   101,
		Title: "0/1 knapsack (dynamic programming)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder -- filled in once the math and picture exist.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "step", Label: "Step (fill the table one cell at a time)", Min: 0, Max: 18, Step: 1, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(700, 460, 0, 1, 0, 1).String()
}
