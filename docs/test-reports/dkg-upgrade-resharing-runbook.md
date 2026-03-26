# DKG Upgrade Resharing — Reproducible Runbook

This runbook walks through testing the DKG kernel upgrade resharing flow on a devnet.

---

## Prerequisites

- 3+ Azure VMs with SGX support (validators)
- 1 bootnode VM (no SGX required)
- SSH access to all nodes
- Foundry (`cast`) installed on at least one node
- Two kernel binaries with different MRENCLAVE values (old and new)

## Step 1: Deploy Devnet

### 1.1 Build and Start story-geth

On each node:
```bash
cd ~/story-geth
git pull origin <branch>
make geth
sudo systemctl restart story-geth
```

### 1.2 Build and Start story

On each node:
```bash
cd ~/story
git pull origin <branch>
make build
sudo systemctl restart story
```

### 1.3 Build and Start Old Kernel

On each validator:
```bash
cd ~/story-kernel
git checkout <old-kernel-branch>
make build
# Start on port 50051
./story-kernel start --home ~/.story-kernel --port 50051
```

### 1.4 Build and Start New Kernel

On each validator:
```bash
cd ~/story-kernel
git checkout <new-kernel-branch>
make build-new  # Or use a separate directory
# Start on port 50052
./story-kernel start --home ~/.story-kernel-new --port 50052
```

### 1.5 Verify Both Kernels Running

```bash
# Check old kernel
curl -s http://localhost:50051/health

# Check new kernel
curl -s http://localhost:50052/health
```

Note both kernels' MRENCLAVE values:
- Old: appears in story logs as `code_commitment` during normal DKG rounds
- New: appears when new kernel connects (or check build output)

## Step 2: Wait for Normal DKG Round to Complete

Wait until at least one normal DKG round completes successfully:
```bash
journalctl -u story --no-pager | grep "DKG process completed successfully"
```

Verify the active round:
```bash
journalctl -u story --no-pager | grep "Emitted DKGFinalized"
```

Record the GlobalPubKey from the finalized round — it should be preserved after resharing.

## Step 3: Whitelist New Enclave Type

**IMPORTANT**: This must be done BEFORE scheduling the upgrade.

### 3.1 Identify the Validation Hook Address

```bash
# Query enclave type 1 to get the existing validation hook address
PATH=$PATH:$HOME/.foundry/bin
cast call 0xCcCcCC0000000000000000000000000000000004 \
  "enclaveTypeData(bytes32)" \
  $(cast --to-bytes32 1) \
  --rpc-url http://localhost:8545
```

The validation hook address is in the second 32-byte slot of the response.

### 3.2 Whitelist New Enclave Type

```bash
NEW_MRENCLAVE=<new kernel MRENCLAVE hex, no 0x prefix>
HOOK_ADDR=<validation hook address from step 3.1>
DKG_OWNER_PRIVKEY=<DKG contract owner private key>

PATH=$PATH:$HOME/.foundry/bin
cast send 0xCcCcCC0000000000000000000000000000000004 \
  "whitelistEnclaveType(bytes32,(bytes32,address),bool)" \
  $(cast --to-bytes32 2) \
  "(0x${NEW_MRENCLAVE},${HOOK_ADDR})" \
  true \
  --private-key $DKG_OWNER_PRIVKEY \
  --rpc-url http://localhost:8545 \
  --gas-price 17000000000 \
  --legacy
```

### 3.3 Verify Whitelisting

```bash
cast call 0xCcCcCC0000000000000000000000000000000004 \
  "isEnclaveTypeWhitelisted(bytes32)" \
  $(cast --to-bytes32 2) \
  --rpc-url http://localhost:8545
```

Expected: `0x...01` (true)

## Step 4: Schedule Upgrade

### 4.1 Determine Activation Height

Check current block height:
```bash
journalctl -u story --no-pager | tail -5 | grep "height="
```

Choose an activation height at least 50-100 blocks in the future to allow time for the transaction to be processed.

### 4.2 Call scheduleUpgrade

```bash
ACTIVATION_HEIGHT=<chosen height>

PATH=$PATH:$HOME/.foundry/bin
cast send 0xCcCcCC0000000000000000000000000000000004 \
  "scheduleUpgrade(uint256,string)" \
  $ACTIVATION_HEIGHT \
  "v3.0.0" \
  --private-key $DKG_OWNER_PRIVKEY \
  --rpc-url http://localhost:8545 \
  --gas-price 17000000000 \
  --legacy
```

### 4.3 Verify Upgrade Scheduled

