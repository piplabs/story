#!/usr/bin/env bash
# deploy_dcap.sh — Deploy Automata DCAP contracts and upload Intel collateral.
#
# Deploys the full Automata DCAP attestation stack for real on-chain SGX quote
# verification. After this script, DKG.register() calls go through:
#   SGXValidationHook → AutomataDcapAttestationFee → V3QuoteVerifier
#   → X.509 chain + CRL + QE Identity + TCB Info verification (~7.7M gas)
#
# Usage:
#   source tests/integration/dkg/config.env
#   bash scripts/dkg_e2e/deploy_dcap.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/scenarios/_common.sh"

# --- Configuration ---
DCAP_REPO="${DKG_DCAP_REPO_PATH:-/home/ubuntu/automata-dcap-attestation}"
DCAP_BRANCH="${DKG_DCAP_REPO_BRANCH:-main}"
TCB_EVAL="${DKG_TCB_EVAL_NUMBER:-18}"
FMSPC="${DKG_FMSPC:-00606A000000}"
OWNER_KEY="${DKG_OWNER_KEY:-${DKG_SIGNER_PRIVATE_KEY:-ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}}"
COLLATERAL_DIR="${DKG_DCAP_COLLATERAL_DIR:-}"
GAS_PRICE="${DKG_GAS_PRICE:-20000000000}"
RPC_URL="http://localhost:8545"
DKG_ADDR="0xCcCcCC0000000000000000000000000000000004"

# --- Helper: SSH directly (bypass _ssh_cmd for complex commands) ---
val1_ssh() {
  ssh ${DKG_SSH_KEY:+-i "$DKG_SSH_KEY"} ${DKG_SSH_OPTS:-} "$(_ssh_target 1)" "$@"
}

# Auto-detect EL chain ID from geth (ALWAYS query geth, ignore stale env vars)
CHAIN_ID=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin cast chain-id --rpc-url ${RPC_URL}" 2>/dev/null | tr -d '[:space:]')
if [ -z "$CHAIN_ID" ] || ! [[ "$CHAIN_ID" =~ ^[0-9]+$ ]]; then
  echo "  [ERROR] Failed to detect EL chain ID from geth (got: '${CHAIN_ID}')"
  exit 1
fi
echo "  Auto-detected EL chain ID: ${CHAIN_ID}"

# Derive owner address from private key
OWNER_ADDR=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast wallet address --private-key 0x${OWNER_KEY} 2>/dev/null" | tr -d '[:space:]')

echo "=========================================="
echo "  DCAP Attestation Deployment"
echo "=========================================="
echo "  Owner: ${OWNER_ADDR}"
echo "  Chain ID: ${CHAIN_ID}"
echo "  TCB Eval Number: ${TCB_EVAL}"
echo "  FMSPC: ${FMSPC}"
echo ""

# --- Helper: cast send with retry ---
val1_cast_send() {
  local ARGS="$*"
  local GAS="${GAS_PRICE}"
  for retry in 1 2 3; do
    local RESULT
    RESULT=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin && cast send ${ARGS} --private-key 0x${OWNER_KEY} --rpc-url ${RPC_URL} --gas-price ${GAS} --legacy 2>&1" 2>&1 || echo "SEND_ERROR")
    if echo "$RESULT" | grep -q "status.*1"; then
      echo "OK"
      return 0
    fi
    GAS=$((GAS * 2))
    echo "[retry ${retry}] $(echo "$RESULT" | grep -oE 'error.*|Error.*' | head -1)"
    sleep 3
  done
  echo "FAILED"
  return 0  # non-fatal
}

# --- Helper: cast send with remote file data ---
val1_cast_send_with_file() {
  local FILE="$1"; shift
  local ARGS="$*"
  local GAS="${GAS_PRICE}"
  for retry in 1 2 3; do
    local RESULT
    RESULT=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin && HEX=\$(xxd -p ${FILE} | tr -d '\n') && cast send ${ARGS} 0x\${HEX} --private-key 0x${OWNER_KEY} --rpc-url ${RPC_URL} --gas-price ${GAS} --legacy 2>&1" 2>&1 || echo "SEND_ERROR")
    if echo "$RESULT" | grep -q "status.*1"; then
      echo "OK"
      return 0
    fi
    GAS=$((GAS * 2))
    echo "[retry ${retry}] $(echo "$RESULT" | grep -oE 'error.*|Error.*' | head -1)"
    sleep 3
  done
  echo "FAILED"
  return 0
}

