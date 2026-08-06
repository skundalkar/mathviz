// Package normalskew visualizes skewness and (excess) kurtosis as knobs that
// reshape a standard normal curve. It uses the Gram-Charlier A series, a
// closed-form correction to the normal density built from Hermite
// polynomials: a small skew or kurtosis term tilts or reshapes the bell curve
// while leaving it centered and (to first order) still normalized.
package normalskew

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "normal-vs-skew",
		Title: "Normal vs. skewed & fat-tailed",
		Blurb: "A standard normal curve is symmetric with zero excess kurtosis: its skewness " +
			"and kurtosis knobs are both at zero. Turn the skew knob and the curve leans — one " +
			"tail lengthens while the other shortens, exactly what a real-world skewed dataset " +
			"looks like. Turn the kurtosis knob and the curve stays symmetric but its peak " +
			"sharpens and its tails fatten (positive excess kurtosis, 'leptokurtic') or the " +
			"opposite (negative, 'platykurtic'). The dashed curve is the plain normal for " +
			"comparison.",
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
	return viz.New(680, 340, -5, 5, 0, 1).String()
}
