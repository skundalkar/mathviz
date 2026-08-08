// Package confusionmatrix visualizes the four confusion-matrix cells — true
// positive, false positive, false negative, true negative — as a literal 2x2
// grid whose shading tracks each cell's share of the population. It sits next
// to precision-recall and roc-auc as another view of the same underlying
// setup (two overlapping normal score distributions, a threshold), but reads
// off the raw counts and the metrics built from them directly.
package confusionmatrix

import (
	"math"

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
	_ = p
	return viz.New(680, 400, 0, 1, 0, 1).String()
}
