# Bug: Bad dealer bypasses invalidation, joins active committee, submits invalid partial decryptions, and receives rewards

## Summary

A validator who sends cryptographically invalid deals (Pedersen VSS verification fails) is **never marked as `Invalidated` on-chain**. It can successfully finalize, join the active decryption committee, submit garbage partial decryptions that pass all on-chain checks, and receive committee rewards — degrading CDR decryption availability and extracting unearned economic value.

## Severity

**High** — Impacts security (decryption availability), economic integrity (unearned rewards), and system fault tolerance (silent committee degradation).

## Environment

- **story**: branch `dkg/dev`, commit `62ce0fd`
- **story-kernel**: branch `main`, commit `10f9856`

## Root Cause

`invalidateDealerRegistration()` is defined at `client/x/dkg/keeper/dkg_justification.go:206` but is **never called from production code**. It is only referenced in unit tests.

```go
// dkg_justification.go:198-205
// NOTE: This function is NOT called from ProcessJustifications because justification
// processing must not affect on-chain state. It is retained as a utility for
// potential future use (e.g., explicit slashing proposals).
func (k *Keeper) invalidateDealerRegistration(ctx context.Context, ...) error {
```

When a justification's VSS verification fails (proving the deal was indeed invalid), the justification is silently dropped with only a log message — no on-chain state change occurs:

```go
// dkg_svc_dealing.go:463-468
if !valid {
    log.Info(ctx, "Justification VSS verification failed (deal was invalid), dropping",
        "dealer_index", j.Index,
    )
    continue  // ← No call to invalidateDealerRegistration(). Status stays Verified.
}
```

## Attack Path

### Step 1: Attacker sends invalid deals → recipients generate complaints

Attacker's kernel returns deals with valid DH encryption but invalid VSS polynomial shares.

```
story-kernel/service/dkg_process_deals.go:75
  → distKeyGen.ProcessDeal(deal)
  → Kyber decrypts deal, checks share * G == C_0 + i*C_1 + ...
  → VSS fails → returns Response.Status = false (complaint)
```

Complaint is enqueued for broadcast via Vote Extension (`client/x/dkg/keeper/dkg_svc_dealing.go:207`).

### Step 2: Attacker receives complaint → generates justification → VSS verification fails → silently dropped

Attacker's kernel produces a justification revealing the plaintext deal:

```
story-kernel/service/dkg_process_responses.go:90
  → distKeyGen.ProcessResponse(resp)
  → Detects Status=false → generates Justification (plaintext deal + Schnorr signature)
```

Consensus layer verifies the justification:

```
client/x/dkg/keeper/dkg_svc_dealing.go:434
  → verifyJustificationSignature() → PASSES (attacker's real Ed25519 key)

client/x/dkg/keeper/dkg_svc_dealing.go:454
  → verifyJustification() → valid=false (share doesn't match commitments)

client/x/dkg/keeper/dkg_svc_dealing.go:464
  → "Justification VSS verification failed (deal was invalid), dropping"
  → continue  ← NO on-chain state change. Attacker stays DKGRegStatusVerified.
```

### Step 3: Attacker calls `finalize()` → passes all on-chain checks → becomes `Finalized`

```
client/x/dkg/keeper/dkg_handler.go:117-148

  line 117: if reg.Status == DKGRegStatusFinalized → PASSES (first attempt)
  line 122: if reg.Status == DKGRegStatusInvalidated → PASSES (status is Verified, NEVER set to Invalidated)
  line 126: validateParticipantsRoot() → PASSES (attacker is Verified, included in root hash)
  line 130: verifyFinalizationSignature() → PASSES (attacker's valid ECDSA signature)
  line 134: AddGlobalPubKeyVote() → attacker votes different globalPubKey (minority vote, but does NOT block finalization)
  line 148: finalizeDKGRegistration() → attacker marked as DKGRegStatusFinalized
```

Critical: `AddGlobalPubKeyVote` at line 134 returns vote count, but **finalization continues regardless of whether the attacker's globalPubKey matches consensus**. Lines 134→148 are sequential with no early return on minority vote.

```
client/x/dkg/keeper/dkg_registration.go:78-88

  dkgReg.PubKeyShare = pubKeyShare  // line 84: stores attacker's WRONG pubKeyShare
  dkgReg.Status = types.DKGRegStatusFinalized  // line 85: unconditionally Finalized
```

### Step 4: Attacker counted in `finalizedCount` → round enters Active

```
client/x/dkg/keeper/dkg_finalization.go:30-61

  finalizedCount = countDKGRegistrationsByStatus(round, DKGRegStatusFinalized)
  // finalizedCount includes attacker (3 instead of correct 2)

  line 41: finalizedCount >= MinReqFinalizedParticipants → 3 >= 3 PASSES
  line 53: finalizedCount >= Threshold → 3 >= 2 PASSES
  // → Round becomes Active with inflated committee size
```

