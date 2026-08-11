# Build-cycle instruction

This is the standing task the scheduled build loop runs each time it fires. It
is written to be self-contained: each run starts fresh, clones the repo, builds
exactly one concept as a series of small commits, pushes, and reports back.

Copy this (with the repo URL filled in) into the scheduled task prompt.

---

You are the MathViz build loop. Do ONE concept this run, then stop.

**Setup**
1. Clone `<REPO_URL>` and `cd` into it. (`git` auth is available in this
   environment.) Ensure Go 1.24+ is installed.
2. Set your git identity to the repo owner's before making any commits:
   `git config user.name "Snehal Kundalkar"` and
   `git config user.email "snehal.kundalkar@gmail.com"`. Never let commits
   fall back to a session-default identity (e.g. `Claude <noreply@...>`) —
   that doesn't count toward the owner's GitHub contribution graph. If you
   want to flag a commit as AI-generated, add a `Co-Authored-By: Claude
   <noreply@anthropic.com>` trailer instead of using it as the author.
3. Read `BACKLOG.md`. Pick the FIRST unchecked `- [ ]` item. That is your
   concept: note its `<id>` and description. If there are no unchecked items,
   push nothing and report "backlog empty" — done.

**Build the concept as ~5 small, single-purpose commits.** Follow the existing
pattern in `internal/concepts/stddev/` exactly. Never bundle everything into one
big commit. After each step, `make check` must pass before you commit.

1. `feat(<id>): scaffold concept + register in gallery`
   - Create `internal/concepts/<id>/<id>.go` with the package doc comment, an
     `init()` that calls `concept.Register(...)` with the ID, Title, `Sections`
     (see the writing bar under step 5 — populate this, not the legacy
     `Blurb` field, for every new concept), and the Params (sliders). Render
     may return an empty `viz.New(...).String()` placeholder for now.
   - Set `Seq` to one more than the highest `Seq` currently registered (grep
     `Seq:` across `internal/concepts/*/*.go` to find it, or use
     `concept.Count()+1`). This is what keeps the gallery sidebar sorted
     newest-first — `Register` panics if `Seq` is missing or reused, so a
     forgotten or duplicate value fails loudly at `make check` rather than
     silently burying the new concept in the list.
   - Add the blank import line to `internal/concepts/all/all.go`.
2. `feat(<id>): pure math for <concept>`
   - Add the exported, side-effect-free math functions this concept needs.
     Pure in, pure out — no globals, no time, no randomness.
3. `test(<id>): unit tests for the math`
   - Add `<id>_test.go` asserting the math's known values and invariants, plus a
     `TestRenderProducesSVG` smoke test. Run `make check` — must be green.
4. `feat(<id>): interactive SVG render`
   - Implement `render(...)` using the `viz` helpers so each Param visibly
     changes the picture. Keep it readable.
