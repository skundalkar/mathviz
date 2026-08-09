// Package overfitting shows what "the model fits the training data too well"
// actually looks like: a fixed handful of noisy points, and a polynomial
// curve fit to them. Crank the model's degree up and the curve stops
// approximating the underlying trend and starts weaving through every single
// point instead — perfect on the data it was trained on, wild everywhere in
// between.
package overfitting

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "overfitting",
		Title: "Overfitting",
		Blurb: "Give a student twelve practice problems and their answers, then ask them to " +
			"explain the pattern. One student writes down a short, general rule that gets most " +
			"of the twelve right and should generalize to problem thirteen. Another memorizes " +
			"all twelve exact answers, including the quirks and typos in the answer key — zero " +
			"mistakes on practice, but no idea what to do with a new problem. That second " +
			"student has overfit. Here a handful of noisy data points are the practice " +
			"problems, and a polynomial curve is the 'rule' fit to them. A low-degree curve " +
			"stays smooth and roughly tracks the true pattern (dashed). Push the degree up and " +
			"training error keeps dropping — the curve now passes near every dot — but it does " +
			"so by wildly overshooting between them, chasing noise instead of signal.",
		Params: []concept.ParamSpec{
			{Key: "degree", Label: "Model complexity (degree)", Min: 1, Max: 11, Step: 1, Def: 3},
			{Key: "noise", Label: "Noise", Min: 0, Max: 1, Step: 0.05, Def: 0.4},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 340, -1, 1, -1, 1).String()
}
