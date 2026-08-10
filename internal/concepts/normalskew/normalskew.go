// Package normalskew visualizes skewness and (excess) kurtosis as knobs that
// reshape a standard normal curve. It uses the Gram-Charlier A series, a
// closed-form correction to the normal density built from Hermite
// polynomials: a small skew or kurtosis term tilts or reshapes the bell curve
// while leaving it centered and (to first order) still normalized.
package normalskew

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "normal-vs-skew",
		Seq:   7,
		Title: "Normal vs. skewed & fat-tailed",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Two neighborhoods report the exact same average home price and the exact same " +
						"standard deviation — interchangeable, right? With real numbers: Neighborhood " +
						"A, 5 homes at $380k/$390k/$400k/$410k/$420k — mean $400k, tight and symmetric. " +
						"Neighborhood B: 4 homes at $350k plus one $600k mansion — mean = " +
						"(4x$350k+$600k)/5 = $400k, exactly matching A, even though most of B's " +
						"residents live well below their own neighborhood's 'average.' Mean and " +
						"standard deviation alone can't tell these two neighborhoods apart — so what " +
						"number would?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Skewness is that third number: it measures which way a distribution leans. " +
						"Neighborhood B has positive skew — a long right tail from the one mansion " +
						"dragging the mean up above where most of the data (the four $350k homes) " +
						"actually sits; turn the skew knob and one tail lengthens while the other " +
						"shortens. Kurtosis is a fourth number, answering a different question: are " +
						"extreme, rare outcomes more likely than a plain bell curve predicts ('fat " +
						"tails') — a lesson risk models learned expensively in 2008, when markets " +
						"assumed roughly normal returns and got blindsided by a crash a thin-tailed " +
						"model treated as nearly impossible. The picture builds both effects on top of " +
						"a plain normal curve with a closed-form correction (the Gram-Charlier series): " +
						"at skew = kurtosis = 0 it's exactly the standard normal, and the dashed curve " +
						"stays that plain normal throughout for comparison.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Drag the skewness slider positive and the curve tilts the way Neighborhood B's " +
						"home prices do — one tail stretches out while the other pulls in, even though " +
						"the peak stays centered. Drag the kurtosis slider positive and the curve " +
						"sharpens its peak and fattens its tails — the 2008-style shape where extremes " +
						"are more common than the dashed normal reference predicts; negative kurtosis " +
						"flattens it the other way. Both knobs at zero reproduce the dashed reference " +
						"curve exactly.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can now tell apart two distributions that mean and standard deviation report " +
						"as identical — like Neighborhood A and B, or a stable market and one about to " +
						"crash — by checking the two numbers that actually capture shape instead of " +
						"stopping at 'same average, same spread, must be the same.'",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Financial risk models that assume normality and get blindsided by fat-tailed " +
						"crashes, A/B test metrics that violate the normality assumptions baked into " +
						"some significance tests, and any dashboard that reports only 'mean ± stddev' " +
						"for data that's secretly lopsided or spiky — two numbers that can hide a very " +
						"different-shaped reality.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'That sample is skewed' means the data leans hard toward one " +
						"side — a few extreme values on one end are pulling the average away from " +
						"where most of the data actually sits, even if the 'typical' spread (stddev) " +
						"looks perfectly ordinary.",
					"Not like this: 'The average looks fine, so the data's probably normal' — mean " +
						"and standard deviation alone can't tell you the shape; two very " +
						"differently-shaped distributions can share both numbers exactly.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "skew", Label: "Skewness", Min: -0.9, Max: 0.9, Step: 0.1, Def: 0},
			{Key: "kurt", Label: "Excess kurtosis", Min: -1, Max: 2, Step: 0.1, Def: 0},
		},
		Render: render,
	})
}

// StdNormalPDF is the density of the standard normal distribution N(0,1).
func StdNormalPDF(x float64) float64 {
	return math.Exp(-0.5*x*x) / math.Sqrt(2*math.Pi)
}

// Hermite3 is the third probabilist's Hermite polynomial, He3(x) = x³ - 3x.
func Hermite3(x float64) float64 {
	return x*x*x - 3*x
}

// Hermite4 is the fourth probabilist's Hermite polynomial,
// He4(x) = x⁴ - 6x² + 3.
func Hermite4(x float64) float64 {
	return x*x*x*x - 6*x*x + 3
}

// GramCharlierPDF approximates a distribution with the given skewness and
// excess kurtosis by correcting the standard normal density with a
// Hermite-polynomial expansion (the Gram-Charlier A series):
//
//	f(x) = φ(x) * (1 + skew/6 * He3(x) + kurt/24 * He4(x))
//
// This is exact to first order in skew and kurt and integrates to 1 by
// construction (each Hermite term integrates to zero against φ). For large
// |skew| or |kurt| the correction can push the density negative in the
// tails, which is a known limitation of the expansion; the result is clamped
// to 0 so it always stays a valid (non-negative) curve.
func GramCharlierPDF(x, skew, kurt float64) float64 {
	adj := 1 + skew/6*Hermite3(x) + kurt/24*Hermite4(x)
	y := StdNormalPDF(x) * adj
	if y < 0 {
		return 0
	}
	return y
}

func render(p map[string]float64) string {
	skew, kurt := p["skew"], p["kurt"]

	const xmin, xmax = -5.0, 5.0
	curve := viz.Sample(xmin, xmax, 300, func(x float64) float64 {
		return GramCharlierPDF(x, skew, kurt)
	})
	refCurve := viz.Sample(xmin, xmax, 300, StdNormalPDF)

	peak := 0.0
	for _, pt := range curve {
		if pt[1] > peak {
			peak = pt[1]
		}
	}
	if peak <= 0 {
		peak = StdNormalPDF(0)
	}

	c := viz.New(680, 340, xmin, xmax, 0, peak*1.2)
	c.Axes()
	for x := -4.0; x <= 4.0; x += 2 {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	// Reference standard normal, thin and muted, for visual contrast.
	c.Path(refCurve, viz.Muted, 1.5)
	c.Path(curve, viz.Accent, 2.5)

	c.Text(20, 24, fmt.Sprintf("skewness = %.1f    excess kurtosis = %.1f", skew, kurt),
		14, viz.Ink, "start")
	c.Text(20, 44, "skew > 0: right tail heavier   skew < 0: left tail heavier",
		12, viz.Muted, "start")
	c.Text(20, 62, "kurt > 0: sharper peak, fatter tails   kurt < 0: flatter, thinner tails",
		12, viz.Muted, "start")
	c.Text(c.W-18, c.PadT+14, "— standard normal", 12, viz.Muted, "end")

	return c.String()
}
