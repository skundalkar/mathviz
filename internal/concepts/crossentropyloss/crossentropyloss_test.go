package crossentropyloss

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func approx(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestLossWorkedExample(t *testing.T) {
	cases := []struct {
		name     string
		y, phat  float64
		wantLoss float64
	}{
		{"email1 confident correct spam", 1, 0.9, 0.152},
		{"email2 barely correct spam", 1, 0.55, 0.863},
		{"email3 confident correct not-spam", 0, 0.2, 0.322},
		{"email4 confident wrong spam", 1, 0.1, 3.322},
	}
	for _, c := range cases {
		if got := Loss(c.y, c.phat); !approx(got, c.wantLoss, 5e-3) {
			t.Errorf("%s: Loss(%v,%v) = %v, want ~%v", c.name, c.y, c.phat, got, c.wantLoss)
		}
	}
}

func TestAverageLossWorkedExample(t *testing.T) {
	got := AverageLoss(Labels, Preds)
	if !approx(got, 1.165, 5e-3) {
		t.Errorf("AverageLoss(Labels, Preds) = %v, want ~1.165", got)
	}
}

func TestLossApproachesZeroWhenConfidentAndCorrect(t *testing.T) {
	if got := Loss(1, 0.999); got >= 0.01 {
		t.Errorf("Loss(1, 0.999) = %v, want close to 0", got)
	}
	if got := Loss(0, 0.001); got >= 0.01 {
		t.Errorf("Loss(0, 0.001) = %v, want close to 0", got)
	}
}

func TestLossBlowsUpWhenConfidentAndWrong(t *testing.T) {
	if got := Loss(1, 0.001); got < 9 {
		t.Errorf("Loss(1, 0.001) = %v, want a large penalty", got)
	}
	if got := Loss(0, 0.999); got < 9 {
		t.Errorf("Loss(0, 0.999) = %v, want a large penalty", got)
	}
}

func TestLossIsSymmetricAroundHalf(t *testing.T) {
	// Being 50/50 unsure costs the same 1 bit whichever label turns out true.
	a, b := Loss(1, 0.5), Loss(0, 0.5)
	if !approx(a, 1.0, 1e-9) || !approx(b, 1.0, 1e-9) {
		t.Errorf("Loss(1,0.5)=%v Loss(0,0.5)=%v, want both 1.0", a, b)
	}
}

func TestMoreConfidentCorrectPredictionCostsLess(t *testing.T) {
	if Loss(1, 0.99) >= Loss(1, 0.6) {
		t.Errorf("Loss(1,0.99)=%v should be less than Loss(1,0.6)=%v", Loss(1, 0.99), Loss(1, 0.6))
	}
}

func TestAverageLossEmptyDataset(t *testing.T) {
	if got := AverageLoss(nil, nil); got != 0 {
		t.Errorf("AverageLoss(nil,nil) = %v, want 0", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("cross-entropy-loss")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
