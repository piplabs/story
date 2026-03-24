#!/usr/bin/env bash
# Scenario: tee_down_one_node
# Cases: IT-REG-10/11, IT-DL-08/10, IT-FN-07/08, IT-ACT-12, IT-RES-*, IT-EDGE-06
# Goal: Stop story-kernel on one validator so TEE calls fail → Phase=Failed / MarkFailed
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TARGET_NODE="${DKG_VALIDATOR_COUNT:-3}"
echo "[tee_down_one_node] Pre: stopping story-kernel on validator $TARGET_NODE."

stop_kernel "$TARGET_NODE"

echo "[tee_down_one_node] Pre: done. Validator $TARGET_NODE kernel is down."
