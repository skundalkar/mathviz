package autocorr

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

func TestSeriesLength(t *testing.T) {
	if got := len(Series()); got != SeriesLength {
		t.Errorf("len(Series()) = %d, want %d", got, SeriesLength)
	}
}

func TestSeriesDeterministic(t *testing.T) {
	a, b := Series(), Series()
	for i := range a {
		if a[i] != b[i] {
			t.Errorf("Series()[%d] not deterministic: %v then %v", i, a[i], b[i])
		}
	}
}

// TestPearsonRSelfIsOne checks the base case ACF relies on: a series
// correlated with an unshifted copy of itself is a perfect straight line,
// r = 1.
func TestPearsonRSelfIsOne(t *testing.T) {
	s := Series()
	if got := PearsonR(s, s); !near(got, 1, 1e-9) {
		t.Errorf("PearsonR(s,s) = %v, want 1", got)
	}
}

func TestPearsonRZeroVariance(t *testing.T) {
	flat := []float64{5, 5, 5, 5}
	other := []float64{1, 2, 3, 4}
	if got := PearsonR(flat, other); got != 0 {
		t.Errorf("PearsonR(flat, other) = %v, want 0", got)
	}
}

// TestACFAtLagZeroIsOne checks that every series (not just the fixed
// example) is perfectly autocorrelated with itself at lag 0.
func TestACFAtLagZeroIsOne(t *testing.T) {
	if got := ACF(Series(), 0); !near(got, 1, 1e-9) {
		t.Errorf("ACF(Series(),0) = %v, want 1", got)
	}
	other := []float64{2, 4, 1, 8, 3}
	if got := ACF(other, 0); !near(got, 1, 1e-9) {
		t.Errorf("ACF(other,0) = %v, want 1", got)
	}
}

func TestACFOutOfRangeLagIsZero(t *testing.T) {
	s := Series()
	if got := ACF(s, -1); got != 0 {
		t.Errorf("ACF(s,-1) = %v, want 0", got)
	}
	if got := ACF(s, len(s)); got != 0 {
		t.Errorf("ACF(s,len(s)) = %v, want 0", got)
	}
	if got := ACF(s, len(s)+5); got != 0 {
		t.Errorf("ACF(s,len(s)+5) = %v, want 0", got)
	}
}

// TestACFKnownValues is pinned to the deterministic Series() sequence -- the
// same numbers walked through in LEARNINGS.md and the concept's own
// Sections. Series has a hidden ~7-step period, so ACF should peak near
// lag 7 and 14 (a full period later, still roughly in phase) and dip most
// negative near lag 3-4 (about half a period out of phase).
func TestACFKnownValues(t *testing.T) {
	s := Series()
	cases := map[int]float64{
		0:  1.0000,
		1:  0.5779,
		3:  -0.8527,
		7:  0.9521,
		14: 0.9317,
	}
	for lag, want := range cases {
		if got := ACF(s, lag); !near(got, want, 5e-4) {
			t.Errorf("ACF(s,%d) = %v, want %v", lag, got, want)
		}
	}
}

// TestACFPeaksNearThePeriodNotAtHalfPeriod is a structural check on the same
// claim, independent of the exact pinned decimals above: the lag with the
// highest ACF among 1..MaxLag should be a near-multiple of Period, and it
// should score far higher than the lag around half that period.
func TestACFPeaksNearThePeriodNotAtHalfPeriod(t *testing.T) {
	s := Series()
	best, bestLag := -2.0, -1
	for lag := 1; lag <= MaxLag; lag++ {
		if v := ACF(s, lag); v > best {
			best, bestLag = v, lag
		}
	}
	if bestLag != 7 && bestLag != 14 {
		t.Errorf("best lag = %d, want 7 or 14 (near Period=%.0f)", bestLag, Period)
	}
	halfPeriodACF := ACF(s, 3)
	if best <= halfPeriodACF {
		t.Errorf("ACF at best lag (%d) = %.4f, should exceed ACF at half-period lag 3 = %.4f", bestLag, best, halfPeriodACF)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("autocorrelation")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
