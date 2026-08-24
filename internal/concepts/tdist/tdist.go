// Package tdist visualizes Student's t-distribution: what happens to
// confidence-interval's z-critical-value method once you admit you don't
// actually know σ, only a sample estimate s of it. Small samples make s a
// noisy stand-in for σ, and the t-distribution's fatter tails (controlled by
// degrees of freedom) are the correction for that extra uncertainty — as the
// sample grows, the correction shrinks and t converges to the normal exactly.
package tdist

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "t-distribution",
		Seq:   63,
		Title: "Student's t-distribution",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`confidence-interval` built its margin of error as z* × σ/√n — but that " +
						"formula quietly assumes you already know σ, the true population standard " +
						"deviation, exactly. In practice you almost never do. Say you test a new " +
						"fertilizer on 6 plants and measure their height increase: 5.1, 3.8, 4.9, " +
						"2.7, 4.5, 4.2 cm. The sample mean is x̄ = 4.2cm and the sample standard " +
						"deviation is s = 1.5cm — but s is itself just an estimate, computed from " +
						"the same 6 noisy numbers, and with only 6 plants it could easily be off. " +
						"If you plug s in for σ and still use z* = 1.96 like before, is the " +
						"resulting interval still honest, or does it understate how uncertain you " +
						"really are?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"In 1908 William Gosset (publishing as 'Student') showed that once σ is " +
						"replaced by the sample estimate s, the quantity (x̄-μ)/(s/√n) no longer " +
						"follows a standard normal distribution — it follows a related but " +
						"heavier-tailed curve, Student's t-distribution, shaped by one parameter: " +
						"degrees of freedom, df = n-1 (one degree of freedom is 'used up' " +
						"estimating the mean before s can be computed).",
					"Walk the 6-plant example: n=6, so df = 6-1 = 5. Reading the t-distribution's " +
						"95% two-tailed critical value at df=5 gives t* = 2.571 — noticeably bigger " +
						"than z* = 1.96, about 31% wider. Compare the two margins of error: " +
						"z-margin = 1.96 × 1.5/√6 = 1.96 × 0.612 = 1.200cm; t-margin = 2.571 × " +
						"1.5/√6 = 2.571 × 0.612 = 1.575cm. The z-based interval, 4.2 ± 1.200 = " +
						"(3.0, 5.4), looks tighter and more impressive than the honest t-based " +
						"interval, 4.2 ± 1.575 = (2.625, 5.775) — but it's tighter because it's " +
						"quietly ignoring the extra uncertainty of estimating σ from just 6 points.",
					"That gap between t* and z* shrinks as df grows, because more data makes s a " +
						"more trustworthy stand-in for σ:",
					"• df=5 (n=6): t* = 2.571, about 31% wider than z*=1.96",
					"• df=10 (n=11): t* = 2.228, about 14% wider",
					"• df=30 (n=31): t* = 2.042, about 4% wider",
					"• df→∞: t* → 1.960 exactly, matching z* once the sample is large enough " +
						"that s's own uncertainty stops mattering",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The df slider reshapes the solid t-curve: at low df it's visibly shorter and " +
						"fatter-tailed than the dashed standard-normal curve behind it; drag df up " +
						"and the two curves converge until they're nearly indistinguishable by " +
						"df≈30 — the same shrinking gap from the 5/10/30 numbers above, now drawn. " +
						"The shaded band marks ±t* (the 95% two-tailed critical value at the " +
						"current df) and always covers exactly 95% of the area under the t-curve, " +
						"the same invariant `standard-deviation` showed for ±1σ under the normal " +
						"curve — but here the band's width visibly changes with df even though the " +
						"area inside it doesn't.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Build a small-sample confidence interval that's actually honest about its " +
						"own uncertainty: read off t* for the sample's df instead of reaching for " +
						"z*=1.96 out of habit, and get the correctly wider interval — 4.2 ± 1.575 " +
						"instead of the overconfident 4.2 ± 1.200 for the fertilizer example. The " +
						"rule of thumb this gives you: below about df=30, the choice of t vs z " +
						"visibly changes the answer; above it, z is a fine approximation.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Early-stage clinical trials with a handful of patients, quality checks on a " +
						"short production run, and A/B tests still in their first few days of " +
						"traffic all share the fertilizer example's problem: a small n and an " +
						"unknown true σ. Statistical software defaults to the t-distribution for " +
						"exactly this reason — R's `t.test()` and Python's `scipy.stats.ttest_1samp` " +
						"both compute a t-statistic and compare it against a t-distribution, not a " +
						"normal one, unless you explicitly tell them σ is known.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'n=6, so df=5 — use t*=2.571 for a 95% interval, not " +
						"z*=1.96,' naming both the sample size's effect on df and the resulting " +
						"critical value explicitly.",
					"Not like this: treating 1.96 as a fixed universal multiplier for '95% " +
						"confidence' no matter the sample size. It's only exactly right in the " +
						"df→∞ limit; using it for a small sample produces an interval that's too " +
						"narrow — a false sense of precision, not a wrong-but-safe shortcut. A " +
						"second common slip: writing df=n instead of df=n-1.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "df", Label: "Degrees of freedom (df = n-1)", Min: 1, Max: 50, Step: 1, Def: 5},
		},
		Render: render,
	})
}

