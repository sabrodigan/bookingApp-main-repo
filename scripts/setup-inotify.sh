#!/usr/bin/env bash
# Raises Linux inotify limits so GitKraken and IDEs can watch repositories.
set -euo pipefail

CONF_NAME="99-gitkraken-inotify.conf"
CONF_PATH="/etc/sysctl.d/${CONF_NAME}"

WATCHES="${INOTIFY_MAX_USER_WATCHES:-524288}"
INSTANCES="${INOTIFY_MAX_USER_INSTANCES:-512}"

echo "Current limits:"
sysctl fs.inotify.max_user_watches fs.inotify.max_user_instances 2>/dev/null || true
echo

CONTENT="fs.inotify.max_user_watches=${WATCHES}
fs.inotify.max_user_instances=${INSTANCES}
"

if [[ "$(id -u)" -eq 0 ]]; then
  echo "$CONTENT" > "$CONF_PATH"
  sysctl --system
  echo "Done. Restart GitKraken (and your IDE if needed)."
  exit 0
fi

if command -v sudo >/dev/null 2>&1; then
  echo "Installing ${CONF_PATH} (requires sudo)..."
  echo "$CONTENT" | sudo tee "$CONF_PATH" >/dev/null
  sudo sysctl --system
  echo "Done. Restart GitKraken (and your IDE if needed)."
  exit 0
fi

echo "Could not write ${CONF_PATH}. Run as root or with sudo:"
echo
echo "$CONTENT"
echo "  sudo tee ${CONF_PATH}"
echo "  sudo sysctl --system"
