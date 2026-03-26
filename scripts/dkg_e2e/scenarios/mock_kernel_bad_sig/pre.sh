#!/bin/bash
# mock_kernel_bad_sig/pre.sh — Deploy mock kernel with forged-sig mode on validator 3.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TARGET_NODE=${TARGET_NODE_INDEX:-3}  # 1-based, default: validator 3

deploy_fake_dcap
whitelist_mock_code_commitment
deploy_mock_kernel "$TARGET_NODE" "forged-sig"

echo "=== mock_kernel_bad_sig: ready ==="
