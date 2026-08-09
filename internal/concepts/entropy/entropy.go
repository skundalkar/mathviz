// Package entropy visualizes Shannon entropy for a two-outcome (coin-flip)
// distribution: how "surprising" an outcome is depends on how likely it was,
// and entropy is just the average surprise across both outcomes. A fair coin
// is the most surprising, on average, of any biased coin.
package entropy

import (
	"math"

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

// Surprise returns the information content of an outcome with probability p,
// in bits: -log2(p). Rarer outcomes (small p) carry more surprise; a certain
// outcome (p=1) carries none. p=0 is defined as +Inf (an outcome that never
// happens would be infinitely surprising if it somehow did).
func Surprise(p float64) float64 {
	if p <= 0 {
		return math.Inf(1)
	}
	return -math.Log2(p)
}

// BinaryEntropy returns the Shannon entropy, in bits, of a two-outcome
// distribution with P(heads) = p and P(tails) = 1-p:
//
//	H(p) = -p*log2(p) - (1-p)*log2(1-p)
//
// the probability-weighted average surprise across both outcomes. By
// convention 0*log2(0) = 0, so H is well-defined (and 0) at p=0 and p=1.
// H peaks at exactly 1 bit when p=0.5.
func BinaryEntropy(p float64) float64 {
	term := func(x float64) float64 {
		if x <= 0 {
			return 0
		}
		return -x * math.Log2(x)
	}
	return term(p) + term(1-p)
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 340, -1, 1, -1, 1).String()
}
