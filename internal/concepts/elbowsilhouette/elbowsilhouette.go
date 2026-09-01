// Package elbowsilhouette visualizes two ways to pick k for k-means when
// nobody hands you the right number of clusters: watch where adding another
// cluster stops meaningfully shrinking the within-cluster distance (the
// elbow method), and separately score how cleanly each point sits in its own
// cluster versus the next-closest one (the silhouette score).
package elbowsilhouette

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "elbow-method-silhouette-score",
		Seq:   86,
		Title: "Elbow method & silhouette score (choosing k for k-means)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "k", Label: "Number of clusters (k)", Min: 1, Max: 6, Step: 1, Def: 3},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(680, 460, 0, 1, 0, 1)
	return c.String()
}
