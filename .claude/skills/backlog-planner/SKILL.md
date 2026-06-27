---
name: backlog-planner
description: |
  Decompose phases into a backlog: Phase → Story → Subtask. Each story is a
  vertical slice with INVEST shape; each subtask is roughly one PR and stands
  alone as one GitHub issue the `work-on-issue` skill can execute with zero
  clarifying questions. Produces BACKLOG.md only — creating the GitHub Issues
  is a separate, user-confirmed step (`gh issue create`).

  Triggers: after `phase-planner` completes and PHASES.md is convergence-checklist-clean,
  when the user references an existing PHASES.md, or says "backlog", "stories",
  "subtasks for phases", "decompose phases", "rozplánuj fáze na stories".

  Does NOT trigger for: ad-hoc tasks, single-PR work, bug fixes, when no PHASES.md
  exists (run `phase-planner` first), or when SOLUTION.md / PHASES.md has UNRESOLVED
  items (resolve upstream first).
---

# Backlog Planner Skill

You decompose engineering phases into a backlog with Story and Subtask
granularity. BACKLOG.md is the **execution-unit contract** — each subtask body
must work as a GitHub issue that `work-on-issue` runs autonomously, because that
is exactly what it will be used for.

This skill is **stack-agnostic and brownfield-safe**. It does not assume any
particular language or framework. The tracker is GitHub Issues (the only one
this repo uses). It reads what SOLUTION.md and PHASES.md describe and decomposes
accordingly.

## The principal-architect lens (carried forward)

Same as upstream skills: if you turn "Story 2.3" into an issue and hand it to
`work-on-issue`, does it have everything to finish without asking you anything?
Every implicit default in a subtask body is a question the agent can't ask — a
divergence vector.

The shape is right when both hold:

- **The execution test:** `work-on-issue` can take the subtask as an issue, with
  no other context beyond the upstream artifacts the body references, and ship
  correct output with zero clarifying questions.
- **The owner test:** you (or your one teammate) can read the same subtask and
  recognize what is being built and how to tell it's done.

## Inheritance protocol (mandatory)

Before drafting, read in full:

- `docs/PHASES.md` — required. Specifically extract:
  - **Phase number, name, status, goal** per phase
  - **Acceptance criteria** per phase
  - **Tasks list** per phase (treat as a starting hint — you may reshape)
  - **Entities** referenced per phase
  - **ADRs authored** per phase
  - **Wire formats** referenced per phase
- `docs/SOLUTION.md` — required for context. Extract:
  - **§2 Glossary** — vocabulary; do not introduce synonyms
  - **§6 Data model** — entity status (MVP / DEFERRED) drives what a story may persist
  - **§7 Wire formats** — subtask bodies cite specific rows when relevant
  - **§11 Testing topology** — subtasks include "add N cases to `<test-project>`"
- Parent `CLAUDE.md` if present — operational rules

**Upstream validation.** Refuse to draft if:

- PHASES.md is missing (tell user to run `phase-planner` first)
- PHASES.md has any phase in `UNRESOLVED` status (resolve upstream first)
- PHASES.md cites entities not present in SOLUTION §6 (chain corruption — escalate to user)
- PHASES.md ADR mapping references ADRs not in SOLUTION §13 (chain corruption — escalate)

Soft mode (degrade gracefully if SOLUTION.md is pre-v2):

- Emit warning in Planning notes naming each missing input
- Proceed with available context
- Mark every subtask whose context is uncertain with `**Verify:**` line listing what to confirm against the codebase

## Mode (inherited from PHASES.md)

| PHASES mode | Backlog-planner behavior |
| ----------- | ------------------------ |
| **Vibe** | Skip. Vibe doesn't need backlog. If user invokes anyway, ask once. |
| **Standard** | Produce BACKLOG.md. Convergence checklist run lightly. |
| **Contract** | Full discipline — convergence checklist as TodoWrite tasks, story-to-§section mapping enforced, DAG validated. |

## Decision state model (carried forward)

Every story carries a state:

