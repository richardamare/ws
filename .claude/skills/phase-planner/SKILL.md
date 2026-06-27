---
name: phase-planner
description: |
  Decompose an approved SOLUTION.md into executable phases. Each phase is a
  vertical or horizontal slice that stands alone, has measurable acceptance,
  references the SOLUTION.md sections it implements, and authors specific ADRs.
  Produces PHASES.md + STATE.md. Pushing the
  phases to GitHub Issues (the only tracker this repo uses) is a final step the
  user chooses; the generated issues must be self-contained enough for the
  `work-on-issue` skill to execute with zero clarifying questions.

  Triggers: after `solution-design` completes and SOLUTION.md is convergence-checklist-clean,
  when the user says "rozplánuj", "vytvoř fáze", "phase plan", "phases for", or
  references an existing approved SOLUTION.md.

  Does NOT trigger for: ad-hoc tasks, single-PR work, bug fixes, or when SOLUTION.md
  has UNRESOLVED items (resolve first via `solution-design`).
---

# Phase Planner Skill

You decompose an approved SOLUTION.md into executable phases. Each phase is the
**execution-level convergence contract**: when an agent (or you, three weeks
from now) picks a phase to work on, PHASES.md must pin down the scope, the
acceptance criteria, and the downstream artifacts (ADRs, PRs, tests) so the
result is the same every time — no re-deciding settled questions.

The output feeds two downstreams: `backlog-planner` (optional — decomposes each
phase into stories/subtasks), and the `work-on-issue` skill, which executes a
GitHub issue autonomously with **zero clarifying questions**. Whatever this
skill emits as issues must therefore be self-contained and measurable.
PHASES.md is the joint of the chain — its shape must be stable.

## The principal-architect lens (carried forward)

Same lens as `solution-design`: if you hand Phase 2 to `work-on-issue` and walk
away, does it have everything it needs to finish without asking you anything?
Every implicit default in PHASES.md is a question the agent can't ask — so it's
a divergence vector. Pin it down.

In addition, phase-planner enforces a *temporal* contract: phases are ordered,
and dependency order is part of the doc. Skipping a dependency is not "an agent
got creative" — it is a contract violation.

## Inheritance protocol (mandatory)

Before drafting, read in full:

- `docs/SOLUTION.md` — the upstream contract. Specifically extract:
  - **§5 Stack decisions** — every Inherited / LOCKED row drives execution defaults
  - **§6 Data model** — entities marked `MVP` are buildable now; `DEFERRED` cannot land in any phase without escalation; `reserved-table-only` ships schema only
  - **§7 Wire formats** — every phase that ships an endpoint must comply
  - **§11 Testing topology** — phase tasks include "extend `<test-project>` with N cases"
  - **§13 ADR roster** — each ADR has a Trigger column naming the phase that authors it
  - **§14 Risk register** — risks with named mitigation that bind to a phase
- Parent `CLAUDE.md` if present — for operational rules, commit conventions, branch naming
- Any existing `docs/PHASES.md` (idempotency — see Step 0)

**Upstream validation.** If SOLUTION.md has `UNRESOLVED` items, **do not draft
phases.** Stop and tell the user: "SOLUTION.md has N UNRESOLVED items blocking
phase planning. Resolve via `solution-design` (option a/b/c per handoff)." This
prevents inheriting ambiguity.

If SOLUTION.md lacks the v2 sections (§6 status markers, §13 ADR roster, §7
wire formats), the doc is pre-v2. Two options:

- **Soft mode:** proceed but emit warnings in PHASES.md "Planning notes" naming each missing input.
- **Strict mode:** tell the user "SOLUTION.md is pre-v2 shape; re-run `solution-design` against current cwd to bring it forward, or invoke me again with `--soft` to proceed anyway."

Default to soft. Always emit the warnings.

## Mode (inherited from SOLUTION.md)

Phase-planner runs in the same mode as the SOLUTION.md it consumes:

