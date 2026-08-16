package montecarlo

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestSamplePointStaysInSquare(t *testing.T) {
	for i := 1; i <= 2000; i++ {
		x, y := SamplePoint(i, 0)
		if x < -1 || x > 1 || y < -1 || y > 1 {
			t.Fatalf("SamplePoint(%d, 0) = (%v, %v), want within [-1,1]x[-1,1]", i, x, y)
		}
	}
}

func TestSamplePointIsDeterministic(t *testing.T) {
	for i := 1; i <= 20; i++ {
		x1, y1 := SamplePoint(i, 3)
		x2, y2 := SamplePoint(i, 3)
		if x1 != x2 || y1 != y2 {
			t.Errorf("SamplePoint(%d,3) not deterministic: (%v,%v) then (%v,%v)", i, x1, y1, x2, y2)
		}
	}
}

func TestSamplePointRunsDontCollide(t *testing.T) {
	seen := make(map[[2]float64]bool)
	for run := 0; run < 3; run++ {
		for i := 1; i <= 50; i++ {
			x, y := SamplePoint(i, run)
			key := [2]float64{x, y}
			if seen[key] {
				t.Fatalf("run=%d i=%d produced a point already seen in an earlier run: (%v,%v)", run, i, x, y)
			}
			seen[key] = true
		}
	}
}

func TestInsideKnownPoints(t *testing.T) {
	cases := []struct {
		x, y float64
		want bool
	}{
		{0, 0, true},      // center
		{1, 0, true},      // exactly on the boundary
		{0, 1, true},      // exactly on the boundary
		{0.7, 0.7, true},  // 0.49+0.49=0.98 < 1, just inside
		{0.8, 0.8, false}, // 0.64+0.64=1.28 > 1, just outside
		{1, 1, false},     // corner of the square, well outside
		{-1, -1, false},   // opposite corner
	}
	for _, c := range cases {
		if got := Inside(c.x, c.y); got != c.want {
			t.Errorf("Inside(%v,%v) = %v, want %v", c.x, c.y, got, c.want)
		}
	}
}

func TestEstimatePiKnownValues(t *testing.T) {
	// Pinned to the deterministic SamplePoint sequence at run=0, so a
	// change to the sequence (or a broken estimate calculation) shows up
	// immediately. See LEARNINGS.md for the same worked numbers.
	cases := map[int]float64{5: 3.20, 20: 3.20, 100: 3.16, 200: 3.34}
	for n, want := range cases {
		if got := EstimatePi(n, 0); !almostEqual(got, want, 1e-9) {
			t.Errorf("EstimatePi(%d,0) = %v, want %v", n, got, want)
		}
	}
}

func TestEstimatePiSeriesLength(t *testing.T) {
	if got := len(EstimatePiSeries(2000, 0)); got != 2000 {
		t.Errorf("len(EstimatePiSeries(2000,0)) = %d, want 2000", got)
	}
	if got := len(EstimatePiSeries(0, 0)); got != 1 {
		t.Errorf("len(EstimatePiSeries(0,0)) = %d, want 1 (clamped)", got)
	}
}

func TestEstimatePiMatchesSeriesLastEntry(t *testing.T) {
	for _, n := range []int{1, 5, 50, 2000} {
		series := EstimatePiSeries(n, 0)
		want := series[len(series)-1]
		if got := EstimatePi(n, 0); !almostEqual(got, want, 1e-12) {
			t.Errorf("EstimatePi(%d,0) = %v, want %v (series' last entry)", n, got, want)
		}
	}
}

func TestEstimatePiStaysInPlausibleRange(t *testing.T) {
	// 4x a fraction in [0,1] is always in [0,4]; that's the hard bound.
	for run := 0; run < 5; run++ {
		for _, n := range []int{10, 100, 1000} {
			est := EstimatePi(n, run)
			if est < 0 || est > 4 {
				t.Errorf("EstimatePi(%d,%d) = %v, out of [0,4]", n, run, est)
			}
		}
	}
}

func TestEstimatePiVarianceShrinksWithN(t *testing.T) {
	// The defining Monte Carlo behavior: the running estimate bounces
	// around a lot early on and settles down later, and by a lot more
	// than the law-of-large-numbers dilution alone -- check it on this
	// concrete sequence rather than just asserting it in prose.
	series := EstimatePiSeries(2000, 0)
	variance := func(vals []float64) float64 {
		mean := 0.0
		for _, v := range vals {
			mean += v
		}
		mean /= float64(len(vals))
		sq := 0.0
		for _, v := range vals {
			sq += (v - mean) * (v - mean)
		}
		return sq / float64(len(vals))
	}
	early := variance(series[19:39])    // n=20..39
	late := variance(series[1980:2000]) // n=1981..2000
	if late >= early {
		t.Errorf("variance did not shrink: early(n=20..39)=%v, late(n=1981..2000)=%v", early, late)
	}
	// Not just smaller -- multiple orders of magnitude, matching the
	// 1/sqrt(n) rate over a 50x larger n.
	if late*100 >= early {
		t.Errorf("variance shrank less than 100x: early=%v, late=%v", early, late)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("monte-carlo-estimation")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
