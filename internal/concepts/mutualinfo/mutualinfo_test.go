package mutualinfo

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

const eps = 1e-9

func approx(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestNewJointSumsToOne(t *testing.T) {
	cases := [][3]float64{{0.3, 0.9, 0.1}, {0.5, 0.5, 0.5}, {0.01, 0.99, 0.01}, {0.7, 0.2, 0.8}}
	for _, c := range cases {
		j := NewJoint(c[0], c[1], c[2])
		sum := j.P11 + j.P10 + j.P01 + j.P00
		if !approx(sum, 1, eps) {
			t.Errorf("NewJoint(%.2f,%.2f,%.2f) probabilities sum to %v, want 1", c[0], c[1], c[2], sum)
		}
	}
}

func TestWorkedExampleRainUmbrella(t *testing.T) {
	// P(rain)=0.3, P(umbrella|rain)=0.9, P(umbrella|no rain)=0.1 -- the
	// LEARNINGS.md worked example.
	j := NewJoint(0.3, 0.9, 0.1)

	want := Joint{P11: 0.27, P10: 0.03, P01: 0.07, P00: 0.63}
	if !approx(j.P11, want.P11, 1e-9) || !approx(j.P10, want.P10, 1e-9) ||
		!approx(j.P01, want.P01, 1e-9) || !approx(j.P00, want.P00, 1e-9) {
		t.Fatalf("NewJoint(0.3,0.9,0.1) = %+v, want %+v", j, want)
	}

	if px := j.MarginalX(); !approx(px, 0.3, 1e-9) {
		t.Errorf("MarginalX() = %v, want 0.3", px)
	}
	if py := j.MarginalY(); !approx(py, 0.34, 1e-9) {
		t.Errorf("MarginalY() = %v, want 0.34", py)
	}

	hx, hy := BinaryEntropy(j.MarginalX()), BinaryEntropy(j.MarginalY())
	if !approx(hx, 0.881, 5e-4) {
		t.Errorf("H(X) = %v, want ~0.881", hx)
	}
	if !approx(hy, 0.925, 5e-4) {
		t.Errorf("H(Y) = %v, want ~0.925", hy)
	}

	hxy := j.JointEntropy()
	if !approx(hxy, 1.350, 5e-4) {
		t.Errorf("H(X,Y) = %v, want ~1.350", hxy)
	}

	mi := j.MutualInformation()
	if !approx(mi, 0.456, 5e-4) {
		t.Errorf("I(X;Y) = %v, want ~0.456", mi)
	}
}

func TestMutualInformationMatchesKLFormula(t *testing.T) {
	cases := [][3]float64{{0.3, 0.9, 0.1}, {0.5, 0.5, 0.5}, {0.2, 0.7, 0.3}, {0.6, 0.1, 0.95}}
	for _, c := range cases {
		j := NewJoint(c[0], c[1], c[2])
		a, b := j.MutualInformation(), j.MutualInformationViaKL()
		if !approx(a, b, 1e-6) {
			t.Errorf("NewJoint(%.2f,%.2f,%.2f): MutualInformation()=%v, MutualInformationViaKL()=%v, want equal",
				c[0], c[1], c[2], a, b)
		}
	}
}

func TestMutualInformationIsZeroWhenIndependent(t *testing.T) {
	// P(Y=1|X=1) == P(Y=1|X=0) means Y doesn't depend on X at all.
	cases := [][3]float64{{0.3, 0.5, 0.5}, {0.5, 0.2, 0.2}, {0.9, 0.7, 0.7}}
	for _, c := range cases {
		j := NewJoint(c[0], c[1], c[2])
		if mi := j.MutualInformation(); !approx(mi, 0, 1e-9) {
			t.Errorf("NewJoint(%.2f,%.2f,%.2f) MutualInformation() = %v, want 0 (X,Y independent)",
				c[0], c[1], c[2], mi)
		}
	}
}

func TestMutualInformationIsMaximalForPerfectDependence(t *testing.T) {
	// P(Y=1|X=1)=1, P(Y=1|X=0)=0: Y is perfectly determined by X, so
	// knowing X removes all uncertainty about Y -- I(X;Y) should equal
	// H(Y) exactly (and H(X), since X and Y then carry identical
	// information).
	j := NewJoint(0.4, 1, 0)
	mi := j.MutualInformation()
	hx, hy := BinaryEntropy(j.MarginalX()), BinaryEntropy(j.MarginalY())
	if !approx(mi, hx, 1e-6) {
		t.Errorf("perfect dependence: I(X;Y)=%v, want H(X)=%v", mi, hx)
	}
	if !approx(mi, hy, 1e-6) {
		t.Errorf("perfect dependence: I(X;Y)=%v, want H(Y)=%v", mi, hy)
	}
}

func TestMutualInformationIsNeverNegative(t *testing.T) {
	for px := 0.05; px < 1; px += 0.1 {
		for py1 := 0.05; py1 < 1; py1 += 0.1 {
			for py0 := 0.05; py0 < 1; py0 += 0.1 {
				j := NewJoint(px, py1, py0)
				if mi := j.MutualInformation(); mi < -1e-9 {
					t.Fatalf("NewJoint(%.2f,%.2f,%.2f) MutualInformation() = %v, want >= 0", px, py1, py0, mi)
				}
			}
		}
	}
}

func TestMutualInformationIsSymmetric(t *testing.T) {
	// I(X;Y) should equal I(Y;X) -- rebuild the joint with X and Y's roles
	// swapped and confirm the same number comes out.
	px, py1, py0 := 0.3, 0.9, 0.1
	j := NewJoint(px, py1, py0)

	// P(X=1|Y=1) and P(X=1|Y=0), derived from the same joint table, define
	// the swapped-roles distribution.
	py := j.MarginalY()
	pxGivenY1 := j.P11 / py
	pxGivenY0 := j.P10 / (1 - py)
	swapped := NewJoint(py, pxGivenY1, pxGivenY0)

	if a, b := j.MutualInformation(), swapped.MutualInformation(); !approx(a, b, 1e-6) {
		t.Errorf("I(X;Y)=%v, I(Y;X)=%v, want equal", a, b)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("mutual-information")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
