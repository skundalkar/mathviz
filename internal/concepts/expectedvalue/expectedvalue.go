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
