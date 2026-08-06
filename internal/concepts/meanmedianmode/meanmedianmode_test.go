package meanmedianmode

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestKnownValuesAtMuZeroSigmaOne(t *testing.T) {
	// Standard log-normal(0,1): well-known closed forms.
	if got, want := Median(0, 1), 1.0; math.Abs(got-want) > 1e-9 {
		t.Errorf("Median(0,1) = %v, want %v", got, want)
	}
	if got, want := Mean(0, 1), math.Exp(0.5); math.Abs(got-want) > 1e-9 {
		t.Errorf("Mean(0,1) = %v, want %v", got, want)
	}
	if got, want := Mode(0, 1), math.Exp(-1); math.Abs(got-want) > 1e-9 {
		t.Errorf("Mode(0,1) = %v, want %v", got, want)
	}
}

func TestOrderingForPositiveSkew(t *testing.T) {
	// For any sigma > 0, log-normal is right-skewed: mode < median < mean.
	for _, sigma := range []float64{0.1, 0.5, 1.0, 1.2} {
		mode, median, mean := Mode(0, sigma), Median(0, sigma), Mean(0, sigma)
		if !(mode < median && median < mean) {
			t.Errorf("sigma=%v: expected mode < median < mean, got %v < %v < %v", sigma, mode, median, mean)
		}
	}
}

func TestConvergeAsSkewShrinks(t *testing.T) {
	// As sigma -> 0, the log-normal collapses toward a spike at exp(mu): all
	// three statistics should converge to the same point.
	mu, sigma := 0.3, 1e-4
	mode, median, mean := Mode(mu, sigma), Median(mu, sigma), Mean(mu, sigma)
	want := math.Exp(mu)
	for name, got := range map[string]float64{"mode": mode, "median": median, "mean": mean} {
		if math.Abs(got-want) > 1e-3 {
			t.Errorf("%s = %v, want ≈ %v as sigma -> 0", name, got, want)
		}
	}
}

func TestPDFPeaksAtMode(t *testing.T) {
	mu, sigma := 0.2, 0.7
	mode := Mode(mu, sigma)
	peak := LogNormalPDF(mode, mu, sigma)
	for _, x := range []float64{0.05, 0.3, 0.6, 1.0, 2.0, 4.0} {
		if LogNormalPDF(x, mu, sigma) > peak+1e-9 {
			t.Errorf("density at x=%v exceeded density at mode (%v)", x, mode)
		}
	}
}

func TestLogNormalPDFNonPositiveXIsZero(t *testing.T) {
	if got := LogNormalPDF(0, 0, 1); got != 0 {
		t.Errorf("LogNormalPDF(0,...) = %v, want 0", got)
	}
	if got := LogNormalPDF(-1, 0, 1); got != 0 {
		t.Errorf("LogNormalPDF(-1,...) = %v, want 0", got)
	}
}

func TestLogNormalPDFIntegratesToOne(t *testing.T) {
	// Trapezoidal integral of the PDF over a wide window ≈ 1.
	mu, sigma := 0.1, 0.5
	lo, hi, n := 0.001, 20.0, 20000
	h := (hi - lo) / float64(n)
	sum := 0.0
	for i := 0; i <= n; i++ {
		x := lo + float64(i)*h
		w := 1.0
		if i == 0 || i == n {
			w = 0.5
		}
		sum += w * LogNormalPDF(x, mu, sigma)
	}
	if area := sum * h; math.Abs(area-1) > 1e-2 {
		t.Fatalf("integral = %v, want ≈ 1", area)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("mean-median-mode")
	if !ok {
		t.Fatal("concept \"mean-median-mode\" not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
