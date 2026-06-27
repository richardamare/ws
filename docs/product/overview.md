# Overview

`ws` is a small Go CLI that sets up an entire per-project developer workspace in one command:
opens the cmux workspace + tabs, logs into Azure as a **scoped, Reader-only** service principal,
and tracks Claude Code session IDs so good-context sessions can be reused instead of kept open.

## Goals

- **One command to start working** on a project: `ws up <project>`.
- **Local config only** — one YAML/JSON file per project. No database, no server.
- **Safe by default** — everyday sessions authenticate to Azure as a Reader SP scoped to a single
  resource group; they cannot delete or modify anything, anywhere. See `../security/README.md`.
- **Short to type** — name is `ws`; every command also works interactively (pick from a list).
- **LLM-friendly output** — structured text by default for agents, `--json` when strict parsing is needed.

## What it does per project

1. **Workspace setup** — generate/drive a cmux workspace with predefined tabs (Claude Code terminal,
   manual shell, browser tab on the GitHub repo / docs).
2. **Scoped Azure login** — `az login` as that project's Reader SP into an isolated `AZURE_CONFIG_DIR`.
3. **Session bookmarks** — a short, named list of Claude session IDs to resume on demand.

## Non-goals

- Not a reimplementation of cmux (cmux owns layout, restore, agent resume).
- Not a container manager (dev containers are **optional** per project, not required).
- Not a secrets store (Azure tokens live in the CLI's own credential cache; certs on disk, `chmod 600`).

## Status

Design complete (2026-06-27). Next: scaffold the Cobra + huh skeleton. See `../roadmap.md`.
