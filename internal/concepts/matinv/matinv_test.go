package matinv

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestIdentity(t *testing.T) {
	id := Identity(3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			want := 0.0
			if i == j {
				want = 1
			}
			if id[i][j] != want {
				t.Errorf("Identity(3)[%d][%d] = %v, want %v", i, j, id[i][j], want)
			}
		}
	}
}

func TestAugment(t *testing.T) {
	a := Matrix{{4, 7}, {2, 6}}
	aug := Augment(a)
	want := Matrix{{4, 7, 1, 0}, {2, 6, 0, 1}}
	for i := range want {
		for j := range want[i] {
			if aug[i][j] != want[i][j] {
				t.Errorf("Augment(a)[%d][%d] = %v, want %v", i, j, aug[i][j], want[i][j])
			}
		}
	}
}

func TestGaussJordanMatchesWorkedExample(t *testing.T) {
	// A = [[4,7],[2,6]] -> A^-1 = [[0.6,-0.7],[-0.2,0.4]] (see LEARNINGS.md).
	steps := GaussJordan(Augment(System(6)))
	final := steps[len(steps)-1].Matrix
	inv, ok := ExtractInverse(final, 2)
	if !ok {
		t.Fatal("ExtractInverse reported not invertible for a known-invertible matrix")
	}
	want := Matrix{{0.6, -0.7}, {-0.2, 0.4}}
	for i := range want {
		for j := range want[i] {
			if !almostEqual(inv[i][j], want[i][j], 1e-9) {
				t.Errorf("inv[%d][%d] = %v, want %v", i, j, inv[i][j], want[i][j])
			}
		}
	}
}

func TestGaussJordanFailsOnSingularMatrix(t *testing.T) {
	// At k=3.5, det(A) = 4*3.5-14 = 0: A is singular, no inverse exists.
	steps := GaussJordan(Augment(System(3.5)))
	final := steps[len(steps)-1].Matrix
	if _, ok := ExtractInverse(final, 2); ok {
		t.Error("ExtractInverse reported invertible for a singular matrix (k=3.5)")
	}
}

func TestInverseSatisfiesAInvATimesAIsIdentity(t *testing.T) {
	for _, k := range []float64{1, 2, 4, 6, 8, 10} {
		a := System(k)
		steps := GaussJordan(Augment(a))
		inv, ok := ExtractInverse(steps[len(steps)-1].Matrix, 2)
		if !ok {
			t.Fatalf("k=%v: expected invertible, got not invertible", k)
		}
		// A * A^-1 should be the identity.
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				sum := 0.0
				for m := 0; m < 2; m++ {
					sum += a[i][m] * inv[m][j]
				}
				want := 0.0
				if i == j {
					want = 1
				}
				if !almostEqual(sum, want, 1e-9) {
					t.Errorf("k=%v: (A*A^-1)[%d][%d] = %v, want %v", k, i, j, sum, want)
				}
			}
		}
	}
}

func TestDeterminant2x2KnownValues(t *testing.T) {
	if got := Determinant2x2(System(6)); !almostEqual(got, 10, 1e-9) {
		t.Errorf("Determinant2x2(System(6)) = %v, want 10", got)
	}
	if got := Determinant2x2(System(3.5)); !almostEqual(got, 0, 1e-9) {
		t.Errorf("Determinant2x2(System(3.5)) = %v, want 0", got)
	}
}

func TestGaussJordanStartsWithUnmodifiedMatrix(t *testing.T) {
	steps := GaussJordan(Augment(System(6)))
	if steps[0].Description != "Start" {
		t.Errorf("first step description = %q, want %q", steps[0].Description, "Start")
	}
	if steps[0].PivotRow != -1 || steps[0].PivotCol != -1 {
		t.Errorf("first step pivot = (%d,%d), want (-1,-1)", steps[0].PivotRow, steps[0].PivotCol)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("matrix-inverse")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
