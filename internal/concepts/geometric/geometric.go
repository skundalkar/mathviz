// Package geometric visualizes the geometric distribution: the probability
// that the first success in a sequence of independent Bernoulli(p) trials
// lands exactly on trial k. binomial-distribution answers "out of a fixed n
// trials, how many succeed"; this concept answers the open-ended question
// "how many trials until the first success happens at all" -- flipping a
// coin until the first heads is the running example.
package geometric

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "geometric-distribution",
		Seq:   90,
		Title: "Geometric distribution (trials until first success)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"binomial-distribution answers 'out of a fixed n trials, how many " +
						"succeed' — but it has to fix n before it can say anything at all. Plenty " +
						"of real questions don't come with a fixed trial count: how many times do " +
						"you need to flip a coin before you see the first heads? How many sales " +
						"calls before the first yes? Someone new to this often reasons 'well, it's " +
						"a 50/50 coin, so it should take about 2 flips' and stops there, without " +
						"ever asking how likely it is to take exactly 1 flip, exactly 5, or more " +
						"than 10. Binomial's whole machinery — pick n, then ask about a count within " +
						"it — has nothing to say about 'how long until' questions where the number " +
						"of trials is exactly what's unknown. Is there a version of the same " +
						"trial-by-trial reasoning that instead answers 'what's the probability the " +
						"first success takes exactly k tries'?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Flip a fair coin (p=0.5) repeatedly until the first heads. For the first " +
						"heads to land exactly on flip k, every one of the first k-1 flips must be " +
						"tails, and flip k itself must be heads — independent events, so multiply:",
					"• First heads on flip 1: just need heads immediately — P = 0.5.",
					"• First heads on flip 2: tails, then heads — P = 0.5 × 0.5 = 0.25.",
					"• First heads on flip 3: tails, tails, then heads — P = 0.5³ = 0.125.",
					"• First heads on flip 4: three tails, then heads — P = 0.5⁴ = 0.0625.",
					"Each step is exactly half the one before it — hence 'geometric.' The general " +
						"formula for any success probability p: P(X=k) = (1-p)^(k-1) × p — the first " +
						"k-1 trials all fail (probability (1-p) each, multiplied k-1 times) and the " +
						"k-th succeeds (probability p). Summed over all k this series adds to exactly " +
						"1, and its probability-weighted average — the same Σk×P(X=k) definition " +
						"`expected-value` introduced — works out to a clean closed form: E[X] = 1/p. " +
						"A fair coin (p=0.5) takes 2 flips on average to see the first heads; a " +
						"1-in-6 event (rolling a specific number on a die, p=1/6) takes 6 rolls on " +
						"average.",
					"|success probability p|mean trials until first success (1/p)|",
					"|0.5 (fair coin)|2|",
					"|0.3|3.33|",
					"|0.1|10|",
					"|0.05|20|",
					"Rarer successes push the mean wait up fast — dropping p from 0.5 to 0.05 " +
						"doesn't just double the wait, it makes it 10x longer.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Each bar is P(X=k) for one trial count k, from 1 to 25. p tilts how fast the " +
						"bars shrink — a high p (close to 1) makes the first bar towering and the " +
						"rest nearly invisible, since success usually comes almost immediately; a low " +
						"p spreads real probability out across many more bars, and the tallest bar is " +
						"still trial 1 (it always is — see the common mistake below), just barely " +
						"taller than trial 2. The k slider highlights one bar in orange; the readout " +
						"gives its exact probability alongside P(X<=k) (first success by trial k) and " +
						"P(X>k) (it takes longer than k tries). The gray line marks the mean, 1/p.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Answer 'how many tries until the first success' questions directly from p, " +
						"with no simulation — for the exact count ('what's the chance it takes " +
						"exactly 5 tries'), for a range ('what's the chance it takes at most 5' or " +
						"'more than 10'), and for the long-run average wait (1/p) all at once. You " +
						"can also now precisely quantify a feeling most people already have vaguely " +
						"right — that rare events take a long time on average to show up at all — " +
						"and say by how much: at p=0.05 you'd expect a 20-trial wait on average, and " +
						"there's still a real (P(X>20)=36%) chance it takes even longer than that.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Number of sales calls before the first yes, number of job applications before " +
						"the first offer, number of coin flips or die rolls before a specific result " +
						"in a game, number of at-bats before a hitter's next hit, and quality control " +
						"('how many units come off the line, on average, before the first defect'). " +
						"Anywhere the question is 'how long do I have to keep trying before this " +
						"specific, independent, repeatable thing finally happens.'",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: trial 1 is always the single most likely count for the first " +
						"success, no matter how small p is — the bars only ever shrink from there, " +
						"they never rise to a peak partway through. It feels wrong ('shouldn't " +
						"success be 'due' after a few failures?') but it isn't: each trial is an " +
						"independent fresh 1-in-p shot, and P(X=1)=p is always at least as large as " +
						"P(X=k) for any later k.",
					"Not like this: assuming that after several failures, a success is 'due' soon " +
						"— the classic gambler's fallacy. The geometric distribution is memoryless: " +
						"given that you've already failed the first 5 trials, the probability the " +
						"next success takes m more trials is exactly PMF(p,m), identical to the " +
						"probability it would've taken m trials starting from scratch. Five failures " +
						"in a row don't shorten the wait for the sixth trial even slightly — they " +
						"only tell you five trials are already spent, not that the next one owes you " +
						"anything.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "p", Label: "Success probability (p)", Min: 0.05, Max: 0.95, Step: 0.05, Def: 0.3},
			{Key: "k", Label: "Highlighted trial (k)", Min: 1, Max: 25, Step: 1, Def: 3},
		},
		Render: render,
	})
}