- **LOCKED** — scope frozen; subtasks may be executed.
- **DEFERRED** — punted to a later phase (named) or post-MVP.
- **UNRESOLVED** — story scope ambiguous; listed in Planning notes with owner.

Subtask order within a story is implicit (top-to-bottom). Subtasks do not have
status — only their parent story does.

**Field-name note.** BACKLOG.md `Status:` per story = *scope decision* (LOCKED
/ DEFERRED / UNRESOLVED), same semantics as PHASES.md's per-phase Status. The
upstream PHASES.md scope decision flows down: a phase's LOCKED status doesn't
override a story's UNRESOLVED, but a phase's DEFERRED forces all its stories
DEFERRED too.

## Core principles

**Phase = Epic.** No re-decomposition of phase boundaries here. Phase-planner
owns phase scope; backlog-planner just expands downward.

**Story = vertical slice.** Independently demoable. INVEST: Independent,
Negotiable, Valuable, Estimable, Small, Testable. Title pattern:
`"<role> can <action>"` or `"<system> <behavior>"`. User/system-centric, never
technology-centric.

**Subtask ≈ one PR.** Imperative title. One area (folder / module).
Self-contained body — measurable acceptance criteria, every decision made, no
vague language — so `work-on-issue` runs it zero-question (`docs/product/definition-of-done.md`).
A subtask is what one agent session finishes in one sitting.

**One subtask = one GitHub issue.** The subtask body maps straight onto a GitHub
issue (`.github/ISSUE_TEMPLATE/feature.md` shape). No translation layer.

**No tracker writes here.** BACKLOG.md is Markdown only. Creating issues
(`gh issue create` per subtask) is a separate, user-confirmed step.

**Story-to-§section mapping is mandatory in Contract mode.** Every story body
cites which SOLUTION.md sections it implements (entities from §6, wire formats
from §7, glossary terms from §2 it uses). This is the inheritance trace.

## Process

### Step 0 — Idempotency check

If `docs/BACKLOG.md` already exists:

1. Read it.
2. Diff against PHASES.md — phase count or names changed?
3. Do **not** overwrite. Bump version, append change log entry, edit in place.
4. If PHASES.md version moved (structural change), re-run full decomposition
   for affected phases only; preserve unaffected phases verbatim.

### Step 1 — Upstream read (silent)

Run the inheritance protocol. Hold these mental indexes:

- `phases = [{n, name, status, goal, acceptance, tasks, entities, adrs, wire_refs}, ...]`
- `entities_mvp_status = {name: status}` from SOLUTION §6
- `glossary_terms = {term: definition}` from SOLUTION §2
- `wire_rows = [...]` from SOLUTION §7

### Step 2 — Decompose each phase into Stories

For each phase, produce 2–6 stories. Rules:

- Each story is a vertical slice through the phase's scope.
- Each story has clear acceptance criteria (Given/When/Then or bullet list).
- Title is user/system-centric, not technology-centric.
- Story body uses glossary terms verbatim; never synonyms.
- If a phase needs >6 stories → phase is oversized; note in Planning notes,
  still emit but flag.
- If a phase produces only 1 story → phase is undersized; note, still emit.

Number stories `<phase>.<n>`: `1.1`, `1.2`, `2.1`.

### Step 3 — Decompose each Story into Subtasks

For each story, produce 1–5 subtasks. Rules:

- Each subtask ≈ one PR.
- Imperative title: `"Add users table migration"`, `"Wire POST /auth/signup"`.
- Subtasks within a story are ordered top-to-bottom (implicit dependency).
- Subtask body: Context + Area + Acceptance + (optionally) Verify line.
- If a story needs >5 subtasks → split the story.

Number subtasks `<phase>.<story>.<n>`.

### Step 4 — Sizing

Each story gets a size:

- **S** — fits in an afternoon (1 subtask, maybe 2)
- **M** — fits in a day (2–3 subtasks)
- **L** — multiple days (4+ subtasks, or one subtask of unusual depth)

Sizing is sanity-check, not sprint-planning poker. **Active check:** if every
story in a phase is S, the phase is mis-sized — surface in Planning notes.

