// Package crossvalidation visualizes k-fold cross-validation: instead of
// judging a model on one train/test split, rotate which slice of the data is
// held out, score the model on each slice in turn, and average the scores.
// The picture drives home why that average is a far more honest number than
// whatever a single lucky (or unlucky) split happens to hand you.
package crossvalidation

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "cross-validation",
		Seq:   85,
		Title: "Cross-validation (rotating the held-out fold)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "k", Label: "Number of folds (k)", Min: 2, Max: 10, Step: 1, Def: 5},
			{Key: "highlightFold", Label: "Held-out fold", Min: 0, Max: 9, Step: 1, Def: 3},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(680, 440, 0, 1, 0, 1)
	return c.String()
}
