package decisiontree

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestEntropyKnownValues(t *testing.T) {
	if got := Entropy(0.5); !near(got, 1.0, 1e-9) {
		t.Errorf("Entropy(0.5) = %v, want 1.0", got)
	}
	if got := Entropy(0); got != 0 {
		t.Errorf("Entropy(0) = %v, want 0", got)
	}
	if got := Entropy(1); got != 0 {
		t.Errorf("Entropy(1) = %v, want 0", got)
	}
	if got := Entropy(0.6); !near(got, 0.970951, 1e-5) {
		t.Errorf("Entropy(0.6) = %v, want ~0.970951", got)
	}
}

func TestSplitPartitionsByThreshold(t *testing.T) {
	left, right := Split(Hours, Pass, 2.75)
	if want := []int{0, 0, 0, 0}; !equalInts(left, want) {
		t.Errorf("left = %v, want %v", left, want)
	}
	if want := []int{1, 0, 1, 1, 1, 1}; !equalInts(right, want) {
		t.Errorf("right = %v, want %v", right, want)
	}
}

func TestParentEntropyIsMaximal(t *testing.T) {
	// The full dataset is 5 pass / 5 fail -> maximum uncertainty, 1 bit.
	full := Entropy(5.0 / 10.0)
	if !near(full, 1.0, 1e-9) {
		t.Errorf("whole-dataset entropy = %v, want 1.0 (5 of 10 passed)", full)
	}
}

func TestInfoGainKnownValues(t *testing.T) {
	cases := []struct {
		threshold float64
		wantGain  float64
	}{
		{3.25, 0.2781}, // the "obvious" midpoint - not actually best
		{2.75, 0.6100}, // the actual best split (tied with 3.75)
		{3.75, 0.6100},
	}
	for _, c := range cases {
		left, right := Split(Hours, Pass, c.threshold)
		gain := InfoGain(Pass, left, right)
		if !near(gain, c.wantGain, 1e-3) {
			t.Errorf("InfoGain at threshold=%v = %v, want ~%v", c.threshold, gain, c.wantGain)
		}
	}
}

func TestInfoGainNeverNegative(t *testing.T) {
	for th := 1.25; th < 5.5; th += 0.5 {
		left, right := Split(Hours, Pass, th)
		if gain := InfoGain(Pass, left, right); gain < -1e-9 {
			t.Errorf("InfoGain at threshold=%v = %v, want >= 0", th, gain)
		}
	}
}

func TestInfoGainPerfectSeparationEqualsParentEntropy(t *testing.T) {
	labels := []int{0, 0, 1, 1}
	left := []int{0, 0}
	right := []int{1, 1}
	gain := InfoGain(labels, left, right)
	parent := classEntropy(labels)
	if !near(gain, parent, 1e-9) {
		t.Errorf("perfectly separating split gain = %v, want parent entropy %v", gain, parent)
	}
	if !near(gain, 1.0, 1e-9) {
		t.Errorf("gain = %v, want 1.0 (parent was a 50/50 split)", gain)
	}
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("decision-trees")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
