// Package evalplaybook answers the question every other classifier concept
// in this gallery leaves implicit: once you've trained a model, in what
// order do you actually look at the numbers, and what does each pattern of
// numbers tell you to do next? Loss, ROC-AUC, PR-AUC, precision and recall
// each isolate one thing — this concept is about reading them together.
// Unlike the rest of the gallery, this one has no sliders: it's a fixed
// reference table of the handful of patterns that come up constantly, each
// with its own diagnosis and action, not a single formula to explore.
package evalplaybook

import (
	"fmt"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "eval-playbook",
		Seq:   20,
		Title: "Model evaluation playbook",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder.",
				},
			},
		},
		Params: []concept.ParamSpec{},
		Render: render,
	})
}

// Scenario is one canonical training-diagnosis pattern: the numbers you'd
// see together at the same time, in practice, plus the diagnosis those
// numbers point to and the concrete next action. TrainLoss/ValLoss are
// arbitrary loss units (lower is better, no fixed scale); ROCAUC, PRAUC,
// Precision and Recall are all in [0,1], with Precision and Recall read at
// whatever single operating threshold the model has been set to.
type Scenario struct {
	Name      string
	TrainLoss float64
	ValLoss   float64
	ROCAUC    float64
	PRAUC     float64
	Precision float64
	Recall    float64
	Diagnosis string
	Action    string
}

// LossGap is validation loss minus training loss — the single number that,
// once both losses have otherwise settled to a low value, separates a
// healthy fit (gap near zero) from overfitting (gap growing while training
// loss keeps dropping).
func LossGap(trainLoss, valLoss float64) float64 {
	return valLoss - trainLoss
}

// IsOverfitting flags overfitting specifically: training loss already low
// (below lowThreshold) but validation loss trailing it by more than
// tolerance. This deliberately does NOT fire when both losses are high and
// close together — that pattern is underfitting, a different diagnosis
// with the opposite fix (more capacity, not more regularization).
func IsOverfitting(trainLoss, valLoss, lowThreshold, tolerance float64) bool {
	return trainLoss < lowThreshold && LossGap(trainLoss, valLoss) > tolerance
}

// F1 is the harmonic mean of precision and recall — 0 if both are 0, to
// avoid dividing by zero.
func F1(precision, recall float64) float64 {
	if precision+recall <= 0 {
		return 0
	}
	return 2 * precision * recall / (precision + recall)
}

// Scenarios returns the four fixed reference patterns this concept teaches,
// in the order they're worth checking: is the fit healthy, is it
// overfitting, is it underfitting, and — only once loss and AUC both look
// fine — is precision hiding an imbalance problem AUC alone would miss.
// Pure, deterministic data: same call, same four rows, every time.
func Scenarios() []Scenario {
	return []Scenario{
		{
			Name:      "Healthy fit",
			TrainLoss: 0.30, ValLoss: 0.33,
			ROCAUC: 0.93, PRAUC: 0.89,
			Precision: 0.91, Recall: 0.90,
			Diagnosis: "Train and validation loss converged close together; both ROC-AUC " +
				"and PR-AUC are high; precision and recall are balanced at the chosen " +
				"threshold. Nothing here is contradicting anything else.",
			Action: "Ship it. Pick the exact operating threshold based on which error " +
				"(false alarm vs. missed case) costs more in your situation — that's a " +
				"business call, not a modeling one.",
		},
		{
			Name:      "Overfitting",
			TrainLoss: 0.06, ValLoss: 0.58,
			ROCAUC: 0.81, PRAUC: 0.74,
			Precision: 0.78, Recall: 0.70,
			Diagnosis: "Training loss is excellent, validation loss is not — a wide, " +
				"growing gap. The model has started memorizing training-set specifics " +
				"instead of the general pattern; every number below loss is measuring " +
				"that memorized version, not a trustworthy one.",
			Action: "Regularize, simplify the model, add data (especially more of what " +
				"it's overfitting to), or stop training earlier. Don't act on precision/ " +
				"recall yet — retrain first, then re-read this table.",
		},
		{
			Name:      "Underfitting",
			TrainLoss: 0.61, ValLoss: 0.63,
			ROCAUC: 0.68, PRAUC: 0.55,
			Precision: 0.60, Recall: 0.58,
			Diagnosis: "Training loss is high too, not just validation loss — and the two " +
				"are close together. This is not a generalization gap, it's a capacity " +
				"or signal gap: the model can't fit even the data it's training on.",
			Action: "A bigger/more expressive model or better features. More data alone " +
				"rarely fixes this — a model too weak to fit 1,000 examples usually can't " +
				"fit 100,000 of the same kind either.",
		},
		{
			Name:      "Imbalance trap",
			TrainLoss: 0.29, ValLoss: 0.31,
			ROCAUC: 0.95, PRAUC: 0.42,
			Precision: 0.26, Recall: 0.90,
			Diagnosis: "Loss looks healthy and ROC-AUC looks great — but PR-AUC is " +
				"mediocre and precision at the operating threshold has collapsed. This " +
				"is the false-positive-rate-dilution pattern from the roc-auc/pr-auc " +
				"concepts: with rare positives, ROC-AUC can stay high while precision " +
				"quietly falls apart, because FPR divides false alarms by a huge pool of " +
				"real negatives and precision doesn't get that same cover.",
			Action: "Don't retrain — this isn't a capacity problem, loss and ROC-AUC both " +
				"say the model is fine. Move the threshold, or if PR-AUC itself is the " +
				"ceiling, get more positive-class examples or better features. Never " +
				"judge an imbalanced problem by ROC-AUC alone.",
		},
	}
}

