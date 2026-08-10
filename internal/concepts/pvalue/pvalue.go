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
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Your friend hands you a coin and claims it's fair. You flip it 10 times and " +
						"get 8 heads. Suspicious — but is it PROOF the coin is rigged, or could a fair " +
						"coin just get lucky?",
					"A fair coin doesn't land on exactly 5 heads every single time either; sometimes " +
						"it runs hot by pure chance. So the real question isn't whether 8 heads sounds " +
						"like a lot in isolation — it's: if the coin really were fair, how often would " +
						"10 flips produce a result this lopsided just from ordinary luck?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Flip a fair coin 10 times and there are 2^10=1,024 equally likely head/tail " +
						"sequences. Count how many land on 8 or more heads:",
					"• 45 sequences give exactly 8 heads, 10 give exactly 9, and 1 gives all 10 — " +
						"56 out of 1,024.",
					"• By symmetry, another 56 sequences land on 8 or more tails.",
					"• That's 112 of 1,024 sequences this lopsided or worse in either direction: " +
						"112/1024 = 10.9%, about 11% — no simulation, just counting equally likely " +
						"outcomes.",
					"That 11% is the p-value: start from a 'null' distribution — what results would " +
						"look like if nothing suspicious were going on — run the actual test, and ask " +
						"how much of that distribution's probability lies at least as far out as what " +
						"you actually observed, in either direction. That's the whole definition, " +
						"nothing more.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The picture generalizes the coin-counting idea to a continuous test statistic " +
						"instead of discrete sequences. The curve is that same 'null' distribution: " +
						"what a test statistic looks like if there's truly no effect. The vertical " +
						"line is what you actually observed (z); the shaded area — everything at " +
						"least as extreme — IS the p-value, the same idea as counting 112 of 1,024 " +
						"sequences, just computed as an area instead of a count. The dashed lines mark " +
						"the critical boundary implied by your chosen significance level α: push the " +
						"observed statistic past a dashed line and the shaded area drops below α, " +
						"flipping the verdict from 'fail to reject H₀' to 'reject H₀'.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can now put a number on how surprising a result would be if nothing real " +
						"were going on, instead of eyeballing whether '8 heads' feels suspicious. " +
						"Picking α in advance turns that number into a repeatable, agreed-upon " +
						"yes/no call — 'statistically significant' — rather than a post-hoc gut " +
						"decision about whether the result looks impressive.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Clinical trials, A/B tests, and scientific studies all lean on this same logic " +
						"to decide whether an observed effect is worth taking seriously or could " +
						"plausibly be noise. It's also why running enough different tests, subgroups, " +
						"or thresholds against the same data — 'p-hacking' — is a real problem: ask " +
						"enough different questions of the same data and a suspicious-looking result " +
						"shows up eventually even when nothing real is going on.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the result was statistically significant' means that if " +
						"nothing real were going on, a result this extreme would have been unlikely " +
						"to show up by chance alone (below an agreed-on threshold, usually 5%) — " +
						"nothing more, and nothing about how big or important the effect is.",
					"Not like this: 'p < 0.05, so there's a 95% chance the effect is real' — a " +
						"p-value is not the probability your hypothesis is true; it's the probability " +
						"of seeing data this extreme if your hypothesis were false.",
				},
			},
		},
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
