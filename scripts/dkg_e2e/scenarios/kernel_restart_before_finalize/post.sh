#!/bin/bash
# kernel_restart_after_deals/post.sh
# 确保 validator 3 的 story-kernel 已恢复运行。

set -euo pipefail

TARGET_NODE_INDEX=${TARGET_NODE_INDEX:-2}
SSH_TARGETS=${DKG_SSH_TARGETS:-}
SSH_KEY=${DKG_SSH_KEY:-}

if [ -z "$SSH_TARGETS" ] || [ -z "$SSH_KEY" ]; then
    exit 0
fi

IFS=',' read -ra TARGETS <<< "$SSH_TARGETS"
TARGET=${TARGETS[$TARGET_NODE_INDEX]}

# 确认 kernel 正在运行
STATUS=$(ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
    "systemctl is-active story-kernel 2>/dev/null || echo inactive")

if [ "$STATUS" != "active" ]; then
    echo "=== kernel_restart_after_deals: Kernel not active ($STATUS), restarting ==="
    ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no "$TARGET" \
        "sudo systemctl restart story-kernel"
    sleep 30
fi

echo "=== kernel_restart_after_deals: Teardown complete ==="
