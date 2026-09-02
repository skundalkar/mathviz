package attention

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) < tol }

func TestDotKnownValue(t *testing.T) {
	if got := Dot(0, 2, 0, 1); got != 2 {
		t.Errorf("Dot(0,2,0,1) = %v, want 2", got)
	}
}

func TestScaledScoresDefaultQuery(t *testing.T) {
	// q=(0,2), worked by hand in LEARNINGS.md: scores 0, 1.414, -0.707.
	got := ScaledScores(0, 2)
	want := []float64{0, 1.41421356, -0.70710678}
	for i := range want {
		if !near(got[i], want[i], 1e-6) {
			t.Errorf("ScaledScores(0,2)[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestSoftmaxWeightsSumToOne(t *testing.T) {
	for _, temp := range []float64{0.3, 1, 3} {
		w := Softmax(ScaledScores(0, 2), temp)
		var sum float64
		for _, wi := range w {
			sum += wi
		}
		if !near(sum, 1, 1e-9) {
			t.Errorf("Softmax weights at temp=%v sum to %v, want 1", temp, sum)
		}
	}
}

func TestOutputDefaultQueryMostlyAnimal(t *testing.T) {
	// q=(0,2) worked by hand: weights ~17.8%, 73.4%, 8.8%, output ~0.646.
	weights := Softmax(ScaledScores(0, 2), 1)
	if !near(weights[1], 0.7337, 1e-3) {
		t.Errorf("weight on 'animal' = %v, want ~0.7337", weights[1])
	}
	if got := Output(weights); !near(got, 0.6457, 1e-3) {
		t.Errorf("Output(weights) = %v, want ~0.6457", got)
	}
}

func TestOutputOppositeQueryMostlyStreet(t *testing.T) {
	weights := Softmax(ScaledScores(0, -2), 1)
	if got := Output(weights); got >= 0 {
		t.Errorf("Output for q=(0,-2) = %v, want negative (leaning toward 'street')", got)
	}
}

func TestOutputAmbiguousQueryNearZero(t *testing.T) {
	weights := Softmax(ScaledScores(0, 0), 1)
	if got := Output(weights); !near(got, 0, 0.05) {
		t.Errorf("Output for q=(0,0) = %v, want near 0 (ambiguous)", got)
	}
	for _, w := range weights {
		if !near(w, 1.0/3, 0.01) {
			t.Errorf("weights for q=(0,0) = %v, want all close to 1/3", weights)
		}
	}
}

func TestLowTemperatureSharpensTowardTopScore(t *testing.T) {
	scores := ScaledScores(0, 2) // 'animal' has the top score
	sharp := Softmax(scores, 0.3)
	flat := Softmax(scores, 3)
	if sharp[1] <= flat[1] {
		t.Errorf("low-temp weight on top token (%v) should exceed high-temp weight (%v)", sharp[1], flat[1])
	}
	if sharp[1] < 0.9 {
		t.Errorf("low-temp (0.3) weight on top token = %v, want close to 1", sharp[1])
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("attention-mechanism")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
