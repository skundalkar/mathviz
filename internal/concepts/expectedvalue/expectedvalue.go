// Package expectedvalue visualizes expected value: the probability-weighted
// average of a random variable's outcomes. The running example is a $5
// carnival game — flip a weighted coin; heads pays a net $15, tails nets a
// $5 loss — and the question is whether that game is worth playing on
// average, not on any single play.
package expectedvalue

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "expected-value",
		Seq:   89,
		Title: "Expected value (the number a bet centers on)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"bayes-theorem and binomial-distribution both showed how to find the " +
						"probability of one particular outcome. But knowing a probability doesn't " +
						"by itself tell you whether a bet is worth taking. A carnival booth charges " +
						"$5 to play: flip a weighted coin that lands heads 20% of the time; heads " +
						"pays $20 (a net $15 win after the $5 you paid), tails pays nothing (a net " +
						"$5 loss). Most people's gut reaction latches onto one branch and stops " +
						"there — either 'it pays $20, that's a big win' or '80% of the time I lose, " +
						"skip it.' Is there a single number that honestly accounts for *both* " +
						"branches at once, weighted by how often each one actually happens, so this " +
						"game (and any other bet) can be judged on equal footing?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Work the carnival game by hand. There are two outcomes: win $15 net with " +
						"probability 0.20, or lose $5 net with probability 0.80. Weight each outcome " +
						"by how often it happens, and add:",
					"• P(win)=0.20, payoff=+$15 → contributes 0.20×15 = $3.00 to the average.",
					"• P(lose)=0.80, payoff=−$5 → contributes 0.80×(−5) = −$4.00 to the average.",
					"• Total: $3.00 + (−$4.00) = −$1.00 — the expected value, E[X].",
					"So on average you lose $1 per play, even though the $15 win looks tempting on " +
						"any single flip. That's the general recipe for any random variable with " +
						"outcomes x1..xn and probabilities p1..pn (which must sum to 1): " +
						"E[X] = Σ xi·pi — a probability-weighted average, not a plain average across " +
						"the branches (a plain average of $15 and −$5 would give $5, a wildly " +
						"different and wrong answer). The weighting is everything: flip the win " +
						"probability to 50% and the same game becomes E[X] = 0.5×15 + 0.5×(−5) = " +
						"7.5 − 2.5 = $5, a good bet instead of a bad one, with nothing about the " +
						"payoffs themselves having changed.",
					"|win probability p|E[X]|verdict|",
					"|10%|0.10×15 + 0.90×(−5) = −3.00|unfavorable|",
					"|20% (default)|0.20×15 + 0.80×(−5) = −1.00|unfavorable|",
					"|25%|0.25×15 + 0.75×(−5) = 0.00|breakeven|",
					"|50%|0.50×15 + 0.50×(−5) = 5.00|favorable|",
					"The breakeven row isn't a coincidence: solving p×15 + (1−p)×(−5) = 0 for p " +
						"gives p = 5/20 = 0.25 — below a 25% win chance this game loses money on " +
						"average at these payoffs; above it, it's profitable on average.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Two bars sit along a number line of net dollar outcomes: a green bar at the " +
						"win amount, height equal to the win probability p, and a red bar at the " +
						"lose amount, height equal to 1−p. p, win amount, and lose amount are all " +
						"sliders. The dashed orange line marks E[X] — the probability-weighted " +
						"balance point between the two bars, the same way a see-saw balances around " +
						"a point that depends on both how heavy each side is (the probability) and " +
						"how far out it sits (the payoff). Push p up and the line slides toward the " +
						"win bar; push it down and the line slides toward the lose bar — at the " +
						"default 20% it sits closer to the loss, which is exactly why the verdict " +
						"below reads unfavorable.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Compare two different gambles, deals, or business decisions on equal footing " +
						"with a single number, instead of an intuition that only looks at one " +
						"branch. This is also the bridge to two concepts already in this gallery: " +
						"variance-vs-stddev measured spread around a sample's mean — expected value " +
						"is the theoretical center that spread is measured around in the first " +
						"place, Var(X) = E[(X−E[X])²]. And mean-median-mode's 'mean' is the average " +
						"of data you've already collected; expected value is the average you'd " +
						"converge to over infinitely many repetitions of a random process you " +
						"haven't run yet — law-of-large-numbers is precisely the statement that the " +
						"first approaches the second as your sample grows.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Insurance companies set premiums above the expected payout so their expected " +
						"value across many policies is positive, even though any single policyholder " +
						"might file a large claim. State lotteries are built so a ticket's expected " +
						"value is always negative — that shortfall is exactly how the lottery funds " +
						"its payouts and overhead. Casino games are engineered the same way: the " +
						"house's expected value per bet is always slightly positive, which is why a " +
						"casino profits reliably even though any individual gambler might win big on " +
						"any individual night. The same arithmetic applies to everyday choices: " +
						"deciding whether to buy an extended warranty, whether a coin-flip signing " +
						"bonus is worth taking over a smaller guaranteed one, or whether a marketing " +
						"campaign with an uncertain payoff is worth its cost.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: expected value is a long-run average across many " +
						"repetitions, not a prediction about what happens on any single play — " +
						"you'll never actually observe a −$1.00 outcome playing the carnival game " +
						"once; you'll observe either +$15 or −$5, and −$1 is only where those two " +
						"outcomes average out over many plays.",
					"Not like this: assuming a positive expected value guarantees a win this time, " +
						"or that the expected value is close to what will 'usually' happen. " +
						"binomial-distribution already made this exact point about its mean n·p — " +
						"the single most likely count still isn't a promise about any one trial — " +
						"and the same caution applies here, now measured in dollars instead of " +
						"counts.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "p", Label: "Win probability (p)", Min: 0.01, Max: 0.99, Step: 0.01, Def: 0.2},
			{Key: "winAmt", Label: "Net win amount", Min: 0, Max: 50, Step: 1, Def: 15},
			{Key: "loseAmt", Label: "Net lose amount", Min: -50, Max: 0, Step: 1, Def: -5},
		},
		Render: render,
	})
}

