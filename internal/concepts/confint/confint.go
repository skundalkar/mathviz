// Package confint visualizes what "95% confidence" actually means: not that
// there's a 95% chance the true mean lies in one particular interval, but
// that if you repeated the experiment many times, about 95% of the intervals
// you'd build would contain the true mean. To make that concrete without any
// randomness, the picture draws a fixed set of 20 hypothetical sample means
// placed at evenly spaced quantiles of the sampling distribution — a
// deterministic stand-in for "20 repeated experiments" — and shows which of
// their confidence intervals happen to capture the true mean.
package confint

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

// NumSamples is how many hypothetical repeated experiments the picture
// draws — a fixed stand-in for "if you did this many times".
const NumSamples = 20

// PopulationSigma is the (known) population standard deviation used to build
// the sampling distribution. TrueMean is where the population is centered;
// the whole point of the picture is that a real experimenter never sees it.
const (
	PopulationSigma = 1.0
	TrueMean        = 0.0
)

func init() {
	concept.Register(concept.Concept{
		ID:    "confidence-interval",
		Seq:   9,
		Title: "Confidence interval",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You want the average commute time across an entire city of 100,000 commuters " +
						"— the true population mean — but you can't measure everyone, so you survey a " +
						"random n=64 of them. Their average comes out to exactly 32 minutes, no doubt " +
						"about it — it's the arithmetic average of the 64 people you actually asked. " +
						"But 32 minutes isn't the number you actually wanted: it's an estimate of a " +
						"different, bigger number — the true citywide average — that you'll never get " +
						"to measure directly. How do you honestly say how close your one estimate " +
						"probably is to a number you can never check?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Two steps turn your sample's numbers into an interval that answers that " +
						"question. Individual commute times vary by standard deviation s=8 minutes " +
						"among PEOPLE — that's spread among people, not uncertainty about the average:",
					"• Standard error — how much the SAMPLE MEAN itself would wobble if you " +
						"resurveyed a different random 64 people: SE = s/sqrt(n) = 8/sqrt(64) = 1 " +
						"minute.",
					"• Margin of error — scale that wobble by a multiplier for your confidence level " +
						"(z=1.96 for 95%, the width in standard errors that captures the middle 95% " +
						"of a normal curve): margin = 1.96x1 ~= 2 minutes.",
					"• Interval = 32 +/- 2, about [30, 34] — 'we're 95% confident the true citywide " +
						"average is between 30 and 34 minutes.'",
					"32 itself isn't uncertain — it's the exact average of your 64 people, the way a " +
						"thrown dart's landing spot is never in question once it's stuck in the board. " +
						"The uncertainty is about a different, hidden target: the true citywide " +
						"average, which your dart only lands near, not on. That's why the interval is " +
						"built from standard error (1 minute — how much the sample MEAN wobbles) and " +
						"never the raw standard deviation (8 minutes — how spread individual commutes " +
						"are): stddev doesn't shrink as n grows, but standard error does. Quadruple n " +
						"to 256 and standard error only drops to 0.5 — 4x the data buys 2x the " +
						"precision (1/sqrt(n)).",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The dashed line is the true mean — the number a real experimenter, like our " +
						"city-commute surveyor, never gets to see. Each of the 20 rows is one " +
						"hypothetical repeat of the same 64-person survey: a fresh sample, its own " +
						"mean (the square marker), its own net cast around that mean — green if the " +
						"net captured the true mean, red if it missed. Raise the confidence slider and " +
						"every interval widens like casting a bigger net, turning red rows green — the " +
						"same 1.96 multiplier from the worked example above growing for a higher " +
						"confidence level. Raise the sample size slider and every interval narrows " +
						"around its own mean, the same 1/sqrt(n) shrinkage that took standard error " +
						"from 1 minute at n=64 to 0.5 minutes at n=256 — narrower nets, not more of " +
						"them landing on target: coverage is the confidence level's job, not sample " +
						"size's.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can now honestly report a range for a number you can never directly measure " +
						"— 'we're 95% confident the true citywide average commute is between 30 and 34 " +
						"minutes' — and say correctly what that 95% actually promises: not that this " +
						"one interval has a 95% chance of being right, but that the method, repeated " +
						"many times, brackets the truth about 95% of the time.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"'This poll has a ±3 point margin of error at 95% confidence' means 95% of polls " +
						"run this way would bracket the true value — not that this specific poll has a " +
						"95% chance of being right. It's also why a narrower interval from a bigger " +
						"sample is a genuine improvement (a tighter net, same success rate), while " +
						"cherry-picking a higher confidence level just to get a cleaner-looking " +
						"one-off result buys you nothing on the one survey that matters to you.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'This poll has a 3-point margin of error at 95% confidence' " +
						"means the polling METHOD, repeated many times, would bracket the true number " +
						"about 95% of the time.",
					"Not like this: 'There's a 95% chance the true value is in this specific " +
						"interval' — once the interval is drawn, it either contains the true value or " +
						"it doesn't; the 95% describes the method's long-run success rate, not this " +
						"one instance.",
					"Also not like this: assuming the guarantee holds no matter how the 64 people " +
						"were chosen — sampling bias (e.g. only surveying one train platform) breaks " +
						"it in a way no formula can detect or fix; the confidence interval can be tight " +
						"and confidently wrong at the same time.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "confidence", Label: "Confidence level", Min: 50, Max: 99, Step: 1, Def: 90, Unit: "%"},
			{Key: "n", Label: "Sample size (n)", Min: 2, Max: 100, Step: 1, Def: 20},
		},
		Render: render,
	})
}

