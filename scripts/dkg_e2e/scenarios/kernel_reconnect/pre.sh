#!/usr/bin/env bash
# kernel_reconnect/pre.sh — Stop kernel on validator 3, keep story running.
# Tests PR #725 auto-reconnect: story should reconnect to kernel without restart.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TARGET_NODE="${DKG_VALIDATOR_COUNT:-3}"
echo "[kernel_reconnect] Pre: stopping kernel on validator ${TARGET_NODE} (story stays running)."

stop_kernel "$TARGET_NODE"
sleep 5

echo "[kernel_reconnect] Pre: kernel stopped. Story should report 'kernel unavailable'."
