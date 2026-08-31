// Package huffman visualizes Huffman coding: greedily merging the two
// least-frequent symbols over and over builds a binary tree whose leaves
// are the symbols, with frequent symbols landing shallow (short codes) and
// rare symbols landing deep (long codes) -- and the resulting average code
// length always comes within 1 bit of `entropy`'s theoretical floor.
package huffman

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

// Node is one node of a Huffman tree: a leaf (Label set, Left/Right nil)
// carries one symbol's frequency, an internal node (Label empty) carries
// the combined frequency of everything beneath it.
type Node struct {
	Freq        float64
	Label       string
	Left, Right *Node
	order       int // insertion order, used only to break Freq ties deterministically while building
}

// Frequencies returns a fixed 4-symbol alphabet (A, B, C, D) with
// frequencies shaped by skew: raw weight skew^i for the i-th symbol,
// normalized to sum to 1. skew=1 gives a uniform distribution (all four
// equally likely); skew near 0 concentrates almost everything on A.
func Frequencies(skew float64) (labels []string, freqs []float64) {
	labels = []string{"A", "B", "C", "D"}
	raw := make([]float64, len(labels))
	w, sum := 1.0, 0.0
	for i := range raw {
		raw[i] = w
		sum += w
		w *= skew
	}
	freqs = make([]float64, len(raw))
	for i, r := range raw {
		freqs[i] = r / sum
	}
	return labels, freqs
}

// BuildHuffmanTree builds the Huffman tree for the given labels and their
// matching frequencies: repeatedly take the two lowest-frequency nodes
// (ties broken by whichever was created first) and merge them into a new
// parent node, until one node -- the root -- remains. Frequent symbols end
// up merged in late, near the root (short codes); rare symbols get merged
// in early, burying them deep (long codes). Returns nil for an empty
// alphabet, and the lone leaf itself (never merged) for a single symbol.
func BuildHuffmanTree(labels []string, freqs []float64) *Node {
	if len(labels) == 0 {
		return nil
	}
	nodes := make([]*Node, len(labels))
	for i := range labels {
		nodes[i] = &Node{Freq: freqs[i], Label: labels[i], order: i}
	}
	next := len(nodes)
	for len(nodes) > 1 {
		sort.SliceStable(nodes, func(i, j int) bool {
			if nodes[i].Freq != nodes[j].Freq {
				return nodes[i].Freq < nodes[j].Freq
			}
			return nodes[i].order < nodes[j].order
		})
		a, b := nodes[0], nodes[1]
		merged := &Node{Freq: a.Freq + b.Freq, Left: a, Right: b, order: next}
		next++
		nodes = append(nodes[2:], merged)
	}
	return nodes[0]
}

// Codes walks the tree and returns each leaf's binary codeword: the path
// from the root, "0" for every left branch and "1" for every right one. A
// single-leaf tree (nothing to distinguish) is given the codeword "0"
// rather than the empty string.
func Codes(root *Node) map[string]string {
	out := map[string]string{}
	var walk func(n *Node, prefix string)
	walk = func(n *Node, prefix string) {
		if n == nil {
			return
		}
		if n.Left == nil && n.Right == nil {
			if prefix == "" {
				prefix = "0"
			}
			out[n.Label] = prefix
			return
		}
		walk(n.Left, prefix+"0")
		walk(n.Right, prefix+"1")
	}
	walk(root, "")
	return out
}

// CodeLengths returns each leaf's codeword length in bits -- its depth in
// the tree -- without building the codeword strings themselves.
func CodeLengths(root *Node) map[string]int {
	out := map[string]int{}
	var walk func(n *Node, depth int)
	walk = func(n *Node, depth int) {
		if n == nil {
			return
		}
		if n.Left == nil && n.Right == nil {
			if depth == 0 {
				depth = 1
			}
			out[n.Label] = depth
			return
		}
		walk(n.Left, depth+1)
		walk(n.Right, depth+1)
	}
	walk(root, 0)
	return out
}

// AverageCodeLength returns the expected number of bits per symbol when
// encoding with the given code lengths: the frequency-weighted average
// Σ freq_i * length_i.
func AverageCodeLength(labels []string, freqs []float64, lengths map[string]int) float64 {
	avg := 0.0
	for i, l := range labels {
		avg += freqs[i] * float64(lengths[l])
	}
	return avg
}

