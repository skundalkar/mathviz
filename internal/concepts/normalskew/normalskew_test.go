package normalskew

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestZeroSkewZeroKurtIsStandardNormal(t *testing.T) {
	for _, x := range []float64{-3, -1.5, -0.5, 0, 0.5, 1.5, 3} {
		got := GramCharlierPDF(x, 0, 0)
		want := StdNormalPDF(x)
		if math.Abs(got-want) > 1e-12 {
			t.Errorf("GramCharlierPDF(%v,0,0) = %v, want %v", x, got, want)
		}
	}
}

func TestHermitePolynomialsKnownValues(t *testing.T) {
	if got, want := Hermite3(0), 0.0; got != want {
		t.Errorf("Hermite3(0) = %v, want %v", got, want)
	}
	if got, want := Hermite4(0), 3.0; got != want {
		t.Errorf("Hermite4(0) = %v, want %v", got, want)
	}
	// He3(√3) = 3√3 - 3√3 = 0.
	if got := Hermite3(math.Sqrt(3)); math.Abs(got) > 1e-9 {
		t.Errorf("Hermite3(√3) = %v, want 0", got)
	}
}

func TestPositiveKurtosisRaisesPeakDensity(t *testing.T) {
	// At x=0, He3(0)=0 and He4(0)=3, so density = φ(0)*(1 + kurt/8):
	// positive excess kurtosis makes the peak taller (leptokurtic).
	peak0 := GramCharlierPDF(0, 0, 0)
	peak1 := GramCharlierPDF(0, 0, 1)
	if peak1 <= peak0 {
		t.Errorf("peak at kurt=1 (%v) should exceed peak at kurt=0 (%v)", peak1, peak0)
	}
	want := StdNormalPDF(0) * (1 + 1.0/8)
	if math.Abs(peak1-want) > 1e-9 {
		t.Errorf("GramCharlierPDF(0,0,1) = %v, want %v", peak1, want)
	}
}

func TestKurtosisAloneStaysSymmetric(t *testing.T) {
	// With skew=0, only the even Hermite4 term contributes, so the curve must
	// stay symmetric about zero.
	for _, x := range []float64{0.3, 1.1, 2.4} {
		for _, kurt := range []float64{-0.8, 0.5, 1.5} {
			left, right := GramCharlierPDF(-x, 0, kurt), GramCharlierPDF(x, 0, kurt)
			if math.Abs(left-right) > 1e-12 {
				t.Errorf("kurt=%v: PDF(-%v)=%v != PDF(%v)=%v", kurt, x, left, x, right)
			}
		}
	}
}

func TestDensityNeverNegative(t *testing.T) {
	for _, skew := range []float64{-0.9, -0.4, 0, 0.4, 0.9} {
		for _, kurt := range []float64{-1, -0.3, 0, 1, 2} {
			for x := -6.0; x <= 6.0; x += 0.25 {
				if got := GramCharlierPDF(x, skew, kurt); got < 0 {
					t.Fatalf("GramCharlierPDF(%v,%v,%v) = %v, negative", x, skew, kurt, got)
				}
			}
		}
	}
}

func TestIntegratesApproximatelyToOne(t *testing.T) {
	// For modest skew/kurt the expansion stays non-negative across the bulk
	// of the range, so the integral should still be close to 1.
	skew, kurt := 0.4, 0.5
	lo, hi, n := -10.0, 10.0, 20000
	h := (hi - lo) / float64(n)
	sum := 0.0
	for i := 0; i <= n; i++ {
		x := lo + float64(i)*h
		w := 1.0
		if i == 0 || i == n {
			w = 0.5
		}
		sum += w * GramCharlierPDF(x, skew, kurt)
	}
	if area := sum * h; math.Abs(area-1) > 3e-2 {
		t.Fatalf("integral = %v, want ≈ 1", area)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("normal-vs-skew")
	if !ok {
		t.Fatal("concept \"normal-vs-skew\" not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
