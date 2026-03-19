#!/usr/bin/env bash
# deploy-dcap-devnet.sh — Deploy full Automata DCAP attestation stack to Story devnet
#
# Prerequisites:
#   - forge (foundry) installed
#   - automata-on-chain-pccs repo cloned at $PCCS_REPO
#   - automata-dcap-attestation repo cloned at $DCAP_REPO
#
# Usage:
#   export DEVNET_RPC_URL="http://<devnet-rpc>:8545"
#   export DEPLOYER_PRIVATE_KEY="0x..."
#   ./deploy-dcap-devnet.sh [step]
#
# Steps (run in order, or specify one):
#   all              — Run all steps (default)
#   1-helpers        — Deploy PCCS helper contracts
#   2-dao            — Deploy PCCS DAO contracts
#   3-versioned      — Deploy versioned DAO + TcbEvalDao
#   4-router         — Deploy PCCSRouter
#   5-attestation    — Deploy AutomataDcapAttestationFee
#   6-verifier       — Deploy V3QuoteVerifier + register
#   7-summary        — Print deployment summary

set -euo pipefail

# ── Configuration ──────────────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PCCS_REPO="${PCCS_REPO:-/Users/hanslee/repos/storyprotocol/automata-on-chain-pccs}"
DCAP_REPO="${DCAP_REPO:-/Users/hanslee/repos/storyprotocol/automata-dcap-attestation}"
DCAP_EVM="${DCAP_REPO}/evm"

# TCB evaluation data number — use latest known (21 as of 2026-03)
# Set to 0 in SGXValidationHook to accept any TCB level for testing
TCB_EVAL_NUMBER="${TCB_EVAL_NUMBER:-21}"

# Output file for deployment addresses
DEPLOY_OUTPUT="${SCRIPT_DIR}/deployment-addresses.json"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()   { echo -e "${RED}[ERROR]${NC} $1"; }
step()  { echo -e "\n${CYAN}═══ $1 ═══${NC}"; }

# ── Preflight checks ──────────────────────────────────────────────────────────
preflight() {
    if [ -z "${DEVNET_RPC_URL:-}" ]; then
        err "DEVNET_RPC_URL is not set"
        echo "  export DEVNET_RPC_URL=\"http://<devnet-rpc>:8545\""
        exit 1
    fi
    if [ -z "${DEPLOYER_PRIVATE_KEY:-}" ]; then
        err "DEPLOYER_PRIVATE_KEY is not set"
        exit 1
    fi
    if ! command -v forge &>/dev/null; then
        err "forge not found. Install foundry: https://book.getfoundry.sh"
        exit 1
    fi
    if ! command -v cast &>/dev/null; then
        err "cast not found. Install foundry: https://book.getfoundry.sh"
        exit 1
    fi
    if [ ! -d "$PCCS_REPO" ]; then
        err "PCCS repo not found at $PCCS_REPO"
        exit 1
    fi
    if [ ! -d "$DCAP_EVM" ]; then
        err "DCAP EVM repo not found at $DCAP_EVM"
        exit 1
    fi

    CHAIN_ID=$(cast chain-id --rpc-url "$DEVNET_RPC_URL")
    OWNER=$(cast wallet address --private-key "$DEPLOYER_PRIVATE_KEY")
    info "Chain ID: $CHAIN_ID"
    info "Deployer/Owner: $OWNER"
    info "RPC: $DEVNET_RPC_URL"

    # Create deployment dir for PCCS
    mkdir -p "${PCCS_REPO}/deployment"
    # Ensure a JSON exists for this chain
    if [ ! -f "${PCCS_REPO}/deployment/${CHAIN_ID}.json" ]; then
        echo "{}" > "${PCCS_REPO}/deployment/${CHAIN_ID}.json"
    fi

    # Create deployment dir for DCAP (the scripts expect this path structure)
    DCAP_DEPLOY_DIR="${DCAP_EVM}/../rust-crates/libraries/network-registry/deployment/current/${CHAIN_ID}"
    mkdir -p "$DCAP_DEPLOY_DIR"
    if [ ! -f "$DCAP_DEPLOY_DIR/onchain_pccs.json" ]; then
        echo "{}" > "$DCAP_DEPLOY_DIR/onchain_pccs.json"
    fi
    if [ ! -f "$DCAP_DEPLOY_DIR/dcap.json" ]; then
        echo "{}" > "$DCAP_DEPLOY_DIR/dcap.json"
    fi
}

