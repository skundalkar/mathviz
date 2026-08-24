// Package logreg visualizes logistic regression: fitting an S-shaped curve
// that models P(pass) directly, instead of stretching linear-regression's
// straight line across data whose answer is only ever 0 or 1. Built from
// sigmoid-softmax's squashing function and gradient-descent's step-downhill
// mechanism, now aimed at cross-entropy (log) loss instead of squared error.
package logreg

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "logistic-regression",
		Seq:   64,
		Title: "Logistic regression",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`linear-regression` fits the best straight line through a scatter of " +
						"points — great when the answer is a continuous number like a quiz score. " +
						"But now the outcome is binary: 10 students study for 0.5 to 5 hours, and " +
						"each either passes (1) or fails (0) an exam. Fit a straight line through " +
						"that 0/1 scatter anyway and it happily predicts things like -0.15 or 1.30 " +
						"'probability of passing' for students near either end of the hours range " +
						"— numbers that aren't probabilities at all. Is there a fitted curve that " +
						"stays inside [0,1] everywhere, so its output can actually be read as a " +
						"probability?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"`sigmoid-softmax` already has the squashing function: sigmoid(z) = " +
						"1/(1+e^-z) takes any real number and returns something in (0,1). Logistic " +
						"regression's model is P(pass|hours) = sigmoid(w·hours + b) — the same " +
						"linear combination linear regression used, just squashed through sigmoid " +
						"before it comes out, so it's always a legal probability.",
					"'Best fit' now can't mean smallest squared error (that's still the straight " +
						"-0.15-or-1.30 problem) — it means the w,b that make the model's predicted " +
						"probability as close as possible to what actually happened, measured by " +
						"log-loss: for one student, -log(p) if they passed, -log(1-p) if they " +
						"failed, where p = sigmoid(w·hours+b). A confident-and-right prediction " +
						"(p=0.98 for an actual pass) costs almost nothing (-log(0.98)≈0.02); a " +
						"confident-and-wrong one (p=0.98 for an actual fail) costs a lot " +
						"(-log(0.02)≈3.91) — log-loss punishes confident mistakes far harder than " +
						"cautious ones.",
					"`gradient-descent`'s downhill-stepping mechanism finds the w,b that minimize " +
						"the average log-loss across all 10 students, the same way it found the " +
						"bottom of a bowl-shaped valley elsewhere — just aimed at this loss instead " +
						"of squared error. Starting from w=0, b=0 and taking 10,000 small downhill " +
						"steps on this exact dataset converges to w≈2.59, b≈-8.43, average log-loss " +
						"≈0.251 — versus 0.693 (=-log(0.5)) at the untrained w=0,b=0 starting point, " +
						"where the model has no information and guesses 50/50 for everyone.",
					"Two students break the overall trend: at 3.0 hours the student passed even " +
						"though fewer-hours students failed, and at 3.5 hours the student failed " +
						"even though more-hours students passed. The fitted curve gives them p=0.34 " +
						"and p=0.66 — on the wrong side of the 0.5 threshold both times. That's not " +
						"a bug: the fitted curve traces the overall trend across all 10 students, " +
						"not the two individual flips against it, the same 'smooth trend beats " +
						"chasing every point' idea `overfitting` covers for continuous data.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Ten black squares are the fixed (hours studied, passed) data points. The " +
						"green S-curve is the fitted model from the worked example (w≈2.59, " +
						"b≈-8.43); the blue S-curve is your own guess, built from the slope and " +
						"intercept sliders — drag either one and watch both the curve's shape and " +
						"the 'your log-loss' readout change immediately. The horizontal gray line at " +
						"0.5 is the decision threshold (the same idea `precision-recall` used to " +
						"turn a probability into a yes/no call); the vertical dashed line marks " +
						"where your current curve crosses it — the hours value your model treats as " +
						"the pass/fail cutoff.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Turn a continuous input into a calibrated probability of a yes/no outcome, " +
						"and read off a genuine decision boundary from it — 'below this many study " +
						"hours, the model predicts fail' — instead of eyeballing a threshold on a " +
						"scatter plot or forcing a straight line to do a job it structurally can't: " +
						"stay inside [0,1] everywhere.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Spam filters, credit approval, medical screening ('does this scan show the " +
						"condition'), and churn prediction ('will this customer cancel') are all " +
						"classic logistic regression use cases: one probability of a yes/no event, " +
						"built from a handful of input features. It's also the same math sitting " +
						"underneath one layer of a neural network classifier — a single sigmoid " +
						"neuron over a weighted sum of inputs is exactly this model.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the model outputs a probability of passing; we call it a " +
						"predicted pass once that probability crosses 0.5' — probability and " +
						"decision are two separate steps.",
					"Not like this: fitting a plain straight line to 0/1 data and reading its " +
						"output directly as a probability, or treating logistic regression's output " +
						"as a hard yes/no by itself. The model always outputs a number in (0,1); " +
						"turning it into a decision needs an explicit threshold, and that threshold " +
						"doesn't have to be 0.5 — `precision-recall` already showed why you'd move " +
						"it.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "guessW", Label: "Your slope (w)", Min: -3, Max: 4, Step: 0.05, Def: 0.3},
			{Key: "guessB", Label: "Your intercept (b)", Min: -10, Max: 5, Step: 0.1, Def: -1},
		},
		Render: render,
	})
}

// HoursStudied and Passed are ten fixed (hours studied, passed exam)
// observations -- the worked example every Section above and LEARNINGS.md
// refer to. Passed entries are 0 (failed) or 1 (passed).
var (
	HoursStudied = []float64{0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0}
	Passed       = []float64{0, 0, 0, 0, 0, 1, 0, 1, 1, 1}
)