# --- Helper: flush stale txpool ---
flush_pending_txs() {
  local GETH_SVC="${DKG_SYSTEMD_GETH_SERVICE:-node-geth}"
  local QUEUED
  QUEUED=$(val1_ssh "curl -s -X POST -H 'Content-Type: application/json' --data '{\"jsonrpc\":\"2.0\",\"method\":\"txpool_status\",\"id\":1}' ${RPC_URL} 2>/dev/null | python3 -c 'import json,sys; d=json.load(sys.stdin).get(\"result\",{}); print(int(d.get(\"queued\",\"0x0\"),16))' 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "0")
  if [ "${QUEUED:-0}" -gt 0 ]; then
    echo "  Stale txpool (queued=${QUEUED}). Restarting geth..."
    _ssh_cmd 1 "sudo systemctl restart ${GETH_SVC}" 2>/dev/null || true
    sleep 5
  fi
}

# ============================================================================
# Step 1: Validate prerequisites
# ============================================================================
echo "=== Step 1: Validating prerequisites ==="

if [ "${DKG_DCAP_ENABLED:-}" != "true" ]; then
  echo "[ERROR] DKG_DCAP_ENABLED is not 'true'."
  exit 1
fi

_ssh_cmd 1 "
  if [ ! -d '${DCAP_REPO}' ]; then
    echo 'Cloning automata-dcap-attestation...'
    git clone --depth 1 -b '${DCAP_BRANCH}' https://github.com/automata-network/automata-dcap-attestation.git '${DCAP_REPO}' 2>&1 | tail -3
  fi
  [ -d '${DCAP_REPO}/evm' ] && echo 'DCAP repo: OK' || { echo '[ERROR] DCAP repo missing'; exit 1; }
"
_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin && forge --version 2>/dev/null | head -1 || { echo '[ERROR] forge not found'; exit 1; }"

HEIGHT=$(_ssh_cmd 1 "curl -s http://localhost:26657/status 2>/dev/null | jq -r .result.sync_info.latest_block_height" 2>/dev/null || echo "0")
if [ "${HEIGHT:-0}" -lt 10 ]; then
  echo "[ERROR] Chain height ${HEIGHT} < 10."
  exit 1
fi
echo "  Chain height: ${HEIGHT}, forge OK"

# ============================================================================
# Step 2: Deploy 15 DCAP contracts
# ============================================================================
echo "=== Step 2: Deploying DCAP contracts ==="

flush_pending_txs

START_NONCE=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin cast nonce ${OWNER_ADDR} --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "0")
echo "  Starting nonce: ${START_NONCE}"