| SOLUTION mode | Phase-planner behavior |
| ------------- | ---------------------- |
| **Vibe** | Skip. Vibe SOLUTION.md doesn't need formal phases — the user runs tasks directly. If the user invokes phase-planner anyway, ask once to confirm. |
| **Standard** | Produce PHASES.md + STATE.md. Convergence checklist run lightly. GitHub Issues push optional. |
| **Contract** | Full discipline — convergence checklist mandatory as TodoWrite tasks, ADR mapping enforced, upstream validation strict. |

## Decision state model (carried forward)

Every scope decision in PHASES.md is in one of three states:

- **LOCKED** — scope frozen for execution. Implementation must follow.
- **DEFERRED** — explicitly punted to a later phase (named) or post-MVP.
- **UNRESOLVED** — blocks the phase. Listed at top of PHASES.md "Planning notes" with named owner.

There is no fourth state.

**Field-name note.** PHASES.md `Status:` per phase = *scope decision* (LOCKED /
DEFERRED / UNRESOLVED). STATE.md `Status:` per phase = *lifecycle state* (not
started / in progress / blocked / done). Different files, different semantics,
deliberate. Do not conflate.

## Process

### Step 0 — Idempotency check

If `docs/PHASES.md` already exists in cwd:

1. Read it and STATE.md.
2. Diff against SOLUTION.md — has the upstream contract changed since PHASES.md
   was authored (compare SOLUTION.md `Version`)?
3. Do **not** overwrite. Author a new version (bump `Version` field), append a
   change log entry, edit in place.
4. If the upstream change is structural (new entity in §6, removed §5 row),
   bump major version and re-run convergence checklist.

### Step 1 — Upstream read (silent)

Run the inheritance protocol. Hold these mental indexes:

- `entities_mvp = [...]` (from §6 — eligible for any phase)
- `entities_deferred = [...]` (cannot land without escalation)
- `wire_formats = {...}` (every phase shipping an endpoint must comply)
- `adr_roster = [{id, title, planned_trigger_phase}, ...]`
- `risks = [...]` (each tied to applicable phase)

### Step 2 — Decomposition

Apply these principles:

**Phase 1 is always a vertical slice.** End-to-end happy path of the primary
use case, fake data where needed, but actually running. After phase 1 the user
can demo the value proposition from SOLUTION.md §1.

Phase 1 scope is constrained by §6 status: it may build only entities marked
`MVP`. If the vertical slice requires an entity not yet marked MVP, surface
that as an UNRESOLVED item — do not silently upgrade an entity's status.

**Phase 2+ is horizontal expansion.** Real data integrations, error handling,
additional use cases, observability hardening, perf, security review.

**Each phase has a measurable acceptance criterion.** Not "feature complete" —
something a test or a script can assert. Quote performance targets from §3
Non-functional where applicable.

**Each phase stands alone.** If we run out of time/budget at phase N-1, that's
a legitimate stopping point. Important for client PoCs.

**Target 3–7 phases.** More than 7 → phases too small (merge). Fewer than 3 →
phases too big (split).

**ADR authorship.** Walk the §13 ADR roster. Each ADR with `Trigger: Phase N`
becomes a deliverable of that phase. If a phase has no ADRs in flight,
verify — most phases produce at least one.

**Risk binding.** Walk §14 risks. Each risk with a mitigation that fits in a
phase becomes a `Risks:` entry on that phase. Risks that don't bind to any
phase surface as cross-cutting concerns at the top of PHASES.md.

### Step 3 — Draft PHASES.md (use template below)

### Step 4 — Self-critique (convergence checklist)

