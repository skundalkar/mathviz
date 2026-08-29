// Package svm visualizes the support vector machine: among the infinitely
// many straight lines that separate two linearly-separable classes, the one
// that maximizes the gap (margin) to the nearest point on each side. The
// running example is four points -- two per class -- where the two closest
// points across classes turn out to be the only ones that matter: the
// "support vectors" that pin down the boundary.
package svm

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "support-vector-machine",
		Seq:   77,
		Title: "Support vector machine (the maximum-margin boundary)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`logistic-regression` fit one specific S-curve by minimizing log-loss over " +
						"every point, and `decision-trees` picked one specific split by scoring " +
						"information gain -- both land on a single answer because their scoring " +
						"rule prefers it. But once two classes are cleanly separable, there isn't " +
						"just one line that gets 100% right: rotate or shift a separating line a " +
						"little and it can still classify every training point correctly. Two " +
						"students working the same four-point dataset by hand could draw two " +
						"different 'perfect' lines and both be technically right. If accuracy on " +
						"the training points can't tell those lines apart, what can?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take four points: class +1 at (3,3) and (5,4); class -1 at (1,1) and " +
						"(-1,0). Instead of scoring a line by how many points it gets right (all " +
						"of these lines get all four right), score it by its margin: the distance " +
						"from the line to the single nearest point of either class, doubled to " +
						"cover both sides. A line that grazes past a point has almost no buffer " +
						"before a slightly different point of that class would land on the wrong " +
						"side; a line with room on both sides is safer.",
					"The two closest points across the two classes are (3,3) and (1,1) -- " +
						"distance sqrt((3-1)²+(3-1)²) = sqrt(8) ≈ 2.83 apart. The maximum-margin " +
						"line is their perpendicular bisector: the direction straight from (1,1) " +
						"to (3,3) is (2,2), so the line's normal vector (after scaling to length " +
						"1) is n ≈ (0.707, 0.707), passing through their midpoint (2,2). Checking " +
						"the other two points against this line: (5,4) sits 3.54 units out and " +
						"(-1,0) sits 3.54 units out on the correct sides -- both comfortably " +
						"farther from the line than (3,3) and (1,1) are (1.41 units each). Moving " +
						"(5,4) or (-1,0) anywhere in that safe zone never changes the boundary at " +
						"all; only (3,3) and (1,1) -- the 'support vectors' -- decide where the " +
						"line sits, which is exactly the pattern `lagrange-multipliers` pointed at: " +
						"the points with a nonzero multiplier at the solution are the ones that " +
						"actually constrain it.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"theta and offset draw your own candidate line (blue) by its angle and " +
						"distance from the origin; the green line and shaded band are the true " +
						"maximum-margin solution from the worked example above, with (3,3) and " +
						"(1,1) circled as its support vectors. The readout reports your line's " +
						"achieved margin -- or flags it as not separating the classes at all, if " +
						"you rotate it too far -- next to the best possible margin (2.83) for " +
						"comparison.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Pick, out of infinitely many training-accuracy-perfect lines, provably the " +
						"one that leaves the most breathing room before a new point near the " +
						"boundary would get misclassified -- and know in advance which handful of " +
						"points (the support vectors) that decision actually depends on, so " +
						"collecting more of the other points wouldn't change the answer at all.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Text classification (spam vs. not-spam from word-count features) and " +
						"bioinformatics (classifying tumor samples from a handful of gene-" +
						"expression readings) were classic SVM strongholds, especially when there " +
						"are more features than examples -- a regime where margin, not just " +
						"training accuracy, matters a lot for how well the boundary generalizes. " +
						"The 'kernel trick' (not shown here) lets the same max-margin idea draw " +
						"curved boundaries by measuring distance in a reshaped space instead of " +
						"the raw one.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the boundary is decided by the closest points across " +
						"classes, so only those points are support vectors' -- moving any other " +
						"point, as long as it stays on its own side of the margin, leaves the " +
						"line exactly where it was.",
					"Not like this: assuming the line that looks centered by eye, or the one a " +
						"different training method (like logistic regression's log-loss fit) " +
						"happens to produce, is automatically the max-margin one -- two lines can " +
						"both classify every point correctly while having very different margins, " +
						"and only measuring the distance to the nearest point tells them apart.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "theta", Label: "your line's angle (theta)", Min: 1, Max: 179, Step: 1, Def: 80, Unit: "°"},
			{Key: "offset", Label: "your line's offset from origin (d)", Min: -2, Max: 5, Step: 0.1, Def: 1.5},
		},
		Render: render,
	})
}

