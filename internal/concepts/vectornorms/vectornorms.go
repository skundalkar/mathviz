// Package vectornorms visualizes the L1, L2, and L∞ norms -- three
// different, equally legitimate ways to total up a vector's components
// into a single "size" number -- plus the general Lp norm that unifies
// them into one family. Built on top of the same 2D vector `vectors`
// introduced: where that concept measured size with a single fixed ruler
// (L2), this one asks what changes when you measure it differently.
package vectornorms

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "vector-norms",
		Seq:   93,
		Title: "Vector norms (L1, L2, and L∞ ways to measure size)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"vectors measured a pull's strength with one fixed formula, √(x²+y²) — " +
						"the length of the arrow, exactly the way a ruler would measure it. But " +
						"'how big is this vector' isn't actually a single fixed question: a " +
						"delivery robot moving from its depot at (0,0) to a drop-off at (3,4) " +
						"blocks away can't cut diagonally through buildings the way that " +
						"ruler-length arrow does — it has to travel 3 blocks one way and then 4 " +
						"blocks the other way, 7 blocks total, not the diagonal ruler-distance of " +
						"5. And a warehouse elevator lifting a crate 3 meters sideways and 4 " +
						"meters up at the same time doesn't care about either total — it only " +
						"cares about whichever single direction takes the longest, 4 meters, " +
						"since the other direction rides along for free. Same vector, (3,4), " +
						"three different honest answers (7, 5, 4) to 'how big is it' — so which " +
						"one is the real size, and how do you pin the other two down precisely " +
						"enough to compute them on purpose?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take vectors's own worked pull, v=(3,4), and total up its two components " +
						"into a single size number three different ways:",
					"• L2 norm (ruler length, the same formula vectors used): ||v||₂ = " +
						"√(3²+4²) = √25 = 5.",
					"• L1 norm (the delivery robot's city-block distance, also called " +
						"'Manhattan distance'): ||v||₁ = |3|+|4| = 7 — just add the absolute " +
						"values, no squaring or rooting.",
					"• L∞ norm (the elevator's slowest single direction): ||v||∞ = max(|3|,|4|) " +
						"= 4 — throw away every component except the largest.",
					"|Norm|Formula|Value at v=(3,4)|",
					"|L1|Σ abs(xᵢ) — sum of absolute values|7|",
					"|L2|√(Σxᵢ²) — Euclidean length|5|",
					"|L∞|max abs(xᵢ) — largest component|4|",
					"All three are one family in disguise, the Lp norm: ||v||_p = (|x|^p + " +
						"|y|^p)^(1/p). Set p=1 and the formula collapses to L1 exactly; set p=2 " +
						"and it collapses to L2 exactly. Slide p up further and watch it slide " +
						"toward L∞ instead of jumping there:",
					"• p=1: 7.0000 (this is L1). p=1.5: 5.5843. p=2: 5.0000 (this is L2).",
					"• p=3: 4.4979. p=5: 4.1740. p=10: 4.0220 — closing in on 4.",
					"Once p is large, the bigger component (4) raised to the p-th power " +
						"completely swamps the smaller one (3) inside the sum, so the 1/p root " +
						"effectively just hands back the largest component alone — L∞ is the " +
						"p→∞ limiting case of the exact same formula, not a separately invented " +
						"rule.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The black arrow is the vector v, positioned by the x and y sliders. " +
						"Around it are three fixed reference shapes, one per norm, each scaled " +
						"up from its own 'radius-1' shape until it passes exactly through v's " +
						"tip: a green diamond for L1, a blue circle for L2, and a red square for " +
						"L∞ — v's tip sits exactly on all three at once, but the shapes " +
						"themselves are different sizes, because the three norms don't agree on " +
						"how big v is (the diamond has to stretch to 7 to reach v, the circle " +
						"only to 5, the square only to 4). The orange curve is the general Lp " +
						"shape for whatever p the slider is set to, also always stretched to pass " +
						"through v: at p=1 it sits exactly on top of the green diamond, at p=2 " +
						"exactly on top of the blue circle, and as you push p up past about 8-10 " +
						"it hugs the red square more and more closely without ever quite becoming " +
						"a perfect corner.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Choose the size measure that actually matches the situation instead of " +
						"defaulting to ruler-length out of habit — city-block travel time, " +
						"worst-case single-axis strain, or straight-line distance — and know " +
						"exactly how they relate: L∞(v) ≤ L2(v) ≤ L1(v), always, for any vector. " +
						"You can see that ordering directly in the nested reference shapes: the " +
						"square is always innermost, the diamond outermost, the circle in " +
						"between. This ordering is also what regularization leans on without " +
						"spelling it out: an L1 penalty's diamond-shaped boundary has sharp " +
						"corners sitting exactly on the coordinate axes, which is why Lasso " +
						"regression can shrink a coefficient all the way to exactly zero, while an " +
						"L2 penalty's round boundary has no corners anywhere, which is why Ridge " +
						"regression can only shrink a coefficient toward zero, never all the way " +
						"there.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"GPS and mapping apps use L1-style 'Manhattan distance' to estimate driving " +
						"time on a city grid, since you can't cut through buildings. Structural " +
						"and electrical engineers care about L∞-style 'worst case on any single " +
						"wire/beam' when a system fails if any one component is overloaded, " +
						"regardless of the others. Machine learning uses all three as the " +
						"geometry underneath regularization: L1 regularization (Lasso) for " +
						"automatic feature selection, L2 regularization (Ridge) for smooth " +
						"shrinkage, and L∞ shows up in robustness/adversarial-example analysis, " +
						"where 'how much can an attacker change any single pixel' is exactly an " +
						"L∞ constraint. Nearest-neighbor search and clustering both have to pick " +
						"a distance metric, and the norm you choose there changes which points " +
						"count as 'close.'",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: a vector doesn't have one single 'size' — L1, L2, and " +
						"L∞ are three different, equally legitimate ways to total up the same " +
						"components, and picking the wrong one for the situation (using " +
						"straight-line distance for a city-grid delivery route, say) gives an " +
						"answer that's actually wrong for the question being asked, not just a " +
						"rounding difference. Not like this: assuming L1 ≤ L2 ≤ L∞ because " +
						"'higher p sounds like it should mean a bigger norm' — it's the opposite: " +
						"raising p in the Lp formula always shrinks or holds steady the norm of " +
						"any fixed vector with more than one nonzero component, never grows it, " +
						"because a higher power lets the single largest component increasingly " +
						"dominate the sum while the smaller ones fade toward irrelevant. L∞ is " +
						"the smallest of the three, not the largest, despite the '∞' looking like " +
						"it should mean 'biggest.'",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "x", Label: "Vector v — x", Min: -5, Max: 5, Step: 0.5, Def: 3},
			{Key: "y", Label: "Vector v — y", Min: -5, Max: 5, Step: 0.5, Def: 4},
			{Key: "p", Label: "Lp exponent (p)", Min: 1, Max: 20, Step: 0.5, Def: 2},
		},
		Render: render,
	})
}