5. `docs(<id>): explain <concept> + tick backlog`
   - Prepend a new section to `LEARNINGS.md` (newest on top) explaining the
     concept in plain language.
   - Check off the item in `BACKLOG.md` (move it under `## Done`).
   - **The in-app explanation is `concept.Section` blocks, not one Blurb
     paragraph.** `Concept.Sections []Section` (each a `Heading` and a
     `Body []string`) renders in the browser as separate headed blocks
     instead of one wall of text — the WASM front-end only falls back to
     the old single-paragraph `Blurb` when `Sections` is empty, which is a
     compatibility path for concepts built before this existed, not
     something a new concept should use. A `Body` entry prefixed with
     `"• "` renders as a bullet-list item instead of a paragraph — use it
     for a short worked sequence (e.g. sweeping a threshold through a few
     settings) so it reads as steps rather than a run-on sentence. A `Body`
     entry shaped like `"|cell|cell|cell|"` renders as one row of a real
     `<table>` instead — consecutive table rows form one table, first row
     as the header — use it whenever the concept was worked out as a
     pattern/diagnosis/action table in conversation (see
     `internal/concepts/evalplaybook/evalplaybook.go`'s "How does it
     actually work?" section): a table reads better than prose for that
     shape of content, and it's how the reference material looked in the
     chat it was transcribed from. See `internal/concepts/prauc/prauc.go`
     for a full worked bullet-list example, and load `web/index.html` after
     `make build && make serve` to see how either renders.
   - **Writing bar for the explanation.** The specific bullets below were
     reverse-engineered from a long series of docs-only fixup commits
     across nearly every existing concept (`git log --oneline -- LEARNINGS.md`
     going back to 2026-08-08) and a live walkthrough of `confidence-interval`
     — but a flat checklist of individually-satisfiable bullets turned out
     to be insufficient on its own: `pr-auc`'s first rewrite ticked every
     box below (real numbers, a scenario, a Say-it-like-this block) and
     *still* failed to explain why anyone would want it, because no single
     bullet tested for that. So before the checklist, the actual reusable
     framework, restated in one sentence so it transfers to concepts not
     yet built: **every entry answers, in this order, (1) why would anyone
     want this — the gap or tension a simpler prior approach leaves open,
     (2) how does it actually work — a mechanism derived from a small
     hands-on case, and (3) what can you now do that you couldn't before —
     and no entry may open with question 2.** The bullets below are how to
     execute each of those three well, not a replacement for asking them
     explicitly every time. Match this bar on the first pass instead of
     needing a follow-up fixup commit later:
     - **The gate every entry must pass before anything else: motivation
       before mechanism.** This is not one item in the list below — it's a
       precondition on the first two items, and a checklist of "has a
       scenario / has real numbers / has a Say-it-like-this block" can be
       satisfied in full while still failing it, which is exactly what
       happened to `pr-auc`'s first rewrite. The test: read only the first
       1-2 sentences. Can a total beginner now say, in their own words, why
       anyone would want this concept — *before* they understand a single
       thing about how it works? If the honest reaction to those sentences
       is "okay, why would I do that?", it's mechanism first and needs to
       be rewritten. Concretely, that means:
       - The one-line idea must name a **tension, a gap, or a question**
         the concept resolves — not a description of what the concept
         does. Bad (mechanism-first, what `pr-auc` actually shipped):
         "sweep every threshold and plot (recall, precision) — the area
         under that curve grades the classifier." Good (names the gap
         first): "a single precision/recall reading only tells you how one
         threshold performs, not whether a different threshold — or a
         different model entirely — would serve you better." Same
         concept, but only the second version gives the reader a reason to
         keep going.
       - The opening scenario must follow **situation → complication →
         open question**: what was the reader already doing (often a
         *simpler prior concept in this same gallery* — precision-recall
         itself is the prior concept `pr-auc` should have leaned on),
         where does that approach fall short or leave something
         unresolved, and what question does that gap leave open — all
         stated *before* any part of the new concept's mechanism shows up.
         `precision-recall`'s own opening already does this correctly
         (spam filter, "turn the dial all the way up," which visibly
         fails before precision/recall are even named) — use it as the
         model, not just `bayes-theorem`.
     - **The standard shape most entries now follow: six Sections, each
       heading a literal question, in this order** (deviate only when a
       concept genuinely doesn't need one — don't force a square peg in,
       but don't skip one just because it's easier not to write):
       1. **`"Why would you need this?"`** — held to the motivation-before-
          mechanism gate above: situation → complication → open question,
          pitched at a first-time learner with no assumed background, not a
          formula or niche-domain setup. If people commonly get this
          concept wrong on first instinct, say the wrong instinctive answer
          out loud in the reader's own likely words, *then* show why it's
          wrong — don't just state the correct version next to it unlabeled
          (model: `bayes-theorem`'s "gut instinct says 99%... that instinct
          is wrong"; `pr-auc`'s "which threshold do you even compare two
          filters at?" for a gallery-internal complication instead of a
          misconception).
       2. **`"How does it actually work?"`** — derive any formula or
          constant from the scenario in section 1 *before* stating it
          formally: walk a small, concrete, hands-on case (a guessing game,
          a hand-counted example, a short `"• "`-bulleted sequence of
          steps) that makes it feel inevitable, computed to real plugged-in
          figures, not left abstract. Don't open with the equation and
          explain it afterward; entropy needed three separate rewrites
          specifically because early passes kept leading with `-log2(p)`
          and bits before the reader had any reason to care
          (`docs(entropy)`'s three-commit chain — "20 questions" →
          "first-time learner" → "drop the math, explain the word itself"
          — is the cautionary example, not just a model to copy from). If a
          degenerate/trivial case would expose the concept's core trap
          (like a "do-nothing" classifier scoring 99% accuracy), work
          through that *first*, then a realistic, non-degenerate case
          (model: `confusion-matrix`). If the example produces several
          related numbers, say explicitly what to read off them *together*
          — don't stop at "here are the numbers, they're all fine" (model:
          `confusion-matrix`'s precision-vs-recall-gap paragraph).
       3. **`"What does the picture show?"`** — tied back to the exact same
          concrete numbers from section 2, not a separate abstract
          description; say what each Param/knob visibly does to the same
          worked example.
       4. **`"What can you do now that you couldn't before?"`** — the
          payoff, explicitly closing the loop section 1 opened. This is
          the section most likely to get skipped because sections 1-3 feel
          like "enough" — it's the one that turns "here's a mechanism" back
          into "here's why that mechanism mattered." (model: `pr-auc`'s
          rewrite — section 1 asks how to compare two models or pick a
          threshold; section 4 answers both explicitly, by name.)
       5. **`"Where does this show up in real life?"`** — real domains, but
          calibrated to what a generalist reader already has context for as
          the *primary* examples (Wordle, sports upsets, weather forecasts,
          studying for a test — not decision trees or cross-entropy loss as
          the lead example unless the concept itself is inherently an ML
          concept for an ML audience, like `sigmoid-softmax`).
          Specialist/niche applications are fine as a brief closing pointer
          ("if you want to go further...") but shouldn't carry the main
          explanation. If the concept's name is also an everyday English
          word people already use loosely outside math (entropy, variance,
          significant, confidence, bias), this is also where to quote a few
          real sentences using it that way, unpack each into what it
          actually means, and say when the word does *not* apply (model:
          `entropy`'s "this password has high entropy" list and its
          explicit "when the word doesn't fit" callout).
       6. **`"What's the common mistake here?"`** — the Say-it-like-this /
          Not-like-this contrast, now as its own headed section instead of
          a closing paragraph. This is standard on essentially every entry;
          include it every time, not just when a mistake feels obvious.
     - `LEARNINGS.md` keeps its existing free-form markdown structure (bold
       labels, not literal `Section`s) — the six-question shape above still
       applies there, just written as `**Why would you need this?**` etc.
       instead of Go struct literals.
     - Ground it in one concrete, worked numeric example with actual
       plugged-in numbers — a reader should be able to follow the
       arithmetic, not just the shape of the idea. Never assert a number or
       a conclusion without showing the arithmetic that produced it.
     - When a worked example compares two things side by side, describe
       both using the *same kind* of measurement — don't give one side
       "within 2 points of the mean" (a distance from center) and the
       other "a 40-point range" (a total span); pick one framing and apply
       it to both, or the comparison won't actually add up on the page
       (`standard-deviation`'s original two-student example had exactly
       this mismatch).
     - Define any jargon term the first time it's used, in the same
       sentence or the next one — don't assume the reader already knows it.
     - If you reach for a physical-world analogy, make sure its structure
       actually matches the math (don't describe a 2D area/region for a
       concept that's really a 1D range, etc.) — a mismatched analogy
       actively misleads rather than helping (this is exactly what went
       wrong with `confidence-interval`'s original fish-in-a-lake analogy).
     - If the concept reuses one word for two different things (e.g. a
       *sample* statistic vs. the *true/population* value it estimates;
       standard deviation vs. standard error; a model's raw score vs. a
       calibrated probability), name both explicitly and say plainly why
       they differ — don't assume the reader tracks an implicit swap.
     - If the math leans on a constant that looks arbitrary (a threshold, a
       critical value, a magic number), say in one sentence what question
       it's the answer to, not just that it gets used.
     - Where a parameter's effect is easiest to see side by side (e.g. how
       an answer shifts across a few settings of a knob), add a small
       comparison table rather than only prose.
     - When it's genuinely useful and doesn't force it, add a "from a
       reading to an action" table — given a specific symptom in the
       numbers, what would you actually *do* differently (model:
       `precision-recall`'s "from a reading to an action" and diagnostic
       tables). Also fine to be explicit about what the numbers *can't* tell
       you on their own, if that's true (same entry: "Can precision, recall,
       F1, and AUC themselves tell you which cause it is? No —").
     - If this concept's setup or picture builds directly on another
       already-shipped concept (e.g. precision-recall and confusion-matrix
       share the same threshold/separation setup), say so and point at it
       by name instead of re-explaining it from scratch or leaving the
       connection implicit.
     - **Keep the in-app `Sections` (in `<id>.go`) and the `LEARNINGS.md`
       entry numerically consistent** — same worked example, same numbers,
       in both places. `LEARNINGS.md` can go deeper (extra rows in a table,
       a footnote), but never contradict the shorter in-app version. If you
       improve one, mirror the fix in the other in the same commit (or
       immediately after) rather than letting them drift.
   - If a later real conversation surfaces a fix or a gap in how an already-
     shipped concept is explained, apply it immediately as its own commit
     (`docs(<id>): ...`) rather than letting the gap sit — see the
     `docs(confidence-interval)` commit history for the pattern.
   - **When fixing an already-shipped entry, read it first and prefer an
     additive, surgical fix over a wholesale rewrite.** Applying this whole
     checklist to an existing entry does not mean tearing it down and
     starting over — a prior pass discovered that all 15 other entries
     already had a grounded real-world scenario from an even earlier
     session, and only needed the "Say it like this" block *added*; a full
     rebuild would have risked discarding a worked example that was already
     load-bearing (`docs: apply the entropy pattern's 'say it like this'
     block to every concept` is the model for this: additive, not a
     rewrite). Only rebuild a section from scratch when it's demonstrably
     missing the thing you're fixing, not on general principle.

**Finish**
- `make check` one final time. If anything fails, fix it before pushing; never
  push red.
- Push to `main` (auto-commit mode) OR open a pull request (review mode),
  whichever the task is configured for.
- Run `go run ./cmd/digest` and put its output in your reply so the daily
  summary reaches the user.

**Rules**
- One concept per run. Small commits. Green before every commit.
- Keep the code style consistent with existing concepts.
- If the concept is ambiguous, make a reasonable choice and note it in the
  LEARNINGS entry rather than stalling.
