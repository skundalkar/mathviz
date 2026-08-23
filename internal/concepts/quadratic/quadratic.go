// Package quadratic visualizes the quadratic formula: the single closed-form
// expression x = (−b ± √(b²−4ac)) / (2a) that solves ax²+bx+c=0 directly,
// derived by completing the square. The sign of the discriminant b²−4ac
// alone says how many times the parabola y=ax²+bx+c crosses the x-axis,
// before a single root is computed.
package quadratic

import (
	"fmt"
	"math"
	"math/cmplx"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "quadratic-formula",
		Seq:   61,
		Title: "Quadratic formula (roots via the discriminant)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`newtons-method` found a function's root by guessing a starting point and " +
						"following tangent lines closer and closer, one step at a time — it took " +
						"several iterations to close in on √2. That works for almost any function. " +
						"But for the specific shape ax²+bx+c=0, is there a single-shot formula, " +
						"computed straight from a, b, and c, that jumps directly to the exact root(s) " +
						"— no guessing, no iterating, no graphing required?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Yes — derive it once by 'completing the square,' then reuse it forever. " +
						"Start from ax²+bx+c=0, divide every term by a: x² + (b/a)x + c/a = 0. Move " +
						"the constant over: x² + (b/a)x = −c/a. Add (b/2a)² to both sides — chosen " +
						"specifically because it makes the left side a perfect square: (x + b/2a)² = " +
						"(b/2a)² − c/a = (b²−4ac)/(4a²). Take the square root of both sides (± because " +
						"both a positive and negative square root solve it) and subtract b/2a: x = " +
						"(−b ± √(b²−4ac)) / (2a). The quantity under the square root, D = b²−4ac, is " +
						"called the discriminant.",
					"Try it on x²−5x+6=0 (a=1, b=−5, c=6): D = (−5)² − 4×1×6 = 25−24 = 1, so x = " +
						"(5 ± 1)/2 = 3 or 2. Check both by plugging back in: 3²−5×3+6 = 9−15+6 = 0 ✓, " +
						"2²−5×2+6 = 4−10+6 = 0 ✓.",
					"The discriminant's sign alone — before finishing the arithmetic — says how " +
						"many real roots to expect:",
					"• D>0 (like the example above, D=1): the ± produces two different real " +
						"numbers, so the parabola crosses the x-axis twice.",
					"• D=0, e.g. x²+2x+1=0 (a=1,b=2,c=1): D = 4−4 = 0, so + and − both give the " +
						"same value, x = −2/2 = −1 — one repeated root, where the parabola's vertex " +
						"sits exactly on the x-axis, touching it without crossing.",
					"• D<0, e.g. x²+x+1=0 (a=1,b=1,c=1): D = 1−4 = −3, a negative number under a " +
						"square root — no real number solves it, so the parabola never touches the " +
						"x-axis at all. The formula still produces an answer, using i=√−1: x = " +
						"(−1 ± i√3)/2 ≈ −0.5 ± 0.866i, a pair of complex roots.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The three sliders a, b, c are the parabola's coefficients directly. The blue " +
						"curve is y=ax²+bx+c; the dashed vertical line is its axis of symmetry, x = " +
						"−b/(2a); the marked point is the vertex, sitting exactly on that line. When " +
						"D≥0, the real root(s) are marked where the curve crosses y=0 — one point when " +
						"D=0 (the vertex itself touches the axis), two when D>0. When D<0 the curve " +
						"floats entirely above or below the x-axis and no root markers appear — the " +
						"readout reports the complex pair instead.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Solve any quadratic in one step, without guessing factors or running " +
						"`newtons-method`'s iteration — and, just from the sign of b²−4ac, know in " +
						"advance whether a quadratic model even has a real solution before computing " +
						"one. That second part matters whenever 'solve for x' is really a physical " +
						"question: does a projectile's height ever reach a given target, does a " +
						"break-even curve actually cross zero profit, does a beam under load reach a " +
						"critical stress — D<0 answers 'no, never' immediately, without solving " +
						"anything further.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A ball's height over time under gravity is a downward parabola, h(t) = " +
						"−½gt²+v₀t+h₀; solving h(t)=0 for t (when does it land?) is exactly this " +
						"formula. Break-even analysis in business often models profit as a quadratic " +
						"in price or quantity, and the roots mark the exact prices where profit " +
						"crosses zero. Satellite dishes and telescope mirrors are shaped as parabolas " +
						"specifically because of a reflective property tied to this same curve. And " +
						"the discriminant's sign shows up again well beyond quadratics — the same " +
						"'does this even have a real solution' question, computed a different way, " +
						"is exactly what `newtons-method` can silently fail to answer if you start it " +
						"chasing a root that was never real to begin with.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: the ± is not optional bookkeeping — 'x = (−b+√D)/(2a)' alone " +
						"is only half the answer whenever D>0, since −√D solves the equation exactly " +
						"as validly as +√D does.",
					"Not like this: assuming D<0 means 'no solution, full stop.' It means no REAL " +
						"solution — the formula still produces two valid complex roots, and whether " +
						"that counts as 'no answer' depends entirely on whether the problem you're " +
						"modeling only makes sense for real numbers (a ball's landing time can't be " +
						"complex) or genuinely wants the complex pair (some physics and engineering " +
						"problems do).",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "a", Min: -3, Max: 3, Step: 0.1, Def: 1},
			{Key: "b", Label: "b", Min: -6, Max: 6, Step: 0.1, Def: -5},
			{Key: "c", Label: "c", Min: -6, Max: 6, Step: 0.1, Def: 6},
		},
		Render: render,
	})
}

