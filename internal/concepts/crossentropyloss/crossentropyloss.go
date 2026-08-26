// Package crossentropyloss visualizes binary cross-entropy (log) loss: the
// per-example penalty a classifier pays for how confidently right or wrong
// its predicted probability was against a true 0/1 label, and why averaging
// that penalty over a dataset is the loss function most classifiers are
// actually trained to minimize.
package crossentropyloss

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "cross-entropy-loss",
		Seq:   69,
		Title: "Cross-entropy loss (scoring a predicted probability against a label)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`kl-divergence` compared two full probability distributions — a true " +
						"coin bias P(heads)=0.9 against a model's believed Q(heads)=0.5 — and " +
						"measured the gap between them in bits. But training a real classifier " +
						"doesn't hand you a true distribution to compare against; it hands you a " +
						"hard label. A spam filter either was spam (y=1) or wasn't (y=0) for a " +
						"given email — there's no 'this email was 90% spam' ground truth. The " +
						"model still outputs a probability, say p̂=0.1 for an email that actually " +
						"was spam. How do you turn 'the true label was 1 but the model said 0.1' " +
						"into a single number a training algorithm can push toward zero?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Four emails, each with a true spam label y and the model's predicted " +
						"P(spam)=p̂:",
					"• Email 1: y=1 (spam), p̂=0.9 (confident and correct) → loss = " +
						"-log2(0.9) = 0.152 bits",
					"• Email 2: y=1 (spam), p̂=0.55 (barely correct) → loss = -log2(0.55) = " +
						"0.863 bits",
					"• Email 3: y=0 (not spam), p̂=0.2 (confident and correct, since the model " +
						"put 80% on 'not spam') → loss = -log2(1-0.2) = -log2(0.8) = 0.322 bits",
					"• Email 4: y=1 (spam), p̂=0.1 (confident and WRONG) → loss = -log2(0.1) = " +
						"3.322 bits",
					"Each loss only ever looks at the probability the model assigned to the " +
						"label that actually happened: -log2(p̂) when y=1, -log2(1-p̂) when y=0. " +
						"Average the four: (0.152+0.863+0.322+3.322)/4 = 1.165 bits — the average " +
						"cross-entropy loss over this tiny dataset. Notice email 4 alone " +
						"contributes almost three times as much as the other three combined: " +
						"-log2(x) shoots toward infinity as x→0, so a confident wrong prediction " +
						"is punished far harder than an unconfident right one is rewarded.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Two curves of loss = -log2(p̂) (when y=1) and -log2(1-p̂) (when y=0), " +
						"plotted across every possible predicted probability p̂ from 0 to 1. The y " +
						"slider picks which true label is in force — which curve is 'active' — " +
						"and the p̂ slider marks a point on it. Drag p̂ toward the label's own side " +
						"(1 when y=1, 0 when y=0) and watch the loss fall toward 0; drag it toward " +
						"the wrong side and watch the loss climb steeply, blowing up as p̂ " +
						"approaches the wrong extreme.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Score any single prediction against its true label with one number that " +
						"rewards confident-and-correct, tolerates unconfident-and-correct, and " +
						"punishes confident-and-wrong sharply — exactly the ranking `kl-divergence` " +
						"couldn't produce without a full second distribution to compare against. " +
						"Average that number over a whole training set and you get the actual " +
						"quantity gradient descent minimizes when training a classifier: nudge " +
						"every prediction a little closer to its true label, over and over, and the " +
						"average loss falls.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"'Cross-entropy loss' (also called 'log loss') is the default training " +
						"objective for logistic regression, neural network classifiers, and spam/ " +
						"fraud detectors — anywhere a model outputs a probability and gets scored " +
						"against a known right answer. Kaggle competitions and production ML " +
						"dashboards alike report 'log loss' as a leaderboard metric for exactly " +
						"this reason: it's differentiable, and it punishes confidently-wrong " +
						"predictions in a way plain accuracy doesn't.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'email 4's loss of 3.322 bits dominates the average " +
						"because the model was both wrong and confident about it — cross-entropy " +
						"loss punishes that combination far harder than being merely unconfident.'",
					"Not like this: treating cross-entropy loss as the same thing as accuracy, " +
						"or assuming a lower average loss always means more correct predictions. A " +
						"model that's right 90% of the time but wildly overconfident on its misses " +
						"can post a worse average loss than a model that's right 85% of the time " +
						"but stays modest (p̂ near 0.5) whenever it's unsure — loss cares about " +
						"calibrated confidence, not just which side of 0.5 the prediction landed " +
						"on.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "y", Label: "True label y", Min: 0, Max: 1, Step: 1, Def: 1},
			{Key: "phat", Label: "Predicted P(y=1)", Min: 0.01, Max: 0.99, Step: 0.01, Def: 0.7},
		},
		Render: render,
	})
}

