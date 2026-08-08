# MathViz

An **interactive math & statistics gallery** that teaches concepts by letting
you turn the knobs — and a self-building project that grows one concept at a
time, committing a clean history to GitHub and sending you a daily lesson.

Every concept is a small, pure Go function that renders live **SVG** in the
browser via **WebAssembly**. Drag a slider and the picture updates instantly.
No server, no JavaScript framework, no database.

## What's inside right now

| Concept | What you learn by dragging the sliders |
|---|---|
| **Standard deviation (σ)** | The mean ± 1σ band always holds ~68% of the data, however wide you make it. |
| **Logarithms** | log_b(x) is "b to what power is x", and multiplying x by b only *adds 1* to the log. |
| **Precision vs. recall** | Moving the threshold trades precision against recall; only a better model improves both. |

More are queued in [`BACKLOG.md`](BACKLOG.md). The plain-language explanations
live in [`LEARNINGS.md`](LEARNINGS.md). How the daily build loop itself is
wired together — the cron trigger, the plugin architecture, the commit
gate — is documented in [`AUTOMATION.md`](AUTOMATION.md).

## Run it locally

Requires Go 1.24+ and Python 3 (for a static file server).

```bash
make serve         # builds the WASM and serves http://localhost:8080
```

Open http://localhost:8080 and pick a concept from the sidebar.

```bash
make check         # go vet + all unit tests (the loop's commit gate)
make digest        # print the "what got built + today's lesson" summary
```

## How the code is organized

```
cmd/
  wasm/       the browser front-end (compiled to WebAssembly) — DOM only, no logic
  digest/     prints the daily lesson summary from git log + LEARNINGS.md
internal/
  concept/    the plugin contract (Concept, ParamSpec) + global registry
  viz/        a tiny, dependency-free SVG builder (pure, testable)
  concepts/   one folder per lesson; each registers itself via init()
    stddev/  logscale/  precisionrecall/  all/
web/          index.html gallery shell + wasm_exec.js loader
scripts/      BUILD_CYCLE.md — the standing instruction for the build loop
```

**The key design choice:** all the math and all the SVG generation are *pure
Go* in the `concept` packages, so `go test ./...` validates every lesson. The
WASM layer (`cmd/wasm`) only reads params and injects the returned SVG — it
holds no logic worth testing. That split is what lets an autonomous loop build
features safely: the tests are the guardrail.

## Adding a concept by hand

1. Copy `internal/concepts/stddev/` to `internal/concepts/<id>/`.
2. Change the `init()` registration (ID, Title, Blurb, Params) and write the
   math + `render`.
3. Add `<id>_test.go` for the math.
4. Add one blank-import line to `internal/concepts/all/all.go`.
5. `make check`, then `make serve` to look at it.

That's the same recipe the build loop follows automatically.

## The self-building loop

This repo is designed to grow on its own on a schedule. Each cycle:

1. clones the repo and reads `BACKLOG.md`,
2. picks the first unchecked concept,
3. builds it as **~5 small, conventional commits** (`feat` → `feat` → `test` →
   `feat` → `docs`), running `make check` green before every commit,
4. pushes to `main` (or opens a PR), and
5. runs `make digest` and sends you the summary — your daily lesson.

The full standing instruction is in [`scripts/BUILD_CYCLE.md`](scripts/BUILD_CYCLE.md).
Commit granularity is deliberate: one concept is never a single huge diff, it's
a short readable story of small steps.

## Commit convention

`<type>(<concept-id>): <summary>` using `feat`, `test`, `docs`, `refactor`.
This keeps `git log` readable and makes each concept easy to review or revert.

## Optional: view it online

Enable **Settings → Pages → Build and deployment: GitHub Actions**. The included
[`.github/workflows/pages.yml`](.github/workflows/pages.yml) builds the WASM and
publishes the gallery to `https://<you>.github.io/mathviz/` on every push.
