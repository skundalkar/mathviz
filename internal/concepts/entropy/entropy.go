// Package entropy visualizes Shannon entropy for a two-outcome (coin-flip)
// distribution: how "surprising" an outcome is depends on how likely it was,
// and entropy is just the average surprise across both outcomes. A fair coin
// is the most surprising, on average, of any biased coin.
package entropy

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "entropy",
		Title: "Entropy",
		Blurb: "A weather app that says '50% chance of rain' is telling you something genuinely " +
			"uncertain — either outcome would be news. One that says '99% chance of rain' has " +
			"basically already told you the answer; if it then actually rains, that's barely " +
			"surprising at all. Entropy is that intuition made precise: rarer outcomes carry " +
			"more 'surprise' (information), quantified as -log2(p), and entropy is the average " +
			"surprise you should expect across all outcomes, weighted by how often each happens. " +
			"For a coin with P(heads) = p, entropy peaks at exactly 1 bit when p = 0.5 — maximum " +
			"uncertainty — and falls toward 0 as p approaches 0 or 1, where the outcome is all " +
			"but certain and learning it tells you almost nothing.",
		Params: []concept.ParamSpec{
			{Key: "p", Label: "P(heads)", Min: 0.01, Max: 0.99, Step: 0.01, Def: 0.5},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 340, -1, 1, -1, 1).String()
}
