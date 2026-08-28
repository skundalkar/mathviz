// Package lagrange visualizes Lagrange multipliers: finding a function's
// maximum or minimum while being forced to stay on a constraint curve, by
// matching the objective's gradient to the constraint's gradient instead of
// searching every point on the curve. The running example maximizes and
// minimizes f(x,y)=xy over the circle x²+y²=r², where the optimal points
// are exactly where the two gradients point the same (or opposite) way.
package lagrange

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "lagrange-multipliers",
		Seq:   75,
		Title: "Lagrange multipliers (matching gradients to find a constrained max/min)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`gradient-descent` hunts for a function's minimum by following its " +
						"gradient downhill, free to move anywhere in the plane. But plenty of " +
						"real optimization isn't free to move anywhere: split a fixed budget " +
						"between two purchases to maximize satisfaction, or find the point on a " +
						"circle of radius r that makes x·y as large as possible. Setting the " +
						"ordinary partial derivatives of f(x,y)=xy to zero — the unconstrained " +
						"way to hunt for a max — only finds the origin, which is both the wrong " +
						"answer (it's a saddle point, not a max) and, the moment a constraint " +
						"like x²+y²=r² is added, not even a point you're allowed to be at unless " +
						"r=0. How do you find where a function is largest or smallest while being " +
						"forced to stay on a curve?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Picture walking along the constraint curve g(x,y)=x²+y²−r²=0, " +
						"parameterized by time t as (x(t),y(t)). By the chain rule (`chain-rule`), " +
						"how f changes as you walk is d/dt f(x(t),y(t)) = ∇f·(dx/dt, dy/dt) — the " +
						"objective's gradient dotted with the curve's tangent direction. At a max " +
						"or min of f along the curve, that rate of change must hit exactly zero, " +
						"same as any 1-D max or min — so ∇f must be perpendicular to the curve's " +
						"tangent right there. The constraint's own gradient ∇g is also always " +
						"perpendicular to the curve g=0's tangent — that's what a gradient does, " +
						"it points straight across a level curve. Two vectors both perpendicular " +
						"to the same line have to be parallel to each other: ∇f = λ∇g for some " +
						"number λ, the Lagrange multiplier. Instead of checking every point on the " +
						"curve, only the points where the two gradients line up can possibly be " +
						"the answer.",
					"Work it out for f(x,y)=xy on the circle x²+y²=4 (r=2): ∇f=(y,x), " +
						"∇g=(2x,2y). Setting ∇f=λ∇g gives y=2λx and x=2λy; dividing one by the " +
						"other gives y/x=x/y, so y²=x², i.e. y=±x. On the circle with y=x: " +
						"2x²=4, so x=y=√2≈1.414 and f=xy=2 — a maximum (λ=0.5, from y=2λx with " +
						"x=y). With y=−x: x=√2, y=−√2, f=−2 — a minimum (λ=−0.5). Only these four " +
						"points (the two solutions and their mirror images) are candidates at all; " +
						"every other point on the circle has ∇f and ∇g pointing in visibly " +
						"different directions.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"r sets the constraint circle's radius; theta drags a point around it. A few " +
						"contour curves of f(x,y)=xy (hyperbolas) are drawn faintly for scale; the " +
						"blue arrow at the current point is ∇f, the orange arrow is ∇g. Drag theta " +
						"toward 45° or 225° and the two arrows swing into pointing the same way " +
						"(the maximum, f=r²/2); toward 135° or 315° and they point opposite ways " +
						"(the minimum, f=−r²/2). Everywhere else the two arrows visibly diverge, " +
						"and the readout's 'alignment' number — the arrows' 2-D cross product, " +
						"zero exactly when they're parallel — confirms it numerically instead of " +
						"just by eye.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Solve a constrained max/min exactly with algebra — four candidate points on " +
						"an entire continuous circle — instead of scanning every point on the " +
						"curve for the best one. And λ itself isn't just an algebra byproduct: it's " +
						"the constrained max's sensitivity to loosening the constraint. Here, the " +
						"maximum value is exactly r²/2 (plug r=2: 4/2=2, matching the worked " +
						"example), so d(max f)/d(r²) = 1/2 — precisely the λ=0.5 found above. " +
						"Loosen the budget (raise r²) by a little, and λ tells you, to first order, " +
						"how much the best achievable outcome improves per unit of extra room — " +
						"before re-solving anything.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Economics: maximizing satisfaction subject to a fixed budget gives λ as the " +
						"marginal utility per extra dollar — literally the 'shadow price' of " +
						"loosening the budget constraint. Engineering: minimizing the material used " +
						"in a beam subject to a required strength. Machine learning: a support " +
						"vector machine maximizes its margin subject to every point being " +
						"correctly classified, and the points with a nonzero λ at the solution are " +
						"exactly the 'support vectors' that define the boundary. Physics: " +
						"constrained mechanics (a bead forced to stay on a wire) uses the same " +
						"gradient-matching condition.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: '∇f=λ∇g only narrows the search to a handful of " +
						"candidate points — check each one's actual f value to tell the max from " +
						"the min,' since the condition alone doesn't say which is which (both " +
						"(√2,√2) and (√2,−√2) satisfy it here, one a max, one a min).",
					"Not like this: eyeballing two arrows on a picture and calling them 'basically " +
						"parallel' — near a solution the arrows can look close without the cross " +
						"product actually being zero, and a closed constraint curve like a circle " +
						"generally has more than one gradient-matching point, so 'found one' isn't " +
						"the same as 'found the best one.'",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "r", Label: "r (constraint radius)", Min: 0.5, Max: 3, Step: 0.1, Def: 2},
			{Key: "theta", Label: "theta (angle around the circle)", Min: 0, Max: 360, Step: 1, Def: 30, Unit: "°"},
		},
		Render: render,
	})
}

