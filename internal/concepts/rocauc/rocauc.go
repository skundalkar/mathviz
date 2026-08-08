// Package rocauc visualizes an ROC curve: sweep a classifier's decision
// threshold from "call nothing positive" to "call everything positive" and
// trace out (false-positive rate, true-positive rate) at every setting. A
// classifier that separates its two classes well bows the curve toward the
// top-left corner; a coin-flip classifier can't do better than the diagonal.
// The area under that curve (AUC) summarizes how good the ranking is across
// every possible threshold at once, not just the one you happen to pick.
package rocauc

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "roc-auc",
		Title: "ROC & AUC",
		Blurb: "Positive- and negative-class scores are two overlapping normal " +
			"distributions, 'separation' apart. Sweeping the decision threshold from high " +
			"to low traces the ROC curve: x is the false-positive rate, y is the " +
			"true-positive rate (recall). The diagonal is what a coin-flip classifier " +
			"achieves; a better classifier bows the curve toward the top-left corner. AUC " +
			"is the area under that curve — the probability a random positive example " +
			"scores higher than a random negative one, across every threshold at once.",
		Params: []concept.ParamSpec{
			{Key: "thresh", Label: "Threshold", Min: -4, Max: 4, Step: 0.1, Def: 0},
			{Key: "sep", Label: "Class separation", Min: 0.5, Max: 5, Step: 0.1, Def: 2.5},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 380, 0, 1, 0, 1).String()
}
