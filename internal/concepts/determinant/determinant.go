// Package determinant visualizes a 2x2 matrix's determinant as exactly what
// it is: the signed factor by which that matrix scales area. Apply the
// matrix to the unit square and the resulting parallelogram's area is
// |det|; a negative determinant means the transformation also flips
// orientation, and a zero determinant collapses the square onto a line.
package determinant

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "determinant",
		Seq:   59,
		Title: "Determinant (area/volume scaling factor)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`matrix-multiplication` showed a matrix as a machine that sends the unit " +
						"square somewhere else — rotated, sheared, stretched. But 'somewhere " +
						"else' hides a question that composing transformations never had to " +
						"answer: did that transformation make shapes bigger, smaller, or exactly " +
						"the same size? A rotation obviously doesn't grow or shrink a shape's " +
						"area; a pure horizontal stretch obviously does. But for an arbitrary " +
						"matrix that mixes rotation, shearing, and stretching all at once, how " +
						"much does area change, exactly — is there one number, computed straight " +
						"from a matrix's four entries, that answers this without actually " +
						"transforming a shape and measuring the result?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take M = [[2,1],[0,1]] and apply it to the unit square's four corners: " +
						"(0,0)→(0,0), (1,0)→(2,0), (0,1)→(1,1), (1,1)→(3,1). The result is a " +
						"parallelogram sitting on a base of length 2 (along the x-axis, from " +
						"(0,0) to (2,0)) with a 'height' of 1 (how far the top edge sits above " +
						"the bottom) — base × height = 2×1 = 2. Compute the same thing algebraically " +
						"instead, straight from M's entries a=2, b=1, c=0, d=1: ad−bc = 2×1 − 1×0 " +
						"= 2. Same answer — this ad−bc formula is the determinant, det(M).",
					"Two more matrices show what the formula catches that eyeballing area alone " +
						"would miss. Swap i and j with M = [[0,1],[1,0]]: det = 0×0 − 1×1 = −1. " +
						"The resulting shape still has area 1, but i and j have traded places — " +
						"what used to be the counterclockwise order (i, then j) is now clockwise. " +
						"A NEGATIVE determinant means the transformation flips orientation, on top " +
						"of scaling area by its magnitude. And M = [[2,1],[4,2]] (second row is " +
						"exactly double the first): det = 2×2 − 1×4 = 0. Every corner of the unit " +
						"square lands somewhere on the single line y=2x — the square has been " +
						"flattened to a line segment with zero area, and there's no way to undo " +
						"that flattening (many different starting points now land on the same " +
						"spot), so a matrix with det=0 has no inverse.",
					"One more property, straight from `matrix-multiplication`'s own worked " +
						"example: A = 90° rotation, [[0,−1],[1,0]], has det(A) = 0×0 − (−1)×1 = 1 " +
						"(rotations never change area). B = shear k=1, [[1,1],[0,1]], has det(B) = " +
						"1×1 − 1×0 = 1. Their product C = A·B = [[0,−1],[1,1]] has det(C) = 0×1 − " +
						"(−1)×1 = 1 — and 1 = det(A) × det(B) exactly. Determinants multiply the " +
						"same way the matrices themselves combine: det(A·B) = det(A)·det(B), always.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The four sliders a, b, c, d are M's entries directly. The faint square is " +
						"the original unit square; the colored outline is M applied to it — " +
						"green when det>0 (orientation preserved), red when det<0 (orientation " +
						"flipped), grey and flattened when det≈0 (collapsed to a line). The two " +
						"green arrows are M's columns — where it sends i=(1,0) and j=(0,1) — " +
						"which are literally the two sides of that parallelogram. The readout " +
						"reports det = ad−bc alongside the shape's area computed independently " +
						"by the shoelace formula (½|Σ(x_i·y_{i+1} − x_{i+1}·y_i)|), so you can " +
						"see the two methods land on the exact same number every time.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Read off exactly how a linear transformation scales area (or, in 3D, " +
						"volume) without applying it to anything and measuring — and tell " +
						"instantly whether a matrix is invertible: det=0 exactly means the " +
						"transformation destroys information (multiple inputs collapse to the " +
						"same output), so no inverse can exist, while any nonzero determinant " +
						"guarantees one does. For a chain of transformations, det(A·B)=det(A)det(B) " +
						"means the combined area-scaling factor can be read off by multiplying " +
						"each step's determinant, without first computing the combined matrix C " +
						"the way `matrix-multiplication` had to.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A system of linear equations has a unique solution exactly when its " +
						"coefficient matrix's determinant is nonzero — det=0 is the algebraic " +
						"signature of the same 'collapsed to a line' picture above, meaning the " +
						"equations don't pin down a single answer. Map projections and other " +
						"coordinate changes use the determinant of their local linear " +
						"approximation (the Jacobian) to say exactly how much a region's area " +
						"gets stretched or shrunk at each point — why Greenland looks so much " +
						"bigger than it really is on a Mercator map. Physical simulations " +
						"(fluid flow, material deformation) track a deformation's determinant to " +
						"detect when a material has been compressed to zero volume, which signals " +
						"a breakdown in the simulation.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'det=0 means the transformation destroys information — " +
						"points that used to be spread across the whole square now sit on a " +
						"single line, so the original position can't be recovered from the " +
						"result' — anchoring the claim on what the picture actually shows " +
						"collapsing.",
					"Not like this: 'a bigger determinant always means a matrix with bigger " +
						"entries,' or the reverse. M = [[2,1],[4,2]] above has entries up to 4 but " +
						"det=0 exactly, because its rows point in the same direction; a matrix " +
						"whose rows are long and close to perpendicular can have a large " +
						"determinant with modest-looking entries. The determinant measures how " +
						"the rows/columns relate to each other geometrically, not how large any " +
						"single entry is.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "a", Min: -2.5, Max: 2.5, Step: 0.1, Def: 2},
			{Key: "b", Label: "b", Min: -2.5, Max: 2.5, Step: 0.1, Def: 1},
			{Key: "c", Label: "c", Min: -2.5, Max: 2.5, Step: 0.1, Def: 0},
			{Key: "d", Label: "d", Min: -2.5, Max: 2.5, Step: 0.1, Def: 1},
		},
		Render: render,
	})
}

