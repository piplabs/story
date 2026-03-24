#!/bin/bash
# mock_kernel_2_bad_dealers/pre.sh
# 在 Validator 2 和 3 上同时部署 mock kernel (invalid VSS)，仅 Validator 1 用真实 TEE。

set -euo pipefail

SSH_TARGETS=${DKG_SSH_TARGETS:-}
SSH_KEY=${DKG_SSH_KEY:-}
MOCK_BIN=${MOCK_KERNEL_BINARY:-/tmp/mock-kernel-server}
MOCK_PORT=${MOCK_KERNEL_PORT:-50051}

if [ -z "$SSH_TARGETS" ] || [ -z "$SSH_KEY" ]; then
    echo "WARNING: DKG_SSH_TARGETS or DKG_SSH_KEY not set; skipping"
    exit 0
fi

IFS=',' read -ra TARGETS <<< "$SSH_TARGETS"

for IDX in 1 2; do
    TARGET=${TARGETS[$IDX]}
    echo "=== mock_kernel_2_bad_dealers: Deploying on node $((IDX+1)) ($TARGET) ==="
    ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
        "sudo systemctl stop story-kernel 2>/dev/null || true"
    ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
        "nohup $MOCK_BIN --listen=:$MOCK_PORT --mode=invalid-vss-only-deals > /tmp/mock-kernel.log 2>&1 &"
done

echo "=== mock_kernel_2_bad_dealers: 2 mock kernels started (nodes 2,3) ==="
sleep 5
