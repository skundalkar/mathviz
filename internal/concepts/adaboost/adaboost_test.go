package adaboost

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestAlphaKnownValues(t *testing.T) {
	cases := []struct {
		err  float64
		want float64
	}{
		{0.25, 0.549306},      // 0.5*ln(3)
		{1.0 / 6.0, 0.804719}, // 0.5*ln(5)
		{0.2, 0.693147},       // 0.5*ln(4) == ln(2)
		{0.5, 0},              // a coin-flip stump earns zero vote
	}
	for _, c := range cases {
		if got := Alpha(c.err); !near(got, c.want, 1e-5) {
			t.Errorf("Alpha(%v) = %v, want ~%v", c.err, got, c.want)
		}
	}
}

func TestAlphaGrowsAsErrorShrinks(t *testing.T) {
	if !(Alpha(0.05) > Alpha(0.25) && Alpha(0.25) > Alpha(0.45)) {
		t.Errorf("Alpha should strictly decrease as err grows toward 0.5")
	}
}

func TestWeightedErrorKnownValue(t *testing.T) {
	weights := []float64{0.125, 0.125, 0.125, 0.125, 0.125, 0.125, 0.125, 0.125}
	st := Stump{Threshold: 2.5, LeftLabel: 1, RightLabel: -1}
	err := WeightedError(Ys, predictAll(st, Xs), weights)
	if !near(err, 0.25, 1e-9) {
		t.Errorf("WeightedError at threshold=2.5 = %v, want 0.25 (x=5,6 misclassified)", err)
	}
}

func TestFitWeightedStumpRoundOneFindsThreshold25(t *testing.T) {
	uniform := []float64{0.125, 0.125, 0.125, 0.125, 0.125, 0.125, 0.125, 0.125}
	st, err := FitWeightedStump(Xs, Ys, uniform)
	if !near(st.Threshold, 2.5, 1e-9) {
		t.Errorf("round-1 best threshold = %v, want 2.5", st.Threshold)
	}
	if st.LeftLabel != 1 || st.RightLabel != -1 {
		t.Errorf("round-1 stump labels = (%d,%d), want (1,-1)", st.LeftLabel, st.RightLabel)
	}
	if !near(err, 0.25, 1e-9) {
		t.Errorf("round-1 weighted error = %v, want 0.25", err)
	}
}

func TestRunThreeRoundsMatchesWorkedExample(t *testing.T) {
	rounds, weightsHistory := Run(Xs, Ys, 3)
	if len(rounds) != 3 || len(weightsHistory) != 4 {
		t.Fatalf("Run(3) returned %d rounds and %d weight snapshots, want 3 and 4", len(rounds), len(weightsHistory))
	}

	wantThresholds := []float64{2.5, 6.5, 4.5}
	wantErrs := []float64{0.25, 1.0 / 6.0, 0.2}
	wantAlphas := []float64{0.549306, 0.804719, 0.693147}
	for i, r := range rounds {
		if !near(r.Stump.Threshold, wantThresholds[i], 1e-9) {
			t.Errorf("round %d threshold = %v, want %v", i+1, r.Stump.Threshold, wantThresholds[i])
		}
		if !near(r.Err, wantErrs[i], 1e-9) {
			t.Errorf("round %d err = %v, want %v", i+1, r.Err, wantErrs[i])
		}
		if !near(r.Alpha, wantAlphas[i], 1e-5) {
			t.Errorf("round %d alpha = %v, want %v", i+1, r.Alpha, wantAlphas[i])
		}
	}

	// After round 1, the two points it got wrong (x=5, x=6) should have
	// doubled their starting weight of 1/8 to 1/4, while every other point
	// shrinks to 1/12.
	after1 := weightsHistory[1]
	for i, x := range Xs {
		want := 1.0 / 12.0
		if x == 5 || x == 6 {
			want = 0.25
		}
		if !near(after1[i], want, 1e-9) {
			t.Errorf("weight of x=%v after round 1 = %v, want %v", x, after1[i], want)
		}
	}
}

func TestPredictSingleStumpImperfect(t *testing.T) {
	rounds, _ := Run(Xs, Ys, 1)
	wrong := 0
	for i, x := range Xs {
		if Predict(rounds, x) != Ys[i] {
			wrong++
		}
	}
	if wrong != 2 {
		t.Errorf("round-1-only ensemble got %d wrong, want 2 (x=5, x=6)", wrong)
	}
}

func TestPredictThreeRoundsPerfect(t *testing.T) {
	rounds, _ := Run(Xs, Ys, 3)
	for i, x := range Xs {
		if got := Predict(rounds, x); got != Ys[i] {
			t.Errorf("3-round ensemble Predict(%v) = %d, want %d", x, got, Ys[i])
		}
	}
}

func TestUpdateWeightsSumsToOne(t *testing.T) {
	weights := []float64{0.125, 0.125, 0.125, 0.125, 0.125, 0.125, 0.125, 0.125}
	st := Stump{Threshold: 2.5, LeftLabel: 1, RightLabel: -1}
	preds := predictAll(st, Xs)
	alpha := Alpha(WeightedError(Ys, preds, weights))
	updated := UpdateWeights(Ys, preds, weights, alpha)
	var sum float64
	for _, w := range updated {
		sum += w
	}
	if !near(sum, 1.0, 1e-9) {
		t.Errorf("updated weights sum to %v, want 1.0", sum)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("adaboost")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
