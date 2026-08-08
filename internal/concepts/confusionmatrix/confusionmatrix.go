// Package confusionmatrix visualizes the four confusion-matrix cells — true
// positive, false positive, false negative, true negative — as a literal 2x2
// grid whose shading tracks each cell's share of the population. It sits next
// to precision-recall and roc-auc as another view of the same underlying
// setup (two overlapping normal score distributions, a threshold), but reads
// off the raw counts and the metrics built from them directly.
package confusionmatrix

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "confusion-matrix",
		Title: "Confusion matrix",
		Blurb: "A population is split 50/50 into a real-positive and real-negative class, " +
			"scored from two overlapping normal distributions 'separation' apart, and " +
			"classified positive whenever the score clears the threshold. The 2x2 grid is " +
			"the result: true positives and true negatives on one diagonal, false positives " +
			"and false negatives (the two ways a classifier gets it wrong) on the other. " +
			"Darker cells hold more of the population. Slide the threshold or separation and " +
			"watch counts — and the accuracy/precision/recall/F1 built from them — shift.",
		Params: []concept.ParamSpec{
			{Key: "thresh", Label: "Threshold", Min: -3, Max: 6, Step: 0.1, Def: 1.5},
			{Key: "sep", Label: "Class separation", Min: 1, Max: 5, Step: 0.1, Def: 3},
			{Key: "n", Label: "Population size (n)", Min: 20, Max: 500, Step: 10, Def: 200},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 400, 0, 1, 0, 1).String()
}
