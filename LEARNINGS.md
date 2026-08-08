# Learnings — the tutoring log

Append-only. Every concept the loop builds adds one entry here in the same
commit set. This is the file to read when you want the *lesson* rather than the
code. Newest entries go at the top.

---

## confusion-matrix — the four ways a classifier can be right or wrong
**The idea in one line:** every classification a model makes lands in exactly
one of four buckets — correct-positive, correct-negative, or one of two
distinct kinds of *wrong* — and nearly every classification metric is just
some ratio built from those four counts.

Rows are ground truth (actually positive / actually negative), columns are
the prediction (predicted positive / predicted negative), and the two
diagonals mean very different things: the TP/TN diagonal is where the
classifier agreed with reality, the FP/FN diagonal is where it didn't — and
those two kinds of wrong are not interchangeable. A false positive (FP) is a
false alarm: flagging something that was actually fine. A false negative
(FN) is a miss: letting something real slip through. The picture shades each
cell by its share of the population, so at a glance you can see not just
which mistake is more common but how the four numbers move as the knobs
change: raising the threshold shrinks both "predicted positive" cells (TP
and FP) and grows both "predicted negative" cells (FN and TN), because
fewer examples clear a higher bar.

Accuracy, precision, recall and F1 are all just different ratios over these
same four counts — accuracy asks "what fraction did I get right overall"
(the diagonal over everything), precision asks "of what I flagged, how much
was real" (TP over the predicted-positive column), recall asks "of what was
real, how much did I catch" (TP over the actual-positive row). None of them
is "the" right metric — which one matters depends on whether a false alarm
or a miss costs you more.

**Where it bites in real life:** a 99%-accurate fraud model can still be
useless if fraud is rare — nearly all its correct predictions are the easy
"not fraud" calls, and accuracy alone hides whether it's catching any real
fraud (recall) or drowning investigators in false alarms (precision). Read
the actual 2x2 grid, not just one summary number, before trusting a
classifier's headline accuracy.

---

## roc-auc — the whole trade-off curve, not just one threshold
**The idea in one line:** an ROC curve is what you get by drawing precision
and recall's threshold trade-off (see precision-recall) as a path instead of
a single point — sweep the threshold across every possible setting and plot
false-positive rate against true-positive rate at each one.

Positive- and negative-class scores are modeled as two overlapping normal
distributions, "separation" apart — the same setup as precision-recall.
Instead of picking one threshold and reading off precision and recall, this
picture sweeps the threshold from "call nothing positive" down to "call
everything positive" and traces every (false-positive rate, true-positive
rate) pair along the way. At one extreme the classifier flags nothing —
(0, 0). At the other it flags everyone — (1, 1). Where the curve goes between
those two corners is the whole story: a classifier that can't tell the
classes apart traces the diagonal (at any threshold, whatever fraction of
negatives it catches, it catches the same fraction of positives — no
better than a coin flip), while a classifier with real separation bows the
curve up and to the left, catching positives while dragging along far fewer
false alarms.

AUC — the shaded area under that curve — compresses the entire curve into
one number: it's exactly the probability that a randomly chosen positive
example scores higher than a randomly chosen negative one. That's why AUC
0.5 means random guessing (no separation) and AUC 1.0 means perfect
separation, and why it doesn't require picking a threshold at all — it grades
the ranking a model produces, before any operating point is chosen. Slide
"class separation" and watch the curve peel away from the diagonal as AUC
climbs; slide "threshold" and watch the orange marker slide along the fixed
curve, because changing the threshold moves *where on the curve* you're
operating, not the curve's shape.

**Where it bites in real life:** two models can have the same accuracy at
their default threshold yet very different AUCs — the one with higher AUC
has more headroom no matter where you eventually set the cutoff. It's also
why AUC is popular for comparing models before deployment, but a poor
substitute for precision/recall once you actually have to pick one operating
threshold and live with its false-positive rate.

---

## bayes-theorem — a positive test is only as good as how rare the thing was
**The idea in one line:** Bayes' theorem is just bookkeeping — it takes how
common a condition was *before* the test (the prior) and updates it using how
reliable the test is, and for rare conditions the "before" number dominates
far more than intuition expects.

