#!/usr/bin/env bash
# configure_periods.sh — Change DKG period parameters via genesis JSON + chain reset.
#
# Since #726, DKG params come from genesis InitGenesis, not the upgrade handler.
# This script patches the bootnode genesis and resets the chain.
# Binary rebuild is only done when code has actually changed (auto-detected).
#
# Usage: ./configure_periods.sh <registration_period> <dealing_period> <finalization_period> <active_period>
#   All values are in BLOCKS (not seconds).
#
# Example:
#   ./configure_periods.sh 10 10 10 20       # Short cycle for testing
#   ./configure_periods.sh 200 300 300 200   # Devnet default

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

set -euo pipefail

if [ "$#" -ne 4 ]; then
  echo "Usage: $0 <registration_period> <dealing_period> <finalization_period> <active_period>"
  exit 1
fi

REG_PERIOD="$1"
DEAL_PERIOD="$2"
FIN_PERIOD="$3"
ACTIVE_PERIOD="$4"

TOTAL="${DKG_VALIDATOR_COUNT:-3}"
BOOTNODE_IP="${DKG_BOOTNODE_IP:-23.102.71.16}"
BOOTNODE_USER="${DKG_BOOTNODE_USER:-ubuntu}"
STORY_SRC="${DKG_STORY_SRC:-/home/ubuntu/story}"
UPGRADE_FILE="${STORY_SRC}/client/app/upgrades/v_2_0_0/upgrades.go"

echo "=== configure_periods: reg=${REG_PERIOD} deal=${DEAL_PERIOD} fin=${FIN_PERIOD} active=${ACTIVE_PERIOD} (blocks) ==="

# ---------------------------------------------------------------------------
# Detect whether binary rebuild is needed (code changed on remote)
# ---------------------------------------------------------------------------
NEED_REBUILD=false
STORY_BRANCH="${DKG_STORY_BRANCH:-dkg/dev}"
STORY_PIN="${DKG_STORY_PIN_COMMIT:-}"
DEVNET_COMMIT="${DKG_DEVNET_COMMIT:-c61b74a4}"

if [ "${DKG_FORCE_REBUILD:-}" = "true" ]; then
  NEED_REBUILD=true
  echo "[configure_periods] DKG_FORCE_REBUILD=true, will rebuild"
elif [ -n "$STORY_PIN" ]; then
  echo "[configure_periods] story pinned to ${STORY_PIN}, skipping pull check"
else
  echo "[configure_periods] Checking for code changes on ${STORY_BRANCH}..."
  STORY_CHANGED=$(_ssh_cmd 1 "cd ${STORY_SRC} && git fetch origin ${STORY_BRANCH} 2>/dev/null && LOCAL=\$(git rev-parse HEAD) && REMOTE=\$(git rev-parse origin/${STORY_BRANCH}) && if [ \"\$LOCAL\" != \"\$REMOTE\" ]; then echo changed; else echo up-to-date; fi" 2>/dev/null || echo "up-to-date")

  if [ "$STORY_CHANGED" = "changed" ]; then
    NEED_REBUILD=true
    echo "  story has new commits, will pull + rebuild"
  else
    echo "  story up-to-date, skipping rebuild"
  fi
fi

