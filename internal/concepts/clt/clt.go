// Package clt visualizes the central limit theorem: draw samples of size n
// from a distinctly non-normal population (an exponential — sharply peaked at
// zero, long right tail) and look at the distribution of the sample mean. For
// n=1 that's just the population itself, but as n grows the sample-mean
// distribution tightens around the true mean and its shape converges to a
// normal curve — no matter how skewed the population you started from. The
// sum of n iid Exponential(λ) draws is exactly Gamma(n, λ), so the sample
// mean's distribution has a closed form and needs no random sampling.
package clt

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "central-limit-theorem",
		Seq:   8,
		Title: "Central limit theorem",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"At a school fair, a jar of jellybeans sits on a table and everyone guesses how " +
						"many are inside. Look at any one guess and it's basically useless — kids " +
						"guess 50, guess 900, guess a suspiciously round 1,000; individually, wildly " +
						"scattered. Is there any way to turn a pile of unreliable individual guesses " +
						"into something you can actually trust?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"There is: average them. Say the jar actually holds 620 jellybeans:",
					"• Five kids guess 50, 900, 1000, 300, and 750 — wildly scattered, off by as much " +
						"as 570 in either direction.",
					"• Average those five guesses: (50+900+1000+300+750)/5 = 3000/5 = 600, only 20 " +
						"off from the true 620 — closer than four of the five individual guesses, even " +
						"though every guess going into it looked like noise.",
					"That's the 'wisdom of crowds,' and it isn't magic — it's the central limit " +
						"theorem: no matter how skewed the population of individual values is, the " +
						"distribution of the AVERAGE of many draws from it tightens around the true " +
						"mean and turns bell-shaped as you average more of them. The population here " +
						"(n=1) is deliberately about as far from normal as it gets: an exponential " +
						"distribution, sharply peaked at zero with a long right tail — think 'the " +
						"population of guesses, with a few extreme overestimates.'",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Raise the sample size slider n and watch the sample-mean distribution (thick " +
						"curve) visibly tighten and reshape into a normal curve (the thin reference " +
						"line), regardless of the population rate λ. At n=1 the picture shows the " +
						"skewed exponential population itself — as far from normal as the jellybean " +
						"guesses were scattered; raise n and it converges toward normal exactly the " +
						"way averaging the five guesses landed closer to 620 than almost any one guess " +
						"did. Because the sum of n exponential draws has an exact closed form (a Gamma " +
						"distribution), this is the true shape at every n, not a simulation with " +
						"sampling noise of its own.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can now trust an average built from individually unreliable, skewed inputs " +
						"— treating it as approximately normal and applying standard normal-based " +
						"statistical tools to it — as long as you're averaging enough independent " +
						"draws, no matter how weird any one input looked on its own.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Polling averages, A/B test metrics, and quality-control sample means can be " +
						"treated as roughly normal — and standard statistical tests applied — even " +
						"when individual data points are skewed or weird, as long as the sample is big " +
						"enough. It's also why 'the sample size is too small' is a real objection: " +
						"with just a handful of draws from a skewed population, the average is still " +
						"skewed too, and normal-based confidence intervals can quietly mislead.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'Let's average a bunch of independent estimates, the noise " +
						"should wash out' — individual guesses can be wildly off, but averaging enough " +
						"of them (the 'wisdom of crowds') reliably lands close to the truth, " +
						"regardless of how weird any one guess looked.",
					"Not like this: 'One really confident estimate beats an average of many rough " +
						"ones' — often false when the rough estimates are independent; the average's " +
						"error shrinks as you add more of them, while one confident-but-biased guess " +
						"never self-corrects at all.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Sample size (n)", Min: 1, Max: 30, Step: 1, Def: 1},
			{Key: "lambda", Label: "Population rate (λ)", Min: 0.3, Max: 3, Step: 0.1, Def: 1},
		},
		Render: render,
	})
}

// ExponentialPDF is the density of an Exponential(λ) population: sharply
// peaked at zero with a long right tail, mean 1/λ, variance 1/λ².
func ExponentialPDF(x, lambda float64) float64 {
	if x < 0 || lambda <= 0 {
		return 0
	}
	return lambda * math.Exp(-lambda*x)
}

// SampleMeanPDF is the exact density of the mean of n iid Exponential(λ)
// draws. The sum of n such draws is Gamma(shape=n, rate=λ); dividing by n
// rescales that to Gamma(shape=n, rate=n·λ), a closed form with no sampling
// required:
//
//	f(m) = (nλ)^n · m^(n-1) · e^(-nλm) / (n-1)!
//
// Computed in log-space (via math.Lgamma) so it stays accurate for the large
// n this concept's slider allows.
func SampleMeanPDF(m float64, n int, lambda float64) float64 {
	if n < 1 || lambda <= 0 || m <= 0 {
		return 0
	}
	rate := float64(n) * lambda
	lgammaN, _ := math.Lgamma(float64(n))
	logf := float64(n)*math.Log(rate) - rate*m - lgammaN
	if n > 1 {
		logf += float64(n-1) * math.Log(m)
	}
	return math.Exp(logf)
}

// SampleMeanVariance is Var(mean of n iid draws) = Var(population)/n, with
// Var(Exponential(λ)) = 1/λ².
func SampleMeanVariance(n int, lambda float64) float64 {
	if n < 1 || lambda <= 0 {
		return 0
	}
	return 1 / (float64(n) * lambda * lambda)
}

// NormalPDF is the density of a normal distribution — used to draw the
// reference curve the sample-mean distribution converges to.
func NormalPDF(x, mu, sigma float64) float64 {
	if sigma <= 0 {
		return 0
	}
	z := (x - mu) / sigma
	return math.Exp(-0.5*z*z) / (sigma * math.Sqrt(2*math.Pi))
}

func render(p map[string]float64) string {
	n := int(p["n"] + 0.5)
	if n < 1 {
		n = 1
	}
	lambda := p["lambda"]
	if lambda <= 0 {
		lambda = 0.3
	}

	mean := 1 / lambda
	sigma := math.Sqrt(SampleMeanVariance(n, lambda))

	// Fixed x window sized off the population (n=1) so the curve visibly
	// narrows within the same frame as n grows, instead of the frame
	// rescaling with it.
	const xmax = 5.0
	xhi := xmax / lambda

	curve := viz.Sample(1e-4, xhi, 300, func(x float64) float64 {
		return SampleMeanPDF(x, n, lambda)
	})
	refCurve := viz.Sample(0, xhi, 300, func(x float64) float64 {
		return NormalPDF(x, mean, sigma)
	})

	peak := 0.0
	for _, pt := range curve {
		if pt[1] > peak {
			peak = pt[1]
		}
	}
	if peak <= 0 {
		peak = lambda
	}

	c := viz.New(680, 340, 0, xhi, 0, peak*1.2)
	c.Axes()
	step := xhi / 5
	for x := step; x <= xhi; x += step {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	// Reference normal with the same mean & variance, thin and muted — the
	// shape the sample-mean curve converges to.
	c.Path(refCurve, viz.Muted, 1.5)
	c.Path(curve, viz.Accent, 2.5)
	c.VLine(mean, viz.Ink, false)

	c.Text(20, 24, fmt.Sprintf("n = %d    λ = %.1f", n, lambda), 14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("sample mean: μ = %.2f (fixed)    σ = %.3f (shrinks as n grows)", mean, sigma),
		12, viz.Muted, "start")
	c.Text(20, 62, "thick: exact distribution of the sample mean   thin: matching normal curve",
		12, viz.Muted, "start")

	return c.String()
}
