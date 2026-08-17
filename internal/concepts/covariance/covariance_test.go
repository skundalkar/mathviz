package covariance

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// worked example from LEARNINGS.md: hours studied vs. test score.
var studyXs = []float64{1, 2, 3, 4, 5}
var studyYs = []float64{2, 4, 5, 4, 5}

func TestCovarianceKnownValue(t *testing.T) {
	if got := Covariance(studyXs, studyYs); !almostEqual(got, 1.2, 1e-9) {
		t.Errorf("Covariance(studyXs, studyYs) = %v, want 1.2", got)
	}
}

func TestCovarianceIsSymmetric(t *testing.T) {
	if got, want := Covariance(studyXs, studyYs), Covariance(studyYs, studyXs); !almostEqual(got, want, 1e-9) {
		t.Errorf("Covariance(x,y) = %v, Covariance(y,x) = %v, want equal", got, want)
	}
}

func TestCovarianceWithSelfIsVariance(t *testing.T) {
	// Covariance(x,x) should equal the population variance of x: 2.0 for
	// studyXs (mean 3, squared deviations 4,1,0,1,4 -> average 2).
	if got := Covariance(studyXs, studyXs); !almostEqual(got, 2.0, 1e-9) {
		t.Errorf("Covariance(studyXs, studyXs) = %v, want 2.0", got)
	}
}

func TestCovarianceEmptyOrMismatched(t *testing.T) {
	if got := Covariance(nil, nil); got != 0 {
		t.Errorf("Covariance(nil,nil) = %v, want 0", got)
	}
	if got := Covariance([]float64{1, 2}, []float64{1}); got != 0 {
		t.Errorf("Covariance of mismatched lengths = %v, want 0", got)
	}
}

func TestScaleScalesCovarianceProportionally(t *testing.T) {
	// Relabeling x's units by a factor of 10 should scale covariance by
	// exactly 10, per the worked example in LEARNINGS.md.
	scaledXs := Scale(studyXs, 10)
	base := Covariance(studyXs, studyYs)
	scaled := Covariance(scaledXs, studyYs)
	if !almostEqual(scaled, base*10, 1e-9) {
		t.Errorf("Covariance(10x, y) = %v, want %v (10x base %v)", scaled, base*10, base)
	}
}

func TestScaleLeavesUnscaledInputUnmodified(t *testing.T) {
	orig := []float64{1, 2, 3}
	_ = Scale(orig, 5)
	want := []float64{1, 2, 3}
	for i := range orig {
		if orig[i] != want[i] {
			t.Errorf("Scale mutated its input: orig = %v, want %v", orig, want)
		}
	}
}

func TestStdDevKnownValues(t *testing.T) {
	if got := StdDev(studyXs); !almostEqual(got, math.Sqrt(2.0), 1e-9) {
		t.Errorf("StdDev(studyXs) = %v, want %v", got, math.Sqrt(2.0))
	}
	// studyYs: mean 4, squared deviations 4,0,1,0,1 -> average 1.2.
	if got := StdDev(studyYs); !almostEqual(got, math.Sqrt(1.2), 1e-9) {
		t.Errorf("StdDev(studyYs) = %v, want %v", got, math.Sqrt(1.2))
	}
}

func TestCorrelationKnownValue(t *testing.T) {
	if got := Correlation(studyXs, studyYs); !almostEqual(got, 0.7745966692414833, 1e-9) {
		t.Errorf("Correlation(studyXs, studyYs) = %v, want ~0.7746", got)
	}
}

func TestCorrelationIsInvariantToScaling(t *testing.T) {
	// The central claim this package exists to demonstrate: correlation
	// stays fixed under rescaling, even though covariance moves.
	base := Correlation(studyXs, studyYs)
	for _, factor := range []float64{0.1, 2, 10, 100} {
		scaled := Correlation(Scale(studyXs, factor), studyYs)
		if !almostEqual(scaled, base, 1e-9) {
			t.Errorf("Correlation(%vx, y) = %v, want unchanged %v", factor, scaled, base)
		}
	}
}

func TestCorrelationZeroVarianceIsZero(t *testing.T) {
	constant := []float64{5, 5, 5, 5}
	if got := Correlation(constant, studyYs[:4]); got != 0 {
		t.Errorf("Correlation with zero-variance series = %v, want 0", got)
	}
}

func TestGeneratePointsMatchesTargetCorrelation(t *testing.T) {
	for _, r := range []float64{-0.9, -0.5, 0, 0.5, 0.9} {
		xs, ys := GeneratePoints(r, 400)
		got := Correlation(xs, ys)
		if !almostEqual(got, r, 0.15) {
			t.Errorf("GeneratePoints(%v, 400): sample correlation = %v, want close to %v", r, got, r)
		}
	}
}

func TestGeneratePointsIsDeterministic(t *testing.T) {
	xs1, ys1 := GeneratePoints(0.6, 50)
	xs2, ys2 := GeneratePoints(0.6, 50)
	for i := range xs1 {
		if xs1[i] != xs2[i] || ys1[i] != ys2[i] {
			t.Fatalf("GeneratePoints(0.6, 50) is not deterministic at index %d", i)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("covariance")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
