#!/usr/bin/env bash
# kernel_restart_active_decrypt/pre.sh — Stop kernel on validator 3 during Active stage.
# Tests PR #727: decrypt worker should resume after kernel restart.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TARGET_NODE="${DKG_VALIDATOR_COUNT:-3}"
echo "[kernel_restart_active_decrypt] Pre: stopping kernel on validator ${TARGET_NODE} during Active stage."

stop_kernel "$TARGET_NODE"
sleep 5

echo "[kernel_restart_active_decrypt] Pre: kernel stopped. Decrypt worker interrupted."
