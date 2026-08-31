// Package huffman visualizes Huffman coding: greedily merging the two
// least-frequent symbols over and over builds a binary tree whose leaves
// are the symbols, with frequent symbols landing shallow (short codes) and
// rare symbols landing deep (long codes) -- and the resulting average code
// length always comes within 1 bit of `entropy`'s theoretical floor.
package huffman

import (
	"math"
	"sort"

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
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "skew", Label: "Frequency skew", Min: 0.15, Max: 0.95, Step: 0.05, Def: 0.5},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(720, 460, 0, 1, 0, 1)
	return c.String()
}
