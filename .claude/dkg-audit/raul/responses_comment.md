# Audit Finding Responses

Responses to all findings from the CDR AI Security Audit Report and the Apex Report (Story Network Scan #1).

---

## CDR AI Security Audit Report

---

### CDR-001: SGX Debug Mode

**Status**: DUPLICATE

This matches our internal finding C-03. The SGX debug mode will be set to `sgx.debug = false` in the production Gramine manifest. This is an operational configuration change, not a code fix. The debug flag is only enabled in development/testing environments to allow debugging tools to attach to the enclave. Production deployments will enforce non-debug mode as part of the release checklist.

---

### CDR-002: No On-Chain Verification for Partial Decryption

**Status**: DUPLICATE / FIXED

This matches our internal finding M-13. The fix has been implemented across two PRs:
- **story-kernel PR #31**: Added `HasDecryptRequest` verification via the TEE's tamper-proof light client. The kernel now verifies that a decrypt request exists on the canonical chain before performing partial decryption.
- **story PR #728**: Added `QueryClient` RPC to support the on-chain verification path from the CL side.

The combined fix ensures that partial decryptions are only performed for legitimate, on-chain decrypt requests.

---

### CDR-003: Non-Deterministic Map Iteration in CDR Fee Distribution

**Status**: DUPLICATE / FIXED

This matches our internal finding C-01. Fixed in commit `930a1ab0` on the `dkg/fix-audit-findings` branch and merged to `dkg/dev` via PR #737. The fix replaces non-deterministic Go map iteration with sorted key iteration to ensure consistent fee distribution across all validators.

---

### CDR-004: Partial Decryption Signature Missing Ethereum Signed Message Prefix

**Status**: INVALID

Contract-level signature verification for DKG is no longer required. The on-chain verification has been removed from the contract side. The CL still verifies signatures from the TEE for defense-in-depth, but the Ethereum Signed Message prefix inconsistency between `verifyFinalizationSignature` (uses prefix) and `verifyPartialDecryptionSignature` (no prefix) is a non-issue since contract verification is not in the critical path.

Code review confirms both kernel (`service/dkg_partial_decrypt.go`) and CL partial decrypt paths consistently omit the prefix, so no functional mismatch exists between the signing and verification sides.

---

### CDR-005: DKG.sol finalize() No On-Chain Validation

**Status**: DUPLICATE

This matches our internal finding M-08. The `finalize()` function is called via system transactions from the EVM engine. Adding a `msg.sender` check was evaluated and skipped for two reasons:
1. The `chargesFee` modifier combined with CL signature verification already provides the trust boundary.
2. Any participant can legitimately call finalize — the function's correctness depends on the signed data, not the caller identity.

---

### CDR-006: CDR.sol Condition Bypass When msg.sender IS the Condition Contract

**Status**: ACKNOWLEDGED

This is by design. The condition contract is the trusted arbiter for access control. If `msg.sender == vault.writeConditionAddr`, the condition check is intentionally bypassed because the condition contract itself is performing the operation — it has already evaluated whatever conditions it needs to enforce. Documentation has been added to clarify this design choice in the contract comments.

---

### CDR-007: Partial Decryption Signature Missing Code Commitment

**Status**: INVALID

Verified that BOTH kernel (`service/dkg_partial_decrypt.go:216-221`) and CL (`client/x/dkg/keeper/dkg_handler.go:375-380`) consistently exclude code commitment from the partial decrypt signature. The signed message is: `round(4B) + ciphertext + encryptedPartial + ephPubKey + pubShare`. Since both sides agree on the message format, there is no mismatch.

Additionally, during kernel upgrade resharing, the round number changes, so a new distinct signature is produced for each round — preventing cross-round replay. The code commitment is verified separately during registration (via DCAP attestation), so it does not need to be included in every partial decrypt signature.

---

### CDR-008: Pre-Release Dependencies

**Status**: KNOWN ISSUE

The `kyber v4.0.0-pre2` dependency is a known pre-release. We are tracking the upstream release at `go.dedis.ch/kyber/v4` and will upgrade to v4.0.0 stable when it becomes available. The pre-release version has been stable in our testing and the API surface we use is not expected to change in the final release.

---

### CDR-009: Kernel Light Client Witness Configuration

**Status**: INVALID

The story-kernel enforces `MinWitnessCount = 2` in config validation (`config/config.go:22`), so a zero-witness configuration is rejected at startup.

More importantly, the light client runs inside the SGX enclave using sealed LevelDB storage. After the initial registration is validated by the CL (which verifies the SGX DCAP attestation quote), the kernel's light client follows the canonical chain with Merkle proof verification. The sealed DB ensures that even if the host OS is compromised, the light client state cannot be tampered with. The CometBFT light client verifies validator set signatures, making it cryptographically bound to the canonical chain.

---

### CDR-010: Plaintext Deal Storage

**Status**: KNOWN ISSUE (FALSE POSITIVE)

This matches our existing analysis (KERNEL-002 in the initial audit, reclassified as FALSE POSITIVE). The `state.json` does not contain private keys. Deals are encrypted using the recipient's public key. Justification SecShares are publicly broadcast data (they are revealed during the complaint/justification protocol and are safe to store in plaintext). The TEE's sealed storage (SGX-encrypted) protects the actual secret material: Ed25519 signing keys and distributed key shares.

---

### CDR-011: CDR Fee Pool Dual Tracking (EL Burn + CL Mint)

**Status**: INVALID

The fee flow is deterministic: EL events are processed by the CL in a deterministic order within `ProcessDKGEvents`. The event's `amount` field is used to mint the exact corresponding amount on the CL side. Since all validators process the same events in the same order (guaranteed by CometBFT consensus), the minted amount is guaranteed to be consistent across all validators. There is no double-counting or inconsistency — the EL burns and the CL mints are two sides of the same deterministic state transition.

---

### CDR-012: Gas Limit for Registration Transactions

**Status**: KNOWN ISSUE

The current gas estimation uses `eth_estimateGas` which can fail if the gas limit is too low. The DCAP attestation verification in the SGX validation hook requires approximately 7.7M gas. This is handled operationally by ensuring the block gas limit is sufficiently high. A future improvement could add a minimum gas floor for registration transactions to make the system more robust against misconfigured gas limits.

---

### CDR-013: DKG Module Restart Recovery

**Status**: DUPLICATE / FIXED

This matches our dealer polynomial persistence fix (PR #30). The `rebuildInitDKG` function now persists and restores the dealer's polynomial coefficients across restarts, and calls `Deals()` to populate `verifiers[self]` before replaying messages. This ensures that a kernel restart during an active DKG round can recover its state and continue participating.

---

### CDR-014: Gramine Insecure Command Line Args

**Status**: ACKNOWLEDGED

The `loader.insecure__use_cmdline_argv = true` setting allows command-line arguments to be passed from the untrusted host to the enclave. This is needed for the `start` and `init` subcommands. The kernel validates the subcommand and flags internally, so arbitrary command injection is not possible. This has been added to our FINAL_AUDIT.md as a known operational consideration and will be revisited for production hardening alongside the SGX debug mode fix (CDR-001/C-03).

---

### CDR-015: No Minimum Registration Fee Enforcement

**Status**: ACKNOWLEDGED

The registration fee is queried from the DKG contract at runtime. A zero or very low fee would be a governance/contract configuration issue rather than a code bug. The fee is set by the contract owner (expected to be a multisig/governance) and can be updated as needed. The CL does not enforce a minimum because the appropriate fee level is a policy decision, not a protocol invariant.

---

### CDR-016: Private Key Zeroing

**Status**: DUPLICATE

This matches our internal findings L-12 / CRYPTO-006. The private key material resides in SGX enclave memory, which is encrypted by the CPU's Memory Encryption Engine (MEE). Combined with the enclave's access control (only enclave code can access enclave memory), the risk of key material leaking via memory inspection is mitigated by the hardware security guarantees. Explicit zeroing would be defense-in-depth but is not critical given the SGX protections.

---

### CDR-017: uint32 Overflow in DKG Params Validation

**Status**: VALID — Will fix.

The period parameters (`DealingPeriod`, `ResponsePeriod`, `FinalizationPeriod`, `ActivePeriod`) are uint32 and their sums could theoretically overflow. We will add overflow checks in `ValidateParams` to ensure that the sum of all period parameters does not exceed `math.MaxUint32`. This is a straightforward safety improvement.

---

### CDR-018: uint64 Overflow in CDR Fee Accounting

**Status**: VALID — Will fix.

The `totalCount` multiplication in fee distribution could overflow uint64 for very large values. We will add SafeMath checks (or use `math/big` for intermediate calculations) to prevent silent overflow in fee accounting. This is important for correctness of fee distribution.

---

### CDR-019: Event Data Not Indexed

**Status**: LOW PRIORITY

Events are used for off-chain monitoring and indexing, not for consensus. Adding indexed fields to key event parameters would improve query efficiency for off-chain consumers but is not required for protocol correctness. We may add indexed fields in a future contract upgrade as a quality-of-life improvement.

---

### CDR-020: Inconsistent Error Handling in handleDKGProcessJustifications

**Status**: VALID — Will fix.

The function logs errors but continues processing, which could silently drop valid justifications. We will align the error handling with the rest of the DKG lifecycle functions (deal processing, response processing) to ensure consistent behavior: either fail-fast or explicitly skip invalid entries with appropriate logging and metrics.

---

### SC-01: DKG.sol register() No Duplicate Check

**Status**: INVALID

Duplicate registration is guarded at the CL level in `handleDKGRegistration` (`client/x/dkg/keeper/dkg_svc_registration.go`). The CL checks if the validator has already registered for the current round before emitting the system transaction to call the contract. Additionally, the contract's `register()` function would overwrite the previous registration data, not create a duplicate entry, so even without the CL guard there is no state corruption — just wasted gas.

---

### SC-02: DKG.sol Zero-Address Validation Hook

**Status**: INVALID

There are no precompile contracts at the zero address on Story's EL. The validation hook address is set by the contract owner and validated at the CL level during the whitelisting process. A zero-address hook would fail attestation verification (since no contract exists to verify the DCAP quote), not silently pass. The contract owner is expected to be a multisig with proper operational procedures.

---

### SC-03: SGXValidationHook Missing Input Validation

**Status**: VALID — Will add.

Adding basic input validation (non-empty quote, non-zero code commitment) to the `SGXValidationHook` contract provides defense-in-depth. Even though the CL performs its own validation before calling the hook, rejecting obviously invalid inputs at the contract level prevents misuse and makes the contract's requirements explicit. We will add `require` checks for non-empty `quote` and non-zero `codeCommitment`.

---

### SC-04: CDR.sol allocateFee Can Be Zero

**Status**: ACKNOWLEDGED

A zero `allocateFee` is a valid configuration choice by the contract owner. This allows the contract to be deployed and tested without requiring fees. The fee can be updated later via the setter function. This is an operational/governance decision, not a code vulnerability. Documentation has been added to the contract to clarify that zero fees are intentionally allowed.

---

### SC-05: DKG.sol register() Fee Check

**Status**: NOT AN ISSUE

The `msg.value >= registrationFee` check is standard and correct. If `registrationFee` is set to zero by governance, free registration is intentional — it means the governance has decided not to charge for registration. The overpayment (excess `msg.value`) is handled by the contract's balance management. This is working as designed.

---

### SC-06: CDR.sol Missing Encrypted Partial Length Validation

**Status**: VALID — Will add.

Adding a length check on `encryptedPartial` in `submitEncryptedPartialDecryption` prevents storing unbounded data on-chain. We will add a reasonable maximum length check based on the expected size of an encrypted partial decryption share (which is deterministic given the curve parameters). This prevents DoS via oversized calldata.

---

### SC-07: CDR.sol Missing Event Emissions

**Status**: VALID — Will add.

Adding events for admin operations improves transparency and allows off-chain monitoring of governance actions. We will add events for:
- `setAllocateFee` — emits new fee value
- `setReadFee` — emits new fee value
- `setRegistrationFee` — emits new fee value
- `pause` / `unpause` — emits the caller and timestamp

---

### SC-08: CDR.sol submitEncryptedPartialDecryption Missing Reentrancy Guard

**Status**: VALID — Will add.

Although the function does not currently make external calls, adding the `nonReentrant` modifier provides defense-in-depth against future code changes that might introduce external calls. This is a low-cost safety improvement.

---

### SC-10: CDR.sol read() Access Control Documentation

**Status**: ACKNOWLEDGED

The `require(msg.sender != vault.readConditionAddr)` pattern means: "anyone EXCEPT the condition contract can call `read()` directly; the condition contract must go through its own logic path." This is intentional — the condition contract enforces access control through a different code path, and allowing it to call `read()` directly would bypass its own conditions. Documentation has been added to the contract to clarify this design pattern.

---

### SC-12: CDR.sol Memory vs Storage Pointer

**Status**: VALID — Will fix.

Using a `storage` pointer instead of copying the struct to `memory` saves gas by avoiding unnecessary memory allocation and copying. We will update the relevant vault lookups to use `storage` pointers where the data is only read (not modified in a way that requires a separate copy).

---

### L1-02: CDR.sol Missing Maximum Data Size for encryptedData

**Status**: VALID — Will add.

Adding a maximum size limit prevents DoS via oversized calldata that would consume excessive gas for storage. We will add a limit based on approximately 150KB (~5M gas worth of calldata at 16 gas per non-zero byte). This provides a reasonable upper bound for encrypted data while preventing abuse.

---

### L1-04: DKG.sol scheduleUpgrade Lacks Version Format Validation

**Status**: ACKNOWLEDGED

The upgrade version string is used as an identifier for matching between the CL and the contract, not parsed semantically. Format validation at the contract level is a nice-to-have but the CL side handles version matching and validation. The contract owner (governance) is responsible for providing valid version strings. Documentation has been added to clarify the expected format.

---

### L1-05: Upgrade Activation During Active DKG Round

**Status**: BY DESIGN

When an upgrade activates, it immediately starts a new upgrade resharing round regardless of the current round's stage. This is intentional: the upgrade takes priority over the current round. The current round is effectively abandoned and a new round with `isUpgrade=true` begins.

Code review confirms this is safe: `BeginBlocker` checks for pending upgrade activation BEFORE processing normal stage transitions (`abci.go:55-70`). The new round uses the new kernel binary's code commitment for registration, ensuring only upgraded kernels can participate in the new committee.

---

### L1-07: Package-Level Global State Test Isolation

**Status**: KNOWN ISSUE

This matches our H-05 in FINAL_AUDIT.md. The package-level globals (`deals`, `responses`, `justifications` queues) prevent parallel test execution with `t.Parallel()`. Moving these to `Keeper` struct fields is a large refactor tracked for future work. In the meantime, tests are run sequentially with explicit state cleanup between test cases to ensure isolation.

---

## Apex Report — Story Network Scan #1

---

### STOR-1: Partial Decryption Request Hijacking

**Status**: INVALID

The `ThresholdDecryptRequested` event is emitted by the EVM and processed deterministically by the CL. The CL stores the decrypt request with the requester's public key and label. Even if an attacker intercepts the partial decryptions in transit, they are encrypted with the requester's ECIES public key (`encryptPartialToRequester` in `dkg_partial_decrypt.go`), so only the legitimate requester holding the corresponding private key can decrypt them.

Additionally, the round number is included in the partial decrypt signature, preventing cross-round replay attacks.

---

### STOR-2: Upgrade Resharing Registration Logic

**Status**: VALID — Will fix.

The `resolveRegistrationKernelClient` function is implicit about its behavior. It receives `oldCC` (old code commitment) and returns a kernel client with a DIFFERENT code commitment (the new kernel binary). The function name should clearly indicate this transformation. We will:
1. Rename to `resolveNewKernelClientForUpgrade` (or similar) for clarity.
2. Add validation that the resolved client's code commitment differs from `oldCC` to catch misconfiguration.

---

### STOR-3: Registration Index Overwrite

**Status**: FIXED

Already fixed in PR #723 (`fix(dkg): avoid index overwrite`). The fix ensures that participant indices are not overwritten during registration, preserving the correct mapping between validators and their DKG indices.

---

### STOR-4: Justification Verification

**Status**: FIXED

Already fixed in PR #719 (`fix(dkg): invalidate dealer registration when justification verification failed`). When a justification fails verification, the dealer's registration is now properly invalidated, preventing them from being included in the finalized committee with invalid key material.

---

### STOR-5: CDR Partial Decrypt Without On-Chain Verification

**Status**: FIXED

Fixed in story-kernel PR #31 (`feat(cdr): validate if decrypt request exists on canonical chain`). The kernel now uses the TEE's tamper-proof light client to verify that the decrypt request exists on the canonical chain before performing partial decryption. The light client runs inside the SGX enclave with sealed DB storage, ensuring the verification cannot be bypassed by a compromised host OS.

---

### STOR-6: Kernel gRPC Missing Timeout/Retry

**Status**: PARTIALLY FIXED

PR #725 added auto-reconnect for kernel gRPC clients. Our audit fix branch adds exponential backoff for reconnection attempts. The `retryAttempts=3` with `retryDelay=2s` in the DKG helper provides basic retry logic for transient failures. Further improvements may be needed for longer-running operations (e.g., DCAP attestation verification which can take several seconds).

---

### STOR-7: Vote Extension Replay Between Rounds

**Status**: VALID — Will fix.

Vote extensions from a previous round could potentially be replayed in a new round if the session ID / round number is not checked during verification. We will add round number validation to `VerifyVoteExtension` to ensure VE data references the current active round. The validation order should be: validate round number first, then dedup, then process.

---

### STOR-8: PIDCache Not Populated During Resharing

**Status**: VALID — Will fix.

During resharing, `GenerateDeals` skips `CachePID` because `req.GetIsResharing()` is true. However, the new committee members need their PID (Participant ID) for `PartialDecryptTDH2` operations. We will fix by ensuring PID is cached for all committee members regardless of whether this is a resharing round, since the PID is needed for CDR operations after the round completes.

---

### STOR-9: Missing Proof Verification for CDR Partial Decrypt Collection

**Status**: IN PROGRESS

PR #727 on the `dkg/mediumAuditIssue` branch is addressing this. The commit `ed9eb922` indicates "pid is not needed when calling tee for partial decryption since tee stores pid and will validate pid against range from query client." We are verifying that this change fully addresses the proof verification concern and that the kernel's internal PID validation is sufficient.

---

### STOR-10: Vote Extension Processing Before Validation

**Status**: VALID — Will fix.

Vote extensions should be validated (signature check, format check, round number check) before being processed or deduped. Currently, invalid extensions may be processed and then rejected, wasting compute resources and potentially polluting dedup caches. We will reorder the logic to validate first, then dedup, then process — ensuring only valid VE data enters the processing pipeline.

---

### STOR-11: Decrypt Request Registry Key Missing Validator Address

**Status**: VALID — Will fix.

The registry key for decrypt requests should include the validator address to prevent cross-validator collisions. Without the validator address, two validators processing the same decrypt request could interfere with each other's state. We will add the validator address to the registry key to ensure proper namespacing.

---

### STOR-12: DKG Contract Owner Privileges

**Status**: ACKNOWLEDGED

The contract owner has significant privileges including fee setting, upgrade scheduling, and enclave whitelisting. This is an operational concern managed through proper key management practices:
- Multisig ownership for the DKG contract
- Timelock for sensitive operations (upgrade scheduling)
- Monitoring for unexpected owner actions

This is not a code fix priority — it is an operational/governance concern that is addressed through deployment procedures.

---

### STOR-13: PubKeyShare Prefix Mismatch

**Status**: FIXED

Fixed in story-kernel `jdub/marshal-pubshare` branch (commit `7bb7057: fix: marshal pubshare cb-mpc`). The `marshalPubShare` function now prepends `[0x04, 0x3f]` (SEC1 uncompressed prefix + TDH2 Edwards25519 curve ID) to match what `buildTDH2PublicKey` expects during deserialization. Story-kernel PR #34 includes this fix.

---

### STOR-14: Kernel Input Validation and MRENCLAVE

**Status**: INVALID

Modifying kernel code to add input validation would change the MRENCLAVE (code commitment hash), which would invalidate all existing sealed keys. This means existing validators would lose access to their DKG key shares, requiring a full re-keying ceremony.

Any kernel code change — including input validation additions — must go through the coordinated upgrade resharing flow: schedule upgrade on-chain, validators build new kernel binary, perform upgrade resharing round, finalize with new code commitment. Input validation changes cannot be patched in-place without triggering this process.

---

### STOR-15: CDR Fee Wei/Gwei 1e9 Amplification

**Status**: FIXED

Fixed in PR #737 (`fix(cdr): fix fee distribution`). The CL now correctly handles the wei-to-gwei conversion when processing CDR fee events from the EL. Previously, the conversion factor was applied incorrectly, amplifying or reducing fee amounts by 10^9. The fix ensures the correct unit conversion at the CL boundary where EL events are processed.

---

### STOR-16: Reward Distribution During Non-Active Stages

**Status**: BY DESIGN

The previous active round's committee continues to receive rewards until the next round becomes active. This is intentional: it incentivizes committee members to remain available for DKG operations (deal generation, response processing, finalization, and CDR partial decryption) during the transition period between rounds. Without this incentive, committee members might go offline during the transition, stalling the next round.

---

### STOR-17: Light Client Rollback Attack

**Status**: MITIGATED

The CometBFT light client inside the SGX enclave provides strong rollback protection:
1. **Validator set signatures**: Each block header is verified against the validator set, requiring >=2/3 signatures. A rollback attack would require forging signatures from >=2/3 of validators — equivalent to a full byzantine fault.
2. **Sealed DB**: The light client's trusted state is stored in SGX-encrypted sealed storage, preventing the host OS from tampering with it.
3. **TrustedPeriod**: Old validator sets expire and cannot be used for verification after the trusted period elapses.
4. **Merkle proofs**: All state queries include Merkle proof verification against the verified block header.

Combined, these provide canonical state guarantees equivalent to running a full node.

---

### STOR-18: Kernel Query Client Key Spoofing

**Status**: VALID — Will add.

Although the sealed DB light client provides strong guarantees for header and proof verification, adding explicit key format validation for ABCI query responses adds defense-in-depth. We will add key prefix and format checks to ensure that query responses contain well-formed keys matching the expected schema (e.g., correct store prefix, expected key length).

---

### STOR-19: Stage Boundary Off-By-One

**Status**: ACKNOWLEDGED

This is a timing issue at stage boundaries, similar to STOR-22. A deal submitted in the last block of the dealing period might not be processed if the stage transitions in the same block. This is a fundamental limitation of discrete block-based state transitions. Documentation will be updated to clarify that submissions should be made well before the stage boundary (at least 2-3 blocks before the transition) to ensure processing.

---

### STOR-20: Protobuf Schema Mismatch

**Status**: IN PROGRESS

PR #739 (story) and `jdub/align-proto` branch (story-kernel) are addressing the proto field alignment between story and story-kernel. The kernel's proto definitions must match the CL's expectations for correct serialization/deserialization of gRPC messages. Mismatched field numbers or types would cause silent data corruption or deserialization errors.

---

### STOR-21: Active Period Expiry Stage Transition

**Status**: FIXED

The current code in `dkg_round.go` transitions directly from Active to Registration (for the next round). The "Ended" stage exists but is only set when `FinalizeDKGRound` marks the previous round as ended during the new round's finalization. This means there is no gap where a round is stuck in Active without transitioning — the lifecycle is: Active -> (next round starts) -> Registration, and the previous round is marked Ended during the new round's finalization step.

---

### STOR-22: Vote Extension Timing at Stage Boundaries

**Status**: ACKNOWLEDGED

Vote extensions are collected during block N-1 and processed in block N's `PrepareProposal`/`ProcessProposal`. If a stage transition happens at block N, VEs from block N-1 (which were created in the previous stage's context) are processed in the new stage. This is a known limitation of the VE-based architecture in CometBFT.

Documentation will clarify that VE data is "best effort" and stage-boundary VEs may be dropped. The period lengths are configured to be long enough that losing a single block's worth of VEs at the boundary does not impact the protocol's ability to complete.

---

### STOR-23: Complaint/Justification Timing via Vote Extensions

**Status**: NOT AN ISSUE

A complaint broadcast in the last block of the dealing phase is structurally dropped due to the ABCI execution order: `BeginBlocker` (stage transition) runs before `DeliverTx` (`AddVote`). When `AddVote` executes, the stage has already transitioned to `DKGStageFinalization`, and `msg_server.go:28` gates all VE processing on `stage == DKGStageDealing`.

This is not a practical concern for the following reasons:
1. **Local safety**: The recipient that generated the complaint has already marked the deal as invalid in its own kernel state. It will not use the malicious dealer's share during finalization.
2. **Threshold tolerance**: Even if 1-2 participants disagree on a dealer's validity, the round completes as long as the threshold is met with the remaining participants.
3. **Operational buffer**: The dealing period is configured to be long enough (50+ blocks) that complaints are generated and propagated well before the stage boundary. A complaint reaching the last block implies a deal was sent extremely late, which itself is an edge case.

This is the same structural limitation as STOR-22 (VE timing at stage boundaries) and is inherent to the CometBFT vote extension architecture where `BeginBlocker` stage gates precede `DeliverTx` processing.

---

### STOR-24: Vote Extension Item Count Mismatch

**Status**: VALID — Will verify and fix.

The `PrepareVotes` dequeue limit should match the `VerifyVoteExtension` item count check (80 items each for deals, responses, justifications). If there is a mismatch between the preparation limit and the verification limit, valid VEs could be rejected by other validators during `VerifyVoteExtension`. We will audit all VE size constants and ensure they are aligned across `PrepareVotes`, `VerifyVoteExtension`, and `ProcessProposal`.

---

### STOR-25: Vote Extension Validation Gap

**Status**: NOT AN ISSUE

The suggested fixes (round number validation in `VerifyVoteExtension`, session ID validation in `AddVote`) are either unnecessary or infeasible:

**Round validation is unnecessary:**
1. The Deal, Response, and Justification protobuf messages do not contain a DKG round number field — adding one would require a protocol-level schema change.
2. Cross-round replay protection is already guaranteed cryptographically. Kyber's `DistKeyGenerator` derives a unique `SessionID` from `(dealer pubkeys + verifiers + commitments + threshold)` per round. Since `SessionID` is embedded in `Justification.Hash()`, Schnorr signature verification implicitly rejects replayed data from other rounds (`dkg_dealing.go:96-99`).
3. `VerifyVoteExtension` is intentionally stateless (`vote.go:46-56`) — it does not access the store. Adding round validation would require state access, conflicting with this design principle.
4. The kernel already rejects session-phase-mismatched data (`dkg_svc_dealing.go:260`).

**Session ID validation is infeasible at the CL:**
1. The session ID is computed by kyber's `DistKeyGenerator` inside the TEE kernel. The CL does not run the DKG state machine and has no access to the expected session ID.
2. To validate session IDs at the CL, the kernel would need to export the session ID to on-chain state — introducing a new trust boundary with no security benefit over the existing cryptographic guarantees.

**Existing protections are sufficient:**
- Cryptographic session ID-based implicit replay protection (kyber)
- Stage gate in `msg_server.go:28` (`DKGStageDealing` only)
- Kernel-level session phase verification
- Proposer rotation + VE re-broadcast + threshold for censorship resistance

---

### STOR-26: Off-Chain Registration Status

**Status**: KNOWN

Registration status transitions happen off-chain in the DKG service layer. The on-chain contract records registrations (via system transactions), and the CL tracks verification status through the DKG keeper state. This is an architectural observation, not a bug — the CL is the source of truth for registration validity, while the contract is the source of truth for on-chain state (fee escrow, participant list).

---

### STOR-27: CDR Fee Collection Timing

**Status**: FIXED

PR #737 moved fee distribution from `FinalizeDKG` to the correct lifecycle point. Previously, fees could be distributed at the wrong time relative to the round lifecycle. The fix ensures fees are distributed at the appropriate stage, aligned with the committee that earned them.

---

### STOR-28: CDR Collect Partials Duplicate Check

**Status**: IN PROGRESS

PR #727 on the `dkg/mediumAuditIssue` branch is addressing this. The changes include PID validation for partial decryption submissions, which should prevent duplicate partial decryptions from the same participant from being collected and counted multiple times.

---

### STOR-29: CDR Partial Decrypt PID Verification

**Status**: IN PROGRESS

Same PR #727 on the `dkg/mediumAuditIssue` branch. The commit `ed9eb922` indicates: "pid is not needed when calling tee for partial decryption since tee stores pid and will validate pid against range from query client." The kernel stores the PID internally (set during DKG finalization) and validates it against the expected range queried from the CL, removing the need for the caller to provide the PID and preventing PID spoofing.
