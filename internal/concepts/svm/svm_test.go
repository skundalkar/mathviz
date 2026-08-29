package svm

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func approxEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestNormalIsUnitLength(t *testing.T) {
	for _, theta := range []float64{0, 30, 45, 90, 150, 179} {
		nx, ny := Normal(theta)
		if got := math.Hypot(nx, ny); !approxEqual(got, 1) {
			t.Errorf("Normal(%v) has length %v, want 1", theta, got)
		}
	}
}

func TestMaxMarginFromPairMatchesWorkedExample(t *testing.T) {
	nx, ny, d, margin := MaxMarginFromPair(Points[0], Points[2]) // (3,3) and (1,1)
	want := 1 / math.Sqrt2
	if !approxEqual(nx, want) || !approxEqual(ny, want) {
		t.Errorf("MaxMarginFromPair normal = (%v,%v), want (%v,%v)", nx, ny, want, want)
	}
	if wantD := 2 * math.Sqrt2; !approxEqual(d, wantD) {
		t.Errorf("MaxMarginFromPair d = %v, want %v", d, wantD)
	}
	if wantMargin := math.Sqrt(8); !approxEqual(margin, wantMargin) {
		t.Errorf("MaxMarginFromPair margin = %v, want %v", margin, wantMargin)
	}
}

func TestBestHyperplaneSeparatesAllFourPoints(t *testing.T) {
	nx, ny, d, _ := MaxMarginFromPair(Points[0], Points[2])
	if !Separates(Points, Labels, nx, ny, d) {
		t.Fatal("best hyperplane should separate all four points")
	}
	if got := Margin(Points, Labels, nx, ny, d); !approxEqual(got, math.Sqrt(8)) {
		t.Errorf("Margin on full dataset = %v, want %v (support vectors are the tightest pair)", got, math.Sqrt(8))
	}
}

func TestSupportVectorsHaveTheSmallestFunctionalMargin(t *testing.T) {
	nx, ny, d, _ := MaxMarginFromPair(Points[0], Points[2])
	svMargin := FunctionalMargin(nx, ny, d, Points[0][0], Points[0][1], Labels[0])
	for i, pt := range Points {
		fm := FunctionalMargin(nx, ny, d, pt[0], pt[1], Labels[i])
		if fm < svMargin-1e-9 {
			t.Errorf("point %v has functional margin %v, smaller than support vector's %v", pt, fm, svMargin)
		}
	}
	// The two non-support points should sit strictly farther from the line.
	for _, i := range []int{1, 3} {
		fm := FunctionalMargin(nx, ny, d, Points[i][0], Points[i][1], Labels[i])
		if fm <= svMargin+1e-9 {
			t.Errorf("non-support point %v has functional margin %v, want strictly greater than %v", Points[i], fm, svMargin)
		}
	}
}

func TestSeparatesFalseWhenLineMisclassifiesAPoint(t *testing.T) {
	nx, ny := Normal(45)
	// d=-5 pushes the line far enough that (1,1) and (-1,0) end up on the
	// wrong side.
	if Separates(Points, Labels, nx, ny, -5) {
		t.Fatal("expected a badly-placed line not to separate the classes")
	}
	if got := Margin(Points, Labels, nx, ny, -5); got >= 0 {
		t.Errorf("Margin = %v, want negative (misclassifies a point)", got)
	}
}

func TestLinePointsLieOnTheLine(t *testing.T) {
	nx, ny := Normal(30)
	const d = 1.5
	for _, pt := range LinePoints(nx, ny, d, -4, 4, 10) {
		if got := nx*pt[0] + ny*pt[1]; !approxEqual(got, d) {
			t.Errorf("LinePoints point %v: n·x = %v, want %v", pt, got, d)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("support-vector-machine")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
