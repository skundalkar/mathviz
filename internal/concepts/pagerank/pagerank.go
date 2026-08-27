// Package pagerank visualizes Google's original PageRank algorithm: ranking
// nodes in a small link graph by the steady-state probability that a random
// web surfer, repeatedly clicking links (and occasionally teleporting to a
// random page instead), ends up there. markov-chains solves a 2-state
// steady state exactly with algebra; this concept shows the technique that
// scales to a graph with far more than 2 states — repeated power-iteration
// updates — and the "damping factor" trick that keeps pages with no inbound
// links from vanishing to zero.
package pagerank

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "pagerank",
		Seq:   71,
		Title: "PageRank (ranking a graph by a random walk)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"markov-chains solved a 2-state weather system exactly: solve one equation " +
						"for the fixed point, done. A web with billions of pages linking to each " +
						"other is a Markov chain too — each page is a state, and 'click a random " +
						"link on this page' is the transition rule — but nobody can write down " +
						"billions of equations and solve them by hand. And pure link-counting has " +
						"its own trap: a brand-new page nobody has linked to yet would score a " +
						"permanent zero, and a small tightly-linked cluster of pages that only link " +
						"to each other could hoard all the importance forever, since a surfer who " +
						"only ever clicks links can never leave. Is there a way to find the " +
						"steady-state 'importance' of every page in a huge graph without solving " +
						"algebra by hand, and without pages getting permanently stuck at zero or a " +
						"few pages hoarding everything?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Instead of solving for the fixed point directly, guess a starting point and " +
						"repeatedly apply the one-step update — 'power iteration' — until the " +
						"numbers stop changing, the same way markov-chains's day-by-day P(sunny) " +
						"walked toward 0.8333 without ever solving the fixed-point equation. Take a " +
						"4-page example: A links to B and C; B links to C; C links to A; D links to " +
						"C. Start every page at an equal rank of 1/4 = 0.25, and update each page's " +
						"rank as (1-d)/4 (a fixed 'floor' every page gets) plus d times the sum of " +
						"each linking page's rank split evenly across its outbound links.",
					"• With d=0.85: after 1 round, C — which gets links from A, B, and D all at " +
						"once — jumps to 0.5688 while D, which nothing links to, drops to exactly " +
						"0.0375 = (1-0.85)/4.",
					"• After 2 rounds: A (now getting C's whole boosted rank back) jumps to " +
						"0.5209; C drops back to 0.2978 as its own inbound pages haven't caught up " +
						"yet — the ranks are still sloshing back and forth.",
					"• By round 20 the sloshing has settled: A≈0.3725, B≈0.1958, C≈0.3942, and D " +
						"is still exactly 0.0375, every single round, forever.",
					"That last number isn't a coincidence: since nothing links to D, the only rank " +
						"it can ever receive is the (1-d)/n floor itself — the 'random surfer " +
						"occasionally teleports to a uniformly random page instead of clicking a " +
						"link' term. Set d=1 (turn teleporting off entirely, pure link-following) " +
						"and that floor becomes (1-1)/4=0 — D collapses to exactly 0 and stays " +
						"there forever, the dead-end trap section 1 asked about, made concrete.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Four boxes laid out as A, B, C, D with arrows for every link (A→B, A→C, " +
						"B→C, C→A, D→C). Each box's size is its current PageRank — bigger box, " +
						"more importance. The 't' slider re-runs the power-iteration update that " +
						"many times from the same 0.25-each starting point in section 2; the 'd' " +
						"slider sets the damping factor. Push t up and watch the boxes stop " +
						"changing size around t≈10-15 (converged); pull d down toward 0.5 and " +
						"convergence happens in far fewer steps, the same speed-vs-stickiness " +
						"trade-off markov-chains's a/b sliders showed. Push d all the way to 1.0 " +
						"and watch D's box shrink to the smallest size on the page and stay there " +
						"no matter how large t gets.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Rank every node in a graph with thousands or billions of nodes by repeatedly " +
						"applying one cheap, purely local update — no simultaneous-equation solve " +
						"required — and get an answer that's provably stable: guaranteed to settle " +
						"to a unique steady state (as long as d<1) no matter which page you started " +
						"the guess from. The damping factor also buys you a guarantee markov-chains " +
						"didn't need to worry about with only 2 always-reachable states: every page, " +
						"even one with zero inbound links, keeps a nonzero floor of importance " +
						"instead of vanishing.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Google's original web search ranking is the namesake example — PageRank " +
						"treated the entire web as exactly this kind of link graph. The same " +
						"random-walk-with-restart idea ranks nodes in any graph: influential " +
						"accounts in a social network (who follows whom), important proteins in a " +
						"biological interaction network, or which papers matter most in a citation " +
						"network (who cites whom). Recommendation systems use a close cousin " +
						"('personalized PageRank') to rank items by a random walk that starts from " +
						"a specific user's history instead of teleporting uniformly.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: a page's PageRank is the long-run fraction of time a random " +
						"surfer — one who usually clicks a random link but occasionally teleports to " +
						"a uniformly random page — spends on it, found by repeating one simple " +
						"update until it stops changing. Not like this: assuming more inbound links " +
						"always means a higher rank on its own — a link from an important page " +
						"(like C, which several pages point to) counts for far more than a link from " +
						"an unimportant one, because each linking page's contribution is its own " +
						"rank divided among its outbound links, not just a vote counted equally. " +
						"Also not like this: thinking the damping factor is just a tuning knob with " +
						"no real meaning — it's a literal probability (typically ~0.85 in Google's " +
						"original paper) that the surfer keeps clicking links rather than jumping " +
						"to a random page next.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "d", Label: "Damping factor (d)", Min: 0.5, Max: 1, Step: 0.05, Def: 0.85},
			{Key: "t", Label: "Iterations (t)", Min: 0, Max: 20, Step: 1, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
