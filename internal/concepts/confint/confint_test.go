package confint

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestStdNormalQuantileKnownValues(t *testing.T) {
	if got := StdNormalQuantile(0.5); math.Abs(got) > 1e-9 {
		t.Errorf("StdNormalQuantile(0.5) = %v, want 0", got)
	}
	if got, want := StdNormalQuantile(0.975), 1.959963985; math.Abs(got-want) > 1e-6 {
		t.Errorf("StdNormalQuantile(0.975) = %v, want %v", got, want)
	}
	for _, p := range []float64{0, 1, -0.1, 1.1} {
		if got := StdNormalQuantile(p); !math.IsNaN(got) {
			t.Errorf("StdNormalQuantile(%v) = %v, want NaN", p, got)
		}
	}
}

func TestCriticalZKnownValues(t *testing.T) {
	cases := map[float64]float64{
		0.90: 1.6448536,
		0.95: 1.9599640,
		0.99: 2.5758293,
	}
	for conf, want := range cases {
		if got := CriticalZ(conf); math.Abs(got-want) > 1e-5 {
			t.Errorf("CriticalZ(%v) = %v, want %v", conf, got, want)
		}
	}
}

func TestStandardErrorShrinksWithN(t *testing.T) {
	if got, want := StandardError(1, 4), 0.5; math.Abs(got-want) > 1e-12 {
		t.Errorf("StandardError(1,4) = %v, want %v", got, want)
	}
	se1, se4, se100 := StandardError(2, 1), StandardError(2, 4), StandardError(2, 100)
	if !(se1 > se4 && se4 > se100) {
		t.Fatalf("standard error should shrink as n grows: se1=%v se4=%v se100=%v", se1, se4, se100)
	}
	if got := StandardError(1, 0); got != 0 {
		t.Errorf("StandardError(1,0) = %v, want 0 (invalid n)", got)
	}
}

func TestSampleMeansAreSymmetricAroundTrueMean(t *testing.T) {
	means := SampleMeans(10)
	if len(means) != NumSamples {
		t.Fatalf("len(SampleMeans) = %d, want %d", len(means), NumSamples)
	}
	for i := 0; i < NumSamples/2; i++ {
		lo, hi := means[i], means[NumSamples-1-i]
		if math.Abs((lo+hi)/2-TrueMean) > 1e-9 {
			t.Errorf("means[%d]=%v and means[%d]=%v not symmetric about %v", i, lo, NumSamples-1-i, hi, TrueMean)
		}
	}
	// Strictly increasing: evenly spaced quantiles of an increasing CDF.
	for i := 1; i < len(means); i++ {
		if means[i] <= means[i-1] {
			t.Fatalf("SampleMeans not increasing at index %d: %v <= %v", i, means[i], means[i-1])
		}
	}
}

func TestCoverageCountKnownValues(t *testing.T) {
	// These land solidly inside grid cells (not on a quantile boundary), so
	// the count is exact and reproducible.
	cases := map[float64]int{
		0.5: 10,
		0.8: 16,
		0.9: 18,
	}
	for conf, want := range cases {
		if got := CoverageCount(conf, 20); got != want {
			t.Errorf("CoverageCount(%v, 20) = %d, want %d", conf, got, want)
		}
	}
}

func TestCoverageCountIsIndependentOfN(t *testing.T) {
	// Scaling n rescales every mean and interval by the same standard error,
	// so which grid points fall inside the band shouldn't change.
	for _, conf := range []float64{0.5, 0.8, 0.9} {
		want := CoverageCount(conf, 5)
		for _, n := range []int{2, 20, 100} {
			if got := CoverageCount(conf, n); got != want {
				t.Errorf("CoverageCount(%v, %d) = %d, want %d (n-invariant)", conf, n, got, want)
			}
		}
	}
}

func TestCoverageCountNondecreasingInConfidence(t *testing.T) {
	prev := -1
	for conf := 0.5; conf <= 0.99; conf += 0.05 {
		got := CoverageCount(conf, 20)
		if got < prev {
			t.Fatalf("coverage decreased at confidence=%v: %d < %d", conf, got, prev)
		}
		if got < 0 || got > NumSamples {
			t.Fatalf("coverage out of range at confidence=%v: %d", conf, got)
		}
		prev = got
	}
}

func TestIntervalCentersOnMean(t *testing.T) {
	lo, hi := Interval(2, 1.5, 0.5)
	if want := 2 - 1.5*0.5; math.Abs(lo-want) > 1e-12 {
		t.Errorf("Interval lo = %v, want %v", lo, want)
	}
	if want := 2 + 1.5*0.5; math.Abs(hi-want) > 1e-12 {
		t.Errorf("Interval hi = %v, want %v", hi, want)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("confidence-interval")
	if !ok {
		t.Fatal("concept \"confidence-interval\" not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
