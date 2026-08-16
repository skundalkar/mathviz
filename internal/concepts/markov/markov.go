// Package markov visualizes a two-state Markov chain: a simple sunny/rainy
// weather model where tomorrow's weather depends only on today's, not on an
// independent coin flip. Starting from a fixed initial state, the
// probability of being sunny on day t converges to a steady-state value --
// and the gap to that steady state shrinks by the exact same fixed factor
// every single step, unlike the noisy convergence of law-of-large-numbers or
// monte-carlo-estimation.
package markov

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "markov-chains",
		Seq:   41,
		Title: "Markov chains (weather that remembers yesterday)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"law-of-large-numbers and monte-carlo-estimation both lean on the same " +
						"assumption: each trial is independent of the last one — a coin flip has no " +
						"memory, and neither does a random dart thrown at a square. Weather doesn't " +
						"work that way. A sunny day is much more likely to be followed by another " +
						"sunny day than a fresh coin flip would predict — today's weather visibly " +
						"depends on yesterday's. If tomorrow's outcome depends on today's instead of " +
						"being drawn fresh and independent, does the whole idea of 'a stable " +
						"long-run average' break down — or is there still a predictable long-run " +
						"fraction of sunny days, just reached a different way?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Build the simplest version: two states, Sunny and Rainy, with fixed " +
						"one-step transition odds — P(sunny tomorrow | sunny today) = a = 0.9, and " +
						"P(rainy tomorrow | rainy today) = b = 0.5. Start on a day you know for " +
						"certain is Rainy (P(sunny) = 0) and update day by day using P(sunny)_{t+1} " +
						"= P(sunny)_t·a + P(rainy)_t·(1-b):",
					"• t=0: P(sunny) = 0.0000 (certain rain). t=1: 0.5000. t=2: 0.7000.",
					"• t=3: 0.7800. t=4: 0.8120. t=5: 0.8248.",
					"• t=10: 0.8332. t=20 and t=30: 0.8333 — settled.",
					"That 0.8333 is the steady state: solving P(sunny) = P(sunny)·a + (1-P(sunny))" +
						"·(1-b) for a fixed point gives steady-state P(sunny) = (1-b)/(2-a-b) = " +
						"0.5/0.6 = 0.8333, no simulation required. And unlike a coin-flip average's " +
						"noisy wobble, the gap to that steady state shrinks by the exact same " +
						"factor λ = a+b-1 = 0.4 every single step: gap at t=0 is -0.8333, at t=1 is " +
						"-0.3333 (ratio 0.400), at t=2 is -0.1333 (ratio 0.400 again), at t=3 is " +
						"-0.0533 (ratio 0.400) — a clean geometric decay, not a bouncing one.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The a and b sliders set how 'sticky' each state is — how likely sunny stays " +
						"sunny, and rainy stays rainy. The dashed line is the steady-state P(sunny) " +
						"those two numbers imply; the solid curve is the actual day-by-day P(sunny) " +
						"starting from a certain-rain day 0. The t slider marks a specific day on " +
						"that curve, with a readout of the exact probability and its remaining gap " +
						"to the steady state. Push a and b both toward 1 (very sticky weather) and " +
						"the curve visibly takes longer to flatten out; pull them toward 0.5 and it " +
						"snaps to the steady state almost immediately.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Predict a system's long-run behavior directly from its local, one-step rules " +
						"— no simulation, no waiting years of data — by solving for the fixed point " +
						"of the transition probabilities, the same way this page solves for 0.8333 " +
						"from just a and b. You can also predict *how fast* that long-run behavior " +
						"kicks in: since the gap shrinks by exactly λ = a+b-1 each step, you can say " +
						"in advance how many days it takes to get within any given tolerance of the " +
						"steady state, rather than discovering it by trial and error.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Google's original PageRank algorithm modeled a web surfer randomly clicking " +
						"links as a Markov chain, and ranked pages by their steady-state probability " +
						"of being visited. Predictive-text and early chatbots generate the next word " +
						"from a Markov chain over the current word (or last few words), not from " +
						"scratch each time. Board games with dice and fixed squares (Monopoly's " +
						"famous long-run bias toward Illinois Avenue and jail) are Markov chains over " +
						"board positions. Credit-rating agencies model a company's rating (AAA, AA, " +
						"..., default) as a Markov chain to estimate long-run default risk. Even a " +
						"basic weather forecast model, like the toy one here, is a real (if much " +
						"simplified) technique used in climate and queueing models.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: the long-run fraction of time spent in each state is fixed " +
						"entirely by the one-step transition probabilities, and (for a chain like " +
						"this one) the same steady state is reached no matter which state you start " +
						"in. Not like this: assuming that because a chain 'remembers' yesterday, it " +
						"can never settle into a single stable long-run average the way independent " +
						"trials do — it still does, just via geometric decay toward a fixed point " +
						"instead of the noisy shrinking wobble of an average of independent draws. " +
						"Also not like this: confusing the steady-state probability (the long-run " +
						"*fraction of days* that turn out sunny, 0.8333 above) with the one-step " +
						"transition probability (the chance tomorrow is sunny *given* today already " +
						"is, a=0.9 above) — they're two different numbers computed from the same " +
						"matrix, not the same number under two names.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "P(stay sunny)", Min: 0.05, Max: 0.95, Step: 0.05, Def: 0.9},
			{Key: "b", Label: "P(stay rainy)", Min: 0.05, Max: 0.95, Step: 0.05, Def: 0.5},
			{Key: "t", Label: "Day (t)", Min: 0, Max: 30, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// SteadyState returns the long-run probability of being Sunny and Rainy for
// a two-state chain with P(stay sunny)=a and P(stay rainy)=b: the fixed
// point of the one-step update, solved directly rather than simulated.
// a+b==2 (both states absorbing, no way to leave either) is degenerate and
// has no unique fixed point; that combination never occurs given the
// Params' Max of 0.95 each, so it isn't guarded against here.
func SteadyState(a, b float64) (sunny, rainy float64) {
	denom := 2 - a - b
	sunny = (1 - b) / denom
	rainy = (1 - a) / denom
	return sunny, rainy
}

// Step advances one day: given today's probability of being sunny and the
// chain's stickiness a (P(stay sunny)) and b (P(stay rainy)), returns
// tomorrow's probability of being sunny.
func Step(pSunny, a, b float64) float64 {
	return pSunny*a + (1-pSunny)*(1-b)
}

// Trajectory returns the probability of being sunny on each of days
// 0..steps, starting from pSunny0 on day 0 and applying Step once per day.
// steps < 0 is treated as 0.
func Trajectory(pSunny0, a, b float64, steps int) []float64 {
	if steps < 0 {
		steps = 0
	}
	out := make([]float64, steps+1)
	out[0] = pSunny0
	for t := 1; t <= steps; t++ {
		out[t] = Step(out[t-1], a, b)
	}
	return out
}

func render(params map[string]float64) string {
	_ = params
	return viz.New(680, 420, 0, 30, 0, 1).String()
}