DEPLOY_MAX_RETRIES=3
for deploy_attempt in $(seq 1 "$DEPLOY_MAX_RETRIES"); do
  echo "  Deploy attempt ${deploy_attempt}/${DEPLOY_MAX_RETRIES}..."

  # Clean broadcast cache
  _ssh_cmd 1 "rm -rf '${DCAP_REPO}/evm/broadcast/DeployAllDevnet.s.sol/${CHAIN_ID}/' '${DCAP_REPO}/evm/cache_forge/DeployAllDevnet.s.sol/${CHAIN_ID}/' 2>/dev/null || true"

  # Run forge in background, output to remote file
  val1_ssh "cd '${DCAP_REPO}/evm' && PATH=\$PATH:\$HOME/.foundry/bin && [ -d node_modules ] || npm install 2>&1 | tail -3 && OWNER=${OWNER_ADDR} TCB_EVAL_NUMBER=${TCB_EVAL} forge script forge-script/DeployAllDevnet.s.sol:DeployAllDevnet --rpc-url ${RPC_URL} --private-key 0x${OWNER_KEY} --broadcast --legacy --gas-price ${GAS_PRICE} --slow -vv > /tmp/forge_deploy.log 2>&1" &
  FORGE_PID=$!

  # Progress: poll nonce
  while kill -0 "$FORGE_PID" 2>/dev/null; do
    CUR=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin cast nonce ${OWNER_ADDR} --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "$START_NONCE")
    echo "  ... ${CUR} txs confirmed (nonce ${START_NONCE}→${CUR})"
    sleep 5
  done
  wait "$FORGE_PID" || true

  # Show forge result
  val1_ssh "tail -5 /tmp/forge_deploy.log 2>/dev/null" 2>/dev/null || true

  # Check success by nonce
  END_NONCE=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin cast nonce ${OWNER_ADDR} --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "0")
  CONFIRMED=$((END_NONCE - START_NONCE))
  echo "  Confirmed: ${CONFIRMED} txs (nonce ${START_NONCE}→${END_NONCE})"

  if [ "$CONFIRMED" -ge 20 ]; then
    echo "  Deployment successful."
    break
  fi

  if [ "$deploy_attempt" -lt "$DEPLOY_MAX_RETRIES" ]; then
    echo "  Only ${CONFIRMED} txs, retrying..."
    flush_pending_txs
    START_NONCE="$END_NONCE"
    sleep 5
  else
    echo "  [WARN] Only ${CONFIRMED} txs after ${DEPLOY_MAX_RETRIES} attempts."
  fi
done

# ============================================================================
# Parse contract addresses
# ============================================================================
echo "  Parsing contract addresses..."

# Auto-discover broadcast JSON (forge may use actual geth chain ID, not our CHAIN_ID variable)
BROADCAST_JSON=$(val1_ssh "ls -t ${DCAP_REPO}/evm/broadcast/DeployAllDevnet.s.sol/*/run-latest.json 2>/dev/null | head -1" 2>/dev/null | tr -d '[:space:]')
if [ -z "$BROADCAST_JSON" ]; then
  BROADCAST_JSON="${DCAP_REPO}/evm/broadcast/DeployAllDevnet.s.sol/${CHAIN_ID}/run-latest.json"
fi
echo "  Broadcast JSON: ${BROADCAST_JSON}"

ADDRESSES=$(val1_ssh "jq -r '.transactions[] | select(.transactionType==\"CREATE\") | (.contractName // \"Unknown\") + \"=\" + .contractAddress' '${BROADCAST_JSON}'" 2>&1)
echo "$ADDRESSES"

DCAP_ATTESTATION=$(echo "$ADDRESSES" | grep "AutomataDcapAttestationFee=" | cut -d= -f2)
PCS_DAO=$(echo "$ADDRESSES" | grep "AutomataPcsDao=" | cut -d= -f2)
ENCLAVE_ID_DAO=$(echo "$ADDRESSES" | grep "AutomataEnclaveIdentityDaoVersioned=" | cut -d= -f2)
FMSPC_TCB_DAO=$(echo "$ADDRESSES" | grep "AutomataFmspcTcbDaoVersioned=" | cut -d= -f2)
ROUTER=$(echo "$ADDRESSES" | grep "PCCSRouter=" | cut -d= -f2)
V3_VERIFIER=$(echo "$ADDRESSES" | grep "V3QuoteVerifier=" | cut -d= -f2)

if [ -z "$DCAP_ATTESTATION" ]; then
  echo "[ERROR] Could not parse contract addresses."
  exit 1
fi

# Verify contracts have code
echo "  Verifying on-chain..."
for name_addr in "DcapAttestation=${DCAP_ATTESTATION}" "PcsDao=${PCS_DAO}" "EnclaveIdDao=${ENCLAVE_ID_DAO}" "FmspcTcbDao=${FMSPC_TCB_DAO}" "V3Verifier=${V3_VERIFIER}"; do
  NAME="${name_addr%%=*}"; ADDR="${name_addr#*=}"
  [ -z "$ADDR" ] && continue
  CODE_LEN=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin cast code ${ADDR} --rpc-url ${RPC_URL} 2>/dev/null | wc -c" 2>/dev/null | tr -d '[:space:]')
  if [ "${CODE_LEN:-0}" -gt 100 ]; then
    echo "  ${NAME}: ${ADDR} (${CODE_LEN}c)"
  else
    echo "  [ERROR] ${NAME}: ${ADDR} no code (${CODE_LEN}c)!"
    exit 1
  fi
