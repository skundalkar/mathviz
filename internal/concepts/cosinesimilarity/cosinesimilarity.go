// Package cosinesimilarity visualizes cosine similarity: comparing two
// vectors by the angle between them instead of by how far apart their tips
// are, so two vectors pointing the same way score as identical no matter
// how different their lengths are — the comparison embedding search runs
// millions of times a second.
package cosinesimilarity

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "cosine-similarity",
		Seq:   33,
		Title: "Cosine similarity (comparing direction, not distance)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "ux", Label: "u_x", Min: -6, Max: 6, Step: 1, Def: 4},
			{Key: "uy", Label: "u_y", Min: -6, Max: 6, Step: 1, Def: 2},
			{Key: "vx", Label: "v_x", Min: -6, Max: 6, Step: 1, Def: 2},
			{Key: "vy", Label: "v_y", Min: -6, Max: 6, Step: 1, Def: 1},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(560, 480, -7, 7, -7, 7)
	return c.String()
}
