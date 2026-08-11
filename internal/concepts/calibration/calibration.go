// Package calibration answers a question the eval-playbook concept doesn't:
// even once a model separates the classes well (good loss, good ROC-AUC,
// good precision/recall), does its output NUMBER mean what it claims? A
// model that says "90% confident" should be right about 90% of the time it
// says that — ranking ability (AUC) says nothing about whether that's true.
// Temperature scaling — the same knob the sigmoid-softmax concept exposes —
// is the standard, cheap fix, and it never touches the ranking at all.
package calibration

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "calibration",
		Seq:   21,
		Title: "Calibration, ECE & temperature",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"Placeholder."},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "sep", Label: "Class separation", Min: 1, Max: 5, Step: 0.1, Def: 3},
			{Key: "temp", Label: "Temperature", Min: 0.3, Max: 3, Step: 0.1, Def: 0.5},
			{Key: "markP", Label: "Inspect confidence", Min: 0.5, Max: 0.99, Step: 0.01, Def: 0.9},
		},
		Render: render,
	})
}

// sigmoid squashes any real number into (0,1): 1/(1+e^-x).
func sigmoid(x float64) float64 { return 1 / (1 + math.Exp(-x)) }

// logit is sigmoid's inverse: ln(p/(1-p)).
func logit(p float64) float64 { return math.Log(p / (1 - p)) }

// normPDF is the unit-variance normal density centered at mu.
func normPDF(x, mu float64) float64 {
	d := x - mu
	return math.Exp(-0.5*d*d) / math.Sqrt(2*math.Pi)
}

// TrueLogit is the EXACT log-odds of a raw score z being positive, for two
// equal-variance classes — negatives ~ N(0,1), positives ~ N(sep,1), equal
// priors — the same generative setup precision-recall, roc-auc and pr-auc
// all use. Derived directly from Bayes' theorem (not fit, not estimated):
//
//	log[P(pos|z)/P(neg|z)] = log[f_pos(z)/f_neg(z)] = sep*z - sep²/2
//
// This is what a PERFECTLY calibrated model would report, in log-odds form.
func TrueLogit(z, sep float64) float64 {
	return sep*z - sep*sep/2
}

// TrueProbability is TrueLogit passed through sigmoid: the exact,
// perfectly-calibrated probability that score z came from the positive
// class. This is the number a model SHOULD report if calibration is the
// only thing being tested — its ranking of z is untouched either way.
func TrueProbability(z, sep float64) float64 {
	return sigmoid(TrueLogit(z, sep))
}

// ReportedProbability is what a model with temperature `temp` actually
// reports for score z. temp=1 reports TrueProbability exactly — perfectly
// calibrated. temp<1 sharpens the true logit before squashing (pushes
// probabilities toward 0/1 — overconfident, the common failure mode for
// modern neural nets). temp>1 flattens it toward 0.5 (underconfident).
// This is the identical temperature knob the sigmoid-softmax concept
// exposes, applied here to a logit with a known ground truth instead of an
// arbitrary one — which is what makes miscalibration measurable at all.
func ReportedProbability(z, sep, temp float64) float64 {
	if temp <= 0 {
		temp = 1e-6
	}
	return sigmoid(TrueLogit(z, sep) / temp)
}

// ObservedFrequency is the reliability-diagram curve itself: given a
// reported probability p from a model with temperature `temp`, what
// fraction of examples reported at that p are ACTUALLY positive? Derived by
// inverting ReportedProbability (it's monotonic in z for temp>0, so each p
// maps back to exactly one z): sigmoid(temp * logit(p)). At temp=1 this is
// the identity function (observed = predicted, p in, p out) — perfect
// calibration is exactly the diagonal line.
func ObservedFrequency(p, temp float64) float64 {
	if p <= 0 {
		p = 1e-6
	} else if p >= 1 {
		p = 1 - 1e-6
	}
	return sigmoid(temp * logit(p))
}

// ExpectedCalibrationError numerically integrates the population-weighted
// gap between a model's reported probability and the true, calibrated
// probability, across the whole population of raw scores z (weighted by how
// often each z actually occurs — half from N(0,1), half from N(sep,1)).
// This is the same quantity binned ECE estimates from a finite sample of
// predictions; computed here exactly via a fine grid instead of bins, since
// the generative model is known in closed form. 0 = perfectly calibrated
// (temp=1); grows as temp moves away from 1 in either direction.
func ExpectedCalibrationError(sep, temp float64, steps int) float64 {
	if steps < 2 {
		steps = 2
	}
	const lo, hi = -8.0, 8.0
	step := (hi - lo) / float64(steps)
	var totalGap, totalWeight float64
	for i := 0; i <= steps; i++ {
		z := lo + float64(i)*step
		weight := 0.5*normPDF(z, 0) + 0.5*normPDF(z, sep)
		gap := math.Abs(ReportedProbability(z, sep, temp) - TrueProbability(z, sep))
		totalGap += gap * weight
		totalWeight += weight
	}
	if totalWeight <= 0 {
		return 0
	}
	return totalGap / totalWeight
}

func render(p map[string]float64) string {
	sep, temp, markP := p["sep"], p["temp"], p["markP"]
	if temp <= 0.01 {
		temp = 0.01
	}
	if markP <= 0 {
		markP = 0.01
	}
	if markP >= 1 {
		markP = 0.99
	}

	c := viz.New(680, 460, 0, 1, 0, 1)
	c.PadT = 96 // room for four lines of header text above the diagram
	c.Axes()
	for x := 0.0; x <= 1.0; x += 0.2 {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	// Perfect calibration is exactly the diagonal.
	c.Path([][2]float64{{0, 0}, {1, 1}}, viz.Muted, 1.5)

	// The reliability curve itself, avoiding p=0/1 where logit blows up.
	curve := viz.Sample(0.01, 0.99, 240, func(pp float64) float64 {
		return ObservedFrequency(pp, temp)
	})
	c.Path(curve, viz.Accent, 2.5)

	// Mark the inspected confidence level: predicted vs. what it actually
	// means, with a dashed guide down to the x-axis.
	observed := ObservedFrequency(markP, temp)
	c.VLine(markP, viz.Warm, true)
	mx, my := c.X(markP), c.Y(observed)
	c.Rect(mx-4, my-4, 8, 8, viz.Warm, 1)

	ece := ExpectedCalibrationError(sep, temp, 400)
	verdict := "well-calibrated (temp ≈ 1)"
	if temp < 0.95 {
		verdict = "overconfident (temp < 1) — the common failure for trained neural nets"
	} else if temp > 1.05 {
		verdict = "underconfident (temp > 1)"
	}

	c.Text(20, 24, fmt.Sprintf("temperature = %.2f    ECE ≈ %.3f    separation = %.1f", temp, ece, sep),
		14, viz.Ink, "start")
	c.Text(20, 44, verdict, 13, viz.Warm, "start")
	c.Text(20, 64, fmt.Sprintf("when the model says %.0f%% confident, it's actually right about %.0f%% of the time",
		markP*100, observed*100), 13, viz.Muted, "start")
	c.Text(20, 84, "x = predicted probability, y = observed frequency    grey diagonal = perfect calibration",
		12, viz.Muted, "start")

	return c.String()
}
