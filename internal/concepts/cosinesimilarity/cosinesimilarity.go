// Package cosinesimilarity visualizes cosine similarity: comparing two
// vectors by the angle between them instead of by how far apart their tips
// are, so two vectors pointing the same way score as identical no matter
// how different their lengths are — the comparison embedding search runs
// millions of times a second.
package cosinesimilarity

import (
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
	c := viz.New(560, 480, -7, 7, -7, 7)
	return c.String()
}
