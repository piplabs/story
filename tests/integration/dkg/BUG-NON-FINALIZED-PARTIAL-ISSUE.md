# Bug: Non-finalized validator can submit partial decryptions and earn CDR fee rewards

## Summary

A validator that registered but never finalized (status = `Verified`, not `Finalized`) can submit partial decryption responses that pass all on-chain validation checks. Each successful submission increments the validator's `CDRPartialSubmitCount`, entitling them to a share of the CDR fee pool — without having contributed any valid key share to the DKG committee.

## Severity

**Medium-High** — Economic extraction with zero useful contribution. No exploit or transaction flooding required; the non-finalized validator simply follows normal decrypt request traffic.

## Environment

- **story**: branch `dkg/dev`, commit `dd43702`
- **story-kernel**: branch `main`, commit `10f9856`

## Root Cause

Two missing checks:

1. `PartialDecryptionSubmitted` does not verify `reg.Status == DKGRegStatusFinalized` after retrieving the registration
2. The `pubShare` comparison uses `bytes.Equal(pubShare, reg.PubKeyShare)` where both values are `nil` for a non-finalized validator, and `bytes.Equal(nil, nil)` returns `true` in Go

## Code Evidence

### 1. PubKeyShare is only set during finalization

**`client/x/dkg/keeper/dkg_registration.go:78-88`**:
```go
func (k *Keeper) finalizeDKGRegistration(ctx context.Context, round uint32, validatorAddr common.Address, pubKeyShare []byte) error {
    dkgReg, err := k.getDKGRegistration(ctx, round, validatorAddr)
    dkgReg.PubKeyShare = pubKeyShare                    // only written here
    dkgReg.Status = types.DKGRegStatusFinalized
    return k.setDKGRegistration(ctx, validatorAddr, dkgReg)
}
```

A validator that never calls this function has `PubKeyShare = nil` and `Status = Verified`.

### 2. PartialDecryptionSubmitted has no status check, and nil pubShare passes

**`client/x/dkg/keeper/dkg_handler.go:541-555`**:
```go
reg, err := k.getDKGRegistration(ctx, req.Round, validator)
if err != nil {
    return errors.Wrap(err, "failed to get DKG registration for signature verification")
}
// ← No check: reg.Status == DKGRegStatusFinalized

if !bytes.Equal(pubShare, reg.PubKeyShare) {
    return errors.New("pubShare mismatch: submitted pubShare does not match stored pubKeyShare", ...)
}
// pubShare = attacker submits nil/empty
// reg.PubKeyShare = nil (never finalized)
// bytes.Equal(nil, nil) = true → PASSES

if err := verifyPartialDecryptionSignature(reg.CommPubKey, round, ciphertext, ...); err != nil {
    return errors.Wrap(err, "partial decryption signature verification failed")
}
// CommPubKey was set during registration → attacker signs with valid key → PASSES
```

All checks pass. `PartialDecryptionSubmitted` returns `nil` (success).

### 3. Successful submission increments CDR fee count and refunds baseFee

**`client/x/evmengine/keeper/cdr.go:174-181`**:
```go
if partialErr == nil {
    if err := k.dkgKeeper.IncrementCDRPartialSubmitCount(cachedCtx, ev.Validator); err != nil {
        partialErr = errors.Wrap(err, "increment CDR submit count")
    } else if ev.Fee != nil && ev.Fee.Sign() > 0 {
        if err := k.dkgKeeper.RefundCDRFee(cachedCtx, ev.Validator, ev.Fee); err != nil {
            partialErr = errors.Wrap(err, "refund CDR fee")
        }
    }
}
```

`partialErr` is `nil` → count incremented → baseFee refunded from pool.

### 4. CDR fee pool distribution does not filter by registration status

**`client/x/dkg/keeper/dkg_cdr_fees.go:132-198`**:
```go
iter, err := k.CDRPartialSubmitCount.Iterate(ctx, nil)
counts := map[string]uint64{}
var totalCount uint64

for ; iter.Valid(); iter.Next() {
    key, _ := iter.Key()     // validator address string
    count, _ := iter.Value()
    counts[key] += count     // no registration status check
    totalCount += count
}

for addr, count := range counts {
    share := poolBalance.Mul(math.NewInt(int64(count))).Quo(totalCountInt)
    // non-finalized validator receives share proportional to their count
    k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.CDRFeePoolName, recipient, coins)
}
```

Compare with UBI reward distribution which correctly filters:

**`client/x/dkg/keeper/dkg_rewards.go:59`**:
```go
finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, round.Round, types.DKGRegStatusFinalized)
// ← UBI rewards are only distributed to Finalized validators
```

CDR fee distribution has no equivalent filter.

## Attack Path

No exploit needed. The non-finalized validator simply follows normal traffic:

```
Setup: 3 validators. B and C finalize normally. A registers but does NOT finalize.

1. A user calls CDR.read(uuid) → ThresholdDecryptRequested event

2. B submits valid partial → count_B += 1
   C submits valid partial → count_C += 1
   A submits partial with empty pubShare → passes all checks → count_A += 1

3. This repeats for every user read request throughout the Active round.
   After 100 user reads:
     count_A = 100, count_B = 100, count_C = 100

4. distributeCDRRewardPool():
     A receives: poolBalance × 100/300 = 33.3%
     B receives: poolBalance × 100/300 = 33.3%
     C receives: poolBalance × 100/300 = 33.3%
```

A earns an equal share of CDR fees as honest validators, despite never completing DKG finalization and never providing a usable key share.

## Impact

- **Unearned economic reward**: Non-finalized validator extracts CDR fee pool share without contributing any useful decryption capability
- **Dilutes honest validators' income**: With 1 non-finalized validator out of 3, honest validators each lose 1/6 of their expected CDR fee share (from 50% each down to 33.3%)
- **Garbage partial occupies dedup slot**: The stored partial is cryptographically useless (empty pubShare, no valid DistKeyShare). If a client selects this partial for decryption combining, it fails
- **No special tooling needed**: Attacker only needs a registered validator that skips finalization — no mock kernel, no modified binary

## Suggested Fix

Add a status check in `PartialDecryptionSubmitted` after retrieving the registration:

```go
reg, err := k.getDKGRegistration(ctx, req.Round, validator)
if err != nil {
    return errors.Wrap(err, "failed to get DKG registration")
}

if reg.Status != types.DKGRegStatusFinalized {
    return errors.New("validator has not finalized for this round",
        "validator", validator.Hex(),
        "round", req.Round,
        "status", reg.Status.String(),
    )
}

if len(pubShare) == 0 || len(reg.PubKeyShare) == 0 {
    return errors.New("pubShare or registered pubKeyShare is empty",
        "validator", validator.Hex(),
        "round", req.Round,
    )
}
```

This blocks non-finalized validators from submitting partials, which also prevents their CDR fee count from incrementing.
