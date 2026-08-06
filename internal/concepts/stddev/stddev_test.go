package stddev

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestNormalPDFPeakAtMean(t *testing.T) {
	mu, sigma := 1.5, 0.8
	peak := NormalPDF(mu, mu, sigma)
	// Density anywhere else must be lower than at the mean.
	for _, x := range []float64{-1, 0, 1, 2, 3} {
		if NormalPDF(x, mu, sigma) > peak+1e-12 {
			t.Fatalf("density at %v exceeded peak", x)
		}
	}
	// Analytic peak value = 1/(σ√(2π)).
	want := 1 / (sigma * math.Sqrt(2*math.Pi))
	if math.Abs(peak-want) > 1e-9 {
		t.Fatalf("peak = %v, want %v", peak, want)
	}
}

func TestAreaWithinKnownSigmaCoverage(t *testing.T) {
	cases := []struct {
		k    float64
		want float64
	}{
		{1, 0.6827}, {2, 0.9545}, {3, 0.9973},
	}
	for _, tc := range cases {
		got := AreaWithin(tc.k)
		if math.Abs(got-tc.want) > 5e-4 {
			t.Errorf("AreaWithin(%v) = %.4f, want ≈ %.4f", tc.k, got, tc.want)
		}
	}
}

func TestNormalizesToOne(t *testing.T) {
	// Trapezoidal integral of the PDF over a wide window ≈ 1.
	mu, sigma := 0.3, 1.2
	lo, hi, n := -12.0, 12.0, 20000
	h := (hi - lo) / float64(n)
	sum := 0.0
	for i := 0; i <= n; i++ {
		x := lo + float64(i)*h
		w := 1.0
		if i == 0 || i == n {
			w = 0.5
		}
		sum += w * NormalPDF(x, mu, sigma)
	}
	if area := sum * h; math.Abs(area-1) > 1e-3 {
		t.Fatalf("integral = %v, want ≈ 1", area)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("standard-deviation")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
