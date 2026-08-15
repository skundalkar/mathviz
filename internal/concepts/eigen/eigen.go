// Package eigen visualizes eigenvectors and eigenvalues of a 2x2 symmetric
// matrix: most directions get bent sideways by the transformation, but two
// special, perpendicular directions only get stretched (or flipped) — never
// rotated off their own line.
package eigen

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "eigenvectors-eigenvalues",
		Seq:   35,
		Title: "Eigenvectors & eigenvalues (directions that only stretch)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "Matrix a (x-axis stretch)", Min: 0.5, Max: 2.5, Step: 0.1, Def: 2},
			{Key: "d", Label: "Matrix d (y-axis stretch)", Min: 0.5, Max: 2.5, Step: 0.1, Def: 1},
			{Key: "b", Label: "Matrix b (shear / coupling)", Min: -1.2, Max: 1.2, Step: 0.1, Def: 0.8},
			{Key: "angle", Label: "Test vector angle θ", Min: 0, Max: 360, Step: 5, Def: 30, Unit: "°"},
		},
		Render: render,
	})
}

// Eigenvalues returns the two eigenvalues of the symmetric 2x2 matrix
// [[a, b], [b, d]], larger first (lambda1 >= lambda2). Every real symmetric
// matrix has real eigenvalues, so this never needs complex numbers: the
// discriminant ((a-d)/2)^2 + b^2 is a sum of squares and can't go negative.
func Eigenvalues(a, b, d float64) (lambda1, lambda2 float64) {
	mid := (a + d) / 2
	half := (a - d) / 2
	disc := math.Sqrt(half*half + b*b)
	return mid + disc, mid - disc
}

// EigenvectorAngle returns the direction (in radians, in [0, π)) of the
// eigenvector belonging to lambda1 — the larger eigenvalue from
// Eigenvalues — using the closed-form 2x2 symmetric-matrix formula
// theta = atan2(2b, a-d) / 2. The other eigenvector (belonging to
// lambda2) always sits exactly perpendicular to this one, at theta+π/2 —
// a property specific to symmetric matrices, proven by the spectral
// theorem, not a coincidence of the examples this package ships with.
// When a==d and b==0 the matrix is a scalar multiple of the identity
// (every direction is an eigenvector with the same eigenvalue); the
// formula still returns a defined angle (0) rather than panicking.
func EigenvectorAngle(a, b, d float64) float64 {
	return math.Atan2(2*b, a-d) / 2
}

// Apply returns A*v for the symmetric matrix [[a, b], [b, d]] and the
// vector (vx, vy).
func Apply(a, b, d, vx, vy float64) (x, y float64) {
	return a*vx + b*vy, b*vx + d*vy
}

// AngleBetweenDeg returns the angle in degrees, in [0, 180], between the
// vectors (ux,uy) and (vx,vy). Returns 0 if either vector has zero length.
func AngleBetweenDeg(ux, uy, vx, vy float64) float64 {
	mu, mv := math.Hypot(ux, uy), math.Hypot(vx, vy)
	if mu == 0 || mv == 0 {
		return 0
	}
	cos := (ux*vx + uy*vy) / (mu * mv)
	switch {
	case cos > 1:
		cos = 1
	case cos < -1:
		cos = -1
	}
	return math.Acos(cos) * 180 / math.Pi
}

// vecLen is the fixed length drawn for the test vector v -- long enough to
// read clearly against the unit circle, short enough that A*v stays
// comfortably inside the canvas across the full slider ranges.
const vecLen = 1.6

// snapDeg is how close (in degrees) the angle between v and Av has to be to
// 0 or 180 before the picture calls v an eigenvector. It's set just under
// half the angle slider's 5-degree step, so exactly landing on an
// eigenvector's own angle always counts, and its two immediate slider
// neighbors never falsely do.
const snapDeg = 2.4

