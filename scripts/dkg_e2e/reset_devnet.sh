#!/usr/bin/env bash
# reset_devnet.sh — Full devnet reset following the operational runbook.
#
# This script resets all chain state and restores the devnet to a fresh state.
# It follows the runbook Part 1 steps: Stop → Clean → Init geth → Start → Config kernel → Start kernel.
#
# Prerequisites:
#   - SSH access configured (see config.env for DKG_SSH_TARGETS, DKG_SSH_KEY)
#   - Original genesis files exist on bootnode or are provided
#   - source tests/integration/dkg/config.env before running
#
# Usage:
#   source tests/integration/dkg/config.env
#   bash scripts/dkg_e2e/reset_devnet.sh
#
# Optional: provide CL genesis path as argument (otherwise copies from bootnode)
#   bash scripts/dkg_e2e/reset_devnet.sh /path/to/original/genesis.json

set -euo pipefail

# --- Load common helpers ---
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/scenarios/_common.sh"

# --- Configuration ---
TOTAL="${DKG_VALIDATOR_COUNT:-3}"
STORY_SVC="${DKG_SYSTEMD_STORY_SERVICE:-story}"
GETH_SVC="${DKG_SYSTEMD_GETH_SERVICE:-node-geth}"
KERNEL_SVC="${DKG_SYSTEMD_KERNEL_SERVICE:-story-kernel}"
STORY_HOME="${DKG_STORY_HOME:-/home/ubuntu/.story/story}"
STORY_SRC="${DKG_STORY_SRC:-/home/ubuntu/story}"
GETH_DATA_DIR="${DKG_GETH_DATA_DIR:-/home/ubuntu/.story/geth/data}"
GETH_GENESIS="${DKG_GETH_GENESIS:-/home/ubuntu/config/genesis-geth.json}"
KERNEL_HOME="${DKG_KERNEL_HOME:-/home/ubuntu/.story-kernel}"
BOOTNODE_IP="${DKG_BOOTNODE_IP:-23.102.71.16}"
BOOTNODE_USER="${DKG_BOOTNODE_USER:-ubuntu}"

# Validator IPs (parsed from DKG_SSH_TARGETS)
IFS=',' read -ra SSH_TARGETS <<< "${DKG_SSH_TARGETS}"

echo "=========================================="
echo "  DKG Devnet Full Reset"
echo "=========================================="
echo "Validators: ${TOTAL}"
echo "Bootnode: ${BOOTNODE_IP}"
echo ""

# --- Step 1: Stop all processes on all machines (runbook Step 1) ---
echo "=== Step 1: Stopping all processes ==="

# Stop all validators in parallel
STOP_VALIDATOR_CMD="
  sudo systemctl stop ${KERNEL_SVC} 2>/dev/null || true
  sudo systemctl stop ${STORY_SVC} 2>/dev/null || true
  sudo systemctl stop ${GETH_SVC} 2>/dev/null || true
  sleep 1
  sudo pkill -9 -x gramine-sgx 2>/dev/null || true
  sudo pkill -9 -x loader 2>/dev/null || true
  for port in 50051 50052; do
    pid=\$(sudo lsof -t -i :\$port 2>/dev/null || true)
    [ -n \"\$pid\" ] && sudo kill -9 \$pid 2>/dev/null || true
  done
  sudo pkill -9 -x geth 2>/dev/null || true
  sudo pkill -9 -x story 2>/dev/null || true
"

for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "$STOP_VALIDATOR_CMD" &
done

# Stop bootnode in parallel
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "
  sudo systemctl stop ${STORY_SVC} 2>/dev/null || true
  sudo systemctl stop ${GETH_SVC} 2>/dev/null || true
  sleep 1
  sudo pkill -9 -x geth 2>/dev/null || true
  sudo pkill -9 -x story 2>/dev/null || true
" &

wait
echo "  All processes stopped."
sleep 1

# --- Step 2: Restore original genesis (from bootnode or argument) ---
echo "=== Step 2: Restoring original genesis ==="

CL_GENESIS="${1:-}"
if [ -n "$CL_GENESIS" ] && [ -f "$CL_GENESIS" ]; then
  echo "  Using provided CL genesis: $CL_GENESIS"
else
  # Use genesis from bootnode's STORY_HOME (prepared by prepare_devnet.sh with correct
  # execution_block_hash and resolved template variables). Bootnode genesis is not modified
  # by configure_periods.sh, so it's the authoritative copy.
  echo "  Copying CL genesis from bootnode (${BOOTNODE_USER}@${BOOTNODE_IP})..."
  CL_GENESIS="/tmp/genesis_cl_backup.json"
  scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} \
    ${BOOTNODE_USER}@${BOOTNODE_IP}:${STORY_HOME}/config/genesis.json \
    "$CL_GENESIS"
  echo "  CL genesis saved to $CL_GENESIS"
