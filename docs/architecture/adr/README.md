# Architecture Decision Records

One file per architecturally significant decision, in
[Michael Nygard's format](https://github.com/joelparkerhenderson/architecture-decision-record/blob/main/locales/en/templates/decision-record-template-by-michael-nygard/index.md).

## Conventions

- Filename: `NNNN-short-title.md`, zero-padded sequential (`0001-`, `0002-`…).
- Numbers are never reused; a superseded ADR stays in place with status
  `Superseded by [ADR-NNNN](NNNN-….md)`.
- Status is one of: `Proposed` · `Accepted` · `Deprecated` · `Superseded`.
- Keep it short. An ADR captures *one* decision and its consequences, not a
  design doc — that's `docs/SOLUTION.md`.

## Authoring

Copy [`0000-template.md`](0000-template.md) to the next number and fill it in.
The `phase-planner` skill names which ADRs a phase authors; write them as that
phase lands.

When a decision is settled, set status to `Accepted` and link it from the
relevant `docs/SOLUTION.md` section and any domain folder it affects.
