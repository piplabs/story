#!/usr/bin/env bash
# Scenario: complaint_justification
# Cases: IT-E2E-05
# Goal: Inject invalid deal data to trigger complaint → response → justification path
# NOTE: This requires TEE mock mode or data tampering capability.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[complaint_justification] Pre: triggering invalid deal / complaint path."
echo "  This scenario requires one of:"
echo "  1. story-kernel in test/mock mode that produces invalid deals"
echo "  2. Data tampering on one validator's deal submission"
echo "  3. Custom TEE mock binary that returns corrupted deal shares"
echo "  Ensure your environment supports one of these approaches."
