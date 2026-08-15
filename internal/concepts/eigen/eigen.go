// Package eigen visualizes eigenvectors and eigenvalues of a 2x2 symmetric
// matrix: most directions get bent sideways by the transformation, but two
// special, perpendicular directions only get stretched (or flipped) — never
// rotated off their own line.
package eigen

import (
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

func render(p map[string]float64) string {
	c := viz.New(560, 560, -7, 7, -7, 7)
	return c.String()
}
