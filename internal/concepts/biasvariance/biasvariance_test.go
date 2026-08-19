package biasvariance

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// worked example from LEARNINGS.md: x0=1.5, noise=0.4, numRuns=30.
const (
	worksedX0    = 1.5
	worksedNoise = 0.4
)

func TestTrueFunctionKnownValue(t *testing.T) {
	if got := TrueFunction(1.5); !almostEqual(got, 0.9974949866, 1e-9) {
		t.Errorf("TrueFunction(1.5) = %v, want ~0.9974949866", got)
	}
}

func TestDataPointsDifferentRunsDiffer(t *testing.T) {
	_, ys0 := DataPoints(numPoints, 0.4, 0)
	_, ys1 := DataPoints(numPoints, 0.4, 1)
	same := true
	for i := range ys0 {
		if ys0[i] != ys1[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("DataPoints with different runs produced identical y's, want them to differ")
	}
}

func TestDataPointsIsDeterministic(t *testing.T) {
	_, ysA := DataPoints(numPoints, 0.4, 5)
	_, ysB := DataPoints(numPoints, 0.4, 5)
	for i := range ysA {
		if ysA[i] != ysB[i] {
			t.Fatalf("DataPoints(_, _, 5) is not deterministic at index %d", i)
		}
	}
}

func TestPredictionsCountMatchesNumRuns(t *testing.T) {
	preds := Predictions(worksedX0, 3, worksedNoise)
	if len(preds) != numRuns {
		t.Errorf("len(Predictions(...)) = %d, want %d", len(preds), numRuns)
	}
}

// Known values from LEARNINGS.md's worked example: x0=1.5 (true value
// sin(1.5)~=0.9975), noise=0.4, numRuns=30.
func TestBiasAndVarianceKnownValues(t *testing.T) {
	trueVal := TrueFunction(worksedX0)
	cases := []struct {
		degree       int
		wantBias2    float64
		wantVariance float64
	}{
		{1, 0.4576082948, 0.0076283121},
		{3, 0.0025707953, 0.0187633171},
		{9, 0.0004202690, 0.0334014724},
	}
	for _, c := range cases {
		preds := Predictions(worksedX0, c.degree, worksedNoise)
		bias2 := BiasSquared(preds, trueVal)
		vari := Variance(preds)
		if !almostEqual(bias2, c.wantBias2, 1e-6) {
			t.Errorf("degree=%d: BiasSquared = %v, want ~%v", c.degree, bias2, c.wantBias2)
		}
		if !almostEqual(vari, c.wantVariance, 1e-6) {
			t.Errorf("degree=%d: Variance = %v, want ~%v", c.degree, vari, c.wantVariance)
		}
	}
}

func TestDegree3HasLowestTotalError(t *testing.T) {
	// The central claim this package exists to demonstrate: the moderate
	// degree beats both the too-simple and too-flexible extremes on total
	// expected error (bias^2 + variance), at the worked example's x0.
	trueVal := TrueFunction(worksedX0)
	totalFor := func(degree int) float64 {
		preds := Predictions(worksedX0, degree, worksedNoise)
		return BiasSquared(preds, trueVal) + Variance(preds)
	}
	total1, total3, total9 := totalFor(1), totalFor(3), totalFor(9)
	if total3 >= total1 {
		t.Errorf("degree=3 total error %v should be less than degree=1 total error %v", total3, total1)
	}
	if total3 >= total9 {
		t.Errorf("degree=3 total error %v should be less than degree=9 total error %v", total3, total9)
	}
}

func TestVarianceIsZeroForIdenticalPredictions(t *testing.T) {
	preds := []float64{0.5, 0.5, 0.5, 0.5}
	if got := Variance(preds); got != 0 {
		t.Errorf("Variance of identical predictions = %v, want 0", got)
	}
}

func TestBiasSquaredIsZeroWhenMeanMatchesTrueValue(t *testing.T) {
	preds := []float64{0.4, 0.6} // mean 0.5
	if got := BiasSquared(preds, 0.5); !almostEqual(got, 0, 1e-9) {
		t.Errorf("BiasSquared with mean == trueVal = %v, want 0", got)
	}
}

func TestPolyFitAndEvalRoundTrip(t *testing.T) {
	// A degree-1 fit to a perfectly linear, noiseless dataset should
	// reproduce that line exactly.
	xs := []float64{0, 1, 2, 3}
	ys := []float64{1, 3, 5, 7} // y = 1 + 2x
	coeffs := PolyFit(xs, ys, 1)
	for _, x := range xs {
		want := 1 + 2*x
		if got := PolyEval(coeffs, x); !almostEqual(got, want, 1e-6) {
			t.Errorf("PolyEval(coeffs, %v) = %v, want %v", x, got, want)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("bias-variance-tradeoff")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
