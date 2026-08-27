// Package hypergeometric visualizes the hypergeometric distribution: the
// probability of drawing exactly k successes when sampling n items without
// replacement from a finite population of N items that contains K
// successes. binomial-distribution assumes every trial has the same fixed
// success probability p, which is only exactly true when you sample with
// replacement (or from an effectively infinite population); this concept
// shows what changes once the population is finite and each draw is
// removed from it.
package hypergeometric

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "hypergeometric-distribution",
		Seq:   73,
		Title: "Hypergeometric distribution (sampling without replacement)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"binomial-distribution assumes every trial has the exact same success " +
						"probability p, independent of every other trial — flip the same coin " +
						"again, draw the same card back out of the deck before the next draw. But " +
						"plenty of real sampling doesn't put anything back: draw 5 cards from a " +
						"52-card deck, pick 5 marbles from a bag without returning any, or survey " +
						"20 students from a class of 30 without picking the same student twice. " +
						"Once you keep the first success out of the pool, doesn't that change the " +
						"odds for every draw after it — and if binomial's fixed-p assumption is " +
						"technically wrong here, does it actually matter, or is it close enough?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"A bag holds 20 marbles, 8 red (a 'success') and 12 blue. Draw 5 without " +
						"putting any back, and ask for the probability of exactly k red marbles. " +
						"Unlike binomial, you can't just multiply the same p by itself: the number " +
						"of ways to choose k reds out of the 8 available (C(8,k)) times the number " +
						"of ways to fill the rest with blues (C(12, 5-k)), divided by the total " +
						"number of ways to draw any 5 of the 20 (C(20,5)) = 15,504.",
					"• P(X=0) = C(8,0)xC(12,5)/15504 = 792/15504 = 0.0511.",
					"• P(X=1) = C(8,1)xC(12,4)/15504 = 3960/15504 = 0.2554.",
					"• P(X=2) = C(8,2)xC(12,3)/15504 = 6160/15504 = 0.3973 -- the single most " +
						"likely count.",
					"• P(X=3)=0.2384, P(X=4)=0.0542, P(X=5)=0.0036 -- and all six add to 1.0000.",
					"Now compare against the naive binomial approximation that pretends each of " +
						"the 5 draws independently has probability p=8/20=0.4 of being red (as if " +
						"you put every marble back before the next draw): P(X=5)=0.0102, nearly 3x " +
						"the true 0.0036. Both distributions share the exact same mean (n x K/N = 5 " +
						"x 8/20 = 2.0), but the true variance is 0.9474 while binomial's naive " +
						"variance is 1.2000 -- hypergeometric is measurably *less* spread out, " +
						"because running low on red marbles as you draw makes an extreme run " +
						"(5-for-5 red) harder to sustain than independent draws would.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Solid bars are the true hypergeometric PMF for the current N (population), K " +
						"(successes in it), and n (draw size); the dashed line traces what binomial " +
						"would predict at the same average success rate p=K/N, drawn over the same " +
						"bars for direct comparison. The k slider highlights one bar and reads off " +
						"its exact probability. Watch the gap between the bars and the dashed line " +
						"shrink as you drag N up while keeping K/N and n fixed -- draw 5 from a " +
						"population of 20 vs. 5 from 2,000, and 'without replacement' barely changes " +
						"the odds once the population dwarfs the sample.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Get an exact answer for 'how many successes in a fixed-size sample drawn " +
						"without replacement' instead of an approximation that quietly assumes an " +
						"infinite population -- which matters most for small populations or large " +
						"sample fractions (auditing 5 of 20 invoices is a very different draw than " +
						"5 of 20,000), and matters less and less as the population grows relative to " +
						"the sample. You can also now name exactly which correction factor to apply " +
						"to binomial's variance formula -- the finite population correction, (N-n)/" +
						"(N-1) -- to get the true, smaller hypergeometric variance instead of " +
						"binomial's overestimate.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Card games: the probability of being dealt a specific number of a card type " +
						"(hearts, aces) is hypergeometric, not binomial, because the deck doesn't " +
						"refill between cards. Quality control audits a fixed sample of units from " +
						"a finished batch without putting inspected units back into the batch. " +
						"Ecologists estimate a population's total size with 'capture-recapture': tag " +
						"K animals, release them, then catch n more later and count how many are " +
						"already tagged -- exactly this distribution. Lottery and raffle odds (how " +
						"many of your 5 tickets match the 6 drawn numbers) are hypergeometric too.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: hypergeometric is what binomial becomes once you stop " +
						"putting each draw back -- same idea of counting successes in a fixed " +
						"number of draws, but the odds shift after every draw instead of staying " +
						"fixed at p. Not like this: assuming the difference from binomial is always " +
						"negligible -- it's small when the sample is a tiny fraction of a huge " +
						"population (survey 1,000 people out of a country of millions), but large " +
						"when the sample is a big chunk of a small population (this concept's own " +
						"5-of-20 example), which is exactly what the finite population correction " +
						"factor (N-n)/(N-1) measures: it's close to 1 (negligible correction) when n " +
						"is tiny relative to N, and shrinks well below 1 as n approaches N.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "N", Label: "Population size (N)", Min: 10, Max: 40, Step: 1, Def: 20},
			{Key: "K", Label: "Successes in population (K)", Min: 1, Max: 40, Step: 1, Def: 8},
			{Key: "n", Label: "Sample size drawn (n)", Min: 1, Max: 20, Step: 1, Def: 5},
			{Key: "k", Label: "Highlighted count (k)", Min: 0, Max: 20, Step: 1, Def: 2},
		},
		Render: render,
	})
}