// Sigmoid squashes any real-valued z into (0,1): 1/(1+e^-z).
func Sigmoid(z float64) float64 {
	return 1 / (1 + math.Exp(-z))
}

// Predict returns the model's predicted probability of the positive class
// at x: sigmoid(w*x + b).
func Predict(w, b, x float64) float64 {
	return Sigmoid(w*x + b)
}

// logLossEps keeps LogLoss finite when a prediction saturates to exactly 0
// or 1 (log(0) is -Inf) -- clamped just far enough that it never changes a
// realistic loss value, only guards the boundary.
const logLossEps = 1e-12

// LogLoss returns the average binary cross-entropy (log) loss of
// Predict(w,b,xs[i]) against labels ys[i] (each 0 or 1): -log(p) for an
// actual positive, -log(1-p) for an actual negative, averaged over every
// example. Confident-and-wrong predictions cost far more than
// confident-and-right ones cost little -- e.g. -log(0.98)≈0.02 vs
// -log(0.02)≈3.91.
func LogLoss(xs, ys []float64, w, b float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	var total float64
	for i := range xs {
		p := Predict(w, b, xs[i])
		if p < logLossEps {
			p = logLossEps
		}
		if p > 1-logLossEps {
			p = 1 - logLossEps
		}
		if ys[i] == 1 {
			total -= math.Log(p)
		} else {
			total -= math.Log(1 - p)
		}
	}
	return total / float64(len(xs))
}

// FitLogisticRegression fits w,b by batch gradient descent on the average
// log-loss, starting from w=0,b=0 and taking a fixed `iters` steps of size
// `lr`. Pure and deterministic -- same inputs always produce the same
// output, no randomness -- so it can stand in for "the best fit" the same
// way linreg's closed-form least-squares line does, just reached by
// gradient-descent's stepping mechanism instead of a formula (logistic
// regression's loss has no closed-form minimizer).
func FitLogisticRegression(xs, ys []float64, iters int, lr float64) (w, b float64) {
	n := float64(len(xs))
	if n == 0 {
		return 0, 0
	}
	for i := 0; i < iters; i++ {
		var gw, gb float64
		for j := range xs {
			err := Predict(w, b, xs[j]) - ys[j]
			gw += err * xs[j]
			gb += err
		}
		w -= lr * gw / n
		b -= lr * gb / n
	}
	return w, b
}

// DecisionBoundary returns the x at which Predict(w,b,x) == threshold, i.e.
// where the model's classification flips. Returns NaN if w==0, since a flat
// prediction never crosses a threshold strictly between 0 and 1.
func DecisionBoundary(w, b, threshold float64) float64 {
	if w == 0 {
		return math.NaN()
	}
	logit := math.Log(threshold / (1 - threshold))
	return (logit - b) / w
}

// fitIters and fitLR are the fixed gradient-descent settings used to compute
// the "best fit" reference curve on every render -- deterministic (same
// inputs every call), so Render stays pure. See LEARNINGS.md for the
// resulting w,b this converges to.
const (
	fitIters = 10000
	fitLR    = 0.5
)

func render(p map[string]float64) string {
	guessW, guessB := p["guessW"], p["guessB"]
	bestW, bestB := FitLogisticRegression(HoursStudied, Passed, fitIters, fitLR)

	const xmin, xmax = 0.0, 6.0
	const ymin, ymax = -0.12, 1.12
	c := viz.New(680, 420, xmin, xmax, ymin, ymax)
	c.PadT = 90
	c.Axes()
	for x := 0.0; x <= xmax; x++ {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	// The 0.5 decision threshold, as a horizontal reference line.
	c.Path([][2]float64{{xmin, 0.5}, {xmax, 0.5}}, viz.Muted, 1)

	bestCurve := viz.Sample(xmin, xmax, 120, func(x float64) float64 { return Predict(bestW, bestB, x) })
	guessCurve := viz.Sample(xmin, xmax, 120, func(x float64) float64 { return Predict(guessW, guessB, x) })
	c.Path(bestCurve, viz.Good, 2.5)
	c.Path(guessCurve, viz.Accent, 2)

	// Your guess's decision boundary: where your curve crosses p=0.5.
	if guessBoundary := DecisionBoundary(guessW, guessB, 0.5); !math.IsNaN(guessBoundary) &&
		guessBoundary >= xmin && guessBoundary <= xmax {
		c.VLine(guessBoundary, viz.Accent, true)
	}

	for i, x := range HoursStudied {
		px, py := c.X(x), c.Y(Passed[i])
		c.Rect(px-4, py-4, 8, 8, viz.Ink, 1)
	}

	lossBest := LogLoss(HoursStudied, Passed, bestW, bestB)
	lossGuess := LogLoss(HoursStudied, Passed, guessW, guessB)

	c.Text(20, 22, fmt.Sprintf("best fit (green): w=%.2f b=%.2f    log-loss(best) = %.3f",
		bestW, bestB, lossBest), 14, viz.Good, "start")
	c.Text(20, 44, fmt.Sprintf("your guess (blue): w=%.2f b=%.2f    log-loss(yours) = %.3f",
		guessW, guessB, lossGuess), 14, viz.Accent, "start")
	c.Text(20, 64, "black squares = 10 students (0=failed, 1=passed)    gray line = 0.5 decision threshold",
		12, viz.Muted, "start")
	c.Text(20, 84, "lower log-loss is better; try to close the gap to log-loss(best)",
		12, viz.Muted, "start")

	return c.String()
}
