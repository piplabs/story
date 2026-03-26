#!/usr/bin/env bash
# Restore after upgrade: if upgrade succeeded, the chain now uses the noop kernel's CC.
# Stop noop kernel on 50052, swap noop binary into /home/ubuntu/story-kernel/ (primary),
# restart story-kernel on 50051 with the noop binary so subsequent rounds use the new CC.
# If upgrade failed, just stop noop kernel and restore single-endpoint config.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

OLD_KERNEL_DIR="/home/ubuntu/story-kernel-noop"
OLD_KERNEL_PORT="50052"
KERNEL_HOME="${DKG_KERNEL_HOME:-/opt/story-kernel}"
STORY_TOML="${DKG_STORY_TOML_PATH:-/home/ubuntu/.story/story/config/story.toml}"

echo "[upgrade_scheduled] Post: cleaning up after upgrade test."

# 1. Remove 50052 from story.toml (restore single endpoint)
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  _ssh_cmd "$i" "sed -i 's|, \"127.0.0.1:${OLD_KERNEL_PORT}\"||' ${STORY_TOML}" 2>/dev/null
done

# 2. Cancel any pending upgrade (may fail if already consumed — that's OK)
OWNER_KEY="${DKG_OWNER_KEY:-${DKG_SIGNER_PRIVATE_KEY:-ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}}"
DKG_ADDR="0xCcCcCC0000000000000000000000000000000004"
RPC="http://localhost:8545"
VERSION="${DKG_UPGRADE_VERSION:-v2.0.0-test}"
_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast send --private-key 0x${OWNER_KEY} --rpc-url ${RPC} --gas-price 20000000000 --legacy ${DKG_ADDR} 'cancelUpgrade(string)' '${VERSION}' 2>/dev/null" 2>/dev/null || true

# 3. Stop noop kernel on 50052, stop main kernel, swap binary, restart
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  _ssh_cmd "$i" "
    # Stop noop kernel (kill by port since it's not a systemd service)
    sudo lsof -ti :${OLD_KERNEL_PORT} | xargs -r sudo kill -9 2>/dev/null || true
    pkill -f 'story-kernel-noop' 2>/dev/null || true

    # Stop main kernel
    sudo systemctl stop story-kernel 2>/dev/null || true
    sleep 2

    # Verify both kernels stopped
    for port in 50051 ${OLD_KERNEL_PORT}; do
      pid=\$(sudo lsof -ti :\$port 2>/dev/null || true)
      [ -n \"\$pid\" ] && sudo kill -9 \$pid 2>/dev/null || true
    done
    sleep 1

    # Swap: copy noop kernel files over the current kernel
    if [ -f ${OLD_KERNEL_DIR}/story-kernel.manifest ]; then
      cp -f ${OLD_KERNEL_DIR}/story-kernel.manifest /home/ubuntu/story-kernel/story-kernel.manifest
      cp -f ${OLD_KERNEL_DIR}/story-kernel.manifest.sgx /home/ubuntu/story-kernel/story-kernel.manifest.sgx
      cp -f ${OLD_KERNEL_DIR}/story-kernel.sig /home/ubuntu/story-kernel/story-kernel.sig
      cp -f ${OLD_KERNEL_DIR}/build/story-kernel /home/ubuntu/story-kernel/build/story-kernel
      echo '  binary swap done'
    else
      echo '  [WARN] noop kernel files not found, skipping swap'
    fi

    # Ensure config listen_addr is 50051 for primary kernel
    sed -i 's|listen_addr = \":${OLD_KERNEL_PORT}\"|listen_addr = \":50051\"|' ${KERNEL_HOME}/config.toml 2>/dev/null || true

    # Start kernel (now running noop binary on 50051)
    sudo systemctl start story-kernel
  " 2>/dev/null
  echo "  node $i: swapped to noop kernel on port 50051"
done

# Wait for kernels to start
sleep 10

# Verify kernels are running
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  STATUS=$(_ssh_cmd "$i" "systemctl is-active story-kernel" 2>/dev/null | tr -d '[:space:]')
  if [ "$STATUS" = "active" ]; then
    echo "  node $i: kernel active"
  else
    echo "  [WARN] node $i: kernel status=$STATUS, retrying..."
    _ssh_cmd "$i" "sudo systemctl start story-kernel" 2>/dev/null || true
  fi
done

# 4. Restart story with single endpoint
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  _ssh_cmd "$i" "sudo systemctl restart story" 2>/dev/null &
done
wait
sleep 5

echo "[upgrade_scheduled] Post: noop kernel now primary on 50051, upgrade cleaned up."
