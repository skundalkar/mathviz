package overfitting

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestPolyEvalKnownValues(t *testing.T) {
	// y = 2 + 3x
	coeffs := []float64{2, 3}
	if got := PolyEval(coeffs, 4); math.Abs(got-14) > 1e-9 {
		t.Fatalf("PolyEval = %v, want 14", got)
	}
	if got := PolyEval(coeffs, 0); math.Abs(got-2) > 1e-9 {
		t.Fatalf("PolyEval(0) = %v, want 2", got)
	}
}

func TestPolyFitRecoversExactLinear(t *testing.T) {
	// Noiseless points on y = 2 - 1.5x should fit a degree-1 poly exactly.
	xs := []float64{-2, -1, 0, 1, 2, 3}
	ys := make([]float64, len(xs))
	for i, x := range xs {
		ys[i] = 2 - 1.5*x
	}
	coeffs := PolyFit(xs, ys, 1)
	if len(coeffs) != 2 {
		t.Fatalf("expected 2 coefficients, got %d", len(coeffs))
	}
	if math.Abs(coeffs[0]-2) > 1e-6 || math.Abs(coeffs[1]-(-1.5)) > 1e-6 {
		t.Fatalf("coeffs = %v, want [2, -1.5]", coeffs)
	}
}

func TestPolyFitRecoversExactQuadratic(t *testing.T) {
	// Noiseless points on y = 1 - 2x + 0.5x^2 should fit a degree-2 poly exactly.
	xs := []float64{-3, -2, -1, 0, 1, 2, 3, 4}
	ys := make([]float64, len(xs))
	for i, x := range xs {
		ys[i] = 1 - 2*x + 0.5*x*x
	}
	coeffs := PolyFit(xs, ys, 2)
	want := []float64{1, -2, 0.5}
	for i, w := range want {
		if math.Abs(coeffs[i]-w) > 1e-6 {
			t.Fatalf("coeffs = %v, want %v", coeffs, want)
		}
	}
}

func TestTrainingMSENonIncreasingWithDegree(t *testing.T) {
	xs, ys := DataPoints(12, 0.4)
	prevMSE := math.Inf(1)
	for degree := 1; degree <= 10; degree++ {
		coeffs := PolyFit(xs, ys, degree)
		mse := MSE(xs, ys, coeffs)
		if mse > prevMSE+1e-9 {
			t.Fatalf("degree %d: training MSE %v rose above degree %d's %v", degree, mse, degree-1, prevMSE)
		}
		prevMSE = mse
	}
}

func TestHighDegreeOverfitsMoreThanLowDegree(t *testing.T) {
	// The core lesson: a high-degree fit should win on training error but
	// lose on true (generalization) error, once there's noise to chase.
	xs, ys := DataPoints(12, 0.5)
	low := PolyFit(xs, ys, 2)
	high := PolyFit(xs, ys, 11)

	if MSE(xs, ys, high) > MSE(xs, ys, low)+1e-9 {
		t.Fatalf("high-degree training MSE (%v) should be <= low-degree's (%v)",
			MSE(xs, ys, high), MSE(xs, ys, low))
	}
	if TrueError(high) <= TrueError(low) {
		t.Fatalf("high-degree true error (%v) should exceed low-degree's (%v) — that's the overfit",
			TrueError(high), TrueError(low))
	}
}

func TestDataPointsDeterministic(t *testing.T) {
	xs1, ys1 := DataPoints(10, 0.3)
	xs2, ys2 := DataPoints(10, 0.3)
	for i := range xs1 {
		if xs1[i] != xs2[i] || ys1[i] != ys2[i] {
			t.Fatal("DataPoints is not deterministic")
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("overfitting")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
