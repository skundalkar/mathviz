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
		Seq:   15,
		Title: "Gradient descent",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Drop a ball on the side of a bowl-shaped valley and let physics take over: " +
						"it rolls toward the bottom, guided continuously at every instant by " +
						"which way is downhill. An algorithm can't move continuously — it has to " +
						"take discrete steps — so it needs a rule physics never had to bother " +
						"with: how big should each step be?",
					"That size is the 'learning rate,' and it's not a minor tuning knob — get it " +
						"wrong in either direction and the algorithm fails in a completely " +
						"different way: creep toward the bottom so slowly it never practically " +
						"arrives, or overshoot the valley so badly it ends up further away with " +
						"every step instead of closer.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Gradient descent is the ball-rolling idea turned into an algorithm: at each " +
						"step, look at the slope where you're standing (the gradient) and move " +
						"against it by learning_rate × slope. Starting at x=4.5 on the bowl " +
						"f(x)=x² (slope = 2x):",
					"• Learning rate 0.3: 4.5 → 1.8 → 0.72 → 0.29 — each step lands closer to the " +
						"minimum at x=0, cautious and steady.",
					"• Learning rate 1.1 — just past this bowl's critical value of 1.0 — the same " +
						"start goes 4.5 → -5.4 → 6.48 → -7.78: the sign flips every step, " +
						"crossing back and forth over the bottom, and the distance from zero " +
						"grows each time instead of shrinking. That's divergence: the algorithm " +
						"didn't get stuck, it actively made things worse by taking steps too " +
						"large for the curvature of the valley it's in.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The grey curve is the valley, f(x)=x², and the dots trace the ball's " +
						"position at each step starting from x=4.5. Dragging the learning rate " +
						"slider is exactly the 0.3-vs-1.1 comparison above: below the critical " +
						"value the dots march steadily toward the bottom and settle there; past " +
						"it they alternate sides and visibly fly further from the minimum with " +
						"every bounce, eventually off the edge of the chart. The steps slider " +
						"controls how many of those moves get drawn.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can now tell, from watching a handful of steps, whether a learning rate " +
						"is safely below the critical value (steady convergence), needlessly " +
						"small (correct but slow), or past it entirely (diverging) — instead of " +
						"only discovering which one it was after a training run had already " +
						"failed.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"This is exactly how neural network training fails when the learning rate " +
						"is set too high — the loss doesn't plateau, it blows up to NaN within " +
						"the first few steps. It's also why training schedules often start with " +
						"a larger learning rate for speed and shrink it over time: big steps to " +
						"cover ground fast early on, small steps to settle precisely into the " +
						"minimum once you're close.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'We're overcorrecting — dial back how much we change each " +
						"time' is a learning-rate problem in plain English, whether it's a " +
						"thermostat, a steering wheel, or a training run: the adjustment per " +
						"attempt is too big, so it overshoots the target and the next correction " +
						"overshoots back the other way.",
					"Not like this: 'Just make bigger adjustments, we'll get there faster' — " +
						"true only up to a point. Past that point, bigger steps don't just take " +
						"longer, they actively make each round worse than the last, not better.",
				},
			},
		},
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
