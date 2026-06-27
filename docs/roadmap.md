# Roadmap

Pre-skill design intent. Once `solution-design` / `phase-planner` run, `SOLUTION.md` / `PHASES.md` /
`STATE.md` supersede this.

## Resolved decisions

- Name: **ws** · Language: **Go** (ADR 0001) · cmux control: **hybrid** (ADR 0003)
- Azure: **Reader SP per project**, scoped to one RG, isolated via `AZURE_CONFIG_DIR`; personal admin
  untouched; write via deliberate `ws elevate` (ADR 0002)
- Sandbox: dev containers **optional**, host-first
- Output: pretty / **structured-text (LLM default)** / `--json`
- Interactive: **Cobra + huh**

## Needed before/while coding

- Go **module path** (e.g. `github.com/richardamare/ws`).
- Subscription id + the resource-group names → generate the `az ad sp create-for-rbac` commands.
- Tenant allows `az ad sp create-for-rbac` (confirmed: can register apps + assign roles).

## Build order (MVP first)

1. **Scaffold** — Cobra + huh, global `--json`/`--plain`, output layer, config load/save.
2. **`ws up`** for one project: read YAML → generate `cmux.json` template → open tabs. **No Azure yet** —
   prove the workspace loop end to end (the "it opened everything" moment).
3. **Scoped Azure** — `ws auth`, `AZURE_CONFIG_DIR` isolation, login in tabs.
4. **`ws new`** — create SP + cert + config.
5. **Sessions** — `ws save` / `sessions` / `resume`.
6. **Lifecycle polish** — `down`, `ls`, `status`, `rm` (+ SP cleanup), `rotate`, `elevate`.

## Later / optional

- `container:` projects (compose overlay + `docker exec` tabs).
- Richer TUI dashboard (bubbletea) if huh pickers aren't enough.
- `ws --json schema` self-documentation.

## Hard constraints (do not violate)

- Never remove or modify the maintainer's existing Azure role assignments.
- The per-project SP is **Reader-only**; write only via deliberate personal `az login`.
- No database; local files only.
- Keep the CLI short to type; interactive by default.