done

# Save addresses
val1_ssh "printf '%s\n' 'DCAP_ATTESTATION=${DCAP_ATTESTATION}' 'PCS_DAO=${PCS_DAO}' 'ENCLAVE_ID_DAO=${ENCLAVE_ID_DAO}' 'FMSPC_TCB_DAO=${FMSPC_TCB_DAO}' 'ROUTER=${ROUTER}' 'V3_VERIFIER=${V3_VERIFIER}' > /tmp/dcap_addresses.env"

# ============================================================================
# Step 3: Configure SGXValidationHook
# ============================================================================
# Genesis (prepared by prepare_devnet.sh) deploys a real SGXValidationHook proxy
# with DKG owner = TEST_DKG_OWNER (EOA). We just set the DCAP attestation address.
echo "=== Step 3: Configuring SGXValidationHook ==="

SGX_HOOK="${DKG_SGX_HOOK_PROXY:-}"
if [ -z "$SGX_HOOK" ]; then
  SGX_HOOK=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin cast call ${DKG_ADDR} 'enclaveTypeData(bytes32)(bytes32,address)' 0x0000000000000000000000000000000000000000000000000000000000000001 --rpc-url ${RPC_URL} 2>/dev/null | tail -1" 2>/dev/null | tr -d '[:space:]')
fi
echo "  SGXValidationHook: ${SGX_HOOK}"

echo -n "  setAutomataValidationAddr → ${DCAP_ATTESTATION}: "
val1_cast_send "${SGX_HOOK}" "'setAutomataValidationAddr(address)'" "${DCAP_ATTESTATION}"

echo -n "  setTcbEvaluationDataNumber → ${TCB_EVAL}: "
val1_cast_send "${SGX_HOOK}" "'setTcbEvaluationDataNumber(uint32)'" "${TCB_EVAL}"

# Verify
CONFIGURED_ADDR=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin cast call ${SGX_HOOK} 'automataValidationAddr()(address)' --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "?")
CONFIGURED_TCB=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin cast call ${SGX_HOOK} 'tcbEvaluationDataNumber()(uint32)' --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "?")
echo "  Verified: addr=${CONFIGURED_ADDR}"
echo "  Verified: tcb=${CONFIGURED_TCB}"

# ============================================================================
# Step 4: Fetch Intel PCS collateral
# ============================================================================
echo "=== Step 4: Fetching Intel PCS collateral ==="

val1_ssh "
  set -e
  WORK=/tmp/dcap_collateral
  rm -rf \${WORK} && mkdir -p \${WORK} && cd \${WORK}

  echo 'Fetching PCK CRL...'
  curl -sD headers.txt 'https://api.trustedservices.intel.com/sgx/certification/v4/pckcrl?ca=platform&encoding=der' -o pck_crl.der
  CHAIN=\$(grep -i 'SGX-PCK-CRL-Issuer-Chain' headers.txt | sed 's/.*: //' | tr -d '\r')
  python3 << 'PY1'
import urllib.parse, sys, os
chain = os.environ.get('CHAIN','')
if not chain:
    with open('headers.txt') as f:
        for line in f:
            if 'SGX-PCK-CRL-Issuer-Chain' in line:
                chain = line.split(': ',1)[1].strip()
chain = urllib.parse.unquote(chain)
certs = chain.split('-----END CERTIFICATE-----')
names = ['platform_ca.pem', 'root_ca.pem']
idx = 0
for c in certs:
    c = c.strip()
    if c:
        c += '\n-----END CERTIFICATE-----\n'
        if idx < len(names):
            open(names[idx],'w').write(c)
        idx += 1
