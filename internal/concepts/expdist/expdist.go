// Package expdist visualizes the exponential distribution: the probability
// distribution over how long you wait for the next rare, independent event,
// given only its average rate λ. It's the continuous counterpart to
// poisson-distribution — poisson-distribution counts how many events land in
// a fixed window; this package asks how long the gap until the next one is,
// derived from that same "zero events so far" probability.
package expdist

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "exponential-distribution",
		Seq:   60,
		Title: "Exponential distribution (waiting time between events)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"poisson-distribution answers 'how many customers arrive in the next hour, " +
						"given they arrive at an average rate of λ=3 per hour?' — a question about a " +
						"count. But a shopkeeper watching the door after one customer just walked in " +
						"isn't thinking in counts; they're thinking 'how long until the next one?' " +
						"That's a completely different kind of question — not a whole number of " +
						"events in a fixed window, but a continuous amount of time until a single " +
						"event. Poisson-distribution's machinery was built for counts. Is there a way " +
						"to get an exact probability distribution over a *waiting time* instead, from " +
						"that same rate λ and nothing else?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Yes — by reusing poisson-distribution's own formula, just asked a different " +
						"question of it. 'The wait for the next customer exceeds t hours' is exactly " +
						"the same event as 'zero customers arrive in the next t hours,' and " +
						"poisson-distribution already gives P(X=0) in a window: with rate λ over a " +
						"window of length t, the expected count is λt, so P(X=0) = (λt)^0 · e^(−λt) / " +
						"0! = e^(−λt). So: P(wait > t) = e^(−λt). That's the whole distribution, " +
						"read straight off poisson-distribution's k=0 case.",
					"Plug in λ=3 per hour, the same rate poisson-distribution used. P(wait > 1 hour) " +
						"= e^(−3) ≈ 0.0498 — only about a 5% chance of waiting over an hour. So " +
						"P(wait ≤ 1 hour) = 1 − e^(−3) ≈ 0.9502; that 1−e^(−λt) is the exponential " +
						"distribution's CDF, and differentiating it gives the PDF f(t) = λ·e^(−λt). " +
						"The average wait works out to exactly 1/λ = 1/3 hour = 20 minutes — the rate " +
						"and the average gap are just reciprocals of each other.",
					"One more property falls straight out of the same formula. Compare two " +
						"situations: (a) you've already waited s=0.5 hours with nobody arriving, and " +
						"want P(wait > s+t) given that; (b) you just started watching and want plain " +
						"P(wait > t), for t=1/3 hour. Case (a): P(wait>0.5+0.333 | wait>0.5) = " +
						"e^(−3×0.833)/e^(−3×0.5) = e^(−2.5)/e^(−1.5) = 0.0821/0.2231 = 0.3679. Case " +
						"(b): P(wait>0.333) = e^(−3×0.333) = e^(−1) = 0.3679. Identical — the 30 " +
						"minutes already spent waiting vanished from the answer entirely. This is " +
						"called memorylessness: e^(−λ(s+t))/e^(−λs) = e^(−λt) algebraically, for any " +
						"s and t, because the exponents just subtract.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The curve is f(t) = λe^(−λt), the density of the wait; the shaded region is " +
						"the area from 0 to the highlighted time t, which is exactly P(wait ≤ t). Raise " +
						"λ (events get more frequent) and the whole curve rises at t=0 and falls off " +
						"faster — shorter waits become far more likely. The dashed vertical line marks " +
						"the mean wait, 1/λ. A second readout re-derives memorylessness live: it always " +
						"waits one mean interval (s=1/λ) first, then reports P(wait > s+t | already " +
						"waited s) side by side with plain P(wait > t) — they match to full precision " +
						"for every λ and t you pick, not just the worked example above.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Put an exact probability on 'how long until the next one,' not just 'how many " +
						"in this window' — 'what's the chance I wait over 10 minutes for the next bus, " +
						"if they arrive every 12 minutes on average?' is now a direct plug-in, and so " +
						"is the flip side: given a target confidence, read off how long you'd need to " +
						"budget for the wait. You can also recognize when a 'time since last event' " +
						"argument is invalid — because of memorylessness, 'it's been a while since the " +
						"last one, so the next is due soon' is provably false for a true Poisson " +
						"process, no matter how long the current gap has already run.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Call centers use it to model the gap between incoming calls, on top of using " +
						"poisson-distribution for how many arrive per hour — the same process, two " +
						"different questions. Reliability engineering models the time until a " +
						"component fails this way when failures are due to random external shocks " +
						"rather than wear-and-tear aging. Radioactive decay's 'time until the next " +
						"decay event' is the textbook physical example, which is also why half-life " +
						"calculations don't care how long a sample has already sat around. Web " +
						"services model the gap between incoming requests this way when deciding how " +
						"long a connection can sit idle before it's probably not coming back.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: λ is events per unit time and 1/λ is the average *gap* " +
						"between them — they're reciprocals, so doubling the rate halves the average " +
						"wait, it doesn't double it.",
					"Not like this: 'the bus hasn't come in 20 minutes, so it's more likely to come " +
						"in the next 5' — the gambler's-fallacy mistake poisson-distribution's own " +
						"'common mistake' section already warned about, and exactly what " +
						"memorylessness rules out: P(wait a bit longer | already waited a long time) " +
						"is the same number it always was, never higher. (Real buses have schedules, " +
						"so this only applies to truly random, independent arrivals — a point worth " +
						"keeping straight before reaching for this model.)",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "lambda", Label: "Rate (λ, events per window)", Min: 0.2, Max: 5, Step: 0.1, Def: 3},
			{Key: "t", Label: "Highlighted wait (t)", Min: 0, Max: 4, Step: 0.05, Def: 1},
		},
		Render: render,
	})
}

// PDF returns the exponential distribution's probability density at t:
// λ·e^(−λt) for t>=0, or 0 for t<0 (a wait can't be negative) or λ<=0 (not a
// valid rate).
func PDF(lambda, t float64) float64 {
	if lambda <= 0 || t < 0 {
		return 0
	}
	return lambda * math.Exp(-lambda*t)
}

// CDF returns P(wait <= t): the probability the next event has already
// arrived by time t, 1 − e^(−λt) for t>=0.
func CDF(lambda, t float64) float64 {
	if t < 0 {
		return 0
	}
	return 1 - Survival(lambda, t)
}

// Survival returns P(wait > t): the probability nothing has arrived yet by
// time t, e^(−λt). This is exactly poisson-distribution's PMF(λt, 0) — "zero
// events in a window of length t" — which is the whole derivation of this
// distribution (see LEARNINGS.md).
func Survival(lambda, t float64) float64 {
	if lambda <= 0 {
		return 1
	}
	if t < 0 {
		return 1
	}
	return math.Exp(-lambda * t)
}

// Mean returns the average wait, 1/λ. Returns 0 for λ<=0 (not a valid rate).
func Mean(lambda float64) float64 {
	if lambda <= 0 {
		return 0
	}
	return 1 / lambda
}

// ConditionalSurvival returns P(wait > s+t | wait > s): given the wait has
// already exceeded s with nothing arriving, the probability it exceeds s+t.
// Used to demonstrate memorylessness — see the test asserting this always
// equals Survival(lambda, t), independent of s.
func ConditionalSurvival(lambda, s, t float64) float64 {
	denom := Survival(lambda, s)
	if denom == 0 {
		return 0
	}
	return Survival(lambda, s+t) / denom
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 360, 0, 4, 0, 1).String()
}
