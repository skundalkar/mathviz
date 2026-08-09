package gradientdescent

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestGradMatchesFiniteDifference(t *testing.T) {
	const h = 1e-6
	for _, x := range []float64{-3, -0.5, 0, 1, 4.2} {
		approx := (F(x+h) - F(x-h)) / (2 * h)
		if got := Grad(x); math.Abs(got-approx) > 1e-4 {
			t.Fatalf("Grad(%v) = %v, want ≈ %v (finite difference)", x, got, approx)
		}
	}
}

func TestDescendMatchesClosedForm(t *testing.T) {
	// x_t = x0 * (1 - 2*lr)^t for F(x) = x^2.
	x0, lr, steps := StartX, 0.2, 10
	path := Descend(x0, lr, steps)
	if len(path) != steps+1 {
		t.Fatalf("len(path) = %d, want %d", len(path), steps+1)
	}
	ratio := 1 - 2*lr
	for t2 := 0; t2 <= steps; t2++ {
		want := x0 * math.Pow(ratio, float64(t2))
		if got := path[t2]; math.Abs(got-want) > 1e-9 {
			t.Fatalf("path[%d] = %v, want %v", t2, got, want)
		}
	}
}

func TestDescendConvergesForSmallLearningRate(t *testing.T) {
	path := Descend(StartX, 0.1, 20)
	for i := 1; i < len(path); i++ {
		if math.Abs(path[i]) >= math.Abs(path[i-1]) {
			t.Fatalf("step %d: |x|=%v did not shrink from |x|=%v", i, path[i], path[i-1])
		}
	}
	if math.Abs(path[len(path)-1]) > 0.1 {
		t.Fatalf("final position %v did not approach the minimum at 0", path[len(path)-1])
	}
}

func TestDescendDivergesForLargeLearningRate(t *testing.T) {
	// lr > 1 makes |1 - 2*lr| > 1, so each step should move further from 0.
	path := Descend(StartX, 1.1, 10)
	for i := 1; i < len(path); i++ {
		if math.Abs(path[i]) <= math.Abs(path[i-1]) {
			t.Fatalf("step %d: |x|=%v did not grow from |x|=%v (expected divergence)", i, path[i], path[i-1])
		}
	}
}

func TestDescendStableAtOptimalLearningRate(t *testing.T) {
	// lr = 0.5 makes 1 - 2*lr = 0: the very next step lands exactly on the
	// minimum for a pure quadratic.
	path := Descend(StartX, 0.5, 3)
	for i := 1; i < len(path); i++ {
		if math.Abs(path[i]) > 1e-9 {
			t.Fatalf("step %d: x=%v, want 0 at lr=0.5", i, path[i])
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("gradient-descent")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}

func TestRenderHandlesDivergentAndZeroSteps(t *testing.T) {
	c, ok := concept.Get("gradient-descent")
	if !ok {
		t.Fatal("concept not registered")
	}
	for _, p := range []map[string]float64{
		{"lr": 1.2, "steps": 30},
		{"lr": 0.01, "steps": 1},
	} {
		out := c.Render(p)
		if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
			t.Fatalf("render(%v) did not produce a well-formed svg", p)
		}
	}
}
