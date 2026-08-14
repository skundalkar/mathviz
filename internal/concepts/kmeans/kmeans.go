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
				Body:    []string{"placeholder"},
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