// NormL1 returns the L1 (Manhattan / city-block) norm of the vector (x,y):
// the sum of the absolute values of its components.
func NormL1(x, y float64) float64 {
	return math.Abs(x) + math.Abs(y)
}

// NormL2 returns the L2 (Euclidean / ruler) norm of the vector (x,y): its
// straight-line length.
func NormL2(x, y float64) float64 {
	return math.Hypot(x, y)
}

// NormLInf returns the L∞ (Chebyshev / max) norm of the vector (x,y): its
// largest single component's absolute value.
func NormLInf(x, y float64) float64 {
	return math.Max(math.Abs(x), math.Abs(y))
}

// NormLp returns the general Lp norm of the vector (x,y) for any p>=1:
// (|x|^p + |y|^p)^(1/p), which collapses to NormL1 at p=1 and NormL2 at
// p=2, and converges to NormLInf as p grows without bound. The largest
// component is factored out before raising to the p-th power so the
// computation stays well-behaved (no overflow to +Inf) even for very
// large p, where a naive |x|^p would blow past float64's range long
// before the 1/p root could bring it back down.
func NormLp(x, y, p float64) float64 {
	ax, ay := math.Abs(x), math.Abs(y)
	m := math.Max(ax, ay)
	if m == 0 {
		return 0
	}
	rx, ry := ax/m, ay/m
	return m * math.Pow(math.Pow(rx, p)+math.Pow(ry, p), 1/p)
}

