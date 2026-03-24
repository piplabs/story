#!/usr/bin/env bash
# Scenario: upgrade_scheduled
# Cases: IT-BB-03/09, IT-REG-03/13/14, IT-DL-09, IT-ACT-05, IT-E2E-06, IT-UPG-01~05
# Goal: Call ScheduleUpgrade(activationHeight, version) on DKG contract
# Requires: STORY_ETH_RPC_URL + DKG_SIGNER_PRIVATE_KEY configured
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[upgrade_scheduled] Pre: scheduling DKG upgrade via EthChainClient."
echo "  The Go test will call h.ChainClient.ScheduleDKGUpgrade() with:"
echo "    version=${DKG_UPGRADE_VERSION:-v2.0.0-test}"
echo "    activationOffset=${DKG_UPGRADE_ACTIVATION_OFFSET:-50} blocks ahead"
echo "  Ensure STORY_ETH_RPC_URL and DKG_SIGNER_PRIVATE_KEY are set in config.env"

# The actual ScheduleUpgrade call is done in Go code (scenarios_run.go)
# This script just validates prerequisites
if [ -z "${STORY_ETH_RPC_URL:-}" ]; then
  echo "ERROR: STORY_ETH_RPC_URL not set. Cannot schedule upgrade."
  exit 1
fi
if [ -z "${DKG_SIGNER_PRIVATE_KEY:-}" ]; then
  echo "ERROR: DKG_SIGNER_PRIVATE_KEY not set. Cannot schedule upgrade."
  exit 1
fi
echo "[upgrade_scheduled] Pre: prerequisites OK."
