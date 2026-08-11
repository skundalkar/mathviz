package derivative

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestSecantSlopeKnownValues(t *testing.T) {
	cases := []struct {
		x0, h float64
		want  float64
	}{
		{1, 1, 3},       // (4-1)/1
		{1, 0.5, 2.5},   // (2.25-1)/0.5
		{1, 0.1, 2.1},   // (1.21-1)/0.1
		{1, 0.01, 2.01}, // (1.0201-1)/0.01
		{0, 1, 1},       // (1-0)/1
	}
	for _, tc := range cases {
		got := SecantSlope(tc.x0, tc.h)
		if math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("SecantSlope(%v, %v) = %v, want %v", tc.x0, tc.h, got, tc.want)
		}
	}
}

func TestSecantSlopeApproachesDerivative(t *testing.T) {
	x0 := 1.5
	want := Derivative(x0)
	prevGap := math.Inf(1)
	for _, h := range []float64{1, 0.1, 0.01, 0.001, 0.0001} {
		gap := math.Abs(SecantSlope(x0, h) - want)
		if gap >= prevGap {
			t.Fatalf("secant slope did not get closer to the derivative as h shrank: gap=%v at h=%v, previous gap=%v", gap, h, prevGap)
		}
		prevGap = gap
	}
	if prevGap > 1e-3 {
		t.Fatalf("secant slope at smallest h still far from derivative: gap=%v", prevGap)
	}
}

func TestDerivativeKnownValues(t *testing.T) {
	cases := []struct{ x0, want float64 }{
		{0, 0}, {1, 2}, {-1, -2}, {2.5, 5},
	}
	for _, tc := range cases {
		if got := Derivative(tc.x0); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("Derivative(%v) = %v, want %v", tc.x0, got, tc.want)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("derivative")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
