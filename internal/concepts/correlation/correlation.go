// Package correlation shows what the Pearson correlation coefficient r
// actually looks like: a scatter cloud that goes from a tight downhill line
// (r = -1) through a shapeless blob (r = 0) to a tight uphill line (r = +1).
// The picture also carries the standard warning label — a strong r only says
// two variables move together, never that one causes the other.
package correlation

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "correlation",
		Title: "Correlation (r)",
		Blurb: "Hours studied 1,2,3,4 vs. test score 60,70,80,90 — each extra hour is worth " +
			"exactly +10 points, points sit dead on a straight uphill line: that's r = +1. " +
			"Hours of sleep lost 0,1,2,3 vs. focus score 90,80,70,60 — same line, downhill: " +
			"r = -1. Shoe size 8,9,10,11 vs. test score 72,58,81,65 — no pattern, r works out " +
			"to about 0.03: essentially r = 0. That's the whole scale, in real numbers. Now: " +
			"ice cream sales and drowning deaths both spike every summer, giving a strong " +
			"positive r — same tight-line signature as the hours-studied example. Does ice " +
			"cream cause drowning? Obviously not — both are pulled along by a third thing " +
			"(hot weather: more swimming, more ice cream). r only ever measures how tightly " +
			"points hug a straight line; it is blind to *why* the line is there. Slide r and " +
			"watch the scatter cloud tighten into a line or spread into a shapeless blob; the " +
			"orange line is the trend r implies. A strong r is a real, useful signal that " +
			"something is going on — it just never tells you what.",
		Params: []concept.ParamSpec{
			{Key: "r", Label: "Target correlation (r)", Min: -1, Max: 1, Step: 0.05, Def: 0.6},
			{Key: "n", Label: "Sample size (n)", Min: 20, Max: 300, Step: 10, Def: 150},
		},
		Render: render,
	})
}

// hash01 turns an integer seed into a deterministic, evenly-scattered value in
// [0, 1). It is a fixed function of its input — not a random-number generator
// — so GeneratePoints stays pure: same (r, n) always produces the same cloud.
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// GeneratePoints returns n (x, y) pairs whose population correlation is r.
// Both x and y are standardized (mean 0, variance 1): x is standard-normal
// noise, and y = r*x + sqrt(1-r^2)*noise, so the correlation is exactly r in
// expectation and the theoretical trend line has slope r through the origin.
// The "noise" comes from hash01 fed through a Box-Muller transform, so the
// whole function is pure and deterministic — no math/rand, no time.
func GeneratePoints(r float64, n int) (xs, ys []float64) {
	if n < 1 {
		n = 1
	}
	if r > 1 {
		r = 1
	}
	if r < -1 {
		r = -1
	}
	xs = make([]float64, n)
	ys = make([]float64, n)
	tail := math.Sqrt(1 - r*r)
	for i := 0; i < n; i++ {
		u1 := hash01(2*i + 1)
		u2 := hash01(2*i + 2)
		if u1 < 1e-9 {
			u1 = 1e-9 // keep log() finite
		}
		mag := math.Sqrt(-2 * math.Log(u1))
		z1 := mag * math.Cos(2*math.Pi*u2)
		z2 := mag * math.Sin(2*math.Pi*u2)
		xs[i] = z1
		ys[i] = r*z1 + tail*z2
	}
	return xs, ys
}

// PearsonR is the sample Pearson correlation coefficient of xs and ys. Pure
// math, unit-tested below. Returns 0 if either series has zero variance
// (undefined correlation, rather than dividing by zero).
func PearsonR(xs, ys []float64) float64 {
	n := len(xs)
	if n == 0 || len(ys) != n {
		return 0
	}
	var sx, sy float64
	for i := 0; i < n; i++ {
		sx += xs[i]
		sy += ys[i]
	}
	mx, my := sx/float64(n), sy/float64(n)

	var cov, vx, vy float64
	for i := 0; i < n; i++ {
		dx, dy := xs[i]-mx, ys[i]-my
		cov += dx * dy
		vx += dx * dx
		vy += dy * dy
	}
	if vx <= 0 || vy <= 0 {
		return 0
	}
	return cov / math.Sqrt(vx*vy)
}

func render(p map[string]float64) string {
	r, n := p["r"], int(p["n"])
	if r > 1 {
		r = 1
	}
	if r < -1 {
		r = -1
	}
	if n < 2 {
		n = 2
	}

	xs, ys := GeneratePoints(r, n)
	sampleR := PearsonR(xs, ys)

	const lim = 3.6
	c := viz.New(680, 400, -lim, lim, -lim, lim)
	c.Axes()
	for x := -3.0; x <= 3.0; x++ {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	// Reference lines through the origin (both axes are standardized).
	c.Path([][2]float64{{-lim, 0}, {lim, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, -lim}, {0, lim}}, viz.Muted, 1)

	// The scatter cloud itself: one small square per point.
	for i := range xs {
		px, py := c.X(xs[i]), c.Y(ys[i])
		c.Rect(px-2.5, py-2.5, 5, 5, viz.Accent, 0.55)
	}

	// The trend line r implies: slope r through the origin, since x and y are
	// both standardized to mean 0, variance 1.
	c.Path([][2]float64{{-lim, -lim * r}, {lim, lim * r}}, viz.Warm, 2.5)

	c.Text(20, 24, fmt.Sprintf("target r = %.2f    sample r ≈ %.2f    n = %d", r, sampleR, n),
		14, viz.Ink, "start")
	c.Text(20, 44, "r = -1 downhill line   r = 0 no linear pattern   r = +1 uphill line",
		12, viz.Muted, "start")
	c.Text(20, 62, "a strong r means the points move together — it never proves one causes the other",
		12, viz.Muted, "start")
	return c.String()
}
