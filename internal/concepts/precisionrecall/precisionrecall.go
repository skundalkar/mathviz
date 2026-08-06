// Package precisionrecall shows the precision/recall trade-off as a picture.
// Two overlapping bell curves are the real negatives (left) and real positives
// (right). A classifier calls everything to the RIGHT of the threshold
// "positive". Slide the threshold: move it right and precision climbs while
// recall falls; move it left and you catch more positives but drag in false
// alarms. Precision and recall pull in opposite directions — that tension is
// the lesson.
package precisionrecall

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "precision-recall",
		Title: "Precision vs. recall",
		Blurb: "Right of the threshold = predicted positive. Recall = of all the real " +
			"positives (right curve), what fraction did we catch? Precision = of everything " +
			"we flagged, what fraction was actually positive? Slide the threshold right and " +
			"precision goes up but recall drops; slide it left and you catch more but let in " +
			"false positives. Increase separation and both improve — that is what a 'better " +
			"model' does. You rarely get both at once from the threshold alone.",
		Params: []concept.ParamSpec{
			{Key: "thresh", Label: "Threshold", Min: -3, Max: 6, Step: 0.1, Def: 1.5},
			{Key: "sep", Label: "Class separation", Min: 1, Max: 5, Step: 0.1, Def: 3},
		},
		Render: render,
	})
}

// normPDF is the unit-variance normal density centered at mu.
func normPDF(x, mu float64) float64 {
	z := x - mu
	return math.Exp(-0.5*z*z) / math.Sqrt(2*math.Pi)
}

// tailAbove returns P(X > t) for X ~ N(mu, 1) = 0.5 * erfc((t-mu)/√2).
func tailAbove(t, mu float64) float64 {
	return 0.5 * math.Erfc((t-mu)/math.Sqrt2)
}

// Metrics computes precision, recall and F1 for a threshold t with the negative
// class at 0 and positive class at `sep`, assuming equal class sizes. Pure math.
func Metrics(t, sep float64) (precision, recall, f1 float64) {
	tp := tailAbove(t, sep) // real positives we flagged
	fp := tailAbove(t, 0)   // real negatives we flagged
	fn := 1 - tp            // real positives we missed

	recall = tp / (tp + fn) // == tp since tp+fn == 1
	if tp+fp > 0 {
		precision = tp / (tp + fp)
	}
	if precision+recall > 0 {
		f1 = 2 * precision * recall / (precision + recall)
	}
	return
}

func render(p map[string]float64) string {
	t, sep := p["thresh"], p["sep"]

	const xmin, xmax = -4.0, 8.0
	peak := normPDF(0, 0)
	c := viz.New(680, 340, xmin, xmax, 0, peak*1.2)
	c.Axes()
	for x := -4.0; x <= 8.0; x += 2 {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	neg := viz.Sample(xmin, xmax, 240, func(x float64) float64 { return normPDF(x, 0) })
	pos := viz.Sample(xmin, xmax, 240, func(x float64) float64 { return normPDF(x, sep) })

	// Shade what the classifier flags positive (everything right of t):
	//   true positives (from the positive curve) in green,
	//   false positives (from the negative curve) in red.
	c.Area(pos, t, xmax, viz.Good, 0.30)
	c.Area(neg, t, xmax, viz.Bad, 0.30)

	c.Path(neg, viz.Muted, 2)
	c.Path(pos, viz.Ink, 2)
	c.VLine(t, viz.Accent, false)

	prec, rec, f1 := Metrics(t, sep)
	c.Text(20, 24, "left = real negatives   right = real positives", 12, viz.Muted, "start")
	c.Text(c.X(t), c.PadT+12, "threshold", 12, viz.Accent, "middle")
	c.Text(20, 300, fmt.Sprintf("precision = %.0f%%", prec*100), 15, viz.Good, "start")
	c.Text(230, 300, fmt.Sprintf("recall = %.0f%%", rec*100), 15, viz.Ink, "start")
	c.Text(420, 300, fmt.Sprintf("F1 = %.0f%%", f1*100), 15, viz.Accent, "start")
	return c.String()
}
