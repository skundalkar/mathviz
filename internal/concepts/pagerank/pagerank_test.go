package pagerank

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func approx(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestIterateRoundOneWorkedExample(t *testing.T) {
	// Section 2's hand-worked example: start uniform at 0.25 each, d=0.85.
	ranks := []float64{0.25, 0.25, 0.25, 0.25}
	next := Iterate(ranks, Outlinks, 0.85)
	want := []float64{0.25, 0.14375, 0.56875, 0.0375}
	for i, w := range want {
		if !approx(next[i], w, 1e-5) {
			t.Errorf("round 1 rank[%d] = %v, want %v", i, next[i], w)
		}
	}
}

func TestIterateRoundTwoWorkedExample(t *testing.T) {
	ranks := []float64{0.25, 0.25, 0.25, 0.25}
	ranks = Iterate(ranks, Outlinks, 0.85)
	ranks = Iterate(ranks, Outlinks, 0.85)
	want := []float64{0.5209375, 0.14375, 0.2978125, 0.0375}
	for i, w := range want {
		if !approx(ranks[i], w, 1e-6) {
			t.Errorf("round 2 rank[%d] = %v, want %v", i, ranks[i], w)
		}
	}
}

func TestIterateAlwaysSumsToOne(t *testing.T) {
	ranks := []float64{0.25, 0.25, 0.25, 0.25}
	for i := 0; i < 20; i++ {
		ranks = Iterate(ranks, Outlinks, 0.85)
		sum := 0.0
		for _, r := range ranks {
			sum += r
		}
		if !approx(sum, 1.0, 1e-9) {
			t.Fatalf("round %d: ranks sum to %v, want 1.0", i, sum)
		}
	}
}

func TestConvergeMatchesRepeatedIterate(t *testing.T) {
	got := Converge(Outlinks, 0.85, 3)
	ranks := []float64{0.25, 0.25, 0.25, 0.25}
	for i := 0; i < 3; i++ {
		ranks = Iterate(ranks, Outlinks, 0.85)
	}
	for i := range got {
		if !approx(got[i], ranks[i], 1e-9) {
			t.Errorf("Converge[%d] = %v, want %v", i, got[i], ranks[i])
		}
	}
}

func TestConvergeZeroIterationsIsUniform(t *testing.T) {
	got := Converge(Outlinks, 0.85, 0)
	for i, r := range got {
		if !approx(r, 0.25, 1e-9) {
			t.Errorf("Converge(t=0)[%d] = %v, want 0.25", i, r)
		}
	}
}

func TestConvergePageDIsAlwaysExactlyTheFloor(t *testing.T) {
	// Nothing links to D (index 3): every round it can only ever receive
	// the (1-d)/n teleport floor, never any link-following contribution.
	const d = 0.85
	want := (1 - d) / 4
	for iters := 1; iters <= 20; iters++ {
		ranks := Converge(Outlinks, d, iters)
		if !approx(ranks[3], want, 1e-9) {
			t.Errorf("t=%d: D's rank = %v, want exactly %v", iters, ranks[3], want)
		}
	}
}

func TestConvergeDWithFullDampingCollapsesDToZero(t *testing.T) {
	// d=1 turns off teleporting entirely: pure link-following. D has no
	// inbound links, so it collapses to exactly zero, permanently.
	ranks := Converge(Outlinks, 1.0, 10)
	if !approx(ranks[3], 0, 1e-9) {
		t.Errorf("D's rank with d=1 = %v, want 0", ranks[3])
	}
}

func TestConvergeApproachesKnownSteadyState(t *testing.T) {
	// Round-20 values from the hand-worked example in section 2.
	got := Converge(Outlinks, 0.85, 20)
	want := []float64{0.3725, 0.1958, 0.3942, 0.0375}
	for i, w := range want {
		if !approx(got[i], w, 1e-3) {
			t.Errorf("t=20 rank[%d] = %v, want ~%v", i, got[i], w)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("pagerank")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
