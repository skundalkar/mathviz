# Learnings — the tutoring log

Append-only. Every concept the loop builds adds one entry here in the same
commit set. This is the file to read when you want the *lesson* rather than the
code. Newest entries go at the top.

---

## exponential-growth — why compounding looks boring right up until it isn't
**The idea in one line:** exponential growth compounds off of whatever's
already there, so a fixed *rate* produces an accelerating *amount* — the
curve looks tame for a long stretch and then rockets, catching straight-line
intuition off guard every time.

The lily-pad puzzle makes the trap concrete: a patch that doubles every day
covers a whole pond on day 30. What day is it half-covered? The honest
answer is day 29 — one day before full, because doubling from half to whole
is exactly one more doubling, full stop. Nothing about "halfway through 30
days" enters into it. For roughly the first 25 of those 30 days the pond
looks nearly empty (under 10% covered) — plausible to write off as "not
much happening" — then the last handful of days it visibly explodes,
because the same percentage growth applied to a much bigger base produces a
much bigger absolute jump. Nothing about the *rate* changed; only the base
it's compounding off of did.

**What the knobs show:** the growth-rate slider is the entire story —
doubling time = ln(2)/ln(1+rate) falls as the rate climbs
(`TestDoublingTimeShrinksAsRateGrows`), so a small rate bump buys a
surprisingly large speedup in how often the value doubles. The periods
slider just widens the window so slower rates get enough runway to show the
same eventual pattern a fast rate shows quickly. The thin gray line is the
straight-line guess sharing the curve's exact starting slope — near t=0 the
two are nearly indistinguishable (`TestLinearMatchesValueNearZero`), and
that overlap is exactly why "it's growing about the same amount each
period" feels true right up until the exponential curve visibly pulls away.
The orange dots mark each doubling; watch how tightly they bunch together
at high rates and how far apart they spread at low ones.

**Where it matters:** compound interest, population growth, viral spread,
Moore's-law-style tech scaling, unchecked technical debt — anywhere a
quantity's growth is proportional to its current size rather than a fixed
amount per period. The practical trap this concept is built around: judging
an exponential trend by its early, unremarkable-looking segment (like the
first 25 days of the pond) systematically underestimates how close the
"sudden" acceleration actually is.

**A design choice worth flagging:** `Rule70` (70 ÷ rate) is included as a
mental-math cross-check next to the exact `DoublingTime`, since it's the
shortcut most people reach for outside a calculator. `TestRule70ApproximatesDoublingTime`
only asserts it stays within 10% for modest rates (≤10%); the concept's
slider goes up to 50%, where the approximation visibly drifts further off
from the exact value — shown side by side in the render rather than hidden,
since watching a mental-math shortcut degrade is itself part of the lesson.

---

## sigmoid-softmax — turning raw scores into numbers that behave like probabilities
**The idea in one line:** sigmoid and softmax are the same operation at
different scales — sigmoid squashes one logit into a single probability for
two classes, softmax squashes any number of logits into a full probability
distribution that sums to 1, and running softmax on two logits recovers
sigmoid exactly.

A model's last layer hands back logits — raw, unbounded real numbers like
2.3 or -0.5 — not probabilities. Nothing stops a logit from being 1000 or
-1000; it's just a score, bigger meaning "more likely" and nothing more.
Sigmoid is the two-class fix: sigmoid(z) = 1/(1+e^-z) takes any real z and
returns something in (0,1), crossing exactly 0.5 at z=0 (a coin-flip logit),
and flattening out — saturating — toward 0 or 1 as z runs to the extremes.
That flattening is worth noticing on the curve: past about z=±5 the output
barely moves even though the logit is still changing a lot, which is exactly
why very confident (very large-magnitude) logits stop giving useful gradient
signal during training.

Softmax generalizes the same trick to any number of classes at once. Given
logits (2, 0.5, -1) for cat/dog/fox: exponentiate each one (e^2≈7.39,
e^0.5≈1.65, e^-1≈0.37), sum them (≈9.41), divide each by that sum —
probabilities ≈78%, 18%, 4%, adding to 100%. Every logit pulls probability
mass away from every other one; that's the mechanism that makes it a genuine
probability distribution over classes rather than three independent
squashed numbers. And the connection to sigmoid isn't just an analogy — it's
exact: run softmax on (z, 0) and its first output is algebraically
1/(1+e^-z), the same formula as `Sigmoid(z)`. `TestSoftmaxTwoClassMatchesSigmoid`
checks exactly this.

**What the knobs show:** the "Logit A" slider drives both panels with the
same z, so you can watch the sigmoid curve's point and class A's softmax bar
move together — the shared logit is the thread connecting the two pictures.
"Logit B" only affects the softmax panel, giving a third class (C, pinned at
logit 0 as a fixed reference) to compete against. Temperature divides every
logit before exponentiating: push it below 1 and whichever class has the
highest logit gets pulled toward 100% (the model "commits" harder than its
raw logit gap would suggest); push it above 1 and even a clear logit lead
gets pulled back toward a uniform 33/33/33% (the model hedges even when it
has a real signal). Temperature is exactly the knob systems expose as "more
creative / more conservative" sampling in language models — same math, just
applied to a much bigger softmax.

**Where it matters:** this squashing step is the last mile of essentially
every classifier — binary logistic regression, multi-class neural nets,
attention weights, next-token prediction in language models. The failure
mode worth internalizing: softmax outputs always sum to 1 and are always
positive, so a model that's never seen anything like the current input
still has to hand back some full distribution — softmax can't express
"I have no idea," only "here's how my confidence splits across the classes
I know about." Low-confidence-looking outputs (probabilities close to
uniform) are the closest softmax gets to saying that.

**A scope choice worth flagging:** the softmax panel fixes class C's logit
at 0 rather than exposing a third slider — three independent logit sliders
would clutter the interaction without adding a new idea, since softmax is
shift-invariant (`TestSoftmaxShiftInvariant`) and only the *differences*
between logits matter. Pinning one class as a zero reference keeps the two
sliders' effects legible while still showing genuine three-way competition.

---

## pr-auc — the metric that doesn't get fooled by imbalance
**The idea in one line:** sweep the classifier's threshold from strictest to
loosest, plot (recall, precision) at every setting, and the area under that
curve — PR-AUC — grades the model across all thresholds at once, without
ROC-AUC's blind spot for rare positives.

Take the spam filter from the precision/recall lesson, but stretch it to a
more realistic imbalance: 1,000 emails, only 20 really spam. A "flag ≥ t"
threshold has two very different report cards depending on which curve you
read. ROC's axes are false-positive rate and true-positive rate — both are
computed *within* a class (FPR out of the 980 real hams, TPR out of the 20
real spams), so a threshold that wrongly flags 50 of those 980 hams only
nudges FPR to 50/980 ≈ 5%, and the ROC curve barely notices. But look at
precision instead: out of everything flagged (18 real spam caught + 50 false
alarms = 68 flagged), precision is only 18/68 ≈ 26% — the inbox is now
*mostly* wrong flags, and PR-AUC feels that pain directly because precision
is computed against everything flagged, not diluted by the 980 easy true
negatives sitting in the background.

**What the knobs show:** the threshold slider sweeps the same point across
the curve that the precision/recall lesson fixed in place — dragging it
right (stricter) climbs toward the top-left (high precision, low recall);
dragging it left (looser) slides down toward the bottom-right. The
separation slider makes the two classes' score distributions easier or
harder to tell apart; more separation bows the curve up and to the right and
raises PR-AUC, because a threshold can now catch more real positives without
also catching more negatives. The flat line at y=0.5 (with equal-sized
classes here) is the floor: a classifier that ranks completely at random
still lands its flags right, on average, exactly as often as positives occur
in the data.

