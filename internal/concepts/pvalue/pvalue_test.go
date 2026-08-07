package pvalue

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestStdNormalPDFKnownValues(t *testing.T) {
	got := StdNormalPDF(0)
	want := 1 / math.Sqrt(2*math.Pi)
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("StdNormalPDF(0) = %v, want %v", got, want)
	}
	if StdNormalPDF(0) <= StdNormalPDF(1) {
		t.Error("density at 0 should exceed density at 1 (peak at zero)")
	}
}

func TestPValueAtZeroIsOne(t *testing.T) {
	// An observed statistic of exactly 0 is the most "unsurprising" possible
	// result under the null — everything is at least as extreme, so p=1.
	if got := PValue(0); math.Abs(got-1) > 1e-12 {
		t.Errorf("PValue(0) = %v, want 1", got)
	}
}

func TestPValueSymmetric(t *testing.T) {
	for _, z := range []float64{0.3, 1.0, 1.96, 3.5} {
		pos, neg := PValue(z), PValue(-z)
		if math.Abs(pos-neg) > 1e-12 {
			t.Errorf("PValue(%v)=%v != PValue(%v)=%v", z, pos, -z, neg)
		}
	}
}

func TestPValueKnownValue(t *testing.T) {
	// The classic 1.96 <-> 0.05 two-sided correspondence.
	if got, want := PValue(1.959963985), 0.05; math.Abs(got-want) > 1e-6 {
		t.Errorf("PValue(1.96) = %v, want %v", got, want)
	}
}

func TestPValueDecreasesAsStatisticMovesFromZero(t *testing.T) {
	prev := PValue(0)
	for _, z := range []float64{0.5, 1, 1.5, 2, 3} {
		got := PValue(z)
		if got >= prev {
			t.Fatalf("PValue should decrease as |z| grows: PValue(%v)=%v >= previous %v", z, got, prev)
		}
		prev = got
	}
}

func TestCriticalZKnownValues(t *testing.T) {
	cases := map[float64]float64{
		0.10: 1.6448536,
		0.05: 1.9599640,
		0.01: 2.5758293,
	}
	for alpha, want := range cases {
		if got := CriticalZ(alpha); math.Abs(got-want) > 1e-5 {
			t.Errorf("CriticalZ(%v) = %v, want %v", alpha, got, want)
		}
	}
}

func TestCriticalZMatchesPValueBoundary(t *testing.T) {
	// By construction, the p-value of the critical z should equal alpha.
	for _, alpha := range []float64{0.2, 0.1, 0.05, 0.01} {
		z := CriticalZ(alpha)
		if got := PValue(z); math.Abs(got-alpha) > 1e-6 {
			t.Errorf("PValue(CriticalZ(%v))=%v, want %v", alpha, got, alpha)
		}
	}
}

func TestSignificant(t *testing.T) {
	if !Significant(0.03, 0.05) {
		t.Error("p=0.03 should be significant at alpha=0.05")
	}
	if Significant(0.10, 0.05) {
		t.Error("p=0.10 should not be significant at alpha=0.05")
	}
	if Significant(0.05, 0.05) {
		t.Error("p exactly equal to alpha should not count as significant (strict <)")
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("p-value")
	if !ok {
		t.Fatal("concept \"p-value\" not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
