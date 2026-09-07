// Package hmm visualizes a Hidden Markov Model: the classic "guess the
// weather from your neighbor's chores" example, where the true state
// (Sunny or Rainy) is never observed directly -- only a noisy side-effect
// of it (whether your neighbor went for a Walk, went Shopping, or did
// Cleaning) is. markov-chains modeled directly-observable weather evolving
// day to day; this concept adds a layer of indirection on top of that same
// two-state chain and asks: given only the observed activities, what's the
// single most likely sequence of hidden weather that produced them? The
// Viterbi algorithm answers that exactly, without ever seeing the weather.
package hmm

import (
	"fmt"
	"strings"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "hidden-markov-models",
		Seq:   92,
		Title: "Hidden Markov models (guessing the weather you can't see)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"markov-chains modeled weather you could always see directly — Sunny or " +
						"Rainy — and asked how the *known* state evolves day to day. But suppose " +
						"you're stuck inside all week with no window, and the only information " +
						"reaching you is a daily text from your neighbor across the street saying " +
						"whether they went for a Walk, went Shopping, or did some Cleaning. You " +
						"know, from watching them for months, that they tend to walk on sunny days " +
						"and clean on rainy ones — but it's only a tendency: they might still walk " +
						"on a rainy day, or clean on a sunny one, just less often. A single day's " +
						"activity doesn't prove the weather. If the state you actually care about " +
						"(Sunny or Rainy) is hidden, and all you ever observe is a noisy side-effect " +
						"of it, can you still say anything precise about which days were actually " +
						"sunny — or are you stuck guessing?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Reuse markov-chains's two states, Sunny and Rainy, with the exact same " +
						"stickiness convention: P(stay Sunny)=a and P(stay Rainy)=b. Add one new " +
						"ingredient — an emission probability for each state, how likely it is to " +
						"produce each observed activity:",
					"|State|Walk|Shop|Clean|",
					"|Sunny|0.60|0.30|0.10|",
					"|Rainy|0.10|0.40|0.50|",
					"Start with a=0.6, b=0.7 (a rainy day is slightly stickier than a sunny one) " +
						"and a prior P(Sunny on day 0)=0.4, P(Rainy on day 0)=0.6. Your neighbor " +
						"reports three days in a row: Walk, Shop, Clean. The Viterbi algorithm " +
						"tracks, for every day and every hidden state, the probability of the " +
						"single best path that ends there — call it δ (delta) — updating it day by " +
						"day as δ_t(s) = [max over previous state s' of δ_{t-1}(s')·P(s'→s)] × " +
						"P(observation_t | s):",
					"• Day 0 (Walk): δ(Sunny) = 0.4×0.60 = 0.2400. δ(Rainy) = 0.6×0.10 = 0.0600. " +
						"Sunny is way ahead — a walk is 6x more likely on a sunny day.",
					"• Day 1 (Shop): the best path into Sunny comes from Sunny " +
						"(0.2400×0.6=0.1440, beating Rainy→Sunny's 0.0600×0.3=0.0180), then ×0.30 " +
						"for the Shop emission → δ(Sunny)=0.0432. The best path into Rainy also " +
						"comes from Sunny (0.2400×0.4=0.0960, beating Rainy→Rainy's " +
						"0.0600×0.7=0.0420), then ×0.40 → δ(Rainy)=0.0384.",
					"• Day 2 (Clean): the best path into Sunny again comes from Sunny " +
						"(0.0432×0.6=0.02592), ×0.10 → δ(Sunny)=0.002592. The best path into Rainy " +
						"comes from Rainy this time (0.0384×0.7=0.02688, beating Sunny→Rainy's " +
						"0.0432×0.4=0.01728), ×0.50 → δ(Rainy)=0.01344.",
					"δ(Rainy)=0.01344 beats δ(Sunny)=0.002592 on the final day, so the most " +
						"likely explanation ends on Rainy. Walking the backpointers back — day 2 " +
						"came from Rainy, day 1 (feeding into that Rainy) came from Sunny, day 0 is " +
						"Sunny by construction — gives the single most likely hidden sequence: " +
						"Sunny, Rainy, Rainy, with overall probability 0.01344. (These are the " +
						"textbook numbers used to introduce the Viterbi algorithm.)",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"a and b are the same stickiness sliders as markov-chains, but now driving " +
						"the hidden layer instead of a directly-visible one. The three columns are " +
						"the three observed days (Walk, Shop, Clean); the two rows are the hidden " +
						"states, Sunny on top and Rainy on bottom. Each circle's size is that " +
						"day/state's δ value — how likely the single best explanation ending there " +
						"is — normalized against the largest δ on the board. The 'Day revealed (t)' " +
						"slider reveals the trellis one day at a time: at t=0 only day 0's two " +
						"circles show; pushing t up reveals the next day's circles along with a " +
						"thin line from each to whichever previous-day state actually produced its " +
						"best δ (the 'surviving' backpointer, not every possible transition). At " +
						"t=2 the fully-revealed trellis also bolds the one path Viterbi actually " +
						"picked by backtracking from the largest final δ — Sunny→Rainy→Rainy in the " +
						"default a=0.6, b=0.7 case — in the same warm orange markov-chains used for " +
						"its highlighted day.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Recover the specific hidden sequence of states that most likely produced a " +
						"run of only-indirectly-related observations — not just 'it's probably " +
						"more sunny than rainy on average' the way a plain prior would say, but a " +
						"concrete day-by-day best guess (Sunny, Rainy, Rainy) backed by an exact " +
						"probability (0.01344), computed without ever seeing the weather directly. " +
						"That closes the loop section 1 opened: yes, a hidden state can be decoded " +
						"from noisy side-effects alone, as long as you know how the hidden state " +
						"evolves (the transition probabilities) and how it leaks into what you " +
						"observe (the emission probabilities).",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Speech recognition systems use an HMM where the hidden states are the " +
						"sounds (phonemes) someone intended to make and the observations are the " +
						"noisy audio actually recorded. Predictive keyboards and grammar checkers " +
						"tag each word's part of speech (noun, verb, ...) — hidden — from the " +
						"sequence of words actually typed — observed. Bioinformatics tools scan DNA " +
						"for genes by treating 'coding region' vs 'non-coding region' as a hidden " +
						"state and the actual sequence of A/C/G/T letters as the observation. Even " +
						"something as everyday as a fitness tracker guessing whether you're " +
						"walking, running, or asleep from noisy accelerometer readings is the same " +
						"shape of problem: infer a hidden state from an indirect, noisy trace of " +
						"it.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: the most likely hidden sequence is the one whole path " +
						"with the highest combined probability across every day at once — " +
						"transitions and emissions multiplied together — found by keeping, at each " +
						"day, only the single best route into each state and discarding the rest " +
						"(that's exactly what the δ update above does). Not like this: picking each " +
						"day's guess independently by whichever state's emission probability is " +
						"higher that day alone. Day 1's Shop observation has emission probability " +
						"0.40 under Rainy versus only 0.30 under Sunny — a day-by-day greedy guess " +
						"would flip day 1 over to Rainy. But the actual Viterbi decode keeps day 1 " +
						"on Sunny (δ(Sunny)=0.0432 > δ(Rainy)=0.0384), because it also weighs how " +
						"well day 1 connects to the strong Sunny evidence already banked from day " +
						"0 — a single day's isolated emission number can point the wrong way once " +
						"the whole sequence is considered together.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "P(stay Sunny)", Min: 0.05, Max: 0.95, Step: 0.05, Def: 0.6},
			{Key: "b", Label: "P(stay Rainy)", Min: 0.05, Max: 0.95, Step: 0.05, Def: 0.7},
			{Key: "t", Label: "Day revealed (t)", Min: 0, Max: 2, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// Hidden-state indices, in the same order as StateNames.
const (
	Sunny = 0
	Rainy = 1
)

// Observed-activity indices, in the same order as ObservationNames.
const (
	Walk  = 0
	Shop  = 1
	Clean = 2
)

// StateNames labels the two hidden states, indexed the same way as Sunny/Rainy.
var StateNames = []string{"Sunny", "Rainy"}

// ObservationNames labels the three observable activities, indexed the same
// way as Walk/Shop/Clean.
var ObservationNames = []string{"Walk", "Shop", "Clean"}

// Emission is the fixed emission matrix: Emission[state][observation] is
// P(observation | state). These are the textbook values -- a sunny day is
// six times more likely to produce a Walk than a rainy one is, and a rainy
// day is five times more likely to produce Cleaning.
var Emission = [2][3]float64{
	{0.6, 0.3, 0.1}, // Sunny
	{0.1, 0.4, 0.5}, // Rainy
}

// InitialDist is the fixed prior over the hidden state on day 0: P(Sunny),
// P(Rainy).
var InitialDist = [2]float64{0.4, 0.6}

// Sequence is the fixed 3-day observed activity sequence every Section
// walks through: Walk, Shop, Clean.
var Sequence = []int{Walk, Shop, Clean}

// Transition returns the 2x2 transition matrix for stickiness a (P(stay
// Sunny)) and b (P(stay Rainy)), using the exact same convention as
// markov-chains.Step: Transition[from][to].
func Transition(a, b float64) [2][2]float64 {
	return [2][2]float64{
		{a, 1 - a},
		{1 - b, b},
	}
}

// Viterbi runs the Viterbi dynamic-programming algorithm: for every day t
// and hidden state s, it finds delta[t][s], the probability of the single
// most likely path of hidden states ending in s at time t that is
// consistent with the observations seen so far. It returns that full delta
// table (for visualizing how the trellis fills in), the overall best
// hidden-state path (found by backtracking from the largest final-day
// delta), and that path's probability. obs must be non-empty; an empty obs
// returns a nil path and zero probability.
func Viterbi(trans [2][2]float64, emit [2][3]float64, initial [2]float64, obs []int) (delta [][2]float64, path []int, prob float64) {
	n := len(obs)
	if n == 0 {
		return nil, nil, 0
	}

	delta = make([][2]float64, n)
	backptr := make([][2]int, n)
	for s := 0; s < 2; s++ {
		delta[0][s] = initial[s] * emit[s][obs[0]]
		backptr[0][s] = -1 // no predecessor on day 0
	}
	for t := 1; t < n; t++ {
		for s := 0; s < 2; s++ {
			best, bestPrev := -1.0, 0
			for sp := 0; sp < 2; sp++ {
				v := delta[t-1][sp] * trans[sp][s]
				if v > best {
					best, bestPrev = v, sp
				}
			}
			delta[t][s] = best * emit[s][obs[t]]
			backptr[t][s] = bestPrev
		}
	}

	path = make([]int, n)
	if delta[n-1][Sunny] >= delta[n-1][Rainy] {
		path[n-1], prob = Sunny, delta[n-1][Sunny]
	} else {
		path[n-1], prob = Rainy, delta[n-1][Rainy]
	}
	for t := n - 2; t >= 0; t-- {
		path[t] = backptr[t+1][path[t+1]]
	}
	return delta, path, prob
}

// backpointers recomputes, for every day t>=1 and state s, which previous
// day's state actually produced the highest delta[t-1][sp]*trans[sp][s] --
// the single "surviving" predecessor arrow Viterbi kept while filling in
// delta[t][s]. It's derived directly from delta (already computed by
// Viterbi) purely for drawing the trellis; day 0 has no predecessor and is
// left zeroed.
func backpointers(delta [][2]float64, trans [2][2]float64) [][2]int {
	n := len(delta)
	plink := make([][2]int, n)
	for day := 1; day < n; day++ {
		for s := 0; s < 2; s++ {
			if delta[day-1][Sunny]*trans[Sunny][s] >= delta[day-1][Rainy]*trans[Rainy][s] {
				plink[day][s] = Sunny
			} else {
				plink[day][s] = Rainy
			}
		}
	}
	return plink
}

// columnX is the day-0/1/2 x positions on the unit-square canvas, and rowY
// the Sunny/Rainy y positions -- Sunny drawn on top, Rainy on the bottom,
// matching markov-chains's "sunny is the tracked state" convention.
var columnX = []float64{0.15, 0.5, 0.85}
var rowY = [2]float64{0.75, 0.25}

func render(p map[string]float64) string {
	a := p["a"]
	b := p["b"]
	t := int(p["t"] + 0.5)
	maxT := len(Sequence) - 1
	if t < 0 {
		t = 0
	}
	if t > maxT {
		t = maxT
	}

	trans := Transition(a, b)
	delta, path, prob := Viterbi(trans, Emission, InitialDist, Sequence)
	plink := backpointers(delta, trans)

	maxDelta := 0.0
	for _, d := range delta {
		for s := 0; s < 2; s++ {
			if d[s] > maxDelta {
				maxDelta = d[s]
			}
		}
	}

	c := viz.New(680, 460, 0, 1, 0, 1)

	// The surviving backpointer arrow into each revealed day -- the single
	// predecessor delta actually kept, not every possible transition.
	for day := 1; day <= t; day++ {
		for s := 0; s < 2; s++ {
			from := [2]float64{columnX[day-1], rowY[plink[day][s]]}
			to := [2]float64{columnX[day], rowY[s]}
			c.Path([][2]float64{from, to}, viz.Muted, 1.5)
		}
	}

	// Once every day is revealed, bold the one path Viterbi actually
	// decoded by backtracking from the largest final-day delta.
	if t == maxT {
		pathPts := make([][2]float64, len(path))
		for day, s := range path {
			pathPts[day] = [2]float64{columnX[day], rowY[s]}
		}
		c.Path(pathPts, viz.Warm, 3)
	}

	stateColor := [2]string{viz.Accent, viz.Bad} // Sunny, Rainy
	for day := 0; day <= t; day++ {
		c.Text(c.X(columnX[day]), 28, ObservationNames[Sequence[day]], 13, viz.Muted, "middle")
		for s := 0; s < 2; s++ {
			px, py := c.X(columnX[day]), c.Y(rowY[s])
			side := 14.0
			if maxDelta > 0 {
				side = 14 + 46*(delta[day][s]/maxDelta)
			}
			color := stateColor[s]
			if t == maxT && path[day] == s {
				color = viz.Warm
			}
			c.Rect(px-side/2, py-side/2, side, side, color, 0.85)
			c.Text(px, py+side/2+16, fmt.Sprintf("%.4f", delta[day][s]), 11, viz.Muted, "middle")
		}
	}

	c.Text(16, c.H-52, StateNames[Sunny], 12, viz.Accent, "start")
	c.Text(16, c.H-32, StateNames[Rainy], 12, viz.Bad, "start")
	c.Text(16, 44, fmt.Sprintf("P(stay Sunny)=%.2f    P(stay Rainy)=%.2f", a, b), 14, viz.Ink, "start")

	if t == maxT {
		names := make([]string, len(path))
		for i, s := range path {
			names[i] = StateNames[s]
		}
		c.Text(16, 64, fmt.Sprintf("most likely path: %s (probability %.5f)", strings.Join(names, " → "), prob),
			14, viz.Warm, "start")
	} else {
		c.Text(16, 64, "drag the day slider to reveal each day's trellis and see the decode", 12, viz.Muted, "start")
	}

	return c.String()
}
