// Package bonferroni visualizes the Bonferroni correction: when you run many
// hypothesis tests at once, even if every single null hypothesis is
// perfectly true, testing each one at the usual significance level (like
// α=0.05) lets false positives pile up across the whole batch. Dividing α by
// the number of tests, m, before applying it to any one test keeps the
// probability of even one false positive across the whole family bounded by
// the original α.
package bonferroni

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "bonferroni-correction",
		Seq:   104,
		Title: "Bonferroni correction (adjusting for multiple tests)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder -- filled in once the math and picture exist.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "m", Label: "Number of tests (m)", Min: 1, Max: 50, Step: 1, Def: 20},
			{Key: "alpha", Label: "Significance level (α)", Min: 0.01, Max: 0.20, Step: 0.01, Def: 0.05},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(700, 460, 0, 1, 0, 1).String()
}
