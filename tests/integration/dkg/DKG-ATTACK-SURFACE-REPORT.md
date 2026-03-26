# DKG Module Attack Surface Report

> Excludes the "bad dealer bypass" bug (already filed and being fixed).

## Summary

| Priority | ID | Attack | Severity | Core Issue |
|----------|----|--------|----------|------------|
| **P1** | F | CDR fee farming by legitimate validator | Medium-High | Validator self-generates read requests to inflate CDRPartialSubmitCount, extracting disproportionate CDR fee pool share |
| **P1** | D | Stage transition purely time-driven | Medium | Malicious late deal broadcast repeatedly fails rounds |
| **P1** | E | Non-finalized validator can submit partial | Medium | Missing `reg.Status == Finalized` check in PartialDecryptionSubmitted |
| **P1** | B | Stale ActiveValSet snapshot | Medium | Jailed validator can still participate for ~3 days |
| **P2** | C | Registry bloat DoS | Medium | Low-cost accumulation of expired entries impacts cleanup block |
| **P2** | A | Timeout deletion race | Low | delete-on-timeout design flaw; not reliably exploitable in practice |

---

## Attack F: CDR Fee Farming by Legitimate Validator (P1 — Medium-High)

### Description

A legitimate active committee member can inflate their `CDRPartialSubmitCount` by self-generating decrypt requests via `CDR.read()`, then submitting partials for each. Since CDR fee pool distribution is proportional to submit count, the attacker extracts a disproportionate share of the pool.

### Code Evidence

**Fee pool distribution is proportional to submit count** — `client/x/dkg/keeper/dkg_cdr_fees.go:182-198`:
```go
for addr, count := range counts {
    share := poolBalance.Mul(math.NewInt(int64(count))).Quo(totalCountInt)
    // Validator with highest count gets the largest share
}
```

**Each successful partial submission increments count** — `client/x/evmengine/keeper/cdr.go:174-175`:
```go
if partialErr == nil {
    k.dkgKeeper.IncrementCDRPartialSubmitCount(cachedCtx, ev.Validator)  // count+1 per submission
```

**baseFee is refunded from pool after successful submission** — `client/x/evmengine/keeper/cdr.go:177-179`:
```go
    if ev.Fee != nil && ev.Fee.Sign() > 0 {
        k.dkgKeeper.RefundCDRFee(cachedCtx, ev.Validator, ev.Fee)  // baseFee returned from pool
    }
}
```

**CDR.sol `submitEncryptedPartialDecryption` has no caller restriction** — `contracts/src/protocol/CDR.sol:223-251`:
```solidity
function submitEncryptedPartialDecryption(...) external payable whenNotPaused {
    uint256 fee = _getCDRStorage().baseFee;
    _collectFee(fee, ICDR.FeeType.SubmitPartial);  // baseFee burned at contract level
    emit EncryptedPartialDecryptionSubmitted(msg.sender, ...);  // no committee membership check
}
```

**Dedup is per (requesterPubKey, label, ciphertext, round, validator)** — `client/x/dkg/keeper/dkg_partial_decryption.go:17-27`:
```go
func dkgPartialDecryptKey(requesterPubKey []byte, label []byte, ciphertext []byte, round uint32, validator common.Address) string {
    // Each unique decrypt request allows one submission per validator
}
```

### Attack Path

```
Attacker = Validator A (legitimate finalized committee member)

1. A calls CDR.sol read(uuid) N times using different addresses or same address
   → Pays readFee × N (burned)
   → Each call emits VaultRead event with unique (requesterPubKey, ciphertext) pair
   → Creates N distinct decrypt requests in DecryptRequestRegistry

2. A's kernel receives N ThresholdDecryptRequested events
   → Generates N valid partial decryptions (A has a correct DistKeyShare)

3. A calls CDR.sol submitEncryptedPartialDecryption() for each
   → Pays baseFee × N at contract layer (burned)
   → Consensus layer: each passes PartialDecryptionSubmitted validation
   → IncrementCDRPartialSubmitCount: A's count += N
   → RefundCDRFee: baseFee × N returned to A from CDR fee pool

4. End of Active round → distributeCDRRewardPool():
   count_A = 100, count_B = 5, count_C = 5, totalCount = 110
   A receives: poolBalance × 100/110 = 90.9% of CDR fee pool
   B receives: poolBalance × 5/110 = 4.5%
   C receives: poolBalance × 5/110 = 4.5%
```

### Cost-Benefit Analysis

| Item | Cost |
|------|------|
| Issue N read requests | `readFee × N` (burned at contract) |
| Submit N partials | `baseFee × N` (refunded from pool) |
| Gas | `~2 txs × N × gas_price` |
| **Net cost** | `readFee × N + gas × 2N` |
| **Revenue** | `CDR_fee_pool × N / (N + others_count)` |

