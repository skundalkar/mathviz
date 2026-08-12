package primesieve

import (
	"reflect"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestIsPrimeKnownValues(t *testing.T) {
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 97}
	for _, n := range primes {
		if !IsPrime(n) {
			t.Errorf("IsPrime(%d) = false, want true", n)
		}
	}
	composites := []int{0, 1, 4, 6, 8, 9, 15, 25, 91, 100}
	for _, n := range composites {
		if IsPrime(n) {
			t.Errorf("IsPrime(%d) = true, want false", n)
		}
	}
}

func TestIsPrimeNegativeIsFalse(t *testing.T) {
	if IsPrime(-7) {
		t.Error("IsPrime(-7) = true, want false")
	}
}

func TestSmallestPrimeFactorKnownValues(t *testing.T) {
	cases := []struct{ n, want int }{
		{4, 2}, {9, 3}, {15, 3}, {25, 5}, {91, 7}, {97, 97}, {2, 2},
	}
	for _, tc := range cases {
		if got := SmallestPrimeFactor(tc.n); got != tc.want {
			t.Errorf("SmallestPrimeFactor(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestMarkedAtStepMatchesWorkedExample(t *testing.T) {
	// From the LEARNINGS/Sections worked example, sieving 2..30: after
	// sweeping p=2, every even number from 4 is marked, odd numbers are not.
	if !MarkedAtStep(4, 2) {
		t.Error("MarkedAtStep(4, 2) = false, want true (4 = 2×2)")
	}
	if MarkedAtStep(9, 2) {
		t.Error("MarkedAtStep(9, 2) = true, want false (9's smallest factor is 3 > 2)")
	}
	if !MarkedAtStep(9, 3) {
		t.Error("MarkedAtStep(9, 3) = false, want true (9 = 3×3)")
	}
	if MarkedAtStep(25, 3) {
		t.Error("MarkedAtStep(25, 3) = true, want false (25's smallest factor is 5 > 3)")
	}
	if !MarkedAtStep(25, 5) {
		t.Error("MarkedAtStep(25, 5) = false, want true (25 = 5×5)")
	}
}

func TestMarkedAtStepNeverMarksAPrime(t *testing.T) {
	for _, n := range PrimesUpTo(50) {
		if MarkedAtStep(n, 50) {
			t.Errorf("MarkedAtStep(%d, 50) = true, want false — %d is prime", n, n)
		}
	}
}

func TestMarkedAtStepIsMonotonicInStep(t *testing.T) {
	// Once crossed off, a number stays crossed off as the sweep advances.
	for n := 4; n <= 100; n++ {
		wasMarked := false
		for step := 1; step <= 100; step++ {
			marked := MarkedAtStep(n, step)
			if wasMarked && !marked {
				t.Errorf("MarkedAtStep(%d, %d) = false after being true at a smaller step", n, step)
			}
			wasMarked = wasMarked || marked
		}
	}
}

func TestPrimesUpToKnownValues(t *testing.T) {
	got := PrimesUpTo(30)
	want := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PrimesUpTo(30) = %v, want %v", got, want)
	}
}

func TestPrimesUpToBelowTwoIsEmpty(t *testing.T) {
	if got := PrimesUpTo(1); len(got) != 0 {
		t.Errorf("PrimesUpTo(1) = %v, want empty", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("prime-sieve")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
