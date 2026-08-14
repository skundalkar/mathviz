// Package kmeans visualizes k-means clustering: repeatedly assigning each
// point to its nearest centroid, then moving each centroid to the average
// of the points now assigned to it, until the groups stop changing.
package kmeans

import (
	"fmt"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "k-means-clustering",
		Seq:   32,
		Title: "K-means clustering (grouping points by nearest centroid)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You're handed a spreadsheet of customers, each row just (visits per " +
						"month, average spend) — no labels, nobody has told you which customers " +
						"are 'similar' — and you want to group them so three different email " +
						"campaigns can go to three different types of shopper. Gut instinct: 'just " +
						"plot it and eyeball where the clumps are.' That works fine on a napkin, " +
						"for a handful of points, on exactly two measurements. It breaks down " +
						"fast: with three or more measurements per customer you can't even draw " +
						"the scatter plot anymore, with thousands of points the clumps blur into " +
						"one smear, and two different people staring at the same cloud will circle " +
						"different boundaries. Is there a mechanical, repeatable rule — one a " +
						"computer could run with nobody circling anything by hand — that keeps " +
						"producing the same grouping from the same data every time?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Guess a number of groups, k, and k starting 'centers' anywhere at all — " +
						"they don't need to be good guesses. Then repeat two steps until nothing " +
						"changes: (1) assign every point to whichever center is currently closest " +
						"to it, (2) move each center to the average position of the points just " +
						"assigned to it. Work a small example with k=3: nine points forming three " +
						"visible clumps — {A(1,2), B(2,1), C(3,3)} near (2,2), {D(7,1), E(8,2), " +
						"F(9,3)} near (8,2), {G(4,7), H(5,9), I(6,8)} near (5,8) — and three " +
						"starting centers guessed badly on purpose, to make the process visible: " +
						"center 1 at (1,1), center 2 at (2,2), center 3 at (9,9) (two guesses " +
						"crammed into the first blob's corner, one guess covering both far blobs " +
						"at once).",
					"• Assign step 1: compare each point's distance to all three centers. A(1,2) " +
						"is equally close to center 1 and center 2 (distance² = 1 either way) and " +
						"goes to whichever is checked first, center 1. G(4,7) is equally close to " +
						"center 2 and center 3 (distance² = 29 either way) and goes to center 2. " +
						"Working through all nine points this way gives a messy first grouping: " +
						"center 1 = {A, B}; center 2 = {C, D, E, G}; center 3 = {F, H, I} — center " +
						"2 grabbed the nearby point C, but also D and E from the second blob and G " +
						"from the third.",
					"• Update step 1: move each center to the mean of its current group. Center 1 " +
						"→ mean(A,B) = (1.5, 1.5). Center 2 → mean(C,D,E,G) = " +
						"((3+7+8+4)/4, (3+1+2+7)/4) = (5.5, 3.25). Center 3 → mean(F,H,I) = " +
						"((9+5+6)/3, (3+9+8)/3) ≈ (6.67, 6.67).",
					"• Assign step 2: recompute nearest-center with the moved centers. Every " +
						"point now lands in its real blob: center 1 = {A, B, C}; center 2 = " +
						"{D, E, F}; center 3 = {G, H, I} — one update cycle already recovered the " +
						"three real groups, even though the starting guesses were bad.",
					"• Update step 2: recompute the means one more time — center 1 → (2, 2), " +
						"center 2 → (8, 2), center 3 → (5, 8), exactly the three blobs' true " +
						"centers. Reassigning again changes nothing: that 'nothing changed' is " +
						"exactly the signal to stop.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Each of the nine points (labeled A-I, matching the worked example above) is " +
						"a small square colored by whichever center it's currently assigned to; " +
						"the three larger bordered squares are the centers themselves; a thin " +
						"line connects every point to its current center, so a point switching " +
						"groups shows up as a line jumping, not just a color changing. The " +
						"Iteration slider steps through the process above exactly: at 0 you see " +
						"the messy first assignment (center 2's lines reaching out to grab D, E, " +
						"and G); move it to 1 and the lines redraw around the newly-moved centers, " +
						"correctly regrouping into the three blobs; move it to 2 and the centers " +
						"visibly slide the rest of the way onto the true blob centers computed " +
						"above; moving it further changes nothing at all — convergence, made " +
						"visible. The starting-centroids toggle swaps in the other seed set (three " +
						"well-spread guesses, one per blob), where the very first assignment is " +
						"already correct and only the centers' exact positions still need to catch " +
						"up — worth flipping back and forth to see how much a bad starting guess " +
						"changes the trajectory, even when (as here) it doesn't change the final " +
						"answer.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Automatically and repeatably sort unlabeled points into groups — for as " +
						"many points and as many measurements per point as you have, since " +
						"'distance between two points' and 'average of a group' both keep working " +
						"past two dimensions even though a picture can't — with no eyeballing, and " +
						"the same input always producing the same output. This is exactly how a " +
						"real customer-segmentation tool sorts thousands of shoppers into a " +
						"handful of groups worth targeting differently, without anyone drawing " +
						"circles on a scatter plot.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Customer segmentation (grouping shoppers by behavior so different groups " +
						"get different marketing), image compression (grouping similar pixel " +
						"colors down to a small palette), grouping news articles or documents by " +
						"rough topic, and as a common first step inside bigger machine-learning " +
						"pipelines (picking a handful of 'typical' examples to represent a much " +
						"larger dataset). It's one of the first tools reached for whenever the " +
						"goal is 'let the data suggest the groups' instead of defining categories " +
						"by hand ahead of time.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'k-means repeatedly assigns each point to its nearest " +
						"center, then moves each center to the mean of the points now assigned to " +
						"it, until nothing changes' — the centers are not real data points, " +
						"they're computed averages that happen to land somewhere in the middle of " +
						"a group. Not like this: assuming the algorithm figures out how many " +
						"groups exist on its own — you choose k yourself ahead of time (fixed at 3 " +
						"here to match the three visible blobs; k=2 would force two of these real " +
						"blobs to merge, k=4 would split one blob in half, and either produces a " +
						"technically-valid but wrong-feeling answer); or assuming k-means always " +
						"converges to the same grouping no matter the starting guesses — it does " +
						"here because the three blobs are so well separated, but on messier or " +
						"overlapping data, different starting centers can converge to genuinely " +
						"different final groupings (a 'local optimum'), not just take a different " +
						"number of iterations to arrive at the same one.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "iteration", Label: "Iteration", Min: 0, Max: 3, Step: 1, Def: 0},
			{Key: "seed", Label: "Starting centroids", Min: 0, Max: 1, Step: 1, Def: 1},
		},
		Render: render,
	})
}

