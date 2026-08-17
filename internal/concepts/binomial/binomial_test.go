package binomial

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestChooseKnownValues(t *testing.T) {
	cases := []struct {
		n, k int
		want float64
	}{
		{10, 5, 252},
		{4, 2, 6},
		{5, 0, 1},
		{5, 5, 1},
		{50, 25, 1.26410606437752e+14},
	}
	for _, c := range cases {
		if got := Choose(c.n, c.k); !almostEqual(got, c.want, c.want*1e-9+1e-9) {
			t.Errorf("Choose(%d,%d) = %v, want %v", c.n, c.k, got, c.want)
		}
	}
}

func TestChooseOutOfRange(t *testing.T) {
	if got := Choose(5, -1); got != 0 {
		t.Errorf("Choose(5,-1) = %v, want 0", got)
	}
	if got := Choose(5, 6); got != 0 {
		t.Errorf("Choose(5,6) = %v, want 0", got)
	}
}

func TestPMFKnownValues(t *testing.T) {
	// Pinned to the worked n=10, p=0.5 example. See LEARNINGS.md.
	cases := []struct {
		k    int
		want float64
	}{
		{0, 0.0010}, {1, 0.0098}, {2, 0.0439}, {3, 0.1172}, {4, 0.2051},
		{5, 0.2461}, {6, 0.2051}, {7, 0.1172}, {8, 0.0439}, {9, 0.0098}, {10, 0.0010},
	}
	for _, c := range cases {
		if got := PMF(10, 0.5, c.k); !almostEqual(got, c.want, 1e-4) {
			t.Errorf("PMF(10,0.5,%d) = %v, want %v", c.k, got, c.want)
		}
	}
}

func TestPMFSymmetricDistribution(t *testing.T) {
	// n=4, p=0.5: 1,4,6,4,1 out of 16.
	want := []float64{0.0625, 0.25, 0.375, 0.25, 0.0625}
	for k, w := range want {
		if got := PMF(4, 0.5, k); !almostEqual(got, w, 1e-9) {
			t.Errorf("PMF(4,0.5,%d) = %v, want %v", k, got, w)
		}
	}
}

func TestPMFOutOfRangeIsZero(t *testing.T) {
	if got := PMF(10, 0.5, -1); got != 0 {
		t.Errorf("PMF(10,0.5,-1) = %v, want 0", got)
	}
	if got := PMF(10, 0.5, 11); got != 0 {
		t.Errorf("PMF(10,0.5,11) = %v, want 0", got)
	}
}

func TestPMFSumsToOne(t *testing.T) {
	for _, n := range []int{1, 10, 30, 50} {
		for _, p := range []float64{0.05, 0.3, 0.5, 0.9} {
			sum := 0.0
			for k := 0; k <= n; k++ {
				sum += PMF(n, p, k)
			}
			if !almostEqual(sum, 1, 1e-9) {
				t.Errorf("sum of PMF(%d,%v,.) = %v, want 1", n, p, sum)
			}
		}
	}
}

func TestCDFKnownValue(t *testing.T) {
	// CDF(10,0.5,5) = sum of PMF(0..5) = 638/1024.
	if got := CDF(10, 0.5, 5); !almostEqual(got, 0.623046875, 1e-9) {
		t.Errorf("CDF(10,0.5,5) = %v, want 0.623046875", got)
	}
}

func TestCDFIsNonDecreasingAndEndsAtOne(t *testing.T) {
	n, p := 20, 0.4
	prev := 0.0
	for k := 0; k <= n; k++ {
		got := CDF(n, p, k)
		if got < prev-1e-12 {
			t.Errorf("CDF(%d,%v,%d) = %v is less than CDF at k-1 = %v", n, p, k, got, prev)
		}
		prev = got
	}
	if !almostEqual(prev, 1, 1e-9) {
		t.Errorf("CDF(%d,%v,%d) = %v, want 1", n, p, n, prev)
	}
}

func TestCDFClampsNegativeAndAboveN(t *testing.T) {
	if got := CDF(10, 0.5, -1); got != 0 {
		t.Errorf("CDF(10,0.5,-1) = %v, want 0", got)
	}
	if got := CDF(10, 0.5, 20); !almostEqual(got, 1, 1e-9) {
		t.Errorf("CDF(10,0.5,20) = %v, want 1", got)
	}
}

func TestMeanAndVarianceKnownValues(t *testing.T) {
	if got := Mean(100, 0.3); !almostEqual(got, 30, 1e-9) {
		t.Errorf("Mean(100,0.3) = %v, want 30", got)
	}
	if got := Variance(100, 0.3); !almostEqual(got, 21, 1e-9) {
		t.Errorf("Variance(100,0.3) = %v, want 21", got)
	}
}

func TestExactCountShareShrinksAsNGrows(t *testing.T) {
	// The single most likely count still owns a shrinking share of the
	// probability as n grows, even at p=0.5 (see LEARNINGS.md).
	p10 := PMF(10, 0.5, 5)
	p100 := PMF(100, 0.5, 50)
	if p100 >= p10 {
		t.Errorf("PMF(100,0.5,50)=%v should be less than PMF(10,0.5,5)=%v", p100, p10)
	}
	if !almostEqual(p100, 0.0796, 1e-3) {
		t.Errorf("PMF(100,0.5,50) = %v, want ~0.0796", p100)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("binomial-distribution")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
