// Package taylorseries visualizes a Taylor series: building a polynomial
// approximation of a curve out of nothing but its derivatives at a single
// point, and watching that approximation improve — near the center
// immediately, and over a widening range as more derivative terms are
// added. The example function is sin(x), expanded around x0=0.
package taylorseries

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "taylor-series",
		Seq:   45,
		Title: "Taylor series (approximating a curve near a point)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`derivative` gives you the tangent line at a point — a straight line that " +
						"matches sin(x)'s value and slope at x0=0, so it's a great approximation " +
						"very close to 0 (say x=0.1). But it's still just a straight line, and " +
						"sin(x) curves away from it fast: by x=1 the tangent line says sin(1)≈1, " +
						"while the real value is 0.841 — already off by 0.159. What if, instead of " +
						"only matching the slope at x0, you also matched the curvature (the second " +
						"derivative), and the rate the curvature itself changes (the third " +
						"derivative), and so on? Does piling on more derivatives — all still " +
						"measured at that same single point, x0=0 — make a *polynomial* " +
						"approximation good over a wider range, not just closer at that one exact " +
						"point?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"sin(x)'s derivatives at 0 cycle through four values forever: f(0)=0, f'(0)=1, " +
						"f''(0)=0, f'''(0)=-1, f''''(0)=0, f'''''(0)=1, ... The Taylor polynomial of " +
						"degree n is built from exactly those numbers: sum over k=0..n of " +
						"f^(k)(0)/k! times x^k. Try it at x=1, where the true value is " +
						"sin(1)=0.841471, adding one more term at a time:",
					"• degree 1 (just the tangent line, x itself): approx = 1.000000, error = 0.158529.",
					"• degree 3 (add -x³/3! = -x³/6): approx = 1 - 0.166667 = 0.833333, error = -0.008138.",
					"• degree 5 (add +x⁵/5! = +x⁵/120): approx = 0.841667, error = 0.000196.",
					"• degree 7 (add -x⁷/5040): approx = 0.841468, error = -0.000003 — already " +
						"accurate to 5 decimal places.",
					"Each new term is smaller than the last (x=1 keeps every power modest) and the " +
						"error shrinks fast. But try x=3, farther from the center, where true " +
						"sin(3)=0.141120: degree 1 gives approx=3.000000 (error 2.858880 — wildly " +
						"wrong), degree 5 gives 0.525000 (error 0.383880 — still well off), and it " +
						"takes degree 9 to reach 0.145312 (error 0.004192) — the same handful of " +
						"terms that nailed x=1 to five decimal places barely gets x=3 into the right " +
						"neighborhood. Matching more derivatives at x0=0 widens the range where the " +
						"polynomial is trustworthy, but a fixed, small number of terms is only ever " +
						"trustworthy up to some finite distance from the center — go far enough out " +
						"and you need more terms to catch up.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The dashed curve is the true sin(x). The solid curve is the Taylor polynomial " +
						"of the degree set by the order slider, built entirely from sin's derivatives " +
						"at x0=0 (marked with a small circle). The eval x slider marks one point; the " +
						"readout reports the true value, the polynomial's approximation there, and " +
						"the gap between them. Push order up and watch the solid curve hug the dashed " +
						"one over a wider and wider stretch around 0, while — at any fixed order — " +
						"dragging eval x far enough from 0 eventually finds where the solid curve " +
						"peels away.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Compute a curved function's value to any precision you like using nothing but " +
						"addition, multiplication, and a handful of derivative values at one point — " +
						"no calculator's sin button required, which is close to how a calculator's " +
						"sin button actually works internally. You can also judge, for a given order, " +
						"roughly how far from the center you can trust the approximation before it's " +
						"time to add more terms, rather than assuming any polynomial that matches at " +
						"the center is automatically trustworthy everywhere.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Software math libraries (libm, and every calculator built on top of it) compute " +
						"sin, cos, and exp using truncated Taylor-like polynomial approximations, not " +
						"a lookup table of every possible input. Physics leans on the very first " +
						"nonzero term alone: the small-angle approximation sin(θ)≈θ (a degree-1 " +
						"Taylor polynomial) is what makes a pendulum's period formula solvable in " +
						"introductory physics, valid because real pendulums swing through small " +
						"angles close to the center. Engineers use the same trick — linearizing a " +
						"nonlinear system near an operating point — to analyze circuits and control " +
						"systems with straightforward linear-algebra tools instead of the full " +
						"nonlinear equations.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: a Taylor polynomial is built from derivatives at ONE point " +
						"and is most trustworthy near that point, with the trustworthy range " +
						"widening — but never becoming unlimited — as more terms are added at a " +
						"fixed order.",
					"Not like this: assuming a Taylor polynomial that matches a function's value and " +
						"several derivatives at x0 must therefore fit the function well everywhere " +
						"— the x=3 numbers above show a degree-5 polynomial that's excellent at x=1 " +
						"is still off by 0.38 at x=3. Also not like this: treating the derivative " +
						"page's tangent line as a special, different tool from a Taylor polynomial — " +
						"it IS one, the degree-1 case, and every higher-order Taylor polynomial is " +
						"the same idea with more derivative terms folded in.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "order", Label: "Polynomial degree (order)", Min: 0, Max: 11, Step: 1, Def: 3},
			{Key: "evalX", Label: "Eval point x", Min: -4, Max: 4, Step: 0.1, Def: 1},
		},
		Render: render,
	})
}

