package ludecomp

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func almostEqualSlice(t *testing.T, got, want []float64, tol float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (got %v, want %v)", len(got), len(want), got, want)
	}
	for i := range want {
		if !almostEqual(got[i], want[i], tol) {
			t.Errorf("[%d] = %v, want %v (full: got %v, want %v)", i, got[i], want[i], got, want)
		}
	}
}

// matMul multiplies two square matrices, for reconstructing A from L*U in
// tests -- independent of any Decompose internals.
func matMul(a, b Matrix) Matrix {
	n := len(a)
	out := make(Matrix, n)
	for i := 0; i < n; i++ {
		out[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += a[i][k] * b[k][j]
			}
			out[i][j] = sum
		}
	}
	return out
}

// matVec multiplies a square matrix by a vector.
func matVec(a Matrix, x []float64) []float64 {
	n := len(a)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			sum += a[i][j] * x[j]
		}
		out[i] = sum
	}
	return out
}

// TestDecomposeKnownValues is pinned to the L/U worked step by step in
// LEARNINGS.md and the concept's own Sections.
func TestDecomposeKnownValues(t *testing.T) {
	L, U, ok := Decompose(A)
	if !ok {
		t.Fatal("Decompose(A) ok = false, want true")
	}
	wantL := Matrix{
		{1, 0, 0},
		{0.5, 1, 0},
		{0.25, 1.0 / 6, 1},
	}
	wantU := Matrix{
		{4, 3, 2},
		{0, 1.5, 0},
		{0, 0, 1.5},
	}
	for i := range wantL {
		almostEqualSlice(t, L[i], wantL[i], 1e-6)
		almostEqualSlice(t, U[i], wantU[i], 1e-6)
	}
}

// TestDecomposeReconstructsA checks L*U = a for several matrices,
// independent of the pinned known-values example.
func TestDecomposeReconstructsA(t *testing.T) {
	cases := []Matrix{
		A,
		{{2, 1}, {4, 3}},
		{{9, 3, 6}, {3, 5, 1}, {6, 1, 15}},
	}
	for _, a := range cases {
		L, U, ok := Decompose(a)
		if !ok {
			t.Fatalf("Decompose(%v) ok = false, want true", a)
		}
		got := matMul(L, U)
		for i := range a {
			almostEqualSlice(t, got[i], a[i], 1e-9)
		}
	}
}

// TestForwardSubstituteKnownSystem checks forward substitution against a
// hand-solved unit-lower-triangular system independent of Decompose.
func TestForwardSubstituteKnownSystem(t *testing.T) {
	L := Matrix{
		{1, 0, 0},
		{2, 1, 0},
		{3, 4, 1},
	}
	// L y = [1, 5, 20] -> y0=1, y1=5-2*1=3, y2=20-3*1-4*3=5.
	got := ForwardSubstitute(L, []float64{1, 5, 20})
	almostEqualSlice(t, got, []float64{1, 3, 5}, 1e-9)
}

// TestBackSubstituteKnownSystem checks back substitution against a
// hand-solved upper-triangular system independent of Decompose.
func TestBackSubstituteKnownSystem(t *testing.T) {
	U := Matrix{
		{2, 1, 1},
		{0, 3, 1},
		{0, 0, 4},
	}
	// U x = [5, 10, 8] -> x2=2, x1=(10-1*2)/3=8/3, x0=(5-1*(8/3)-1*2)/2.
	x2 := 2.0
	x1 := (10.0 - 1*x2) / 3
	x0 := (5.0 - 1*x1 - 1*x2) / 2
	got := BackSubstitute(U, []float64{5, 10, 8})
	almostEqualSlice(t, got, []float64{x0, x1, x2}, 1e-9)
}

// TestSolveKnownValues is pinned to the b=[9,7,6] example worked in
// LEARNINGS.md and the concept's own Sections.
func TestSolveKnownValues(t *testing.T) {
	x, ok := Solve(A, []float64{9, 7, 6})
	if !ok {
		t.Fatal("Solve(A, b) ok = false, want true")
	}
	almostEqualSlice(t, x, []float64{-1.0 / 9, 5.0 / 3, 20.0 / 9}, 1e-6)
}

// TestSolveUnitVectorsInvertA checks the concept's claim that solving with
// each unit vector eᵢ recovers the corresponding column of A⁻¹: A times
// that column should reproduce eᵢ.
func TestSolveUnitVectorsInvertA(t *testing.T) {
	n := len(A)
	for i := 0; i < n; i++ {
		e := make([]float64, n)
		e[i] = 1
		x, ok := Solve(A, e)
		if !ok {
			t.Fatalf("Solve(A, e%d) ok = false, want true", i)
		}
		got := matVec(A, x)
		almostEqualSlice(t, got, e, 1e-9)
	}
}

// TestDecomposeSingularNoPivotReturnsFalse checks that Decompose reports
// ok=false rather than dividing by ~zero when a pivot is missing and this
// no-pivoting version has no row to swap in.
func TestDecomposeSingularNoPivotReturnsFalse(t *testing.T) {
	singular := Matrix{
		{0, 1},
		{1, 0},
	}
	_, _, ok := Decompose(singular)
	if ok {
		t.Error("Decompose(singular) ok = true, want false")
	}
}

// TestDeterminantKnownValue is pinned to the det(A)=9 byproduct claimed in
// the concept's Sections.
func TestDeterminantKnownValue(t *testing.T) {
	got, ok := Determinant(A)
	if !ok {
		t.Fatal("Determinant(A) ok = false, want true")
	}
	if !almostEqual(got, 9, 1e-6) {
		t.Errorf("Determinant(A) = %v, want 9", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("lu-decomposition")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
