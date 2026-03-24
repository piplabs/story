#!/bin/bash
# mock_kernel_bad_dealer_finalize/pre.sh
# 在 Validator 3 (index=2) 上部署 mock kernel，GenerateDeals 返回无效 VSS share。
# 其他 RPC (ProcessDeals, ProcessResponses, FinalizeDKG, PartialDecryptTDH2) 正常行为。
# 目的: 验证 bad dealer 在 justification 失败后是否被 invalidate。

set -euo pipefail

TARGET_NODE_INDEX=${TARGET_NODE_INDEX:-2}

SSH_TARGETS=${DKG_SSH_TARGETS:-}
SSH_KEY=${DKG_SSH_KEY:-}
MOCK_BIN=${MOCK_KERNEL_BINARY:-/tmp/mock-kernel-server}
MOCK_PORT=${MOCK_KERNEL_PORT:-50051}

if [ -z "$SSH_TARGETS" ] || [ -z "$SSH_KEY" ]; then
    echo "WARNING: DKG_SSH_TARGETS or DKG_SSH_KEY not set; skipping"
    exit 0
fi

IFS=',' read -ra TARGETS <<< "$SSH_TARGETS"
TARGET=${TARGETS[$TARGET_NODE_INDEX]}

echo "=== mock_kernel_bad_dealer_finalize: Deploying on $TARGET ==="
echo "    Mode: invalid-vss (GenerateDeals returns invalid VSS shares)"
echo "    Other RPCs: normal behavior (ProcessDeals, FinalizeDKG, PartialDecryptTDH2)"

ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl stop story-kernel 2>/dev/null || true"

ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "nohup $MOCK_BIN --listen=:$MOCK_PORT --mode=invalid-vss-only-deals > /tmp/mock-kernel.log 2>&1 &"

echo "=== mock_kernel_bad_dealer_finalize: Mock kernel started ==="
sleep 5
