// Package backprop visualizes backpropagation: running the chain rule
// backward through a small network — one input, two hidden neurons, one
// output — to compute every weight's gradient from a single forward pass
// followed by a single backward pass, instead of re-deriving each gradient
// from scratch.
package backprop

import (
	"fmt"
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
				Body: []string{
					"`chain-rule` showed you how to multiply derivatives through a chain of " +
						"two composed functions, and `gradient-descent` showed you how to nudge a " +
						"single number downhill once you have its derivative. A real network " +
						"isn't one composed function with one input and one derivative, though — " +
						"it's dozens (in a real network, millions) of weights, all feeding forward " +
						"through several layers into one final number, the loss — and gradient " +
						"descent needs a derivative for *every single weight* before it can take " +
						"even one step. Do you really have to re-derive each weight's gradient " +
						"completely separately, from scratch, one at a time?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take the smallest network that still has this problem: one input x, two " +
						"hidden neurons (each a weight, a bias, and a sigmoid squashing the " +
						"result into (0,1)), and one output neuron (same shape) whose output y is " +
						"compared against a target with squared-error loss L = ½(y−target)². " +
						"Fixed weights: w1=0.60,b1=0.10 (feeds hidden neuron 1), " +
						"w2=−0.40,b2=0.20 (feeds hidden neuron 2), v1=0.90,v2=−0.70,b3=−0.10 " +
						"(combines both hidden neurons into the output). With x=1.00 and " +
						"target=1.00, the forward pass is: z1=0.70 → h1=0.668, z2=−0.20 → " +
						"h2=0.450, z3 = 0.9(0.668)+(−0.7)(0.450)+(−0.1) = 0.19 → y=0.546, giving " +
						"Loss = ½(0.546−1.00)² = 0.1029.",
					"Now go backward, one layer at a time, computing each layer's \"local " +
						"error\" term (δ) exactly once and reusing it for every weight that " +
						"touches it. Start at the output: δ3 = (y−target)·y·(1−y) = " +
						"(−0.454)(0.546)(0.454) = −0.112. Every weight feeding *into* y uses δ3 " +
						"directly: ∂L/∂v1 = δ3·h1 = −0.112×0.668 = −0.075, ∂L/∂v2 = δ3·h2 = " +
						"−0.112×0.450 = −0.051, ∂L/∂b3 = δ3 = −0.112.",
					"• The chain-rule step that gives backpropagation its name: to get hidden " +
						"neuron 1's local error δ1, take δ3 (already computed — never " +
						"recomputed), multiply by the weight it flows back through (v1=0.90) to " +
						"get ∂L/∂h1 = −0.112×0.90 = −0.101, then multiply by h1's own local " +
						"derivative h1(1−h1) = 0.668×0.332 = 0.222: δ1 = −0.101×0.222 = −0.022. " +
						"That's the chain rule (∂L/∂h1 × ∂h1/∂z1), applied exactly the way " +
						"`chain-rule` multiplies rates through a composition — just one more " +
						"link added to the chain.",
					"• Hidden neuron 2 gets the identical treatment through v2 instead of v1: " +
						"∂L/∂h2 = δ3·v2 = −0.112×(−0.70) = 0.079, δ2 = 0.079×h2(1−h2) = " +
						"0.079×0.248 = 0.019.",
					"• Every remaining weight is now just δ times whatever fed it: ∂L/∂w1 = " +
						"δ1·x = −0.022×1.00 = −0.022, ∂L/∂b1 = δ1 = −0.022 (a bias's \"input\" " +
						"is always 1, so its gradient is always exactly its neuron's δ). Same " +
						"pattern for neuron 2: ∂L/∂w2 = δ2·x = 0.019, ∂L/∂b2 = δ2 = 0.019.",
					"Tally it up: 3 δ's computed (δ3, δ1, δ2), each computed once, produced all " +
						"7 gradients (∂L/∂w1, ∂L/∂b1, ∂L/∂w2, ∂L/∂b2, ∂L/∂v1, ∂L/∂v2, ∂L/∂b3) — " +
						"one forward pass followed by one backward pass, not seven separate " +
						"derivations from scratch.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The diagram is the network itself: x on the left, the two hidden neurons " +
						"h1/h2 in the middle, y on the right, arrows for every weight. Each arrow " +
						"is labeled with its weight (black) and its gradient (red for negative, " +
						"green for positive) directly below it. Each neuron box shows its " +
						"activation value (black) and its δ (orange) underneath — the exact same " +
						"numbers from the worked example above, recomputed live. The \"Input " +
						"(x)\" slider re-runs the whole forward-then-backward pass for a " +
						"different x, so you can watch every δ and every gradient shift together " +
						"as the input changes. The \"Learning rate\" slider feeds the computed " +
						"gradients into one `gradient-descent` step on every weight at once, and " +
						"the bottom line shows the loss actually dropping — the payoff for having " +
						"computed all 7 gradients in the first place.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Train an actual multi-layer network: repeat forward pass → backward pass " +
						"(every gradient, in one shot) → `gradient-descent` step on every weight → " +
						"repeat, which is the entire training loop behind every neural network, " +
						"no matter how many layers or weights it has. Without backpropagation's " +
						"δ-reuse trick, getting a gradient for every weight in a real network " +
						"(millions of them, not 7) by separately re-deriving each one would be " +
						"computationally hopeless; with it, the cost of the whole backward pass " +
						"is about the same as one forward pass.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Every deep learning framework's `.backward()` call (PyTorch, TensorFlow, " +
						"JAX) runs exactly this algorithm, automated for networks with millions " +
						"or billions of weights instead of this example's 7 — \"backprop\" is " +
						"quite literally what trains every image classifier, chatbot, and " +
						"recommendation system in production today. The everyday intuition: it's " +
						"like tracing back through a long chain of \"because of that, because of " +
						"that, because of that\" to find how much blame each early decision " +
						"deserves for a final outcome, computing each link's share of the blame " +
						"once and passing it backward, rather than re-tracing the entire chain " +
						"from scratch for every single decision you want to evaluate.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: \"backpropagation doesn't invent a new kind of " +
						"derivative — it applies the ordinary chain rule, but computes each " +
						"layer's local error (δ) exactly once and reuses it for every weight " +
						"that touches it, instead of expanding the full chain-rule product by " +
						"hand for each weight separately.\" Not like this: treating each weight's " +
						"gradient as its own unrelated computation — δ1 and δ2 above both " +
						"started from the same δ3, and a real network with millions of weights " +
						"leans on exactly that sharing to stay computationally feasible at all. " +
						"Also not like this: confusing the backward pass with training itself — " +
						"backward pass only tells you which direction, and how steeply, to nudge " +
						"each weight; actually improving the network still requires taking a " +
						"`gradient-descent` step with those gradients (Loss 0.1029 → 0.0919 " +
						"above, at lr=0.5), and in practice repeating that over many examples.",
				},
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

// Node positions in the 0..1 x 0..1 diagram space: input on the left, the
// two hidden neurons stacked in the middle, output on the right.
var (
	posX  = [2]float64{0.08, 0.50}
	posH1 = [2]float64{0.42, 0.80}
	posH2 = [2]float64{0.42, 0.20}
	posY  = [2]float64{0.80, 0.50}
)

// drawEdge draws a straight arrow from src to dst, stopping short of both
// node boxes so the line doesn't run under them, with a small arrowhead at
// the dst end.
func drawEdge(c *viz.Canvas, src, dst [2]float64) {
	dx, dy := dst[0]-src[0], dst[1]-src[1]
	length := math.Hypot(dx, dy)
	if length == 0 {
		return
	}
	ux, uy := dx/length, dy/length
	const margin = 0.075
	sx, sy := src[0]+ux*margin, src[1]+uy*margin
	ex, ey := dst[0]-ux*margin, dst[1]-uy*margin
	c.Path([][2]float64{{sx, sy}, {ex, ey}}, viz.Muted, 1.5)

	const wingLen, wingWidth = 0.03, 0.014
	perpX, perpY := -uy, ux
	w1x, w1y := ex-ux*wingLen+perpX*wingWidth, ey-uy*wingLen+perpY*wingWidth
	w2x, w2y := ex-ux*wingLen-perpX*wingWidth, ey-uy*wingLen-perpY*wingWidth
	c.Path([][2]float64{{w1x, w1y}, {ex, ey}, {w2x, w2y}}, viz.Muted, 1.5)
}

// edgeLabel writes a 2-line label (weight, then its gradient) near the
// midpoint of the src-dst edge, offset toward dy so it doesn't sit on top
// of the line itself.
func edgeLabel(c *viz.Canvas, src, dst [2]float64, dyPx float64, weightLine, gradLine string, gradColor string) {
	mx, my := (src[0]+dst[0])/2, (src[1]+dst[1])/2
	px, py := c.X(mx), c.Y(my)
	c.Text(px, py+dyPx, weightLine, 12, viz.Ink, "middle")
	c.Text(px, py+dyPx+15, gradLine, 12, gradColor, "middle")
}

// node draws one neuron's box, its label, and (if given) its activation
// value and local error term underneath.
func node(c *viz.Canvas, pos [2]float64, label, valueLine, deltaLine string) {
	px, py := c.X(pos[0]), c.Y(pos[1])
	const side = 44.0
	c.Rect(px-side/2, py-side/2, side, side, viz.Accent, 0.9)
	c.Text(px, py+5, label, 15, "white", "middle")
	if valueLine != "" {
		c.Text(px, py+side/2+16, valueLine, 12, viz.Ink, "middle")
	}
	if deltaLine != "" {
		c.Text(px, py+side/2+32, deltaLine, 12, viz.Warm, "middle")
	}
}

func gradColorFor(grad float64) string {
	if grad < 0 {
		return viz.Bad
	}
	return viz.Good
}

func render(p map[string]float64) string {
	x, lr := p["x"], p["lr"]
	w := FixedWeights

	a := Forward(w, x, Target)
	g := Backward(w, a, x, Target)
	updated := GradientDescentStep(w, g, lr)
	newLoss := Forward(updated, x, Target).Loss

	c := viz.New(700, 560, 0, 1, 0, 1)
	c.PadT = 116
	c.PadB = 130

	drawEdge(c, posX, posH1)
	drawEdge(c, posX, posH2)
	drawEdge(c, posH1, posY)
	drawEdge(c, posH2, posY)

	edgeLabel(c, posX, posH1, -14, fmt.Sprintf("w1=%.2f", w.W1), fmt.Sprintf("∂L/∂w1=%.3f", g.DW1), gradColorFor(g.DW1))
	edgeLabel(c, posX, posH2, 14, fmt.Sprintf("w2=%.2f", w.W2), fmt.Sprintf("∂L/∂w2=%.3f", g.DW2), gradColorFor(g.DW2))
	edgeLabel(c, posH1, posY, -14, fmt.Sprintf("v1=%.2f", w.V1), fmt.Sprintf("∂L/∂v1=%.3f", g.DV1), gradColorFor(g.DV1))
	edgeLabel(c, posH2, posY, 14, fmt.Sprintf("v2=%.2f", w.V2), fmt.Sprintf("∂L/∂v2=%.3f", g.DV2), gradColorFor(g.DV2))

	node(c, posX, "x", fmt.Sprintf("x=%.2f", x), "")
	node(c, posH1, "h1", fmt.Sprintf("h1=%.3f", a.H1), fmt.Sprintf("δ1=%.3f", g.Delta1))
	node(c, posH2, "h2", fmt.Sprintf("h2=%.3f", a.H2), fmt.Sprintf("δ2=%.3f", g.Delta2))
	node(c, posY, "y", fmt.Sprintf("y=%.3f", a.Y), fmt.Sprintf("δ3=%.3f", g.Delta3))

	c.Text(20, 22, fmt.Sprintf("forward: x=%.2f → z1=%.2f,h1=%.3f | z2=%.2f,h2=%.3f → z3=%.2f,y=%.3f", x, a.Z1, a.H1, a.Z2, a.H2, a.Z3, a.Y), 13, viz.Ink, "start")
	c.Text(20, 42, fmt.Sprintf("target=%.2f    Loss = ½(y−target)² = %.4f", Target, a.Loss), 13, viz.Muted, "start")
	c.Text(20, 64, "one backward pass computes every δ once; each weight's gradient = its δ × what fed it", 12, viz.Muted, "start")
	c.Text(20, 84, "white = neuron label   black = activation   orange = δ, the term backprop reuses", 12, viz.Muted, "start")

	c.Text(20, 490, fmt.Sprintf("gradients: w1=%.3f  b1=%.3f  w2=%.3f  b2=%.3f  v1=%.3f  v2=%.3f  b3=%.3f",
		g.DW1, g.DB1, g.DW2, g.DB2, g.DV1, g.DV2, g.DB3), 12, viz.Ink, "start")
	c.Text(20, 512, fmt.Sprintf("gradient-descent step (lr=%.2f): Loss %.4f → %.4f", lr, a.Loss, newLoss), 13, viz.Good, "start")
	c.Text(20, 534, "bias gradients equal δ1, δ2, δ3 exactly — a bias's \"input\" is always 1", 12, viz.Muted, "start")

	return c.String()
}
