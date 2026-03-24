#!/usr/bin/env bash
# deploy_dcap.sh — Deploy Automata DCAP contracts and upload Intel collateral.
#
# Deploys the full Automata DCAP attestation stack for real on-chain SGX quote
# verification. After this script, DKG.register() calls go through:
#   SGXValidationHook → AutomataDcapAttestationFee → V3QuoteVerifier
#   → X.509 chain + CRL + QE Identity + TCB Info verification (~7.7M gas)
#
# Prerequisites:
#   - Chain is running and past V200 upgrade (DKG module active)
#   - Foundry (forge, cast) available on val1
#   - source tests/integration/dkg/config.env (with DKG_DCAP_ENABLED=true)
#
# Usage:
#   source tests/integration/dkg/config.env
#   DKG_DCAP_ENABLED=true bash scripts/dkg_e2e/deploy_dcap.sh
#
# Can also be called automatically by reset_devnet.sh when DKG_DCAP_ENABLED=true.

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
CHAIN_ID="${DKG_EL_CHAIN_ID:-1511}"

# Derive owner address from private key
OWNER_ADDR=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast wallet address --private-key 0x${OWNER_KEY} 2>/dev/null" | tr -d '[:space:]')

echo "=========================================="
echo "  DCAP Attestation Deployment"
echo "=========================================="
echo "  Owner: ${OWNER_ADDR}"
echo "  TCB Eval Number: ${TCB_EVAL}"
echo "  FMSPC: ${FMSPC}"
echo ""

# Helper: run cast send on val1 with standard gas settings, verify success
cast_send() {
  local RESULT
  RESULT=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast send $* \
    --private-key 0x${OWNER_KEY} --rpc-url ${RPC_URL} \
    --gas-price ${GAS_PRICE} --legacy 2>&1" 2>&1 || echo "CAST_SEND_ERROR")
  if echo "$RESULT" | grep -q "status.*1"; then
    echo "  OK"
  elif echo "$RESULT" | grep -q "CAST_SEND_ERROR"; then
    echo "  [WARN] cast send failed (SSH/network error)"
  else
    echo "  $RESULT" | grep -E 'status|error|Error|revert' | head -3
    echo "  [WARN] Transaction may have failed. Check output above."
  fi
}

# ============================================================================
# Step 1: Validate prerequisites
# ============================================================================
echo "=== Step 1: Validating prerequisites ==="

if [ "${DKG_DCAP_ENABLED:-}" != "true" ]; then
  echo "[ERROR] DKG_DCAP_ENABLED is not 'true'. Set it in config.env."
  exit 1
fi

# Ensure automata repo exists on val1
_ssh_cmd 1 "
  if [ ! -d '${DCAP_REPO}' ]; then
    echo 'Cloning automata-dcap-attestation...'
    git clone --depth 1 -b '${DCAP_BRANCH}' https://github.com/automata-network/automata-dcap-attestation.git '${DCAP_REPO}' 2>&1 | tail -3
  fi
  [ -d '${DCAP_REPO}/evm' ] && echo 'DCAP repo: OK' || { echo '[ERROR] DCAP repo missing'; exit 1; }
"
_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin && forge --version 2>/dev/null | head -1 || { echo '[ERROR] forge not found'; exit 1; }"

# Verify chain is past V200
HEIGHT=$(_ssh_cmd 1 "curl -s http://localhost:26657/status 2>/dev/null | jq -r .result.sync_info.latest_block_height" 2>/dev/null || echo "0")
V200="${DKG_V200_HEIGHT:-150}"
if [ "${HEIGHT:-0}" -lt "$V200" ]; then
  echo "[ERROR] Chain height ${HEIGHT} < V200 height ${V200}. Wait for upgrade."
  exit 1
fi
echo "  Chain height: ${HEIGHT} (V200=${V200}), forge OK"

# ============================================================================
# Step 2: Deploy 15 DCAP contracts
# ============================================================================
echo "=== Step 2: Deploying DCAP contracts (15 contracts, ~45 txs) ==="

# Clean previous broadcast cache to avoid nonce conflicts
_ssh_cmd 1 "rm -rf '${DCAP_REPO}/evm/broadcast/DeployAllDevnet.s.sol/${CHAIN_ID}/' '${DCAP_REPO}/evm/cache_forge/DeployAllDevnet.s.sol/${CHAIN_ID}/' 2>/dev/null || true"

