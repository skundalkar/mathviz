// Package cosinesimilarity visualizes cosine similarity: comparing two
// vectors by the angle between them instead of by how far apart their tips
// are, so two vectors pointing the same way score as identical no matter
// how different their lengths are — e.g. two movie fans with identical
// taste who simply rate everything at different volumes.
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
					"A movie app wants to recommend 'users with taste like yours also loved...' " +
						"— which means it first has to answer 'which two users actually have " +
						"similar taste?' Represent each user as a vector: how much they tend to " +
						"enjoy Action movies, how much they tend to enjoy Romance movies, on a " +
						"scale from -5 (can't stand it) to +5 (loves it). Priya rates Action +4 " +
						"and Romance +2; Sam rates Action +2 and Romance +1 — exactly half of " +
						"Priya's numbers, straight across, because Sam is simply a more reserved " +
						"rater who never scores anything as high as Priya does. Gut instinct: " +
						"'just measure how far apart their two rating pairs are.' That instinct " +
						"actively misleads here — by that raw distance, Priya and Sam's ratings " +
						"look meaningfully apart, even though anyone glancing at both profiles can " +
						"tell their taste is identical; Sam just rates everything lower, across " +
						"the board. Is there a way to compare two people's ratings that captures " +
						"'do they like the same things, in the same proportion' without being " +
						"thrown off by 'is one of them just a harsher or more enthusiastic rater " +
						"overall'?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Compare directions instead of raw positions. Priya's ratings as a vector, " +
						"(Action, Romance): u = (4, 2). Sam's: v = (2, 1) — exactly half of u.",
					"• Dot product: u·v = (4)(2) + (2)(1) = 8 + 2 = 10 — multiply matching " +
						"components and add, same operation the vectors concept elsewhere in this " +
						"gallery uses.",
					"• Magnitudes (lengths): |u| = √(4²+2²) = √20 ≈ 4.47, |v| = √(2²+1²) = √5 ≈ " +
						"2.24.",
					"• Cosine similarity = (u·v)/(|u|·|v|) = 10/(4.47×2.24) = 10/10.00 = 1.00 — " +
						"dividing out both magnitudes cancels out exactly how generously each of " +
						"them rates, leaving only the ratio between their two genre scores. That " +
						"ratio (2-to-1, Action over Romance) is identical for both of them, so the " +
						"score hits its maximum, 1.00, even though Priya's numbers run twice as " +
						"high as Sam's.",
					"• Contrast with Euclidean (straight-line) distance between the same two " +
						"rating pairs: √((4-2)²+(2-1)²) = √5 ≈ 2.24 — very much not zero. Distance " +
						"says 'these two rating profiles are apart'; cosine similarity says 'these " +
						"two people like the same things, the same amount relative to each other.' " +
						"Both are true at once, and for 'should the app treat these two as having " +
						"similar taste,' direction is the one that matters.",
					"• A third user, Jordan, rates Action +1, Romance +5 — mostly here for the " +
						"romance: u·Jordan = 4(1)+2(5) = 4+10 = 14, |Jordan's vector| = √26 ≈ 5.10, " +
						"cosine similarity = 14/(4.47×5.10) ≈ 0.61, angle ≈ 52° — related to " +
						"Priya's taste, but noticeably less so than Sam's.",
					"• A fourth, Max, rates Action +2, Romance -4 — likes action, actively " +
						"dislikes romance: u·Max = 4(2)+2(-4) = 8-8 = 0 — cosine similarity exactly " +
						"0, a 90° angle: orthogonal, the vector-space way of saying 'no shared " +
						"taste direction with Priya at all.'",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The blue arrow is u — Priya's ratings by default, (Action, Romance) = " +
						"(4, 2); the orange arrow is v — Sam's, (2, 1) — both drawn from the " +
						"origin, and dragging either vector's sliders moves its arrow to compare " +
						"any two rating profiles you like. The green arc traces θ, the angle " +
						"between them, and the numbers above report u·v, cos θ, and θ itself, " +
						"live. The dashed circle has radius 1; the small blue and orange squares " +
						"sitting on it are u and v after dividing each by its own length — 'ignore " +
						"how loud a rater each of them is' made literal, since every vector lands " +
						"on the same circle regardless of how enthusiastic or reserved that rater " +
						"tends to be. With the defaults, Priya's and Sam's squares land in exactly " +
						"the same spot on the circle (cos θ = 1.00) even though Priya's arrow " +
						"reaches twice as far out as Sam's — and the gray 'Euclidean distance' " +
						"line underneath stays a stubborn 2.24, proof the two measures are " +
						"answering genuinely different questions. Scale v further out along the " +
						"same direction (say to (4,2), matching Priya exactly, or beyond to (8,4)) " +
						"and watch cos θ hold at 1.00 the whole time while the distance keeps " +
						"changing.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Match people (or, in a modern system, embeddings — vectors a model produces " +
						"to stand in for taste, meaning, or style) by what they actually prefer, " +
						"without a naturally more enthusiastic or more reserved rater throwing off " +
						"the comparison. This is exactly the comparison a recommendation engine " +
						"runs between your taste profile and every other user's (or between your " +
						"profile and every movie's) to answer 'who else has taste like mine' or " +
						"'what should we suggest you watch next.'",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Recommendation systems — 'people with taste like yours also loved...' — " +
						"comparing rating or interaction profiles exactly this way, whether the " +
						"categories are movie genres, songs, or products. The same comparison, " +
						"generalized to hundreds or thousands of dimensions instead of two, is also " +
						"the core of embedding-based semantic search — the retrieval step behind " +
						"modern AI assistants that look up relevant documents before answering, or " +
						"a search engine ranking results by relevance instead of raw keyword " +
						"overlap. Outside tech, everyday 'similar taste' already means roughly " +
						"this: two friends can have the exact same taste in movies while one of " +
						"them is simply a much harsher critic across the board — cosine similarity " +
						"is a precise version of exactly that intuition.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'cosine similarity is the cosine of the angle between two " +
						"rating vectors — it's blind to how generous or harsh each rater is, on " +
						"purpose.' Not like this: treating it like a distance, where smaller means " +
						"more similar — cosine similarity runs the other way, with 1.00 (0°, " +
						"identical taste direction) as the most similar and -1.00 (180°, exactly " +
						"opposite taste) as the least; or assuming a score of 0 always means 'no " +
						"relationship at all' in a real-world sense — it precisely means " +
						"'perpendicular in whatever rating space these vectors were built from,' " +
						"which for a two-genre example like Priya-vs-Max is a fairly literal 'no " +
						"shared taste direction,' but in a real system built on hundreds of genres " +
						"or a learned embedding space it's a subtler geometric statement, not a " +
						"plain-English verdict; and don't assume the score is bounded 0 to 1 the " +
						"way a percentage or probability is — it legitimately runs from -1 to 1, " +
						"matching the fact that taste can be not just 'unrelated' but 'opposite.'",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "ux", Label: "u — Action rating (-5 hate, +5 love)", Min: -5, Max: 5, Step: 1, Def: 4},
			{Key: "uy", Label: "u — Romance rating (-5 hate, +5 love)", Min: -5, Max: 5, Step: 1, Def: 2},
			{Key: "vx", Label: "v — Action rating (-5 hate, +5 love)", Min: -5, Max: 5, Step: 1, Def: 2},
			{Key: "vy", Label: "v — Romance rating (-5 hate, +5 love)", Min: -5, Max: 5, Step: 1, Def: 1},
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

	c.Text(16, 22, fmt.Sprintf("u = (Action %.0f, Romance %.0f)   |u| = %.2f", ux, uy, magU), 14, viz.Accent, "start")
	c.Text(16, 42, fmt.Sprintf("v = (Action %.0f, Romance %.0f)   |v| = %.2f", vx, vy, magV), 14, viz.Warm, "start")
	c.Text(16, 66, fmt.Sprintf("u · v = %.2f    cos θ = %.3f    θ = %.1f°", dot, cos, angle), 14, viz.Ink, "start")
	c.Text(16, 86, fmt.Sprintf("Euclidean distance |u - v| = %.2f  (cosine similarity ignores this entirely)", dist),
		13, viz.Muted, "start")
	c.Text(16, 500, "small squares on the dashed circle = u and v after dividing out how loud a rater each is",
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
