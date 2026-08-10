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
		Seq:   17,
		Title: "Precision-Recall curve & PR-AUC",
		Blurb: "Same spam filter as the precision-recall lesson, but instead of picking one " +
			"threshold, try three, strictest to loosest, and watch (recall, precision) move: " +
			"Flag only the 8 most obvious spam emails: all 8 really are spam, so recall = " +
			"8/20 = 40%, precision = 8/8 = 100%. Loosen it to flag 22: 18 are real spam, 4 are " +
			"false alarms, so recall = 18/20 = 90%, precision = 18/22 = 82%. Loosen it all the " +
			"way to flag 200, catching every last spam email: recall = 20/20 = 100%, but " +
			"precision craters to 20/200 = 10% — 9 of every 10 flagged emails are now false " +
			"alarms. Plot each (recall, precision) pair as a point and connect them: that's the " +
			"PR curve. Sweep EVERY threshold instead of just these 3 and the curve becomes " +
			"continuous; the area under it (PR-AUC) grades the classifier across every " +
			"threshold at once, the same way ROC-AUC does for the ROC curve. The difference " +
			"matters when positives are rare: ROC-AUC can look great even when precision is " +
			"terrible, because it's diluted by a huge pool of easy true negatives. PR-AUC has " +
			"nowhere to hide from false positives piling up, so it's the sharper read when the " +
			"classes are imbalanced.",
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
	// The sweep has to run comfortably above the positive class's mean
	// (sep) so recall actually reaches ~0 at the strict end. A fixed
	// window that ignored sep left recall stuck around 0.16 (not 0) at the
	// top of the allowed separation range, silently truncating the
	// high-precision tip of the curve and undercounting PR-AUC for exactly
	// the classifiers that should score best. tLow=-6 is already far below
	// every allowed sep (negative class is fixed at 0), so only the high
	// end needs to track sep.
	tHigh := sep + 6
	const tLow = -6.0
	pts := make([][2]float64, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := tHigh - (tHigh-tLow)*float64(i)/float64(steps)
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

	c := viz.New(680, 420, 0, 1, 0, 1)
	// Extra headroom for three lines of header text — with good separation
	// the curve itself hugs y=1 (precision near 100%) across most of the
	// plot, which would otherwise run straight through the caption text.
	c.PadT = 70
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
	c.Text(20, 62, "x = recall, y = precision    grey flat line = a classifier that ranks "+
		"randomly (always lands right ~half the time)", 12, viz.Muted, "start")

	return c.String()
}
