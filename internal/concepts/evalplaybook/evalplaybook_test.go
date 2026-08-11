package evalplaybook

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestLossGap(t *testing.T) {
	if g := LossGap(0.30, 0.33); math.Abs(g-0.03) > 1e-9 {
		t.Errorf("LossGap(0.30, 0.33) = %v, want 0.03", g)
	}
	// A model that fits validation data better than training data (unusual,
	// but the formula should still just report the (negative) gap honestly).
	if g := LossGap(0.5, 0.4); math.Abs(g-(-0.1)) > 1e-9 {
		t.Errorf("LossGap(0.5, 0.4) = %v, want -0.1", g)
	}
}

func TestIsOverfitting(t *testing.T) {
	// Overfitting's own numbers: train loss low, gap wide.
	if !IsOverfitting(0.06, 0.58, 0.2, 0.1) {
		t.Error("expected overfitting scenario to be flagged as overfitting")
	}
	// Healthy fit: train loss low-ish but the gap is small.
	if IsOverfitting(0.30, 0.33, 0.2, 0.1) {
		t.Error("healthy fit (train loss above lowThreshold) should not be flagged as overfitting")
	}
	// Underfitting: train loss itself is high, even though val loss is close.
	if IsOverfitting(0.61, 0.63, 0.2, 0.1) {
		t.Error("underfitting should not be flagged as overfitting")
	}
}

func TestF1(t *testing.T) {
	if f := F1(0.91, 0.90); math.Abs(f-0.9049723756906077) > 1e-9 {
		t.Errorf("F1(0.91, 0.90) = %v, want ~0.905", f)
	}
	if f := F1(0, 0); f != 0 {
		t.Errorf("F1(0, 0) = %v, want 0", f)
	}
	// Precision and recall both 1: F1 must be exactly 1.
	if f := F1(1, 1); math.Abs(f-1) > 1e-9 {
		t.Errorf("F1(1, 1) = %v, want 1", f)
	}
}

func TestScenariosBounded(t *testing.T) {
	scenarios := Scenarios()
	if len(scenarios) == 0 {
		t.Fatal("Scenarios() returned no rows")
	}
	for _, s := range scenarios {
		if s.Name == "" {
			t.Error("scenario with empty Name")
		}
		if s.Diagnosis == "" || s.Action == "" {
			t.Errorf("%s: Diagnosis and Action must both be set", s.Name)
		}
		for label, v := range map[string]float64{
			"ROCAUC": s.ROCAUC, "PRAUC": s.PRAUC,
			"Precision": s.Precision, "Recall": s.Recall,
		} {
			if v < 0 || v > 1 {
				t.Errorf("%s: %s = %v, want in [0,1]", s.Name, label, v)
			}
		}
		if s.TrainLoss < 0 || s.ValLoss < 0 {
			t.Errorf("%s: loss values must be non-negative (train=%v, val=%v)", s.Name, s.TrainLoss, s.ValLoss)
		}
	}
}

func TestOverfittingScenarioHasWidestLossGap(t *testing.T) {
	// The one scenario literally named "Overfitting" should have the largest
	// train/val loss gap of the four — otherwise the table is internally
	// inconsistent with its own diagnosis.
	scenarios := Scenarios()
	var overfitGap float64
	maxOtherGap := math.Inf(-1)
	found := false
	for _, s := range scenarios {
		gap := LossGap(s.TrainLoss, s.ValLoss)
		if s.Name == "Overfitting" {
			overfitGap = gap
			found = true
			continue
		}
		if gap > maxOtherGap {
			maxOtherGap = gap
		}
	}
	if !found {
		t.Fatal(`no scenario named "Overfitting" found`)
	}
	if overfitGap <= maxOtherGap {
		t.Errorf("Overfitting's loss gap (%.3f) should exceed every other scenario's (max %.3f)", overfitGap, maxOtherGap)
	}
}

func TestImbalanceTrapHasHighROCLowPR(t *testing.T) {
	// The scenario this whole concept was requested for: high ROC-AUC,
	// noticeably lower PR-AUC and precision, despite healthy loss.
	scenarios := Scenarios()
	for _, s := range scenarios {
		if s.Name != "Imbalance trap" {
			continue
		}
		if s.ROCAUC < 0.9 {
			t.Errorf("Imbalance trap: ROCAUC = %v, want a high (>=0.9) value to sell the contrast", s.ROCAUC)
		}
		if s.PRAUC >= s.ROCAUC-0.2 {
			t.Errorf("Imbalance trap: PRAUC (%v) should sit well below ROCAUC (%v)", s.PRAUC, s.ROCAUC)
		}
		if s.Precision >= 0.5 {
			t.Errorf("Imbalance trap: Precision = %v, want a collapsed (<0.5) value", s.Precision)
		}
		return
	}
	t.Fatal(`no scenario named "Imbalance trap" found`)
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("eval-playbook")
	if !ok {
		t.Fatal("eval-playbook not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
