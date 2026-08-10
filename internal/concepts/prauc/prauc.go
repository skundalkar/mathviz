// Package prauc visualizes the precision-recall curve: sweep a classifier's
// decision threshold from "call nothing positive" to "call everything
// positive" and trace out (recall, precision) at every setting. Unlike ROC,
// which stays optimistic when negatives vastly outnumber positives, the PR
// curve reacts directly to false positives piling up against a small
// positive class — the area under it (PR-AUC) is the metric to reach for on
// imbalanced problems.
package prauc

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "pr-auc",
		Title: "Precision-Recall curve & PR-AUC",
		Blurb: "Same spam filter, same two classes, but now we sweep every possible " +
			"threshold instead of eyeballing one. At each threshold, plot (recall, " +
			"precision) as a single point: recall on the x-axis, precision on the y-axis. " +
			"Flag almost nothing and recall sits near 0 while precision is high (whatever you " +
			"do flag is probably right); flag almost everything and recall climbs toward 1 " +
			"while precision collapses toward the positive class's share of the data. Trace " +
			"every threshold in between and you get the PR curve; the area under it (PR-AUC) " +
			"summarizes the classifier across every threshold at once, the same way ROC-AUC " +
			"does for the ROC curve. The difference matters when positives are rare: ROC-AUC " +
			"can look great even when precision is terrible, because it's diluted by a huge " +
			"pool of true negatives. PR-AUC has nowhere to hide from false positives, so it's " +
			"the sharper read when the classes are imbalanced.",
		Params: []concept.ParamSpec{
			{Key: "thresh", Label: "Threshold", Min: -3, Max: 6, Step: 0.1, Def: 1.5},
			{Key: "sep", Label: "Class separation", Min: 1, Max: 5, Step: 0.1, Def: 3},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 380, 0, 1, 0, 1).String()
}
