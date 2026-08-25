package svd

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

const eps = 1e-9

func approx(a, b float64) bool { return math.Abs(a-b) < eps }

func TestDecomposeDiagonalMatrix(t *testing.T) {
	// A diagonal matrix is already its own SVD: V and U are the standard
	// basis, and the singular values are the (sorted) diagonal entries.
	s := Decompose(Matrix{A: 3, B: 0, C: 0, D: 2})
	if !approx(s.Sigma1, 3) || !approx(s.Sigma2, 2) {
		t.Fatalf("Sigma1,Sigma2 = %v,%v, want 3,2", s.Sigma1, s.Sigma2)
	}
	want := []Vec2{{1, 0}, {0, 1}}
	for i, v := range []Vec2{s.V1, s.V2} {
		if !approx(v.X, want[i].X) || !approx(v.Y, want[i].Y) {
			t.Errorf("V%d = %v, want %v", i+1, v, want[i])
		}
	}
	for i, v := range []Vec2{s.U1, s.U2} {
		if !approx(v.X, want[i].X) || !approx(v.Y, want[i].Y) {
			t.Errorf("U%d = %v, want %v", i+1, v, want[i])
		}
	}
}

func TestDecomposeSymmetricMatrixHasUEqualV(t *testing.T) {
	// For a symmetric positive-definite matrix, SVD reduces to
	// eigenvectors-eigenvalues: U == V and the singular values equal the
	// (positive) eigenvalues -- the same [[2,1],[1,2]] example that
	// package eigen uses, with eigenvalues 3 and 1.
	s := Decompose(Matrix{A: 2, B: 1, C: 1, D: 2})
	if !approx(s.Sigma1, 3) || !approx(s.Sigma2, 1) {
		t.Fatalf("Sigma1,Sigma2 = %v,%v, want 3,1", s.Sigma1, s.Sigma2)
	}
	if !approx(s.U1.X, s.V1.X) || !approx(s.U1.Y, s.V1.Y) {
		t.Errorf("U1 = %v, want equal to V1 = %v", s.U1, s.V1)
	}
	if !approx(s.U2.X, s.V2.X) || !approx(s.U2.Y, s.V2.Y) {
		t.Errorf("U2 = %v, want equal to V2 = %v", s.U2, s.V2)
	}
}

func TestDecomposeShearMatrixGivesGoldenRatioSingularValues(t *testing.T) {
	// The shear matrix [[1,1],[0,1]] is defective (only one eigenvalue, one
	// eigenvector direction) but its SVD exists and gives the golden ratio
	// and its reciprocal as singular values.
	s := Decompose(Matrix{A: 1, B: 1, C: 0, D: 1})
	phi := (1 + math.Sqrt(5)) / 2
	if !approx(s.Sigma1, phi) {
		t.Errorf("Sigma1 = %v, want phi = %v", s.Sigma1, phi)
	}
	if !approx(s.Sigma2, 1/phi) {
		t.Errorf("Sigma2 = %v, want 1/phi = %v", s.Sigma2, 1/phi)
	}
}

func TestSigmaOrderingAndNonNegative(t *testing.T) {
	matrices := []Matrix{
		{1, 1, 0, 1},
		{2, 1, 1, 2},
		{1, 2, 2, 4}, // rank-deficient (row 2 = 2 * row 1)
		{0, 0, 0, 0},
		{-1, 2, 3, -4},
		{5, 0, 0, 0},
	}
	for _, m := range matrices {
		s := Decompose(m)
		if s.Sigma1 < 0 || s.Sigma2 < 0 {
			t.Errorf("Decompose(%+v) gave a negative singular value: %v, %v", m, s.Sigma1, s.Sigma2)
		}
		if s.Sigma1 < s.Sigma2 {
			t.Errorf("Decompose(%+v) Sigma1=%v < Sigma2=%v, want Sigma1 >= Sigma2", m, s.Sigma1, s.Sigma2)
		}
	}
}

func TestUAndVAreOrthonormal(t *testing.T) {
	matrices := []Matrix{
		{1, 1, 0, 1},
		{2, 1, 1, 2},
		{1, 2, 2, 4},
		{-1, 2, 3, -4},
		{5, 0, 0, 0}, // rank-deficient, exercises the sigma2~0 fallback
		{0, 0, 0, 0}, // zero matrix, exercises the sigma1~0 fallback
	}
	for _, m := range matrices {
		s := Decompose(m)
		for name, v := range map[string]Vec2{"V1": s.V1, "V2": s.V2, "U1": s.U1, "U2": s.U2} {
			if length := math.Hypot(v.X, v.Y); !approx(length, 1) {
				t.Errorf("Decompose(%+v).%s has length %v, want 1", m, name, length)
			}
		}
		if dot := s.V1.X*s.V2.X + s.V1.Y*s.V2.Y; math.Abs(dot) > 1e-6 {
			t.Errorf("Decompose(%+v) V1.V2 = %v, want ~0", m, dot)
		}
		if dot := s.U1.X*s.U2.X + s.U1.Y*s.U2.Y; math.Abs(dot) > 1e-6 {
			t.Errorf("Decompose(%+v) U1.U2 = %v, want ~0", m, dot)
		}
	}
}

func TestReconstructRecoversOriginalMatrix(t *testing.T) {
	matrices := []Matrix{
		{1, 1, 0, 1},
		{2, 1, 1, 2},
		{1, 2, 2, 4},
		{3, 0, 0, 2},
		{-1, 2, 3, -4},
	}
	for _, m := range matrices {
		got := Reconstruct(Decompose(m))
		if math.Abs(got.A-m.A) > 1e-6 || math.Abs(got.B-m.B) > 1e-6 ||
			math.Abs(got.C-m.C) > 1e-6 || math.Abs(got.D-m.D) > 1e-6 {
			t.Errorf("Reconstruct(Decompose(%+v)) = %+v, want %+v", m, got, m)
		}
	}
}

func TestApplyMatchesSigmaTimesU(t *testing.T) {
	// The defining property of the decomposition: M*V_i == Sigma_i * U_i.
	matrices := []Matrix{{1, 1, 0, 1}, {2, 1, 1, 2}, {-1, 2, 3, -4}}
	for _, m := range matrices {
		s := Decompose(m)
		mv1 := Apply(m, s.V1)
		if !approx(mv1.X, s.Sigma1*s.U1.X) || !approx(mv1.Y, s.Sigma1*s.U1.Y) {
			t.Errorf("Apply(%+v, V1) = %v, want Sigma1*U1 = (%v,%v)",
				m, mv1, s.Sigma1*s.U1.X, s.Sigma1*s.U1.Y)
		}
		mv2 := Apply(m, s.V2)
		if !approx(mv2.X, s.Sigma2*s.U2.X) || !approx(mv2.Y, s.Sigma2*s.U2.Y) {
			t.Errorf("Apply(%+v, V2) = %v, want Sigma2*U2 = (%v,%v)",
				m, mv2, s.Sigma2*s.U2.X, s.Sigma2*s.U2.Y)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("singular-value-decomposition")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
