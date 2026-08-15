// Package lln visualizes the law of large numbers: a running average of
// repeated coin flips, jumping around wildly at small n and settling in on
// the true probability as n grows -- without ever guaranteeing that any
// single next flip "corrects" the average.
package lln

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "law-of-large-numbers",
		Seq:   37,
		Title: "Law of large numbers (the running average settles down)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You flip a fair coin and it comes up heads three times in a row. Gut " +
						"instinct — the 'gambler's fallacy' — says tails is 'due': the coin must " +
						"somehow correct itself soon, or the long-run 50/50 split everyone knows " +
						"about wouldn't hold up. That instinct assumes something is nudging future " +
						"flips to balance out past ones. But a coin has no memory of its last flip " +
						"— so if nothing is correcting anything, how does the fraction of heads " +
						"reliably end up near 50% after enough flips? What's actually doing the " +
						"settling, if not some correcting force?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Flip a fair coin (p=0.5) over and over and track the running average — " +
						"heads-so-far divided by flips-so-far — after each one:",
					"• n=1: average = 0% (a single tails). n=5: average = 60%. n=10: average = " +
						"50%, exactly — not because anything corrected, just where the running " +
						"count happened to land.",
					"• n=20: 45%. n=50: 52% — closer to 50% than n=20 was.",
					"• n=100: 60% — farther from 50% than n=50's 52% was! More flips didn't " +
						"guarantee a closer answer at every single step.",
					"• n=200: 53%. n=300: 53% — settled in close, and staying close.",
					"No flip corrected anything — what actually shrinks is how much any new flip " +
						"can move the average. A running average updates as new_average = " +
						"old_average + (this_flip − old_average)/n: one more flip is only ever " +
						"weighted 1/n against everything already counted. At n=10, one flip can " +
						"swing the average by up to 10%; at n=300, one flip can only swing it by up " +
						"to about 0.33%. That's the whole mechanism: dilution by a growing " +
						"denominator, not correction. It shows up directly in this run's own " +
						"numbers: the running average's spread (variance) over its first 20 flips " +
						"is about 0.0153; over its last 20 (flips 281-300) it's about 0.0000072 — " +
						"roughly 2,100 times smaller, purely because n is 15x larger by then.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The p slider sets the coin's true probability of heads; the n slider marks " +
						"a point on the running-average curve. The blue curve plots the running " +
						"average against the number of flips so far, from 1 up to 300 — wide, " +
						"jagged swings on the left where each flip is a big fraction of the total, " +
						"smoothing into a flat line on the right. The dashed horizontal line is the " +
						"true probability p; the orange marker and readout show exactly where the " +
						"current n sits, and how far its running average still is from p. Drag n " +
						"back and forth and watch the marker ride the curve's actual bumps — " +
						"including places where it moves briefly farther from p, not closer.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Tell apart two claims that sound similar but aren't: 'the average of many " +
						"flips converges to p' (true — the law of large numbers) and 'each " +
						"individual flip's odds shift to help the average converge' (false — the " +
						"gambler's fallacy). You can also judge, roughly, how much to trust an " +
						"average from a small sample: since one new data point can only move an " +
						"n-point average by up to 1/n, an average built from 10 data points is " +
						"still easy to knock around, while one built from 1,000 is nearly " +
						"immovable by any single new observation — a concrete reason to distrust a " +
						"strong claim ('this coin is rigged!') based on a handful of flips.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Polling and surveys: asking 10 people how they'll vote gives an estimate " +
						"that one or two people can swing wildly, while asking 1,000 gives an " +
						"estimate that's much harder for a handful of outliers to move — which is " +
						"exactly why pollsters report a sample size and a margin of error together, " +
						"not just a percentage. Casinos and insurers rely on the same fact from the " +
						"other side: a casino doesn't need to win any single hand, only enough " +
						"hands that its built-in edge dominates the long-run average payout, and an " +
						"insurer prices policies assuming that across enough policyholders, actual " +
						"claims settle in near the predicted average rate. In machine learning, this " +
						"is why evaluating a model on a tiny test set is unreliable — its measured " +
						"accuracy can swing several points from noise alone — while a large test set " +
						"gives a number you can actually trust.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the average converges because the denominator n keeps " +
						"growing, diluting any run of heads or tails — never because any individual " +
						"flip's odds changed.' Not like this: the gambler's fallacy — treating three " +
						"heads in a row as making tails 'due' on the next flip. Flip 4 has no " +
						"memory of flips 1-3; its probability is p, full stop, every single time. " +
						"Also not like this: expecting the running average to get strictly closer " +
						"to p with every additional flip — this run's own n=100 point (60%, farther " +
						"from 50% than n=50's 52%) shows that isn't guaranteed at any particular " +
						"step; the law of large numbers only promises the long-run trend, not a " +
						"monotonic march toward p.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "p", Label: "True probability p", Min: 0.05, Max: 0.95, Step: 0.05, Def: 0.5},
			{Key: "n", Label: "Number of flips (n)", Min: 1, Max: 300, Step: 1, Def: 100},
		},
		Render: render,
	})
}

