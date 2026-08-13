// Package naivebayes visualizes the naive Bayes classifier: chaining
// several independent clues into one verdict by multiplying each clue's own
// likelihood ratio into a running posterior odds, the way bayes-theorem
// (elsewhere in this gallery) updates on a single clue, repeated per word.
package naivebayes

import (
	"fmt"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "naive-bayes",
		Seq:   31,
		Title: "Naive Bayes (chaining independent clues)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"The bayes-theorem concept elsewhere in this gallery updates one belief " +
						"using one piece of evidence — a single positive test result. A real spam " +
						"filter doesn't get just one clue: an email might contain 'free,' not " +
						"contain 'meeting,' contain 'winner,' and so on, dozens of words at once, " +
						"and you want to combine all of them into one verdict. Gut instinct: just " +
						"build one giant table of 'how often do spam emails contain exactly this " +
						"combination of words' and look up the answer directly. That falls apart " +
						"fast — with even 20 words to check, there are over a million possible " +
						"present/absent combinations, and almost none of them will show up often " +
						"enough in any real training set to estimate reliably. Is there a way to " +
						"combine many independent clues into one verdict without needing training " +
						"data for every possible combination of them?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Make one simplifying assumption: treat each word's presence as independent " +
						"of every other word's presence, once you already know whether the email " +
						"is spam or not (the 'naive' in naive Bayes). That assumption turns an " +
						"impossible combinatorial lookup into something you actually can estimate " +
						"— the odds-of-spam version of Bayes' theorem, updated once per word: " +
						"posterior odds = prior odds × (likelihood ratio for word 1) × (likelihood " +
						"ratio for word 2) × ..., where each word's likelihood ratio is P(word's " +
						"status | spam) / P(word's status | not spam) — exactly the same " +
						"odds-updating step bayes-theorem uses for one clue, just chained across " +
						"several.",
					"Work a small example. Start with a 50/50 prior (odds = 1:1), and two words " +
						"estimated from a training set: 'free' appears in 60% of spam but only 10% " +
						"of non-spam; 'meeting' appears in 10% of spam but 50% of non-spam. A new " +
						"email contains 'free' but not 'meeting':",
					"• likelihood ratio for seeing 'free': 0.60/0.10 = 6 — six times more likely " +
						"under spam than under not-spam.",
					"• likelihood ratio for not seeing 'meeting': P(no 'meeting'|spam)/P(no " +
						"'meeting'|not spam) = 0.90/0.50 = 1.8 — its absence is itself evidence, and " +
						"it also points (mildly) toward spam.",
					"• posterior odds = 1 × 6 × 1.8 = 10.8 : 1 in favor of spam → posterior " +
						"probability = 10.8/(10.8+1) ≈ 91.5%.",
					"Each word only ever needed its own two numbers (how often it shows up in " +
						"spam, how often in non-spam) — nothing about how it interacts with any " +
						"other word — and the combination still produced a confident, specific " +
						"verdict.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Three stacked bars, read top to bottom as evidence accumulates. The top bar " +
						"is the prior — P(spam) before checking any word. The middle bar folds in " +
						"whatever the 'free' toggle says (present or absent), the same way section " +
						"2's worked example does; the bottom bar folds in the 'meeting' toggle on " +
						"top of that. Each bar's fill length is that stage's running P(spam), and " +
						"it turns orange once it crosses 50% (the email would be classified as " +
						"spam at that point) or stays green below it. Flip either word toggle and " +
						"every bar from that word downward updates — flipping 'free' off, for " +
						"instance, removes its ×6 pull toward spam entirely, changing where the " +
						"final bar lands.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Combine as many independent clues as you have — not just two — into one " +
						"confident classification, using only per-clue statistics that are cheap " +
						"to estimate from a modest amount of training data, instead of needing a " +
						"combinatorially huge table covering every possible combination of clues " +
						"at once. This is exactly how a naive Bayes spam filter reaches a verdict " +
						"on a whole email: fold in every word's individual likelihood ratio, one at " +
						"a time, and read the final odds.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Early and still widely-used spam filters are literally this — the 'naive " +
						"Bayes spam filter.' Sentiment analysis (is a product review positive or " +
						"negative?) and basic document classification (which newsgroup or topic " +
						"does this article belong to?) use the same word-by-word odds-multiplying " +
						"approach. It's popular specifically because it's cheap to train, works " +
						"surprisingly well even though the independence assumption is rarely " +
						"exactly true in real language, and degrades gracefully rather than needing " +
						"exponentially more data as you add more clues.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'naive Bayes multiplies each clue's own likelihood ratio " +
						"into the running odds, assuming the clues don't influence each other once " +
						"you know the class' — the 'naive' part is exactly that independence " +
						"assumption, not a flaw in Bayes' theorem itself. Not like this: assuming " +
						"the independence assumption has to be true for the method to be useful at " +
						"all (in practice it's often somewhat wrong — 'free' and 'winner' really do " +
						"tend to show up together in spam — and naive Bayes still classifies well " +
						"regardless), or forgetting that an absent word is still evidence (a normal " +
						"work email is often more identifiable by which common words it's missing " +
						"than by which unusual ones it contains).",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "prior", Label: "Prior P(spam)", Min: 5, Max: 95, Step: 5, Def: 50, Unit: "%"},
			{Key: "freePresent", Label: "Contains 'free'?", Min: 0, Max: 1, Step: 1, Def: 1},
			{Key: "meetingPresent", Label: "Contains 'meeting'?", Min: 0, Max: 1, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// Fixed per-word statistics for the worked example every Section walks
// through, as if estimated from a training set of spam and non-spam email:
// how often each word appears in spam vs. in non-spam.
const (
	PFreeGivenSpam    = 0.6
	PFreeGivenHam     = 0.1
	PMeetingGivenSpam = 0.1
	PMeetingGivenHam  = 0.5
)

// Odds converts a probability into odds: p / (1-p). p must be in (0, 1).
func Odds(p float64) float64 {
	return p / (1 - p)
}

// ProbabilityFromOdds is Odds's inverse: odds / (1 + odds).
func ProbabilityFromOdds(odds float64) float64 {
	return odds / (1 + odds)
}

// WordLikelihoodRatio returns one word's contribution to the posterior
// odds: how much more likely its observed status (present or absent) is
// under spam than under not-spam. present=true uses the word's own
// probabilities directly; present=false uses their complements, since a
// word's absence is itself evidence.
func WordLikelihoodRatio(present bool, pGivenSpam, pGivenHam float64) float64 {
	if present {
		return pGivenSpam / pGivenHam
	}
	return (1 - pGivenSpam) / (1 - pGivenHam)
}

// PosteriorProbability folds a prior probability and any number of
// independent likelihood ratios into one posterior probability: convert
// the prior to odds, multiply in every ratio (the naive-Bayes "naive"
// step — treating each ratio's evidence as independent of the others),
// then convert back to a probability. Called with no ratios, it returns
// the prior unchanged.
func PosteriorProbability(prior float64, ratios ...float64) float64 {
	odds := Odds(prior)
	for _, r := range ratios {
		odds *= r
	}
	return ProbabilityFromOdds(odds)
}

// bar draws one row's running P(spam): a label above a horizontal bar
// filled proportionally to prob, colored orange once prob crosses 50% (the
// email would be classified spam at that point) or green below it, with a
// thin tick marking the 50% threshold itself. Returns the y to start the
// next row at.
func bar(c *viz.Canvas, y float64, label string, prob float64) float64 {
	const x0, w, h = 20.0, 460.0, 24.0

	color := viz.Good
	if prob >= 0.5 {
		color = viz.Warm
	}

	c.Text(x0, y, label, 13, viz.Ink, "start")
	barY := y + 8
	c.Rect(x0, barY, w, h, viz.Faint, 1)
	c.Rect(x0, barY, w*prob, h, color, 1)
	c.Rect(x0+w*0.5-0.5, barY-3, 1, h+6, viz.Muted, 1)
	c.Text(x0+w+10, barY+h/2+5, fmt.Sprintf("%.1f%%", prob*100), 13, viz.Ink, "start")

	return barY + h + 34
}

func render(p map[string]float64) string {
	prior := p["prior"] / 100
	freePresent := p["freePresent"] >= 0.5
	meetingPresent := p["meetingPresent"] >= 0.5

	ratioFree := WordLikelihoodRatio(freePresent, PFreeGivenSpam, PFreeGivenHam)
	ratioMeeting := WordLikelihoodRatio(meetingPresent, PMeetingGivenSpam, PMeetingGivenHam)
	afterFree := PosteriorProbability(prior, ratioFree)
	afterBoth := PosteriorProbability(prior, ratioFree, ratioMeeting)

	c := viz.New(640, 380, 0, 1, 0, 1)

	freeWord, meetingWord := "doesn't contain", "doesn't contain"
	if freePresent {
		freeWord = "contains"
	}
	if meetingPresent {
		meetingWord = "contains"
	}

	c.Text(20, 24, fmt.Sprintf("P(free|spam)=%.2f  P(free|ham)=%.2f    P(meeting|spam)=%.2f  P(meeting|ham)=%.2f",
		PFreeGivenSpam, PFreeGivenHam, PMeetingGivenSpam, PMeetingGivenHam), 12, viz.Muted, "start")

	y := 70.0
	y = bar(c, y, "Prior — before checking any word", prior)
	y = bar(c, y, fmt.Sprintf("Email %s 'free'  (likelihood ratio ×%.2f)", freeWord, ratioFree), afterFree)
	y = bar(c, y, fmt.Sprintf("...and %s 'meeting'  (likelihood ratio ×%.2f)", meetingWord, ratioMeeting), afterBoth)

	verdict, vColor := "HAM", viz.Good
	if afterBoth >= 0.5 {
		verdict, vColor = "SPAM", viz.Warm
	}
	c.Text(20, y+8, fmt.Sprintf("Final verdict: %s — P(spam) = %.1f%%", verdict, afterBoth*100), 15, vColor, "start")
	c.Text(20, y+30, "gray tick = the 50% classification boundary each bar crosses",
		12, viz.Muted, "start")

	return c.String()
}
