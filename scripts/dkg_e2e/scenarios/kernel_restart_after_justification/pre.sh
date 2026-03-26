#!/bin/bash
# kernel_restart_after_justification/pre.sh
# 综合场景：
#   1. 在 validator 3 上部署 mock kernel 产出 bad deal → 触发 complaint
#   2. 等待 justification 处理完成
#   3. 停止 mock kernel，恢复真实 story-kernel（模拟重启）
# rebuildInitDKG 需要 replay deals + responses + justifications。

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

echo "=== Step 1: Deploy mock kernel (bad deal) on $TARGET ==="
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl stop story-kernel 2>/dev/null || true"
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "nohup $MOCK_BIN --listen=:$MOCK_PORT --mode=invalid-vss > /tmp/mock-kernel.log 2>&1 &"
sleep 5

echo "=== Step 2: Waiting for justification processing ==="
# 等待 DKG 进入 Dealing → complaint → justification 流程
sleep 60

echo "=== Step 3: Stop mock kernel, restart real story-kernel on $TARGET ==="
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "pkill -9 -f mock-kernel-server 2>/dev/null || true; sleep 3; \
     sudo lsof -ti :$MOCK_PORT | xargs -r sudo kill -9 2>/dev/null || true; sleep 2"
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl start story-kernel" || true

echo "=== kernel_restart_after_justification: Real kernel restarted (or restart attempted) ==="
sleep 30
echo "=== Ready for test verification ==="
