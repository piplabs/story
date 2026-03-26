# DKG E2E Test Workflow by @Hansol Lee

# DKG End-to-End Test Workflow

Step-by-step instructions for running DKG happy path and unhappy path e2e tests on a 3-validator devnet.
Each test includes commands, verification steps, and expected outputs for reliable reproduction.

---

## 1. Overview

### 1.1 Scope

- **Happy path**: HP-01 through HP-05
- **Unhappy path**: TC-01 through TC-07
- Each test starts from a clean devnet state (full reset)
- Tests are independent and can be run in any order

### 1.2 Repositories

| Component | Repository | Branch |
| --- | --- | --- |
| story (CL) | `storyprotocol/story` | `dkg/devnet` |
| story-geth (EL) | `storyprotocol/story-geth` | `v1.2.1` |
| story-kernel (TEE) | `storyprotocol/story-kernel` | `main` |

### 1.3 Key Constants

| Parameter | Value | Notes |
| --- | --- | --- |
| DKG Contract | `0xCcCcCC0000000000000000000000000000000004` | Predeployed |
| OperationalThreshold | 500 (devnet) | 2-of-3 threshold |
| MinReqRegistered | 2 | Min validators to proceed |
| MinReqFinalized | 2 | Min validators to finalize |
| RegistrationPeriod | 200 blocks (~6.5 min) |  |
| DealingPeriod | 300 blocks (~10 min) |  |
| FinalizationPeriod | 300 blocks (~10 min) |  |
| ActivePeriod | 600 blocks (~20 min) |  |
| V2.0.0 Upgrade Height | 110 (devnet) | DKG activates here |

### 1.4 Block Timing Reference

For a round starting at block S (startBlockHeight):
- Registration: blocks S to S+199
- Dealing: blocks S+200 to S+499
- Finalization: blocks S+500 to S+799
- Active: blocks S+800 to S+1399

### 1.5 Devnet Setup and Reset

All infrastructure setup, binary build/deploy, devnet reset procedures, and service start/stop instructions are in **DEVNET_RUNBOOK.md** (`.claude/DEVNET_RUNBOOK.md`).

Every test below assumes:
1. Devnet has been reset per the runbook (Part 1, Steps 1-8)
2. All services are running (story, node-geth, story-kernel on validators)
3. Chain is producing blocks and kernel connections are verified

### 1.6 Environment Variables

```bash
export VAL1_IP="20.46.165.193"
export VAL2_IP="20.48.25.208"
export VAL3_IP="40.115.139.113"
export BOOTNODE_IP="23.102.71.16"
export ALL_VALS="$VAL1_IP$VAL2_IP$VAL3_IP"
export ALL_NODES="$BOOTNODE_IP$VAL1_IP$VAL2_IP$VAL3_IP"
export SSH_USER="ubuntu"
```

---

## 2. Happy Path Tests

### HP-01: DKG Initialization (Round 1)

**Objective**: Verify that a fresh DKG round completes successfully with all 3 validators.

**Prerequisites**: Clean devnet, all services running, block height < 110.

**Steps**:

```bash
# Step 1: Wait for V2.0.0 upgrade activation (block 110)
while true; do
  HEIGHT=$(ssh $SSH_USER@$VAL1_IP "curl -s localhost:26657/status | python3 -c\"import sys,json; print(json.load(sys.stdin)['result']['sync_info']['latest_block_height'])\"")
  echo "Current height:$HEIGHT"
  if [ "$HEIGHT" -ge 110 ]; then break; fi
  sleep 5
done

# Step 2: Verify DKG round initiated
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '5 minutes ago' | grep 'Initiated new DKG round'"

# Step 3: Wait for all 3 registrations
for IP in $ALL_VALS; do
  echo "===$IP ==="
  timeout 420 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"DKG initialization complete\"'; do sleep 10; done" && echo "REGISTERED" || echo "TIMEOUT"
done

# Step 4: Wait for dealing phase
for IP in $ALL_VALS; do
  timeout 600 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"15 minutes ago\" | grep -q\"DKG deals are generated successfully\"'; do sleep 10; done" && echo "DEALS GENERATED on$IP" || echo "TIMEOUT on$IP"
done

# Step 5: Wait for finalization
for IP in $ALL_VALS; do
  timeout 600 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"25 minutes ago\" | grep -q\"DKG finalization phase complete\"'; do sleep 10; done" && echo "FINALIZED on$IP" || echo "TIMEOUT on$IP"
done
```

**Verification Checklist**:
- [ ] “Initiated new DKG round” log on at least one validator
- [ ] “DKG initialization complete” log on all 3 validators
- [ ] “DKG deals are generated successfully” on all 3
- [ ] “Process deals complete” on all 3
- [ ] “Process responses complete” on all 3
- [ ] “DKG finalization phase complete” on all 3
- [ ] On-chain DKG network has GlobalPublicKey set
- [ ] Stage transitioned to Active

**Success Criteria**: All 3 validators registered, dealt, and finalized. GlobalPublicKey set on-chain.

---

### HP-02: DKG Resharing (Round 2 — Natural)

**Objective**: Verify that when the Active period ends, a new resharing round starts and completes.

**Prerequisites**: HP-01 completed. Round 1 in Active stage.

**Steps**:

