// Package matmul visualizes matrix multiplication as composing two linear
// transformations into one: apply B to the unit square, then A to the
// result, and see that it lands in exactly the same place as applying the
// single combined matrix C directly. It also reads a matrix by its columns
// — where it sends the basis vectors i=(1,0) and j=(0,1) — and shows that
// applying A after B is generally different from applying B after A.
package matmul

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "matrix-multiplication",
		Seq:   51,
		Title: "Matrix multiplication (composing transformations)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`eigenvectors-eigenvalues` already showed one matrix bending a vector — " +
						"multiply v by A and it lands somewhere new. Now say a shape (a video-game " +
						"model, an image, a robot arm's joints) needs two transformations in a row: " +
						"rotate it, then shear it. The straightforward way is to transform every " +
						"one of its thousand points twice — once through the rotation matrix, then " +
						"the result through the shear matrix — two matrix-vector multiplies per " +
						"point, two thousand total, every single frame the shape moves. Is there a " +
						"way to boil 'rotate, then shear' down into a single matrix, computed once, " +
						"so transforming each point afterward costs one multiply instead of two?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take A, a 90° rotation, and B, a shear that pushes each point sideways by " +
						"its own height: A = [[0,-1],[1,0]], B = [[1,1],[0,1]]. Apply B then A to " +
						"the point v=(1,0), one step at a time:",
					"• Bv = (1·1+1·0, 0·1+1·0) = (1, 0) — the shear leaves a point on the x-axis " +
						"untouched (it has height 0, so it doesn't get pushed sideways at all).",
					"• A(Bv) = A(1,0) = (0·1+(-1)·0, 1·1+0·0) = (0, 1) — the rotation then swings " +
						"that point 90°.",
					"Now compute the single combined matrix C = A·B directly, using the row-times-" +
						"column rule (each output entry is a row of A dotted with a column of B — " +
						"the same dot product from `vectors`, applied four times): C = [[0·1+(-1)·0, " +
						"0·1+(-1)·1], [1·1+0·0, 1·1+0·1]] = [[0,-1], [1,1]]. Apply C to v=(1,0) " +
						"directly: (0·1+(-1)·0, 1·1+1·0) = (0, 1) — exactly matching A(Bv) above, " +
						"in one multiply instead of two.",
					"Now swap the order and compute B·A instead: B·A = [[1·0+1·1, 1·(-1)+1·0], " +
						"[0·0+1·1, 0·(-1)+1·0]] = [[1,-1], [1,0]]. Apply that to the same v=(1,0): " +
						"(1·1+(-1)·0, 1·1+0·0) = (1, 1) — a different point than C's (0,1). A·B and " +
						"B·A are, in general, two different matrices: order matters.",
					"One more reading of C = [[0,-1],[1,1]]: its first column, (0,1), is exactly " +
						"where C sends i=(1,0) — matches the (0,1) computed above. Its second " +
						"column, (-1,1), is exactly where C sends j=(0,1): C·j = (0·0+(-1)·1, " +
						"1·0+1·1) = (-1,1). A matrix's columns aren't just numbers in a grid — " +
						"they're literally the images of the basis vectors under that " +
						"transformation.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"angleA sets A, a pure rotation; shearB sets B, a shear along x. order picks " +
						"which composition C represents: 0 for A after B (C = A·B, shear first, " +
						"then rotate), 1 for B after A (C = B·A, rotate first, then shear). The " +
						"faint square is the original unit square; the blue square is it after the " +
						"first transformation alone; the orange square is C applied directly to " +
						"the original — and it always lands exactly where chaining the two steps " +
						"by hand would put it, because that's what C is built to do. The two green " +
						"arrows are C's columns: where i and j end up, read straight off the matrix " +
						"printed in the readout. Flip order at the same angleA/shearB and watch the " +
						"orange square land somewhere different — the same two building blocks, " +
						"combined in the other order, is a genuinely different transformation.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Precompute one matrix that captures an entire chain of transformations, then " +
						"apply that single matrix to every point cheaply instead of re-running the " +
						"whole chain per point — and chain arbitrarily many steps the same way, " +
						"since matrix products combine just as C = A·B did above (D·C = D·(A·B) " +
						"folds a third transformation in without touching how A and B were built).",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A game engine multiplies together a model's rotation, scale, and its camera's " +
						"view and projection matrices into one combined matrix per object, once per " +
						"frame, instead of chaining four separate transforms per vertex; a robot " +
						"arm's forward kinematics multiplies each joint's local rotation matrix " +
						"together, base joint to fingertip, to get one matrix for where the hand " +
						"ends up; and a map app rotating and zooming the view at once is applying " +
						"one matrix that's really 'rotate' and 'zoom' multiplied together ahead of " +
						"time.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'A·B means apply B first, then A' — matrix products read " +
						"right-to-left when you think of them as acting on a vector, the opposite " +
						"of how the symbols sit on the page. Not like this: assuming A·B and B·A " +
						"always give the same answer because ordinary number multiplication " +
						"(3×5 = 5×3) does — the worked example above shows A·B sending (1,0) to " +
						"(0,1) while B·A sends the very same point to (1,1); matrix multiplication " +
						"only commutes for special pairs (like two pure rotations), never as a " +
						"general rule.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "angleA", Label: "A: rotate by", Min: -180, Max: 180, Step: 5, Def: 90, Unit: "°"},
			{Key: "shearB", Label: "B: shear by", Min: -1.5, Max: 1.5, Step: 0.1, Def: 1},
			{Key: "order", Label: "Order (0 = A after B, 1 = B after A)", Min: 0, Max: 1, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// Matrix is a 2x2 matrix [[A, B], [C, D]] (row-major).
type Matrix struct {
	A, B, C, D float64
}

// Rotation returns the matrix that rotates a vector counterclockwise by
// angleDeg degrees.
func Rotation(angleDeg float64) Matrix {
	rad := angleDeg * math.Pi / 180
	cos, sin := math.Cos(rad), math.Sin(rad)
	return Matrix{A: cos, B: -sin, C: sin, D: cos}
}

// Shear returns the matrix that pushes each point sideways along x by k
// times its own height (y stays fixed): (x, y) -> (x+k*y, y).
func Shear(k float64) Matrix {
	return Matrix{A: 1, B: k, C: 0, D: 1}
}

// Apply returns M*v for the vector (x, y).
func Apply(m Matrix, x, y float64) (float64, float64) {
	return m.A*x + m.B*y, m.C*x + m.D*y
}

// Mul returns the matrix product m1*m2 — the single matrix that applies m2
// first, then m1. For any vector v, Apply(Mul(m1, m2), v) equals
// Apply(m1, Apply(m2, v)) exactly; Mul is generally NOT commutative
// (Mul(m1, m2) != Mul(m2, m1)).
func Mul(m1, m2 Matrix) Matrix {
	return Matrix{
		A: m1.A*m2.A + m1.B*m2.C,
		B: m1.A*m2.B + m1.B*m2.D,
		C: m1.C*m2.A + m1.D*m2.C,
		D: m1.C*m2.B + m1.D*m2.D,
	}
}

// unitSquare is the shape transformed on the canvas: the four corners of the
// unit square, closed back to the start so Path draws a complete outline.
var unitSquare = [][2]float64{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}

func transform(m Matrix, pts [][2]float64) [][2]float64 {
	out := make([][2]float64, len(pts))
	for i, pt := range pts {
		x, y := Apply(m, pt[0], pt[1])
		out[i] = [2]float64{x, y}
	}
	return out
}

func render(p map[string]float64) string {
	angleA, shearB, order := p["angleA"], p["shearB"], p["order"]

	A := Rotation(angleA)
	B := Shear(shearB)

	// first is whichever matrix is applied first; second is applied to
	// first's result. C is the single combined matrix — Mul(second, first)
	// either way, since Mul(m1,m2) applies m2 first.
	var first, second Matrix
	var label string
	if order < 0.5 {
		first, second, label = B, A, "C = A·B  (B first, then A)"
	} else {
		first, second, label = A, B, "C = B·A  (A first, then B)"
	}
	C := Mul(second, first)

	c := viz.New(560, 560, -3.2, 3.2, -3.2, 3.2)
	c.Path([][2]float64{{c.XMin, 0}, {c.XMax, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, c.YMin}, {0, c.YMax}}, viz.Muted, 1)

	afterFirst := transform(first, unitSquare)
	afterC := transform(C, unitSquare) // == second applied to afterFirst, pointwise

	c.Path(unitSquare, viz.Faint, 1.5)
	c.Path(afterFirst, viz.Accent, 2)
	c.Path(afterC, viz.Warm, 2.5)

	// C's columns: where it sends i=(1,0) and j=(0,1).
	arrow(c, 0, 0, C.A, C.C, viz.Good, 2.5)
	arrow(c, 0, 0, C.B, C.D, viz.Good, 2.5)

	c.Text(16, 24, label, 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("A = [[%.2f,%.2f],[%.2f,%.2f]] (rotate %.0f°)",
		A.A, A.B, A.C, A.D, angleA), 13, viz.Muted, "start")
	c.Text(16, 62, fmt.Sprintf("B = [[%.2f,%.2f],[%.2f,%.2f]] (shear %.1f)",
		B.A, B.B, B.C, B.D, shearB), 13, viz.Muted, "start")
	c.Text(16, 86, fmt.Sprintf("C = [[%.2f,%.2f],[%.2f,%.2f]]", C.A, C.B, C.C, C.D), 14, viz.Ink, "start")
	c.Text(16, 106, fmt.Sprintf("C·i = (%.2f,%.2f)   C·j = (%.2f,%.2f)  (green arrows = C's columns)",
		C.A, C.C, C.B, C.D), 13, viz.Muted, "start")
	c.Text(16, 536, "faint = original square   blue = after first step   orange = C applied directly",
		12, viz.Muted, "start")

	return c.String()
}

// arrow draws a straight line from (x0,y0) to (x1,y1) in data space, with a
// small V-shaped arrowhead at the end.
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
