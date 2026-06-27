# Schemas

Data and persistence schemas for `ws`.

- **`config.md`** — the per-project config file (`~/.config/ws/projects/<name>.yaml`). This is the only
  persisted data; there is no database.

`ws` stores nothing else of its own beyond config files and SP certs under `~/.config/ws/`. Azure
tokens live in the Azure CLI's own credential cache (isolated per project via `AZURE_CONFIG_DIR`).
