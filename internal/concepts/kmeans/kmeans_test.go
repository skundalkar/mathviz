package kmeans

import (
	"strings"
	"testing"

	"mathviz/internal/concept"
)

const epsilon = 1e-9

func almostEqual(a, b, tol float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= tol
}

func almostEqualPoint(a, b Point, tol float64) bool {
	return almostEqual(a.X, b.X, tol) && almostEqual(a.Y, b.Y, tol)
}

func TestSquaredDistanceKnownValues(t *testing.T) {
	if got := SquaredDistance(Point{0, 0}, Point{3, 4}); !almostEqual(got, 25, epsilon) {
		t.Errorf("SquaredDistance((0,0),(3,4)) = %v, want 25", got)
	}
	if got := SquaredDistance(Point{2, 2}, Point{2, 2}); !almostEqual(got, 0, epsilon) {
		t.Errorf("SquaredDistance of a point with itself = %v, want 0", got)
	}
}

func TestAssignPicksNearestCentroid(t *testing.T) {
	centroids := []Point{{0, 0}, {10, 0}}
	assignments := Assign([]Point{{1, 0}, {9, 0}}, centroids)
	if assignments[0] != 0 {
		t.Errorf("point near (0,0) assigned to centroid %d, want 0", assignments[0])
	}
	if assignments[1] != 1 {
		t.Errorf("point near (10,0) assigned to centroid %d, want 1", assignments[1])
	}
}

func TestAssignTiesBreakToLowerIndex(t *testing.T) {
	// (1,1) is equidistant from (1,1) itself... use a point exactly midway
	// between two centroids instead: (5,0) is equidistant from (0,0) and (10,0).
	centroids := []Point{{0, 0}, {10, 0}}
	assignments := Assign([]Point{{5, 0}}, centroids)
	if assignments[0] != 0 {
		t.Errorf("tied point assigned to centroid %d, want 0 (lower index)", assignments[0])
	}
}

func TestUpdateCentroidsIsTheMean(t *testing.T) {
	points := []Point{{1, 2}, {2, 1}, {3, 3}}
	assignments := []int{0, 0, 0}
	prev := []Point{{99, 99}}
	got := UpdateCentroids(points, assignments, prev)
	want := Point{2, 2}
	if !almostEqualPoint(got[0], want, epsilon) {
		t.Errorf("UpdateCentroids mean = %v, want %v", got[0], want)
	}
}

func TestUpdateCentroidsKeepsPreviousWhenEmpty(t *testing.T) {
	points := []Point{{1, 2}, {2, 1}}
	assignments := []int{0, 0} // nothing assigned to centroid 1
	prev := []Point{{0, 0}, {5, 5}}
	got := UpdateCentroids(points, assignments, prev)
	if !almostEqualPoint(got[1], prev[1], epsilon) {
		t.Errorf("empty cluster centroid = %v, want unchanged %v", got[1], prev[1])
	}
}

func TestLloydIterationZeroIsInitialAssignment(t *testing.T) {
	// From the worked example: with the "crammed corner" seeds, A and B
	// (both very close to seed centroid 0 and 1) tie-break to centroid 0
	// at iteration 0, before anything has moved.
	centroids, assignments := Lloyd(DataPoints, SeedCentroids(1), 0)
	if len(centroids) != 3 {
		t.Fatalf("Lloyd returned %d centroids, want 3", len(centroids))
	}
	if !almostEqualPoint(centroids[0], Point{1, 1}, epsilon) {
		t.Errorf("iteration-0 centroid 0 = %v, want unchanged seed (1,1)", centroids[0])
	}
	if assignments[0] != 0 || assignments[1] != 0 {
		t.Errorf("A,B assignments at iteration 0 = %v,%v, want 0,0", assignments[0], assignments[1])
	}
}

func TestLloydConvergesToTrueBlobMeans(t *testing.T) {
	// From the worked example: whichever seed set is used, by iteration 2
	// the three centroids have settled on the true blob means (2,2),
	// (8,2), (5,8), and running a further iteration changes nothing.
	for _, seed := range []int{0, 1} {
		centroids, assignments := Lloyd(DataPoints, SeedCentroids(seed), 2)
		more, moreAssignments := Lloyd(DataPoints, SeedCentroids(seed), 3)
		for k := range centroids {
			if !almostEqualPoint(centroids[k], more[k], epsilon) {
				t.Errorf("seed %d: centroid %d moved from iteration 2 (%v) to 3 (%v), expected convergence",
					seed, k, centroids[k], more[k])
			}
		}
		for i := range assignments {
			if assignments[i] != moreAssignments[i] {
				t.Errorf("seed %d: point %d assignment changed from iteration 2 to 3, expected convergence", seed, i)
			}
		}
	}
}

func TestLloydSameFinalGroupingRegardlessOfSeed(t *testing.T) {
	// Both seed sets should converge on the same {A,B,C}/{D,E,F}/{G,H,I}
	// partition (just possibly with the three clusters in a different
	// centroid-index order), since the blobs are well separated.
	_, a0 := Lloyd(DataPoints, SeedCentroids(0), 3)
	_, a1 := Lloyd(DataPoints, SeedCentroids(1), 3)
	sameGroup := func(assignments []int, i, j int) bool { return assignments[i] == assignments[j] }
	pairs := [][2]int{{0, 1}, {1, 2}, {3, 4}, {4, 5}, {6, 7}, {7, 8}} // within each blob
	for _, pr := range pairs {
		if sameGroup(a0, pr[0], pr[1]) != sameGroup(a1, pr[0], pr[1]) {
			t.Errorf("seed 0 and seed 1 disagree on whether points %d,%d are grouped together", pr[0], pr[1])
		}
		if !sameGroup(a0, pr[0], pr[1]) {
			t.Errorf("points %d,%d (same blob) not grouped together at convergence", pr[0], pr[1])
		}
	}
}

func TestSeedCentroidsClamps(t *testing.T) {
	if got := SeedCentroids(-5); len(got) != len(SeedCentroids(0)) {
		t.Errorf("SeedCentroids(-5) len = %d, want clamped to seed 0's length", len(got))
	}
	if got := SeedCentroids(99); len(got) != len(SeedCentroids(1)) {
		t.Errorf("SeedCentroids(99) len = %d, want clamped to seed 1's length", len(got))
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("k-means-clustering")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
