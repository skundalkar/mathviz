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
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You've probably heard 'entropy' used loosely outside math — 'this password " +
						"has high entropy,' 'there's a lot of entropy in this org chart,' 'the " +
						"model's predictions have high entropy.' All of those point at the same " +
						"shape of thing: could this go a lot of different ways, with no clear " +
						"favorite, or is it pretty much locked in on one obvious answer?",
					"A coin is the simplest place to check that intuition. A fair coin (p=0.5) is " +
						"the ultimate toss-up — nobody has a good guess in advance. A trick coin " +
						"that lands heads 90% of the time (p=0.9) is nearly locked in — you can " +
						"already guess 'probably heads' and mostly be right. Is there a way to say " +
						"precisely how much more 'toss-up' the fair coin is than the trick coin, " +
						"instead of just eyeballing it?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"A fair coin (p=0.5) is high-entropy: nobody has a good guess in advance, so " +
						"watching it land genuinely tells you something. A trick coin at p=0.9 is " +
						"low-entropy: you can already guess 'probably heads' and mostly be right, " +
						"so watching it land barely tells you anything new.",
					"Entropy is the average of how surprising each outcome would be, weighted by " +
						"how often it actually happens: the rare outcome would genuinely surprise " +
						"you if it happened, the common one wouldn't, and entropy blends those two " +
						"together by how often each actually occurs.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Drag p and watch the curve trace exactly that pattern: entropy peaks at 1 " +
						"bit when p = 0.5 (true toss-up) and falls toward 0 as p nears 0 or 1 (the " +
						"outcome is basically already decided).",
					"The two bars below the curve show why: the rarer outcome's bar (how " +
						"surprising it would be if it happened) is taller, the common outcome's " +
						"bar is short since it barely tells you anything new, and entropy is the " +
						"probability-weighted blend of those two bars.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"'This password has high entropy' or 'the model's predictions have high " +
						"entropy' isn't just a vibe anymore — it's a specific, computable number " +
						"(in bits) you can use to actually compare how mixed-up two passwords, " +
						"models, or coins are, instead of eyeballing which one feels more locked " +
						"in.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"• 'This password has high entropy' — there's no discernible pattern to " +
						"exploit; a guesser has no shortcut, every character genuinely could've " +
						"been anything. A low-entropy password ('password123') is the opposite — " +
						"heavily concentrated on a small, guessable set of likely choices.",
					"• 'There's a lot of entropy in this org chart right now' — nobody's really " +
						"sure who owns what; responsibilities are scattered, not concentrated on " +
						"clear owners.",
					"• 'Left alone, codebases/garages/inboxes tend toward entropy' — the physics " +
						"metaphor directly: things drift toward disorganized unless someone spends " +
						"effort keeping them organized.",
					"• 'The model's predictions have high entropy here' — the model is genuinely " +
						"unsure on this input, spreading its confidence across several possible " +
						"answers instead of piling it onto one.",
					"If you've heard 'cross-entropy loss' in an AI context, it's the same core " +
						"idea extended one step further: how well a model's stated confidence " +
						"matches what actually turned out to be true.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'That rigged coin that always lands heads is extremely " +
						"low-entropy — utterly predictable, even though it's unfair.' Entropy is " +
						"about how spread-out the possibilities are, not about whether something " +
						"is right, fair, or good.",
					"Not like this: 'That's so biased, it must have high entropy' — if something " +
						"is just wrong, or leans hard in one predictable direction, that's " +
						"actually the opposite of entropy: it's low entropy, just low entropy " +
						"pointed at the wrong answer. Reach for 'entropy' specifically when the " +
						"honest description is 'spread out, could go several ways' — not just " +
						"'off' or 'biased.'",
				},
			},
		},
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
