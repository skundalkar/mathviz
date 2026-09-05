// Package expectedvalue visualizes expected value: the probability-weighted
// average of a random variable's outcomes. The running example is a $5
// carnival game — flip a weighted coin; heads pays a net $15, tails nets a
// $5 loss — and the question is whether that game is worth playing on
// average, not on any single play.
package expectedvalue

import (
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
				Body:    []string{"placeholder"},
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

func render(p map[string]float64) string {
	return viz.New(680, 400, 0, 1, 0, 1).String()
}