### Step 5 — Dependencies

At **story level only**. Subtask order is implicit.

Each story declares `Depends on:` — either other story IDs (`1.2, 2.1`) or `—`.

**DAG validation:** dependencies must form a DAG. If you detect a cycle, stop
and report under Planning notes. **No story may depend on a story in a later
phase**, since phase ordering already implies that — flag as chain corruption.

### Step 6 — Planning notes (at top of BACKLOG.md)

Enumerate:

- Oversized phases (>6 stories, or many L stories)
- Undersized phases (1 story, or all-S stories)
- Forced vertical slices — where the decomposition felt artificial
- Missing context from SOLUTION.md (soft-mode warnings)
- Detected dependency anomalies
- UNRESOLVED stories with owner

If nothing notable: `_None — backlog decomposition was clean._`

**Never modify PHASES.md.** Notes are advisory; the human decides whether to
re-run `phase-planner`.

### Step 7 — Self-critique (convergence checklist)

Run the **Convergence checklist** (below). Contract mode: TodoWrite tasks per
line item. Standard: lightly. Vibe: skipped.

### Step 8 — Handoff

Output the message at the bottom of this skill. Do not create GitHub Issues
without the user confirming. Do not run a subtask automatically (that requires
the user picking a subtask and dispatching it via `work-on-issue`).

## Output template

```markdown
# Backlog: <project title from PHASES.md>

**Version:** v0.1
**Source:** docs/PHASES.md (v<X>), docs/SOLUTION.md (v<Y>)
**Generated:** <ISO date>
**Mode:** Standard | Contract

## Planning notes

- (Empty if clean.)
- **WARNING:** SOLUTION.md missing §13 ADR roster — story bodies cannot cite ADRs. (Soft-mode example.)
- **UNRESOLVED:** Story <id> — owner: <name>, blocker: <one line>.

---

## Phase 1: <name from PHASES.md>  (Epic)

**Phase status:** LOCKED (per PHASES.md)
**Phase goal:** <copied verbatim from PHASES.md>
**Phase acceptance:** <copied verbatim from PHASES.md>

### Story 1.1: <user/system-centric title>

**Status:** LOCKED
**Size:** S | M | L
**Depends on:** —
**Implements (from SOLUTION.md):** §6 entities: `Customer`, `UserCustomerAccess`; §7 wire rows: pagination, error shape; §2 glossary: Customer, Consultant.
**Description:** <2–4 sentences. What this story delivers and why. Uses glossary terms verbatim.>
**Acceptance:**
- <criterion 1 — testable>
- <criterion 2 — testable>

#### Subtask 1.1.1: <imperative title>

**Context:** <one or two sentences. Reference the SOLUTION.md section that authoritatively answers any ambiguity. E.g. "See SOLUTION.md §6 for the Customer entity fields and §7 for the error response shape.">
**Area:** <folder / module / surface hint — or "—" if cross-cutting>
**Acceptance:**
- <criterion>
- <criterion>
**Verify (soft-mode only):** <what to confirm against the codebase before starting — omit in strict mode>

#### Subtask 1.1.2: <imperative title>

[same shape]

### Story 1.2: <title>

[same shape]

---

## Phase 2: <name>  (Epic)

[same structure for each phase]

---

**Change log**

| Version | Date | Author | Change |
| ------- | ---- | ------ | ------ |
| v0.1    |      |        | Initial draft |
```

## Convergence checklist

**Upstream**
- [ ] PHASES.md was read in full and is convergence-checklist-clean.
- [ ] SOLUTION.md was read in full and is convergence-checklist-clean.
- [ ] PHASES.md and SOLUTION.md versions cited in Source field.
- [ ] No story persists an entity marked DEFERRED in SOLUTION §6.
- [ ] No story references an ADR not in SOLUTION §13.

**Story shape**
- [ ] Every phase has 2–6 stories.
- [ ] Every story title is user/system-centric (`<role> can <action>` or `<system> <behavior>`).
- [ ] Every story has Status, Size, Depends-on, Implements, Description, Acceptance.
- [ ] Story bodies use glossary terms verbatim (search for synonyms).

