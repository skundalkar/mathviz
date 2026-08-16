package fourier

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestSquareWaveKnownValues(t *testing.T) {
	cases := map[float64]float64{
		math.Pi / 2:  1,  // middle of the +1 plateau
		-math.Pi / 2: -1, // middle of the -1 plateau
		0:            0,  // exactly at the jump from -1 to +1
		// wrapPi's half-open [-π,π) range maps both +π and -π to the
		// same representative point -π, on the -1 side of that jump.
		math.Pi:  -1,
		-math.Pi: -1,
	}
	for x, want := range cases {
		if got := SquareWave(x); got != want {
			t.Errorf("SquareWave(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestSquareWaveIsPeriodic(t *testing.T) {
	for _, x := range []float64{0.3, 1.7, -2.1, 2.9} {
		a := SquareWave(x)
		b := SquareWave(x + 2*math.Pi)
		if a != b {
			t.Errorf("SquareWave(%v)=%v, SquareWave(%v+2π)=%v, want equal", x, a, x, b)
		}
	}
}

func TestSawtoothWaveKnownValues(t *testing.T) {
	cases := map[float64]float64{
		math.Pi / 2:  0.5,
		-math.Pi / 2: -0.5,
		0:            0,
	}
	for x, want := range cases {
		if got := SawtoothWave(x); !almostEqual(got, want, 1e-9) {
			t.Errorf("SawtoothWave(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestSawtoothWaveStaysInRange(t *testing.T) {
	for i := 0; i < 200; i++ {
		x := -10.0 + float64(i)*0.1
		v := SawtoothWave(x)
		if v < -1 || v > 1 {
			t.Errorf("SawtoothWave(%v) = %v, out of [-1,1]", x, v)
		}
	}
}

func TestSquareWavePartialSumKnownValues(t *testing.T) {
	// Pinned to the closed-form series at x=π/2. See LEARNINGS.md for the
	// same worked numbers.
	cases := map[int]float64{1: 1.2732, 2: 0.8488, 3: 1.1035, 5: 1.0631, 10: 0.9682, 20: 0.9841}
	for terms, want := range cases {
		if got := SquareWavePartialSum(math.Pi/2, terms); !almostEqual(got, want, 1e-4) {
			t.Errorf("SquareWavePartialSum(π/2,%d) = %v, want %v", terms, got, want)
		}
	}
}

func TestSquareWavePartialSumClampsTerms(t *testing.T) {
	if got, want := SquareWavePartialSum(math.Pi/2, 0), SquareWavePartialSum(math.Pi/2, 1); got != want {
		t.Errorf("SquareWavePartialSum(π/2,0) = %v, want %v (clamped to terms=1)", got, want)
	}
}

func TestSquareWavePartialSumApproachesTargetAwayFromJump(t *testing.T) {
	// Away from the jump, more terms should get strictly closer to the
	// true target value of 1 at x=π/2... not every single step (it
	// oscillates, see LEARNINGS.md), but by a lot more terms it should be
	// much closer than by a few.
	few := math.Abs(SquareWavePartialSum(math.Pi/2, 2) - 1)
	many := math.Abs(SquareWavePartialSum(math.Pi/2, 40) - 1)
	if many >= few {
		t.Errorf("error did not shrink with more terms: 2 terms=%v, 40 terms=%v", few, many)
	}
}

func TestGibbsOvershootPersists(t *testing.T) {
	// The defining Fourier-series behavior: near the jump, the partial
	// sum's peak overshoot does NOT shrink toward 0 as terms grow -- it
	// stays around 18% no matter how many terms are added. Check the
	// peak found on a fine grid near x=0 for two very different term
	// counts and confirm neither is close to the true jump height of 1.
	peakOvershoot := func(terms int) float64 {
		max := 0.0
		for i := 1; i < 2000; i++ {
			x := float64(i) / 2000 * math.Pi
			if v := SquareWavePartialSum(x, terms); v > max {
				max = v
			}
		}
		return max - 1
	}
	o20 := peakOvershoot(20)
	o50 := peakOvershoot(50)
	for _, o := range []float64{o20, o50} {
		if o < 0.15 || o > 0.20 {
			t.Errorf("Gibbs overshoot = %v, want roughly 0.15-0.20 (~18%%)", o)
		}
	}
}

func TestPartialSumDispatch(t *testing.T) {
	x, terms := 1.234, 7
	if got, want := PartialSum(x, terms, 0), SquareWavePartialSum(x, terms); got != want {
		t.Errorf("PartialSum(_,_,0) = %v, want SquareWavePartialSum = %v", got, want)
	}
	if got, want := PartialSum(x, terms, 1), SawtoothPartialSum(x, terms); got != want {
		t.Errorf("PartialSum(_,_,1) = %v, want SawtoothPartialSum = %v", got, want)
	}
}

func TestTargetDispatch(t *testing.T) {
	x := 1.234
	if got, want := Target(x, 0), SquareWave(x); got != want {
		t.Errorf("Target(_,0) = %v, want SquareWave = %v", got, want)
	}
	if got, want := Target(x, 1), SawtoothWave(x); got != want {
		t.Errorf("Target(_,1) = %v, want SawtoothWave = %v", got, want)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("fourier-series")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