// hash01 turns an integer seed into a deterministic pseudo-random number in
// [0, 1) via a fixed irrational-multiplier trick. It's not a real random
// number generator -- it's a fixed, reproducible sequence that merely looks
// patternless, which is exactly what Render needs (pure: same input, same
// output, no time, no math/rand).
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// Outcome returns the result of the `trial`-th coin flip of a coin that
// comes up heads with true probability p: 1 for heads, 0 for tails.
// Deterministic in trial and p -- the same trial number always reproduces
// the same flip, so a whole run of flips can be recomputed exactly.
func Outcome(trial int, p float64) float64 {
	if hash01(trial) < p {
		return 1
	}
	return 0
}

// RunningAverages returns the running average of Outcome(1,p), Outcome(2,p),
// ..., Outcome(maxN,p) after each flip: result[i] is the average of the
// first i+1 flips. maxN < 1 is treated as 1.
func RunningAverages(maxN int, p float64) []float64 {
	if maxN < 1 {
		maxN = 1
	}
	out := make([]float64, maxN)
	sum := 0.0
	for i := 0; i < maxN; i++ {
		sum += Outcome(i+1, p)
		out[i] = sum / float64(i+1)
	}
	return out
}

// RunningAverage returns the average of the first n flips of a coin with
// true probability p -- the single value RunningAverages(n,p)'s last entry
// holds, computed directly for callers that only need one point.
func RunningAverage(n int, p float64) float64 {
	if n < 1 {
		n = 1
	}
	sum := 0.0
	for i := 1; i <= n; i++ {
		sum += Outcome(i, p)
	}
	return sum / float64(n)
}

// maxFlips is the fixed right edge of the curve's x-axis, matching the
// "Number of flips" slider's own Max so the whole slider range stays
// visible on the plot.
const maxFlips = 300

func render(params map[string]float64) string {
	trueP := params["p"]
	n := int(params["n"])
	if n < 1 {
		n = 1
	}

	series := RunningAverages(maxFlips, trueP)

	c := viz.New(680, 420, 0, maxFlips, 0, 1)
	c.Axes()
	for x := 0.0; x <= maxFlips; x += 50 {
		c.Tick(x, fmt.Sprintf("%.0f", x))
	}

	// The true probability, as a dashed reference line the running
	// average is converging toward.
	c.Path([][2]float64{{0, trueP}, {maxFlips, trueP}}, viz.Muted, 1)

	// The running-average curve itself: wide, noisy swings at small n,
	// settling toward the flat line as n grows.
	curve := make([][2]float64, maxFlips)
	for i, avg := range series {
		curve[i] = [2]float64{float64(i + 1), avg}
	}
	c.Path(curve, viz.Accent, 2)

	// The current n, highlighted.
	avg := series[n-1]
	px, py := c.X(float64(n)), c.Y(avg)
	c.Path([][2]float64{{float64(n), 0}, {float64(n), avg}}, viz.Warm, 1)
	c.Rect(px-4, py-4, 8, 8, viz.Warm, 1)

	c.Text(16, 24, fmt.Sprintf("true probability p = %.2f", trueP), 14, viz.Muted, "start")
	c.Text(16, 44, fmt.Sprintf("after n=%d flips: running average = %.3f", n, avg), 15, viz.Warm, "start")
	c.Text(16, 64, fmt.Sprintf("error |average - p| = %.3f", math.Abs(avg-trueP)), 13, viz.Ink, "start")

	return c.String()
}
