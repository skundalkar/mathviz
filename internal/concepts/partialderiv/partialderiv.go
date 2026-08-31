// Package partialderiv visualizes partial derivatives and the gradient:
// freezing all but one variable of f(x,y) reduces it back to the familiar
// single-variable slope from `derivative`, and packaging the two resulting
// slopes into one vector, the gradient, points in the direction f increases
// fastest — with its length telling you exactly how steep that direction is.
package partialderiv

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "partial-derivatives-gradient",
		Seq:   82,
		Title: "Partial derivatives & the gradient",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "x0", Label: "Point x0", Min: -3, Max: 3, Step: 0.25, Def: 1},
			{Key: "y0", Label: "Point y0", Min: -3, Max: 3, Step: 0.25, Def: 2},
			{Key: "theta", Label: "Direction theta", Min: 0, Max: 360, Step: 1, Def: 0, Unit: "°"},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(560, 560, -4.5, 4.5, -4.5, 4.5)
	c.Axes()
	return c.String()
}
