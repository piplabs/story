#!/usr/bin/env bash
# Common helpers for DKG scenario scripts.
# Source this at the top of every pre.sh / post.sh:
#   source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"
#
# Requires: config.env sourced (or env vars set) before running tests.

set -euo pipefail

# ---------------------------------------------------------------------------
# SSH helper
# ---------------------------------------------------------------------------

# _ssh_target <1-based index>  — extract SSH target from comma-separated DKG_SSH_TARGETS
_ssh_target() {
  local idx="$1"
  echo "$DKG_SSH_TARGETS" | cut -d',' -f"$idx"
}

# _ssh_cmd <node_index> <remote_command>  — run a command on a remote validator via SSH
_ssh_cmd() {
  local idx="$1"; shift
  local target; target=$(_ssh_target "$idx")
  local ssh_opts="${DKG_SSH_OPTS:-}"
  ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} $ssh_opts "$target" "$@"
}

# _node_index returns 1-based index from $1 (default = last validator)
_node_index() { echo "${1:-${DKG_VALIDATOR_COUNT:-3}}"; }

# ---------------------------------------------------------------------------
# Deploy-mode dispatch: stop/start a validator's story process or kernel
# ---------------------------------------------------------------------------

# stop_kernel <node_index>  — stop story-kernel (TEE sidecar) on one node
stop_kernel() {
  local idx; idx=$(_node_index "$1")
  local mode="${DKG_DEPLOY_MODE:-ssh}"
  local svc="${DKG_SYSTEMD_KERNEL_SERVICE:-story-kernel}"
  echo "[_common] stop_kernel node=$idx mode=$mode"
  case "$mode" in
    docker)
      docker stop "${DKG_DOCKER_KERNEL_PREFIX:-story-kernel}${idx}" 2>/dev/null || true
      ;;
    systemd)
      sudo systemctl stop "$svc" || true
      ;;
    ssh)
      _ssh_cmd "$idx" "sudo systemctl stop $svc" || true
      ;;
    *) echo "unknown DKG_DEPLOY_MODE=$mode"; exit 1 ;;
  esac
}

# start_kernel <node_index>  — start story-kernel on one node
start_kernel() {
  local idx; idx=$(_node_index "$1")
  local mode="${DKG_DEPLOY_MODE:-ssh}"
  local svc="${DKG_SYSTEMD_KERNEL_SERVICE:-story-kernel}"
  echo "[_common] start_kernel node=$idx mode=$mode"
  case "$mode" in
    docker)
      docker start "${DKG_DOCKER_KERNEL_PREFIX:-story-kernel}${idx}" 2>/dev/null || true
      ;;
    systemd)
      sudo systemctl start "$svc" || true
      ;;
    ssh)
      _ssh_cmd "$idx" "sudo systemctl start $svc" || true
      ;;
  esac
}

# stop_story <node_index>  — stop the story consensus process on one node
stop_story() {
  local idx; idx=$(_node_index "$1")
  local mode="${DKG_DEPLOY_MODE:-ssh}"
  local svc="${DKG_SYSTEMD_STORY_SERVICE:-story}"
  echo "[_common] stop_story node=$idx mode=$mode"
  case "$mode" in
    docker)
      docker stop "${DKG_DOCKER_VALIDATOR_PREFIX:-validator}${idx}" 2>/dev/null || true
      ;;
    systemd)
      sudo systemctl stop "$svc" || true
      ;;
    ssh)
      _ssh_cmd "$idx" "sudo systemctl stop $svc" || true
      ;;
  esac
}

# start_story <node_index>  — start the story consensus process on one node
start_story() {
  local idx; idx=$(_node_index "$1")
  local mode="${DKG_DEPLOY_MODE:-ssh}"
  local svc="${DKG_SYSTEMD_STORY_SERVICE:-story}"
  echo "[_common] start_story node=$idx mode=$mode"
  case "$mode" in
    docker)
      docker start "${DKG_DOCKER_VALIDATOR_PREFIX:-validator}${idx}" 2>/dev/null || true
      ;;
    systemd)
      sudo systemctl start "$svc" || true
      ;;
    ssh)
      _ssh_cmd "$idx" "sudo systemctl start $svc" || true
      ;;
  esac
}

