# Definition of Done

The pre-merge quality bar. With no human reviewer, `/code-review` is the gate.

A change is done when:

1. **Builds & formatted** — `go build ./...` clean, `gofmt` applied, `go vet ./...` clean.
2. **Tested** — new logic has tests; `go test ./...` passes. Idempotent paths (`up`/`auth`/`new`) have
   "run twice = no-op" coverage.
3. **Both output modes work** — the command renders correctly as pretty, `--plain`, and `--json`;
   missing input under `--json` errors instead of prompting.
4. **Safe** — no new path can use the Reader SP for write, change role assignments, or auto-run
   `az group delete` / `terraform apply`. Destructive ops confirm or require `--yes`.
5. **Issue-linked** — branch and commits reference the GitHub issue; PR body has `Closes #<n>` /
   `Refs #<n>`.
6. **Docs updated** — if behavior/config/commands changed, the relevant `docs/` file is updated in the
   same PR.
7. **`/code-review` run** — every finding fixed or dismissed-with-reason in the PR body.
8. **Progress comment** posted to the issue (Done / In progress / Blocked / Next).