func render(p map[string]float64) string {
	a, b, d := p["a"], p["b"], p["d"]
	angle := p["angle"] * math.Pi / 180

	l1, l2 := Eigenvalues(a, b, d)
	theta1 := EigenvectorAngle(a, b, d)
	theta2 := theta1 + math.Pi/2

	vx, vy := vecLen*math.Cos(angle), vecLen*math.Sin(angle)
	avx, avy := Apply(a, b, d, vx, vy)
	between := AngleBetweenDeg(vx, vy, avx, avy)
	isEigenvector := between <= snapDeg || between >= 180-snapDeg
	// Which eigenvalue is the current direction closest to? Whichever
	// eigenvector line (theta1 or theta2, each valid mod 180) the test
	// angle sits nearer to.
	lambda := l1
	if angDist(angle, theta2) < angDist(angle, theta1) {
		lambda = l2
	}

	c := viz.New(560, 560, -7, 7, -7, 7)

	// Axes through the origin.
	c.Path([][2]float64{{c.XMin, 0}, {c.XMax, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, c.YMin}, {0, c.YMax}}, viz.Muted, 1)

	// Unit circle: every possible direction for v starts on this circle.
	const circleSteps = 96
	pts := make([][2]float64, circleSteps+1)
	for i := range pts {
		t := 2 * math.Pi * float64(i) / circleSteps
		pts[i] = [2]float64{math.Cos(t), math.Sin(t)}
	}
	c.Path(pts, viz.Faint, 1)

	// The two eigenvector directions, drawn as full lines through the
	// origin (each is a whole line, not just a ray -- -v is an eigenvector
	// too, with the same eigenvalue) so their special angles are visible
	// before the slider ever reaches them.
	drawEigenLine(c, theta1)
	drawEigenLine(c, theta2)

	// v (blue) and Av (orange), or both in green when v currently is an
	// eigenvector -- Av lands exactly on v's own line, just longer/shorter
	// or flipped.
	vColor, avColor := viz.Accent, viz.Warm
	if isEigenvector {
		vColor, avColor = viz.Good, viz.Good
	}
	arrow(c, 0, 0, vx, vy, vColor, 2.5)
	arrow(c, 0, 0, avx, avy, avColor, 2.5)

	c.Text(16, 24, fmt.Sprintf("A = [[%.1f, %.1f], [%.1f, %.1f]]", a, b, b, d), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("v = (%.2f, %.2f) at θ=%.0f°", vx, vy, p["angle"]), 14, viz.Accent, "start")
	c.Text(16, 64, fmt.Sprintf("Av = (%.2f, %.2f)", avx, avy), 14, viz.Warm, "start")
	c.Text(16, 84, fmt.Sprintf("angle between v and Av: %.1f°", between), 14, viz.Ink, "start")

	if isEigenvector {
		c.Text(16, 108, fmt.Sprintf("eigenvector! Av = %.2f × v", lambda), 14, viz.Good, "start")
	} else {
		c.Text(16, 108, "not an eigenvector -- Av points off v's line", 14, viz.Muted, "start")
	}

	c.Text(16, 536, fmt.Sprintf("eigenvalues: λ₁=%.2f (θ=%.0f°)   λ₂=%.2f (θ=%.0f°)",
		l1, normDeg(theta1), l2, normDeg(theta2)), 13, viz.Muted, "start")

	return c.String()
}

// angDist returns the smallest angular distance (radians) between two
// angles, treating each as a full line rather than a ray -- so 0 and π
// count as the same direction, matching how an eigenvector's line has two
// ends that are both "the same eigenvector."
func angDist(a, b float64) float64 {
	d := math.Mod(math.Abs(a-b), math.Pi)
	if d > math.Pi/2 {
		d = math.Pi - d
	}
	return d
}

// normDeg wraps an angle in radians into [0, 180) degrees, since an
// eigenvector's direction is a line, not a ray.
func normDeg(rad float64) float64 {
	deg := math.Mod(rad*180/math.Pi, 180)
	if deg < 0 {
		deg += 180
	}
	return deg
}

// drawEigenLine draws a faint dashed line through the origin at angle theta
// (radians), extended in both directions to the edge of the canvas.
func drawEigenLine(c *viz.Canvas, theta float64) {
	const reach = 9.0
	dx, dy := math.Cos(theta), math.Sin(theta)
	c.Path([][2]float64{{-reach * dx, -reach * dy}, {reach * dx, reach * dy}}, viz.Good, 1)
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
	const headLen = 0.28
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