// F is the function this concept expands: f(x) = sin(x).
func F(x float64) float64 {
	return math.Sin(x)
}

// SinDerivativeAtZero returns sin's k-th derivative evaluated at x0=0. The
// derivatives of sin cycle through four values forever — sin, cos, -sin,
// -cos, sin, ... — so at 0 they cycle through 0, 1, 0, -1 forever. Negative
// k returns 0 (undefined).
func SinDerivativeAtZero(k int) float64 {
	if k < 0 {
		return 0
	}
	cycle := [4]float64{0, 1, 0, -1}
	return cycle[k%4]
}

// Factorial returns k! for k>=0 as a float64 (the Taylor coefficient below
// divides by it, and k stays small enough — this package's order Param
// caps at 11 — that float64 precision is never a concern). Negative k
// returns 0 — undefined.
func Factorial(k int) float64 {
	if k < 0 {
		return 0
	}
	f := 1.0
	for i := 2; i <= k; i++ {
		f *= float64(i)
	}
	return f
}

// TaylorTerm returns the k-th term of sin's Taylor series centered at
// x0=0, evaluated at x: f^(k)(0)/k! * x^k.
func TaylorTerm(k int, x float64) float64 {
	return SinDerivativeAtZero(k) / Factorial(k) * math.Pow(x, float64(k))
}

// TaylorPolynomial returns the degree-`order` Taylor polynomial of sin
// (centered at x0=0), evaluated at x: the sum of TaylorTerm(k, x) for
// k=0..order. Negative order is treated as 0 (the constant term alone).
func TaylorPolynomial(order int, x float64) float64 {
	if order < 0 {
		order = 0
	}
	sum := 0.0
	for k := 0; k <= order; k++ {
		sum += TaylorTerm(k, x)
	}
	return sum
}

// plot window: about two full periods of sin(x) each side of the center,
// with a y range that clips a diverging polynomial at the edges rather than
// letting it shoot off-canvas.
const (
	xlim = 6.5
	ylim = 2.5
)

func clamp(y, lo, hi float64) float64 {
	if y < lo {
		return lo
	}
	if y > hi {
		return hi
	}
	return y
}

func render(params map[string]float64) string {
	order := int(params["order"] + 0.5)
	if order < 0 {
		order = 0
	}
	evalX := params["evalX"]
	if evalX < -xlim {
		evalX = -xlim
	}
	if evalX > xlim {
		evalX = xlim
	}

	c := viz.New(680, 420, -xlim, xlim, -ylim, ylim)
	c.Axes()
	for x := -xlim; x <= xlim+1e-9; x += xlim / 3 {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	// True sin(x), thin and muted -- the target the polynomial is chasing.
	truth := viz.Sample(-xlim, xlim, 300, F)
	c.Path(truth, viz.Muted, 1.5)

	// The Taylor polynomial, clamped so a diverging approximation flattens
	// against the top/bottom edge instead of shooting off-canvas.
	approxPts := viz.Sample(-xlim, xlim, 300, func(x float64) float64 {
		return clamp(TaylorPolynomial(order, x), -ylim-0.3, ylim+0.3)
	})
	c.Path(approxPts, viz.Accent, 2.5)

	// The center point, x0=0, where every derivative used was measured.
	cx, cy := c.X(0), c.Y(0)
	c.Rect(cx-4, cy-4, 8, 8, viz.Ink, 1)

	// The eval point: true value vs. this order's approximation, connected
	// by a vertical line showing the gap between them.
	trueVal := F(evalX)
	approxVal := TaylorPolynomial(order, evalX)
	c.Path([][2]float64{{evalX, trueVal}, {evalX, clamp(approxVal, -ylim, ylim)}}, viz.Warm, 1.5)
	tpx, tpy := c.X(evalX), c.Y(trueVal)
	c.Rect(tpx-3.5, tpy-3.5, 7, 7, viz.Muted, 1)
	apx, apy := c.X(evalX), c.Y(clamp(approxVal, -ylim, ylim))
	c.Rect(apx-3.5, apy-3.5, 7, 7, viz.Warm, 1)

	c.Text(16, 24, fmt.Sprintf("Taylor polynomial of sin(x), degree %d, centered at x0=0", order), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("at x=%.2f: true sin(x) = %.6f    approx = %.6f    error = %.6f",
		evalX, trueVal, approxVal, approxVal-trueVal), 13, viz.Warm, "start")
	c.Text(16, 64, "gray = true sin(x)    blue = Taylor polynomial    black square = center x0=0",
		12, viz.Muted, "start")

	return c.String()
}
