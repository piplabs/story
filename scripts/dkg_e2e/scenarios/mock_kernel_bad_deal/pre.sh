#!/bin/bash
# mock_kernel_bad_deal/pre.sh — Deploy mock kernel (bad-vss-deal) on validator 3.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"
TARGET_NODE=${TARGET_NODE_INDEX:-3}
deploy_fake_dcap
whitelist_mock_code_commitment
deploy_mock_kernel "$TARGET_NODE" "bad-vss-deal"
echo "=== mock_kernel_bad_deal: ready ==="
