package lagrange

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func approxEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestMaximumPointMatchesWorkedExample(t *testing.T) {
	x, y := math.Sqrt(2), math.Sqrt(2) // r=2's max candidate, y=x
	if got := F(x, y); !approxEqual(got, 2) {
		t.Errorf("F(√2,√2) = %v, want 2", got)
	}
	if got := Cross(x, y); math.Abs(got) > 1e-9 {
		t.Errorf("Cross(√2,√2) = %v, want ~0 (gradients parallel at the max)", got)
	}
	if got := Lambda(x, y); !approxEqual(got, 0.5) {
		t.Errorf("Lambda(√2,√2) = %v, want 0.5", got)
	}
}

func TestMinimumPointMatchesWorkedExample(t *testing.T) {
	x, y := math.Sqrt(2), -math.Sqrt(2) // r=2's min candidate, y=-x
	if got := F(x, y); !approxEqual(got, -2) {
		t.Errorf("F(√2,-√2) = %v, want -2", got)
	}
	if got := Cross(x, y); math.Abs(got) > 1e-9 {
		t.Errorf("Cross(√2,-√2) = %v, want ~0 (gradients parallel at the min)", got)
	}
	if got := Lambda(x, y); !approxEqual(got, -0.5) {
		t.Errorf("Lambda(√2,-√2) = %v, want -0.5", got)
	}
}

func TestCrossIsNonzeroAwayFromTheDiagonals(t *testing.T) {
	// theta=0 -> (r,0): a point on the circle that is not a candidate.
	x, y := PointOnCircle(2, 0)
	if got := Cross(x, y); math.Abs(got) < 1e-6 {
		t.Errorf("Cross(r,0) = %v, want clearly nonzero (not a Lagrange point)", got)
	}
}

func TestCrossFormulaMatchesTwiceYSquaredMinusXSquared(t *testing.T) {
	for _, pt := range [][2]float64{{1, 1}, {2, -1}, {0.5, 3}, {-2, -2}} {
		x, y := pt[0], pt[1]
		want := 2 * (y*y - x*x)
		if got := Cross(x, y); !approxEqual(got, want) {
			t.Errorf("Cross(%v,%v) = %v, want %v", x, y, got, want)
		}
	}
}

func TestPointOnCircleLiesOnTheConstraint(t *testing.T) {
	for _, theta := range []float64{0, 30, 45, 90, 200, 359} {
		for _, r := range []float64{0.5, 1, 2, 3} {
			x, y := PointOnCircle(r, theta)
			if got := x*x + y*y; !approxEqual(got, r*r) {
				t.Errorf("PointOnCircle(r=%v,theta=%v): x²+y²=%v, want %v", r, theta, got, r*r)
			}
		}
	}
}

func TestPointOnCircleKnownAngles(t *testing.T) {
	x, y := PointOnCircle(2, 0)
	if !approxEqual(x, 2) || !approxEqual(y, 0) {
		t.Errorf("PointOnCircle(2,0) = (%v,%v), want (2,0)", x, y)
	}
	x, y = PointOnCircle(2, 90)
	if !approxEqual(x, 0) || !approxEqual(y, 2) {
		t.Errorf("PointOnCircle(2,90) = (%v,%v), want (0,2)", x, y)
	}
}

func TestContourPointsSatisfyXYEqualsK(t *testing.T) {
	pts := ContourPoints(4, 0.5, 3, 20)
	if len(pts) == 0 {
		t.Fatal("ContourPoints returned no points")
	}
	for _, pt := range pts {
		if got := pt[0] * pt[1]; !approxEqual(got, 4) {
			t.Errorf("ContourPoints(4,...) point (%v,%v): x*y = %v, want 4", pt[0], pt[1], got)
		}
	}
}

func TestContourPointsSkipsNearZeroX(t *testing.T) {
	pts := ContourPoints(1, -0.02, 0.02, 4)
	for _, pt := range pts {
		if math.Abs(pt[0]) < 0.05 {
			t.Errorf("ContourPoints returned a point with |x|=%v < 0.05, should have been skipped", pt[0])
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("lagrange-multipliers")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
