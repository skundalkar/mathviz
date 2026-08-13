package modulararithmetic

import (
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestModKnownValues(t *testing.T) {
	cases := []struct{ a, n, want int }{
		{14, 12, 2},  // the section-1 clock scenario: 9+5=14 -> 2 o'clock
		{0, 12, 0},
		{12, 12, 0},
		{25, 12, 1},
		{11, 12, 11},
	}
	for _, tc := range cases {
		if got := Mod(tc.a, tc.n); got != tc.want {
			t.Errorf("Mod(%d, %d) = %d, want %d", tc.a, tc.n, got, tc.want)
		}
	}
}

func TestModNegativeWrapsIntoRange(t *testing.T) {
	// -1 mod 12 is 11 (the hour right before midnight), not -1 — this is
	// exactly the "common mistake" the concept calls out.
	if got := Mod(-1, 12); got != 11 {
		t.Errorf("Mod(-1, 12) = %d, want 11", got)
	}
	if got := Mod(-13, 12); got != 11 {
		t.Errorf("Mod(-13, 12) = %d, want 11", got)
	}
}

func TestModAlwaysInRange(t *testing.T) {
	for _, n := range []int{3, 7, 12, 30} {
		for a := -50; a <= 50; a++ {
			got := Mod(a, n)
			if got < 0 || got >= n {
				t.Fatalf("Mod(%d, %d) = %d, want a value in [0, %d)", a, n, got, n)
			}
		}
	}
}

func TestAddModWorkedExample(t *testing.T) {
	// The section-1 clock scenario: it's 9 o'clock, wait 5 hours -> 2.
	if got := AddMod(9, 5, 12); got != 2 {
		t.Errorf("AddMod(9, 5, 12) = %d, want 2", got)
	}
}

func TestMulModMatchesWorkedTimesTable(t *testing.T) {
	// From the LEARNINGS/Sections worked example: n=12, multiplier=5.
	cases := []struct{ k, want int }{
		{0, 0}, {1, 5}, {2, 10}, {3, 3}, {4, 8}, {5, 1},
	}
	for _, tc := range cases {
		if got := MulMod(tc.k, 5, 12); got != tc.want {
			t.Errorf("MulMod(%d, 5, 12) = %d, want %d", tc.k, got, tc.want)
		}
	}
}

func TestMulModMultiplierOneIsIdentity(t *testing.T) {
	// multiplier=1 should map every point to itself (no chords in the
	// picture) - see the "What does the picture show?" section.
	for k := 0; k < 12; k++ {
		if got := MulMod(k, 1, 12); got != k {
			t.Errorf("MulMod(%d, 1, 12) = %d, want %d", k, got, k)
		}
	}
}

func TestPowModWorkedExample(t *testing.T) {
	// From the "What can you do now" section: 3^4 mod 7 = 4.
	if got := PowMod(3, 4, 7); got != 4 {
		t.Errorf("PowMod(3, 4, 7) = %d, want 4", got)
	}
}

func TestPowModMatchesNaiveExponentiation(t *testing.T) {
	for base := 2; base <= 6; base++ {
		for exp := 0; exp <= 6; exp++ {
			n := 13
			naive := 1
			for i := 0; i < exp; i++ {
				naive = (naive * base) % n
			}
			if got := PowMod(base, exp, n); got != naive {
				t.Errorf("PowMod(%d, %d, %d) = %d, want %d", base, exp, n, got, naive)
			}
		}
	}
}

func TestPowModExpZeroIsOne(t *testing.T) {
	if got := PowMod(5, 0, 7); got != 1 {
		t.Errorf("PowMod(5, 0, 7) = %d, want 1", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("modular-arithmetic")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
