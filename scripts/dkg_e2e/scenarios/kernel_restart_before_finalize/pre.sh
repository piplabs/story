#!/bin/bash
# kernel_restart_before_finalize/pre.sh
# 在 Finalization 阶段开始前（ProcessResponses 已完成），SSH 重启 validator 3 的 story-kernel。
# 这是最关键的重启场景：rebuildInitDKG 需要正确 replay deals + responses。
# 如果 PrivatePoly 未持久化，重启后 NewDistKeyGenerator 会生成新的随机多项式，
# 导致 DistKeyShare() 在不同节点产出不一致的 global_pub_key。

set -euo pipefail

TARGET_NODE_INDEX=${TARGET_NODE_INDEX:-2}

SSH_TARGETS=${DKG_SSH_TARGETS:-}
SSH_KEY=${DKG_SSH_KEY:-}

if [ -z "$SSH_TARGETS" ] || [ -z "$SSH_KEY" ]; then
    echo "WARNING: DKG_SSH_TARGETS or DKG_SSH_KEY not set; skipping"
    exit 0
fi

IFS=',' read -ra TARGETS <<< "$SSH_TARGETS"
TARGET=${TARGETS[$TARGET_NODE_INDEX]}

echo "=== kernel_restart_before_finalize: Waiting for Finalization stage ==="
# 等待足够时间确保 ProcessResponses 已完成（stage 转换到 Finalization）
sleep 10

echo "=== kernel_restart_before_finalize: Restarting story-kernel on $TARGET ==="
ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "sudo systemctl restart story-kernel"

echo "=== kernel_restart_before_finalize: Kernel restarted (PrivatePoly persistence test) ==="
sleep 30
echo "=== kernel_restart_before_finalize: Ready for FinalizeDKG verification ==="
