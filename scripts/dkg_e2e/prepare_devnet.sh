#!/usr/bin/env bash
# prepare_devnet.sh — One-time devnet preparation from dkg/dev branch.
#
# This script prepares ALL machines for devnet operation:
#   1. Switches ALL machines to dkg/dev + cherry-picks devnet setup (c61b74a4)
#   2. Gets kernel mrenclave and updates SGX_CODE_COMMITMENT on ALL machines
#   3. Generates EL alloc on val1 (forge script)
#   4. Builds EL genesis + updates CL genesis execution_block_hash
#   5. Builds story binary on val1
#   6. Distributes binary + genesis to all machines
#
# After this script, all machines have:
#   - dkg/dev source + devnet patches (MockValidationHook, short periods, EOA DKG owner)
#   - Correct SGX_CODE_COMMITMENT matching kernel mrenclave
#   - Same story binary and genesis files
#
# Subsequent calls to configure_periods.sh will work correctly because
# all machines have the patched source code.
#
# Usage:
#   source tests/integration/dkg/config.env
#   bash scripts/dkg_e2e/prepare_devnet.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/scenarios/_common.sh"

TOTAL="${DKG_VALIDATOR_COUNT:-3}"
BOOTNODE_IP="${DKG_BOOTNODE_IP:-23.102.71.16}"
BOOTNODE_USER="${DKG_BOOTNODE_USER:-ubuntu}"
STORY_HOME="${DKG_STORY_HOME:-/home/ubuntu/.story/story}"
GETH_GENESIS="${DKG_GETH_GENESIS:-/home/ubuntu/config/genesis-geth.json}"
STORY_SRC="/home/ubuntu/story"
CHAIN_ID="1511"

IFS=',' read -ra SSH_TARGETS <<< "${DKG_SSH_TARGETS}"

echo "=========================================="
echo "  DKG Devnet Preparation (all machines)"
echo "=========================================="

# Devnet setup commit (fixes genesis template addresses, MockValidationHook, EOA owner, etc.)
DEVNET_COMMIT="${DKG_DEVNET_COMMIT:-c61b74a4}"

# --- Step 1: Switch ALL machines to dkg/dev + cherry-pick devnet patch ---
echo "=== Step 1: Switching all machines to dkg/dev + devnet patch ==="

# Pin to a specific commit if DKG_STORY_PIN_COMMIT is set, otherwise use latest branch.
# Use pin when latest branch has breaking changes (e.g. #726 removes V200 upgrade path).
STORY_BRANCH="${DKG_STORY_BRANCH:-dkg/dev}"
STORY_PIN="${DKG_STORY_PIN_COMMIT:-}"
echo "  Branch: ${STORY_BRANCH}${STORY_PIN:+ (pinned: $STORY_PIN)}"
if [ -n "$STORY_PIN" ]; then
  SWITCH_CMD="cd ${STORY_SRC} && \
    git fetch origin ${STORY_BRANCH} && \
    git checkout -f ${STORY_BRANCH} && \
    git reset --hard ${STORY_PIN} && \
    git clean -fd 2>/dev/null || true && \
    git cherry-pick --no-commit ${DEVNET_COMMIT} 2>/dev/null || true"
else
  SWITCH_CMD="cd ${STORY_SRC} && \
    git fetch origin ${STORY_BRANCH} && \
    git checkout -f ${STORY_BRANCH} && \
    git reset --hard origin/${STORY_BRANCH} && \
    git clean -fd 2>/dev/null || true && \
    git cherry-pick --no-commit ${DEVNET_COMMIT} 2>/dev/null || true"
fi

for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "$SWITCH_CMD"
  echo "  node $i: dkg/dev + devnet patch"
done

ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$SWITCH_CMD" 2>/dev/null || true
echo "  bootnode: dkg/dev + devnet patch"

# --- Step 1.5: Patch V200 upgrade height for LocalChainID ---
# The devnet uses LocalChainID (story-1001511) which has V200=2000 by default.
# We need V200 to be low so the DKG module activates quickly after chain reset.
V200_HEIGHT="${DKG_V200_HEIGHT:-150}"
echo "=== Step 1.5: Patching V200 height to ${V200_HEIGHT} for LocalChainID ==="

