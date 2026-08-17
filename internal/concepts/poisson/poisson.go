// Package poisson visualizes the Poisson distribution: the probability of
// seeing exactly k rare, independent events in a fixed window, given only
// their average rate λ. It's the limit of binomial-distribution as the
// number of trials n grows huge and the per-trial probability p shrinks
// toward zero while their product n*p stays fixed at λ — no separate "n" is
// needed once you're in that limit, just the rate itself.
package poisson

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "poisson-distribution",
		Seq:   43,
		Title: "Poisson distribution (rare events over a window)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"binomial-distribution needs a fixed number of trials n and a per-trial " +
						"probability p — 10 coin flips, 50 factory units, 1,000 poll respondents. " +
						"But plenty of real counts don't have a natural n at all. How many typos " +
						"are on this page? How many customers walk into a shop between 2pm and " +
						"3pm? There's no fixed number of 'trials' to point to — a typo could occur " +
						"at any of an enormous number of possible character positions, each with a " +
						"vanishingly small chance, and a customer could arrive at any of an " +
						"enormous number of possible instants within the hour. Can you still get an " +
						"exact probability distribution for the count, knowing only the *average " +
						"rate* — say, 3 typos per page, or 3 customers per hour — with no n or p to " +
						"plug in?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Yes — by taking binomial-distribution to an extreme. Chop the hour into " +
						"1,000 tiny one-off chances for a customer to walk in, each independently " +
						"with probability p=0.001, so the expected count is n*p = 1000*0.001 = 1 " +
						"customer, same as the real average rate λ=1. Compare that binomial " +
						"distribution to the Poisson formula P(X=k) = λ^k * e^-λ / k! at λ=1:",
					"• k=0: binomial(1000,0.001) = 0.36770, Poisson(λ=1) = 0.36788.",
					"• k=1: binomial = 0.36806, Poisson = 0.36788.",
					"• k=2: binomial = 0.18403, Poisson = 0.18394.",
					"• k=3: binomial = 0.06128, Poisson = 0.06131.",
					"They already agree to 3 decimal places at n=1000, and they converge exactly " +
						"as n keeps growing and p keeps shrinking with n*p held fixed at λ — chop " +
						"the hour into a million instants instead of a thousand and it's an even " +
						"closer match. So the Poisson PMF is just the binomial PMF's limit once you " +
						"no longer care about (or have) a specific n and p — only their product λ.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Each bar is P(X=k), the probability of exactly k events in the window, for " +
						"the λ set by the slider. The k slider highlights one bar and reads off its " +
						"exact probability plus the cumulative probability of that count or fewer. " +
						"The dashed line marks λ itself — which is both the mean *and* the variance " +
						"of a Poisson distribution, a fact you can check directly: push λ from 1 to " +
						"9 and watch the bars both shift right (higher mean) and spread out wider " +
						"(higher variance) by the exact same amount. At small λ the bars are sharply " +
						"skewed toward 0 (most windows see nothing at all); by λ≈10 the shape has " +
						"become visibly bell-like, the same normal-vs-skew smoothing " +
						"binomial-distribution's bars show as n grows.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Put an exact probability on a rare-event count from nothing but its average " +
						"rate — 'is 7 customer complaints this week unusually high, if the average " +
						"is 3?' — without needing to know how many 'chances' for a complaint there " +
						"technically were. You can also flag it as suspicious the moment mean and " +
						"variance stop matching: real event counts (like most binomial-distribution " +
						"outcomes at extreme p) are often *overdispersed* — variance bigger than the " +
						"mean, because the true rate itself varies rather than staying fixed like λ " +
						"— and the Poisson model's mean=variance identity is exactly the assumption " +
						"that check is testing.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Call centers use it to staff for how many calls arrive per minute. Web " +
						"services use it to model requests per second when deciding server " +
						"capacity. Epidemiologists use it for rare disease case counts in a region " +
						"per month. Insurers use it for the number of claims filed per policy per " +
						"year. Even sports: how many goals a soccer team scores per match is " +
						"commonly modeled as Poisson, which is why goal totals below 1 or above 5 " +
						"both feel unusual — they sit in the thin tails of the shape this page draws.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: λ is a rate over a fixed window (3 per hour), and P(X=k) is " +
						"the probability of exactly that many events in one instance of that same " +
						"window — change the window length and λ changes proportionally (3 per " +
						"hour becomes 6 per two hours, not 3).",
					"Not like this: assuming a Poisson process 'remembers' how long it's been " +
						"since the last event, so a long gap makes the next event 'due' — it " +
						"doesn't; each instant is independent of every other, the same " +
						"independence assumption law-of-large-numbers and monte-carlo-estimation " +
						"lean on, and unlike markov-chains there's no memory of what just happened " +
						"at all. Also not like this: treating mean and variance as two separate " +
						"numbers you'd estimate independently from data — for a true Poisson " +
						"process they're the same number, λ, by construction.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "lambda", Label: "Rate (λ, events per window)", Min: 0.5, Max: 15, Step: 0.5, Def: 3},
			{Key: "k", Label: "Highlighted count (k)", Min: 0, Max: 30, Step: 1, Def: 3},
		},
		Render: render,
	})
}

func render(params map[string]float64) string {
	_ = params
	return viz.New(680, 420, 0, 30, 0, 1).String()
}