**Where it matters:** anywhere positives are rare — fraud detection, disease
screening, defect detection, security alerts. In those settings ROC-AUC can
report 0.95+ while the model is still unusable in practice, because it's
graded on a sea of easy true negatives it isn't even trying hard to get
right. PR-AUC is the metric that keeps the "how many of my alerts are real"
question front and center, which is usually the question a human actually
cares about.

**A design choice worth flagging:** PR-AUC here is computed with a simple
trapezoid rule over the sampled curve (`TrapezoidalPRAUC`), the same
technique `roc-auc` uses for its curve. This is a close cousin of "average
precision," which some libraries compute slightly differently (precision
interpolated only at the points where a positive is added, avoiding
double-counting jagged detail). On the smooth curve two Gaussian classes
produce, trapezoidal integration tracks the true area well; on jagged
real-world curves it can run a touch optimistic. Good enough for the
intuition this concept is teaching, worth knowing if you reach for the exact
number in a real evaluation.

---

## entropy — what people actually mean when they say the word
**The idea in one line:** "entropy" is just a fancy word for how mixed-up,
spread-out, or hard-to-call something is — high entropy means "could go a
lot of different ways, no clear favorite," low entropy means "pretty much
locked in, one obvious answer."

Forget the math entirely for a second. The word started in physics: a messy
room has high entropy, a tidy one has low entropy. Left alone, a tidy room
drifts toward messy — never the other way around — unless someone actively
puts in effort to re-organize it. Hot coffee left on a counter cools down
and its heat spreads into the room until everything's the same temperature;
it never spontaneously un-mixes back into "hot coffee, cool room." That
drift — from concentrated and orderly toward spread-out and mixed — is what
"entropy increasing" originally meant, and it's the intuition every other
use of the word is quietly leaning on, even when no physics is involved.

