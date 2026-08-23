package quadratic

import (
	"math"
	"math/cmplx"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestDiscriminantKnownValues(t *testing.T) {
	// x^2-5x+6=0, worked by hand in LEARNINGS.md: D = 25-24 = 1.
	if got := Discriminant(1, -5, 6); got != 1 {
		t.Errorf("Discriminant(1,-5,6) = %v, want 1", got)
	}
	// x^2+2x+1=0: D = 4-4 = 0.
	if got := Discriminant(1, 2, 1); got != 0 {
		t.Errorf("Discriminant(1,2,1) = %v, want 0", got)
	}
	// x^2+x+1=0: D = 1-4 = -3.
	if got := Discriminant(1, 1, 1); got != -3 {
		t.Errorf("Discriminant(1,1,1) = %v, want -3", got)
	}
}

func TestRootsTwoRealSolveTheEquation(t *testing.T) {
	// x^2-5x+6=0 -> roots 2 and 3, worked by hand in LEARNINGS.md.
	r1, r2 := Roots(1, -5, 6)
	if imag(r1) != 0 || imag(r2) != 0 {
		t.Fatalf("Roots(1,-5,6) = %v, %v, want purely real (D>0)", r1, r2)
	}
	got := []float64{real(r1), real(r2)}
	want := []float64{3, 2} // (-b+sqrtD)/2a first, then (-b-sqrtD)/2a
	for i := range got {
		if !near(got[i], want[i], 1e-9) {
			t.Errorf("root %d = %v, want %v", i, got[i], want[i])
		}
	}
	// Every real root must zero the expression.
	for _, r := range got {
		if v := Evaluate(1, -5, 6, r); !near(v, 0, 1e-9) {
			t.Errorf("Evaluate(1,-5,6, %v) = %v, want 0", r, v)
		}
	}
}

func TestRootsRepeatedWhenDiscriminantZero(t *testing.T) {
	// x^2+2x+1=0 -> repeated root -1.
	r1, r2 := Roots(1, 2, 1)
	if !near(real(r1), -1, 1e-9) || !near(real(r2), -1, 1e-9) {
		t.Errorf("Roots(1,2,1) = %v, %v, want both -1", r1, r2)
	}
	if imag(r1) != 0 || imag(r2) != 0 {
		t.Errorf("Roots(1,2,1) = %v, %v, want purely real (D=0)", r1, r2)
	}
}

func TestRootsComplexWhenDiscriminantNegative(t *testing.T) {
	// x^2+x+1=0 -> (-1 +/- i*sqrt3)/2, worked by hand in LEARNINGS.md.
	r1, r2 := Roots(1, 1, 1)
	wantRe, wantIm := -0.5, math.Sqrt(3)/2
	if !near(real(r1), wantRe, 1e-9) || !near(imag(r1), wantIm, 1e-9) {
		t.Errorf("Roots(1,1,1) r1 = %v, want %v+%vi", r1, wantRe, wantIm)
	}
	if !near(real(r2), wantRe, 1e-9) || !near(imag(r2), -wantIm, 1e-9) {
		t.Errorf("Roots(1,1,1) r2 = %v, want %v-%vi", r2, wantRe, wantIm)
	}
	// The complex pair must still zero ax^2+bx+c in complex arithmetic.
	for _, r := range []complex128{r1, r2} {
		v := complex(1, 0)*r*r + complex(1, 0)*r + complex(1, 0)
		if cmplx.Abs(v) > 1e-9 {
			t.Errorf("complex root %v does not zero x^2+x+1: got %v", r, v)
		}
	}
}

func TestVertexMatchesWorkedExample(t *testing.T) {
	// x^2-5x+6=0: h=-b/2a=2.5, k=Evaluate(1,-5,6,2.5)=6.25-12.5+6=-0.25.
	h, k := Vertex(1, -5, 6)
	if !near(h, 2.5, 1e-9) {
		t.Errorf("Vertex h = %v, want 2.5", h)
	}
	if !near(k, -0.25, 1e-9) {
		t.Errorf("Vertex k = %v, want -0.25", k)
	}
}

func TestVertexSitsBetweenTheTwoRealRoots(t *testing.T) {
	// For any a,b,c with D>0, the axis of symmetry h must sit exactly
	// halfway between the two real roots.
	cases := [][3]float64{{1, -5, 6}, {2, 3, -2}, {-1, 4, 5}}
	for _, cse := range cases {
		a, b, c := cse[0], cse[1], cse[2]
		h, _ := Vertex(a, b, c)
		r1, r2 := Roots(a, b, c)
		mid := (real(r1) + real(r2)) / 2
		if !near(h, mid, 1e-9) {
			t.Errorf("Vertex(%v,%v,%v) h=%v, want midpoint of roots %v", a, b, c, h, mid)
		}
	}
}

func TestEvaluateKnownValue(t *testing.T) {
	if got := Evaluate(1, -5, 6, 0); got != 6 {
		t.Errorf("Evaluate(1,-5,6,0) = %v, want 6", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("quadratic-formula")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
