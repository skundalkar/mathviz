package randomforest

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func approx(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestForestFirstTreeMatchesDecisionTreesOwnSplit(t *testing.T) {
	trees := Forest(1)
	if len(trees) != 1 {
		t.Fatalf("Forest(1) has %d trees, want 1", len(trees))
	}
	tr := trees[0]
	if !approx(tr.Threshold, 2.75, 1e-9) {
		t.Errorf("tree 0 threshold = %v, want 2.75 (decision-trees's own best split)", tr.Threshold)
	}
	if !approx(tr.Gain, 0.610, 5e-3) {
		t.Errorf("tree 0 gain = %v, want ~0.610", tr.Gain)
	}
	if tr.LeftLabel != 0 || tr.RightLabel != 1 {
		t.Errorf("tree 0 labels = (%d,%d), want (0,1)", tr.LeftLabel, tr.RightLabel)
	}
}

func TestForestThresholdsWorkedExample(t *testing.T) {
	trees := Forest(9)
	want := []float64{2.75, 2.75, 3.75, 3.25, 2.75, 2.75, 3.75, 2.75, 2.75}
	if len(trees) != len(want) {
		t.Fatalf("Forest(9) has %d trees, want %d", len(trees), len(want))
	}
	for i, tr := range trees {
		if !approx(tr.Threshold, want[i], 1e-9) {
			t.Errorf("tree %d threshold = %v, want %v", i, tr.Threshold, want[i])
		}
	}
}

func TestForestGainsWorkedExample(t *testing.T) {
	trees := Forest(9)
	wantGains := []float64{0.610, 0.971, 0.881, 1.000, 0.396, 0.881, 0.610, 0.610, 0.610}
	for i, tr := range trees {
		if !approx(tr.Gain, wantGains[i], 5e-3) {
			t.Errorf("tree %d gain = %v, want ~%v", i, tr.Gain, wantGains[i])
		}
	}
}

func TestVoteFractionWorkedExample(t *testing.T) {
	cases := []struct {
		numTrees int
		hours    float64
		want     float64
	}{
		{1, 3.0, 1.0},
		{1, 2.5, 0.0},
		{3, 3.0, 2.0 / 3.0},
		{3, 3.5, 2.0 / 3.0},
		{3, 4.0, 1.0},
		{9, 3.0, 6.0 / 9.0},
		{9, 3.5, 7.0 / 9.0},
		{9, 4.0, 1.0},
	}
	for _, c := range cases {
		f := Forest(c.numTrees)
		if got := VoteFraction(f, c.hours); !approx(got, c.want, 5e-3) {
			t.Errorf("VoteFraction(Forest(%d), %v) = %v, want %v", c.numTrees, c.hours, got, c.want)
		}
	}
}

func TestMoreTreesNeverDecreaseThresholdGranularity(t *testing.T) {
	// Adding trees can only add vote-fraction levels between 0 and 1, never
	// remove one that a smaller forest already showed.
	f3 := Forest(3)
	f9 := Forest(9)
	if VoteFraction(f3, 4.0) != VoteFraction(f9, 4.0) {
		t.Errorf("both forests should be unanimous by hours=4.0")
	}
}

func TestForestClampsSizeToAvailableTrees(t *testing.T) {
	f := Forest(100)
	if len(f) != len(treeSamples) {
		t.Errorf("Forest(100) has %d trees, want %d (clamped)", len(f), len(treeSamples))
	}
	f0 := Forest(0)
	if len(f0) != 1 {
		t.Errorf("Forest(0) has %d trees, want 1 (clamped)", len(f0))
	}
}

func TestPredictUsesThresholdSide(t *testing.T) {
	tr := Tree{Threshold: 3.0, LeftLabel: 0, RightLabel: 1}
	if Predict(tr, 2.9) != 0 {
		t.Errorf("Predict below threshold should return LeftLabel")
	}
	if Predict(tr, 3.0) != 0 {
		t.Errorf("Predict at threshold should return LeftLabel (<=)")
	}
	if Predict(tr, 3.1) != 1 {
		t.Errorf("Predict above threshold should return RightLabel")
	}
}

func TestVoteFractionEmptyForest(t *testing.T) {
	if got := VoteFraction(nil, 3.0); got != 0 {
		t.Errorf("VoteFraction(nil, 3.0) = %v, want 0", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("random-forest")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