# Helper: read address from PCCS deployment JSON
read_pccs_addr() {
    local key="$1"
    jq -r ".[\"$key\"]" "${PCCS_REPO}/deployment/${CHAIN_ID}.json"
}

# Helper: read address from DCAP deployment JSON
read_dcap_addr() {
    local key="$1"
    jq -r ".[\"$key\"]" "$DCAP_DEPLOY_DIR/dcap.json"
}

# Helper: write to combined output
write_output() {
    local key="$1"
    local value="$2"
    if [ -f "$DEPLOY_OUTPUT" ]; then
        local tmp=$(mktemp)
        jq --arg k "$key" --arg v "$value" '. + {($k): $v}' "$DEPLOY_OUTPUT" > "$tmp"
        mv "$tmp" "$DEPLOY_OUTPUT"
    else
        echo "{\"$key\": \"$value\"}" > "$DEPLOY_OUTPUT"
    fi
}

# ── Step 1: Deploy PCCS Helpers ───────────────────────────────────────────────
deploy_helpers() {
    step "Step 1: Deploy PCCS Helper Contracts"
    cd "$PCCS_REPO"

    OWNER="$OWNER" forge script script/helper/DeployHelpers.s.sol:DeployHelpers \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --broadcast --skip-simulation \
        --legacy \
        -vv

    info "PCCS helpers deployed. Addresses saved to deployment/${CHAIN_ID}.json"

    # Copy to DCAP's expected location
    cp "${PCCS_REPO}/deployment/${CHAIN_ID}.json" "$DCAP_DEPLOY_DIR/onchain_pccs.json"

    for name in EnclaveIdentityHelper FmspcTcbHelper PCKHelper X509CRLHelper TcbEvalHelper; do
        addr=$(read_pccs_addr "$name")
        if [ "$addr" != "null" ] && [ -n "$addr" ]; then
            write_output "$name" "$addr"
            info "  $name: $addr"
        fi
    done
}

# ── Step 2: Deploy PCCS DAOs ─────────────────────────────────────────────────
deploy_dao() {
    step "Step 2: Deploy PCCS DAO Contracts"
    cd "$PCCS_REPO"

    OWNER="$OWNER" forge script script/automata/DeployAutomataDao.s.sol:DeployAutomataDao \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --broadcast --skip-simulation \
        --legacy \
        -vv

    info "PCCS DAOs deployed."

    # Update DCAP's copy
    cp "${PCCS_REPO}/deployment/${CHAIN_ID}.json" "$DCAP_DEPLOY_DIR/onchain_pccs.json"

    for name in AutomataDaoStorage AutomataPcsDao AutomataPckDao; do
        addr=$(read_pccs_addr "$name")
        if [ "$addr" != "null" ] && [ -n "$addr" ]; then
            write_output "$name" "$addr"
            info "  $name: $addr"
        fi
    done
}