The picture starts with a population of 100 people at the given base rate
(prior) and runs a test with a given sensitivity (catches true cases) and
specificity (clears true negatives) over all of them. It then regroups
everyone by test result: the "tested positive" strip and the "tested
negative" strip. Inside "tested positive", green squares are true positives
and red squares are false alarms — people who don't have the condition but
tested positive anyway. At the default 1% base rate with a 99%-sensitive,
95%-specific test, the positive group is mostly red: about 5 false alarms
for every true positive, even though the test looks accurate on paper. Raise
the prior and green takes over the positive strip; raise specificity and the
red squares vanish, because specificity controls exactly how many healthy
people get swept up into a positive result.

The reason is arithmetic, not surprise: false positives come from the (huge)
healthy population times a small false-positive rate, while true positives
come from the (tiny) diseased population times a large sensitivity. When the
healthy population dwarfs the diseased one, even a small false-positive rate
applied to it can outnumber the true positives drawn from a small prior.
Bayes' theorem is the formula that adds these two contributions correctly:
P(condition | positive) = (sensitivity × prior) / (sensitivity × prior +
false-positive-rate × (1 - prior)).

**Where it bites in real life:** screening for a rare disease, a rare fraud
pattern, or a rare security alert — a "99% accurate" test sounds like a
near-certainty, but if the thing it's testing for is rare, most positives
are still noise. It's why doctors ask a second, more specific test before
acting on a positive screen, and why "the model flagged it" needs a base
rate attached before anyone should trust the flag.

---

## p-value — how surprising is this, if nothing's really going on?
**The idea in one line:** a p-value is just the area under a "no effect"
curve that's at least as extreme as what you actually observed — nothing
more, and importantly, not "the probability the null hypothesis is true".

Start from the null distribution: what your test statistic would look like
if there were truly no effect, no difference, nothing going on — here, a
standard normal. Run your real experiment and get an observed statistic z.
The p-value asks one narrow question: under the null distribution, how much
of the probability lives at least as far from zero as z, in either
direction? That shaded area is the whole definition. A small p-value means
your observed result would be a rare, surprising draw if the null were true —
which is evidence against the null, but is not itself "the probability the
null is true" (that's a different, more slippery quantity that a single
p-value can't give you).

The significance level α turns that shaded area into a yes/no decision: pick
a threshold in advance (0.05 is conventional, not magic), find the critical
boundary ±z* where the null distribution's tail area equals exactly α (the
dashed lines), and check whether your observed z fell past it. Push z further
from zero and the shaded p-value area shrinks; once it drops below α, the
observed statistic has crossed the dashed boundary and the result counts as
"statistically significant" at that α — a label about the *procedure*
crossing a line you drew beforehand, not a certificate of truth.

**Where it bites in real life:** "p < 0.05" gets read as "there's a 95%
chance the effect is real," which the picture makes clear it isn't — it's a
statement about how the null distribution behaves, built before you saw the
data. It's also why p-hacking (trying many thresholds, subgroups, or metrics
until one clears α by chance) is a real problem: with enough looks, a
surprising-looking shaded area shows up eventually even when there's truly
nothing going on.

---

## confidence-interval — what "95% confident" actually covers
**The idea in one line:** a confidence level doesn't describe one interval —
it describes how often intervals built this way would capture the true value
if you repeated the experiment many times.

It's tempting to read "95% confidence interval" as "there's a 95% chance the
true mean is in this interval," but once an interval is built, the true mean
either is or isn't in it — there's no more randomness left to assign a
probability to. What the 95% actually describes is the *procedure*: draw a
sample, compute its mean, build an interval of mean ± z·(σ/√n). Do that many
times and about 95% of the resulting intervals will contain the true mean,
purely because 95% of sample means land close enough to the true mean for
their interval to reach it.

The picture makes this concrete without leaning on randomness at all: 20
"hypothetical experiments" are placed at evenly spaced quantiles of the
sampling distribution (an exact, reproducible stand-in for repeated random
sampling), and each gets its own confidence interval — green if it captures
the true mean (dashed line), red if it misses. Raise the confidence knob and
every interval widens, turning red rows green; raise the sample size and
every interval narrows around the same mean, because standard error shrinks
as 1/√n. Coverage — how many rows are green — depends only on the confidence
level, never on n: n changes how *tight* the net is, not how *often* it's cast
wide enough to catch the fish.

**Where it bites in real life:** "the poll's margin of error is ±3 points at
95% confidence" means 95% of polls run this way would bracket the true
value — not that this specific poll has a 95% chance of being right. It's
also why a *narrower* interval from a bigger sample is a real improvement
(tighter net) while cherry-picking a higher confidence level just to get a
cleaner-looking one-off result is not (it only pays off in the long run,
across many repeats).

