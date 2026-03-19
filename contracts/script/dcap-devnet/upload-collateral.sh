#!/usr/bin/env bash
# upload-collateral.sh — Upload Intel DCAP collateral to on-chain PCCS on Story devnet
#
# This script fetches Intel SGX certificates, CRLs, TCB info, and QE identity
# from the Intel PCS API and uploads them to the on-chain PCCS DAO contracts.
#
# Prerequisites:
#   - deployment-addresses.json from deploy-dcap-devnet.sh
#   - cast (foundry) installed
#   - openssl installed (for PEM→DER conversion)
#   - curl installed
#   - jq installed
#
# Required env vars:
#   DEVNET_RPC_URL       — devnet EL RPC URL
#   DEPLOYER_PRIVATE_KEY — private key with DAO writer access
#   FMSPC                — FMSPC hex of devnet SGX machines (e.g., "00906ED50000")
#
# Optional env vars:
#   INTEL_PCS_URL        — Intel PCS API base URL (default: https://api.trustedservices.intel.com)
#   PCK_CA               — platform or processor (default: platform)
#   TCB_EVAL_NUMBER      — TCB evaluation data number (default: 21)
#
# To find your machine's FMSPC, run on the SGX machine:
#   # Method 1: from PCK cert
#   sudo pccs_retrieval_tool -f /tmp/pck_cert.pem && openssl x509 -in /tmp/pck_cert.pem -noout -text | grep -A1 FMSPC
#   # Method 2: from /dev/attestation/quote
#   # Parse the quote header to extract FMSPC from the certification data

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_FILE="${SCRIPT_DIR}/deployment-addresses.json"
WORK_DIR="${SCRIPT_DIR}/.collateral-cache"

INTEL_PCS_URL="${INTEL_PCS_URL:-https://api.trustedservices.intel.com}"
PCK_CA="${PCK_CA:-platform}"
TCB_EVAL_NUMBER="${TCB_EVAL_NUMBER:-21}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()   { echo -e "${RED}[ERROR]${NC} $1"; }
step()  { echo -e "\n${CYAN}═══ $1 ═══${NC}"; }

# ── Preflight ────────────────────────────────────────────────────────────────
preflight() {
    for var in DEVNET_RPC_URL DEPLOYER_PRIVATE_KEY FMSPC; do
        if [ -z "${!var:-}" ]; then
            err "$var is not set"
            exit 1
        fi
    done
    for cmd in cast curl openssl jq; do
        if ! command -v "$cmd" &>/dev/null; then
            err "$cmd not found"
            exit 1
        fi
    done
    if [ ! -f "$DEPLOY_FILE" ]; then
        err "deployment-addresses.json not found. Run deploy-dcap-devnet.sh first."
        exit 1
    fi

    mkdir -p "$WORK_DIR"
    CHAIN_ID=$(cast chain-id --rpc-url "$DEVNET_RPC_URL")
    info "Chain ID: $CHAIN_ID"
    info "FMSPC: $FMSPC"
    info "Intel PCS URL: $INTEL_PCS_URL"
    info "PCK CA: $PCK_CA"
}

# Helper: read address from deployment
read_addr() {
    jq -r ".[\"$1\"]" "$DEPLOY_FILE"
}

# Helper: convert PEM to DER hex
pem_to_der_hex() {
    local pem_file="$1"
    openssl x509 -in "$pem_file" -outform DER 2>/dev/null | xxd -p | tr -d '\n'
}

# Helper: convert CRL PEM to DER hex
crl_pem_to_der_hex() {
    local pem_file="$1"
    openssl crl -in "$pem_file" -outform DER 2>/dev/null | xxd -p | tr -d '\n'
}

