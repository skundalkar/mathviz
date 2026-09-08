package kerneltrick

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// TestOriginalSpaceNotSeparable checks the concept's core claim: no single
// threshold on the raw x values separates class A from class B.
func TestOriginalSpaceNotSeparable(t *testing.T) {
	xs := make([]float64, len(Points))
	labels := make([]int, len(Points))
	for i, pt := range Points {
		xs[i], labels[i] = pt.X, pt.Label
	}
	if SeparableByThreshold(xs, labels) {
		t.Error("SeparableByThreshold(raw x) = true, want false")
	}
}

// TestMappedSpaceSeparable checks the flip side: thresholding on x²
// (Phi's second coordinate) does separate the two classes.
func TestMappedSpaceSeparable(t *testing.T) {
	ys := make([]float64, len(Points))
	labels := make([]int, len(Points))
	for i, pt := range Points {
		_, y := Phi(pt.X)
		ys[i], labels[i] = y, pt.Label
	}
	if !SeparableByThreshold(ys, labels) {
		t.Error("SeparableByThreshold(x^2) = false, want true")
	}
}

// TestPhiKnownValues is pinned to the x² values worked in LEARNINGS.md and
// the concept's own Sections.
func TestPhiKnownValues(t *testing.T) {
	cases := []struct {
		x, wantX2 float64
	}{
		{-3, 9}, {-2, 4}, {-1, 1}, {0, 0}, {1, 1}, {2, 4}, {3, 9},
	}
	for _, c := range cases {
		x, y := Phi(c.x)
		if x != c.x {
			t.Errorf("Phi(%v) first coord = %v, want %v", c.x, x, c.x)
		}
		if y != c.wantX2 {
			t.Errorf("Phi(%v) second coord = %v, want %v", c.x, y, c.wantX2)
		}
	}
}

// TestKernelMatchesExplicitPhiDot checks the concept's central claim: the
// kernel formula and explicitly building both Phi vectors and dotting them
// always agree.
func TestKernelMatchesExplicitPhiDot(t *testing.T) {
	vals := []float64{-3, -2.5, -1, 0, 0.5, 1, 2, 3}
	for _, a := range vals {
		for _, b := range vals {
			k := Kernel(a, b)
			want := PhiDot(a, b)
			if !almostEqual(k, want, 1e-9) {
				t.Errorf("Kernel(%v,%v) = %v, want %v (PhiDot)", a, b, k, want)
			}
		}
	}
}

// TestKernelKnownValue is pinned to the K(2,-1)=2 example worked in the
// concept's Sections.
func TestKernelKnownValue(t *testing.T) {
	got := Kernel(2, -1)
	if !almostEqual(got, 2, 1e-9) {
		t.Errorf("Kernel(2,-1) = %v, want 2", got)
	}
}

// TestMarginThresholdKnownValue is pinned to the threshold=2.5 example
// worked in LEARNINGS.md and the concept's own Sections.
func TestMarginThresholdKnownValue(t *testing.T) {
	threshold, innerMax, outerMin := MarginThreshold(Points)
	if !almostEqual(innerMax, 1, 1e-9) {
		t.Errorf("innerMax = %v, want 1", innerMax)
	}
	if !almostEqual(outerMin, 4, 1e-9) {
		t.Errorf("outerMin = %v, want 4", outerMin)
	}
	if !almostEqual(threshold, 2.5, 1e-9) {
		t.Errorf("threshold = %v, want 2.5", threshold)
	}
}

// TestClassifyMatchesTrueLabels checks that Classify, using
// MarginThreshold's own threshold, correctly recovers every point's true
// label -- i.e. the mapped-space split actually works as a classifier.
func TestClassifyMatchesTrueLabels(t *testing.T) {
	threshold, _, _ := MarginThreshold(Points)
	for _, pt := range Points {
		got := Classify(pt.X, threshold)
		if got != pt.Label {
			t.Errorf("Classify(%v, %v) = %d, want %d", pt.X, threshold, got, pt.Label)
		}
	}
}

// TestSeparableByThresholdEmpty checks the degenerate empty-input case.
func TestSeparableByThresholdEmpty(t *testing.T) {
	if !SeparableByThreshold(nil, nil) {
		t.Error("SeparableByThreshold(empty) = false, want true (vacuously separable)")
	}
}

// TestSeparableByThresholdSingleClass checks a trivially separable case
// (every point the same label) as a sanity check on the general algorithm.
func TestSeparableByThresholdSingleClass(t *testing.T) {
	vals := []float64{1, 2, 3}
	labels := []int{1, 1, 1}
	if !SeparableByThreshold(vals, labels) {
		t.Error("SeparableByThreshold(single class) = false, want true")
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("kernel-trick")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
