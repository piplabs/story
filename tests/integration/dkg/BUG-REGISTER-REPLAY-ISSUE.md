# Bug: DKG register() replay attack — third party can cause duplicate index by replaying on-chain registration parameters

## Summary

DKG.sol `register()` does not check `msg.sender == enclaveInstanceData.validatorAddr` and has no duplicate registration guard. Since all registration parameters (including the SGX quote) are publicly visible on-chain, any third party can replay another validator's `register()` call. This causes the consensus layer to overwrite the victim's registration with a new index, leading to duplicate indices across validators and breaking kyber's DKG deal routing.

## Severity

**High** — Any external party (not necessarily a validator) can disrupt an entire DKG round at low cost. Causes round failure and CDR service unavailability.

## Environment

- **story**: branch `dkg/dev`, commit `dd43702`

## Root Cause

Three missing checks across two layers:

1. **DKG.sol**: `register()` does not verify `msg.sender == enclaveInstanceData.validatorAddr` — anyone can call it with any validator's data
2. **DKG.sol**: No mapping or flag to prevent duplicate registration for the same (round, validatorAddr) pair
3. **Consensus layer**: `Registered()` handler at `dkg_handler.go:24` does not check if the validator already has an existing registration for this round

## Code Evidence

### 1. Contract: no msg.sender check, no duplicate guard

**`contracts/src/protocol/DKG.sol:176-217`**:
```solidity
function register(
    bytes calldata enclaveReport,
    EnclaveInstanceData calldata enclaveInstanceData,
    uint256 startBlockHeight,
    bytes32 startBlockHash,
    bytes calldata validationContext
) external payable chargesFee whenNotPaused {
    require(enclaveReport.length != 0, "DKG: Enclave report cannot be empty");
    require(enclaveInstanceData.round != 0, "DKG: Round cannot be zero");
    require(enclaveInstanceData.validatorAddr != address(0), "DKG: Validator address cannot be empty");
    // ...
    // ← No: require(msg.sender == enclaveInstanceData.validatorAddr)
    // ← No: require(!alreadyRegistered[round][validatorAddr])

    _authenticateEnclaveReport(enclaveReport, enclaveInstanceData, expectedDataCommitment, validationContext);
    // ↑ SGX quote has no nonce — identical parameters produce identical quote → replay passes

    emit Registered(
        enclaveReport,
        enclaveInstanceData.round,
        enclaveInstanceData.validatorAddr,  // ← from parameter, NOT msg.sender
        // ...
    );
}
```

### 2. Consensus layer: no existing registration check

**`client/x/dkg/keeper/dkg_handler.go:24-93`**:
```go
func (k *Keeper) Registered(ctx context.Context, validator common.Address, ...) error {
    // Checks: round match, stage == Registration, validator in ActiveValSet
    // ← No check: "does this validator already have a registration for this round?"

    index, err := k.getNextDKGRegistrationIndex(ctx, round)  // line 51
    // index = len(all_registrations_for_round) + 1

    dkgReg := &types.DKGRegistration{
        Index: uint32(index),  // NEW index, different from original
        // ...
    }

    k.setDKGRegistration(ctx, validator, dkgReg)  // line 68
    // key = "round_validatorAddr" → OVERWRITES existing record with new index
}
```

### 3. Index assignment is count-based, not per-validator

**`client/x/dkg/keeper/dkg_registration.go:48-55`**:
```go
func (k *Keeper) getNextDKGRegistrationIndex(ctx context.Context, round uint32) (int, error) {
    registrations, err := k.getDKGRegistrationsByRound(ctx, round)
    return len(registrations) + 1, nil
    // When replaying: record count is unchanged (overwrite, not insert)
    // → same len → same "next" index for different validators
}
```

### 4. All parameters are publicly visible on-chain

The original `register()` transaction contains all parameters in calldata, and the `Registered` event emits them. An attacker can read these from any block explorer or event log.

## Attack Path

```
Setup: 3 validators (A, B, C), Registration stage

Normal flow:
  A calls register() → Registered event → consensus: A gets index=1
  B calls register() → Registered event → consensus: B gets index=2
  C calls register() → Registered event → consensus: C gets index=3

Attack (by any external address, not a validator):
  Step 1: Attacker reads A's register() parameters from on-chain tx data
  Step 2: Attacker calls DKG.sol register() with A's exact parameters
    → Contract: _authenticateEnclaveReport passes (identical SGX quote)
    → emit Registered(validatorAddr=A, ...)
    → Consensus: getNextDKGRegistrationIndex → len([A,B,C]) + 1 = 4
    → setDKGRegistration(A, {index: 4}) → overwrites A's record
    → State: A(index=4), B(index=2), C(index=3)

  Step 3: Attacker calls register() with B's exact parameters
    → Consensus: getNextDKGRegistrationIndex → len([A,B,C]) + 1 = 4  (still 3 records)
    → setDKGRegistration(B, {index: 4}) → overwrites B's record
    → State: A(index=4), B(index=4), C(index=3)

Result: A and B both have index=4. Index 1 and 2 are ghost indices.
```

## Impact

- **Duplicate index**: Two validators share the same kyber PID. Deals intended for one are routed to both, or kyber rejects them → DKG state corruption
- **Ghost indices**: Original indices (1, 2) have no owner. Kyber expects consecutive indices → deal generation/processing fails
- **Round failure**: Corrupted deal routing → VSS verification failures → insufficient finalizations → round skipped
- **Repeatable**: Attacker can do this every round, permanently preventing an active DKG committee from forming
- **Low cost**: Only DKG registration fee per replay + gas. No validator status required.
- **Amplifiable**: Attacker can replay ALL validators' registrations in one round, corrupting every index

## Relation to #703

Issue #703 identified duplicate registration caused by CL resume logic. Fix #715 added a CL-side guard (`alreadyRegistered` flag) to prevent the CL from self-replaying. However, #715 does not protect against an **external caller** replaying registration via the contract directly.

## Suggested Fix

### Option A: Contract-level guard (Recommended)

```solidity
// DKG.sol
mapping(uint32 => mapping(address => bool)) private _registered;

function register(...) external payable chargesFee whenNotPaused {
    require(
        !_registered[enclaveInstanceData.round][enclaveInstanceData.validatorAddr],
        "DKG: Validator already registered for this round"
    );
    // ... existing checks ...
    _registered[enclaveInstanceData.round][enclaveInstanceData.validatorAddr] = true;
    emit Registered(...);
}
```

### Option B: msg.sender check

```solidity
require(msg.sender == enclaveInstanceData.validatorAddr, "DKG: Caller must be the validator");
```

### Option C: Consensus-layer guard (Defense in depth)

```go
// dkg_handler.go Registered()
existing, err := k.getDKGRegistration(ctx, round, validator)
if err == nil && existing != nil {
    return errors.New("validator already registered for this round")
}
```

Option A + C together provides defense at both layers.