// Point is a 2D data point or centroid.
type Point struct{ X, Y float64 }

// DataPoints is the fixed worked-example dataset every Section walks
// through: nine points (labeled A-I, in this order) forming three visible
// blobs — {A,B,C} near (2,2), {D,E,F} near (8,2), {G,H,I} near (5,8).
var DataPoints = []Point{
	{1, 2}, {2, 1}, {3, 3}, // A, B, C
	{7, 1}, {8, 2}, {9, 3}, // D, E, F
	{4, 7}, {5, 9}, {6, 8}, // G, H, I
}

// PointLabels names DataPoints in order, for the picture and the docs.
var PointLabels = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

// seedSets holds two starting-centroid choices for the same three-blob
// dataset, selected by the "seed" Param: 0 is spread across the whole data
// range and lands on the right grouping immediately; 1 crams two of the
// three guesses into the same corner (near {A,B,C}) and needs an
// iteration to untangle, the more instructive default.
var seedSets = [][]Point{
	{{0, 0}, {10, 0}, {5, 10}},
	{{1, 1}, {2, 2}, {9, 9}},
}

// SeedCentroids returns a copy of the starting centroids for the given
// seed choice (0 or 1; out-of-range values clamp to the nearest end).
func SeedCentroids(seed int) []Point {
	if seed < 0 {
		seed = 0
	}
	if seed >= len(seedSets) {
		seed = len(seedSets) - 1
	}
	return append([]Point(nil), seedSets[seed]...)
}

// SquaredDistance returns the squared Euclidean distance between two
// points. Squared, not the true distance, because k-means only ever needs
// to compare distances against each other — the square root cancels out
// of every comparison and is just wasted work.
func SquaredDistance(a, b Point) float64 {
	dx, dy := a.X-b.X, a.Y-b.Y
	return dx*dx + dy*dy
}

