// Package prauc visualizes the precision-recall curve: sweep a classifier's
// decision threshold from "call nothing positive" to "call everything
// positive" and trace out (recall, precision) at every setting. Unlike ROC,
// which stays optimistic when negatives vastly outnumber positives, the PR
// curve reacts directly to false positives piling up against a small
// positive class — the area under it (PR-AUC) is the metric to reach for on
// imbalanced problems.
package prauc

import (
	"fmt"
	"math"
	"sort"

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

// tailAbove is P(X > t) for X ~ N(mu, 1): the fraction of a class's scores a
// threshold t calls positive.
func tailAbove(t, mu float64) float64 {
	return 0.5 * math.Erfc((t-mu)/math.Sqrt2)
}

// Recall is the fraction of the positive class (scores ~ N(sep, 1)) that
// clears threshold t. Equivalent to TPR in the ROC picture.
func Recall(t, sep float64) float64 {
	return tailAbove(t, sep)
}

// Precision is the fraction of everything flagged positive (score > t, drawn
// from either class) that is truly positive. By convention, when nothing is
// flagged (both classes' tails are ~0) precision is defined as 1 — "we made
// no mistakes because we made no calls" — matching the usual scikit-learn
// convention at recall=0.
func Precision(t, sep float64) float64 {
	tp := tailAbove(t, sep)
	fp := tailAbove(t, 0)
	if tp+fp <= 0 {
		return 1
	}
	return tp / (tp + fp)
}

// CurvePoints sweeps the decision threshold from "call everything positive"
// up to "call nothing positive" and returns the resulting (recall,
// precision) points, ordered by increasing recall from ≈0 to ≈1 — the PR
// curve.
func CurvePoints(sep float64, steps int) [][2]float64 {
	if steps < 2 {
		steps = 2
	}
	pts := make([][2]float64, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := 6 - 12*float64(i)/float64(steps) // threshold sweeps +6 down to -6
		pts = append(pts, [2]float64{Recall(t, sep), Precision(t, sep)})
	}
	sort.Slice(pts, func(i, j int) bool { return pts[i][0] < pts[j][0] })
	return pts
}

// TrapezoidalPRAUC numerically integrates precision over recall along a PR
// curve using the trapezoid rule. This is a close cousin of "average
// precision" and, like it, can be a touch optimistic on a jagged curve — but
// for the smooth curves two Gaussian classes produce here it tracks the true
// area well and is simple to compute directly from the sampled points.
func TrapezoidalPRAUC(pts [][2]float64) float64 {
	var area float64
	for i := 1; i < len(pts); i++ {
		x0, y0 := pts[i-1][0], pts[i-1][1]
		x1, y1 := pts[i][0], pts[i][1]
		area += (x1 - x0) * (y0 + y1) / 2
	}
	return area
}

func render(p map[string]float64) string {
	t, sep := p["thresh"], p["sep"]
	if sep < 0.01 {
		sep = 0.01
	}

	pts := CurvePoints(sep, 300)
	auc := TrapezoidalPRAUC(pts)

	c := viz.New(680, 380, 0, 1, 0, 1)
	c.Axes()
	for x := 0.0; x <= 1.0; x += 0.2 {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	// The area under the curve IS the PR-AUC — shade it directly.
	c.Area(pts, 0, 1, viz.Accent, 0.15)
	// A classifier that ranks randomly still lands every flagged item right
	// about `prevalence` of the time regardless of how many it flags — a flat
	// line at the positive class's share (0.5 with our equal-sized classes),
	// unlike the ROC diagonal.
	c.Path([][2]float64{{0, 0.5}, {1, 0.5}}, viz.Muted, 1.5)
	c.Path(pts, viz.Ink, 2.5)

	// Mark where the current threshold sits on the curve.
	curRec, curPrec := Recall(t, sep), Precision(t, sep)
	c.VLine(curRec, viz.Warm, true)
	mx, my := c.X(curRec), c.Y(curPrec)
	c.Rect(mx-4, my-4, 8, 8, viz.Warm, 1)

	c.Text(20, 24, fmt.Sprintf("PR-AUC ≈ %.3f    separation = %.1f", auc, sep),
		14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("at threshold t=%.1f: recall = %.2f, precision = %.2f", t, curRec, curPrec),
		13, viz.Muted, "start")
	c.Text(20, 62, "x = recall, y = precision   flat line = a classifier that ranks randomly",
		12, viz.Muted, "start")

	return c.String()
}
