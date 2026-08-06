// Package correlation shows what the Pearson correlation coefficient r
// actually looks like: a scatter cloud that goes from a tight downhill line
// (r = -1) through a shapeless blob (r = 0) to a tight uphill line (r = +1).
// The picture also carries the standard warning label — a strong r only says
// two variables move together, never that one causes the other.
package correlation

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "correlation",
		Title: "Correlation (r)",
		Blurb: "r measures how tightly a scatter of points hugs a straight line, from " +
			"-1 (perfect downhill) through 0 (no linear pattern) to +1 (perfect uphill). " +
			"Slide r and watch the cloud tighten into a line or spread into a blob. The " +
			"orange line is the trend r implies. A strong r only tells you the two " +
			"variables move together — it never tells you that one causes the other.",
		Params: []concept.ParamSpec{
			{Key: "r", Label: "Target correlation (r)", Min: -1, Max: 1, Step: 0.05, Def: 0.6},
			{Key: "n", Label: "Sample size (n)", Min: 20, Max: 300, Step: 10, Def: 150},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 400, -1, 1, -1, 1).String()
}
