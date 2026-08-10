// Package sigmoidsoftmax visualizes how raw model scores ("logits" — any real
// number, unbounded) get squashed into probabilities. Sigmoid handles the
// two-class case: one curve, one number in (0, 1). Softmax generalizes it to
// any number of classes at once, turning a list of logits into a probability
// distribution that sums to 1 — and sigmoid turns out to be exactly the
// two-class special case of softmax, just written differently.
package sigmoidsoftmax

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "sigmoid-softmax",
		Title: "Sigmoid & softmax",
		Blurb: "A model doesn't output probabilities directly — it outputs 'logits', raw " +
			"scores that can be any real number, positive or negative, with no built-in upper " +
			"limit. Something has to turn '2.3 for cat, -0.5 for dog' into numbers that behave " +
			"like probabilities: between 0 and 1, summing to 1 across every option. For two " +
			"classes, sigmoid does it: sigmoid(z) = 1/(1+e^-z) squashes any logit into (0,1), " +
			"crossing 0.5 exactly at z=0, saturating toward 0 or 1 out at the extremes. For " +
			"three classes, say logits (2, 0.5, -1) for cat/dog/fox, softmax exponentiates " +
			"each one (e^2≈7.39, e^0.5≈1.65, e^-1≈0.37) and divides by their sum (≈9.41): " +
			"probabilities ≈78%, 18%, 4% — every logit contributes, and they add to 100%. " +
			"Sigmoid is exactly softmax's two-class case: run softmax on (z, 0) and the first " +
			"slot works out to 1/(1+e^-z), the same formula. The temperature knob divides every " +
			"logit before exponentiating: below 1 it sharpens the distribution toward a single " +
			"winner, above 1 it flattens the probabilities toward uniform even when the logits " +
			"disagree.",
		Params: []concept.ParamSpec{
			{Key: "z", Label: "Logit A (sigmoid & softmax)", Min: -6, Max: 6, Step: 0.1, Def: 1},
			{Key: "logitB", Label: "Logit B (softmax only)", Min: -6, Max: 6, Step: 0.1, Def: -1},
			{Key: "temp", Label: "Temperature", Min: 0.2, Max: 3, Step: 0.1, Def: 1},
		},
		Render: render,
	})
}

// Sigmoid squashes any real-valued logit z into (0, 1): 1/(1+e^-z).
func Sigmoid(z float64) float64 {
	return 1 / (1 + math.Exp(-z))
}

// Softmax turns a list of logits into a probability distribution: each entry
// is exponentiated and divided by the sum of all the exponentials, so the
// result is entrywise positive and sums to 1. temp (temperature) divides
// every logit before exponentiating: temp<1 sharpens the distribution toward
// whichever logit is largest, temp>1 flattens it toward uniform. Logits are
// shifted by their max before exponentiating — a standard numerical-stability
// trick that doesn't change the result (softmax is shift-invariant) but keeps
// math.Exp from overflowing on large logits.
func Softmax(logits []float64, temp float64) []float64 {
	out := make([]float64, len(logits))
	if len(logits) == 0 {
		return out
	}
	if temp <= 0 {
		temp = 1e-9
	}
	max := logits[0]
	for _, z := range logits[1:] {
		if z > max {
			max = z
		}
	}
	var sum float64
	for i, z := range logits {
		e := math.Exp((z - max) / temp)
		out[i] = e
		sum += e
	}
	if sum <= 0 {
		return out
	}
	for i := range out {
		out[i] /= sum
	}
	return out
}

func render(p map[string]float64) string {
	z, logitB, temp := p["z"], p["logitB"], p["temp"]
	if temp < 0.05 {
		temp = 0.05
	}

	c := viz.New(680, 520, -6, 6, 0, 1)
	// Push the curve down from the top (room for the header text) and
	// compress it away from the bottom (room for the softmax bars below,
	// drawn in raw pixel coordinates so they're unaffected by this).
	c.PadT = 60
	c.PadB = 250
	c.Axes()
	for x := -6.0; x <= 6.0; x += 2 {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	curve := viz.Sample(-6, 6, 240, Sigmoid)
	c.Path(curve, viz.Accent, 2.5)
	c.VLine(z, viz.Warm, true)

	zpx, zpy := c.X(z), c.Y(Sigmoid(z))
	c.Rect(zpx-4, zpy-4, 8, 8, viz.Warm, 0.9)

	c.Text(20, 24, fmt.Sprintf("sigmoid(z) = 1/(1+e^-z)    z = %.1f    sigmoid(z) = %.3f", z, Sigmoid(z)),
		14, viz.Ink, "start")
	c.Text(20, 44, "crosses 0.5 at z=0, saturates toward 0 or 1 at the extremes",
		12, viz.Muted, "start")

	// Below the curve: softmax over three classes. Class A shares the
	// sigmoid's logit z, class C is pinned at 0 as a reference, so this
	// panel and the curve above are the same z telling two related stories.
	probs := Softmax([]float64{z, logitB, 0}, temp)
	labels := []string{"A", "B", "C"}
	logits := []float64{z, logitB, 0}
	colors := []string{viz.Warm, viz.Accent, viz.Muted}

	const barBase, barMaxH, barW, gap = 460.0, 100.0, 110.0, 60.0
	const firstX = 130.0
	c.Text(20, 312, fmt.Sprintf("softmax over 3 logits (A=%.1f, B=%.1f, C=0), temperature = %.1f", z, logitB, temp),
		13, viz.Ink, "start")
	c.Text(20, 330, "each bar's height is that class's probability — all three add to 100%",
		12, viz.Muted, "start")

	for i, label := range labels {
		x := firstX + float64(i)*(barW+gap)
		h := probs[i] * barMaxH
		c.Rect(x, barBase-h, barW, h, colors[i], 0.75)
		c.Text(x+barW/2, barBase-h-10, fmt.Sprintf("%.0f%%", probs[i]*100), 13, viz.Ink, "middle")
		c.Text(x+barW/2, barBase+20, fmt.Sprintf("%s  (logit=%.1f)", label, logits[i]), 12, viz.Muted, "middle")
	}

	return c.String()
}