**Subtask shape**
- [ ] Every story has 1–5 subtasks.
- [ ] Every subtask title is imperative.
- [ ] Every subtask body has Context, Area, Acceptance.
- [ ] Subtask Context cites the SOLUTION.md section that resolves any ambiguity.
- [ ] No subtask body assumes context from a sibling subtask — each stands alone.
- [ ] Every subtask is executable by `work-on-issue` with zero questions: acceptance criteria are measurable, every decision is made, no vague language (`docs/product/definition-of-done.md`).

**Decision state**
- [ ] Every story carries Status (LOCKED / DEFERRED / UNRESOLVED).
- [ ] UNRESOLVED stories have owner in Planning notes.

**Sizing**
- [ ] Every story has S / M / L.
- [ ] No phase is all-S (would imply mis-sized phase).
- [ ] No story has >5 subtasks (would imply mis-sized story).

**Dependencies**
- [ ] All story dependencies are story IDs in this BACKLOG.md.
- [ ] Dependencies form a DAG (no cycles).
- [ ] No story depends on a story in a later phase.

**Versioning**
- [ ] BACKLOG.md has Version + Source (citing PHASES.md and SOLUTION.md versions) + change log.

## Anti-patterns

- **Skipping upstream validation.** A BACKLOG.md built on an UNRESOLVED PHASES.md inherits ambiguity and amplifies it 5×.
- **Re-deciding phase boundaries.** PHASES.md owns that. If a phase is wrong, escalate to `phase-planner`.
- **Synonyms for glossary terms.** "Proposal" instead of "Draft" in a subtask body — agents grep for "Draft", miss "proposal", diverge.
- **Subtask bodies that depend on sibling subtasks.** Each subtask is handed independently to `work-on-issue` as its own GitHub issue. Self-contain.
- **Technology-flavored story titles.** "Add users table" is a subtask, not a story. Story is "Consultant can sign up".
- **Creating issues unprompted.** Not this skill's job — the user confirms, then `gh issue create` per subtask.
- **Modifying PHASES.md.** Planning notes are advisory only.
- **Sizing in story points or hours.** S/M/L is qualitative on purpose.
- **Subtask >5 ceiling violations.** Five is the maximum, not a target. If you're at five, ask whether the story should split.

## Downstream

A subtask becomes a GitHub issue (user-confirmed `gh issue create`), which
`work-on-issue` then executes. Map subtask fields onto the issue template
(`.github/ISSUE_TEMPLATE/feature.md`):

- Subtask title → issue title.
- `**Context:**` + `**Implements (from SOLUTION.md):**` → issue Context.
- `**Acceptance:**` → issue Acceptance criteria (must be measurable, or
  `work-on-issue` rejects the issue as underspecified).
- `**Area:**` → a hint in Solution approach / Notes.

The issue must satisfy `docs/product/definition-of-done.md` discipline at creation time,
not after.

## Handoff message

After writing BACKLOG.md, post:

```
BACKLOG.md vytvořen (v0.1, mode: <Standard|Contract>, source: PHASES.md v<X> + SOLUTION.md v<Y>).

Convergence checklist: <PASS — N/N> | <FAIL — list failed line items>

<N> phases → <M> stories → <K> subtasks
Sizes: <S count> S, <M count> M, <L count> L

UNRESOLVED stories (blokuje execution) — <N> položek:
1. Story <id>: <one line> — owner: <name>
   (Resolve teď / Defer to later phase / Accept block)

Planning notes:
- <note 1>
- <note 2>

Next steps:
- Review BACKLOG.md (tvůj checkpoint)
- Vytvořit issues: potvrď a pustím `gh issue create` per subtask (acceptance criteria měřitelná)
- Execution: `work on issue <N>` na kterýkoli vytvořený issue
```

## Length scaling

- **Vibe:** skip — vibe projects don't need backlog.
- **Standard:** BACKLOG.md with all sections, convergence checklist light.
- **Contract:** full discipline — convergence checklist as TodoWrite, story-to-§section mapping enforced, DAG actively validated.
