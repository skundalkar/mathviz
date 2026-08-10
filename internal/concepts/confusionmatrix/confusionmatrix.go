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
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"A hospital announces a new test for a rare disease is '99% accurate.' That " +
						"sounds excellent — until you notice that a test which does nothing at " +
						"all, just prints 'healthy' for every single patient without looking at " +
						"anything, would also score 99% accurate, as long as only 1% of patients " +
						"actually have the disease. That do-nothing test catches zero real cases, " +
						"yet posts the exact same headline number as a test that's genuinely good " +
						"at its job.",
					"Accuracy alone can't tell these two situations apart, because it collapses " +
						"four very different outcomes — caught, missed, false alarm, correctly " +
						"cleared — into a single ratio. What would actually distinguish a " +
						"genuinely good test from a useless one posting the identical number?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Every classification lands in exactly one of four buckets: rows are the " +
						"ground truth (actually positive / actually negative), columns are the " +
						"prediction (predicted positive / predicted negative). Work through the " +
						"do-nothing test's own numbers first — 100 patients, 1 actually sick, it " +
						"predicts 'healthy' for all 100:",
					"• TP = 0 — it never predicts positive, so it can't catch the 1 real case.",
					"• FN = 1 — that one real case, missed.",
					"• FP = 0 — it never wrongly flags anyone.",
					"• TN = 99 — every healthy patient correctly cleared.",
					"Accuracy = (TP+TN)÷100 = (0+99)÷100 = 99% — the identical headline number a " +
						"genuinely good test would post, produced here by a test that caught " +
						"literally zero real cases.",
					"Underneath the grid, every example gets a score: real negatives cluster " +
						"around 0, real positives cluster around whatever 'class separation' is " +
						"set to (small = heavily overlapping, hard to tell apart; large = barely " +
						"overlapping, easy). 'Threshold' is just where you draw the line — score " +
						"above it, called positive.",
					"Now a realistic, non-degenerate setting: threshold=1.5, separation=2.1, " +
						"n=200 (100 real positive / 100 real negative). 73 of the 100 real " +
						"positives clear 1.5 and get caught (TP=73, FN=27 missed); 7 of the 100 " +
						"real negatives clear it too by chance (FP=7, false alarms; TN=93). " +
						"That's accuracy=(73+93)÷200=83%, precision=73÷80=91%, " +
						"recall=73÷100=73%, F1≈81%.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The grid IS those four counts: rows are ground truth, columns are what the " +
						"classifier predicted, and each cell is shaded darker the more of the " +
						"population lands there. At threshold=1.5, separation=2.1, n=200 the TP " +
						"cell holds 73 (37% of the population), FP holds 7, FN holds 27, and TN " +
						"holds 93 (47%) — read straight off the same numbers worked through " +
						"above. Slide the threshold and watch predicted-positive shrink or grow " +
						"as fewer or more examples clear the bar; slide separation and watch how " +
						"cleanly the two classes can be told apart in the first place, independent " +
						"of where the threshold sits.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can look at the grid instead of the headline number and catch a " +
						"do-nothing test instantly — its TP cell sits at zero no matter how green " +
						"its TN cell looks. And with a real classifier, you can stop at 'all " +
						"pretty high' or actually read the numbers: the 91-vs-73 gap between " +
						"precision and recall here says this setting is conservative — " +
						"trustworthy when it flags something, but missing over a quarter of real " +
						"positives — and whether that's acceptable depends entirely on what " +
						"you're screening for, not on the numbers alone.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Fraud detection, rare-disease screening, security alerting — anywhere the " +
						"thing you're trying to catch is rare, a model can hit sky-high accuracy " +
						"by mostly predicting the common outcome and still be worthless at the " +
						"one job it exists to do. Reading the actual 2x2 grid, not just the " +
						"headline accuracy, is the only way to tell a genuinely skilled " +
						"classifier from one that's just exploiting an imbalanced dataset.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'That's a false positive' — flagged as a problem, but it " +
						"wasn't one (a spam filter catching a real email, a smoke alarm going off " +
						"from toast). 'We had a false negative' — a real problem, missed " +
						"entirely. These are not interchangeable mistakes; which one you'd rather " +
						"live with depends entirely on what's actually at stake.",
					"Not like this: 'It's 99% accurate, so it's basically fine' — accuracy alone " +
						"can't tell you whether the mistakes it does make are harmless false " +
						"alarms or costly misses, and those two failures are almost never equally " +
						"bad.",
				},
			},
		},
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
