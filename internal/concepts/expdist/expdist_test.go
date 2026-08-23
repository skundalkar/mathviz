package expdist

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestPDFKnownValue(t *testing.T) {
	// f(0) = lambda * e^0 = lambda, exactly, for any valid rate.
	if got := PDF(3, 0); !near(got, 3, 1e-9) {
		t.Errorf("PDF(3, 0) = %v, want 3", got)
	}
	// f(1) at lambda=3, worked by hand in LEARNINGS.md: 3*e^-3 ~= 0.14936.
	if got := PDF(3, 1); !near(got, 0.149361, 1e-5) {
		t.Errorf("PDF(3, 1) = %v, want ~0.149361", got)
	}
}

func TestPDFInvalidInputsAreZero(t *testing.T) {
	if got := PDF(3, -1); got != 0 {
		t.Errorf("PDF(3, -1) = %v, want 0 (a wait can't be negative)", got)
	}
	if got := PDF(-1, 1); got != 0 {
		t.Errorf("PDF(-1, 1) = %v, want 0 (not a valid rate)", got)
	}
}

func TestSurvivalMatchesWorkedExample(t *testing.T) {
	// P(wait > 1 hour) at lambda=3, worked by hand in LEARNINGS.md: e^-3 ~= 0.0498.
	if got := Survival(3, 1); !near(got, 0.049787, 1e-5) {
		t.Errorf("Survival(3, 1) = %v, want ~0.049787", got)
	}
}

func TestCDFAndSurvivalSumToOne(t *testing.T) {
	for _, lambda := range []float64{0.5, 1, 3, 5} {
		for _, tt := range []float64{0, 0.5, 1, 2, 10} {
			cdf, surv := CDF(lambda, tt), Survival(lambda, tt)
			if !near(cdf+surv, 1, 1e-9) {
				t.Errorf("CDF(%v,%v)+Survival(%v,%v) = %v, want 1", lambda, tt, lambda, tt, cdf+surv)
			}
		}
	}
}

func TestCDFMatchesWorkedExample(t *testing.T) {
	// P(wait <= 1 hour) at lambda=3, worked by hand in LEARNINGS.md: 1-e^-3 ~= 0.9502.
	if got := CDF(3, 1); !near(got, 0.950213, 1e-5) {
		t.Errorf("CDF(3, 1) = %v, want ~0.950213", got)
	}
}

func TestMeanIsReciprocalOfRate(t *testing.T) {
	if got := Mean(3); !near(got, 1.0/3.0, 1e-9) {
		t.Errorf("Mean(3) = %v, want 1/3", got)
	}
	if got := Mean(0.5); !near(got, 2, 1e-9) {
		t.Errorf("Mean(0.5) = %v, want 2", got)
	}
	if got := Mean(0); got != 0 {
		t.Errorf("Mean(0) = %v, want 0 (not a valid rate)", got)
	}
}

func TestMemorylessness(t *testing.T) {
	// P(wait > s+t | wait > s) must equal plain P(wait > t), for any s and t
	// -- the wait already spent has to vanish from the answer completely.
	cases := []struct{ lambda, s, t float64 }{
		{3, 0.5, 1.0 / 3.0},
		{1, 2, 3},
		{5, 0.1, 0.4},
		{0.5, 4, 1},
	}
	for _, c := range cases {
		got := ConditionalSurvival(c.lambda, c.s, c.t)
		want := Survival(c.lambda, c.t)
		if !near(got, want, 1e-9) {
			t.Errorf("ConditionalSurvival(%v,%v,%v) = %v, want %v (= Survival(lambda,t), memorylessness)",
				c.lambda, c.s, c.t, got, want)
		}
	}
}

func TestSurvivalMatchesWorkedMemorylessExample(t *testing.T) {
	// The two-case arithmetic worked by hand in LEARNINGS.md: both land on
	// ~0.3679 (= e^-1) even though case (a) already spent 0.5 hours waiting.
	caseA := ConditionalSurvival(3, 0.5, 1.0/3.0)
	caseB := Survival(3, 1.0/3.0)
	if !near(caseA, 0.367879, 1e-5) {
		t.Errorf("case (a) = %v, want ~0.367879", caseA)
	}
	if !near(caseB, 0.367879, 1e-5) {
		t.Errorf("case (b) = %v, want ~0.367879", caseB)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("exponential-distribution")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
