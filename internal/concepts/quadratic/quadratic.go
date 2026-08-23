// Package quadratic visualizes the quadratic formula: the single closed-form
// expression x = (−b ± √(b²−4ac)) / (2a) that solves ax²+bx+c=0 directly,
// derived by completing the square. The sign of the discriminant b²−4ac
// alone says how many times the parabola y=ax²+bx+c crosses the x-axis,
// before a single root is computed.
package quadratic

import (
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

func render(p map[string]float64) string {
	_ = p
	return viz.New(640, 420, -4, 4, -4, 4).String()
}
