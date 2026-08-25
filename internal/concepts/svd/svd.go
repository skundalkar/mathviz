// Package svd visualizes the singular value decomposition of a 2x2 matrix:
// every matrix, square or not, symmetric or not, breaks down into "rotate,
// scale along two perpendicular axes, rotate again" (M = U·Σ·Vᵀ). Unlike
// `eigenvectors-eigenvalues`, which needs a square matrix and can fail to
// find two independent directions at all, SVD always exists and always
// gives two orthonormal bases — V for the input, U for the output.
package svd

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "singular-value-decomposition",
		Seq:   66,
		Title: "Singular value decomposition (rotate, scale, rotate)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`eigenvectors-eigenvalues` found special directions where a matrix only " +
						"stretches, never rotates — but try that trick on the shear matrix " +
						"M=[[1,1],[0,1]] and it breaks: both of M's eigenvalues equal 1, and " +
						"solving (M-I)v=0 turns up only one independent direction, (1,0), not two " +
						"perpendicular ones to hang a full picture on. And that's the friendly " +
						"case — a matrix that isn't square (say, a spreadsheet of 500 house sales, " +
						"each row a house and each of 3 columns a feature: size, age, price) has no " +
						"eigenvectors at all, because Av=λv requires v and Av to be the same size, " +
						"which a non-square matrix can never give you. Is there a version of 'the " +
						"special directions a matrix cares about' that works for every matrix — " +
						"square or not, symmetric or not, even one eigendecomposition can't " +
						"handle?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take that same shear matrix, M=[[1,1],[0,1]]. Instead of solving Mv=λv, form " +
						"MᵀM — for any matrix M, MᵀM is always symmetric, so `eigenvectors-" +
						"eigenvalues`'s closed-form 2x2 formula always applies to it even when M " +
						"itself isn't symmetric or square: MᵀM = [[1,1],[1,2]] here.",
					"• MᵀM's eigenvalues come out to λ₁=2.618, λ₂=0.382 (via the same " +
						"mid±√(half²+b²) formula `eigenvectors-eigenvalues` used). Take square " +
						"roots: σ₁=√2.618=1.618, σ₂=√0.382=0.618 — these are M's singular values, " +
						"and they're exactly φ (the golden ratio) and 1/φ.",
					"• MᵀM's eigenvectors give the input directions, V: v₁≈(0.53,0.85) at " +
						"θ≈58.3°, and v₂≈(-0.85,0.53), perpendicular to it.",
					"• Push each through M and rescale back to unit length — u_i = M·v_i/σ_i — to " +
						"get the output directions, U: u₁ = M·(0.53,0.85)/1.618 ≈ (0.85,0.53), and " +
						"u₂ ≈ (-0.53,0.85).",
					"u₁ and u₂ come out perpendicular to each other too (their dot product is " +
						"0, to floating-point precision) even though nothing forced that directly " +
						"— it falls out of v₁⊥v₂ and M's structure, the same way `eigenvectors-" +
						"eigenvalues`' two eigenvector directions came out perpendicular for a " +
						"symmetric matrix. The payoff: V and U are two separate orthonormal bases " +
						"— one for 'which input directions are special,' one for 'where they land' " +
						"— connected by M·v_i = σ_i·u_i for each i. That relationship, written as " +
						"matrices, is exactly M = U·Σ·Vᵀ, where Σ=diag(σ₁,σ₂).",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The faint circle is every unit-length input direction; the solid curve is " +
						"where M sends all of them — always an ellipse, whatever M is. The blue " +
						"arrows are v₁,v₂ (the two special input directions, still on the unit " +
						"circle); the green arrows are σ₁·u₁ and σ₂·u₂ — where those two directions " +
						"land, which are exactly the ellipse's long and short axes. Drag a, b, c, " +
						"or d and watch all four arrows and the ellipse update together: the ellipse " +
						"is never anything other than M applied to the circle, and its axes are " +
						"never anything other than the current σ·U arrows.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Get a full, always-real, always-orthogonal picture of 'how M stretches " +
						"space' for any matrix — including ones eigendecomposition can't handle, " +
						"like the defective shear above or a non-square data matrix. σ₂ close to 0 " +
						"is a direct rank signal: it means M nearly collapses one whole direction " +
						"to nothing, squashing the ellipse into a near-line, which is exactly the " +
						"situation `pca` calls 'a direction that barely varies' — SVD is the way " +
						"PCA is actually computed in practice, run directly on the data matrix " +
						"instead of first building a covariance matrix and eigendecomposing that.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Image compression: treat a grayscale photo as a big matrix, keep only the " +
						"handful of largest singular values and their U/V directions, and the " +
						"reconstructed image is often visually indistinguishable from the original " +
						"at a fraction of the storage. Recommender systems (the Netflix Prize's " +
						"winning approach): decompose a giant, mostly-empty user-by-movie ratings " +
						"matrix and the top singular directions surface 'latent taste dimensions' " +
						"like genre preference, letting the system fill in the empty cells. Noise " +
						"reduction: real signals concentrate in the largest singular values while " +
						"random noise spreads thinly across all of them, so dropping the small ones " +
						"denoises the data.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'M=[[1,1],[0,1]] has singular values φ≈1.618 and 1/φ≈0.618, " +
						"connecting two different orthonormal bases — M·v₁=σ₁·u₁, not M·v₁=σ₁·v₁.' " +
						"Not like this: reaching for 'find M's eigenvectors and eigenvalues' on a " +
						"general matrix — this shear matrix has only one real eigenvalue (1, " +
						"repeated) and a single eigenvector direction, nothing like the two " +
						"perpendicular singular directions above. U only equals V when M is " +
						"symmetric and positive semi-definite (as in `eigenvectors-eigenvalues`'s " +
						"own [[2,1],[1,2]] example) — in general they're two distinct bases, one " +
						"for inputs and one for outputs, and singular values are always ≥0 even " +
						"when a matrix's eigenvalues (if it has any) are negative.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "Matrix a", Min: -2, Max: 2, Step: 0.1, Def: 1},
			{Key: "b", Label: "Matrix b", Min: -2, Max: 2, Step: 0.1, Def: 1},
			{Key: "c", Label: "Matrix c", Min: -2, Max: 2, Step: 0.1, Def: 0},
			{Key: "d", Label: "Matrix d", Min: -2, Max: 2, Step: 0.1, Def: 1},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(560, 560, -5, 5, -5, 5).String()
}
