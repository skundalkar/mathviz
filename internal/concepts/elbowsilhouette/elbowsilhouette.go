// Package elbowsilhouette visualizes two ways to pick k for k-means when
// nobody hands you the right number of clusters: watch where adding another
// cluster stops meaningfully shrinking the within-cluster distance (the
// elbow method), and separately score how cleanly each point sits in its own
// cluster versus the next-closest one (the silhouette score).
package elbowsilhouette

import (
	"fmt"
	"math"
	"strings"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "elbow-method-silhouette-score",
		Seq:   86,
		Title: "Elbow method & silhouette score (choosing k for k-means)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "k", Label: "Number of clusters (k)", Min: 1, Max: 6, Step: 1, Def: 3},
		},
		Render: render,
	})
}

// Point is a 2D data point or centroid.
type Point struct{ X, Y float64 }

// DataPoints and PointLabels are the exact same 9-point, three-blob worked
// example k-means-clustering uses (A-I; {A,B,C} near (2,2), {D,E,F} near
// (8,2), {G,H,I} near (5,8)) — reused unchanged, since "how many clusters"
// only makes sense as a follow-up question about a clustering you already
// know how to run.
var DataPoints = []Point{
	{1, 2}, {2, 1}, {3, 3}, // A, B, C
	{7, 1}, {8, 2}, {9, 3}, // D, E, F
	{4, 7}, {5, 9}, {6, 8}, // G, H, I
}

var PointLabels = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

// SquaredDistance returns the squared Euclidean distance between two points.
func SquaredDistance(a, b Point) float64 {
	dx, dy := a.X-b.X, a.Y-b.Y
	return dx*dx + dy*dy
}

// Distance returns the Euclidean distance between two points.
func Distance(a, b Point) float64 {
	return math.Sqrt(SquaredDistance(a, b))
}

// InitialCentroids picks k deterministic starting centroids for
// DataPoints: the points at k evenly-spaced indices through the dataset.
// Deterministic on purpose — this concept is about choosing k, not about
// how sensitive k-means is to its starting guesses (that's k-means-clustering's
// lesson), so the same k always produces the same starting point.
func InitialCentroids(k int) []Point {
	n := len(DataPoints)
	out := make([]Point, k)
	for i := 0; i < k; i++ {
		idx := int(math.Round((float64(i) + 0.5) * float64(n) / float64(k)))
		if idx >= n {
			idx = n - 1
		}
		out[i] = DataPoints[idx]
	}
	return out
}

