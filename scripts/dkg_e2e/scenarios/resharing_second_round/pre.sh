#!/usr/bin/env bash
# Scenario: resharing_second_round
# Cases: IT-REG-02, IT-BB-07, IT-E2E-02
# Goal: First round reaches Active, then wait for active period to end → new round (IsResharing=true)
# NOTE: All 3 validators stay running; we just need to wait for the active period.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[resharing_second_round] Pre: ensuring first round is Active, then waiting for active period end."
echo "  The Go test will use WaitForRound(round+1) to detect the new round."
echo "  No node changes needed — all validators stay running."