PATCH_V200_CMD="cd ${STORY_SRC} && sed -i '/LocalChainID: {/,/}/{s/V200:.*[0-9]\+/V200:    ${V200_HEIGHT}/}' lib/netconf/upgrades.go"

for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "$PATCH_V200_CMD"
done
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$PATCH_V200_CMD" 2>/dev/null || true

# Verify
_ssh_cmd 1 "cd ${STORY_SRC} && grep -A6 'LocalChainID:' lib/netconf/upgrades.go | head -7"
echo "  V200 height patched on all machines."

# --- Step 1.6: Patch short DKG periods on all machines ---
REG_PERIOD="${DKG_SHORT_REGISTRATION:-20}"
DEAL_PERIOD="${DKG_SHORT_DEALING:-80}"
FIN_PERIOD="${DKG_SHORT_FINALIZATION:-50}"
ACTIVE_PERIOD="${DKG_SHORT_ACTIVE:-20}"
UPGRADE_FILE="${STORY_SRC}/client/app/upgrades/v_2_0_0/upgrades.go"

echo "=== Step 1.6: Patching DKG periods (reg=${REG_PERIOD} deal=${DEAL_PERIOD} fin=${FIN_PERIOD} active=${ACTIVE_PERIOD}) ==="

PATCH_PERIOD_CMD="cd ${STORY_SRC} && \
  sed -i 's/^\(\t\+\)[0-9]\+,\(\s*\/\/ registration:.*\)/\1${REG_PERIOD},\2/' ${UPGRADE_FILE} && \
  sed -i 's/^\(\t\+\)[0-9]\+,\(\s*\/\/ dealing:.*\)/\1${DEAL_PERIOD},\2/' ${UPGRADE_FILE} && \
  sed -i 's/^\(\t\+\)[0-9]\+,\(\s*\/\/ finalization:.*\)/\1${FIN_PERIOD},\2/' ${UPGRADE_FILE} && \
  sed -i 's/^\(\t\+\)[0-9]\+,\(\s*\/\/ active:.*\)/\1${ACTIVE_PERIOD},\2/' ${UPGRADE_FILE}"

for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "$PATCH_PERIOD_CMD"
done
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$PATCH_PERIOD_CMD" 2>/dev/null || true
echo "  DKG periods patched on all machines."

# --- Step 2: Get kernel mrenclave and update SGX_CODE_COMMITMENT on ALL machines ---
echo "=== Step 2: Getting kernel mrenclave ==="
MRENCLAVE=$(_ssh_cmd 1 "gramine-sgx-sigstruct-view /home/ubuntu/story-kernel/story-kernel.sig 2>/dev/null | grep mr_enclave | head -1 | awk '{print \$2}'")
if [ -z "$MRENCLAVE" ]; then
  echo "[ERROR] Could not get mrenclave from kernel signature"
  exit 1
fi
echo "  mrenclave: ${MRENCLAVE}"

# Update SGX_CODE_COMMITMENT on ALL machines (so configure_periods.sh rebuilds are correct)
UPDATE_CC_CMD="cd ${STORY_SRC}/contracts && \
  sed -i '/SGX_CODE_COMMITMENT/{n;s|hex\"[0-9a-f]*\"|hex\"${MRENCLAVE}\"|;}' script/GenerateAlloc.s.sol"

for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "$UPDATE_CC_CMD"
  echo "  node $i: SGX_CODE_COMMITMENT updated"
done
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$UPDATE_CC_CMD" 2>/dev/null || true
echo "  bootnode: SGX_CODE_COMMITMENT updated"

# Verify on val1
UPDATED=$(_ssh_cmd 1 "cd ${STORY_SRC}/contracts && grep -A1 'SGX_CODE_COMMITMENT =' script/GenerateAlloc.s.sol | head -2")
echo "  Verified: ${UPDATED}"

