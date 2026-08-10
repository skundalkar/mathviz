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
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "p-value",
		Seq:   10,
		Title: "P-value",
		Blurb: "Your friend hands you a coin and claims it's fair. You flip it 10 times and " +
			"get 8 heads. Suspicious — but is it PROOF the coin is rigged, or could a fair " +
			"coin just get lucky? There are 2^10=1,024 equally likely head/tail sequences from " +
			"10 flips. Count how many land on 8+ heads: 45+10+1=56. By symmetry, another 56 " +
			"land on 8+ tails. That's 112 of 1,024 sequences this lopsided or worse in either " +
			"direction: 112/1024 = 10.9%, about 11% — no simulation, just counting equally " +
			"likely outcomes. That 11% is the p-value. The curve is that 'null' distribution: " +
			"what a test statistic looks like if there's truly no effect. The vertical line is " +
			"what you actually observed; the shaded area — everything at least as extreme — IS " +
			"the p-value. The dashed lines mark the critical boundary for your chosen " +
			"significance level α — cross it and the result is 'statistically significant'.",
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
	z := p["z"]
	alpha := p["alpha"]
	if alpha <= 0 {
		alpha = 0.01
	}
	if alpha >= 1 {
		alpha = 0.99
	}

	pval := PValue(z)
	zCrit := CriticalZ(alpha)
	sig := Significant(pval, alpha)
	az := math.Abs(z)

	const xmin, xmax = -4.0, 4.0
	curve := viz.Sample(xmin, xmax, 300, StdNormalPDF)
	peak := StdNormalPDF(0)

	c := viz.New(680, 340, xmin, xmax, 0, peak*1.2)
	c.Axes()
	for x := -4.0; x <= 4.0; x += 1 {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	// Shade the p-value: everything at least as extreme as the observed |z|,
	// in both tails.
	c.Area(curve, xmin, -az, viz.Warm, 0.35)
	c.Area(curve, az, xmax, viz.Warm, 0.35)
	c.Path(curve, viz.Ink, 2)

	// Critical boundary implied by alpha, dashed — cross it and the shaded
	// area (the p-value) has dropped below alpha.
	c.VLine(zCrit, viz.Muted, true)
	c.VLine(-zCrit, viz.Muted, true)

	// Observed statistic, solid.
	c.VLine(z, viz.Accent, false)

	verdict, cmp, vColor := "fail to reject H₀", "≥", viz.Muted
	if sig {
		verdict, cmp, vColor = "reject H₀", "<", viz.Good
	}

	c.Text(20, 24, fmt.Sprintf("z = %.2f    p = %.3f    α = %.2f", z, pval, alpha), 14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("p %s α → %s", cmp, verdict), 13, vColor, "start")
	c.Text(20, 62, "dashed: critical boundary from α   solid vertical: observed z",
		12, viz.Muted, "start")

	return c.String()
}