fi

# Verify genesis is valid JSON with chain_id
echo "  Verifying genesis..."
CHAIN_ID=$(jq -r '.chain_id' "$CL_GENESIS" 2>/dev/null || echo "")
if [ -z "$CHAIN_ID" ] || [ "$CHAIN_ID" = "null" ]; then
  echo "[ERROR] Genesis is not valid (missing chain_id)"
  exit 1
fi
echo "  chain_id=${CHAIN_ID}"
# Note: Since #726, DKG params come from genesis (InitGenesis), not upgrade handler.
# prepare_devnet.sh / configure_periods.sh set DKG params in bootnode genesis.

# Distribute to all validators
echo "  Distributing CL genesis to all validators..."
for i in $(seq 1 "$TOTAL"); do
  target=$(_ssh_target "$i")
  scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} \
    "$CL_GENESIS" "${target}:${STORY_HOME}/config/genesis.json"
  echo "    node $i: genesis updated"
done

# Also restore to bootnode
scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} \
  "$CL_GENESIS" "${BOOTNODE_USER}@${BOOTNODE_IP}:${STORY_HOME}/config/genesis.json" 2>/dev/null || true

# Resolve genesis template variables (e.g. {{LOCAL_ACCOUNT_ADDRESS}}) on all validators
echo "  Resolving genesis template variables..."
for i in $(seq 1 "$TOTAL"); do
  resolve_genesis_templates "$i"
done

# --- Step 3: Clean chain data on all machines (runbook Step 6) ---
echo "=== Step 3: Cleaning chain data ==="

CLEAN_CMD="
  rm -rf ${STORY_HOME}/data/*
  echo '{\"height\": \"0\", \"round\": 0, \"step\": 0}' > ${STORY_HOME}/data/priv_validator_state.json
  rm -f ${STORY_HOME}/config/write-file-atomic-* ${STORY_HOME}/config/addrbook.json
  sudo rm -rf \
    ${GETH_DATA_DIR}/geth/chaindata \
    ${GETH_DATA_DIR}/geth/lightchaindata \
    ${GETH_DATA_DIR}/geth/blobpool \
    ${GETH_DATA_DIR}/geth/nodes \
    ${GETH_DATA_DIR}/geth/triecache
"

# Clean kernel state + update chain_id in kernel config to match CL genesis
CLEAN_KERNEL_CMD="
  rm -rf ${KERNEL_HOME}/keys/ ${KERNEL_HOME}/dkg_state/ ${KERNEL_HOME}/light_client/ ${KERNEL_HOME}/data/
  sed -i 's/chain_id = .*/chain_id = \"${CHAIN_ID}\"/' ${KERNEL_HOME}/config.toml 2>/dev/null || true
"

for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "$CLEAN_CMD"
  _ssh_cmd "$i" "$CLEAN_KERNEL_CMD" 2>/dev/null || true
  echo "  node $i: data cleaned"
done

# Clean bootnode (no kernel)
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$CLEAN_CMD" || true
echo "  bootnode: data cleaned"

# --- Step 4: Initialize geth on all machines (runbook Step 7) ---
echo "=== Step 4: Initializing geth ==="

for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "geth --state.scheme=hash init --datadir=${GETH_DATA_DIR} ${GETH_GENESIS} 2>&1 | tail -1"
  echo "  node $i: geth initialized"
done

ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} \
  "geth --state.scheme=hash init --datadir=${GETH_DATA_DIR} ${GETH_GENESIS} 2>&1 | tail -1" || true
echo "  bootnode: geth initialized"

# --- Step 4.5: Configure geth.toml FilterLogCacheSize (DCAP mode) ---
if [ "${DKG_DCAP_ENABLED:-}" = "true" ]; then
  echo "=== Step 4.5: Setting FilterLogCacheSize=0 in geth.toml (DCAP mode) ==="
  GETH_TOML_PATH="${DKG_GETH_TOML_PATH:-/home/ubuntu/.story/geth/config/geth.toml}"
  FLC_CMD="
    if [ -f '${GETH_TOML_PATH}' ]; then
      if grep -q 'FilterLogCacheSize' '${GETH_TOML_PATH}'; then
        sed -i 's/FilterLogCacheSize = .*/FilterLogCacheSize = 0/' '${GETH_TOML_PATH}'
      else
        echo 'FilterLogCacheSize = 0' >> '${GETH_TOML_PATH}'
      fi
      echo 'FilterLogCacheSize set to 0'
    else
      echo 'geth.toml not found, skipping'
    fi
  "
  for i in $(seq 1 "$TOTAL"); do
    _ssh_cmd "$i" "$FLC_CMD" 2>/dev/null || true
    echo "  node $i: FilterLogCacheSize configured"
  done
  ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$FLC_CMD" 2>/dev/null || true
  echo "  bootnode: FilterLogCacheSize configured"
