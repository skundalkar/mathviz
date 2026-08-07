// Package confint visualizes what "95% confidence" actually means: not that
// there's a 95% chance the true mean lies in one particular interval, but
// that if you repeated the experiment many times, about 95% of the intervals
// you'd build would contain the true mean. To make that concrete without any
// randomness, the picture draws a fixed set of 20 hypothetical sample means
// placed at evenly spaced quantiles of the sampling distribution — a
// deterministic stand-in for "20 repeated experiments" — and shows which of
// their confidence intervals happen to capture the true mean.
package confint

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "confidence-interval",
		Title: "Confidence interval",
		Blurb: "The dashed line is the true population mean, which in real life you never " +
			"get to see. Each row is a hypothetical experiment: draw a sample, compute its mean, " +
			"and build an interval around it. A green interval happened to capture the true " +
			"mean; a red one missed. \"95% confidence\" doesn't describe any single interval — " +
			"it describes this whole process: about 95% of intervals built this way capture the " +
			"true mean. Raise the confidence level and every interval widens (fewer misses); " +
			"raise the sample size and every interval narrows (more precision).",
		Params: []concept.ParamSpec{
			{Key: "confidence", Label: "Confidence level", Min: 50, Max: 99, Step: 1, Def: 95, Unit: "%"},
			{Key: "n", Label: "Sample size (n)", Min: 2, Max: 100, Step: 1, Def: 20},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(680, 460, -1, 1, 0, 1).String()
}
