// Package zscore visualizes standardization: converting a raw value into
// "how many standard deviations above or below its own mean" so values
// pulled from two different distributions (different means, different
// spreads) can be compared on one common scale.
package zscore

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "z-score",
		Seq:   29,
		Title: "Z-scores (standardizing onto one common scale)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Two students take different tests. On Test A (class average 70, typical " +
						"spread 10 points), Priya scores 90. On Test B (class average 85, typical " +
						"spread 10 points), Raj scores 92. Whose result was more impressive? Gut " +
						"instinct: just compare the two scores — 92 > 90, so Raj did better. That " +
						"instinct ignores something important: a 90 on a test where most people " +
						"scored near 70 is a very different kind of result than a 92 on a test where " +
						"most people already scored near 85. Comparing two raw numbers only makes " +
						"sense when they're measured on the same scale, starting from the same " +
						"baseline — and these two tests aren't. Is there a way to compare two scores " +
						"from two completely different distributions on equal footing?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Do it in two separate steps, and see why each one is needed on its own.",
					"Step 1 — recenter, by subtracting the mean: this turns a raw score into " +
						"'how far above or below average,' in the test's own points.",
					"• Priya: 90 − 70 = 20 points above Test A's average.",
					"• Raj: 92 − 85 = 7 points above Test B's average.",
					"Just from re-centering, the comparison already flips — Priya's result sits " +
						"20 raw points above her class's baseline while Raj's sits only 7 points " +
						"above his. Recentering alone was enough to change the ranking, but raw " +
						"points aren't quite the finish line either — 20 points is a big deal on a " +
						"test where scores rarely stray more than 10 points from the mean, and " +
						"barely notable on a test where scores routinely swing by 40. That's where " +
						"the standard deviation comes in.",
					"Step 2 — rescale, by dividing by the standard deviation: this turns 'how many " +
						"raw points above average' into 'how many typical-sized swings above " +
						"average,' so gaps become comparable even when two tests have different " +
						"amounts of spread. Compare two more students who both land exactly 10 " +
						"points above their own class's average, but on tests with very different " +
						"spreads:",
					"• Class C (mean 70, stddev 5): Zara scores 80 → 10 raw points above average " +
						"→ 10/5 = 2.0 standard deviations above average — a big deal, since scores " +
						"in Class C rarely stray more than 5 points either way.",
					"• Class D (mean 70, stddev 20): Wes also scores 80 → the same 10 raw points " +
						"above average → 10/20 = 0.5 standard deviations above average — pretty " +
						"ordinary, since scores in Class D routinely swing by 20 points.",
					"Put both steps together and you get the z-score: z = (x − mean) / stddev. " +
						"Back to the opening scenario: Priya's z = (90−70)/10 = 2.0, Raj's z = " +
						"(92−85)/10 = 0.7 — Priya's result, standing 2 full standard deviations " +
						"above her class's average, is the far more exceptional one, even though " +
						"her raw score was lower.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"A bell curve centered at the mean μ with width set by the standard deviation " +
						"σ — the shape of scores in one particular class or test. A dashed vertical " +
						"line marks μ itself; a solid line marks the raw score x you've dialed in, " +
						"and the shaded band between them is exactly the gap step 2 measures. The " +
						"z-score readout above the curve is that same gap, expressed as a count of " +
						"standard deviations instead of raw points. Drag the mean and the whole " +
						"curve slides sideways without changing shape; drag the standard deviation " +
						"and the curve widens or narrows — in both cases the shaded gap and the " +
						"z-score update to reflect exactly how far out on this particular curve the " +
						"score x really sits.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Rank or compare measurements pulled from different distributions on one " +
						"common scale, instead of only being able to compare numbers measured the " +
						"exact same way. A z-score also converts directly into a percentile under a " +
						"normal distribution — z = 2.0 lands at roughly the 97.7th percentile, " +
						"meaning about 97.7% of the class scored below that point — turning 'how " +
						"many standard deviations above average' into 'better than roughly this " +
						"percentage of everyone else,' a number that's meaningful on its own " +
						"without ever needing to know the original scale.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Standardized test scores (SAT, IQ tests) are reported as scaled scores " +
						"precisely so a result from one test date or version can be compared fairly " +
						"to another. Growth charts for children's height and weight use z-scores " +
						"(often shown as 'percentiles' on the chart) to flag values that are " +
						"unusually far from what's typical for a given age. Manufacturing and " +
						"quality control flag a part as defective when a measurement's z-score " +
						"crosses a threshold, since 'how many typical-sized deviations from spec' " +
						"matters more than the raw measurement alone. And 'standardizing' or " +
						"'normalizing' data before feeding it into many machine learning models is " +
						"literally computing a z-score for every feature, so a feature measured in " +
						"thousands (like income) doesn't dominate one measured in single digits " +
						"(like years of education) purely because of its raw scale.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'a z-score tells you how many standard deviations a value " +
						"sits above or below its own mean' — z=1.5 means 1.5 typical-sized steps " +
						"above average, on whatever scale that particular distribution uses. Not " +
						"like this: comparing two raw scores directly when they come from different " +
						"distributions (as the opening scenario shows, the higher raw score isn't " +
						"always the better result), or forgetting that a z-score is meaningless " +
						"without knowing which distribution's mean and standard deviation produced " +
						"it — z=2.0 always means 'unusually high relative to its own group,' but " +
						"the raw value that corresponds to depends entirely on that group's mean " +
						"and spread.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "mean", Label: "Class mean μ", Min: 40, Max: 100, Step: 1, Def: 70},
			{Key: "stddev", Label: "Class stddev σ", Min: 2, Max: 25, Step: 1, Def: 10},
			{Key: "score", Label: "Raw score x", Min: 0, Max: 140, Step: 1, Def: 90},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(680, 380, 0, 1, 0, 1)
	return c.String()
}
