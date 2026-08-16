// Package montecarlo visualizes Monte Carlo estimation: scatter random points
// across a square, check the one thing you *can* check cheaply for each point
// (does it land inside the inscribed circle?), and turn the fraction that do
// into an estimate of π -- a number with no way to "look up" the answer from
// a single random point, only from many of them combined. The running
// fraction is exactly a law-of-large-numbers average in disguise, and its
// error shrinks like 1/sqrt(n), the central-limit-theorem rate.
package montecarlo

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "monte-carlo-estimation",
		Seq:   39,
		Title: "Monte Carlo estimation (guessing π by throwing darts)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"A circle's area has a clean formula: πr². So does a square's, a triangle's, " +
						"almost any shape you drew in school. But plenty of real shapes and " +
						"questions don't have one: the outline of a lake on a map, the region of a " +
						"rocket's parameters where its landing stays safe, the odds a poker hand " +
						"wins against a random opponent's, an integral in 20 variables no textbook " +
						"formula covers. In every one of those cases, the one thing you *can* " +
						"cheaply do is check a single random point or scenario and get a yes/no " +
						"answer: 'did this land inside the lake's outline?', 'did this simulated " +
						"hand win?'. Is a big pile of cheap yes/no checks like that actually enough " +
						"to pin down a real number you couldn't compute directly?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take a case where you already know the true answer, so you can grade the " +
						"technique: a circle of radius 1 sitting inside a 2×2 square. The circle's " +
						"area is π; the square's is 4; so a point dropped uniformly at random into " +
						"the square lands inside the circle with probability π/4. Flip that around: " +
						"drop a lot of random points, count what fraction land inside the circle " +
						"(check x²+y² ≤ 1), and multiply that fraction by 4 to estimate π. No " +
						"formula for π was used anywhere in the estimate itself.",
					"• n=5 points: 4 land inside → estimate = 4×(4/5) = 3.20.",
					"• n=20 points: 16 inside → estimate = 4×(16/20) = 3.20 — same fraction, no " +
						"progress yet.",
					"• n=100 points: 79 inside → estimate = 3.16 — closer to the true 3.14159...",
					"• n=200 points: 167 inside → estimate = 3.34 — farther from π than n=100 was, " +
						"despite doubling the points!",
					"• n=1000 points: 809 inside → estimate = 3.24.",
					"• n=2000 points: 1605 inside → estimate = 3.21 — settling in, but still " +
						"noisy, not locked on.",
					"That n=200 stumble is the same lesson as the law of large numbers: this " +
						"running fraction is a law-of-large-numbers average (see " +
						"law-of-large-numbers) of a sequence of 0/1 checks, so more samples shrink " +
						"its *typical* wobble without promising every individual step lands closer " +
						"than the last one. What's new here is the rate: the central limit theorem " +
						"(see central-limit-theorem) says that wobble shrinks proportional to " +
						"1/√n, not 1/n — a fact with a very concrete consequence, covered below.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The square is the 2×2 sampling region; the dashed circle is the radius-1 " +
						"circle whose area you're estimating. Each small square is one random " +
						"sample point, colored green if it landed inside the circle and red if it " +
						"landed outside. The n slider sets how many points are drawn (more points, " +
						"denser scatter, and the green-to-red ratio settling closer to π/4 ≈ 78.5% " +
						"green). The run slider swaps in a completely different batch of random " +
						"points at the same n — flip it back and forth to see that the estimate " +
						"itself is random: two runs at the same n land on visibly different " +
						"estimates, not the same one.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Estimate a quantity you have no formula for, by turning it into 'what " +
						"fraction of random checks come back yes?' — the same recipe works whether " +
						"the checks are points in a square, simulated card hands, or randomized " +
						"trials of a physics model. You can also budget for it in advance: because " +
						"the typical error shrinks proportional to 1/√n, quadrupling your sample " +
						"count only halves your typical error, not eliminates it.",
					"|samples (n)|typical estimate error (∝ 1/√n)|",
					"|100|≈0.10|",
					"|400|≈0.05|",
					"|1,600|≈0.025|",
					"|6,400|≈0.0125|",
					"Squeezing out one more decimal digit of accuracy costs roughly 100× the " +
						"samples — worth knowing before you set a Monte Carlo simulation running " +
						"overnight expecting it to land within 0.001.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Weather services run dozens of slightly-perturbed simulations of the same " +
						"storm and report 'a 70% chance of rain' from the fraction that rained in " +
						"the simulation, not from a closed-form storm equation. Poker and blackjack " +
						"strategy calculators estimate a hand's win probability by dealing out " +
						"thousands of random simulated hands and counting wins, because computing " +
						"the exact probability by hand is impractical. Retirement calculators run " +
						"thousands of simulated versions of the stock market to estimate the odds " +
						"your savings last 30 years. Video game ray tracers estimate how much light " +
						"reaches a pixel by firing random light rays and averaging what bounces " +
						"back. The technique's name literally comes from a casino — it was coined " +
						"by physicists modeling nuclear chain reactions, who needed a code name for " +
						"work built on repeated random chance, the same idea as a roulette wheel.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: a Monte Carlo estimate is one random draw from a range of " +
						"answers you'd get from different batches of samples, and its typical error " +
						"shrinks like 1/√n — so to halve the error you need 4× the samples, not 2×. " +
						"Not like this: treating the number that comes out of one run as the exact, " +
						"fixed answer — this page's own run slider shows two different batches at " +
						"the same n landing on different estimates. Also not like this: assuming " +
						"more samples guarantees a strictly closer estimate at every step — the " +
						"n=100-to-n=200 stumble above (3.16 to 3.34) shows a doubled sample count " +
						"can still land farther from the truth than a smaller one did, exactly as " +
						"the law of large numbers predicts for any single run.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Number of random points (n)", Min: 10, Max: 2000, Step: 10, Def: 200},
			{Key: "run", Label: "Which random batch (run)", Min: 0, Max: 9, Step: 1, Def: 0},
		},
		Render: render,
	})
}

func render(params map[string]float64) string {
	_ = params
	return viz.New(680, 420, -1.2, 1.2, -1.2, 1.2).String()
}