// StdNormalQuantile is the inverse standard-normal CDF: the z such that
// P(Z <= z) = p, for 0 < p < 1. This is what turns an evenly spaced grid of
// probabilities into an evenly spaced set of "hypothetical" draws from a
// normal sampling distribution, with no randomness needed.
func StdNormalQuantile(p float64) float64 {
	if p <= 0 || p >= 1 {
		return math.NaN()
	}
	return math.Sqrt2 * math.Erfinv(2*p-1)
}

// CriticalZ returns the z-value such that a standard normal variable falls
// within [-z, z] with probability `confidence` (0 < confidence < 1) — the
// familiar 1.96 for confidence=0.95.
func CriticalZ(confidence float64) float64 {
	return StdNormalQuantile(0.5 + confidence/2)
}

// StandardError is the standard deviation of the sample mean for samples of
// size n drawn from a population with standard deviation sigma: σ/√n.
func StandardError(sigma float64, n int) float64 {
	if n < 1 {
		return 0
	}
	return sigma / math.Sqrt(float64(n))
}

// SampleMeans returns NumSamples deterministic stand-ins for "sample means
// from repeated experiments": the sampling distribution N(TrueMean, se²)
// evaluated at NumSamples evenly spaced quantiles ((i-0.5)/NumSamples). This
// is exact and reproducible, unlike drawing NumSamples random samples would
// be, while still spreading out the way real repeated sampling would.
func SampleMeans(n int) []float64 {
	se := StandardError(PopulationSigma, n)
	means := make([]float64, NumSamples)
	for i := 1; i <= NumSamples; i++ {
		q := (float64(i) - 0.5) / float64(NumSamples)
		means[i-1] = TrueMean + se*StdNormalQuantile(q)
	}
	return means
}

// Interval returns the confidence interval [lo, hi] built around a sample
// mean with critical value z and standard error se: mean ± z·se.
func Interval(mean, z, se float64) (lo, hi float64) {
	return mean - z*se, mean + z*se
}

// CoverageCount reports how many of the NumSamples deterministic intervals
// (at the given confidence level) contain TrueMean. Because SampleMeans is
// n-invariant relative to se (scaling se scales both the mean's offset and
// the interval half-width by the same factor), CoverageCount does not depend
// on n — only on how the evenly-spaced quantile grid lines up with the
// confidence band, which is exactly the frequentist guarantee in miniature.
func CoverageCount(confidence float64, n int) int {
	se := StandardError(PopulationSigma, n)
	z := CriticalZ(confidence)
	count := 0
	for _, m := range SampleMeans(n) {
		lo, hi := Interval(m, z, se)
		if lo <= TrueMean && TrueMean <= hi {
			count++
		}
	}
	return count
}

func render(p map[string]float64) string {
	confidence := p["confidence"] / 100
	if confidence <= 0 {
		confidence = 0.5
	}
	if confidence >= 1 {
		confidence = 0.99
	}
	n := int(p["n"] + 0.5)
	if n < 1 {
		n = 1
	}

	se := StandardError(PopulationSigma, n)
	z := CriticalZ(confidence)
	means := SampleMeans(n)

	// Symmetric x window sized to the widest interval, with a little margin.
	xlim := 0.0
	for _, m := range means {
		lo, hi := Interval(m, z, se)
		if -lo > xlim {
			xlim = -lo
		}
		if hi > xlim {
			xlim = hi
		}
	}
	xlim *= 1.15

	c := viz.New(680, 480, -xlim, xlim, 0, NumSamples+1)
	// Extra headroom for four lines of header/legend text, kept entirely
	// above the plot so it never crowds the top row of intervals.
	c.PadT = 96
	c.Axes()
	step := xlim / 4
	for x := -xlim + step; x < xlim; x += step {
		label := x
		if math.Abs(label) < 1e-9 {
			label = 0
		}
		c.Tick(x, fmt.Sprintf("%.2f", label))
	}

	// True mean, dashed — the thing a real experimenter never gets to see.
	// It never moves row to row: it's one fixed, hidden target every
	// hypothetical repeat is aiming at.
	c.VLine(TrueMean, viz.Ink, true)

	covered := 0
	for i, m := range means {
		y := float64(NumSamples - i) // top row = first sample, for reading order
		lo, hi := Interval(m, z, se)
		hit := lo <= TrueMean && TrueMean <= hi
		color := viz.Bad
		if hit {
			color = viz.Good
			covered++
		}
		c.Path([][2]float64{{lo, y}, {hi, y}}, color, 2.5)
		mx, my := c.X(m), c.Y(y)
		c.Rect(mx-2, my-2, 4, 4, color, 1)
	}

	c.Text(20, 24, fmt.Sprintf("confidence = %.0f%%    n = %d    SE = %.3f    z = %.2f",
		confidence*100, n, se, z), 13, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("%d of %d intervals capture the true mean (dashed line) — "+
		"close to the %.0f%% you'd expect over many repeats", covered, NumSamples, confidence*100),
		12, viz.Muted, "start")
	c.Text(20, 64, "each row = one hypothetical repeat: a fresh sample, its own mean, its own net "+
		"cast around that mean", 12, viz.Muted, "start")
	c.Text(20, 84, "■ = a sample's mean (a different draw every row)    dashed | = true mean "+
		"(one fixed target, hidden in real life)", 12, viz.Muted, "start")

	return c.String()
}