# ── Step 3: Deploy Versioned DAOs + TcbEvalDao ──────────────────────────────
deploy_versioned() {
    step "Step 3: Deploy TcbEvalDao + Versioned DAOs (tcbEval=$TCB_EVAL_NUMBER)"
    cd "$PCCS_REPO"

    # 3a: Deploy TcbEvalDao
    info "Deploying AutomataTcbEvalDao..."
    OWNER="$OWNER" forge script \
        script/automata/versioned/DeployAutomataVersioned.s.sol:DeployAutomataVersioned \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --broadcast --skip-simulation \
        --legacy \
        --sig "deployTcbEvalDao()" \
        -vv

    cp "${PCCS_REPO}/deployment/${CHAIN_ID}.json" "$DCAP_DEPLOY_DIR/onchain_pccs.json"
    addr=$(read_pccs_addr "AutomataTcbEvalDao")
    write_output "AutomataTcbEvalDao" "$addr"
    info "  AutomataTcbEvalDao: $addr"

    # 3b: Deploy EnclaveIdentityDaoVersioned
    info "Deploying AutomataEnclaveIdentityDaoVersioned (tcbEval=$TCB_EVAL_NUMBER)..."
    OWNER="$OWNER" forge script \
        script/automata/versioned/DeployAutomataVersioned.s.sol:DeployAutomataVersioned \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --broadcast --skip-simulation \
        --legacy \
        --sig "deployEnclaveIdDaoVersioned(uint32)" "$TCB_EVAL_NUMBER" \
        -vv

    cp "${PCCS_REPO}/deployment/${CHAIN_ID}.json" "$DCAP_DEPLOY_DIR/onchain_pccs.json"
    addr=$(read_pccs_addr "AutomataEnclaveIdentityDaoVersioned_tcbeval_${TCB_EVAL_NUMBER}")
    write_output "AutomataEnclaveIdentityDaoVersioned_tcbeval_${TCB_EVAL_NUMBER}" "$addr"
    info "  EnclaveIdentityDaoVersioned (tcbeval_${TCB_EVAL_NUMBER}): $addr"

    # 3c: Deploy FmspcTcbDaoVersioned
    info "Deploying AutomataFmspcTcbDaoVersioned (tcbEval=$TCB_EVAL_NUMBER)..."
    OWNER="$OWNER" forge script \
        script/automata/versioned/DeployAutomataVersioned.s.sol:DeployAutomataVersioned \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --broadcast --skip-simulation \
        --legacy \
        --sig "deployFmspcTcbDaoVersioned(uint32)" "$TCB_EVAL_NUMBER" \
        -vv

    cp "${PCCS_REPO}/deployment/${CHAIN_ID}.json" "$DCAP_DEPLOY_DIR/onchain_pccs.json"
    addr=$(read_pccs_addr "AutomataFmspcTcbDaoVersioned_tcbeval_${TCB_EVAL_NUMBER}")
    write_output "AutomataFmspcTcbDaoVersioned_tcbeval_${TCB_EVAL_NUMBER}" "$addr"
    info "  FmspcTcbDaoVersioned (tcbeval_${TCB_EVAL_NUMBER}): $addr"
}

# ── Step 4: Deploy PCCSRouter ────────────────────────────────────────────────
deploy_router() {
    step "Step 4: Deploy PCCSRouter"
    cd "$DCAP_EVM"

    OWNER="$OWNER" forge script forge-script/DeployRouter.s.sol:DeployRouter \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --broadcast --skip-simulation \
        --legacy \
        -vv

    local router_addr=$(read_dcap_addr "PCCSRouter")
    write_output "PCCSRouter" "$router_addr"
    info "  PCCSRouter: $router_addr"

    # Grant PCCSRouter access to AutomataDaoStorage
    local storage_addr=$(read_pccs_addr "AutomataDaoStorage")
    info "Granting PCCSRouter access to DaoStorage..."
    cast send "$storage_addr" \
        "setCallerAuthorization(address,bool)" \
        "$router_addr" true \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY"
    info "  PCCSRouter authorized on DaoStorage"

    # Register versioned DAOs on the router
    local enclave_addr=$(read_pccs_addr "AutomataEnclaveIdentityDaoVersioned_tcbeval_${TCB_EVAL_NUMBER}")
    local fmspc_addr=$(read_pccs_addr "AutomataFmspcTcbDaoVersioned_tcbeval_${TCB_EVAL_NUMBER}")

    if [ "$enclave_addr" != "null" ] && [ -n "$enclave_addr" ]; then
        info "Registering versioned EnclaveIdentityDao on router (tcbEval=$TCB_EVAL_NUMBER)..."
        cast send "$router_addr" \
            "setQeIdDaoVersionedAddr(uint32,address)" \
            "$TCB_EVAL_NUMBER" "$enclave_addr" \
            --rpc-url "$DEVNET_RPC_URL" \
            --private-key "$DEPLOYER_PRIVATE_KEY"
    fi

    if [ "$fmspc_addr" != "null" ] && [ -n "$fmspc_addr" ]; then
        info "Registering versioned FmspcTcbDao on router (tcbEval=$TCB_EVAL_NUMBER)..."
        cast send "$router_addr" \
            "setFmspcTcbDaoVersionedAddr(uint32,address)" \
            "$TCB_EVAL_NUMBER" "$fmspc_addr" \
            --rpc-url "$DEVNET_RPC_URL" \
            --private-key "$DEPLOYER_PRIVATE_KEY"
    fi

    info "Router configuration complete"
}

