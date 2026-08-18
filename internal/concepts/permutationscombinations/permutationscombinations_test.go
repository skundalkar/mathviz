package permutationscombinations

import (
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestFactorialKnownValues(t *testing.T) {
	cases := map[int]int64{0: 1, 1: 1, 3: 6, 5: 120, 8: 40320}
	for n, want := range cases {
		if got := Factorial(n); got != want {
			t.Errorf("Factorial(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestFactorialNegative(t *testing.T) {
	if got := Factorial(-1); got != 0 {
		t.Errorf("Factorial(-1) = %v, want 0", got)
	}
}

func TestPermutationsKnownValue(t *testing.T) {
	// From LEARNINGS.md's worked example: 8 racers, medal for 3.
	if got := Permutations(8, 3); got != 336 {
		t.Errorf("Permutations(8, 3) = %v, want 336", got)
	}
}

func TestCombinationsKnownValue(t *testing.T) {
	// Same 8-racer example: unordered relay team of 3.
	if got := Combinations(8, 3); got != 56 {
		t.Errorf("Combinations(8, 3) = %v, want 56", got)
	}
}

func TestPermutationsEqualsCombinationsTimesFactorial(t *testing.T) {
	for n := 0; n <= 10; n++ {
		for k := 0; k <= n; k++ {
			p, c, f := Permutations(n, k), Combinations(n, k), Factorial(k)
			if p != c*f {
				t.Errorf("n=%d k=%d: Permutations=%v, want Combinations(%v)*Factorial(%v)=%v", n, k, p, c, k, c*f)
			}
		}
	}
}

func TestPermutationsAndCombinationsEdgeCases(t *testing.T) {
	for n := 0; n <= 6; n++ {
		if got := Permutations(n, 0); got != 1 {
			t.Errorf("Permutations(%d, 0) = %v, want 1", n, got)
		}
		if got := Combinations(n, 0); got != 1 {
			t.Errorf("Combinations(%d, 0) = %v, want 1", n, got)
		}
		if got := Combinations(n, n); got != 1 {
			t.Errorf("Combinations(%d, %d) = %v, want 1", n, n, got)
		}
	}
}

func TestPermutationsAndCombinationsOutOfRange(t *testing.T) {
	if got := Permutations(3, 4); got != 0 {
		t.Errorf("Permutations(3, 4) = %v, want 0", got)
	}
	if got := Combinations(3, 4); got != 0 {
		t.Errorf("Combinations(3, 4) = %v, want 0", got)
	}
	if got := Permutations(3, -1); got != 0 {
		t.Errorf("Permutations(3, -1) = %v, want 0", got)
	}
	if got := Combinations(-1, 0); got != 0 {
		t.Errorf("Combinations(-1, 0) = %v, want 0", got)
	}
}

func TestCombinationsMatchesPascalsTriangleFormula(t *testing.T) {
	// C(n,k) = n! / (k! * (n-k)!) — cross-check against the raw factorial
	// formula, the same formula pascals-triangle names in its own docs.
	for n := 0; n <= 8; n++ {
		for k := 0; k <= n; k++ {
			want := Factorial(n) / (Factorial(k) * Factorial(n-k))
			if got := Combinations(n, k); got != want {
				t.Errorf("Combinations(%d, %d) = %v, want %v (raw factorial formula)", n, k, got, want)
			}
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("permutations-vs-combinations")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