fi

# --- Step 5: Start geth on all machines ---
echo "=== Step 5: Starting geth ==="

# Reload systemd units in case binaries or unit files changed
for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "sudo systemctl daemon-reload" 2>/dev/null || true
done
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} \
  "sudo systemctl daemon-reload" 2>/dev/null || true

for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "sudo systemctl start ${GETH_SVC}"
done
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} \
  "sudo systemctl start ${GETH_SVC}" || true
echo "  All geth started. Waiting 5s..."
sleep 5

# --- Step 6: Start story on all machines ---
echo "=== Step 6: Starting story ==="

# Start bootnode first (so validators can connect)
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} \
  "sudo systemctl start ${STORY_SVC}" || true
sleep 3

for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "sudo systemctl start ${STORY_SVC}"
done
echo "  All story started."

# --- Step 7: Wait for blocks ---
echo "=== Step 7: Waiting for chain to produce blocks ==="

MAX_WAIT=120
ELAPSED=0
while [ $ELAPSED -lt $MAX_WAIT ]; do
  HEIGHT=$(_ssh_cmd 1 "curl -s http://localhost:26657/status 2>/dev/null | jq -r .result.sync_info.latest_block_height 2>/dev/null" 2>/dev/null || echo "0")
  if [ -n "$HEIGHT" ] && [ "$HEIGHT" != "null" ] && [ "$HEIGHT" -ge 5 ] 2>/dev/null; then
    echo "  Chain is producing blocks: height=$HEIGHT"
    break
  fi
  echo "  Waiting... (height=${HEIGHT:-0}, elapsed=${ELAPSED}s)"
  sleep 10
  ELAPSED=$((ELAPSED + 10))
done

if [ $ELAPSED -ge $MAX_WAIT ]; then
  echo "[WARNING] Chain did not reach height 5 within ${MAX_WAIT}s. Check logs manually."
  echo "  ssh val1 'journalctl -u story --no-pager -n 20'"
  exit 1
fi

# --- Step 8: Get trusted block for kernel light client (runbook Step 8.1) ---
echo "=== Step 8: Configuring kernel light clients ==="

BLOCK_JSON=$(_ssh_cmd 1 "curl -s 'http://localhost:26657/block?height=5'" 2>/dev/null || echo "")
TRUSTED_HEIGHT=$(echo "$BLOCK_JSON" | jq -r '.result.block.header.height' 2>/dev/null || echo "")
TRUSTED_HASH=$(echo "$BLOCK_JSON" | jq -r '.result.block_id.hash' 2>/dev/null || echo "")

if [ -z "$TRUSTED_HEIGHT" ] || [ "$TRUSTED_HEIGHT" = "null" ] || [ -z "$TRUSTED_HASH" ] || [ "$TRUSTED_HASH" = "null" ]; then
  echo "[ERROR] Could not get trusted block info. Manual kernel config required."
  echo "  Run: curl -s 'http://<val1>:26657/block?height=5' | jq '.result.block_id.hash, .result.block.header.height'"
  exit 1
fi

echo "  Trusted block: height=${TRUSTED_HEIGHT} hash=${TRUSTED_HASH}"

# Update kernel config on each validator
for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "
    sed -i 's/^trusted_height = .*/trusted_height = ${TRUSTED_HEIGHT}/' ${KERNEL_HOME}/config.toml
    sed -i 's/^trusted_hash = .*/trusted_hash = \"${TRUSTED_HASH}\"/' ${KERNEL_HOME}/config.toml
  " 2>/dev/null || true
  echo "  node $i: kernel config updated"
done

