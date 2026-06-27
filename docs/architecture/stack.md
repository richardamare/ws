# Stack

## Shape

`ws` is a **thin orchestrator** over three things it does not own:

- **cmux** (the terminal, v0.64.16) — windows/workspaces/panes/surfaces, session restore, agent resume.
- **Docker / devcontainer** (optional) — container lifecycle, only when a project opts in.
- **Azure CLI / gh** — auth; `ws` drives `az login` / `gh` and checks state.

Its value is the 20% the substrate doesn't cover: per-project config, scoped Azure auth, session
bookmarks, one-command startup.

## Language & libraries

| Concern | Choice |
| --- | --- |
| Language | **Go** — single static binary, trivial distribution, good at shelling out |
| Commands / flags | **Cobra** |
| Interactive pickers / confirms | **huh** (Charm) |
| Richer TUI (future, optional) | bubbletea (Charm) — not needed for v1 |

Go was chosen over the usual Bun + Effect-TS stack (see `adr/0001-go-language.md`) because this is
mostly subprocess orchestration where a static binary wins.

## Layout (intended)

```
ws/
  main.go
  cmd/                 # Cobra commands: new, up, down, ls, status, auth, rotate, elevate, rm, sessions, resume, save
  internal/
    config/            # load/save ~/.config/ws/projects/<name>.yaml
    cmux/              # render cmux.json template + live `cmux` calls
    azure/             # az login as Reader SP, AZURE_CONFIG_DIR isolation
    docker/            # optional container lifecycle
    session/           # Claude session bookmarks
    output/            # pretty / structured-text / json render layer
```

## Storage

```
~/.config/ws/
  projects/<name>.yaml     # one file per project, hand-editable
  certs/<name>.pem         # SP cert, chmod 600
```

No database. Per-project `AZURE_CONFIG_DIR` isolates each Reader login from the personal admin login.

## Commands

Full surface in `../product/workflow.md`. Summary:
`new · up · down · ls · status · rm · auth · rotate · elevate · sessions · save · resume`.