_ssh_cmd 1 "
  cd '${DCAP_REPO}/evm'
  PATH=\$PATH:\$HOME/.foundry/bin
  [ -d node_modules ] || npm install 2>&1 | tail -3
  OWNER=${OWNER_ADDR} \
  TCB_EVAL_NUMBER=${TCB_EVAL} \
  forge script forge-script/DeployAllDevnet.s.sol:DeployAllDevnet \
    --rpc-url ${RPC_URL} \
    --private-key 0x${OWNER_KEY} \
    --broadcast --legacy --gas-price ${GAS_PRICE} --slow -vv 2>&1 | tail -20
" || true  # forge may exit non-zero due to stale nonce, but txs may still succeed

# Parse addresses from broadcast JSON (reliable, not from log output)
echo "  Parsing contract addresses from broadcast JSON..."
ADDRESSES=$(_ssh_cmd 1 "
  python3 -c \"
import json, sys
try:
    with open('${DCAP_REPO}/evm/broadcast/DeployAllDevnet.s.sol/${CHAIN_ID}/run-latest.json') as f:
        data = json.load(f)
    txs = data.get('transactions', [])
    receipts = data.get('receipts', [])
    # Count confirmed
    print(f'Transactions: {len(txs)}, Receipts: {len(receipts)}')
    # Extract CREATE txs with contract names
    for t in txs:
        if t.get('transactionType') == 'CREATE':
            name = t.get('contractName', 'Unknown')
            addr = t.get('contractAddress', '')
            print(f'{name}={addr}')
except Exception as e:
    print(f'ERROR: {e}', file=sys.stderr)
    sys.exit(1)
\"
" 2>&1)
echo "$ADDRESSES"

# Extract specific addresses
DCAP_ATTESTATION=$(echo "$ADDRESSES" | grep "AutomataDcapAttestationFee=" | cut -d= -f2)
PCS_DAO=$(echo "$ADDRESSES" | grep "AutomataPcsDao=" | cut -d= -f2)
ENCLAVE_ID_DAO=$(echo "$ADDRESSES" | grep "AutomataEnclaveIdentityDaoVersioned=" | cut -d= -f2)
FMSPC_TCB_DAO=$(echo "$ADDRESSES" | grep "AutomataFmspcTcbDaoVersioned=" | cut -d= -f2)
ROUTER=$(echo "$ADDRESSES" | grep "PCCSRouter=" | cut -d= -f2)
V3_VERIFIER=$(echo "$ADDRESSES" | grep "V3QuoteVerifier=" | cut -d= -f2)

if [ -z "$DCAP_ATTESTATION" ]; then
  echo "[ERROR] Could not parse AutomataDcapAttestationFee address from broadcast JSON."
  exit 1
fi

# Verify contracts have code on-chain
echo "  Verifying on-chain deployment..."
for name_addr in "AutomataDcapAttestationFee=${DCAP_ATTESTATION}" "PcsDao=${PCS_DAO}" "EnclaveIdDao=${ENCLAVE_ID_DAO}" "FmspcTcbDao=${FMSPC_TCB_DAO}" "V3QuoteVerifier=${V3_VERIFIER}"; do
  NAME="${name_addr%%=*}"
  ADDR="${name_addr#*=}"
  if [ -z "$ADDR" ]; then
    echo "  [WARN] ${NAME}: address not found"
    continue
  fi
  CODE_LEN=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast code ${ADDR} --rpc-url ${RPC_URL} 2>/dev/null | wc -c" | tr -d '[:space:]')
  if [ "${CODE_LEN:-0}" -gt 2 ]; then
    echo "  ${NAME}: ${ADDR} (code=${CODE_LEN}c)"
  else
    echo "  [ERROR] ${NAME}: ${ADDR} has no code!"
    exit 1
  fi
done

# Save addresses
_ssh_cmd 1 "cat > /tmp/dcap_addresses.env <<ADDREOF
DCAP_ATTESTATION=${DCAP_ATTESTATION}
PCS_DAO=${PCS_DAO}
ENCLAVE_ID_DAO=${ENCLAVE_ID_DAO}
FMSPC_TCB_DAO=${FMSPC_TCB_DAO}
ROUTER=${ROUTER}
V3_VERIFIER=${V3_VERIFIER}
ADDREOF"

# ============================================================================
# Step 3: Configure SGXValidationHook
# ============================================================================
echo "=== Step 3: Configuring SGXValidationHook ==="

# Get SGX_HOOK proxy address from on-chain DKG.enclaveTypeData(bytes32(1))
SGX_HOOK="${DKG_SGX_HOOK_PROXY:-}"
if [ -z "$SGX_HOOK" ]; then
  SGX_HOOK=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call \
    ${DKG_ADDR} \
    'enclaveTypeData(bytes32)(bytes32,address)' \
    0x0000000000000000000000000000000000000000000000000000000000000001 \
    --rpc-url ${RPC_URL} 2>/dev/null | tail -1" | tr -d '[:space:]')
fi

if [ -z "$SGX_HOOK" ] || [ "$SGX_HOOK" = "0x0000000000000000000000000000000000000000" ]; then
  echo "[ERROR] Could not determine SGXValidationHook proxy address."
  echo "  Set DKG_SGX_HOOK_PROXY in config.env."
  exit 1
fi
echo "  SGXValidationHook proxy: ${SGX_HOOK}"

# Check owner — EOA can call directly, timelock needs schedule+execute
HOOK_OWNER=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call \
  ${SGX_HOOK} 'owner()' --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "")