```bash
# Step 1: Get the current round's start block and calculate when Active ends
# Active stage ends at: startBlock + 200 + 300 + 300 + 600 = startBlock + 1400
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '1 hour ago' | grep 'Initiated new DKG round' | tail -1"

# Step 2: Wait for new round initiation (round 2)
timeout 2400 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"30 minutes ago\" | grep -q\"Initiated new DKG round.*round.*2\"'; do sleep 15; done" && echo "ROUND 2 STARTED" || echo "TIMEOUT"

# Step 3: Verify resharing flag
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '5 minutes ago' | grep -E 'is_resharing|IsResharing'"

# Step 4: Wait for round 2 completion
for IP in $ALL_VALS; do
  timeout 1200 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"30 minutes ago\" | grep -q\"DKG finalization phase complete.*round.*2\"'; do sleep 15; done" && echo "ROUND 2 FINALIZED on$IP" || echo "TIMEOUT on$IP"
done
```

**Verification Checklist**:
- [ ] Round 2 initiated with IsResharing=true
- [ ] All 3 validators registered, dealt, finalized for round 2
- [ ] Round 1 marked as completed/inactive

**Success Criteria**: Round 2 resharing completes with threshold-met finalization.

---

### HP-03: DKG Resharing with New Validator Member

**Objective**: Verify resharing when the active validator set changes (a new validator joins).

**Prerequisites**: HP-01 completed. A 4th validator machine prepared and bonded.

**Note**: This test requires a pre-configured 4th validator. If not available, skip this test.

**Steps**:

```bash
# Step 1: Bond 4th validator (Val4) to the network via staking transaction
export VAL4_IP="<val4-ip>"
ssh $SSH_USER@$VAL4_IP "story validator create-validator --stake-amount 1000000000000000000 --from validator"

# Step 2: Wait for Val4 to appear in active validator set (may take 1+ epochs)

# Step 3: Wait for next resharing round — ActiveValSet should include Val4

# Step 4: Verify all 4 validators participate
for IP in $ALL_VALS $VAL4_IP; do
  timeout 1200 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"30 minutes ago\" | grep -q\"DKG finalization phase complete\"'; do sleep 15; done" && echo "FINALIZED on$IP" || echo "TIMEOUT on$IP"
done
```

**Verification Checklist**:
- [ ] Round ActiveValSet contains 4 addresses
- [ ] Round Threshold updated for 4 validators: `CalculateThreshold(4, 500) = 3`
- [ ] All 4 validators register, deal, and finalize
- [ ] GlobalPublicKey set

**Success Criteria**: Resharing with expanded committee completes successfully.

---

### HP-04: Kernel Upgrade Resharing

**Objective**: Verify the kernel upgrade workflow: schedule → activate → upgrade resharing round.

**Prerequisites**: HP-01 completed. New kernel binary built and deployed to all validators.

See **DEVNET_RUNBOOK.md Part 2** for the full kernel upgrade procedure (build new binary, start on port 50052, update story config, whitelist enclave type, schedule upgrade).

**Steps**:

```bash
# Step 1: Follow DEVNET_RUNBOOK.md Part 2, Steps 1-5 to set up dual-binary mode and schedule upgrade

# Step 2: Wait for activation height
while true; do
  HEIGHT=$(ssh $SSH_USER@$VAL1_IP "curl -s localhost:26657/status | python3 -c\"import sys,json; print(json.load(sys.stdin)['result']['sync_info']['latest_block_height'])\"")
  echo "Height:$HEIGHT /$ACTIVATION_HEIGHT"
  if [ "$HEIGHT" -ge "$ACTIVATION_HEIGHT" ]; then break; fi
  sleep 5
done

# Step 3: Verify upgrade resharing round initiated
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '5 minutes ago' | grep 'Kernel upgrade activated'"
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '5 minutes ago' | grep 'is_upgrade.*true'"

# Step 4: Wait for upgrade round completion
for IP in $ALL_VALS; do
  timeout 1200 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"30 minutes ago\" | grep -q\"DKG finalization phase complete\"'; do sleep 15; done" && echo "UPGRADE ROUND FINALIZED on$IP" || echo "TIMEOUT on$IP"
done
```

**Verification Checklist**:
- [ ] “Kernel upgrade activated, initiating upgrade resharing round” log
- [ ] Round has IsUpgrade=true
- [ ] Deals generated via old binary (dealer uses OldCodeCommitment)
- [ ] Responses sent to BOTH old and new binaries
- [ ] KernelUpgradeInfo deleted after successful finalization
- [ ] GlobalPublicKey set for new round

**Success Criteria**: Upgrade resharing round completes. Key shares migrated from old to new kernel binary.

---

### HP-05: DKG Rewards from UBI Pool

**Objective**: Verify that committee members receive the correct portion of UBI rewards during ProcessUbiWithdrawal (EndBlock).

**Prerequisites**: HP-01 completed. UBI pool has sufficient balance.

**Reward distribution logic** (two paths):
1. **FinalizeDKGRound** (`settleRewardsForPreviousCommittee`): Withdraws ALL UBI → distributes `portion * withdrawn / memberCount` to each member → stores remainder in SettlementBalance
2. **EndBlock ProcessUbiWithdrawal** (`DistributeRewardsToActiveCommittee`): Each cycle, if UBI balance ≥ `MinPartialWithdrawalAmount`, withdraws UBI → distributes `portion * withdrawn / memberCount` to active committee → burns remainder

**Formula**:
- `dkgReward = DkgCommitteeRewardPortion (default 0.10) * withdrawnAmount`
- `perMemberReward = dkgReward / memberCount` (truncated)
- `totalDistributed = perMemberReward * memberCount`
- `remaining = withdrawnAmount - totalDistributed` → SettlementBalance or burned

**Setup**: Lower `MinPartialWithdrawalAmount` so UBI withdrawals trigger frequently (default 8 IP is too high for devnet).

