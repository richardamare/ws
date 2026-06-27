# Go

Conventions for `ws`.

## Layout — interface in `cmd/`, logic in `internal/`

```
cmd/ws/                 # the application interface (CLI). package main.
  main.go               #   entrypoint — calls Execute()
  root.go               #   cobra root, global --json/--plain, format/TTY plumbing
  version.go ls.go up.go#   one file per command; thin — parse input, call internal, render
  picker.go             #   huh interactive pickers (TTY only)
internal/               # business logic / services. NO cobra, NO interface concerns.
  config/               #   per-project YAML store (load/save/list)
  workspace/            #   compute the `up` plan from a project (pure, testable)
  output/               #   render results: pretty / plain / json
  azure/ cmux/ session/ #   (planned) adapters that shell out to az / cmux / claude
```

**Rule:** the cobra command tree and any interactive UI live in `cmd/ws/`. `internal/` holds only
services and domain logic and must not import cobra/huh. Commands are thin; logic is in `internal/` so
it's testable without a terminal.

## Style

- **Single static binary.** `go build ./cmd/ws`.
- **Shell out, don't reimplement.** `cmux`, `az`, `docker`, `gh`, `claude` are driven via `os/exec`.
  Parse their stdout; never reimplement their protocols. (cmux socket is a last-resort fallback.)
- **Errors:** wrap with `fmt.Errorf("...: %w", err)` and return them; the `cmd/ws` layer renders the
  error and sets the exit code. Don't `log.Fatal` deep in the tree.
- **No comments** except a non-obvious *why*.
- **Standard library first.** Deps: cobra (commands), huh (pickers), yaml.v3 (config), x/term (TTY).
  Justify additions.
- **Idempotency is the design center.** Every action = "check actual state → act only on drift."
  `up`, `auth`, `new` must all be safe to run repeatedly.

## Output

Never `fmt.Println` a result. Route every result through `internal/output` (`Record` for one object,
`Table` for rows) so pretty / `--plain` / `--json` stay consistent. Format is resolved once in
`cmd/ws/root.go` from the flags + TTY state. See `../product/output-modes.md`.

## Interactive vs scripted

`cmd/ws/picker.go` provides huh selectors, used only when `interactive()` is true (TTY and not
`--json`). Missing input under `--json`/non-TTY is an error, never a prompt.

## Testing

- Pure logic packages (`config`, `workspace`, `output`) have table-driven unit tests; no TTY needed.
- Command tests build the root via `newRootCmd()`, `SetArgs(...)`, `SetOut(buf)` and assert on rendered
  output. Pass flags as args (`--json`) — defining a cobra flag resets its bound variable to the
  default, so setting the global directly does not stick.
- `go test -race -count=1 ./...` is the bar (also what CI runs).

## Tooling

- `make check` = `go vet` + race tests + gofmt check (mirrors CI).
- `make build` → `bin/ws`. `make fmt` formats. `make install` installs `ws`.
- CI: `.github/workflows/ci.yml` runs gofmt-check, vet, build, and `go test -race` on PRs and pushes to
  `master`.
- `gofmt -w` + `go build ./...` pass before every commit (CLAUDE.md).