# --- Step 8.5: Fix engine-chain-id in story.toml ---
# The DKG contract client signs EL transactions; engine-chain-id must match the geth chain ID.
# Auto-detect from geth and update story.toml on all validators.
echo "=== Step 8.5: Fixing engine-chain-id in story.toml ==="
GETH_CHAIN_ID=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast chain-id --rpc-url http://localhost:8545 2>/dev/null" 2>/dev/null | tr -d '[:space:]')
if [ -n "$GETH_CHAIN_ID" ] && [ "$GETH_CHAIN_ID" != "0" ]; then
  STORY_TOML="${DKG_STORY_TOML_PATH:-${STORY_HOME}/config/story.toml}"
  for i in $(seq 1 "$TOTAL"); do
    _ssh_cmd "$i" "sed -i 's/^engine-chain-id = .*/engine-chain-id = ${GETH_CHAIN_ID}/' ${STORY_TOML}" 2>/dev/null || true
    echo "  node $i: engine-chain-id=${GETH_CHAIN_ID}"
  done
else
  echo "  [WARNING] Could not detect geth chain ID. engine-chain-id not updated."
fi

# --- Step 9: Start kernels (runbook Step 8.3) ---
echo "=== Step 9: Starting kernels (enclave loading takes 1-3 min) ==="

# Start all kernels
for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "sudo systemctl start ${KERNEL_SVC}" 2>/dev/null || true
done

# --- Step 10: Wait for kernel gRPC (poll immediately, no hard sleep) ---
# story has auto-reconnect (#725), so no need to restart story.
# Poll for gRPC readiness on all nodes in parallel.
echo "=== Step 10: Waiting for kernel gRPC ports (polling, up to 5 min) ==="
for i in $(seq 1 "$TOTAL"); do
  (
    for attempt in $(seq 1 150); do
      if _ssh_cmd "$i" "ss -tlnp | grep -q 50051" 2>/dev/null; then
        echo "  node $i: kernel gRPC ready (${attempt}×2s)"
        exit 0
      fi
      sleep 2
    done
    echo "  [ERROR] node $i: kernel gRPC not ready after 5 min"
  ) &
done
wait

# --- Step 10.5: Wait for v2.0.0 upgrade and configure attestation ---
echo "=== Step 10.5: Waiting for v2.0.0 upgrade (block ${DKG_V200_HEIGHT:-150}) ==="

# Wait for block 150+ (V200 upgrade handler activates DKG module)
# LocalChainID has Horace=100, V200=150
V200_HEIGHT="${DKG_V200_HEIGHT:-150}"
for attempt in $(seq 1 180); do
  HEIGHT=$(_ssh_cmd 1 "curl -s http://localhost:26657/status 2>/dev/null | jq -r .result.sync_info.latest_block_height" 2>/dev/null || echo "0")
  if [ "${HEIGHT:-0}" -ge "$V200_HEIGHT" ]; then
    echo "  v2.0.0 upgrade activated at height=$HEIGHT"
    break
  fi
  # Print status every 5th attempt (~10s)
  if [ $((attempt % 5)) -eq 0 ]; then
    echo "  Waiting... (height=${HEIGHT}, elapsed=$((attempt*2))s)"
  fi
  sleep 2
done

if [ "${DKG_DCAP_ENABLED:-}" = "true" ]; then
  # --- DCAP mode: deploy Automata DCAP contracts + Intel collateral ---
  echo "=== Step 10.5a: Deploying DCAP attestation stack ==="
  bash "${SCRIPT_DIR}/deploy_dcap.sh"
