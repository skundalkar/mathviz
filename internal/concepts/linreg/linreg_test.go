package linreg

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

const epsilon = 1e-9

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestMeanKnownValues(t *testing.T) {
	if got := Mean(HoursStudied); !almostEqual(got, 3, epsilon) {
		t.Errorf("Mean(HoursStudied) = %v, want 3", got)
	}
	if got := Mean(QuizScore); !almostEqual(got, 64, epsilon) {
		t.Errorf("Mean(QuizScore) = %v, want 64", got)
	}
}

func TestSlopeInterceptWorkedExample(t *testing.T) {
	// From the section-2 walkthrough: slope = 7.5, intercept = 41.5.
	if got := Slope(HoursStudied, QuizScore); !almostEqual(got, 7.5, epsilon) {
		t.Errorf("Slope(HoursStudied, QuizScore) = %v, want 7.5", got)
	}
	if got := Intercept(HoursStudied, QuizScore); !almostEqual(got, 41.5, epsilon) {
		t.Errorf("Intercept(HoursStudied, QuizScore) = %v, want 41.5", got)
	}
}

func TestSlopeHorizontalLineIsZero(t *testing.T) {
	xs := []float64{1, 2, 3, 4, 5}
	ys := []float64{10, 10, 10, 10, 10}
	if got := Slope(xs, ys); got != 0 {
		t.Errorf("Slope with constant y = %v, want 0", got)
	}
}

func TestPredictWorkedExample(t *testing.T) {
	// The opening scenario's payoff: a student who studies 3.5 hours.
	got := Predict(7.5, 41.5, 3.5)
	if !almostEqual(got, 67.75, epsilon) {
		t.Errorf("Predict(7.5, 41.5, 3.5) = %v, want 67.75", got)
	}
}

func TestPredictAtDataPointsMatchesWorkedResiduals(t *testing.T) {
	slope, intercept := Slope(HoursStudied, QuizScore), Intercept(HoursStudied, QuizScore)
	wantResiduals := []float64{1, -1.5, 1, -1.5, 1}
	sum := 0.0
	for i := range HoursStudied {
		pred := Predict(slope, intercept, HoursStudied[i])
		residual := QuizScore[i] - pred
		if !almostEqual(residual, wantResiduals[i], epsilon) {
			t.Errorf("residual at x=%v = %v, want %v", HoursStudied[i], residual, wantResiduals[i])
		}
		sum += residual
	}
	if !almostEqual(sum, 0, epsilon) {
		t.Errorf("sum of residuals = %v, want 0 (least-squares line passes through (x̄, ȳ))", sum)
	}
}

func TestSumSquaredErrorIsZeroForPerfectFit(t *testing.T) {
	xs := []float64{0, 1, 2, 3}
	ys := []float64{1, 3, 5, 7} // exactly y = 2x + 1
	if got := SumSquaredError(xs, ys, 2, 1); !almostEqual(got, 0, epsilon) {
		t.Errorf("SumSquaredError on a perfect fit = %v, want 0", got)
	}
}

func TestLeastSquaresMinimizesSumSquaredError(t *testing.T) {
	slope, intercept := Slope(HoursStudied, QuizScore), Intercept(HoursStudied, QuizScore)
	best := SumSquaredError(HoursStudied, QuizScore, slope, intercept)

	// Any other slope/intercept pair should do no better than the
	// least-squares line — this is the picture's central claim: no matter
	// how you drag the guess sliders, you can't beat SSE(best).
	guesses := []struct{ slope, intercept float64 }{
		{5, 50}, {10, 30}, {0, 64}, {7.5, 40}, {7.5, 43}, {8, 41.5}, {-2, 70},
	}
	for _, g := range guesses {
		sse := SumSquaredError(HoursStudied, QuizScore, g.slope, g.intercept)
		if sse < best-epsilon {
			t.Errorf("SumSquaredError(%v, %v) = %v, less than least-squares SSE %v",
				g.slope, g.intercept, sse, best)
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("linear-regression")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