```bash
# Step 0 (BEFORE devnet reset): Set MinPartialWithdrawalAmount low in genesis
# In genesis.json, under evmstaking.params:
#   "min_partial_withdrawal_amount": "600000"
# This ensures UBI withdrawal triggers every EndBlock cycle when balance is sufficient.

# Step 1: After DKG round completes, get each validator's EVM address
DKG_CONTRACT="0xCcCcCC0000000000000000000000000000000004"
for IP in $ALL_VALS; do
  EVM_ADDR=$(ssh $SSH_USER@$IP "cat ~/.story/story/config/priv_validator_key.json | python3 -c\"
import sys,json,hashlib
from eth_keys import keys
# Extract validator's EVM address from story logs or config
\"")
  echo "$IP:$EVM_ADDR"
done

# Simpler: get EVM addresses from story logs
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Register succeeded' | head -3"

# Step 2: Query validator EVM balances BEFORE reward distribution
# Use cast to query each validator's IP balance
for ADDR in $VAL1_EVM $VAL2_EVM $VAL3_EVM; do
  echo "===$ADDR ==="
  ssh $SSH_USER@$VAL1_IP "cast balance$ADDR --rpc-url http://localhost:8545"
done

# Step 3: Query current UBI pool balance
ssh $SSH_USER@$VAL1_IP "curl -s localhost:1317/cosmos/distribution/v1beta1/ubi_balance"

# Step 4: Wait for UBI withdrawal to trigger (EndBlock)
# With low MinPartialWithdrawalAmount, this should happen within a few blocks
timeout 120 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"2 minutes ago\" | grep -q\"Distributed DKG committee rewards\"'; do sleep 5; done" && echo "REWARDS DISTRIBUTED" || echo "TIMEOUT"

# Step 5: Get the reward event details
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '5 minutes ago' | grep 'Distributed DKG committee rewards'"
# Expected log fields: round=N, member_count=3, total_distributed=X, per_member=Y

# Step 6: Query validator EVM balances AFTER reward distribution
for ADDR in $VAL1_EVM $VAL2_EVM $VAL3_EVM; do
  echo "===$ADDR ==="
  ssh $SSH_USER@$VAL1_IP "cast balance$ADDR --rpc-url http://localhost:8545"
done

# Step 7: Verify expected reward amount
# From the log: per_member = portion * withdrawnAmount / memberCount
# Each validator's balance should have increased by exactly per_member
```

**Verification**:

```bash
# Check reward distribution log exists with correct member count
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Distributed DKG committee rewards.*member_count=3'" && echo "PASS: 3 members rewarded"

# Verify per_member matches expected formula
# Extract values from log
REWARD_LOG=$(ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Distributed DKG committee rewards' | tail -1")
echo "Reward log:$REWARD_LOG"
# Parse: total_distributed, per_member, member_count
# Verify: total_distributed == per_member * member_count

# Check SettlementBalance (remainder from FinalizeDKGRound)
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'settlement'" || echo "No settlement (evenly divisible)"

# Check DKGCommitteeRewarded event emitted
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Emitted DKGCommitteeRewarded event'" && echo "PASS: Event emitted"

# Verify balance increase matches per_member
# BALANCE_AFTER - BALANCE_BEFORE should equal per_member for each validator
```

**Verification Checklist**:
- [ ] `MinPartialWithdrawalAmount` set low enough for frequent UBI withdrawal
- [ ] “Distributed DKG committee rewards” log with `member_count=3`
- [ ] `per_member` = `floor(0.10 * withdrawnAmount / 3)`
- [ ] `total_distributed` = `per_member * 3`
- [ ] Each finalized validator’s balance increased by exactly `per_member`
- [ ] DKGCommitteeRewarded event emitted with correct values
- [ ] SettlementBalance stores remainder if any (withdrawn - totalDistributed)
- [ ] Non-finalized validators received nothing

**Success Criteria**: Rewards distributed correctly — each validator received exactly `per_member` amount, matching the formula `floor(DkgCommitteeRewardPortion * withdrawnAmount / memberCount)`.

---

## 3. Unhappy Path Tests

### TC-01: ResumeDKG — Registration Phase Crash (Kernel Kill)

**Objective**: Verify `ResumeDKGService` correctly resumes when kernel crashes during registration.

**How it works**:
1. DKG Round 1 starts at V2.0.0 upgrade height (block 110)
2. `handleDKGRegistration` calls `GenerateAndSealKey` on kernel via gRPC
3. Kernel killed → gRPC fails after 3 retries (2s each) → `MarkFailed`
4. `BeginBlocker` → `ResumeDKGService` detects PhaseFailed → retries registration
5. Kernel restarted → gRPC reconnects → registration completes

**Steps**:

