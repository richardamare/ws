# ADR-0003: Hybrid cmux control — generated template plus live calls

- **Status:** Accepted
- **Date:** 2026-06-27
- **Deciders:** ws authors

## Context

cmux owns workspace layout, restore, and agent resume. `ws` needs to open a project's tabs and keep
them recoverable after an accidental close or crash, while also doing dynamic per-launch work (scoped
`az login`, attaching a container). Two pure approaches exist: generate a static `cmux.json` template,
or drive everything live via the cmux CLI/socket. Each alone is insufficient — a static template can't
do per-launch auth, and pure live control loses the tabs if `ws` isn't running when cmux restores.

## Decision

We will use a **hybrid**: generate a `cmux.json` `commands[]` workspace template per project for
durable restore, and drive only the dynamic parts (scoped `az login`, container attach, focus) live at
`ws up` by shelling out to the `cmux` binary.

## Consequences

- Easier: after an accidental close or crash, cmux's own restore reopens the tabs without `ws` running.
  Dynamic state (auth, container) is still handled at launch where a static template can't.
- Harder / accepted: `ws` must keep the generated template in sync when a project's `tabs` change.
- Follow-on: the raw cmux socket (`cmux rpc`) remains a fallback only, used if the CLI can't express
  something.
