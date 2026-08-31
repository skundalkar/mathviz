package partialderiv

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestFKnownValues(t *testing.T) {
	if got := F(1, 2); !almostEqual(got, 5, 1e-9) {
		t.Errorf("F(1,2) = %v, want 5", got)
	}
	if got := F(0, 0); !almostEqual(got, 0, 1e-9) {
		t.Errorf("F(0,0) = %v, want 0", got)
	}
	if got := F(-3, 4); !almostEqual(got, 25, 1e-9) {
		t.Errorf("F(-3,4) = %v, want 25", got)
	}
}

func TestPartialsMatchWorkedExample(t *testing.T) {
	// f(x,y) = x²+y² at (1,2): freezing y=2 gives g(x)=x²+4, g'(1)=2.
	// Freezing x=1 gives h(y)=1+y², h'(2)=4. See LEARNINGS.md.
	if got := PartialX(1, 2); !almostEqual(got, 2, 1e-9) {
		t.Errorf("PartialX(1,2) = %v, want 2", got)
	}
	if got := PartialY(1, 2); !almostEqual(got, 4, 1e-9) {
		t.Errorf("PartialY(1,2) = %v, want 4", got)
	}
}

func TestGradientAndMagnitude(t *testing.T) {
	fx, fy := Gradient(1, 2)
	if !almostEqual(fx, 2, 1e-9) || !almostEqual(fy, 4, 1e-9) {
		t.Errorf("Gradient(1,2) = (%v,%v), want (2,4)", fx, fy)
	}
	want := math.Sqrt(20)
	if got := GradientMagnitude(1, 2); !almostEqual(got, want, 1e-9) {
		t.Errorf("GradientMagnitude(1,2) = %v, want %v", got, want)
	}
}

func TestGradientAngleDeg(t *testing.T) {
	want := math.Atan2(4, 2) * 180 / math.Pi // ≈63.4349°
	if got := GradientAngleDeg(1, 2); !almostEqual(got, want, 1e-9) {
		t.Errorf("GradientAngleDeg(1,2) = %v, want %v", got, want)
	}
	if got := GradientAngleDeg(0, 0); got != 0 {
		t.Errorf("GradientAngleDeg(0,0) = %v, want 0 (zero gradient)", got)
	}
}

func TestDirectionalDerivativeMatchesAxisPartials(t *testing.T) {
	// Straight along +x (theta=0) must match PartialX exactly; straight
	// along +y (theta=90) must match PartialY exactly.
	if got, want := DirectionalDerivative(1, 2, 0), PartialX(1, 2); !almostEqual(got, want, 1e-9) {
		t.Errorf("DirectionalDerivative(1,2,0) = %v, want %v (PartialX)", got, want)
	}
	if got, want := DirectionalDerivative(1, 2, 90), PartialY(1, 2); !almostEqual(got, want, 1e-9) {
		t.Errorf("DirectionalDerivative(1,2,90) = %v, want %v (PartialY)", got, want)
	}
}

func TestDirectionalDerivativeMaximizedAlongGradient(t *testing.T) {
	x, y := 1.0, 2.0
	mag := GradientMagnitude(x, y)
	gradAngle := GradientAngleDeg(x, y)

	// At the gradient's own angle, the directional derivative must equal
	// the gradient's magnitude (the steepest possible rate of increase).
	if got := DirectionalDerivative(x, y, gradAngle); !almostEqual(got, mag, 1e-9) {
		t.Errorf("DirectionalDerivative at gradient angle = %v, want %v", got, mag)
	}
	// Directly opposite the gradient, it must equal the negative of that
	// magnitude (the steepest possible rate of decrease).
	if got := DirectionalDerivative(x, y, gradAngle+180); !almostEqual(got, -mag, 1e-9) {
		t.Errorf("DirectionalDerivative opposite gradient angle = %v, want %v", got, -mag)
	}
	// No other sampled direction may exceed the gradient's own magnitude.
	for theta := 0.0; theta < 360; theta++ {
		if got := DirectionalDerivative(x, y, theta); got > mag+1e-9 {
			t.Errorf("DirectionalDerivative(%v,%v,%v) = %v exceeds gradient magnitude %v", x, y, theta, got, mag)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("partial-derivatives-gradient")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