**Profitable when**: `CDR_fee_pool / (N + others) > readFee + 2 × gas_per_tx`

If `readFee` is low relative to pool accumulation, the attack is profitable.

### Impact

- Attacker extracts up to ~90%+ of CDR fee pool
- Honest validators' CDR fee income diluted to near-zero
- Creates perverse incentive: all validators race to self-generate read requests → pool gets drained by gas/readFee waste
- Attack is invisible from on-chain state — looks like legitimate read activity

### Suggested Fix

Distribute CDR fee pool equally among all finalized committee members (flat distribution), or cap submissions per validator per round, or weight by unique decrypt requests that actually reached threshold completion.

---

## Attack D: Stage Transition Purely Time-Driven (P1 — Medium)

### Description

DKG stage transitions are based solely on elapsed block height, with no consideration of actual protocol completion. An attacker can exploit this by deliberately delaying deal broadcasts until the final blocks of the Dealing stage.

### Code Evidence

**`client/x/dkg/keeper/dkg_round.go:10-48`**:
```go
func (*Keeper) shouldTransitionStage(currentHeight int64, dkgNetwork *types.DKGNetwork, params types.Params) (types.DKGStage, bool) {
    elapsed := currentHeight - dkgNetwork.StartBlockHeight
    dealingEnd := registrationEnd + int64(params.DealingPeriod)

    case types.DKGStageDealing:
        if elapsed >= dealingEnd {
            return types.DKGStageFinalization, true  // forced transition regardless of deal completion
        }
}
```

No check like `"have all deals been exchanged?"` or `"have all responses been processed?"`.

### Attack Path

1. Dealing stage = 86400 blocks (~1 day)
2. Attacker controls 1+ validators, delays `GenerateDeals()` until block 86350 (50 blocks before deadline)
3. Deals arrive via Vote Extension at block ~86360
4. Honest validators receive deals → need to `ProcessDeals` + `ProcessResponses` + broadcast via VE
5. At block 86400, stage force-transitions to Finalization
6. Honest validators' responses/justifications haven't been fully processed
7. Incomplete validators can't finalize → `finalizedCount < threshold` → round fails
8. Repeat every round → CDR service permanently unavailable

### Impact

- Attacker can reliably force DKG round failures
- No economic penalty for the attacker (no slashing for late deals)
- Repeated failures prevent any active committee from forming

### Suggested Fix

Consider adding a minimum deal completion requirement before transitioning from Dealing to Finalization, or extend the stage deadline if deal exchange is still in progress.

---

## Attack E: Non-Finalized Validator Can Submit Partial Decryption and Earn CDR Fee Rewards (P1 — Medium)

### Description

`PartialDecryptionSubmitted` does not verify that the submitting validator has `DKGRegStatusFinalized`. A validator that only registered (status = `Verified`) but never finalized can submit partial decryptions, pass all on-chain checks, and have their `CDRPartialSubmitCount` incremented — earning a share of the CDR fee pool at the end of the Active round.

### Code Evidence

**Missing status check** — `client/x/dkg/keeper/dkg_handler.go:541-551`:
```go
reg, err := k.getDKGRegistration(ctx, req.Round, validator)
if err != nil {
    return errors.Wrap(err, "failed to get DKG registration for signature verification")
}
// ↑ Only checks registration EXISTS. No check for reg.Status == DKGRegStatusFinalized.

if !bytes.Equal(pubShare, reg.PubKeyShare) {
    return errors.New("pubShare mismatch...")
}
// ↑ For non-finalized registration, reg.PubKeyShare is nil/empty (proto3 default)
// ↑ If attacker submits empty pubShare → bytes.Equal(nil, nil) = true in Go
// ↑ Verified: bytes.Equal(nil, nil) == true, bytes.Equal(nil, []byte{}) == true
```

**PubKeyShare only set during finalization** — `client/x/dkg/keeper/dkg_registration.go:84`:
```go
func (k *Keeper) finalizeDKGRegistration(..., pubKeyShare []byte) error {
    dkgReg.PubKeyShare = pubKeyShare  // only set here — never set for non-finalized registrations
    dkgReg.Status = types.DKGRegStatusFinalized
}
```

**Successful submission increments CDR fee count** — `client/x/evmengine/keeper/cdr.go:174-175`:
```go
if partialErr == nil {  // PartialDecryptionSubmitted returns nil = success
    k.dkgKeeper.IncrementCDRPartialSubmitCount(cachedCtx, ev.Validator)  // count+1
```