### Step 5: Attacker receives committee rewards

**UBI rewards**:

```
client/x/dkg/keeper/dkg_rewards.go:59
  finalizedRegs = getDKGRegistrationsByStatus(round, DKGRegStatusFinalized)
  // ↑ includes attacker

client/x/dkg/keeper/dkg_rewards.go:79
  perMemberReward = dkgReward / memberCount  // memberCount=3 instead of correct 2
  // Attacker receives 1/3 of reward. Honest validators each lose 1/6.
```

**CDR fee pool** (`client/x/dkg/keeper/dkg_cdr_fees.go:123-218`): fees distributed proportional to partial submission count — attacker's garbage submissions increment their count.

### Step 6: Attacker submits invalid partial decryption → accepted on-chain

Attacker's kernel produces partial decryption from WRONG `DistKeyShare`:

```
story-kernel/service/dkg_partial_decrypt.go:90
  priShare := distKeyShare.PriShare()  // WRONG private share (derived from different globalPubKey)

story-kernel/service/dkg_partial_decrypt.go:114
  mpc.TDH2PartialDecrypt(ownPID, privShare, pubKey, ct, label)
  // privShare doesn't correspond to consensus globalPubKey → GARBAGE output
```

On-chain validation passes:

```
client/x/dkg/keeper/dkg_handler.go:546
  bytes.Equal(pubShare, reg.PubKeyShare) → PASSES
  // reg.PubKeyShare was SET BY THE ATTACKER during finalize (Step 3, line 84)
  // This is a self-referential check — attacker writes expected value, then matches it

client/x/dkg/keeper/dkg_handler.go:553
  verifyPartialDecryptionSignature(reg.CommPubKey, ...) → PASSES
  // Only verifies ECDSA signature source, not cryptographic correctness of partial

client/x/dkg/keeper/dkg_handler.go:557
  setPartialDecryptionSubmission() → GARBAGE partial stored on-chain
```

## Impact

### 1. CDR Decryption Availability Degradation

With 3 validators (threshold=2), valid partial combinations:
- **Expected**: C(3,2) = 3 combinations, all valid → **100% availability**
- **With 1 bad dealer**: only B+C works → **33% availability** (A+B and A+C both fail)
- **With 2 bad dealers**: 0 valid combinations → **0% availability**

### 2. Silent Fault Tolerance Reduction

System reports `finalizedCount=3` but effective committee size is 2. If one honest validator goes offline, the system expects 2 remaining to handle decryption, but only 1 can produce valid partials. The degradation is **invisible** — no on-chain state distinguishes bad dealers from honest ones after finalization.

### 3. Unearned Economic Rewards

Bad dealer receives:
- UBI committee rewards: equal share with honest validators (`dkg_rewards.go:104`)
- CDR fee pool distributions: proportional to submission count (`dkg_cdr_fees.go:182-198`)
- This also **dilutes** honest validators' rewards (1/3 each instead of 1/2)

## Suggested Fix

### Option A: Gate finalization on globalPubKey consistency (Minimal)

In `Finalized()` handler (`dkg_handler.go`), after `AddGlobalPubKeyVote`, defer `finalizeDKGRegistration` until the validator's submitted `globalPubKey` matches the consensus-selected value. If the consensus `GlobalPublicKey` is already set and the validator's submission doesn't match, reject the finalization.

### Option B: Implement on-chain invalidation (Correct)

Create a deterministic invalidation path so `invalidateDealerRegistration()` is actually called:
- Record proven-bad dealer indices during justification processing
- Propose invalidation via a consensus message (e.g., include in next block's `MsgAddDkgVote`)
- Mark as `DKGRegStatusInvalidated` on-chain

### Option C: Defense in depth (Recommended)

Combine A + B, plus:
- Exclude `Invalidated` registrations from `validateParticipantsRoot()`
- Exclude `Invalidated` registrations from reward distribution (`dkg_rewards.go`)
- Verify `pubKeyShare` against consensus `publicCoeffs` in `PartialDecryptionSubmitted` instead of self-referential comparison

## Reproduction

Deploy a 3-validator devnet with one validator running a mock kernel using `WithInvalidVSSDeal()` behavior (defined in `tests/integration/dkg/mock_kernel_server.go`). Observe that:
1. After Dealing stage, attacker's registration status remains `Verified` (not `Invalidated`)
2. Attacker successfully finalizes and reaches `Finalized` status
3. `finalizedCount` is 3 instead of 2
4. Attacker receives committee rewards during Active stage
5. Attacker's partial decryption submissions pass on-chain validation
6. CDR decryption using attacker's partial + one honest partial fails
