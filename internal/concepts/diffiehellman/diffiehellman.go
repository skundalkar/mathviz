// Package diffiehellman visualizes Diffie-Hellman key exchange: two parties
// who agree on nothing secret in advance -- only a public prime p and a
// public base g -- each combine their own private number with the public
// values to land on the exact same shared secret, without either private
// number ever crossing the public channel. The running example is the
// classic small-number walkthrough (p=23, g=5), so every value in the
// pipeline can be checked by hand.
package diffiehellman

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "diffie-hellman-key-exchange",
		Seq:   83,
		Title: "Diffie-Hellman key exchange (a shared secret over a public channel)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "Alice's private number a", Min: 1, Max: 21, Step: 1, Def: 6},
			{Key: "b", Label: "Bob's private number b", Min: 1, Max: 21, Step: 1, Def: 15},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(760, 480, 0, 1, 0, 1)
	return c.String()
}
