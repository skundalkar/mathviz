package variancestddev

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestMeanVarianceStdDevKnownValues(t *testing.T) {
	// Hand-computed: mean=3, population variance=2, stddev=√2.
	xs := []float64{1, 2, 3, 4, 5}
	if got, want := Mean(xs), 3.0; math.Abs(got-want) > 1e-9 {
		t.Errorf("Mean = %v, want %v", got, want)
	}
	if got, want := Variance(xs), 2.0; math.Abs(got-want) > 1e-9 {
		t.Errorf("Variance = %v, want %v", got, want)
	}
	if got, want := StdDev(xs), math.Sqrt(2); math.Abs(got-want) > 1e-9 {
		t.Errorf("StdDev = %v, want %v", got, want)
	}
}

func TestEmptySliceIsZero(t *testing.T) {
	if got := Mean(nil); got != 0 {
		t.Errorf("Mean(nil) = %v, want 0", got)
	}
	if got := Variance(nil); got != 0 {
		t.Errorf("Variance(nil) = %v, want 0", got)
	}
	if got := StdDev(nil); got != 0 {
		t.Errorf("StdDev(nil) = %v, want 0", got)
	}
}

func TestStdDevIsSqrtOfVariance(t *testing.T) {
	xs := Sample(1.3, 2)
	if got, want := StdDev(xs), math.Sqrt(Variance(xs)); math.Abs(got-want) > 1e-12 {
		t.Errorf("StdDev = %v, want √Variance = %v", got, want)
	}
}

func TestSampleAtDefaultsMatchesBase(t *testing.T) {
	xs := Sample(1, 0)
	if len(xs) != len(baseSample) {
		t.Fatalf("got %d points, want %d", len(xs), len(baseSample))
	}
	for i, v := range baseSample {
		if xs[i] != v {
			t.Errorf("Sample(1,0)[%d] = %v, want %v", i, xs[i], v)
		}
	}
}

func TestVarianceScalesWithSpreadSquared(t *testing.T) {
	// Scaling every value by k scales variance by k².
	base := Variance(Sample(1, 0))
	for _, k := range []float64{0.5, 1.5, 2, 2.5} {
		got := Variance(Sample(k, 0))
		want := base * k * k
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("k=%v: Variance = %v, want %v (base*k²)", k, got, want)
		}
	}
}

func TestStdDevScalesLinearlyWithSpread(t *testing.T) {
	base := StdDev(Sample(1, 0))
	for _, k := range []float64{0.5, 1.5, 2, 2.5} {
		got := StdDev(Sample(k, 0))
		want := base * k
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("k=%v: StdDev = %v, want %v (base*k)", k, got, want)
		}
	}
}

func TestOutlierIncreasesVariance(t *testing.T) {
	var last float64 = -1
	for _, outlier := range []float64{0, 2, 4, 6, 8, 10} {
		v := Variance(Sample(1, outlier))
		if v < last-1e-9 {
			t.Errorf("outlier=%v: variance %v dropped below previous %v", outlier, v, last)
		}
		last = v
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("variance-vs-stddev")
	if !ok {
		t.Fatal("concept \"variance-vs-stddev\" not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
