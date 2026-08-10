// Package confusionmatrix visualizes the four confusion-matrix cells — true
// positive, false positive, false negative, true negative — as a literal 2x2
// grid whose shading tracks each cell's share of the population. It sits next
// to precision-recall and roc-auc as another view of the same underlying
// setup (two overlapping normal score distributions, a threshold), but reads
// off the raw counts and the metrics built from them directly.
package confusionmatrix

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "confusion-matrix",
		Seq:   13,
		Title: "Confusion matrix",
		Blurb: "A hospital brags its new test is '99% accurate.' Meaningless on its own: a " +
			"do-nothing test that says 'healthy' for all 100 patients (1 actually sick) scores " +
			"TP=0, FN=1, FP=0, TN=99 — accuracy 99%, identical to a genuinely good test, while " +
			"catching zero real cases. Underneath the grid, every example gets a score: real " +
			"negatives cluster around 0, real positives cluster around whatever 'class " +
			"separation' is set to (small = heavily overlapping, hard to tell apart; large = " +
			"barely overlapping, easy). 'Threshold' is just where you draw the line — score " +
			"above it, called positive. Try threshold=1.5, separation=2.1, n=200 (100/100 " +
			"split): 73 of 100 real positives clear 1.5 and get caught (TP=73, FN=27 missed), " +
			"7 of 100 real negatives clear it too by chance (FP=7, TN=93). That's accuracy=83%, " +
			"precision=73/80=91%, recall=73/100=73%, F1≈81%. Don't stop at 'all pretty high' — " +
			"the 91-vs-73 gap says this setting is conservative: trustworthy when it flags " +
			"something, but missing over a quarter of real positives. Whether that's fine " +
			"depends entirely on what you're screening for. Slide threshold or separation and " +
			"watch all four counts — and the metrics built from them — shift.",
		Params: []concept.ParamSpec{
			{Key: "thresh", Label: "Threshold", Min: -3, Max: 6, Step: 0.1, Def: 1.5},
			{Key: "sep", Label: "Class separation", Min: 1, Max: 5, Step: 0.1, Def: 3},
			{Key: "n", Label: "Population size (n)", Min: 20, Max: 500, Step: 10, Def: 200},
		},
		Render: render,
	})
}

// tailAbove is P(X > t) for X ~ N(mu, 1): the fraction of a class's scores a
// threshold t calls positive.
func tailAbove(t, mu float64) float64 {
	return 0.5 * math.Erfc((t-mu)/math.Sqrt2)
}

// Counts splits a population of n examples, half truly positive and half
// truly negative, into the four confusion-matrix cells. Positive-class
// scores are drawn from N(sep, 1), negative-class scores from N(0, 1), and
// anything at or above t is classified positive. Counts are whole examples;
// within each real class the two possible outcomes (caught vs. missed, for
// example) always sum exactly to that class's size, so tp+fn+fp+tn == n.
func Counts(t, sep float64, n int) (tp, fp, fn, tn int) {
	if n < 0 {
		n = 0
	}
	positives := n / 2
	negatives := n - positives

	tp = roundHalfUp(tailAbove(t, sep) * float64(positives))
	fn = positives - tp
	fp = roundHalfUp(tailAbove(t, 0) * float64(negatives))
	tn = negatives - fp
	return
}

func roundHalfUp(x float64) int {
	if x < 0 {
		return 0
	}
	return int(math.Floor(x + 0.5))
}

// Metrics computes the standard classification metrics from confusion counts.
// Any metric whose denominator is zero is reported as 0 rather than NaN.
func Metrics(tp, fp, fn, tn int) (accuracy, precision, recall, f1 float64) {
	total := tp + fp + fn + tn
	if total > 0 {
		accuracy = float64(tp+tn) / float64(total)
	}
	if tp+fp > 0 {
		precision = float64(tp) / float64(tp+fp)
	}
	if tp+fn > 0 {
		recall = float64(tp) / float64(tp+fn)
	}
	if precision+recall > 0 {
		f1 = 2 * precision * recall / (precision + recall)
	}
	return
}

func render(p map[string]float64) string {
	t, sep := p["thresh"], p["sep"]
	n := int(p["n"])
	if n < 2 {
		n = 2
	}

	tp, fp, fn, tn := Counts(t, sep, n)
	acc, prec, rec, f1 := Metrics(tp, fp, fn, tn)
	total := tp + fp + fn + tn

	c := viz.New(680, 400, 0, 1, 0, 1)

	const cellW, cellH, gap = 220.0, 110.0, 14.0
	const originX, originY = 170.0, 76.0

	cell := func(col, row int, label string, count int, color string) {
		x := originX + float64(col)*(cellW+gap)
		y := originY + float64(row)*(cellH+gap)
		frac := 0.0
		if total > 0 {
			frac = float64(count) / float64(total)
		}
		// Darker fill = more of the population landed in this cell.
		c.Rect(x, y, cellW, cellH, color, 0.15+0.75*frac)
		c.Text(x+cellW/2, y+cellH/2-6, label, 16, viz.Ink, "middle")
		c.Text(x+cellW/2, y+cellH/2+16, fmt.Sprintf("%d  (%.0f%%)", count, frac*100), 13, viz.Ink, "middle")
	}

	cell(0, 0, "TP", tp, viz.Good)
	cell(1, 0, "FP", fp, viz.Bad)
	cell(0, 1, "FN", fn, viz.Warm)
	cell(1, 1, "TN", tn, viz.Muted)

	c.Text(originX+cellW/2, originY-16, "predicted positive", 12, viz.Muted, "middle")
	c.Text(originX+cellW+gap+cellW/2, originY-16, "predicted negative", 12, viz.Muted, "middle")
	c.Text(originX-16, originY+cellH/2, "actual positive", 12, viz.Muted, "end")
	c.Text(originX-16, originY+cellH+gap+cellH/2, "actual negative", 12, viz.Muted, "end")

	c.Text(20, 24, fmt.Sprintf("threshold = %.1f    separation = %.1f    n = %d", t, sep, n),
		14, viz.Ink, "start")
	c.Text(20, 44, "rows = ground truth, columns = what the classifier predicted",
		12, viz.Muted, "start")
	c.Text(20, 372, fmt.Sprintf("accuracy = %.0f%%   precision = %.0f%%   recall = %.0f%%   F1 = %.0f%%",
		acc*100, prec*100, rec*100, f1*100), 14, viz.Ink, "start")

	return c.String()
}