```bash
# Step 1: Wait for DKG activation (block 110)
while true; do
  HEIGHT=$(ssh $SSH_USER@$VAL1_IP "curl -s localhost:26657/status | python3 -c\"import sys,json; print(json.load(sys.stdin)['result']['sync_info']['latest_block_height'])\"")
  if [ "$HEIGHT" -ge 110 ]; then break; fi
  sleep 3
done

# Step 2: Wait for DKG round to start on Val1
timeout 120 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"2 minutes ago\" | grep -q\"Initiated new DKG round\"'; do sleep 5; done"

# Step 3: Kill kernel on Val1 ONLY (not story, not geth)
ssh $SSH_USER@$VAL1_IP "sudo kill -9\$(sudo lsof -t -i :50051) 2>/dev/null; sudo pkill -9 -x gramine-sgx 2>/dev/null; sudo pkill -9 -x loader 2>/dev/null"

# Step 4: Verify kernel is dead
ssh $SSH_USER@$VAL1_IP "sudo lsof -i :50051 | grep LISTEN || echo 'KERNEL DEAD'"

# Step 5: Watch for failure and retry logs (30 seconds)
sleep 30
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '1 minute ago' | grep -E '(Failed to generate|MarkFailed|Recovering|ResumeDKGService)'"

# Step 6: Restart kernel on Val1
ssh $SSH_USER@$VAL1_IP "cd ~/story-kernel && nohup gramine-sgx story-kernel start --home ~/.story-kernel > /tmp/kernel.log 2>&1 &"
sleep 60  # SGX enclave loading takes 1-3 min

# Step 7: Verify kernel is running
ssh $SSH_USER@$VAL1_IP "sudo lsof -i :50051 | grep LISTEN && echo 'KERNEL RUNNING'"

# Step 8: Restart story to re-establish kernel gRPC connection
ssh $SSH_USER@$VAL1_IP "sudo systemctl restart story"
sleep 5

# Step 9: Wait for recovery
timeout 300 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"5 minutes ago\" | grep -q\"DKG initialization complete\"'; do sleep 10; done" && echo "RECOVERY SUCCESS" || echo "RECOVERY TIMEOUT"

# Step 10: Wait for full round completion
timeout 1800 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"30 minutes ago\" | grep -q\"DKG finalization phase complete\"'; do sleep 15; done" && echo "ROUND COMPLETE" || echo "ROUND TIMEOUT"
```

**Verification**:

```bash
# Check failure log exists
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep -c 'Failed to generate'" | xargs test 0 -lt && echo "PASS: Failure detected"

# Check recovery log exists
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep -c 'DKG initialization complete'" | xargs test 0 -lt && echo "PASS: Recovery confirmed"

# Check round completed
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep -c 'DKG finalization phase complete'" | xargs test 0 -lt && echo "PASS: Round completed"
```

**Success Criteria**:
- Kernel kill triggers `MarkFailed` in session
- BeginBlocker ResumeDKGService retries registration on every block
- Kernel restart + story restart → gRPC reconnect → registration completes
- Full DKG round completes

---

### TC-02: ResumeDKG — Dealing Phase Crash (Kernel Kill)

**Objective**: Verify recovery when kernel crashes during deal generation.

**Steps**:

```bash
# Step 1: Wait for all 3 validators to register
for IP in $ALL_VALS; do
  timeout 420 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"DKG initialization complete\"'; do sleep 10; done" && echo "REGISTERED:$IP" || echo "TIMEOUT:$IP"
done

# Step 2: Watch for Dealing phase start on Val1
timeout 600 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"Handling DKG dealing\"'; do sleep 5; done"

# Step 3: Kill kernel IMMEDIATELY after seeing dealing log
ssh $SSH_USER@$VAL1_IP "sudo kill -9\$(sudo lsof -t -i :50051) 2>/dev/null; sudo pkill -9 -x gramine-sgx 2>/dev/null; sudo pkill -9 -x loader 2>/dev/null"

# Step 4: Watch failure logs (20 seconds)
sleep 20
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '1 minute ago' | grep -E '(Failed to generate deals|MarkFailed)'"

# Step 5: Restart kernel
ssh $SSH_USER@$VAL1_IP "cd ~/story-kernel && nohup gramine-sgx story-kernel start --home ~/.story-kernel > /tmp/kernel.log 2>&1 &"
sleep 60

# Step 6: Restart story to re-establish kernel gRPC connection
ssh $SSH_USER@$VAL1_IP "sudo systemctl restart story"
sleep 5

# Step 7: Wait for deal generation recovery
timeout 600 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"DKG deals are generated successfully\"'; do sleep 10; done" && echo "DEALS RECOVERED" || echo "TIMEOUT"

# Step 8: Wait for round completion
timeout 1800 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"30 minutes ago\" | grep -q\"DKG finalization phase complete\"'; do sleep 15; done" && echo "ROUND COMPLETE" || echo "TIMEOUT"
```

**Verification**:

```bash
# If kernel died before GenerateDeals completed:
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Failed to generate deals'" && echo "PASS: Dealing failure detected"

# If GenerateDeals completed before kill (fast path):
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'DKG deals are generated successfully'" && echo "PASS: Deals generated (possibly before kill)"

# Either way, round should complete (threshold=2)
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'DKG finalization phase complete'" && echo "PASS: Round completed"
```

**Success Criteria**:
- If kernel killed before GenerateDeals: failure detected → resume → deals generated after recovery
- If kernel killed after GenerateDeals: round completes normally
- Round succeeds with threshold=2 (even if Val1 fails permanently)

---

### TC-03: ResumeDKG — Finalization Phase Crash (Kernel Kill)

**Objective**: Verify recovery when kernel crashes during finalization.

**Steps**:

```bash
# Step 1: Wait for registration + dealing complete
for IP in $ALL_VALS; do
  timeout 1200 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"20 minutes ago\" | grep -q\"Process responses complete\"'; do sleep 15; done" && echo "DEALING DONE:$IP" || echo "TIMEOUT:$IP"
done

# Step 2: Watch for Finalization phase
timeout 600 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"15 minutes ago\" | grep -q\"Finalize call to kernel\"'; do sleep 5; done"

# Step 3: Kill kernel immediately
ssh $SSH_USER@$VAL1_IP "sudo kill -9\$(sudo lsof -t -i :50051) 2>/dev/null; sudo pkill -9 -x gramine-sgx 2>/dev/null; sudo pkill -9 -x loader 2>/dev/null"

# Step 4: Wait and check logs
sleep 15
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '1 minute ago' | grep -E '(Failed|MarkFailed|finalization)'"

# Step 5: Restart kernel
ssh $SSH_USER@$VAL1_IP "cd ~/story-kernel && nohup gramine-sgx story-kernel start --home ~/.story-kernel > /tmp/kernel.log 2>&1 &"
sleep 60

# Step 6: Restart story to re-establish kernel gRPC connection
ssh $SSH_USER@$VAL1_IP "sudo systemctl restart story"
sleep 5

# Step 7: Wait for finalization completion
timeout 600 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"15 minutes ago\" | grep -q\"DKG finalization phase complete\"'; do sleep 10; done" && echo "FINALIZATION RECOVERED" || echo "TIMEOUT"
```

