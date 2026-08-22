package regularization

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

// A tiny, hand-verifiable single-feature dataset (x: 1,2,3,4; y: 3,5,7,8),
// independent of the concept's own 6-feature synthetic dataset. Every
// expected value below was derived by hand from the least-squares formulas
// and cross-checked against a standalone reference implementation.
var (
	tinyX = [][]float64{{1}, {2}, {3}, {4}}
	tinyY = []float64{3, 5, 7, 8}
)

func TestRidgeFitAtZeroLambdaMatchesOLS(t *testing.T) {
	// x̄=2.5, ȳ=5.75, Σ(x-x̄)(y-ȳ)=8.5, Σ(x-x̄)²=5 -> slope=1.7, intercept=1.5.
	got := RidgeFit(tinyX, tinyY, 0)
	if !near(got[0], 1.5, 1e-6) || !near(got[1], 1.7, 1e-6) {
		t.Errorf("RidgeFit(lambda=0) = %v, want [1.5, 1.7]", got)
	}
}

func TestRidgeFitShrinksTowardZero(t *testing.T) {
	// (XtX + 5I)c = Xty on the tiny dataset solves to intercept=3.625,
	// slope=0.85 -- worked by hand in LEARNINGS.md.
	got := RidgeFit(tinyX, tinyY, 5)
	if !near(got[0], 3.625, 1e-6) || !near(got[1], 0.85, 1e-6) {
		t.Errorf("RidgeFit(lambda=5) = %v, want [3.625, 0.85]", got)
	}
	// Larger lambda must not increase the slope's magnitude.
	got2 := RidgeFit(tinyX, tinyY, 20)
	if math.Abs(got2[1]) > math.Abs(got[1]) {
		t.Errorf("RidgeFit(lambda=20) slope %v should be smaller in magnitude than lambda=5 slope %v", got2[1], got[1])
	}
}

func TestLassoFitAtZeroLambdaMatchesOLS(t *testing.T) {
	got := LassoFit(tinyX, tinyY, 0)
	if !near(got[0], 1.5, 1e-4) || !near(got[1], 1.7, 1e-4) {
		t.Errorf("LassoFit(lambda=0) = %v, want ~[1.5, 1.7]", got)
	}
}

func TestLassoFitSoftThresholdsKnownValue(t *testing.T) {
	// soft(8.5, 3) = 5.5 -> slope = 5.5/5 = 1.1, intercept = 5.75-1.1*2.5 = 3.0.
	got := LassoFit(tinyX, tinyY, 3)
	if !near(got[0], 3.0, 1e-4) || !near(got[1], 1.1, 1e-4) {
		t.Errorf("LassoFit(lambda=3) = %v, want ~[3.0, 1.1]", got)
	}
}

func TestLassoFitZeroesCoefficientPastThreshold(t *testing.T) {
	// |Sxy|=8.5, so any lambda > 8.5 soft-thresholds the single feature to
	// exactly 0, leaving the intercept-only model (predict the mean, 5.75).
	got := LassoFit(tinyX, tinyY, 10)
	if got[1] != 0 {
		t.Errorf("LassoFit(lambda=10) slope = %v, want exactly 0", got[1])
	}
	if !near(got[0], 5.75, 1e-6) {
		t.Errorf("LassoFit(lambda=10) intercept = %v, want 5.75 (the mean)", got[0])
	}
}

func TestSoftThreshold(t *testing.T) {
	cases := []struct{ z, gamma, want float64 }{
		{10, 3, 7},
		{-10, 3, -7},
		{2, 3, 0},
		{-2, 3, 0},
		{3, 3, 0},
	}
	for _, c := range cases {
		if got := softThreshold(c.z, c.gamma); got != c.want {
			t.Errorf("softThreshold(%v, %v) = %v, want %v", c.z, c.gamma, got, c.want)
		}
	}
}

func TestOLSOverfitsNoiseFeatures(t *testing.T) {
	// On the concept's own dataset, unpenalized OLS should recover the two
	// real coefficients reasonably closely while still assigning some
	// nonzero weight to every one of the 4 pure-noise features.
	X := Features(numSamples)
	y := Targets(X, 0.6)
	ols := RidgeFit(X, y, 0)

	if !near(ols[1], TrueCoeffs[0], 0.2) {
		t.Errorf("OLS x1 coefficient = %v, want close to true %v", ols[1], TrueCoeffs[0])
	}
	if !near(ols[2], TrueCoeffs[1], 0.2) {
		t.Errorf("OLS x2 coefficient = %v, want close to true %v", ols[2], TrueCoeffs[1])
	}
	for i := 3; i <= 6; i++ {
		if ols[i] == 0 {
			t.Errorf("OLS coefficient x%d = 0 exactly; expected some spurious nonzero weight", i)
		}
	}
}

func TestLassoRecoversSparsityAtModerateLambda(t *testing.T) {
	X := Features(numSamples)
	y := Targets(X, 0.6)
	lasso := LassoFit(X, y, 2.0)

	if zeros := CountNearZero(lasso, 3, 1e-9); zeros != 4 {
		t.Errorf("Lasso(lambda=2) zeroed %d of the 4 noise coefficients, want 4", zeros)
	}
	// The two real coefficients must survive, staying reasonably close to
	// the true values (3, -2).
	if !near(lasso[1], 3, 0.3) || !near(lasso[2], -2, 0.3) {
		t.Errorf("Lasso(lambda=2) signal coefficients = [%v, %v], want close to [3, -2]", lasso[1], lasso[2])
	}
}

func TestRidgeNeverExactlyZeroesNoiseFeatures(t *testing.T) {
	X := Features(numSamples)
	y := Targets(X, 0.6)
	ridge := RidgeFit(X, y, 2.0)
	if zeros := CountNearZero(ridge, 3, 1e-9); zeros != 0 {
		t.Errorf("Ridge(lambda=2) zeroed %d noise coefficients exactly, want 0 (ridge only shrinks)", zeros)
	}
}

func TestL2PenaltyShrinksAsLambdaGrows(t *testing.T) {
	X := Features(numSamples)
	y := Targets(X, 0.6)
	prev := L2Penalty(RidgeFit(X, y, 0))
	for _, lam := range []float64{0.5, 1, 2, 4, 8} {
		cur := L2Penalty(RidgeFit(X, y, lam))
		if cur > prev+1e-9 {
			t.Errorf("L2Penalty at lambda=%v (%v) exceeds the penalty at the previous, smaller lambda (%v); ridge's total penalty must not increase with lambda", lam, cur, prev)
		}
		prev = cur
	}
}

func TestPenaltyFunctionsExcludeIntercept(t *testing.T) {
	coeffs := []float64{100, 3, -4} // huge intercept, small real coefficients
	if got := L1Penalty(coeffs); !near(got, 7, 1e-9) {
		t.Errorf("L1Penalty(%v) = %v, want 7 (intercept excluded)", coeffs, got)
	}
	if got := L2Penalty(coeffs); !near(got, 25, 1e-9) {
		t.Errorf("L2Penalty(%v) = %v, want 25 (intercept excluded)", coeffs, got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("regularization")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
