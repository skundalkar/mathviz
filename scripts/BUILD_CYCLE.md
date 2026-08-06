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
2. Read `BACKLOG.md`. Pick the FIRST unchecked `- [ ]` item. That is your
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
     concept in plain language: the one-line idea, what the knobs show, and
     where it matters in real life.
   - Check off the item in `BACKLOG.md` (move it under `## Done`).

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