# ---------------------------------------------------------------------------
# Only rebuild if code actually changed
# ---------------------------------------------------------------------------
if [ "$NEED_REBUILD" = "true" ]; then
  # --- Pull latest code ---
  echo "[configure_periods] Pulling latest story (${STORY_BRANCH})..."
  PULL_STORY_CMD="cd ${STORY_SRC} && git fetch origin ${STORY_BRANCH} && git reset --hard origin/${STORY_BRANCH} && git cherry-pick --no-commit ${DEVNET_COMMIT} 2>/dev/null || true"

  for i in $(seq 1 "$TOTAL"); do
    RESULT=$(_ssh_cmd "$i" "$PULL_STORY_CMD" 2>&1 | tail -1)
    echo "  node $i: ${RESULT}"
  done
  ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$PULL_STORY_CMD" 2>&1 | tail -1 || true
  echo "  bootnode: updated"

  # --- Patch upgrade handler (backward compat with pre-#726 code) ---
  echo "[configure_periods] Patching upgrade handler on all machines..."
  PATCH_CMD="cd ${STORY_SRC} && \
    sed -i 's/^\(\t\+\)[0-9]\+,\(\s*\/\/ registration:.*\)/\1${REG_PERIOD},\2/' ${UPGRADE_FILE} && \
    sed -i 's/^\(\t\+\)[0-9]\+,\(\s*\/\/ dealing:.*\)/\1${DEAL_PERIOD},\2/' ${UPGRADE_FILE} && \
    sed -i 's/^\(\t\+\)[0-9]\+,\(\s*\/\/ finalization:.*\)/\1${FIN_PERIOD},\2/' ${UPGRADE_FILE} && \
    sed -i 's/^\(\t\+\)[0-9]\+,\(\s*\/\/ active:.*\)/\1${ACTIVE_PERIOD},\2/' ${UPGRADE_FILE}"

  for i in $(seq 1 "$TOTAL"); do
    _ssh_cmd "$i" "$PATCH_CMD"
  done
  ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$PATCH_CMD" || true

  # --- Patch V200 height ---
  V200_HEIGHT="${DKG_V200_HEIGHT:-150}"
  PATCH_V200_CMD="cd ${STORY_SRC} && sed -i '/LocalChainID: {/,/}/{s/V200:.*[0-9]\+/V200:    ${V200_HEIGHT}/}' lib/netconf/upgrades.go"
  for i in $(seq 1 "$TOTAL"); do
    _ssh_cmd "$i" "$PATCH_V200_CMD"
  done
  ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$PATCH_V200_CMD" 2>/dev/null || true

  # --- Rebuild story binary in parallel ---
  echo "[configure_periods] Rebuilding story binary on all machines..."
  echo "  (this takes ~1-2 minutes per machine)"
  STORY_BIN="${DKG_STORY_BIN:-/usr/local/bin/story}"
  BUILD_CMD="cd ${STORY_SRC} && make build 2>&1 | tail -3"

  for i in $(seq 1 "$TOTAL"); do
    (
      echo "  node $i: building..."
      _ssh_cmd "$i" "$BUILD_CMD"
      echo "  node $i: build done"
    ) &
  done
  (
    echo "  bootnode: building..."
    ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$BUILD_CMD" || true
    echo "  bootnode: build done"
  ) &
  wait
  echo "  All builds complete."

  # --- Stop services and install new binary ---
  echo "[configure_periods] Stopping services and installing new binary..."
  STOP_CMD="sudo systemctl stop story story-kernel 2>/dev/null; sudo systemctl stop ${DKG_SYSTEMD_GETH_SERVICE:-node-geth} 2>/dev/null; true"
  INSTALL_CMD="sudo cp ${STORY_SRC}/build/story ${STORY_BIN} && echo installed"

  for i in $(seq 1 "$TOTAL"); do
    _ssh_cmd "$i" "$STOP_CMD && $INSTALL_CMD"
    echo "  node $i: binary installed"
  done
  ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$STOP_CMD && $INSTALL_CMD" || true
  echo "  bootnode: binary installed"
else
  echo "[configure_periods] No code changes, skipping rebuild (saving ~3 min)"
fi

# ---------------------------------------------------------------------------
# Always: patch DKG params in bootnode genesis (used by reset_devnet.sh)
# ---------------------------------------------------------------------------
echo "[configure_periods] Patching DKG params in bootnode genesis..."
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "python3 - <<'PYEOF'
import json

genesis_path = '${DKG_STORY_HOME:-/home/ubuntu/.story/story}/config/genesis.json'
with open(genesis_path) as f:
    genesis = json.load(f)

dkg_params = {
    'registration_period': ${REG_PERIOD},
    'dealing_period': ${DEAL_PERIOD},
    'finalization_period': ${FIN_PERIOD},
    'active_period': ${ACTIVE_PERIOD},
    'dkg_committee_reward_portion': '0.100000000000000000',
    'min_req_registered_participants': 3,
    'min_req_finalized_participants': 3,
    'operational_threshold': 500
}
if 'dkg' not in genesis['app_state']:
    genesis['app_state']['dkg'] = {}
genesis['app_state']['dkg']['params'] = dkg_params

with open(genesis_path, 'w') as f:
    json.dump(genesis, f, indent=2)

print(f'  DKG params: reg={dkg_params[\"registration_period\"]} deal={dkg_params[\"dealing_period\"]} fin={dkg_params[\"finalization_period\"]} active={dkg_params[\"active_period\"]}')
PYEOF
" || { echo "[ERROR] Failed to patch bootnode genesis DKG params"; exit 1; }

# ---------------------------------------------------------------------------
# Reset devnet (clean data + restart chain)
# ---------------------------------------------------------------------------
echo "[configure_periods] Resetting devnet..."
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
bash "${SCRIPT_DIR}/../reset_devnet.sh"

echo "=== configure_periods: done ==="
echo "  registration=${REG_PERIOD} dealing=${DEAL_PERIOD} finalization=${FIN_PERIOD} active=${ACTIVE_PERIOD}"
