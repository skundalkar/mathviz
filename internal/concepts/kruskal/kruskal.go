// Package kruskal visualizes Kruskal's algorithm: building the cheapest
// possible tree that still connects every node in a weighted graph, by
// sorting every road cheapest-first and greedily taking each one unless it
// would reconnect two cities that a cheaper road already connected --
// which would only create a wasteful cycle, never a shorter tree. It works
// on the exact same 5-city network dijkstras-algorithm walks through, but
// answers a different question: not "what's the cheapest route between two
// cities" but "what's the cheapest set of roads that keeps every city
// reachable from every other one."
package kruskal

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "kruskals-mst",
		Seq:   100,
		Title: "Kruskal's algorithm (cheapest network connecting every node)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder -- filled in once the math and picture exist.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "step", Label: "Step (consider one edge at a time)", Min: 0, Max: 7, Step: 1, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(700, 460, 0, 1, 0, 1).String()
}
