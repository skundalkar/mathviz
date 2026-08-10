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
     `init()` that calls `concept.Register(...)` with the ID, Title, a clear
     plain-language Blurb, and the Params (sliders). Render may return an empty
     `viz.New(...).String()` placeholder for now.
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
   - **Writing bar for the explanation.** This whole checklist was reverse-
     engineered from two things: a long series of docs-only fixup commits
     across nearly every existing concept (`git log --oneline -- LEARNINGS.md`
     going back to 2026-08-08 — the `docs(confusion-matrix)`,
     `docs(bayes-theorem)`, `docs(precision-recall)`, `docs(standard-
     deviation)`, and `docs: apply the entropy pattern... to every concept`
     commits especially), and a live walkthrough of `confidence-interval`
     that surfaced further gaps. Match this bar on the first pass instead of
     needing a follow-up fixup commit later:
     - **The standard shape most entries now follow** (deviate only when a
       concept genuinely doesn't need a piece of it):
       1. `**The idea in one line:**` — one sentence.
       2. A relatable real-world scenario. If people commonly get this
          concept wrong on first instinct, say the wrong instinctive answer
          out loud in the reader's own likely words, *then* show why it's
          wrong — don't just state the correct version next to it unlabeled
          (model: `bayes-theorem`'s "gut instinct says 99%... that instinct
          is wrong" opening).
       3. `**With real numbers:**` a full worked example computed step by
          step to concrete figures, not formulas left abstract. If a
          degenerate/trivial case would expose the concept's core trap (like
          a "do-nothing" classifier scoring 99% accuracy), work through that
          *first*, then follow with a realistic, non-degenerate worked case
          (model: `confusion-matrix`).
       4. What the interactive picture/knobs show, tied back to the same
          concrete numbers from step 3, not a separate abstract description.
       5. When the example produces several related numbers, say explicitly
          what to read off them *together* — don't stop at "here are the
          numbers, they're all fine" (model: `confusion-matrix`'s precision-
          vs-recall-gap paragraph).
       6. `**Where it bites in real life:**` concrete domains/situations.
       7. `**Say it like this:** / **Not like this:**` — a closing block
          contrasting correct phrasing against the natural-sounding wrong
          phrasing. This is now standard on essentially every entry; include
          it every time, not just when a mistake feels obvious.
     - Ground it in one concrete, worked numeric example with actual
       plugged-in numbers — a reader should be able to follow the
       arithmetic, not just the shape of the idea. Never assert a number or
       a conclusion without showing the arithmetic that produced it.
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
     - **Keep the in-app `Blurb` (in `<id>.go`) and the `LEARNINGS.md` entry
       numerically consistent** — same worked example, same numbers, in
       both places. If you improve one, mirror the fix in the other in the
       same commit (or immediately after) rather than letting them drift.
   - If a later real conversation surfaces a fix or a gap in how an already-
     shipped concept is explained, apply it immediately as its own commit
     (`docs(<id>): ...`) rather than letting the gap sit — see the
     `docs(confidence-interval)` commit history for the pattern.

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
