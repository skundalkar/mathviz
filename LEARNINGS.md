# Learnings — the tutoring log

Append-only. Every concept the loop builds adds one entry here in the same
commit set. This is the file to read when you want the *lesson* rather than the
code. Newest entries go at the top.

---

## mean-median-mode — three averages, one skew
**The idea in one line:** for a symmetric distribution, mean, median, and mode
all coincide; skew the distribution and they peel apart in a fixed order.

The picture uses a log-normal curve because all three statistics have exact
closed forms in terms of its two parameters (μ, σ): median = e^μ, mean =
e^(μ+σ²/2), mode = e^(μ-σ²). No sampling needed — the lesson is exact. Drag σ
up from near-zero and watch the three vertical lines split apart from a single
point into **mode < median < mean**, always in that order for a right-skewed
curve. The mode marks the peak (most likely single value), the median splits
the area in half, and the mean is pulled hardest because it's an average of
*every* value including the ones way out in the long tail — a few big values
drag it further than they drag the median.

This is also why "average" is an ambiguous word: for skewed real-world
quantities (income, home prices, response times, city populations) the mean is
almost always higher than the median, sometimes dramatically. A "mean income"
headline can be pulled way up by a handful of very high earners, while the
median (what a typical person actually makes) barely moves — which is why
income statistics are usually reported as medians.

**Where it bites in real life:** household income and wealth reporting, home
prices, response-time/latency dashboards (mean is dragged by rare slow
requests — p50 vs. mean tells a different story), city/company size
distributions, and any "average" claim about data you suspect is skewed.

---

## correlation — how tightly the dots hug a line
**The idea in one line:** r measures how tightly a scatter of points hugs a
straight line, from -1 (perfect downhill) through 0 (no linear pattern) to +1
(perfect uphill).

Both axes are standardized (mean 0, variance 1), so the trend line r implies
always passes through the origin with slope r — that's why the orange guide
line visibly steepens as you drag r toward ±1 and flattens toward 0. At the
extremes (r = ±1) there's no noise left at all: every point falls exactly on
the line. In between, the cloud is a mix of "along the line" and independent
noise, weighted by r and √(1-r²) respectively.

The picture also shows the sample r alongside the target r — with a small `n`
they can drift apart noticeably, the same "sample vs. population" lesson as
standard deviation's ±1σ band: a statistic estimated from a sample is not the
same thing as the true population value. **Correlation is not causation** is
the other half of the lesson: r only says two variables move together, it says
nothing about mechanism — a third variable, reverse causation, or pure
coincidence can all produce a strong r.

**Where it bites in real life:** spurious correlations (ice cream sales and
drowning both track summer heat), feature selection in ML (a feature can
correlate with the label without being causal), and any headline that reports
a correlation as if it were a proven cause.

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
