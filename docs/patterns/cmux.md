# cmux integration

cmux is the terminal `ws` drives (installed: v0.64.16). `ws` does **not** reimplement cmux — it uses it.

## What cmux already gives us (don't rebuild)

- Window / workspace / pane / surface layout, restore on relaunch, scrollback, browser URL.
- **Agent resume** for Claude Code / Codex after a crash/close — via `cmux hooks setup` +
  `terminal.autoResumeAgentSessions: true`.
- Per-project **workspace templates** in `~/.config/cmux/cmux.json` `commands[]` (a layout tree of
  panes → surfaces with per-terminal startup `command` and per-browser `url`).

## Hybrid control (what `ws` does)

1. **Generate** a `cmux.json` `commands[]` entry from the project config → durable: cmux's own restore
   reopens the tabs after a crash without `ws` running.
2. **Drive live** at `up`: scoped `az login`, attach container (if any), open/focus the workspace.

## Useful cmux verbs (shell out to the `cmux` binary)

```bash
cmux ssh root@localhost --port <n> --name <proj>       # if a project uses SSH
cmux new-surface --type terminal --command "claude"
cmux new-surface --type terminal --command "docker exec -it <c> zsh"  # local container, no sshd
cmux new-surface --type browser  --url "https://github.com/me/proj1"
cmux new-workspace --name <proj> --cwd <path> --command "..."
cmux surface resume show        # focused agent's resume id (session bookmarks)
cmux config validate            # before reload
cmux reload-config              # apply cmux.json + ghostty config, no restart
```

Raw socket fallback: `~/.local/state/cmux/cmux.sock`, `cmux rpc <method>`, `cmux capabilities`.

## Container reachability (optional projects)

When a project opts into a dev container with docker-compose, add the dev container as a service in the
compose file (or join the project's external network) so it reaches `db`/`api`/etc. by name. Attach via
`docker exec` (no sshd). Note: `cmux vm` is the **cloud** product (hosted, billed); local work uses
plain Docker + `cmux` tabs that `docker exec` in.
