package randomwalk

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestStepIsAlwaysPlusOrMinusOne(t *testing.T) {
	for i := 1; i <= 500; i++ {
		s := Step(i, 1)
		if s != 1 && s != -1 {
			t.Errorf("Step(%d, 1) = %v, want 1 or -1", i, s)
		}
	}
}

func TestStepIsDeterministic(t *testing.T) {
	for i := 1; i <= 50; i++ {
		if Step(i, 1) != Step(i, 1) {
			t.Fatalf("Step(%d, 1) is not deterministic", i)
		}
	}
}

func TestDifferentRunsDontAlwaysAgree(t *testing.T) {
	// Two different runs should produce at least one differing step over a
	// reasonably long stretch -- otherwise "run" wouldn't actually change
	// anything, defeating the whole point of the run slider.
	differs := false
	for i := 1; i <= 50; i++ {
		if Step(i, 1) != Step(i, 2) {
			differs = true
			break
		}
	}
	if !differs {
		t.Error("Step(i, 1) matched Step(i, 2) for every i in 1..50, want at least one difference")
	}
}

func TestPositionsKnownValuesRun1(t *testing.T) {
	// From LEARNINGS.md's worked example: run 1's traced-out position.
	pos := Positions(300, 1)
	cases := map[int]float64{10: 6, 100: 8, 300: 6}
	for n, want := range cases {
		if got := pos[n-1]; got != want {
			t.Errorf("Positions(300, 1)[%d] = %v, want %v", n-1, got, want)
		}
	}
}

func TestPositionMatchesPositionsLastEntry(t *testing.T) {
	for _, n := range []int{1, 10, 50, 300} {
		pos := Positions(n, 2)
		if got, want := Position(n, 2), pos[n-1]; got != want {
			t.Errorf("Position(%d, 2) = %v, want %v (Positions' last entry)", n, got, want)
		}
	}
}

func TestPositionsIsCumulative(t *testing.T) {
	pos := Positions(20, 3)
	for i := 1; i < len(pos); i++ {
		step := pos[i] - pos[i-1]
		if step != 1 && step != -1 {
			t.Errorf("Positions(20,3)[%d]-Positions(20,3)[%d] = %v, want +-1", i, i-1, step)
		}
	}
}

func TestExpectedSpreadKnownValues(t *testing.T) {
	cases := map[int]float64{0: 0, 1: 1, 4: 2, 100: 10, 300: math.Sqrt(300)}
	for n, want := range cases {
		if got := ExpectedSpread(n); !almostEqual(got, want, 1e-9) {
			t.Errorf("ExpectedSpread(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestExpectedSpreadNegative(t *testing.T) {
	if got := ExpectedSpread(-1); got != 0 {
		t.Errorf("ExpectedSpread(-1) = %v, want 0", got)
	}
}

func TestExpectedSpreadGrowsSlowerThanN(t *testing.T) {
	// The whole point of the concept: spread grows like sqrt(n), so
	// quadrupling n only doubles it, not quadruples it.
	if got, want := ExpectedSpread(400), 2*ExpectedSpread(100); !almostEqual(got, want, 1e-9) {
		t.Errorf("ExpectedSpread(400) = %v, want %v (2x ExpectedSpread(100))", got, want)
	}
}

// TestManyRunsRMSMatchesExpectedSpread is a statistical sanity check (like
// GeneratePoints's correlation test in correlation): the root-mean-square
// position across many independent runs at a fixed n should land close to
// ExpectedSpread(n), the same check LEARNINGS.md's worked example reports.
func TestManyRunsRMSMatchesExpectedSpread(t *testing.T) {
	const n, runs = 100, 500
	sumSq := 0.0
	for r := 1; r <= runs; r++ {
		p := Position(n, r+1000)
		sumSq += p * p
	}
	rms := math.Sqrt(sumSq / runs)
	want := ExpectedSpread(n)
	if !almostEqual(rms, want, want*0.15) {
		t.Errorf("RMS over %d runs at n=%d = %v, want close to %v", runs, n, rms, want)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("random-walk")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
