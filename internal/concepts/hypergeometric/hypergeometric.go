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

func render(p map[string]float64) string {
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
