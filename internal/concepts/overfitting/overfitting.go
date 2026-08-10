// Package overfitting shows what "the model fits the training data too well"
// actually looks like: a fixed handful of noisy points, and a polynomial
// curve fit to them. Crank the model's degree up and the curve stops
// approximating the underlying trend and starts weaving through every single
// point instead — perfect on the data it was trained on, wild everywhere in
// between.
package overfitting

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "overfitting",
		Seq:   14,
		Title: "Overfitting",
		Blurb: "Give a student twelve practice problems and their answers, then ask them to " +
			"explain the pattern. One student writes down a short, general rule that gets most " +
			"of the twelve right and should generalize to problem thirteen. Another memorizes " +
			"all twelve exact answers, including the quirks and typos in the answer key — zero " +
			"mistakes on practice, but no idea what to do with a new problem. That second " +
			"student has overfit. Here a handful of noisy data points are the practice " +
			"problems, and a polynomial curve is the 'rule' fit to them. A low-degree curve " +
			"stays smooth and roughly tracks the true pattern (dashed). Push the degree up and " +
			"training error keeps dropping — the curve now passes near every dot — but it does " +
			"so by wildly overshooting between them, chasing noise instead of signal.",
		Params: []concept.ParamSpec{
			{Key: "degree", Label: "Model complexity (degree)", Min: 1, Max: 11, Step: 1, Def: 3},
			{Key: "noise", Label: "Noise", Min: 0, Max: 1, Step: 0.05, Def: 0.4},
		},
		Render: render,
	})
}

// TrueFunction is the smooth underlying pattern the noisy data is secretly
// following. It's what a model that generalizes well should end up tracking.
func TrueFunction(x float64) float64 {
	return math.Sin(x)
}

// hash01 turns an integer seed into a deterministic pseudo-random value in
// [0, 1). It is a fixed function of its input, not a random-number generator,
// so every function built on it stays pure: same inputs, same picture.
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// DataPoints returns n fixed, evenly-spaced x's over [-3.2, 3.2] and their
// noisy y's: TrueFunction(x) plus deterministic pseudo-random noise in
// [-noiseAmp, noiseAmp].
func DataPoints(n int, noiseAmp float64) (xs, ys []float64) {
	if n < 2 {
		n = 2
	}
	const lo, hi = -3.2, 3.2
	xs = make([]float64, n)
	ys = make([]float64, n)
	for i := 0; i < n; i++ {
		x := lo + (hi-lo)*float64(i)/float64(n-1)
		noise := (hash01(i+1)*2 - 1) * noiseAmp
		xs[i] = x
		ys[i] = TrueFunction(x) + noise
	}
	return xs, ys
}

// PolyFit fits a degree-d polynomial y ≈ c0 + c1*x + ... + cd*x^d to (xs, ys)
// by ordinary least squares, returning the coefficients c0..cd. It solves the
// normal equations (VᵀV)c = Vᵀy, where V is the Vandermonde matrix of xs,
// via Gaussian elimination with partial pivoting. Pure math: same inputs
// always produce the same coefficients.
func PolyFit(xs, ys []float64, degree int) []float64 {
	n := len(xs)
	if degree < 0 {
		degree = 0
	}
	if degree > n-1 {
		degree = n - 1 // can't fit more free parameters than data points
	}
	k := degree + 1

	// Build the normal-equation matrix A (k x k) and vector b (k), where
	// A[i][j] = sum(x^(i+j)) and b[i] = sum(y * x^i).
	A := make([][]float64, k)
	b := make([]float64, k)
	for i := 0; i < k; i++ {
		A[i] = make([]float64, k)
	}
	for _, x := range xs {
		pow := make([]float64, 2*k-1)
		pow[0] = 1
		for p := 1; p < len(pow); p++ {
			pow[p] = pow[p-1] * x
		}
		for i := 0; i < k; i++ {
			for j := 0; j < k; j++ {
				A[i][j] += pow[i+j]
			}
		}
	}
	for idx, x := range xs {
		xp := 1.0
		for i := 0; i < k; i++ {
			b[i] += ys[idx] * xp
			xp *= x
		}
	}

	return solveLinearSystem(A, b)
}

