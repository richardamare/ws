# ADR-0001: Go for the ws CLI

- **Status:** Accepted
- **Date:** 2026-06-27
- **Deciders:** ws authors

## Context

`ws` orchestrates external tools (`cmux`, `az`, `docker`, `gh`, `claude`): spawn a process, wait for it
to be healthy, branch on exit code, render output. The maintainer's usual stack for personal CLIs is Bun +
Effect-TS. Effect shines for deep async error-recovery trees; this tool is moderate orchestration plus
a strong need for trivial, dependency-free distribution. The binary should drop onto any machine with
no runtime so a project can be started instantly.

## Decision

We will build `ws` in **Go**, using **Cobra** for the command tree and **huh** for interactive pickers.

## Consequences

- Easier: distribution — a single static binary, no runtime. `os/exec` + `net` cover everything ws needs.
- Easier: shelling out to external tools, which is the bulk of the work.
- Harder / accepted: diverges from the Bun + Effect-TS default, so less shared code with other personal
  tools. Acceptable because this is subprocess glue, not complex async domain logic.
- Ruled out for now: Bun + Effect-TS. Revisit only if `ws` grows into a long-running daemon with
  complex concurrent state.
