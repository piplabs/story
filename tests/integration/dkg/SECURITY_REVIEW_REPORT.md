# DKG Consensus Layer Security Review Report

## Summary

**6 attack vectors analyzed · 🔴 2 confirmed bugs filed · ✅ 2 fixed · 🟡 4 assessed as low/no risk**

Branch: [`dkg/dev`](https://github.com/piplabs/story/tree/dkg/dev) · Commit: [`dd43702`](https://github.com/piplabs/story/commit/dd43702) · Scope: `client/x/dkg/keeper/`, `client/x/evmengine/keeper/cdr.go`, `contracts/src/protocol/DKG.sol`, `contracts/src/protocol/CDR.sol`

### Review Approach

This security review analyzed the DKG module from an attacker's perspective, tracing every economic reward path and state transition to identify exploitable vulnerabilities. The review covered the full stack: Solidity contracts (DKG.sol, CDR.sol), consensus layer keeper logic (ABCI hooks, vote extensions, event handlers), and the trust boundary with story-kernel (TEE/SGX). Each potential attack was validated against the actual code to confirm exploitability — distinguishing real bugs from theoretical concerns.

### Analysis Scope

The review focused on 6 primary attack surfaces:

- **Complaint/Justification flow**: Dealer invalidation after VSS verification failure, justification broadcast and verification pipeline
- **Finalization integrity**: GlobalPubKey voting, registration status gating, pubKeyShare self-referential validation
- **Economic extraction**: CDR fee pool distribution mechanism, UBI committee reward distribution, partial decryption submission incentives
- **Registration replay**: Contract-level and consensus-level duplicate registration protection, index assignment consistency
- **TEE trust boundary**: Post-registration attestation, SGX sealing guarantees, mock kernel substitution
- **Protocol liveness**: Stage transition timing, ActiveValSet freshness, decrypt request registry management

### Key Finding

The two confirmed bugs share a common pattern: **the consensus layer trusts data written by validators themselves** (self-referential validation) rather than deriving expected values from consensus state. The register replay bug additionally lacked basic idempotency guards at both the contract and consensus layers.

---

## Issues

| **Issue** | **Status** | **Priority** | **Description** | **Fix PR** |
| --- | --- | --- | --- | --- |
| [#717](https://github.com/piplabs/story/issues/717) | ✅ Fixed | Critical | Bad dealer not invalidated after justification VSS failure — can finalize, join committee, submit invalid partials, receive rewards | [#719](https://github.com/piplabs/story/pull/719) |
| [#721](https://github.com/piplabs/story/issues/721) | ✅ Fixed | High | Register replay attack — third party replays on-chain registration parameters to corrupt validator indices, causing duplicate PID and round failure | [#723](https://github.com/piplabs/story/pull/723) |

---

## Detailed Findings

### 1. Bad Dealer Bypass Invalidation (Critical — Fixed)

**Issue**: `invalidateDealerRegistration()` was defined at `dkg_justification.go:206` but never called from production code. When a justification's VSS verification confirmed a deal was invalid, the dealer's status remained `Verified` — allowing it to finalize, join the active committee, submit garbage partial decryptions, and receive UBI + CDR fee rewards.

**Root cause**: Justification verification ran in an async goroutine that couldn't modify on-chain state. The `!valid` branch at `dkg_svc_dealing.go:463` only logged and continued.

**Fix** ([#719](https://github.com/piplabs/story/pull/719)): Moved justification verification (Schnorr + VSS) into `ProcessJustifications` within the FinalizeBlock context. When VSS fails, `invalidateDealerRegistration` is now called synchronously, marking the dealer as `Invalidated` on-chain.

**Impact before fix**: A single malicious validator could degrade CDR decryption availability from 100% to 33% (3 validators, threshold=2) while earning equal committee rewards.

### 2. Register Replay Attack (High — Fixed)

**Issue**: DKG.sol `register()` did not check `msg.sender == enclaveInstanceData.validatorAddr` and had no duplicate registration guard. The consensus layer `Registered()` handler also lacked an existence check. Any third party could replay validators' registration parameters (visible on-chain) to overwrite their indices, creating duplicate PIDs.

**Root cause**: No idempotency guard at either layer. `getNextDKGRegistrationIndex` computed `len(registrations) + 1`, but `setDKGRegistration` used `round_address` as key — overwrites didn't change the count, producing the same "next index" for different replay victims.

**Fix** ([#723](https://github.com/piplabs/story/pull/723)):
- Contract: Added `require(enclaveInstanceData.validatorAddr == msg.sender)`
- Consensus: Added `hasDKGRegistration()` check before index assignment

**Impact before fix**: Any EOA (non-validator) could cause entire DKG rounds to fail repeatedly at the cost of only the DKG registration fee per replay.

### 3. Non-Finalized Validator Partial Submission (Assessed — TEE Protected)

**Finding**: `PartialDecryptionSubmitted` at `dkg_handler.go:541` does not check `reg.Status == DKGRegStatusFinalized`. For non-finalized validators, `reg.PubKeyShare` is `nil`, and `bytes.Equal(nil, nil)` returns `true` in Go, theoretically allowing empty pubShare submissions to pass validation.

**Assessment**: Not exploitable in practice. The partial decryption requires a valid ECDSA signature from `commPubKey`, whose private key is sealed inside the SGX enclave. A non-finalized validator's kernel cannot produce a partial (no `DistKeyShare`), and the sealed key cannot be extracted to sign externally.

**Recommendation**: Add `reg.Status == Finalized` check as defense-in-depth. Low priority — SGX sealing is the effective protection.

### 4. CDR Fee Farming by Legitimate Validator (Assessed — Not Exploitable)

**Finding**: CDR fee pool is distributed proportional to `CDRPartialSubmitCount`. Initial analysis suggested a finalized validator could self-generate `read()` requests to inflate their count.

**Assessment**: Not exploitable. All validators' kernels automatically respond to every `ThresholdDecryptRequested` event. Attacker-generated read requests increase ALL validators' counts equally, maintaining the same distribution ratio. Net result: attacker loses `readFee × N` with no proportional gain.

### 5. Stage Transition Timing Attack (Assessed — Not Exploitable)

**Finding**: `shouldTransitionStage` at `dkg_round.go:10` is purely time-based with no deal completion check. A validator could delay deal broadcast to the end of the Dealing stage.

**Assessment**: Not exploitable by a single validator. In Pedersen DKG with threshold=2 and 3 validators, the two honest validators exchange deals among themselves and finalize independently. The attacker only prevents their own participation. Round succeeds with `finalizedCount >= threshold`.

### 6. ActiveValSet Snapshot Staleness (Assessed — Design Tradeoff)

**Finding**: `ActiveValSet` is snapshotted once at round initiation (`dkg_initialization.go:44`) and not refreshed during the ~3 day round. A validator jailed mid-round can still participate.

**Assessment**: Low practical risk. A jailed validator's TEE still holds valid sealed keys, but the validator has already been economically penalized (slash). Refreshing `ActiveValSet` mid-round would introduce consensus complexity and potential liveness issues. Accepted as a design tradeoff.

---

## Attack Surface Summary

| # | Attack Vector | Exploitable? | Economic Impact | Status |
|---|--------------|-------------|-----------------|--------|
| 1 | Bad dealer bypass invalidation | ✅ Yes | Unearned rewards + availability degradation | **Fixed** (#719) |
| 2 | Register replay (third party) | ✅ Yes | Round failure (DoS) | **Fixed** (#723) |
| 3 | Non-finalized validator partial | ❌ TEE protected | Would allow unearned CDR fees | Defense-in-depth suggestion |
| 4 | CDR fee farming (self-read) | ❌ Not profitable | Equal count increase for all | No action needed |
| 5 | Late deal broadcast | ❌ Single validator insufficient | Cannot block round | No action needed |
| 6 | Stale ActiveValSet | ❌ Low practical risk | Jailed validator continues | Design tradeoff |

---

## Post-Registration TEE Trust Model

During the review, we assessed whether a validator could substitute a mock kernel after registration. Finding:

- **Registration**: One-time SGX remote attestation via DCAP (verified on-chain by `_authenticateEnclaveReport`)
- **Post-registration**: Zero ongoing attestation. Each kernel RPC only performs a self-check (`enclave.ValidateCodeCommitment`) which a mock kernel could bypass
- **Effective protection**: SGX sealing binds keys to MRENCLAVE. A different binary cannot unseal `commPubKey` or `dkgPubKey` private keys, preventing signature forgery

**Conclusion**: The system's post-registration security relies entirely on SGX sealing, not on continuous attestation. This is a single-point trust assumption — if SGX sealing is compromised (e.g., hardware side-channel attacks), the entire DKG security model breaks. Documented as an architectural observation, not a code bug.
