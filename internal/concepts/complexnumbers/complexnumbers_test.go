package complexnumbers

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestMultiplyByIIsQuarterTurn(t *testing.T) {
	re, im := Multiply(3, 4, 0, 1)
	if re != -4 || im != 3 {
		t.Errorf("Multiply(3,4,0,1) = (%v,%v), want (-4,3)", re, im)
	}
}

func TestMultiplyKnownValue(t *testing.T) {
	re, im := Multiply(3, 4, 4, 3)
	if re != 0 || im != 25 {
		t.Errorf("Multiply(3,4,4,3) = (%v,%v), want (0,25)", re, im)
	}
}

func TestMultiplyByOneIsIdentity(t *testing.T) {
	re, im := Multiply(3, 4, 1, 0)
	if re != 3 || im != 4 {
		t.Errorf("Multiply(3,4,1,0) = (%v,%v), want (3,4)", re, im)
	}
}

func TestMultiplyMultipliesModuli(t *testing.T) {
	a, b, c, d := 3.0, 4.0, -2.0, 5.0
	re, im := Multiply(a, b, c, d)
	got := Modulus(re, im)
	want := Modulus(a, b) * Modulus(c, d)
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("|Multiply(3,4,-2,5)| = %v, want %v", got, want)
	}
}

func TestModulusKnownValues(t *testing.T) {
	if got, want := Modulus(3, 4), 5.0; math.Abs(got-want) > 1e-9 {
		t.Errorf("Modulus(3,4) = %v, want %v", got, want)
	}
	if got, want := Modulus(0, 0), 0.0; got != want {
		t.Errorf("Modulus(0,0) = %v, want %v", got, want)
	}
}

func TestArgumentKnownValues(t *testing.T) {
	cases := []struct {
		a, b, want float64
	}{
		{1, 0, 0}, {0, 1, math.Pi / 2}, {-1, 0, math.Pi}, {0, -1, -math.Pi / 2},
	}
	for _, tc := range cases {
		if got := Argument(tc.a, tc.b); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("Argument(%v,%v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestFromPolarKnownValues(t *testing.T) {
	re, im := FromPolar(5, 0)
	if math.Abs(re-5) > 1e-9 || math.Abs(im-0) > 1e-9 {
		t.Errorf("FromPolar(5,0) = (%v,%v), want (5,0)", re, im)
	}
	re, im = FromPolar(1, math.Pi/2)
	if math.Abs(re-0) > 1e-9 || math.Abs(im-1) > 1e-9 {
		t.Errorf("FromPolar(1,π/2) = (%v,%v), want (0,1)", re, im)
	}
}

func TestFromPolarRoundTripsThroughModulusArgument(t *testing.T) {
	a, b := 3.0, -4.0
	r, theta := Modulus(a, b), Argument(a, b)
	re, im := FromPolar(r, theta)
	if math.Abs(re-a) > 1e-9 || math.Abs(im-b) > 1e-9 {
		t.Errorf("FromPolar(Modulus,Argument) = (%v,%v), want (%v,%v)", re, im, a, b)
	}
}

func TestRotateByHalfPiMatchesMultiplyByI(t *testing.T) {
	re, im := Rotate(3, 4, math.Pi/2, 1)
	if math.Abs(re-(-4)) > 1e-9 || math.Abs(im-3) > 1e-9 {
		t.Errorf("Rotate(3,4,π/2,1) = (%v,%v), want (-4,3)", re, im)
	}
}

func TestRotateByZeroAngleOnlyScales(t *testing.T) {
	re, im := Rotate(3, 4, 0, 2)
	if math.Abs(re-6) > 1e-9 || math.Abs(im-8) > 1e-9 {
		t.Errorf("Rotate(3,4,0,2) = (%v,%v), want (6,8)", re, im)
	}
}

func TestRotatePreservesLengthWhenScaleIsOne(t *testing.T) {
	a, b := 3.0, 4.0
	want := Modulus(a, b)
	for _, theta := range []float64{0, 0.7, math.Pi, 4.2, 2 * math.Pi} {
		re, im := Rotate(a, b, theta, 1)
		if got := Modulus(re, im); math.Abs(got-want) > 1e-9 {
			t.Errorf("Rotate(3,4,%v,1) length = %v, want %v", theta, got, want)
		}
	}
}

func TestRotateAddsAngles(t *testing.T) {
	a, b, theta := 3.0, 4.0, 0.9
	re, im := Rotate(a, b, theta, 1)
	got := Argument(re, im)
	want := Argument(a, b) + theta
	// Normalize both into the same [-π, π] wraparound before comparing.
	diff := math.Mod(got-want+math.Pi, 2*math.Pi) - math.Pi
	if math.Abs(diff) > 1e-9 {
		t.Errorf("Argument(Rotate(3,4,%v,1)) = %v, want %v", theta, got, want)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("complex-numbers")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
