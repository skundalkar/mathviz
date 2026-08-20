// Package newton visualizes Newton's method: at the current guess, replace
// the curve by its tangent line and jump to where that line crosses zero.
// Repeating this a handful of times on f(x) = x^2 - 2 converges on √2 far
// faster than narrowing a bracket by trial and error.
package newton

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "newtons-method",
		Seq:   52,
		Title: "Newton's method (chasing a root with tangent lines)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Say your calculator has +, -, ×, ÷ but no √ button, and you need √2. Gut " +
						"instinct: guess a number, square it, see if you're too high or too low, " +
						"and narrow the guess — 1 is too low (1²=1), 1.5 is too high (1.5²=2.25), " +
						"1.4 is too low (1.4²=1.96), 1.42 is too high (1.42²=2.0164)... This works, " +
						"but each guess only tells you which direction to move, not how far — " +
						"getting a few more correct decimal digits means repeating the same " +
						"halving-the-gap process many more times. Is there a way to use more " +
						"information than just 'too high or too low' at each guess, so each step " +
						"jumps you dramatically closer instead of only halving the remaining gap?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Reframe '√2' as 'the root of f(x) = x² - 2' (the x where f(x)=0, since " +
						"x²-2=0 means x²=2). At any guess x, you already know two things: f(x) " +
						"itself (how far off zero you are) and the slope f'(x) (from `derivative`) " +
						"— which way and how steeply the curve is rising there. Draw the tangent " +
						"line at (x, f(x)) and see where THAT straight line crosses zero, and use " +
						"that crossing as the next guess: x_new = x - f(x)/f'(x). Starting from " +
						"x0=1:",
					"• x0=1: f(1)=1-2=-1, f'(1)=2·1=2. Next guess: 1 - (-1)/2 = 1.5.",
					"• x1=1.5: f(1.5)=2.25-2=0.25, f'(1.5)=2·1.5=3. Next guess: 1.5 - 0.25/3 ≈ " +
						"1.41667.",
					"• x2≈1.41667: f(x2)≈0.00694, f'(x2)≈2.83333. Next guess: 1.41667 - " +
						"0.00694/2.83333 ≈ 1.41422.",
					"• x3≈1.41422: already matches √2 ≈ 1.41421356 to 4 decimal places, after " +
						"just three steps — the guess-and-narrow approach above would still be " +
						"stuck around 2-3 correct digits.",
					"Look at how fast the error shrinks: from x0 it's off by about 0.41, at x1 " +
						"about 0.086, at x2 about 0.0024, at x3 about 0.0000019 — each step's " +
						"error is roughly the SQUARE of the previous one, not just a fraction of " +
						"it. That squaring is why so few steps get so many correct digits.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"x0 sets the starting guess; steps sets how many tangent-line jumps to draw. " +
						"The black curve is f(x) = x² - 2; the vertical dashed line marks the true " +
						"root, √2. Each orange line is the tangent at the current guess — follow " +
						"it down to where it crosses the x-axis, and that crossing becomes the dot " +
						"marking the next guess. Watch the dots visibly leap toward the dashed line " +
						"faster and faster with each step, exactly matching the shrinking-error " +
						"pattern from the worked example above.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Find a root of pretty much any smooth function to as many correct digits as " +
						"you want, in only a handful of steps, using nothing but the ability to " +
						"evaluate the function and its derivative at a point — no algebraic " +
						"formula for the root required, which matters because most equations " +
						"(anything beyond a quadratic, in general) don't have one.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Calculators and computer math libraries compute square roots, cube roots, " +
						"and reciprocals with Newton's method (or a close variant) instead of a " +
						"lookup table, because it converges to full precision in only a handful of " +
						"steps; video games and CAD software use it to find where a ray hits a " +
						"curved surface; and `gradient-descent`'s cousin (sometimes called Newton's " +
						"method in optimization) finds where a function's derivative is exactly " +
						"zero — its minimum — by running this same tangent-line trick one " +
						"derivative up.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'Newton's method jumps to where the TANGENT LINE crosses " +
						"zero, not where the curve itself crosses zero' — those only coincide once " +
						"it's already converged; every step in between is chasing a straight-line " +
						"approximation, and that's exactly why a step can occasionally overshoot. " +
						"Not like this: assuming it always converges from any starting guess — a " +
						"guess where f'(x) is 0 or very close to 0 (a flat tangent) sends the next " +
						"guess flying off to a huge or undefined value, and starting far from the " +
						"intended root can converge to a different root than the one you meant.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "x0", Label: "Starting guess x0", Min: 0.3, Max: 3, Step: 0.05, Def: 1},
			{Key: "steps", Label: "Iterations shown", Min: 0, Max: 6, Step: 1, Def: 3},
		},
		Render: render,
	})
}

// F is the function this concept finds a root of: f(x) = x^2 - 2, whose
// positive root is √2.
func F(x float64) float64 {
	return x*x - 2
}

// FPrime is the exact derivative of F: f'(x) = 2x.
func FPrime(x float64) float64 {
	return 2 * x
}

// NewtonStep returns the next guess: where the tangent line at (x, F(x))
// crosses zero. Returns x unchanged if the tangent is flat (f'(x) == 0),
// since the line never crosses zero in that case.
func NewtonStep(x float64) float64 {
	slope := FPrime(x)
	if slope == 0 {
		return x
	}
	return x - F(x)/slope
}

// Iterates returns x0 followed by `steps` successive Newton steps, so the
// result always has length steps+1.
func Iterates(x0 float64, steps int) []float64 {
	if steps < 0 {
		steps = 0
	}
	out := make([]float64, steps+1)
	out[0] = x0
	for i := 1; i <= steps; i++ {
		out[i] = NewtonStep(out[i-1])
	}
	return out
}

func render(p map[string]float64) string {
	c := viz.New(560, 400, -0.5, 3, -2, 5)
	c.Axes()
	return c.String()
}
