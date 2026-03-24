#!/usr/bin/env bash
# Restore: restart kernel on the downed node
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TARGET_NODE="${DKG_VALIDATOR_COUNT:-3}"
echo "[tee_down_one_node] Post: restarting story-kernel on validator $TARGET_NODE."

start_kernel "$TARGET_NODE"

echo "[tee_down_one_node] Post: done."
