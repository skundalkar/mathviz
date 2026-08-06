// Package meanmedianmode visualizes how a skewed distribution pulls the mean,
// median, and mode apart. All three coincide for a symmetric distribution, but
// as skew increases they peel off in a fixed order: mode stays at the peak,
// the median sits at the halfway point, and the mean gets dragged toward the
// long tail. A log-normal curve is used because all three statistics have
// simple closed forms, so the picture is exact rather than sampled.
package meanmedianmode

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "mean-median-mode",
		Title: "Mean, median & mode",
		Blurb: "For a symmetric distribution, the mean, median, and mode all land in the " +
			"same place. Skew the distribution and they peel apart in a fixed order: the mode " +
			"stays at the peak, the median sits at the 50/50 point, and the mean gets dragged " +
			"toward the long tail — because the mean, unlike the median, is sensitive to extreme " +
			"values. Drag the skew slider and watch mean > median > mode open up as the right " +
			"tail lengthens.",
		Params: []concept.ParamSpec{
			{Key: "sigma", Label: "Skew (σ)", Min: 0.05, Max: 1.2, Step: 0.05, Def: 0.6},
			{Key: "mu", Label: "Log-location (μ)", Min: -0.5, Max: 0.5, Step: 0.05, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(680, 340, 0, 1, 0, 1).String()
}
