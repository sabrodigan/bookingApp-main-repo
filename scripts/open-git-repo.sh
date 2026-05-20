#!/usr/bin/env bash
# Prints or opens the real git repository root (for GitKraken, etc.)
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Git repository root:"
echo "  ${REPO_ROOT}"
echo

if command -v gitkraken >/dev/null 2>&1; then
  echo "Opening in GitKraken..."
  exec gitkraken -p "${REPO_ROOT}"
fi

if command -v gtk-launch >/dev/null 2>&1 && grep -q 'gitkraken' ~/.local/share/applications/*.desktop 2>/dev/null; then
  echo "Try: gitkraken -p \"${REPO_ROOT}\""
fi

echo "In GitKraken: File → Open Repo → select the path above."