**Verification**:

```bash
# Check finalization completed (Val1 or threshold met by Val2+Val3)
FINALIZED_COUNT=0
for IP in $ALL_VALS; do
  if ssh $SSH_USER@$IP "journalctl -u story --since '30 minutes ago' | grep -q 'DKG finalization phase complete'"; then
    FINALIZED_COUNT=$((FINALIZED_COUNT + 1))
  fi
done
echo "Finalized validators:$FINALIZED_COUNT"
test "$FINALIZED_COUNT" -ge 2 && echo "PASS: Threshold met" || echo "FAIL: Insufficient finalizations"
```

**Success Criteria**:
- Finalization completes on Val1 after recovery, OR
- Threshold=2 met by Val2+Val3 even if Val1 fails
- Round transitions to Active

---

### TC-04: Non-Failed Phase Recovery After Process Crash

**Objective**: Verify recovery when the story process crashes (SIGKILL) before `MarkFailed` runs, leaving session in intermediate phase.

**How it works**:
1. Kill story process with SIGKILL during registration (Phase=PhaseInitializing)
2. `MarkFailed` never runs — session file stays at PhaseInitializing
3. On restart, `dkgSvcRound` is 0 (no goroutines)
4. `isSessionStuckForStage(PhaseInitializing, Registration)` returns true
5. `tryAcquireDKGSvc` succeeds → confirms no active goroutine
6. `MarkFailed` called → picked up as PhaseFailed on next block

**Steps**:

```bash
# Step 1: Wait for DKG round to start
timeout 180 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"3 minutes ago\" | grep -q\"Initiated new DKG round\"'; do sleep 5; done"

# Step 2: Kill STORY PROCESS (not kernel) with SIGKILL during registration
# This must happen BEFORE "DKG initialization complete" appears
ssh $SSH_USER@$VAL1_IP "sudo systemctl stop story"

# Step 3: Check session file — should show intermediate phase
ssh $SSH_USER@$VAL1_IP "cat ~/.story/story/data/dkg/session_1.json 2>/dev/null | python3 -c\"import sys,json; d=json.load(sys.stdin); print('phase:', d.get('phase', 'N/A'))\"" || echo "No session file (registration may not have started)"

# Step 4: Verify kernel is still running on Val1
ssh $SSH_USER@$VAL1_IP "sudo lsof -i :50051 | grep LISTEN && echo 'KERNEL STILL RUNNING'"

# Step 5: Restart story (kernel stays running)
ssh $SSH_USER@$VAL1_IP "sudo systemctl start story"
sleep 10

# Step 6: Watch for stuck recovery log
timeout 120 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"2 minutes ago\" | grep -q\"Recovering stuck DKG session\"'; do sleep 5; done" && echo "STUCK RECOVERY DETECTED" || echo "NO RECOVERY LOG (session may have completed before kill)"

# Step 7: Wait for registration to complete
timeout 300 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"5 minutes ago\" | grep -q\"DKG initialization complete\"'; do sleep 10; done" && echo "REGISTRATION RECOVERED" || echo "TIMEOUT"

# Step 8: Wait for full round completion
timeout 1800 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"30 minutes ago\" | grep -q\"DKG finalization phase complete\"'; do sleep 15; done" && echo "ROUND COMPLETE" || echo "TIMEOUT"
```

**Verification**:

```bash
# Check stuck recovery log
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Recovering stuck DKG session: no active goroutine for intermediate phase'" && echo "PASS: Stuck recovery triggered"

# Check round completed
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'DKG finalization phase complete'" && echo "PASS: Round completed after recovery"
```

**Success Criteria**:
- Log: “Recovering stuck DKG session: no active goroutine for intermediate phase”
- Session MarkFailed → resume → registration completes
- Full round completes

---

### TC-05: Justification Flow — Complaint and Recovery

**Objective**: Verify complaint → justification → broadcast → Schnorr + Pedersen VSS verification pipeline.

**How it works**:
1. Kernel’s `--dkg-test-corrupt-deal` flag causes GenerateDeals to produce one deal with corrupted share value (V+1)
2. The deal passes Schnorr signature verification but fails `VerifyDeal` (share doesn’t match commitments)
3. Recipient returns StatusComplaint response
4. Dealer processes complaint → generates justification with plaintext deal
5. Justification broadcast via vote extension
6. All validators verify justification: Schnorr signature + Pedersen VSS
7. If deal was genuinely invalid, VSS verification confirms “No valid justifications” (expected)

**Prerequisites**:
- story-kernel with `--dkg-test-corrupt-deal` flag implemented (see “Kernel Code Changes Required” below)
- All 3 validators must use the SAME binary (same mrenclave) — only Val1 is started with the flag

### Kernel Code Changes Required

The `--dkg-test-corrupt-deal` flag is NOT in the `main` branch. A reference implementation is available on the `dkg/test-corrupt-deal` branch. Cherry-pick this commit before building:

