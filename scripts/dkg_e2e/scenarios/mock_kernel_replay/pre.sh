#!/bin/bash
# mock_kernel_replay/pre.sh
# 在 validator 3 上部署 mock kernel，配置为重放旧 round 的 deals/justifications/finalization。
# 需要先运行一轮完整 DKG 以捕获 round N 数据，然后在 round N+1 重放。

set -euo pipefail

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

echo "=== mock_kernel_replay: Deploying to $TARGET ==="

ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl stop story-kernel 2>/dev/null || true"

# 启动 mock kernel（replay 模式：记录当前 round 数据，下一轮重放）
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "nohup $MOCK_BIN --listen=:$MOCK_PORT --mode=replay > /tmp/mock-kernel.log 2>&1 &"

echo "=== mock_kernel_replay: Mock kernel started in replay mode ==="
sleep 5
