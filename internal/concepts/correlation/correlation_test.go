package correlation

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestPearsonRKnownValues(t *testing.T) {
	// Perfect uphill line.
	if r := PearsonR([]float64{1, 2, 3, 4}, []float64{2, 4, 6, 8}); math.Abs(r-1) > 1e-9 {
		t.Errorf("perfect uphill line: got r=%v, want 1", r)
	}
	// Perfect downhill line.
	if r := PearsonR([]float64{1, 2, 3, 4}, []float64{8, 6, 4, 2}); math.Abs(r-(-1)) > 1e-9 {
		t.Errorf("perfect downhill line: got r=%v, want -1", r)
	}
	// Zero variance in y is undefined correlation, not a divide-by-zero panic.
	if r := PearsonR([]float64{1, 2, 3, 4}, []float64{5, 5, 5, 5}); r != 0 {
		t.Errorf("constant y: got r=%v, want 0", r)
	}
}

func TestPearsonRMismatchedLengths(t *testing.T) {
	if r := PearsonR([]float64{1, 2}, []float64{1}); r != 0 {
		t.Errorf("mismatched lengths: got r=%v, want 0", r)
	}
}

func TestGeneratePointsLength(t *testing.T) {
	xs, ys := GeneratePoints(0.5, 137)
	if len(xs) != 137 || len(ys) != 137 {
		t.Fatalf("got %d/%d points, want 137/137", len(xs), len(ys))
	}
}

func TestGeneratePointsIsDeterministic(t *testing.T) {
	xs1, ys1 := GeneratePoints(0.4, 50)
	xs2, ys2 := GeneratePoints(0.4, 50)
	for i := range xs1 {
		if xs1[i] != xs2[i] || ys1[i] != ys2[i] {
			t.Fatalf("GeneratePoints(0.4, 50) was not deterministic at i=%d", i)
		}
	}
}

func TestGeneratePointsExtremesAreExactLines(t *testing.T) {
	// At r=+1/-1 there is no noise term left (sqrt(1-r^2)=0), so y must track
	// x exactly.
	xs, ys := GeneratePoints(1, 40)
	for i := range xs {
		if math.Abs(ys[i]-xs[i]) > 1e-9 {
			t.Fatalf("r=1: y[%d]=%v != x[%d]=%v", i, ys[i], i, xs[i])
		}
	}
	xs, ys = GeneratePoints(-1, 40)
	for i := range xs {
		if math.Abs(ys[i]+xs[i]) > 1e-9 {
			t.Fatalf("r=-1: y[%d]=%v != -x[%d]=%v", i, ys[i], i, xs[i])
		}
	}
}

func TestSampleCorrelationTracksTargetR(t *testing.T) {
	// The empirical Pearson r of a generated cloud should land close to the
	// target r, and should rise as the target r rises.
	const n = 400
	var last float64 = -2
	for _, target := range []float64{-0.9, -0.5, 0, 0.5, 0.9} {
		xs, ys := GeneratePoints(target, n)
		got := PearsonR(xs, ys)
		if math.Abs(got-target) > 0.15 {
			t.Errorf("target r=%.2f: sample r=%.3f, too far off", target, got)
		}
		if got < last-1e-9 {
			t.Errorf("sample r did not rise with target r: target=%.2f got=%.3f after %.3f", target, got, last)
		}
		last = got
	}
}

func TestPearsonRBounded(t *testing.T) {
	for _, target := range []float64{-1, -0.3, 0, 0.3, 1} {
		xs, ys := GeneratePoints(target, 100)
		r := PearsonR(xs, ys)
		if r < -1-1e-9 || r > 1+1e-9 {
			t.Errorf("target=%.2f: sample r=%v out of [-1,1]", target, r)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("correlation")
	if !ok {
		t.Fatal("concept \"correlation\" not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
