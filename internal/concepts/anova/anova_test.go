package anova

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func approx(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestGroupSamplesCentersOnMean(t *testing.T) {
	xs := GroupSamples(70, 3)
	want := []float64{64, 67, 70, 73, 76}
	if len(xs) != len(want) {
		t.Fatalf("GroupSamples(70,3) has %d samples, want %d", len(xs), len(want))
	}
	for i := range want {
		if !approx(xs[i], want[i], 1e-9) {
			t.Errorf("GroupSamples(70,3)[%d] = %v, want %v", i, xs[i], want[i])
		}
	}
	if m := Mean(xs); !approx(m, 70, 1e-9) {
		t.Errorf("Mean(GroupSamples(70,3)) = %v, want 70", m)
	}
}

func TestWorkedExampleThreeGroups(t *testing.T) {
	groups := [][]float64{
		GroupSamples(70, 3), // A
		GroupSamples(75, 3), // B
		GroupSamples(82, 3), // C
	}
	r := Run(groups)

	if !approx(r.SSBetween, 363.333, 5e-3) {
		t.Errorf("SSBetween = %v, want ~363.333", r.SSBetween)
	}
	if !approx(r.SSWithin, 270, 1e-6) {
		t.Errorf("SSWithin = %v, want 270", r.SSWithin)
	}
	if r.DFBetween != 2 {
		t.Errorf("DFBetween = %d, want 2", r.DFBetween)
	}
	if r.DFWithin != 12 {
		t.Errorf("DFWithin = %d, want 12", r.DFWithin)
	}
	if !approx(r.MSBetween, 181.667, 5e-3) {
		t.Errorf("MSBetween = %v, want ~181.667", r.MSBetween)
	}
	if !approx(r.MSWithin, 22.5, 1e-6) {
		t.Errorf("MSWithin = %v, want 22.5", r.MSWithin)
	}
	if !approx(r.F, 8.074, 5e-3) {
		t.Errorf("F = %v, want ~8.074", r.F)
	}
}

func TestSSTotalEqualsSSBetweenPlusSSWithin(t *testing.T) {
	// The sum-of-squares decomposition SST = SSB + SSW should hold for any
	// set of groups -- check it directly against the total sum of squares
	// computed from the grand mean.
	groups := [][]float64{
		GroupSamples(10, 2),
		GroupSamples(20, 1),
		GroupSamples(15, 4),
		GroupSamples(30, 0.5),
	}
	var pooled []float64
	for _, g := range groups {
		pooled = append(pooled, g...)
	}
	grand := Mean(pooled)
	sst := 0.0
	for _, x := range pooled {
		d := x - grand
		sst += d * d
	}

	ssb, ssw := SSBetween(groups), SSWithin(groups)
	if !approx(ssb+ssw, sst, 1e-6) {
		t.Errorf("SSBetween+SSWithin = %v, want SSTotal = %v", ssb+ssw, sst)
	}
}

func TestFIsZeroWhenAllGroupMeansEqual(t *testing.T) {
	groups := [][]float64{
		GroupSamples(50, 2),
		GroupSamples(50, 2),
		GroupSamples(50, 2),
	}
	r := Run(groups)
	if !approx(r.SSBetween, 0, 1e-9) {
		t.Errorf("SSBetween = %v, want 0 (identical group means)", r.SSBetween)
	}
	if !approx(r.F, 0, 1e-9) {
		t.Errorf("F = %v, want 0 (identical group means)", r.F)
	}
}

func TestFIsInfiniteWhenNoWithinGroupSpreadButMeansDiffer(t *testing.T) {
	groups := [][]float64{
		GroupSamples(10, 0),
		GroupSamples(20, 0),
		GroupSamples(30, 0),
	}
	r := Run(groups)
	if r.SSWithin != 0 {
		t.Errorf("SSWithin = %v, want 0 (zero spread)", r.SSWithin)
	}
	if !math.IsInf(r.F, 1) {
		t.Errorf("F = %v, want +Inf (real gap between means, zero within-group noise)", r.F)
	}
}

func TestFIsZeroWhenEverythingIsIdentical(t *testing.T) {
	groups := [][]float64{
		GroupSamples(40, 0),
		GroupSamples(40, 0),
	}
	r := Run(groups)
	if !approx(r.F, 0, 1e-9) {
		t.Errorf("F = %v, want 0 (no spread anywhere)", r.F)
	}
}

func TestLargerMeanGapGivesLargerF(t *testing.T) {
	closeGroups := [][]float64{GroupSamples(50, 3), GroupSamples(52, 3), GroupSamples(51, 3)}
	farGroups := [][]float64{GroupSamples(50, 3), GroupSamples(90, 3), GroupSamples(70, 3)}
	fClose, fFar := Run(closeGroups).F, Run(farGroups).F
	if fFar <= fClose {
		t.Errorf("F for widely separated means (%v) should exceed F for close means (%v)", fFar, fClose)
	}
}

func TestRunPanicsWithFewerThanTwoGroups(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Run([one group]) did not panic")
		}
	}()
	Run([][]float64{GroupSamples(10, 1)})
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("anova")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
