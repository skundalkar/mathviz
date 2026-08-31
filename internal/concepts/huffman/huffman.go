// Package huffman visualizes Huffman coding: greedily merging the two
// least-frequent symbols over and over builds a binary tree whose leaves
// are the symbols, with frequent symbols landing shallow (short codes) and
// rare symbols landing deep (long codes) -- and the resulting average code
// length always comes within 1 bit of `entropy`'s theoretical floor.
package huffman

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "huffman-coding",
		Seq:   84,
		Title: "Huffman coding (an optimal code built from symbol frequencies)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "skew", Label: "Frequency skew", Min: 0.15, Max: 0.95, Step: 0.05, Def: 0.5},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(720, 460, 0, 1, 0, 1)
	return c.String()
}
