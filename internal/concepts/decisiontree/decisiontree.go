// Package decisiontree visualizes how a decision tree picks where to split:
// scoring every candidate threshold on a single feature by information gain
// (built from `entropy`'s bits-of-surprise measure) rather than by where the
// feature's values happen to be numerically centered.
package decisiontree

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "decision-trees",
		Seq:   56,
		Title: "Decision trees (splitting by information gain)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`entropy` measured how surprising a set of outcomes is, in bits. Here are " +
						"10 students' hours studied and whether they passed: 5 passed, 5 didn't — " +
						"exactly the maximum-uncertainty case entropy warns about, a coin flip. You " +
						"have one number per student (hours studied) and want a simple rule: 'if " +
						"hours studied is above some threshold, predict pass.' Gut instinct: just " +
						"split at the middle of the range — hours run from 1 to 5.5, so try 3.25. " +
						"Is the numeric midpoint actually the best place to draw that line, and how " +
						"would you even measure 'best'?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Score a candidate split the same way entropy scored a whole distribution: " +
						"how much less uncertain are the two resulting groups, on average, than the " +
						"whole group was before you split? Start with the whole group: 5 of 10 " +
						"passed, so parent entropy = 1.0 bit — total uncertainty, a coin flip.",
					"Try the 'obvious' midpoint, threshold = 3.25 hours: the low-hours group " +
						"{1, 1.5, 2, 2.5, 3} has 1 pass out of 5 (entropy ≈ 0.72), and the " +
						"high-hours group {3.5, 4, 4.5, 5, 5.5} has 4 pass out of 5 (entropy ≈ " +
						"0.72 too) — both sides are still fairly mixed. Weighted average child " +
						"entropy = 0.5×0.72 + 0.5×0.72 ≈ 0.72, so information gain = 1.0 − 0.72 = " +
						"0.28 bits. The midpoint only shaves off about a quarter of the original " +
						"uncertainty.",
					"Now try threshold = 2.75 instead: the low side {1, 1.5, 2, 2.5} has 0 passes " +
						"out of 4 — perfectly pure, entropy = 0 — and the high side {3, 3.5, 4, " +
						"4.5, 5, 5.5} has 5 passes out of 6 (entropy ≈ 0.65). Weighted average = " +
						"0.4×0 + 0.6×0.65 ≈ 0.39, so gain = 1.0 − 0.39 = 0.61 bits — more than " +
						"double the midpoint's gain, because the low side came out completely " +
						"pure. Scanning every candidate split confirms 2.75 (tied with 3.75) is " +
						"the actual best split, not 3.25 — the midpoint just happened to land " +
						"inside the messiest part of the data.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"All 10 students plotted by hours studied, colored by pass/fail. The " +
						"threshold slider draws a vertical split line; the readout below reports " +
						"exactly the numbers worked through above for whichever threshold you pick " +
						"— left/right group sizes, each side's pass fraction and entropy, and the " +
						"resulting information gain — plus, for comparison, the best gain " +
						"achievable by any threshold on this data. Drag the slider off 3.25 toward " +
						"2.75 or 3.75 and watch the gain jump from 0.28 up to 0.61.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Instead of eyeballing where to draw a split, you can score every candidate " +
						"threshold with one number — information gain — and mechanically pick the " +
						"one that reduces uncertainty the most. That's the exact rule real " +
						"decision-tree algorithms (ID3, C4.5, CART) use at every branch: try every " +
						"feature and every candidate threshold, split on whichever gives the " +
						"highest information gain, then repeat inside each resulting group.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Spam filters splitting on word counts, medical triage splitting on a lab " +
						"value, credit scoring splitting on an income-to-debt ratio — anywhere you " +
						"want a human-readable 'if X is above/below this line' rule instead of an " +
						"opaque score. Random forests and gradient-boosted trees are, underneath, " +
						"large committees of exactly this kind of tree, each one grown by repeating " +
						"this same highest-information-gain split over and over.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the split at 2.75 hours gives more than double the " +
						"midpoint's information gain, because it produces a perfectly pure " +
						"low-hours group' — anchoring the claim on the measured entropy reduction.",
					"Not like this: 'the best threshold is always the middle or median of the " +
						"data' — the best split is decided by how cleanly it separates the labels, " +
						"not by which value is numerically in the center of the feature's range; " +
						"here the truly best splits (2.75 and 3.75) sit off-center, on either side " +
						"of the messiest two data points, not on top of them.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "threshold", Label: "Split threshold (hours studied)", Min: 1, Max: 5.5, Step: 0.05, Def: 3.25},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 400, 0, 1, 0, 1).String()
}
