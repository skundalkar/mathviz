package prauc

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestRecallMonotonicWithThreshold(t *testing.T) {
	// Raising the threshold should never increase recall.
	sep := 3.0
	last := 2.0
	for thr := -3.0; thr <= 6.0; thr += 0.25 {
		rec := Recall(thr, sep)
		if rec > last+1e-9 {
			t.Errorf("recall rose as threshold rose (thr=%.2f)", thr)
		}
		last = rec
	}
}

func TestPrecisionRecallBounded(t *testing.T) {
	for _, thr := range []float64{-6, -3, 0, 1.5, 3, 6} {
		rec := Recall(thr, 3)
		prec := Precision(thr, 3)
		if rec < -1e-9 || rec > 1+1e-9 {
			t.Errorf("recall = %v out of [0,1] at thr=%v", rec, thr)
		}
		if prec < -1e-9 || prec > 1+1e-9 {
			t.Errorf("precision = %v out of [0,1] at thr=%v", prec, thr)
		}
	}
}

func TestPrecisionAtExtremeThresholds(t *testing.T) {
	// Flag (almost) nothing: whatever gets through is almost certainly a true
	// positive, so precision should sit near 1.
	if p := Precision(6, 3); p < 0.9 {
		t.Errorf("precision at a very high threshold = %.3f, want close to 1", p)
	}
	// Flag everything: precision converges to the positive class's share,
	// which is 0.5 with equal-sized classes here.
	if p := Precision(-6, 3); math.Abs(p-0.5) > 0.05 {
		t.Errorf("precision at a very low threshold = %.3f, want close to 0.5", p)
	}
}

func TestCurvePointsOrderedByRecall(t *testing.T) {
	pts := CurvePoints(3, 50)
	if len(pts) != 51 {
		t.Fatalf("got %d points, want 51", len(pts))
	}
	for i := 1; i < len(pts); i++ {
		if pts[i][0] < pts[i-1][0]-1e-9 {
			t.Fatalf("points not sorted by recall at index %d: %v then %v", i, pts[i-1], pts[i])
		}
	}
	if pts[0][0] > 0.01 || pts[len(pts)-1][0] < 0.99 {
		t.Errorf("curve should span recall ~0 to ~1, got [%.3f, %.3f]", pts[0][0], pts[len(pts)-1][0])
	}
}

// TestCurvePointsSpanFullRecallAcrossSepRange guards against a real bug: the
// threshold sweep used to run a fixed [-6, 6] regardless of sep, but the
// positive class is centered at sep, so at the top of the allowed slider
// range (sep=5) the sweep never got high enough above the positive mean for
// recall to actually reach ~0 — silently truncating the high-precision tip
// of the curve and undercounting PR-AUC for exactly the classifiers that
// should score best. Check every sep the slider allows, not just one.
func TestCurvePointsSpanFullRecallAcrossSepRange(t *testing.T) {
	for sep := 1.0; sep <= 5.0; sep += 0.5 {
		pts := CurvePoints(sep, 500)
		if pts[0][0] > 0.01 {
			t.Errorf("sep=%.1f: curve's lowest recall = %.4f, want ~0", sep, pts[0][0])
		}
		if pts[len(pts)-1][0] < 0.99 {
			t.Errorf("sep=%.1f: curve's highest recall = %.4f, want ~1", sep, pts[len(pts)-1][0])
		}
	}
}

func TestPRAUCBetterSeparationIsHigher(t *testing.T) {
	// More separation between the classes = an easier problem = more area
	// under the PR curve. Check monotonicity across the full slider range,
	// not just two endpoints — the truncated-sweep bug above only showed up
	// once separation got large enough, so PR-AUC actually *dropped* between
	// sep=3 and sep=5 despite the classifier genuinely getting easier.
	var last float64 = -1
	for sep := 1.0; sep <= 5.0; sep += 0.5 {
		auc := TrapezoidalPRAUC(CurvePoints(sep, 2000))
		if auc < last-1e-9 {
			t.Errorf("PR-AUC dropped from %.4f to %.4f as separation rose to %.1f", last, auc, sep)
		}
		last = auc
	}
}

func TestPRAUCBounded(t *testing.T) {
	for _, sep := range []float64{0.5, 1, 2.5, 5} {
		auc := TrapezoidalPRAUC(CurvePoints(sep, 300))
		if auc < 0.4 || auc > 1.0 {
			t.Errorf("PR-AUC = %.3f out of expected [0.4, 1.0] range at sep=%.1f", auc, sep)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("pr-auc")
	if !ok {
		t.Fatal("pr-auc not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