// solveLinearSystem solves Ax = b for a small square system via Gaussian
// elimination with partial pivoting. Returns a zero vector if A is singular.
func solveLinearSystem(A [][]float64, b []float64) []float64 {
	k := len(b)
	// Work on copies so callers' data is never mutated.
	m := make([][]float64, k)
	for i := range m {
		m[i] = append([]float64(nil), A[i]...)
	}
	v := append([]float64(nil), b...)

	for col := 0; col < k; col++ {
		// Partial pivot: swap in the row with the largest value in this column.
		pivot := col
		for r := col + 1; r < k; r++ {
			if math.Abs(m[r][col]) > math.Abs(m[pivot][col]) {
				pivot = r
			}
		}
		m[col], m[pivot] = m[pivot], m[col]
		v[col], v[pivot] = v[pivot], v[col]

		if math.Abs(m[col][col]) < 1e-12 {
			continue // singular in this column; leave the coefficient at 0
		}
		for r := col + 1; r < k; r++ {
			f := m[r][col] / m[col][col]
			for c := col; c < k; c++ {
				m[r][c] -= f * m[col][c]
			}
			v[r] -= f * v[col]
		}
	}

	x := make([]float64, k)
	for i := k - 1; i >= 0; i-- {
		sum := v[i]
		for j := i + 1; j < k; j++ {
			sum -= m[i][j] * x[j]
		}
		if math.Abs(m[i][i]) < 1e-12 {
			x[i] = 0
			continue
		}
		x[i] = sum / m[i][i]
	}
	return x
}

// PolyEval evaluates a polynomial with the given coefficients (c0..cd, low
// degree first) at x.
func PolyEval(coeffs []float64, x float64) float64 {
	result := 0.0
	xp := 1.0
	for _, c := range coeffs {
		result += c * xp
		xp *= x
	}
	return result
}

// MSE returns the mean squared error of the fitted polynomial against the
// given (xs, ys), i.e. the training error.
func MSE(xs, ys, coeffs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sum := 0.0
	for i, x := range xs {
		d := PolyEval(coeffs, x) - ys[i]
		sum += d * d
	}
	return sum / float64(len(xs))
}

// TrueError returns the mean squared error between the fitted polynomial and
// the noiseless TrueFunction, sampled densely across the data's x-range. This
// is the generalization gap the training MSE alone can't see: a high-degree
// fit can drive training MSE near zero while TrueError climbs, because it's
// chasing noise instead of the pattern.
func TrueError(coeffs []float64) float64 {
	const lo, hi, n = -3.2, 3.2, 200
	sum := 0.0
	for i := 0; i <= n; i++ {
		x := lo + (hi-lo)*float64(i)/float64(n)
		d := PolyEval(coeffs, x) - TrueFunction(x)
		sum += d * d
	}
	return sum / float64(n+1)
}

// numPoints is the fixed size of the "practice problem" dataset — small
// enough that a high enough degree can pass through every point exactly.
const numPoints = 12

func render(p map[string]float64) string {
	degree, noise := int(p["degree"]), p["noise"]

	xs, ys := DataPoints(numPoints, noise)
	coeffs := PolyFit(xs, ys, degree)

	const xmin, xmax = -3.5, 3.5
	const ymin, ymax = -2.5, 2.5
	c := viz.New(680, 340, xmin, xmax, ymin, ymax)
	c.Axes()
	for x := -3.0; x <= 3.0; x++ {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	trueCurve := viz.Sample(xmin, xmax, 240, TrueFunction)
	c.Path(trueCurve, viz.Muted, 1.5)

	fitCurve := viz.Sample(xmin, xmax, 240, func(x float64) float64 {
		return PolyEval(coeffs, x)
	})
	c.Path(fitCurve, viz.Accent, 2.5)

	// The noisy "practice problems" the model was trained on.
	for i := range xs {
		px, py := c.X(xs[i]), c.Y(ys[i])
		c.Rect(px-3.5, py-3.5, 7, 7, viz.Warm, 0.85)
	}

	trainErr := MSE(xs, ys, coeffs)
	trueErr := TrueError(coeffs)
	c.Text(20, 24, fmt.Sprintf("degree = %d    noise = %.2f", degree, noise), 14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("training error ≈ %.3f    true error ≈ %.3f", trainErr, trueErr),
		13, viz.Muted, "start")
	c.Text(20, 62, "dots = noisy data    grey = true pattern    blue = fitted curve",
		12, viz.Muted, "start")

	return c.String()
}
