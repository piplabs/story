#!/bin/bash
# kernel_restart_after_deals/pre.sh
# 在 Dealing 阶段，validator 3 的 ProcessDeals 已完成后，SSH 重启其 story-kernel。
# 目的：验证 rebuildInitDKG 能从磁盘 replay deals 恢复 DKG 状态。
#
# 注意：此脚本由 ScriptDriver 在测试进入 Dealing 阶段后自动调用。
# 脚本需要等待足够时间确保 ProcessDeals 已完成（通过检查日志或等待固定时间）。

set -euo pipefail

TARGET_NODE_INDEX=${TARGET_NODE_INDEX:-2}  # 0-indexed, validator 3

SSH_TARGETS=${DKG_SSH_TARGETS:-}
SSH_KEY=${DKG_SSH_KEY:-}

if [ -z "$SSH_TARGETS" ] || [ -z "$SSH_KEY" ]; then
    echo "WARNING: DKG_SSH_TARGETS or DKG_SSH_KEY not set; skipping"
    exit 0
fi

IFS=',' read -ra TARGETS <<< "$SSH_TARGETS"
TARGET=${TARGETS[$TARGET_NODE_INDEX]}

echo "=== kernel_restart_after_deals: Waiting for ProcessDeals to complete ==="
# 等待 validator 日志出现 ProcessDeals 完成的标志
# 或简单等待一段时间（Dealing 阶段初期 deals 会很快处理）
sleep 15

echo "=== kernel_restart_after_deals: Restarting story-kernel on $TARGET ==="
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl restart story-kernel"

echo "=== kernel_restart_after_deals: Kernel restarted, enclave will reload ==="
# 等待 enclave 加载完成（通常 20-60 秒）
sleep 30
echo "=== kernel_restart_after_deals: Ready for test verification ==="
