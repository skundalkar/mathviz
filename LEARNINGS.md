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

## mean-median-mode — three "typical" values that agree until they don't
**The idea in one line:** mean, median, and mode all answer "what's typical
here," but they answer it differently enough that a skewed dataset pulls
them apart — and knowing which one you're looking at matters.

Bill Gates walks into a small bar with 20 regulars, each with a net worth
around $80,000. Before he sits down, the bar's *average* (mean) net worth is
about $80,000 — a fair summary of "the typical person here." The moment he
sits down, that mean net worth rockets into the billions. Nothing about the
20 regulars changed — not one of them got richer or poorer — but "the
average person in this bar" now sounds like a billionaire. That's obviously
a lie about what's typical, and it's the mean's fault: it's an average of
*every* value, including the one wildly extreme one, so a single outlier can
drag it anywhere.

The median doesn't have this problem: it's just "the middle value when
everyone's lined up in order." With Bill Gates added, the middle of a line
of 21 people barely shifts — the median net worth is still roughly $80,000,
because Bill Gates is just one more person standing at the far end of the
line, not a value that gets *averaged in*. The mode — the single most common
value — doesn't move at all either.

For a symmetric distribution these three "typical value" measures all
coincide, but skew the underlying distribution — as the picture does,
dragging a log-normal curve's tail out to the right — and they peel apart in
a fixed order: **mode < median < mean**. The mode sits at the peak (most
likely single outcome), the median splits the population exactly in half,
and the mean gets pulled hardest toward whichever direction the long tail
stretches, because it's the only one of the three that actually does
arithmetic with the extreme values instead of just counting or ranking them.

**Where it bites in real life:** "average household income" headlines are
almost always higher than what a typical household actually earns, because a
small number of very high earners drag the mean up while the median (what
statisticians usually mean by "typical") barely moves — exactly why income
and home-price statistics are usually reported as medians, and why a
company's "average" salary can look great in a press release while most
employees make noticeably less.

---

## correlation — moving together is not the same as causing
**The idea in one line:** r measures how tightly two variables move
together — from -1 (perfect opposite) through 0 (no linear relationship) to
+1 (perfect together) — and moving together is a fact about the data, not an
explanation for it.

Every summer, ice cream sales climb. Every summer, drowning deaths climb too.
Plot one against the other across the year and you'll find a strong positive
correlation — high r, tight cloud. It's tempting to read that as "ice cream
causes drowning" (or, mixed up the other way, "drowning makes people crave
ice cream"). Neither is true. Both are being dragged along by a third thing
entirely: hot weather. Heat means more people swim (more drowning risk) *and*
more people buy ice cream — the two variables never touch each other
causally, they're just both downstream of summer.

That's the trap r sets by design: it only ever measures how tightly a
scatter of points hugs a straight line. Both axes here are standardized, so
the trend line r implies always runs through the origin with slope r — why
the orange guide line visibly steepens as r moves toward ±1. At the extremes
there's no noise left: every point sits exactly on the line. But "sits
exactly on a line" and "one causes the other" are different claims, and the
math only ever makes the first one. A hidden third variable (like summer
heat), reverse causation, or even pure coincidence in a small sample can all
produce a strong r with zero causal link behind it — the number can't
distinguish between them; only more information (a controlled experiment, a
mechanism, timing) can.

The picture also shows sample r next to target r: with a small n they
visibly drift apart, a reminder that a correlation measured from one sample
is itself just an estimate of the true relationship, not the relationship.

**Where it bites in real life:** "cities with more police have more crime"
(both driven by population density), a stock-picking strategy that
correlates with past returns purely by chance, or a feature in an ML model
that correlates with the label without causing it — swap it out in
production and the correlation, having no causal footing, quietly
disappears.

---

## precision-recall — two ways to be wrong, and a knob that trades between them
**The idea in one line:** a classifier can fail in two different ways —
flagging things that aren't real (false alarms) or missing things that are
(misses) — and no single threshold setting minimizes both at once.

You're building a spam filter. The obvious plan: make it aggressive, flag
anything suspicious, catch every last piece of spam. Turn that dial all the
way up and you *do* catch 100% of the spam — recall is perfect — but you also
start flagging your boss's emails, meeting invites, and password-reset links,
because "suspicious" was cast too wide. Loosen the dial to stop burying real
mail and now actual spam slips through the net untouched. There's no setting
of this one dial that avoids both problems — that's not a bug in your filter,
it's structural.

The picture makes the structure visible: real negatives (legitimate email)
and real positives (spam) are two overlapping bell curves, and "flag as
spam" means "score above the threshold." **Recall** asks, of everything that
really was spam, how much did you catch — TP/(TP+FN). **Precision** asks, of
everything you flagged, how much was actually spam — TP/(TP+FP). Slide the
threshold right and you only flag the most obvious spam: precision climbs
(what you flag is almost always right) but recall falls (borderline spam
slips through). Slide it left and the opposite happens. **F1**, their
harmonic mean, only rises when *both* rise together — which the threshold
alone can't deliver.

