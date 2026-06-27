# ws

> Start working on any project with one command — your terminal workspace, a **scoped** Azure login, and your Claude Code sessions, all set up for you.

[![CI](https://github.com/richardamare/ws/actions/workflows/ci.yml/badge.svg)](https://github.com/richardamare/ws/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/richardamare/ws?sort=semver)](https://github.com/richardamare/ws/releases)
[![Go](https://img.shields.io/badge/go-1.23%2B-00ADD8?logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

`ws` is a small CLI that turns "set up my project" into `ws up`. For each project it:

- 🪟 **opens your cmux workspace** — the tabs you want (Claude Code, a shell, a browser on the repo), and writes a durable template so a crash or accidental close restores them.
- 🔐 **logs into Azure safely** — as a per-project, **Reader-only** service principal scoped to a single resource group, isolated from your personal admin account. A stray `az group delete` simply can't reach anything else.
- 🔖 **remembers your Claude sessions** — bookmark good-context sessions by name and `ws resume` them, instead of hoarding open windows.

One local YAML file per project. No database, no server.

---

## Install

**One-liner (recommended)** — builds from source with Go:

```bash
curl -fsSL https://raw.githubusercontent.com/richardamare/ws/master/scripts/install.sh | bash
```

**With Go directly:**

```bash
go install github.com/richardamare/ws/cmd/ws@latest
```

Both drop the `ws` binary in `$(go env GOPATH)/bin` — make sure that's on your `PATH`:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"   # add to ~/.zshrc
```

<details>
<summary>Prebuilt binaries (from a Release)</summary>

Each tagged release publishes macOS/Linux archives. On an Apple-Silicon Mac:

```bash
curl -fsSL https://github.com/richardamare/ws/releases/latest/download/ws_*_darwin_arm64.tar.gz -o ws.tar.gz
tar xzf ws.tar.gz ws
sudo mv ws /usr/local/bin/
```

(Use `darwin_amd64` on Intel Macs.)
</details>

### Prerequisites

`ws` drives tools you already use — have these installed and signed in:

| Tool | For |
| --- | --- |
| [Go 1.23+](https://go.dev/dl/) | building/installing `ws` (`brew install go`) |
| [cmux](https://github.com/manaflow-ai/cmux) | the terminal workspaces it opens |
| [Azure CLI](https://learn.microsoft.com/cli/azure/) (`az`) | the scoped logins (optional per project) |
| [Claude Code](https://claude.com/claude-code) (`claude`) | the agent tabs and session resume |
| [GitHub CLI](https://cli.github.com/) (`gh`) | optional, for repo workflows |

---

## Quick start

```bash
# 1. Create a project — scaffolds a Reader service principal scoped to one RG,
#    stores its cert, and writes ~/.config/ws/projects/<name>.yaml
ws new myapp --cwd ~/Developer/Personal/myapp \
             --sub <subscription-id> --rg rg-myapp \
             --repo https://github.com/me/myapp

# 2. Start working — scoped Azure login + your cmux workspace with tabs
ws up myapp

# 3. Bookmark a good Claude session, then bring it back later
ws save myapp auth-refactor
ws resume myapp auth-refactor
```

Run any command with no arguments to pick interactively (e.g. `ws up`, `ws resume`).

---

## Commands

| Command | What it does |
| --- | --- |
| `ws new <name>` | Create a project: scoped Reader SP + cert + config |
| `ws up [name]` | Scoped Azure login + open the cmux workspace (`--dry-run` to preview) |
| `ws down [name]` | Close the project's cmux workspace |
| `ws ls` | List projects |
| `ws status [name]` | Show config + Azure login state |
| `ws template [name]` | (Re)write the cmux.json workspace template |
| `ws auth [name]` | Re-login the Reader SP |
| `ws rotate [name]` | Issue a fresh SP certificate |
| `ws elevate` | Open a marked personal-admin tab for write/Terraform |
| `ws sessions [name]` | List Claude session bookmarks |
| `ws save <name> <label>` | Bookmark the current Claude session |
| `ws resume <name> [label]` | Resume a bookmarked session |
| `ws rm <name>` | Remove a project (`--purge` also deletes the cert) |

**Output for humans and agents:** pretty tables on a terminal, structured text when piped, and strict `--json` on demand. Pass `--json` to get machine-readable output and disable prompts.

---

## How the Azure safety works

Your everyday sessions log in as a **Reader service principal confined to one resource group**, kept in its own `AZURE_CONFIG_DIR` — so they can read, but never delete or modify anything, anywhere. Your real admin identity is never touched; write work (Terraform, etc.) is the deliberate `ws elevate` path, using your personal login. Full design in [`docs/security/README.md`](docs/security/README.md).

---

## Configuration

One file per project at `~/.config/ws/projects/<name>.yaml` — hand-editable. See [`docs/architecture/schemas/config.md`](docs/architecture/schemas/config.md) for the full schema, and the [`docs/`](docs/README.md) index for everything else (architecture, patterns, release process).

---

## Development

```bash
git clone https://github.com/richardamare/ws && cd ws
make check     # gofmt + go vet + go test -race
make build     # -> bin/ws
make snapshot  # local GoReleaser dry-run (no publish)
```

Contributions go through PRs (Issues are disabled). See [`CLAUDE.md`](CLAUDE.md) for the working rules and [`docs/release-engineering/`](docs/release-engineering/README.md) for cutting releases.

---

## License

[MIT](LICENSE) © ws authors
