package vectors

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestAddKnownValues(t *testing.T) {
	x, y := Add(3, 4, 5, 0)
	if x != 8 || y != 4 {
		t.Errorf("Add(3,4,5,0) = (%v,%v), want (8,4)", x, y)
	}
}

func TestAddIsCommutative(t *testing.T) {
	x1, y1 := Add(3, -2, -1, 5)
	x2, y2 := Add(-1, 5, 3, -2)
	if x1 != x2 || y1 != y2 {
		t.Errorf("Add not commutative: (%v,%v) vs (%v,%v)", x1, y1, x2, y2)
	}
}

func TestDotKnownValues(t *testing.T) {
	if got, want := Dot(3, 4, 5, 0), 15.0; got != want {
		t.Errorf("Dot(3,4,5,0) = %v, want %v", got, want)
	}
}

func TestDotOfPerpendicularVectorsIsZero(t *testing.T) {
	// (3,4) and (4,-3) are perpendicular (rotate one by 90°).
	if got := Dot(3, 4, 4, -3); math.Abs(got) > 1e-9 {
		t.Errorf("Dot(3,4,4,-3) = %v, want 0", got)
	}
}

func TestMagnitudeKnownValues(t *testing.T) {
	if got, want := Magnitude(3, 4), 5.0; math.Abs(got-want) > 1e-9 {
		t.Errorf("Magnitude(3,4) = %v, want %v", got, want)
	}
	if got, want := Magnitude(0, 0), 0.0; got != want {
		t.Errorf("Magnitude(0,0) = %v, want %v", got, want)
	}
}

func TestAngleBetweenKnownValues(t *testing.T) {
	// u=(3,4), v=(5,0): cosθ = 15/(5*5) = 0.6, θ = arccos(0.6).
	got := AngleBetween(3, 4, 5, 0)
	want := math.Acos(0.6)
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("AngleBetween(3,4,5,0) = %v, want %v", got, want)
	}
}

func TestAngleBetweenParallelVectorsIsZero(t *testing.T) {
	if got := AngleBetween(2, 0, 5, 0); math.Abs(got) > 1e-9 {
		t.Errorf("AngleBetween(2,0,5,0) = %v, want 0", got)
	}
}

func TestAngleBetweenOppositeVectorsIsPi(t *testing.T) {
	if got, want := AngleBetween(2, 0, -5, 0), math.Pi; math.Abs(got-want) > 1e-9 {
		t.Errorf("AngleBetween(2,0,-5,0) = %v, want %v", got, want)
	}
}

func TestAngleBetweenPerpendicularVectorsIsHalfPi(t *testing.T) {
	if got, want := AngleBetween(1, 0, 0, 1), math.Pi/2; math.Abs(got-want) > 1e-9 {
		t.Errorf("AngleBetween(1,0,0,1) = %v, want %v", got, want)
	}
}

func TestAngleBetweenZeroVectorIsZero(t *testing.T) {
	if got := AngleBetween(0, 0, 5, 0); got != 0 {
		t.Errorf("AngleBetween(0,0,5,0) = %v, want 0", got)
	}
}

func TestScalarProjectionKnownValue(t *testing.T) {
	// u=(3,4) onto v=(5,0): (u·v)/|v| = 15/5 = 3.
	if got, want := ScalarProjection(3, 4, 5, 0), 3.0; math.Abs(got-want) > 1e-9 {
		t.Errorf("ScalarProjection(3,4,5,0) = %v, want %v", got, want)
	}
}

func TestScalarProjectionNegativeWhenOpposed(t *testing.T) {
	if got := ScalarProjection(3, 4, -5, 0); got >= 0 {
		t.Errorf("ScalarProjection(3,4,-5,0) = %v, want < 0", got)
	}
}

func TestScalarProjectionOntoZeroVectorIsZero(t *testing.T) {
	if got := ScalarProjection(3, 4, 0, 0); got != 0 {
		t.Errorf("ScalarProjection(3,4,0,0) = %v, want 0", got)
	}
}

func TestVectorProjectionKnownValue(t *testing.T) {
	// u=(3,4) onto v=(5,0) should land exactly on (3,0): the x-axis
	// component of u, since v points straight along x.
	x, y := VectorProjection(3, 4, 5, 0)
	if math.Abs(x-3) > 1e-9 || math.Abs(y-0) > 1e-9 {
		t.Errorf("VectorProjection(3,4,5,0) = (%v,%v), want (3,0)", x, y)
	}
}

func TestVectorProjectionMagnitudeMatchesScalarProjection(t *testing.T) {
	px, py := VectorProjection(3, 4, 2, -1)
	got := Magnitude(px, py)
	want := math.Abs(ScalarProjection(3, 4, 2, -1))
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("|VectorProjection| = %v, want %v (|ScalarProjection|)", got, want)
	}
}

func TestVectorProjectionOntoZeroVectorIsZero(t *testing.T) {
	x, y := VectorProjection(3, 4, 0, 0)
	if x != 0 || y != 0 {
		t.Errorf("VectorProjection(3,4,0,0) = (%v,%v), want (0,0)", x, y)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("vectors")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
