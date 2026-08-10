// Package entropy visualizes Shannon entropy for a two-outcome (coin-flip)
// distribution: how "surprising" an outcome is depends on how likely it was,
// and entropy is just the average surprise across both outcomes. A fair coin
// is the most surprising, on average, of any biased coin.
package entropy

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "entropy",
		Seq:   16,
		Title: "Entropy",
		Blurb: "Forget the math for a second: 'entropy' just means how mixed-up or hard-to-call " +
			"something is. A messy room has high entropy; a tidy one has low entropy. 'This " +
			"password has high entropy' means it's got no guessable pattern. 'The model's " +
			"predictions have high entropy' means it's genuinely unsure, spreading its confidence " +
			"across several answers instead of committing to one. A coin is the simplest version of " +
			"that same idea: a fair coin (p=0.5) is the ultimate toss-up — nobody has a good guess in " +
			"advance, so it's high-entropy. A trick coin that lands heads 90% of the time (p=0.9) is " +
			"low-entropy — you can already guess 'probably heads' and mostly be right, so watching it " +
			"land barely tells you anything new. Drag p and watch the curve trace exactly that: " +
			"entropy peaks at p=0.5 (true toss-up) and falls toward 0 as p nears 0 or 1 (the outcome " +
			"is basically already decided). The two bars below show why: the rare outcome would " +
			"genuinely surprise you if it happened, the common one wouldn't — entropy blends those " +
			"two by how often each actually occurs.",
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

// surpriseCap bounds how tall a surprise bar can grow. The p slider stops at
// 0.01, whose surprise is -log2(0.01) ≈ 6.64 bits, so 7 comfortably covers
// the slider's full range without a bar ever clipping.
const surpriseCap = 7.0

func render(params map[string]float64) string {
	pHeads := params["p"]
	pTails := 1 - pHeads

	c := viz.New(680, 420, 0, 1, 0, 1.1)
	// Push the curve down from the top (room for the header text) and
	// compress it away from the bottom (room for the surprise bars below,
	// drawn in raw pixel coordinates so they're unaffected by this).
	c.PadT = 70
	c.PadB = 230
	c.Axes()
	for x := 0.0; x <= 1.0; x += 0.25 {
		c.Tick(x, fmt.Sprintf("%.2g", x))
	}

	curve := viz.Sample(0, 1, 200, BinaryEntropy)
	c.Path(curve, viz.Accent, 2.5)
	c.VLine(pHeads, viz.Warm, true)

	hpx, hpy := c.X(pHeads), c.Y(BinaryEntropy(pHeads))
	c.Rect(hpx-4, hpy-4, 8, 8, viz.Warm, 0.9)

	c.Text(20, 24, fmt.Sprintf("P(heads) = %.2f    entropy H(p) ≈ %.3f bits", pHeads, BinaryEntropy(pHeads)),
		14, viz.Ink, "start")
	c.Text(20, 44, "entropy peaks at 1 bit when p = 0.5 — that's when the outcome is least predictable",
		12, viz.Muted, "start")

	// Below the curve: the surprise (-log2 p) of each outcome, if it happens.
	const barBase, barMaxH, barW = 372.0, 120.0, 70.0
	c.Text(20, 234, "surprise if that outcome actually happens: -log2(p)", 12, viz.Muted, "start")

	drawBar := func(x float64, prob float64, label string, color string) {
		s := Surprise(prob)
		h := s / surpriseCap * barMaxH
		if h > barMaxH {
			h = barMaxH
		}
		c.Rect(x, barBase-h, barW, h, color, 0.75)
		c.Text(x+barW/2, barBase-h-8, fmt.Sprintf("%.2f bits", s), 12, viz.Ink, "middle")
		c.Text(x+barW/2, barBase+18, fmt.Sprintf("%s (p=%.2f)", label, prob), 12, viz.Muted, "middle")
	}
	drawBar(180, pHeads, "heads", viz.Accent)
	drawBar(420, pTails, "tails", viz.Warm)

	return c.String()
}
