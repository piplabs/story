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

echo "[_common] loaded (mode=${DKG_DEPLOY_MODE:-ssh}, validators=${DKG_VALIDATOR_COUNT:-3})"