The one thing that genuinely improves both is separating the two curves
further apart: a better model that scores real spam and real mail more
differently in the first place, leaving more room between "clearly not
spam" and "clearly spam" for a threshold to land cleanly.

**Where it bites in real life:** cancer screening (a false negative means
missed disease — worth tolerating more false positives to avoid), spam
filters (a false positive means a lost important email), fraud detection,
search ranking — anywhere "how sensitive should this be" is really a
business decision about which mistake costs more, dressed up as a slider.

---

## logarithms — the question "how many doublings does it take?"
**The idea in one line:** a logarithm answers the reverse of multiplication —
not "what do I get if I multiply b by itself n times," but "how many times do
I need to multiply by b to reach x."

Fold a 0.1mm-thick sheet of paper in half, over and over. Each fold doesn't
add a fixed amount — it doubles whatever thickness you already have: 0.1 →
0.2 → 0.4 → 0.8mm, and so on. How many folds until the stack passes a
100-meter building?

The tempting shortcut is division: convert 100m to 100,000mm, divide by the
0.1mm sheet, get 1,000,000. Done — right? That's actually the answer to a
*different* problem: stacking a million loose, unfolded sheets on top of each
other, where each new sheet just adds another fixed 0.1mm. That process is
additive — n sheets give you n × 0.1mm — and plain division solves it fine.

Folding is not that. Each fold *multiplies* your current thickness by 2, so
after n folds you have 0.1 × 2ⁿ mm, and the question becomes 0.1 × 2ⁿ =
100,000. Try to isolate n the way you would in the stacking problem and you
get stuck — n isn't multiplying anything here, it's counting how many times 2
gets multiplied by itself, buried up in an exponent. Division has no move for
that. The one thing that can reach into an exponent and pull it back down to
an ordinary number is a logarithm: n = log₂(1,000,000) ≈ 20. Twenty folds
passes the building — not a million sheets, twenty folds — because doubling
compounds.

That gap between "20 folds" and "1,000,000 sheets" for the *same* size
increase is exactly why logarithms exist, and why they show up wherever
growth is multiplicative rather than additive. The picture makes the
mechanism visible: mark the powers of the base (b⁰, b¹, b², b³, …) on the
curve, and while the x-values are multiplying by b each step, the log values
are only ever stepping up by exactly **+1**. Multiplying the input becomes
adding to the output — that's the whole trick, and it's why logarithms were
invented centuries before calculators (to turn slow multiplication into fast
addition), and why decibels, pH, star magnitudes, and earthquake scales are
all logarithmic — each compresses an enormous multiplicative range into a
scale humans can actually read.

**Handy anchors:** log₁₀(1000) = 3, log₂(8) = 3, log_b(1) = 0 for any base,
and log_b(b) = 1.

**Where it bites in real life:** "this investment compounds annually" or
"this outbreak doubles every few days" are folding stories, not stacking
stories — reach for a logarithm, not division, to get from a growth rate to
"how long until X." It's also why a Richter magnitude of 7 isn't "a bit more"
than 6 — it's about 32x more energy, because each whole step up the scale is
another multiplication, not another addition. And it's the reason the
instinct to just divide is worth pausing on: division answers "how many
times do I add this," a logarithm answers "how many times do I multiply by
this" — different questions that happen to look similar until you write them
out.

---

## standard-deviation — a ruler that means the same thing everywhere
**The idea in one line:** "10 points above average" tells you nothing by
itself — standard deviation is what turns a raw distance from average into a
number you can actually compare across totally different situations.

Two students each score 10 points above the class average on their tests.
Student A: sounds impressive — until you learn everyone in class A scored
within 2 points of the mean, so being 10 points up puts them off the charts,
practically alone at the top. Student B scored 10 points up too, but that
class's scores were scattered across a 40-point range — 10 points above
average there barely clears the middle of the pack. Same raw number, wildly
different meaning, because "10 points" only means something relative to how
spread out everyone else is.

Standard deviation (σ) is that spread, and it comes with a guarantee raw
point-differences never do: for a bell-shaped distribution, the band from
mean − 1σ to mean + 1σ always contains about **68%** of everyone — not
"usually," always, whether σ is tiny (a tightly clustered class) or huge (an
all-over-the-place one). That fixed 68% is what makes "z standard deviations
above average" a portable ruler: student A's +10 might be +5σ (astonishingly
rare), student B's +10 might be +0.5σ (unremarkable) — the same raw gap,
translated into a number that actually means the same rarity everywhere it's
used. (Standard deviation itself is just the square root of variance — see
variance-vs-stddev for why we square first, then square-root back.)

**Where it bites in real life:** a 101°F fever is routine for some people and
alarming for others depending on their normal baseline and its variability;
manufacturing tolerances are set in sigmas, not raw units, because "off by
2mm" means something different for a bolt than for a bridge; and "1-in-20"
scientific thresholds (p < 0.05) are really a statement about how many sigmas
out a result has to land before it counts as surprising.