// F is the objective being optimized, f(x,y) = xy.
func F(x, y float64) float64 { return x * y }

// GradF returns the objective's gradient (∂f/∂x, ∂f/∂y) = (y, x).
func GradF(x, y float64) (float64, float64) { return y, x }

// GradG returns the gradient of the constraint g(x,y) = x²+y²-r², which is
// (2x, 2y) regardless of r (r only shifts g's level, not its shape).
func GradG(x, y float64) (float64, float64) { return 2 * x, 2 * y }

// PointOnCircle returns the point at angle thetaDeg (measured counter-
// clockwise from the positive x-axis) on the circle of radius r centered
// at the origin -- the constraint curve x²+y²=r².
func PointOnCircle(r, thetaDeg float64) (float64, float64) {
	rad := thetaDeg * math.Pi / 180
	return r * math.Cos(rad), r * math.Sin(rad)
}

// Cross returns the 2-D cross product of ∇f and ∇g at (x,y):
// fx*gy - fy*gx. It is exactly zero when the two gradients are parallel
// (pointing the same way or directly opposite) -- the Lagrange condition
// ∇f=λ∇g -- and its magnitude otherwise measures how far the point is from
// satisfying it.
func Cross(x, y float64) float64 {
	fx, fy := GradF(x, y)
	gx, gy := GradG(x, y)
	return fx*gy - fy*gx
}

// Lambda estimates the Lagrange multiplier λ at (x,y) as the least-squares
// projection of ∇f onto ∇g: λ = (∇f·∇g)/(∇g·∇g). This equals the exact λ
// from ∇f=λ∇g whenever the gradients are actually parallel there, and is
// otherwise the closest scalar multiple of ∇g to ∇f -- useful for showing
// a continuously-varying readout as theta sweeps past the exact solutions.
func Lambda(x, y float64) float64 {
	fx, fy := GradF(x, y)
	gx, gy := GradG(x, y)
	denom := gx*gx + gy*gy
	if denom == 0 {
		return 0
	}
	return (fx*gx + fy*gy) / denom
}

// ContourPoints samples the curve x*y=k (one branch of a hyperbola) at n+1
// evenly spaced x values across [xlo, xhi], skipping x values too close to
// 0 (where y=k/x blows up). Mirrors viz.Sample's signature and skip-bad-
// values behavior, but for an implicit curve instead of a plain function.
func ContourPoints(k, xlo, xhi float64, n int) [][2]float64 {
	if n < 1 {
		n = 1
	}
	const epsX = 0.05
	pts := make([][2]float64, 0, n+1)
	for i := 0; i <= n; i++ {
		x := xlo + (xhi-xlo)*float64(i)/float64(n)
		if math.Abs(x) < epsX {
			continue
		}
		pts = append(pts, [2]float64{x, k / x})
	}
	return pts
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(560, 560, -3, 3, -3, 3).String()
}
