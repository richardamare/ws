# Conventions

How code is written in `ws`. Stack: `../architecture/stack.md`. Per-area detail: `../patterns/`.

## General

- **Idempotency first.** Every command checks actual state and acts only on drift. `up`/`auth`/`new`
  are safe to run repeatedly. Most of the real logic lives here.
- **Thin commands, testable internals.** `cmd/` files parse flags and call `internal/<area>`. No logic
  in Cobra handlers beyond wiring.
- **One output layer.** Never `fmt.Println` a result directly — route every result through
  `internal/output` so pretty / `--plain` / `--json` stay consistent. See `output-modes.md`.
- **Shell out to tools**, parse stdout; don't reimplement `cmux`/`az`/`docker`.

## Safety (see `../security/README.md`)

- The per-project SP is **Reader-only**. Code must never request more, nor use it for write.
- Write/Terraform flows go through `ws elevate` (personal login) with explicit confirmation.
- Mirror the `.claude/settings.json` deny list in code: confirm-or-refuse destructive Azure/git ops.

## Config

- Per-project YAML under `~/.config/ws/projects/`. Hand-editable; `ws` round-trips it without
  destroying comments where practical. No database.
- Optional blocks (`azure`, `sessions`, `container`) absent = feature off.

## Style

- Go conventions in `../patterns/go.md`. `gofmt` + `go build ./...` before commit.
- No comments except a non-obvious *why*.
