// Package vectors visualizes 2D vector addition, the dot product, and
// projection — three operations built from the same two arrows drawn from
// the origin, so a change to either vector updates the sum, the dot
// product, and the projection all at once.
package vectors

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "vectors",
		Seq:   25,
		Title: "Vectors (addition, dot product & projection)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You and a friend are dragging a heavy crate, each pulling on your own rope. " +
						"You're not standing in exactly the same spot, so the two ropes pull at " +
						"slightly different angles. Gut instinct: 'I'm pulling with 5 lbs of force " +
						"and you're pulling with 5 lbs, so together we're putting 10 lbs of pull on " +
						"the crate.' That instinct is wrong unless the two ropes point in exactly " +
						"the same direction — a pull has a direction as well as a size, and the " +
						"sideways part of your pull can fight the sideways part of your friend's " +
						"instead of adding to it. Two 5-lb pulls at an angle to each other combine " +
						"into somewhere between 0 lbs (pulling directly against each other) and 10 " +
						"lbs (pulling exactly together), and 'plain addition of the numbers' only " +
						"ever gives the one special case at the very top of that range.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Represent each pull as a vector — a pair of numbers, (how far forward, how " +
						"far sideways) — instead of a single 'how hard' number:",
					"• Your pull: u = (3, 4) — 3 lbs forward, 4 lbs sideways. Its actual strength " +
						"(magnitude) is √(3²+4²) = √25 = 5 lbs, the '5 lbs' from the gut-instinct " +
						"story.",
					"• Your friend's pull, standing more directly ahead of the crate: v = (5, 0) — " +
						"all 5 lbs forward, 0 sideways.",
					"• Adding: match forward-with-forward and sideways-with-sideways separately: " +
						"u + v = (3+5, 4+0) = (8, 4). Combined strength = √(8²+4²) = √80 ≈ 8.94 lbs " +
						"— short of the naive 5+5=10 lbs, because your 4 lbs of sideways pull " +
						"doesn't point the same way your friend's pull does, so it can't fully add " +
						"to it.",
					"• How much of your pull actually helps drag the crate in your friend's " +
						"direction? Multiply matching components and add — the dot product: " +
						"u·v = (3)(5) + (4)(0) = 15.",
					"• Divide by your friend's pull's own length to turn that 15 into an actual " +
						"length measured along v's direction: 15/5 = 3 lbs. So only 3 of your 5 lbs " +
						"(60%) is doing anything useful toward your friend's direction; the " +
						"remaining 4 lbs is entirely sideways to it here and contributes 0.",
					"• The angle between the two ropes falls out for free: cos θ = (u·v)/(|u||v|) " +
						"= 15/25 = 0.6, so θ = arccos(0.6) ≈ 53.13°.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The blue arrow is your pull, u; the orange arrow is your friend's pull, v; " +
						"both start at the origin. The dashed green arrow is u+v traced tip-to-tail " +
						"— drag either vector and watch the sum arrow swing to match, always ending " +
						"at the far corner of the parallelogram the two arrows sketch out. The " +
						"dotted segment along the orange arrow marks the scalar projection of u " +
						"onto v (3 units in the example above), and the thin dashed line dropping " +
						"from the tip of u down to that point shows where the 'drop a " +
						"perpendicular' definition of projection comes from. The readout at the top " +
						"prints both vectors' components, the dot product, the angle between them, " +
						"and the projection length live as you drag.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Combine any two directional quantities — forces, velocities, displacements — " +
						"into their true net effect instead of guessing at a total; and pull out " +
						"exactly how much of one vector is 'pointing the same way' as another, which " +
						"is the question underneath comparing two directions, checking whether two " +
						"things are working with or against each other, and — once you also divide " +
						"by both lengths — measuring how similar two directions are at all, from " +
						"'exactly aligned' (cos θ = 1) to 'perpendicular, no shared direction' " +
						"(cos θ = 0) to 'directly opposed' (cos θ = −1).",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Two tug-of-war teams, a boat's engine thrust fighting a sideways current, a " +
						"plane's airspeed combined with wind to give its actual ground track, and " +
						"two ropes holding up a hammock all combine exactly like u+v here. The dot " +
						"product's 'how aligned are these two directions' shows up as physical work " +
						"(force · distance moved — a force perpendicular to the motion, like gravity " +
						"on someone walking on flat ground, does zero work), as 'cosine similarity' " +
						"for comparing two documents or two recommendation profiles, and as the " +
						"lighting calculation that shades 3D graphics (a surface facing the light " +
						"directly is bright; one at 90° to the light gets none).",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'vectors add component-by-component — forward with forward, " +
						"sideways with sideways — giving (8,4) here, whose actual length is √80 ≈ " +
						"8.94, not 5+5=10.' Not like this: adding the two vectors' plain lengths " +
						"(5+5=10) as if direction didn't matter, or treating the dot product's " +
						"result (a single number, 15 here) as if it were itself a vector or a length " +
						"— it only becomes a length, in the direction of v, after dividing by |v| " +
						"(15/5=3). A negative dot product doesn't mean a 'negative length' either; " +
						"it just means the angle between the two vectors is more than 90°, i.e., " +
						"they point more against each other than with each other.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "ux", Label: "Vector u — x", Min: -5, Max: 5, Step: 0.5, Def: 3},
			{Key: "uy", Label: "Vector u — y", Min: -5, Max: 5, Step: 0.5, Def: 4},
			{Key: "vx", Label: "Vector v — x", Min: -5, Max: 5, Step: 0.5, Def: 5},
			{Key: "vy", Label: "Vector v — y", Min: -5, Max: 5, Step: 0.5, Def: 0},
		},
		Render: render,
	})
}

