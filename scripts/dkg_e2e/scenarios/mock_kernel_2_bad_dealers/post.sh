#!/bin/bash
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"
stop_mock_kernel 2
stop_mock_kernel 3
restore_real_code_commitment
restore_real_dcap
echo "=== mock_kernel_2_bad_dealers: restored ==="
