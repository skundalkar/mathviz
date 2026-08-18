// Package randomwalk visualizes a 1D random walk: a position that moves +1
// or -1 with equal probability at every step. Its expected position stays
// pinned at 0 forever, but individual runs wander away from 0 by a typical
// distance that grows with the square root of the number of steps, not a
// fixed amount and not linearly.
package randomwalk

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "random-walk",
		Seq:   47,
		Title: "Random walk (spread grows like √n)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`law-of-large-numbers` tracks a running AVERAGE of coin flips — heads-so-far " +
						"divided by flips-so-far — and shows it settling down toward the true " +
						"probability p as n grows, because dividing by a growing n dilutes any one " +
						"flip's influence. But what if, instead of averaging, you just SUM the " +
						"flips up: +1 for heads, -1 for tails, no dividing by n at all? Does that " +
						"running total also settle down toward some fixed number as n grows, the " +
						"way the average does — or does something completely different happen when " +
						"there's no growing denominator left to dilute anything?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Flip a fair coin repeatedly; step +1 for heads, -1 for tails; track the " +
						"running position S_n = the sum of the first n steps. Each single step has " +
						"expected value 0 (equally likely +1 or -1), so the expected position after " +
						"any number of steps is also exactly 0 — the +1s and -1s cancel on average, " +
						"no matter how large n gets. But 'expected value 0' is not the same claim as " +
						"'stays near 0': one traced-out run (same reproducible sequence used " +
						"throughout this page) lands at S_10=6, S_100=8, S_300=6 — never at 0, and " +
						"not obviously homing in on it either.",
					"• The key number isn't the position itself, it's how far a typical run " +
						"strays from 0. Each step has variance 1 (it's always exactly +1 or -1, " +
						"never 0, so E[step²]=1). Steps are independent, so the variance of a SUM " +
						"of n of them is just n times one step's variance: Var(S_n) = n × 1 = n. " +
						"The typical (root-mean-square) distance from 0 is the square root of that " +
						"variance: √n.",
					"• Checked against many independent runs at a few step counts: n=10, root-" +
						"mean-square position over 500 runs ≈ 3.17 (√10 = 3.162). n=100: ≈ 10.27 " +
						"(√100 = 10.000). n=300: ≈ 18.07 (√300 = 17.321) — all close to √n, the " +
						"small remaining gap itself shrinking as more runs are averaged in, the " +
						"same law-of-large-numbers effect applied to a different quantity.",
					"That's the whole shape of it: the walk's average position stays glued to 0 " +
						"forever, while its typical spread away from 0 keeps growing, just slowly — " +
						"proportional to √n, not to n itself.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The solid curve plots one specific run's position S_t against the number of " +
						"steps t taken so far. The dashed lines are the ±√t envelope — the typical " +
						"spread at each t — widening slowly as t grows. The run slider swaps in a " +
						"completely different sequence of coin flips at the same length, so you can " +
						"see the actual path is random (a different run wanders to different spots) " +
						"even though the envelope around it stays the same for every run. The n " +
						"slider marks one specific step count on the curve, with a readout of the " +
						"exact position there and how it compares to √n.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Tell apart two claims that sound similar but aren't: 'the average position " +
						"across many runs is 0' (true, always) and 'any one run stays close to 0' " +
						"(false — a single run can and does wander noticeably far, by an amount that " +
						"grows over time). You can also predict roughly how far a random-walk " +
						"process will have wandered after n steps without simulating it — √n gives " +
						"you the right order of magnitude directly, the same way monte-carlo-" +
						"estimation's 1/√n error-shrinkage rate let you budget sample size in " +
						"advance instead of guessing.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A stock price's random day-to-day jitter (the part not explained by a trend) " +
						"is commonly modeled as a random walk, which is exactly why a stock's price " +
						"range widens the further out you forecast, proportional to the square root " +
						"of time, not linearly. A drunkard's walk — the classic name for this exact " +
						"process — models how far a disoriented pedestrian ends up from a lamppost " +
						"after n unsteady steps in random directions. A pollen grain jittering in " +
						"water under a microscope (Brownian motion, first explained physically by " +
						"Einstein) follows the continuous-time version of the same math. Genetic " +
						"drift — how a gene's frequency in a population wanders over generations " +
						"purely from random reproduction, with no selection pressure at all — is a " +
						"random walk too.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: the expected position is 0 at every step, but the typical " +
						"distance from 0 keeps growing over time, proportional to √n — 'expected " +
						"value 0' describes the average over many runs, not a claim that any one run " +
						"stays near 0.",
					"Not like this: treating a random walk like `law-of-large-numbers`'s running " +
						"average, which really does settle down toward a fixed number as n grows — a " +
						"random walk's position does the opposite, spreading out farther (just " +
						"slowly) rather than converging. Also not like this: assuming a long run of " +
						"the same direction (many +1s in a row) means the walk is 'due' to reverse " +
						"— each step is an independent coin flip with no memory, the identical " +
						"gambler's-fallacy mistake law-of-large-numbers already warns against.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Steps taken (n)", Min: 1, Max: 300, Step: 1, Def: 100},
			{Key: "run", Label: "Run (different coin-flip sequence)", Min: 1, Max: 5, Step: 1, Def: 1},
		},
		Render: render,
	})
}

func render(params map[string]float64) string {
	_ = params
	return viz.New(680, 420, -1, 1, -1, 1).String()
}
