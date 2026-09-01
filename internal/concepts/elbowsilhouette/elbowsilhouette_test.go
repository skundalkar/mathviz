package elbowsilhouette

import (
	"math"
	"testing"
)

func approxEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// TestKMeansConverges checks that RunKMeans recovers the three real blobs
// exactly when k=3, matching k-means-clustering's own worked example.
func TestKMeansConverges(t *testing.T) {
	_, assignments := RunKMeans(3, 20)
	groups := map[int][]string{}
	for i, a := range assignments {
		groups[a] = append(groups[a], PointLabels[i])
	}
	if len(groups) != 3 {
		t.Fatalf("k=3 produced %d non-empty clusters, want 3", len(groups))
	}
	// Every point within a real blob must share a cluster.
	blobs := [][]string{{"A", "B", "C"}, {"D", "E", "F"}, {"G", "H", "I"}}
	clusterOf := map[string]int{}
	for i, a := range assignments {
		clusterOf[PointLabels[i]] = a
	}
	for _, blob := range blobs {
		want := clusterOf[blob[0]]
		for _, label := range blob {
			if clusterOf[label] != want {
				t.Errorf("blob %v split across clusters: %v", blob, clusterOf)
			}
		}
	}
}

func TestInertiaDecreasesWithK(t *testing.T) {
	want := map[int]float64{1: 138.00, 2: 66.00, 3: 12.00, 4: 9.00, 5: 6.00, 6: 4.50}
	var prev float64 = math.Inf(1)
	for k := 1; k <= 6; k++ {
		centroids, assignments := RunKMeans(k, 20)
		got := Inertia(centroids, assignments)
		if !approxEqual(got, want[k], 0.01) {
			t.Errorf("k=%d inertia = %.2f, want %.2f", k, got, want[k])
		}
		if got > prev {
			t.Errorf("k=%d inertia %.2f increased over k=%d's %.2f — inertia must be non-increasing in k", k, got, k-1, prev)
		}
		prev = got
	}
}

// TestElbowAtThree checks that inertia's drop from k=2->3 dwarfs its drop
// from k=3->4 onward — the "elbow" the method is named for, landing right
// where the data actually has three blobs.
func TestElbowAtThree(t *testing.T) {
	inertiaAt := func(k int) float64 {
		c, a := RunKMeans(k, 20)
		return Inertia(c, a)
	}
	drop23 := inertiaAt(2) - inertiaAt(3)
	drop34 := inertiaAt(3) - inertiaAt(4)
	if drop23 <= drop34*3 {
		t.Errorf("k=2->3 drop (%.2f) isn't much sharper than k=3->4 drop (%.2f) — expected a clear elbow at k=3", drop23, drop34)
	}
}

func TestSilhouetteScoreKnownValues(t *testing.T) {
	want := map[int]float64{1: 0.0, 2: 0.463, 3: 0.679, 4: 0.623, 5: 0.561, 6: 0.430}
	for k := 1; k <= 6; k++ {
		_, assignments := RunKMeans(k, 20)
		got := SilhouetteScore(assignments, k)
		if !approxEqual(got, want[k], 0.005) {
			t.Errorf("k=%d silhouette = %.3f, want %.3f", k, got, want[k])
		}
	}
}

// TestSilhouettePeaksAtThree checks the silhouette score, unlike inertia, is
// actually highest at k=3 -- not just "big enough", the single best k among
// those tried, agreeing with the elbow method's answer.
func TestSilhouettePeaksAtThree(t *testing.T) {
	best, bestScore := -1, math.Inf(-1)
	for k := 1; k <= 6; k++ {
		_, assignments := RunKMeans(k, 20)
		s := SilhouetteScore(assignments, k)
		if s > bestScore {
			best, bestScore = k, s
		}
	}
	if best != 3 {
		t.Errorf("silhouette score peaks at k=%d, want k=3", best)
	}
}

func TestDistanceAndSquaredDistance(t *testing.T) {
	a, b := Point{0, 0}, Point{3, 4}
	if !approxEqual(SquaredDistance(a, b), 25, 1e-9) {
		t.Errorf("SquaredDistance = %v, want 25", SquaredDistance(a, b))
	}
	if !approxEqual(Distance(a, b), 5, 1e-9) {
		t.Errorf("Distance = %v, want 5", Distance(a, b))
	}
}

func TestInitialCentroidsCount(t *testing.T) {
	for k := 1; k <= 6; k++ {
		if got := len(InitialCentroids(k)); got != k {
			t.Errorf("InitialCentroids(%d) returned %d centroids, want %d", k, got, k)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	svg := render(map[string]float64{"k": 3})
	if len(svg) == 0 {
		t.Fatal("render returned empty string")
	}
	if svg[:4] != "<svg" {
		t.Errorf("render output doesn't start with <svg: %q...", svg[:20])
	}
}
