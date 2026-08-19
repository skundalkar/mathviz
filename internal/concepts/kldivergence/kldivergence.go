// Package kldivergence visualizes Kullback-Leibler divergence for two
// two-outcome (coin-flip) distributions: how many extra bits, on average,
// you waste describing outcomes drawn from a true distribution P using a
// code built for a different, approximating distribution Q — and why that
// number changes depending on which of the two you call "true."
package kldivergence

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "kl-divergence",
		Seq:   49,
		Title: "KL divergence (the cost of believing the wrong distribution)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`entropy`'s real-life section already named 'cross-entropy loss' as 'the " +
						"same core idea extended one step further' — but it didn't say what that " +
						"extension actually measures. Say a trick coin truly lands heads 90% of " +
						"the time (the real distribution, call it P), but a model that's never " +
						"seen this coin assumes it's fair and believes heads and tails are 50/50 " +
						"(its belief, call it Q) — the model is confidently wrong. Entropy alone " +
						"can tell you P is fairly predictable and Q is maximally spread out, but " +
						"it can't tell you the two disagree, or by how much. Is there a single " +
						"number that measures the gap between what a model believes and what's " +
						"actually true?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Same setup: true P(heads)=0.9, model's believed Q(heads)=0.5. First, " +
						"`entropy`'s own formula gives P's true entropy: H(P) = " +
						"-0.9·log2(0.9) - 0.1·log2(0.1) = 0.137 + 0.332 = 0.469 bits — the best " +
						"any code could possibly do, on average, if it were built knowing P " +
						"exactly.",
					"Next, cross-entropy H(P,Q) measures the average bits actually spent if you " +
						"build your code assuming Q instead, but outcomes keep coming from the " +
						"real P: H(P,Q) = -0.9·log2(0.5) - 0.1·log2(0.5) = -log2(0.5)·(0.9+0.1) = " +
						"1 bit exactly (both terms use q=0.5, so they collapse to -log2(0.5)·1).",
					"KL divergence is exactly that wasted gap: KL(P‖Q) = H(P,Q) - H(P) = " +
						"1 - 0.469 = 0.531 bits — the extra bits per outcome you pay for trusting " +
						"Q=0.5 when reality is really P=0.9. Now swap which distribution you call " +
						"'true': KL(Q‖P) = H(Q,P) - H(Q). H(Q)=1 bit (a fair coin is already " +
						"maximally uncertain). H(Q,P) = -0.5·log2(0.9) - 0.5·log2(0.1) = 0.076 + " +
						"1.661 = 1.737 bits. KL(Q‖P) = 1.737 - 1 = 0.737 bits — a different number " +
						"from 0.531, even though it's built from the exact same two coins.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The p slider sets the true distribution P(heads); the q slider sets the " +
						"model's believed distribution Q(heads). Two bar pairs show P and Q's " +
						"heads/tails probabilities side by side. The readout reports H(P) (P's " +
						"own entropy), H(P,Q) (the cross-entropy cost of using Q's code on P's " +
						"outcomes), and KL(P‖Q) = H(P,Q)-H(P) highlighted as the wasted-bits gap — " +
						"alongside KL(Q‖P) for direct comparison. Drag q away from p and watch " +
						"both KL numbers grow from 0; drag q back to equal p and both collapse " +
						"back to exactly 0, since there's no gap left to measure once the two " +
						"distributions agree.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Quantify miscalibration in bits instead of a vague 'the model's confidence " +
						"is off': not just 'q=0.5 seems wrong' but 'this model wastes 0.531 bits " +
						"per outcome trusting q=0.5 when reality is p=0.9.' This is also exactly " +
						"what a classifier's 'cross-entropy loss' is doing during training: since " +
						"H(P) is fixed by the (unchangeable) true label distribution, minimizing " +
						"cross-entropy loss H(P,Q) over the model's Q is mathematically identical " +
						"to minimizing KL(P‖Q) — training literally is squeezing this gap toward " +
						"zero.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Every classifier trained with 'cross-entropy loss' — image classifiers, spam " +
						"filters, next-word predictors — is minimizing exactly this gap between the " +
						"true label distribution and its own predicted distribution. Language " +
						"model 'perplexity,' a standard benchmark number, is a repackaged version " +
						"of cross-entropy on real text. Recommendation systems compare a user's " +
						"actual click/watch distribution against the system's predicted preference " +
						"distribution with KL divergence to flag when a model's assumptions have " +
						"drifted from real behavior. A/B testers and pollsters use KL divergence to " +
						"quantify how different two survey-response distributions are, beyond just " +
						"eyeballing two bar charts.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: KL(P‖Q) and KL(Q‖P) are generally different numbers — " +
						"order matters, so always say which distribution is the true/reference one " +
						"(listed first) and which is the approximation (listed second).",
					"Not like this: treating KL divergence like a symmetric 'distance' between two " +
						"distributions, the way Euclidean distance between two points doesn't care " +
						"which point you start from. The worked example already shows this fails: " +
						"swapping which coin is 'true' changed 0.531 bits into 0.737 bits — because " +
						"KL(P‖Q) weights every surprise term by P's own probabilities specifically, " +
						"not by some neutral halfway blend of P and Q.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "p", Label: "True P(heads)", Min: 0.01, Max: 0.99, Step: 0.01, Def: 0.9},
			{Key: "q", Label: "Believed Q(heads)", Min: 0.01, Max: 0.99, Step: 0.01, Def: 0.5},
		},
		Render: render,
	})
}

func render(params map[string]float64) string {
	_ = params
	return viz.New(680, 460, -1, 1, -1, 1).String()
}
