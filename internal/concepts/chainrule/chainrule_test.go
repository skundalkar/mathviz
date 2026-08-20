package chainrule

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestWorkedExampleKnownValues(t *testing.T) {
	rate, t0 := 0.5, 4.0

	if r := Radius(rate, t0); math.Abs(r-2) > 1e-9 {
		t.Fatalf("Radius(0.5, 4) = %v, want 2", r)
	}
	if dr := RadiusPrime(rate); dr != rate {
		t.Fatalf("RadiusPrime(0.5) = %v, want 0.5", dr)
	}

	wantDVdr := 4 * math.Pi * 2 * 2 // 16π
	if dv := VolumePrime(2); math.Abs(dv-wantDVdr) > 1e-9 {
		t.Fatalf("VolumePrime(2) = %v, want %v", dv, wantDVdr)
	}

	wantDVdt := wantDVdr * rate // ≈ 25.133
	if got := ComposedPrime(rate, t0); math.Abs(got-wantDVdt) > 1e-9 {
		t.Fatalf("ComposedPrime(0.5, 4) = %v, want %v", got, wantDVdt)
	}
	if got := ComposedPrime(rate, t0); math.Abs(got-25.132741228718345) > 1e-6 {
		t.Fatalf("ComposedPrime(0.5, 4) = %v, want ≈ 25.1327", got)
	}
}

// TestChainRuleMatchesFiniteDifference checks the chain rule's derivative
// against a direct numerical derivative of Composed itself, at several
// (rate, t) points — the same cross-check the docs mention.
func TestChainRuleMatchesFiniteDifference(t *testing.T) {
	cases := []struct{ rate, t float64 }{
		{0.5, 4}, {1.2, 1.5}, {0.3, 5.5}, {0.8, 0.5},
	}
	const h = 1e-5
	for _, tc := range cases {
		numeric := (Composed(tc.rate, tc.t+h) - Composed(tc.rate, tc.t-h)) / (2 * h)
		analytic := ComposedPrime(tc.rate, tc.t)
		if math.Abs(numeric-analytic) > 1e-3 {
			t.Errorf("rate=%v t=%v: numeric derivative %v vs chain-rule %v differ too much",
				tc.rate, tc.t, numeric, analytic)
		}
	}
}

func TestRadiusPrimeConstantAcrossTime(t *testing.T) {
	// dr/dt = rate at every t, since Radius is linear in t.
	rate := 0.7
	for _, tt := range []float64{0, 1, 3.5, 10} {
		if got := RadiusPrime(rate); got != rate {
			t.Errorf("RadiusPrime(%v) at t=%v = %v, want %v", rate, tt, got, rate)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("chain-rule")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
