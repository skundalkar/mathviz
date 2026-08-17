package poisson

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestPMFKnownValues(t *testing.T) {
	// Pinned to the worked lambda=3 example. See LEARNINGS.md.
	cases := []struct {
		k    int
		want float64
	}{
		{0, 0.0498}, {1, 0.1494}, {2, 0.2240}, {3, 0.2240},
		{4, 0.1680}, {5, 0.1008}, {6, 0.0504}, {7, 0.0216}, {8, 0.0081},
	}
	for _, c := range cases {
		if got := PMF(3, c.k); !almostEqual(got, c.want, 1e-4) {
			t.Errorf("PMF(3,%d) = %v, want %v", c.k, got, c.want)
		}
	}
}

func TestPMFNegativeInputsAreZero(t *testing.T) {
	if got := PMF(3, -1); got != 0 {
		t.Errorf("PMF(3,-1) = %v, want 0", got)
	}
	if got := PMF(-1, 0); got != 0 {
		t.Errorf("PMF(-1,0) = %v, want 0", got)
	}
}

func TestPMFSumsToOne(t *testing.T) {
	for _, lambda := range []float64{0.5, 1, 3, 9, 15} {
		sum := 0.0
		for k := 0; k <= 200; k++ {
			sum += PMF(lambda, k)
		}
		if !almostEqual(sum, 1, 1e-6) {
			t.Errorf("sum of PMF(%v,.) over k=0..200 = %v, want 1", lambda, sum)
		}
	}
}

func TestPMFApproximatesBinomialInTheLimit(t *testing.T) {
	// The defining relationship: Binomial(n,p,k) -> Poisson(n*p,k) as n
	// grows and p shrinks with n*p held fixed. See LEARNINGS.md for the
	// n=1000, p=0.001, lambda=1 worked example.
	n, p := 1000, 0.001
	lambda := float64(n) * p
	for k := 0; k <= 5; k++ {
		binom := choose(n, k) * math.Pow(p, float64(k)) * math.Pow(1-p, float64(n-k))
		pois := PMF(lambda, k)
		if !almostEqual(binom, pois, 5e-4) {
			t.Errorf("k=%d: binomial(%d,%v,%d)=%v, Poisson(%v,%d)=%v, want close", k, n, p, k, binom, lambda, k, pois)
		}
	}
}

// choose is a small local helper (not exported by this package) used only to
// cross-check TestPMFApproximatesBinomialInTheLimit against
// binomial-distribution's own PMF formula, without importing that package.
func choose(n, k int) float64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	result := 1.0
	for i := 0; i < k; i++ {
		result *= float64(n-i) / float64(i+1)
	}
	return result
}

func TestCDFKnownValue(t *testing.T) {
	// CDF(3,3) = sum of PMF(0..3) = 0.0498+0.1494+0.2240+0.2240.
	if got := CDF(3, 3); !almostEqual(got, 0.6472, 1e-4) {
		t.Errorf("CDF(3,3) = %v, want 0.6472", got)
	}
}

func TestCDFIsNonDecreasing(t *testing.T) {
	lambda := 5.0
	prev := 0.0
	for k := 0; k <= 40; k++ {
		got := CDF(lambda, k)
		if got < prev-1e-12 {
			t.Errorf("CDF(%v,%d) = %v is less than CDF at k-1 = %v", lambda, k, got, prev)
		}
		prev = got
	}
	if !almostEqual(prev, 1, 1e-6) {
		t.Errorf("CDF(%v,40) = %v, want ~1", lambda, prev)
	}
}

func TestCDFNegativeKIsZero(t *testing.T) {
	if got := CDF(3, -1); got != 0 {
		t.Errorf("CDF(3,-1) = %v, want 0", got)
	}
}

func TestMeanEqualsVariance(t *testing.T) {
	// The defining identity this package's LEARNINGS.md entry calls out:
	// for a true Poisson distribution, mean and variance are the same
	// number, lambda, not two separately-estimated quantities.
	for _, lambda := range []float64{0.5, 1, 3, 9, 15} {
		if got := Mean(lambda); got != lambda {
			t.Errorf("Mean(%v) = %v, want %v", lambda, got, lambda)
		}
		if got := Variance(lambda); got != lambda {
			t.Errorf("Variance(%v) = %v, want %v", lambda, got, lambda)
		}
	}
}

func TestFactorialKnownValues(t *testing.T) {
	cases := []struct {
		k    int
		want float64
	}{{0, 1}, {1, 1}, {5, 120}, {10, 3628800}}
	for _, c := range cases {
		if got := factorial(c.k); got != c.want {
			t.Errorf("factorial(%d) = %v, want %v", c.k, got, c.want)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("poisson-distribution")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
