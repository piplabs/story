#!/usr/bin/env bash
# Scenario: upgrade_scheduled
# Cases: IT-BB-03/09, IT-REG-03/13/14, IT-DL-09, IT-ACT-05, IT-E2E-06, IT-UPG-01~05
#
# Upgrade flow: new kernel (release/0.1) → old kernel (de329d8)
# Pre: start old kernel on port 50052 as the "upgrade target", add to story endpoints.
# The current kernel on 50051 is the "old" (from DKG's perspective, the one being upgraded FROM).
# The old binary on 50052 has a different mrenclave, so story sees it as the "new" binary.
#
# NOTE: Gramine ignores command-line args — loader.argv in the manifest is the real argv.
# The noop kernel manifest hardcodes --home /opt/story-kernel, so we must temporarily
# change the config's listen_addr to 50052 before starting, then restore it.
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
KERNEL_CONFIG="${KERNEL_HOME}/config.toml"

# Start noop kernel binary on port 50052 on all validators
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  HAS_OLD=$(_ssh_cmd "$i" "[ -f ${OLD_KERNEL_DIR}/story-kernel.sig ] && echo yes || echo no" 2>/dev/null | tr -d '[:space:]')
  if [ "$HAS_OLD" != "yes" ]; then
    echo "  [WARN] node $i: old kernel not found at ${OLD_KERNEL_DIR}, upgrade round may fail"
    continue
  fi

  # Gramine manifest hardcodes --home /opt/story-kernel, ignoring CLI args.
  # Temporarily change config listen_addr to 50052, start noop kernel, then restore.
  _ssh_cmd "$i" "
    pkill -f 'story-kernel-noop' 2>/dev/null || true
    sudo lsof -ti :${OLD_KERNEL_PORT} | xargs -r sudo kill -9 2>/dev/null || true
    sleep 1
    # Swap config port: 50051 → 50052
    sed -i 's|listen_addr = \":50051\"|listen_addr = \":${OLD_KERNEL_PORT}\"|' ${KERNEL_CONFIG}
    cd ${OLD_KERNEL_DIR}
    nohup gramine-sgx story-kernel > /tmp/story-kernel-old.log 2>&1 &
    NOOP_PID=\$!
    echo \"  noop kernel PID=\$NOOP_PID\"
    # Wait for noop kernel to read config and start binding
    sleep 5
    # Restore config port: 50052 → 50051 (main kernel already running, won't re-read)
    sed -i 's|listen_addr = \":${OLD_KERNEL_PORT}\"|listen_addr = \":50051\"|' ${KERNEL_CONFIG}
  " 2>/dev/null || true
  echo "  node $i: noop kernel starting on port ${OLD_KERNEL_PORT}"
done

# Health check: wait for noop kernel to be connectable on 50052
echo "  Waiting for noop kernel on port ${OLD_KERNEL_PORT}..."
MAX_WAIT=60
for i in $(seq 1 "${DKG_VALIDATOR_COUNT:-3}"); do
  ELAPSED=0
  while [ $ELAPSED -lt $MAX_WAIT ]; do
    LISTENING=$(_ssh_cmd "$i" "ss -tlnp | grep ':${OLD_KERNEL_PORT}' | wc -l" 2>/dev/null | tr -d '[:space:]')
    if [ "${LISTENING:-0}" -gt 0 ]; then
      echo "  node $i: noop kernel listening on ${OLD_KERNEL_PORT} (${ELAPSED}s)"
      break
    fi
    sleep 5
    ELAPSED=$((ELAPSED + 5))
  done
  if [ $ELAPSED -ge $MAX_WAIT ]; then
    echo "  [ERROR] node $i: noop kernel NOT listening on ${OLD_KERNEL_PORT} after ${MAX_WAIT}s"
    _ssh_cmd "$i" "tail -20 /tmp/story-kernel-old.log" 2>/dev/null || true
  fi
done

# Add noop kernel endpoint to story.toml
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
