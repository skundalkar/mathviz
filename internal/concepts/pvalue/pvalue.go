// Package pvalue visualizes a p-value as literally what it is: the area under
// a "null" distribution (what test statistics would look like if there were
// no real effect) that lies as far or farther from zero than the statistic
// you actually observed. Small shaded area = the observed value would be
// surprising under the null = small p-value. The picture also marks the
// critical boundary implied by a chosen significance level α, so you can see
// directly whether the shaded p-value area is smaller than α (reject H0) or
// not (fail to reject).
package pvalue

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "p-value",
		Title: "P-value",
		Blurb: "The curve is the null distribution: what a test statistic would look like " +
			"if there were truly no effect. The vertical line is the statistic you actually " +
			"observed; the shaded area is everything at least as extreme as it in either " +
			"direction — that shaded area IS the p-value. It answers one narrow question: " +
			"'if the null hypothesis were true, how often would we see a result this extreme " +
			"just from noise?' The dashed lines mark the critical boundary for your chosen " +
			"significance level α — if the shaded area is smaller than α, the observed " +
			"statistic falls outside the boundary and the result is 'statistically significant'.",
		Params: []concept.ParamSpec{
			{Key: "z", Label: "Observed statistic (z)", Min: -4, Max: 4, Step: 0.1, Def: 2.0},
			{Key: "alpha", Label: "Significance level (α)", Min: 0.01, Max: 0.20, Step: 0.01, Def: 0.05},
		},
		Render: render,
	})
}

// StdNormalPDF is the density of the null distribution: a standard normal
// N(0,1), representing what a test statistic looks like under "no effect".
func StdNormalPDF(x float64) float64 {
	return math.Exp(-0.5*x*x) / math.Sqrt(2*math.Pi)
}

// PValue is the two-sided p-value for an observed statistic zObs under the
// standard normal null distribution: the probability of seeing something at
// least as extreme (in either direction) purely by chance.
//
//	p = P(|Z| >= |zObs|) = erfc(|zObs| / √2)
func PValue(zObs float64) float64 {
	return math.Erfc(math.Abs(zObs) / math.Sqrt2)
}

// CriticalZ returns the two-sided critical value z* such that a standard
// normal falls outside [-z*, z*] with probability exactly alpha — the
// boundary an observed statistic must cross for the result to count as
// "statistically significant" at that significance level.
func CriticalZ(alpha float64) float64 {
	return math.Sqrt2 * math.Erfinv(1-alpha)
}

// Significant reports whether a p-value is small enough to reject the null
// hypothesis at significance level alpha.
func Significant(pValue, alpha float64) bool {
	return pValue < alpha
}

func render(p map[string]float64) string {
	return viz.New(680, 340, -4, 4, 0, 1).String()
}