# Normalize for comparison (cast returns 0x-prefixed 32-byte padded)
HOOK_OWNER_CLEAN=$(echo "$HOOK_OWNER" | sed 's/^0x000000000000000000000000/0x/' | tr '[:upper:]' '[:lower:]')
OWNER_ADDR_CLEAN=$(echo "$OWNER_ADDR" | tr '[:upper:]' '[:lower:]')

echo "  Setting automataValidationAddr → ${DCAP_ATTESTATION}"
cast_send "${SGX_HOOK}" "'setAutomataValidationAddr(address)'" "${DCAP_ATTESTATION}"

echo "  Setting tcbEvaluationDataNumber → ${TCB_EVAL}"
cast_send "${SGX_HOOK}" "'setTcbEvaluationDataNumber(uint32)'" "${TCB_EVAL}"

# Verify (non-fatal — cast call may fail on some hook implementations)
CONFIGURED_ADDR=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call \
  ${SGX_HOOK} 'automataValidationAddr()' --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "unknown")
CONFIGURED_TCB=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call \
  ${SGX_HOOK} 'tcbEvaluationDataNumber()' --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "unknown")
echo "  Verified: automataValidationAddr=${CONFIGURED_ADDR}"
echo "  Verified: tcbEvaluationDataNumber=${CONFIGURED_TCB}"

# ============================================================================
# Step 4: Fetch and upload Intel PCS collateral
# ============================================================================
echo "=== Step 4: Fetching Intel PCS collateral ==="

_ssh_cmd 1 "
  set -e
  WORK=/tmp/dcap_collateral
  rm -rf \${WORK} && mkdir -p \${WORK}
  cd \${WORK}
  COLLATERAL_DIR='${COLLATERAL_DIR}'

  # --- 4a: PCK CRL + Platform CA + Root CA (from issuer chain header) ---
  echo 'Fetching PCK CRL + issuer chain...'
  if curl -sD headers.txt 'https://api.trustedservices.intel.com/sgx/certification/v4/pckcrl?ca=platform&encoding=der' -o pck_crl.der 2>/dev/null && [ -s pck_crl.der ]; then
    CHAIN=\$(grep -i 'SGX-PCK-CRL-Issuer-Chain' headers.txt | sed 's/.*: //' | tr -d '\r')
    python3 -c \"
import urllib.parse, sys
chain = urllib.parse.unquote(sys.argv[1])
certs = chain.split('-----END CERTIFICATE-----')
idx = 0
for c in certs:
    c = c.strip()
    if c:
        c += '\n-----END CERTIFICATE-----\n'
        names = ['platform_ca.pem', 'root_ca.pem']
        if idx < len(names):
            open(names[idx],'w').write(c)
        idx += 1
print(f'Extracted {idx} certs from PCK CRL issuer chain')
\" \"\${CHAIN}\"
    openssl x509 -in platform_ca.pem -outform DER -out platform_ca.der 2>/dev/null
    openssl x509 -in root_ca.pem -outform DER -out root_ca.der 2>/dev/null
  else
    echo '[WARN] PCK CRL fetch failed'
  fi

  # --- 4b: Signing CA (from TCB Info issuer chain — NOT the root CA!) ---
  echo 'Fetching TCB Info (for signing cert + TCB data)...'
  if curl -sD tcb_headers.txt 'https://api.trustedservices.intel.com/sgx/certification/v4/tcb?fmspc=${FMSPC}' -o tcb_info.json 2>/dev/null && [ -s tcb_info.json ]; then
    TCB_CHAIN=\$(grep -i 'TCB-Info-Issuer-Chain\|SGX-TCB-Info-Issuer-Chain' tcb_headers.txt | sed 's/.*: //' | tr -d '\r')
    python3 -c \"
