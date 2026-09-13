// Package bandit visualizes the multi-armed bandit problem: repeatedly
// choosing among several options with unknown reward rates, where every pull
// spent finding out which option is best is a pull that couldn't go to the
// option you already suspect is best. Epsilon-greedy is the simplest
// strategy that faces that trade-off head-on: pull a uniformly random arm
// with probability epsilon (explore), otherwise pull whichever arm currently
// looks best (exploit).
package bandit

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "multi-armed-bandit",
		Seq:   102,
		Title: "Multi-armed bandit (explore vs. exploit)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder -- filled in once the math and picture exist.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "epsilon", Label: "Exploration rate (ε)", Min: 0, Max: 1, Step: 0.05, Def: 0.2},
			{Key: "step", Label: "Pulls so far (step)", Min: 0, Max: MaxSteps, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// ArmNames labels the three slot machines every Section walks through.
var ArmNames = []string{"A", "B", "C"}

// TrueProbs is each arm's true, hidden win probability -- its expected
// reward per pull (a Bernoulli payoff of 1 on a win, 0 on a loss), in the
// same sense expected-value gave a probability-weighted average payoff.
// Arm B is deliberately the best (0.50), but not so far ahead of C (0.35)
// that a handful of unlucky pulls couldn't make C look better for a while --
// exactly the situation that makes exploration worth its cost.
var TrueProbs = []float64{0.20, 0.50, 0.35}

// MaxSteps is the number of pulls every Section's step slider scrubs
// through -- large enough that, at the default exploration rate, the
// picture visibly flips from "a worse arm looks best" to "the true best
// arm dominates" as the slider is dragged to its far end (see
// LEARNINGS.md).
const MaxSteps = 300

// hash01 turns an integer seed into a deterministic pseudo-random number in
// [0, 1) via a fixed irrational-multiplier trick -- not a real random number
// generator, just a fixed, reproducible sequence that merely looks
// patternless, which is exactly what Render needs (pure: same input, same
// output, no math/rand, no time). Same trick as internal/concepts/lln.
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// ArmOutcome returns whether the (occurrence+1)-th pull of the given arm
// (occurrence is 0-indexed: how many times that arm has already been pulled
// before this one) is a win. It's keyed by (arm, occurrence) rather than by
// an overall pull counter, so the same arm always produces the same
// sequence of wins/losses regardless of which other arms were pulled in
// between -- pulling arm B five times in a row and pulling it once every
// twelve pulls both see the same five outcomes, in the same order.
func ArmOutcome(arm, occurrence int) bool {
	return hash01(arm*1_000_003+occurrence) < TrueProbs[arm]
}

// exploreDecision reports whether pull t (0-indexed, only ever called for
// t >= len(TrueProbs)) explores: a fixed pseudo-random draw compared to
// epsilon, independent of exploreArm below so the two decisions don't
// correlate.
func exploreDecision(t int, epsilon float64) bool {
	return hash01(t*104_729+7) < epsilon
}

// exploreArm picks which arm an exploring pull lands on: a uniformly random
// arm among ArmNames, independent of exploreDecision.
func exploreArm(t int) int {
	idx := int(hash01(t*7_919+13) * float64(len(ArmNames)))
	if idx >= len(ArmNames) {
		idx = len(ArmNames) - 1
	}
	return idx
}

// argmaxEstimate returns the arm with the highest current estimated mean
// (sums[i]/counts[i]), breaking ties by lowest index -- the "exploit"
// choice.
func argmaxEstimate(sums []float64, counts []int) int {
	best, bestMean := 0, math.Inf(-1)
	for i := range sums {
		mean := 0.0
		if counts[i] > 0 {
			mean = sums[i] / float64(counts[i])
		}
		if mean > bestMean {
			best, bestMean = i, mean
		}
	}
	return best
}

// Step is a single pull's outcome and the running state right after it, so
// Render can show exactly what changed at the current step without
// re-deriving it from Steps' full history.
type Step struct {
	Arm       int        // which arm was pulled
	Explored  bool       // true if this pull explored (random arm), false if it exploited (best estimate) or was a forced first pull
	Reward    int        // 0 or 1, this pull's payoff
	Counts    [3]int     // pulls per arm, including this one
	Estimates [3]float64 // sums[i]/counts[i] per arm, including this one (0 if unpulled)
	CumRegret float64    // sum so far of (best true prob - chosen arm's true prob), the cost of not always pulling the best arm
}

// Simulate runs `steps` pulls of epsilon-greedy bandit play against
// TrueProbs and returns one Step per pull, in order. The first
// len(TrueProbs) pulls are forced -- one pull per arm, in order -- so every
// arm has at least one estimate before epsilon-greedy has to choose among
// them; after that, pull t explores with probability epsilon (a uniformly
// random arm) or exploits (the arm with the current highest estimate,
// ties to the lowest index). Deterministic and pure: same (epsilon, steps)
// always reproduces the same sequence.
func Simulate(epsilon float64, steps int) []Step {
	if steps < 0 {
		steps = 0
	}
	var counts [3]int
	var sums [3]float64
	best := TrueProbs[0]
	for _, p := range TrueProbs {
		if p > best {
			best = p
		}
	}

	results := make([]Step, steps)
	cumRegret := 0.0
	for t := 0; t < steps; t++ {
		var arm int
		explored := false
		switch {
		case t < len(TrueProbs):
			arm = t // forced first pull of each arm
		case exploreDecision(t, epsilon):
			explored = true
			arm = exploreArm(t)
		default:
			arm = argmaxEstimate(sums[:], counts[:])
		}

		occurrence := counts[arm]
		win := ArmOutcome(arm, occurrence)
		reward := 0
		if win {
			reward = 1
		}
		counts[arm]++
		sums[arm] += float64(reward)
		cumRegret += best - TrueProbs[arm]

		var estimates [3]float64
		for i := range estimates {
			if counts[i] > 0 {
				estimates[i] = sums[i] / float64(counts[i])
			}
		}
		results[t] = Step{
			Arm:       arm,
			Explored:  explored,
			Reward:    reward,
			Counts:    counts,
			Estimates: estimates,
			CumRegret: cumRegret,
		}
	}
	return results
}

// barColors: default fill, then fill when this arm is the one just pulled,
// split by whether that pull explored or exploited.
const (
	barDefault = viz.Accent
	barExplore = viz.Warm
	barExploit = viz.Good
)

func render(p map[string]float64) string {
	epsilon := p["epsilon"]
	if epsilon < 0 {
		epsilon = 0
	}
	if epsilon > 1 {
		epsilon = 1
	}
	step := int(p["step"] + 0.5)
	if step < 0 {
		step = 0
	}
	if step > MaxSteps {
		step = MaxSteps
	}

	all := Simulate(epsilon, MaxSteps)
	maxRegret := all[len(all)-1].CumRegret
	if maxRegret < 1 {
		maxRegret = 1
	}

	// The line chart lives in the lower band of the canvas (PadT pushed
	// down past the bar chart and header text above it); the bar chart is
	// drawn in raw pixel space, same technique dynamic-programming-knapsack
	// uses for its table grid, so the two panels never fight over one
	// coordinate system.
	c := viz.New(700, 520, 0, float64(MaxSteps), 0, maxRegret*1.15)
	c.PadT = 340
	c.PadB = 40
	c.Axes()
	for x := 0; x <= MaxSteps; x += 50 {
		c.Tick(float64(x), fmt.Sprintf("%d", x))
	}
	c.Text(52, c.PadT-14, "cumulative regret so far (cost of not always pulling the best arm)", 12, viz.Muted, "start")

	const barTop, barBottom = 110.0, 280.0
	const barW, gap = 90.0, 90.0
	startX := (700.0 - 3*barW - 2*gap) / 2

	var current *Step
	if step > 0 {
		current = &all[step-1]
	}

	for i, name := range ArmNames {
		x := startX + float64(i)*(barW+gap)
		est := 0.0
		count := 0
		if current != nil {
			est = current.Estimates[i]
			count = current.Counts[i]
		}
		barH := est * (barBottom - barTop)
		color, opacity := barDefault, 0.6
		if current != nil && current.Arm == i {
			opacity = 0.95
			if current.Explored {
				color = barExplore
			} else if step > len(TrueProbs) {
				color = barExploit
			}
		}
		c.Rect(x, barBottom-barH, barW, barH, color, opacity)

		// True win rate, shown as a thin reference line spanning a little
		// past the bar's edges -- what the estimate is trying to converge
		// to, kept visible whether it falls inside a tall bar or floats
		// above a short one.
		trueY := barBottom - TrueProbs[i]*(barBottom-barTop)
		c.Rect(x-6, trueY-1.5, barW+12, 3, viz.Ink, 1)

		c.Text(x+barW/2, barTop-10, fmt.Sprintf("%s: true p=%.2f", name, TrueProbs[i]), 12, viz.Muted, "middle")
		label := "unpulled"
		if count > 0 {
			label = fmt.Sprintf("n=%d  est=%.2f", count, est)
		}
		c.Text(x+barW/2, barBottom+18, label, 12, viz.Ink, "middle")
	}

	// The regret curve so far -- only the prefix up to the current step, so
	// dragging the slider visibly grows the line, and its cost only ever
	// increases (a flat line would mean an oracle that always pulled the
	// best arm).
	pts := make([][2]float64, step+1)
	pts[0] = [2]float64{0, 0}
	for i := 0; i < step; i++ {
		pts[i+1] = [2]float64{float64(i + 1), all[i].CumRegret}
	}
	c.Path(pts, viz.Accent, 2.5)

	if step == 0 {
		c.Text(16, 24, fmt.Sprintf("Step 0: no pulls yet -- epsilon = %.2f", epsilon), 13, viz.Ink, "start")
	} else {
		s := all[step-1]
		verdict := "forced first pull"
		if step > len(TrueProbs) {
			if s.Explored {
				verdict = "explored: uniformly random arm"
			} else {
				verdict = "exploited: highest current estimate"
			}
		}
		outcome := "lost"
		if s.Reward == 1 {
			outcome = "won"
		}
		c.Text(16, 24, fmt.Sprintf("Pull %d/%d: arm %s (%s) -- %s", step, MaxSteps, ArmNames[s.Arm], verdict, outcome),
			13, viz.Ink, "start")
		c.Text(16, 44, fmt.Sprintf("cumulative regret = %.2f (avg %.3f/pull)    epsilon = %.2f",
			s.CumRegret, s.CumRegret/float64(step), epsilon), 13, viz.Muted, "start")
	}
	c.Text(16, 64, "orange bar = just explored   green bar = just exploited   thin gray line = arm's true win rate",
		12, viz.Muted, "start")

	return c.String()
}
