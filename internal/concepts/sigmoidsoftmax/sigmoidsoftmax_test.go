package sigmoidsoftmax

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestSigmoidAtZero(t *testing.T) {
	if got := Sigmoid(0); math.Abs(got-0.5) > 1e-9 {
		t.Errorf("Sigmoid(0) = %v, want 0.5", got)
	}
}

func TestSigmoidBounded(t *testing.T) {
	for _, z := range []float64{-20, -3, -1, 0, 1, 3, 20} {
		s := Sigmoid(z)
		if s <= 0 || s >= 1 {
			t.Errorf("Sigmoid(%v) = %v, want strictly inside (0,1)", z, s)
		}
	}
}

func TestSigmoidMonotonicIncreasing(t *testing.T) {
	last := Sigmoid(-6)
	for z := -5.9; z <= 6; z += 0.1 {
		s := Sigmoid(z)
		if s < last {
			t.Fatalf("Sigmoid not monotonic increasing at z=%.1f", z)
		}
		last = s
	}
}

func TestSigmoidSaturates(t *testing.T) {
	if s := Sigmoid(20); s < 0.9999 {
		t.Errorf("Sigmoid(20) = %v, want close to 1", s)
	}
	if s := Sigmoid(-20); s > 0.0001 {
		t.Errorf("Sigmoid(-20) = %v, want close to 0", s)
	}
}

func TestSoftmaxSumsToOne(t *testing.T) {
	for _, logits := range [][]float64{
		{2, 0.5, -1},
		{0, 0, 0},
		{100, -100, 0},
		{-5},
	} {
		probs := Softmax(logits, 1)
		var sum float64
		for _, p := range probs {
			sum += p
		}
		if math.Abs(sum-1) > 1e-9 {
			t.Errorf("Softmax(%v) sums to %v, want 1", logits, sum)
		}
	}
}

func TestSoftmaxAllPositive(t *testing.T) {
	probs := Softmax([]float64{5, -5, 0, 3}, 1)
	for i, p := range probs {
		if p <= 0 || p >= 1 {
			t.Errorf("Softmax probability[%d] = %v, want strictly inside (0,1)", i, p)
		}
	}
}

func TestSoftmaxKnownValues(t *testing.T) {
	// cat/dog/fox example from the blurb: logits (2, 0.5, -1).
	probs := Softmax([]float64{2, 0.5, -1}, 1)
	want := []float64{0.7847, 0.1746, 0.0407}
	for i, w := range want {
		if math.Abs(probs[i]-w) > 5e-3 {
			t.Errorf("Softmax probability[%d] = %.4f, want ≈%.4f", i, probs[i], w)
		}
	}
}

func TestSoftmaxTwoClassMatchesSigmoid(t *testing.T) {
	// Sigmoid is the two-class special case of softmax: softmax([z, 0])[0]
	// should equal Sigmoid(z) for any z.
	for _, z := range []float64{-4, -1, 0, 1.5, 3, 5} {
		probs := Softmax([]float64{z, 0}, 1)
		want := Sigmoid(z)
		if math.Abs(probs[0]-want) > 1e-9 {
			t.Errorf("Softmax([%v,0])[0] = %v, want Sigmoid(%v) = %v", z, probs[0], z, want)
		}
	}
}

func TestSoftmaxTemperatureSharpensAndFlattens(t *testing.T) {
	logits := []float64{2, 0.5, -1}
	sharp := Softmax(logits, 0.3)
	base := Softmax(logits, 1)
	flat := Softmax(logits, 3)

	// Lower temperature should push the max probability higher (more
	// confident); higher temperature should pull it toward uniform.
	if sharp[0] <= base[0] {
		t.Errorf("lower temperature did not sharpen the top probability: sharp=%.4f base=%.4f", sharp[0], base[0])
	}
	if flat[0] >= base[0] {
		t.Errorf("higher temperature did not flatten the top probability: flat=%.4f base=%.4f", flat[0], base[0])
	}
	uniform := 1.0 / float64(len(logits))
	if math.Abs(flat[0]-uniform) > math.Abs(base[0]-uniform) {
		t.Errorf("high-temperature distribution not closer to uniform than base")
	}
}

func TestSoftmaxShiftInvariant(t *testing.T) {
	// Adding the same constant to every logit shouldn't change the result.
	a := Softmax([]float64{2, 0.5, -1}, 1)
	b := Softmax([]float64{102, 100.5, 99}, 1)
	for i := range a {
		if math.Abs(a[i]-b[i]) > 1e-6 {
			t.Errorf("softmax not shift-invariant at index %d: %v vs %v", i, a[i], b[i])
		}
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("sigmoid-softmax")
	if !ok {
		t.Fatal("sigmoid-softmax not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