**CDR fee pool distributed by count, no status filtering** — `client/x/dkg/keeper/dkg_cdr_fees.go:132-157`:
```go
iter, err := k.CDRPartialSubmitCount.Iterate(ctx, nil)
for ; iter.Valid(); iter.Next() {
    key, _ := iter.Key()    // key = validator address string
    count, _ := iter.Value()
    counts[key] += count    // aggregated by address only, no registration status check
    totalCount += count
}
// line 182-183: share = poolBalance * count / totalCount
```

### Attack Path

```
Validator A registers in DKG round but intentionally does NOT finalize.
reg.Status = DKGRegStatusVerified, reg.PubKeyShare = nil

1. A calls CDR.sol read(uuid) to trigger a decrypt request
   → Pays readFee

2. A calls CDR.sol submitEncryptedPartialDecryption() with empty pubShare:
   → CDR.sol: collects baseFee, emits event
   → dkg_handler.go:541: getDKGRegistration succeeds (A is registered)
   → dkg_handler.go:546: bytes.Equal([]byte{}, nil) = true → PASSES
   → dkg_handler.go:553: verifyPartialDecryptionSignature(reg.CommPubKey, ...) → PASSES
     (A has valid commPubKey from registration, signs correctly)
   → dkg_handler.go:557: setPartialDecryptionSubmission → garbage partial stored
   → PartialDecryptionSubmitted returns nil (success)

3. evmengine/keeper/cdr.go:175:
   → IncrementCDRPartialSubmitCount(A) → A's count += 1
   → RefundCDRFee(A, baseFee) → baseFee returned from pool

4. A repeats steps 1-3 for N different read requests
   → A's CDRPartialSubmitCount = N

5. End of Active round → distributeCDRRewardPool():
   → A's count is included in distribution (no Finalized status check)
   → A receives poolBalance × N / totalCount of CDR fee pool
```

### Impact

- **Free money**: Non-finalized validator earns CDR fee pool share without contributing to decryption
- **Zero useful work**: The submitted partials are garbage (empty pubShare, no valid DistKeyShare) — they cannot be used for actual decryption
- **Dedup slot pollution**: Each garbage submission occupies the (validator, request) dedup slot, preventing any later legitimate submission by the same validator
- **Combines with Attack F**: Non-finalized validator can self-generate read requests to farm CDR fees at scale, without even needing to complete DKG finalization
- **UBI rewards unaffected**: `dkg_rewards.go:59` filters by `DKGRegStatusFinalized`, so UBI committee rewards are not extractable — only CDR fee pool is vulnerable

### Suggested Fix

Add status check in `PartialDecryptionSubmitted` after retrieving registration:
```go
reg, err := k.getDKGRegistration(ctx, req.Round, validator)
if err != nil {
    return errors.Wrap(err, "failed to get DKG registration")
}

// Reject submissions from non-finalized validators
if reg.Status != types.DKGRegStatusFinalized {
    return errors.New("validator has not finalized for this round",
        "validator", validator.Hex(),
        "status", reg.Status.String(),
    )
}
```

Additionally, add a nil/empty check for `pubShare` before the `bytes.Equal` comparison:
```go
if len(pubShare) == 0 || len(reg.PubKeyShare) == 0 {
    return errors.New("pubShare or registered pubKeyShare is empty")
}
```

---

## Attack B: Stale ActiveValSet Snapshot (P1 — Medium)

### Description

`ActiveValSet` is snapshotted once at round initiation and never refreshed. A validator jailed or unbonded during the round can still participate in all subsequent stages.

### Code Evidence

**Snapshot at round start** — `client/x/dkg/keeper/dkg_initialization.go:44-60`:
```go
activeValidators, err := k.GetActiveValidators(ctx)  // Bonded && !Jailed at this moment
dkgNetwork := types.DKGNetwork{
    ActiveValSet: activeValidators,  // frozen snapshot
}
```

**Membership check uses stale snapshot** — `client/x/dkg/keeper/dkg_handler.go:47`:
```go
if !slices.Contains(latest.ActiveValSet, strings.ToLower(validator.Hex())) {
    return errors.New("msg sender is not in the active validator set")
}
```

### Attack Window

| Stage | Duration | Cumulative |
|-------|----------|------------|
| Registration | ~1 day | 1 day |
| Dealing | ~1 day | 2 days |
| Finalization | ~1 day | 3 days |

A validator jailed at any point during these 3 days remains in `ActiveValSet` and can complete the entire DKG flow.

### Impact

- Jailed validator has no economic stake at risk (already slashed) → zero cost for malicious behavior
- Can register with compromised keys, submit bad deals, finalize with wrong globalPubKey
- Combines with other attacks (e.g., late deal broadcast, garbage partial submission)

