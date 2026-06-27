# Changelog

All notable changes to `ws` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project aims to
follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.1] - 2026-06-27

### Changed

- Republished after scrubbing personal data from the repository and rewriting
  git history. `v0.1.0` is retracted in `go.mod`; use `v0.1.1`+.


## [0.1.0] - 2026-06-27

First public release.

### Added

- Per-project workspace CLI (`cmd/ws`): `new`, `up`, `down`, `ls`, `status`,
  `auth`, `rotate`, `elevate`, `template`, `sessions`, `save`, `resume`, `rm`.
- Scoped Azure auth: each project logs in as a Reader-only service principal
  confined to a single resource group, isolated in its own `AZURE_CONFIG_DIR`.
  Write/Terraform is the deliberate `ws elevate` path only.
- cmux integration: opens the project workspace + tabs live, and generates a
  durable `cmux.json` workspace template so a crash/close can restore the tabs.
- Claude Code session bookmarks: `save` / `sessions` / `resume`.
- Output modes: pretty (TTY), structured text (`--plain`, default off-TTY),
  and strict `--json`; interactive pickers (huh) on a TTY.
- Account inference for `~/Developer/Personal` vs `~/Developer/work`.

[Unreleased]: https://github.com/richardamare/ws/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/richardamare/ws/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/richardamare/ws/releases/tag/v0.1.0
