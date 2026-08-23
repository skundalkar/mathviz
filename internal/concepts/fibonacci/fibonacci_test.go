package fibonacci

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestFibonacciKnownValues(t *testing.T) {
	cases := map[int]int64{
		0: 0, 1: 1, 2: 1, 3: 2, 4: 3, 5: 5, 10: 55, 15: 610, 16: 987,
	}
	for n, want := range cases {
		if got := Fibonacci(n); got != want {
			t.Errorf("Fibonacci(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestFibonacciNegativeIsZero(t *testing.T) {
	if got := Fibonacci(-3); got != 0 {
		t.Errorf("Fibonacci(-3) = %v, want 0", got)
	}
}

func TestRatioMatchesWorkedExamples(t *testing.T) {
	// Worked by hand in LEARNINGS.md.
	cases := []struct {
		n    int
		want float64
	}{
		{1, 1.0},
		{5, 1.6},
		{10, 89.0 / 55.0},
		{15, 987.0 / 610.0},
	}
	for _, c := range cases {
		if got := Ratio(c.n); !near(got, c.want, 1e-9) {
			t.Errorf("Ratio(%d) = %v, want %v", c.n, got, c.want)
		}
	}
}

func TestRatioUndefinedAtZeroIsZero(t *testing.T) {
	if got := Ratio(0); got != 0 {
		t.Errorf("Ratio(0) = %v, want 0", got)
	}
}

func TestRatioObeysFixedPointRecurrence(t *testing.T) {
	// r(n) = 1 + 1/r(n-1) must hold exactly for every n>=2, straight from
	// the Fibonacci recurrence F(n+1)=F(n)+F(n-1).
	for n := 2; n <= 20; n++ {
		got := Ratio(n)
		want := FixedPointStep(Ratio(n - 1))
		if !near(got, want, 1e-9) {
			t.Errorf("Ratio(%d) = %v, want FixedPointStep(Ratio(%d)) = %v", n, got, n-1, want)
		}
	}
}

func TestPhiSolvesItsOwnQuadratic(t *testing.T) {
	// Phi must solve x^2 - x - 1 = 0, the equation derived in LEARNINGS.md.
	if got := Phi*Phi - Phi - 1; !near(got, 0, 1e-9) {
		t.Errorf("Phi^2 - Phi - 1 = %v, want 0", got)
	}
}

func TestPhiIsTheFixedPoint(t *testing.T) {
	if got := FixedPointStep(Phi); !near(got, Phi, 1e-9) {
		t.Errorf("FixedPointStep(Phi) = %v, want Phi = %v", got, Phi)
	}
}

func TestRatioConvergesToPhiWithShrinkingError(t *testing.T) {
	// The magnitude of the gap to Phi must shrink every single step -- the
	// oscillating-but-narrowing convergence described in LEARNINGS.md.
	prevErr := math.Abs(Ratio(1) - Phi)
	for n := 2; n <= 20; n++ {
		err := math.Abs(Ratio(n) - Phi)
		if err >= prevErr {
			t.Errorf("|Ratio(%d)-Phi| = %v did not shrink from |Ratio(%d)-Phi| = %v", n, err, n-1, prevErr)
		}
		prevErr = err
	}
	if prevErr > 1e-3 {
		t.Errorf("Ratio(20) still %v away from Phi, want a much tighter approximation by n=20", prevErr)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("fibonacci-golden-ratio")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
