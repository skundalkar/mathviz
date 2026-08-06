# Learnings — the tutoring log

Append-only. Every concept the loop builds adds one entry here in the same
commit set. This is the file to read when you want the *lesson* rather than the
code. Newest entries go at the top.

---

## precision-recall — the threshold trade-off
**The idea in one line:** Precision and recall measure two different failure
modes, and the classification threshold trades one for the other.

Picture two overlapping bell curves: real negatives on the left, real positives
on the right. Your model scores each example, and you call everything above a
threshold "positive."

- **Recall** = of all the *real positives*, how many did you catch? = TP / (TP + FN)
- **Precision** = of everything you *flagged*, how many were right? = TP / (TP + FP)

Slide the threshold right: you only flag the most confident cases, so precision
rises — but you miss more real positives, so recall falls. Slide it left: you
catch almost everything (high recall) but drag in false alarms (low precision).
**F1** is their harmonic mean, high only when *both* are high. The one way to
improve both at once is to separate the classes better — i.e. a better model,
not just a better threshold.

**Where it bites in real life:** spam filters (false positives = lost email),
cancer screening (false negatives = missed disease), search ranking, fraud
detection. Every one is a choice about which mistake is cheaper.

---

## logarithms — turning multiplication into addition
**The idea in one line:** log_b(x) is the exponent — "b to what power gives x?"

Mark the powers of the base on the curve: b⁰=1, b¹, b², b³. The x-values
*multiply* by b each step, but the log values only step up by **+1**. That is
the whole magic: logs convert multiplying into adding. It's why slide rules
worked, why decibels, pH, star magnitudes, and earthquake scales are all
logarithmic (they compress enormous ranges), and why a "log scale" chart turns
exponential growth into a straight line.

**Handy anchors:** log₁₀(1000) = 3, log₂(8) = 3, log_b(1) = 0 for any base,
and log_b(b) = 1.

---

## standard-deviation — a universal ruler for spread
**The idea in one line:** σ is the width of the "typical" band around the mean.

Drag σ and the bell curve gets fatter or skinnier, but the shaded mean ± 1σ band
*always* holds about **68%** of the data (95% within 2σ, 99.7% within 3σ — the
"68–95–99.7 rule"). That fixed coverage is what makes σ a portable ruler: saying
something is "2 sigma out" means the same rarity whether you're measuring heights
or test scores. Standard deviation is just the square root of the variance (the
average squared distance from the mean) — we square to make distances positive,
then square-root to get back to the original units.

**Where it bites in real life:** z-scores, control charts in manufacturing,
"1-in-20" thresholds in science, risk/volatility in finance.
