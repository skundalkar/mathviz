package geometric

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestPMFFairCoinKnownValues(t *testing.T) {
	// Fair coin, worked by hand: P(1)=0.5, P(2)=0.25, P(3)=0.125, P(4)=0.0625.
	want := []float64{0.5, 0.25, 0.125, 0.0625}
	for i, w := range want {
		k := i + 1
		if got := PMF(0.5, k); !near(got, w, 1e-9) {
			t.Errorf("PMF(0.5, %d) = %v, want %v", k, got, w)
		}
	}
}

func TestPMFInvalidK(t *testing.T) {
	if got := PMF(0.5, 0); got != 0 {
		t.Errorf("PMF(0.5, 0) = %v, want 0 (k must be >= 1)", got)
	}
}

func TestCDFMatchesSummedPMF(t *testing.T) {
	p := 0.3
	sum := 0.0
	for k := 1; k <= 10; k++ {
		sum += PMF(p, k)
		if got := CDF(p, k); !near(got, sum, 1e-9) {
			t.Errorf("CDF(%v, %d) = %v, want %v (sum of PMF 1..%d)", p, k, got, sum, k)
		}
	}
}

func TestCDFApproachesOne(t *testing.T) {
	if got := CDF(0.3, 50); !near(got, 1, 1e-6) {
		t.Errorf("CDF(0.3, 50) = %v, want close to 1", got)
	}
}

func TestMeanFairCoin(t *testing.T) {
	// A fair coin takes 2 flips on average to see the first heads.
	if got := Mean(0.5); !near(got, 2, 1e-9) {
		t.Errorf("Mean(0.5) = %v, want 2", got)
	}
}

func TestMeanMatchesExpectedValueDefinition(t *testing.T) {
	// Sum_{k=1..N} k*PMF(p,k) should converge to Mean(p) = 1/p as N grows.
	p := 0.4
	sum := 0.0
	for k := 1; k <= 200; k++ {
		sum += float64(k) * PMF(p, k)
	}
	if got, want := sum, Mean(p); !near(got, want, 1e-3) {
		t.Errorf("Sum k*PMF(0.4,k) = %v, want ~%v (Mean(0.4))", got, want)
	}
}

func TestPMFMonotonicDecreasing(t *testing.T) {
	p := 0.3
	prev := PMF(p, 1)
	for k := 2; k <= 20; k++ {
		cur := PMF(p, k)
		if cur >= prev {
			t.Errorf("PMF(%v, %d)=%v should be less than PMF(%v, %d)=%v", p, k, cur, p, k-1, prev)
		}
		prev = cur
	}
}

func TestMemorylessness(t *testing.T) {
	// P(X = k+m | X > k) should equal PMF(p, m) -- having already failed k
	// times doesn't change the distribution of how many more trials it
	// takes, the geometric distribution's defining "no memory" property.
	p := 0.3
	k, m := 5, 3
	probSurviveK := 1 - CDF(p, k) // P(X > k)
	probJoint := PMF(p, k+m)      // P(X = k+m)
	got := probJoint / probSurviveK
	want := PMF(p, m)
	if !near(got, want, 1e-9) {
		t.Errorf("P(X=%d | X>%d) = %v, want %v (PMF(p,%d), memorylessness)", k+m, k, got, want, m)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("geometric-distribution")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
