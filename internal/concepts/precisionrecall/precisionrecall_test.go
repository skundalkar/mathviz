package precisionrecall

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestThresholdTradeoff(t *testing.T) {
	// Moving the threshold right should not decrease precision and should not
	// increase recall — the core precision/recall tension.
	sep := 3.0
	var lastPrec, lastRec float64 = -1, 2
	for thr := -3.0; thr <= 6.0; thr += 0.25 {
		prec, rec, _ := Metrics(thr, sep)
		if prec < lastPrec-1e-9 {
			t.Errorf("precision dropped as threshold rose (thr=%.2f)", thr)
		}
		if rec > lastRec+1e-9 {
			t.Errorf("recall rose as threshold rose (thr=%.2f)", thr)
		}
		lastPrec, lastRec = prec, rec
	}
}

func TestMetricsBounded(t *testing.T) {
	for _, thr := range []float64{-3, 0, 1.5, 3, 6} {
		prec, rec, f1 := Metrics(thr, 3)
		for name, v := range map[string]float64{"precision": prec, "recall": rec, "f1": f1} {
			if v < -1e-9 || v > 1+1e-9 {
				t.Errorf("%s = %v out of [0,1] at thr=%v", name, v, thr)
			}
		}
	}
}

func TestBetterSeparationRaisesF1(t *testing.T) {
	// At a fixed sensible threshold, more separation = an easier problem = higher F1.
	_, _, f1Low := Metrics(1.5, 2)
	_, _, f1High := Metrics(1.5, 4)
	if f1High <= f1Low {
		t.Fatalf("more separation did not raise F1: %.3f vs %.3f", f1High, f1Low)
	}
}

func TestSymmetricThresholdBalances(t *testing.T) {
	// With the threshold exactly between the classes, precision == recall.
	sep := 3.0
	prec, rec, _ := Metrics(sep/2, sep)
	if math.Abs(prec-rec) > 1e-9 {
		t.Fatalf("at the midpoint precision (%.4f) != recall (%.4f)", prec, rec)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, _ := concept.Get("precision-recall")
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
