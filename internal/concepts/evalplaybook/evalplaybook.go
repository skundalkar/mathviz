// Package evalplaybook answers the question every other classifier concept
// in this gallery leaves implicit: once you've trained a model, in what
// order do you actually look at the numbers, and what does each pattern of
// numbers tell you to do next? Loss, ROC-AUC, PR-AUC, precision and recall
// each isolate one thing — this concept is about reading them together.
// Unlike the rest of the gallery, this one has no sliders: it's a fixed
// reference table of the handful of patterns that come up constantly, each
// with its own diagnosis and action, not a single formula to explore.
package evalplaybook

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "eval-playbook",
		Seq:   20,
		Title: "Model evaluation playbook",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder.",
				},
			},
		},
		Params: []concept.ParamSpec{},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 400, 0, 1, 0, 1).String()
}
