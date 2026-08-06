// Package logscale answers "what does log actually mean?" in one picture:
// log_b(x) is the exponent you raise b to in order to get x. The curve shows
// that every time x MULTIPLIES by b, the log only ADDS 1 — which is why logs
// turn multiplication into addition and why log scales tame huge ranges.
package logscale

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "logarithms",
		Title: "What a logarithm means",
		Blurb: "log base b of x asks: 'b to what power gives me x?' Slide the base. " +
			"Notice the marked points: each time x multiplies by the base (b, then b², " +
			"then b³ …), the log value only steps up by 1. Multiplying the input becomes " +
			"adding to the output — that is the entire trick behind slide rules, decibels, " +
			"pH, earthquake magnitudes, and log-scale charts.",
		Params: []concept.ParamSpec{
			{Key: "base", Label: "Base (b)", Min: 2, Max: 10, Step: 1, Def: 2},
			{Key: "mark", Label: "Highlight x", Min: 1, Max: 64, Step: 1, Def: 8},
		},
		Render: render,
	})
}

// LogBase returns log_b(x). Pure math, unit-tested.
func LogBase(x, b float64) float64 {
	return math.Log(x) / math.Log(b)
}

func render(p map[string]float64) string {
	b := p["base"]
	if b < 2 {
		b = 2
	}
	mark := p["mark"]
	if mark < 1 {
		mark = 1
	}

	const xmin, xmax = 1.0, 64.0
	ymax := LogBase(xmax, b) + 0.5
	c := viz.New(680, 320, xmin, xmax, 0, ymax)
	c.Axes()

	// Curve y = log_b(x).
	curve := viz.Sample(xmin, xmax, 320, func(x float64) float64 {
		if x <= 0 {
			return 0
		}
		return LogBase(x, b)
	})
	c.Path(curve, viz.Accent, 2.5)

	// Mark integer powers of the base: b^0, b^1, b^2, ... showing +1 steps.
	for k := 0; ; k++ {
		x := math.Pow(b, float64(k))
		if x > xmax {
			break
		}
		y := float64(k)
		px, py := c.X(x), c.Y(y)
		c.Rect(px-3, py-3, 6, 6, viz.Warm, 1) // dot
		c.Text(px+6, py-6, fmt.Sprintf("b^%d = %g → %d", k, x, k), 11, viz.Muted, "start")
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	// The user-highlighted x and its (possibly fractional) log.
	yv := LogBase(mark, b)
	c.VLine(mark, viz.Good, true)
	c.Rect(c.X(mark)-3.5, c.Y(yv)-3.5, 7, 7, viz.Good, 1)

	c.Text(20, 24, fmt.Sprintf("base b = %g", b), 14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("log_%g(%g) = %.3f", b, mark, yv), 13, viz.Good, "start")
	c.Text(20, 62, "each ×b in x  →  +1 in log", 12, viz.Muted, "start")
	return c.String()
}
