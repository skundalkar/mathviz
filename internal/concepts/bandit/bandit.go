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
					"expected-value showed how to judge a single gamble once you already know its " +
						"odds. But say you walk up to three slot machines -- A, B, C -- and none of " +
						"them come with their win probabilities printed on the front. The only way to " +
						"learn one is to spend a pull actually playing it, and you only have a limited " +
						"number of pulls before you have to walk away.",
					"Here's the trap either extreme falls into. Play only whichever machine looks " +
						"best so far, based on very few pulls: an early unlucky streak on the actual " +
						"best machine can make it look like the worst one, and a purely 'stick with " +
						"the leader' rule then never gives it another chance to prove otherwise -- one " +
						"bad first impression, permanently misjudged. Or play all three equally the " +
						"whole time, just to be thorough: now you're spending just as many pulls on " +
						"machines you've already gathered good evidence are worse, when those pulls " +
						"could have gone to the one that's actually paying off.",
					"Is there a rule for choosing which machine to pull next that keeps testing " +
						"enough to avoid the first trap, without wasting pulls the way the second one " +
						"does?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Fix three machines with true win probabilities A=0.20, B=0.50, C=0.35 -- kept " +
						"secret from the player, revealed here only to grade the strategy. Each pull " +
						"pays 1 on a win and 0 on a loss, so a machine's true win probability IS its " +
						"expected value per pull (E[X] = p\u00b71 + (1\u2212p)\u00b70 = p, the same formula " +
						"expected-value defined). Epsilon-greedy: pull each machine once to start (a " +
						"forced pull, so every machine has at least one data point), then for every " +
						"later pull, explore with probability \u03b5 (a uniformly random machine) or " +
						"otherwise exploit (whichever machine currently has the highest observed win " +
						"rate, ties toward the earliest one).",
					"\u2022 Pull 1 (forced): A -- wins. A's estimate: 1/1 = 1.00.",
					"\u2022 Pull 2 (forced): B -- loses. B's estimate: 0/1 = 0.00.",
					"\u2022 Pull 3 (forced): C -- wins. C's estimate: 1/1 = 1.00.",
					"Already, the first trap from Section 1 is visible: B is secretly the BEST " +
						"machine (true rate 0.50), but its one data point was a loss, so it now looks " +
						"like the worst of the three. A player with \u03b5=0 (pure exploit, never explore) " +
						"would never touch B again -- nothing changes its estimate away from 0.00 " +
						"except another pull, and a purely greedy rule never chooses a machine that " +
						"isn't currently in the lead. At \u03b5=0 in this exact simulation, that's exactly " +
						"what happens: B is pulled exactly once, forever, no matter how many pulls " +
						"follow.",
					"At the default \u03b5=0.20, pulls 4-14 all happen to exploit (A or C, whichever is " +
						"ahead) -- until pull 15, which explores: a uniformly random draw lands on B. " +
						"B wins, its estimate jumps to 1/2 = 0.50, and it's back in contention. That's " +
						"exactly what \u03b5 buys: without pull 15's exploration, B would never have gotten " +
						"a second chance.",
					"|pulls so far|A pulls (true 0.20)|B pulls (true 0.50)|C pulls (true 0.35)|B's estimate|",
					"|60|25|9|26|0.33|",
					"|150|34|49|67|0.37|",
					"|300|45|155|100|0.46|",
					"|1000|91|758|151|0.53|",
					"By pull 60, B has barely been explored (9 pulls) and C's early lucky wins have " +
						"it looking best. By pull 150, more exploring pulls have found B more often, " +
						"but C (67 pulls, estimate 0.40) still edges it out by chance. Only by pull " +
						"300 does the noise wash out enough for B's true 0.50 rate to show through -- " +
						"155 pulls to B versus 100 to C -- and by pull 1000 B dominates outright. " +
						"Nothing about the RULE changed between pull 60 and pull 1000; more pulls just " +
						"gave the sampling noise more chances to average out, the same " +
						"shrinking-uncertainty effect law-of-large-numbers describes for one running " +
						"average, now racing across three of them at once.",
					"Regret measures the running cost of not always pulling the best machine: each " +
						"pull adds (best true rate \u2212 chosen arm's true rate) to a cumulative total, " +
						"0.30 for a pull of A, 0 for a pull of B, 0.15 for a pull of C. In this run, " +
						"cumulative regret is 11.40 by pull 60 (0.19/pull), 28.50 by pull 300 " +
						"(0.095/pull), and 49.95 by pull 1000 (0.05/pull) -- still growing, but slower " +
						"and slower per pull as the strategy locks onto B, which is exactly the " +
						"signature of a strategy that's working.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Three bars, one per machine, height equal to its current estimated win rate; " +
						"the thin black line across each bar marks that machine's TRUE win rate -- the " +
						"target the bar is trying to converge to. The bar just pulled lights up orange " +
						"(this pull explored) or green (this pull exploited); the step slider scrubs " +
						"through pulls 1-300 of the exact run worked above. Drag it past pull 15 and " +
						"watch B's bar jump for the first time since pull 2; drag it toward 300 and " +
						"watch B's bar climb past C's as the counts in the table above play out live. " +
						"Below, the cumulative regret curve only ever rises, but its slope visibly " +
						"flattens as more pulls go to B instead of A or C.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Choose defensibly among several options of unknown quality under a limited " +
						"budget of trials, without falling into either trap Section 1 named: never " +
						"permanently writing off an option that had one unlucky early result, and " +
						"never spending equal effort on options you already have good evidence are " +
						"worse. The \u03b5 dial makes that trade-off explicit and tunable instead of " +
						"implicit: raise it and you find the best option faster on average (more " +
						"exploring pulls, like B's early rescue at pull 15) but permanently spend more " +
						"pulls on options you already know are worse; lower it and you exploit your " +
						"current best guess more often, at the risk of an early bad impression sticking " +
						"the way it did for B at \u03b5=0. There's no setting that avoids the trade-off " +
						"entirely -- only ways to match it to how many pulls you actually have.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"The name comes from exactly this picture: a row of casino slot machines " +
						"('one-armed bandits'), each with an unknown payout rate, and a gambler with a " +
						"limited bankroll deciding which lever to keep pulling. The same everyday " +
						"decision shows up outside a casino: sticking with your usual lunch spot versus " +
						"trying the new place down the street, or always taking the same route to work " +
						"versus occasionally testing an alternate one that might be faster. At larger " +
						"scale, a website deciding which of several headlines or ad creatives to show " +
						"more often as click data comes in, a streaming service deciding which shows to " +
						"promote to more viewers, and a clinical trial deciding how to allocate the " +
						"next patient between two treatments as early results come in are all running " +
						"this exact explore/exploit trade-off, usually with more than three arms and " +
						"smarter strategies than plain epsilon-greedy.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: epsilon-greedy lowers the AVERAGE regret over many pulls and " +
						"finds the best arm often enough given enough exploration -- it does not " +
						"guarantee that after any fixed, finite number of pulls, whichever arm " +
						"currently looks best (or has been pulled the most) actually IS the best one.",
					"Not like this: assuming whichever option 'looks best' after a modest amount of " +
						"testing must be the true best -- this exact run is the counterexample, at " +
						"pull 150 C (not B) had the highest estimate, purely from lucky results. The " +
						"fix isn't to stop testing once something looks good; it's to keep enough " +
						"exploration going that a close-looking alternative still gets tested before " +
						"you commit to it.",
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
