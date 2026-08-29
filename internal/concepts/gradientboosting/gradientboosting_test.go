package gradientboosting

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func approxEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestBoostZeroStagesIsJustTheMean(t *testing.T) {
	base, stumps := Boost(Xs, Ys, 0, 1.0)
	if !approxEqual(base, 5) {
		t.Errorf("base = %v, want 5", base)
	}
	if len(stumps) != 0 {
		t.Errorf("len(stumps) = %d, want 0", len(stumps))
	}
}

func TestBoostStage1MatchesWorkedExample(t *testing.T) {
	base, stumps := Boost(Xs, Ys, 1, 1.0)
	if len(stumps) != 1 {
		t.Fatalf("len(stumps) = %d, want 1", len(stumps))
	}
	s := stumps[0]
	if !approxEqual(s.Threshold, 3.5) {
		t.Errorf("stage1 threshold = %v, want 3.5", s.Threshold)
	}
	if want := -7.0 / 3.0; !approxEqual(s.LeftVal, want) {
		t.Errorf("stage1 left = %v, want %v", s.LeftVal, want)
	}
	if !approxEqual(s.RightVal, 3.5) {
		t.Errorf("stage1 right = %v, want 3.5", s.RightVal)
	}

	preds := make([]float64, len(Xs))
	for i, x := range Xs {
		preds[i] = Predict(base, stumps, 1.0, x)
	}
	if got := SSE(Ys, preds); !approxEqual(got, 7.0/6.0) {
		t.Errorf("stage1 SSE = %v, want %v", got, 7.0/6.0)
	}
}

func TestBoostStage2MatchesWorkedExample(t *testing.T) {
	base, stumps := Boost(Xs, Ys, 2, 1.0)
	if len(stumps) != 2 {
		t.Fatalf("len(stumps) = %d, want 2", len(stumps))
	}
	s := stumps[1]
	if !approxEqual(s.Threshold, 1.5) {
		t.Errorf("stage2 threshold = %v, want 1.5", s.Threshold)
	}
	if want := -2.0 / 3.0; !approxEqual(s.LeftVal, want) {
		t.Errorf("stage2 left = %v, want %v", s.LeftVal, want)
	}
	if want := 1.0 / 6.0; !approxEqual(s.RightVal, want) {
		t.Errorf("stage2 right = %v, want %v", s.RightVal, want)
	}

	if got := Predict(base, stumps, 1.0, 1); !approxEqual(got, 2) {
		t.Errorf("stage2 prediction at x=1 = %v, want 2 (exact fit)", got)
	}

	preds := make([]float64, len(Xs))
	for i, x := range Xs {
		preds[i] = Predict(base, stumps, 1.0, x)
	}
	if got := SSE(Ys, preds); !approxEqual(got, 11.0/18.0) {
		t.Errorf("stage2 SSE = %v, want %v", got, 11.0/18.0)
	}
}

func TestTrainingSSEShrinksMonotonicallyWithMoreStages(t *testing.T) {
	base, allStumps := Boost(Xs, Ys, 4, 1.0)
	prevSSE := math.Inf(1)
	for stage := 0; stage <= 4; stage++ {
		preds := make([]float64, len(Xs))
		for i, x := range Xs {
			preds[i] = Predict(base, allStumps[:stage], 1.0, x)
		}
		sse := SSE(Ys, preds)
		if sse > prevSSE+1e-9 {
			t.Errorf("stage %d SSE = %v, want <= previous stage's %v", stage, sse, prevSSE)
		}
		prevSSE = sse
	}
}

func TestLearningRateShrinksTheFirstStep(t *testing.T) {
	base, stumps := Boost(Xs, Ys, 1, 0.5)
	// The first stump is identical regardless of lr (it's fit against the
	// same base residual), but the applied correction is halved.
	want := base + 0.5*stumps[0].LeftVal
	if got := Predict(base, stumps, 0.5, Xs[0]); !approxEqual(got, want) {
		t.Errorf("Predict with lr=0.5 = %v, want %v", got, want)
	}
	full := base + stumps[0].LeftVal
	if got := Predict(base, stumps, 0.5, Xs[0]); approxEqual(got, full) {
		t.Errorf("Predict with lr=0.5 should differ from full-strength lr=1.0 prediction %v", full)
	}
}

func TestMeanAndSSEBasics(t *testing.T) {
	if got := Mean([]float64{2, 3, 3, 8, 9}); !approxEqual(got, 5) {
		t.Errorf("Mean = %v, want 5", got)
	}
	if got := SSE([]float64{2, 3}, []float64{2, 3}); !approxEqual(got, 0) {
		t.Errorf("SSE of identical slices = %v, want 0", got)
	}
	if got := SSE([]float64{0, 0}, []float64{1, 2}); !approxEqual(got, 5) {
		t.Errorf("SSE = %v, want 5", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("gradient-boosting")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
