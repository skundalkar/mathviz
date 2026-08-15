// Package eigen visualizes eigenvectors and eigenvalues of a 2x2 symmetric
// matrix: most directions get bent sideways by the transformation, but two
// special, perpendicular directions only get stretched (or flipped) — never
// rotated off their own line.
package eigen

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "eigenvectors-eigenvalues",
		Seq:   35,
		Title: "Eigenvectors & eigenvalues (directions that only stretch)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "Matrix a (x-axis stretch)", Min: 0.5, Max: 2.5, Step: 0.1, Def: 2},
			{Key: "d", Label: "Matrix d (y-axis stretch)", Min: 0.5, Max: 2.5, Step: 0.1, Def: 1},
			{Key: "b", Label: "Matrix b (shear / coupling)", Min: -1.2, Max: 1.2, Step: 0.1, Def: 0.8},
			{Key: "angle", Label: "Test vector angle θ", Min: 0, Max: 360, Step: 5, Def: 30, Unit: "°"},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(560, 560, -7, 7, -7, 7)
	return c.String()
}
