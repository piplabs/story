#!/usr/bin/env bash
# Restore: re-enable DKG on the disabled node
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TARGET_NODE="${DKG_VALIDATOR_COUNT:-3}"
echo "[dkg_disabled_one_node] Post: re-enabling DKG on validator $TARGET_NODE."

toggle_dkg "$TARGET_NODE" "true"

echo "[dkg_disabled_one_node] Post: done."
