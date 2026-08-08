package confusionmatrix

import (
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestCountsSumToPopulation(t *testing.T) {
	cases := []struct {
		thresh, sep float64
		n           int
	}{
		{1.5, 3, 200},
		{0, 2, 101}, // odd n
		{-3, 1, 20},
		{6, 5, 500},
	}
	for _, c := range cases {
		tp, fp, fn, tn := Counts(c.thresh, c.sep, c.n)
		if got := tp + fp + fn + tn; got != c.n {
			t.Errorf("Counts(%v,%v,%d) sums to %d, want %d", c.thresh, c.sep, c.n, got, c.n)
		}
		if tp < 0 || fp < 0 || fn < 0 || tn < 0 {
			t.Errorf("Counts(%v,%v,%d) produced a negative count: tp=%d fp=%d fn=%d tn=%d",
				c.thresh, c.sep, c.n, tp, fp, fn, tn)
		}
	}
}

func TestCountsExtremeThresholds(t *testing.T) {
	// A threshold far above both distributions calls almost nothing positive.
	tp, fp, fn, tn := Counts(50, 3, 200)
	if tp != 0 || fp != 0 {
		t.Errorf("Counts(50, 3, 200) = tp=%d fp=%d, want both 0", tp, fp)
	}
	if fn != 100 || tn != 100 {
		t.Errorf("Counts(50, 3, 200) = fn=%d tn=%d, want both 100", fn, tn)
	}

	// A threshold far below both distributions calls almost everything positive.
	tp, fp, fn, tn = Counts(-50, 3, 200)
	if fn != 0 || tn != 0 {
		t.Errorf("Counts(-50, 3, 200) = fn=%d tn=%d, want both 0", fn, tn)
	}
	if tp != 100 || fp != 100 {
		t.Errorf("Counts(-50, 3, 200) = tp=%d fp=%d, want both 100", tp, fp)
	}
}

func TestCountsMonotonicAsThresholdFalls(t *testing.T) {
	// Lowering the threshold should never call fewer examples positive.
	prevTP, prevFP := 0, 0
	for _, thresh := range []float64{6, 3, 1.5, 0, -1, -3} {
		tp, fp, _, _ := Counts(thresh, 3, 200)
		if tp < prevTP || fp < prevFP {
			t.Fatalf("thresh=%v: tp=%d fp=%d should be >= previous tp=%d fp=%d",
				thresh, tp, fp, prevTP, prevFP)
		}
		prevTP, prevFP = tp, fp
	}
}

func TestMetricsKnownValues(t *testing.T) {
	acc, prec, rec, f1 := Metrics(40, 10, 10, 40)
	want := 0.8
	for name, got := range map[string]float64{
		"accuracy": acc, "precision": prec, "recall": rec, "f1": f1,
	} {
		if diff := got - want; diff > 1e-9 || diff < -1e-9 {
			t.Errorf("%s = %v, want %v", name, got, want)
		}
	}
}

func TestMetricsPerfectClassifier(t *testing.T) {
	acc, prec, rec, f1 := Metrics(50, 0, 0, 50)
	if acc != 1 || prec != 1 || rec != 1 || f1 != 1 {
		t.Errorf("Metrics(50,0,0,50) = acc=%v prec=%v rec=%v f1=%v, want all 1", acc, prec, rec, f1)
	}
}

func TestMetricsAllZeroNoDivideByZero(t *testing.T) {
	acc, prec, rec, f1 := Metrics(0, 0, 0, 0)
	if acc != 0 || prec != 0 || rec != 0 || f1 != 0 {
		t.Errorf("Metrics(0,0,0,0) = acc=%v prec=%v rec=%v f1=%v, want all 0", acc, prec, rec, f1)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("confusion-matrix")
	if !ok {
		t.Fatal("concept \"confusion-matrix\" not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
