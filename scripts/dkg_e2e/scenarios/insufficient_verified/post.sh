#!/usr/bin/env bash
# Restore: restart all validators after insufficient_verified scenario
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[insufficient_verified] Post: restarting all kernels."
# Story was never stopped, only restart kernels
start_all_kernels
echo "[insufficient_verified] Post: done."
