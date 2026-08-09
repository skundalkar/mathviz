package entropy

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestSurpriseKnownValues(t *testing.T) {
	// A coin flip outcome with p=0.5 carries exactly 1 bit of surprise.
	if got := Surprise(0.5); math.Abs(got-1) > 1e-9 {
		t.Fatalf("Surprise(0.5) = %v, want 1", got)
	}
	// A certain outcome carries zero surprise.
	if got := Surprise(1); math.Abs(got) > 1e-9 {
		t.Fatalf("Surprise(1) = %v, want 0", got)
	}
	// An outcome with p=0.25 carries 2 bits (log2(4)).
	if got := Surprise(0.25); math.Abs(got-2) > 1e-9 {
		t.Fatalf("Surprise(0.25) = %v, want 2", got)
	}
	if !math.IsInf(Surprise(0), 1) {
		t.Fatalf("Surprise(0) = %v, want +Inf", Surprise(0))
	}
}

func TestSurpriseDecreasesAsProbabilityRises(t *testing.T) {
	prev := math.Inf(1)
	for p := 0.05; p <= 1.0; p += 0.05 {
		s := Surprise(p)
		if s > prev+1e-9 {
			t.Fatalf("Surprise not monotonically decreasing: Surprise(%.2f)=%v > previous %v", p, s, prev)
		}
		prev = s
	}
}

func TestBinaryEntropyPeaksAtHalf(t *testing.T) {
	peak := BinaryEntropy(0.5)
	if math.Abs(peak-1) > 1e-9 {
		t.Fatalf("BinaryEntropy(0.5) = %v, want 1", peak)
	}
	for _, p := range []float64{0.01, 0.1, 0.3, 0.7, 0.9, 0.99} {
		if h := BinaryEntropy(p); h > peak+1e-12 {
			t.Fatalf("BinaryEntropy(%v) = %v exceeded the peak at 0.5 (%v)", p, h, peak)
		}
	}
}

func TestBinaryEntropyZeroAtExtremes(t *testing.T) {
	if got := BinaryEntropy(0); math.Abs(got) > 1e-9 {
		t.Fatalf("BinaryEntropy(0) = %v, want 0", got)
	}
	if got := BinaryEntropy(1); math.Abs(got) > 1e-9 {
		t.Fatalf("BinaryEntropy(1) = %v, want 0", got)
	}
}

func TestBinaryEntropySymmetric(t *testing.T) {
	for _, p := range []float64{0.05, 0.2, 0.35, 0.5} {
		a, b := BinaryEntropy(p), BinaryEntropy(1-p)
		if math.Abs(a-b) > 1e-9 {
			t.Fatalf("BinaryEntropy(%v)=%v != BinaryEntropy(%v)=%v", p, a, 1-p, b)
		}
	}
}

func TestBinaryEntropyMatchesWeightedSurprise(t *testing.T) {
	// H(p) should equal the probability-weighted average of both outcomes'
	// surprise: p*Surprise(p) + (1-p)*Surprise(1-p).
	for _, p := range []float64{0.1, 0.3, 0.5, 0.8} {
		want := p*Surprise(p) + (1-p)*Surprise(1-p)
		got := BinaryEntropy(p)
		if math.Abs(got-want) > 1e-9 {
			t.Fatalf("BinaryEntropy(%v) = %v, want weighted-surprise %v", p, got, want)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("entropy")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
