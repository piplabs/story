#!/bin/bash
# mock_kernel_invalid_argument/pre.sh
# 部署 mock kernel，GenerateDeals 返回 gRPC codes.InvalidArgument。
# 验证 story 是否正确分类为 non-retryable error → session marked Failed immediately。

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

echo "=== mock_kernel_invalid_argument: Deploying on $TARGET ==="
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl stop story-kernel 2>/dev/null || true"
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "nohup $MOCK_BIN --listen=:$MOCK_PORT --mode=invalid-argument > /tmp/mock-kernel.log 2>&1 &"
echo "=== mock_kernel_invalid_argument: Mock kernel started ==="
sleep 5
