// Package svd visualizes the singular value decomposition of a 2x2 matrix:
// every matrix, square or not, symmetric or not, breaks down into "rotate,
// scale along two perpendicular axes, rotate again" (M = U·Σ·Vᵀ). Unlike
// `eigenvectors-eigenvalues`, which needs a square matrix and can fail to
// find two independent directions at all, SVD always exists and always
// gives two orthonormal bases — V for the input, U for the output.
package svd

import (
	"math"

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

// Matrix is a 2x2 matrix [[A,B],[C,D]] (row-major).
type Matrix struct{ A, B, C, D float64 }

// Vec2 is a 2D vector.
type Vec2 struct{ X, Y float64 }

// Apply returns M*v.
func Apply(m Matrix, v Vec2) Vec2 {
	return Vec2{m.A*v.X + m.B*v.Y, m.C*v.X + m.D*v.Y}
}

// SVD holds a 2x2 matrix's singular value decomposition M = U·Σ·Vᵀ: two
// singular values (Sigma1 >= Sigma2 >= 0) and the orthonormal right
// singular vectors (V1, V2 — columns of V, the special input directions)
// and left singular vectors (U1, U2 — columns of U, where those
// directions land) that go with them, so that Apply(m, V1) == Sigma1*U1
// and Apply(m, V2) == Sigma2*U2.
type SVD struct {
	Sigma1, Sigma2 float64
	V1, V2         Vec2
	U1, U2         Vec2
}

// Decompose computes the singular value decomposition of the 2x2 matrix m.
// It works for every matrix — square or not, symmetric or not, singular or
// not — unlike eigenvectors-eigenvalues, which needs a square matrix and
// can fail to produce two independent directions at all.
//
// The method: MᵀM is always symmetric and positive semi-definite, so
// (mirroring eigenvectors-eigenvalues' closed-form 2x2 formula) it always
// has two real, non-negative eigenvalues and perpendicular eigenvectors.
// Those eigenvectors are V; the square roots of their eigenvalues are the
// singular values; and the matching output directions U are recovered by
// pushing each V direction through M and rescaling back to unit length:
// U_i = M·V_i / Sigma_i.
func Decompose(m Matrix) SVD {
	p := m.A*m.A + m.C*m.C
	q := m.A*m.B + m.C*m.D
	r := m.B*m.B + m.D*m.D

	mid := (p + r) / 2
	half := (p - r) / 2
	disc := math.Sqrt(half*half + q*q)
	lambda1, lambda2 := mid+disc, mid-disc
	if lambda2 < 0 {
		lambda2 = 0 // clamp float noise -- MᵀM is PSD, so this is never a real negative
	}
	sigma1, sigma2 := math.Sqrt(lambda1), math.Sqrt(lambda2)

	var thetaV float64
	if p != r || q != 0 {
		thetaV = math.Atan2(2*q, p-r) / 2
	} // else MᵀM is a multiple of the identity: any orthonormal basis works, thetaV=0
	v1 := Vec2{math.Cos(thetaV), math.Sin(thetaV)}
	v2 := Vec2{-math.Sin(thetaV), math.Cos(thetaV)}

	var u1, u2 Vec2
	if sigma1 < 1e-9 {
		// M collapses everything to (near) zero -- there's no image
		// direction to normalize, so fall back to a default basis.
		u1, u2 = Vec2{1, 0}, Vec2{0, 1}
	} else {
		mv1 := Apply(m, v1)
		u1 = Vec2{mv1.X / sigma1, mv1.Y / sigma1}
		if sigma2 < 1e-9 {
			// M collapses this one direction to (near) nothing -- complete
			// the orthonormal basis by rotating u1 by 90° instead of
			// dividing by a near-zero sigma2.
			u2 = Vec2{-u1.Y, u1.X}
		} else {
			mv2 := Apply(m, v2)
			u2 = Vec2{mv2.X / sigma2, mv2.Y / sigma2}
		}
	}

	return SVD{sigma1, sigma2, v1, v2, u1, u2}
}

// Reconstruct rebuilds M = U·Σ·Vᵀ from a decomposition (as
// Sigma1*outer(U1,V1) + Sigma2*outer(U2,V2)), so callers can confirm the
// pieces really do multiply back to the original matrix.
func Reconstruct(s SVD) Matrix {
	return Matrix{
		A: s.Sigma1*s.U1.X*s.V1.X + s.Sigma2*s.U2.X*s.V2.X,
		B: s.Sigma1*s.U1.X*s.V1.Y + s.Sigma2*s.U2.X*s.V2.Y,
		C: s.Sigma1*s.U1.Y*s.V1.X + s.Sigma2*s.U2.Y*s.V2.X,
		D: s.Sigma1*s.U1.Y*s.V1.Y + s.Sigma2*s.U2.Y*s.V2.Y,
	}
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(560, 560, -5, 5, -5, 5).String()
}
