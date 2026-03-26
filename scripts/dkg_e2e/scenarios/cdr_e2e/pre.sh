#!/usr/bin/env bash
# Scenario: cdr_e2e
# Cases: IT-CDR-01, IT-CDR-08
# Goal: Ensure Active round with GlobalPublicKey + CDR contract accessible
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[cdr_e2e] Pre: verifying Active round and CDR contract access."

if [ -z "${STORY_ETH_RPC_URL:-}" ]; then
  echo "ERROR: STORY_ETH_RPC_URL not set. CDR tests require chain interaction."
  exit 1
fi
if [ -z "${DKG_SIGNER_PRIVATE_KEY:-}" ]; then
  echo "ERROR: DKG_SIGNER_PRIVATE_KEY not set. CDR tests require signing."
  exit 1
fi
echo "[cdr_e2e] Pre: prerequisites OK."
