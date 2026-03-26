#!/usr/bin/env bash
# Scenario: upgrade_scheduled
# Cases: IT-BB-03/09, IT-REG-03/13/14, IT-DL-09, IT-ACT-05, IT-E2E-06, IT-UPG-01~05
#
# Upgrade flow: new kernel (release/0.1) → old kernel (de329d8)
# Pre: start old kernel on port 50052 as the "upgrade target", add to story endpoints.
# The current kernel on 50051 is the "old" (from DKG's perspective, the one being upgraded FROM).
# The old binary on 50052 has a different mrenclave, so story sees it as the "new" binary.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[upgrade_scheduled] Pre: setting up dual kernel for upgrade resharing."

# Validate prerequisites
if [ -z "${STORY_ETH_RPC_URL:-}" ]; then
  echo "ERROR: STORY_ETH_RPC_URL not set."
  exit 1
fi
if [ -z "${DKG_SIGNER_PRIVATE_KEY:-}" ]; then
  echo "ERROR: DKG_SIGNER_PRIVATE_KEY not set."
  exit 1
fi

OLD_KERNEL_DIR="/home/ubuntu/story-kernel-noop"
OLD_KERNEL_PORT="50052"
KERNEL_HOME="${DKG_KERNEL_HOME:-/opt/story-kernel}"
STORY_TOML="${DKG_STORY_TOML_PATH:-/home/ubuntu/.story/story/config/story.toml}"

# Start old kernel binary on port 50052 on all validators
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  HAS_OLD=$(_ssh_cmd "$i" "[ -f ${OLD_KERNEL_DIR}/story-kernel.sig ] && echo yes || echo no" 2>/dev/null | tr -d '[:space:]')
  if [ "$HAS_OLD" != "yes" ]; then
    echo "  [WARN] node $i: old kernel not found at ${OLD_KERNEL_DIR}, upgrade round may fail"
    continue
  fi

  _ssh_cmd "$i" "
    pkill -f 'gramine-sgx.*${OLD_KERNEL_PORT}' 2>/dev/null || true
    sleep 1
    cd ${OLD_KERNEL_DIR}
    nohup gramine-sgx story-kernel start --home ${KERNEL_HOME} --grpc.listen_addr :${OLD_KERNEL_PORT} > /tmp/story-kernel-old.log 2>&1 &
  " 2>/dev/null || true
  echo "  node $i: old kernel started on port ${OLD_KERNEL_PORT}"
done

sleep 10

# Add old kernel endpoint to story.toml
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  HAS_EP=$(_ssh_cmd "$i" "grep '${OLD_KERNEL_PORT}' ${STORY_TOML} 2>/dev/null | wc -l" 2>/dev/null | tr -d '[:space:]')
  if [ "${HAS_EP:-0}" -eq 0 ]; then
    _ssh_cmd "$i" "sed -i 's|\"127.0.0.1:50051\"|\"127.0.0.1:50051\", \"127.0.0.1:${OLD_KERNEL_PORT}\"|' ${STORY_TOML}" 2>/dev/null
    echo "  node $i: added ${OLD_KERNEL_PORT} to kernel-endpoints"
  fi
done

# Restart story to pick up dual kernel endpoints
echo "  Restarting story..."
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  _ssh_cmd "$i" "sudo systemctl restart story" 2>/dev/null &
done
wait
sleep 5

echo "[upgrade_scheduled] Pre: dual kernel running (50051=current, ${OLD_KERNEL_PORT}=upgrade target)."
