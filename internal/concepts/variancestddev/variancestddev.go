// Package variancestddev visualizes why we square deviations before averaging
// them, then take a square root to get back to the original units. Raw
// deviations from the mean always sum to zero (positive and negative cancel),
// so squaring is what lets them accumulate into a real measure of spread —
// variance. But variance is in squared units, so stddev = √variance converts
// it back to something comparable to the data itself.
package variancestddev

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

// baseSample is a small, fixed dataset (deviations from its own mean, roughly)
// that every rendered sample is derived from. Keeping it as a literal — rather
// than generating it — keeps the picture concrete: these are "quiz score"
// style deltas, not synthetic noise.
var baseSample = []float64{-4, -3, -1, 0, 1, 1, 2, 4, 6}

// Sample scales baseSample by spread (so its variance scales with spread²)
// and, if outlier is nonzero, pushes the last point further out — showing how
// much a single extreme value can inflate variance versus a robust measure.
// Pure function: same inputs always produce the same slice.
func Sample(spread, outlier float64) []float64 {
	xs := make([]float64, len(baseSample))
	for i, v := range baseSample {
		xs[i] = v * spread
	}
	if len(xs) > 0 {
		xs[len(xs)-1] += outlier
	}
	return xs
}

// Mean is the arithmetic mean of xs. Returns 0 for an empty slice.
func Mean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

// Variance is the population variance of xs: the mean of the squared
// deviations from the mean. Returns 0 for an empty slice.
func Variance(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	mean := Mean(xs)
	sum := 0.0
	for _, x := range xs {
		d := x - mean
		sum += d * d
	}
	return sum / float64(len(xs))
}

// StdDev is the population standard deviation of xs: √Variance(xs). Same
// units as the data, unlike Variance.
func StdDev(xs []float64) float64 {
	return math.Sqrt(Variance(xs))
}

func init() {
	concept.Register(concept.Concept{
		ID:    "variance-vs-stddev",
		Seq:   6,
		Title: "Variance vs. standard deviation",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Two dart throwers each take 4 throws. Measure how far each throw landed " +
						"from the bullseye, in inches — negative if left, positive if right. " +
						"Thrower 1: -2, -1, +1, +2. Thrower 2: -8, -4, +4, +8 — visibly about 4x " +
						"wilder just from looking at the numbers.",
					"Try the obvious way to turn that visible difference into one number: average " +
						"how far off each throw was. Thrower 1: (-2-1+1+2)/4 = 0. Thrower 2: " +
						"(-8-4+4+8)/4 = 0. Both average to exactly zero, even though one is clearly " +
						"far more scattered — deviations from the mean always cancel, every time, " +
						"for any data. So how do you actually turn 'more scattered' into a single " +
						"number?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Square each deviation before averaging it. Squaring erases the minus sign " +
						"(both -2 and +2 become 4), so nothing cancels — and a throw twice as far " +
						"off ends up counting four times as much, punishing wild misses harder than " +
						"close ones:",
					"• Thrower 1: squared deviations 4,1,1,4 sum to 10; average = 10/4 = 2.5 — " +
						"that's the variance, in square inches.",
					"• Thrower 2: squared deviations 64,16,16,64 sum to 160; average = 160/4 = " +
						"40 — variance, in square inches.",
					"Finally distinguishable: 40 is a lot bigger than 2.5. But 'square inches' " +
						"isn't a unit anyone uses to describe a dart thrower, so square-root back " +
						"down to undo exactly that: stddev 1.58 vs. 6.32, in real inches again — " +
						"the same 4x gap the raw numbers showed at a glance.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Each bar is one point's squared deviation from the mean (green if the point " +
						"sits above the mean, red if below) — the same numbers computed above; a " +
						"taller bar means the point is further from the mean. Drag spread to watch " +
						"stddev scale linearly while variance scales with the square (this sample " +
						"is literally Thrower 1's numbers × spread, so a 4x stretch is exactly the " +
						"Thrower 1 → Thrower 2 comparison above); drag the outlier slider to see " +
						"how much a single extreme point inflates both, since squaring punishes the " +
						"far-out point hardest of all.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You now have a single number that actually distinguishes a tight, consistent " +
						"set of throws from a wild one — something the 'obvious' average-of-" +
						"deviations approach couldn't do because it always lands on zero — plus a " +
						"version of that number (standard deviation) back in the data's own units " +
						"instead of squared ones.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Why one wild outlier can wreck a variance-based estimate more than you'd " +
						"guess (median/IQR are more robust because they don't square anything), " +
						"why error is usually reported as RMSE — root-mean-squared-error, 'square, " +
						"average, then square-root back' — instead of raw MSE, and why standard " +
						"deviation, not variance, is the number that actually gets quoted next to a " +
						"mean.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'that's a high-variance strategy' (common in startups, " +
						"poker, investing) means outcomes are spread wide — could go great, could " +
						"go badly, hard to call which.",
					"Not like this: 'high variance means it's worse' — variance measures spread, " +
						"not direction; a high-variance bet can have a better average outcome than " +
						"a safe one, just with far less certainty about any single try.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "spread", Label: "Spread", Min: 0.3, Max: 2.5, Step: 0.1, Def: 1},
			{Key: "outlier", Label: "Outlier push", Min: 0, Max: 10, Step: 0.5, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	spread, outlier := p["spread"], p["outlier"]
	if spread <= 0 {
		spread = 0.1
	}
	xs := Sample(spread, outlier)
	mean := Mean(xs)
	variance := Variance(xs)
	stddev := StdDev(xs)

	lo, hi := xs[0], xs[0]
	for _, x := range xs {
		if x < lo {
			lo = x
		}
		if x > hi {
			hi = x
		}
	}
	pad := (hi-lo)*0.15 + 1
	lo -= pad
	hi += pad

	sq := make([]float64, len(xs))
	maxSq := 0.0
	for i, x := range xs {
		d := x - mean
		sq[i] = d * d
		if sq[i] > maxSq {
			maxSq = sq[i]
		}
	}
	if maxSq == 0 {
		maxSq = 1
	}

	c := viz.New(680, 360, lo, hi, 0, maxSq*1.3)
	c.Axes()
	step := 2.0
	start := math.Ceil(lo/step) * step
	for x := start; x <= hi; x += step {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	// Squared-deviation bars: taller bar = point further from the mean. Bars
	// never go negative, unlike the raw signed deviations they replace.
	for i, x := range xs {
		px := c.X(x)
		top := c.Y(sq[i])
		bottom := c.Y(0)
		color := viz.Good
		if x < mean {
			color = viz.Bad
		}
		c.Rect(px-4, top, 8, bottom-top, color, 0.55)
	}

	c.VLine(mean, viz.Ink, false)

	// The data points themselves, as a row of markers above the tallest bar.
	markerY := maxSq * 1.16
	for _, x := range xs {
		px, py := c.X(x), c.Y(markerY)
		c.Rect(px-3, py-3, 6, 6, viz.Ink, 0.85)
	}

	c.Text(c.X(mean), c.PadT+12, "mean", 12, viz.Ink, "middle")
	c.Text(20, 24, fmt.Sprintf("mean = %.2f", mean), 14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("variance = %.2f (squared units — bar heights above)", variance),
		13, viz.Muted, "start")
	c.Text(20, 62, fmt.Sprintf("stddev = √variance = %.2f (back to the data's own units)", stddev),
		13, viz.Muted, "start")

	return c.String()
}
