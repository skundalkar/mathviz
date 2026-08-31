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
				Body: []string{
					"`derivative` found the exact slope of a curve at one point — but there " +
						"was only ever one direction to move in, along the x-axis. Standing on a " +
						"hillside, that's not true anymore: you could head north, east, or any " +
						"compass bearing in between, and the ground gets steeper or flatter " +
						"depending on which way you pick. If f(x,y) reports elevation at map " +
						"coordinates (x,y), which single direction should you walk to gain " +
						"elevation the fastest — and how much do you actually gain per step going " +
						"that way?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take f(x,y) = x²+y² and the point (1,2), where f(1,2) = 1+4 = 5:",
					"• Freeze y at 2 and let only x vary: g(x) = f(x,2) = x²+4. That's a plain " +
						"single-variable curve again, so `derivative` applies directly: " +
						"g'(x) = 2x, so g'(1) = 2. This is the partial derivative of f with " +
						"respect to x at (1,2), written ∂f/∂x(1,2) = 2 — moving east (increasing " +
						"x, holding y=2) raises f by about 2 units per unit of x.",
					"• Freeze x at 1 and let only y vary: h(y) = f(1,y) = 1+y². h'(y) = 2y, so " +
						"h'(2) = 4: ∂f/∂y(1,2) = 4 — moving north raises f by about 4 per unit of " +
						"y, steeper than east.",
					"• Package the two slopes into one vector, the gradient: " +
						"∇f(1,2) = (∂f/∂x, ∂f/∂y) = (2,4). Its length, |∇f| = √(2²+4²) = √20 ≈ " +
						"4.47, isn't an arbitrary combination of the two numbers — it's about to " +
						"turn out to be the actual steepest slope available at this point, in any " +
						"direction at all, not just east or north.",
					"• Check that claim with `vectors`' dot product: the rate f changes moving " +
						"along any unit direction u = (cos θ, sin θ) is the directional " +
						"derivative D_u f = ∇f · u. Straight east (θ=0°, u=(1,0)): " +
						"D_u f = 2(1)+4(0) = 2, matching ∂f/∂x exactly, as it must — 'freeze y' " +
						"is just 'move along the x-axis' restated. Straight along the gradient's " +
						"own direction (θ = atan2(4,2) ≈ 63.43°, u ≈ (0.447,0.894)): " +
						"D_u f = 2(0.447)+4(0.894) ≈ 0.894+3.578 ≈ 4.47 — exactly |∇f|, the " +
						"biggest a directional derivative can get here, because a dot product " +
						"u·v is maximized exactly when u points the same way as v.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"x0 and y0 drag the point around; theta swings a direction arrow (orange) " +
						"around that point. The faint circles are contours of f = x²+y² (level " +
						"sets, always circles here since f only depends on distance from the " +
						"origin); the dashed lines mark the y=y0 and x=x0 slices ∂f/∂x and ∂f/∂y " +
						"come from. The blue arrow is the gradient ∇f, fixed at the point's own " +
						"steepest-ascent direction. Drag theta until the orange arrow lines up " +
						"with the blue one and the readout's D_u f climbs to match |∇f| exactly, " +
						"flipping the status line to 'ALIGNED' — every other angle gives a " +
						"smaller (or, pointing the other way, negative) number.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Name a single best direction to move from any point of a multivariable " +
						"function, instead of only ever being able to ask 'how steep is it in " +
						"this one direction I already picked' — and know exactly how much you'd " +
						"gain moving that way. It's also the compass `gradient-descent` follows: " +
						"stepping opposite the gradient walks straight downhill by the steepest " +
						"available route at every point along the way.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A hiking app can report the grade in the direction of the trail you're on, " +
						"but the steepest slope at your exact GPS pin might point off into the " +
						"woods on a completely different bearing — the gradient is that bearing. " +
						"Weather maps' pressure-gradient force pushes air straight down the " +
						"steepest drop in pressure, not along whatever direction the wind " +
						"happened to be blowing before. Heat flows down a temperature gradient " +
						"the same way. And `gradient-descent`, already in this gallery, is " +
						"exactly this idea applied to training a model: repeatedly step opposite " +
						"∇(loss) to reduce the loss as fast as locally possible.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: '∂f/∂x only measures the slope along one axis with " +
						"everything else frozen; the gradient combines both partials into the " +
						"one vector that points the actual steepest way up.'",
					"Not like this: treating a single partial derivative, like ∂f/∂x(1,2)=2, as " +
						"'the' slope at that point — it's only the slope in that one axis-aligned " +
						"direction, and a steeper direction exists unless the gradient happens to " +
						"point exactly along that axis (compare: 4.47 available along the " +
						"gradient here, versus only 2 heading due east). Also not like this: " +
						"reading the gradient's components as literal step sizes to walk — they " +
						"give a direction and a rate (units of f gained per unit of distance " +
						"moved), not a distance to travel.",
				},
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