Run the **Convergence checklist** (below). Contract mode: create a TodoWrite
task per line item, work them in order. Standard: run lightly. Vibe: skipped
(you shouldn't be in phase-planner in Vibe mode).

### Step 5 — Initialize STATE.md

(Template below.)

### Step 6 — (no bespoke prompt)

Don't author a per-phase run prompt. Once a phase becomes a GitHub issue
(Step 7), the issue *is* the prompt — `work-on-issue` reads it directly. Run it
foreground (`work on issue <N>`) or background (`claude --bg "work on issue <N>"`).
A bespoke prompt block would just duplicate the issue; put that effort into the
issue's acceptance criteria instead. See `docs/product/workflow.md` → "The prompt is the
issue".

### Step 7 — Handoff (GitHub Issues push is a separate user decision)

Output the handoff message (below). **Do not auto-create issues.** The user
picks one of:

- **GitHub Issues** — phase-planner offers to run `gh issue create` per phase.
  Each issue body MUST carry a measurable **Acceptance criteria** section
  (`.github/ISSUE_TEMPLATE/feature.md` shape) and reference `PHASES.md#phase-<n>`,
  so `work-on-issue` can execute it with zero questions. A phase too coarse for
  one issue is split into several, or run through `backlog-planner` first.
  **Issue titles describe the work, never the phase position** — `feat: customer
  onboarding API`, not `Phase 1: customer onboarding`. The `Phase N:` label is
  internal PHASES.md structure; the body links back via `PHASES.md#phase-<n>`.
- **Markdown-only** — PHASES.md + STATE.md are the source of truth; no tracker.

Issue creation is one command after the user confirms. Never create without
explicit confirmation — issues are visible and hard to clean up.

## Output template — PHASES.md

```markdown
# Development phases: <project title from SOLUTION.md §1>

**Version:** v0.1
**Source:** docs/SOLUTION.md (v<X>)
**Generated:** <ISO date>
**Mode:** Standard | Contract
**Tracker:** GitHub Issues | Markdown-only

## Planning notes

- (Empty if clean.)
- **WARNING:** SOLUTION.md missing §13 ADR roster — ADR authorship cannot be assigned. (Soft-mode example.)
- **UNRESOLVED:** <name> — owner: <person>, blocks: Phase <n>.

## Cross-cutting risks (not bound to a single phase)

| Risk | Likelihood | Impact | Mitigation | Owner |
| ---- | ---------- | ------ | ---------- | ----- |

---

## Phase 1: <name>

**Status:** LOCKED | DEFERRED | UNRESOLVED
**Goal:** <one sentence — quotable by an agent>

**Scope:**
- **In:**
  - <feature/module from SOLUTION §3 — name M-number>
- **Out:**
  - <explicitly punted to phase N or post-MVP>

**Entities (from SOLUTION §6):**
- `Customer` (MVP)
- `Scenario` (MVP) — partial: core + scoring fields only; computed-fields persistence
- `ScenarioEmbedding` (DEFERRED to Phase 3)

**Wire formats (from SOLUTION §7) this phase ships:**
- Endpoints comply with: pagination shape (skip/take), error shape (RFC 7807 + code), enum casing (snake_case), ID format (Guid), date format (ISO UTC Z).

**ADRs authored in this phase (from SOLUTION §13):**
- ADR-001 Platform
- ADR-002 Database
- ADR-006 Scoring deterministic
- ADR-007 IDs plain Guid
- ADR-008 Pagination
- ADR-009 Errors
- ADR-010 Enums snake_case

**Tasks:**
- [ ] (each task is roughly one PR. If a task becomes a GitHub issue, it must be
      self-contained and measurable enough for `work-on-issue` to run with zero
      questions — otherwise expand it via `backlog-planner` first.)

**Acceptance criteria:**
- <Testable assertion 1 — quote § NFR targets where relevant>
- <Testable assertion 2>

**Tests delivered (from SOLUTION §11):**
- `<test-project-1>`: N new cases covering <surface>
- `<test-project-2>`: N new cases covering <surface>

**Effort:** X–Y MD

**Dependencies:**
- Upstream: SOLUTION.md sections §5, §6, §7
- Phase dependencies: none (Phase 1 is the vertical slice)
- External: <e.g. Entra app reg provisioned by IT>

**Risks (bound from SOLUTION §14):**
- <Risk + mitigation tied to this phase's acceptance>

**Tracker:** Issue / Epic <id> (filled after Step 7 if tracker push happens)

---

## Phase 2: <name>

[same structure]

---

[etc.]
```

## Output template — STATE.md

```markdown
# Project state — <project title>

**Last updated:** <ISO timestamp>
**Current focus:** Phase <n>
**Active workstreams:** <list of phase numbers being worked in parallel, if any>

## Phases

| # | Name | Status | Tracker | Last update | Notes |
| - | ---- | ------ | ------- | ----------- | ----- |
| 1 | <name> | not started \| in progress \| blocked \| done | <Issue/Epic id> | <ISO> | <one line> |

## Latest session summary

<Appended at the end of each session per CLAUDE.md workflow rule. Each entry:
**Date / Phase / Done / In progress (exactly where) / Blocked / Next.**>

## Blocked items (cross-phase)

- <Phase N>: <blocker>, owner: <name>, since: <date>
```

## How a phase becomes runnable work

No bespoke prompt template. A phase is executed by turning it into a GitHub
issue (Step 7) whose **Acceptance criteria** are measurable and whose decisions
are all made — then running `work on issue <N>` (foreground) or
`claude --bg "work on issue <N>"` (background). The issue carries the upstream
contract by reference: its body cites `PHASES.md#phase-<n>` and the relevant
SOLUTION.md sections, and `work-on-issue` reads them. If the issue isn't
self-contained enough to run zero-question, the fix is a better issue, not a
longer prompt. See `docs/product/workflow.md` → "The prompt is the issue".

## Convergence checklist

Run before declaring PHASES.md done. Each line gets a tick.

**Upstream**
- [ ] SOLUTION.md was read in full (not skimmed) and is convergence-checklist-clean.
- [ ] SOLUTION.md has zero UNRESOLVED items, OR the unresolved items are explicitly listed in this PHASES.md "Planning notes" as blocking.
- [ ] Every phase's "Entities" line references only entities present in SOLUTION §6.
- [ ] No phase ships an entity marked DEFERRED in §6 without an UNRESOLVED escalation entry.

**Phase shape**
- [ ] 3–7 phases total.
- [ ] Phase 1 is a vertical slice (end-to-end happy path).
- [ ] Every phase has a measurable acceptance criterion (a test or script can assert).
- [ ] Every phase stands alone — stopping at the end of Phase N-1 is a valid product state.

**Decision state**
- [ ] Every phase carries Status (LOCKED / DEFERRED / UNRESOLVED).
- [ ] No phase is in "TBD" or "we'll see" state.

**ADRs**
- [ ] Every ADR in SOLUTION §13 is mapped to exactly one phase under "ADRs authored in this phase".
- [ ] No phase claims an ADR not in §13.

**Wire formats**
- [ ] Every phase that ships an endpoint cites the wire-format rows from §7 it complies with.

**Tests**
- [ ] Every phase names the test projects from §11 it extends and what it adds.

**Risks**
- [ ] Every §14 risk is either bound to a phase OR listed in PHASES.md cross-cutting risks at the top.

**Dependencies**
- [ ] Phase dependencies form a DAG (no cycles).
- [ ] Phase N never depends on Phase >N.

**State + handoff**
- [ ] STATE.md initialized with all phases listed as "not started".
- [ ] GitHub Issues push is **not** auto-executed — user decision pending.
- [ ] Any phase intended to become a GitHub issue has acceptance criteria measurable enough for `work-on-issue` to run zero-question.

**Versioning**
- [ ] PHASES.md has Version field + change log.
- [ ] Source field cites SOLUTION.md version explicitly.

## Anti-patterns

- **Skipping upstream validation.** Drafting phases on top of an UNRESOLVED SOLUTION.md inherits the ambiguity. Refuse and surface.
- **Re-deciding upstream choices.** PHASES.md does not change stack, wire formats, or entity status. If a phase wants to, it's a SOLUTION.md change first.
- **Aspirational entities in early phases.** Phase 1 cannot persist `ScenarioEmbedding` if §6 marks it DEFERRED. Escalate, don't quietly upgrade.
- **Authoring a bespoke run prompt per phase.** The issue is the prompt — put the detail in its acceptance criteria, run `work on issue <N>`.
- **Auto-creating GitHub issues** without explicit user confirmation. Tracker mutations are visible and hard to clean up.
- **`Phase N:` in issue titles.** Issue title describes the work (`feat: <thing>`), not its position in the plan. Phase linkage lives in the body (`PHASES.md#phase-<n>`).
- **Phases that depend on later phases.** That's a cycle. Re-order or split.
- **Risk-free phases.** Every phase has at least one risk — if none, you missed something. Re-walk §14.
- **Generic acceptance criteria.** "Feature complete" / "tests pass" — not assertable. Quote §3 NFRs where you can. Vague criteria make `work-on-issue` reject the resulting issue as underspecified (see `docs/product/definition-of-done.md`).
- **Letting STATE.md drift from PHASES.md.** Status changes go through STATE.md, never silent edits to PHASES.md.
- **Treating PHASES.md as immutable.** Phases get re-scoped (descope, split, merge) as work proceeds. Bump version, log the change.

## Downstream contract

Two consumers read PHASES.md:

- **`work-on-issue`** (via GitHub Issues) — the real execution path. It never
  reads PHASES.md directly; it reads the issue you create from a phase/task. So
  any issue derived from a phase must be self-contained and measurable per
  `docs/product/definition-of-done.md`, or `work-on-issue` will refuse it.
- **`backlog-planner`** (optional) — consumes PHASES.md and expects these stable
  anchors:

- `## Phase <N>: <name>` — heading shape
- `**Status:**` field
- `**Goal:**` field
- `**Acceptance criteria:**` list
- `**Tasks:**` list
- `**ADRs authored in this phase:**` list
- `**Entities (from SOLUTION §6):**` list

Renaming any of these silently breaks `backlog-planner`. If you need to evolve
the shape, coordinate the change across both skills.

## Handoff message

After writing PHASES.md + STATE.md, post:

```
PHASES.md vytvořen (v0.1, mode: <Standard|Contract>, source: SOLUTION.md v<X>).

Convergence checklist: <PASS — N/N> | <FAIL — list failed line items>

Phases:
- Phase 1: <name> — <vertical slice, M-numbers from §3>
- Phase 2: <name>
- ...
- Phase N: <name>

UNRESOLVED (blokuje execution) — <N> položek:
1. <item> — owner: <name>
   (Resolve teď / Defer to later phase / Accept block)

DEFERRED (záměrně později):
- <decision> → Phase <n>

GitHub Issues push:
  (a) GitHub Issues — pustím `gh issue create` per phase (každý issue dostane
      měřitelná acceptance criteria, aby ho `work-on-issue` zvládl bez doptávání)
  (b) Markdown-only — PHASES.md je source of truth, žádný tracker

Nepouštím nic automaticky — review PHASES.md je tvůj checkpoint.

Next: vyber (a/b). Po vytvoření issue spusť `work on issue <N>` (foreground)
nebo `claude --bg "work on issue <N>"` (background). Issue je prompt.
```

## Length scaling

- **Vibe:** skip — vibe SOLUTION.md doesn't need formal phases.
- **Standard:** PHASES.md with all phase sections, convergence checklist run lightly, GitHub Issues push optional.
- **Contract:** full discipline — convergence checklist as TodoWrite tasks, ADR mapping enforced, GitHub Issues push proposed but never auto-executed.
