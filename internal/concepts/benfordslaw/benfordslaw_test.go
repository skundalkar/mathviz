package benfordslaw

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestLeadingDigitKnownValues(t *testing.T) {
	cases := []struct {
		x    float64
		want int
	}{
		{1, 1}, {2, 2}, {9, 9}, {10, 1}, {99, 9}, {100, 1},
		{314, 3}, {0.0072, 7}, {9999, 9}, {1000, 1},
	}
	for _, c := range cases {
		if got := LeadingDigit(c.x); got != c.want {
			t.Errorf("LeadingDigit(%v) = %d, want %d", c.x, got, c.want)
		}
	}
}

func TestLeadingDigitOfPowersOfTwo(t *testing.T) {
	// 2^0=1, 2^1=2, ..., 2^10=1024, 2^11=2048, matching the worked example.
	want := []int{1, 2, 4, 8, 1, 3, 6, 1, 2, 5, 1, 2}
	for i, w := range want {
		x := math.Pow(2, float64(i))
		if got := LeadingDigit(x); got != w {
			t.Errorf("LeadingDigit(2^%d=%v) = %d, want %d", i, x, got, w)
		}
	}
}

func TestBenfordProbabilityKnownValues(t *testing.T) {
	// From the worked example: digit 1 owns log10(2)=0.301, digit 9 owns
	// log10(10/9)=0.0458.
	if got := BenfordProbability(1); !almostEqual(got, 0.30103, 1e-4) {
		t.Errorf("BenfordProbability(1) = %v, want ~0.30103", got)
	}
	if got := BenfordProbability(9); !almostEqual(got, 0.045757, 1e-5) {
		t.Errorf("BenfordProbability(9) = %v, want ~0.045757", got)
	}
}

func TestBenfordProbabilityOutOfRangeIsZero(t *testing.T) {
	if got := BenfordProbability(0); got != 0 {
		t.Errorf("BenfordProbability(0) = %v, want 0", got)
	}
	if got := BenfordProbability(10); got != 0 {
		t.Errorf("BenfordProbability(10) = %v, want 0", got)
	}
}

func TestBenfordProbabilitiesSumToOne(t *testing.T) {
	sum := 0.0
	for d := 1; d <= 9; d++ {
		sum += BenfordProbability(d)
	}
	if !almostEqual(sum, 1, 1e-9) {
		t.Errorf("sum of BenfordProbability(1..9) = %v, want 1", sum)
	}
}

func TestBenfordProbabilityIsDecreasing(t *testing.T) {
	for d := 1; d < 9; d++ {
		if BenfordProbability(d) <= BenfordProbability(d+1) {
			t.Errorf("BenfordProbability(%d)=%v should exceed BenfordProbability(%d)=%v",
				d, BenfordProbability(d), d+1, BenfordProbability(d+1))
		}
	}
}

func TestPowersOfBase(t *testing.T) {
	got := PowersOfBase(2, 5)
	want := []float64{1, 2, 4, 8, 16}
	for i, w := range want {
		if !almostEqual(got[i], w, 1e-9) {
			t.Errorf("PowersOfBase(2,5)[%d] = %v, want %v", i, got[i], w)
		}
	}
	if len(PowersOfBase(3, 0)) != 1 {
		t.Errorf("PowersOfBase(3,0) should clamp n to at least 1")
	}
}

func TestLinearSamples(t *testing.T) {
	got := LinearSamples(1, 999, 3)
	want := []float64{1, 500, 999}
	for i, w := range want {
		if !almostEqual(got[i], w, 1e-9) {
			t.Errorf("LinearSamples(1,999,3)[%d] = %v, want %v", i, got[i], w)
		}
	}
	if got := LinearSamples(5, 10, 1); len(got) != 1 || got[0] != 5 {
		t.Errorf("LinearSamples(5,10,1) = %v, want [5]", got)
	}
}

func TestLeadingDigitFrequenciesMatchesWorkedExample(t *testing.T) {
	// 30 powers of 2: digit 1 appears 9/30 times (see LEARNINGS.md).
	vals := PowersOfBase(2, 30)
	freq := LeadingDigitFrequencies(vals)
	if !almostEqual(freq[0], 9.0/30.0, 1e-9) {
		t.Errorf("LeadingDigitFrequencies(30 powers of 2)[1] = %v, want %v", freq[0], 9.0/30.0)
	}
	sum := 0.0
	for _, f := range freq {
		sum += f
	}
	if !almostEqual(sum, 1, 1e-9) {
		t.Errorf("frequencies should sum to 1, got %v", sum)
	}
}

func TestLeadingDigitFrequenciesEmptyInput(t *testing.T) {
	freq := LeadingDigitFrequencies(nil)
	for i, f := range freq {
		if f != 0 {
			t.Errorf("LeadingDigitFrequencies(nil)[%d] = %v, want 0", i, f)
		}
	}
}

func TestLinearSamplesStayCloseToUniform(t *testing.T) {
	// The contrasting non-multiplicative dataset: leading digits should
	// stay roughly flat (~11% each), nowhere near Benford's ~30% for digit 1.
	freq := LeadingDigitFrequencies(LinearSamples(1, 999, 300))
	if freq[0] > 0.15 {
		t.Errorf("digit-1 frequency for evenly spaced samples = %v, want close to uniform (~0.11)", freq[0])
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("benfords-law")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