// Matrix is a 2x2 matrix [[A, B], [C, D]] (row-major). Same convention as
// matrix-multiplication.Matrix, duplicated here so this package stays
// self-contained (see BUILD_CYCLE.md).
type Matrix struct {
	A, B, C, D float64
}

// Determinant is ad-bc: the signed factor by which M scales area. Its
// magnitude is how much area (or, in 3D, volume) any shape gets scaled by
// after applying M; its sign says whether M also flips orientation
// (negative) or preserves it (positive). Exactly 0 means M collapses every
// shape onto a lower-dimensional line (or point), so M has no inverse.
func Determinant(m Matrix) float64 {
	return m.A*m.D - m.B*m.C
}

// Rotation returns the matrix that rotates a vector counterclockwise by
// angleDeg degrees. Identical to matrix-multiplication.Rotation.
func Rotation(angleDeg float64) Matrix {
	rad := angleDeg * math.Pi / 180
	cos, sin := math.Cos(rad), math.Sin(rad)
	return Matrix{A: cos, B: -sin, C: sin, D: cos}
}

// Shear returns the matrix that pushes each point sideways along x by k
// times its own height. Identical to matrix-multiplication.Shear.
func Shear(k float64) Matrix {
	return Matrix{A: 1, B: k, C: 0, D: 1}
}

// Apply returns M*v for the vector (x, y). Identical to
// matrix-multiplication.Apply.
func Apply(m Matrix, x, y float64) (float64, float64) {
	return m.A*x + m.B*y, m.C*x + m.D*y
}

// Mul returns the matrix product m1*m2 -- the single matrix that applies m2
// first, then m1. Identical to matrix-multiplication.Mul.
func Mul(m1, m2 Matrix) Matrix {
	return Matrix{
		A: m1.A*m2.A + m1.B*m2.C,
		B: m1.A*m2.B + m1.B*m2.D,
		C: m1.C*m2.A + m1.D*m2.C,
		D: m1.C*m2.B + m1.D*m2.D,
	}
}

// unitSquareCorners is the unit square's four corners, in order, NOT closed
// back to the start -- the shape used for area math (ShoelaceArea assumes
// this open form and wraps around internally). See unitSquare below for the
// closed form used when drawing.
var unitSquareCorners = [][2]float64{{0, 0}, {1, 0}, {1, 1}, {0, 1}}

// unitSquare is unitSquareCorners closed back to the start, for drawing a
// complete outline with Path (which just connects consecutive points --
// without the repeated closing point, the last edge wouldn't be drawn).
var unitSquare = [][2]float64{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}

