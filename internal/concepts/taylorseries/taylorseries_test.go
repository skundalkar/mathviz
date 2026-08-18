package taylorseries

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestSinDerivativeAtZeroCycles(t *testing.T) {
	want := []float64{0, 1, 0, -1, 0, 1, 0, -1}
	for k, w := range want {
		if got := SinDerivativeAtZero(k); got != w {
			t.Errorf("SinDerivativeAtZero(%d) = %v, want %v", k, got, w)
		}
	}
}

func TestSinDerivativeAtZeroNegative(t *testing.T) {
	if got := SinDerivativeAtZero(-1); got != 0 {
		t.Errorf("SinDerivativeAtZero(-1) = %v, want 0", got)
	}
}

func TestFactorialKnownValues(t *testing.T) {
	cases := map[int]float64{0: 1, 1: 1, 2: 2, 3: 6, 5: 120, 7: 5040}
	for k, want := range cases {
		if got := Factorial(k); got != want {
			t.Errorf("Factorial(%d) = %v, want %v", k, got, want)
		}
	}
}

func TestTaylorPolynomialKnownValuesAtXEquals1(t *testing.T) {
	// From LEARNINGS.md's worked example: sin(1) = 0.841471.
	cases := map[int]float64{
		1: 1.000000,
		3: 0.833333,
		5: 0.841667,
		7: 0.841468,
	}
	for order, want := range cases {
		if got := TaylorPolynomial(order, 1); !almostEqual(got, want, 1e-5) {
			t.Errorf("TaylorPolynomial(%d, 1) = %v, want %v", order, got, want)
		}
	}
}

func TestTaylorPolynomialKnownValuesAtXEquals3(t *testing.T) {
	// From LEARNINGS.md's worked example: sin(3) = 0.141120.
	cases := map[int]float64{
		1: 3.000000,
		5: 0.525000,
		9: 0.145312,
	}
	for order, want := range cases {
		if got := TaylorPolynomial(order, 3); !almostEqual(got, want, 1e-5) {
			t.Errorf("TaylorPolynomial(%d, 3) = %v, want %v", order, got, want)
		}
	}
}

func TestTaylorPolynomialConvergesToF(t *testing.T) {
	// Higher order at a fixed x should get closer to the true value.
	for _, x := range []float64{0.5, 1.5, 2.5} {
		errLow := math.Abs(TaylorPolynomial(3, x) - F(x))
		errHigh := math.Abs(TaylorPolynomial(11, x) - F(x))
		if errHigh > errLow {
			t.Errorf("x=%v: higher-order error %v not smaller than lower-order error %v", x, errHigh, errLow)
		}
	}
}

func TestTaylorPolynomialAtCenterMatchesFExactly(t *testing.T) {
	// At x0=0 itself, every order matches F(0)=0 exactly (all terms are 0
	// since x^k=0 for k>=1, and the k=0 term is SinDerivativeAtZero(0)=0).
	for order := 0; order <= 11; order++ {
		if got := TaylorPolynomial(order, 0); got != 0 {
			t.Errorf("TaylorPolynomial(%d, 0) = %v, want 0", order, got)
		}
	}
}

func TestTaylorPolynomialNegativeOrderTreatedAsZero(t *testing.T) {
	if got, want := TaylorPolynomial(-1, 1), TaylorPolynomial(0, 1); got != want {
		t.Errorf("TaylorPolynomial(-1, 1) = %v, want %v (same as order 0)", got, want)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("taylor-series")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
