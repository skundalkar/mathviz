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
					"You've trained a model: split your data, ran training, and now you have a " +
						"handful of numbers — training loss, validation loss, ROC-AUC, PR-AUC, " +
						"precision, recall. Each one, read alone, tells you something narrow. Is the " +
						"model actually good? Do you need more data, a bigger model, a different " +
						"threshold, or is it fine to ship right now? None of these numbers answers " +
						"that by itself — you have to know which combination of them means what.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Five steps, in order, because a later step's numbers can't be trusted until " +
						"the earlier ones check out.",
					"• Step 1 — loss curves (train vs. validation, over training):",
					"|Pattern|Diagnosis|Action|",
					"|Both high, both flat|Underfitting — model/features can't capture the pattern|Bigger/more expressive model, better features, train longer|",
					"|Train drops, validation stalls or rises|Overfitting — memorizing train-set noise|Regularize, simplify, more data, early stopping (the overfitting concept)|",
					"|Both drop, converge close together|Healthy fit|Move to step 2|",
					"|Both noisy, not really decreasing|Optimization problem, not a data/capacity problem|Lower the learning rate, check init/bugs (the gradient-descent concept — a learning rate past the critical value looks exactly like this)|",
					"Only leave step 1 once train and validation loss both look reasonable and close together.",
					"• Step 2 — threshold-independent separation: ROC-AUC and PR-AUC. Regardless of " +
						"where you eventually set the threshold, can this model tell the classes apart " +
						"at all?",
					"|Reading|Diagnosis|Action|",
					"|High ROC-AUC and high PR-AUC|Model genuinely separates the classes|Move to step 3 — this is a threshold-picking problem now, not a model problem|",
					"|Low AUC (ROC ≈ 0.5, PR ≈ class prevalence)|No threshold will save this — the ranking itself is bad|Go back to the model/data, not the dial (step 4)|",
					"|High ROC-AUC, disappointing PR-AUC|Model separates fine, but you have the FPR-dilution imbalance-trap problem|Don't touch the model — this is exactly what PR-AUC is built to catch; go to step 3 knowing precision is the metric to watch closely|",
					"• Step 3 — if separation is good, pick the operating threshold. Sweep it " +
						"(precision-recall and confusion-matrix's job), read off precision and recall " +
						"at a few candidates, and pick based on which error costs more in your " +
						"situation — a missed case vs. a false alarm are not equally bad, and no " +
						"formula picks that for you. A free, instant move — no retraining — always " +
						"worth trying first once AUC is already good.",
					"• Step 4 — if separation itself is bad (or good-enough-but-not-there), " +
						"diagnose why. Precision, recall, F1 and AUC alone can't tell you the cause, " +
						"only that there's a ceiling — a few extra, cheap moves, cheapest first:",
					"|Diagnostic|What it tells you|",
					"|Compare train-set metric vs. validation-set metric (not just loss — AUC/precision itself)|Train ≫ validation → overfitting (back to step 1's fixes). Both mediocre → genuine signal shortage, keep going down this table|",
					"|Manual error review — pull 20-30 false positives/negatives, read them by hand|A lot are actually mislabeled → noisy labels, clean them (often the cheapest fix). Correctly labeled but genuinely ambiguous → a features/capacity problem|",
					"|Learning curve — retrain at 25%/50%/100% of your data, plot validation performance vs. data size|Still climbing at 100% → yes, more data will help. Flat/plateaued → more of the same data won't move it|",
					"|Feature ablation — add a feature you suspect carries signal, retrain|Meaningful jump → features were the bottleneck. No change → keep looking (try a more expressive model on the same features)|",
					"• Step 5 — touch the test set exactly once. Everything above uses train/" +
						"validation, potentially many times, as you iterate. The test set gets " +
						"evaluated once, right before shipping, for an honest final number. If test " +
						"performance is notably worse than validation, that's itself a signal: " +
						"repeated tuning quietly overfit to the validation set — the fix is a fresh " +
						"validation split going forward, not more tuning against the one you already " +
						"burned.",
					"Four situations from steps 1-2 come up constantly enough to recognize on " +
						"sight — the table below is exactly that, with real numbers:",
					"• Healthy fit — train loss 0.30, val loss 0.33 (close together); ROC-AUC 0.93, " +
						"PR-AUC 0.89 (both high); precision 91%, recall 90% (balanced). Nothing here " +
						"contradicts anything else.",
					"• Overfitting — train loss 0.06, val loss 0.58: a wide, growing gap. Every " +
						"number below loss (ROC-AUC 0.81, PR-AUC 0.74) is measuring an " +
						"over-memorized model, not a trustworthy one.",
					"• Underfitting — train loss 0.61, val loss 0.63: both high AND close together. " +
						"Not a generalization gap — a capacity gap. The model can't even fit its own " +
						"training data.",
					"• Imbalance trap — train loss 0.29, val loss 0.31 (healthy); ROC-AUC 0.95 " +
						"(looks great) — but PR-AUC 0.42 and precision at the chosen threshold, only " +
						"26%, reveal the model floods you with false alarms relative to how rare " +
						"positives are, exactly the FPR-dilution pattern from the roc-auc and pr-auc " +
						"concepts, something ROC-AUC alone hides.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"A fixed reference table, not something to turn a knob on — every row is a " +
						"pattern worth recognizing on sight. The left color tag is the fast read: " +
						"green means ship it, red means something in the fit itself is broken and " +
						"needs retraining, orange means the fit is fine and the problem is " +
						"downstream (the threshold, or the data's imbalance) — a different, usually " +
						"cheaper, fix. Read down a column across all four rows, not just across one " +
						"row, to see which single number actually distinguishes each situation from " +
						"'Healthy fit.'",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can match your own numbers against these four patterns and go straight to " +
						"the action instead of guessing: a wide, growing loss gap means retrain with " +
						"regularization; both losses high and close together means retrain with more " +
						"capacity or better features; a high-ROC-AUC/low-PR-AUC combination means " +
						"don't retrain at all — just move the threshold, or go get more of the rare " +
						"class.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"This is the actual sequence behind 'is my model ready to ship' for any binary " +
						"classifier — fraud detection, spam filtering, medical screening, defect " +
						"detection. Teams that skip straight to 'accuracy is 95%, ship it' are the " +
						"ones who get blindsided by the Imbalance trap row once it's in production.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Judging a model by one metric in isolation. 'ROC-AUC is 0.95, this is a great " +
						"model' is exactly the Imbalance trap row's setup, and it's wrong specifically " +
						"because it never checked precision. The fix isn't finding a better single " +
						"metric — there isn't one. It's checking the right numbers in the right " +
						"order, and knowing what a mismatch between them means.",
					"One more question this table can't answer even for the 'Healthy fit' row: is " +
						"the model's stated confidence trustworthy — is '90% confident' actually " +
						"right 90% of the time? A model can rank perfectly and still be badly " +
						"miscalibrated. See the calibration concept.",
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
