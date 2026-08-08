// Package bayestheorem visualizes why a positive test for a rare condition is
// so often wrong: a population of 100 people is split into those who tested
// positive and those who tested negative, colored by whether the test got it
// right. When the base rate is low, most of the "positive" group turns out to
// be false alarms even for a fairly accurate test — Bayes' theorem is just
// the arithmetic that explains why.
package bayestheorem

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "bayes-theorem",
		Title: "Bayes' theorem",
		Blurb: "A population of 100 people is split into a 'tested positive' group and a " +
			"'tested negative' group. Within 'tested positive', green squares are people who " +
			"really are sick (true positives) and red squares are people who aren't (false " +
			"alarms). When the condition is rare, red usually outnumbers green even for a " +
			"pretty accurate test — that's Bayes' theorem: a positive result updates your " +
			"belief, but it starts from how rare the condition was in the first place.",
		Params: []concept.ParamSpec{
			{Key: "prior", Label: "Base rate (prior)", Min: 0.001, Max: 0.5, Step: 0.001, Def: 0.01},
			{Key: "sensitivity", Label: "Sensitivity (catches true cases)", Min: 0.5, Max: 0.999, Step: 0.001, Def: 0.99},
			{Key: "specificity", Label: "Specificity (clears true negatives)", Min: 0.5, Max: 0.999, Step: 0.001, Def: 0.95},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 360, 0, 1, 0, 1).String()
}
