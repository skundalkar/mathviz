package chisquared

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestObservedTableKnownValues(t *testing.T) {
	// The worked example: n=50, 30% vs 50% success rates.
	got := ObservedTable(50, 0.30, 0.50)
	want := Table2x2{{15, 35}, {25, 25}}
	if got != want {
		t.Errorf("ObservedTable(50, 0.30, 0.50) = %v, want %v", got, want)
	}
}

func TestExpectedUnderIndependence(t *testing.T) {
	observed := ObservedTable(50, 0.30, 0.50)
	got := Expected(observed)
	want := Table2x2{{20, 30}, {20, 30}}
	if got != want {
		t.Errorf("Expected(%v) = %v, want %v", observed, got, want)
	}
}

func TestExpectedPreservesRowAndColTotals(t *testing.T) {
	observed := Table2x2{{12, 8}, {3, 27}} // an asymmetric table
	expected := Expected(observed)
	for r := 0; r < 2; r++ {
		if !near(RowTotal(observed, r), RowTotal(expected, r), 1e-9) {
			t.Errorf("row %d total changed: observed=%v expected=%v", r, RowTotal(observed, r), RowTotal(expected, r))
		}
	}
	for c := 0; c < 2; c++ {
		if !near(ColTotal(observed, c), ColTotal(expected, c), 1e-9) {
			t.Errorf("col %d total changed: observed=%v expected=%v", c, ColTotal(observed, c), ColTotal(expected, c))
		}
	}
}

func TestChiSquareStatKnownValue(t *testing.T) {
	// Hand-computed in LEARNINGS.md: (15-20)²/20+(35-30)²/30+(25-20)²/20+(25-30)²/30
	// = 1.25+0.8333+1.25+0.8333 = 4.1667.
	observed := ObservedTable(50, 0.30, 0.50)
	expected := Expected(observed)
	got := ChiSquareStat(observed, expected)
	if !near(got, 4.1667, 1e-3) {
		t.Errorf("ChiSquareStat = %v, want ~4.1667", got)
	}
}

func TestChiSquareStatSameSampleSmallerN(t *testing.T) {
	// Same 30%/50% rates at n=10: chi2 drops to ~0.8333.
	observed := ObservedTable(10, 0.30, 0.50)
	expected := Expected(observed)
	got := ChiSquareStat(observed, expected)
	if !near(got, 0.8333, 1e-3) {
		t.Errorf("ChiSquareStat(n=10) = %v, want ~0.8333", got)
	}
}

func TestChiSquareStatZeroWhenObservedEqualsExpected(t *testing.T) {
	// Equal rates in both groups -> observed IS the independence table ->
	// chi2 must be exactly 0.
	observed := ObservedTable(40, 0.5, 0.5)
	expected := Expected(observed)
	got := ChiSquareStat(observed, expected)
	if got != 0 {
		t.Errorf("ChiSquareStat with matched rates = %v, want exactly 0", got)
	}
}

func TestChiSquareCDFDF1KnownCriticalValues(t *testing.T) {
	// Textbook chi-squared(1) critical values: 3.841 -> p=0.05, 6.635 -> p=0.01.
	if got := PValueDF1(3.841); !near(got, 0.05, 1e-3) {
		t.Errorf("PValueDF1(3.841) = %v, want ~0.05", got)
	}
	if got := PValueDF1(6.635); !near(got, 0.01, 1e-3) {
		t.Errorf("PValueDF1(6.635) = %v, want ~0.01", got)
	}
}

func TestPValueDF1AtWorkedExamples(t *testing.T) {
	if got := PValueDF1(4.1667); !near(got, 0.0412, 1e-3) {
		t.Errorf("PValueDF1(4.1667) = %v, want ~0.0412", got)
	}
	if got := PValueDF1(0.8333); !near(got, 0.3613, 1e-3) {
		t.Errorf("PValueDF1(0.8333) = %v, want ~0.3613", got)
	}
}

func TestPValueDF1BoundsAndMonotonicity(t *testing.T) {
	if got := PValueDF1(0); got != 1 {
		t.Errorf("PValueDF1(0) = %v, want 1 (no deviation at all is completely unsurprising)", got)
	}
	prev := 1.0
	for x := 0.5; x <= 20; x += 0.5 {
		p := PValueDF1(x)
		if p > prev {
			t.Errorf("PValueDF1(%v) = %v is greater than the value at a smaller x (%v); must be non-increasing", x, p, prev)
		}
		if p < 0 || p > 1 {
			t.Errorf("PValueDF1(%v) = %v out of [0,1]", x, p)
		}
		prev = p
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("chi-squared-test")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