// Labels and Preds are the small worked dataset used throughout this
// concept: 4 emails' true spam labels (1=spam, 0=not spam) and a model's
// predicted P(spam), index-aligned.
var (
	Labels = []float64{1, 1, 0, 1}
	Preds  = []float64{0.9, 0.55, 0.2, 0.1}
)

// log2Safe returns -log2(x), treating x<=0 as +Inf (an infinitely surprising,
// infinitely costly event) instead of NaN, so Loss stays well-defined even at
// the extreme ends of a probability slider.
func log2Safe(x float64) float64 {
	if x <= 0 {
		return math.Inf(1)
	}
	return -math.Log2(x)
}

// Loss returns the binary cross-entropy ("log") loss, in bits, of predicting
// P(y=1)=phat against a true hard label y (0 or 1):
//
//	Loss(1, phat) = -log2(phat)
//	Loss(0, phat) = -log2(1-phat)
//
// It's 0 only in the limit of a perfectly confident, perfectly correct
// prediction, and grows without bound as phat confidently disagrees with y —
// a single wrong-and-confident prediction can outweigh many merely-unsure
// correct ones (see the worked example above).
func Loss(y, phat float64) float64 {
	if y >= 0.5 {
		return log2Safe(phat)
	}
	return log2Safe(1 - phat)
}

// AverageLoss returns the mean Loss across every (label, prediction) pair —
// the actual quantity a classifier's training loop tries to drive toward
// zero. labels and preds must be the same length and index-aligned; returns
// 0 for an empty dataset.
func AverageLoss(labels, preds []float64) float64 {
	if len(labels) == 0 {
		return 0
	}
	sum := 0.0
	for i, y := range labels {
		sum += Loss(y, preds[i])
	}
	return sum / float64(len(labels))
}

// lossCap bounds how tall the loss curves can climb -- Loss shoots toward
// infinity as a prediction confidently disagrees with the label, so capping
// the curve keeps the picture readable instead of the y-axis needing to
// stretch to infinity to fit it (same idea as kldivergence's klCap).
const lossCap = 6.0

func clampCap(v float64) float64 {
	if v > lossCap {
		return lossCap
	}
	return v
}

func render(p map[string]float64) string {
	y, phat := p["y"], p["phat"]
	active := y >= 0.5

	c := viz.New(680, 440, 0, 1, 0, lossCap)
	c.PadT = 90
	c.PadB = 90
	c.Axes()
	for x := 0.0; x <= 1.0; x += 0.25 {
		c.Tick(x, fmt.Sprintf("%.2g", x))
	}

	y1Curve := viz.Sample(0.005, 0.995, 200, func(x float64) float64 { return clampCap(Loss(1, x)) })
	y0Curve := viz.Sample(0.005, 0.995, 200, func(x float64) float64 { return clampCap(Loss(0, x)) })

	// Draw the inactive curve first, muted, so the active one visibly sits on
	// top of it in the accent color.
	if active {
		c.Path(y0Curve, viz.Muted, 1.5)
		c.Path(y1Curve, viz.Accent, 2.5)
	} else {
		c.Path(y1Curve, viz.Muted, 1.5)
		c.Path(y0Curve, viz.Accent, 2.5)
	}

	c.VLine(phat, viz.Warm, true)
	loss := Loss(y, phat)
	mx, my := c.X(phat), c.Y(clampCap(loss))
	c.Rect(mx-4, my-4, 8, 8, viz.Warm, 0.9)

	// The 4 worked-example emails, plotted on whichever curve their own
	// label belongs to, so the picture ties directly back to the numbers in
	// "How does it actually work?".
	for i, lbl := range Labels {
		ex, ey := Preds[i], clampCap(Loss(lbl, Preds[i]))
		px, py := c.X(ex), c.Y(ey)
		color := viz.Accent
		if lbl < 0.5 {
			color = viz.Good
		}
		c.Rect(px-3, py-3, 6, 6, color, 0.85)
		c.Text(px+6, py-6, fmt.Sprintf("e%d", i+1), 11, viz.Muted, "start")
	}

	label := "y=0 (curve: -log2(1-p̂))"
	if active {
		label = "y=1 (curve: -log2(p̂))"
	}
	lossStr := fmt.Sprintf("%.3f", loss)
	if math.IsInf(loss, 1) {
		lossStr = "+Inf"
	}
	c.Text(20, 24, fmt.Sprintf("true label %s   p̂=%.2f   loss=%s bits", label, phat, lossStr), 14, viz.Ink, "start")
	c.Text(20, 44, "accent curve = active label's loss; muted curve = the other label, for comparison", 12, viz.Muted, "start")
	c.Text(20, 64, fmt.Sprintf("worked example (e1-e4): average loss over the 4 emails = %.3f bits", AverageLoss(Labels, Preds)), 12, viz.Muted, "start")

	return c.String()
}
