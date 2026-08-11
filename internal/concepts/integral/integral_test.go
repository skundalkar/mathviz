package integral

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestRiemannSumKnownValues(t *testing.T) {
	cases := []struct {
		n         int
		samplePos float64
		want      float64
	}{
		{1, 0, 0},   // v(0)*2 = 0
		{2, 0, 2},   // v(0)*1 + v(1)*1 = 0 + 2
		{4, 0, 3},   // 0.5*(v(0)+v(0.5)+v(1)+v(1.5)) = 0.5*(0+1+2+3)
		{1, 1, 8},   // v(2)*2 = 4*2
		{2, 1, 6},   // v(1)*1 + v(2)*1 = 2 + 4
		{1, 0.5, 4}, // midpoint of a straight line is exact even at n=1: v(1)*2 = 2*2
	}
	for _, tc := range cases {
		got := RiemannSum(0, 2, tc.n, tc.samplePos)
		if math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("RiemannSum(0,2,%d,%v) = %v, want %v", tc.n, tc.samplePos, got, tc.want)
		}
	}
}

func TestRiemannSumClampsNBelowOne(t *testing.T) {
	if got, want := RiemannSum(0, 2, 0, 0), RiemannSum(0, 2, 1, 0); got != want {
		t.Fatalf("RiemannSum with n=0 = %v, want same as n=1 (%v)", got, want)
	}
}

func TestLeftSumUnderestimatesRightSumOverestimates(t *testing.T) {
	// V is strictly increasing, so a left-edge sample sits below the true
	// area and a right-edge sample sits above it, for any n.
	trueVal := TrueIntegral(0, 2)
	for _, n := range []int{1, 2, 4, 10} {
		left := RiemannSum(0, 2, n, 0)
		right := RiemannSum(0, 2, n, 1)
		if left > trueVal+1e-9 {
			t.Errorf("n=%d: left sum %v exceeds true integral %v", n, left, trueVal)
		}
		if right < trueVal-1e-9 {
			t.Errorf("n=%d: right sum %v is below true integral %v", n, right, trueVal)
		}
	}
}

func TestRiemannSumConvergesToTrueIntegral(t *testing.T) {
	trueVal := TrueIntegral(0, 2)
	prevGap := math.Inf(1)
	for _, n := range []int{1, 10, 100, 1000, 10000} {
		gap := math.Abs(RiemannSum(0, 2, n, 0) - trueVal)
		if gap > prevGap+1e-12 {
			t.Fatalf("left sum did not get closer to the true integral as n grew: gap=%v at n=%d, previous gap=%v", gap, n, prevGap)
		}
		prevGap = gap
	}
	if prevGap > 1e-3 {
		t.Fatalf("left sum at largest n still far from true integral: gap=%v", prevGap)
	}
}

func TestTrueIntegralKnownValues(t *testing.T) {
	cases := []struct{ a, b, want float64 }{
		{0, 2, 4}, {0, 1, 1}, {1, 3, 8}, {0, 0, 0},
	}
	for _, tc := range cases {
		if got := TrueIntegral(tc.a, tc.b); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("TrueIntegral(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("integral")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
