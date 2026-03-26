#!/usr/bin/env bash
# Scenario: insufficient_finalized
# Cases: IT-E2E-04, IT-ACT-02, IT-ACT-03, IT-SKIP-02
# Goal: Registration+Dealing normal (3 nodes), then stop 2 kernels so only 1 can finalize
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[insufficient_finalized] Pre: stopping kernel on validators 2..${DKG_VALIDATOR_COUNT:-3} to prevent finalization."

# Stop story-kernel (TEE) on nodes 2..N
stop_kernels_except 1

echo "[insufficient_finalized] Pre: done. Only validator 1 kernel is running for finalization."