### Suggested Fix

Re-validate validator status (bonded, not jailed) at each stage transition (Registration → Dealing, Dealing → Finalization, Finalization → Active). Remove jailed validators from `ActiveValSet` during `BeginDealing` or `BeginFinalization`.

---

## Attack C: Decrypt Request Registry Bloat DoS (P2 — Medium)

### Description

An attacker can cheaply accumulate thousands of expired decrypt request entries in `DecryptRequestRegistry`, which are only pruned every 1000 blocks via a full-table scan in `BeginBlocker`.

### Code Evidence

**Cleanup interval** — `client/x/dkg/types/keys.go:23`:
```go
const DecryptRequestRegistryCleanupInterval int64 = 1000
```

**Full iteration on cleanup** — `client/x/dkg/keeper/dkg_decrypt_registry.go:66-94`:
```go
func (k *Keeper) pruneTimedOutDecryptRequests(ctx context.Context, currentHeight uint64) error {
    iter, err := k.DecryptRequestRegistry.Iterate(ctx, nil)  // scan ALL entries
    for ; iter.Valid(); iter.Next() {
        // check each entry for timeout
    }
    // then delete expired ones
}
```

**Each `read()` call creates a registry entry** — `client/x/evmengine/keeper/cdr.go:99`:
```go
k.dkgKeeper.ThresholdDecryptRequested(cachedCtx, round, ev.RequesterPubKey, ev.Ciphertext, label[:], ...)
```

### Attack Path

1. Attacker calls `CDR.read()` thousands of times across 999 blocks (cost: `readFee` per call)
2. All requests expire after 200 blocks but remain in registry until cleanup
3. At block N×1000, `pruneTimedOutDecryptRequests` scans all entries
4. With 10,000+ expired entries, the BeginBlocker iteration is expensive
5. If cleanup takes longer than block production time → block production delay

### Impact

- BeginBlocker latency spike every 1000 blocks
- Cost to attacker: `readFee × N` (potentially low)
- Worst case: temporary chain liveness impact during cleanup blocks

### Suggested Fix

- Add incremental cleanup (process max N entries per block instead of all at once)
- Or add a max registry size cap with LRU eviction
- Or increase `readFee` to make mass-creation expensive

---

## Attack A: Timeout Deletion Race (P2 — Low)

### Description

When a partial decryption submission arrives after the 200-block timeout, `PartialDecryptionSubmitted` deletes the decrypt request registry entry inline. In theory, this could cause other submissions in the same block (ordered after the timed-out one) to silently fail because the registry entry no longer exists.

### Code Evidence

**`client/x/dkg/keeper/dkg_handler.go:527-538`** — delete-on-timeout:
```go
currentHeight := uint64(sdk.UnwrapSDKContext(ctx).BlockHeight())
if currentHeight-req.Height > types.PartialDecryptionTimeoutBlocks {  // timeout = 200 blocks
    if err := k.deleteDecryptRequest(ctx, requesterPubKey, label, round, ciphertext); err != nil {
        return errors.Wrap(err, "failed to delete expired decrypt request registry entry")
    }
    return nil  // registry entry gone — subsequent submissions in same block see !found
}
```

**`client/x/dkg/keeper/dkg_handler.go:497-507`** — subsequent lookup fails:
```go
req, found, err := k.getDecryptRequest(ctx, requesterPubKey, label, round, ciphertext)
if !found {
    log.Info(ctx, "Partial decryption submitted for unknown or cleaned-up request", ...)
    return nil  // silently discarded
}
```

### Why This Is Low Severity

While the delete-on-timeout design is suboptimal, it is **not reliably exploitable** in practice:

1. **Validators submit partials immediately** — honest validators submit within a few blocks of receiving the decrypt request, not near the 200-block timeout boundary
2. **Threshold is typically met early** — with threshold=2 and 3 validators, 2 partials are collected within ~10 blocks, well before timeout
3. **Transaction ordering is not attacker-controlled** — requires being block proposer at exactly block H+201, while another validator also submits in the same block
4. **The scenario requires an unlikely conjunction** — a request stuck at threshold-1 partials for 200 blocks AND an attacker precisely timing a submission AND a late honest submission in the same block

### Impact

- Design flaw rather than exploitable attack
- In normal operation, the timeout boundary is never the bottleneck
- Could cause issues in extreme edge cases (network partition recovery near timeout)

### Suggested Fix

Remove inline `deleteDecryptRequest` from the timeout branch — just reject the submission and let the existing periodic `pruneTimedOutDecryptRequests` (runs every 1000 blocks) handle cleanup.
