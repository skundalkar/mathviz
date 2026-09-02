package attention

import "math"

// Token is one fixed key/value entry the query attends over: a 2D key
// vector (a stand-in for that token's embedding direction) and a scalar
// value (the "meaning" it contributes to the output if attention lands on
// it).
type Token struct {
	Name  string
	Key   [2]float64
	Value float64
}

// Tokens is the running example's fixed vocabulary, modeling "The animal
// didn't cross the street because it was too tired" — the query below
// stands in for the pronoun "it", deciding which noun it actually refers
// to. "The" carries no content of its own (Value 0); "animal" and "street"
// pull the output toward +1 ("it" = the creature) or −1 ("it" = the
// place) depending on which key direction the query lines up with.
var Tokens = []Token{
	{Name: "The", Key: [2]float64{1, 0}, Value: 0},
	{Name: "animal", Key: [2]float64{0, 1}, Value: 1},
	{Name: "street", Key: [2]float64{-0.866, -0.5}, Value: -1},
}

// dim is the key/query vector dimension (2 here) -- ScaledScores divides
// by sqrt(dim), the "scaled" in scaled dot-product attention.
const dim = 2

// Dot returns the dot product of two 2D vectors: ax*bx + ay*by.
// cosine-similarity's Dot, duplicated here so this package stays
// self-contained (see BUILD_CYCLE.md).
func Dot(ax, ay, bx, by float64) float64 {
	return ax*bx + ay*by
}

// ScaledScores returns one similarity score per token for query (qx,qy):
// the query dotted with that token's key, divided by sqrt(dim). Dividing
// by sqrt(dim) is the "scaled" part of scaled dot-product attention: raw
// dot products grow with vector dimension, which would otherwise push
// softmax toward saturating on whichever key has the largest magnitude
// rather than whichever one the query actually aligns with.
func ScaledScores(qx, qy float64) []float64 {
	scale := math.Sqrt(dim)
	scores := make([]float64, len(Tokens))
	for i, t := range Tokens {
		scores[i] = Dot(qx, qy, t.Key[0], t.Key[1]) / scale
	}
	return scores
}

// Softmax turns scores into attention weights that sum to 1: each score is
// divided by temp, exponentiated, and normalized by the total. This is
// sigmoid-softmax's Softmax, duplicated here so this package stays
// self-contained. Lower temp sharpens the weights toward whichever score
// is largest; higher temp flattens them toward uniform (1/n each).
// Scores are shifted by their max before exponentiating -- shift-invariant
// for softmax, but keeps math.Exp from overflowing on large scores.
func Softmax(scores []float64, temp float64) []float64 {
	out := make([]float64, len(scores))
	if len(scores) == 0 {
		return out
	}
	if temp <= 0 {
		temp = 1e-9
	}
	max := scores[0]
	for _, s := range scores[1:] {
		if s > max {
			max = s
		}
	}
	var sum float64
	for i, s := range scores {
		e := math.Exp((s - max) / temp)
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

// Output returns the attention-weighted sum of every token's value -- the
// single number attention actually hands downstream, blending each
// token's contribution by how much the query attended to it.
func Output(weights []float64) float64 {
	sum := 0.0
	for i, w := range weights {
		sum += w * Tokens[i].Value
	}
	return sum
}