```bash
# On ALL 3 validators (mrenclave must match):
cd ~/story-kernel
git fetch origin
git cherry-pick origin/dkg/test-corrupt-deal
make clean && make build-with-cpp && make all-gramine

# Verify mrenclave matches across all validators
cat ~/story-kernel/story-kernel.manifest.sgx.d/mrenclave.txt
```

**What it does**: Uses `reflect` + `unsafe.Pointer` to access kyber’s unexported `vss.Dealer` inside `DistKeyGenerator`, modifies one deal’s share value (V+1) in-place, re-encrypts via `dealer.EncryptedDeal()`, and re-signs the DKG-level Schnorr signature. The corrupted deal passes signature verification but fails VSS share verification → `StatusComplaint`.

**Why simple cipher XOR doesn’t work**: The Schnorr signature covers the encrypted cipher bytes. XOR-ing the cipher invalidates the signature, causing deal REJECTION (not complaint).

Only Val1 is started with the `--dkg-test-corrupt-deal` flag. Val2 and Val3 use the same binary but start normally (without the flag).

**Steps**:

```bash
# Step 1: Restart kernel on Val1 ONLY with corruption flag
ssh $SSH_USER@$VAL1_IP "sudo kill -9\$(sudo lsof -t -i :50051) 2>/dev/null; sudo pkill -9 -x gramine-sgx 2>/dev/null; sudo pkill -9 -x loader 2>/dev/null; sleep 3"
ssh $SSH_USER@$VAL1_IP "cd ~/story-kernel && nohup gramine-sgx story-kernel start --home ~/.story-kernel --dkg-test-corrupt-deal > /tmp/kernel.log 2>&1 &"
sleep 60

# Step 2: Restart story to re-establish kernel connection
ssh $SSH_USER@$VAL1_IP "sudo systemctl restart story"
sleep 5

# Step 3: Wait for round 1 to start and complete registration
for IP in $ALL_VALS; do
  timeout 420 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"DKG initialization complete\"'; do sleep 10; done" && echo "REGISTERED:$IP" || echo "TIMEOUT:$IP"
done

# Step 4: Wait for dealing phase — Val1 generates corrupted deal
timeout 600 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"DKG deals are generated successfully\"'; do sleep 10; done"

# Step 5: Watch for complaint response on Val2 or Val3
for IP in $VAL2_IP $VAL3_IP; do
  ssh $SSH_USER@$IP "journalctl -u story --since '10 minutes ago' | grep -E '(complaint|justification)' || echo 'No complaint yet on$IP'"
done

# Step 6: Watch for justification broadcast from Val1 (the dealer)
timeout 300 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"Enqueued justifications for broadcast\"'; do sleep 10; done" && echo "JUSTIFICATIONS ENQUEUED" || echo "TIMEOUT"

# Step 7: Watch for justification processing on all validators
for IP in $ALL_VALS; do
  ssh $SSH_USER@$IP "journalctl -u story --since '15 minutes ago' | grep 'Handling DKG process justifications'" && echo "JUSTIFICATION PROCESSED:$IP" || echo "NOT YET:$IP"
done

# Step 8: Check VSS verification result
for IP in $ALL_VALS; do
  ssh $SSH_USER@$IP "journalctl -u story --since '15 minutes ago' | grep -E '(No valid justifications|Justification VSS verification)'"
done

# Step 9: After testing, restart kernel without corruption flag
ssh $SSH_USER@$VAL1_IP "sudo kill -9\$(sudo lsof -t -i :50051) 2>/dev/null; sudo pkill -9 -x gramine-sgx 2>/dev/null; sudo pkill -9 -x loader 2>/dev/null; sleep 3"
ssh $SSH_USER@$VAL1_IP "cd ~/story-kernel && nohup gramine-sgx story-kernel start --home ~/.story-kernel > /tmp/kernel.log 2>&1 &"
```

**Verification**:

```bash
# Check kernel logged corruption
ssh $SSH_USER@$VAL1_IP "cat /tmp/kernel.log | grep 'Corrupted deal share'" && echo "PASS: Deal corruption triggered"

# Check justification enqueued
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '20 minutes ago' | grep 'Enqueued justifications for broadcast'" && echo "PASS: Justification enqueued"

# Check justification processed by all
for IP in $ALL_VALS; do
  ssh $SSH_USER@$IP "journalctl -u story --since '20 minutes ago' | grep 'Handling DKG process justifications'" && echo "PASS: Justification processed on$IP"
done

# Check VSS verification ran
ssh $SSH_USER@$VAL2_IP "journalctl -u story --since '20 minutes ago' | grep -E '(No valid justifications|Justification VSS verification failed)'" && echo "PASS: VSS verification completed (deal was genuinely invalid)"
```

**Success Criteria**:
- Kernel produces deal with corrupted share (valid signature, invalid VSS)
- Recipient generates StatusComplaint response
- Dealer enqueues justification for broadcast
- All validators process justification
- Schnorr signature verification passes
- Pedersen VSS verification runs (result: “No valid justifications” because deal was genuinely bad)
- DKG round still completes via threshold (other deals valid)

---

### TC-06: Story Client Restart During Dealing — Block Resync

**Objective**: Verify that when story stops and restarts during dealing, it processes other validators’ deals from replayed blocks (vote extensions).

**How it works (different from TC-07)**:
1. Story process stops during dealing (before GenerateDeals)
2. Val2/Val3 generate and broadcast deals via vote extensions
3. Deals committed in blocks while Val1 is down
4. Val1 restarts → CometBFT resyncs missed blocks → FinalizeBlock replays VE data
5. VE replay delivers Val2/Val3 deals → `handleDKGProcessDeals`
6. ResumeDKGService detects incomplete dealing → `handleDKGDealing` → GenerateDeals

