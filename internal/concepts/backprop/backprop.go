// Package backprop visualizes backpropagation: running the chain rule
// backward through a small network — one input, two hidden neurons, one
// output — to compute every weight's gradient from a single forward pass
// followed by a single backward pass, instead of re-deriving each gradient
// from scratch.
package backprop

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "backpropagation",
		Seq:   87,
		Title: "Backpropagation (the chain rule through a small network)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "x", Label: "Input (x)", Min: -2, Max: 2, Step: 0.1, Def: 1},
			{Key: "lr", Label: "Learning rate", Min: 0, Max: 1.5, Step: 0.1, Def: 0.5},
		},
		Render: render,
	})
}

// Target is the fixed label every worked example trains toward: the network
// should learn to output close to 1.0 for whatever x the slider is set to.
const Target = 1.0

// Weights holds the network's 7 learnable parameters: a 1-input, 2-hidden,
// 1-output network. (W1,B1) and (W2,B2) feed the two hidden neurons from the
// single input; (V1,V2,B3) combine the two hidden activations into the
// output.
type Weights struct {
	W1, B1     float64
	W2, B2     float64
	V1, V2, B3 float64
}

// FixedWeights is the one starting point every Section and the picture walk
// through — chosen asymmetric on purpose (W1 != -W2, V1 != V2) so the two
// hidden neurons' gradients come out visibly different, not coincidentally
// equal.
var FixedWeights = Weights{
	W1: 0.6, B1: 0.1,
	W2: -0.4, B2: 0.2,
	V1: 0.9, V2: -0.7, B3: -0.1,
}

// Sigmoid squashes any real number into (0, 1).
func Sigmoid(z float64) float64 {
	return 1 / (1 + math.Exp(-z))
}

// SigmoidPrime is the derivative of Sigmoid, expressed in terms of the
// sigmoid's own output h = Sigmoid(z): h·(1-h). Backprop always has h
// already on hand from the forward pass, so this is the form actually used.
func SigmoidPrime(h float64) float64 {
	return h * (1 - h)
}

// Activations holds every intermediate value computed on the way from an
// input x to the loss — everything the forward pass produces, and every
// value the backward pass will need.
type Activations struct {
	Z1, H1 float64 // first hidden neuron: pre-activation, activation
	Z2, H2 float64 // second hidden neuron
	Z3, Y  float64 // output neuron
	Loss   float64 // 0.5*(y-target)^2
}

// Forward runs the network on input x and returns every intermediate value.
func Forward(w Weights, x, target float64) Activations {
	z1 := w.W1*x + w.B1
	h1 := Sigmoid(z1)
	z2 := w.W2*x + w.B2
	h2 := Sigmoid(z2)
	z3 := w.V1*h1 + w.V2*h2 + w.B3
	y := Sigmoid(z3)
	loss := 0.5 * (y - target) * (y - target)
	return Activations{Z1: z1, H1: h1, Z2: z2, H2: h2, Z3: z3, Y: y, Loss: loss}
}

// Gradients holds ∂Loss/∂(each weight), plus the three "local error" terms
// (Delta1, Delta2, Delta3) backpropagation computes on the way there. Every
// weight's gradient below is one of these deltas times whatever fed into
// that weight — the chain rule, applied once per layer, backward.
type Gradients struct {
	DW1, DB1               float64
	DW2, DB2               float64
	DV1, DV2, DB3          float64
	Delta1, Delta2, Delta3 float64
}

// Backward computes every weight's gradient in one backward pass over the
// Activations Forward already produced: start from ∂Loss/∂y at the output
// and multiply back through each layer's local derivative, reusing each
// delta for every weight that feeds into it instead of recomputing it.
func Backward(w Weights, a Activations, x, target float64) Gradients {
	// Output layer: how loss changes with the output, and with z3.
	dLdy := a.Y - target
	delta3 := dLdy * SigmoidPrime(a.Y)
	dv1 := delta3 * a.H1
	dv2 := delta3 * a.H2
	db3 := delta3

	// Hidden neuron 1: the chain rule step that gives backprop its name --
	// delta3 (already computed) times the weight it flows back through
	// (V1), times this neuron's own local derivative.
	dLdh1 := delta3 * w.V1
	delta1 := dLdh1 * SigmoidPrime(a.H1)
	dw1 := delta1 * x
	db1 := delta1

	// Hidden neuron 2: same pattern, flowing back through V2 instead.
	dLdh2 := delta3 * w.V2
	delta2 := dLdh2 * SigmoidPrime(a.H2)
	dw2 := delta2 * x
	db2 := delta2

	return Gradients{
		DW1: dw1, DB1: db1,
		DW2: dw2, DB2: db2,
		DV1: dv1, DV2: dv2, DB3: db3,
		Delta1: delta1, Delta2: delta2, Delta3: delta3,
	}
}

// GradientDescentStep applies one step of gradient descent (the update rule
// from `gradient-descent`) to every weight at once, using the gradients a
// single backward pass already computed: new weight = old weight - lr *
// gradient.
func GradientDescentStep(w Weights, g Gradients, lr float64) Weights {
	return Weights{
		W1: w.W1 - lr*g.DW1, B1: w.B1 - lr*g.DB1,
		W2: w.W2 - lr*g.DW2, B2: w.B2 - lr*g.DB2,
		V1: w.V1 - lr*g.DV1, V2: w.V2 - lr*g.DV2, B3: w.B3 - lr*g.DB3,
	}
}

func render(p map[string]float64) string {
	c := viz.New(700, 560, 0, 1, 0, 1)
	return c.String()
}
