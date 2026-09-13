// Package autocorr visualizes autocorrelation: how well a time series
// matches up with a lagged copy of itself, computed the same way
// correlation measures how two variables move together -- except here the
// "second variable" is the same series, shifted in time. Sweeping the lag
// traces a correlogram, and where that correlogram peaks reveals a series's
// hidden period even when the naked eye can't quite pin it down through the
// noise.
package autocorr

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "autocorrelation",
		Seq:   103,
		Title: "Autocorrelation (correlating a series with its lagged self)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder -- filled in once the math and picture exist.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "lag", Label: "Lag (k)", Min: 0, Max: 20, Step: 1, Def: 3},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(700, 460, 0, 1, 0, 1).String()
}