// Points and Labels are the four fixed observations the worked example
// walks through: class +1 at (3,3) and (5,4), class -1 at (1,1) and
// (-1,0). (3,3) and (1,1) are the closest pair across the two classes, so
// they turn out to be the only points that decide the boundary -- the
// support vectors.
var (
	Points = [][2]float64{{3, 3}, {5, 4}, {1, 1}, {-1, 0}}
	Labels = []int{1, 1, -1, -1}
)

// Normal returns the unit vector at angle thetaDeg (measured counter-
// clockwise from the positive x-axis), used as a candidate hyperplane's
// normal direction.
func Normal(thetaDeg float64) (nx, ny float64) {
	rad := thetaDeg * math.Pi / 180
	return math.Cos(rad), math.Sin(rad)
}

// FunctionalMargin returns how correctly and how confidently a unit-normal
// hyperplane {n·x = d} classifies one labeled point: label*(n·point - d).
// It is positive when the point sits on the side its label predicts,
// negative when the hyperplane misclassifies it, and its magnitude is the
// point's perpendicular distance to the line.
func FunctionalMargin(nx, ny, d, x, y float64, label int) float64 {
	return float64(label) * (nx*x + ny*y - d)
}

// MinFunctionalMargin returns the smallest FunctionalMargin over every
// (point, label) pair -- positive only when the hyperplane separates every
// point correctly, in which case it's also the perpendicular distance from
// the line to the single nearest point of either class.
func MinFunctionalMargin(points [][2]float64, labels []int, nx, ny, d float64) float64 {
	min := math.Inf(1)
	for i, pt := range points {
		if fm := FunctionalMargin(nx, ny, d, pt[0], pt[1], labels[i]); fm < min {
			min = fm
		}
	}
	return min
}

// Margin returns the full width of the gap a unit-normal hyperplane leaves
// between the two classes: twice the distance from the line to the nearest
// point. It is negative when the hyperplane misclassifies at least one
// point (Separates reports false in that case).
func Margin(points [][2]float64, labels []int, nx, ny, d float64) float64 {
	return 2 * MinFunctionalMargin(points, labels, nx, ny, d)
}

// Separates reports whether a unit-normal hyperplane classifies every point
// correctly (allowing points to sit exactly on the line).
func Separates(points [][2]float64, labels []int, nx, ny, d float64) bool {
	return MinFunctionalMargin(points, labels, nx, ny, d) >= 0
}

// MaxMarginFromPair returns the maximum-margin hyperplane determined by two
// support vectors, one from each class: the perpendicular bisector of the
// segment joining them. margin is the distance between pos and neg, which
// equals the full gap width whenever no other point lies closer to the
// resulting line than pos and neg do.
func MaxMarginFromPair(pos, neg [2]float64) (nx, ny, d, margin float64) {
	dx, dy := pos[0]-neg[0], pos[1]-neg[1]
	margin = math.Hypot(dx, dy)
	if margin == 0 {
		return 0, 0, 0, 0
	}
	nx, ny = dx/margin, dy/margin
	midx, midy := (pos[0]+neg[0])/2, (pos[1]+neg[1])/2
	d = nx*midx + ny*midy
	return nx, ny, d, margin
}

// LinePoints samples the implicit line {n·x = d} at n points spread across
// t in [tLo, tHi], walking along the line's direction (-ny, nx) from the
// point on the line closest to the origin (d*nx, d*ny). Mirrors viz.Sample's
// signature for an explicit function, but for a line given by its normal
// form instead.
func LinePoints(nx, ny, d, tLo, tHi float64, n int) [][2]float64 {
	if n < 1 {
		n = 1
	}
	pts := make([][2]float64, 0, n+1)
	ox, oy := d*nx, d*ny
	for i := 0; i <= n; i++ {
		t := tLo + (tHi-tLo)*float64(i)/float64(n)
		pts = append(pts, [2]float64{ox - t*ny, oy + t*nx})
	}
	return pts
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(560, 420, -3, 6, -3, 6).String()
}