# restart_story <node_index>  — restart the story consensus process
restart_story() {
  local idx; idx=$(_node_index "$1")
  local mode="${DKG_DEPLOY_MODE:-ssh}"
  local svc="${DKG_SYSTEMD_STORY_SERVICE:-story}"
  echo "[_common] restart_story node=$idx mode=$mode"
  case "$mode" in
    docker)
      docker restart "${DKG_DOCKER_VALIDATOR_PREFIX:-validator}${idx}" 2>/dev/null || true
      ;;
    systemd)
      sudo systemctl restart "$svc" || true
      ;;
    ssh)
      _ssh_cmd "$idx" "sudo systemctl restart $svc" || true
      ;;
  esac
}

# toggle_dkg <node_index> <true|false>  — set dkg.enable on one node and restart story
# Uses scoped sed to only modify the enable line under the [dkg] section.
toggle_dkg() {
  local idx; idx=$(_node_index "$1")
  local enable="${2:-true}"
  local mode="${DKG_DEPLOY_MODE:-ssh}"
  local cfg="${DKG_STORY_TOML_PATH:-/home/ubuntu/.story/story/config/story.toml}"
  echo "[_common] toggle_dkg node=$idx enable=$enable mode=$mode cfg=$cfg"

  # Scoped sed: only modify "enable = ..." under the [dkg] section header
  local sed_cmd="sed -i '/^\\[dkg\\]/,/^\\[/ s/^\\(\\s*\\)enable = .*/\\1enable = $enable/' $cfg"

  case "$mode" in
    docker)
      local cname="${DKG_DOCKER_VALIDATOR_PREFIX:-validator}${idx}"
      docker exec "$cname" bash -c "$sed_cmd" 2>/dev/null || true
      docker restart "$cname" 2>/dev/null || true
      ;;
    systemd)
      eval "$sed_cmd" || true
      sudo systemctl restart "${DKG_SYSTEMD_STORY_SERVICE:-story}" || true
      ;;
    ssh)
      _ssh_cmd "$idx" "$sed_cmd && sudo systemctl restart ${DKG_SYSTEMD_STORY_SERVICE:-story}" || true
      ;;
  esac
}

# force_kill_kernel <node_index>  — multi-pass kill for SGX enclaves that resist cleanup
force_kill_kernel() {
  local idx; idx=$(_node_index "$1")
  local mode="${DKG_DEPLOY_MODE:-ssh}"
  echo "[_common] force_kill_kernel node=$idx mode=$mode"

  local kill_script='
    sudo systemctl stop story-kernel 2>/dev/null || true
    sleep 1
    # Kill by port
    for port in 50051 50052; do
      pid=$(sudo lsof -ti :$port 2>/dev/null || true)
      [ -n "$pid" ] && sudo kill -9 $pid 2>/dev/null || true
    done
    # Kill gramine-sgx and loader processes
    sudo pkill -9 -f gramine-sgx 2>/dev/null || true
    sudo pkill -9 -f "loader.*story-kernel" 2>/dev/null || true
    sleep 1
    # Verify
    if sudo lsof -ti :50051 >/dev/null 2>&1 || sudo lsof -ti :50052 >/dev/null 2>&1; then
      echo "[WARN] ports 50051/50052 still in use after force kill"
    else
      echo "[OK] kernel processes cleaned up"
    fi
  '

  case "$mode" in
    ssh)
      _ssh_cmd "$idx" "$kill_script" || true
      ;;
    systemd)
      eval "$kill_script" || true
      ;;
    *)
      stop_kernel "$idx"
      ;;
  esac
}

