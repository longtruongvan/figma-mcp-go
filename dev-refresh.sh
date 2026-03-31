#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_GO=true
BUILD_PLUGIN=true
CHECK_MCP=true
RESTART_CODEX=false

usage() {
  cat <<'EOF'
Usage: ./dev-refresh.sh [options]

Options:
  --skip-go         Skip rebuilding the Go MCP server binary
  --skip-plugin     Skip rebuilding the Figma plugin
  --no-mcp-check    Skip printing Codex MCP status
  --restart-codex   Restart the Codex macOS app after rebuilding
  -h, --help        Show this help message

Default behavior:
  1. Build bin/figma-mcp-go
  2. Build plugin/code.js and plugin/index.html
  3. Show current Codex MCP config when Codex CLI is available
  4. Print the next manual reload steps for Codex and Figma
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --skip-go)
      BUILD_GO=false
      ;;
    --skip-plugin)
      BUILD_PLUGIN=false
      ;;
    --no-mcp-check)
      CHECK_MCP=false
      ;;
    --restart-codex)
      RESTART_CODEX=true
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      echo
      usage
      exit 1
      ;;
  esac
  shift
done

find_codex_bin() {
  if command -v codex >/dev/null 2>&1; then
    command -v codex
    return 0
  fi

  local app_bin="/Applications/Codex.app/Contents/Resources/codex"
  if [[ -x "$app_bin" ]]; then
    echo "$app_bin"
    return 0
  fi

  return 1
}

require_cmd() {
  local cmd="$1"
  local hint="$2"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "Missing required command: $cmd" >&2
    echo "$hint" >&2
    exit 1
  fi
}

echo "==> Working directory: $ROOT_DIR"

if $BUILD_GO; then
  require_cmd "go" "Install Go first, for example: brew install go"
  echo "==> Building Go MCP server binary"
  (
    cd "$ROOT_DIR"
    go build -o bin/figma-mcp-go ./cmd/figma-mcp-go
  )
fi

if $BUILD_PLUGIN; then
  require_cmd "npm" "Install Node.js first so npm is available"
  echo "==> Building Figma plugin"
  (
    cd "$ROOT_DIR/plugin"
    npm run build
  )
fi

CODEX_BIN=""
if $CHECK_MCP || $RESTART_CODEX; then
  if CODEX_BIN="$(find_codex_bin)"; then
    :
  else
    CODEX_BIN=""
  fi
fi

if $CHECK_MCP; then
  echo "==> Codex MCP status"
  if [[ -n "$CODEX_BIN" ]]; then
    "$CODEX_BIN" mcp list || true
    echo
    "$CODEX_BIN" mcp get figma-mcp-go || true
  else
    echo "Codex CLI not found in PATH and not found in /Applications/Codex.app/Contents/Resources/codex"
  fi
fi

if $RESTART_CODEX; then
  echo "==> Restarting Codex"
  if [[ -n "$CODEX_BIN" ]] && [[ "$OSTYPE" == darwin* ]]; then
    osascript -e 'tell application "Codex" to quit' >/dev/null 2>&1 || true
    sleep 1
    open -a Codex || true
  else
    echo "Skipping Codex restart because Codex app was not found or this is not macOS"
  fi
fi

cat <<'EOF'

Done.

Next steps:
1. If the Go server changed, open a new Codex thread or rerun this script with --restart-codex.
2. In Figma Desktop, close and reopen the "Figma MCP Go" plugin.
3. Confirm the badge shows "Connected".
4. Smoke test with: get_metadata

Tips:
- Plugin/source only changed: ./dev-refresh.sh --skip-go
- Go server only changed: ./dev-refresh.sh --skip-plugin
- Full refresh + app restart: ./dev-refresh.sh --restart-codex
EOF
