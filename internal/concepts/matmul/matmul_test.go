package matmul

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func approxEq(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestRotationKnownValues(t *testing.T) {
	x, y := Apply(Rotation(90), 1, 0)
	if !approxEq(x, 0) || !approxEq(y, 1) {
		t.Fatalf("Rotation(90)*(1,0) = (%v,%v), want (0,1)", x, y)
	}
	x, y = Apply(Rotation(180), 1, 0)
	if !approxEq(x, -1) || !approxEq(y, 0) {
		t.Fatalf("Rotation(180)*(1,0) = (%v,%v), want (-1,0)", x, y)
	}
}

func TestShearKnownValues(t *testing.T) {
	x, y := Apply(Shear(2), 1, 1)
	if !approxEq(x, 3) || !approxEq(y, 1) {
		t.Fatalf("Shear(2)*(1,1) = (%v,%v), want (3,1)", x, y)
	}
	// A point on the x-axis (height 0) is untouched by a shear along x.
	x, y = Apply(Shear(5), 1, 0)
	if !approxEq(x, 1) || !approxEq(y, 0) {
		t.Fatalf("Shear(5)*(1,0) = (%v,%v), want (1,0)", x, y)
	}
}

// TestMulMatchesChainedApply is the core invariant this concept teaches:
// applying the product matrix directly must equal chaining the two applies.
func TestMulMatchesChainedApply(t *testing.T) {
	cases := []struct{ angleA, shearB, x, y float64 }{
		{90, 1, 1, 0},
		{45, -0.7, 2, 3},
		{-30, 1.5, -1, 4},
		{0, 0, 5, -5},
	}
	for _, tc := range cases {
		A, B := Rotation(tc.angleA), Shear(tc.shearB)

		chainedX, chainedY := Apply(B, tc.x, tc.y)
		chainedX, chainedY = Apply(A, chainedX, chainedY)

		C := Mul(A, B)
		directX, directY := Apply(C, tc.x, tc.y)

		if !approxEq(chainedX, directX) || !approxEq(chainedY, directY) {
			t.Fatalf("angleA=%v shearB=%v v=(%v,%v): chained (%v,%v) != direct (%v,%v)",
				tc.angleA, tc.shearB, tc.x, tc.y, chainedX, chainedY, directX, directY)
		}
	}
}

// TestOrderMatters reproduces the worked example from the "How does it
// actually work?" section: A*B and B*A send (1,0) to different points.
func TestOrderMatters(t *testing.T) {
	A, B := Rotation(90), Shear(1)

	abX, abY := Apply(Mul(A, B), 1, 0)
	if !approxEq(abX, 0) || !approxEq(abY, 1) {
		t.Fatalf("A*B*(1,0) = (%v,%v), want (0,1)", abX, abY)
	}

	baX, baY := Apply(Mul(B, A), 1, 0)
	if !approxEq(baX, 1) || !approxEq(baY, 1) {
		t.Fatalf("B*A*(1,0) = (%v,%v), want (1,1)", baX, baY)
	}

	if approxEq(abX, baX) && approxEq(abY, baY) {
		t.Fatal("A*B and B*A produced the same result — order should matter here")
	}
}

// TestColumnsAreBasisImages checks the "read a matrix by its columns"
// claim: column 1 of M is M*i, column 2 is M*j.
func TestColumnsAreBasisImages(t *testing.T) {
	M := Mul(Rotation(90), Shear(1)) // = [[0,-1],[1,1]]
	ix, iy := Apply(M, 1, 0)
	if !approxEq(ix, M.A) || !approxEq(iy, M.C) {
		t.Fatalf("M*i = (%v,%v), want column 1 (%v,%v)", ix, iy, M.A, M.C)
	}
	jx, jy := Apply(M, 0, 1)
	if !approxEq(jx, M.B) || !approxEq(jy, M.D) {
		t.Fatalf("M*j = (%v,%v), want column 2 (%v,%v)", jx, jy, M.B, M.D)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("matrix-multiplication")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
