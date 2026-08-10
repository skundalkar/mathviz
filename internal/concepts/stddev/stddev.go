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
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Two students each score 10 points above their class average. Same number, " +
						"'10 points above average' — but do they deserve equal bragging rights? To " +
						"answer that you need one more fact about each class: how far scores " +
						"typically land from the average there, high or low. That's the piece the " +
						"raw '+10' is missing.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Once you have that missing fact, watch what it does to each student's result:",
					"• Class A: scores typically land only about 2 points from the average — a " +
						"tightly bunched class. Student A's +10 is 10÷2 = 5 times further from " +
						"average than it's normal to land here — off the charts.",
					"• Class B: scores typically land about 20 points from the average — a much " +
						"more spread-out class. Student B's +10 is only 10÷20 = 0.5, half of how " +
						"far it's completely normal to land there — unremarkable.",
					"That 'typically lands X points from average' number is standard deviation " +
						"(σ): take every score, measure how far it sits from the mean, and boil " +
						"those distances down to one typical distance. Class A's σ ≈ 2, class B's σ " +
						"≈ 20 — the exact numbers used above. Counting how many σ's away a score " +
						"sits gives it a z-score: student A's +10 is +5σ, student B's +10 is +0.5σ " +
						"— same raw gap, now comparable.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Drag σ and the bell curve gets fatter or skinnier while the shaded band — " +
						"always mean ± 1σ — keeps covering about 68% of the area no matter how wide " +
						"or narrow the curve is; ±2σ covers about 95%, ±3σ about 99.7%. Set σ small, " +
						"like class A's ≈2, and the curve looks tightly bunched; set σ large, like " +
						"class B's ≈20, and it spreads thin — the same two classes from the worked " +
						"example, now literally narrow or wide on the page.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"'10 points above average' is no longer just a raw number you have to " +
						"eyeball — you can convert it to a z-score and read off an actual rarity " +
						"from the 68/95/99.7 rule, so a claim like '5 sigma out' means the same " +
						"thing no matter what's being measured.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A 101°F fever is routine for some people and alarming for others depending " +
						"on their normal baseline and its variability; manufacturing tolerances are " +
						"set in sigmas, not raw units, because 'off by 2mm' means something " +
						"different for a bolt than for a bridge; and '1-in-20' scientific thresholds " +
						"(p < 0.05) are really a statement about how many sigmas out a result has to " +
						"land before it counts as surprising.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"'That was a 3-sigma event' (said after market crashes, unusual test scores, " +
						"quality-control failures) means this extreme should happen only about 3 " +
						"times in 1,000 if the usual pattern of variation held — the 99.7% rule made " +
						"concrete, not just a vague 'very rare.'",
					"Calling a fixed number of points or dollars 'a lot' or 'a little' on its own, " +
						"without asking relative to what, is the common mistake — the same 10-point " +
						"gap can be unremarkable in one dataset and the single most extreme value in " +
						"another.",
				},
			},
		},
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
