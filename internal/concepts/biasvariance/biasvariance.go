// Package biasvariance visualizes the bias-variance decomposition:
// retraining the same-degree polynomial on many independently resampled
// noisy datasets and looking at how its predictions at one fixed point
// behave, not just once (as `overfitting` does) but across many retrainings
// — separating "consistently wrong" (bias) from "unpredictable" (variance).
package biasvariance

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "bias-variance-tradeoff",
		Seq:   50,
		Title: "Bias-variance tradeoff",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`overfitting` already showed that a high-degree curve chases noise instead " +
						"of the true pattern — but that was one look at one noisy dataset. " +
						"Retrain that exact same-degree model on a DIFFERENT random sample of the " +
						"same underlying pattern, and does it land in the same wrong place every " +
						"time, or does it land somewhere different each time? Those are two very " +
						"different kinds of 'wrong,' and a single-dataset picture can't tell them " +
						"apart — is there a way to separate 'the model is consistently biased " +
						"toward the wrong answer' from 'the model is unpredictable, landing " +
						"somewhere different every time it's retrained'?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Same setup as `overfitting`: 12 noisy points sampled from the true pattern " +
						"sin(x), a degree-d polynomial fit by least squares. But now retrain that " +
						"same-degree model on 30 independently resampled noisy datasets (same true " +
						"pattern, fresh noise each time), and look at all 30 fitted curves' " +
						"predictions at one fixed point, x0=1.5 (true value sin(1.5)≈0.997):",
					"• degree=1 (straight line): the 30 predictions cluster tightly together " +
						"(variance≈0.0076 — barely any spread run to run) but their average, " +
						"0.321, sits far from the true 0.997 — bias²≈0.458. A straight line " +
						"structurally can't bend to follow sin's curve near its peak, so it's " +
						"consistently, predictably wrong, no matter which noisy sample it's " +
						"retrained on.",
					"• degree=9 (very flexible): the 30 predictions scatter widely " +
						"(variance≈0.033) but their average, 0.977, lands close to the true 0.997 " +
						"— bias²≈0.0004. Flexible enough to track the true curve on average, but " +
						"exactly which quirks of noise it latches onto varies wildly from one " +
						"training sample to the next.",
					"• degree=3 (moderate): bias²≈0.0026 and variance≈0.0188 — both small, and " +
						"their sum, 0.0213, beats both extremes (0.465 for degree 1, 0.034 for " +
						"degree 9).",
					"Total expected error decomposes exactly into these two pieces: " +
						"E[(prediction - true)²] = bias² + variance. Degree 1 fails almost " +
						"entirely from bias; degree 9 fails almost entirely from variance; " +
						"degree 3 balances both to land at the lowest total.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"degree sets model complexity, the same slider as `overfitting`; noise sets " +
						"how much each training sample's y-values scatter from the true pattern; " +
						"x0 moves the vertical evaluation line where bias and variance are " +
						"measured. Every faint curve is one polynomial fit to one independently " +
						"resampled noisy dataset — 30 of them drawn together as a 'spaghetti " +
						"plot.' At low degree, all 30 strands bunch tightly on top of each other " +
						"but visibly miss the true dashed curve wherever it bends the most (high " +
						"bias, low variance). At high degree, the strands fan out wildly, each " +
						"chasing its own dataset's particular noise (low bias, high variance). The " +
						"dots mark where each strand crosses the x0 line — their spread is the " +
						"variance number in the readout, and how far their average sits from the " +
						"true curve is the bias.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Diagnose why a model is underperforming, not just that it is: a model with " +
						"high bias and low variance is consistently missing the pattern, and needs " +
						"more flexibility (more features, higher degree, a bigger model) to fix; a " +
						"model with low bias and high variance is capable of finding the pattern " +
						"but is unstable, and needs regularization, more training data, or a " +
						"simpler model to fix — opposite prescriptions for what looks, from " +
						"training accuracy alone, like the same kind of failure. Sweep the degree " +
						"slider from 1 to 11 and watch total error trace a U-shape, exactly the " +
						"sweet-spot degree `overfitting` gestured at without naming what was " +
						"actually trading off underneath it.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A linear regression model on a genuinely curved relationship (like " +
						"`overfitting`'s straight-line problem) is high-bias: retrain it on new " +
						"data all year and it keeps making the same kind of systematically wrong " +
						"prediction. A deep decision tree with no depth limit is high-variance: " +
						"retrain it on a slightly different sample of the same customers and it " +
						"can produce a noticeably different tree with different predictions, even " +
						"though the underlying pattern didn't change. 'Ensemble' methods exploit " +
						"this asymmetry directly: bagging (as in random forests) averages many " +
						"high-variance, low-bias trees together, and averaging cancels out " +
						"variance — the same reason `law-of-large-numbers` shrinks toward the " +
						"true mean — while leaving each tree's low bias intact.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: bias and variance are two separate failure modes, and the " +
						"fix for one often makes the other worse — a simpler model reduces " +
						"variance but tends to raise bias; a more flexible model reduces bias but " +
						"tends to raise variance, exactly what the degree slider demonstrates " +
						"directly.",
					"Not like this: assuming 'the model has high error' by itself tells you which " +
						"one is the problem, or that adding more training data always helps — more " +
						"data mainly shrinks variance (each individual fit has less noise to " +
						"overfit to), so it does little for a model whose real problem is high " +
						"bias. Degree 1 in this picture stays badly biased no matter how many " +
						"resampled datasets you average over, because a straight line simply " +
						"can't bend.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "degree", Label: "Model complexity (degree)", Min: 1, Max: 11, Step: 1, Def: 3},
			{Key: "noise", Label: "Noise", Min: 0, Max: 1, Step: 0.05, Def: 0.4},
			{Key: "x0", Label: "Evaluation point (x0)", Min: -3.2, Max: 3.2, Step: 0.2, Def: 1.5},
		},
		Render: render,
	})
}

func render(params map[string]float64) string {
	_ = params
	return viz.New(680, 460, -1, 1, -1, 1).String()
}
