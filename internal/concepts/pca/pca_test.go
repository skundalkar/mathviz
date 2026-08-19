package pca

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// worked example from LEARNINGS.md: hours studied vs. test score, same data
// as covariance's worked example.
var studyXs = []float64{1, 2, 3, 4, 5}
var studyYs = []float64{2, 4, 5, 4, 5}

func TestCovarianceMatrixKnownValues(t *testing.T) {
	varX, covXY, varY := CovarianceMatrix(studyXs, studyYs)
	if !almostEqual(varX, 2.0, 1e-9) {
		t.Errorf("varX = %v, want 2.0", varX)
	}
	if !almostEqual(covXY, 1.2, 1e-9) {
		t.Errorf("covXY = %v, want 1.2", covXY)
	}
	if !almostEqual(varY, 1.2, 1e-9) {
		t.Errorf("varY = %v, want 1.2", varY)
	}
}

func TestEigenvaluesKnownValues(t *testing.T) {
	// From LEARNINGS.md: a=2, b=1.2, d=1.2 -> lambda1~=2.865, lambda2~=0.335.
	l1, l2 := Eigenvalues(2, 1.2, 1.2)
	if !almostEqual(l1, 2.864911064, 1e-6) {
		t.Errorf("lambda1 = %v, want ~2.864911", l1)
	}
	if !almostEqual(l2, 0.335088936, 1e-6) {
		t.Errorf("lambda2 = %v, want ~0.335089", l2)
	}
}

func TestEigenvaluesSumToTraceAndProductToDeterminant(t *testing.T) {
	// Eigenvalue identities that hold for any symmetric 2x2 matrix: their
	// sum is the trace (a+d), their product is the determinant (ad-b^2).
	a, b, d := 2.0, 1.2, 1.2
	l1, l2 := Eigenvalues(a, b, d)
	if got, want := l1+l2, a+d; !almostEqual(got, want, 1e-9) {
		t.Errorf("lambda1+lambda2 = %v, want trace %v", got, want)
	}
	if got, want := l1*l2, a*d-b*b; !almostEqual(got, want, 1e-9) {
		t.Errorf("lambda1*lambda2 = %v, want determinant %v", got, want)
	}
}

func TestPrincipalAngleKnownValue(t *testing.T) {
	theta := PrincipalAngle(2, 1.2, 1.2)
	gotDeg := theta * 180 / math.Pi
	if !almostEqual(gotDeg, 35.7825, 1e-3) {
		t.Errorf("PrincipalAngle degrees = %v, want ~35.7825", gotDeg)
	}
}

func TestPrincipalAngleAxisAlignedWhenNoShear(t *testing.T) {
	// No off-diagonal term (b=0): the principal axis is exactly the x-axis
	// when a > d (more variance in x than y).
	theta := PrincipalAngle(3, 0, 1)
	if !almostEqual(theta, 0, 1e-9) {
		t.Errorf("PrincipalAngle(3,0,1) = %v radians, want 0", theta)
	}
}

func TestExplainedVarianceRatioKnownValue(t *testing.T) {
	l1, l2 := Eigenvalues(2, 1.2, 1.2)
	ratio := ExplainedVarianceRatio(l1, l2)
	if !almostEqual(ratio, 0.895285, 1e-5) {
		t.Errorf("ExplainedVarianceRatio = %v, want ~0.895285", ratio)
	}
}

func TestExplainedVarianceRatioIsotropicIsHalf(t *testing.T) {
	// A perfectly round cloud (equal eigenvalues) splits variance 50/50 --
	// neither direction explains more than the other.
	ratio := ExplainedVarianceRatio(1.5, 1.5)
	if !almostEqual(ratio, 0.5, 1e-9) {
		t.Errorf("ExplainedVarianceRatio(1.5,1.5) = %v, want 0.5", ratio)
	}
}

func TestExplainedVarianceRatioZeroTotalIsZero(t *testing.T) {
	if got := ExplainedVarianceRatio(0, 0); got != 0 {
		t.Errorf("ExplainedVarianceRatio(0,0) = %v, want 0", got)
	}
}

func TestGeneratePointsMatchesTargetCorrelation(t *testing.T) {
	for _, r := range []float64{-0.9, -0.5, 0, 0.5, 0.9} {
		xs, ys := GeneratePoints(r, 400)
		varX, covXY, varY := CovarianceMatrix(xs, ys)
		corr := covXY / math.Sqrt(varX*varY)
		if !almostEqual(corr, r, 0.15) {
			t.Errorf("GeneratePoints(%v, 400): sample correlation = %v, want close to %v", r, corr, r)
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

func TestScaleGrowsVarXByFactorSquared(t *testing.T) {
	// Variance scales with the SQUARE of a linear rescaling, unlike
	// covariance/standard-deviation-adjacent quantities that scale linearly.
	xs, _ := GeneratePoints(0.7, 200)
	baseVar := Covariance(xs, xs)
	scaled := Scale(xs, 3)
	scaledVar := Covariance(scaled, scaled)
	if !almostEqual(scaledVar, baseVar*9, 1e-9) {
		t.Errorf("Var(3x) = %v, want %v (9x base %v)", scaledVar, baseVar*9, baseVar)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("principal-component-analysis")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
