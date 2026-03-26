#!/usr/bin/env bash
# Scenario: insufficient_verified
# Cases: IT-E2E-03, IT-DL-02, IT-SKIP-01
# Goal: Only 1 validator has DKG+TEE running → Verified < MinReq → SkipToNextRound
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[insufficient_verified] Pre: stopping kernel on validators 2..${DKG_VALIDATOR_COUNT:-3}, keeping only node 1's kernel."
echo "  (story stays running on all nodes to maintain consensus)"

# Stop story-kernel (TEE) on nodes 2..N so they can't register
# Do NOT stop story — need 2/3 consensus to keep chain producing blocks
stop_kernels_except 1

echo "[insufficient_verified] Pre: done. Only validator 1 has kernel running."
