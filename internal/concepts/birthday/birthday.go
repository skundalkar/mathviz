// Package birthday visualizes the birthday paradox: how the probability
// that at least two people in a room share a birthday climbs to "more
// likely than not" far sooner than most people's gut instinct expects.
package birthday

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "birthday-paradox",
		Seq:   36,
		Title: "Birthday paradox (collisions sooner than you'd guess)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "People in the room", Min: 1, Max: 100, Step: 1, Def: 23},
			{Key: "days", Label: "Possible birthdays", Min: 2, Max: 365, Step: 1, Def: 365},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(680, 420, 0, 100, 0, 1)
	return c.String()
}
