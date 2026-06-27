# Git

- **Never push to `master`.** `master` is always buildable.
- **Branch:** `feat/issue-<n>-<slug>`, `fix/issue-<n>-<slug>`, `chore/<slug>`.
- **Commit format:** `type(scope): subject (#<issue>)` — Conventional Commits + issue link.
  Atomic commits, one concept each.
- **Before commit:** `gofmt -w` and `go build ./...` pass.
- **Never** `--no-verify`, `--amend` published commits, or force-push `master` (also denied in
  `.claude/settings.json`).
- **PR body references the issue** — `Closes #<n>` (complete) or `Refs #<n>` (partial).
- **Before ending a session**, post a Done / In progress / Blocked / Next comment on the active issue.