// StudentTPDF is the probability density of Student's t-distribution with df
// degrees of freedom:
//
//	f(t) = Γ((df+1)/2) / (√(df·π)·Γ(df/2)) · (1 + t²/df)^(-(df+1)/2)
//
// df must be > 0; StudentTPDF returns 0 otherwise. Computed in log space
// (via math.Lgamma) rather than evaluating Γ directly and dividing — Γ(x)
// overflows float64 once x exceeds about 171, which a large df (a big
// sample) reaches immediately, and Inf/Inf would silently produce NaN.
func StudentTPDF(t, df float64) float64 {
	if df <= 0 {
		return 0
	}
	logNum, _ := math.Lgamma((df + 1) / 2)
	logDen, _ := math.Lgamma(df / 2)
	logCoef := logNum - logDen - 0.5*math.Log(df*math.Pi)
	logDensity := logCoef - (df+1)/2*math.Log(1+t*t/df)
	return math.Exp(logDensity)
}

// StdNormalPDF is the probability density of the standard normal
// distribution — the curve StudentTPDF converges to as df grows.
func StdNormalPDF(x float64) float64 {
	return math.Exp(-x*x/2) / math.Sqrt(2*math.Pi)
}

// tIntegrand is StudentTPDF pulled back through the substitution t = tan(θ),
// θ ∈ (-π/2, π/2): tIntegrand(θ,df) = StudentTPDF(tan θ, df)·sec²θ. This
// substitution maps the t-distribution's infinite domain onto the finite
// interval (-π/2, π/2), so StudentTCDF below can integrate it with a plain
// composite Simpson's rule instead of picking an arbitrary truncation bound
// for the tails (which, for small df, decay too slowly to truncate safely).
// The integrand stays bounded even at the endpoints: as θ→±π/2, t→±∞ and
// StudentTPDF(t,df)·(1+t²) tends to a finite limit (0 for df>1, 1/π for
// df=1, the Cauchy case) rather than blowing up.
func tIntegrand(theta, df float64) float64 {
	t := math.Tan(theta)
	return StudentTPDF(t, df) * (1 + t*t)
}

// StudentTCDF returns P(T ≤ t) for a t-distributed T with df degrees of
// freedom, via composite Simpson's rule on tIntegrand over the substituted,
// finite domain described above.
func StudentTCDF(t, df float64) float64 {
	if df <= 0 {
		return 0
	}
	const n = 2000 // even, for Simpson's rule
	lo, hi := -math.Pi/2+1e-9, math.Atan(t)
	h := (hi - lo) / n
	sum := tIntegrand(lo, df) + tIntegrand(hi, df)
	for i := 1; i < n; i++ {
		theta := lo + float64(i)*h
		weight := 4.0
		if i%2 == 0 {
			weight = 2.0
		}
		sum += weight * tIntegrand(theta, df)
	}
	area := sum * h / 3
	if area < 0 {
		return 0
	}
	if area > 1 {
		return 1
	}
	return area
}

// CriticalValue returns t* such that a two-tailed interval [-t*, t*] holds
// exactly 1-alpha of the t-distribution's probability at the given df —
// e.g. CriticalValue(5, 0.05) ≈ 2.571, the 95% critical value at df=5. Found
// by bisecting StudentTCDF, which is continuous and strictly increasing in t.
func CriticalValue(df, alpha float64) float64 {
	target := 1 - alpha/2
	lo, hi := 0.0, 1000.0
	for i := 0; i < 60; i++ {
		mid := (lo + hi) / 2
		if StudentTCDF(mid, df) < target {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

func render(p map[string]float64) string {
	df := p["df"]
	if df < 1 {
		df = 1
	}

	const xmin, xmax = -8.0, 8.0
	peak := StudentTPDF(0, df) // the mode of a symmetric, unimodal density
	c := viz.New(680, 340, xmin, xmax, 0, peak*1.15)
	c.PadT = 56
	c.Axes()
	for x := -8.0; x <= 8.0; x += 2 {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	// Dashed standard normal behind the solid t-curve, so the two are easy
	// to compare directly -- the whole lesson is how close (or not) they are.
	normalCurve := viz.Sample(xmin, xmax, 240, StdNormalPDF)
	c.Path(normalCurve, viz.Muted, 1.5)

	tCurve := viz.Sample(xmin, xmax, 240, func(x float64) float64 {
		return StudentTPDF(x, df)
	})

	tStar := CriticalValue(df, 0.05)
	c.Area(tCurve, -tStar, tStar, viz.Accent, 0.18)
	c.Path(tCurve, viz.Accent, 2.5)
	c.VLine(-tStar, viz.Warm, true)
	c.VLine(tStar, viz.Warm, true)

	const zStar = 1.959964 // 95% two-tailed critical value of the standard normal
	widerPct := (tStar/zStar - 1) * 100

	c.Text(20, 20, fmt.Sprintf("df = %.0f    t* (95%%, two-tailed) = %.3f    z* = %.3f",
		df, tStar, zStar), 14, viz.Ink, "start")
	c.Text(20, 40, fmt.Sprintf("shaded band ±t* holds 95%% of this curve -- %.1f%% wider than ±z* would be",
		widerPct), 12, viz.Muted, "start")

	return c.String()
}
