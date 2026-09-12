// Package astar visualizes A* search: Dijkstra's algorithm with one change
// -- instead of always expanding whichever not-yet-closed node has the
// smallest cost-so-far, it expands whichever has the smallest cost-so-far
// PLUS a heuristic estimate of the remaining distance to a fixed goal. Both
// are built from the exact same relaxation rule dijkstras-algorithm used;
// the only difference is what the search sorts by, and dijkstras-algorithm
// is exactly the special case where that heuristic is zero everywhere.
package astar

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "a-star-search",
		Seq:   99,
		Title: "A* search (Dijkstra with a distance-to-goal hint)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder -- filled in once the math and picture exist.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "heuristic", Label: "Use heuristic (0=plain Dijkstra, 1=A*)", Min: 0, Max: 1, Step: 1, Def: 1},
			{Key: "step", Label: "Step (expand one node at a time)", Min: 0, Max: 6, Step: 1, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(700, 460, 0, 1, 0, 1).String()
}
