// Package rocauc visualizes an ROC curve: sweep a classifier's decision
// threshold from "call nothing positive" to "call everything positive" and
// trace out (false-positive rate, true-positive rate) at every setting. A
// classifier that separates its two classes well bows the curve toward the
// top-left corner; a coin-flip classifier can't do better than the diagonal.
// The area under that curve (AUC) summarizes how good the ranking is across
// every possible threshold at once, not just the one you happen to pick.
package rocauc

import (
	"fmt"
	"math"

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

// tailAbove is P(X > t) for X ~ N(mu, 1): the fraction of a class's scores a
// threshold t calls positive.
func tailAbove(t, mu float64) float64 {
	return 0.5 * math.Erfc((t-mu)/math.Sqrt2)
}

// TPR is the true-positive rate (recall) at threshold t: the fraction of the
// positive class (scores ~ N(sep, 1)) that clears the threshold.
func TPR(t, sep float64) float64 {
	return tailAbove(t, sep)
}

// FPR is the false-positive rate at threshold t: the fraction of the negative
// class (scores ~ N(0, 1)) that clears the threshold too.
func FPR(t float64) float64 {
	return tailAbove(t, 0)
}

// CurvePoints sweeps the decision threshold from "call nothing positive" down
// to "call everything positive" and returns the resulting (FPR, TPR) points,
// ordered from (≈0, ≈0) to (≈1, ≈1) — the ROC curve.
func CurvePoints(sep float64, steps int) [][2]float64 {
	if steps < 2 {
		steps = 2
	}
	pts := make([][2]float64, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := 6 - 12*float64(i)/float64(steps) // threshold sweeps +6 down to -6
		pts = append(pts, [2]float64{FPR(t), TPR(t, sep)})
	}
	return pts
}

// TrapezoidalAUC numerically integrates TPR over FPR along an ROC curve using
// the trapezoid rule. Used to cross-check the closed-form AUC formula.
func TrapezoidalAUC(pts [][2]float64) float64 {
	var area float64
	for i := 1; i < len(pts); i++ {
		x0, y0 := pts[i-1][0], pts[i-1][1]
		x1, y1 := pts[i][0], pts[i][1]
		area += (x1 - x0) * (y0 + y1) / 2
	}
	return area
}

// AUC is the closed-form area under the ROC curve for two equal-variance
// normal classes separated by sep: the probability that a random positive
// example scores higher than a random negative one, P(X>Y) for X~N(sep,1),
// Y~N(0,1), which works out to Φ(sep/√2).
func AUC(sep float64) float64 {
	return normCDF(sep / math.Sqrt2)
}

// normCDF is the standard normal cumulative distribution function Φ(z).
func normCDF(z float64) float64 {
	return 0.5 * (1 + math.Erf(z/math.Sqrt2))
}

func render(p map[string]float64) string {
	t, sep := p["thresh"], p["sep"]
	if sep < 0.01 {
		sep = 0.01
	}

	pts := CurvePoints(sep, 300)
	auc := AUC(sep)
	numericAUC := TrapezoidalAUC(pts)

	c := viz.New(680, 380, 0, 1, 0, 1)
	c.Axes()
	for x := 0.0; x <= 1.0; x += 0.2 {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	// The area under the curve IS the AUC — shade it directly.
	c.Area(pts, 0, 1, viz.Accent, 0.15)
	// The diagonal is what a coin-flip classifier achieves; a real classifier
	// should bow up and to the left of it.
	c.Path([][2]float64{{0, 0}, {1, 1}}, viz.Muted, 1.5)
	c.Path(pts, viz.Ink, 2.5)

	// Mark where the current threshold sits on the curve.
	curFPR, curTPR := FPR(t), TPR(t, sep)
	c.VLine(curFPR, viz.Warm, true)
	mx, my := c.X(curFPR), c.Y(curTPR)
	c.Rect(mx-4, my-4, 8, 8, viz.Warm, 1)

	c.Text(20, 24, fmt.Sprintf("AUC = %.3f (numerically ≈ %.3f)    separation = %.1f", auc, numericAUC, sep),
		14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("at threshold t=%.1f: FPR = %.2f, TPR (recall) = %.2f", t, curFPR, curTPR),
		13, viz.Muted, "start")
	c.Text(20, 62, "x = false-positive rate, y = true-positive rate   diagonal = random guessing",
		12, viz.Muted, "start")

	return c.String()
}