// Discriminant returns b²−4ac, whose sign says how many real roots
// ax²+bx+c=0 has: positive means two, zero means one repeated, negative
// means none (only a complex pair).
func Discriminant(a, b, c float64) float64 {
	return b*b - 4*a*c
}

// Evaluate returns ax²+bx+c at x — the parabola's height, used both to draw
// the curve and to verify a computed root actually zeroes the expression.
func Evaluate(a, b, c, x float64) float64 {
	return a*x*x + b*x + c
}

// Vertex returns the parabola's turning point (h, k): h = −b/(2a) is the
// axis of symmetry, k = Evaluate(a,b,c,h) is the minimum (a>0) or maximum
// (a<0) value. a=0 isn't a quadratic; callers keep a away from 0.
func Vertex(a, b, c float64) (h, k float64) {
	h = -b / (2 * a)
	k = Evaluate(a, b, c, h)
	return h, k
}

// Roots returns both solutions of ax²+bx+c=0 via the quadratic formula,
// x = (−b ± √(b²−4ac)) / (2a), as complex128 so the same formula handles
// every discriminant sign uniformly: real (possibly equal) roots come back
// with a zero imaginary part when D>=0, a genuine complex conjugate pair
// when D<0. a=0 isn't a quadratic; callers keep a away from 0.
func Roots(a, b, c float64) (r1, r2 complex128) {
	d := complex(Discriminant(a, b, c), 0)
	sqrtD := cmplx.Sqrt(d)
	denom := complex(2*a, 0)
	r1 = (complex(-b, 0) + sqrtD) / denom
	r2 = (complex(-b, 0) - sqrtD) / denom
	return r1, r2
}

func render(p map[string]float64) string {
	a, b, cc := p["a"], p["b"], p["c"]
	if math.Abs(a) < 0.2 {
		a = math.Copysign(0.2, a)
		if a == 0 {
			a = 0.2
		}
	}

	d := Discriminant(a, b, cc)
	h, k := Vertex(a, b, cc)

	// Window wide enough to show both real roots with room to spare, or a
	// fixed default width when there are none to anchor on.
	xHalf := 4.0
	if d >= 0 {
		spread := math.Sqrt(d) / math.Abs(a)
		xHalf = math.Max(4, spread*1.3+1.5)
	}
	xMin, xMax := h-xHalf, h+xHalf

	curve := viz.Sample(xMin, xMax, 200, func(x float64) float64 {
		return Evaluate(a, b, cc, x)
	})

	yLo, yHi := 0.0, 0.0
	for _, pt := range curve {
		if pt[1] < yLo {
			yLo = pt[1]
		}
		if pt[1] > yHi {
			yHi = pt[1]
		}
	}
	if k < yLo {
		yLo = k
	}
	if k > yHi {
		yHi = k
	}
	padY := (yHi - yLo) * 0.12
	if padY == 0 {
		padY = 1
	}
	yLo -= padY
	yHi += padY

	c := viz.New(640, 420, xMin, xMax, yLo, yHi)
	c.Axes()
	xStep := (xMax - xMin) / 8
	for x := xMin; x <= xMax; x += xStep {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	c.Path(curve, viz.Accent, 2.5)
	c.VLine(h, viz.Muted, true)

	vpx, vpy := c.X(h), c.Y(k)
	c.Rect(vpx-4, vpy-4, 8, 8, viz.Ink, 0.9)

	r1, r2 := Roots(a, b, cc)
	var rootText string
	switch {
	case d > 1e-9:
		x1, x2 := real(r1), real(r2)
		for _, rx := range []float64{x1, x2} {
			px, py := c.X(rx), c.Y(0)
			c.Rect(px-4, py-4, 8, 8, viz.Good, 1)
		}
		rootText = fmt.Sprintf("two real roots: x = %.3f, %.3f", x1, x2)
	case d > -1e-9:
		x1 := real(r1)
		px, py := c.X(x1), c.Y(0)
		c.Rect(px-4, py-4, 8, 8, viz.Good, 1)
		rootText = fmt.Sprintf("one repeated real root: x = %.3f (touches the x-axis)", x1)
	default:
		rootText = fmt.Sprintf("no real roots: x = %.3f ± %.3fi", real(r1), math.Abs(imag(r1)))
	}

	c.Text(16, 24, fmt.Sprintf("y = %.2fx² + %.2fx + %.2f", a, b, cc), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("D = b²−4ac = %.2f² − 4×%.2f×%.2f = %.3f", b, a, cc, d), 13, viz.Muted, "start")
	c.Text(16, 64, rootText, 14, viz.Warm, "start")
	c.Text(16, 84, fmt.Sprintf("vertex (h,k) = (%.3f, %.3f)", h, k), 13, viz.Muted, "start")

	return c.String()
}
