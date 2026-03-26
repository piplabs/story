#!/bin/bash
set -uo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"
TARGET_NODE=${TARGET_NODE_INDEX:-3}
stop_mock_kernel "$TARGET_NODE" || true
restore_real_code_commitment
restore_real_dcap
echo "=== mock_kernel_replay: restored ==="
