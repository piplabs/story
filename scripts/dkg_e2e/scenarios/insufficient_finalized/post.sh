#!/usr/bin/env bash
# Restore: restart all kernels after insufficient_finalized scenario
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[insufficient_finalized] Post: restarting all kernels."
start_all_kernels
echo "[insufficient_finalized] Post: done."
