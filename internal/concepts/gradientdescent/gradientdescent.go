// Package gradientdescent shows what "learning rate" actually controls: how
// big a step a ball rolling downhill takes at each move. Too small and it
// creeps toward the bottom forever; too large and it overshoots the valley,
// bounces to the other wall, and can fly further away with every step.
package gradientdescent

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "gradient-descent",
		Title: "Gradient descent",
		Blurb: "Drop a ball on the side of a bowl-shaped valley and let physics take over: " +
			"it rolls toward the bottom, the lowest point, guided at every instant by which way " +
			"is downhill. Gradient descent is that same idea turned into an algorithm — at each " +
			"step, move against the slope (the gradient) by an amount set by the 'learning " +
			"rate'. A small learning rate takes tiny, cautious steps: it will get there, but " +
			"slowly. A learning rate that's too large overshoots the bottom entirely, lands " +
			"higher up the opposite wall, and if it's large enough the overshoot gets worse " +
			"each step instead of better — the 'ball' flies further from the minimum with every " +
			"bounce. Drag the learning rate and watch the same starting point converge, crawl, " +
			"or blow up.",
		Params: []concept.ParamSpec{
			{Key: "lr", Label: "Learning rate", Min: 0.01, Max: 1.2, Step: 0.01, Def: 0.3},
			{Key: "steps", Label: "Steps", Min: 1, Max: 30, Step: 1, Def: 12},
		},
		Render: render,
	})
}

// StartX is the fixed starting position of the "ball" — off-center so its
// path toward (or away from) the minimum at x=0 is visible.
const StartX = 4.5

// F is the "valley" the ball rolls down: a simple bowl with its minimum at
// x=0, f(0)=0.
func F(x float64) float64 {
	return x * x
}

// Grad is the derivative of F, i.e. the slope gradient descent follows.
func Grad(x float64) float64 {
	return 2 * x
}

// Descend runs `steps` iterations of gradient descent on F starting from x0
// with the given learning rate, returning the position at every step
// including the start: [x0, x1, ..., x_steps]. Pure math — same inputs
// always produce the same path.
//
// Because F(x)=x^2 has Grad(x)=2x, each step is x_{t+1} = x_t - lr*2*x_t =
// x_t*(1-2*lr): the path is a geometric sequence, |1-2*lr| < 1 converges,
// |1-2*lr| > 1 diverges. That closed form is what the tests check against.
func Descend(x0, lr float64, steps int) []float64 {
	if steps < 0 {
		steps = 0
	}
	path := make([]float64, steps+1)
	path[0] = x0
	x := x0
	for i := 1; i <= steps; i++ {
		x = x - lr*Grad(x)
		path[i] = x
	}
	return path
}

func render(p map[string]float64) string {
	lr, steps := p["lr"], int(p["steps"])

	const xmin, xmax = -6.5, 6.5
	const ymin, ymax = -3.0, 44.0
	c := viz.New(680, 340, xmin, xmax, ymin, ymax)
	c.Axes()
	for x := -6.0; x <= 6.0; x += 2 {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	// The valley itself.
	c.Path(viz.Sample(xmin, xmax, 200, F), viz.Muted, 2)

	path := Descend(StartX, lr, steps)

	// Draw the ball's trajectory as straight jumps between successive
	// positions on the curve, stopping if it flies off the visible window —
	// that's what a diverging (too-large) learning rate looks like.
	diverged := false
	visible := 0
	for i, x := range path {
		if math.Abs(x) > xmax {
			diverged = true
			break
		}
		visible = i
		if i > 0 {
			prev := path[i-1]
			c.Path([][2]float64{{prev, F(prev)}, {x, F(x)}}, viz.Warm, 1.5)
		}
	}
	for i := 0; i <= visible; i++ {
		x := path[i]
		px, py := c.X(x), c.Y(F(x))
		if i == 0 {
			c.Rect(px-4, py-4, 8, 8, viz.Ink, 0.9) // start
		} else if i == visible && !diverged {
			c.Rect(px-4, py-4, 8, 8, viz.Good, 0.95) // final resting spot
		} else {
			c.Rect(px-3, py-3, 6, 6, viz.Accent, 0.7)
		}
	}

	c.Text(20, 24, fmt.Sprintf("learning rate = %.2f    steps = %d", lr, steps), 14, viz.Ink, "start")
	final := path[len(path)-1]
	status := fmt.Sprintf("x ≈ %.3f    f(x) ≈ %.3f — converging toward the minimum at x=0", final, F(final))
	if diverged {
		status = "diverged — the learning rate is too large, each step overshoots further"
	}
	c.Text(20, 44, status, 13, viz.Muted, "start")

	return c.String()
}
