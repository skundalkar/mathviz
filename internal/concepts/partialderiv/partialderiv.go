// Package partialderiv visualizes partial derivatives and the gradient:
// freezing all but one variable of f(x,y) reduces it back to the familiar
// single-variable slope from `derivative`, and packaging the two resulting
// slopes into one vector, the gradient, points in the direction f increases
// fastest — with its length telling you exactly how steep that direction is.
package partialderiv

import (
	"fmt"
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

// contourLevels are the f=k circles drawn faintly for scale: since
// F(x,y)=x²+y², the level set f=k is the circle of radius √k.
var contourLevels = []float64{1, 4, 9, 16}

// alignTolDeg is how close theta must sit to the gradient's own angle
// before the readout calls the two "ALIGNED" -- comfortably tighter than
// the 1° slider step so a genuine match is never missed by rounding.
const alignTolDeg = 2.0

func render(p map[string]float64) string {
	x0, y0, theta := p["x0"], p["y0"], p["theta"]
	fx, fy := Gradient(x0, y0)
	mag := GradientMagnitude(x0, y0)
	gradAngle := GradientAngleDeg(x0, y0)
	dDeriv := DirectionalDerivative(x0, y0, theta)

	c := viz.New(560, 560, -4.5, 4.5, -4.5, 4.5)
	c.Path([][2]float64{{c.XMin, 0}, {c.XMax, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, c.YMin}, {0, c.YMax}}, viz.Muted, 1)

	// Contours of f: since f=x²+y², each level set f=k is a plain circle
	// of radius √k, so no gradient-free implicit-curve sampling is needed.
	for _, k := range contourLevels {
		r := math.Sqrt(k)
		pts := make([][2]float64, 0, 73)
		for i := 0; i <= 72; i++ {
			rad := float64(i) * 5 * math.Pi / 180
			pts = append(pts, [2]float64{r * math.Cos(rad), r * math.Sin(rad)})
		}
		c.Path(pts, viz.Muted, 1)
	}

	// Dashed cross-section guides: the horizontal one at y=y0 is where
	// PartialX's slice g(x)=F(x,y0) lives, the vertical one at x=x0 is
	// where PartialY's slice h(y)=F(x0,y) lives.
	c.Path([][2]float64{{c.XMin, y0}, {c.XMax, y0}}, viz.Faint, 1)
	c.Path([][2]float64{{x0, c.YMin}, {x0, c.YMax}}, viz.Faint, 1)

	px, py := c.X(x0), c.Y(y0)
	c.Rect(px-3, py-3, 6, 6, viz.Ink, 1)

	// Both arrows are drawn at a fixed visual length from the point so
	// their *directions* -- the thing being compared -- read clearly no
	// matter how their true magnitudes differ; the readout prints the
	// real numbers.
	const armLen = 1.2
	if mag > 1e-9 {
		arrow(c, x0, y0, x0+fx/mag*armLen, y0+fy/mag*armLen, viz.Accent, 2.5)
	}
	rad := theta * math.Pi / 180
	ux, uy := math.Cos(rad), math.Sin(rad)
	arrow(c, x0, y0, x0+ux*armLen, y0+uy*armLen, viz.Warm, 2.5)

	c.Text(16, 24, fmt.Sprintf("point = (%.2f, %.2f)    f(x,y) = x²+y² = %.2f", x0, y0, F(x0, y0)), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("∂f/∂x = %.2f    ∂f/∂y = %.2f    grad f = (%.2f, %.2f) [blue]", fx, fy, fx, fy),
		13, viz.Muted, "start")
	c.Text(16, 64, fmt.Sprintf("|grad f| = %.3f at %.1f°    direction u [orange] = %.1f°    D_u f = %.3f",
		mag, gradAngle, theta, dDeriv), 13, viz.Ink, "start")

	diff := math.Abs(theta - gradAngle)
	if diff > 180 {
		diff = 360 - diff
	}
	status := "not steepest"
	if diff < alignTolDeg {
		status = "ALIGNED -- steepest ascent"
	}
	c.Text(16, 84, status, 13, viz.Good, "start")
	c.Text(16, 536, "faint circles = contours of f    dashed lines = the x0,y0 slices the partials use",
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
	const headAngle = 0.5 // radians, ~29° off the shaft on each side

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
