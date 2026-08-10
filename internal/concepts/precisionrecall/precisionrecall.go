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
		Seq:   3,
		Title: "Precision vs. recall",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You're building a spam filter. The obvious plan: make it aggressive, flag " +
						"anything suspicious, catch every last piece of spam. Turn that dial all " +
						"the way up and you do catch 100% of the spam — recall is perfect — but you " +
						"also start flagging your boss's emails, meeting invites, and " +
						"password-reset links, because 'suspicious' was cast too wide.",
					"Loosen the dial to stop burying real mail and now actual spam slips through " +
						"the net untouched. There's no setting of this one dial that avoids both " +
						"problems — that's not a bug in your filter, it's structural.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"The picture makes the structure visible: real negatives (legitimate email) " +
						"and real positives (spam) are two overlapping bell curves, and 'flag as " +
						"spam' means 'score above the threshold.' Recall asks, of everything that " +
						"really was spam, how much did you catch. Precision asks, of everything you " +
						"flagged, how much was actually spam.",
					"100 emails arrive, 20 really are spam, 80 are legit. At one threshold the " +
						"filter flags 25 emails: 18 truly spam (caught), 7 legit mail wrongly " +
						"flagged, 2 real spam missed. Recall = 18/20 = 90%, precision = 18/25 = 72%.",
					"• Tighten the threshold to flag only 19 emails: 17 caught, just 2 false " +
						"alarms, but now 3 real spam emails slip through instead of 2. Recall drops " +
						"to 17/20 = 85%, precision climbs to 17/19 ≈ 89% — precision up, recall " +
						"down, from the exact same dial.",
					"The one thing that genuinely improves both at once is separating the two " +
						"curves further apart — a better model that scores real spam and real mail " +
						"more differently in the first place.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Everything right of the threshold line counts as 'predicted positive.' " +
						"Slide the threshold right and precision climbs while recall falls — the " +
						"same move as tightening from flagging 25 emails down to 19 above. Slide it " +
						"left and the opposite happens. The separation slider is that 'better " +
						"model': push the two curves further apart and both precision and recall " +
						"rise together, instead of just trading against each other.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can now tell which lever actually fixes a disappointing number: " +
						"nudging the threshold trades precision for recall instantly with no " +
						"retraining, while only separating the classes further — a genuinely " +
						"better model — improves both at once. Section 1's 'no setting avoids " +
						"both problems' is exactly the sign you've maxed out the threshold lever " +
						"and need to reach for the other one.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Cancer screening (missing a real case is worse than a false alarm, so it's " +
						"worth tolerating more false positives), spam filters (a false positive " +
						"means a lost important email), fraud detection, and search ranking — " +
						"anywhere 'how sensitive should this be' is really a business decision " +
						"about which mistake costs more, dressed up as a slider.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"'We need higher recall here, even if precision drops' — for something " +
						"like cancer screening, missing a real case is worse than a false alarm, " +
						"so the tradeoff should be set on purpose, not left to a model's default.",
					"'Just make the model more accurate' misses the point — accuracy alone " +
						"hides which kind of mistake it's making; the actual decision is which " +
						"type of error costs more in your specific situation.",
				},
			},
		},
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