# --- Step 2.5: DCAP preparation (if enabled) ---
if [ "${DKG_DCAP_ENABLED}" = "true" ]; then
  echo "=== Step 2.5: DCAP preparation ==="

  # Clone automata-dcap-attestation on val1 if not present
  DCAP_REPO="${DKG_DCAP_REPO_PATH:-/home/ubuntu/automata-dcap-attestation}"
  DCAP_BRANCH="${DKG_DCAP_REPO_BRANCH:-main}"
  _ssh_cmd 1 "
    if [ ! -d '${DCAP_REPO}' ]; then
      echo 'Cloning automata-dcap-attestation...'
      git clone --depth 1 -b '${DCAP_BRANCH}' https://github.com/automata-network/automata-dcap-attestation.git '${DCAP_REPO}' 2>&1 | tail -3
    else
      echo 'automata-dcap-attestation already present, updating...'
      cd '${DCAP_REPO}' && git fetch origin '${DCAP_BRANCH}' && git checkout -f '${DCAP_BRANCH}' && git reset --hard 'origin/${DCAP_BRANCH}' 2>&1 | tail -3
    fi
  "

  # Verify forge/cast available
  _ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin && command -v forge >/dev/null && echo 'forge: OK' || echo '[ERROR] forge not found'"
  _ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin && command -v cast >/dev/null && echo 'cast: OK' || echo '[ERROR] cast not found'"

  # Update TCB_EVALUATION_DATA_NUMBER in GenerateAlloc.s.sol on ALL machines
  TCB_EVAL="${DKG_TCB_EVAL_NUMBER:-18}"
  UPDATE_TCB_CMD="cd ${STORY_SRC}/contracts && \
    sed -i 's/TCB_EVALUATION_DATA_NUMBER = [0-9]*/TCB_EVALUATION_DATA_NUMBER = ${TCB_EVAL}/' script/GenerateAlloc.s.sol"

  for i in $(seq 1 "$TOTAL"); do
    _ssh_cmd "$i" "$UPDATE_TCB_CMD"
  done
  ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "$UPDATE_TCB_CMD" 2>/dev/null || true

  echo "  DCAP repo cloned, TCB_EVALUATION_DATA_NUMBER=${TCB_EVAL} on all machines."
fi

# --- Step 3: Generate EL alloc on val1 ---
echo "=== Step 3: Generating EL alloc ==="
_ssh_cmd 1 "
  if ! command -v npm &>/dev/null; then
    echo 'Installing Node.js...'
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash - 2>&1 | tail -1
    sudo apt-get install -y nodejs 2>&1 | tail -1
  fi
  cd ${STORY_SRC}/contracts
  [ -d node_modules ] || npm install 2>&1 | tail -3
  export ADMIN_ADDRESS=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
  export TIMELOCK_EXECUTOR_ADDRESS=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
  export TIMELOCK_GUARDIAN_ADDRESS=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
  PATH=\$PATH:\$HOME/.foundry/bin forge script script/GenerateAlloc.s.sol --tc GenerateAlloc -vvv --chain-id ${CHAIN_ID} 2>&1 | tail -5
"
_ssh_cmd 1 "ls -la ${STORY_SRC}/contracts/local-alloc.json | head -1"
echo "  Alloc generated."

# --- Step 4: Build genesis files ---
echo "=== Step 4: Building genesis files ==="

# Build new EL genesis: replace alloc in existing genesis-geth.json
_ssh_cmd 1 "python3 - <<'PYEOF'
import json

with open('${GETH_GENESIS}') as f:
    geth_genesis = json.load(f)

with open('${STORY_SRC}/contracts/local-alloc.json') as f:
    new_alloc = json.load(f)

geth_genesis['alloc'] = new_alloc

# Bump gasLimit to 30M for DCAP mode (TCB Info upload requires ~22M gas)
dcap_enabled = '${DKG_DCAP_ENABLED}' == 'true'
if dcap_enabled:
    geth_genesis['gasLimit'] = '0x1C9C380'  # 30,000,000
    print('DCAP mode: gasLimit set to 0x1C9C380 (30M)')

with open('/tmp/genesis-geth-new.json', 'w') as f:
    json.dump(geth_genesis, f, indent=2)

