# STOR-29: DKG Session Created with Empty CodeCommitment Never Recovers

## Summary

When a DKG session is created but the first `GenerateAndSealKey` call fails (e.g., kernel not yet connected, DCAP not deployed), the session persists with `CodeCommitment=""`. Subsequent blocks detect "Session already exists" and skip re-creation, but each retry generates a new key pair with a new quote. The new quote's data commitment doesn't match the registration parameters from previous attempts, causing permanent `Code commitment does not match` errors until the story process is manually restarted.

## Severity

**HIGH** — Causes a validator to be permanently unable to register for DKG rounds until story is restarted. In a 3-validator devnet with `min_req_registered=3`, this blocks all DKG progress.

## Affected Code

- **`client/x/dkg/keeper/dkg_svc_registration.go:38-69`** — Session creation with empty CodeCommitment
- **`client/x/dkg/keeper/dkg_svc_registration.go:148-149`** — Late CC assignment from kernel client
- **`client/x/dkg/keeper/dkg_svc_registration.go:156-175`** — GenerateAndSealKey regenerates keys on every retry

## Reproduction

1. Start a devnet with 3 validators
2. Ensure kernel is slow to connect or DCAP deployment is in progress during the first DKG round
3. Observe val1 logs:
   ```
   INFO Created DKG session code_commitment="" round=N phase=Initializing
   ERRO Failed to generate the sealed key err="..."
   ```
4. On subsequent blocks:
   ```
   INFO Session already exists with the code commitment and round, skip creating a new session code_commitment="" round=N
   INFO GenerateAndSealKey call to kernel client round=N
   ERRO Failed to call register method err="...Code commitment does not match"
   ```
5. This loops every block for the entire registration period
6. Round skips due to `verified_count < min_req_registered`
7. Next round: same problem (new session created with CC="" again)
8. **Only fix**: `sudo systemctl restart story`

## Root Cause Analysis

### Session Creation (line 38-64)

```go
session := &types.DKGSession{
    Round: dkgNetwork.Round,
    // CodeCommitment is "" (zero value)
}
// Only set for upgrade rounds with old CC
if dkgNetwork.IsUpgrade && len(oldCC) > 0 { ... }
k.stateManager.CreateSession(ctx, session) // Persists CC=""
```

Session is created **before** `GenerateAndSealKey` is called. If the kernel call fails, the session remains in the state manager with `CodeCommitment=""`.

### Retry Path (subsequent blocks)

```go
// Line ~45: checks if session exists
existing := k.stateManager.GetSession(round)
if existing != nil {
    log.Info("Session already exists...skip creating a new session", "code_commitment", "")
    // Uses the existing session with CC=""
}
```

The existing session is reused, but CC is still empty.

### GenerateAndSealKey Regeneration

Each retry calls `GenerateAndSealKey` which:
1. Generates a **new** key pair (or loads existing sealed keys)
2. Generates a **new** SGX quote with data commitment = `keccak256(validatorAddr || round || height || hash || dkgPubKey || commPubKey)`
3. The dkgPubKey/commPubKey may differ from previous attempts if keys were regenerated

The contract's `validateReport` checks:
- MRENCLAVE at quote[112:144] == on-chain code commitment → **PASS** (same kernel)
- Data commitment at quote[368:400] == expected → **FAIL** (different keys each time)

## Suggested Fix

Option A: Delete session on GenerateAndSealKey failure so next block creates a fresh one:

```go
if err := retry(ctx, func(ctx context.Context) error {
    resp, err = client.GenerateAndSealKey(ctx, req)
    return err
}); err != nil {
    k.stateManager.DeleteSession(ctx, session) // Clean up bad session
    return errors.Wrap(err, "kernel client GenerateAndSealKey request failed")
}
```

Option B: Don't persist session until GenerateAndSealKey succeeds — move `CreateSession` after the kernel call returns successfully.

## Related Issues

- **STOR-28**: Decrypt worker timeout (fixed in PR #727)
- **Deal retry infinite loop**: `retry()` in `helpers.go` doesn't distinguish `InvalidArgument` (non-retryable) from `Internal` (retryable), causing cached deals to loop every block during dealing phase

## Environment

- **story**: `release/1.6` (commit `8307f9c2`)
- **story-kernel**: `release/0.1` (commit `ba2f33c`)
- **Trigger**: Reset devnet → V160 upgrade → DCAP deployment in progress → first registration attempt fails → session stuck with CC=""
