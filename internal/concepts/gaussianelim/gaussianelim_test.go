package gaussianelim

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func approxEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestEliminateMatchesWorkedExample(t *testing.T) {
	steps := Eliminate(system(2))
	// Start, R2<-R2-(-1.50)*R1, R3<-R3-(-1.00)*R1, R3<-R3-(4.00)*R2 = 4 steps.
	if len(steps) != 4 {
		t.Fatalf("Eliminate(e=2) produced %d steps, want 4: %+v", len(steps), steps)
	}
	final := steps[len(steps)-1].Matrix
	want := Matrix{
		{2, 1, -1, 8},
		{0, 0.5, 0.5, 1},
		{0, 0, -1, 1},
	}
	for i := range want {
		for j := range want[i] {
			if !approxEqual(final[i][j], want[i][j]) {
				t.Errorf("final[%d][%d] = %v, want %v (full: %+v)", i, j, final[i][j], want[i][j], final)
			}
		}
	}
}

func TestEliminateSingularCaseProducesContradictionRow(t *testing.T) {
	steps := Eliminate(system(3))
	final := steps[len(steps)-1].Matrix
	last := final[2]
	if !approxEqual(last[0], 0) || !approxEqual(last[1], 0) || !approxEqual(last[2], 0) {
		t.Fatalf("e=3 last row coefficients = %v, want all ~0", last[:3])
	}
	if approxEqual(last[3], 0) {
		t.Fatalf("e=3 last row RHS = %v, want nonzero (0=%.2f is a contradiction)", last[3], last[3])
	}
}

func TestBackSubstituteMatchesWorkedExample(t *testing.T) {
	steps := Eliminate(system(2))
	final := steps[len(steps)-1].Matrix
	sol, ok := BackSubstitute(final)
	if !ok {
		t.Fatal("BackSubstitute(e=2) reported no solution, want (2,3,-1)")
	}
	want := []float64{2, 3, -1}
	for i := range want {
		if !approxEqual(sol[i], want[i]) {
			t.Errorf("solution[%d] = %v, want %v", i, sol[i], want[i])
		}
	}
}

func TestBackSubstituteReportsNoSolutionWhenSingular(t *testing.T) {
	steps := Eliminate(system(3))
	final := steps[len(steps)-1].Matrix
	if _, ok := BackSubstitute(final); ok {
		t.Error("BackSubstitute(e=3) reported a solution, want ok=false (system is inconsistent)")
	}
}

func TestRankDropsAtTheSingularCoefficients(t *testing.T) {
	steps2 := Eliminate(system(2))
	final2 := steps2[len(steps2)-1].Matrix
	if got := Rank(final2, 3); got != 3 {
		t.Errorf("Rank(e=2 coefficients) = %d, want 3", got)
	}
	if got := Rank(final2, 4); got != 3 {
		t.Errorf("Rank(e=2 augmented) = %d, want 3", got)
	}

	steps3 := Eliminate(system(3))
	final3 := steps3[len(steps3)-1].Matrix
	if got := Rank(final3, 3); got != 2 {
		t.Errorf("Rank(e=3 coefficients) = %d, want 2", got)
	}
	if got := Rank(final3, 4); got != 3 {
		t.Errorf("Rank(e=3 augmented) = %d, want 3 (inconsistent: coefficient rank < augmented rank)", got)
	}
}

func TestBackSubstituteRejectsNonSquareInput(t *testing.T) {
	m := Matrix{{1, 2, 3, 4}, {5, 6, 7, 8}} // 2 rows, 4 cols -- not n x (n+1) for n=2
	if _, ok := BackSubstitute(m); ok {
		t.Error("BackSubstitute on a non-square augmented matrix reported ok=true, want false")
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("gaussian-elimination")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