// EntropyBits returns the Shannon entropy, in bits, of a discrete
// distribution given as a list of probabilities -- the same quantity
// `entropy` computes for two outcomes, generalized to any number of them:
// Σ -p_i*log2(p_i). This is the theoretical floor no code, Huffman
// included, can beat on average.
func EntropyBits(freqs []float64) float64 {
	h := 0.0
	for _, p := range freqs {
		if p > 0 {
			h += -p * math.Log2(p)
		}
	}
	return h
}

// FixedLengthBits returns the number of bits a fixed-width code needs per
// symbol to distinguish n symbols: ceil(log2(n)) -- what you'd use without
// bothering to look at the frequencies at all.
func FixedLengthBits(n int) float64 {
	return math.Ceil(math.Log2(float64(n)))
}

func init() {
	concept.Register(concept.Concept{
		ID:    "huffman-coding",
		Seq:   84,
		Title: "Huffman coding (an optimal code built from symbol frequencies)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`entropy` measured how many bits of genuine surprise a distribution " +
						"carries on average -- but it never said how to actually build a code " +
						"that gets close to that number. The obvious approach is a fixed-length " +
						"code: with 4 symbols, ⌈log2(4)⌉=2 bits apiece, no exceptions, whether a " +
						"symbol is the most common one in the alphabet or the rarest. That's " +
						"simple, but it spends the same 2 bits on a symbol that shows up over " +
						"half the time as on one that shows up a fraction of that -- treating a " +
						"near-certainty and a rare surprise as equally expensive to send, even " +
						"though `entropy` already said they're not. Can you build an actual, " +
						"usable code -- one where a stream of bits still unambiguously decodes " +
						"back to the original symbols -- that spends fewer bits on frequent " +
						"symbols and more on rare ones, pulling the average down toward entropy's " +
						"floor?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take the 4-symbol alphabet A, B, C, D with frequencies 8/15≈0.533, " +
						"4/15≈0.267, 2/15≈0.133, and 1/15≈0.067 (skew=0.5 below). Build the tree " +
						"bottom-up by repeatedly merging the two least-frequent nodes into one:",
					"• Merge the two smallest: C(0.133) and D(0.067) → a new combined node " +
						"CD=0.200.",
					"• Merge the two smallest again: CD(0.200) and B(0.267) → a new combined " +
						"node BCD=0.467.",
					"• Merge the last two: BCD(0.467) and A(0.533) → the root, 1.000.",
					"Read each symbol's code off its path down from the root -- 0 for every " +
						"left branch, 1 for every right one. A sits just one merge below the " +
						"root: code '1', 1 bit. B is two merges down: code '01', 2 bits. C and D " +
						"are each three merges down: codes '001' and '000', 3 bits apiece.",
					"Weighted average: 8/15×1 + 4/15×2 + 2/15×3 + 1/15×3 = (8+8+6+3)/15 = " +
						"25/15 ≈ 1.667 bits/symbol -- noticeably under the fixed 2 bits/symbol, " +
						"and only 0.026 bits above this distribution's 1.640-bit entropy floor.",
					"Why merging-the-smallest-first works: every merge buries that round's two " +
						"least-likely candidates one level deeper, and one level deeper costs " +
						"exactly one more bit. The algorithm always spends that extra bit where " +
						"it's cheapest -- on whatever is least frequent at that moment -- never on " +
						"something more common that an earlier merge has already promoted higher " +
						"up the tree.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"skew reshapes the four frequencies (A always the largest, D always the " +
						"smallest); the tree below rebuilds itself merge by merge to match. Each " +
						"leaf square shows its symbol, frequency, and Huffman codeword; each " +
						"branch is labeled with the bit it contributes, so tracing a leaf's path " +
						"back to the root spells out its codeword directly. Drag skew toward 0.95 " +
						"and the frequencies flatten out, the tree balances into every leaf sitting " +
						"2 levels down, and the Huffman average converges on the fixed-length cost " +
						"it was trying to beat. Drag skew toward 0.15 and A dominates so heavily " +
						"that it keeps its 1-bit code while B, C, and D get pushed progressively " +
						"deeper. The bottom readout tracks the Huffman average, the entropy floor, " +
						"the fixed-length cost, and the gap between the first two, live.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Build an actual prefix-free code from nothing but symbol frequencies, " +
						"provably as good as any code can get to within 1 bit of entropy (Shannon's " +
						"source coding theorem) -- turning entropy from a description of a limit " +
						"into a limit you can actually approach with a constructive algorithm. And " +
						"the savings are real, not just theoretical: at skew=0.5 above, encoding " +
						"1000 symbols costs about 1667 bits with Huffman versus 2000 bits " +
						"fixed-length, a 16.7% reduction, for free, just by looking at the " +
						"frequencies before choosing the codes.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"ZIP, DEFLATE, and gzip all run Huffman coding on the back end, after a " +
						"separate pass finds repeated sequences to shorten first. JPEG and MP3 " +
						"Huffman-encode their final quantized values before writing the file. " +
						"Morse code is a hand-built approximation of the exact same idea, " +
						"invented a century before Huffman's algorithm: E, the most common letter " +
						"in English, gets a single dot; Q, one of the rarest, gets four symbols.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'Huffman coding assigns shorter codes to more frequent " +
						"symbols and always lands within 1 bit of entropy per symbol -- it reads " +
						"the frequencies and builds the optimal tree from them directly, instead " +
						"of guessing ahead of time how much a symbol will be used.'",
					"Not like this: assuming a shorter Huffman code for a frequent symbol means " +
						"it somehow compressed that individual symbol -- one occurrence of A " +
						"still costs exactly 1 bit here, not a fraction of one; the savings only " +
						"show up in aggregate, averaged across many symbols, which is exactly why " +
						"entropy (itself an average) is the right thing to compare against, not " +
						"any single code length in isolation. Also not like this: expecting " +
						"Huffman to ever beat entropy -- the gap can shrink toward 0 (skew=0.5 " +
						"above gets within 0.026 bits) but the average length can never dip below " +
						"the entropy floor itself, for any distribution.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "skew", Label: "Frequency skew", Min: 0.15, Max: 0.95, Step: 0.05, Def: 0.5},
		},
		Render: render,
	})
}