Key difference from TC-07: here story restarts (block replay), not kernel crash (in-memory cache).

**Steps**:

```bash
# Step 1: Wait for all 3 registered
for IP in $ALL_VALS; do
  timeout 420 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"DKG initialization complete\"'; do sleep 10; done" && echo "REGISTERED:$IP" || echo "TIMEOUT:$IP"
done

# Step 2: Watch for dealing phase start on Val1
timeout 600 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"Handling DKG dealing\"'; do sleep 5; done"

# Step 3: IMMEDIATELY stop story on Val1 (before GenerateDeals finishes, ~2s window)
ssh $SSH_USER@$VAL1_IP "sudo systemctl stop story"

# Step 4: Wait 30s — Val2/Val3 generate and broadcast deals
sleep 30

# Step 5: Verify kernel still running on Val1
ssh $SSH_USER@$VAL1_IP "sudo lsof -i :50051 | grep LISTEN && echo 'KERNEL RUNNING'"

# Step 6: Restart story on Val1
ssh $SSH_USER@$VAL1_IP "sudo systemctl start story"
sleep 10

# Step 7: Watch for deal processing during resync
timeout 300 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"5 minutes ago\" | grep -q\"Handling DKG process deals\"'; do sleep 10; done" && echo "DEALS PROCESSED DURING RESYNC" || echo "TIMEOUT"

# Step 8: Watch for Val1's own deal generation
timeout 300 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"5 minutes ago\" | grep -q\"DKG deals are generated successfully\"'; do sleep 10; done" && echo "OWN DEALS GENERATED" || echo "TIMEOUT"

# Step 9: Wait for round completion
timeout 1800 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"30 minutes ago\" | grep -q\"DKG finalization phase complete\"'; do sleep 15; done" && echo "ROUND COMPLETE" || echo "TIMEOUT"
```

**Verification**:

```bash
# Check Val1 processed deals during resync
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Handling DKG process deals'" && echo "PASS: Deals processed during resync"

# Check Val1 generated its own deals
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'DKG deals are generated successfully'" && echo "PASS: Own deals generated"

# Check round completed
FINALIZED_COUNT=0
for IP in $ALL_VALS; do
  if ssh $SSH_USER@$IP "journalctl -u story --since '30 minutes ago' | grep -q 'DKG finalization phase complete'"; then
    FINALIZED_COUNT=$((FINALIZED_COUNT + 1))
  fi
done
echo "Finalized:$FINALIZED_COUNT/3"
test "$FINALIZED_COUNT" -ge 2 && echo "PASS: Threshold met" || echo "FAIL"
```

**Success Criteria**:
- Val1 processes Val2/Val3 deals during block resync (“Handling DKG process deals”)
- Val1 generates its own deals after resync (“DKG deals are generated successfully”)
- Round completes (2/3 or 3/3 finalized)

---

### TC-07: Deal/Response/Justification Caching — Kernel Down During Deal Processing

**Objective**: Verify in-memory caching of deals, responses, and justifications when kernel is temporarily unreachable, and ordered replay when kernel recovers.

**Key difference from TC-06**: Story stays running (in-memory cache preserved). Only kernel is blocked.

**How it works**:
1. Block kernel gRPC port with iptables (kernel stays alive, just unreachable)
2. Deals arrive via VE → kernel gRPC fails → `cachePendingIncomingDeals`
3. Responses arrive → kernel gRPC fails → `cachePendingIncomingResponses`
4. Justifications arrive (if any complaints) → kernel gRPC fails → `cachePendingIncomingJustifications`
5. Remove iptables rules → kernel becomes reachable
6. `BeginBlocker` → `reprocessPendingIncomingData` → drain in order: deals → responses → justifications
7. Cached data processed successfully

**Why iptables instead of pkill**: SGX enclave processes (gramine-sgx) are difficult to kill reliably (zombie processes). iptables blocks at network level while keeping both processes alive.

**Steps**:

```bash
# Step 1: Wait for all 3 registered
for IP in $ALL_VALS; do
  timeout 420 bash -c "while ! ssh$SSH_USER@$IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"DKG initialization complete\"'; do sleep 10; done" && echo "REGISTERED:$IP" || echo "TIMEOUT:$IP"
done

# Step 2: Apply iptables BEFORE dealing starts (block gRPC to kernel on Val1)
ssh $SSH_USER@$VAL1_IP "
  sudo iptables -A INPUT -p tcp --dport 50051 -j DROP
  sudo iptables -A OUTPUT -p tcp --sport 50051 -j DROP
"

# Step 3: Verify iptables rules applied
ssh $SSH_USER@$VAL1_IP "sudo iptables -L -n | grep 50051"

# Step 4: Wait for Dealing phase → deals will fail with DeadlineExceeded
timeout 600 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"10 minutes ago\" | grep -q\"cached for retry\"'; do sleep 10; done" && echo "DEALS CACHED" || echo "TIMEOUT"

# Step 5: Verify caching logs (deals, responses, and optionally justifications)
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '10 minutes ago' | grep -E 'cached for retry|Cached responses for retry|Cached justifications for retry'"

# Step 6: Wait 30 seconds, then remove iptables rules (unblock kernel)
sleep 30
ssh $SSH_USER@$VAL1_IP "
  sudo iptables -D INPUT -p tcp --dport 50051 -j DROP
  sudo iptables -D OUTPUT -p tcp --sport 50051 -j DROP
"

# Step 7: Verify iptables rules removed
ssh $SSH_USER@$VAL1_IP "sudo iptables -L -n | grep 50051 || echo 'RULES REMOVED'"

# Step 8: Watch for replay logs
timeout 120 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"2 minutes ago\" | grep -q\"Replaying cached\"'; do sleep 5; done" && echo "REPLAY DETECTED" || echo "TIMEOUT"

# Step 9: Check replay ordering (deals → responses → justifications)
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '5 minutes ago' | grep -E 'Replaying cached (deals|responses|justifications)'"

# Step 10: Wait for round completion
timeout 1800 bash -c "while ! ssh$SSH_USER@$VAL1_IP 'journalctl -u story --since\"30 minutes ago\" | grep -q\"DKG finalization phase complete\"'; do sleep 15; done" && echo "ROUND COMPLETE" || echo "TIMEOUT"
```

