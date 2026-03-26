#!/usr/bin/env bash
# kernel_restart_active_decrypt/post.sh — Restart kernel on validator 3 after Active-stage stop.
# Tests PR #727: decrypt worker should resume after kernel restart.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TARGET_NODE="${DKG_VALIDATOR_COUNT:-3}"
echo "[kernel_restart_active_decrypt] Post: restarting kernel on validator ${TARGET_NODE}."

start_kernel "$TARGET_NODE"
sleep 60

echo "[kernel_restart_active_decrypt] Post: kernel restarted. Waiting for decrypt worker to resume."
