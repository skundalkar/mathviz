package determinant

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestDeterminantKnownValue(t *testing.T) {
	// M = [[2,1],[0,1]], worked by hand in LEARNINGS.md: det = 2*1-1*0 = 2.
	m := Matrix{A: 2, B: 1, C: 0, D: 1}
	if got := Determinant(m); got != 2 {
		t.Errorf("Determinant(%v) = %v, want 2", m, got)
	}
}

func TestDeterminantNegativeMeansOrientationFlip(t *testing.T) {
	swap := Matrix{A: 0, B: 1, C: 1, D: 0} // sends i->(0,1), j->(1,0)
	if got := Determinant(swap); got != -1 {
		t.Errorf("Determinant(swap) = %v, want -1", got)
	}
}

func TestDeterminantZeroWhenRowsProportional(t *testing.T) {
	// Second row is exactly double the first -> collapses to a line.
	m := Matrix{A: 2, B: 1, C: 4, D: 2}
	if got := Determinant(m); got != 0 {
		t.Errorf("Determinant(%v) = %v, want 0", m, got)
	}
}

func TestDeterminantOfIdentityAndRotationsIsOne(t *testing.T) {
	identity := Matrix{A: 1, B: 0, C: 0, D: 1}
	if got := Determinant(identity); got != 1 {
		t.Errorf("Determinant(identity) = %v, want 1", got)
	}
	for _, deg := range []float64{0, 30, 90, 180, 271} {
		r := Rotation(deg)
		if got := Determinant(r); !near(got, 1, 1e-9) {
			t.Errorf("Determinant(Rotation(%v)) = %v, want 1 (rotations preserve area)", deg, got)
		}
	}
}

func TestDeterminantMultiplicativity(t *testing.T) {
	// det(A*B) must equal det(A)*det(B) for any pair of matrices.
	pairs := [][2]Matrix{
		{Rotation(90), Shear(1)},
		{{A: 2, B: 1, C: 0, D: 1}, {A: 0, B: 1, C: 1, D: 0}},
		{{A: 3, B: -1, C: 2, D: 4}, {A: 1, B: 2, C: -1, D: 1}},
	}
	for _, pair := range pairs {
		m1, m2 := pair[0], pair[1]
		prod := Mul(m1, m2)
		gotDet := Determinant(prod)
		wantDet := Determinant(m1) * Determinant(m2)
		if !near(gotDet, wantDet, 1e-6) {
			t.Errorf("Determinant(Mul(%v, %v)) = %v, want %v (= det(m1)*det(m2))", m1, m2, gotDet, wantDet)
		}
	}
}

func TestApplyUnitSquareAreaMatchesAbsDeterminant(t *testing.T) {
	matrices := []Matrix{
		{A: 2, B: 1, C: 0, D: 1},
		{A: 0, B: 1, C: 1, D: 0},
		{A: 3, B: -1, C: 2, D: 4},
		Rotation(37),
		Shear(-0.7),
	}
	for _, m := range matrices {
		image := transform(m, unitSquareCorners)
		area := ShoelaceArea(image)
		want := math.Abs(Determinant(m))
		if !near(area, want, 1e-6) {
			t.Errorf("ShoelaceArea(image of unit square under %v) = %v, want |det| = %v", m, area, want)
		}
	}
}

func TestShoelaceAreaOfUnitSquareIsOne(t *testing.T) {
	if got := ShoelaceArea(unitSquareCorners); got != 1 {
		t.Errorf("ShoelaceArea(unit square) = %v, want 1", got)
	}
}

func TestShoelaceAreaDegenerateIsZero(t *testing.T) {
	// All four points collinear (on y = 2x) -> zero area.
	m := Matrix{A: 2, B: 1, C: 4, D: 2}
	image := transform(m, unitSquareCorners)
	if got := ShoelaceArea(image); !near(got, 0, 1e-9) {
		t.Errorf("ShoelaceArea(collapsed square) = %v, want 0", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("determinant")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