// messageLen is the message length used to translate "bits per symbol" into
// a concrete total-bits comparison in the footer readout.
const messageLen = 1000

func render(p map[string]float64) string {
	skew := p["skew"]
	labels, freqs := Frequencies(skew)
	root := BuildHuffmanTree(labels, freqs)
	lengths := CodeLengths(root)
	codes := Codes(root)
	avg := AverageCodeLength(labels, freqs, lengths)
	ent := EntropyBits(freqs)
	fixed := FixedLengthBits(len(labels))

	// Depth of every node (root=0), and the tree's overall depth, so leaves
	// and internal nodes alike can be placed on the right horizontal row.
	depths := map[*Node]int{}
	maxDepth := 0
	var walkDepth func(n *Node, d int)
	walkDepth = func(n *Node, d int) {
		if n == nil {
			return
		}
		depths[n] = d
		if d > maxDepth {
			maxDepth = d
		}
		walkDepth(n.Left, d+1)
		walkDepth(n.Right, d+1)
	}
	walkDepth(root, 0)
	if maxDepth == 0 {
		maxDepth = 1 // guards the (unused here) single-leaf case from a divide-by-zero below
	}

	// x position: leaves spread left-to-right in the order the tree visits
	// them; every internal node sits above the midpoint of its children.
	nLeaves := len(labels)
	leafIdx := 0
	xPos := map[*Node]float64{}
	var assignX func(n *Node) float64
	assignX = func(n *Node) float64 {
		if n.Left == nil && n.Right == nil {
			x := 0.08 + (float64(leafIdx)+0.5)/float64(nLeaves)*0.84
			leafIdx++
			xPos[n] = x
			return x
		}
		lx := assignX(n.Left)
		rx := assignX(n.Right)
		x := (lx + rx) / 2
		xPos[n] = x
		return x
	}
	assignX(root)

	// y position: root at the top (large data-y, since Canvas.Y flips so
	// larger data-y renders higher), leaves at the bottom.
	yFor := func(depth int) float64 {
		return 0.95 - float64(depth)/float64(maxDepth)*0.80
	}

	c := viz.New(760, 580, 0, 1, 0, 1)
	c.PadT = 130 // room for three lines of header text above the tree
	c.PadB = 130 // room for the bits-per-symbol comparison below it

	// Edges, drawn before nodes so the node markers sit on top of them.
	// Each edge is labeled with the bit it contributes (0 for the left
	// branch, 1 for the right) -- the same path Codes walks to build each
	// leaf's codeword.
	var drawEdges func(n *Node)
	drawEdges = func(n *Node) {
		if n == nil || (n.Left == nil && n.Right == nil) {
			return
		}
		for _, child := range []*Node{n.Left, n.Right} {
			c.Path([][2]float64{{xPos[n], yFor(depths[n])}, {xPos[child], yFor(depths[child])}}, viz.Muted, 1.5)
			bit := "0"
			if child == n.Right {
				bit = "1"
			}
			mx := (xPos[n] + xPos[child]) / 2
			my := (yFor(depths[n]) + yFor(depths[child])) / 2
			c.Text(c.X(mx)+6, c.Y(my), bit, 12, viz.Accent, "middle")
		}
		drawEdges(n.Left)
		drawEdges(n.Right)
	}
	drawEdges(root)

	// Nodes: a tinted square for every leaf, its symbol, frequency, and
	// Huffman codeword printed below it; a small dot with just the
	// combined frequency for every internal (merge) node.
	var drawNodes func(n *Node)
	drawNodes = func(n *Node) {
		if n == nil {
			return
		}
		px, py := c.X(xPos[n]), c.Y(yFor(depths[n]))
		if n.Left == nil && n.Right == nil {
			c.Rect(px-16, py-16, 32, 32, viz.Accent, 0.16)
			c.Text(px, py+5, n.Label, 15, viz.Ink, "middle")
			c.Text(px, py+30, fmt.Sprintf("p=%.3f", n.Freq), 11, viz.Muted, "middle")
			c.Text(px, py+46, fmt.Sprintf("code %s (%d bit)", codes[n.Label], lengths[n.Label]), 11, viz.Ink, "middle")
			return
		}
		c.Rect(px-5, py-5, 10, 10, viz.Muted, 0.5)
		c.Text(px, py-12, fmt.Sprintf("%.3f", n.Freq), 10, viz.Muted, "middle")
		drawNodes(n.Left)
		drawNodes(n.Right)
	}
	drawNodes(root)

	freqStrs := make([]string, len(labels))
	for i, l := range labels {
		freqStrs[i] = fmt.Sprintf("%s=%.3f", l, freqs[i])
	}
	c.Text(20, 24, fmt.Sprintf("Symbol frequencies (skew=%.2f): %s", skew, strings.Join(freqStrs, "  ")),
		14, viz.Ink, "start")
	c.Text(20, 46, "Merge the two least-frequent nodes repeatedly (greedy), building the tree bottom-up",
		12, viz.Muted, "start")
	c.Text(20, 66, "left branch = bit 0, right branch = bit 1 -- a symbol's codeword is its path from the root",
		12, viz.Muted, "start")

	saved := 0.0
	if fixed > 0 {
		saved = (fixed - avg) / fixed * 100
	}
	c.Text(20, c.H-92, fmt.Sprintf("Huffman average = %.3f bits/symbol    entropy floor = %.3f bits    fixed-length = %.0f bits",
		avg, ent, fixed), 13, viz.Ink, "start")
	c.Text(20, c.H-70, fmt.Sprintf("gap above entropy = %.3f bits (Huffman is always within 1 bit of entropy)", avg-ent),
		12, viz.Muted, "start")
	c.Text(20, c.H-48, fmt.Sprintf("over %d symbols: fixed-length ~%.0f bits vs Huffman ~%.0f bits (%.1f%% smaller)",
		messageLen, fixed*messageLen, avg*messageLen, saved), 12, viz.Muted, "start")
	c.Text(20, c.H-20, "square = leaf symbol    dot = merged (internal) node, labeled with its combined frequency",
		12, viz.Muted, "start")

	return c.String()
}
