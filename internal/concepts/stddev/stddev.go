// Package stddev visualizes what standard deviation (σ) actually measures:
// the width of the "typical spread" around the mean. Drag σ and the bell curve
// gets fatter or skinnier; the shaded band is always mean ± 1σ and always holds
// ~68% of the area — that invariance is the whole lesson.
package stddev

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "standard-deviation",
		Seq:   1,
		Title: "Standard deviation (σ)",
		Blurb: "Two students each score '10 points above average' — equally impressive? Depends " +
			"how far scores typically land from average in each class. In class A, scores " +
			"typically land only 2 points from average, so +10 is 5 times further out than normal " +
			"— off the charts. In class B, scores typically land 20 points from average, so +10 is " +
			"just half of how far it's normal to land there — unremarkable. Standard deviation (σ) " +
			"IS that 'typically lands X points from average' number: measure how far every score " +
			"sits from the mean, boil those distances down to one typical distance, and that's σ. " +
			"For a bell-shaped curve it comes with a fixed rule too: mean ± 1σ always covers about " +
			"68% of everyone, ±2σ about 95%, ±3σ about 99.7% — narrow curve or wide, doesn't " +
			"matter. Drag σ and watch the band stretch or shrink while always holding that same " +
			"68% — that's what makes '2 sigma above average' a portable, comparable rarity instead " +
			"of just a raw number.",
		Params: []concept.ParamSpec{
			{Key: "mu", Label: "Mean (μ)", Min: -4, Max: 4, Step: 0.1, Def: 0},
			{Key: "sigma", Label: "Std dev (σ)", Min: 0.4, Max: 3, Step: 0.05, Def: 1, Unit: "σ"},
		},
		Render: render,
	})
}

// NormalPDF is the probability density of a normal distribution — pure math,
// unit-tested below.
func NormalPDF(x, mu, sigma float64) float64 {
	if sigma <= 0 {
		return 0
	}
	z := (x - mu) / sigma
	return math.Exp(-0.5*z*z) / (sigma * math.Sqrt(2*math.Pi))
}

// AreaWithin returns the fraction of the distribution within k standard
// deviations of the mean, i.e. erf(k/√2). For k=1 this is ~0.6827.
func AreaWithin(k float64) float64 {
	return math.Erf(k / math.Sqrt2)
}

func render(p map[string]float64) string {
	mu, sigma := p["mu"], p["sigma"]
	if sigma <= 0 {
		sigma = 0.4
	}

	// Fixed x window so the curve visibly moves/resizes as μ,σ change.
	const xmin, xmax = -8.0, 8.0
	peak := NormalPDF(mu, mu, sigma) // max density is at the mean
	c := viz.New(680, 320, xmin, xmax, 0, peak*1.15)
	c.Axes()

	// x ticks every 2 units.
	for x := -8.0; x <= 8.0; x += 2 {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	curve := viz.Sample(xmin, xmax, 240, func(x float64) float64 {
		return NormalPDF(x, mu, sigma)
	})

	// Shade mean ± 1σ.
	c.Area(curve, mu-sigma, mu+sigma, viz.Accent, 0.18)
	c.Path(curve, viz.Accent, 2.5)

	// Mean line + σ boundary lines.
	c.VLine(mu, viz.Ink, false)
	c.VLine(mu-sigma, viz.Muted, true)
	c.VLine(mu+sigma, viz.Muted, true)

	pct := AreaWithin(1) * 100
	c.Text(20, 24, fmt.Sprintf("μ = %.2f    σ = %.2f", mu, sigma), 14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("shaded band (μ ± 1σ) ≈ %.1f%% of the data", pct), 13, viz.Muted, "start")
	c.Text(c.X(mu), c.PadT+12, "μ", 13, viz.Ink, "middle")

	return c.String()
}
