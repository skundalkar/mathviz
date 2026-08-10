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
		Seq:   2,
		Title: "What a logarithm means",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Fold a 0.1mm-thick sheet of paper in half, over and over — each fold doubles " +
						"whatever thickness you already have (0.1 → 0.2 → 0.4 → 0.8mm …). How many " +
						"folds until the stack passes a 100m building? Your instinct might be to " +
						"divide: 100,000mm ÷ 0.1mm = 1,000,000 — but that answers a different " +
						"question, stacking a million loose sheets where each one just adds a fixed " +
						"0.1mm.",
					"Folding multiplies instead of adding, so the fold count is stuck up in an " +
						"exponent: 0.1×2ⁿ = 100,000. Plain division has no move for pulling n out of " +
						"an exponent.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"A logarithm is exactly the tool that reaches into an exponent and pulls it " +
						"back down to an ordinary number: log₂(1,000,000) ≈ 20 — just 20 folds, not " +
						"a million sheets. More generally, log_b(x) is the exponent you raise b to " +
						"in order to get x, and you can see it directly by marking the powers of a " +
						"base like b=2 and watching what happens to x versus the log:",
					"• 2^0 = 1 → log = 0",
					"• 2^1 = 2 → log = 1",
					"• 2^2 = 4 → log = 2",
					"• 2^3 = 8 → log = 3",
					"Each time x multiplies by the base, the log only steps up by exactly +1 — " +
						"multiplying the input becomes adding to the output.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The curve is y = log_b(x); sliding the base changes b and re-spaces the " +
						"marked powers along it, while the 'highlight x' slider drags a marker to " +
						"any x and reads off its (possibly fractional) log — with the defaults " +
						"(base 2, highlight x=8) the marker sits right on log₂(8) = 3, the same " +
						"2^3 = 8 point marked above. The caption 'each ×b in x → +1 in log' is the " +
						"paper-folding trick, shown for any base, not just doublings.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Any 'how many multiplications/doublings does it take to reach X' question " +
						"now has a direct answer instead of getting stuck trying to divide your way " +
						"to an exponent — the same move behind slide rules, and why decibels, pH, " +
						"star magnitudes, and log-scale charts can compress an enormous " +
						"multiplicative range onto a scale you can actually read.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"'This investment compounds annually' or 'this outbreak doubles every few " +
						"days' are folding stories, not stacking stories — reach for a logarithm, " +
						"not division, to get from a growth rate to 'how long until X.' It's also " +
						"why a Richter magnitude of 7 isn't 'a bit more' than 6 — it's about 32x " +
						"more energy, because each whole step up the scale is another " +
						"multiplication, not another addition.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"'This is growing exponentially' is a genuine folding story — each step " +
						"multiplies what came before, not adds a fixed amount — and it's exactly " +
						"when a logarithm, not a raw number, gives you a sane scale to reason about " +
						"it (compound interest, an outbreak's doubling time).",
					"'Exponential' gets used constantly to just mean 'a lot' or 'suddenly' — most " +
						"things people call exponential are actually growing at a steady additive " +
						"rate that merely felt surprising; real exponential growth compounds, and " +
						"compounding is what actually gets out of hand.",
				},
			},
		},
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