// TwoOutcome returns the expected value of a random variable with exactly
// two possible outcomes: win with probability p, lose with probability
// 1-p. E[X] = p*win + (1-p)*lose -- the probability-weighted average of the
// two outcomes, not a plain average of win and lose.
func TwoOutcome(win, lose, p float64) float64 {
	return p*win + (1-p)*lose
}

// Discrete returns the expected value of a random variable with an
// arbitrary number of outcomes: E[X] = Sum(values[i] * probs[i]). This is
// the general form TwoOutcome is a two-outcome special case of. Mismatched
// slice lengths or a nil/empty input return 0.
func Discrete(values, probs []float64) float64 {
	sum := 0.0
	n := len(values)
	if len(probs) < n {
		n = len(probs)
	}
	for i := 0; i < n; i++ {
		sum += values[i] * probs[i]
	}
	return sum
}

// Breakeven returns the win probability p at which TwoOutcome(win, lose, p)
// is exactly 0 -- the threshold above which the bet is favorable on average
// and below which it isn't. Solving p*win + (1-p)*lose = 0 for p gives
// p = -lose / (win - lose). Returns NaN if win == lose (every p gives the
// same expected value, so no single breakeven point exists).
func Breakeven(win, lose float64) float64 {
	if win == lose {
		return math.NaN()
	}
	return -lose / (win - lose)
}

func render(params map[string]float64) string {
	p := params["p"]
	if p < 0.01 {
		p = 0.01
	}
	if p > 0.99 {
		p = 0.99
	}
	win := params["winAmt"]
	lose := params["loseAmt"]

	ev := TwoOutcome(win, lose, p)

	xMin, xMax := lose-5, win+5
	if xMin > -5 {
		xMin = -5
	}
	if xMax < 5 {
		xMax = 5
	}

	c := viz.New(680, 400, xMin, xMax, 0, 1.08)
	c.Axes()
	c.Tick(lose, fmt.Sprintf("%.0f", lose))
	c.Tick(0, "0")
	c.Tick(win, fmt.Sprintf("%.0f", win))

	// Zero reference line, so it's easy to read which bar sits on which
	// side of "break even on this single play".
	c.VLine(0, viz.Muted, false)

	const barHalfPx = 34.0
	drawBar := func(x, height float64, color string) {
		x0, x1 := c.X(x)-barHalfPx, c.X(x)+barHalfPx
		y0, y1 := c.Y(0), c.Y(height)
		c.Rect(x0, y1, x1-x0, y0-y1, color, 0.85)
		c.Text(c.X(x), y1-8, fmt.Sprintf("%.0f%%", height*100), 13, viz.Ink, "middle")
	}
	drawBar(win, p, viz.Good)
	drawBar(lose, 1-p, viz.Bad)

	// E[X] is the probability-weighted "balance point" between the two
	// bars -- the dashed line showing where it lands.
	c.VLine(ev, viz.Warm, true)
	c.Text(c.X(ev), 24, fmt.Sprintf("E[X] = %.2f", ev), 14, viz.Warm, "middle")

	c.Text(16, 340, fmt.Sprintf("win $%.0f w.p. %.0f%%, lose $%.0f w.p. %.0f%%", win, p*100, math.Abs(lose), (1-p)*100),
		13, viz.Ink, "start")
	c.Text(16, 360, fmt.Sprintf("E[X] = %.2f×%.0f + %.2f×%.0f = %.2f", p, win, 1-p, lose, ev),
		13, viz.Muted, "start")

	verdict := "a break-even bet"
	verdictColor := viz.Muted
	switch {
	case ev > 0.01:
		verdict = "favorable on average -- play it (on average, over many plays)"
		verdictColor = viz.Good
	case ev < -0.01:
		verdict = "unfavorable on average -- skip it (on average, over many plays)"
		verdictColor = viz.Bad
	}
	c.Text(16, 384, verdict, 13, verdictColor, "start")

	return c.String()
}
