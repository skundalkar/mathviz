package zscore

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

const epsilon = 1e-9

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestZScoreWorkedExample(t *testing.T) {
	// From the opening scenario: Priya scores 90 on a test with mean 70,
	// stddev 10; Raj scores 92 on a test with mean 85, stddev 10.
	if got := ZScore(90, 70, 10); !almostEqual(got, 2.0, epsilon) {
		t.Errorf("ZScore(90, 70, 10) = %v, want 2.0", got)
	}
	if got := ZScore(92, 85, 10); !almostEqual(got, 0.7, epsilon) {
		t.Errorf("ZScore(92, 85, 10) = %v, want 0.7", got)
	}
}

func TestZScoreDifferingSpreadExample(t *testing.T) {
	// From the section-2 walkthrough: same raw gap (10 points above the
	// mean), different spreads, different z-scores.
	if got := ZScore(80, 70, 5); !almostEqual(got, 2.0, epsilon) {
		t.Errorf("ZScore(80, 70, 5) = %v, want 2.0", got)
	}
	if got := ZScore(80, 70, 20); !almostEqual(got, 0.5, epsilon) {
		t.Errorf("ZScore(80, 70, 20) = %v, want 0.5", got)
	}
}

func TestZScoreAtMeanIsZero(t *testing.T) {
	if got := ZScore(50, 50, 12); got != 0 {
		t.Errorf("ZScore(50, 50, 12) = %v, want 0", got)
	}
}

func TestZScoreBelowMeanIsNegative(t *testing.T) {
	if got := ZScore(60, 70, 10); !almostEqual(got, -1.0, epsilon) {
		t.Errorf("ZScore(60, 70, 10) = %v, want -1.0", got)
	}
}

func TestRawFromZIsZScoresInverse(t *testing.T) {
	cases := []struct{ x, mean, sd float64 }{
		{90, 70, 10}, {60, 70, 10}, {80, 70, 5}, {123.4, 50, 3.2},
	}
	for _, tc := range cases {
		z := ZScore(tc.x, tc.mean, tc.sd)
		got := RawFromZ(z, tc.mean, tc.sd)
		if !almostEqual(got, tc.x, 1e-6) {
			t.Errorf("RawFromZ(ZScore(%v, %v, %v), %v, %v) = %v, want %v",
				tc.x, tc.mean, tc.sd, tc.mean, tc.sd, got, tc.x)
		}
	}
}

func TestNormalPDFPeaksAtMean(t *testing.T) {
	// The density at the mean should be strictly higher than at any point
	// away from it, for a fixed stddev.
	peak := NormalPDF(70, 70, 10)
	for _, x := range []float64{50, 60, 80, 90, 100} {
		if NormalPDF(x, 70, 10) >= peak {
			t.Errorf("NormalPDF(%v, 70, 10) = %v, want < peak %v", x, NormalPDF(x, 70, 10), peak)
		}
	}
}

func TestNormalPDFIsSymmetric(t *testing.T) {
	left := NormalPDF(60, 70, 10)
	right := NormalPDF(80, 70, 10)
	if !almostEqual(left, right, 1e-9) {
		t.Errorf("NormalPDF(60,70,10) = %v, NormalPDF(80,70,10) = %v, want equal (symmetric about the mean)", left, right)
	}
}

func TestNormalPDFKnownPeakValue(t *testing.T) {
	// The peak of a standard normal (mean 0, stddev 1) is 1/sqrt(2*pi).
	want := 1 / math.Sqrt(2*math.Pi)
	if got := NormalPDF(0, 0, 1); !almostEqual(got, want, 1e-9) {
		t.Errorf("NormalPDF(0, 0, 1) = %v, want %v", got, want)
	}
}

func TestPercentileFromZKnownValues(t *testing.T) {
	if got := PercentileFromZ(0); !almostEqual(got, 50, 1e-9) {
		t.Errorf("PercentileFromZ(0) = %v, want 50", got)
	}
	// z=2.0 is the "why would you need this" section's payoff figure:
	// about the 97.7th percentile.
	if got := PercentileFromZ(2.0); !almostEqual(got, 97.72, 0.01) {
		t.Errorf("PercentileFromZ(2.0) = %v, want ~97.72", got)
	}
}

func TestPercentileFromZIsMonotonic(t *testing.T) {
	prev := PercentileFromZ(-4)
	for z := -3.9; z <= 4; z += 0.1 {
		cur := PercentileFromZ(z)
		if cur < prev {
			t.Fatalf("PercentileFromZ(%v) = %v, less than PercentileFromZ(%v) = %v — not monotonic", z, cur, z-0.1, prev)
		}
		prev = cur
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("z-score")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
