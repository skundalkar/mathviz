package simpsonsparadox

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

const tol = 1e-9

func near(a, b float64) bool { return math.Abs(a-b) < tol }

func TestCombinedRateEndpoints(t *testing.T) {
	// mixLarge=0 -> pure small-stone rate; mixLarge=1 -> pure large-stone rate.
	if got := CombinedRate(0, 0.93, 0.73); !near(got, 0.93) {
		t.Errorf("CombinedRate(0, ...) = %v, want 0.93", got)
	}
	if got := CombinedRate(1, 0.93, 0.73); !near(got, 0.73) {
		t.Errorf("CombinedRate(1, ...) = %v, want 0.73", got)
	}
}

func TestCombinedRateMidpoint(t *testing.T) {
	// Equal mix is a plain average of the two rates.
	got := CombinedRate(0.5, 0.8, 0.6)
	if want := 0.7; !near(got, want) {
		t.Errorf("CombinedRate(0.5, 0.8, 0.6) = %v, want %v", got, want)
	}
}

func TestCombinedRateClampsMix(t *testing.T) {
	if got := CombinedRate(-0.5, 0.9, 0.5); !near(got, 0.9) {
		t.Errorf("CombinedRate(-0.5, ...) = %v, want clamped to mixLarge=0 -> 0.9", got)
	}
	if got := CombinedRate(1.5, 0.9, 0.5); !near(got, 0.5) {
		t.Errorf("CombinedRate(1.5, ...) = %v, want clamped to mixLarge=1 -> 0.5", got)
	}
}

// TestReproducesRealParadox checks that the exact case-mix from the 1986
// Charig et al. study reverses the ranking: Treatment A beats Treatment B on
// both stone sizes, but Treatment B's overall rate (computed from that real
// mix) comes out higher.
func TestReproducesRealParadox(t *testing.T) {
	if !(RateASmall > RateBSmall) || !(RateALarge > RateBLarge) {
		t.Fatalf("expected A to beat B on both subgroups")
	}

	const mixA = 263.0 / 350.0 // A: 263 of 350 patients had large stones
	const mixB = 80.0 / 350.0  // B: 80 of 350 patients had large stones

	overallA := CombinedRate(mixA, RateASmall, RateALarge)
	overallB := CombinedRate(mixB, RateBSmall, RateBLarge)

	if !near(round2(overallA), 0.78) {
		t.Errorf("overall A = %v, want ~0.78 (matches the published 78%%)", overallA)
	}
	if !near(round2(overallB), 0.83) {
		t.Errorf("overall B = %v, want ~0.83 (matches the published 83%%)", overallB)
	}
	if overallA >= overallB {
		t.Errorf("expected the paradox: overall A (%v) < overall B (%v)", overallA, overallB)
	}
}

// TestBalancedMixDissolvesParadox checks that once both treatments face the
// same case mix, the subgroup winner (A) is also the overall winner.
func TestBalancedMixDissolvesParadox(t *testing.T) {
	overallA := CombinedRate(0.5, RateASmall, RateALarge)
	overallB := CombinedRate(0.5, RateBSmall, RateBLarge)
	if overallA <= overallB {
		t.Errorf("with equal case mix expected A > B overall, got A=%v B=%v", overallA, overallB)
	}
}

func round2(x float64) float64 {
	return math.Round(x*100) / 100
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("simpsons-paradox")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document: %q", svg[:min(40, len(svg))])
	}
}