# verify_kernel_running <node_index>  — check if story-kernel is responding
verify_kernel_running() {
  local idx; idx=$(_node_index "$1")
  local mode="${DKG_DEPLOY_MODE:-ssh}"
  case "$mode" in
    ssh)
      _ssh_cmd "$idx" "systemctl is-active story-kernel" || { echo "[WARN] kernel not running on node $idx"; return 1; }
      ;;
    systemd)
      systemctl is-active story-kernel || { echo "[WARN] kernel not running"; return 1; }
      ;;
  esac
  echo "[OK] kernel running on node $idx"
}

# verify_story_running <node_index>  — check if story consensus is responding
verify_story_running() {
  local idx; idx=$(_node_index "$1")
  local mode="${DKG_DEPLOY_MODE:-ssh}"
  case "$mode" in
    ssh)
      _ssh_cmd "$idx" "systemctl is-active story" || { echo "[WARN] story not running on node $idx"; return 1; }
      ;;
    systemd)
      systemctl is-active story || { echo "[WARN] story not running"; return 1; }
      ;;
  esac
  echo "[OK] story running on node $idx"
}

# ---------------------------------------------------------------------------
# Bulk operations
# ---------------------------------------------------------------------------

# stop_kernels_except <keep_index>  — stop kernel on all nodes except <keep_index>
stop_kernels_except() {
  local keep="${1:-1}"
  local total="${DKG_VALIDATOR_COUNT:-3}"
  for i in $(seq 1 "$total"); do
    if [ "$i" -ne "$keep" ]; then
      stop_kernel "$i"
    fi
  done
}

# start_all_kernels  — start kernel on all nodes
start_all_kernels() {
  local total="${DKG_VALIDATOR_COUNT:-3}"
  for i in $(seq 1 "$total"); do
    start_kernel "$i"
  done
}

# stop_stories_except <keep_index>  — stop story on all nodes except <keep_index>
stop_stories_except() {
  local keep="${1:-1}"
  local total="${DKG_VALIDATOR_COUNT:-3}"
  for i in $(seq 1 "$total"); do
    if [ "$i" -ne "$keep" ]; then
      stop_story "$i"
    fi
  done
}

# start_all_stories  — start story on all nodes
start_all_stories() {
  local total="${DKG_VALIDATOR_COUNT:-3}"
  for i in $(seq 1 "$total"); do
    start_story "$i"
  done
}

# enable_all_dkg  — set dkg.enable=true on all nodes and restart
enable_all_dkg() {
  local total="${DKG_VALIDATOR_COUNT:-3}"
  for i in $(seq 1 "$total"); do
    toggle_dkg "$i" "true"
  done
}

# wait_blocks <count>  — wait for approximately <count> blocks (rough: 2s each)
wait_blocks() {
  local count="${1:-10}"
  local delay=$(( count * 2 ))
  echo "[_common] waiting ~${delay}s for $count blocks..."
  sleep "$delay"
}

