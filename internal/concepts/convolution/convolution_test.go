package convolution

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func almostEqualSlice(t *testing.T, got, want []float64, tol float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (got %v, want %v)", len(got), len(want), got, want)
	}
	for i := range want {
		if !almostEqual(got[i], want[i], tol) {
			t.Errorf("[%d] = %v, want %v (full: got %v, want %v)", i, got[i], want[i], got, want)
		}
	}
}

// TestConvolveSmoothingKnownValues is pinned to the moving-average example
// worked step by step in LEARNINGS.md and the concept's own Sections.
func TestConvolveSmoothingKnownValues(t *testing.T) {
	got := Convolve(Signal, Kernels[0])
	want := []float64{1.0 / 3, 2.0 / 3, 1, 1, 2.0 / 3, 1.0 / 3}
	almostEqualSlice(t, got, want, 1e-9)
}

// TestConvolveEdgeKnownValues is pinned to the true-convolution
// (kernel-flipped) edge-detector example.
func TestConvolveEdgeKnownValues(t *testing.T) {
	got := Convolve(Signal, Kernels[1])
	want := []float64{-1, -1, 0, 0, 1, 1}
	almostEqualSlice(t, got, want, 1e-9)
}

// TestCrossCorrelateEdgeKnownValues is pinned to the no-flip
// cross-correlation example -- the exact sign-reversal of Convolve's
// result for this asymmetric kernel.
func TestCrossCorrelateEdgeKnownValues(t *testing.T) {
	got := CrossCorrelate(Signal, Kernels[1])
	want := []float64{1, 1, 0, 0, -1, -1}
	almostEqualSlice(t, got, want, 1e-9)
}

// TestSymmetricKernelFlipInvariant checks the concept's specific claim
// that flipping doesn't matter for the symmetric averaging kernel.
func TestSymmetricKernelFlipInvariant(t *testing.T) {
	conv := Convolve(Signal, Kernels[0])
	corr := CrossCorrelate(Signal, Kernels[0])
	almostEqualSlice(t, conv, corr, 1e-9)
}

// TestAsymmetricKernelFlipNegates checks the concept's other specific
// claim: for the edge kernel, flipping negates every output.
func TestAsymmetricKernelFlipNegates(t *testing.T) {
	conv := Convolve(Signal, Kernels[1])
	corr := CrossCorrelate(Signal, Kernels[1])
	if len(conv) != len(corr) {
		t.Fatalf("len(conv)=%d, len(corr)=%d", len(conv), len(corr))
	}
	for i := range conv {
		if !almostEqual(conv[i], -corr[i], 1e-9) {
			t.Errorf("[%d]: Convolve=%v, -CrossCorrelate=%v, want equal", i, conv[i], -corr[i])
		}
	}
}

func TestConvolveOutputLength(t *testing.T) {
	x := make([]float64, 8)
	k := make([]float64, 3)
	if got := len(Convolve(x, k)); got != 6 {
		t.Errorf("len(Convolve) = %d, want 6", got)
	}
	if got := len(CrossCorrelate(x, k)); got != 6 {
		t.Errorf("len(CrossCorrelate) = %d, want 6", got)
	}
}

func TestConvolveKernelLongerThanSignal(t *testing.T) {
	x := []float64{1, 2}
	k := []float64{1, 1, 1}
	if got := Convolve(x, k); got != nil {
		t.Errorf("Convolve(short x, long k) = %v, want nil", got)
	}
	if got := CrossCorrelate(x, k); got != nil {
		t.Errorf("CrossCorrelate(short x, long k) = %v, want nil", got)
	}
}

// TestConvolveMatchesBruteForceDefinition checks Convolve directly against
// the textbook sum-of-products-with-reversed-kernel definition, for
// several signals and kernels, independent of the pinned example.
func TestConvolveMatchesBruteForceDefinition(t *testing.T) {
	cases := []struct {
		x, k []float64
	}{
		{[]float64{1, 2, 3, 4, 5}, []float64{1, 0, -1}},
		{[]float64{0, 1, 0, -1, 0, 1}, []float64{0.5, 0.5}},
		{[]float64{2, 4, 6, 8}, []float64{1, 2, 3}},
	}
	for _, c := range cases {
		got := Convolve(c.x, c.k)
		n := len(c.x) - len(c.k) + 1
		want := make([]float64, n)
		for i := 0; i < n; i++ {
			sum := 0.0
			for j := 0; j < len(c.k); j++ {
				sum += c.x[i+j] * c.k[len(c.k)-1-j]
			}
			want[i] = sum
		}
		almostEqualSlice(t, got, want, 1e-9)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("convolution")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
