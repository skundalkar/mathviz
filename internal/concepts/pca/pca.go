// Package pca visualizes Principal Component Analysis: given a cloud of
// correlated points, it finds the one direction the cloud spreads out along
// the most (not necessarily the x-axis or the y-axis — those were just an
// arbitrary choice of how to draw the plot) and reports how much of the
// cloud's total variance that single direction accounts for.
package pca

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "principal-component-analysis",
		Seq:   48,
		Title: "Principal Component Analysis (the direction that explains the most)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`covariance` and `correlation` only ever describe two variables against " +
						"each other, and even then the answer depends on which two axes you " +
						"happened to draw the scatter on. Take a cloud of (hours studied, test " +
						"score) points stretched out diagonally, leaning toward neither axis — the " +
						"x-spread alone (variance of hours) and the y-spread alone (variance of " +
						"score) each miss most of the story, because the cloud isn't actually " +
						"spread out along x or along y, it's spread out along its own diagonal. " +
						"`eigenvectors-eigenvalues` already hinted that a covariance matrix has " +
						"'special directions that only stretch' — but which direction is that, for " +
						"an actual cloud of data, and how much of the cloud's total spread does " +
						"that one direction alone account for?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Same five students as `covariance`: hours studied x = 1,2,3,4,5, test score " +
						"(in tens of points) y = 2,4,5,4,5. Their covariance matrix collects three " +
						"numbers into one 2x2 table — variance of x, variance of y, and their " +
						"covariance filling both off-diagonal slots (the matrix is symmetric, since " +
						"Cov(x,y) = Cov(y,x)): [[Var(x), Cov(x,y)], [Cov(x,y), Var(y)]] = " +
						"[[2, 1.2], [1.2, 1.2]] — Var(x)=2 and Cov(x,y)=1.2 exactly as computed in " +
						"`covariance`'s worked example, plus Var(y)=1.2.",
					"Feed that matrix into the exact same closed-form eigenvalue/eigenvector " +
						"formulas `eigenvectors-eigenvalues` used (λ = mid ± disc, θ = " +
						"atan2(2b, a−d)/2): with a=2, b=1.2, d=1.2, mid=(2+1.2)/2=1.6, " +
						"half=(2−1.2)/2=0.4, disc=√(0.4²+1.2²)=√1.6≈1.265, so λ1≈2.865, λ2≈0.335, " +
						"and θ1 = atan2(2.4, 0.8)/2 ≈ 35.8°.",
					"λ1's direction (35.8° off the x-axis — not 0° or 90°, i.e. not aligned to " +
						"either original axis) is the first principal component: the single " +
						"direction the cloud spreads out along the most. Its eigenvalue, 2.865, is " +
						"exactly the variance of the data once you re-measure it along that " +
						"direction instead of along x or y. Add both eigenvalues: 2.865+0.335=3.2 " +
						"— exactly Var(x)+Var(y)=2+1.2=3.2, confirming the eigen-decomposition just " +
						"redistributes the same total variance across two new, uncorrelated " +
						"directions rather than creating or destroying any of it. 2.865/3.2 ≈ 89.5% " +
						"of the cloud's total spread lives along that one 35.8° direction alone.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The r slider sets how correlated x and y are (same meaning as in " +
						"`covariance` and `correlation`); the stretch slider re-stretches the " +
						"x-axis, growing Var(x) the same way `covariance`'s scale slider did. The " +
						"green arrow is the first principal component (PC1) — the direction of " +
						"maximum spread — drawn one standard deviation long (√λ1); the shorter " +
						"perpendicular arrow is the second principal component (PC2), one √λ2 " +
						"long. Watch PC1 rotate as you change the sliders: increasing stretch alone " +
						"drags PC1 toward the x-axis, because x now carries more of the variance; " +
						"increasing r alone drags it toward the 45° diagonal. The readout reports " +
						"what fraction of the total variance PC1 alone explains — climbing toward " +
						"100% as the cloud gets thinner and more line-like, and dropping toward 50% " +
						"as the cloud rounds out toward a circle (r near 0, stretch near 1).",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Summarize a cloud of correlated points with one number and one direction " +
						"instead of two separate, axis-tied variances: instead of reporting " +
						"'Var(x)=2, Var(y)=1.2' — two numbers that depend on which way you happened " +
						"to draw the axes — report 'this data spreads out about 89.5% along one " +
						"specific direction, and barely at all perpendicular to it,' a fact about " +
						"the cloud itself, not about the axes. That's also literally how PCA " +
						"compresses data with many variables: instead of keeping all of them, keep " +
						"only the first few principal components, whichever few directions capture " +
						"the most variance, and accept a small, quantifiable loss (the leftover " +
						"eigenvalues) in exchange for far fewer numbers per data point.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A face-recognition system built on 'eigenfaces' represents each face image " +
						"(thousands of raw pixel values) as a handful of principal-component " +
						"coordinates instead, because most of the pixel-to-pixel variation across a " +
						"set of faces packs into a small number of directions. A genetics study " +
						"with measurements across thousands of genes uses just the first two or " +
						"three principal components to plot samples on a 2D scatterplot a human can " +
						"actually look at, trusting that those few directions carry most of the " +
						"real signal. Finance runs PCA on many correlated stock returns to find a " +
						"small number of independent 'risk factors' driving the whole market, " +
						"instead of tracking every stock's covariance with every other stock " +
						"individually.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: PC1 is the direction of maximum variance, and its " +
						"eigenvalue tells you exactly how much of the total variance that direction " +
						"accounts for — a number you can quote precisely, like '89.5% of variance " +
						"explained by PC1.'",
					"Not like this: assuming the first principal component always lines up with " +
						"one of your original axes, or that it always means 'the answer' in some " +
						"deeper sense — PC1 is only ever the direction of maximum spread in the " +
						"data you fed it; it has no idea which axis you'd find meaningful. Feed it " +
						"a perfectly round, uncorrelated cloud (r=0, stretch=1) and PC1 and PC2 tie " +
						"at 50% each — the 'first' direction becomes arbitrary, since any direction " +
						"is as good as any other once the cloud has no elongation left to point " +
						"along.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "r", Label: "Target correlation (r)", Min: -1, Max: 1, Step: 0.05, Def: 0.7},
			{Key: "n", Label: "Sample size (n)", Min: 20, Max: 300, Step: 10, Def: 150},
			{Key: "stretch", Label: "x-axis stretch (×)", Min: 1, Max: 4, Step: 0.25, Def: 2},
		},
		Render: render,
	})
}

func render(params map[string]float64) string {
	_ = params
	return viz.New(680, 460, -1, 1, -1, 1).String()
}