// Choose returns n-choose-k, the number of ways to pick an unordered
// k-element subset out of n items. Computed multiplicatively (rather than
// via factorials, which overflow fast) and returns 0 for out-of-range k.
// Duplicated from binomial-distribution's own Choose rather than imported,
// matching this gallery's convention of keeping each concept self-contained.
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

// PMF returns the hypergeometric probability of drawing exactly k successes
// when drawing n items without replacement from a population of N items
// that contains K successes: the number of ways to pick k of the K
// successes and the rest from the N-K non-successes, divided by the total
// number of ways to draw any n of the N items.
func PMF(N, K, n, k int) float64 {
	if k < 0 || k > n || k > K || n-k > N-K {
		return 0
	}
	return Choose(K, k) * Choose(N-K, n-k) / Choose(N, n)
}

// CDF returns P(X<=k), the probability of k or fewer successes, by summing
// PMF over every count from 0 to k.
func CDF(N, K, n, k int) float64 {
	if k < 0 {
		return 0
	}
	sum := 0.0
	for i := 0; i <= k; i++ {
		sum += PMF(N, K, n, i)
	}
	return sum
}

// Mean returns the expected number of successes, n*K/N -- the same mean a
// binomial(n, p=K/N) distribution has.
func Mean(N, K, n int) float64 {
	return float64(n) * float64(K) / float64(N)
}

// Variance returns the hypergeometric variance: binomial's n*p*(1-p),
// scaled down by the finite population correction (N-n)/(N-1), which
// accounts for successes being removed from the pool as they're drawn.
// N<=1 has no well-defined correction and returns 0.
func Variance(N, K, n int) float64 {
	if N <= 1 {
		return 0
	}
	p := float64(K) / float64(N)
	fpc := float64(N-n) / float64(N-1)
	return float64(n) * p * (1 - p) * fpc
}

// BinomialPMF is the naive with-replacement approximation -- p held fixed
// at every one of the n draws -- kept here only for comparison against the
// true hypergeometric PMF above.
func BinomialPMF(n int, p float64, k int) float64 {
	if k < 0 || k > n {
		return 0
	}
	return Choose(n, k) * math.Pow(p, float64(k)) * math.Pow(1-p, float64(n-k))
}

func render(params map[string]float64) string {
	N := int(params["N"] + 0.5)
	if N < 1 {
		N = 1
	}
	K := int(params["K"] + 0.5)
	if K < 0 {
		K = 0
	}
	if K > N {
		K = N
	}
	n := int(params["n"] + 0.5)
	if n < 0 {
		n = 0
	}
	if n > N {
		n = N
	}
	k := int(params["k"] + 0.5)
	if k < 0 {
		k = 0
	}
	if k > n {
		k = n
	}

	pApprox := float64(K) / float64(N)

	// PMF for every count 0..n, both the true hypergeometric and the naive
	// binomial comparison, plus the tallest bar so the y-axis fits this
	// particular N/K/n instead of always reserving room for the worst case.
	hyperPMF := make([]float64, n+1)
	binomPMF := make([]float64, n+1)
	maxP := 0.0
	for i := 0; i <= n; i++ {
		hyperPMF[i] = PMF(N, K, n, i)
		binomPMF[i] = BinomialPMF(n, pApprox, i)
		if hyperPMF[i] > maxP {
			maxP = hyperPMF[i]
		}
		if binomPMF[i] > maxP {
			maxP = binomPMF[i]
		}
	}
	yMax := maxP * 1.2
	if yMax <= 0 {
		yMax = 1
	}

	c := viz.New(680, 420, -0.5, float64(n)+0.5, 0, yMax)
	c.Axes()
	tickStep := 1.0
	if n > 12 {
		tickStep = math.Ceil(float64(n) / 10)
	}
	for x := 0.0; x <= float64(n); x += tickStep {
		c.Tick(x, fmt.Sprintf("%.0f", x))
	}

	barHalfW := 0.42
	for i := 0; i <= n; i++ {
		color := viz.Accent
		if i == k {
			color = viz.Warm
		}
		x0, x1 := c.X(float64(i)-barHalfW), c.X(float64(i)+barHalfW)
		y0, y1 := c.Y(0), c.Y(hyperPMF[i])
		c.Rect(x0, y1, x1-x0, y0-y1, color, 0.8)
	}

	// The naive with-replacement (binomial) comparison, as a line over the
	// same counts, so the gap from the true hypergeometric bars is visible
	// directly rather than needing a second set of bars.
	binomLine := make([][2]float64, n+1)
	for i := 0; i <= n; i++ {
		binomLine[i] = [2]float64{float64(i), binomPMF[i]}
	}
	c.Path(binomLine, viz.Warm, 2)

	mean := Mean(N, K, n)
	variance := Variance(N, K, n)
	binomVar := float64(n) * pApprox * (1 - pApprox)
	fpc := 0.0
	if N > 1 {
		fpc = float64(N-n) / float64(N-1)
	}

	c.Text(16, 24, fmt.Sprintf("N=%d    K=%d    n=%d    p=K/N=%.2f", N, K, n, pApprox), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("P(X=%d) = %.4f    P(X<=%d) = %.4f", k, PMF(N, K, n, k), k, CDF(N, K, n, k)),
		15, viz.Warm, "start")
	c.Text(16, 64, fmt.Sprintf("mean = n*K/N = %.2f (same for both)    variance: hypergeometric = %.4f, binomial = %.4f (x%.4f correction)",
		mean, variance, binomVar, fpc), 12, viz.Muted, "start")
	c.Text(16, 400, "solid bars = true hypergeometric PMF    orange line = naive binomial(n, p=K/N) comparison",
		12, viz.Muted, "start")

	return c.String()
}