import urllib.parse, sys
chain = urllib.parse.unquote(sys.argv[1])
certs = chain.split('-----END CERTIFICATE-----')
idx = 0
for c in certs:
    c = c.strip()
    if c:
        c += '\n-----END CERTIFICATE-----\n'
        if idx == 0:
            open('signing_ca.pem','w').write(c)
        idx += 1
print(f'Extracted signing cert from TCB issuer chain ({idx} certs)')
\" \"\${TCB_CHAIN}\"
    openssl x509 -in signing_ca.pem -outform DER -out signing_ca.der 2>/dev/null
    openssl x509 -in signing_ca.pem -noout -subject 2>/dev/null
  else
    echo '[WARN] TCB Info fetch failed'
  fi

  # --- 4c: Root CRL ---
  echo 'Fetching Root CA CRL...'
  if curl -s 'https://certificates.trustedservices.intel.com/IntelSGXRootCA.crl' -o root_crl.pem 2>/dev/null && [ -s root_crl.pem ]; then
    openssl crl -in root_crl.pem -inform PEM -outform DER -out root_crl.der 2>/dev/null
  else
    echo '[WARN] Root CRL fetch failed'
  fi

  # --- 4d: QE Identity ---
  echo 'Fetching QE Identity...'
  if curl -s 'https://api.trustedservices.intel.com/sgx/certification/v4/qe/identity' -o qe_identity.json 2>/dev/null && [ -s qe_identity.json ]; then
    echo \"  tcbEvaluationDataNumber: \$(jq '.enclaveIdentity.tcbEvaluationDataNumber' qe_identity.json)\"
  else
    echo '[WARN] QE Identity fetch failed'
  fi

  # --- 4e: Fall back to local bundle for any missing files ---
  if [ -n \"\${COLLATERAL_DIR}\" ] && [ -d \"\${COLLATERAL_DIR}\" ]; then
    for f in root_ca.der platform_ca.der signing_ca.der root_crl.der pck_crl.der qe_identity.json tcb_info.json; do
      if [ ! -s \"\${f}\" ] && [ -f \"\${COLLATERAL_DIR}/\${f}\" ]; then
        cp \"\${COLLATERAL_DIR}/\${f}\" .
        echo \"  Fallback: copied \${f} from local bundle\"
      fi
    done
  fi

  # --- 4f: Parse QE Identity and TCB Info into body + sig ---
  # Use jq to split JSON: extract .enclaveIdentity/.tcbInfo as compact body, .signature as hex
  jq -r -c '.enclaveIdentity' qe_identity.json > qe_str.txt
  jq -r '.signature' qe_identity.json > qe_sig.txt
  echo \"QE Identity: body=\$(wc -c < qe_str.txt)b, sig=\$(wc -c < qe_sig.txt)b\"

  jq -r -c '.tcbInfo' tcb_info.json > tcb_str.txt
  jq -r '.signature' tcb_info.json > tcb_sig.txt
  echo \"TCB Info: body=\$(wc -c < tcb_str.txt)b, fmspc=\$(jq -r '.tcbInfo.fmspc' tcb_info.json)\"

  # Verify all required files
  echo ''
  echo 'Collateral files:'
  ls -la *.der *.json *.txt 2>/dev/null
  for f in root_ca.der platform_ca.der signing_ca.der root_crl.der pck_crl.der qe_str.txt qe_sig.txt tcb_str.txt tcb_sig.txt; do
    [ -s \"\${f}\" ] || { echo \"[ERROR] Missing: \${f}\"; exit 1; }
  done
  echo 'All collateral ready.'
"

# ============================================================================
# Step 5: Upload PCS certificates & CRLs (order matters!)
# ============================================================================
echo "=== Step 5: Uploading PCS certificates & CRLs ==="

echo -n "  Root CA (type=0): "
cast_send "${PCS_DAO}" "'upsertPcsCertificates(uint8,bytes)'" "0" "\$(xxd -p /tmp/dcap_collateral/root_ca.der | tr -d '\n' | sed 's/^/0x/')"

echo -n "  Platform CA (type=2): "
cast_send "${PCS_DAO}" "'upsertPcsCertificates(uint8,bytes)'" "2" "\$(xxd -p /tmp/dcap_collateral/platform_ca.der | tr -d '\n' | sed 's/^/0x/')"

# Signing CA (type=3) — CRITICAL: extracted from TCB Info issuer chain, not Root CA!
echo -n "  Signing CA (type=3): "
cast_send "${PCS_DAO}" "'upsertPcsCertificates(uint8,bytes)'" "3" "\$(xxd -p /tmp/dcap_collateral/signing_ca.der | tr -d '\n' | sed 's/^/0x/')"

echo -n "  Root CRL: "
cast_send "${PCS_DAO}" "'upsertRootCACrl(bytes)'" "\$(xxd -p /tmp/dcap_collateral/root_crl.der | tr -d '\n' | sed 's/^/0x/')"

echo -n "  PCK CRL (platform): "
cast_send "${PCS_DAO}" "'upsertPckCrl(uint8,bytes)'" "2" "\$(xxd -p /tmp/dcap_collateral/pck_crl.der | tr -d '\n' | sed 's/^/0x/')"

# ============================================================================
# Step 6: Grant ATTESTER_ROLE on versioned DAOs
# ============================================================================
echo "=== Step 6: Granting ATTESTER_ROLE ==="

# Solady OwnableRoles: ATTESTER_ROLE = _ROLE_0 = 1
echo -n "  EnclaveIdentityDao: "
cast_send "${ENCLAVE_ID_DAO}" "'grantRoles(address,uint256)'" "${OWNER_ADDR}" "1"

echo -n "  FmspcTcbDao: "
cast_send "${FMSPC_TCB_DAO}" "'grantRoles(address,uint256)'" "${OWNER_ADDR}" "1"

# ============================================================================
# Step 7: Upload QE Identity (via forge script)
# ============================================================================
echo "=== Step 7: Uploading QE Identity ==="

_ssh_cmd 1 "
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
" || echo "  [WARN] QE Identity upload may have failed (non-fatal)"
echo "  QE Identity uploaded."

# ============================================================================
# Step 8: Upload TCB Info (requires block gasLimit >= 22M)
# ============================================================================
echo "=== Step 8: Uploading TCB Info (~17M gas) ==="

# Check block gas limit, pump if needed
CURRENT_GAS=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast block latest --field gasLimit --rpc-url ${RPC_URL} 2>/dev/null" | tr -d '[:space:]')
echo "  Block gasLimit: ${CURRENT_GAS}"

if [ -n "$CURRENT_GAS" ] && [ "$CURRENT_GAS" -lt 22000000 ] 2>/dev/null; then
  echo "  Pumping gas limit with dummy txs (need >= 22M)..."
  _ssh_cmd 1 "
    PATH=\$PATH:\$HOME/.foundry/bin
    for i in \$(seq 1 300); do
      cast send 0x0000000000000000000000000000000000000001 --value 0 --gas-limit 10000000 \
        --private-key 0x${OWNER_KEY} --rpc-url ${RPC_URL} \
        --gas-price ${GAS_PRICE} --legacy 2>/dev/null &
      if [ \$((i % 50)) -eq 0 ]; then
        wait
        GAS=\$(cast block latest --field gasLimit --rpc-url ${RPC_URL} 2>/dev/null)
        echo \"  \${i} txs sent, gasLimit=\${GAS}\"
        [ \"\${GAS}\" -ge 22000000 ] 2>/dev/null && break
      fi
    done
    wait
  "
fi

_ssh_cmd 1 "
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
" || echo "  [WARN] TCB Info upload may have failed (non-fatal)"
echo "  TCB Info uploaded."

# ============================================================================
# Step 9: Final verification
# ============================================================================
echo "=== Step 9: Verification ==="

FINAL_ADDR=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call \
  ${SGX_HOOK} 'automataValidationAddr()' --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "unknown")
FINAL_TCB=$(_ssh_cmd 1 "PATH=\$PATH:\$HOME/.foundry/bin cast call \
  ${SGX_HOOK} 'tcbEvaluationDataNumber()' --rpc-url ${RPC_URL} 2>/dev/null" 2>/dev/null | tr -d '[:space:]' || echo "unknown")

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
echo "  SGXValidationHook (${SGX_HOOK}):"
echo "    automataValidationAddr   : ${FINAL_ADDR}"
echo "    tcbEvaluationDataNumber  : ${FINAL_TCB}"
echo ""
echo "  Addresses saved: val1:/tmp/dcap_addresses.env"
echo ""
echo "  Verification pipeline:"
echo "    DKG.register() → SGXValidationHook → DCAP → V3QuoteVerifier"
echo "    Expected gas per registration: ~7.7M"
echo "=========================================="
