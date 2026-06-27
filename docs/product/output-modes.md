# Output modes

`ws` must be usable by a human **and** by an LLM (Claude). Three render modes, one internal struct per
command result — no duplicated logic.

| Mode | Flag | When | For |
| --- | --- | --- | --- |
| Pretty | (default on TTY) | interactive terminal | human — tables, color, huh pickers |
| **Structured text** | `--plain` / **auto when non-TTY** | piped / agent | **LLM default** — labeled blocks or TSV |
| JSON | `--json` | needs strict parsing | machines / `jq` |

## Why not JSON-always for LLMs

For many cases **flat labeled text / TSV beats JSON** for an LLM: fewer tokens, no escaping or nesting
noise, easier to read. So the agent default is **structured text**, with `--json` available when
something genuinely needs to parse it.

## Structured text examples

List → aligned columns / TSV:

```
ws sessions proj1 --plain

LABEL              ID         NOTE
auth-refactor      3ee3...    RBAC context
terraform-boot     9ab1...    infra setup
```

Single record → key:value labels:

```
ws status proj1 --plain

name:    proj1
rg:      rg-proj1
sp:      sp-ws-proj1-reader
azure:   logged-in (Reader)
status:  up
```

## JSON

```bash
ws ls --json
# [{"name":"proj1","rg":"rg-proj1","sp":"sp-ws-proj1-reader","status":"up"}]
```

## Rules

- `--json` and non-TTY **disable interactive prompts**; missing input becomes an error:
  `{"error":"label required","code":"missing_arg"}` (stderr, non-zero exit).
- Keep the JSON **schema stable** so the LLM can rely on it.
- `ws --json schema` dumps every output shape.