# ── Step 1: Fetch Intel Certificates ─────────────────────────────────────────
fetch_certs() {
    step "Fetching Intel certificates from PCS API"

    # Fetch PCK CRL (contains issuer chain in response header)
    info "Fetching PCK CRL (ca=${PCK_CA})..."
    curl -s -D "${WORK_DIR}/pckcrl_headers.txt" \
        "${INTEL_PCS_URL}/sgx/certification/v4/pckcrl?ca=${PCK_CA}&encoding=pem" \
        -o "${WORK_DIR}/pck_crl.pem"

    # Extract issuer chain from header (URL-encoded PEM certs)
    local issuer_chain
    issuer_chain=$(grep -i "SGX-PCK-CRL-Issuer-Chain" "${WORK_DIR}/pckcrl_headers.txt" | \
        sed 's/SGX-PCK-CRL-Issuer-Chain: //i' | tr -d '\r')

    if [ -n "$issuer_chain" ]; then
        # URL-decode the issuer chain
        local decoded_chain
        decoded_chain=$(python3 -c "import urllib.parse; print(urllib.parse.unquote('$issuer_chain'))" 2>/dev/null || \
            printf '%b' "${issuer_chain//%/\\x}")

        # Split the chain into individual certificates
        echo "$decoded_chain" | awk '/-----BEGIN CERTIFICATE-----/{f=1; n++} f{print > ("'"${WORK_DIR}"'/issuer_cert_" n ".pem")} /-----END CERTIFICATE-----/{f=0}'
        info "  Extracted $(ls ${WORK_DIR}/issuer_cert_*.pem 2>/dev/null | wc -l | tr -d ' ') certificates from issuer chain"
    fi

    # Fetch Root CA CRL
    info "Fetching Root CA CRL..."
    curl -s "${INTEL_PCS_URL}/sgx/certification/v4/rootcacrl?encoding=pem" \
        -o "${WORK_DIR}/root_ca_crl.pem"

    # Fetch QE Identity
    info "Fetching QE Identity..."
    curl -s -D "${WORK_DIR}/qeidentity_headers.txt" \
        "${INTEL_PCS_URL}/sgx/certification/v4/qe/identity" \
        -o "${WORK_DIR}/qe_identity.json"

    # Extract signing cert chain from QE Identity response
    local qe_issuer_chain
    qe_issuer_chain=$(grep -i "SGX-Enclave-Identity-Issuer-Chain" "${WORK_DIR}/qeidentity_headers.txt" | \
        sed 's/SGX-Enclave-Identity-Issuer-Chain: //i' | tr -d '\r')

    if [ -n "$qe_issuer_chain" ]; then
        local decoded_qe_chain
        decoded_qe_chain=$(python3 -c "import urllib.parse; print(urllib.parse.unquote('$qe_issuer_chain'))" 2>/dev/null || echo "")
        echo "$decoded_qe_chain" | awk '/-----BEGIN CERTIFICATE-----/{f=1; n++} f{print > ("'"${WORK_DIR}"'/signing_cert_" n ".pem")} /-----END CERTIFICATE-----/{f=0}'
    fi

    # Fetch TCB Info for our FMSPC
    info "Fetching TCB Info for FMSPC=$FMSPC..."
    curl -s -D "${WORK_DIR}/tcbinfo_headers.txt" \
        "${INTEL_PCS_URL}/sgx/certification/v4/tcb?fmspc=${FMSPC}" \
        -o "${WORK_DIR}/tcb_info.json"

    info "All collateral fetched to ${WORK_DIR}/"
}

