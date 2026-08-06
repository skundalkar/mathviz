// Package correlation shows what the Pearson correlation coefficient r
// actually looks like: a scatter cloud that goes from a tight downhill line
// (r = -1) through a shapeless blob (r = 0) to a tight uphill line (r = +1).
// The picture also carries the standard warning label — a strong r only says
// two variables move together, never that one causes the other.
package correlation

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "correlation",
		Title: "Correlation (r)",
		Blurb: "r measures how tightly a scatter of points hugs a straight line, from " +
			"-1 (perfect downhill) through 0 (no linear pattern) to +1 (perfect uphill). " +
			"Slide r and watch the cloud tighten into a line or spread into a blob. The " +
			"orange line is the trend r implies. A strong r only tells you the two " +
			"variables move together — it never tells you that one causes the other.",
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
	_ = p
	return viz.New(680, 400, -1, 1, -1, 1).String()
}