// boundary samples the scaled boundary of a norm's unit ball: the curve of
// points whose value under norm equals target. Since every norm here is
// homogeneous (norm(r*u) == r*norm(u) for r>=0), scaling the direction
// (cos θ, sin θ) by target/norm(cos θ, sin θ) lands exactly on that curve,
// for every θ -- so the boundary always passes exactly through whichever
// vector target was measured from.
func boundary(norm func(x, y float64) float64, target float64, steps int) [][2]float64 {
	if target <= 0 {
		return nil
	}
	pts := make([][2]float64, 0, steps+1)
	for i := 0; i <= steps; i++ {
		theta := 2 * math.Pi * float64(i) / float64(steps)
		dx, dy := math.Cos(theta), math.Sin(theta)
		r := target / norm(dx, dy)
		pts = append(pts, [2]float64{r * dx, r * dy})
	}
	return pts
}

func render(p map[string]float64) string {
	x, y := p["x"], p["y"]
	pExp := p["p"]

	l1 := NormL1(x, y)
	l2 := NormL2(x, y)
	linf := NormLInf(x, y)
	lp := NormLp(x, y, pExp)

	c := viz.New(534, 520, -11, 11, -11, 11)

	// Axes through the origin.
	c.Path([][2]float64{{c.XMin, 0}, {c.XMax, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, c.YMin}, {0, c.YMax}}, viz.Muted, 1)

	const steps = 128
	// Three fixed reference shapes, each scaled to pass exactly through v.
	c.Path(boundary(NormL1, l1, steps), viz.Good, 1.5)
	c.Path(boundary(NormL2, l2, steps), viz.Accent, 1.5)
	c.Path(boundary(NormLInf, linf, steps), viz.Bad, 1.5)
	// The general Lp shape for the current slider, also through v; at
	// p=1/p=2 it sits exactly on top of the L1/L2 shapes above.
	c.Path(boundary(func(dx, dy float64) float64 { return NormLp(dx, dy, pExp) }, lp, steps), viz.Warm, 3)

	arrow(c, 0, 0, x, y, viz.Ink, 2.5)

	c.Text(16, 24, fmt.Sprintf("v = (%.2f, %.2f)", x, y), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("L1 = %.4f (green diamond)   L2 = %.4f (blue circle)   L∞ = %.4f (red square)",
		l1, l2, linf), 13, viz.Muted, "start")
	c.Text(16, 64, fmt.Sprintf("Lp with p=%.1f → %.4f (orange curve, always through v)", pExp, lp), 14, viz.Warm, "start")

	return c.String()
}

// arrow draws a straight line from (x0,y0) to (x1,y1) in data space, with a
// small V-shaped arrowhead at the end. Mirrors vectors's own arrow helper:
// the data range here is kept square too (see render), so a fixed
// data-space arrowhead reads correctly on both axes.
func arrow(c *viz.Canvas, x0, y0, x1, y1 float64, color string, width float64) {
	c.Path([][2]float64{{x0, y0}, {x1, y1}}, color, width)

	dx, dy := x1-x0, y1-y0
	length := math.Hypot(dx, dy)
	if length < 1e-9 {
		return
	}
	ux, uy := dx/length, dy/length
	const headLen = 0.45
	const headAngle = 0.5 // radians, ~29° off the shaft on each side

	barb := func(theta float64) (float64, float64) {
		cos, sin := math.Cos(theta), math.Sin(theta)
		bx, by := -ux, -uy // pointing back along the shaft
		return bx*cos - by*sin, bx*sin + by*cos
	}
	b1x, b1y := barb(headAngle)
	b2x, b2y := barb(-headAngle)
	c.Path([][2]float64{{x1, y1}, {x1 + headLen*b1x, y1 + headLen*b1y}}, color, width)
	c.Path([][2]float64{{x1, y1}, {x1 + headLen*b2x, y1 + headLen*b2y}}, color, width)
}
