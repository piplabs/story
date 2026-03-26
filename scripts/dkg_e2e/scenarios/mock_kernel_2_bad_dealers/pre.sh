#!/bin/bash
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"
deploy_fake_dcap
whitelist_mock_code_commitment
deploy_mock_kernel 2 "bad-vss-deal"
deploy_mock_kernel 3 "bad-vss-deal"
echo "=== mock_kernel_2_bad_dealers: ready ==="
