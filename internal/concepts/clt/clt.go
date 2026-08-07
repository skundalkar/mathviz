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

func render(p map[string]float64) string {
	return viz.New(680, 340, 0, 5, 0, 1).String()
}
