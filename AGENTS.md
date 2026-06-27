# AGENTS.md

Operational guide for Claude Code in the `ws` repo. Read once per session.

`ws` is a small Go CLI that sets up a per-project developer workspace: opens the cmux workspace + tabs,
logs into Azure as a scoped Reader service principal, and tracks Claude Code session IDs.

## Communication

- Code, identifiers, commit messages, PR titles, issue titles: **English**. Conversation: match the user.
- Be terse — the user reads diffs, not summaries.
- Output design matters: `ws` itself must be human- and LLM-friendly (see `docs/product/output-modes.md`).

## Where things live

| Topic | File |
| --- | --- |
| What ws is, goals, non-goals | `docs/product/overview.md` |
| Stack, libraries, binary layout, commands | `docs/architecture/stack.md` |
| Per-project config (YAML/JSON) schema | `docs/architecture/schemas/config.md` |
| Command surface (new/up/down/auth/sessions…) | `docs/product/workflow.md` |
| Session bookmarks (reuse Claude sessions) | `docs/product/sessions.md` |
| Output modes (pretty / structured-text / json) | `docs/product/output-modes.md` |
| Pre-merge quality bar | `docs/product/definition-of-done.md` |
| **Azure security model (read before any az work)** | `docs/security/README.md` |
| Roadmap & build order | `docs/roadmap.md` |

Patterns:

| Topic | File |
| --- | --- |
| Go conventions | `docs/patterns/go.md` |
| Cobra + huh CLI patterns | `docs/patterns/cli.md` |
| cmux integration | `docs/patterns/cmux.md` |
| Azure CLI / SP usage | `docs/patterns/azure.md` |
| Git commit & branch hygiene | `docs/patterns/git.md` |

Read the relevant file before working in that area. `docs/README.md` is the index.

## Non-negotiable rules

These apply to every change. Full rationale in `docs/security/README.md` and `docs/product/workflow.md`.

1. **Never modify the maintainer's Azure role assignments.** They are slow to regain. Scope sessions with a
   *separate* Reader SP, never by changing his account.
2. **The per-project SP is Reader-only.** Never grant it write, never use it for write/Terraform.
3. **Write / Terraform only via deliberate personal `az login`** (`ws elevate`) — never automated, never
   in a default session. `terraform apply` requires human approval.
4. **Never push to `master`.** `master` is always buildable.
5. **Every non-trivial change starts with a GitHub issue.** `gh issue create` if none exists.
6. **Branch naming:** `feat/issue-<n>-<slug>`, `fix/issue-<n>-<slug>`, `chore/<slug>`.
7. **Commit format:** `type(scope): subject (#<issue>)` — Conventional Commits + issue link. Atomic.
8. **`gofmt` + `go build ./...` pass before commit.** Never `--no-verify`, `--amend`, or force-push.
9. **No database. Local files only** (`~/.config/ws/`).
10. **Every PR meets `docs/product/definition-of-done.md`.** `/code-review` is the quality gate.

## Editing defaults

- Go. Single static binary. Shell out to `cmux` / `az` / `docker`; don't reimplement them.
- Default to no comments — only a comment for a non-obvious *why*.
- Don't add error handling for cases that can't happen. Trust the framework.
- Every command works both interactively (huh) and arg-driven; missing input under `--json` is an error.
- Keep the CLI short to type — the user explicitly wants minimal typing.
- Skill-generated planning docs (`docs/SOLUTION.md`, `docs/PHASES.md`, `docs/STATE.md`) are allowed.
  Don't create other planning docs unprompted.

## Tooling

- `solution-design` skill — generate `docs/SOLUTION.md` from a natural-language request.
- `phase-planner` skill — decompose `SOLUTION.md` into phases + self-contained GitHub issues.
- `work-on-issue` skill — "work on issue N" runs the full dev cycle from one issue, zero clarifying
  questions. Thin issue → it comments what's missing and stops.
