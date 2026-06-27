---
name: work-on-issue
description: |
  Execute the full development cycle for a single GitHub issue in one go, with
  ZERO clarifying questions. The issue is the contract: everything the work
  needs must already be written there. Reads the issue, branches, implements
  against the acceptance criteria, verifies each one, runs /code-review, opens a
  DoD-compliant PR, and posts a progress comment.

  Triggers: "work on issue N", "implement issue N", "do issue N", "start issue N",
  "udělej issue N", "vyřeš issue N", or any phrasing that names a single issue
  number to be taken from open to PR.

  Does NOT trigger for: planning ("plan issue N"), questions about an issue,
  multi-issue requests, or work with no issue number. If the issue is too thin
  to execute, the skill comments what's missing and stops — it never negotiates
  requirements in chat.
---

# Work On Issue Skill

Take one GitHub issue from open to PR in a single autonomous pass. **No
clarifying questions.** The instant requirements get discussed in chat instead
of read from the issue, the work stops being reproducible and quality drops —
so this skill treats the issue as the sole source of truth and pushes any gap
back into the issue, never into conversation.

This is the enforcement mechanism for `docs/product/definition-of-done.md`: a thin issue
cannot produce good output, so the skill refuses to code against one.

## Hard rules

1. **Ask the user nothing.** No clarifying questions, no "should I…?", no
   confirmations mid-run. Decide from the issue or stop.
2. **The issue is the contract.** If a decision isn't answerable from the issue
   body, the acceptance criteria, or this repo's docs/conventions, the issue is
   underspecified — go to the *Underspecified issue* path.
3. **Obey every non-negotiable rule in `CLAUDE.md`** — never push to `master`,
   branch naming, atomic issue-linked commits, `bun run gen:api` on contract
   changes, never `--no-verify`/`--amend`/force-push.
4. **The PR must satisfy `docs/product/definition-of-done.md`** before you stop.

## Procedure

Create a TodoWrite item per step and work them in order.

### 1. Read the issue

```bash
gh issue view <N> --json number,title,body,labels,comments
```

Extract: acceptance criteria, context, solution approach, dependencies, notes.
Read any `docs/SOLUTION.md` section the issue links. Read the repo docs relevant
to the work area (the `CLAUDE.md` "Where things live" table points to them).

### 2. Gate: is the issue executable?

The issue is executable only if **all** hold:

- It has an **Acceptance criteria** section with at least one criterion.
- **Every criterion is measurable** — checkable by a test, a command, or a
  browser observation. "Works well", "looks good", "is fast" are not.
- No criterion depends on a decision that isn't made in the issue or repo
  conventions (e.g. "use the right colour" with no value, "add the field" with
  no name/type).

If any fails → *Underspecified issue* path. Otherwise continue.

### 3. Branch

Use a worktree per `docs/product/workflow.md`:

```bash
git worktree add ./.claude/worktrees/feat-issue-<N> -b feat/issue-<N>-<slug>
```

Pick `feat/`, `fix/`, or `chore/` from the issue's label/type. Slug from the
title. Never edit from the main checkout.

### 4. Implement

- Follow repo conventions (`docs/product/conventions.md`, `docs/patterns/*`).
- TDD where it pays off — especially API tests against real Postgres.
- One concept per commit, message `type(scope): subject (#<N>)`.
- API contract change (DTO, endpoint, enum) → `bun run gen:api` in the **same**
  commit.
- Keep context lean: delegate broad code search to subagents, don't read the
  whole tree into context (see `docs/product/workflow.md` → Context discipline).

### 5. Verify every acceptance criterion

For each criterion, produce evidence: a passing test, command output, or a
browser check (start the dev server for UI changes — don't assume from the
diff). Record what proved each one; it goes in the PR body.

### 6. Review gate

Run `/code-review` on the diff. Fix every finding, or dismiss it with a written
reason. This replaces the human reviewer — it is not optional.

**Grounding pass.** `/code-review` alone can cite a line that no longer exists or
flag a change without knowing why it was made. Before trusting its findings, run
two checks — in parallel subagents when the diff is non-trivial:

- **Code-grounder** — for every finding, confirm the cited `file:line` still
  exists and the quoted code matches the current diff. A finding that cites stale
  or wrong code is discarded, not fixed.
- **Historian** — `git blame` the changed lines and search for any `docs/` ADR or
  prior issue/PR touching them. If the diff reverses a deliberate earlier
  decision, that is a blocker: stop and surface it, don't silently re-undo it.

Each check returns a confidence `0–100`; take the **minimum** of the two as the
diff's grounding score. Anchor each confidence to the strongest evidence actually
in hand, ranked: executable reproduction → source inspection → dependency
docs/types → stale comment or old CI. Stale comments and old behavior are hints
until rechecked. When the trail is too weak to verify a finding, label it
**"not proven"** and score it low rather than guessing high — a forced verdict on
a weak trail is itself a false positive. **Trigger the grounding pass when** the
diff is ≥10 changed lines OR carries any of the labels `security`, `migration`,
`payments`, `auth`. Below the trigger, a single inline grounding read is enough.

If the grounding score is below 70, do **not** treat the change as merge-ready:
record the lowest-confidence finding and resolve it (fix, or dismiss with a
written reason) before step 7. Note the score and how findings were handled in
the PR body.

### 7. PR

```bash
git push -u origin feat/issue-<N>-<slug>
gh pr create --fill-first
```

Body follows `.github/PULL_REQUEST_TEMPLATE.md`: each acceptance criterion with
its verification, the `/code-review` result, test coverage, `Closes #<N>`
(complete) or `Refs #<N>` (partial). The PR must meet every item in
`docs/product/definition-of-done.md`.

### 8. Progress comment

Post to the issue: **Done** / **In progress** (exactly where) / **Blocked** /
**Next**, plus the PR link.

### 9. Report

Tell the user what shipped, the PR URL, and any criterion that couldn't be
verified (with why). Do not claim done for anything unverified.

## Underspecified issue path

Do **not** ask the user in chat. Instead:

1. Post a single issue comment listing precisely what's missing — which
   acceptance criteria are absent or unmeasurable, which decisions aren't made,
   what inputs (values, names, types, repro steps) are needed to proceed.
2. Stop. Do not create a branch or write code.
3. Tell the user, in one line: the issue is underspecified, a comment listing
   the gaps was posted, re-run `work on issue <N>` once the issue is updated.

The comment must be specific enough that editing the issue to answer it makes
the issue executable — no generic "please add more detail".

## Why no questions

A self-contained issue is reproducible: any agent or person, today or in three
weeks, gets the same result. A chat conversation is not — it evaporates, it
drifts, and it lets vague requirements through because the human fills gaps in
their head. Forcing every gap back into the issue keeps the contract intact and
the output quality high. That trade — slower issues, faster and better
execution — is the entire point of this repo's workflow.
