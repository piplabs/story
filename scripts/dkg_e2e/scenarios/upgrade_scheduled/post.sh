#!/usr/bin/env bash
# Restore after upgrade: if upgrade succeeded, the chain now uses the old kernel's CC.
# Stop the new kernel (50051), move old kernel from 50052 to 50051, clean up endpoints.
# If upgrade failed, just stop old kernel and restore single-endpoint config.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

OLD_KERNEL_DIR="/home/ubuntu/story-kernel-noop"
OLD_KERNEL_PORT="50052"
KERNEL_HOME="${DKG_KERNEL_HOME:-/opt/story-kernel}"
STORY_TOML="${DKG_STORY_TOML_PATH:-/home/ubuntu/.story/story/config/story.toml}"

echo "[upgrade_scheduled] Post: cleaning up after upgrade test."

# Remove 50052 from story.toml (restore single endpoint)
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  _ssh_cmd "$i" "sed -i 's|, \"127.0.0.1:${OLD_KERNEL_PORT}\"||' ${STORY_TOML}" 2>/dev/null
done

# Cancel any pending upgrade
OWNER_KEY="${DKG_OWNER_KEY:-${DKG_SIGNER_PRIVATE_KEY:-ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}}"
DKG_ADDR="0xCcCcCC0000000000000000000000000000000004"
RPC="http://localhost:8545"
VERSION="${DKG_UPGRADE_VERSION:-v2.0.0-test}"
_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast send --private-key 0x${OWNER_KEY} --rpc-url ${RPC} --gas-price 20000000000 --legacy ${DKG_ADDR} 'cancelUpgrade(string)' '${VERSION}' 2>/dev/null" 2>/dev/null || true

# Stop both kernels, then start the old kernel as the primary on port 50051.
# After a successful upgrade, the chain's code commitment is the old kernel's mrenclave.
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  _ssh_cmd "$i" "
    sudo systemctl stop story-kernel 2>/dev/null || true
    pkill -f 'gramine-sgx.*${OLD_KERNEL_PORT}' 2>/dev/null || true
    sleep 2
    # Swap: copy old kernel files over the current kernel
    cp ${OLD_KERNEL_DIR}/story-kernel.manifest /home/ubuntu/story-kernel/story-kernel.manifest
    cp ${OLD_KERNEL_DIR}/story-kernel.manifest.sgx /home/ubuntu/story-kernel/story-kernel.manifest.sgx
    cp ${OLD_KERNEL_DIR}/story-kernel.sig /home/ubuntu/story-kernel/story-kernel.sig
    cp ${OLD_KERNEL_DIR}/build/story-kernel /home/ubuntu/story-kernel/build/story-kernel
    sudo systemctl start story-kernel
  " 2>/dev/null
  echo "  node $i: swapped to old kernel on port 50051"
done

# Restart story with single endpoint
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  _ssh_cmd "$i" "sudo systemctl restart story" 2>/dev/null &
done
wait
sleep 5

echo "[upgrade_scheduled] Post: old kernel now primary on 50051, upgrade cleaned up."
