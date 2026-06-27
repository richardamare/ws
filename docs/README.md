# Documentation

Project documentation, technical conventions, and operational guides. Read the
root `CLAUDE.md` first — it carries the non-negotiable rules and points here for
detail.

Docs are grouped by concern: **architecture** (what we build), **product** (how
we build and what the tool does), **security** (the Azure trust model — the most
important doc in this repo), and **patterns** (prescriptive code conventions).

## architecture/

What the system is.

- **`stack.md`** — language, libraries, binary/package layout, command reference.
- **`adr/`** — Architecture Decision Records in [Michael Nygard format](architecture/adr/README.md). One file per decision, numbered `NNNN-title.md`.
- **`diagrams/`** — Mermaid (`.mmd`) diagrams; render in-place on GitHub.
- **`schemas/`** — the per-project **`config.md`** (YAML/JSON schema for `~/.config/ws/projects/*`).

## product/

How we build, and what `ws` does for the user.

- **`overview.md`** — what ws is, goals, non-goals, status.
- **`workflow.md`** — the command surface and how a project goes from `new` to daily `up`.
- **`sessions.md`** — Claude Code session bookmarks (reuse good-context sessions instead of hoarding workspaces).
- **`output-modes.md`** — pretty / structured-text (LLM default) / `--json`.
- **`conventions.md`** — code conventions and editing rules.
- **`definition-of-done.md`** — the pre-merge quality bar; `/code-review` as the gate.

## security/

**The most important doc in the repo.** [`security/README.md`](security/README.md) — the Azure trust
model: scoped Reader service principal per project, `AZURE_CONFIG_DIR` isolation, the deliberate write
path, and the hard rule that the maintainer's role assignments are never touched.

## patterns/

Prescriptive conventions for consistent implementation:

- **`go.md`** — Go style, project layout, error handling, shelling out to external tools.
- **`cli.md`** — Cobra command structure + huh interactive pickers; the "works typed or interactive" rule.
- **`cmux.md`** — how ws drives cmux (hybrid template + live), useful verbs, the socket fallback.
- **`azure.md`** — `az` invocation, the Reader SP, `AZURE_CONFIG_DIR`, cert handling and rotation.
- **`git.md`** — commit hygiene, message format, branching.

## Skill-generated docs

The `solution-design`, `phase-planner`, and `backlog-planner` skills (`.claude/skills/`) write these
into `docs/` during planning — the sanctioned exception to the "no planning docs" rule:

- **`SOLUTION.md`** — the project's solution design (reserved name; not shipped yet).
- **`PHASES.md`** — phase decomposition derived from `SOLUTION.md`.
- **`STATE.md`** — live project state, updated at the end of each session.

The current pre-skill design intent lives in [`roadmap.md`](roadmap.md).
