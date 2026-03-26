#!/usr/bin/env bash
# Restore: return TEE to normal mode
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[complaint_justification] Post: restoring normal TEE operation."
# If you modified kernel config for mock mode, restore it
start_all_kernels
echo "[complaint_justification] Post: done."