**Verification**:

```bash
# Check deals were cached
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'cached for retry'" && echo "PASS: Deals cached"

# Check replay happened in correct order
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Replaying cached deals'" && echo "PASS: Deals replayed"
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Replaying cached responses'" && echo "PASS: Responses replayed"
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Replaying cached justifications'" && echo "PASS: Justifications replayed" || echo "INFO: No justifications cached (no complaints occurred)"

# Verify replay ordering: deals before responses before justifications
DEAL_TIME=$(ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Replaying cached deals' | head -1 | awk '{print\$1,\$2,\$3}'")
RESP_TIME=$(ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Replaying cached responses' | head -1 | awk '{print\$1,\$2,\$3}'")
JUST_TIME=$(ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'Replaying cached justifications' | head -1 | awk '{print\$1,\$2,\$3}'")
echo "Deals replayed at:$DEAL_TIME"
echo "Responses replayed at:$RESP_TIME"
echo "Justifications replayed at:$JUST_TIME"

# Check round completed
ssh $SSH_USER@$VAL1_IP "journalctl -u story --since '30 minutes ago' | grep 'DKG finalization phase complete'" && echo "PASS: Round completed"
```

**Cleanup** (always run, even on failure):

```bash
ssh $SSH_USER@$VAL1_IP "
  sudo iptables -D INPUT -p tcp --dport 50051 -j DROP 2>/dev/null
  sudo iptables -D OUTPUT -p tcp --sport 50051 -j DROP 2>/dev/null
  echo 'Cleanup done'
"
```

**Success Criteria**:
- Deals cached: “cached for retry” log with cached_deals > 0
- Responses cached: “Cached responses for retry” log
- Justifications cached (if complaints occurred): “Cached justifications for retry” log
- Replay order: “Replaying cached deals” → “Replaying cached responses” → “Replaying cached justifications”
- “Process deals complete” + “Process responses complete” after replay
- Round completes

---

## 4. Success Criteria Summary

| Test | Type | Key Verification | Minimum Pass |
| --- | --- | --- | --- |
| HP-01 | Happy | 3/3 finalized, GlobalPubKey set | 2/3 finalized |
| HP-02 | Happy | Resharing round completes, IsResharing=true | 2/3 finalized |
| HP-03 | Happy | 4/4 participate in resharing | 3/4 finalized |
| HP-04 | Happy | Upgrade resharing, dual binary | 2/3 finalized |
| HP-05 | Happy | Rewards distributed from UBI pool | Balances increased |
| TC-01 | Unhappy | Kernel kill → resume → registration | Round completes |
| TC-02 | Unhappy | Kernel kill → resume → dealing | Round completes (threshold) |
| TC-03 | Unhappy | Kernel kill → resume → finalization | Round completes (threshold) |
| TC-04 | Unhappy | Process crash → stuck recovery | “Recovering stuck” log |
| TC-05 | Unhappy | Corrupt deal → complaint → justification pipeline | “Enqueued justifications” log |
| TC-06 | Unhappy | Story restart → block resync → deal replay | Deals processed during resync |
| TC-07 | Unhappy | iptables block → cache → replay | “Replaying cached” logs |

### Automated Pass/Fail Pattern

```bash
check_log() {
  local IP=$1
  local PATTERN=$2
  local SINCE=$3
  ssh $SSH_USER@$IP "journalctl -u story --since '$SINCE' | grep -q '$PATTERN'" && return 0 || return 1
}

# Example: TC-01 pass criteria
TC01_PASS=true
check_log $VAL1_IP "DKG initialization complete" "30 minutes ago" || TC01_PASS=false
check_log $VAL1_IP "DKG finalization phase complete" "30 minutes ago" || TC01_PASS=false
echo "TC-01:$TC01_PASS"
```

---

## 5. Test Cleanup

Every test MUST have a cleanup step that runs on success AND failure:

```bash
cleanup() {
  # Remove iptables rules on all validators
  for IP in $ALL_VALS; do
    ssh $SSH_USER@$IP "
      sudo iptables -D INPUT -p tcp --dport 50051 -j DROP 2>/dev/null || true
      sudo iptables -D OUTPUT -p tcp --sport 50051 -j DROP 2>/dev/null || true
    " &
  done
  wait

  # Restart kernel without test flags
  for IP in $ALL_VALS; do
    ssh $SSH_USER@$IP "
      sudo kill -9\$(sudo lsof -t -i :50051) 2>/dev/null
      sudo pkill -9 -x gramine-sgx 2>/dev/null
      sudo pkill -9 -x loader 2>/dev/null
      sleep 2
      cd ~/story-kernel && nohup gramine-sgx story-kernel start --home ~/.story-kernel > /tmp/kernel.log 2>&1 &
    " &
  done
  wait

  # Restart story and geth if stopped
  for IP in $ALL_NODES; do
    ssh $SSH_USER@$IP "sudo systemctl start story 2>/dev/null; sudo systemctl start node-geth 2>/dev/null" &
  done
  wait
}

trap cleanup EXIT
```