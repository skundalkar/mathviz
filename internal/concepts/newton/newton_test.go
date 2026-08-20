package newton

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestNewtonStepKnownValues(t *testing.T) {
	cases := []struct {
		x, want float64
	}{
		{1, 1.5},
		{1.5, 1.0 + 5.0/12.0}, // 1.5 - 0.25/3 = 1.41666...
	}
	for _, tc := range cases {
		got := NewtonStep(tc.x)
		if math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("NewtonStep(%v) = %v, want %v", tc.x, got, tc.want)
		}
	}
}

func TestIteratesConvergesToSqrt2(t *testing.T) {
	iters := Iterates(1, 5)
	if len(iters) != 6 {
		t.Fatalf("len(iters) = %d, want 6", len(iters))
	}
	if iters[0] != 1 {
		t.Fatalf("iters[0] = %v, want 1 (unchanged starting guess)", iters[0])
	}
	last := iters[len(iters)-1]
	if math.Abs(last-math.Sqrt2) > 1e-9 {
		t.Fatalf("after 5 steps from x0=1, got %v, want ≈ √2 = %v", last, math.Sqrt2)
	}
}

// TestErrorShrinksQuadratically checks the "each error is roughly the
// square of the previous one" claim from the docs, once iterates are
// already close to the root (quadratic convergence is a limiting
// behavior, not exact from the very first step).
func TestErrorShrinksQuadratically(t *testing.T) {
	iters := Iterates(1, 4)
	errAt := func(i int) float64 { return math.Abs(iters[i] - math.Sqrt2) }

	e2, e3, e4 := errAt(2), errAt(3), errAt(4)
	// e3 should be on the order of e2^2, and e4 on the order of e3^2 —
	// allow generous slack since this is an asymptotic property.
	if e3 > e2*e2*10 {
		t.Errorf("error at step 3 (%v) not roughly the square of step 2's error (%v)", e3, e2)
	}
	if e4 > e3*e3*10 {
		t.Errorf("error at step 4 (%v) not roughly the square of step 3's error (%v)", e4, e3)
	}
}

func TestNewtonStepFlatTangentDoesNotPanic(t *testing.T) {
	// f'(0) = 0, a flat tangent; NewtonStep must not divide by zero.
	got := NewtonStep(0)
	if got != 0 {
		t.Fatalf("NewtonStep(0) = %v, want 0 (unchanged, flat tangent)", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("newtons-method")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
