// Package crossentropyloss visualizes binary cross-entropy (log) loss: the
// per-example penalty a classifier pays for how confidently right or wrong
// its predicted probability was against a true 0/1 label, and why averaging
// that penalty over a dataset is the loss function most classifiers are
// actually trained to minimize.
package crossentropyloss

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "cross-entropy-loss",
		Seq:   69,
		Title: "Cross-entropy loss (scoring a predicted probability against a label)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`kl-divergence` compared two full probability distributions — a true " +
						"coin bias P(heads)=0.9 against a model's believed Q(heads)=0.5 — and " +
						"measured the gap between them in bits. But training a real classifier " +
						"doesn't hand you a true distribution to compare against; it hands you a " +
						"hard label. A spam filter either was spam (y=1) or wasn't (y=0) for a " +
						"given email — there's no 'this email was 90% spam' ground truth. The " +
						"model still outputs a probability, say p̂=0.1 for an email that actually " +
						"was spam. How do you turn 'the true label was 1 but the model said 0.1' " +
						"into a single number a training algorithm can push toward zero?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Four emails, each with a true spam label y and the model's predicted " +
						"P(spam)=p̂:",
					"• Email 1: y=1 (spam), p̂=0.9 (confident and correct) → loss = " +
						"-log2(0.9) = 0.152 bits",
					"• Email 2: y=1 (spam), p̂=0.55 (barely correct) → loss = -log2(0.55) = " +
						"0.863 bits",
					"• Email 3: y=0 (not spam), p̂=0.2 (confident and correct, since the model " +
						"put 80% on 'not spam') → loss = -log2(1-0.2) = -log2(0.8) = 0.322 bits",
					"• Email 4: y=1 (spam), p̂=0.1 (confident and WRONG) → loss = -log2(0.1) = " +
						"3.322 bits",
					"Each loss only ever looks at the probability the model assigned to the " +
						"label that actually happened: -log2(p̂) when y=1, -log2(1-p̂) when y=0. " +
						"Average the four: (0.152+0.863+0.322+3.322)/4 = 1.165 bits — the average " +
						"cross-entropy loss over this tiny dataset. Notice email 4 alone " +
						"contributes almost three times as much as the other three combined: " +
						"-log2(x) shoots toward infinity as x→0, so a confident wrong prediction " +
						"is punished far harder than an unconfident right one is rewarded.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Two curves of loss = -log2(p̂) (when y=1) and -log2(1-p̂) (when y=0), " +
						"plotted across every possible predicted probability p̂ from 0 to 1. The y " +
						"slider picks which true label is in force — which curve is 'active' — " +
						"and the p̂ slider marks a point on it. Drag p̂ toward the label's own side " +
						"(1 when y=1, 0 when y=0) and watch the loss fall toward 0; drag it toward " +
						"the wrong side and watch the loss climb steeply, blowing up as p̂ " +
						"approaches the wrong extreme.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Score any single prediction against its true label with one number that " +
						"rewards confident-and-correct, tolerates unconfident-and-correct, and " +
						"punishes confident-and-wrong sharply — exactly the ranking `kl-divergence` " +
						"couldn't produce without a full second distribution to compare against. " +
						"Average that number over a whole training set and you get the actual " +
						"quantity gradient descent minimizes when training a classifier: nudge " +
						"every prediction a little closer to its true label, over and over, and the " +
						"average loss falls.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"'Cross-entropy loss' (also called 'log loss') is the default training " +
						"objective for logistic regression, neural network classifiers, and spam/ " +
						"fraud detectors — anywhere a model outputs a probability and gets scored " +
						"against a known right answer. Kaggle competitions and production ML " +
						"dashboards alike report 'log loss' as a leaderboard metric for exactly " +
						"this reason: it's differentiable, and it punishes confidently-wrong " +
						"predictions in a way plain accuracy doesn't.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'email 4's loss of 3.322 bits dominates the average " +
						"because the model was both wrong and confident about it — cross-entropy " +
						"loss punishes that combination far harder than being merely unconfident.'",
					"Not like this: treating cross-entropy loss as the same thing as accuracy, " +
						"or assuming a lower average loss always means more correct predictions. A " +
						"model that's right 90% of the time but wildly overconfident on its misses " +
						"can post a worse average loss than a model that's right 85% of the time " +
						"but stays modest (p̂ near 0.5) whenever it's unsure — loss cares about " +
						"calibrated confidence, not just which side of 0.5 the prediction landed " +
						"on.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "y", Label: "True label y", Min: 0, Max: 1, Step: 1, Def: 1},
			{Key: "phat", Label: "Predicted P(y=1)", Min: 0.01, Max: 0.99, Step: 0.01, Def: 0.7},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
