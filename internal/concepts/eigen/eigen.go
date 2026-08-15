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
				Body: []string{
					"You're modeling something that gets transformed over and over — a webpage-" +
						"ranking algorithm bouncing traffic between linked pages one click at a " +
						"time, a population's age groups scaled by fixed birth and survival rates " +
						"each year, a photo filter applied 20 times in a row. Gut instinct: 'just " +
						"track where each point goes, one step at a time.' That works for one step " +
						"— but a single 2x2 matrix multiply can rotate a direction as well as " +
						"resize it, so after 20 chained steps almost every starting direction has " +
						"been bent to some new, unpredictable angle, and the only way to know where " +
						"it ends up is to grind through all 20 multiplications by hand. Is there a " +
						"shortcut — some special starting direction where you already know, without " +
						"doing the multiplication step by step, exactly what happens each round, " +
						"because that direction never rotates at all — it only grows or shrinks by " +
						"a fixed factor?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take the matrix A = [[2,1],[1,2]] and multiply it by a few different " +
						"direction vectors, one at a time, matching component with component and " +
						"summing (the same rule as Vectors' dot product, applied twice per output " +
						"coordinate):",
					"• v = (1, 0), pointing right. Av = (2·1+1·0, 1·1+2·0) = (2, 1) — that new " +
						"vector points up-and-right, rotated away from v's original direction.",
					"• v = (0, 1), pointing straight up. Av = (2·0+1·1, 1·0+2·1) = (1, 2) — " +
						"rotated again, this time swung toward the diagonal from the other side.",
					"• v = (1, 1), the 45° diagonal. Av = (2·1+1·1, 1·1+2·1) = (3, 3) = 3×(1, 1). " +
						"No rotation at all — Av lands exactly on v's own line, just 3 times as " +
						"long.",
					"• v = (1, −1), the other diagonal. Av = (2·1+1·(−1), 1·1+2·(−1)) = (1, −1) = " +
						"1×(1, −1). Again no rotation — this direction doesn't even change length, " +
						"since it's multiplied by exactly 1.",
					"Those last two directions are eigenvectors of A, and 3 and 1 are their " +
						"eigenvalues — 'eigen' is German for 'own' or 'characteristic': each " +
						"eigenvector keeps its own line under A, scaled by its own fixed number. " +
						"Why 45° and 135° specifically, and not some other pair of directions? For " +
						"a symmetric matrix [[a,b],[b,d]] the eigenvector angle has a closed-form " +
						"answer, θ = atan2(2b, a−d) / 2. Here a=d=2, so a−d=0 and θ = atan2(2,0)/2 " +
						"= 90°/2 = 45° — the shear term b alone tilts the special directions onto " +
						"the diagonal, and the other eigenvector always sits exactly 90° from the " +
						"first one, a fact proven for every symmetric matrix, not just this one.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The a, d, and b sliders set the symmetric matrix A = [[a,b],[b,d]]; the θ " +
						"slider sweeps a test vector v (blue arrow) all the way around the unit " +
						"circle. The orange arrow is Av. For almost every angle the two arrows " +
						"visibly point in different directions — watch them separate as you sweep " +
						"θ. The two faint green lines mark the exact eigenvector directions the " +
						"current matrix has right now; whenever θ lands on (or very near) one of " +
						"them, both arrows turn green and the readout reports 'eigenvector! Av = " +
						"λ×v', because Av has snapped onto v's own line. Drag a or d, or b, and " +
						"watch the two green lines themselves swing to new angles — the eigenvector " +
						"directions are a property of the matrix, and change the instant the matrix " +
						"does.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Predict the long-run behavior of a repeated transformation without " +
						"simulating it step by step: apply A to almost any starting vector n times " +
						"in a row, and it bends toward the eigenvector with the larger eigenvalue " +
						"(here, the 45° line, eigenvalue 3) and grows by roughly that eigenvalue " +
						"each round — 3ⁿ, not some unpredictable smear. That single fact is the " +
						"engine behind 'power iteration,' the standard trick for finding a " +
						"dominant eigenvector by just multiplying by A repeatedly and watching " +
						"where the direction settles, instead of solving for it algebraically.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A bridge or a guitar string doesn't vibrate in some arbitrary shape when " +
						"disturbed — it settles into a small set of natural vibration patterns " +
						"('modes'), each with its own natural frequency; those patterns are the " +
						"eigenvectors of the structure's stiffness matrix, and pushing the bridge " +
						"at exactly its lowest mode's frequency is how resonance can shake it apart " +
						"(the Tacoma Narrows Bridge collapse). In machine learning, PageRank ranks " +
						"webpages by the dominant eigenvector of the page-link matrix — repeatedly " +
						"following links is literally power iteration — and Principal Component " +
						"Analysis finds the directions data varies most along by taking the " +
						"eigenvectors of the data's covariance matrix, the same 'special directions " +
						"that only stretch' idea applied to a cloud of data points instead of a " +
						"single shape.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'Av = λv only holds exactly on an eigenvector's own line — " +
						"off that line, Av points somewhere else entirely, not just \"a bit off\" " +
						"from λv.' Not like this: assuming every direction has some eigenvalue " +
						"describing what A did to it — only the (at most two, for a 2x2 matrix) " +
						"eigenvector directions get a clean scale factor; every other direction " +
						"gets rotated, and 'rotated' has no single number that plays the same role " +
						"as an eigenvalue. Also not like this: reading a matrix's diagonal entries " +
						"as its eigenvalues — that shortcut only works when there's no shear (b=0); " +
						"the moment b≠0, the eigenvalues shift away from a and d, exactly as this " +
						"example shows: diagonal entries 2 and 2, but eigenvalues 3 and 1.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "Matrix a (x-axis stretch)", Min: 0.5, Max: 2.5, Step: 0.1, Def: 2},
			{Key: "d", Label: "Matrix d (y-axis stretch)", Min: 0.5, Max: 2.5, Step: 0.1, Def: 2},
			{Key: "b", Label: "Matrix b (shear / coupling)", Min: -1.2, Max: 1.2, Step: 0.1, Def: 1},
			{Key: "angle", Label: "Test vector angle θ", Min: 0, Max: 360, Step: 5, Def: 45, Unit: "°"},
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