print('Extracted ' + str(idx) + ' certs')
PY1
  openssl x509 -in platform_ca.pem -outform DER -out platform_ca.der 2>/dev/null
  openssl x509 -in root_ca.pem -outform DER -out root_ca.der 2>/dev/null

  echo 'Fetching TCB Info...'
  curl -sD tcb_headers.txt 'https://api.trustedservices.intel.com/sgx/certification/v4/tcb?fmspc=${FMSPC}' -o tcb_info.json
  python3 << 'PY2'
import urllib.parse, os
with open('tcb_headers.txt') as f:
    for line in f:
        if 'TCB-Info-Issuer-Chain' in line or 'SGX-TCB-Info-Issuer-Chain' in line:
            chain = urllib.parse.unquote(line.split(': ',1)[1].strip())
            certs = chain.split('-----END CERTIFICATE-----')
            for c in certs:
                c = c.strip()
                if c:
                    open('signing_ca.pem','w').write(c + '\n-----END CERTIFICATE-----\n')
                    break
            break
print('Signing cert extracted')
PY2
  openssl x509 -in signing_ca.pem -outform DER -out signing_ca.der 2>/dev/null

  echo 'Fetching Root CRL...'
  curl -s 'https://certificates.trustedservices.intel.com/IntelSGXRootCA.crl' -o root_crl.pem
  openssl crl -in root_crl.pem -inform PEM -outform DER -out root_crl.der 2>/dev/null

  echo 'Fetching QE Identity...'
  curl -s 'https://api.trustedservices.intel.com/sgx/certification/v4/qe/identity' -o qe_identity.json
  jq -r -c '.enclaveIdentity' qe_identity.json > qe_str.txt
  jq -r '.signature' qe_identity.json > qe_sig.txt
  jq -r -c '.tcbInfo' tcb_info.json > tcb_str.txt
  jq -r '.signature' tcb_info.json > tcb_sig.txt

  echo 'Collateral files:'
  ls -la *.der *.txt 2>/dev/null
  echo 'All collateral ready.'
"

# ============================================================================
# Step 5: Upload PCS certificates & CRLs
# ============================================================================
echo "=== Step 5: Uploading PCS certificates ==="

echo -n "  Root CA (type=0): "
val1_cast_send_with_file "/tmp/dcap_collateral/root_ca.der" "${PCS_DAO} 'upsertPcsCertificates(uint8,bytes)' 0"

echo -n "  Platform CA (type=2): "
val1_cast_send_with_file "/tmp/dcap_collateral/platform_ca.der" "${PCS_DAO} 'upsertPcsCertificates(uint8,bytes)' 2"

echo -n "  Signing CA (type=3): "
val1_cast_send_with_file "/tmp/dcap_collateral/signing_ca.der" "${PCS_DAO} 'upsertPcsCertificates(uint8,bytes)' 3"

echo -n "  Root CRL: "
val1_cast_send_with_file "/tmp/dcap_collateral/root_crl.der" "${PCS_DAO} 'upsertRootCACrl(bytes)'"

echo -n "  PCK CRL: "
val1_cast_send_with_file "/tmp/dcap_collateral/pck_crl.der" "${PCS_DAO} 'upsertPckCrl(uint8,bytes)' 2"

# ============================================================================
# Step 6: Grant ATTESTER_ROLE
# ============================================================================
echo "=== Step 6: Granting ATTESTER_ROLE ==="

echo -n "  EnclaveIdentityDao: "
val1_cast_send "${ENCLAVE_ID_DAO}" "'grantRoles(address,uint256)'" "${OWNER_ADDR}" "1"

echo -n "  FmspcTcbDao: "
val1_cast_send "${FMSPC_TCB_DAO}" "'grantRoles(address,uint256)'" "${OWNER_ADDR}" "1"

# ============================================================================
# Step 7: Upload QE Identity
# ============================================================================
echo "=== Step 7: Uploading QE Identity ==="