// assign returns, for each point, the index of its nearest centroid.
func assign(points, centroids []Point) []int {
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

// update moves each centroid to the mean of the points currently assigned
// to it, keeping a centroid's previous position if nothing is assigned to it.
func update(points []Point, assignments []int, prev []Point) []Point {
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

// RunKMeans runs Lloyd's algorithm from InitialCentroids(k) until the
// assignments stop changing (or maxIterations is reached — the well
// separated blobs here always converge in well under 10), and returns the
// final centroids and per-point cluster assignments.
func RunKMeans(k, maxIterations int) (centroids []Point, assignments []int) {
	centroids = InitialCentroids(k)
	assignments = assign(DataPoints, centroids)
	for i := 0; i < maxIterations; i++ {
		centroids = update(DataPoints, assignments, centroids)
		next := assign(DataPoints, centroids)
		same := true
		for j := range next {
			if next[j] != assignments[j] {
				same = false
				break
			}
		}
		assignments = next
		if same {
			break
		}
	}
	return centroids, assignments
}

// Inertia is the within-cluster sum of squared distances: every point's
// squared distance to its own cluster's centroid, totaled. This is the
// quantity the elbow method plots against k — it can only go down (or stay
// flat) as k grows, since more centroids can only get points closer to
// their nearest one.
func Inertia(centroids []Point, assignments []int) float64 {
	var total float64
	for i, pt := range DataPoints {
		total += SquaredDistance(pt, centroids[assignments[i]])
	}
	return total
}

// SilhouetteScore averages, over every point, (b-a)/max(a,b), where a is
// the point's mean distance to the other points in its own cluster and b is
// its mean distance to the points of the nearest *other* cluster. It runs
// from -1 (points look closer to a different cluster than their own) to +1
// (perfectly separated clusters); k=1 always scores 0 since "the nearest
// other cluster" doesn't exist.
func SilhouetteScore(assignments []int, k int) float64 {
	n := len(DataPoints)
	if k <= 1 {
		return 0
	}
	var total float64
	for i, pt := range DataPoints {
		own := assignments[i]

		var aSum float64
		aCount := 0
		for j, other := range DataPoints {
			if j == i || assignments[j] != own {
				continue
			}
			aSum += Distance(pt, other)
			aCount++
		}
		a := 0.0
		if aCount > 0 {
			a = aSum / float64(aCount)
		}

		bBest := math.Inf(1)
		for cluster := 0; cluster < k; cluster++ {
			if cluster == own {
				continue
			}
			var bSum float64
			bCount := 0
			for j, other := range DataPoints {
				if assignments[j] != cluster {
					continue
				}
				bSum += Distance(pt, other)
				bCount++
			}
			if bCount == 0 {
				continue
			}
			if avg := bSum / float64(bCount); avg < bBest {
				bBest = avg
			}
		}

		s := 0.0
		if !math.IsInf(bBest, 1) {
			m := math.Max(a, bBest)
			if m > 0 {
				s = (bBest - a) / m
			}
		}
		total += s
	}
	return total / float64(n)
}

const maxK = 6

// groupsLabel returns a compact "{A,B,C} {D,E,F} ..." rendering of the
// current cluster assignments, so the curve can be tied back to the exact
// same labeled points k-means-clustering uses.
func groupsLabel(assignments []int, k int) string {
	groups := make([][]string, k)
	for i, a := range assignments {
		groups[a] = append(groups[a], PointLabels[i])
	}
	parts := make([]string, 0, k)
	for _, g := range groups {
		if len(g) == 0 {
			continue
		}
		parts = append(parts, "{"+strings.Join(g, ",")+"}")
	}
	return strings.Join(parts, " ")
}

func render(p map[string]float64) string {
	k := int(p["k"])
	if k < 1 {
		k = 1
	}
	if k > maxK {
		k = maxK
	}

	inertias := make([]float64, maxK+1) // 1-indexed by k
	silhouettes := make([]float64, maxK+1)
	var assignmentsAtK []int
	var maxInertia float64
	for kk := 1; kk <= maxK; kk++ {
		centroids, assignments := RunKMeans(kk, 20)
		inertias[kk] = Inertia(centroids, assignments)
		silhouettes[kk] = SilhouetteScore(assignments, kk)
		if inertias[kk] > maxInertia {
			maxInertia = inertias[kk]
		}
		if kk == k {
			assignmentsAtK = assignments
		}
	}

	const xmin, xmax = 0.5, float64(maxK) + 0.5
	const ymin, ymax = 0.0, 1.08
	c := viz.New(680, 460, xmin, xmax, ymin, ymax)
	c.PadT = 108
	c.Axes()
	for kk := 1; kk <= maxK; kk++ {
		c.Tick(float64(kk), fmt.Sprintf("%d", kk))
	}

	// Normalize inertia to [0,1] (divide by k=1's inertia) so it shares an
	// axis with the silhouette score, which is already 0-1 here.
	inertiaCurve := make([][2]float64, maxK)
	silhouetteCurve := make([][2]float64, maxK)
	for kk := 1; kk <= maxK; kk++ {
		inertiaCurve[kk-1] = [2]float64{float64(kk), inertias[kk] / maxInertia}
		silhouetteCurve[kk-1] = [2]float64{float64(kk), silhouettes[kk]}
	}

	c.VLine(float64(k), viz.Ink, true)
	c.Path(inertiaCurve, viz.Accent, 2.5)
	c.Path(silhouetteCurve, viz.Warm, 2.5)

	for kk := 1; kk <= maxK; kk++ {
		px, py := c.X(float64(kk)), c.Y(inertias[kk]/maxInertia)
		c.Rect(px-3, py-3, 6, 6, viz.Accent, 1)
		px2, py2 := c.X(float64(kk)), c.Y(silhouettes[kk])
		c.Rect(px2-3, py2-3, 6, 6, viz.Warm, 1)
	}

	c.Text(20, 22, fmt.Sprintf("k = %d    inertia = %.2f (raw, unnormalized)    silhouette = %.3f",
		k, inertias[k], silhouettes[k]), 14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("clusters at k=%d: %s", k, groupsLabel(assignmentsAtK, k)), 13, viz.Muted, "start")
	c.Text(20, 66, "blue = inertia, normalized to k=1's value    orange = silhouette score (already 0-1)", 12, viz.Muted, "start")
	c.Text(20, 86, "inertia keeps falling as k grows; the elbow is where it stops falling fast — here, k=3", 12, viz.Muted, "start")

	return c.String()
}
