package calibration

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestObservedFrequencyIdentityAtTempOne(t *testing.T) {
	// Perfect calibration is the diagonal: reported probability equals
	// observed frequency exactly when temp=1.
	for _, p := range []float64{0.05, 0.1, 0.3, 0.5, 0.7, 0.9, 0.95} {
		if got := ObservedFrequency(p, 1); math.Abs(got-p) > 1e-9 {
			t.Errorf("ObservedFrequency(%.2f, 1) = %.6f, want %.2f", p, got, p)
		}
	}
}

func TestObservedFrequencyOverconfidentIsSteeper(t *testing.T) {
	// temp<1 (overconfident): a high predicted probability corresponds to a
	// LOWER actual frequency, and a low predicted probability to a HIGHER
	// one — the curve is steeper than the diagonal through its center.
	if got := ObservedFrequency(0.9, 0.5); got >= 0.9 {
		t.Errorf("ObservedFrequency(0.9, 0.5) = %.3f, want < 0.9 (overconfident)", got)
	}
	if got := ObservedFrequency(0.1, 0.5); got <= 0.1 {
		t.Errorf("ObservedFrequency(0.1, 0.5) = %.3f, want > 0.1 (overconfident)", got)
	}
}

func TestObservedFrequencyUnderconfidentIsFlatter(t *testing.T) {
	// temp>1 (underconfident): the opposite pattern.
	if got := ObservedFrequency(0.9, 2); got <= 0.9 {
		t.Errorf("ObservedFrequency(0.9, 2) = %.3f, want > 0.9 (underconfident)", got)
	}
	if got := ObservedFrequency(0.1, 2); got >= 0.1 {
		t.Errorf("ObservedFrequency(0.1, 2) = %.3f, want < 0.1 (underconfident)", got)
	}
}

func TestObservedFrequencyMatchesWorkedExample(t *testing.T) {
	// The exact "90% predicted, 75% actual" example this concept was built
	// around: sep=3, temp=0.5.
	if got := ObservedFrequency(0.9, 0.5); math.Abs(got-0.75) > 1e-2 {
		t.Errorf("ObservedFrequency(0.9, 0.5) = %.3f, want ~0.75", got)
	}
}

func TestReportedMatchesTrueAtTempOne(t *testing.T) {
	for _, z := range []float64{-3, -1, 0, 1.5, 3, 5} {
		sep := 3.0
		rep := ReportedProbability(z, sep, 1)
		tru := TrueProbability(z, sep)
		if math.Abs(rep-tru) > 1e-9 {
			t.Errorf("at temp=1, ReportedProbability(%.1f) = %.6f != TrueProbability = %.6f", z, rep, tru)
		}
	}
}

func TestProbabilitiesBounded(t *testing.T) {
	// z=-4..4 comfortably covers this concept's plotted range without
	// saturating. Sharpening a logit at the most extreme (z, temp)
	// combinations this concept allows (e.g. z=8, temp=0.3) legitimately
	// rounds sigmoid to exactly 1 in float64 once the logit exceeds ~37 —
	// correct floating-point behavior, not a bug — so this test stays
	// inside the range that shouldn't saturate.
	for _, z := range []float64{-4, -2, 0, 1.5, 3, 4} {
		for _, temp := range []float64{0.3, 1, 3} {
			if v := ReportedProbability(z, 3, temp); v <= 0 || v >= 1 {
				t.Errorf("ReportedProbability(%.1f, 3, %.1f) = %v, want strictly in (0,1)", z, temp, v)
			}
		}
		if v := TrueProbability(z, 3); v <= 0 || v >= 1 {
			t.Errorf("TrueProbability(%.1f, 3) = %v, want strictly in (0,1)", z, v)
		}
	}
}

func TestECEZeroAtTempOne(t *testing.T) {
	for _, sep := range []float64{1, 2, 3, 5} {
		if ece := ExpectedCalibrationError(sep, 1, 500); ece > 1e-6 {
			t.Errorf("ExpectedCalibrationError(sep=%.1f, temp=1) = %.6f, want ~0", sep, ece)
		}
	}
}

func TestECEGrowsAwayFromTempOne(t *testing.T) {
	// ECE should increase monotonically as temp moves further from 1 in
	// EITHER direction — checked separately below 1 and above 1.
	sep := 3.0
	below := []float64{0.9, 0.7, 0.5, 0.3}
	last := -1.0
	for _, temp := range below {
		ece := ExpectedCalibrationError(sep, temp, 1000)
		if ece < last {
			t.Errorf("ECE did not increase as temp fell to %.1f: got %.4f after %.4f", temp, ece, last)
		}
		last = ece
	}
	above := []float64{1.1, 1.5, 2, 3}
	last = -1.0
	for _, temp := range above {
		ece := ExpectedCalibrationError(sep, temp, 1000)
		if ece < last {
			t.Errorf("ECE did not increase as temp rose to %.1f: got %.4f after %.4f", temp, ece, last)
		}
		last = ece
	}
}

func TestECENonNegative(t *testing.T) {
	for _, sep := range []float64{1, 3, 5} {
		for _, temp := range []float64{0.3, 0.5, 1, 2, 3} {
			if ece := ExpectedCalibrationError(sep, temp, 500); ece < 0 {
				t.Errorf("ExpectedCalibrationError(%.1f, %.1f) = %v, want >= 0", sep, temp, ece)
			}
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("calibration")
	if !ok {
		t.Fatal("calibration not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
