package vectornorms

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// TestNormsKnownValues is pinned to the worked v=(3,4) example in
// LEARNINGS.md and the concept's own Sections.
func TestNormsKnownValues(t *testing.T) {
	if got := NormL1(3, 4); !almostEqual(got, 7, 1e-9) {
		t.Errorf("NormL1(3,4) = %v, want 7", got)
	}
	if got := NormL2(3, 4); !almostEqual(got, 5, 1e-9) {
		t.Errorf("NormL2(3,4) = %v, want 5", got)
	}
	if got := NormLInf(3, 4); !almostEqual(got, 4, 1e-9) {
		t.Errorf("NormLInf(3,4) = %v, want 4", got)
	}
}

// TestNormLpMatchesL1AndL2 checks that the general formula collapses to
// the named special cases at p=1 and p=2, for several vectors.
func TestNormLpMatchesL1AndL2(t *testing.T) {
	vectors := [][2]float64{{3, 4}, {1, 1}, {-2, 5}, {0, 7}, {-3, -3}}
	for _, v := range vectors {
		x, y := v[0], v[1]
		if got, want := NormLp(x, y, 1), NormL1(x, y); !almostEqual(got, want, 1e-9) {
			t.Errorf("NormLp(%v,%v,1) = %v, want NormL1 = %v", x, y, got, want)
		}
		if got, want := NormLp(x, y, 2), NormL2(x, y); !almostEqual(got, want, 1e-9) {
			t.Errorf("NormLp(%v,%v,2) = %v, want NormL2 = %v", x, y, got, want)
		}
	}
}

// TestNormLpKnownValues pins the LEARNINGS.md p-sweep table for v=(3,4).
func TestNormLpKnownValues(t *testing.T) {
	cases := []struct {
		p, want float64
	}{
		{1, 7.0000},
		{1.5, 5.5843},
		{2, 5.0000},
		{3, 4.4979},
		{5, 4.1740},
		{10, 4.0220},
	}
	for _, c := range cases {
		if got := NormLp(3, 4, c.p); !almostEqual(got, c.want, 1e-4) {
			t.Errorf("NormLp(3,4,%v) = %v, want %v", c.p, got, c.want)
		}
	}
}

// TestNormLpConvergesToLInf checks that, without overflowing, NormLp keeps
// sliding toward NormLInf as p grows very large -- the limiting behavior
// the concept's "How does it actually work?" section describes.
func TestNormLpConvergesToLInf(t *testing.T) {
	x, y := 3.0, 4.0
	linf := NormLInf(x, y)
	prev := NormLp(x, y, 2)
	for _, p := range []float64{5, 10, 50, 200, 1000, 1e6} {
		got := NormLp(x, y, p)
		if math.IsInf(got, 0) || math.IsNaN(got) {
			t.Fatalf("NormLp(%v,%v,%v) = %v, want a finite number", x, y, p, got)
		}
		if got > prev+1e-9 {
			t.Errorf("NormLp(%v,%v,%v) = %v is larger than NormLp at a smaller p (%v); want non-increasing in p", x, y, p, got, prev)
		}
		if got < linf-1e-6 {
			t.Errorf("NormLp(%v,%v,%v) = %v dropped below NormLInf = %v", x, y, p, got, linf)
		}
		prev = got
	}
	if got := NormLp(x, y, 1e6); !almostEqual(got, linf, 1e-3) {
		t.Errorf("NormLp(%v,%v,1e6) = %v, want within 1e-3 of NormLInf = %v", x, y, got, linf)
	}
}

// TestNormOrdering checks the defining inequality L∞ <= L2 <= L1 that
// holds for every vector, not just the pinned (3,4) example.
func TestNormOrdering(t *testing.T) {
	vectors := [][2]float64{{3, 4}, {1, 1}, {-2, 5}, {0, 7}, {-3, -3}, {0, 0}, {10, -1}}
	for _, v := range vectors {
		x, y := v[0], v[1]
		l1, l2, linf := NormL1(x, y), NormL2(x, y), NormLInf(x, y)
		if linf > l2+1e-9 {
			t.Errorf("(%v,%v): NormLInf = %v > NormL2 = %v", x, y, linf, l2)
		}
		if l2 > l1+1e-9 {
			t.Errorf("(%v,%v): NormL2 = %v > NormL1 = %v", x, y, l2, l1)
		}
	}
}

func TestNormsAtOrigin(t *testing.T) {
	if got := NormL1(0, 0); got != 0 {
		t.Errorf("NormL1(0,0) = %v, want 0", got)
	}
	if got := NormL2(0, 0); got != 0 {
		t.Errorf("NormL2(0,0) = %v, want 0", got)
	}
	if got := NormLInf(0, 0); got != 0 {
		t.Errorf("NormLInf(0,0) = %v, want 0", got)
	}
	if got := NormLp(0, 0, 3.5); got != 0 {
		t.Errorf("NormLp(0,0,3.5) = %v, want 0", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("vector-norms")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