print(f'New EL genesis written with {len(new_alloc)} alloc entries')
PYEOF
"

# Compute block 0 hash
NEW_BLOCK_HASH=$(_ssh_cmd 1 "
  TMPDIR=\$(mktemp -d)
  geth --state.scheme=hash init --datadir=\${TMPDIR} /tmp/genesis-geth-new.json 2>/dev/null
  HASH=\$(geth --state.scheme=hash --datadir=\${TMPDIR} --port 0 --http.port 0 --authrpc.port 0 console --exec 'eth.getBlock(0).hash' 2>/dev/null | tr -d '\"')
  rm -rf \${TMPDIR}
  echo \${HASH}
" | tr -d '[:space:]')
echo "  Block 0 hash: ${NEW_BLOCK_HASH}"

if [ -z "$NEW_BLOCK_HASH" ] || [ "$NEW_BLOCK_HASH" = "" ]; then
  echo "[ERROR] Could not compute new block hash"
  exit 1
fi

# Convert hex hash to base64 for CL genesis
NEW_HASH_B64=$(_ssh_cmd 1 "echo -n '${NEW_BLOCK_HASH}' | sed 's/^0x//' | xxd -r -p | base64")
echo "  Base64 execution_block_hash: ${NEW_HASH_B64}"

# Resolve all {{...}} template vars in CL genesis (addresses, validator key, etc.)
echo "  Resolving genesis template variables..."
resolve_genesis_templates 1
echo "  Genesis template resolved on val1."

# Update CL genesis: execution_block_hash + DKG params (short periods for devnet)
_ssh_cmd 1 "python3 - <<'PYEOF'
import json

with open('${STORY_HOME}/config/genesis.json') as f:
    cl_genesis = json.load(f)

# Update execution_block_hash
old_hash = cl_genesis['app_state']['evmengine']['params']['execution_block_hash']
cl_genesis['app_state']['evmengine']['params']['execution_block_hash'] = '${NEW_HASH_B64}'
new_hash = '${NEW_HASH_B64}'
print(f'Updated execution_block_hash: {old_hash[:20]}... -> {new_hash[:20]}...')

# Set DKG params with short periods for devnet testing.
# On #726+, DKG params come from genesis (InitGenesis), not upgrade handler.
# Without this, default params are 1-day periods which is unusable for testing.
dkg_params = {
    'registration_period': ${REG_PERIOD},
    'dealing_period': ${DEAL_PERIOD},
    'finalization_period': ${FIN_PERIOD},
    'active_period': ${ACTIVE_PERIOD},
    'dkg_committee_reward_portion': '0.100000000000000000',
    'min_req_registered_participants': 3,
    'min_req_finalized_participants': 3,
    'operational_threshold': 500
}
if 'dkg' not in cl_genesis['app_state']:
    cl_genesis['app_state']['dkg'] = {}
cl_genesis['app_state']['dkg']['params'] = dkg_params
print(f'Set DKG params: reg={dkg_params[\"registration_period\"]} deal={dkg_params[\"dealing_period\"]} fin={dkg_params[\"finalization_period\"]} active={dkg_params[\"active_period\"]}')

with open('/tmp/genesis-cl-new.json', 'w') as f:
    json.dump(cl_genesis, f, indent=2)
PYEOF
"

_ssh_cmd 1 "cp /tmp/genesis-geth-new.json ${GETH_GENESIS} && cp /tmp/genesis-cl-new.json ${STORY_HOME}/config/genesis.json"
echo "  Genesis files updated on val1."

# --- Step 4.5: Add all validators to genesis ---
# story init only creates gentx for val1. We need to add val2/val3 as validators.
echo "=== Step 4.5: Adding all validators to genesis ==="

# Collect pubkeys from all validators
VAL_PUBKEYS=""
VAL_EVM_ADDRS=""
for i in $(seq 1 "$TOTAL"); do
  PK=$(_ssh_cmd "$i" "python3 -c \"import json; print(json.load(open('${STORY_HOME}/config/priv_validator_key.json'))['pub_key']['value'])\"")
  EVM=$(_ssh_cmd "$i" "story validator export 2>/dev/null | grep 'EVM Address' | awk '{print \$3}'")
  VAL_PUBKEYS="${VAL_PUBKEYS}${PK},"
  VAL_EVM_ADDRS="${VAL_EVM_ADDRS}${EVM},"
  echo "  val$i: pubkey=${PK:0:20}... evm=$EVM"
done

# Add validators to genesis on val1
_ssh_cmd 1 "python3 - <<'PYEOF'
import json, copy

pubkeys = '${VAL_PUBKEYS}'.rstrip(',').split(',')
evm_addrs = '${VAL_EVM_ADDRS}'.rstrip(',').split(',')

with open('${STORY_HOME}/config/genesis.json') as f:
    genesis = json.load(f)

# bech32 encode helper
CHARSET = 'qpzry9x8gf2tvdw0s3jn54khce6mua7l'
def bech32_polymod(values):
    GEN = [0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3]
    chk = 1
    for v in values:
        b = (chk >> 25)
        chk = (chk & 0x1ffffff) << 5 ^ v
        for i in range(5):
            chk ^= GEN[i] if ((b >> i) & 1) else 0
    return chk
def bech32_hrp_expand(h):
    return [ord(x) >> 5 for x in h] + [0] + [ord(x) & 31 for x in h]
def convertbits(data, fb, tb, pad=True):
    acc=0; bits=0; ret=[]
    for v in data:
        acc=(acc<<fb)|v; bits+=fb
        while bits>=tb: bits-=tb; ret.append((acc>>bits)&((1<<tb)-1))
    if pad and bits: ret.append((acc<<(tb-bits))&((1<<tb)-1))
    return ret
def bech32_encode(hrp, addr_bytes):
    d5 = convertbits(list(addr_bytes), 8, 5)
    vals = bech32_hrp_expand(hrp) + d5
    polymod = bech32_polymod(vals + [0]*6) ^ 1
    cs = [(polymod >> 5*(5-i)) & 31 for i in range(6)]
    return hrp + '1' + ''.join([CHARSET[x] for x in d5+cs])

def evm_to_bech32(evm_hex, prefix):
    addr_bytes = bytes.fromhex(evm_hex.replace('0x','').lower())
    return bech32_encode(prefix, addr_bytes)

# Get existing gentx as template
existing_gentxs = genesis['app_state']['genutil']['gen_txs']
template_gentx = existing_gentxs[0]

# Get existing bank balances
balances = genesis['app_state']['bank']['balances']
template_balance = balances[0] if balances else None

new_gentxs = [template_gentx]  # keep val1's gentx

# Get existing auth accounts as template
accounts = genesis['app_state']['auth']['accounts']
template_account = accounts[0] if accounts else None

for i in range(1, len(pubkeys)):
    pk = pubkeys[i]
    evm = evm_addrs[i]
    acc_addr = evm_to_bech32(evm, 'story')
    val_addr = evm_to_bech32(evm, 'storyvaloper')

    # Create gentx for this validator
    gentx = copy.deepcopy(template_gentx)
    msg = gentx['body']['messages'][0]
    msg['delegator_address'] = acc_addr
    msg['validator_address'] = val_addr
    msg['pubkey']['key'] = pk
    msg['description']['moniker'] = f'Validator-{i+1}'
    new_gentxs.append(gentx)

    # Add auth account (required before bank balance and gentx)
    if template_account:
        existing_accs = [a.get('address', '') for a in accounts]
        if acc_addr not in existing_accs:
            new_account = copy.deepcopy(template_account)
            new_account['address'] = acc_addr
            # Increment account_number
            new_account['account_number'] = str(len(accounts))
            accounts.append(new_account)

    # Add bank balance if template exists
    if template_balance:
        existing_addrs = [b['address'] for b in balances]
        if acc_addr not in existing_addrs:
            new_balance = copy.deepcopy(template_balance)
            new_balance['address'] = acc_addr
            balances.append(new_balance)

    print(f'  Added val{i+1}: acc={acc_addr} valoper={val_addr}')

genesis['app_state']['auth']['accounts'] = accounts
genesis['app_state']['genutil']['gen_txs'] = new_gentxs
genesis['app_state']['bank']['balances'] = balances

# Update bank supply to match total balances (prevents "genesis supply is incorrect" panic)
supply = genesis['app_state']['bank'].get('supply', [])
if supply and template_balance:
    # Sum all balances per denom
    denom_totals = {}
    for b in balances:
        for coin in b.get('coins', []):
            d = coin['denom']
            denom_totals[d] = denom_totals.get(d, 0) + int(coin['amount'])
    # Update supply
    new_supply = [{'denom': d, 'amount': str(a)} for d, a in sorted(denom_totals.items())]
    genesis['app_state']['bank']['supply'] = new_supply
    print(f'  Updated bank supply: {new_supply}')

with open('${STORY_HOME}/config/genesis.json', 'w') as f:
    json.dump(genesis, f, indent=2)

print(f'  Genesis now has {len(new_gentxs)} gentxs and {len(balances)} bank balances')
PYEOF
"

# --- Step 5: Build story binary on val1 ---
echo "=== Step 5: Building story binary ==="
_ssh_cmd 1 "cd ${STORY_SRC} && make build 2>&1 | tail -5"
_ssh_cmd 1 "sudo systemctl stop story 2>/dev/null || true"
_ssh_cmd 1 "sudo cp ${STORY_SRC}/build/story /usr/local/bin/story && story version 2>&1 | head -1"
echo "  val1: binary built and installed."

# --- Step 6: Distribute binary + genesis to all machines ---
echo "=== Step 6: Distributing binary and genesis ==="

# Download from val1 to local
echo "  Downloading from val1..."
scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} \
  $(_ssh_target 1):${STORY_SRC}/build/story /tmp/story-dkg-dev
scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} \
  $(_ssh_target 1):${GETH_GENESIS} /tmp/genesis-geth-new.json
# Download the FINAL genesis from val1's STORY_HOME (has DKG params + all validator gentxs)
scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} \
  $(_ssh_target 1):${STORY_HOME}/config/genesis.json /tmp/genesis-cl-new.json

# Deploy to val2..N
for i in $(seq 2 "$TOTAL"); do
  target=$(_ssh_target "$i")
  echo "  Deploying to node $i..."
  _ssh_cmd "$i" "sudo systemctl stop story 2>/dev/null || true"
  scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} /tmp/story-dkg-dev ${target}:/tmp/story-new
  _ssh_cmd "$i" "sudo cp /tmp/story-new /usr/local/bin/story && rm /tmp/story-new && story version 2>&1 | head -1"
  scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} /tmp/genesis-geth-new.json ${target}:${GETH_GENESIS}
  scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} /tmp/genesis-cl-new.json ${target}:${STORY_HOME}/config/genesis.json
  echo "  node $i: done"
done

# Deploy to bootnode
echo "  Deploying to bootnode..."
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "sudo systemctl stop story 2>/dev/null || true"
scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} /tmp/story-dkg-dev ${BOOTNODE_USER}@${BOOTNODE_IP}:/tmp/story-new
ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} ${BOOTNODE_USER}@${BOOTNODE_IP} "sudo cp /tmp/story-new /usr/local/bin/story && rm /tmp/story-new"
scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} /tmp/genesis-geth-new.json ${BOOTNODE_USER}@${BOOTNODE_IP}:${GETH_GENESIS}
scp ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} /tmp/genesis-cl-new.json ${BOOTNODE_USER}@${BOOTNODE_IP}:${STORY_HOME}/config/genesis.json
echo "  bootnode: done"

echo ""
echo "=========================================="
echo "  Preparation complete!"
echo "  mrenclave: ${MRENCLAVE}"
echo "  block0 hash: ${NEW_BLOCK_HASH}"
echo "  All machines: dkg/dev + devnet patches"
echo "=========================================="
echo ""
echo "Now run reset_devnet.sh to start the chain:"
echo "  bash scripts/dkg_e2e/reset_devnet.sh"
