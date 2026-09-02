// Package attention visualizes scaled dot-product attention: comparing a
// query vector against several candidate key vectors with the same dot
// product `cosine-similarity` uses, turning those scores into weights with
// `sigmoid-softmax`'s Softmax, and blending every candidate's value by
// that weight into one output. The running example is the pronoun "it" in
// "The animal didn't cross the street because it was too tired," deciding
// how much "it" should be read as the animal versus the street.
package attention

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "attention-mechanism",
		Seq:   88,
		Title: "Attention mechanism (weighting what a query attends to)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`cosine-similarity` showed how to measure whether two vectors point the " +
						"same way — handy for comparing exactly two things at once. But language " +
						"rarely resolves one word using just one other word: in 'The animal didn't " +
						"cross the street because it was too tired,' the pronoun 'it' has to be " +
						"matched against every earlier noun to figure out who's tired — and forcing " +
						"a single winner ('it' = 100% animal, 0% anything else) throws away real " +
						"ambiguity that a good reader (or model) should be able to represent when a " +
						"sentence is genuinely unclear. Is there a way to compare one 'query' — " +
						"what a word is currently trying to resolve — against several candidates at " +
						"once, and blend their information together in proportion to how well each " +
						"one matches, rather than picking exactly one or averaging all of them " +
						"equally?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Represent 'it' as a query vector q, and three candidate tokens as key " +
						"vectors it's compared against:",
					"• 'The' — key=(1,0), value=0 (a function word: no content to contribute " +
						"even if attention landed on it)",
					"• 'animal' — key=(0,1), value=1 (attention landing here contributes '+1, " +
						"it's the creature')",
					"• 'street' — key=(−0.87,−0.5), value=−1 (attention landing here contributes " +
						"'−1, it's the place')",
					"With q=(0,2) (leaning toward 'animal's key direction), `cosine-similarity`'s " +
						"dot product measures how well q lines up with each key: q·key_The=0, " +
						"q·key_animal=2, q·key_street=−1. Divide each by √2 (the key dimension) — " +
						"the 'scaled' in scaled dot-product attention, which keeps scores from " +
						"growing with vector dimension and swamping softmax before it even runs — " +
						"giving scores 0, 1.414, −0.707.",
					"Feed those three scores into `sigmoid-softmax`'s Softmax: weights = 17.8%, " +
						"73.4%, 8.8% — all three add to 100%, and 'animal' dominates because its " +
						"key best matches q, without the other two being zeroed out entirely. The " +
						"output is the weighted sum of every token's value: 0.178×0 + 0.734×1 + " +
						"0.088×(−1) = 0.646 — a number close to +1 ('it' mostly resolves to " +
						"'animal'), not exactly 1, honestly reflecting that a sliver of weight " +
						"still sits on 'street'.",
					"Flip the query to q=(0,−2) (aligned with 'street' instead) and the same " +
						"pipeline gives weights 30.6%, 7.4%, 62.0%, output −0.546 — 'it' now mostly " +
						"resolves the other way. A perfectly ambiguous query, q≈(0,0), lands close " +
						"to equal thirds on every token and an output near 0: genuine uncertainty, " +
						"represented as a number instead of a forced guess.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"qx and qy position the query vector; temp is `sigmoid-softmax`'s softmax " +
						"temperature. Three bars, one per token, show that token's attention " +
						"weight — tallest is whichever key the query currently lines up with best. " +
						"The readout beneath gives each token's scaled score and the blended " +
						"output number, colored toward green (near +1, 'animal') or red (near −1, " +
						"'street') by its sign. Drag temp down and the tallest bar shoots toward " +
						"100% (nearly all weight on the best-aligned token, e.g. 99% at temp=0.3); " +
						"drag it up and the bars flatten toward each other (weights closer to 33% " +
						"each at temp=3) — the exact same sharpen/flatten behavior " +
						"`sigmoid-softmax` already showed, now deciding how confidently attention " +
						"commits to one token instead of spreading across several.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Blend information from several sources at once, weighted by relevance " +
						"computed directly from the data (the query and keys) rather than from a " +
						"fixed, hand-set rule — and represent genuine ambiguity as a spread-out " +
						"weighting instead of forcing one hard choice. This is the actual mechanism " +
						"(scaled dot-product attention) inside every transformer-based language " +
						"model: every word in a sentence computes a query, compares it against " +
						"every other word's key exactly this way, and blends their values — " +
						"repeated across many words and many independent query/key/value patterns " +
						"at once (multi-head attention) — to build up context-aware representations " +
						"of meaning.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Every modern large language model — GPT, Claude, and the rest — resolves " +
						"pronouns, disambiguates words with multiple meanings ('bank' the " +
						"riverbank vs. 'bank' the building), and tracks relationships between " +
						"distant words in a sentence using this exact mechanism, scaled up to " +
						"thousands of tokens and hundreds of query/key/value patterns running in " +
						"parallel. It's also used well outside language: image transformers attend " +
						"one patch of an image to every other patch, and recommendation systems " +
						"attend a user's current context to their past behavior, both using the " +
						"same scaled dot-product-then-softmax pipeline.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'attention weights don't pick one token and ignore the " +
						"rest — they're a full probability distribution, so even the losing tokens " +
						"still contribute a little unless their weight rounds all the way down to " +
						"zero.' The 8.8% still sitting on 'street' in the default example isn't a " +
						"rounding error; it's a real (small) contribution to the blended output.",
					"Not like this: assuming a higher raw dot-product score always means 'more " +
						"related' in some absolute sense, independent of the other candidates. " +
						"Attention weights are relative — softmax normalizes every score against " +
						"all the others being compared in the same pass, so the same key can get " +
						"very different attention depending entirely on what else it's competing " +
						"against, not on any property of that key alone.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "qx", Label: "query x", Min: -2, Max: 2, Step: 0.1, Def: 0},
			{Key: "qy", Label: "query y", Min: -2, Max: 2, Step: 0.1, Def: 2},
			{Key: "temp", Label: "softmax temperature", Min: 0.2, Max: 3, Step: 0.1, Def: 1},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(680, 420, 0, 1, 0, 1)
	c.Text(16, 24, "attention-mechanism: scaffold -- render not implemented yet", 14, viz.Ink, "start")
	return c.String()
}
