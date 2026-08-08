package rocauc

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestTPRandFPRAtZeroSeparation(t *testing.T) {
	// With no separation, positives and negatives are the same distribution,
	// so at threshold 0 both rates should be 0.5.
	if got := TPR(0, 0); math.Abs(got-0.5) > 1e-9 {
		t.Errorf("TPR(0, 0) = %v, want 0.5", got)
	}
	if got := FPR(0); math.Abs(got-0.5) > 1e-9 {
		t.Errorf("FPR(0) = %v, want 0.5", got)
	}
}

func TestTPRandFPRDecreaseAsThresholdRises(t *testing.T) {
	prevTPR, prevFPR := TPR(-4, 2), FPR(-4)
	for _, t2 := range []float64{-2, 0, 2, 4} {
		gotTPR, gotFPR := TPR(t2, 2), FPR(t2)
		if gotTPR >= prevTPR {
			t.Errorf("TPR should fall as threshold rises: TPR(%v)=%v >= previous %v", t2, gotTPR, prevTPR)
		}
		if gotFPR >= prevFPR {
			t.Errorf("FPR should fall as threshold rises: FPR(%v)=%v >= previous %v", t2, gotFPR, prevFPR)
		}
		prevTPR, prevFPR = gotTPR, gotFPR
	}
}

func TestAUCNoSeparationIsRandomGuessing(t *testing.T) {
	if got := AUC(0); math.Abs(got-0.5) > 1e-9 {
		t.Errorf("AUC(0) = %v, want 0.5 (a coin flip)", got)
	}
}

func TestAUCIncreasesWithSeparation(t *testing.T) {
	prev := AUC(0)
	for _, sep := range []float64{0.5, 1, 2, 3, 5} {
		got := AUC(sep)
		if got <= prev {
			t.Fatalf("AUC should rise with separation: AUC(%v)=%v <= previous %v", sep, got, prev)
		}
		prev = got
	}
	if got := AUC(5); got <= 0.99 {
		t.Errorf("AUC(5) = %v, want close to 1 for well-separated classes", got)
	}
}

func TestCurvePointsEndpoints(t *testing.T) {
	pts := CurvePoints(2.5, 400)
	if len(pts) != 401 {
		t.Fatalf("CurvePoints(2.5, 400) returned %d points, want 401", len(pts))
	}
	first, last := pts[0], pts[len(pts)-1]
	if first[0] > 0.01 || first[1] > 0.01 {
		t.Errorf("curve should start near (0,0), got (%v, %v)", first[0], first[1])
	}
	if last[0] < 0.99 || last[1] < 0.99 {
		t.Errorf("curve should end near (1,1), got (%v, %v)", last[0], last[1])
	}
}

func TestCurvePointsMonotonic(t *testing.T) {
	pts := CurvePoints(1.5, 200)
	for i := 1; i < len(pts); i++ {
		if pts[i][0] < pts[i-1][0]-1e-9 || pts[i][1] < pts[i-1][1]-1e-9 {
			t.Fatalf("ROC curve should be monotonic in both FPR and TPR: point %d (%v,%v) < point %d (%v,%v)",
				i, pts[i][0], pts[i][1], i-1, pts[i-1][0], pts[i-1][1])
		}
	}
}

func TestTrapezoidalAUCMatchesClosedForm(t *testing.T) {
	for _, sep := range []float64{0, 0.5, 1.5, 2.5, 4} {
		pts := CurvePoints(sep, 2000)
		numeric := TrapezoidalAUC(pts)
		closed := AUC(sep)
		if math.Abs(numeric-closed) > 0.01 {
			t.Errorf("sep=%v: numeric AUC=%v, closed-form AUC=%v, differ by more than 0.01", sep, numeric, closed)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("roc-auc")
	if !ok {
		t.Fatal("concept \"roc-auc\" not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
