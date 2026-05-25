# DCAP On-Chain Remote Attestation — Devnet Deployment Guide

Deploy the full Automata DCAP attestation stack to Story devnet for testing
real on-chain SGX remote attestation with the DKG protocol.

## Architecture

```
DKG.sol
  └─ SGXValidationHook.validateReport()
       └─ AutomataDcapAttestationFee.verifyAndAttestOnChain(rawQuote, tcbEvalNumber)
            └─ V3QuoteVerifier.verifyQuote()
                 ├─ Parse quote header + auth data
                 ├─ Verify QE Report Data (sha256)
                 ├─ Verify QE Identity via PCCSRouter → EnclaveIdentityDao
                 ├─ Verify X.509 cert chain (PCK → Intermediate → Root CA)
                 ├─ Verify CRL (not revoked)
                 ├─ Verify P256 ECDSA signatures
                 └─ Match TCB level via PCCSRouter → FmspcTcbDao
```

## Prerequisites

- `forge` / `cast` (foundry)
- `openssl`, `curl`, `jq`, `python3`
- Repos cloned:
  - `automata-on-chain-pccs` at `$PCCS_REPO` (default: `../../automata-on-chain-pccs`)
  - `automata-dcap-attestation` at `$DCAP_REPO` (default: `../../automata-dcap-attestation`)
- FMSPC of devnet SGX machines (see below)

## Finding Your Machine's FMSPC

Run on one of the devnet SGX machines:

```bash
# Option 1: via pccs_retrieval_tool
sudo pccs_retrieval_tool -f /tmp/pck_cert.pem
openssl x509 -in /tmp/pck_cert.pem -noout -text | grep -A1 FMSPC

# Option 2: via /dev/attestation (Gramine)
# The FMSPC is embedded in the PCK certificate extension
```

## Deployment Steps

### Step 1: Deploy contracts

```bash
export DEVNET_RPC_URL="http://<devnet-geth-rpc>:8545"
export DEPLOYER_PRIVATE_KEY="0x..."    # account with ETH on devnet
export TCB_EVAL_NUMBER=21              # latest TCB eval number

./deploy-dcap-devnet.sh
```

This deploys in order:
1. PCCS Helper contracts (5 contracts)
2. PCCS DAO contracts (DaoStorage + PcsDao + PckDao)
3. Versioned DAOs (TcbEvalDao + EnclaveIdentityDao + FmspcTcbDao)
4. PCCSRouter (connects to all DAOs)
5. AutomataDcapAttestationFee (main entrypoint)
6. V3QuoteVerifier (registered on AttestationFee)

Output: `deployment-addresses.json`

### Step 2: Upload Intel collateral

```bash
export DEVNET_RPC_URL="http://<devnet-geth-rpc>:8545"
export DEPLOYER_PRIVATE_KEY="0x..."
export FMSPC="00906ED50000"            # your machine's FMSPC

./upload-collateral.sh
```

This fetches from Intel PCS API and uploads:
- Root CA / Platform CA / Signing CA certificates
- Root CA CRL + PCK CRL
- QE Identity
- FMSPC TCB Info

### Step 3: Connect SGXValidationHook

```bash
export DEVNET_RPC_URL="http://<devnet-geth-rpc>:8545"
export SGX_HOOK_OWNER_KEY="0x..."              # SGXValidationHook owner (TEST_DKG_OWNER on devnet)
export SGX_VALIDATION_HOOK_PROXY="0x..."       # proxy address from genesis
export TCB_EVAL_NUMBER=0                       # 0 = accept any TCB level (for testing)

./connect-sgx-hook.sh
```

### Step 4: Test

Run DKG ceremony on devnet. The nodes will submit real SGX quotes that get
verified on-chain through the full DCAP verification pipeline.

## Alternative: Use automata-dcap CLI

The CLI can also manage collateral:

```bash
AUTOMATA_DCAP_CLI="/path/to/automata-dcap-attestation/rust-crates/target/release/automata-dcap"

# Check collateral status
$AUTOMATA_DCAP_CLI --rpc-url $DEVNET_RPC_URL qpl status

# Check missing collateral for a specific quote
$AUTOMATA_DCAP_CLI --rpc-url $DEVNET_RPC_URL --quote-path /path/to/quote.bin qpl check

# Upload collateral from Intel PCS API
$AUTOMATA_DCAP_CLI --rpc-url $DEVNET_RPC_URL --private-key $DEPLOYER_PRIVATE_KEY \
    qpl function sgx_ql_get_quote_verification_collateral \
    --source all --fmspc $FMSPC --pck-ca platform --collateral-version v3

# Verify a quote on-chain
$AUTOMATA_DCAP_CLI --rpc-url $DEVNET_RPC_URL --quote-path /path/to/quote.bin verify onchain
```

## Deployed Contract Addresses (Live Networks)

Story Mainnet (1514) and Aeneid testnet (1315) already have the full stack
deployed by Automata. See `automata-on-chain-pccs/deployment/1514.json`.

## TCB Evaluation Data Number

- `0` in SGXValidationHook = accept any TCB level (good for testing)
- Specific number (e.g., 21) = only accept collateral from that TCB recovery event
- Intel increments this number with each TCB Recovery Event
- On-chain versioned DAOs are per-tcbEvalNumber (separate contract per number)

## Troubleshooting

### "Attestation failed" in DKG
1. Check collateral is uploaded: `automata-dcap ... qpl check --quote-path <quote>`
2. Check FMSPC matches: the uploaded TCB info FMSPC must match the quote's FMSPC
3. Check certificates aren't expired: Intel certs have validity periods
4. Check tcbEvaluationDataNumber: if set to specific number, versioned DAO must exist

### Gas costs
- On-chain quote verification: ~4-5M gas
- Each collateral upload: ~3-15M gas (TCB info is the largest)