// Add returns the component-wise sum of the vectors (ux,uy) and (vx,vy).
func Add(ux, uy, vx, vy float64) (x, y float64) {
	return ux + vx, uy + vy
}

// Dot returns the dot product of the vectors (ux,uy) and (vx,vy): the sum
// of the products of their matching components.
func Dot(ux, uy, vx, vy float64) float64 {
	return ux*vx + uy*vy
}

// Magnitude returns the length of the vector (x,y).
func Magnitude(x, y float64) float64 {
	return math.Hypot(x, y)
}

// AngleBetween returns the angle in radians, in [0, π], between the
// vectors (ux,uy) and (vx,vy). Returns 0 if either vector has zero length
// (a zero vector has no direction to measure an angle against).
func AngleBetween(ux, uy, vx, vy float64) float64 {
	mu, mv := Magnitude(ux, uy), Magnitude(vx, vy)
	if mu == 0 || mv == 0 {
		return 0
	}
	cos := Dot(ux, uy, vx, vy) / (mu * mv)
	// Clamp: floating-point rounding can push |cos| a hair past 1, which
	// would make Acos return NaN.
	switch {
	case cos > 1:
		cos = 1
	case cos < -1:
		cos = -1
	}
	return math.Acos(cos)
}

// ScalarProjection returns the signed length of the component of (ux,uy)
// that lies along the direction of (vx,vy) — how far u reaches in v's
// direction. Negative means u leans more against v's direction than with
// it. Returns 0 if v is the zero vector (no direction to project onto).
func ScalarProjection(ux, uy, vx, vy float64) float64 {
	mv := Magnitude(vx, vy)
	if mv == 0 {
		return 0
	}
	return Dot(ux, uy, vx, vy) / mv
}

// VectorProjection returns the component of (ux,uy) that lies along the
// vector (vx,vy), as a vector pointing along v (or its opposite, when the
// scalar projection is negative). Returns (0,0) if v is the zero vector.
func VectorProjection(ux, uy, vx, vy float64) (x, y float64) {
	vv := Dot(vx, vy, vx, vy)
	if vv == 0 {
		return 0, 0
	}
	scale := Dot(ux, uy, vx, vy) / vv
	return scale * vx, scale * vy
}

func render(p map[string]float64) string {
	ux, uy, vx, vy := p["ux"], p["uy"], p["vx"], p["vy"]
	sx, sy := Add(ux, uy, vx, vy)
	dot := Dot(ux, uy, vx, vy)
	mu, mv, msum := Magnitude(ux, uy), Magnitude(vx, vy), Magnitude(sx, sy)
	proj := ScalarProjection(ux, uy, vx, vy)
	projX, projY := VectorProjection(ux, uy, vx, vy)
	angleDeg := AngleBetween(ux, uy, vx, vy) * 180 / math.Pi

	// Square-ish data range: viz.New's fixed padding (48/18 left/right,
	// 18/34 top/bottom) leaves an equal-ish plot area on a 534x520 canvas,
	// so a symmetric data range renders angles without visible skew.
	c := viz.New(534, 520, -11, 11, -11, 11)

	// Axes through the origin.
	c.Path([][2]float64{{c.XMin, 0}, {c.XMax, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, c.YMin}, {0, c.YMax}}, viz.Muted, 1)

	// Parallelogram construction: dashed guides from each vector's tip to
	// the sum, parallel to the other vector, showing where u+v lands.
	c.Path([][2]float64{{ux, uy}, {sx, sy}}, viz.Faint, 1)
	c.Path([][2]float64{{vx, vy}, {sx, sy}}, viz.Faint, 1)

	// Projection: the dotted segment from the origin to the projection
	// point along v, and the perpendicular dropped from the tip of u.
	c.Path([][2]float64{{0, 0}, {projX, projY}}, viz.Good, 3)
	c.Path([][2]float64{{ux, uy}, {projX, projY}}, viz.Muted, 1)

	arrow(c, 0, 0, vx, vy, viz.Warm, 2.5)
	arrow(c, 0, 0, ux, uy, viz.Accent, 2.5)
	arrow(c, 0, 0, sx, sy, viz.Ink, 2)

	c.Text(16, 24, fmt.Sprintf("u = (%.1f, %.1f)  |u| = %.2f", ux, uy, mu), 14, viz.Accent, "start")
	c.Text(16, 44, fmt.Sprintf("v = (%.1f, %.1f)  |v| = %.2f", vx, vy, mv), 14, viz.Warm, "start")
	c.Text(16, 64, fmt.Sprintf("u + v = (%.1f, %.1f)  |u+v| = %.2f", sx, sy, msum), 14, viz.Ink, "start")
	c.Text(16, 84, fmt.Sprintf("u · v = %.2f    angle = %.1f°    proj of u onto v = %.2f", dot, angleDeg, proj),
		13, viz.Good, "start")

	return c.String()
}

// arrow draws a straight line from (x0,y0) to (x1,y1) in data space, with
// a small V-shaped arrowhead at the end. The data range is kept roughly
// square (see render), so a fixed data-space arrowhead reads correctly on
// both axes without needing pixel-space compensation.
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
