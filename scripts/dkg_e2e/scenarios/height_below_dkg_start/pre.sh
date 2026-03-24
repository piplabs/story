#!/usr/bin/env bash
# Scenario: height_below_dkg_start
# Cases: IT-BB-01
# Goal: Chain height < dkgStartBlock so BeginBlocker returns nil
# NOTE: This typically requires a fresh chain or genesis with high dkgStartBlock.
# The script checks current height and logs; actual chain setup is environment-specific.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[height_below_dkg_start] Pre: verifying chain height < dkgStartBlock."
echo "  This scenario requires either:"
echo "  1. A fresh chain that hasn't reached dkgStartBlock yet, OR"
echo "  2. Genesis configured with a high dkgStartBlock value"
echo "  Ensure your devnet meets this condition before running."

# If STORY_RPC_URL is set, check current block height
if [ -n "${STORY_RPC_URL:-}" ]; then
  HEIGHT=$(curl -s "${STORY_RPC_URL}/status" | jq -r '.result.sync_info.latest_block_height' 2>/dev/null || echo "unknown")
  echo "[height_below_dkg_start] Current block height: $HEIGHT"
fi
