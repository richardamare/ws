# CLI — Cobra + huh

## Structure

- **Cobra** owns the command tree, flags, help. One file per command in `cmd/`.
- **huh** (Charm) provides interactive selects/confirms — the user wants minimal typing.

## The dual-mode rule

Every command works **typed** or **interactive**:

- `ws resume proj1 auth-refactor` — fully specified, runs immediately.
- `ws resume` — missing args → huh pickers (project, then session) + confirm.

Resolution order for each argument: explicit arg → (if TTY and not `--json`) interactive picker →
(else) error. Never block on a prompt under `--json` or non-TTY.

```go
var choice string
huh.NewSelect[string]().
    Title("Resume which session?").
    Options(
        huh.NewOption("auth-refactor — RBAC context", "3ee3..."),
        huh.NewOption("terraform-bootstrap — infra", "9ab1..."),
    ).
    Value(&choice).
    Run()
```

## Global flags (persistent)

- `--json` — strict JSON, disables prompts.
- `--plain` — structured text; auto-on when stdout is not a TTY.

Detect TTY once at startup (`term.IsTerminal`) and thread the chosen output mode through a single
`internal/output` renderer. See `../product/output-modes.md`.

## Destructive commands

`rm`, `rotate`, `elevate` and anything touching Azure write require an explicit confirm (huh `Confirm`)
on a TTY, and must be passed `--yes` to run under `--json`/non-TTY. Mirror the `.claude/settings.json`
deny list — never auto-run `az group delete`, role-assignment changes, or `terraform apply`.
