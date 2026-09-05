package lawtotalprob

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestSharesSumToOne(t *testing.T) {
	cases := [][2]float64{{0.90, 0.08}, {0.5, 0.4}, {0.96, 0.4}, {0.6, 0.6}}
	for _, c := range cases {
		a, b, cc := Shares(c[0], c[1])
		sum := a + b + cc
		if !near(sum, 1, 1e-9) {
			t.Errorf("Shares(%v, %v) = (%v, %v, %v), sum %v, want 1", c[0], c[1], a, b, cc, sum)
		}
		if a < 0 || b < 0 || cc < 0 {
			t.Errorf("Shares(%v, %v) = (%v, %v, %v), want all non-negative", c[0], c[1], a, b, cc)
		}
	}
}

func TestSharesDefaultCase(t *testing.T) {
	// The default scenario: 90% / 8% / 2%, no clamping needed.
	a, b, c := Shares(0.90, 0.08)
	if !near(a, 0.90, 1e-9) || !near(b, 0.08, 1e-9) || !near(c, 0.02, 1e-9) {
		t.Errorf("Shares(0.90, 0.08) = (%v, %v, %v), want (0.90, 0.08, 0.02)", a, b, c)
	}
}

func TestTotalProbabilityFactoryExample(t *testing.T) {
	// Worked by hand in LEARNINGS.md: 0.9*0.01 + 0.08*0.20 + 0.02*0.01 = 0.0252.
	shares := []float64{0.90, 0.08, 0.02}
	rates := []float64{0.01, 0.20, 0.01}
	if got := TotalProbability(shares, rates); !near(got, 0.0252, 1e-9) {
		t.Errorf("TotalProbability(factory example) = %v, want 0.0252", got)
	}
}

func TestNaiveAverageFactoryExample(t *testing.T) {
	// (0.01 + 0.20 + 0.01) / 3 = 0.0733...
	rates := []float64{0.01, 0.20, 0.01}
	if got := NaiveAverage(rates); !near(got, 0.073333, 1e-5) {
		t.Errorf("NaiveAverage(factory example) = %v, want ~0.073333", got)
	}
}

func TestNaiveAverageDivergesFromWeighted(t *testing.T) {
	// The whole point of the concept: when shares are uneven, the naive
	// unweighted average is a poor stand-in for the correctly weighted
	// total probability.
	shares := []float64{0.90, 0.08, 0.02}
	rates := []float64{0.01, 0.20, 0.01}
	weighted := TotalProbability(shares, rates)
	naive := NaiveAverage(rates)
	if math.Abs(naive-weighted) < 0.02 {
		t.Errorf("naive (%v) and weighted (%v) should differ substantially for uneven shares", naive, weighted)
	}
}

func TestTotalProbabilityEvenSharesMatchesNaive(t *testing.T) {
	// When every scenario is equally likely, weighting by share IS the
	// same as an unweighted average -- the naive shortcut only fails when
	// shares are uneven.
	shares := []float64{1.0 / 3, 1.0 / 3, 1.0 / 3}
	rates := []float64{0.01, 0.20, 0.01}
	if got, want := TotalProbability(shares, rates), NaiveAverage(rates); !near(got, want, 1e-9) {
		t.Errorf("TotalProbability with even shares = %v, want %v (NaiveAverage)", got, want)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("law-of-total-probability")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
