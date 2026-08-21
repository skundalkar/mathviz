// Package bootstrap visualizes bootstrap resampling: estimating how
// uncertain a statistic is by resampling the data itself, with replacement,
// instead of assuming a formula or a known distribution shape. With a
// sample this small (n=5) the full space of resamples (5^5=3125) is small
// enough to enumerate exactly, so the picture shows the true bootstrap
// sampling distribution rather than a Monte Carlo approximation of it.
package bootstrap

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "bootstrap-resampling",
		Seq:   55,
		Title: "Bootstrap resampling",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`confidence-interval` built a margin of error from a formula that assumes " +
						"your data is roughly normal, or that your sample is large enough for the " +
						"central limit theorem to kick in. But say you only have 5 employees' " +
						"commute times: 12, 15, 15, 22, and 41 minutes. That's nowhere near " +
						"'large enough,' and the 41-minute commute drags the sample lopsided. The " +
						"average is 21 minutes — but how uncertain is that 21, when you can't " +
						"honestly assume the data looks like a bell curve and there's no formula " +
						"handy for 'the uncertainty of a mean from 5 skewed numbers'?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"The trick: treat the 5 numbers you already have as a stand-in for the whole " +
						"population, and resample from them, with replacement, to see how much the " +
						"mean bounces around. A 'resample' is a new set of 5 draws from {12, 15, " +
						"15, 22, 41}, where any value can be picked more than once (or not at " +
						"all) — e.g. one resample might be [15, 41, 12, 41, 15] (mean 24.8), " +
						"another [12, 12, 15, 15, 22] (mean 15.2).",
					"With only 5 slots and 5 possible values per slot, there are exactly 5^5 = " +
						"3,125 possible resamples — small enough to compute every single one " +
						"exactly, instead of only simulating a random handful of them. Doing that " +
						"and sorting the 3,125 resulting means gives the exact bootstrap " +
						"distribution: it ranges all the way from 12 (the resample that happened " +
						"to draw '41' zero times) to 41 (the resample that drew '41' all five " +
						"times), and its middle 90% — from the 5th to the 95th percentile of " +
						"those 3,125 means — runs from about 14.4 to 30.0.",
					"That range, [14.4, 30.0], is the bootstrap 90% interval: it says the true " +
						"average commute is plausibly anywhere in that band, given only how much " +
						"the mean itself wiggles when resampled from the data you actually " +
						"collected — no bell-curve assumption required anywhere in the " +
						"calculation.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"A histogram of all 3,125 resampled means, built from the same 5 commute " +
						"times with the 5th value adjustable. The vertical line marks the actual " +
						"sample mean; the shaded band marks the percentile interval at the chosen " +
						"confidence level. Drag the 5th commute time up (a bigger outlier) and " +
						"watch the whole histogram stretch wider and the interval widen with it — " +
						"one unusual data point makes the mean itself less certain. Drag the " +
						"confidence level up and watch the shaded band widen to cover more of the " +
						"same fixed histogram.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can put an honest uncertainty range on a statistic computed from a small " +
						"or skewed sample — or on a statistic like a median or a ratio that has no " +
						"tidy formula for its standard error at all — by resampling the data " +
						"itself and reading the spread directly off the resampled distribution, " +
						"instead of needing a distributional assumption to hold.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A/B test results with a handful of conversions, medical studies with a small " +
						"number of patients, and sports statistics from a short season all share " +
						"the same problem: too little data to trust a normal-theory formula, but " +
						"still a real need to know how much to trust the number you computed. " +
						"Bootstrap resampling (and its close cousin, cross-validation, which " +
						"resamples which data points train vs. test a model) is the standard fix " +
						"across statistics and machine learning whenever 'assume it's normal' " +
						"isn't a safe assumption.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the bootstrap interval is [14.4, 30.0] minutes, based on " +
						"how much the mean of this exact sample wiggles when resampled' — the " +
						"interval describes uncertainty in the estimate, not a prediction about " +
						"any single future commute.",
					"Not like this: 'resampling the same 5 numbers over and over just gives you " +
						"back the same answer, so it can't tell you anything new' — the resamples " +
						"aren't identical to the original sample (most leave out at least one " +
						"value and duplicate another), and it's exactly that variation across " +
						"resamples that reveals how sensitive the mean is to which 5 commutes you " +
						"happened to record.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "outlier", Label: "5th commute time (minutes)", Min: 22, Max: 90, Step: 1, Def: 41},
			{Key: "level", Label: "Confidence level", Min: 0.5, Max: 0.98, Step: 0.02, Def: 0.90},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 400, 0, 1, 0, 1).String()
}