val1_ssh "
  cd '${DCAP_REPO}/evm'
  PATH=\$PATH:\$HOME/.foundry/bin
  rm -rf broadcast/UploadQEOnly.s.sol cache_forge/UploadQEOnly.s.sol 2>/dev/null
  export OWNER=${OWNER_ADDR}
  export ENCLAVE_ID_DAO=${ENCLAVE_ID_DAO}
  export QE_IDENTITY_STR=\$(cat /tmp/dcap_collateral/qe_str.txt)
  export QE_IDENTITY_SIG=0x\$(cat /tmp/dcap_collateral/qe_sig.txt)
  forge script forge-script/UploadQEOnly.s.sol:UploadQEOnly \
    --rpc-url ${RPC_URL} --private-key 0x${OWNER_KEY} \
    --broadcast --legacy --gas-price ${GAS_PRICE} -vv 2>&1 | tail -5
" || echo "  [WARN] QE upload may have failed"
echo "  QE Identity done."

# ============================================================================
# Step 8: Upload TCB Info
# ============================================================================
echo "=== Step 8: Uploading TCB Info ==="

CURRENT_GAS=$(val1_ssh "PATH=\$PATH:\$HOME/.foundry/bin cast block latest --field gasLimit --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]')
echo "  Block gasLimit: ${CURRENT_GAS}"

if [ -n "$CURRENT_GAS" ] && [ "$CURRENT_GAS" -lt 22000000 ] 2>/dev/null; then
  echo "  Pumping gas limit..."
  val1_ssh "
    PATH=\$PATH:\$HOME/.foundry/bin
    for i in \$(seq 1 300); do
      cast send 0x0000000000000000000000000000000000000001 --value 0 --gas-limit 10000000 \
        --private-key 0x${OWNER_KEY} --rpc-url ${RPC_URL} --gas-price ${GAS_PRICE} --legacy 2>/dev/null &
      if [ \$((i % 50)) -eq 0 ]; then
        wait
        GAS=\$(cast block latest --field gasLimit --rpc-url ${RPC_URL} 2>/dev/null)
        echo \"  \${i} txs, gasLimit=\${GAS}\"
        [ \"\${GAS}\" -ge 22000000 ] 2>/dev/null && break
      fi
    done
    wait
  "
fi

val1_ssh "
  cd '${DCAP_REPO}/evm'
  PATH=\$PATH:\$HOME/.foundry/bin
  rm -rf broadcast/UploadTCBOnly.s.sol cache_forge/UploadTCBOnly.s.sol 2>/dev/null
  export OWNER=${OWNER_ADDR}
  export FMSPC_TCB_DAO=${FMSPC_TCB_DAO}
  export TCB_INFO_STR=\$(cat /tmp/dcap_collateral/tcb_str.txt)
  export TCB_INFO_SIG=0x\$(cat /tmp/dcap_collateral/tcb_sig.txt)
  forge script forge-script/UploadTCBOnly.s.sol:UploadTCBOnly \
    --rpc-url ${RPC_URL} --private-key 0x${OWNER_KEY} \
    --broadcast --legacy --gas-price ${GAS_PRICE} --gas-limit 22000000 -vv 2>&1 | tail -5
" || echo "  [WARN] TCB upload may have failed"
echo "  TCB Info done."

# ============================================================================
# Step 9: Final verification
# ============================================================================
echo "=== Step 9: Verification ==="

echo ""
echo "=========================================="
echo "  DCAP Deployment Complete!"
echo "=========================================="
echo ""
echo "  AutomataDcapAttestationFee : ${DCAP_ATTESTATION}"
echo "  V3QuoteVerifier            : ${V3_VERIFIER}"
echo "  PCCSRouter                 : ${ROUTER}"
echo "  PcsDao                     : ${PCS_DAO}"
echo "  EnclaveIdentityDao         : ${ENCLAVE_ID_DAO}"
echo "  FmspcTcbDao                : ${FMSPC_TCB_DAO}"
echo ""
echo "  SGXValidationHook          : ${SGX_HOOK}"
echo "  Chain ID                   : ${CHAIN_ID}"
echo ""
echo "  Addresses saved: val1:/tmp/dcap_addresses.env"
echo "=========================================="