// rowColor picks the at-a-glance category color for a scenario: green for a
// genuinely healthy fit, red for a real fit problem that needs retraining,
// orange for "the fit itself is fine, something downstream of it isn't."
func rowColor(name string) string {
	switch name {
	case "Healthy fit":
		return viz.Good
	case "Imbalance trap":
		return viz.Warm
	default:
		return viz.Bad
	}
}

// render draws the whole thing as a static table — no Params read at all,
// since every render() call produces the identical picture. This concept is
// a fixed reference, not something to turn a knob on.
func render(p map[string]float64) string {
	_ = p
	scenarios := Scenarios()

	c := viz.New(680, 420, 0, 1, 0, 1)

	c.Text(20, 24, "Same numbers, four situations — read down each column, not just across one row",
		13, viz.Ink, "start")
	c.Text(20, 44, "loss = arbitrary units, lower is better    other columns = a fraction in [0,1]",
		12, viz.Muted, "start")

	const startX, startY, rowH = 20.0, 70.0, 56.0
	colW := []float64{150, 70, 70, 70, 70, 70, 70, 70}
	headers := []string{"Scenario", "Train loss", "Val loss", "ROC-AUC", "PR-AUC", "Precision", "Recall", "F1"}

	colX := make([]float64, len(colW))
	x := startX
	for i, w := range colW {
		colX[i] = x
		x += w
	}

	for i, h := range headers {
		align, tx := "middle", colX[i]+colW[i]/2
		if i == 0 {
			align, tx = "start", colX[i]+4
		}
		c.Text(tx, startY-8, h, 12, viz.Muted, align)
	}
	c.Rect(startX, startY, x-startX, 1.5, viz.Muted, 0.6)

	for r, s := range scenarios {
		y := startY + float64(r)*rowH
		color := rowColor(s.Name)

		c.Rect(startX, y, x-startX, rowH-6, color, 0.06) // tinted row background
		c.Rect(startX, y, 4, rowH-6, color, 0.9)         // category tag strip

		vals := []string{
			s.Name,
			fmt.Sprintf("%.2f", s.TrainLoss),
			fmt.Sprintf("%.2f", s.ValLoss),
			fmt.Sprintf("%.2f", s.ROCAUC),
			fmt.Sprintf("%.2f", s.PRAUC),
			fmt.Sprintf("%.0f%%", s.Precision*100),
			fmt.Sprintf("%.0f%%", s.Recall*100),
			fmt.Sprintf("%.2f", F1(s.Precision, s.Recall)),
		}
		for i, v := range vals {
			align, tx, ink := "middle", colX[i]+colW[i]/2, viz.Ink
			if i == 0 {
				align, tx, ink = "start", colX[i]+12, color
			}
			c.Text(tx, y+rowH/2-2, v, 13, ink, align)
		}
	}

	legendY := startY + float64(len(scenarios))*rowH + 22
	c.Text(20, legendY, "green = healthy  •  red = a real fit problem, retrain  •  orange = fit is "+
		"fine, threshold or data isn't", 12, viz.Muted, "start")

	return c.String()
}
