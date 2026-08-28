package rsa

import (
	"testing"

	"mathviz/internal/concept"
)

func TestKeyGenerationMatchesWorkedExample(t *testing.T) {
	if N != 3233 {
		t.Errorf("N = %d, want 3233", N)
	}
	if Phi != 3120 {
		t.Errorf("Phi = %d, want 3120", Phi)
	}
	if D != 2753 {
		t.Errorf("D = %d, want 2753", D)
	}
}

func TestEncryptDecryptMatchesWorkedExample(t *testing.T) {
	c := Encrypt(65, E, N)
	if c != 2790 {
		t.Errorf("Encrypt(65,17,3233) = %d, want 2790", c)
	}
	m := Decrypt(c, D, N)
	if m != 65 {
		t.Errorf("Decrypt(2790,2753,3233) = %d, want 65", m)
	}
}

func TestEncryptDecryptRoundTripsForManyMessages(t *testing.T) {
	for _, m := range []int{0, 1, 2, 65, 100, 1000, 2790, 3000, 3232} {
		c := Encrypt(m, E, N)
		got := Decrypt(c, D, N)
		if got != m {
			t.Errorf("round trip failed for m=%d: encrypted to %d, decrypted to %d", m, c, got)
		}
	}
}

func TestDecryptWithWrongExponentGenerallyFails(t *testing.T) {
	m := 65
	c := Encrypt(m, E, N)
	mismatches := 0
	for offset := 1; offset <= 5; offset++ {
		if Decrypt(c, D+offset, N) != m {
			mismatches++
		}
	}
	if mismatches != 5 {
		t.Errorf("wrong-exponent decryption accidentally matched m for %d/5 offsets, want 0 matches", 5-mismatches)
	}
}

func TestModPowKnownValues(t *testing.T) {
	cases := []struct{ base, exp, mod, want int }{
		{2, 10, 1000, 24},   // 1024 mod 1000
		{3, 0, 7, 1},        // anything^0 = 1
		{5, 1, 13, 5},       // anything^1 = itself mod m
		{65, 17, 3233, 2790},
		{7, 4, 1, 0}, // mod 1 is always 0
	}
	for _, c := range cases {
		if got := ModPow(c.base, c.exp, c.mod); got != c.want {
			t.Errorf("ModPow(%d,%d,%d) = %d, want %d", c.base, c.exp, c.mod, got, c.want)
		}
	}
}

func TestExtendedGCDSatisfiesBezoutIdentity(t *testing.T) {
	cases := [][2]int{{17, 3120}, {1071, 462}, {13, 8}, {48, 18}, {7, 7}}
	for _, pair := range cases {
		a, b := pair[0], pair[1]
		g, x, y := ExtendedGCD(a, b)
		if got := a*x + b*y; got != g {
			t.Errorf("ExtendedGCD(%d,%d): %d*%d + %d*%d = %d, want %d", a, b, a, x, b, y, got, g)
		}
	}
}

func TestExtendedGCDMatchesEuclideanAlgorithmGCD(t *testing.T) {
	// Same GCD values euclidalg's worked example uses, cross-checked here.
	cases := []struct{ a, b, want int }{
		{1071, 462, 21},
		{17, 3120, 1},
		{48, 18, 6},
	}
	for _, c := range cases {
		if g, _, _ := ExtendedGCD(c.a, c.b); g != c.want {
			t.Errorf("ExtendedGCD(%d,%d) gcd = %d, want %d", c.a, c.b, g, c.want)
		}
	}
}

func TestModInverseKnownValue(t *testing.T) {
	d, ok := ModInverse(17, 3120)
	if !ok {
		t.Fatal("ModInverse(17,3120) reported no inverse, want 2753")
	}
	if d != 2753 {
		t.Errorf("ModInverse(17,3120) = %d, want 2753", d)
	}
	if (17*d)%3120 != 1 {
		t.Errorf("17*%d mod 3120 = %d, want 1", d, (17*d)%3120)
	}
}

func TestModInverseReportsNoneWhenNotCoprime(t *testing.T) {
	if _, ok := ModInverse(4, 8); ok {
		t.Error("ModInverse(4,8) reported an inverse, want ok=false (gcd(4,8)=4, not 1)")
	}
}

func TestModInverseAlwaysReturnsAValueInRange(t *testing.T) {
	for e := 1; e < 20; e++ {
		if d, ok := ModInverse(e, 3120); ok {
			if d < 0 || d >= 3120 {
				t.Errorf("ModInverse(%d,3120) = %d, want a value in [0,3120)", e, d)
			}
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("rsa-encryption")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
