#!/usr/bin/env bash
#
# install.sh — install the `ws` CLI from source via `go install`.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/richardamare/ws/master/scripts/install.sh | bash
#   ./scripts/install.sh            # install latest tagged version
#   WS_VERSION=v0.1.0 ./scripts/install.sh
#   WS_VERSION=master ./scripts/install.sh   # bleeding edge
#
# Requires: Go 1.23+ (https://go.dev/dl/ or `brew install go`).

set -euo pipefail

MODULE="github.com/richardamare/ws/cmd/ws"
VERSION="${WS_VERSION:-latest}"

red()  { printf '\033[31m%s\033[0m\n' "$*"; }
blue() { printf '\033[1;34m%s\033[0m\n' "$*"; }

if ! command -v go >/dev/null 2>&1; then
  red "Go is not installed. Install it first:"
  echo "  macOS:  brew install go"
  echo "  other:  https://go.dev/dl/"
  exit 1
fi

blue "Installing ${MODULE}@${VERSION} ..."
go install "${MODULE}@${VERSION}"

# Resolve where `go install` put the binary.
BINDIR="$(go env GOBIN)"
[ -n "$BINDIR" ] || BINDIR="$(go env GOPATH)/bin"
BIN="${BINDIR}/ws"

if [ ! -x "$BIN" ]; then
  red "Install finished but ${BIN} was not found."
  exit 1
fi

blue "Installed: ${BIN}"
"$BIN" version || true

case ":${PATH}:" in
  *:"${BINDIR}":*) ;;
  *)
    echo
    blue "Add ${BINDIR} to your PATH (e.g. in ~/.zshrc):"
    echo "  export PATH=\"${BINDIR}:\$PATH\""
    ;;
esac
