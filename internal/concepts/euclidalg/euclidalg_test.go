package euclidalg

import (
	"testing"

	"mathviz/internal/concept"
)

func TestGCDKnownValues(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{1071, 462, 21},
		{13, 8, 1}, // consecutive Fibonacci numbers -- the algorithm's worst case
		{48, 18, 6},
		{17, 5, 1}, // coprime
		{7, 7, 7},
	}
	for _, c := range cases {
		if got := GCD(c.a, c.b); got != c.want {
			t.Errorf("GCD(%d,%d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestGCDBaseCases(t *testing.T) {
	if got := GCD(5, 0); got != 5 {
		t.Errorf("GCD(5,0) = %d, want 5", got)
	}
	if got := GCD(0, 5); got != 5 {
		t.Errorf("GCD(0,5) = %d, want 5", got)
	}
	if got := GCD(0, 0); got != 0 {
		t.Errorf("GCD(0,0) = %d, want 0", got)
	}
}

func TestGCDIsSymmetricAndSignIndependent(t *testing.T) {
	if GCD(1071, 462) != GCD(462, 1071) {
		t.Errorf("GCD not symmetric: GCD(1071,462)=%d, GCD(462,1071)=%d", GCD(1071, 462), GCD(462, 1071))
	}
	if GCD(-1071, 462) != 21 {
		t.Errorf("GCD(-1071,462) = %d, want 21", GCD(-1071, 462))
	}
}

func TestGCDDividesBothInputs(t *testing.T) {
	for a := 1; a <= 40; a++ {
		for b := 1; b <= 40; b++ {
			g := GCD(a, b)
			if a%g != 0 || b%g != 0 {
				t.Fatalf("GCD(%d,%d)=%d does not divide both", a, b, g)
			}
			// It must also be the LARGEST common divisor: no bigger
			// candidate up to min(a,b) should divide both.
			max := a
			if b < max {
				max = b
			}
			for cand := g + 1; cand <= max; cand++ {
				if a%cand == 0 && b%cand == 0 {
					t.Fatalf("GCD(%d,%d)=%d but %d is a larger common divisor", a, b, g, cand)
				}
			}
		}
	}
}

func TestStepsMatchesWorkedExample(t *testing.T) {
	want := []Step{
		{1071, 462, 2, 147},
		{462, 147, 3, 21},
		{147, 21, 7, 0},
	}
	got := Steps(1071, 462)
	if len(got) != len(want) {
		t.Fatalf("Steps(1071,462) has %d steps, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Steps(1071,462)[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestStepsLastDivisorIsGCD(t *testing.T) {
	for _, pair := range [][2]int{{1071, 462}, {13, 8}, {48, 18}} {
		steps := Steps(pair[0], pair[1])
		if len(steps) == 0 {
			t.Fatalf("Steps(%d,%d) returned no steps", pair[0], pair[1])
		}
		last := steps[len(steps)-1]
		if last.Remainder != 0 {
			t.Errorf("Steps(%d,%d) last remainder = %d, want 0", pair[0], pair[1], last.Remainder)
		}
		if last.Divisor != GCD(pair[0], pair[1]) {
			t.Errorf("Steps(%d,%d) last divisor = %d, want GCD = %d",
				pair[0], pair[1], last.Divisor, GCD(pair[0], pair[1]))
		}
	}
}

func TestTileSquaresExactlyTilesTheRectangle(t *testing.T) {
	for _, pair := range [][2]int{{1071, 462}, {13, 8}, {48, 18}, {7, 7}} {
		sq := TileSquares(pair[0], pair[1])
		var area float64
		for _, s := range sq {
			area += s.Side * s.Side
		}
		want := float64(pair[0] * pair[1])
		if area != want {
			t.Errorf("TileSquares(%d,%d) total area = %v, want %v", pair[0], pair[1], area, want)
		}
	}
}

func TestTileSquaresLastSquareIsTheGCD(t *testing.T) {
	for _, pair := range [][2]int{{1071, 462}, {13, 8}, {48, 18}} {
		sq := TileSquares(pair[0], pair[1])
		if len(sq) == 0 {
			t.Fatalf("TileSquares(%d,%d) returned nothing", pair[0], pair[1])
		}
		last := sq[len(sq)-1]
		if last.Side != float64(GCD(pair[0], pair[1])) {
			t.Errorf("TileSquares(%d,%d) last square side = %v, want GCD = %d",
				pair[0], pair[1], last.Side, GCD(pair[0], pair[1]))
		}
	}
}

func TestTileSquaresStepCountMatchesSteps(t *testing.T) {
	sq := TileSquares(1071, 462)
	maxStep := -1
	for _, s := range sq {
		if s.Step > maxStep {
			maxStep = s.Step
		}
	}
	if want := len(Steps(1071, 462)) - 1; maxStep != want {
		t.Errorf("TileSquares(1071,462) max Step = %d, want %d", maxStep, want)
	}
}

func TestTileSquaresEmptyForNonPositiveInputs(t *testing.T) {
	if got := TileSquares(0, 5); got != nil {
		t.Errorf("TileSquares(0,5) = %v, want nil", got)
	}
	if got := TileSquares(5, 0); got != nil {
		t.Errorf("TileSquares(5,0) = %v, want nil", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("euclidean-algorithm")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
