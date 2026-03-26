#!/usr/bin/env bash
# kernel_reconnect/post.sh — Restart kernel on validator 3.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TARGET_NODE="${DKG_VALIDATOR_COUNT:-3}"
echo "[kernel_reconnect] Post: restarting kernel on validator ${TARGET_NODE}."

start_kernel "$TARGET_NODE"

echo "[kernel_reconnect] Post: done."
