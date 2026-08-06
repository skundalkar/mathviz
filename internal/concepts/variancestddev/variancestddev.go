// Package variancestddev visualizes why we square deviations before averaging
// them, then take a square root to get back to the original units. Raw
// deviations from the mean always sum to zero (positive and negative cancel),
// so squaring is what lets them accumulate into a real measure of spread —
// variance. But variance is in squared units, so stddev = √variance converts
// it back to something comparable to the data itself.
package variancestddev

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "variance-vs-stddev",
		Title: "Variance vs. standard deviation",
		Blurb: "Raw deviations from the mean always sum to zero — positives and negatives " +
			"cancel out, so they can't measure spread directly. Squaring each deviation makes " +
			"every term positive (and punishes outliers hardest, since they get squared too), " +
			"so the average squared deviation — the variance — actually grows with spread. The " +
			"catch: variance is in squared units. Taking its square root, the standard " +
			"deviation, brings the number back to the same units as the data. Drag spread to " +
			"see stddev scale linearly while variance scales with the square; drag the outlier " +
			"to see how much one extreme point can inflate both.",
		Params: []concept.ParamSpec{
			{Key: "spread", Label: "Spread", Min: 0.3, Max: 2.5, Step: 0.1, Def: 1},
			{Key: "outlier", Label: "Outlier push", Min: 0, Max: 10, Step: 0.5, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(680, 360, 0, 1, 0, 1).String()
}