**Outside physics, the word got borrowed to describe the same shape of
thing: are the possibilities concentrated (predictable) or spread out
(you-can't-really-say)?** That's it — that's the whole transplant. Nothing
mathematical has to be happening for the word to fit; you're just describing
whether something leans one clear way or is genuinely up for grabs.

**Sentences where someone would actually say it, and what they mean:**

- *"There's a lot of entropy in this org chart right now."* — Nobody's
  really sure who owns what; responsibilities are scattered, not
  concentrated on clear owners. Not "bad," just unsettled and hard to
  predict who you'd even ask.
- *"This password has high entropy."* — There's no discernible pattern to
  exploit; a guesser has no shortcut, every character genuinely could've
  been anything. A low-entropy password ("password123") is the opposite —
  heavily concentrated on a small, guessable set of likely choices.
- *"Left alone, codebases/garages/inboxes tend toward entropy."* — This is
  the physics metaphor directly: things drift toward disorganized unless
  someone spends effort keeping them organized. Order doesn't happen by
  itself; mess is the default direction things drift.
- *"The market's been really entropic this week."* — Price moves aren't
  following any readable trend; outcomes are scattered rather than leaning
  one way, so predicting the next move is close to a coin flip.
- *"The model's predictions have high entropy here."* — Someone building AI
  saying the model is genuinely unsure on this input — its confidence is
  spread thinly across several possible answers instead of piled onto one.
  Low entropy there would mean the model is confidently committed to a
  single answer (right or wrong).

**When the word doesn't fit:** if something is just wrong, or leans hard in
one predictable direction, that's actually the *opposite* of entropy — it's
low entropy, just low entropy pointed at the wrong answer. A rigged coin
that always lands heads is extremely low entropy (utterly predictable) even
though it's unfair. Reach for "entropy" specifically when the honest
description is "spread out / could go several ways," not just "off" or
"biased."

**Now, cross-entropy — since it's the other place you'll hear this word.**
In AI conversations you'll hear "we're minimizing cross-entropy loss." Set
aside the math; here's what's actually being scored: how well does the
model's *confidence* match what's *actually true*? If a model says "I'm 90%
sure this photo is a cat" and it is a cat, that's a small penalty — good
call, stay confident. If it says "90% sure it's a cat" and it's actually a
dog, that's a *big* penalty — much bigger than if it had hedged and said
"maybe 50/50." Cross-entropy loss specifically punishes being confidently
wrong harder than it punishes being unsure. So when someone says "the
cross-entropy loss went down during training," what they mean in plain
English is: the model's confidence levels are now doing a better job of
matching reality — it's not just getting answers right more often, it's
getting *appropriately more or less sure of itself* as the evidence
warrants.

---

## gradient-descent — why "learning rate" is the knob that can break everything
**The idea in one line:** gradient descent always knows which way is
downhill, but the learning rate decides whether it walks there or vaults
straight over it.

Drop a ball on the inside wall of a bowl and it rolls to the bottom — at
every instant, gravity pulls it in the direction of steepest descent. Turn
that into an algorithm and you get gradient descent: at each step, look at
the slope where you're standing (the gradient) and move against it. The one
thing physics doesn't have to decide, but the algorithm does, is *how far*
to move at each step — that's the learning rate.

The picture uses the simplest possible valley, f(x) = x², whose minimum sits
at x = 0. A ball starts partway up the wall and takes a fixed number of
steps, each one x_new = x_old − learning_rate × slope. With a small learning
rate the ball creeps downhill, step after cautious step, converging on the
bottom — accurate, but slow. Push the learning rate up and it starts
overshooting: it crosses the bottom and lands partway up the *opposite*
wall, then crosses back, etc. As long as the overshoot each time is smaller
than the last, it still spirals in on the minimum. But past a critical
learning rate the overshoot gets *bigger* each step — the ball doesn't
oscillate toward the bottom, it flies further from it with every bounce,
literally off the edge of the chart. That's divergence: the algorithm didn't
get stuck, it actively made things worse by taking steps too large for the
curvature of the valley it's in.

**Here's what that actually looks like in numbers.** Start at x = 4.5, with
f(x) = x² so the slope at any point is 2x. At learning rate 0.3:

| step | x | slope (2x) | next x = x − 0.3 × slope |
|---|---|---|---|
| 0 | 4.50 | 9.00 | 4.50 − 2.70 = 1.80 |
| 1 | 1.80 | 3.60 | 1.80 − 1.08 = 0.72 |
| 2 | 0.72 | 1.44 | 0.72 − 0.43 = 0.29 |
| 3 | 0.29 | 0.58 | 0.29 − 0.17 = 0.12 |

Each step lands closer to 0, the true minimum — cautious, steady progress.
Now the same starting point at learning rate 1.1, just past this bowl's
critical value of 1.0:

| step | x | slope (2x) | next x = x − 1.1 × slope |
|---|---|---|---|
| 0 | 4.50 | 9.00 | 4.50 − 9.90 = −5.40 |
| 1 | −5.40 | −10.80 | −5.40 + 11.88 = 6.48 |
| 2 | 6.48 | 12.96 | 6.48 − 14.26 = −7.78 |

The sign flips every step — crossing back and forth over the bottom — and
the distance from zero *grows* each time: 4.50 → 5.40 → 6.48 → 7.78. That's
the "flies further from it with every bounce" divergence described above,
now with the actual numbers behind it.

**Where it bites in real life:** this is exactly how neural network training
fails when the learning rate is set too high — the loss doesn't plateau, it
blows up to NaN within the first few steps. It's also why training
schedules often *start* with a larger learning rate for speed and shrink it
over time: big steps to cover ground fast early on, small steps to settle
precisely into the minimum once you're close.

**Say it like this:** "We're overcorrecting — dial back how much we change
each time" is a learning-rate problem in plain English, whether it's a
thermostat, a steering wheel, or a training run: the adjustment per attempt
is too big, so it overshoots the target and the next correction overshoots
back the other way.
**Not like this:** "Just make bigger adjustments, we'll get there faster" —
true only up to a point. Past that point, bigger steps don't just take
longer, they actively make each round worse than the last, not better.

---

## overfitting — memorizing the practice problems instead of learning the rule
**The idea in one line:** a model that scores perfectly on its training data
hasn't necessarily learned anything — it may have just memorized the noise.

Give a student twelve practice problems and their answers, then ask them to
explain the pattern. One student writes down a short, general rule that gets
most of the twelve right and should generalize to problem thirteen. Another
memorizes all twelve exact answers, quirks and all — zero mistakes on
practice, but no real idea what to do with a new problem. The second student
has overfit: their "model" fits the specific data they saw, not the
underlying pattern that generated it.

The picture makes this literal. Twelve noisy data points sit near a smooth
true curve (the pattern a good model would recover). A polynomial is fit to
those twelve points by least squares, and its degree is the "model
complexity" knob. At low degree the curve can't bend much, so it stays close
to the true pattern and mostly ignores the noise. Crank the degree up and the
curve gains enough free parameters to weave through every single point
exactly — training error keeps dropping, all the way to zero once the degree
reaches eleven (one fewer than the number of points, so an exact fit always
exists). But look at the curve *between* the dots: it doesn't glide smoothly
along the true pattern anymore, it swings wildly, overshooting in both
directions to thread each noisy point precisely. That's the "true error"
number climbing even as "training error" falls to zero — the model is now
excellent at reproducing the twelve answers it memorized and worse at the
actual underlying rule.

**Where it bites in real life:** any model complex enough to memorize its
training set — a deep enough decision tree, a large enough neural net, a
high-enough-degree regression — will show this exact split: training
accuracy that looks fantastic and real-world performance that's worse than a
simpler model's. The fix is never "make the model even more powerful"; it's
validating on data the model never saw, and choosing complexity by how well
it does *there*, not on the practice set.

**Say it like this:** "You're overfitting to one bad interview" — treating
the specifics of a single, possibly-unrepresentative experience as if it
were the general rule. "This strategy is overfit to last quarter" — it was
tuned to fit exactly what just happened, not to what tends to happen.
**Not like this:** "It just needs more data to fix the overfitting" — not
automatically true. More of the same noisy data can let an already
too-flexible model memorize even more precisely; the usual fix is
simplifying the model or checking it against data it hasn't seen, not
volume alone.

---

## confusion-matrix — the trap a single accuracy number can hide
**The idea in one line:** "99% accurate" can describe both a genuinely
excellent classifier and a useless one that just guesses the common answer
every time — the confusion matrix is what tells you which one you're
looking at.

A hospital announces a new test for a rare disease is "99% accurate." That
sounds excellent — until you notice that a test which does nothing at all,
just prints "healthy" for every single patient without looking at anything,
would *also* score 99% accurate, as long as only 1% of patients actually
have the disease. That do-nothing test catches zero real cases. It is, for
the one purpose a medical test exists for, completely useless — and yet it
posts the exact same headline number as a test that's genuinely good at its
job. Accuracy alone cannot tell these two situations apart, because it
collapses four very different outcomes into a single ratio.

The confusion matrix refuses to collapse them. Every classification a model
makes lands in exactly one of four buckets: rows are the ground truth
(actually positive / actually negative), columns are the prediction
(predicted positive / predicted negative). The TP/TN diagonal is where the
model agreed with reality; the FP/FN diagonal is where it didn't — and those
two kinds of wrong mean very different things. A false positive is a false
alarm: flagging something that was actually fine. A false negative is a
miss: letting something real slip through undetected — exactly the failure
mode the do-nothing "always healthy" test commits 100% of the time, which
the raw accuracy number never revealed.

**With real counts, not just percentages:** 100 patients, 1 actually has
the disease, 99 don't. The do-nothing test predicts "healthy" for all 100,
no exceptions. Its confusion matrix: TP = 0 (it never predicts positive, so
it can't catch the 1 real case), FN = 1 (that real case, missed), FP = 0
(never wrongly flags anyone), TN = 99 (correctly clears every healthy
patient). Accuracy = (TP+TN) ÷ 100 = (0+99) ÷ 100 = **99%** — the identical
headline number a genuinely good test would post, produced here by a test
that catches literally zero real cases.

Look at the grid instead of the headline number and the do-nothing test is
exposed instantly: its TP cell (real cases actually caught) sits at zero, no
matter how green its TN cell looks. The picture shades each cell by its
share of the population so this is visible at a glance.

**What "threshold" and "class separation" actually control.** Underneath
the grid, every example gets a score — some raw number the classifier
computes for it. Real negatives cluster around a baseline score; real
positives cluster around a higher score, offset from that baseline by
however far **class separation** is set. Small separation means the two
groups' scores overlap heavily — genuinely hard to tell apart no matter
where you draw the line. Large separation means they barely overlap — easy
to tell apart. It's a property of how good the underlying signal is, not
something a threshold can fix. **Threshold** is simply where you draw the
line: anyone scoring above it gets called positive. Raise it and fewer
examples clear the bar (both "predicted positive" cells shrink); lower it
and more do.

**A realistic setting, not the degenerate do-nothing case:** threshold =
1.5, separation = 2.1, population = 200 (always split 100 real positive /
100 real negative). Positives cluster around 2.1, negatives around 0, so a
threshold of 1.5 sits between them, closer to the positive side. Working
through the same two score distributions the do-nothing test skipped
entirely: about 73 of the 100 real positives score above 1.5 and get caught
(TP=73), the other 27 score lower and get missed (FN=27); about 7 of the
100 real negatives happen to score above 1.5 anyway (FP=7, false alarms),
the other 93 correctly clear (TN=93). From those four counts:

- **Accuracy** = (73+93) ÷ 200 = **83%**
- **Precision** = 73 ÷ (73+7) = 73/80 = **91%**
- **Recall** = 73 ÷ (73+27) = 73/100 = **73%**
- **F1** = 2×0.91×0.73 ÷ (0.91+0.73) ≈ **81%**

**What to actually read off those four numbers together — don't stop at
"they're all pretty high."** The gap between precision (91%) and recall
(73%) is the real information: this threshold is *conservative* — when it
flags something, trust it (only 7 wrong out of 80 flags) — but it's
leaving over a quarter of real positives, 27 of 100, uncaught. Whether
that's good news depends entirely on what's being screened for, not on the
numbers alone: for a spam filter, missing 27% of spam while almost never
flagging real mail might be exactly the tradeoff you want; for a cancer
screen or fraud alert, missing 27% of real cases is a serious problem even
at 83% accuracy, and you'd lower the threshold, trading away some of that
91% precision for higher recall. F1 (81%) blends the two into one
comparison-friendly number, but it *hides* this asymmetry — a cautious
classifier and a trigger-happy one can land on the identical F1 while
making completely different mistakes. Precision and recall have to be read
side by side, not replaced by one number, to know which one you're
actually looking at.

Lowering the threshold here is a free, instant trade — try it before
anything more expensive. If *no* threshold gets both numbers where they
need to be, that's a separation problem, not a threshold problem — see the
precision-recall concept's "from a reading to an action" table for what
"improve separation" concretely means in a real system (better features,
cleaner labels, more data of the harder class, or a bigger model —
different fixes for different causes, not interchangeable).

**Where it bites in real life:** fraud detection, rare-disease screening,
security alerting — anywhere the thing you're trying to catch is rare, a
model can hit sky-high accuracy by mostly predicting the common outcome and
still be worthless at the one job it exists to do. Reading the actual 2x2
grid — not just the headline accuracy — is the only way to tell a genuinely
skilled classifier from one that's just exploiting an imbalanced dataset.

**Say it like this:** "That's a false positive" — flagged as a problem, but
it wasn't one (a spam filter catching a real email, a smoke alarm going off
from toast). "We had a false negative" — a real problem, missed entirely.
These are not interchangeable mistakes; which one you'd rather live with
depends entirely on what's actually at stake.
**Not like this:** "It's 99% accurate, so it's basically fine" — accuracy
alone can't tell you whether the mistakes it does make are harmless false
alarms or costly misses, and those two failures are almost never equally
bad.

---

## roc-auc — grading judgment, not just the cutoff someone happened to pick
**The idea in one line:** comparing two classifiers by their accuracy at
whatever threshold they happen to be using compares their *habits*, not
their underlying skill — an ROC curve and its AUC grade the skill directly,
at every possible threshold at once.

Two doctors review the same set of X-rays for a rare condition. Doctor A
calls almost anything even slightly ambiguous "suspicious" — nearly every
real case gets caught, but so does a lot of harmless noise, so Doctor A also
generates a pile of false alarms. Doctor B is the opposite: cautious, only
flags the clearest cases, rarely wrong when they do flag something, but
quietly misses a fair number of real cases along the way. Ask "which doctor
is more accurate" using their current habits and you'll get a misleading
answer, because their thresholds for "suspicious enough to flag" are
personal styles, not a measure of how well they can actually tell a real
case from a healthy one when they look at the same scan.

What you actually want to know is threshold-independent: at every possible
cutoff — from "flag almost nothing" to "flag almost everything" — how well
does each doctor separate the two groups? Sweep an imaginary threshold
across that whole range and, at each setting, plot the false-positive rate
against the true-positive rate. That path is the ROC curve. A doctor with
real diagnostic skill bows that curve up toward the top-left corner —
catching real cases while dragging along comparatively few false alarms at
every threshold, not just their habitual one. A doctor guessing blindly (or
a coin flip) can't do better than the diagonal, because at any cutoff,
whatever fraction of healthy scans they wrongly flag, they catch exactly
that same fraction of real cases too — no separation, no skill.

**One point on that curve, with real numbers.** Four sick patients score 9,
7, 6, 4; four healthy patients score 8, 5, 3, 1 (some overlap — this
doctor isn't perfect). Flag anyone scoring 6 or higher as "suspicious":

- Sick patients ≥6: 9, 7, 6 → caught 3 of 4 → **TPR = 75%**
- Healthy patients ≥6: 8 → 1 false alarm out of 4 → **FPR = 25%**

That's the single point (25% FPR, 75% TPR) this one threshold plots on the
ROC curve. Slide the threshold from "flag almost everyone" down to "flag
almost no one" and every possible cutoff plots its own (FPR, TPR) point —
trace all of them together and that's the full curve.

**AUC**, the shaded area under that curve, compresses the entire
threshold-independent comparison into one number: it's exactly the
probability that a randomly chosen real case scores more "suspicious" than a
randomly chosen healthy one, under that doctor's judgment. AUC 0.5 is a coin
flip; AUC 1.0 is perfect separation — and crucially, it never requires
either doctor to have picked the same cutoff, or any cutoff at all, before
you can compare them.

**Computed directly from that definition, same tiny dataset:** compare
every sick score against every healthy score (4 × 4 = 16 pairs) and count
how often the sick patient scored higher:

| sick score | beats 8 | beats 5 | beats 3 | beats 1 | wins |
|---|---|---|---|---|---|
| 9 | ✓ | ✓ | ✓ | ✓ | 4 |
| 7 | ✗ | ✓ | ✓ | ✓ | 3 |
| 6 | ✗ | ✓ | ✓ | ✓ | 3 |
| 4 | ✗ | ✗ | ✓ | ✓ | 2 |

12 wins out of 16 pairs: AUC = 12 ÷ 16 = **0.75** — a randomly chosen sick
patient outscores a randomly chosen healthy patient 75% of the time, under
this doctor's judgment, computed without ever picking a threshold at all.

**Where it bites in real life:** two models can post identical accuracy at
their default settings and have very different AUCs — the higher-AUC one
has more headroom no matter where you eventually set the operating
threshold, which is why AUC is the standard way to compare models before
deployment. It's a poor substitute for precision/recall once you've actually
picked one operating threshold and have to live with its specific
false-positive rate, though — AUC grades potential, not the one decision
you're actually stuck with in production.

**Say it like this:** "Model A has a higher AUC than Model B" means A has
more underlying skill at telling the two classes apart, at every possible
cutoff — not just at whatever threshold happens to be in use today.
**Not like this:** "Model A is more accurate right now, so it's the better
model" — that's comparing today's threshold setting, a tunable choice, not
either model's real separating power; a lower-AUC model can still look
better at one specific cutoff while being the weaker model overall.

---

## bayes-theorem — why "99% accurate" doesn't mean what it sounds like
**The idea in one line:** a positive test result should update your belief,
but it has to start from how rare the thing was *before* the test — and
skipping that step is why "99% accurate" and "99% likely you have it" are
two very different numbers.

A test for a rare disease is 99% accurate, and you test positive. Gut
instinct says: 99% chance you have it — the test is 99% accurate, after all,
what else would it mean? That instinct is wrong, often dramatically so, and
the gap is exactly what Bayes' theorem measures.

Here's the arithmetic gut instinct skips: imagine 100 people take the test,
and the disease is rare — say only 1 person actually has it. The test almost
certainly catches that 1 true case (99% sensitivity). But the test also has
some false-positive rate among the 99 healthy people — say it wrongly flags
5% of them, roughly 5 people. Count up everyone who tested positive: 1 real
case plus roughly 5 false alarms — about 6 positive results, only 1 of which
is real. A "99% accurate" test just produced a positive-test group that's
more than 80% wrong, because the healthy population was so much bigger than
the sick one that even a small false-positive rate on it outweighs the tiny
number of true cases available to catch.

The picture makes this concrete: split a population of 100 by the given base
rate, sensitivity and specificity, then regroup everyone by test *result*
instead of true status. In the "tested positive" strip, green squares are
true positives and red squares are false alarms — at the default numbers
(1% base rate, 99% sensitive, 95% specific), the strip is mostly red: about
5 false alarms for every true positive. Raise the base rate and green takes
over the strip; raise specificity and the false alarms vanish, because
specificity directly controls how many healthy people get swept into a
positive result in the first place.

Bayes' theorem is just the formula that combines these two contributions
correctly: P(condition | positive) = (sensitivity × prior) / (sensitivity ×
prior + false-positive-rate × (1 − prior)). Sensitivity alone (the "99%
accurate" part everyone fixates on) is only half the story — the prior
matters just as much.

Plug in the exact numbers from the 100-people story above — prior = 1%,
sensitivity = 99%, false-positive rate = 5%: (0.99 × 0.01) ÷ (0.99 × 0.01 +
0.05 × 0.99) = 0.0099 ÷ 0.0594 ≈ **16.7%**. That's the formula, computed
independently of the counting story, landing on the same "about 1 in 6"
answer the population-of-100 walkthrough gave — two different routes to the
same number.

**Where it bites in real life:** screening for a rare disease, a rare fraud
pattern, or a rare security alert — a "99% accurate" flag sounds like a
near-certainty, but if the thing being flagged is rare, most flags are still
noise. It's why doctors order a second, more specific test before acting on
one positive screen, and why "the model flagged it" needs a base rate
attached before anyone should trust the flag.

**Say it like this:** "Don't ignore the base rate" — before updating your
belief off one new piece of evidence, remember how rare or common the thing
was to begin with; a positive test for a rare condition means far less than
gut instinct suggests.
**Not like this:** "The test is 99% accurate, so a positive result means a
99% chance I have it" — that skips the base rate entirely, and for a rare
enough condition, the base rate is often *most* of what determines the real
answer.

---

## p-value — how suspicious is this, assuming nothing's actually going on?
**The idea in one line:** a p-value is the chance of seeing a result at
least this extreme purely by luck, if the thing you suspect isn't actually
true — it is not, and never was, "the probability you're wrong."

Your friend hands you a coin and claims it's fair. You flip it 10 times and
get 8 heads. Suspicious? Maybe — but a fair coin doesn't produce exactly 5
heads every single time either; sometimes it runs hot by pure chance. So the
real question isn't "is 8 heads a lot" in isolation, it's: if the coin
really were fair, how often would 10 flips produce a result this lopsided
(8+ heads, or by symmetry 8+ tails) just from ordinary luck? Work that out
and the answer is about 11% of the time — not rare enough to be damning, but
not nothing either. That 11% is a p-value.

**Where does "about 11%" actually come from?** Flip a fair coin 10 times and
there are 2¹⁰ = 1,024 equally likely head/tail sequences. Count how many of
them land on 8 or more heads: 45 sequences give exactly 8 heads, 10 give
exactly 9, and 1 gives all 10 — 56 sequences out of 1,024. By symmetry,
another 56 sequences give 8 or more tails. That's 112 out of 1,024
sequences landing at least this lopsided in either direction: 112 ÷ 1,024 ≈
**10.9%**, rounding to the "about 11%" above — no simulation, no gut
feeling, just counting equally likely outcomes.

That's the whole recipe: start from the "null" assumption — here, a fair
coin, formally a distribution of what results would look like if nothing
suspicious were going on. Run your actual test and get an observed result.
The p-value asks one narrow question: under that null distribution, how much
of the probability lies at least as far out as what you actually saw, in
either direction? That's the entire definition — nothing more. A small
p-value means your result would be a rare, surprising draw *if* the null
were true, which counts as evidence against the null — but it is emphatically
not "the probability the null is true." That's a different, much slipperier
quantity a single p-value can never hand you (getting 8 heads doesn't tell
you the probability the coin is fair; it tells you how surprising 8 heads
would be from a coin that already is fair).

The significance level α turns that shaded area into a yes/no call: pick a
threshold in advance (0.05 is conventional, not sacred), find the critical
boundary where the null distribution's tail area equals exactly α, and check
whether your observed statistic crossed it. Push the observed statistic
further from what the null predicts and the shaded p-value area shrinks;
once it drops below α, the result is labeled "statistically significant" —
a statement that a pre-drawn line got crossed, not a certificate of truth.

**Where it bites in real life:** "p < 0.05" gets read as "95% chance the
effect is real," which the coin story makes clear it isn't — it's a
statement about how a fair coin (or a no-effect world) behaves, computed
before you ever saw the data. It's also why p-hacking — trying enough
thresholds, subgroups, or metrics until one clears α by chance — is a real
problem: flip enough different coins, or ask enough different questions of
the same data, and a suspicious-looking result shows up eventually even when
nothing real is going on.

**Say it like this:** "The result was statistically significant" means:
if nothing real were going on, a result this extreme would have been
unlikely to show up by chance alone (below an agreed-on threshold, usually
5%) — nothing more, and nothing about how big or important the effect is.
**Not like this:** "p < 0.05, so there's a 95% chance the effect is real" —
a p-value is not the probability your hypothesis is true; it's the
probability of seeing data this extreme *if your hypothesis were false*.

---

## confidence-interval — a statement about the method, not about one answer
**The idea in one line:** "95% confidence" doesn't mean this one interval has
a 95% chance of containing the truth — it means that if you repeated the
whole procedure many times, about 95% of the intervals you built this way
would contain it.

Imagine a fish hiding somewhere in a lake, sitting perfectly still. You can't
see it, but you get an imperfect reading of roughly where it is, and you cast
a net centered on that reading. Cast a generously wide net and you'll
capture the fish almost every time, even with a rough guess. Cast a narrow
net and you'll only succeed when your guess happens to land close to the
truth. Now here's where intuition goes wrong: once you've thrown one
specific net, it's tempting to say "there's a 95% chance the fish is in this
net." But the fish was never moving, and the net has already landed — it
either caught the fish or it didn't. There's no more randomness left to
attach a probability to; you just don't know which case you're in.

What the 95% actually describes is the *casting method*, not any single
cast: draw a sample, compute its mean, build an interval of mean ± z·(σ/√n)
around it — do that many times, and about 95% of the resulting intervals
will contain the true value, purely because 95% of sample means happen to
land close enough to the truth for their interval to reach it. The
confidence level is a statement about how often this whole process succeeds,
made *before* you ever throw a net, not a probability you get to attach to
the specific net sitting in front of you.

**With real numbers:** a sample of n = 25 has mean 50 and standard
deviation σ = 10. The 95% margin is z × (σ/√n) = 1.96 × (10/√25) = 1.96 × 2
= 3.92, so the interval is 50 ± 3.92, roughly **[46.1, 53.9]**. Quadruple
the sample to n = 100 and the margin shrinks to 1.96 × (10/10) = **1.96** —
a much tighter net — because standard error falls as 1/√n: 4× the sample
only buys a 2× tighter interval (√4 = 2), not a 4× tighter one.

The picture makes this concrete without leaning on actual randomness: 20
"hypothetical repeated experiments" sit at evenly spaced quantiles of the
sampling distribution — an exact, reproducible stand-in for "cast the net 20
times" — each with its own interval, green if it captured the true mean
(dashed line), red if it missed. Raise the confidence knob and every
interval widens, turning red rows green, because a wider net catches more
fish regardless of how good your guess was. Raise the sample size instead
and every interval narrows around the same mean, because a bigger sample
means a sharper reading of where the fish actually is — standard error
shrinks as 1/√n. Notice what *doesn't* change coverage: n changes how tight
the net is, never how often it's cast wide enough to succeed — that's the
confidence level's job alone.

**Where it bites in real life:** "this poll has a ±3 point margin of error
at 95% confidence" means 95% of polls run this way would bracket the true
value — not that this specific poll has a 95% chance of being right. It's
also why a narrower interval from a bigger sample is a genuine improvement
(a tighter net, same success rate) while cherry-picking a higher confidence
level just to get a cleaner-looking one-off result buys you nothing (it only
pays off across many repeats, not on the one that matters to you).

**Say it like this:** "This poll has a 3-point margin of error at 95%
confidence" means the polling *method*, repeated many times, would bracket
the true number about 95% of the time.
**Not like this:** "There's a 95% chance the true value is in this specific
interval" — once the interval is drawn, it either contains the true value
or it doesn't; the 95% describes the method's long-run success rate, not
this one instance sitting in front of you.

---

## central-limit-theorem — why averaging fixes wildly unreliable guesses
**The idea in one line:** no matter how skewed or strange a population of
individual values is, the distribution of the *average* of many draws from
it becomes approximately normal and tightens around the true mean as you
average more of them.

At a school fair, a jar of jellybeans sits on a table and everyone takes a
guess at how many are inside. Look at any one guess and it's basically
useless — kids guess 50, guess 900, guess a suspiciously round 1,000. The
guesses as a group don't look like a tidy bell curve at all; they're skewed,
scattered, full of wild outliers. And yet a strange thing happens if you
average every single guess together: that average lands remarkably close to
the true count, far closer than almost any individual guess did. This is the
"wisdom of crowds" effect, and it isn't magic — it's the central limit
theorem.

Here's why it works: each individual guess is one noisy draw from "the
population of everything a person might guess," and that population can be
as lopsided as you like. But the *average* of many draws is a different
quantity entirely, with its own distribution — and that distribution behaves
in a way individual draws never do: it concentrates tightly around the true
mean and its shape becomes normal (bell-curved), regardless of how weird the
underlying population was, as long as you're averaging enough independent
draws. One wildly-off guess barely moves a large average; a handful of
wildly-off guesses in different directions mostly cancel each other out.

**With real numbers:** say the jar actually holds 620 jellybeans, and five
kids guess 50, 900, 1000, 300, and 750 — wildly scattered, off by as much
as 570 in either direction. Average those five guesses: (50 + 900 + 1000 +
300 + 750) ÷ 5 = 3000 ÷ 5 = **600**, only 20 off from the true 620 — closer
than four of the five individual guesses, even though every guess going
into it looked like noise.

The picture makes this exact rather than simulated: the population is an
exponential distribution — sharply peaked at zero with a long right tail,
about as far from a bell curve as it gets (think "time between rare events,"
or a population of guesses with a few extreme overestimates). At n=1, "the
sample mean" is just one draw, so the picture shows the skewed exponential
itself. Raise n and the sample-mean distribution visibly tightens and
reshapes toward the normal reference curve — and because the sum of n
exponential draws has an exact closed form (a Gamma distribution), this
isn't a simulation with sampling noise of its own; it's the true shape at
every n.

**Where it bites in real life:** polling averages, A/B test metrics, and
quality-control sample means can be treated as roughly normal — and standard
statistical tests applied — even when individual data points are skewed or
weird, as long as the sample is big enough. It's also why "the sample size
is too small" is a real objection: with just a handful of draws from a
skewed population, the average is still skewed too, and normal-based
confidence intervals can quietly mislead.

**Say it like this:** "Let's average a bunch of independent estimates, the
noise should wash out" — individual guesses can be wildly off, but averaging
enough of them (the "wisdom of crowds") reliably lands close to the truth,
regardless of how weird any one guess looked.
**Not like this:** "One really confident estimate beats an average of many
rough ones" — often false when the rough estimates are independent; the
average's error shrinks as you add more of them, while one confident-but-
biased guess never self-corrects at all.

---

## normal-vs-skew — same average, same spread, different shape
**The idea in one line:** two distributions can share an identical mean and
an identical standard deviation and still look nothing alike — skewness and
kurtosis are the extra numbers that capture the difference.

Two neighborhoods report the exact same average home price and the exact
same standard deviation. A house hunter might assume they're basically
interchangeable — same "typical" price, same amount of variation.

**With real numbers:** Neighborhood A, 5 homes: $380k, $390k, $400k, $410k,
$420k — mean = $400k, a clean, tight, symmetric spread. Neighborhood B:
four modest homes at $350k each and one mansion at $600k — mean =
(4×$350k + $600k) ÷ 5 = $2,000k ÷ 5 = **$400k**, exactly matching
Neighborhood A. Same reported average, built from two completely different
sets of homes — most of B's residents live well *below* the average their
own neighborhood reports, dragged up entirely by the one mansion. (Matching
the spread too, not just the mean, takes more than five hand-picked points
to pull off cleanly — the interactive picture does exactly that, with a
closed-form correction that keeps mean and stddev locked together while you
drag skew and watch the shape underneath them change.)

The actual experience of house-hunting in the two neighborhoods is
completely different. That gap is invisible to mean and standard deviation
alone; you need a third number.

**Skewness** is that number: it measures which way a distribution leans.
Positive skew (like neighborhood B) means a long right tail — a few big
values stretching out and dragging the mean above where most of the data
actually sits. **Kurtosis** (specifically *excess* kurtosis, measured
against the normal curve's baseline) answers a different question: are
extreme outcomes more or less common than a plain bell curve predicts? A
positive value means "fat tails" — more of the probability sits in a sharp
central peak *and* out in rare extremes, at the expense of the ordinary
middle ground. That sounds abstract until you remember 2008: risk models
built assuming roughly normal (thin-tailed) returns treated a market crash
as a once-in-thousands-of-years event. It happened anyway, because real
returns have fat tails — the "impossible" outcome was simply more likely
than the model's normal-curve assumption allowed for.

The picture builds both effects on top of a plain normal curve using a
closed-form correction, so at skew = kurtosis = 0 it's exactly the standard
bell curve, and each knob's effect is visible in isolation: skew tilts the
curve, lengthening one tail and shortening the other; kurtosis sharpens the
peak and fattens the tails (or the reverse) while keeping the curve
symmetric.

**Where it bites in real life:** financial risk models that assume
normality and get blindsided by fat-tailed crashes, A/B test metrics that
violate the normality assumptions baked into some significance tests, and
any dashboard that reports only "mean ± stddev" for data that's secretly
lopsided or spiky — two numbers that can hide a very different-shaped
reality.

**Say it like this:** "That sample is skewed" means the data leans hard
toward one side — a few extreme values on one end are pulling the average
away from where most of the data actually sits, even if the "typical"
spread (stddev) looks perfectly ordinary.
**Not like this:** "The average looks fine, so the data's probably normal" —
mean and standard deviation alone can't tell you the shape; two very
differently-shaped distributions can share both numbers exactly.

---

## variance-vs-stddev — why "average distance from average" doesn't work
**The idea in one line:** the obvious way to measure spread — average how far
each point is from the mean — always gives exactly zero, so variance squares
first (to stop the cancellation) and standard deviation square-roots back (to
undo the squared units).

Two dart throwers each take 4 throws. You measure how far each throw landed
from the bullseye, in inches — call it negative if the throw landed left of
the bullseye, positive if it landed right, the same way you'd count steps
left or right of home base:

```
Thrower 1 (tight):     -2   -1        +1   +2
                         •    •   |    •    •
                                bullseye

Thrower 2 (wild):  -8        -4        +4        +8
                     •         •   |    •         •
                                bullseye
```

Just looking at the two lines, Thrower 2 is obviously the wilder, less
consistent one — those throws are spread almost 4 times as wide as
Thrower 1's. Now try to turn that visible difference into one number the
"obvious" way: average how far off each throw was.

- **Thrower 1:** (−2) + (−1) + (+1) + (+2) = 0, average = 0 ÷ 4 = **0**
- **Thrower 2:** (−8) + (−4) + (+4) + (+8) = 0, average = 0 ÷ 4 = **0**

Both throwers come out to exactly the same number — zero — even though the
picture above makes it obvious they're nothing alike. That's not bad luck
with these particular throws; it's guaranteed to happen for *any* set of
throws, every time. It's baked into what "average" means: the throws above
the average and the throws below it are, by definition, balanced to cancel
out exactly. "Average raw deviation" is therefore useless here — it can't
tell a tight thrower from a wild one no matter how different they actually
are, because it always lands on zero before it gets the chance to.

**The fix: square each deviation before averaging it.** Squaring erases the
minus sign (−2 and +2 both become 4), so nothing cancels anymore — and as a
bonus, a throw twice as far off ends up counting *four times* as much, not
twice, so wild misses get punished harder than close ones:

| Thrower 1 throw | squared | Thrower 2 throw | squared |
|---|---|---|---|
| −2 | 4 | −8 | 64 |
| −1 | 1 | −4 | 16 |
| +1 | 1 | +4 | 16 |
| +2 | 4 | +8 | 64 |
| **sum** | **10** | **sum** | **160** |

Average those squared numbers and you get **variance**:

- Thrower 1: 10 ÷ 4 = **2.5** square inches
- Thrower 2: 160 ÷ 4 = **40** square inches

Now the two throwers are finally distinguishable — 40 is a lot bigger than
2.5, which matches what your eyes already told you from the number lines.
But look at the units: "2.5 square inches" isn't a sentence anyone actually
says about how consistent a dart thrower is. That's the one wrinkle
squaring introduces — it squares the units too. Taking the square root
undoes exactly that and lands back in ordinary inches:

- Thrower 1: √2.5 ≈ **1.58 inches** — this is the standard deviation
- Thrower 2: √40 ≈ **6.32 inches** — this is the standard deviation

That's standard deviation: the "how far off is a typical throw" question
the raw average was trying and failing to answer, now actually answered in
real inches — Thrower 2's typical miss is about 4 times further from the
bullseye than Thrower 1's, the same 4x gap the number lines showed at a
glance. (Thrower 2's throws here are literally each of Thrower 1's throws
×4 — which is exactly why variance came out ×16 [4²] and standard deviation
came out ×4: variance scales with the *square* of how far you stretch the
data, standard deviation scales with the stretch itself. The app's
interactive picture shows this same k-and-k² relationship directly — drag
its spread slider and watch the variance number shoot up far faster than
the standard deviation does.)

Push one point out as an outlier in the interactive picture and variance
jumps disproportionately too, for the identical reason: that squaring step
punishes the single far-out point hardest of all.

**Where it bites in real life:** why one wild outlier can wreck a
variance-based estimate more than you'd guess (median/IQR are more robust
because they don't square anything), why error is usually reported as RMSE —
root-mean-squared-error, i.e. "square, average, then square-root back" —
instead of raw MSE, and why standard deviation, not variance, is the number
that actually gets quoted next to a mean.

**Say it like this:** "That's a high-variance strategy" (common in startups,
poker, investing) means outcomes are spread wide — could go great, could go
badly, hard to call which. "Low variance" means outcomes cluster tightly
around what you'd expect, for better or worse.
**Not like this:** "High variance means it's worse" — variance measures
spread, not direction; a high-variance bet can have a *better* average
outcome than a safe one, just with far less certainty about any single try.

---

## mean-median-mode — three "typical" values that agree until they don't
**The idea in one line:** mean, median, and mode all answer "what's typical
here," but they answer it differently enough that a skewed dataset pulls
them apart — and knowing which one you're looking at matters.

Bill Gates walks into a small bar with 20 regulars, each with a net worth of
about $80,000. Before he sits down: total net worth in the room is
20 × $80,000 = $1,600,000, so the *average* (mean) is $1,600,000 ÷ 20 =
$80,000 — a fair summary of "the typical person here."

Now Bill Gates sits down. His real net worth is roughly $100 billion — round
it to exactly $100,000,000,000 for the math. The room's total net worth is
now $1,600,000 + $100,000,000,000 = $100,001,600,000, split across 21
people: $100,001,600,000 ÷ 21 ≈ **$4.76 billion**. That's the new mean.
Nothing about the 20 regulars changed — not one of them got richer or
poorer — but "the average person in this bar" now sounds like a
multi-billionaire, because the mean did arithmetic with Gates's number and
one $100-billion value is enough to drag a 21-person average that far on
its own.

The median doesn't have this problem. Sort all 21 net worths from lowest to
highest: twenty $80,000 values, then Gates's $100 billion sitting alone at
the very end. The *middle* of that line — the 11th person out of 21 — is
still one of the ordinary $80,000 regulars, because Gates is just one more
name standing at the far end of the line, not a number that gets *averaged
in*. Median net worth: still $80,000, barely moved. The mode — the single
most common value — doesn't move at all either: $80,000 is still shared by
20 out of 21 people in the room.

That gap between "median barely moves" and "mean rockets to billions" is
the whole lesson, and it generalizes: for a symmetric distribution, mean,
median, and mode all coincide, but skew the data and they peel apart in a
fixed order, **mode < median < mean**, which you can see with a much
smaller example. Five salaries: $40k, $40k, $45k, $50k, $200k.

- **Mode** = $40k (it's the only value that repeats — appears twice, every
  other value appears once)
- **Median** = $45k (sort them — 40k, 40k, 45k, 50k, 200k — and take the
  middle, the 3rd of 5)
- **Mean** = (40k+40k+45k+50k+200k) ÷ 5 = 375k ÷ 5 = $75k

$40k < $45k < $75k — mode, then median, then mean, in that exact order,
because the mean is the only one of the three that actually does arithmetic
with the $200k outlier instead of just counting or ranking it. The picture
shows the same effect continuously: drag the skew slider and watch mode,
median, and mean peel apart from a single starting point in that same fixed
order as the tail stretches out.

**Where it bites in real life:** "average household income" headlines are
almost always higher than what a typical household actually earns, because a
small number of very high earners drag the mean up while the median (what
statisticians usually mean by "typical") barely moves — exactly why income
and home-price statistics are usually reported as medians, and why a
company's "average" salary can look great in a press release while most
employees make noticeably less.

**Say it like this:** whenever someone says "the average," it's worth
asking which one they mean — "average household income" almost always means
the mean, and it's usually well above what a typical household actually
earns.
**Not like this:** "half of people earn less than the average" — only
guaranteed true if "average" means the median (the middle value by
definition); it's often false for the mean, which a handful of outliers can
drag well above where most people actually sit.

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

**Before going further, here's what r actually looks like with real
numbers**, at the three landmark values. Four students, hours studied vs.
test score:

```
hours:  1    2    3    4
score: 60   70   80   90
```
```
score
 90 |                    •
 80 |               •
 70 |          •
 60 |     •
    +----+----+----+----+
         1    2    3    4   hours
```
Every extra hour is worth exactly +10 points — the points sit dead on a
straight uphill line, no exceptions. That's **r = +1**, perfect positive
correlation. Now hours of sleep lost vs. next-day focus score:

```
lost:   0    1    2    3
focus: 90   80   70   60
```
Same perfectly straight line, only downhill this time — every hour of lost
sleep costs exactly 10 focus points. That's **r = -1**, perfect negative
correlation. Now shoe size vs. test score — two things with no real reason
to be related:

```
size:   8    9   10   11
score: 72   58   81   65
```
```
score
 81 |          •
 72 |     •
 65 |                    •
 58 |               •
    +----+----+----+----+
         8    9   10   11   shoe size
```
No line, no pattern — bigger shoe size doesn't reliably predict a higher or
lower score. Run the actual numbers and r comes out to about **0.03**,
essentially zero. That's the whole scale: r isn't measuring whether a
relationship *exists* in some deep sense, only how tightly the points hug a
straight line — and the ice-cream/drowning cloud above sits up near +1 for
exactly the same reason the hours-studied cloud does: real, tight,
predictable co-movement. It just doesn't tell you *why* the line is there.

That's the trap r sets by design: it only ever measures how tightly a
scatter of points hugs a straight line — as the three examples above just
showed directly. In the interactive picture, both axes are standardized, so
the trend line r implies always runs through the origin with slope r — why
the orange guide line visibly steepens as r moves toward ±1, the same way
the hand-plotted lines above got perfectly straight at r = ±1. But "sits
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

**Say it like this:** "X correlates with Y" is a purely descriptive claim —
they tend to move together — and stops exactly there; it says nothing about
which one, if either, is causing the other.
**Not like this:** "it's just correlation" is said so reflexively it's
become the opposite mistake just as often — dismissing a strong, useful
signal as meaningless because no mechanism has been proven yet, when a
strong correlation is usually worth investigating, not shrugging off.

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

**With real numbers:** 100 emails arrive, 20 really are spam, 80 are legit.
At one threshold the filter flags 25 emails as spam: 18 of those are truly
spam (caught), 7 are legit mail wrongly flagged, and 2 of the 20 real spam
emails slip through uncaught. Recall = 18 ÷ (18+2) = 18/20 = **90%** (caught
90% of real spam). Precision = 18 ÷ (18+7) = 18/25 = **72%** (72% of what
got flagged was actually spam). F1 = 2×0.72×0.90 ÷ (0.72+0.90) = 1.296 ÷
1.62 = **80%**.

Tighten the threshold and the filter gets pickier, flagging only 19 emails:
17 true spam caught, just 2 false alarms — but now 3 real spam emails slip
through instead of 2. Recall drops to 17/20 = **85%**, precision climbs to
17/19 ≈ **89%** — precision up, recall down, from the exact same tightened
dial. That's the tradeoff, in numbers instead of just a described shape.

The one thing that genuinely improves both is separating the two curves
further apart: a better model that scores real spam and real mail more
differently in the first place, leaving more room between "clearly not
spam" and "clearly spam" for a threshold to land cleanly.

**From a reading to an action: two knobs, two different costs.** Starting
from threshold=1.5, separation=2.1 (TP=73, FN=27, FP=7, TN=93 — precision
91%, recall 73%), compare what each knob actually buys:

- **Move the threshold to 1.0, separation unchanged:** TP=86, FN=14, FP=16,
  TN=84 → precision 84%, recall 86%. Recall gained 13 points; precision
  paid 7 for it. A trade, made instantly, no retraining — worth trying
  *first* whenever one metric is low and the other is comfortably high.
- **Raise separation to 3.0, threshold unchanged:** TP=93, FN=7, FP=7,
  TN=93 → precision 93%, recall 93%. *Both* improved. Not a trade — the
  underlying signal genuinely got better, which only ever comes from
  changing something upstream of the threshold.

**How to tell which one you need, before spending real effort:** check the
AUC (the concept next to this one). High AUC but disappointing numbers at
your current threshold means you're just standing in the wrong spot on an
otherwise good curve — sweep the threshold and stop. AUC close to 0.5 means
no threshold rescues you; that's when it's actually time to change the
model, not the cutoff.

**What "improve separation" concretely means, once you're actually there**
— "add more data" is the common instinct, and it's sometimes right, but
only for one specific cause:

| The real cause | The fix | Not this |
|---|---|---|
| Model never saw enough examples of the rarer class to learn its pattern | More labeled data *of that class* | More of the class you already have plenty of |
| Features don't contain the information needed to tell the classes apart | New/better features | More rows of the same uninformative columns |
| Some training labels are wrong, dragging one cluster toward the other | Clean the labels — often the cheapest fix, check this first | A bigger model, which will cheerfully learn the bad labels too |
| The true boundary is complex and the model is too simple to represent it | A more expressive model, same data | More data — a linear model doesn't get less linear from volume alone |

One check before any of these: compare performance on training data against
held-out data. Good on training, poor on held-out is overfitting (fix:
simplify or regularize, not enlarge — see the overfitting concept). Poor on
both is a genuine signal shortage, and the table above applies.

**Can precision, recall, F1, and AUC themselves tell you which cause it
is?** No — and it's worth being honest about that instead of pretending
otherwise. Those four numbers, computed once from one trained model, can
only tell you *that* there's a separation ceiling. A low AUC looks
identical whether the cause is missing features, noisy labels, too little
data, or too weak a model — same symptom, different diseases. Telling them
apart takes a few extra, cheap moves beyond the base metrics, cheapest
first:

| Diagnostic | What you do | What it tells you |
|---|---|---|
| Manual error review | Pull 20–30 false positives and false negatives, look at them by hand | A large chunk are actually mislabeled ground truth → noisy labels. Correctly labeled but genuinely look like the other class → a features or capacity problem, not labels |
| Train-vs-validation gap | Compare the same metric on training data vs. held-out data | Train high, val low → overfitting, not a signal shortage. Train low too → the model can't fit data it's already seen — underfitting: features lack signal, or the model needs more capacity |
| Learning curve | Retrain on 25%, 50%, 100% of the data, plot held-out performance against training-set size | Still climbing at 100% → genuinely data-limited, more of the hard class will likely help. Flat/plateaued → more of the same kind of data won't move it |
| Feature ablation | Add one candidate feature you suspect carries signal, retrain | Meaningful jump → the original features were the bottleneck. No change → keep looking |

The train-low-too row hides one more fork: "features lack signal" and
"model too weak" produce the *identical* symptom (low training
performance). The test that splits them: try a substantially more powerful
model on the exact same features. Still can't fit the training data → the
features are the ceiling. Suddenly fits well → it was capacity all along.
The aggregate metrics tell you *whether* to investigate; they don't do the
investigating for you.

**Where it bites in real life:** cancer screening (a false negative means
missed disease — worth tolerating more false positives to avoid), spam
filters (a false positive means a lost important email), fraud detection,
search ranking — anywhere "how sensitive should this be" is really a
business decision about which mistake costs more, dressed up as a slider.

**Say it like this:** "We need higher recall here, even if precision drops"
— for something like cancer screening, missing a real case is worse than a
false alarm, so the tradeoff should be set on purpose, not left to a
model's default.
**Not like this:** "just make the model more accurate" — accuracy alone
hides *which* kind of mistake it's making; the actual decision is which
type of error costs more in your specific situation.

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

**Say it like this:** "this is growing exponentially" is a genuine folding
story — each step multiplies what came before, not adds a fixed amount —
and it's exactly when a logarithm, not a raw number, gives you a sane scale
to reason about it (compound interest, an outbreak's doubling time).
**Not like this:** "exponential" gets used constantly to just mean "a lot"
or "suddenly" — most things people call exponential are actually growing at
a steady additive rate that merely *felt* surprising; real exponential
growth compounds, and compounding is what actually gets out of hand.

---

## standard-deviation — a ruler that means the same thing everywhere
**The idea in one line:** "10 points above average" tells you nothing by
itself — standard deviation is the ruler that turns a raw distance from
average into a number that means the same thing no matter what you're
measuring.

Two students each score 10 points above their class average. Same number,
"10 points above average" — but do they deserve equal bragging rights?
To answer that you need one more fact about each class: how far scores
*typically* land from the average there, high or low. That's the piece the
raw "+10" is missing, and once you have it, watch what it does to each
student's result:

- **Class A:** scores typically land only about 2 points from the average
  in either direction — this is a tightly bunched class, nobody strays far.
  Student A's +10 is 10 ÷ 2 = 5 times further from average than it's normal
  to land here. That's not a little above average, that's off the charts —
  practically nobody else in the class gets anywhere near that far out.
- **Class B:** scores typically land about 20 points from the average in
  either direction — a much more spread-out class, where landing far from
  average is routine. Student B's +10 is only 10 ÷ 20 = 0.5, half of how
  far it's completely normal to land here. Plenty of classmates routinely
  score further from average than that.

Same "+10 points," same raw number — but one student is a genuine outlier
and the other is barely above the middle of the pack, and the only reason
you can tell them apart is that one extra fact: how far a typical score in
each class lands from that class's average.

**So what is standard deviation, concretely?** It's exactly that fact,
made precise: take every score in a class, measure how far it sits from the
average, and boil all of those distances down to one typical, representative
distance. That's standard deviation — written σ (sigma is just the Greek
letter used as its label, nothing more). Class A's σ is about 2; class B's σ
is about 20 — the exact two numbers used above. Once a class has that one
number, you can describe any individual score by *how many σ's away from
average it sits* — that count is called a **z-score**, and it's the same
division just done above, given a name: student A's +10 points is +5σ,
student B's +10 points is +0.5σ. Same raw gap, translated into numbers that
finally mean comparable things.

**Here's the part that makes σ more than just a label — a fixed rule you
can lean on.** For data piled up in the classic bell shape (most values near
the middle, tapering off evenly on both sides), the same three percentages
hold every single time, whether the bell is narrow or wide: about **68%** of
everything sits within 1σ of the average, about **95%** within 2σ, and about
**99.7%** within 3σ. That's not a coincidence you have to re-measure per
situation — it's baked into the bell shape itself, the same way a circle's
circumference is always about 3.14× its width no matter how big you draw the
circle. That's why "+5σ" (student A) sounds astonishing and "+0.5σ" sounds
unremarkable: the σ-to-percentage ladder is what lets you attach an actual
rarity to a gap instead of eyeballing a raw number. (Computing σ exactly
involves one extra step — squaring each distance before averaging, then
square-rooting back — see variance-vs-stddev for why; the plain version
above is close enough to build the right intuition.)

**Where it bites in real life:** a 101°F fever is routine for some people and
alarming for others depending on their normal baseline and its variability;
manufacturing tolerances are set in sigmas, not raw units, because "off by
2mm" means something different for a bolt than for a bridge; and "1-in-20"
scientific thresholds (p < 0.05) are really a statement about how many sigmas
out a result has to land before it counts as surprising.

**Say it like this:** "that was a 3-sigma event" (said after market crashes,
unusual test scores, quality-control failures) means something this extreme
should happen only about 3 times in 1,000 if the usual pattern of variation
held — the 99.7% rule made concrete, not just a vague "very rare."
**Not like this:** calling a fixed number of points or dollars "a lot" or
"a little" on its own, without asking *relative to what* — the same 10-point
gap can be unremarkable in one dataset and the single most extreme value in
another.
