#!/bin/bash
# mock_kernel_non_committee/pre.sh
# 在一个不在 DKG committee 中的 validator 上部署 mock kernel，令其尝试提交 partial decryption。
# 通常用于测试非 committee member 提交 partial 被拒绝的场景。

set -euo pipefail

# 默认选最后一个 validator（如果有额外节点不在 committee 中）
TARGET_NODE_INDEX=${TARGET_NODE_INDEX:-2}

SSH_TARGETS=${DKG_SSH_TARGETS:-}
SSH_KEY=${DKG_SSH_KEY:-}
MOCK_BIN=${MOCK_KERNEL_BINARY:-/tmp/mock-kernel-server}
MOCK_PORT=${MOCK_KERNEL_PORT:-50051}

if [ -z "$SSH_TARGETS" ] || [ -z "$SSH_KEY" ]; then
    echo "WARNING: DKG_SSH_TARGETS or DKG_SSH_KEY not set; skipping mock deployment"
    exit 0
fi

IFS=',' read -ra TARGETS <<< "$SSH_TARGETS"
TARGET=${TARGETS[$TARGET_NODE_INDEX]}

echo "=== mock_kernel_non_committee: Deploying to $TARGET ==="

ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl stop story-kernel 2>/dev/null || true"

ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "nohup $MOCK_BIN --listen=:$MOCK_PORT --mode=force-partial > /tmp/mock-kernel.log 2>&1 &"

echo "=== mock_kernel_non_committee: Mock kernel started in force-partial mode ==="
sleep 5