else
  # --- Mock mode: whitelist enclave type with MockValidationHook ---
  echo "=== Step 10.5b: Whitelisting enclave type (MockValidationHook) ==="

  # Always get the REAL code commitment from the running kernel
  echo "  Getting real code commitment from kernel..."
  CODE_COMMITMENT=$(_ssh_cmd 1 "gramine-sgx-sigstruct-view /home/ubuntu/story-kernel/story-kernel.sig 2>/dev/null | grep mr_enclave | head -1 | awk '{print \$2}'" 2>/dev/null || echo "")

  if [ -z "$CODE_COMMITMENT" ]; then
    # Fallback: try story logs (kernel reports code_commitment on connect)
    sleep 10
    CODE_COMMITMENT=$(_ssh_cmd 1 "journalctl -u story --no-pager -n 500 2>/dev/null | grep -oP 'code_commitment=\K[a-f0-9]+' | tail -1" 2>/dev/null || echo "")
  fi

  if [ -z "$CODE_COMMITMENT" ]; then
    echo "  [ERROR] Could not determine code_commitment from kernel."
    echo "  You must manually whitelist the enclave type before DKG can start."
  else
    echo "  Kernel code commitment: ${CODE_COMMITMENT}"

    # Check what's currently whitelisted on-chain
    ONCHAIN_CC=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call \
      0xCcCcCC0000000000000000000000000000000004 \
      'enclaveTypeData(bytes32)(bytes32,address)' \
      \$(cast --to-bytes32 1) \
      --rpc-url http://localhost:8545 2>/dev/null | head -1" 2>/dev/null || echo "")

    # Normalize: remove 0x prefix, lowercase for comparison
    ONCHAIN_CC_CLEAN=$(echo "$ONCHAIN_CC" | sed 's/^0x//' | tr '[:upper:]' '[:lower:]')
    KERNEL_CC_CLEAN=$(echo "$CODE_COMMITMENT" | tr '[:upper:]' '[:lower:]')

    if [ "$ONCHAIN_CC_CLEAN" = "$KERNEL_CC_CLEAN" ]; then
      echo "  On-chain code commitment matches kernel. No update needed."
    else
      echo "  On-chain code commitment MISMATCH!"
      echo "    on-chain: ${ONCHAIN_CC_CLEAN}"
      echo "    kernel:   ${KERNEL_CC_CLEAN}"
      echo "  Updating whitelist..."

      # Get MockValidationHook address from on-chain data
      MOCK_HOOK=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call \
        0xCcCcCC0000000000000000000000000000000004 \
        'enclaveTypeData(bytes32)(bytes32,address)' \
        \$(cast --to-bytes32 1) \
        --rpc-url http://localhost:8545 2>/dev/null | tail -1" 2>/dev/null || echo "")

      # Fallback to known MockValidationHook address
      if [ -z "$MOCK_HOOK" ] || [ "$MOCK_HOOK" = "0x0000000000000000000000000000000000000000" ]; then
        MOCK_HOOK="${DKG_MOCK_VALIDATION_HOOK:-0x858F0B7C7c9f440F1bbA90a48Ab4A941C4E44E58}"
      fi

      DKG_OWNER_KEY="${DKG_SIGNER_PRIVATE_KEY:-ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}"

      # Retry whitelist with increasing gas price to handle "replacement transaction underpriced"
      for GAS_PRICE in 20000000000 30000000000 50000000000; do
        WHITELIST_CMD="PATH=\$PATH:\$HOME/.foundry/bin cast send \
          0xCcCcCC0000000000000000000000000000000004 \
          'whitelistEnclaveType(bytes32,(bytes32,address),bool)' \
          \$(cast --to-bytes32 1) \
          '(0x${CODE_COMMITMENT},${MOCK_HOOK})' \
          true \
          --private-key 0x${DKG_OWNER_KEY} \
          --rpc-url http://localhost:8545 \
          --gas-price ${GAS_PRICE} \
          --legacy 2>&1 | grep -E 'status|error'"

        RESULT=$(_ssh_cmd 1 "$WHITELIST_CMD" 2>&1 || echo "whitelist failed")
        echo "  Whitelist result (gas=${GAS_PRICE}): $RESULT"
        if echo "$RESULT" | grep -q "status.*1"; then
          break
        fi
        sleep 3
      done
    fi
  fi
fi

# --- Step 11: Verify DKG ---
echo "=== Step 11: Verifying DKG ==="

DKG_LOG=$(_ssh_cmd 1 "journalctl -u story --no-pager -n 30 2>/dev/null | grep -iE 'Connected.*kernel|DKG|Initiated'" 2>/dev/null || echo "")
if echo "$DKG_LOG" | grep -q "Connected to kernel"; then
  echo "  Kernel connected successfully!"
else
  echo "  [WARNING] No 'Connected to kernel' found in logs. Check manually:"
  echo "    ssh val1 'journalctl -u story --no-pager -n 30 | grep -iE kernel'"
fi

if echo "$DKG_LOG" | grep -q "Initiated new DKG round"; then
  echo "  DKG round initiated!"
else
  echo "  [INFO] DKG round not yet initiated (may need to wait for block 10+)"
fi

# Final status
FINAL_HEIGHT=$(_ssh_cmd 1 "curl -s http://localhost:26657/status 2>/dev/null | jq -r .result.sync_info.latest_block_height" 2>/dev/null || echo "?")
echo ""
echo "=========================================="
echo "  Devnet reset complete!"
echo "  Chain height: ${FINAL_HEIGHT}"
echo "  DKG starts after V200 upgrade"
echo "  v2.0.0 upgrade at block ${V200_HEIGHT}"
echo "=========================================="
echo ""
echo "Wait for block 100+ for DKG to be fully active, then run tests:"
echo "  source tests/integration/dkg/config.env"
echo "  go test -tags=integration -v -run TestDKG_Params -timeout 30s ./tests/integration/dkg/..."
