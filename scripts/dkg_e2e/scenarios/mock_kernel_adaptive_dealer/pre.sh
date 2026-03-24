#!/bin/bash
# mock_kernel_adaptive_dealer/pre.sh
# 在 Validator 3 上部署特制 mock kernel：
#   GenerateDeals → 无效 VSS share
#   FinalizeDKG → 排除自己的 deal contribution，匹配诚实多数 globalPubKey
# 这是一个理论攻击场景，验证自适应 bad dealer 是否可行。

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

echo "=== mock_kernel_adaptive_dealer: Deploying on $TARGET ==="
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl stop story-kernel 2>/dev/null || true"
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "nohup $MOCK_BIN --listen=:$MOCK_PORT --mode=adaptive-bad-dealer > /tmp/mock-kernel.log 2>&1 &"

echo "=== mock_kernel_adaptive_dealer: Mock kernel started (adaptive mode) ==="
sleep 5
