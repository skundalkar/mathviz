package tdist

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestStudentTPDFIsSymmetric(t *testing.T) {
	for _, df := range []float64{1, 2, 5, 30} {
		for _, x := range []float64{0.5, 1, 2, 4} {
			pos, neg := StudentTPDF(x, df), StudentTPDF(-x, df)
			if !near(pos, neg, 1e-9) {
				t.Errorf("StudentTPDF(%v,%v)=%v != StudentTPDF(%v,%v)=%v", x, df, pos, -x, df, neg)
			}
		}
	}
}

func TestStudentTPDFNonPositiveDFIsZero(t *testing.T) {
	if got := StudentTPDF(0, 0); got != 0 {
		t.Errorf("StudentTPDF(0,0) = %v, want 0", got)
	}
	if got := StudentTPDF(0, -3); got != 0 {
		t.Errorf("StudentTPDF(0,-3) = %v, want 0", got)
	}
}

func TestStudentTPDFApproachesNormalAsDFGrows(t *testing.T) {
	// The whole point of the concept: at large df, t's density is
	// indistinguishable from the standard normal's.
	for _, x := range []float64{-2, -0.5, 0, 0.5, 2} {
		got, want := StudentTPDF(x, 100000), StdNormalPDF(x)
		if !near(got, want, 1e-3) {
			t.Errorf("StudentTPDF(%v,100000) = %v, want ~StdNormalPDF(%v) = %v", x, got, x, want)
		}
	}
}

func TestStudentTPDFFatterThanNormalAtSmallDF(t *testing.T) {
	// At a fixed, moderately-far-out x, low df must put MORE density (a
	// "fatter tail") than the normal does -- that's the entire lesson.
	x := 3.0
	tDensity := StudentTPDF(x, 2)
	normalDensity := StdNormalPDF(x)
	if tDensity <= normalDensity {
		t.Errorf("StudentTPDF(%v,df=2) = %v, want > StdNormalPDF(%v) = %v (fatter tail)",
			x, tDensity, x, normalDensity)
	}
}

func TestStudentTCDFIsOneHalfAtZero(t *testing.T) {
	for _, df := range []float64{1, 2, 5, 10, 30} {
		if got := StudentTCDF(0, df); !near(got, 0.5, 1e-3) {
			t.Errorf("StudentTCDF(0,%v) = %v, want 0.5", df, got)
		}
	}
}

func TestStudentTCDFIsMonotonic(t *testing.T) {
	prev := StudentTCDF(-10, 5)
	for x := -9.0; x <= 10; x++ {
		got := StudentTCDF(x, 5)
		if got < prev {
			t.Errorf("StudentTCDF not monotonic near x=%v: %v < previous %v", x, got, prev)
		}
		prev = got
	}
}

func TestCriticalValueMatchesKnownTTable(t *testing.T) {
	// Standard two-tailed 95% critical values from any t-table.
	cases := map[float64]float64{
		1:  12.706,
		2:  4.303,
		5:  2.571,
		10: 2.228,
		30: 2.042,
	}
	for df, want := range cases {
		if got := CriticalValue(df, 0.05); !near(got, want, 0.01) {
			t.Errorf("CriticalValue(%v,0.05) = %v, want %v", df, got, want)
		}
	}
}

func TestCriticalValueApproachesZAsDFGrows(t *testing.T) {
	// z* for a 95% two-tailed normal interval is 1.959964...
	if got := CriticalValue(1000000, 0.05); !near(got, 1.95996, 0.01) {
		t.Errorf("CriticalValue(1e6,0.05) = %v, want ~1.960", got)
	}
}

func TestCriticalValueShrinksAsDFGrows(t *testing.T) {
	prev := CriticalValue(1, 0.05)
	for _, df := range []float64{2, 5, 10, 30, 100} {
		got := CriticalValue(df, 0.05)
		if got >= prev {
			t.Errorf("CriticalValue(%v,0.05) = %v did not shrink from previous %v", df, got, prev)
		}
		prev = got
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("t-distribution")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