// Assign returns, for each point, the index of its nearest centroid. Ties
// resolve to the lower centroid index, since a point is only reassigned to
// a strictly closer centroid.
func Assign(points, centroids []Point) []int {
	assignments := make([]int, len(points))
	for i, pt := range points {
		best, bestDist := 0, SquaredDistance(pt, centroids[0])
		for k := 1; k < len(centroids); k++ {
			if d := SquaredDistance(pt, centroids[k]); d < bestDist {
				best, bestDist = k, d
			}
		}
		assignments[i] = best
	}
	return assignments
}

// UpdateCentroids moves each centroid to the mean of the points currently
// assigned to it. A centroid with no points assigned keeps its previous
// position (prev) rather than dividing by zero.
func UpdateCentroids(points []Point, assignments []int, prev []Point) []Point {
	sums := make([]Point, len(prev))
	counts := make([]int, len(prev))
	for i, pt := range points {
		k := assignments[i]
		sums[k].X += pt.X
		sums[k].Y += pt.Y
		counts[k]++
	}
	out := make([]Point, len(prev))
	for k := range prev {
		if counts[k] == 0 {
			out[k] = prev[k]
			continue
		}
		out[k] = Point{sums[k].X / float64(counts[k]), sums[k].Y / float64(counts[k])}
	}
	return out
}

// Lloyd runs Lloyd's algorithm from initCentroids for the given number of
// iterations and returns the resulting centroids and point assignments.
// Iteration 0 is the starting state: the initial centroids with points
// assigned to whichever is nearest, before any centroid has moved. Each
// iteration after that is one assign-then-update-then-reassign cycle:
// move every centroid to the mean of its currently assigned points, then
// recompute every point's nearest centroid against the moved centroids.
func Lloyd(points, initCentroids []Point, iterations int) (centroids []Point, assignments []int) {
	centroids = append([]Point(nil), initCentroids...)
	assignments = Assign(points, centroids)
	for i := 0; i < iterations; i++ {
		centroids = UpdateCentroids(points, assignments, centroids)
		assignments = Assign(points, centroids)
	}
	return centroids, assignments
}

var clusterColor = []string{viz.Accent, viz.Warm, viz.Good}

func render(p map[string]float64) string {
	iteration := int(p["iteration"])
	seed := int(p["seed"])

	initCentroids := SeedCentroids(seed)
	centroids, assignments := Lloyd(DataPoints, initCentroids, iteration)

	c := viz.New(600, 440, -1, 11, -1, 11)
	c.Axes()

	// A thin line from each point to the centroid it's currently assigned
	// to, so a reassignment (a line jumping to a different centroid) is as
	// visible as a centroid moving.
	for i, pt := range DataPoints {
		cen := centroids[assignments[i]]
		c.Path([][2]float64{{pt.X, pt.Y}, {cen.X, cen.Y}}, viz.Faint, 1)
	}

	for i, pt := range DataPoints {
		px, py := c.X(pt.X), c.Y(pt.Y)
		color := clusterColor[assignments[i]%len(clusterColor)]
		c.Rect(px-4, py-4, 8, 8, color, 0.85)
		c.Text(px+7, py-6, PointLabels[i], 11, viz.Muted, "start")
	}

	for k, cen := range centroids {
		px, py := c.X(cen.X), c.Y(cen.Y)
		c.Rect(px-7, py-7, 14, 14, viz.Ink, 0.9)
		c.Rect(px-5, py-5, 10, 10, clusterColor[k%len(clusterColor)], 1)
	}

	seedLabel := "spread-out corners"
	if seed == 1 {
		seedLabel = "two guesses crammed into one corner"
	}
	c.Text(20, 20, "squares = data points, colored by current assignment    bigger boxes = centroids", 12, viz.Muted, "start")
	c.Text(20, 400, fmt.Sprintf("starting centroids: %s", seedLabel), 13, viz.Ink, "start")
	c.Text(20, 420, fmt.Sprintf("iteration %d of Lloyd's algorithm (assign to nearest centroid, then move each centroid to its group's mean)", iteration),
		12, viz.Muted, "start")

	return c.String()
}
