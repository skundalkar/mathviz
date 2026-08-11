// Package derivative visualizes the derivative as the limit of a secant
// slope: pick a point x0 on f(x) = x^2, draw the line through (x0, f(x0))
// and a second point h away, then shrink h. The secant swings until it hugs
// the tangent line — that limit, not the secant at any one fixed h, is the
// derivative.
package derivative

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "derivative",
		Seq:   22,
		Title: "Derivative (slope of the tangent)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"A car's position is p(t) = t² meters at time t seconds (it's accelerating, " +
						"not moving at a constant speed). From t=1s to t=2s it covers 4-1 = 3 meters " +
						"in 1 second, so its average speed over that second is 3 m/s. That's a fine " +
						"answer to 'how fast on average between two moments' — but a speedometer " +
						"doesn't average anything; it reports one number for right now. What is the " +
						"car's speed at the single exact instant t=1s, when there's no 'before' and " +
						"'after' left to compare?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Shrink the gap between the two moments being compared and watch the average " +
						"speed change:",
					"• From t=1 to t=2 (gap h=1): p(2)=4, p(1)=1, average speed = (4-1)/1 = 3 m/s.",
					"• From t=1 to t=1.5 (h=0.5): p(1.5)=2.25, average speed = (2.25-1)/0.5 = 2.5 m/s.",
					"• From t=1 to t=1.1 (h=0.1): p(1.1)=1.21, average speed = (1.21-1)/0.1 = 2.1 m/s.",
					"• From t=1 to t=1.01 (h=0.01): p(1.01)=1.0201, average speed = (1.0201-1)/0.01 = 2.01 m/s.",
					"The gap keeps shrinking and the average speed keeps landing closer to 2 m/s — " +
						"that limit, not any one of the numbers above, is the speed at the exact " +
						"instant t=1s. Algebraically the same pattern holds for every x0: " +
						"[(x0+h)² - x0²]/h = (2·x0·h + h²)/h = 2x0 + h, which → 2x0 as h shrinks to 0 " +
						"— matching the 2.01, 2.001, ... trend above exactly. That limit is the " +
						"derivative, written f'(x0); for f(x)=x² it's f'(x0) = 2x0, so f'(1) = 2.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The blue secant line runs through (x0, f(x0)) and (x0+h, f(x0+h)) — the same " +
						"two points used to compute an average speed above. Drag h down and the " +
						"secant visibly rotates to hug the dashed tangent line, which always has " +
						"slope 2·x0. Drag x0 and both lines slide to a new point on the curve; the " +
						"tangent's slope changes with it, exactly tracking 2x0.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can now name a rate of change at one exact point — 'the speed right now,' " +
						"'the slope right here' — instead of only ever being able to describe an " +
						"average over some stretch before or after it.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A speedometer reads instantaneous speed, not the average speed since the trip " +
						"began; economists distinguish marginal cost (the cost of the very next unit) " +
						"from average cost so far; and the steepness underfoot on a hillside at the " +
						"exact spot you're standing can be much steeper or gentler than the trail's " +
						"average grade end to end.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the derivative is the limit of the secant slope as the gap " +
						"shrinks to zero.' Not like this: treating a secant slope at some small but " +
						"fixed h (like h=0.1) as already exact — it's an approximation that gets " +
						"better as h shrinks, and confusing an average rate of change over a " +
						"non-zero interval with the instantaneous rate at a single point is the " +
						"single most common mix-up here; they only agree when the function is a " +
						"straight line over that interval.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "x0", Label: "Point x0", Min: -2, Max: 2, Step: 0.1, Def: 1},
			{Key: "h", Label: "Secant gap h", Min: 0.05, Max: 2, Step: 0.05, Def: 1},
		},
		Render: render,
	})
}

// F is the sample function this concept visualizes, f(x) = x². Pure in,
// pure out — no globals, no time, no randomness.
func F(x float64) float64 {
	return x * x
}

// SecantSlope is the slope of the line through (x0, F(x0)) and
// (x0+h, F(x0+h)) — an average rate of change over the interval [x0, x0+h].
func SecantSlope(x0, h float64) float64 {
	return (F(x0+h) - F(x0)) / h
}

// Derivative is the exact, analytic derivative of F at x0: f'(x) = 2x. It's
// the limit SecantSlope(x0, h) approaches as h shrinks to 0.
func Derivative(x0 float64) float64 {
	return 2 * x0
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 320, -2.5, 2.5, 0, 4.5).Axes().String()
}
