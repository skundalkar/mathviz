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

## Up next (statistics & ML intuition)
- [ ] gradient-descent — a ball rolling downhill with an adjustable learning rate
- [ ] entropy — how surprise/information changes with probability
- [ ] sigmoid-softmax — squashing logits into probabilities

## Up next (math foundations)
- [ ] exponential-growth — doubling time and why it sneaks up on you
- [ ] derivative — slope of a tangent line as you shrink the interval
- [ ] integral — area under a curve as a sum of shrinking slabs
- [ ] sine-cosine — the unit circle unrolled into waves
- [ ] vectors — addition, dot product, and projection
- [ ] complex-numbers — rotation in the complex plane
- [ ] prime-sieve — the Sieve of Eratosthenes animated
- [ ] modular-arithmetic — the clock that makes cryptography work

## How to add work
Append an unchecked line under the right section. The loop never runs dry as
long as there is one unchecked box.
