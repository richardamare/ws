# Go

Conventions for `ws`.

- **Single static binary.** No CGO unless forced; `go build ./...` produces a droppable executable.
- **Shell out, don't reimplement.** `cmux`, `az`, `docker`, `gh`, `claude` are driven via `os/exec`.
  Parse their stdout; never reimplement their protocols. (cmux socket is a last-resort fallback.)
- **`gofmt -w` + `go build ./...` pass before every commit.** Non-negotiable (CLAUDE.md).
- **Errors:** wrap with `fmt.Errorf("...: %w", err)`; return them, don't `log.Fatal` deep in the call
  tree. The top-level command handler renders the error via the output layer (`internal/output`).
- **No comments** except a non-obvious *why*.
- **Standard library first.** Dependencies limited to Cobra, huh, a YAML parser. Justify additions.
- **Idempotency is the design center.** Every action = "check actual state → act only on drift."
  `up`, `auth`, `new` must all be safe to run repeatedly.
- **Layout:** commands in `cmd/`, logic in `internal/<area>/`. Commands are thin; logic is testable.
