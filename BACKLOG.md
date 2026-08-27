# Backlog — the queue the build loop pulls from

Each build cycle picks the **first unchecked** item, implements it as one new
concept package (following the pattern in `internal/concepts/stddev`), and
checks it off in the same set of commits. Keep this list ordered by what you
want to learn next — reorder freely.

Format: `- [ ] <id> — <one-line description of the concept to visualize>`

## Done
- [x] standard-deviation — σ as the width of the ±1σ band that always holds ~68%
- [x] logarithms — log_b(x) as "b to what power is x", and ×b → +1
- [x] precision-recall — the threshold trade-off between precision and recall
- [x] correlation — scatter clouds from r = -1 to +1, and why r≠causation
- [x] mean-median-mode — how a skewed sample pulls the three apart
- [x] variance-vs-stddev — why we square, then square-root back
- [x] normal-vs-skew — skewness and kurtosis reshaping a distribution
- [x] central-limit-theorem — sample means going normal as n grows
- [x] confidence-interval — what "95% confident" actually covers
- [x] p-value — a null distribution and where the observed statistic falls
- [x] bayes-theorem — base rates and why rare-disease positives mislead
- [x] roc-auc — sweeping the threshold to trace an ROC curve
- [x] confusion-matrix — TP/FP/FN/TN as an interactive grid
- [x] overfitting — a wiggly curve chasing noise vs. a smooth fit
- [x] gradient-descent — a ball rolling downhill with an adjustable learning rate
- [x] entropy — how surprise/information changes with probability
- [x] pr-auc — sweeping the threshold to trace a precision-recall curve, and why it beats ROC-AUC on imbalanced classes
- [x] sigmoid-softmax — squashing logits into probabilities
- [x] exponential-growth — doubling time and why it sneaks up on you
- [x] eval-playbook — reading loss, ROC-AUC, PR-AUC, precision and recall together, in order, as a fixed reference table (no sliders — a synthesis of the six classifier concepts above it)
- [x] calibration — reliability diagrams, Expected Calibration Error, and temperature scaling — why a good ranking (AUC) doesn't mean a trustworthy confidence number
- [x] derivative — slope of a tangent line as you shrink the interval
- [x] integral — area under a curve as a sum of shrinking slabs
- [x] sine-cosine — the unit circle unrolled into waves
- [x] vectors — addition, dot product, and projection
- [x] complex-numbers — rotation in the complex plane
- [x] prime-sieve — the Sieve of Eratosthenes animated
- [x] modular-arithmetic — the clock that makes cryptography work
- [x] z-score — standardizing a value by how many standard deviations it sits from the mean, so scores from different scales become comparable
- [x] linear-regression — fitting the least-squares line through a scatter and reading the residuals as what it doesn't explain
- [x] naive-bayes — classifying by combining per-feature likelihoods under Bayes' theorem's independence assumption
- [x] k-means-clustering — partitioning points into clusters by iteratively updating centroids
- [x] cosine-similarity — comparing vector directions instead of magnitudes, the backbone of embedding search
- [x] pascals-triangle — binomial coefficients built row by row, and why each entry is "n choose k"
- [x] eigenvectors-eigenvalues — the directions a linear transformation only stretches, never rotates
- [x] birthday-paradox — why shared-birthday collisions happen far sooner than intuition expects
- [x] law-of-large-numbers — why a sample average settles down toward the true mean as n grows, and how noisy it still is at small n
- [x] monte-carlo-estimation — estimating a number you can't solve for directly (like π) by sampling randomly and watching the estimate converge
- [x] fourier-series — approximating a wave (even a square wave) by summing sine waves of increasing frequency
- [x] markov-chains — transition probabilities and the steady-state distribution a random walk settles into
- [x] binomial-distribution — the probability of exactly k successes in n independent trials, built from pascals-triangle's n-choose-k
- [x] poisson-distribution — modeling rare, independent events over a fixed window with a single rate λ, as the limit of binomial when n is huge and p is tiny
- [x] covariance — how two variables move together in raw units, and why dividing by both standard deviations turns it into correlation
- [x] taylor-series — approximating a curve near a point using a polynomial built from its derivatives
- [x] permutations-vs-combinations — counting arrangements when order matters vs. when it doesn't
- [x] random-walk — a 1D random walk's expected position stays at 0 while its spread grows like √n
- [x] principal-component-analysis — finding the directions of maximum variance in a cloud of points, built from covariance and eigenvectors
- [x] kl-divergence — measuring how far one probability distribution is from another, built from entropy
- [x] bias-variance-tradeoff — decomposing prediction error into bias, variance, and irreducible noise as model complexity changes
- [x] matrix-multiplication — combining two linear transformations into one, and reading a matrix as what it does to the basis vectors
- [x] newtons-method — finding a function's root by repeatedly following its tangent line
- [x] chain-rule — how the derivative of a composed function multiplies the derivatives of its parts
- [x] simpsons-paradox — a trend that reverses when subgroups are combined, even though one option wins every subgroup on its own
- [x] bootstrap-resampling — estimating a statistic's uncertainty by resampling the data itself instead of assuming a formula/distribution
- [x] decision-trees — splitting a dataset by information gain (built from entropy) into a tree of yes/no questions
- [x] regularization — how an L1/L2 penalty shrinks coefficients to fight overfitting, built from linear-regression and bias-variance-tradeoff
- [x] chi-squared-test — testing whether observed category counts differ from what independence would predict
- [x] determinant — a matrix's determinant as the area/volume scaling factor of its linear transformation, built from matrix-multiplication
- [x] exponential-distribution — modeling the waiting time between rare events, the continuous counterpart to poisson-distribution
- [x] quadratic-formula — completing the square and reading a parabola's discriminant as how many real roots it has
- [x] fibonacci-golden-ratio — the Fibonacci sequence's ratio of consecutive terms converging to φ
- [x] t-distribution — how estimating σ from a small sample fattens the normal's tails, and why small-sample confidence intervals need a wider critical value than z
- [x] logistic-regression — fitting a decision boundary by modeling class probability directly with sigmoid(wx+b) via gradient descent on log-loss, built from sigmoid-softmax and gradient-descent
- [x] euclidean-algorithm — finding the greatest common divisor by repeated remainders, and why it terminates in so few steps, built from modular-arithmetic
- [x] singular-value-decomposition — decomposing any matrix into rotate-scale-rotate (UΣVᵀ), generalizing eigendecomposition to non-square matrices, built from eigenvectors-eigenvalues and matrix-multiplication
- [x] mutual-information — how much knowing one variable reduces uncertainty about another, built from entropy and kl-divergence
- [x] anova — comparing means across three or more groups by partitioning total variance into between-group and within-group components, built from variance-vs-stddev
- [x] cross-entropy-loss — the loss function that scores a predicted probability against a true label, and why minimizing it is the same as minimizing KL-divergence to the true distribution, built from entropy and kl-divergence
- [x] random-forest — bagging many decision trees on bootstrap-resampled data and averaging their votes to cut variance without adding bias, built from decision-trees and bootstrap-resampling
- [x] pagerank — ranking nodes in a graph by the steady-state of a random walk over links, built from markov-chains
- [x] k-nearest-neighbors — classifying a new point by majority vote of its k closest neighbors, and how the choice of k trades bias for variance
- [x] hypergeometric-distribution — the probability of k successes when sampling without replacement, and how it differs from binomial as the population shrinks

## Up next (statistics & ML intuition)

## Up next (math foundations)
- [ ] gaussian-elimination — solving a linear system by row-reducing to echelon form, and reading a matrix's rank off the result
- [ ] lagrange-multipliers — finding a function's constrained max/min by matching gradients with the constraint's gradient, built from chain-rule
- [ ] rsa-encryption — public-key encryption built from modular arithmetic and the euclidean algorithm

## Up next (statistics & ML intuition, batch 2)
- [ ] support-vector-machine — finding the maximum-margin hyperplane that separates two classes
- [ ] gradient-boosting — sequentially fitting trees to the residual errors of the ensemble so far, built from decision-trees

## How to add work
Append an unchecked line under the right section by hand any time — the loop
never runs dry as long as there is one unchecked box. The daily build trigger
also tops this up itself: before building, it checks the unchecked count and,
if it's running low, proposes a few new math/stats/ML items (checked against
`## Done` and the rest of this file for duplicates) and commits them on their
own before touching any concept work.
