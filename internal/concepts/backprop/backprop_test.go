package backprop

import (
	"math"
	"testing"
)

func approxEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestSigmoidBasics(t *testing.T) {
	if !approxEqual(Sigmoid(0), 0.5, 1e-9) {
		t.Errorf("Sigmoid(0) = %v, want 0.5", Sigmoid(0))
	}
	if Sigmoid(10) <= 0.999 {
		t.Errorf("Sigmoid(10) = %v, want close to 1", Sigmoid(10))
	}
	if Sigmoid(-10) >= 0.001 {
		t.Errorf("Sigmoid(-10) = %v, want close to 0", Sigmoid(-10))
	}
}

func TestSigmoidPrimeAtHalf(t *testing.T) {
	// The derivative of sigmoid peaks at h=0.5, where h*(1-h)=0.25.
	if got := SigmoidPrime(0.5); !approxEqual(got, 0.25, 1e-9) {
		t.Errorf("SigmoidPrime(0.5) = %v, want 0.25", got)
	}
}

// TestForwardKnownValues checks the worked example (x=1, FixedWeights,
// target=1) against hand-computed values.
func TestForwardKnownValues(t *testing.T) {
	a := Forward(FixedWeights, 1.0, Target)
	want := Activations{
		Z1: 0.7, H1: 0.668188,
		Z2: -0.2, H2: 0.450166,
		Z3: 0.186253, Y: 0.546429,
		Loss: 0.102863,
	}
	tol := 0.0005
	if !approxEqual(a.Z1, want.Z1, tol) || !approxEqual(a.H1, want.H1, tol) {
		t.Errorf("hidden 1: got Z1=%.6f H1=%.6f, want Z1=%.6f H1=%.6f", a.Z1, a.H1, want.Z1, want.H1)
	}
	if !approxEqual(a.Z2, want.Z2, tol) || !approxEqual(a.H2, want.H2, tol) {
		t.Errorf("hidden 2: got Z2=%.6f H2=%.6f, want Z2=%.6f H2=%.6f", a.Z2, a.H2, want.Z2, want.H2)
	}
	if !approxEqual(a.Z3, want.Z3, tol) || !approxEqual(a.Y, want.Y, tol) {
		t.Errorf("output: got Z3=%.6f Y=%.6f, want Z3=%.6f Y=%.6f", a.Z3, a.Y, want.Z3, want.Y)
	}
	if !approxEqual(a.Loss, want.Loss, tol) {
		t.Errorf("Loss = %.6f, want %.6f", a.Loss, want.Loss)
	}
}

// TestBackwardKnownValues checks the same worked example's gradients.
func TestBackwardKnownValues(t *testing.T) {
	x := 1.0
	a := Forward(FixedWeights, x, Target)
	g := Backward(FixedWeights, a, x, Target)

	tol := 0.0005
	checks := map[string]struct{ got, want float64 }{
		"delta3": {g.Delta3, -0.112415},
		"delta1": {g.Delta1, -0.022431},
		"delta2": {g.Delta2, 0.019477},
		"dv1":    {g.DV1, -0.075114},
		"dv2":    {g.DV2, -0.050605},
		"db3":    {g.DB3, -0.112415},
		"dw1":    {g.DW1, -0.022431},
		"db1":    {g.DB1, -0.022431},
		"dw2":    {g.DW2, 0.019477},
		"db2":    {g.DB2, 0.019477},
	}
	for name, c := range checks {
		if !approxEqual(c.got, c.want, tol) {
			t.Errorf("%s = %.6f, want %.6f", name, c.got, c.want)
		}
	}
}

// TestDeltasNotCoincidentallyEqual guards against a degenerate/symmetric
// choice of weights making the two hidden neurons' gradients identical,
// which would hide the fact backprop computes a genuinely separate
// gradient per weight.
func TestDeltasNotCoincidentallyEqual(t *testing.T) {
	a := Forward(FixedWeights, 1.0, Target)
	g := Backward(FixedWeights, a, 1.0, Target)
	if approxEqual(g.Delta1, g.Delta2, 1e-6) {
		t.Errorf("delta1 (%.6f) and delta2 (%.6f) came out equal — weights should be asymmetric enough to avoid this", g.Delta1, g.Delta2)
	}
}

// TestGradientDescentStepReducesLoss checks that stepping every weight
// opposite its gradient (gradient-descent's update rule) actually lowers
// the loss for the worked example, at a couple of learning rates.
func TestGradientDescentStepReducesLoss(t *testing.T) {
	x := 1.0
	a := Forward(FixedWeights, x, Target)
	g := Backward(FixedWeights, a, x, Target)

	for _, lr := range []float64{0.1, 0.5, 1.0} {
		updated := GradientDescentStep(FixedWeights, g, lr)
		newLoss := Forward(updated, x, Target).Loss
		if newLoss >= a.Loss {
			t.Errorf("lr=%.2f: loss did not decrease (%.6f -> %.6f)", lr, a.Loss, newLoss)
		}
	}
}

// TestGradientDescentStepZeroLR is a no-op: lr=0 must leave every weight
// unchanged.
func TestGradientDescentStepZeroLR(t *testing.T) {
	a := Forward(FixedWeights, 1.0, Target)
	g := Backward(FixedWeights, a, 1.0, Target)
	updated := GradientDescentStep(FixedWeights, g, 0)
	if updated != FixedWeights {
		t.Errorf("lr=0 changed the weights: got %+v, want %+v", updated, FixedWeights)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	svg := render(map[string]float64{"x": 1, "lr": 0.5})
	if len(svg) == 0 {
		t.Fatal("render returned empty string")
	}
	if svg[:4] != "<svg" {
		t.Errorf("render output doesn't start with <svg: %q...", svg[:20])
	}
}
