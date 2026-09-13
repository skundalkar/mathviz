package bonferroni

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

func TestNullPValueDeterministic(t *testing.T) {
	for i := 0; i < 60; i++ {
		if a, b := NullPValue(i), NullPValue(i); a != b {
			t.Errorf("NullPValue(%d) not deterministic: %v then %v", i, a, b)
		}
	}
}

func TestNullPValueStaysInUnitRange(t *testing.T) {
	for i := 0; i < 200; i++ {
		v := NullPValue(i)
		if v < 0 || v >= 1 {
			t.Errorf("NullPValue(%d) = %v, out of [0,1)", i, v)
		}
	}
}

// TestNullPValuesApproximatelyUniform checks the property everything else
// in this concept leans on: under a true null, p-values are uniform on
// [0,1), so roughly alpha's worth of a large batch should fall below any
// given threshold alpha.
func TestNullPValuesApproximatelyUniform(t *testing.T) {
	const n = 4000
	pv := NullPValues(n)
	for _, alpha := range []float64{0.1, 0.3, 0.5} {
		got := float64(CountBelow(pv, alpha)) / float64(n)
		if math.Abs(got-alpha) > 0.03 {
			t.Errorf("fraction below %.2f = %.3f, want close to %.2f", alpha, got, alpha)
		}
	}
}

func TestNullPValuesLength(t *testing.T) {
	if got := len(NullPValues(20)); got != 20 {
		t.Errorf("len(NullPValues(20)) = %d, want 20", got)
	}
	if got := len(NullPValues(0)); got != 0 {
		t.Errorf("len(NullPValues(0)) = %d, want 0", got)
	}
	if got := len(NullPValues(-3)); got != 0 {
		t.Errorf("len(NullPValues(-3)) = %d, want 0 (clamped)", got)
	}
}

func TestNullPValuesMatchesNullPValue(t *testing.T) {
	pv := NullPValues(10)
	for i, v := range pv {
		if want := NullPValue(i); v != want {
			t.Errorf("NullPValues(10)[%d] = %v, want %v (NullPValue(%d))", i, v, want, i)
		}
	}
}

func TestCountBelow(t *testing.T) {
	vals := []float64{0.01, 0.2, 0.03, 0.5, 0.049}
	if got := CountBelow(vals, 0.05); got != 3 {
		t.Errorf("CountBelow(vals,0.05) = %d, want 3 (0.01, 0.03, 0.049)", got)
	}
	if got := CountBelow(vals, 0); got != 0 {
		t.Errorf("CountBelow(vals,0) = %d, want 0", got)
	}
	if got := CountBelow(nil, 0.5); got != 0 {
		t.Errorf("CountBelow(nil,0.5) = %d, want 0", got)
	}
}

// TestFamilyWiseErrorRateKnownValues is pinned to the standard FWER formula
// 1-(1-alpha)^m at alpha=0.05, the same numbers walked through in
// LEARNINGS.md and the concept's own Sections.
func TestFamilyWiseErrorRateKnownValues(t *testing.T) {
	cases := map[int]float64{
		1:  0.0500,
		5:  0.2262,
		10: 0.4013,
		20: 0.6415,
		30: 0.7854,
		50: 0.9231,
	}
	for m, want := range cases {
		if got := FamilyWiseErrorRate(0.05, m); !near(got, want, 5e-4) {
			t.Errorf("FamilyWiseErrorRate(0.05,%d) = %v, want %v", m, got, want)
		}
	}
}

func TestFamilyWiseErrorRateGrowsWithM(t *testing.T) {
	prev := FamilyWiseErrorRate(0.05, 1)
	for m := 2; m <= MaxM; m++ {
		got := FamilyWiseErrorRate(0.05, m)
		if got <= prev {
			t.Errorf("FWER did not increase from m=%d (%.4f) to m=%d (%.4f)", m-1, prev, m, got)
		}
		prev = got
	}
}

func TestFamilyWiseErrorRateAtMEqualsOneIsAlpha(t *testing.T) {
	for _, alpha := range []float64{0.01, 0.05, 0.1, 0.2} {
		if got := FamilyWiseErrorRate(alpha, 1); !near(got, alpha, 1e-12) {
			t.Errorf("FamilyWiseErrorRate(%.2f,1) = %v, want %v (a single test's FWER is just its own alpha)", alpha, got, alpha)
		}
	}
}

func TestBonferroniAlpha(t *testing.T) {
	cases := []struct {
		alpha float64
		m     int
		want  float64
	}{
		{0.05, 1, 0.05},
		{0.05, 20, 0.0025},
		{0.05, 50, 0.001},
		{0.10, 5, 0.02},
	}
	for _, c := range cases {
		if got := BonferroniAlpha(c.alpha, c.m); !near(got, c.want, 1e-9) {
			t.Errorf("BonferroniAlpha(%v,%d) = %v, want %v", c.alpha, c.m, got, c.want)
		}
	}
	if got := BonferroniAlpha(0.05, 0); !near(got, 0.05, 1e-9) {
		t.Errorf("BonferroniAlpha(0.05,0) = %v, want 0.05 (m clamped to 1)", got)
	}
}

// TestCorrectionRemovesTheFalsePositiveRawAlphaLetsThrough is pinned to the
// same 20-test batch every Section walks through: at the raw threshold
// alpha=0.05 exactly one of the 20 (all drawn from a true null) is called
// significant -- a false positive -- but the Bonferroni-corrected threshold
// alpha/m=0.0025 lets none through.
func TestCorrectionRemovesTheFalsePositiveRawAlphaLetsThrough(t *testing.T) {
	pv := NullPValues(20)
	alpha := 0.05
	if got := CountBelow(pv, alpha); got != 1 {
		t.Errorf("CountBelow(pv,0.05) = %d, want 1", got)
	}
	corrected := BonferroniAlpha(alpha, 20)
	if got := CountBelow(pv, corrected); got != 0 {
		t.Errorf("CountBelow(pv,%.4f) = %d, want 0", corrected, got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("bonferroni-correction")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
