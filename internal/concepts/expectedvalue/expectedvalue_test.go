package expectedvalue

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestTwoOutcomeCarnivalGame(t *testing.T) {
	// p=0.2, win=+15, lose=-5, worked by hand in LEARNINGS.md: -1.00.
	if got := TwoOutcome(15, -5, 0.2); !near(got, -1, 1e-9) {
		t.Errorf("TwoOutcome(15, -5, 0.2) = %v, want -1", got)
	}
}

func TestTwoOutcomeAtHigherP(t *testing.T) {
	// Flip win probability to 0.5: worked by hand, +5.00.
	if got := TwoOutcome(15, -5, 0.5); !near(got, 5, 1e-9) {
		t.Errorf("TwoOutcome(15, -5, 0.5) = %v, want 5", got)
	}
}

func TestDiscreteFairDie(t *testing.T) {
	// A fair six-sided die: E[X] = (1+2+3+4+5+6)/6 = 3.5.
	values := []float64{1, 2, 3, 4, 5, 6}
	probs := []float64{1.0 / 6, 1.0 / 6, 1.0 / 6, 1.0 / 6, 1.0 / 6, 1.0 / 6}
	if got := Discrete(values, probs); !near(got, 3.5, 1e-9) {
		t.Errorf("Discrete(fair die) = %v, want 3.5", got)
	}
}

func TestDiscreteMatchesTwoOutcome(t *testing.T) {
	// The general n-outcome formula must agree with the two-outcome
	// special case for the same scenario.
	got := Discrete([]float64{15, -5}, []float64{0.2, 0.8})
	want := TwoOutcome(15, -5, 0.2)
	if !near(got, want, 1e-9) {
		t.Errorf("Discrete two-outcome case = %v, want %v (TwoOutcome)", got, want)
	}
}

func TestBreakevenCarnivalGame(t *testing.T) {
	// Solving 15p + (-5)(1-p) = 0 by hand gives p = 0.25.
	p := Breakeven(15, -5)
	if !near(p, 0.25, 1e-9) {
		t.Errorf("Breakeven(15, -5) = %v, want 0.25", p)
	}
	if got := TwoOutcome(15, -5, p); !near(got, 0, 1e-9) {
		t.Errorf("TwoOutcome at breakeven p=%v = %v, want 0", p, got)
	}
}

func TestBreakevenEqualOutcomesIsNaN(t *testing.T) {
	if got := Breakeven(5, 5); !math.IsNaN(got) {
		t.Errorf("Breakeven(5, 5) = %v, want NaN (no single breakeven point)", got)
	}
}

func TestTwoOutcomeMonotonicInP(t *testing.T) {
	// Raising the win probability should only ever raise E[X] when win > lose.
	low := TwoOutcome(15, -5, 0.1)
	high := TwoOutcome(15, -5, 0.4)
	if high <= low {
		t.Errorf("TwoOutcome at p=0.4 (%v) should exceed p=0.1 (%v)", high, low)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("expected-value")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
