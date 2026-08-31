package diffiehellman

import (
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestPublicValueMatchesWorkedExample(t *testing.T) {
	// p=23, g=5, a=6, b=15 -> A=8, B=19. See LEARNINGS.md.
	if got := PublicValue(6, G, P); got != 8 {
		t.Errorf("PublicValue(6,5,23) = %v, want 8", got)
	}
	if got := PublicValue(15, G, P); got != 19 {
		t.Errorf("PublicValue(15,5,23) = %v, want 19", got)
	}
}

func TestSharedSecretMatchesWorkedExample(t *testing.T) {
	A := PublicValue(6, G, P)
	B := PublicValue(15, G, P)
	if got := SharedSecret(B, 6, P); got != 2 {
		t.Errorf("SharedSecret(B,6,23) = %v, want 2", got)
	}
	if got := SharedSecret(A, 15, P); got != 2 {
		t.Errorf("SharedSecret(A,15,23) = %v, want 2", got)
	}
}

func TestBothSidesAgreeAcrossManyPrivateNumbers(t *testing.T) {
	// The whole point: whatever a and b are, B^a mod p must equal A^b mod
	// p, since both reduce to g^(a*b) mod p.
	for a := 1; a < P; a++ {
		for b := 1; b < P; b++ {
			A := PublicValue(a, G, P)
			B := PublicValue(b, G, P)
			sA := SharedSecret(B, a, P)
			sB := SharedSecret(A, b, P)
			if sA != sB {
				t.Fatalf("a=%d b=%d: Alice's secret %d != Bob's secret %d", a, b, sA, sB)
			}
		}
	}
}

func TestBruteForceDiscreteLogRecoversKnownExponent(t *testing.T) {
	target := PublicValue(6, G, P)
	x, ok := BruteForceDiscreteLog(G, target, P)
	if !ok {
		t.Fatal("BruteForceDiscreteLog reported no match for a known exponent")
	}
	// 6 isn't necessarily the only x with g^x mod p == target (the search
	// returns the smallest), but whatever it returns must itself satisfy
	// the equation.
	if got := ModPow(G, x, P); got != target {
		t.Errorf("ModPow(%d,%d,%d) = %v, want %v (target)", G, x, P, got, target)
	}
}

func TestBruteForceDiscreteLogNoMatchOutsideRange(t *testing.T) {
	// 0 is never a power of a nonzero base mod a prime > 1, so no exponent
	// in [1,p-1] can produce it.
	if _, ok := BruteForceDiscreteLog(G, 0, P); ok {
		t.Error("BruteForceDiscreteLog reported a match for target=0, want none")
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("diffie-hellman-key-exchange")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
