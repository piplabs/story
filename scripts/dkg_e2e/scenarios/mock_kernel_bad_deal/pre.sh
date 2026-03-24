#!/bin/bash
# mock_kernel_bad_deal/pre.sh
# 在 validator 3 上停止真实 story-kernel，启动 mock kernel（无效 VSS deal / oversized / garbage）。

set -euo pipefail

TARGET_NODE_INDEX=${TARGET_NODE_INDEX:-2}
MOCK_MODE=${MOCK_MODE:-invalid-vss}  # invalid-vss | oversized | garbage

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

echo "=== mock_kernel_bad_deal: Deploying to $TARGET (mode=$MOCK_MODE) ==="

ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl stop story-kernel 2>/dev/null || true"

ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "nohup $MOCK_BIN --listen=:$MOCK_PORT --mode=$MOCK_MODE > /tmp/mock-kernel.log 2>&1 &"

echo "=== mock_kernel_bad_deal: Mock kernel started (mode=$MOCK_MODE) ==="
sleep 5
