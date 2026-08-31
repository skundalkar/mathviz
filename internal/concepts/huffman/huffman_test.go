package huffman

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestFrequenciesUniformAtSkewOne(t *testing.T) {
	labels, freqs := Frequencies(1)
	if len(labels) != 4 || len(freqs) != 4 {
		t.Fatalf("Frequencies(1) returned %d labels, %d freqs, want 4 each", len(labels), len(freqs))
	}
	for i, f := range freqs {
		if !almostEqual(f, 0.25, 1e-9) {
			t.Errorf("freqs[%d] = %v, want 0.25 (uniform at skew=1)", i, f)
		}
	}
}

func TestFrequenciesSumToOne(t *testing.T) {
	for _, skew := range []float64{0.15, 0.3, 0.5, 0.7, 0.95} {
		_, freqs := Frequencies(skew)
		sum := 0.0
		for _, f := range freqs {
			sum += f
		}
		if !almostEqual(sum, 1, 1e-9) {
			t.Errorf("skew=%v: frequencies sum to %v, want 1", skew, sum)
		}
	}
}

func TestFrequenciesDescendingForSkewBelowOne(t *testing.T) {
	_, freqs := Frequencies(0.5)
	for i := 1; i < len(freqs); i++ {
		if freqs[i] >= freqs[i-1] {
			t.Errorf("freqs[%d]=%v should be smaller than freqs[%d]=%v for skew<1", i, freqs[i], i-1, freqs[i-1])
		}
	}
}

// TestHuffmanMatchesWorkedExample checks the skew=0.5 worked example from
// LEARNINGS.md: probabilities 8/15, 4/15, 2/15, 1/15 for A,B,C,D produce
// code lengths 1,2,3,3 with average 25/15 = 1.6667 bits.
func TestHuffmanMatchesWorkedExample(t *testing.T) {
	labels, freqs := Frequencies(0.5)
	root := BuildHuffmanTree(labels, freqs)
	lengths := CodeLengths(root)

	want := map[string]int{"A": 1, "B": 2, "C": 3, "D": 3}
	for label, wantLen := range want {
		if got := lengths[label]; got != wantLen {
			t.Errorf("CodeLengths()[%q] = %v, want %v", label, got, wantLen)
		}
	}

	avg := AverageCodeLength(labels, freqs, lengths)
	if want := 25.0 / 15.0; !almostEqual(avg, want, 1e-9) {
		t.Errorf("AverageCodeLength = %v, want %v", avg, want)
	}
}

func TestCodesArePrefixFreeAndMatchLengths(t *testing.T) {
	labels, freqs := Frequencies(0.4)
	root := BuildHuffmanTree(labels, freqs)
	codes := Codes(root)
	lengths := CodeLengths(root)

	for _, l := range labels {
		if len(codes[l]) != lengths[l] {
			t.Errorf("len(Codes()[%q])=%d != CodeLengths()[%q]=%d", l, len(codes[l]), l, lengths[l])
		}
	}
	// Prefix-free: no codeword may be a prefix of another.
	for _, l1 := range labels {
		for _, l2 := range labels {
			if l1 == l2 {
				continue
			}
			if strings.HasPrefix(codes[l2], codes[l1]) {
				t.Errorf("code %q for %q is a prefix of code %q for %q", codes[l1], l1, codes[l2], l2)
			}
		}
	}
}

func TestAverageCodeLengthWithinOneBitOfEntropy(t *testing.T) {
	// Shannon's source coding theorem: Huffman's average length is always
	// in [H, H+1) bits for any distribution.
	for _, skew := range []float64{0.15, 0.3, 0.5, 0.7, 0.95} {
		labels, freqs := Frequencies(skew)
		root := BuildHuffmanTree(labels, freqs)
		lengths := CodeLengths(root)
		avg := AverageCodeLength(labels, freqs, lengths)
		h := EntropyBits(freqs)
		if avg < h-1e-9 {
			t.Errorf("skew=%v: avg=%v below entropy=%v (impossible)", skew, avg, h)
		}
		if avg >= h+1 {
			t.Errorf("skew=%v: avg=%v not within 1 bit of entropy=%v", skew, avg, h)
		}
	}
}

func TestEntropyBitsUniformFourWay(t *testing.T) {
	// A uniform 4-way choice has entropy log2(4) = 2 bits.
	freqs := []float64{0.25, 0.25, 0.25, 0.25}
	if got := EntropyBits(freqs); !almostEqual(got, 2, 1e-9) {
		t.Errorf("EntropyBits(uniform 4-way) = %v, want 2", got)
	}
}

func TestFixedLengthBits(t *testing.T) {
	if got := FixedLengthBits(4); got != 2 {
		t.Errorf("FixedLengthBits(4) = %v, want 2", got)
	}
	if got := FixedLengthBits(5); got != 3 {
		t.Errorf("FixedLengthBits(5) = %v, want 3", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("huffman-coding")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
