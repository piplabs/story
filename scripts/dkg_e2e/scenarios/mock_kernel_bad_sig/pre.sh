#!/bin/bash
# mock_kernel_bad_sig/pre.sh
# 在 validator 3 (index=2) 上停止真实 story-kernel，启动 mock kernel（伪造签名）。
#
# 使用方式：
#   1. 将编译好的 mock-kernel-server binary 复制到目标节点
#   2. 停止 story-kernel systemd service
#   3. 启动 mock kernel 并配置 PartialDecryptTDH2 返回伪造签名
#
# 环境变量：
#   DKG_SSH_TARGETS: user@ip1,user@ip2,user@ip3
#   DKG_SSH_KEY: SSH 私钥路径
#   MOCK_KERNEL_BINARY: mock kernel binary 在目标节点的路径（默认 /tmp/mock-kernel-server）
#   MOCK_KERNEL_PORT: mock kernel 监听端口（默认 50051）

set -euo pipefail

TARGET_NODE_INDEX=${TARGET_NODE_INDEX:-2}  # 0-indexed, default: validator 3

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

echo "=== mock_kernel_bad_sig: Deploying to $TARGET ==="

# 停止真实 story-kernel
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl stop story-kernel 2>/dev/null || true"

# 启动 mock kernel（后台运行，配置为返回伪造签名）
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "nohup $MOCK_BIN --listen=:$MOCK_PORT --mode=forged-sig > /tmp/mock-kernel.log 2>&1 &"

echo "=== mock_kernel_bad_sig: Mock kernel started on $TARGET:$MOCK_PORT ==="
sleep 5  # 等待 mock kernel 启动