// transform applies M to every point in pts, returning the image.
func transform(m Matrix, pts [][2]float64) [][2]float64 {
	out := make([][2]float64, len(pts))
	for i, pt := range pts {
		x, y := Apply(m, pt[0], pt[1])
		out[i] = [2]float64{x, y}
	}
	return out
}

// ShoelaceArea returns the area enclosed by a simple polygon given as an
// ordered (but not explicitly closed -- the last point does not repeat the
// first) list of vertices, via the shoelace formula: half the absolute
// value of the signed sum Σ(x_i·y_{i+1} - x_{i+1}·y_i), wrapping the last
// vertex back to the first. This is a geometric area computation completely
// independent of Determinant, used to cross-check it (see the tests below).
func ShoelaceArea(pts [][2]float64) float64 {
	n := len(pts)
	if n < 3 {
		return 0
	}
	sum := 0.0
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		sum += pts[i][0]*pts[j][1] - pts[j][0]*pts[i][1]
	}
	return math.Abs(sum) / 2
}

func render(p map[string]float64) string {
	M := Matrix{A: p["a"], B: p["b"], C: p["c"], D: p["d"]}
	det := Determinant(M)

	c := viz.New(560, 560, -3.2, 3.2, -3.2, 3.2)
	c.Path([][2]float64{{c.XMin, 0}, {c.XMax, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, c.YMin}, {0, c.YMax}}, viz.Muted, 1)

	image := transform(M, unitSquare)

	color := viz.Good
	switch {
	case math.Abs(det) < 1e-9:
		color = viz.Muted
	case det < 0:
		color = viz.Bad
	}

	c.Path(unitSquare, viz.Faint, 1.5)
	c.Path(image, color, 2.5)

	// M's columns: where it sends i=(1,0) and j=(0,1) -- the parallelogram's
	// own two sides.
	arrow(c, 0, 0, M.A, M.C, viz.Good, 2.5)
	arrow(c, 0, 0, M.B, M.D, viz.Good, 2.5)

	area := ShoelaceArea(transform(M, unitSquareCorners))

	var interp string
	switch {
	case math.Abs(det) < 1e-9:
		interp = "det ≈ 0 — the square collapses onto a line: no inverse exists"
	case det > 0:
		interp = fmt.Sprintf("area scales by %.2f×, orientation preserved", math.Abs(det))
	default:
		interp = fmt.Sprintf("area scales by %.2f×, orientation FLIPPED (negative)", math.Abs(det))
	}

	c.Text(16, 24, fmt.Sprintf("M = [[%.2f,%.2f],[%.2f,%.2f]]", M.A, M.B, M.C, M.D), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("det(M) = ad−bc = %.2f×%.2f − %.2f×%.2f = %.3f", M.A, M.D, M.B, M.C, det),
		13, viz.Muted, "start")
	c.Text(16, 64, fmt.Sprintf("shoelace area of M's image = %.3f    |det(M)| = %.3f", area, math.Abs(det)),
		13, viz.Muted, "start")
	c.Text(16, 88, interp, 14, viz.Ink, "start")
	c.Text(16, 536, "faint = original unit square    colored outline = M applied to it    green arrows = M's columns",
		12, viz.Muted, "start")

	return c.String()
}

// arrow draws a straight line from (x0,y0) to (x1,y1) in data space, with a
// small V-shaped arrowhead at the end. Identical to
// matrix-multiplication.arrow.
func arrow(c *viz.Canvas, x0, y0, x1, y1 float64, color string, width float64) {
	c.Path([][2]float64{{x0, y0}, {x1, y1}}, color, width)

	dx, dy := x1-x0, y1-y0
	length := math.Hypot(dx, dy)
	if length < 1e-9 {
		return
	}
	ux, uy := dx/length, dy/length
	const headLen = 0.18
	const headAngle = 0.5 // radians, ~29 degrees off the shaft on each side

	barb := func(t float64) (float64, float64) {
		cos, sin := math.Cos(t), math.Sin(t)
		bx, by := -ux, -uy
		return bx*cos - by*sin, bx*sin + by*cos
	}
	b1x, b1y := barb(headAngle)
	b2x, b2y := barb(-headAngle)
	c.Path([][2]float64{{x1, y1}, {x1 + headLen*b1x, y1 + headLen*b1y}}, color, width)
	c.Path([][2]float64{{x1, y1}, {x1 + headLen*b2x, y1 + headLen*b2y}}, color, width)
}
