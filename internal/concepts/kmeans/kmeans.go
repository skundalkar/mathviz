// Package kmeans visualizes k-means clustering: repeatedly assigning each
// point to its nearest centroid, then moving each centroid to the average
// of the points now assigned to it, until the groups stop changing.
package kmeans

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "k-means-clustering",
		Seq:   32,
		Title: "K-means clustering (grouping points by nearest centroid)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "iteration", Label: "Iteration", Min: 0, Max: 3, Step: 1, Def: 0},
			{Key: "seed", Label: "Starting centroids", Min: 0, Max: 1, Step: 1, Def: 1},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(600, 420, 0, 1, 0, 1)
	return c.String()
}
