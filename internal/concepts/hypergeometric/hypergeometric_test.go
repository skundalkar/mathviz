package hypergeometric

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func approx(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestChooseKnownValues(t *testing.T) {
	cases := []struct {
		n, k int
		want float64
	}{
		{20, 5, 15504},
		{8, 0, 1},
		{8, 2, 28},
		{12, 3, 220},
		{5, 6, 0}, // k > n
		{5, -1, 0},
	}
	for _, c := range cases {
		got := Choose(c.n, c.k)
		if !approx(got, c.want, 1e-6) {
			t.Errorf("Choose(%d,%d) = %v, want %v", c.n, c.k, got, c.want)
		}
	}
}

func TestPMFWorkedExample(t *testing.T) {
	// Section 2's marble bag: N=20, K=8, n=5.
	want := []float64{0.0511, 0.2554, 0.3973, 0.2384, 0.0542, 0.0036}
	for k, w := range want {
		got := PMF(20, 8, 5, k)
		if !approx(got, w, 5e-4) {
			t.Errorf("PMF(20,8,5,%d) = %v, want ~%v", k, got, w)
		}
	}
}

func TestPMFSumsToOne(t *testing.T) {
	sum := 0.0
	for k := 0; k <= 5; k++ {
		sum += PMF(20, 8, 5, k)
	}
	if !approx(sum, 1.0, 1e-9) {
		t.Errorf("PMF sums to %v, want 1.0", sum)
	}
}

func TestPMFOutOfRangeIsZero(t *testing.T) {
	if PMF(20, 8, 5, 6) != 0 {
		t.Error("PMF(k=6) with n=5 draws should be 0 -- can't draw more successes than trials")
	}
	if PMF(20, 8, 5, -1) != 0 {
		t.Error("PMF(k=-1) should be 0")
	}
	// Can't draw more reds than exist in the population, or more
	// non-reds than exist either.
	if PMF(20, 3, 5, 4) != 0 {
		t.Error("PMF(k=4) with only K=3 successes available should be 0")
	}
}

func TestCDFMatchesSummedPMF(t *testing.T) {
	got := CDF(20, 8, 5, 2)
	want := PMF(20, 8, 5, 0) + PMF(20, 8, 5, 1) + PMF(20, 8, 5, 2)
	if !approx(got, want, 1e-9) {
		t.Errorf("CDF(k=2) = %v, want %v", got, want)
	}
}

func TestMeanMatchesWorkedExample(t *testing.T) {
	if got := Mean(20, 8, 5); !approx(got, 2.0, 1e-9) {
		t.Errorf("Mean(20,8,5) = %v, want 2.0", got)
	}
}

func TestVarianceMatchesWorkedExample(t *testing.T) {
	got := Variance(20, 8, 5)
	want := 0.9474
	if !approx(got, want, 1e-3) {
		t.Errorf("Variance(20,8,5) = %v, want ~%v", got, want)
	}
}

func TestVarianceIsLessThanBinomialVariance(t *testing.T) {
	// The whole point of section 2: the finite population correction
	// makes the true variance strictly smaller than binomial's naive one,
	// for any sample that's a nontrivial fraction of a finite population.
	hyperVar := Variance(20, 8, 5)
	p := 8.0 / 20.0
	binomVar := 5 * p * (1 - p)
	if hyperVar >= binomVar {
		t.Errorf("hypergeometric variance %v should be < binomial variance %v", hyperVar, binomVar)
	}
}

func TestMeanEqualsBinomialMean(t *testing.T) {
	// Both distributions share the same mean, n*K/N == n*p -- only the
	// spread differs.
	hyperMean := Mean(20, 8, 5)
	binomMean := 5 * (8.0 / 20.0)
	if !approx(hyperMean, binomMean, 1e-9) {
		t.Errorf("hypergeometric mean %v should equal binomial mean %v", hyperMean, binomMean)
	}
}

func TestBinomialPMFWorkedExample(t *testing.T) {
	// Section 2's naive with-replacement comparison at p=0.4, k=5.
	got := BinomialPMF(5, 0.4, 5)
	if !approx(got, 0.0102, 1e-4) {
		t.Errorf("BinomialPMF(5,0.4,5) = %v, want ~0.0102", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("hypergeometric-distribution")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
