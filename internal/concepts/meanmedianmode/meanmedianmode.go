// Package meanmedianmode visualizes how a skewed distribution pulls the mean,
// median, and mode apart. All three coincide for a symmetric distribution, but
// as skew increases they peel off in a fixed order: mode stays at the peak,
// the median sits at the halfway point, and the mean gets dragged toward the
// long tail. A log-normal curve is used because all three statistics have
// simple closed forms, so the picture is exact rather than sampled.
package meanmedianmode

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "mean-median-mode",
		Title: "Mean, median & mode",
		Blurb: "For a symmetric distribution, the mean, median, and mode all land in the " +
			"same place. Skew the distribution and they peel apart in a fixed order: the mode " +
			"stays at the peak, the median sits at the 50/50 point, and the mean gets dragged " +
			"toward the long tail — because the mean, unlike the median, is sensitive to extreme " +
			"values. Drag the skew slider and watch mean > median > mode open up as the right " +
			"tail lengthens.",
		Params: []concept.ParamSpec{
			{Key: "sigma", Label: "Skew (σ)", Min: 0.05, Max: 1.2, Step: 0.05, Def: 0.6},
			{Key: "mu", Label: "Log-location (μ)", Min: -0.5, Max: 0.5, Step: 0.05, Def: 0},
		},
		Render: render,
	})
}

// LogNormalPDF is the probability density of a log-normal distribution with
// log-space parameters mu, sigma. Pure math, unit-tested below.
func LogNormalPDF(x, mu, sigma float64) float64 {
	if x <= 0 || sigma <= 0 {
		return 0
	}
	z := (math.Log(x) - mu) / sigma
	return math.Exp(-0.5*z*z) / (x * sigma * math.Sqrt(2*math.Pi))
}

// Mean is the closed-form mean of a log-normal(mu, sigma) distribution.
func Mean(mu, sigma float64) float64 {
	return math.Exp(mu + sigma*sigma/2)
}

// Median is the closed-form median of a log-normal(mu, sigma) distribution.
func Median(mu, sigma float64) float64 {
	return math.Exp(mu)
}

// Mode is the closed-form mode (location of peak density) of a
// log-normal(mu, sigma) distribution.
func Mode(mu, sigma float64) float64 {
	return math.Exp(mu - sigma*sigma)
}

func render(p map[string]float64) string {
	mu, sigma := p["mu"], p["sigma"]
	if sigma <= 0 {
		sigma = 0.05
	}

	mean, median, mode := Mean(mu, sigma), Median(mu, sigma), Mode(mu, sigma)

	xmax := math.Exp(mu+3.2*sigma) * 1.1
	peak := LogNormalPDF(mode, mu, sigma)
	c := viz.New(680, 340, 0, xmax, 0, peak*1.15)
	c.Axes()

	step := xmax / 8
	for x := step; x <= xmax; x += step {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	curve := viz.Sample(0.001, xmax, 300, func(x float64) float64 {
		return LogNormalPDF(x, mu, sigma)
	})
	c.Path(curve, viz.Accent, 2.5)

	c.VLine(mode, viz.Good, true)
	c.VLine(median, viz.Warm, true)
	c.VLine(mean, viz.Bad, false)

	c.Text(c.X(mode), c.PadT+12, "mode", 12, viz.Good, "middle")
	c.Text(c.X(median), c.PadT+28, "median", 12, viz.Warm, "middle")
	c.Text(c.X(mean), c.PadT+44, "mean", 12, viz.Bad, "middle")

	c.Text(20, 24, fmt.Sprintf("mode = %.2f    median = %.2f    mean = %.2f", mode, median, mean),
		14, viz.Ink, "start")
	c.Text(20, 44, "right-skewed: the long tail drags the mean past the median, past the mode",
		12, viz.Muted, "start")

	return c.String()
}
