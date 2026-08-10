package expgrowth

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestValueAtZeroIsOne(t *testing.T) {
	for _, rate := range []float64{1, 10, 50} {
		if v := Value(0, rate); math.Abs(v-1) > 1e-9 {
			t.Errorf("Value(0, %v) = %v, want 1", rate, v)
		}
	}
}

func TestValueMonotonicIncreasing(t *testing.T) {
	rate := 10.0
	last := Value(0, rate)
	for tp := 0.5; tp <= 40; tp += 0.5 {
		v := Value(tp, rate)
		if v < last {
			t.Fatalf("Value not increasing at t=%.1f: %v < %v", tp, v, last)
		}
		last = v
	}
}

func TestValueDoublesAtDoublingTime(t *testing.T) {
	for _, rate := range []float64{1, 5, 10, 25, 50} {
		dt := DoublingTime(rate)
		v := Value(dt, rate)
		if math.Abs(v-2) > 1e-6 {
			t.Errorf("Value(DoublingTime(%v), %v) = %v, want 2", rate, rate, v)
		}
	}
}

func TestDoublingTimeShrinksAsRateGrows(t *testing.T) {
	last := DoublingTime(1)
	for _, rate := range []float64{2, 5, 10, 20, 50} {
		dt := DoublingTime(rate)
		if dt >= last {
			t.Fatalf("doubling time did not shrink as rate rose to %v: %v >= %v", rate, dt, last)
		}
		last = dt
	}
}

func TestDoublingTimeNonPositiveRateIsInfinite(t *testing.T) {
	if dt := DoublingTime(0); !math.IsInf(dt, 1) {
		t.Errorf("DoublingTime(0) = %v, want +Inf", dt)
	}
	if dt := DoublingTime(-5); !math.IsInf(dt, 1) {
		t.Errorf("DoublingTime(-5) = %v, want +Inf", dt)
	}
}

func TestLinearMatchesValueNearZero(t *testing.T) {
	// Both start at 1 and share the same initial slope, so they should
	// nearly agree for a small step.
	rate := 10.0
	v, l := Value(0.01, rate), Linear(0.01, rate)
	if math.Abs(v-l) > 1e-3 {
		t.Errorf("Value and Linear diverge too fast near t=0: %v vs %v", v, l)
	}
}

func TestExponentialOutpacesLinearEventually(t *testing.T) {
	rate := 10.0
	if Value(30, rate) <= Linear(30, rate) {
		t.Fatalf("compounding did not pull ahead of the linear projection by t=30")
	}
}

func TestRule70ApproximatesDoublingTime(t *testing.T) {
	// The rule of 70 is a rough mental-math shortcut; it should be within a
	// reasonable margin of the exact doubling time for modest rates.
	for _, rate := range []float64{2, 5, 10} {
		exact := DoublingTime(rate)
		approx := Rule70(rate)
		if math.Abs(exact-approx)/exact > 0.10 {
			t.Errorf("Rule70(%v) = %.2f too far from exact DoublingTime = %.2f", rate, approx, exact)
		}
	}
}

func TestRule70NonPositiveRateIsInfinite(t *testing.T) {
	if r := Rule70(0); !math.IsInf(r, 1) {
		t.Errorf("Rule70(0) = %v, want +Inf", r)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("exponential-growth")
	if !ok {
		t.Fatal("exponential-growth not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
