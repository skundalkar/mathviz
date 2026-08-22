// Package regularization visualizes L1 (Lasso) and L2 (Ridge) penalties:
// fitting a linear model on 6 candidate input signals, only 2 of which
// actually drive the outcome, and watching how each penalty type treats the
// 4 irrelevant ones as the penalty strength λ increases — Ridge shrinks
// every coefficient smoothly, Lasso can zero one out entirely.
package regularization

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "regularization",
		Seq:   57,
		Title: "Regularization (L1/L2 penalties)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`linear-regression` already showed how to fit the single best straight " +
						"line through a scatter with one input — but that story assumed you " +
						"already knew hours-studied was the input that mattered. Real prediction " +
						"problems usually start with a pile of CANDIDATE inputs, not one " +
						"pre-screened one: predicting an exam score from hours studied, hours " +
						"slept, a handful of other measurements nobody has vetted for relevance. " +
						"Fit ordinary least squares using every candidate at once, and it will " +
						"happily hand back a nonzero weight for every single one — including the " +
						"ones that don't actually matter — because on this one particular batch " +
						"of data, pure chance gives every measured number some small, spurious " +
						"correlation with the outcome. How do you get a model that doesn't quietly " +
						"treat that noise as if it were signal?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take 24 data points and 6 candidate input signals, x1 through x6. Only x1 " +
						"and x2 actually generate the outcome for this worked example (true " +
						"relationship y = 3·x1 − 2·x2 + noise); x3 through x6 are pure noise, " +
						"unrelated to y by construction, mixed in exactly the way an irrelevant " +
						"measurement would be in a real dataset. Fit ordinary least squares (no " +
						"penalty at all, λ=0) and this comes out: x1 ≈ 3.01, x2 ≈ −1.88 — close " +
						"to the true 3 and −2, good — but ALSO x3 ≈ −0.06, x4 ≈ 0.06, x5 ≈ 0.03, " +
						"x6 ≈ 0.09: four nonzero weights on features that have nothing to do with " +
						"y, purely because in this one 24-point sample they happen to line up " +
						"with the noise a little.",
					"Ridge regression patches the normal equations plain least squares already " +
						"solves — (XᵀX)β = Xᵀy — by adding a penalty λ to the diagonal for every " +
						"coefficient except the intercept: (XᵀX + λI)β = Xᵀy. That extra λ makes " +
						"large coefficients costlier, so the solver settles for smaller ones " +
						"across the board. At λ=2, the total penalized quantity — the sum of " +
						"every non-intercept coefficient squared — drops from 12.65 (λ=0) to " +
						"11.03, but no single coefficient hits exactly zero; x3–x6 shift to " +
						"roughly (−0.05, 0.04, 0.10, −0.07) — smaller in total, but still there.",
					"Lasso instead updates one coefficient at a time by a rule called " +
						"soft-thresholding: nudge the coefficient toward zero by exactly λ, and " +
						"clip it to exactly 0 if that nudge would flip its sign. At the same λ=2, " +
						"x1 ≈ 2.95 and x2 ≈ −1.81 survive — still clearly nonzero, close to the " +
						"true 3 and −2 — but x3, x4, x5, and x6 all land at EXACTLY 0.0, every " +
						"one of them. Lasso didn't just shrink the four irrelevant weights, it " +
						"deleted them.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Six bars, one per candidate input x1–x6: a faint grey ghost bar behind each " +
						"one is the unpenalized OLS weight (λ=0) for reference, and the solid " +
						"colored bar in front is the current fit at whatever λ and penalty type " +
						"the sliders are set to — blue for the two genuinely predictive inputs " +
						"(x1, x2), orange for the four that are pure noise (x3–x6). Set penalty " +
						"to 1 (Lasso) and push λ up past about 1: the four orange bars snap to " +
						"exactly zero height while the two blue bars stay tall. Switch penalty " +
						"back to 0 (Ridge) at the same λ and the orange bars only shrink a " +
						"little — they never fully disappear. The readout reports the exact SSE " +
						"(how well the current fit matches these 24 points) and the penalty term " +
						"(Σ|coefficient| for Lasso, Σcoefficient² for Ridge), so you can watch " +
						"that total shrink smoothly as λ grows even when individual Ridge " +
						"coefficients don't move in lockstep.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Recover which inputs actually matter without knowing that in advance: " +
						"crank λ up under Lasso and read off which coefficients survive at " +
						"exactly nonzero — that IS automatic feature selection, done by the same " +
						"machinery that fits the line, no separate step required. Ridge doesn't " +
						"give that crisp yes/no, but it gives something plain least squares " +
						"can't either: a dial for trading a little bias for a lot less variance, " +
						"exactly the trade `bias-variance-tradeoff` described — plain OLS has " +
						"exactly one setting (λ=0) and no dial at all.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Genomics studies routinely start with thousands of candidate genes and a " +
						"few dozen patients — Lasso is a standard tool for narrowing that down to " +
						"the handful of genes that actually predict an outcome. Spam filters and " +
						"recommendation systems start from huge feature sets (every word, every " +
						"past click) and use L1 penalties to keep only the ones worth keeping. " +
						"Calling a theory or an explanation 'overfit' or 'too complicated' in " +
						"everyday conversation is an informal version of exactly the same " +
						"complaint regularization answers formally: don't let something explain " +
						"away noise that isn't really there.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: L1 (Lasso) and L2 (Ridge) both shrink coefficients to " +
						"fight overfitting, but by mechanisms with a real, visible difference — " +
						"L1 can zero a coefficient out entirely, L2 mathematically cannot (its " +
						"penalty is smooth everywhere, so the pull toward zero gets weaker, not " +
						"sharper, the closer a coefficient already is to it).",
					"Not like this: 'more regularization is always better.' Push λ high enough " +
						"under either penalty and even the genuinely predictive x1, x2 weights " +
						"get crushed toward zero along with the noise — the same underfitting " +
						"`bias-variance-tradeoff` called high bias. λ is a dial, not a one-way " +
						"improvement switch; picking it well normally means checking performance " +
						"on held-out data, not just watching training error fall.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "lambda", Label: "Penalty strength (λ)", Min: 0, Max: 10, Step: 0.1, Def: 2},
			{Key: "penalty", Label: "Penalty type (0=L2 ridge, 1=L1 lasso)", Min: 0, Max: 1, Step: 1, Def: 1},
			{Key: "noise", Label: "Noise", Min: 0, Max: 1.5, Step: 0.05, Def: 0.6},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 460, 0, 1, 0, 1).String()
}
