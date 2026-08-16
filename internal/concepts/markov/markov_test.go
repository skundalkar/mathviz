package markov

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestSteadyStateKnownValues(t *testing.T) {
	// Pinned to the worked a=0.9, b=0.5 example. See LEARNINGS.md for the
	// same numbers.
	sunny, rainy := SteadyState(0.9, 0.5)
	if !almostEqual(sunny, 0.8333, 1e-4) {
		t.Errorf("SteadyState(0.9,0.5) sunny = %v, want ~0.8333", sunny)
	}
	if !almostEqual(rainy, 0.1667, 1e-4) {
		t.Errorf("SteadyState(0.9,0.5) rainy = %v, want ~0.1667", rainy)
	}
}

func TestSteadyStateSumsToOne(t *testing.T) {
	for _, a := range []float64{0.1, 0.5, 0.9} {
		for _, b := range []float64{0.1, 0.5, 0.9} {
			sunny, rainy := SteadyState(a, b)
			if !almostEqual(sunny+rainy, 1, 1e-9) {
				t.Errorf("SteadyState(%v,%v) sums to %v, want 1", a, b, sunny+rainy)
			}
		}
	}
}

func TestSteadyStateIsFixedPoint(t *testing.T) {
	// The defining property: applying Step to the steady state should
	// return the steady state unchanged.
	for _, a := range []float64{0.2, 0.6, 0.9} {
		for _, b := range []float64{0.2, 0.6, 0.9} {
			sunny, _ := SteadyState(a, b)
			next := Step(sunny, a, b)
			if !almostEqual(next, sunny, 1e-9) {
				t.Errorf("Step(SteadyState(%v,%v), %v, %v) = %v, want fixed point %v", a, b, a, b, next, sunny)
			}
		}
	}
}

func TestStepKnownValues(t *testing.T) {
	cases := []struct {
		pSunny, a, b, want float64
	}{
		{0, 0.9, 0.5, 0.5},
		{0.5, 0.9, 0.5, 0.7},
		{0.7, 0.9, 0.5, 0.78},
	}
	for _, c := range cases {
		if got := Step(c.pSunny, c.a, c.b); !almostEqual(got, c.want, 1e-9) {
			t.Errorf("Step(%v,%v,%v) = %v, want %v", c.pSunny, c.a, c.b, got, c.want)
		}
	}
}

func TestTrajectoryKnownValues(t *testing.T) {
	traj := Trajectory(0, 0.9, 0.5, 5)
	want := []float64{0, 0.5, 0.7, 0.78, 0.812, 0.8248}
	if len(traj) != len(want) {
		t.Fatalf("len(Trajectory) = %d, want %d", len(traj), len(want))
	}
	for i, w := range want {
		if !almostEqual(traj[i], w, 1e-4) {
			t.Errorf("Trajectory[%d] = %v, want %v", i, traj[i], w)
		}
	}
}

func TestTrajectoryLength(t *testing.T) {
	if got := len(Trajectory(0, 0.9, 0.5, 30)); got != 31 {
		t.Errorf("len(Trajectory(...,30)) = %d, want 31", got)
	}
	if got := len(Trajectory(0, 0.9, 0.5, -5)); got != 1 {
		t.Errorf("len(Trajectory(...,-5)) = %d, want 1 (clamped to 0 steps)", got)
	}
}

func TestTrajectoryConvergesToSteadyStateFromEitherStart(t *testing.T) {
	// The defining Markov-chain behavior: no matter which state you
	// start in, the chain converges to the same steady state.
	a, b := 0.9, 0.5
	sunny, _ := SteadyState(a, b)
	fromRainy := Trajectory(0, a, b, 50)
	fromSunny := Trajectory(1, a, b, 50)
	if !almostEqual(fromRainy[50], sunny, 1e-6) {
		t.Errorf("starting from rainy, day 50 P(sunny) = %v, want steady state %v", fromRainy[50], sunny)
	}
	if !almostEqual(fromSunny[50], sunny, 1e-6) {
		t.Errorf("starting from sunny, day 50 P(sunny) = %v, want steady state %v", fromSunny[50], sunny)
	}
}

func TestTrajectoryGapShrinksGeometrically(t *testing.T) {
	// Unlike law-of-large-numbers' noisy convergence, this gap shrinks by
	// the exact same ratio lambda = a+b-1 every single step.
	a, b := 0.9, 0.5
	lambda := a + b - 1
	sunny, _ := SteadyState(a, b)
	traj := Trajectory(0, a, b, 10)
	for day := 1; day < len(traj); day++ {
		gapPrev := traj[day-1] - sunny
		gapNow := traj[day] - sunny
		predicted := gapPrev * lambda
		if !almostEqual(gapNow, predicted, 1e-9) {
			t.Errorf("day=%d: gap=%v, want gap[day-1]*lambda=%v", day, gapNow, predicted)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("markov-chains")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
