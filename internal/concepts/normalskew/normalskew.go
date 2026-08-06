// Package normalskew visualizes skewness and (excess) kurtosis as knobs that
// reshape a standard normal curve. It uses the Gram-Charlier A series, a
// closed-form correction to the normal density built from Hermite
// polynomials: a small skew or kurtosis term tilts or reshapes the bell curve
// while leaving it centered and (to first order) still normalized.
package normalskew

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "normal-vs-skew",
		Title: "Normal vs. skewed & fat-tailed",
		Blurb: "A standard normal curve is symmetric with zero excess kurtosis: its skewness " +
			"and kurtosis knobs are both at zero. Turn the skew knob and the curve leans — one " +
			"tail lengthens while the other shortens, exactly what a real-world skewed dataset " +
			"looks like. Turn the kurtosis knob and the curve stays symmetric but its peak " +
			"sharpens and its tails fatten (positive excess kurtosis, 'leptokurtic') or the " +
			"opposite (negative, 'platykurtic'). The dashed curve is the plain normal for " +
			"comparison.",
		Params: []concept.ParamSpec{
			{Key: "skew", Label: "Skewness", Min: -0.9, Max: 0.9, Step: 0.1, Def: 0},
			{Key: "kurt", Label: "Excess kurtosis", Min: -1, Max: 2, Step: 0.1, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(680, 340, -5, 5, 0, 1).String()
}