# ── Step 2: Upload PCS Certificates ─────────────────────────────────────────
upload_pcs_certs() {
    step "Uploading PCS Certificates"
    local pcs_dao=$(read_addr "AutomataPcsDao")
    info "PcsDao: $pcs_dao"

    # Upload each certificate found in the issuer chain
    # Typically: cert 1 = Intermediate CA (Platform or Processor), cert 2 = Root CA
    # Plus the signing certificates from QE Identity response

    # Try to identify and upload Root CA
    for cert_file in ${WORK_DIR}/issuer_cert_*.pem ${WORK_DIR}/signing_cert_*.pem; do
        [ -f "$cert_file" ] || continue
        local subject
        subject=$(openssl x509 -in "$cert_file" -noout -subject 2>/dev/null || echo "")
        local der_hex
        der_hex=$(pem_to_der_hex "$cert_file")

        if echo "$subject" | grep -qi "Root CA"; then
            info "Uploading Root CA certificate..."
            cast send "$pcs_dao" \
                "upsertPcsCertificates(uint8,bytes)" \
                0 "0x${der_hex}" \
                --rpc-url "$DEVNET_RPC_URL" \
                --private-key "$DEPLOYER_PRIVATE_KEY" \
                --gas-limit 5000000
            info "  Root CA uploaded"
        elif echo "$subject" | grep -qi "Platform CA"; then
            info "Uploading Platform CA certificate..."
            cast send "$pcs_dao" \
                "upsertPcsCertificates(uint8,bytes)" \
                2 "0x${der_hex}" \
                --rpc-url "$DEVNET_RPC_URL" \
                --private-key "$DEPLOYER_PRIVATE_KEY" \
                --gas-limit 5000000
            info "  Platform CA uploaded"
        elif echo "$subject" | grep -qi "Processor CA"; then
            info "Uploading Processor CA certificate..."
            cast send "$pcs_dao" \
                "upsertPcsCertificates(uint8,bytes)" \
                1 "0x${der_hex}" \
                --rpc-url "$DEVNET_RPC_URL" \
                --private-key "$DEPLOYER_PRIVATE_KEY" \
                --gas-limit 5000000
            info "  Processor CA uploaded"
        elif echo "$subject" | grep -qi "Signing"; then
            info "Uploading TCB Signing certificate..."
            cast send "$pcs_dao" \
                "upsertPcsCertificates(uint8,bytes)" \
                3 "0x${der_hex}" \
                --rpc-url "$DEVNET_RPC_URL" \
                --private-key "$DEPLOYER_PRIVATE_KEY" \
                --gas-limit 5000000
            info "  TCB Signing CA uploaded"
        else
            warn "  Unknown certificate: $subject (skipping)"
        fi
    done
}

# ── Step 3: Upload CRLs ─────────────────────────────────────────────────────
upload_crls() {
    step "Uploading CRLs"
    local pcs_dao=$(read_addr "AutomataPcsDao")

    # Upload Root CA CRL
    if [ -f "${WORK_DIR}/root_ca_crl.pem" ]; then
        info "Uploading Root CA CRL..."
        local root_crl_hex
        root_crl_hex=$(crl_pem_to_der_hex "${WORK_DIR}/root_ca_crl.pem")
        cast send "$pcs_dao" \
            "upsertRootCACrl(bytes)" \
            "0x${root_crl_hex}" \
            --rpc-url "$DEVNET_RPC_URL" \
            --private-key "$DEPLOYER_PRIVATE_KEY" \
            --gas-limit 3000000
        info "  Root CA CRL uploaded"
    fi

    # Upload PCK CRL (Platform or Processor)
    if [ -f "${WORK_DIR}/pck_crl.pem" ]; then
        local ca_enum
        if [ "$PCK_CA" = "platform" ]; then
            ca_enum=2
        else
            ca_enum=1
        fi
        info "Uploading PCK CRL (ca=$PCK_CA, enum=$ca_enum)..."
        local pck_crl_hex
        pck_crl_hex=$(crl_pem_to_der_hex "${WORK_DIR}/pck_crl.pem")
        cast send "$pcs_dao" \
            "upsertPckCrl(uint8,bytes)" \
            "$ca_enum" "0x${pck_crl_hex}" \
            --rpc-url "$DEVNET_RPC_URL" \
            --private-key "$DEPLOYER_PRIVATE_KEY" \
            --gas-limit 3000000
        info "  PCK CRL uploaded"
    fi
}