# resolve_genesis_templates <node_index>
# Replaces all {{...}} template vars in CL genesis using the node's validator key.
# Calls `story init` with a temp home to generate a resolved genesis, then merges
# the resolved fields back into the actual genesis (preserving execution_block_hash).
resolve_genesis_templates() {
  local idx="${1:-1}"
  local story_home="${DKG_STORY_HOME:-/home/ubuntu/.story/story}"
  _ssh_cmd "$idx" "python3 - <<'PYEOF'
import json, os, shutil, subprocess, tempfile

genesis_path = '${story_home}/config/genesis.json'
with open(genesis_path) as f:
    content = f.read()

if '{{' not in content:
    print('  No template vars found, skipping.')
    exit(0)

# Use story init to resolve templates: copy validator key to temp dir, run init
tmpdir = tempfile.mkdtemp()
tmp_config = os.path.join(tmpdir, 'config')
tmp_data = os.path.join(tmpdir, 'data')
os.makedirs(tmp_config, exist_ok=True)
os.makedirs(tmp_data, exist_ok=True)

# Copy the real validator key and state so init uses the same identity
shutil.copy('${story_home}/config/priv_validator_key.json', os.path.join(tmp_config, 'priv_validator_key.json'))
# story init expects priv_validator_state.json to exist
with open(os.path.join(tmp_data, 'priv_validator_state.json'), 'w') as f:
    json.dump({'height': '0', 'round': 0, 'step': 0}, f)

# Run story init to generate resolved genesis (existing key file won't be overwritten)
result = subprocess.run(
    ['story', 'init', '--home', tmpdir, '--network', 'local'],
    capture_output=True, text=True
)

resolved_genesis = os.path.join(tmp_config, 'genesis.json')
if not os.path.exists(resolved_genesis):
    print('[ERROR] story init did not produce genesis.json')
    print('stdout:', result.stdout[-200:] if result.stdout else '')
    print('stderr:', result.stderr[-200:] if result.stderr else '')
    shutil.rmtree(tmpdir)
    exit(1)

# Read resolved genesis to get the real values
with open(resolved_genesis) as f:
    resolved = json.load(f)

# Read original genesis (with templates)
with open(genesis_path) as f:
    original = json.load(f)

# Preserve custom fields from original (set by prepare_devnet.sh)
orig_ebh = original['app_state']['evmengine']['params']['execution_block_hash']
orig_dkg = original['app_state'].get('dkg', {})

# Take the resolved app_state (has real addresses in bank, staking, genutil)
original['app_state'] = resolved['app_state']

# Restore our custom fields
original['app_state']['evmengine']['params']['execution_block_hash'] = orig_ebh
# Restore DKG params (short periods for devnet, set by prepare_devnet.sh)
if orig_dkg and orig_dkg.get('params'):
    original['app_state']['dkg'] = orig_dkg

with open(genesis_path, 'w') as f:
    json.dump(original, f, indent=2)

# Verify no templates remain
with open(genesis_path) as f:
    check = f.read()
remaining = check.count('{{')
if remaining > 0:
    print(f'  [WARNING] {remaining} template vars still remain!')
else:
    print('  All template vars resolved.')

shutil.rmtree(tmpdir)
PYEOF
"
}

# ---------------------------------------------------------------------------
# Fake DCAP helpers (for mock kernel scenarios with real DCAP enabled)
# ---------------------------------------------------------------------------

# File to persist real DCAP address across pre/post scripts
DCAP_ADDR_FILE="/tmp/dcap_real_addr.txt"