// PMF returns P(X=k): the probability the first success lands exactly on
// trial k, for k=1,2,3,... The first k-1 trials must all fail (each with
// probability 1-p) and trial k must succeed (probability p): PMF(p,k) =
// (1-p)^(k-1) * p. Returns 0 for k<1 or p outside (0,1].
func PMF(p float64, k int) float64 {
	if k < 1 || p <= 0 || p > 1 {
		return 0
	}
	return math.Pow(1-p, float64(k-1)) * p
}

// CDF returns P(X<=k): the probability the first success happens by trial k
// (i.e. within the first k trials). Summing the geometric series
// Sum_{i=1..k} (1-p)^(i-1)*p has the closed form 1-(1-p)^k -- "it didn't
// take longer than k trials" is the complement of "every one of the first
// k trials failed".
func CDF(p float64, k int) float64 {
	if k < 1 {
		return 0
	}
	if p <= 0 || p > 1 {
		return 0
	}
	return 1 - math.Pow(1-p, float64(k))
}

// Mean returns the expected number of trials until the first success, 1/p.
// A fair coin (p=0.5) takes 2 flips on average to see the first heads.
func Mean(p float64) float64 {
	if p <= 0 {
		return math.Inf(1)
	}
	return 1 / p
}

// Variance returns the variance of the trial count until first success,
// (1-p)/p^2.
func Variance(p float64) float64 {
	if p <= 0 {
		return math.Inf(1)
	}
	return (1 - p) / (p * p)
}

// displayTrials is how many trials the bar chart shows -- fixed regardless
// of p so the picture's shape (a steadily shrinking tail) is comparable
// across slider settings, matching the k slider's own range.
const displayTrials = 25

func render(params map[string]float64) string {
	p := params["p"]
	if p < 0.01 {
		p = 0.01
	}
	if p > 1 {
		p = 1
	}
	k := int(params["k"] + 0.5)
	if k < 1 {
		k = 1
	}
	if k > displayTrials {
		k = displayTrials
	}

	pmf := make([]float64, displayTrials+1) // index 0 unused, trials are 1-based
	maxP := 0.0
	for i := 1; i <= displayTrials; i++ {
		pmf[i] = PMF(p, i)
		if pmf[i] > maxP {
			maxP = pmf[i]
		}
	}
	yMax := maxP * 1.2
	if yMax <= 0 {
		yMax = 1
	}

	c := viz.New(680, 420, 0.5, float64(displayTrials)+0.5, 0, yMax)
	c.Axes()
	for x := 1.0; x <= displayTrials; x += 5 {
		c.Tick(x, fmt.Sprintf("%.0f", x))
	}

	mean := Mean(p)
	if mean <= float64(displayTrials) {
		c.Path([][2]float64{{mean, 0}, {mean, yMax}}, viz.Muted, 1)
	}

	barHalfW := 0.4
	for i := 1; i <= displayTrials; i++ {
		color := viz.Accent
		if i == k {
			color = viz.Warm
		}
		x0, x1 := c.X(float64(i)-barHalfW), c.X(float64(i)+barHalfW)
		y0, y1 := c.Y(0), c.Y(pmf[i])
		c.Rect(x0, y1, x1-x0, y0-y1, color, 0.8)
	}

	c.Text(16, 24, fmt.Sprintf("p=%.2f    mean trials until first success=1/p=%.2f    variance=%.2f",
		p, mean, Variance(p)), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("P(X=%d)=%.4f    P(X<=%d)=%.4f    P(X>%d)=%.4f",
		k, PMF(p, k), k, CDF(p, k), k, 1-CDF(p, k)), 15, viz.Warm, "start")

	return c.String()
}
