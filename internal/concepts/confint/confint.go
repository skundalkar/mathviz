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
		Title: "Confidence interval",
		Blurb: "A fish hides somewhere in a lake, sitting perfectly still. You get an " +
			"imperfect reading of roughly where it is and cast a net centered on that reading. " +
			"A wide net catches it almost every time, even with a rough guess; a narrow net " +
			"only succeeds when your guess lands close. Once you've thrown one specific net, " +
			"it's tempting to say 'there's a 95% chance the fish is in this net' — but the " +
			"fish never moved and the net already landed; it either caught it or it didn't. " +
			"What 95% actually describes is the METHOD: repeat this cast-around-your-guess " +
			"process many times and about 95% of the nets you throw capture the fish. The " +
			"dashed line here is the true mean, which in real life you never get to see. Each " +
			"row is one hypothetical 'cast' — green if its interval captured the true mean, " +
			"red if it missed. Raise the confidence level and every interval widens (fewer " +
			"misses); raise the sample size and every interval narrows (more precision).",
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

	c := viz.New(680, 460, -xlim, xlim, 0, NumSamples+1)
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

	return c.String()
}
