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
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "central-limit-theorem",
		Title: "Central limit theorem",
		Blurb: "The population here (n=1) is an exponential distribution: sharply peaked at " +
			"zero with a long right tail — about as far from a bell curve as a distribution " +
			"gets. Raise the sample size n and watch the distribution of the *sample mean* " +
			"tighten around the true mean and reshape itself into a normal curve (dashed), " +
			"regardless of how skewed the population you started from was. This is the central " +
			"limit theorem: averages of independent draws go normal as n grows.",
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
	return viz.New(680, 340, 0, 5, 0, 1).String()
}