# deploy_fake_dcap — deploy always-true DCAP contract + switch SGXValidationHook
# Saves real DCAP address to $DCAP_ADDR_FILE for restore_real_dcap.
deploy_fake_dcap() {
  if [ "${DKG_DCAP_ENABLED:-}" != "true" ]; then
    return 0
  fi
  local owner_key="${DKG_OWNER_KEY:-${DKG_SIGNER_PRIVATE_KEY:-ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}}"
  local rpc="http://localhost:8545"
  local dkg_addr="0xCcCcCC0000000000000000000000000000000004"

  echo "[_common] deploying fake DCAP contract..."

  # Get SGXValidationHook proxy address
  local sgx_hook
  sgx_hook=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call ${dkg_addr} 'enclaveTypeData(bytes32)(bytes32,address)' 0x0000000000000000000000000000000000000000000000000000000000000001 --rpc-url ${rpc} 2>/dev/null | tail -1" 2>/dev/null | tr -d '[:space:]')

  # Save real DCAP address
  local real_dcap
  real_dcap=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call ${sgx_hook} 'automataValidationAddr()' --rpc-url ${rpc} 2>/dev/null" 2>/dev/null | tr -d '[:space:]')
  echo "${real_dcap}" > "$DCAP_ADDR_FILE"
  echo "${sgx_hook}" >> "$DCAP_ADDR_FILE"
  echo "[_common] saved real DCAP=${real_dcap} SGX_HOOK=${sgx_hook}"

  # Deploy fake DCAP: verifyAndAttestOnChain(bytes,uint32) returns (true, "")
  # Bytecode: receives any calldata, returns abi.encode(true, bytes(""))
  # Runtime: PUSH1 0x01 PUSH1 0 MSTORE PUSH1 0x40 PUSH1 0 MSTORE ... RETURN
  # Simpler: just return 64 bytes with bool=true at offset 0 + empty bytes at offset 32
  local fake_dcap
  fake_dcap=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin && cast send --private-key 0x${owner_key} --rpc-url ${rpc} --legacy --create '$(cat <<'BYTECODE'
0x608060405234801561001057600080fd5b5060b58061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063098c649914602d575b600080fd5b604080516001815260206060820181905260009082015260800160405180910390f3fea264697066735822122000000000000000000000000000000000000000000000000000000000000000006473
BYTECODE
)' 2>&1 | grep 'contractAddress' | awk '{print \$2}'" 2>/dev/null | tr -d '[:space:]')

  if [ -z "$fake_dcap" ]; then
    # Fallback: deploy minimal always-true contract via raw bytecode
    # Runtime: return (true, empty_bytes) for any call
    # Init code deploys runtime that returns abi.encode(true, 0x40, 0x00) = 96 bytes
    fake_dcap=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast send --private-key 0x${owner_key} --rpc-url ${rpc} --legacy --json --create 0x6080604052348015600f57600080fd5b50607e8061001e6000396000f3fe6080604052600160005260406020526000604052606060006000f3fea164736f6c6343000817000a 2>&1 | python3 -c 'import json,sys; print(json.load(sys.stdin).get(\"contractAddress\",\"\"))'" 2>/dev/null | tr -d '[:space:]')
  fi

  if [ -z "$fake_dcap" ]; then
    echo "[_common] WARN: failed to deploy fake DCAP, mock kernel registration may fail"
    return 1
  fi

  echo "[_common] fake DCAP deployed at ${fake_dcap}"

  # Switch SGXValidationHook to fake DCAP
  _ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast send --private-key 0x${owner_key} --rpc-url ${rpc} --legacy ${sgx_hook} 'setAutomataValidationAddr(address)' ${fake_dcap}" 2>/dev/null || true
  echo "[_common] SGXValidationHook switched to fake DCAP"
}

# restore_real_dcap — restore real DCAP address in SGXValidationHook
restore_real_dcap() {
  if [ "${DKG_DCAP_ENABLED:-}" != "true" ]; then
    return 0
  fi
  if [ ! -f "$DCAP_ADDR_FILE" ]; then
    echo "[_common] no saved DCAP address, skipping restore"
    return 0
  fi
  local owner_key="${DKG_OWNER_KEY:-${DKG_SIGNER_PRIVATE_KEY:-ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}}"
  local rpc="http://localhost:8545"

  local real_dcap sgx_hook
  real_dcap=$(sed -n '1p' "$DCAP_ADDR_FILE")
  sgx_hook=$(sed -n '2p' "$DCAP_ADDR_FILE")

  if [ -n "$real_dcap" ] && [ -n "$sgx_hook" ]; then
    echo "[_common] restoring real DCAP=${real_dcap}"
    _ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast send --private-key 0x${owner_key} --rpc-url ${rpc} --legacy ${sgx_hook} 'setAutomataValidationAddr(address)' ${real_dcap}" 2>/dev/null || true
    echo "[_common] SGXValidationHook restored to real DCAP"
  fi
  rm -f "$DCAP_ADDR_FILE"
}

# ---------------------------------------------------------------------------
# Mock kernel binary auto-build and deploy
# ---------------------------------------------------------------------------

MOCK_KERNEL_LOCAL_BIN="${MOCK_KERNEL_LOCAL_BIN:-/tmp/mock-kernel-server}"
MOCK_KERNEL_REMOTE_BIN="${MOCK_KERNEL_BINARY:-/tmp/mock-kernel-server}"
MOCK_KERNEL_SRC="${MOCK_KERNEL_SRC:-}"  # auto-detect from story-kernel repo

