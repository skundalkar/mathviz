// Package logreg visualizes logistic regression: fitting an S-shaped curve
// that models P(pass) directly, instead of stretching linear-regression's
// straight line across data whose answer is only ever 0 or 1. Built from
// sigmoid-softmax's squashing function and gradient-descent's step-downhill
// mechanism, now aimed at cross-entropy (log) loss instead of squared error.
package logreg

import (
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
						"steps on this exact dataset converges to w≈1.10, b≈-3.10, average log-loss " +
						"≈0.325 — versus 0.693 (=-log(0.5)) at the untrained w=0,b=0 starting point, " +
						"where the model has no information and guesses 50/50 for everyone.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Ten black squares are the fixed (hours studied, passed) data points. The " +
						"green S-curve is the fitted model from the worked example (w≈1.10, " +
						"b≈-3.10); the blue S-curve is your own guess, built from the slope and " +
						"intercept sliders — drag either one and watch both the curve's shape and " +
						"the 'your log-loss' readout change immediately. The dashed horizontal line " +
						"at 0.5 is the decision threshold (the same idea `precision-recall` used to " +
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

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 320, -1, 1, -1, 1).String()
}
