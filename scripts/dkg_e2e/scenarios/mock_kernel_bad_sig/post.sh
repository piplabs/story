#!/bin/bash
# mock_kernel_bad_sig/post.sh — Restore real kernel on validator 3.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TARGET_NODE=${TARGET_NODE_INDEX:-3}  # 1-based

stop_mock_kernel "$TARGET_NODE"
restore_real_code_commitment
restore_real_dcap

echo "=== mock_kernel_bad_sig: restored ==="