# ── Step 4: Upload QE Identity ───────────────────────────────────────────────
upload_qe_identity() {
    step "Uploading QE Identity"
    local enclave_dao_key="AutomataEnclaveIdentityDaoVersioned_tcbeval_${TCB_EVAL_NUMBER}"
    local enclave_dao=$(read_addr "$enclave_dao_key")
    info "EnclaveIdentityDao (tcbeval_${TCB_EVAL_NUMBER}): $enclave_dao"

    if [ ! -f "${WORK_DIR}/qe_identity.json" ]; then
        err "qe_identity.json not found"
        return 1
    fi

    # The QE Identity response has format: {"enclaveIdentity": {...}, "signature": "..."}
    local identity_str
    identity_str=$(jq -c '.enclaveIdentity' "${WORK_DIR}/qe_identity.json")
    local signature
    signature=$(jq -r '.signature' "${WORK_DIR}/qe_identity.json")

    info "Uploading QE Identity (id=0, version=4)..."
    # EnclaveIdentityJsonObj = (string identityStr, bytes signature)
    # upsertEnclaveIdentity(uint256 id, uint256 version, (string,bytes) enclaveIdentityObj)
    cast send "$enclave_dao" \
        "upsertEnclaveIdentity(uint256,uint256,(string,bytes))" \
        0 4 "(\"$(echo "$identity_str" | sed 's/"/\\"/g')\",0x${signature})" \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --gas-limit 10000000
    info "  QE Identity uploaded"
}

# ── Step 5: Upload TCB Info ──────────────────────────────────────────────────
upload_tcb_info() {
    step "Uploading TCB Info for FMSPC=$FMSPC"
    local fmspc_dao_key="AutomataFmspcTcbDaoVersioned_tcbeval_${TCB_EVAL_NUMBER}"
    local fmspc_dao=$(read_addr "$fmspc_dao_key")
    info "FmspcTcbDao (tcbeval_${TCB_EVAL_NUMBER}): $fmspc_dao"

    if [ ! -f "${WORK_DIR}/tcb_info.json" ]; then
        err "tcb_info.json not found"
        return 1
    fi

    # The TCB Info response has format: {"tcbInfo": {...}, "signature": "..."}
    local tcb_info_str
    tcb_info_str=$(jq -c '.tcbInfo' "${WORK_DIR}/tcb_info.json")
    local signature
    signature=$(jq -r '.signature' "${WORK_DIR}/tcb_info.json")

    info "Uploading FMSPC TCB Info..."
    # TcbInfoJsonObj = (string tcbInfoStr, bytes signature)
    # upsertFmspcTcb((string,bytes) tcbInfoObj)
    cast send "$fmspc_dao" \
        "upsertFmspcTcb((string,bytes))" \
        "(\"$(echo "$tcb_info_str" | sed 's/"/\\"/g')\",0x${signature})" \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --gas-limit 15000000
    info "  TCB Info uploaded"
}

# ── Step 6: Summary ──────────────────────────────────────────────────────────
collateral_summary() {
    step "Collateral Upload Summary"

    echo ""
    echo "  Uploaded collateral:"
    echo "    - Intel Root CA certificate"
    echo "    - Intel ${PCK_CA^} CA certificate"
    echo "    - Intel TCB Signing certificate"
    echo "    - Root CA CRL"
    echo "    - PCK CRL (${PCK_CA})"
    echo "    - QE Identity (id=0, version=4)"
    echo "    - FMSPC TCB Info (fmspc=${FMSPC})"
    echo ""
    echo -e "  ${GREEN}Collateral upload complete!${NC}"
    echo ""
    echo "  You can verify collateral status using the automata-dcap CLI:"
    echo "    automata-dcap --rpc-url $DEVNET_RPC_URL qpl status"
    echo ""
    echo "  Or check with a quote:"
    echo "    automata-dcap --rpc-url $DEVNET_RPC_URL --quote-path <quote.bin> qpl check"
}

# ── Main ─────────────────────────────────────────────────────────────────────
preflight
fetch_certs
upload_pcs_certs
upload_crls
upload_qe_identity
upload_tcb_info
collateral_summary
