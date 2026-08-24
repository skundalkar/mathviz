package logreg

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestSigmoidKnownValues(t *testing.T) {
	if got := Sigmoid(0); !near(got, 0.5, 1e-9) {
		t.Errorf("Sigmoid(0) = %v, want 0.5", got)
	}
	if got := Sigmoid(100); !near(got, 1, 1e-9) {
		t.Errorf("Sigmoid(100) = %v, want ~1", got)
	}
	if got := Sigmoid(-100); !near(got, 0, 1e-9) {
		t.Errorf("Sigmoid(-100) = %v, want ~0", got)
	}
}

func TestPredictMatchesSigmoidOfLinearCombination(t *testing.T) {
	w, b, x := 2.0, -1.0, 3.0
	want := Sigmoid(w*x + b)
	if got := Predict(w, b, x); got != want {
		t.Errorf("Predict(%v,%v,%v) = %v, want %v", w, b, x, got, want)
	}
}

func TestLogLossZeroExamplesIsZero(t *testing.T) {
	if got := LogLoss(nil, nil, 1, 1); got != 0 {
		t.Errorf("LogLoss(nil,nil,...) = %v, want 0", got)
	}
}

func TestLogLossAtUntrainedStartingPoint(t *testing.T) {
	// At w=0,b=0 every prediction is sigmoid(0)=0.5, so every example costs
	// exactly -log(0.5), the untrained baseline named in LEARNINGS.md.
	got := LogLoss(HoursStudied, Passed, 0, 0)
	want := -math.Log(0.5)
	if !near(got, want, 1e-9) {
		t.Errorf("LogLoss at w=0,b=0 = %v, want %v", got, want)
	}
}

func TestLogLossPunishesConfidentWrongMoreThanCautious(t *testing.T) {
	xs, ys := []float64{0}, []float64{1}
	// w=0 forces Predict to sigmoid(b) regardless of x, letting b alone set
	// the (wrong, since the label is 1) predicted probability directly.
	confidentWrong := LogLoss(xs, ys, 0, -5) // predicts close to 0 for a true 1
	cautious := LogLoss(xs, ys, 0, 0)        // predicts 0.5
	if confidentWrong <= cautious {
		t.Errorf("confident-wrong loss %v should exceed cautious loss %v", confidentWrong, cautious)
	}
}

func TestFitLogisticRegressionReducesLoss(t *testing.T) {
	w, b := FitLogisticRegression(HoursStudied, Passed, 10000, 0.5)
	fitted := LogLoss(HoursStudied, Passed, w, b)
	baseline := LogLoss(HoursStudied, Passed, 0, 0)
	if fitted >= baseline {
		t.Errorf("fitted loss %v did not improve on baseline %v", fitted, baseline)
	}
	// Regression guard on the exact worked-example numbers in LEARNINGS.md.
	if !near(w, 2.59, 0.05) {
		t.Errorf("fitted w = %v, want ~2.59", w)
	}
	if !near(b, -8.43, 0.05) {
		t.Errorf("fitted b = %v, want ~-8.43", b)
	}
	if !near(fitted, 0.251, 0.01) {
		t.Errorf("fitted loss = %v, want ~0.251", fitted)
	}
}

func TestFitLogisticRegressionSeparatesTheData(t *testing.T) {
	// The fit should predict low probability early (mostly-fail region) and
	// high probability late (mostly-pass region), tracking the data's
	// overall trend even though two individual points buck it.
	w, b := FitLogisticRegression(HoursStudied, Passed, 10000, 0.5)
	if p := Predict(w, b, 0.5); p > 0.1 {
		t.Errorf("Predict at x=0.5 = %v, want small (early, all-fail region)", p)
	}
	if p := Predict(w, b, 5.0); p < 0.9 {
		t.Errorf("Predict at x=5.0 = %v, want large (late, all-pass region)", p)
	}
}

func TestDecisionBoundaryMatchesThresholdCrossing(t *testing.T) {
	w, b := 2.0, -6.0 // crosses p=0.5 at x=3
	x := DecisionBoundary(w, b, 0.5)
	if !near(x, 3.0, 1e-9) {
		t.Errorf("DecisionBoundary(%v,%v,0.5) = %v, want 3.0", w, b, x)
	}
	if p := Predict(w, b, x); !near(p, 0.5, 1e-9) {
		t.Errorf("Predict at the boundary = %v, want 0.5", p)
	}
}

func TestDecisionBoundaryFlatModelIsNaN(t *testing.T) {
	if x := DecisionBoundary(0, 1, 0.5); !math.IsNaN(x) {
		t.Errorf("DecisionBoundary with w=0 = %v, want NaN", x)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("logistic-regression")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
