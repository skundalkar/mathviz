package eigen

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestEigenvaluesDiagonalMatrix(t *testing.T) {
	// A = [[2,0],[0,1]] is already diagonal -- eigenvalues are just the
	// diagonal entries, larger first.
	l1, l2 := Eigenvalues(2, 0, 1)
	if !almostEqual(l1, 2, 1e-9) || !almostEqual(l2, 1, 1e-9) {
		t.Errorf("Eigenvalues(2,0,1) = (%v,%v), want (2,1)", l1, l2)
	}
}

func TestEigenvaluesDiagonalMatrixSwapped(t *testing.T) {
	// A = [[1,0],[0,2]] -- same matrix, axes swapped: still larger first.
	l1, l2 := Eigenvalues(1, 0, 2)
	if !almostEqual(l1, 2, 1e-9) || !almostEqual(l2, 1, 1e-9) {
		t.Errorf("Eigenvalues(1,0,2) = (%v,%v), want (2,1)", l1, l2)
	}
}

func TestEigenvaluesWorkedShearExample(t *testing.T) {
	// A = [[2,1],[1,2]]: trace 4, disc = sqrt(0+1) = 1, eigenvalues 3 and 1.
	l1, l2 := Eigenvalues(2, 1, 2)
	if !almostEqual(l1, 3, 1e-9) || !almostEqual(l2, 1, 1e-9) {
		t.Errorf("Eigenvalues(2,1,2) = (%v,%v), want (3,1)", l1, l2)
	}
}

func TestEigenvaluesRepeatedRoot(t *testing.T) {
	// A = 3*I: both eigenvalues are 3 -- degenerate, but shouldn't panic
	// or produce NaN.
	l1, l2 := Eigenvalues(3, 0, 3)
	if !almostEqual(l1, 3, 1e-9) || !almostEqual(l2, 3, 1e-9) {
		t.Errorf("Eigenvalues(3,0,3) = (%v,%v), want (3,3)", l1, l2)
	}
}

func TestEigenvectorAngleAxisAligned(t *testing.T) {
	// A = [[2,0],[0,1]]: no shear, so the eigenvector for the larger
	// eigenvalue (2) is already the x-axis, angle 0.
	theta := EigenvectorAngle(2, 0, 1)
	if !almostEqual(theta, 0, 1e-9) {
		t.Errorf("EigenvectorAngle(2,0,1) = %v rad, want 0", theta)
	}
}

func TestEigenvectorAngleAxisSwapped(t *testing.T) {
	// A = [[1,0],[0,2]]: the larger eigenvalue (2) belongs to the y-axis,
	// angle 90 degrees.
	theta := EigenvectorAngle(1, 0, 2) * 180 / math.Pi
	if !almostEqual(theta, 90, 1e-6) {
		t.Errorf("EigenvectorAngle(1,0,2) = %v deg, want 90", theta)
	}
}

func TestEigenvectorAngleWorkedShearExample(t *testing.T) {
	// A = [[2,1],[1,2]]: the eigenvector for lambda=3 is (1,1)/sqrt2, a
	// 45 degree direction -- verified directly below via Apply too.
	theta := EigenvectorAngle(2, 1, 2) * 180 / math.Pi
	if !almostEqual(theta, 45, 1e-6) {
		t.Errorf("EigenvectorAngle(2,1,2) = %v deg, want 45", theta)
	}
}

func TestApplyMatchesEigenvalueOnEigenvector(t *testing.T) {
	// For every (a,b,d) tried, applying A to its own eigenvector direction
	// must return exactly lambda*v -- the defining property of an
	// eigenvector, not just something that happens to hold for one example.
	cases := []struct{ a, b, d float64 }{
		{2, 0, 1},
		{1, 0, 2},
		{2, 1, 2},
		{2.4, 0.8, 1.1},
		{0.5, -1.2, 2.5},
	}
	for _, c := range cases {
		l1, _ := Eigenvalues(c.a, c.b, c.d)
		theta := EigenvectorAngle(c.a, c.b, c.d)
		vx, vy := math.Cos(theta), math.Sin(theta)
		ax, ay := Apply(c.a, c.b, c.d, vx, vy)
		wantX, wantY := l1*vx, l1*vy
		if !almostEqual(ax, wantX, 1e-9) || !almostEqual(ay, wantY, 1e-9) {
			t.Errorf("Apply(%v,%v,%v, eigenvector) = (%v,%v), want lambda1*v = (%v,%v)",
				c.a, c.b, c.d, ax, ay, wantX, wantY)
		}
	}
}

func TestApplyOffEigenvectorIsNotParallel(t *testing.T) {
	// A generic direction (not an eigenvector) should get rotated: Av is
	// not a scalar multiple of v, so the angle between them isn't 0 or 180.
	a, b, d := 2.0, 1.0, 2.0
	vx, vy := math.Cos(0), math.Sin(0) // angle 0, not the 45 deg eigenvector
	ax, ay := Apply(a, b, d, vx, vy)
	angle := AngleBetweenDeg(vx, vy, ax, ay)
	if almostEqual(angle, 0, 1) || almostEqual(angle, 180, 1) {
		t.Errorf("AngleBetweenDeg(v, Av) = %v deg for a non-eigenvector direction, want far from 0/180", angle)
	}
}

func TestAngleBetweenDegParallelAndOpposite(t *testing.T) {
	if got := AngleBetweenDeg(1, 0, 2, 0); !almostEqual(got, 0, 1e-9) {
		t.Errorf("AngleBetweenDeg(1,0,2,0) = %v, want 0", got)
	}
	if got := AngleBetweenDeg(1, 0, -1, 0); !almostEqual(got, 180, 1e-9) {
		t.Errorf("AngleBetweenDeg(1,0,-1,0) = %v, want 180", got)
	}
	if got := AngleBetweenDeg(1, 0, 0, 1); !almostEqual(got, 90, 1e-9) {
		t.Errorf("AngleBetweenDeg(1,0,0,1) = %v, want 90", got)
	}
}

func TestAngleBetweenDegZeroVector(t *testing.T) {
	if got := AngleBetweenDeg(0, 0, 1, 1); got != 0 {
		t.Errorf("AngleBetweenDeg(0,0,1,1) = %v, want 0", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("eigenvectors-eigenvalues")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
