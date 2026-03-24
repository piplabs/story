#!/bin/bash
# mock_kernel_bad_sig/post.sh
# 恢复 validator 3 的真实 story-kernel。

set -euo pipefail

TARGET_NODE_INDEX=${TARGET_NODE_INDEX:-2}

SSH_TARGETS=${DKG_SSH_TARGETS:-}
SSH_KEY=${DKG_SSH_KEY:-}
MOCK_PORT=${MOCK_KERNEL_PORT:-50051}

if [ -z "$SSH_TARGETS" ] || [ -z "$SSH_KEY" ]; then
    echo "WARNING: DKG_SSH_TARGETS or DKG_SSH_KEY not set; skipping teardown"
    exit 0
fi

IFS=',' read -ra TARGETS <<< "$SSH_TARGETS"
TARGET=${TARGETS[$TARGET_NODE_INDEX]}

echo "=== mock_kernel_bad_sig: Restoring $TARGET ==="

# 停止 mock kernel
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "pkill -f mock-kernel-server 2>/dev/null || true; sleep 2"

# 重新启动真实 story-kernel
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl start story-kernel"

echo "=== mock_kernel_bad_sig: Real story-kernel restored on $TARGET ==="
sleep 10  # 等待 kernel 重连
