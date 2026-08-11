package sinecosine

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestCirclePointOnUnitCircle(t *testing.T) {
	for _, theta := range []float64{0, math.Pi / 6, math.Pi / 2, math.Pi, 3 * math.Pi / 2, 4} {
		x, y := CirclePoint(theta, 1)
		if got := x*x + y*y; math.Abs(got-1) > 1e-9 {
			t.Errorf("CirclePoint(%v, 1): x²+y² = %v, want 1", theta, got)
		}
	}
}

func TestCirclePointScalesWithRadius(t *testing.T) {
	for _, radius := range []float64{0.5, 1, 2, 3} {
		x, y := CirclePoint(1.2, radius)
		if got, want := x*x+y*y, radius*radius; math.Abs(got-want) > 1e-9 {
			t.Errorf("CirclePoint(1.2, %v): x²+y² = %v, want %v", radius, got, want)
		}
	}
}

func TestCirclePointKnownValues(t *testing.T) {
	x, y := CirclePoint(0, 1)
	if math.Abs(x-1) > 1e-9 || math.Abs(y-0) > 1e-9 {
		t.Errorf("CirclePoint(0, 1) = (%v, %v), want (1, 0)", x, y)
	}
	x, y = CirclePoint(math.Pi/2, 1)
	if math.Abs(x-0) > 1e-9 || math.Abs(y-1) > 1e-9 {
		t.Errorf("CirclePoint(π/2, 1) = (%v, %v), want (0, 1)", x, y)
	}
}

func TestSineWaveKnownValues(t *testing.T) {
	cases := []struct {
		theta, want float64
	}{
		{0, 0}, {math.Pi / 2, 1}, {math.Pi, 0}, {3 * math.Pi / 2, -1},
	}
	for _, tc := range cases {
		if got := SineWave(tc.theta, 1); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("SineWave(%v, 1) = %v, want %v", tc.theta, got, tc.want)
		}
	}
}

func TestCosineWaveKnownValues(t *testing.T) {
	cases := []struct {
		theta, want float64
	}{
		{0, 1}, {math.Pi / 2, 0}, {math.Pi, -1}, {3 * math.Pi / 2, 0},
	}
	for _, tc := range cases {
		if got := CosineWave(tc.theta, 1); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("CosineWave(%v, 1) = %v, want %v", tc.theta, got, tc.want)
		}
	}
}

func TestSineCosineIdentity(t *testing.T) {
	// sin²θ + cos²θ = 1, the Pythagorean identity, follows directly from
	// the point sitting on a circle of radius 1.
	for _, theta := range []float64{0.3, 1.1, 2.7, 4.2, 5.9} {
		s, c := SineWave(theta, 1), CosineWave(theta, 1)
		if got := s*s + c*c; math.Abs(got-1) > 1e-9 {
			t.Errorf("sin²(%v)+cos²(%v) = %v, want 1", theta, theta, got)
		}
	}
}

func TestSineWaveScalesWithRadius(t *testing.T) {
	if got, want := SineWave(math.Pi/2, 2), 2.0; math.Abs(got-want) > 1e-9 {
		t.Errorf("SineWave(π/2, 2) = %v, want %v", got, want)
	}
}

func TestRevolutionsKnownValues(t *testing.T) {
	cases := []struct{ theta, want float64 }{
		{0, 0}, {math.Pi, 0.5}, {2 * math.Pi, 1}, {4 * math.Pi, 2},
	}
	for _, tc := range cases {
		if got := Revolutions(tc.theta); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("Revolutions(%v) = %v, want %v", tc.theta, got, tc.want)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("sine-cosine")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