# ── Step 5: Deploy AutomataDcapAttestationFee ────────────────────────────────
deploy_attestation() {
    step "Step 5: Deploy AutomataDcapAttestationFee"
    cd "$DCAP_EVM"

    OWNER="$OWNER" forge script forge-script/AttestationScript.s.sol:AttestationScript \
        --sig "deployEntrypoint()" \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --broadcast --skip-simulation \
        --legacy \
        -vv

    local attest_addr=$(read_dcap_addr "AutomataDcapAttestationFee")
    write_output "AutomataDcapAttestationFee" "$attest_addr"
    info "  AutomataDcapAttestationFee: $attest_addr"
}

# ── Step 6: Deploy V3QuoteVerifier ───────────────────────────────────────────
deploy_verifier() {
    step "Step 6: Deploy V3QuoteVerifier"
    cd "$DCAP_EVM"

    OWNER="$OWNER" QUOTE_VERIFIER_VERSION=3 forge script \
        forge-script/DeployVerifier.s.sol:DeployVerifier \
        --rpc-url "$DEVNET_RPC_URL" \
        --private-key "$DEPLOYER_PRIVATE_KEY" \
        --broadcast --skip-simulation \
        --legacy \
        -vv

    local verifier_addr=$(read_dcap_addr "V3QuoteVerifier")
    write_output "V3QuoteVerifier" "$verifier_addr"
    info "  V3QuoteVerifier: $verifier_addr"
    info "  (auto-registered on AttestationFee + authorized on Router by deploy script)"
}

# ── Step 7: Print summary ───────────────────────────────────────────────────
print_summary() {
    step "Deployment Summary"

    if [ -f "$DEPLOY_OUTPUT" ]; then
        echo ""
        cat "$DEPLOY_OUTPUT" | jq .
        echo ""
    fi

    local attest_addr=$(read_dcap_addr "AutomataDcapAttestationFee" 2>/dev/null || echo "NOT DEPLOYED")

    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}Next steps:${NC}"
    echo ""
    echo "  1. Upload Intel collateral (certificates, CRL, TCB info, QE identity):"
    echo "     ./upload-collateral.sh"
    echo ""
    echo "  2. Update SGXValidationHook to point to the deployed attestation contract:"
    echo "     cast send <SGX_VALIDATION_HOOK_PROXY> \\"
    echo "       \"setAutomataValidationAddr(address)\" \\"
    echo "       $attest_addr \\"
    echo "       --rpc-url $DEVNET_RPC_URL \\"
    echo "       --private-key <OWNER_KEY>"
    echo ""
    echo "  3. Set tcbEvaluationDataNumber on SGXValidationHook:"
    echo "     cast send <SGX_VALIDATION_HOOK_PROXY> \\"
    echo "       \"setTcbEvaluationDataNumber(uint32)\" \\"
    echo "       $TCB_EVAL_NUMBER \\"
    echo "       --rpc-url $DEVNET_RPC_URL \\"
    echo "       --private-key <OWNER_KEY>"
    echo ""
    echo "  4. Run DKG ceremony to test on-chain attestation"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# ── Main ─────────────────────────────────────────────────────────────────────
STEP="${1:-all}"

preflight

case "$STEP" in
    all)
        deploy_helpers
        deploy_dao
        deploy_versioned
        deploy_router
        deploy_attestation
        deploy_verifier
        print_summary
        ;;
    1-helpers)     deploy_helpers ;;
    2-dao)         deploy_dao ;;
    3-versioned)   deploy_versioned ;;
    4-router)      deploy_router ;;
    5-attestation) deploy_attestation ;;
    6-verifier)    deploy_verifier ;;
    7-summary)     print_summary ;;
    *)
        err "Unknown step: $STEP"
        echo "Valid steps: all, 1-helpers, 2-dao, 3-versioned, 4-router, 5-attestation, 6-verifier, 7-summary"
        exit 1
        ;;
esac
