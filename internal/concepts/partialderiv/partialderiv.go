// Package partialderiv visualizes partial derivatives and the gradient:
// freezing all but one variable of f(x,y) reduces it back to the familiar
// single-variable slope from `derivative`, and packaging the two resulting
// slopes into one vector, the gradient, points in the direction f increases
// fastest — with its length telling you exactly how steep that direction is.
package partialderiv

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "partial-derivatives-gradient",
		Seq:   82,
		Title: "Partial derivatives & the gradient",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "x0", Label: "Point x0", Min: -3, Max: 3, Step: 0.25, Def: 1},
			{Key: "y0", Label: "Point y0", Min: -3, Max: 3, Step: 0.25, Def: 2},
			{Key: "theta", Label: "Direction theta", Min: 0, Max: 360, Step: 1, Def: 0, Unit: "°"},
		},
		Render: render,
	})
}

// F is the sample function this concept visualizes, f(x,y) = x² + y², a
// bowl shape whose contours (level sets) are circles centered on the
// origin. Pure in, pure out — no globals, no time, no randomness.
func F(x, y float64) float64 {
	return x*x + y*y
}

// PartialX is ∂f/∂x at (x,y): freeze y and take the ordinary single-
// variable derivative of g(x) = F(x,y) with respect to x, the same way
// `derivative` differentiates x². For F it happens not to depend on y.
func PartialX(x, y float64) float64 {
	return 2 * x
}

// PartialY is ∂f/∂y at (x,y): freeze x and differentiate h(y) = F(x,y)
// with respect to y.
func PartialY(x, y float64) float64 {
	return 2 * y
}

// Gradient returns ∇f(x,y) = (∂f/∂x, ∂f/∂y), the vector that packages both
// partial derivatives together. It points in the direction f increases
// fastest from (x,y), and its length is how fast f grows moving that way.
func Gradient(x, y float64) (fx, fy float64) {
	return PartialX(x, y), PartialY(x, y)
}

// GradientMagnitude returns |∇f(x,y)|, the steepest possible rate of
// increase of f at (x,y) in any direction.
func GradientMagnitude(x, y float64) float64 {
	fx, fy := Gradient(x, y)
	return math.Hypot(fx, fy)
}

// GradientAngleDeg returns the compass angle, in degrees measured counter-
// clockwise from the positive x-axis, that ∇f(x,y) points along. Returns 0
// at a point where the gradient is the zero vector (no direction to point).
func GradientAngleDeg(x, y float64) float64 {
	fx, fy := Gradient(x, y)
	if fx == 0 && fy == 0 {
		return 0
	}
	return math.Atan2(fy, fx) * 180 / math.Pi
}

// DirectionalDerivative returns D_u f(x,y), the rate f changes at (x,y)
// moving along the unit direction thetaDeg (measured the same way as
// GradientAngleDeg) -- the dot product (see `vectors`) of ∇f(x,y) with
// that direction's unit vector. It reduces to PartialX at theta=0° and to
// PartialY at theta=90°, and it is maximized, at exactly
// GradientMagnitude(x,y), when theta matches GradientAngleDeg(x,y).
func DirectionalDerivative(x, y, thetaDeg float64) float64 {
	fx, fy := Gradient(x, y)
	rad := thetaDeg * math.Pi / 180
	return fx*math.Cos(rad) + fy*math.Sin(rad)
}

func render(p map[string]float64) string {
	c := viz.New(560, 560, -4.5, 4.5, -4.5, 4.5)
	c.Axes()
	return c.String()
}