# ensure_mock_binary <node_index> — build (if needed) and deploy mock-kernel-server to remote node.
# Builds for linux/amd64 via cross-compilation, then SCPs to the target.
ensure_mock_binary() {
  local idx; idx=$(_node_index "$1")

  # Check if already present on remote
  if _ssh_cmd "$idx" "test -x ${MOCK_KERNEL_REMOTE_BIN}" 2>/dev/null; then
    echo "[_common] mock-kernel-server already present on node $idx"
    return 0
  fi

  # Auto-detect story-kernel source
  local src_dir="$MOCK_KERNEL_SRC"
  if [ -z "$src_dir" ]; then
    # Try relative to this script's location (../../story-kernel or ../../../story-kernel)
    local script_dir; script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    for candidate in \
      "${script_dir}/../../../../story-kernel" \
      "${script_dir}/../../../story-kernel" \
      "$HOME/cdr/story-kernel" \
      "$HOME/story-kernel"; do
      if [ -f "${candidate}/cmd/mock-kernel-server/main.go" ]; then
        src_dir="$candidate"
        break
      fi
    done
  fi

  if [ -z "$src_dir" ] || [ ! -f "${src_dir}/cmd/mock-kernel-server/main.go" ]; then
    echo "[_common] WARN: story-kernel source not found, cannot auto-build mock-kernel-server"
    echo "  Set MOCK_KERNEL_SRC=/path/to/story-kernel or pre-deploy binary to ${MOCK_KERNEL_REMOTE_BIN}"
    return 1
  fi

  # Build for linux/amd64 (cross-compile from macOS if needed)
  if [ ! -f "$MOCK_KERNEL_LOCAL_BIN" ] || [ "${MOCK_KERNEL_REBUILD:-}" = "true" ]; then
    echo "[_common] building mock-kernel-server from ${src_dir}..."
    (cd "$src_dir" && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "$MOCK_KERNEL_LOCAL_BIN" ./cmd/mock-kernel-server/)
    echo "[_common] built: $MOCK_KERNEL_LOCAL_BIN ($(du -h "$MOCK_KERNEL_LOCAL_BIN" | cut -f1))"
  fi

  # Deploy to remote node
  local target; target=$(_ssh_target "$idx")
  echo "[_common] deploying mock-kernel-server to node $idx ($target)..."
  scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} \
    "$MOCK_KERNEL_LOCAL_BIN" "${target}:${MOCK_KERNEL_REMOTE_BIN}"
  _ssh_cmd "$idx" "chmod +x ${MOCK_KERNEL_REMOTE_BIN}"
  echo "[_common] mock-kernel-server deployed to node $idx"
}

# deploy_mock_kernel <node_index> <mode> — full sequence: ensure binary + stop real kernel + start mock.
deploy_mock_kernel() {
  local idx; idx=$(_node_index "$1")
  local mock_mode="${2:-normal}"
  local mock_port="${MOCK_KERNEL_PORT:-50051}"

  ensure_mock_binary "$idx"
  stop_kernel "$idx"
  sleep 2

  # Start mock kernel
  local cc="${MOCK_CODE_COMMITMENT:-deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef}"
  _ssh_cmd "$idx" "nohup ${MOCK_KERNEL_REMOTE_BIN} --listen=:${mock_port} --mode=${mock_mode} --code-commitment=${cc} > /tmp/mock-kernel.log 2>&1 &"
  sleep 3

  # Verify it started
  if _ssh_cmd "$idx" "ss -tlnp | grep -q ${mock_port}" 2>/dev/null; then
    echo "[_common] mock-kernel-server running on node $idx port $mock_port mode=$mock_mode"
  else
    echo "[_common] WARN: mock-kernel-server may not have started on node $idx"
  fi
}

