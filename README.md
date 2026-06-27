# ws

A small Go CLI that sets up a per-project developer workspace in one command: opens the cmux workspace
+ tabs, logs into Azure as a **scoped Reader service principal**, and tracks Claude Code session IDs.

```bash
ws new proj1     # once: create the Reader SP + cert + config
ws up proj1      # daily: workspace + scoped login + tabs
ws resume proj1  # reopen a bookmarked Claude session
```

- **Local config only**, one YAML per project (`~/.config/ws/projects/`). No database.
- **Safe by default** — sessions are Reader-only on one resource group; write only via deliberate
  `ws elevate`. See [`docs/security/README.md`](docs/security/README.md).
- **Human- and LLM-friendly** — interactive pickers (huh) for humans, structured text / `--json` for
  agents.

## Docs

Start with [`CLAUDE.md`](CLAUDE.md) (operational rules) and the [`docs/`](docs/README.md) index:
architecture, product, security, patterns.

> Status: design complete (2026-06-27), implementation not started. See
> [`docs/roadmap.md`](docs/roadmap.md).
