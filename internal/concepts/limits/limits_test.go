package limits

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestRemovableHoleIsUndefinedAtTheHole(t *testing.T) {
	got := RemovableHole(HoleAt)
	if !math.IsNaN(got) {
		t.Errorf("RemovableHole(%v) = %v, want NaN (0/0)", HoleAt, got)
	}
}

func TestRemovableHoleMatchesSimplifiedFormNearby(t *testing.T) {
	// Away from the hole, (x^2-1)/(x-1) should exactly equal x+1.
	for _, x := range []float64{0, 0.5, 0.9, 0.99, 1.01, 1.1, 2, 5} {
		got := RemovableHole(x)
		want := x + 1
		if !almostEqual(got, want, 1e-9) {
			t.Errorf("RemovableHole(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestRemovableHoleConvergesToLimitFromBothSides(t *testing.T) {
	limit := RemovableHoleLimit()
	if !almostEqual(limit, 2, 1e-9) {
		t.Errorf("RemovableHoleLimit() = %v, want 2", limit)
	}
	prevLeftGap := math.Inf(1)
	prevRightGap := math.Inf(1)
	for _, h := range []float64{0.1, 0.01, 0.001} {
		leftGap := math.Abs(RemovableHole(HoleAt-h) - limit)
		rightGap := math.Abs(RemovableHole(HoleAt+h) - limit)
		if leftGap >= prevLeftGap {
			t.Errorf("left gap should shrink as h shrinks: h=%v gap=%v, prev=%v", h, leftGap, prevLeftGap)
		}
		if rightGap >= prevRightGap {
			t.Errorf("right gap should shrink as h shrinks: h=%v gap=%v, prev=%v", h, rightGap, prevRightGap)
		}
		prevLeftGap, prevRightGap = leftGap, rightGap
	}
}

func TestStepKnownValues(t *testing.T) {
	cases := []struct {
		x    float64
		want float64
	}{
		{-1, -1}, {0.9, -1}, {0.999, -1}, {1, 1}, {1.001, 1}, {2, 1},
	}
	for _, c := range cases {
		if got := Step(c.x); got != c.want {
			t.Errorf("Step(%v) = %v, want %v", c.x, got, c.want)
		}
	}
}

func TestLeftRightLimitsAgreeForRemovableHole(t *testing.T) {
	if LeftLimit(0) != RightLimit(0) {
		t.Errorf("mode 0: LeftLimit=%v RightLimit=%v, want equal", LeftLimit(0), RightLimit(0))
	}
	if !LimitExists(0) {
		t.Error("LimitExists(0) = false, want true (removable hole's limit exists)")
	}
}

func TestLeftRightLimitsDisagreeForStep(t *testing.T) {
	if LeftLimit(1) == RightLimit(1) {
		t.Errorf("mode 1: LeftLimit=%v RightLimit=%v, want different", LeftLimit(1), RightLimit(1))
	}
	if LimitExists(1) {
		t.Error("LimitExists(1) = true, want false (step function's two-sided limit doesn't exist)")
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("limits")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