---

## central-limit-theorem — averages go normal, no matter where they started
**The idea in one line:** the distribution of a sample *mean* converges to a
normal curve as the sample size n grows, even when the population you're
averaging from looks nothing like a bell curve.

The population here is deliberately about as far from normal as it gets: an
exponential distribution, sharply peaked at zero with a long right tail (think
"time between rare events," or income data with a few huge outliers). At n=1,
"the sample mean" is just one draw from that population, so the picture is the
skewed exponential itself. But average n independent draws together and look
at *that* number's distribution instead of any single draw's: it tightens
around the true mean and its shape reshapes itself toward a normal curve — and
it keeps doing this for any population with finite variance, not just this one.

The picture is exact rather than simulated: the sum of n iid Exponential(λ)
draws is exactly a Gamma(n, λ) distribution, so the sample mean's distribution
has a closed form (`SampleMeanPDF`) with no random sampling involved. Two
things are worth watching as you drag n up: the peak (thick curve) narrows and
tracks a normal curve of the same mean and variance (thin reference line) more
and more closely, while the mean itself never moves — only the *spread*
shrinks, proportional to 1/√n.

**Where it bites in real life:** it's the reason polling averages, A/B test
metrics, and quality-control sample means can be treated as roughly normal
(and t-tests/z-tests applied) even when the underlying data is skewed or
weird — as long as the sample is big enough. It's also why "n is too small"
is a real objection: with a handful of draws from a skewed population, the
sample mean's own distribution is still skewed, and normal-based confidence
intervals can mislead.

---

## normal-vs-skew — two more knobs on the bell curve
**The idea in one line:** skewness and kurtosis are the next two "shape"
numbers after mean and variance, and turning them reshapes a normal curve
into something lopsided or spiky-with-fat-tails.

The mean tells you where a distribution is centered and the variance tells
you how wide it is, but two distributions can share both of those and still
look completely different. **Skewness** measures lopsidedness — positive skew
means a longer right tail (a few big values dragging out to the right, like
income), negative skew means a longer left tail. **Kurtosis** (specifically
*excess* kurtosis, measured relative to the normal's baseline) measures how
"peaked-and-fat-tailed" versus "flat-and-thin-tailed" a distribution is: a
positive value means more of the probability sits in a sharp central peak
*and* in rare extreme outliers, at the expense of the shoulders — real-world
returns and error distributions are almost always like this ("fat tails"
mean the rare disaster is more likely than a plain normal curve would suggest).

The picture builds both from a normal curve using a closed-form correction
(the Gram-Charlier series): skew scales an odd cubic-shaped term and kurtosis
scales an even quartic-shaped term, and at skew=kurt=0 the curve is exactly
the standard normal — no approximation, no sampling.

**Where it bites in real life:** financial returns and risk models (fat
tails mean "once in a century" events happen far more than a normal-curve
model predicts — a lesson learned expensively in 2008), A/B test metrics that
aren't actually normal (skew breaks assumptions behind some significance
tests), and any dashboard that reports only a mean and stddev for data that
is secretly lopsided or spiky.

---

## variance-vs-stddev — why we square, then square-root back
**The idea in one line:** raw deviations from the mean always cancel to zero,
so we square them to measure spread — then square-root the result to get back
to the data's own units.

Take any dataset's deviations from its mean and add them up: by definition the
positives and negatives cancel exactly, giving zero every time. That makes a
raw average deviation useless as a spread measure. Squaring each deviation
fixes the cancellation (every term becomes positive) and, as a side effect,
punishes far-out points much harder than close ones — a point twice as far
from the mean contributes four times the squared deviation. The average of
those squared deviations is the **variance**. Its only flaw: squaring the data
also squares the units (dollars become dollars², seconds become seconds²), so
variance isn't directly comparable to the data. Taking its square root — the
**standard deviation** — undoes that and lands back in the original units.

The picture makes both effects visible: scaling every point away from the
mean by a factor *k* scales variance by *k²* but stddev only by *k* (drag
spread and watch the bars shoot up much faster than the eye expects), and
pushing one point out further as an "outlier" inflates variance disproportionately
— the whole reason variance/stddev are sensitive to outliers while the median
and IQR are not.

**Where it bites in real life:** why one wild outlier can wreck a
variance-based estimate (use median/IQR for robustness instead), why RMSE
(root-mean-squared-error) is reported in the target's own units instead of
raw MSE, and why standard deviation — not variance — is the number quoted
alongside a mean.

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
