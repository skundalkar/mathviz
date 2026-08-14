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
				Body: []string{
					"You've built a document search engine: type a query, get back the most " +
						"similar documents. Represent each document as a vector of word counts — " +
						"the same idea the naive-bayes concept elsewhere in this gallery uses for " +
						"'does this email contain word X,' just counted instead of yes/no. 'Moby " +
						"Dick' mentions 'whale' 900 times; a one-paragraph book review mentions " +
						"'whale' 3 times. Gut instinct: 'just compare the raw counts directly — " +
						"closer counts should mean more related documents.' That instinct actively " +
						"works against you: measured by ordinary straight-line distance between " +
						"word-count vectors, the short review will always look 'far' from the long " +
						"novel even though they're about exactly the same thing, while two " +
						"unrelated documents that happen to be similar lengths can look deceptively " +
						"close. Length alone ends up dominating the comparison. Is there a way to " +
						"compare two vectors that captures 'are these about the same thing' without " +
						"being thrown off by 'how long are they'?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Compare directions instead of raw positions. Take two word-count vectors " +
						"tracking just two words, 'cat' and 'dog': document A, u = (4, 2), and a " +
						"much shorter document B that happens to use the same two words in the " +
						"same ratio, v = (2, 1) — exactly half of u.",
					"• Dot product: u·v = (4)(2) + (2)(1) = 8 + 2 = 10 — multiply matching " +
						"components and add, same operation the vectors concept elsewhere in this " +
						"gallery uses.",
					"• Magnitudes (lengths): |u| = √(4²+2²) = √20 ≈ 4.47, |v| = √(2²+1²) = √5 ≈ " +
						"2.24.",
					"• Cosine similarity = (u·v)/(|u|·|v|) = 10/(4.47×2.24) = 10/10.00 = 1.00 — " +
						"dividing out both magnitudes cancels the length difference completely, " +
						"leaving only 'do these two vectors point the same way.' They do, exactly, " +
						"so the score hits its maximum, 1.00, even though B is half as long as A.",
					"• Contrast with Euclidean (straight-line) distance between the same two " +
						"points: √((4-2)²+(2-1)²) = √5 ≈ 2.24 — very much not zero. Distance says " +
						"'these are apart'; cosine similarity says 'these point the same way.' Both " +
						"are true at once, and for 'is this document about the same topic,' the " +
						"direction is the one that matters.",
					"• A third document, C = (1, 5) — mostly 'dog,' barely any 'cat' — gives u·C = " +
						"4+10 = 14, |C| = √26 ≈ 5.10, cosine similarity = 14/(4.47×5.10) ≈ 0.61, " +
						"angle ≈ 52° — related, but noticeably less so than A and B.",
					"• A fourth, D = (2, -4), gives u·D = 8-8 = 0 — cosine similarity exactly 0, a " +
						"90° angle: orthogonal, the vector-space way of saying 'no direction in " +
						"common at all.'",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The blue arrow is u, the orange arrow is v, both drawn from the origin — " +
						"drag either vector's sliders and its arrow moves. The green arc traces θ, " +
						"the angle between them, and the numbers above report u·v, cos θ, and θ " +
						"itself, live. The dashed circle has radius 1; the small blue and orange " +
						"squares sitting on it are u and v after dividing each by its own length — " +
						"'direction only' made literal, since every vector lands on the same circle " +
						"regardless of how long it started. With the default u=(4,2), v=(2,1), the " +
						"two squares land in exactly the same spot on the circle (cos θ = 1.00) " +
						"even though the arrows themselves are clearly different lengths — and the " +
						"gray 'Euclidean distance' line underneath stays a stubborn 2.24, proof " +
						"that the two measures are answering genuinely different questions. Scale " +
						"v further out along the same direction (say to (4,2) itself, or (8,4)) and " +
						"watch cos θ hold at 1.00 the whole time while the distance number keeps " +
						"changing.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Rank documents (or, in a modern system, embeddings — vectors a neural " +
						"network produces to stand in for meaning) by topic or meaning, with a " +
						"short precise match able to outrank a long meandering one, instead of " +
						"length quietly deciding the ranking for you. This is exactly the " +
						"comparison a semantic search or recommendation engine runs, millions of " +
						"times per query, between a query vector and every candidate in its index.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Search engines and recommendation systems ranking results by relevance " +
						"rather than raw overlap; embedding-based semantic search — the retrieval " +
						"step behind modern AI assistants that look up relevant documents before " +
						"answering — which represents text, images, or audio as vectors and finds " +
						"'similar' ones by exactly this comparison; duplicate and near-duplicate " +
						"detection (two versions of the same document, lightly edited, still point " +
						"the same way); and 'people who liked this also liked...' recommendations, " +
						"comparing users' preference vectors the same way. Outside tech, everyday " +
						"'similar' usually folds in size too ('a similar-sized crowd') — cosine " +
						"similarity deliberately does not; two things can be maximally 'similar' by " +
						"this measure while being wildly different in scale.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'cosine similarity is the cosine of the angle between two " +
						"vectors — it's blind to their lengths on purpose.' Not like this: treating " +
						"it like a distance, where smaller means more similar — cosine similarity " +
						"runs the other way, with 1.00 (0°, pointing the same way) as the most " +
						"similar and -1.00 (180°, pointing exactly opposite) as the least; or " +
						"assuming a score of 0 always means 'completely unrelated' in a real-world " +
						"sense — it precisely means 'perpendicular in whatever feature space these " +
						"vectors were built from,' which for a two-word count vector like this " +
						"example is a fairly literal 'shares nothing,' but in a 300-dimension " +
						"embedding space is a subtler statement about the geometry a model learned, " +
						"not a plain-English verdict; and don't assume the score is bounded 0 to 1 " +
						"the way a probability is — it legitimately runs from -1 to 1.",
				},
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
