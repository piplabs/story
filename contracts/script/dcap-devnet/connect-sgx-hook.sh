#!/usr/bin/env bash
# connect-sgx-hook.sh — Connect SGXValidationHook to deployed DCAP attestation contracts
#
# This script updates the SGXValidationHook proxy on devnet to point to the real
# AutomataDcapAttestationFee contract instead of the mock.
#
# Prerequisites:
#   - deploy-dcap-devnet.sh completed
#   - upload-collateral.sh completed
#   - SGXValidationHook proxy already deployed (via GenerateAlloc or DeploySGXValidationHook)
#
# Required env vars:
#   DEVNET_RPC_URL                — devnet EL RPC URL
#   SGX_HOOK_OWNER_KEY            — private key of SGXValidationHook owner
#   SGX_VALIDATION_HOOK_PROXY     — address of SGXValidationHook proxy
#
# Optional env vars:
#   TCB_EVAL_NUMBER               — TCB evaluation data number (default: 0, accepts any)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_FILE="${SCRIPT_DIR}/deployment-addresses.json"

TCB_EVAL_NUMBER="${TCB_EVAL_NUMBER:-0}"

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
err()   { echo -e "${RED}[ERROR]${NC} $1"; }
step()  { echo -e "\n${CYAN}═══ $1 ═══${NC}"; }

# ── Preflight ────────────────────────────────────────────────────────────────
for var in DEVNET_RPC_URL SGX_HOOK_OWNER_KEY SGX_VALIDATION_HOOK_PROXY; do
    if [ -z "${!var:-}" ]; then
        err "$var is not set"
        exit 1
    fi
done

if [ ! -f "$DEPLOY_FILE" ]; then
    err "deployment-addresses.json not found. Run deploy-dcap-devnet.sh first."
    exit 1
fi

ATTEST_ADDR=$(jq -r '.AutomataDcapAttestationFee' "$DEPLOY_FILE")
if [ "$ATTEST_ADDR" = "null" ] || [ -z "$ATTEST_ADDR" ]; then
    err "AutomataDcapAttestationFee not found in deployment-addresses.json"
    exit 1
fi

# ── Step 1: Read current SGXValidationHook state ────────────────────────────
step "Current SGXValidationHook State"

CURRENT_AUTOMATA=$(cast call "$SGX_VALIDATION_HOOK_PROXY" \
    "automataValidationAddr()(address)" \
    --rpc-url "$DEVNET_RPC_URL" 2>/dev/null || echo "FAILED")

CURRENT_TCB=$(cast call "$SGX_VALIDATION_HOOK_PROXY" \
    "tcbEvaluationDataNumber()(uint32)" \
    --rpc-url "$DEVNET_RPC_URL" 2>/dev/null || echo "FAILED")

info "Hook proxy: $SGX_VALIDATION_HOOK_PROXY"
info "Current automataValidationAddr: $CURRENT_AUTOMATA"
info "Current tcbEvaluationDataNumber: $CURRENT_TCB"
info "New automataValidationAddr: $ATTEST_ADDR"
info "New tcbEvaluationDataNumber: $TCB_EVAL_NUMBER"

# ── Step 2: Update automataValidationAddr ────────────────────────────────────
step "Updating automataValidationAddr"

cast send "$SGX_VALIDATION_HOOK_PROXY" \
    "setAutomataValidationAddr(address)" \
    "$ATTEST_ADDR" \
    --rpc-url "$DEVNET_RPC_URL" \
    --private-key "$SGX_HOOK_OWNER_KEY"

info "automataValidationAddr updated to $ATTEST_ADDR"

# ── Step 3: Update tcbEvaluationDataNumber ───────────────────────────────────
step "Updating tcbEvaluationDataNumber"

cast send "$SGX_VALIDATION_HOOK_PROXY" \
    "setTcbEvaluationDataNumber(uint32)" \
    "$TCB_EVAL_NUMBER" \
    --rpc-url "$DEVNET_RPC_URL" \
    --private-key "$SGX_HOOK_OWNER_KEY"

info "tcbEvaluationDataNumber updated to $TCB_EVAL_NUMBER"

# ── Step 4: Verify ───────────────────────────────────────────────────────────
step "Verification"

NEW_AUTOMATA=$(cast call "$SGX_VALIDATION_HOOK_PROXY" \
    "automataValidationAddr()(address)" \
    --rpc-url "$DEVNET_RPC_URL")

NEW_TCB=$(cast call "$SGX_VALIDATION_HOOK_PROXY" \
    "tcbEvaluationDataNumber()(uint32)" \
    --rpc-url "$DEVNET_RPC_URL")

info "Verified automataValidationAddr: $NEW_AUTOMATA"
info "Verified tcbEvaluationDataNumber: $NEW_TCB"

if [ "$NEW_AUTOMATA" = "$ATTEST_ADDR" ]; then
    echo -e "\n${GREEN}SGXValidationHook is now connected to real DCAP attestation!${NC}"
    echo ""
    echo "  Next: Run DKG ceremony on devnet to test on-chain remote attestation"
    echo "  The DKG nodes will submit SGX quotes that get verified on-chain via:"
    echo "    DKG → SGXValidationHook → AutomataDcapAttestationFee → V3QuoteVerifier → PCCS"
else
    err "Verification failed! automataValidationAddr mismatch"
    exit 1
fi
