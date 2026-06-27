# Per-project config schema

One file per project, stored locally, hand-editable. No database.

```
~/.config/ws/projects/<name>.yaml
```

## Schema

```yaml
name: proj1
cwd: ~/code/proj1

azure:                          # optional — omit for projects with no Azure work
  sp_app_id: <appId>           # the Reader service principal
  tenant: <tenantId>
  cert: ~/.config/ws/certs/proj1.pem
  config_dir: ~/.azure-proj1   # isolated AZURE_CONFIG_DIR for the scoped login
  subscription: <subId>
  resource_group: rg-proj1     # the single RG the SP is scoped to (Reader)

tabs:
  - { type: terminal, name: Claude, run: "claude" }
  - { type: terminal, name: Shell }
  - { type: browser,  name: Repo,  url: "https://github.com/me/proj1" }
  - { type: browser,  name: Docs,  url: "https://learn.microsoft.com/..." }

sessions:                       # curated Claude session bookmarks (see ../../product/sessions.md)
  - { label: auth-refactor,       id: 3ee3..., note: "RBAC context" }
  - { label: terraform-bootstrap, id: 9ab1..., note: "infra setup" }

container:                      # optional, OFF by default — only if this project uses a dev container
  compose: docker-compose.yml          # plus an overlay that adds a 'devcontainer' service
  service: devcontainer
  exec_shell: zsh
```

## Notes

- `azure`, `sessions`, and `container` blocks are all optional.
- `config_dir` is the key isolation mechanism: the project's Reader login lives here, never in the
  personal `~/.azure`. See `../../security/README.md`.
- JSON is also accepted; YAML is the default for hand-editing.
