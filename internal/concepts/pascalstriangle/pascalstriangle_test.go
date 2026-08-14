package pascalstriangle

import (
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestFactorialKnownValues(t *testing.T) {
	cases := map[int]int64{0: 1, 1: 1, 5: 120, 8: 40320}
	for n, want := range cases {
		if got := Factorial(n); got != want {
			t.Errorf("Factorial(%d) = %d, want %d", n, got, want)
		}
	}
}

func TestFactorialNegativeIsUndefined(t *testing.T) {
	if got := Factorial(-1); got != 0 {
		t.Errorf("Factorial(-1) = %d, want 0", got)
	}
}

func TestBinomialCoefficientWorkedExample(t *testing.T) {
	// 5!/(2!3!) = 120/(2*6) = 10.
	if got := BinomialCoefficient(5, 2); got != 10 {
		t.Errorf("BinomialCoefficient(5,2) = %d, want 10", got)
	}
}

func TestBinomialCoefficientEdges(t *testing.T) {
	for n := 0; n <= 8; n++ {
		if got := BinomialCoefficient(n, 0); got != 1 {
			t.Errorf("BinomialCoefficient(%d,0) = %d, want 1", n, got)
		}
		if got := BinomialCoefficient(n, n); got != 1 {
			t.Errorf("BinomialCoefficient(%d,%d) = %d, want 1", n, n, got)
		}
	}
}

func TestBinomialCoefficientOutOfRangeIsZero(t *testing.T) {
	if got := BinomialCoefficient(5, -1); got != 0 {
		t.Errorf("BinomialCoefficient(5,-1) = %d, want 0", got)
	}
	if got := BinomialCoefficient(5, 6); got != 0 {
		t.Errorf("BinomialCoefficient(5,6) = %d, want 0", got)
	}
}

func TestBinomialCoefficientIsSymmetric(t *testing.T) {
	// Choosing k of n leaves the other n-k out -- both are the same count.
	for n := 0; n <= 8; n++ {
		for k := 0; k <= n; k++ {
			a, b := BinomialCoefficient(n, k), BinomialCoefficient(n, n-k)
			if a != b {
				t.Errorf("BinomialCoefficient(%d,%d)=%d != BinomialCoefficient(%d,%d)=%d", n, k, a, n, n-k, b)
			}
		}
	}
}

func TestPascalTriangleKnownRow(t *testing.T) {
	tri := PascalTriangle(6)
	want := []int64{1, 5, 10, 10, 5, 1}
	got := tri[5]
	if len(got) != len(want) {
		t.Fatalf("PascalTriangle(6)[5] has %d entries, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("PascalTriangle(6)[5][%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestPascalTriangleMatchesBinomialCoefficient(t *testing.T) {
	tri := PascalTriangle(MaxRows)
	for n := 0; n < MaxRows; n++ {
		for k := 0; k <= n; k++ {
			if got, want := tri[n][k], BinomialCoefficient(n, k); got != want {
				t.Errorf("PascalTriangle entry (%d,%d) = %d, want BinomialCoefficient = %d", n, k, got, want)
			}
		}
	}
}

func TestPascalTriangleRowIsSumOfParents(t *testing.T) {
	tri := PascalTriangle(MaxRows)
	for n := 1; n < MaxRows; n++ {
		for k := 1; k < n; k++ {
			want := tri[n-1][k-1] + tri[n-1][k]
			if tri[n][k] != want {
				t.Errorf("triangle[%d][%d] = %d, want sum of parents %d", n, k, tri[n][k], want)
			}
		}
	}
}

func TestPascalTriangleRowsBelowOneClampToOne(t *testing.T) {
	tri := PascalTriangle(0)
	if len(tri) != 1 || len(tri[0]) != 1 || tri[0][0] != 1 {
		t.Errorf("PascalTriangle(0) = %v, want a single [1] row", tri)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("pascals-triangle")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
