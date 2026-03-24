#!/usr/bin/env bash
# Scenario: dkg_disabled_one_node
# Cases: IT-REG-04, IT-DL-03, IT-FN-02, IT-ACT-04
# Goal: Disable DKG on node N (set dkg.enable=false + restart), other nodes normal
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TARGET_NODE="${DKG_VALIDATOR_COUNT:-3}"
echo "[dkg_disabled_one_node] Pre: disabling DKG on validator $TARGET_NODE."

toggle_dkg "$TARGET_NODE" "false"

echo "[dkg_disabled_one_node] Pre: done. Validator $TARGET_NODE has dkg.enable=false."