```bash
journalctl -u story --no-pager | grep "UpgradeScheduled"
```

Expected:
```
Upgrade scheduled  activation_height=<ACTIVATION_HEIGHT> upgrade_version=v3.0.0
```

## Step 5: Monitor Upgrade Resharing

### 5.1 Wait for Activation Height

```bash
# Watch for the upgrade activation
journalctl -u story -f | grep -E "Kernel upgrade activated|initiating upgrade resharing"
```

Expected:
```
Kernel upgrade activated, initiating upgrade resharing round  upgrade_version=v3.0.0 activation_height=<HEIGHT>
```

### 5.2 Monitor Registration Phase

```bash
journalctl -u story -f | grep -E "Handling DKG registration|Register succeeded|Resolved enclave type"
```

Key verification:
- `Resolved enclave type for upgrade round  enclave_type=...0002` — must be type 2 (new kernel)
- `Register succeeded` — registration tx mined successfully

### 5.3 Monitor Dealing Phase

```bash
journalctl -u story -f | grep -E "process deals|process responses"
```

### 5.4 Monitor Finalization Phase

```bash
journalctl -u story -f | grep -E "Finalize|finalized|DKGFinalized"
```

Key verification:
- `Calling finalize contract method  enclave_type=...0002` — finalization uses type 2
- `DKG successfully finalized  code_commitment=<NEW_MRENCLAVE>` — all validators finalized with new CC
- `Emitted DKGFinalized event  round=<N>` — on-chain finalization event

### 5.5 Verify Completion

```bash
journalctl -u story --no-pager | grep -E "Upgrade resharing round completed|DKG process completed"
```

Expected:
```
Upgrade resharing round completed, new TEE binary is now active  round=<N>
DKG process completed successfully  round=<N>
```

## Step 6: Verify Results

### 6.1 GlobalPubKey Preservation

Compare GlobalPubKey from the upgrade round with the previous round:
```bash
journalctl -u story --no-pager | grep "global_pub_key=" | tail -5
```

The GlobalPubKey should be **identical** before and after resharing.

### 6.2 Committee Rewards

```bash
journalctl -u story --no-pager | grep "DKGCommitteeRewarded.*round=<N>"
```

### 6.3 Chain Health

Verify the chain is producing blocks normally:
```bash
journalctl -u story --no-pager | tail -5
```

No consensus failures or app hash mismatches should occur.

### 6.4 All Validators Consistent

Check that all 3 validators show the same completion logs:
```bash
for IP in <VAL1_IP> <VAL2_IP> <VAL3_IP>; do
  echo "=== $IP ==="
  ssh -i <KEY> ubuntu@$IP 'journalctl -u story --no-pager | grep "DKG process completed.*round=<N>"'
done
```

## Troubleshooting

### External TX Not Being Mined

Story-geth's miner has a minimum gas price (typically 16 Gwei). Always use:
```bash
--gas-price 17000000000 --legacy
```

### ResolveEnclaveType Failure

If you see `"no whitelisted enclave type found for code commitment"`:
1. Verify the new MRENCLAVE is whitelisted (Step 3.3)
2. Ensure the MRENCLAVE hex matches exactly (compare with kernel build output)
3. If missed, whitelist it now — the next DKG round will use it

### DKG Service Already Running

If you see `"DKG service already running; skipping registration"`:
- This is a race condition from the previous round's goroutine
- Self-resolves: the next round will start normally

### Finalization Signature Mismatch

If a validator shows `"finalization signature address mismatch"`:
- Lower `minReqFinalizedParticipants` if needed:
```bash
cast send 0xCcCcCC0000000000000000000000000000000004 \
  "setMinReqFinalizedParticipants(uint256)" \
  2 \
  --private-key $DKG_OWNER_PRIVKEY \
  --rpc-url http://localhost:8545 \
  --gas-price 17000000000 --legacy
```

### Validator Recovery

If a validator hits app hash mismatch:
```bash
# Roll back 1 block (do NOT wipe state)
./story rollback
sudo systemctl restart story
```

## DKG Stage Timing Reference

With default devnet params (registration=30, dealing=100, finalization=100, active=600):

| Stage | Duration | Cumulative |
|-------|----------|------------|
| Registration | 30 blocks (~1.5 min) | 0-30 |
| Dealing | 100 blocks (~5 min) | 30-130 |
| Finalization | 100 blocks (~5 min) | 130-230 |
| Active | 600 blocks (~30 min) | 230-830 |

Total round duration: ~830 blocks (~42 minutes at 3s block time).
