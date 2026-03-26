#!/bin/bash
set -uo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"
stop_mock_kernel 2 || true
stop_mock_kernel 3 || true
restore_real_code_commitment
restore_real_dcap
echo "=== mock_kernel_2_bad_dealers: restored ==="
