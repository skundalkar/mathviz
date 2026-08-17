// Package binomial visualizes the binomial distribution: the probability of
// getting exactly k successes in n independent trials, each with success
// probability p. It builds directly on pascals-triangle's n-choose-k — the
// binomial PMF is just that count, weighted by how likely each specific
// arrangement of successes and failures is.
package binomial

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "binomial-distribution",
		Seq:   42,
		Title: "Binomial distribution (counting successes)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You already know how to compute the probability of one specific outcome: " +
						"flip a fair coin 10 times and get heads-heads-tails-heads-...-tails, and " +
						"that exact sequence has probability 0.5^10. But that's rarely the question " +
						"anyone actually asks. The real question is almost always about a *count*, " +
						"not a *sequence*: out of 10 flips, what's the probability of exactly 5 " +
						"heads — in any order? A gut instinct says 'about half the time, since " +
						"p=0.5' — is that instinct right, or is 'exactly half' actually rarer than " +
						"it sounds?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Ten fair coin flips, p = 0.5. There are 2^10 = 1024 equally likely sequences " +
						"total, and pascals-triangle already tells you how many of them land exactly " +
						"5 heads: row 10, entry 5, which is 10-choose-5 = 252. So P(exactly 5 heads) " +
						"= 252 x 0.5^5 x 0.5^5 = 252/1024 = 0.2461 — about 24.6%, not 50%.",
					"• P(0 heads) = 1/1024 = 0.0010. P(1) = 10/1024 = 0.0098. P(2) = 45/1024 = 0.0439.",
					"• P(3) = 120/1024 = 0.1172. P(4) = 210/1024 = 0.2051. P(5) = 252/1024 = 0.2461.",
					"• P(6) = 210/1024 = 0.2051. P(7) = 120/1024 = 0.1172. P(8..10) mirror P(2..0).",
					"That's the general recipe, for any n and p: P(X=k) = C(n,k) x p^k x (1-p)^(n-k) " +
						"— C(n,k) counts how many of the n trials could be the k successes, and " +
						"p^k(1-p)^(n-k) is the probability of any one of those specific arrangements. " +
						"5 heads out of 10 is the single most likely count, but it still only owns " +
						"about a quarter of the probability — the other three quarters are spread " +
						"across every other count from 0 to 10.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Each bar is P(X=k) for one count k, from 0 to n. The n slider changes how many " +
						"trials there are, p tilts the distribution toward more or fewer successes " +
						"(p=0.5 keeps it symmetric; push p toward 0 or 1 and it skews hard toward " +
						"one edge), and the k slider highlights one specific bar in orange with its " +
						"exact probability, plus the cumulative probability of that count or fewer. " +
						"The dashed line marks the mean, n x p. Push n up while keeping p fixed and " +
						"watch the jagged bars smooth into the familiar bell shape — that's the " +
						"central-limit-theorem at work on a discrete count instead of a sample mean.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Answer 'exactly k' and 'k or fewer' (or 'k or more') questions about a count of " +
						"successes directly from n and p, with no simulation — the same way " +
						"markov-chains solved for a steady state instead of running a random walk. " +
						"You can also sanity-check gut instincts about counts: 'exactly half' out of " +
						"n trials is the single most likely outcome at p=0.5, but as n grows its own " +
						"share of the probability actually shrinks (P(exactly 5 of 10)=24.6%, but " +
						"P(exactly 50 of 100)≈8.0%), even as the distribution as a whole piles up " +
						"more tightly around the mean.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A basketball player who makes 70% of free throws: how many of their next 10 " +
						"attempts will they make? A factory with a known 2% defect rate: how many " +
						"defective units turn up in a batch of 50? An A/B test with a fixed number " +
						"of visitors and a known baseline conversion rate: how many conversions is " +
						"'normal' vs. surprisingly high? A political poll of 1,000 fixed respondents " +
						"with a known true yes-rate: how many say yes? All four are the same shape — " +
						"a fixed number of independent yes/no trials with one success probability.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: the probability of exactly k successes weights the count of " +
						"arrangements (C(n,k)) by the probability of any one arrangement " +
						"(p^k(1-p)^(n-k)) — both factors matter, and the mean n x p being the most " +
						"likely single count doesn't make it a likely outcome in absolute terms.",
					"Not like this: assuming p=0.5 makes the distribution symmetric no matter what " +
						"— it's only symmetric exactly at p=0.5; any other p skews it, same as it " +
						"would in normal-vs-skew. Also not like this: treating 'expected value' as " +
						"'the thing that will basically always happen' — n x p is an average over " +
						"many repetitions, not a promise about any single one, exactly the same " +
						"distinction law-of-large-numbers draws between a long-run average and a " +
						"single noisy trial.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Trials (n)", Min: 1, Max: 50, Step: 1, Def: 20},
			{Key: "p", Label: "Success probability (p)", Min: 0.05, Max: 0.95, Step: 0.05, Def: 0.5},
			{Key: "k", Label: "Highlighted count (k)", Min: 0, Max: 50, Step: 1, Def: 10},
		},
		Render: render,
	})
}

// Choose returns n-choose-k, the number of ways to pick an unordered k-element
// subset out of n items — the same count pascals-triangle builds row by row.
// It's computed multiplicatively (rather than via factorials, which overflow
// fast) and returns 0 for out-of-range k. Exact for the n<=50 this package's
// Params allow: float64 carries integers exactly up to 2^53, far past the
// largest value this range produces (n=50: C(50,25)~1.26e14).
func Choose(n, k int) float64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k // C(n,k) == C(n,n-k); fewer multiplications either way
	}
	result := 1.0
	for i := 0; i < k; i++ {
		result *= float64(n-i) / float64(i+1)
	}
	return result
}

// PMF returns P(X=k), the probability of exactly k successes in n independent
// trials each with success probability p: C(n,k) ways to place the successes,
// times the probability of any one specific placement, p^k(1-p)^(n-k).
func PMF(n int, p float64, k int) float64 {
	if k < 0 || k > n {
		return 0
	}
	return Choose(n, k) * math.Pow(p, float64(k)) * math.Pow(1-p, float64(n-k))
}

// CDF returns P(X<=k), the probability of k or fewer successes in n trials,
// by summing PMF over every count from 0 to k.
func CDF(n int, p float64, k int) float64 {
	if k < 0 {
		return 0
	}
	if k > n {
		k = n
	}
	sum := 0.0
	for i := 0; i <= k; i++ {
		sum += PMF(n, p, i)
	}
	return sum
}

// Mean returns the expected count of successes, n*p.
func Mean(n int, p float64) float64 {
	return float64(n) * p
}

// Variance returns the variance of the count of successes, n*p*(1-p).
func Variance(n int, p float64) float64 {
	return float64(n) * p * (1 - p)
}

func render(params map[string]float64) string {
	_ = params
	return viz.New(680, 420, 0, 50, 0, 1).String()
}
