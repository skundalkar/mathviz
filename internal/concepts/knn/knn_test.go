package knn

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func approx(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestDistanceKnownValues(t *testing.T) {
	if !approx(Distance(0, 0, 3, 4), 5, 1e-9) {
		t.Errorf("Distance(0,0,3,4) = %v, want 5", Distance(0, 0, 3, 4))
	}
	if !approx(Distance(1, 1, 1, 1), 0, 1e-9) {
		t.Errorf("Distance to itself should be 0")
	}
}

func TestNearestSortsClosestFirst(t *testing.T) {
	// Section 2's worked query: (4,6). The two closest training points are
	// an exact-distance tie -- (2,4,fail) and (2,8,pass), both 2.828 -- and
	// (2,4,fail) comes first in TrainingSet, so the stable sort must favor
	// it as neighbor 0.
	neighbors := Nearest(4, 6, 10, TrainingSet)
	if len(neighbors) != 10 {
		t.Fatalf("Nearest(k=10) returned %d neighbors, want 10", len(neighbors))
	}
	want := []Point{
		{2, 4, 0}, {2, 8, 1}, {7, 8, 1}, {3, 2, 0}, {1, 2, 0},
		{8, 9, 1}, {9, 7, 1}, {2, 1, 0}, {8, 2, 0}, {9, 9, 1},
	}
	for i, w := range want {
		if neighbors[i].Point != w {
			t.Errorf("neighbor[%d] = %v, want %v", i, neighbors[i].Point, w)
		}
	}
	if !approx(neighbors[0].Distance, neighbors[1].Distance, 1e-9) {
		t.Errorf("neighbors 0 and 1 should be an exact distance tie (both 2.828), got %v and %v",
			neighbors[0].Distance, neighbors[1].Distance)
	}
}

func TestClassifyOscillatesAsKGrows(t *testing.T) {
	// The exact flip-flopping sequence from section 2: FAIL, PASS, FAIL,
	// PASS, FAIL as k climbs from 1 to 9 by 2, for query (4,6).
	cases := []struct {
		k    int
		want int
	}{
		{1, 0}, {3, 1}, {5, 0}, {7, 1}, {9, 0},
	}
	for _, c := range cases {
		got := Classify(4, 6, c.k, TrainingSet)
		if got != c.want {
			t.Errorf("Classify(4,6,k=%d) = %d, want %d", c.k, got, c.want)
		}
	}
}

func TestVoteCountsWorkedExample(t *testing.T) {
	neighbors := Nearest(4, 6, 5, TrainingSet)
	label0, label1 := VoteCounts(neighbors)
	if label0 != 3 || label1 != 2 {
		t.Errorf("k=5 vote counts = (%d,%d), want (3,2)", label0, label1)
	}
}

func TestNearestClampsKToTrainingSetSize(t *testing.T) {
	neighbors := Nearest(0, 0, 999, TrainingSet)
	if len(neighbors) != len(TrainingSet) {
		t.Errorf("Nearest with oversized k returned %d neighbors, want %d", len(neighbors), len(TrainingSet))
	}
}

func TestClassifyDeepInsideEachCluster(t *testing.T) {
	if got := Classify(2, 2, 3, TrainingSet); got != 0 {
		t.Errorf("point deep in the fail cluster classified as %d, want 0", got)
	}
	if got := Classify(8, 8, 3, TrainingSet); got != 1 {
		t.Errorf("point deep in the pass cluster classified as %d, want 1", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("k-nearest-neighbors")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
