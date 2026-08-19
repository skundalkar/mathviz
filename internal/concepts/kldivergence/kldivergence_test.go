package kldivergence

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestBinaryEntropyKnownValues(t *testing.T) {
	if got := BinaryEntropy(0.5); !almostEqual(got, 1.0, 1e-9) {
		t.Errorf("BinaryEntropy(0.5) = %v, want 1.0", got)
	}
	// From LEARNINGS.md: H(0.9) = 0.469 bits.
	if got := BinaryEntropy(0.9); !almostEqual(got, 0.468996, 1e-5) {
		t.Errorf("BinaryEntropy(0.9) = %v, want ~0.468996", got)
	}
}

func TestBinaryEntropyBoundaryIsZero(t *testing.T) {
	if got := BinaryEntropy(0); got != 0 {
		t.Errorf("BinaryEntropy(0) = %v, want 0", got)
	}
	if got := BinaryEntropy(1); got != 0 {
		t.Errorf("BinaryEntropy(1) = %v, want 0", got)
	}
}

func TestCrossEntropyKnownValue(t *testing.T) {
	// From LEARNINGS.md: H(P,Q) with p=0.9, q=0.5 = 1 bit exactly.
	if got := CrossEntropy(0.9, 0.5); !almostEqual(got, 1.0, 1e-9) {
		t.Errorf("CrossEntropy(0.9, 0.5) = %v, want 1.0", got)
	}
	// H(Q,P) with p=0.5, q=0.9 = 1.737 bits.
	if got := CrossEntropy(0.5, 0.9); !almostEqual(got, 1.736966, 1e-5) {
		t.Errorf("CrossEntropy(0.5, 0.9) = %v, want ~1.736966", got)
	}
}

func TestCrossEntropyEqualsEntropyWhenPEqualsQ(t *testing.T) {
	// Using the true distribution's own code costs exactly its own entropy
	// -- no waste when p == q.
	for _, p := range []float64{0.1, 0.3, 0.5, 0.7, 0.9} {
		ce, h := CrossEntropy(p, p), BinaryEntropy(p)
		if !almostEqual(ce, h, 1e-9) {
			t.Errorf("CrossEntropy(%v,%v) = %v, want BinaryEntropy(%v) = %v", p, p, ce, p, h)
		}
	}
}

func TestCrossEntropyZeroProbabilityIsInfinite(t *testing.T) {
	// Q assigns 0 probability to an outcome P actually produces: infinite cost.
	if got := CrossEntropy(0.9, 1.0); !math.IsInf(got, 1) {
		t.Errorf("CrossEntropy(0.9, 1.0) = %v, want +Inf", got)
	}
}

func TestKLDivergenceKnownValues(t *testing.T) {
	// From LEARNINGS.md: KL(P||Q) at p=0.9, q=0.5 ~= 0.531 bits.
	if got := KLDivergence(0.9, 0.5); !almostEqual(got, 0.531004, 1e-4) {
		t.Errorf("KLDivergence(0.9, 0.5) = %v, want ~0.531004", got)
	}
	// KL(Q||P), swapped: p=0.5, q=0.9 ~= 0.737 bits -- a different number.
	if got := KLDivergence(0.5, 0.9); !almostEqual(got, 0.736966, 1e-4) {
		t.Errorf("KLDivergence(0.5, 0.9) = %v, want ~0.736966", got)
	}
}

func TestKLDivergenceIsAsymmetric(t *testing.T) {
	// The central claim this package exists to demonstrate: swapping which
	// distribution is "true" gives a different divergence.
	forward := KLDivergence(0.9, 0.5)
	backward := KLDivergence(0.5, 0.9)
	if almostEqual(forward, backward, 1e-6) {
		t.Errorf("KLDivergence(0.9,0.5) = %v and KLDivergence(0.5,0.9) = %v, want them to differ", forward, backward)
	}
}

func TestKLDivergenceIsZeroWhenDistributionsMatch(t *testing.T) {
	for _, p := range []float64{0.1, 0.3, 0.5, 0.7, 0.9} {
		if got := KLDivergence(p, p); !almostEqual(got, 0, 1e-9) {
			t.Errorf("KLDivergence(%v, %v) = %v, want 0", p, p, got)
		}
	}
}

func TestKLDivergenceIsNeverNegative(t *testing.T) {
	// Gibbs' inequality: KL divergence is always >= 0, for any pair of
	// distributions.
	for _, p := range []float64{0.05, 0.2, 0.5, 0.8, 0.95} {
		for _, q := range []float64{0.05, 0.2, 0.5, 0.8, 0.95} {
			if got := KLDivergence(p, q); got < -1e-9 {
				t.Errorf("KLDivergence(%v, %v) = %v, want >= 0", p, q, got)
			}
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("kl-divergence")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
