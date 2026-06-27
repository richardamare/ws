# Workflow — command surface

Name is `ws` (short on purpose). **Every command works two ways:**

- typed/fast: `ws resume proj1 auth-refactor`
- interactive: `ws resume` → pick project → pick session (huh select + confirm)

Missing input on a TTY → interactive picker. Missing input under `--json` → error (no prompts). See
`output-modes.md`.

## Lifecycle

| Command | Does |
| --- | --- |
| `ws new <name>` | Scaffold a project: prompt for RG/sub/repo/tabs, **create the Reader SP + cert**, write `<name>.yaml`. Idempotent (reuse existing SP). |
| `ws up <name>` | Start working: scoped `az login` (idempotent) → open cmux workspace + tabs. |
| `ws down <name>` | Close the project's cmux workspace (and container if any). |
| `ws ls` | List projects + status. |
| `ws status <name>` | Show one project (azure login state, tabs, container). |
| `ws rm <name>` | Remove project; optionally delete the SP + role assignment + cert (keep Entra clean). |

## Azure

| Command | Does |
| --- | --- |
| `ws auth <name>` | Re-login the Reader SP; detect expired cert. |
| `ws rotate <name>` | Rotate the SP cert before/after expiry (~1yr). |
| `ws elevate <name>` | Open a marked elevated tab: `az login` as yourself for write/Terraform. **Never uses the SP.** See `../security/README.md`. |

## Sessions (Claude Code bookmarks)

| Command | Does |
| --- | --- |
| `ws sessions <name>` | List curated session bookmarks (label + note). |
| `ws save <name> <label>` | Bookmark the current Claude session id (read via `cmux surface resume show`). |
| `ws resume <name> <label>` | `claude --resume <id>` for that bookmark. |

Detail in `sessions.md`.

## A project's life

```
ws new proj1        # once: SP + cert + config
ws up proj1         # daily: workspace + scoped login + tabs
ws save proj1 x     # bookmark a good session
ws down proj1       # end of day
ws elevate proj1    # rare: deliberate admin for Terraform
```

## Global flags

- `--json` — strict JSON output, disables all interactive prompts.
- `--plain` — structured text (auto-on when stdout is not a TTY; LLM default).
- `ws --json schema` — dump all output shapes (self-documenting for LLMs).
