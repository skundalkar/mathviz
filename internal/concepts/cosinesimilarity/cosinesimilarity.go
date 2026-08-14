// Package cosinesimilarity visualizes cosine similarity: comparing two
// vectors by the angle between them instead of by how far apart their tips
// are, so two vectors pointing the same way score as identical no matter
// how different their lengths are — the comparison embedding search runs
// millions of times a second.
package cosinesimilarity

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "cosine-similarity",
		Seq:   33,
		Title: "Cosine similarity (comparing direction, not distance)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "ux", Label: "u_x", Min: -6, Max: 6, Step: 1, Def: 4},
			{Key: "uy", Label: "u_y", Min: -6, Max: 6, Step: 1, Def: 2},
			{Key: "vx", Label: "v_x", Min: -6, Max: 6, Step: 1, Def: 2},
			{Key: "vy", Label: "v_y", Min: -6, Max: 6, Step: 1, Def: 1},
		},
		Render: render,
	})
}

// Dot returns the dot product of u=(ux,uy) and v=(vx,vy).
func Dot(ux, uy, vx, vy float64) float64 {
	return ux*vx + uy*vy
}

// Magnitude returns a vector's length: how far its tip is from the origin.
func Magnitude(x, y float64) float64 {
	return math.Hypot(x, y)
}

// CosineSimilarity returns cos(theta), theta the angle between u and v: the
// dot product divided by the product of the two magnitudes. Dividing out
// the magnitudes is exactly what makes this "direction only" — scaling
// either vector by any positive amount leaves the result unchanged.
// Returns 0 for a zero-length vector, where direction is undefined.
func CosineSimilarity(ux, uy, vx, vy float64) float64 {
	mu, mv := Magnitude(ux, uy), Magnitude(vx, vy)
	if mu == 0 || mv == 0 {
		return 0
	}
	return Dot(ux, uy, vx, vy) / (mu * mv)
}

// AngleDegrees returns the angle between u and v in degrees, in [0, 180].
func AngleDegrees(ux, uy, vx, vy float64) float64 {
	cs := CosineSimilarity(ux, uy, vx, vy)
	// Clamp before acos: floating-point rounding can push an
	// exactly-parallel or exactly-opposite pair's cosine a hair outside
	// [-1, 1], which would otherwise make acos return NaN.
	if cs > 1 {
		cs = 1
	}
	if cs < -1 {
		cs = -1
	}
	return math.Acos(cs) * 180 / math.Pi
}

// EuclideanDistance returns the straight-line distance between u's and v's
// tips — the "how far apart" measure cosine similarity deliberately
// ignores, kept here to contrast the two directly.
func EuclideanDistance(ux, uy, vx, vy float64) float64 {
	return math.Hypot(ux-vx, uy-vy)
}

func render(p map[string]float64) string {
	ux, uy := p["ux"], p["uy"]
	vx, vy := p["vx"], p["vy"]

	dot := Dot(ux, uy, vx, vy)
	magU, magV := Magnitude(ux, uy), Magnitude(vx, vy)
	cos := CosineSimilarity(ux, uy, vx, vy)
	angle := AngleDegrees(ux, uy, vx, vy)
	dist := EuclideanDistance(ux, uy, vx, vy)

	c := viz.New(560, 520, -7, 7, -7, 7)

	// Axes through the origin.
	c.Path([][2]float64{{c.XMin, 0}, {c.XMax, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, c.YMin}, {0, c.YMax}}, viz.Muted, 1)

	// Dashed unit circle: after dividing each vector by its own length,
	// its tip lands somewhere on this circle -- direction only, magnitude
	// gone. The two small dots on it are u and v after that division.
	circle := make([][2]float64, 0, 121)
	for i := 0; i <= 120; i++ {
		t := 2 * math.Pi * float64(i) / 120
		circle = append(circle, [2]float64{math.Cos(t), math.Sin(t)})
	}
	c.Path(circle, viz.Faint, 1)

	if magU > 0 {
		px, py := c.X(ux/magU), c.Y(uy/magU)
		c.Rect(px-3, py-3, 6, 6, viz.Accent, 0.6)
	}
	if magV > 0 {
		px, py := c.X(vx/magV), c.Y(vy/magV)
		c.Rect(px-3, py-3, 6, 6, viz.Warm, 0.6)
	}

	// Arc between the two vectors, at a small fixed radius, tracing out
	// the angle theta that cosine similarity is really measuring.
	if magU > 0 && magV > 0 {
		angleU := math.Atan2(uy, ux)
		angleV := math.Atan2(vy, vx)
		delta := angleV - angleU
		for delta > math.Pi {
			delta -= 2 * math.Pi
		}
		for delta < -math.Pi {
			delta += 2 * math.Pi
		}
		const arcRadius = 1.6
		arc := make([][2]float64, 0, 61)
		for i := 0; i <= 60; i++ {
			t := angleU + delta*float64(i)/60
			arc = append(arc, [2]float64{arcRadius * math.Cos(t), arcRadius * math.Sin(t)})
		}
		c.Path(arc, viz.Good, 1.5)
	}

	arrow(c, 0, 0, ux, uy, viz.Accent, 2.5)
	arrow(c, 0, 0, vx, vy, viz.Warm, 2.5)

	c.Text(16, 22, fmt.Sprintf("u = (%.0f, %.0f)   |u| = %.2f", ux, uy, magU), 14, viz.Accent, "start")
	c.Text(16, 42, fmt.Sprintf("v = (%.0f, %.0f)   |v| = %.2f", vx, vy, magV), 14, viz.Warm, "start")
	c.Text(16, 66, fmt.Sprintf("u · v = %.2f    cos θ = %.3f    θ = %.1f°", dot, cos, angle), 14, viz.Ink, "start")
	c.Text(16, 86, fmt.Sprintf("Euclidean distance |u - v| = %.2f  (cosine similarity ignores this entirely)", dist),
		13, viz.Muted, "start")
	c.Text(16, 500, "small squares on the dashed circle = u and v after dividing out their own length",
		12, viz.Muted, "start")

	return c.String()
}

// arrow draws a straight line from (x0,y0) to (x1,y1) in data space, with
// a small V-shaped arrowhead at the end.
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