# stop_mock_kernel <node_index> — stop mock kernel + restart real kernel.
stop_mock_kernel() {
  local idx; idx=$(_node_index "$1")
  _ssh_cmd "$idx" "pkill -f mock-kernel-server 2>/dev/null || true"
  sleep 2
  start_kernel "$idx"
  sleep 5
}

# Mock kernel's well-known fake code commitment (must match mock-kernel-server default)
MOCK_CODE_COMMITMENT="deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
CC_ADDR_FILE="/tmp/.dkg_real_code_commitment"

# whitelist_mock_code_commitment — switch on-chain whitelisted code commitment to mock kernel's fake value.
# Saves the real code commitment for restore_real_code_commitment.
whitelist_mock_code_commitment() {
  local owner_key="${DKG_OWNER_KEY:-${DKG_SIGNER_PRIVATE_KEY:-ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}}"
  local rpc="http://localhost:8545"
  local dkg_addr="0xCcCcCC0000000000000000000000000000000004"
  local enclave_type="0x0000000000000000000000000000000000000000000000000000000000000001"

  # Read current on-chain code commitment + hook address
  local current_cc hook_addr
  current_cc=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call ${dkg_addr} \
    'enclaveTypeData(bytes32)(bytes32,address)' ${enclave_type} \
    --rpc-url ${rpc} 2>/dev/null | head -1" 2>/dev/null | tr -d '[:space:]')
  hook_addr=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call ${dkg_addr} \
    'enclaveTypeData(bytes32)(bytes32,address)' ${enclave_type} \
    --rpc-url ${rpc} 2>/dev/null | tail -1" 2>/dev/null | tr -d '[:space:]')

  # Save real values
  echo "${current_cc}" > "$CC_ADDR_FILE"
  echo "${hook_addr}" >> "$CC_ADDR_FILE"
  echo "[_common] saved real code_commitment=${current_cc} hook=${hook_addr}"

  # Whitelist mock code commitment (same hook address, different code commitment)
  echo "[_common] whitelisting mock code commitment 0x${MOCK_CODE_COMMITMENT}..."
  _ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast send --private-key 0x${owner_key} --rpc-url ${rpc} --legacy \
    ${dkg_addr} 'whitelistEnclaveType(bytes32,(bytes32,address),bool)' \
    ${enclave_type} \
    '(0x${MOCK_CODE_COMMITMENT},${hook_addr})' \
    true" 2>/dev/null || true
  echo "[_common] mock code commitment whitelisted"
}

# restore_real_code_commitment — restore the original on-chain code commitment.
restore_real_code_commitment() {
  if [ ! -f "$CC_ADDR_FILE" ]; then
    echo "[_common] no saved code commitment, skipping restore"
    return 0
  fi
  local owner_key="${DKG_OWNER_KEY:-${DKG_SIGNER_PRIVATE_KEY:-ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}}"
  local rpc="http://localhost:8545"
  local dkg_addr="0xCcCcCC0000000000000000000000000000000004"
  local enclave_type="0x0000000000000000000000000000000000000000000000000000000000000001"

  local real_cc hook_addr
  real_cc=$(sed -n '1p' "$CC_ADDR_FILE")
  hook_addr=$(sed -n '2p' "$CC_ADDR_FILE")

  if [ -n "$real_cc" ] && [ -n "$hook_addr" ]; then
    echo "[_common] restoring real code_commitment=${real_cc}"
    _ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast send --private-key 0x${owner_key} --rpc-url ${rpc} --legacy \
      ${dkg_addr} 'whitelistEnclaveType(bytes32,(bytes32,address),bool)' \
      ${enclave_type} \
      '(${real_cc},${hook_addr})' \
      true" 2>/dev/null || true
    echo "[_common] real code commitment restored"
  fi
  rm -f "$CC_ADDR_FILE"
}

echo "[_common] loaded (mode=${DKG_DEPLOY_MODE:-ssh}, validators=${DKG_VALIDATOR_COUNT:-3})"
