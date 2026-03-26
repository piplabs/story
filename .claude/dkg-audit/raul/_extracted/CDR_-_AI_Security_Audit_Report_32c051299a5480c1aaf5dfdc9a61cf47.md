# CDR - AI Security Audit Report

**Commits audited:**

- `piplabs/story` @ [`11307a5`](https://github.com/piplabs/story/tree/11307a50d2862d325898f2239e3c22bafd3a1565) (branch `dkg/dev`)
- `piplabs/story-kernel` @ [`6eb8e93`](https://github.com/piplabs/story-kernel/tree/6eb8e9345407a58f1cd49ff61a1d7f193a72b178) (branch `main`)

## Executive Summary

**Scope:** Full CDR (Confidential Data Rails) security audit across all three layers — smart contracts (CDR.sol, DKG.sol, SGXValidationHook.sol), TEE kernel (story-kernel: Go gRPC service in Intel SGX via Gramine), and L1 consensus module (x/dkg Cosmos SDK module).

**Methodology:** 8-phase CDR audit protocol with delegated layer audits (web3-security-auditor, tee-security-auditor, L1 module review) plus CDR-specific cryptographic protocol analysis (Pedersen DKG, TDH2), cross-layer boundary security, and vulnerability taxonomy coverage.

**Architecture:** N validators (dynamic, bonded set), threshold t = ceil(N * 667/1000) (~66.7%), SGX enclave type whitelist. Three-layer coordination: EVM contracts emit events -> CL processes deterministically -> TEE performs crypto operations via gRPC.

**Key Risk:** The CDR system's security relies on the composition of TEE isolation with threshold cryptography. Neither alone is sufficient — the system assumes fewer than t enclaves are compromised simultaneously.

### Finding Summary

| Severity | Count | Description |
| --- | --- | --- |
| **P0** | 1 | SGX debug mode enabled in manifest template |
| **P1** | 3 | No on-chain verification before partial decrypt, consensus-breaking map iteration, partial decrypt sig prefix mismatch |
| **P2** | 10 | finalize() gaps, supply chain, witness enforcement, condition bypass, replay silence, CDR fee mint, registry key collision, cache growth, broad Gramine allowed_files, insecure Gramine flags |
| **P3** | 3 | Vault unreadable when readConditionAddr=0, session files 0644, uint64->int64 overflow |
| **P4** | 1 | UUID uint32 overflow |
| **P5** | 2 | DKG state plaintext, GC secret copy risk |
| **Total** | **20** |  |

---

## Findings by Severity

### CDR-001: Gramine SGX Debug Mode Enabled in Manifest Template

- **Severity:** P0 (Key Compromise)
- **Layer:** TEE
- **Status:** Open
- **Location:** [`story-kernel/story-kernel.manifest.template:51`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/story-kernel.manifest.template#L51)
- **Description:** `sgx.debug = true` in the Gramine manifest template. Debug mode disables SGX memory encryption protections and allows debugger attachment.
- **Impact:** Any attacker with host access can dump enclave memory and recover all DKG private key shares, communication keys, and DistKeyShares, enabling full shared secret reconstruction (bypasses threshold property).
- **Recommendation:** Set `sgx.debug = false` for production. Add CI check rejecting `sgx.debug = true` in release artifacts.

### CDR-002: TEE Does Not Verify On-Chain Read Request Before Partial Decryption

- **Severity:** P1 (Confidentiality Breach)
- **Layer:** Cross-Layer (TEE / Contract)
- **Status:** Open
- **Location:** [`story-kernel/service/dkg_partial_decrypt.go:42-43`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_partial_decrypt.go#L42-L43)
- **Description:** `PartialDecryptTDH2` performs partial decryption on any ciphertext via gRPC without verifying a `CDR.read()` transaction exists on-chain. The TODO at line 42 acknowledges this. The TEE only verifies the round matches the latest active network.
- **Impact:** A compromised CL node or attacker with gRPC access can decrypt arbitrary ciphertexts without read conditions being enforced. Enables pre-decryption probing and MEV extraction from encrypted transactions.
- **Recommendation:** Verify via light client that a `VaultRead` event was emitted for `(label, requesterPubKey, ciphertext)` before performing partial decryption.

### CDR-003: Non-Deterministic Map Iteration in CDR Reward Distribution (Consensus Halt Risk)

- **Severity:** P1 (Integrity Violation — Consensus Safety)
- **Layer:** L1 Module
- **Status:** Open
- **Location:** [`story/client/x/dkg/keeper/dkg_cdr_fees.go:182`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_cdr_fees.go#L182)
- **Description:** `distributeCDRRewardPool` iterates a Go `map[string]uint64` (`counts`) without sorting keys. Go map iteration is non-deterministic. If any `SendCoinsFromModuleToAccount` call fails mid-iteration, different nodes may have distributed to different validator subsets, causing **consensus divergence and chain halt**.
- **Impact:** Chain halt. Non-deterministic state transitions in FinalizeBlock.
- **Recommendation:** Sort map keys before iterating, matching the pattern in [`distributeRewardsFromModule`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_rewards.go#L74).

### CDR-004: Partial Decryption Signature Verification Missing Ethereum Signed Message Prefix

- **Severity:** P1 (Integrity Violation — Cryptographic)
- **Layer:** L1 Module
- **Status:** Open
- **Location:** [`story/client/x/dkg/keeper/dkg_handler.go:382`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_handler.go#L382)
- **Description:** `verifyPartialDecryptionSignature` computes `respHash = crypto.Keccak256(encoded)` then calls `crypto.SigToPub(respHash, sig)` **without** the Ethereum signed message prefix (`\\x19Ethereum Signed Message:\\n32`). In contrast, [`verifyFinalizationSignature`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_handler.go#L325-L327) **does** apply this prefix. The kernel's [`signPartialDecryptResponse`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_partial_decrypt.go#L214) also signs without the prefix, so kernel and CL are consistent — but they differ from the finalization path.
- **Impact:** If this asymmetry is unintentional, one path has a signature verification bug. If intentional, the inconsistency means partial decryption signatures are not EIP-191 compliant and could be confused with raw hash signatures.
- **Recommendation:** Align both signature paths to use the same prefix convention. Add comments documenting the intentional difference if it is by design.

### CDR-005: DKG.sol `finalize()` Has No On-Chain Signature, ParticipantsRoot, or Sender Validation

- **Severity:** P2 (Integrity Violation)
- **Layer:** Contract
- **Status:** Open
- **Location:** [`story/contracts/src/protocol/DKG.sol:229-260`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/DKG.sol#L229-L260)
- **Description:** `finalize()` checks only that `enclaveType` is whitelisted and fields are non-empty. No signature verification, no `participantsRoot` validation, and critically **no `validatorAddr == msg.sender` check** (unlike [`register()` at line 186](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/DKG.sol#L186)). Anyone can call `finalize()` for any validator address.
- **Impact:** The CL validates all of this, but the contract provides zero protection. If CL transaction filtering is bypassed, arbitrary finalization data could be injected.
- **Recommendation:** Add `require(validatorAddr == msg.sender)` at minimum. Consider on-chain signature verification as defense-in-depth.

### CDR-006: CDR.sol Condition Bypass When msg.sender IS the Condition Contract

- **Severity:** P2 (Integrity Violation)
- **Layer:** Contract
- **Status:** Open
- **Location:** [`story/contracts/src/protocol/CDR.sol:152`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L152), [`CDR.sol:191`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L191)
- **Description:** If `msg.sender == vault.writeConditionAddr` or `msg.sender == vault.readConditionAddr`, the condition check is bypassed entirely. The condition contract itself gets unconditional access.
- **Impact:** Compromised, upgradeable, or buggy condition contracts gain unrestricted vault access.
- **Recommendation:** Document as intentional or remove bypass.

### CDR-007: Partial Decryption Signature Missing Code Commitment

- **Severity:** P2 (Integrity Violation)
- **Layer:** Cross-Layer
- **Location:** [`story-kernel/service/dkg_partial_decrypt.go:207-214`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_partial_decrypt.go#L207-L214)
- **Description:** Signature covers `round || ciphertext || encryptedPartial || ephPubKey || pubShare` but not `codeCommitment`. During kernel upgrades with two binaries, signatures from different enclave builds are indistinguishable.
- **Recommendation:** Add `codeCommitment` to signed payload.

### CDR-008: Pre-Release kyber v4 and Forked cb-mpc-go Dependencies

- **Severity:** P2 (Integrity Violation)
- **Layer:** TEE
- **Status:** Open
- **Location:** [`story-kernel/go.mod:19`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/go.mod#L19), [`go.mod:25`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/go.mod#L25)
- **Description:** `kyber v4.0.0-pre2` (pre-release) and `cb-mpc-go` (Piplabs fork, commit `7f03db8a8fa1`). Neither covered by advisory databases.
- **Recommendation:** Audit specific commits. Diff fork against upstream Coinbase. Plan migration to stable kyber.

### CDR-009: No Witness Count Enforcement in Light Client Configuration

- **Severity:** P2 (Integrity Violation)
- **Layer:** TEE
- **Status:** Open
- **Location:** [`story-kernel/story/query_client.go:87-94`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/story/query_client.go#L87-L94)
- **Description:** Light client accepts 0 witnesses. With `SkippingVerification` and 0 witnesses, a compromised primary RPC can serve forged state.
- **Recommendation:** Enforce `len(witness_addrs) >= 2` at startup.

### CDR-010: `replayMessages()` Silently Ignores All Errors

- **Severity:** P2 (Integrity Violation)
- **Layer:** TEE
- **Status:** Open
- **Location:** [`story-kernel/service/dist_key_gen.go:531-541`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dist_key_gen.go#L531-L541)
- **Description:** All errors from `ProcessDeal`, `ProcessResponse`, `ProcessJustification` are silently discarded during state reconstruction.
- **Recommendation:** Log and track failed replays. Return error on critical failures.

### CDR-011: CDR Fee Pool Minting Without EVM-Side Burn Verification

- **Severity:** P2 (Economic / Integrity)
- **Layer:** L1 Module
- **Status:** Open
- **Location:** [`story/client/x/dkg/keeper/dkg_cdr_fees.go:24-52`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_cdr_fees.go#L24-L52)
- **Description:** `AddCDRFeeToPool` mints new tokens via `bankKeeper.MintCoins` based on EVM event `Amount` field. The CL trusts the event without verifying the contract actually burned equivalent value.
- **Impact:** If CDR contract fee handling has a bug, this creates inflationary token issuance.
- **Recommendation:** Verify CDR contract's `FeeCollected` event only fires when tokens are deposited and burned.

### CDR-012: Raw Label Bytes in Registry Key Causing Potential Collisions

- **Severity:** P2 (Integrity Violation)
- **Layer:** L1 Module
- **Status:** Open
- **Location:** [`story/client/x/dkg/keeper/dkg_decrypt_registry.go:25`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_decrypt_registry.go#L25)
- **Description:** `label` (`[]byte`) formatted with `%s` in Sprintf key. If label contains `_` (delimiter) or non-printable bytes, key collisions occur. Other code paths use `hex.EncodeToString(label)`.
- **Recommendation:** Use `hex.EncodeToString(label)` consistently.

### CDR-013: Unbounded In-Memory Cache Growth in TEE Kernel

- **Severity:** P2 (Availability)
- **Layer:** TEE
- **Status:** Open
- **Location:** [`story-kernel/store/dkg_cache.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/store/dkg_cache.go), [`story-kernel/service/dkg_server.go:31-33`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_server.go#L31-L33)
- **Description:** All caches and per-round `sync.Map` mutexes grow without bound. No eviction policy. Enclave memory limited to 4GB.
- **Recommendation:** Implement LRU eviction. Evict entries older than `current_round - 2`.

### CDR-014: Broad `allowed_files` and `insecure__` Flags in Gramine Manifest

- **Severity:** P2 (Integrity Violation)
- **Layer:** TEE
- **Status:** Open
- **Location:** [`story-kernel/story-kernel.manifest.template:5`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/story-kernel.manifest.template#L5), [`line 19`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/story-kernel.manifest.template#L19), [`lines 62-69`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/story-kernel.manifest.template#L62-L69)
- **Description:** `file:/home/ubuntu/` allows host full read/write to enclave data directory. `file:/etc/ssl/` lets host control TLS trust store. `insecure__use_cmdline_argv = true` allows host to inject arguments.
- **Recommendation:** Minimize `allowed_files`. Move CA certs to `trusted_files`. Restrict or hardcode command-line arguments.

### CDR-015: CDR.sol `read()` Reverts When readConditionAddr Is address(0)

- **Severity:** P3 (Availability / Liveness)
- **Layer:** Contract
- **Status:** Open
- **Location:** [`story/contracts/src/protocol/CDR.sol:109`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L109), [`CDR.sol:191-199`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L191-L199)
- **Description:** Vaults allocated with only `writeConditionAddr` have permanently unreadable data.
- **Recommendation:** Add `require(vault.readConditionAddr != address(0), "CDR: Read condition not set")` in `read()`.

### CDR-016: Session Files Written World-Readable (0644)

- **Severity:** P3 (Information Disclosure)
- **Layer:** L1 Module
- **Status:** Open
- **Location:** [`story/client/x/dkg/keeper/state_manager.go:257`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/state_manager.go#L257)
- **Description:** Session files containing `DKGPubKey`, `CommPubKey`, and session state written with 0644 permissions.
- **Recommendation:** Use 0600 permissions.

### CDR-017: Integer Overflow in CDR Share Calculation (uint64 -> int64)

- **Severity:** P3 (Economic)
- **Layer:** L1 Module
- **Status:** Open
- **Location:** [`story/client/x/dkg/keeper/dkg_cdr_fees.go:183`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_cdr_fees.go#L183)
- **Description:** `count` is `uint64` cast to `int64` via `math.NewInt(int64(count))`. Overflows silently for large values.
- **Recommendation:** Use `math.NewIntFromBigInt` with proper `big.Int` conversion.

### CDR-018: CDR.sol Vault UUID Overflow (uint32)

- **Severity:** P4 (Economic)
- **Layer:** Contract
- **Status:** Open
- **Location:** [`story/contracts/src/protocol/CDR.sol:115`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L115)
- **Description:** `uint32` UUID counter. In Solidity 0.8.x, overflow reverts, permanently bricking `allocate()` after ~4.3B vaults.
- **Recommendation:** Use `uint256` or add explicit overflow check with descriptive error.

### CDR-019: DKG State Persisted as Plaintext JSON (Not Sealed)

- **Severity:** P5 (Information Disclosure)
- **Layer:** TEE
- **Status:** Open
- **Location:** [`story-kernel/store/dkg_state.go:91-107`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/store/dkg_state.go#L91-L107)
- **Description:** DKG state including justifications with plaintext `SecShare` scalars saved as plaintext JSON. Host operator can read secret share values from justification data.
- **Recommendation:** Seal DKG state files with `enclave.SealToFile()`.

### CDR-020: Go GC May Copy Secrets Before zeroBytes() Wipe

- **Severity:** P5 (Information Disclosure — requires TEE breach)
- **Layer:** TEE
- **Status:** Open
- **Location:** [`story-kernel/service/dkg_partial_decrypt.go:115`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_partial_decrypt.go#L115), [`dkg_finalize.go:159`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_finalize.go#L159)
- **Description:** Go GC may copy heap objects before `defer zeroBytes()` executes. Kyber scalar internals are never zeroed.
- **Recommendation:** Use CGO `C.malloc` for sensitive allocations or `runtime.KeepAlive`.

---

# CDR Supplementary Findings

Findings from delegated audits that are **not** in the main `CDR_SECURITY_AUDIT_REPORT.md`. These are lower-severity or layer-specific issues that did not meet the threshold for the integrated report but are still valid and actionable.

**Commits audited:**

- `piplabs/story` @ [`11307a5`](https://github.com/piplabs/story/tree/11307a50d2862d325898f2239e3c22bafd3a1565)
- `piplabs/story-kernel` @ [`6eb8e93`](https://github.com/piplabs/story-kernel/tree/6eb8e9345407a58f1cd49ff61a1d7f193a72b178)

---

## Smart Contract Findings

### SC-01: DKG.sol `register()` Does Not Prevent Duplicate Registrations On-Chain

- **Severity:** Medium
- **Location:** [`contracts/src/protocol/DKG.sol:176-218`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/DKG.sol#L176-L218)
- **Description:** No on-chain state tracks which validators have registered for a given round. A validator can call `register()` multiple times, each time paying the fee and emitting a `Registered` event with potentially different keys. The CL deduplicates, but the contract provides no protection.
- **Impact:** Event spam, potential CL confusion if it processes the latest event rather than the first. Economic cost deters casual spam but not funded attackers.
- **Recommendation:** Add `mapping(uint32 => mapping(address => bool)) hasRegistered` and enforce uniqueness.

### SC-02: Fee Burning via `transfer()` to address(0) May Brick on Non-Standard Chains

- **Severity:** Medium
- **Location:** [`contracts/src/protocol/CDR.sol:327`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L327), [`contracts/src/protocol/DKG.sol:36`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/DKG.sol#L36)
- **Description:** `payable(address(0x0)).transfer(amount)` forwards only 2300 gas. If Story ever deploys a precompile at address(0), this bricks all fee-paying operations.
- **Recommendation:** Use `.call{value: amount}("")` or a dedicated burn address like `0x...dEaD`.

### SC-03: SGXValidationHook Inline Assembly Lacks Bounds Check

- **Severity:** Medium
- **Location:** [`contracts/src/protocol/SGXValidationHook.sol:153-166`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/SGXValidationHook.sol#L153-L166)
- **Description:** `_extractReportInstanceDataCommitment` uses inline assembly `mload` at offset 368 without bounds-checking that `enclaveReport.length >= 432`. Unlike `_extractReportCodeCommitment` (which uses `BytesUtils.substring` with a bounds check), this function silently reads past allocated memory for short reports.
- **Impact:** Defense-in-depth — Automata DCAP should reject short reports, but the assembly has no safety net.
- **Recommendation:** Add `require(enclaveReport.length >= 432, "Report too short")` at the start of `validateReport()`.

### SC-04: Vault Conditions Immutable But Condition Contracts May Be Upgradeable

- **Severity:** Medium (Design)
- **Location:** [`contracts/src/protocol/CDR.sol:102-133`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L102-L133)
- **Description:** Once allocated, vault conditions cannot be changed. However, if a condition contract is a proxy, its logic can change without the vault's knowledge, effectively mutating conditions.
- **Recommendation:** Document this. Consider warning against upgradeable proxies as condition addresses.

### SC-05: DKG.sol `whitelistEnclaveType()` Requires Valid Data Even When De-Whitelisting

- **Severity:** Low
- **Location:** [`contracts/src/protocol/DKG.sol:99-118`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/DKG.sol#L99-L118)
- **Description:** When `isWhitelisted = false`, the function still validates `codeCommitment != 0` and `validationHookAddr != 0`. Owner must provide dummy valid data to de-whitelist.
- **Recommendation:** Only validate when `isWhitelisted = true`.

### SC-06: `submitEncryptedPartialDecryption()` Has No On-Chain Validation

- **Severity:** Low
- **Location:** [`contracts/src/protocol/CDR.sol:223-251`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L223-L251)
- **Description:** Accepts arbitrary `uuid`, `round`, `pid`, `signature` without checking vault existence, round validity, or sender registration. Fee is the only spam deterrent.
- **Recommendation:** Add `require(vaults[uuid].encryptedData.length > 0)` at minimum.

### SC-07: SGXValidationHook Extraction Functions Not Marked `pure`

- **Severity:** Low
- **Location:** [`contracts/src/protocol/SGXValidationHook.sol:144`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/SGXValidationHook.sol#L144), [`SGXValidationHook.sol:153`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/SGXValidationHook.sol#L153)
- **Description:** `_extractReportCodeCommitment` and `_extractReportInstanceDataCommitment` perform no state modifications but are not marked `pure`.
- **Recommendation:** Mark both `internal pure`.

### SC-08: DKG.sol Missing ReentrancyGuard

- **Severity:** Low
- **Location:** [`contracts/src/protocol/DKG.sol:10`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/DKG.sol#L10)
- **Description:** CDR.sol inherits `ReentrancyGuardUpgradeable`, but DKG.sol does not. DKG.sol makes external calls to `IAttestationReportValidator.validateReport()` in `_authenticateEnclaveReport`. If the validation hook is compromised, it could reenter DKG functions.
- **Recommendation:** Add `ReentrancyGuardUpgradeable` to DKG.sol.

### SC-09: CDR.sol `allocate()` Missing `nonReentrant`

- **Severity:** Low
- **Location:** [`contracts/src/protocol/CDR.sol:108`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L108)
- **Description:** `allocate()` has `whenNotPaused` but not `nonReentrant`, while `write()` and `read()` both have it. Inconsistent.
- **Recommendation:** Add `nonReentrant` for consistency.

### SC-10: No Timelock Enforced in `_authorizeUpgrade`

- **Severity:** Info
- **Location:** [`contracts/src/protocol/CDR.sol:333`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L333), [`DKG.sol:367`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/DKG.sol#L367), [`SGXValidationHook.sol:170`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/SGXValidationHook.sol#L170)
- **Description:** All contracts restrict to `onlyOwner`. Tests show owner is a `TimelockController`, but the contract doesn't enforce it. If ownership is transferred to a non-timelock address, upgrades become instant.
- **Recommendation:** Hardcode timelock check or document ownership requirement.

### SC-11: Test Coverage Critically Insufficient

- **Severity:** Info
- **Location:** [`contracts/test/cdr/cdr.t.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/test/cdr/cdr.t.sol), [`contracts/test/dkg/dkg.t.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/test/dkg/dkg.t.sol)
- **Description:** No tests for `write()`, `read()`, `submitEncryptedPartialDecryption()`, `register()`, `finalize()`, condition bypass, reentrancy, fee edge cases, or pausing.
- **Recommendation:** Add comprehensive test suite covering all entry points and adversarial scenarios.

### SC-12: Gas Inefficiency — Full Vault Struct Copied to Memory

- **Severity:** Info (Gas)
- **Location:** [`contracts/src/protocol/CDR.sol:148`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L148), [`CDR.sol:187`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol#L187)
- **Description:** `Vault memory vault = $.vaults[uuid]` copies the entire struct (including potentially large `encryptedData`) to memory in both `write()` and `read()`.
- **Recommendation:** Use storage pointers or load only needed fields.

---

## TEE Kernel Findings

### TEE-05: HKDF Used Without Salt in Partial Decryption Encryption

- **Severity:** Medium
- **Location:** [`service/dkg_partial_decrypt.go:292`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_partial_decrypt.go#L292)
- **Description:** HKDF key derivation uses nil salt: `hkdf.New(sha256.New, sharedBytes, nil, []byte("dkg-tdh2-partial"))`. Per RFC 5869, nil salt defaults to `HashLen` zeros. Secure when IKM has sufficient entropy (ECDH does), but a random salt provides stronger multi-target resistance.
- **Recommendation:** Use concatenated public keys as salt: `salt := append(ephPub, requesterPubKey...)`.

### TEE-06: TLS Not Enforced by Default

- **Severity:** Low
- **Location:** [`server/server.go:52-54`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/server/server.go#L52-L54)
- **Description:** Without TLS config, gRPC falls back to plaintext. Host-level attacker can observe/modify traffic metadata (round, validator, timing).
- **Recommendation:** Make TLS mandatory or document plaintext as local-only acceptable.

### TEE-07: Panic Recovery Interceptor May Mask Crypto State Corruption

- **Severity:** Low
- **Location:** [`server/server.go:192-208`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/server/server.go#L192-L208)
- **Description:** Catches all panics including those from kyber/cb-mpc. Continuing after a crypto panic may leave DKG state machine inconsistent. Logs full stack trace visible to host operator.
- **Recommendation:** For crypto panics, allow enclave to crash rather than recovering with corrupted state.

### TEE-08: `sync.Map` Per-Round Mutexes Never Pruned

- **Severity:** Low
- **Location:** [`service/dkg_server.go:31-33`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_server.go#L31-L33)
- **Description:** `initDKGMu`, `resharePrevMu`, `reshareNextMu` accumulate entries forever. Each ~16 bytes, practical impact minimal.
- **Recommendation:** Delete mutex entries after DKG round finalization.

---

## L1 Module Findings

### L1-01: `getDKGRegistrationsByRound` Full Table Scan

- **Severity:** Medium (Performance/DoS)
- **Location:** [`client/x/dkg/keeper/dkg_registration.go:74-91`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_registration.go#L74-L91)
- **Description:** `k.DKGRegistrations.Walk(ctx, nil, ...)` iterates ALL registrations across ALL rounds, then filters by prefix. Runs in BeginBlocker/FinalizeBlock hot path. Gets slower as rounds accumulate.
- **Recommendation:** Use `collections.Range` with prefix filtering to bound iteration.

### L1-02: No Size Limit on Ciphertext in Decrypt Registry

- **Severity:** Medium
- **Location:** [`client/x/dkg/keeper/dkg_handler.go:414-487`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_handler.go#L414-L487)
- **Description:** `ThresholdDecryptRequested` stores full `ciphertext` in the `DecryptRequestRegistry` without size limits. Large ciphertexts bloat consensus state.
- **Recommendation:** Enforce a maximum ciphertext size.

### L1-03: `GetAllCodeCommitments` Non-Deterministic Map Iteration (Off-Chain Only)

- **Severity:** Medium (Off-Chain)
- **Location:** [`client/x/dkg/keeper/kernel_router.go:142-158`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/kernel_router.go#L142-L158)
- **Description:** Iterates `r.clients` (Go map) returning code commitments in non-deterministic order. Used in `getRegistrationKernelClient` to pick `allCCs[0]`. Only affects off-chain DKG service, not consensus.
- **Recommendation:** Sort code commitments for deterministic selection.

### L1-04: `participantsRoot` Validation Behavior With Invalidated Dealers

- **Severity:** Medium (Documentation)
- **Location:** [`client/x/dkg/keeper/dkg_handler.go:175-221`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_handler.go#L175-L221)
- **Description:** `validateParticipantsRoot` includes `Verified` and `Finalized` registrations but not `Invalidated`. This is correct (invalidated dealers have status `Invalidated`, so they're excluded), but the behavior is implicit and undocumented.
- **Recommendation:** Add an explicit comment explaining that invalidated dealers are excluded by virtue of their status.

### L1-05: `UpgradeScheduled` Does Not Validate Activation Height vs Current DKG Round End

- **Severity:** Low
- **Location:** [`client/x/dkg/keeper/dkg_handler.go:226-259`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_handler.go#L226-L259)
- **Description:** Validates `activationHeight > currentBlockHeight` but not that it's far enough in the future to complete the current DKG round. Activating mid-round forces a `SkipToNextRound`.
- **Recommendation:** Validate activation height accounts for remaining round duration.

### L1-06: `proposalServer.AddVote` Does Not Validate Vote Content

- **Severity:** Low
- **Location:** [`client/x/dkg/keeper/proposal_server.go:16-23`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/proposal_server.go#L16-L23)
- **Description:** Only checks `Authority`, then returns success. Intentionally permissive per ABCI++ pattern (PrepareProposal is permissive, FinalizeBlock validates). Should be documented.
- **Recommendation:** Add comment documenting the intentional design.

### L1-07: `dkgSvcRound` Global Atomic Shared Across Test Instances

- **Severity:** Low (Testing)
- **Location:** [`client/x/dkg/keeper/dkg_svc.go:18`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_svc.go#L18)
- **Description:** Package-level `atomic.Uint64` shared across test instances running multiple keepers. Not a production issue.
- **Recommendation:** Consider test isolation if parallel keeper tests are needed.

### L1-08: `CleanupExpiredSessions` May Delete Active Sessions

- **Severity:** Low
- **Location:** [`client/x/dkg/keeper/state_manager.go:182-211`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/state_manager.go#L182-L211)
- **Description:** Removes all `PhaseCompleted`/`PhaseFailed` sessions. If the decrypt worker is still running for a completed session, deleting it breaks `GetSession`. The TODO at line 189 acknowledges this.
- **Recommendation:** Check if decrypt worker is active before deleting completed sessions.

### L1-09: `GlobalPubKeyVotes` Not Bounded Per Validator

- **Severity:** Low (False Positive)
- **Location:** [`client/x/dkg/keeper/dkg_votes.go:15-46`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_votes.go#L15-L46)
- **Description:** `AddGlobalPubKeyVote` increments without checking if the same validator already voted. However, the `Finalized` handler checks `reg.Status == DKGRegStatusFinalized` before reaching the vote, preventing double-voting. Protected by the double-finalization check — not an actual vulnerability.

### L1-10: `ProcessDKGEvents` Silently Continues on Error

- **Severity:** Info
- **Location:** [`client/x/evmengine/keeper/dkg.go:20-70`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/evmengine/keeper/dkg.go#L20-L70)
- **Description:** Failed events are silently skipped with `continue`. Each handler uses `CacheContext`/`writeCache` for isolation. Reasonable design choice for EVM event processing, but logging should be `Warn` not `Error` since these are handled conditions.

### L1-11: BeginBlocker Determinism Confirmed Sound

- **Severity:** Info (Positive)
- **Location:** [`client/x/dkg/keeper/abci.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/abci.go)
- **Description:** Uses only deterministic inputs (height, params, round state). Async goroutines gated by `isDKGSvcEnabled` affect only off-chain state. `pruneTimedOutDecryptRequests` uses fixed block intervals. Confirmed safe for consensus.

### L1-12: Async Goroutine Safety Confirmed Sound

- **Severity:** Info (Positive)
- **Location:** [`client/x/dkg/keeper/dkg_svc.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_svc.go)
- **Description:** Async goroutines use `context.Background()` with 1-minute timeout, avoiding CometBFT context cancellation. `tryAcquireDKGSvc`/`releaseDKGSvc` provides round-level deduplication. `dkgKernelMu` serializes kernel operations. Patterns are sound.

### Summaries

### Smart Contract Audit (web3-security-auditor)

- **3 High**: read() reverts on address(0) condition, condition bypass, finalize() no validation
- **4 Medium**: UUID overflow, no duplicate registration check, fee burn to address(0), SGXValidationHook extraction inconsistency
- **4 Low**: whitelistEnclaveType UX, submitPartialDecryption no validation, extraction functions not pure, no reentrancy guard on DKG.sol
- **1 Info**: No timelock enforced in `_authorizeUpgrade`, test coverage critically insufficient

### TEE Kernel Audit (tee-security-auditor)

- **2 Critical**: DKG state plaintext (TEE-001), sgx.debug=true (TEE-002)
- **2 High**: Missing on-chain verification (TEE-003), GC secret copy (TEE-004)
- **3 Medium**: insecure__ flags (TEE-005), unbounded cache (TEE-006), HKDF no salt (TEE-007)
- **3 Low**: TLS optional (TEE-008), panic interceptor (TEE-009), sync.Map leak (TEE-010)
- **3 Info**: kyber pre-release (TEE-011), cb-mpc fork (TEE-012), broad allowed_files (TEE-013)

### L1 Module Audit (general-purpose)

- **2 P1**: Non-deterministic map iteration (F-08), missing Ethereum sig prefix (F-04)
- **5 P2**: Full table scan (F-03), CDR fee minting (F-06), raw label bytes (F-10), ciphertext size (F-11), participantsRoot includes invalidated (F-05)
- **5 P3**: int overflow (F-09), upgrade timing (F-12), proposal no validation (F-13), global atomic (F-14), session cleanup (F-16)
- **3 Info**: BeginBlocker determinism sound (I-01), EVM trust model correct (I-02), async safety sound (I-03)

---

## Cryptographic Protocol Assessment

### Pedersen DKG (Phase 3)

| Check | Status | Notes |
| --- | --- | --- |
| Group parameters (Edwards25519, DDH-hard) | PASS | kyber v4 BlakeSHA256Ed25519 suite |
| Generator provenance | PASS | kyber library standard generators |
| Commitment format (Pedersen) | PASS | Handled by kyber internally |
| Polynomial degree enforcement | PASS | kyber `NewDistKeyGenerator` enforces threshold |
| Share verification equation | PASS | kyber Pedersen VSS verification |
| Encrypted share transport | PASS | kyber `EncryptedDeal` with DH key exchange |
| Complaint mechanism | PASS | Responses + justifications with Schnorr sig + VSS verify |
| Rogue key prevention | PARTIAL | Implicit PoK via Schnorr-signed deals, no standalone NIZK |
| Public key computation (qualified set Q) | PASS | `DistKeyShare.Public()` over qualified set |
| Bias assessment | KNOWN | Pedersen DKG is biased. Acceptable for encryption-only (P2) |
| Threshold consistency across layers | PASS | `operationalThreshold=667/1000` consistent |
| DKG 3 scenarios (fresh/reshare-prev/reshare-next) | PASS | Correctly dispatched |

### TDH2 (Phase 4)

| Check | Status | Notes |
| --- | --- | --- |
| AES-GCM encryption (partial to requester) | PASS | Random nonce, 32-byte key via HKDF |
| TDH2 partial decryption via cb-mpc | PASS | Correct PID and share conversion |
| Ciphertext validity check before decrypt | UNKNOWN | Delegated to cb-mpc C++ internals |
| On-chain request verification | **FAIL** | CDR-002 |
| Share zeroization | PARTIAL | CDR-020 (GC risk) |
| Proactive refresh | NOT IMPLEMENTED | Only resharing exists |

---

## Vulnerability Taxonomy Coverage

| ID | Attack Class | Applicable | Mitigated | Residual Risk |
| --- | --- | --- | --- | --- |
| T1 | DKG Bias | Yes | No (by design) | P2 |
| T2 | Rogue Key | Yes | Partial (implicit PoK) | Low |
| T3 | Polynomial Degree Violation | Yes | Yes (kyber) | None |
| T4 | Inconsistent Dealing | Yes | Yes (VSS + complaints) | None |
| T5 | CCA Attack on TDH2 | Yes | Unknown (cb-mpc) | P1 |
| T6 | Partial Decryption Oracle | Yes | Partial | P1 (CDR-002) |
| T7 | NIZK Forgery | Yes | Yes (cb-mpc) | Low |
| T8 | Side-Channel Share Recovery | Yes | Partial (SGX) | P5/P0 if debug |
| T9 | Adaptive Adversary | Yes | No explicit defense | P1 |
| T10 | Lagrange Manipulation | Yes | N/A (off-chain) | N/A |
| B1 | Event Replay | Yes | Partial (CL dedup) | P2 |
| B2 | Event Omission | Yes | Partial (consensus) | P3 |
| B3 | Front-Running DKG | Yes | Partial | P2 |
| B4 | Gas Griefing | Yes | Partial (fees) | P4 |
| B5 | State Desync | Yes | Partial (light client) | P2 |
| B6 | Condition Bypass | Yes | Partial | P2 (CDR-006) |
| B7 | Forged Attestation | Yes | Yes (DCAP) | Low |
| B8 | Finalization Without Validation | Yes | Partial (CL) | P2 (CDR-005) |
| L1 | Stuck DKG Phase | Yes | Yes (SkipToNextRound) | Low |
| L2 | Committee Below Threshold | Yes | Yes (MinReqFinalized) | Low |
| L3 | Decryption Denial | Yes | Partial (timeout) | P3 |
| L4 | Refresh/Reshare Failure | Yes | Partial (SkipToNextRound) | P3 |
| K1 | CGO Memory Corruption | Yes | Unknown (needs audit) | P0 potential |
| K2 | Go GC Secret Recovery | Yes | Partial (zeroBytes) | P5 (CDR-020) |
| K3 | DKG Cache Poisoning | Yes | Partial (RWMutex) | Low |
| K4 | Serialization Confusion | Yes | Low risk (JSON+protobuf) | Low |
| K5 | Supply Chain Compromise | Yes | Not mitigated | P2 (CDR-008) |

---

## Priority Remediation Order

| Priority | Finding | Action |
| --- | --- | --- |
| **IMMEDIATE** | CDR-001 | Set `sgx.debug = false` in production manifest |
| **IMMEDIATE** | CDR-003 | Sort map keys in `distributeCDRRewardPool` to prevent chain halt |
| **HIGH** | CDR-002 | Implement on-chain read request verification in TEE |
| **HIGH** | CDR-004 | Verify and align partial decryption signature prefix convention |
| **HIGH** | CDR-019 | Seal DKG state files with `enclave.SealToFile()` |
| **MEDIUM** | CDR-005 | Add `validatorAddr == msg.sender` to `finalize()` |
| **MEDIUM** | CDR-006 | Document or remove condition bypass |
| **MEDIUM** | CDR-009 | Enforce minimum 2 witnesses for light client |
| **MEDIUM** | CDR-012 | Use `hex.EncodeToString(label)` in registry key |
| **MEDIUM** | CDR-014 | Restrict Gramine manifest allowed_files |
| **LOW** | CDR-015 | Add readConditionAddr != 0 check in `read()` |
| **LOW** | CDR-016 | Fix session file permissions to 0600 |

---

## Appendix: Files Audited

### Smart Contract Layer (14 files)

- [`contracts/src/protocol/CDR.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/CDR.sol), [`DKG.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/DKG.sol), [`SGXValidationHook.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/protocol/SGXValidationHook.sol)
- [`contracts/src/interfaces/ICDR.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/interfaces/ICDR.sol), [`IDKG.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/interfaces/IDKG.sol), [`ICDRWriteCondition.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/interfaces/ICDRWriteCondition.sol), [`ICDRReadCondition.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/interfaces/ICDRReadCondition.sol), [`ISGXValidationHook.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/interfaces/ISGXValidationHook.sol), [`IAttestationReportValidator.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/interfaces/IAttestationReportValidator.sol)
- [`contracts/src/libraries/BytesUtils.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/libraries/BytesUtils.sol), [`Predeploys.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/src/libraries/Predeploys.sol)
- [`contracts/test/cdr/cdr.t.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/test/cdr/cdr.t.sol), [`test/dkg/dkg.t.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/test/dkg/dkg.t.sol), [`test/upgrades/PredeployUpgrades.t.sol`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/contracts/test/upgrades/PredeployUpgrades.t.sol)

### TEE Kernel Layer (22 files)

- [`service/dkg_server.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_server.go), [`dist_key_gen.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dist_key_gen.go), [`dkg_generate_key.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_generate_key.go), [`dkg_generate_deals.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_generate_deals.go), [`dkg_process_deals.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_process_deals.go), [`dkg_process_responses.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_process_responses.go), [`dkg_process_justification.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_process_justification.go), [`dkg_finalize.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_finalize.go), [`dkg_partial_decrypt.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_partial_decrypt.go), [`dkg_get_code_commitment.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/dkg_get_code_commitment.go), [`round_context.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/round_context.go), [`helper.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/helper.go), [`hasher.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/hasher.go), [`utils.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/service/utils.go)
- [`enclave/seal.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/enclave/seal.go), [`quote.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/enclave/quote.go), [`sealed_leveldb.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/enclave/sealed_leveldb.go), [`noflock_storage.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/enclave/noflock_storage.go)
- [`store/dkg_store.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/store/dkg_store.go), [`dkg_cache.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/store/dkg_cache.go), [`dkg_state.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/store/dkg_state.go), [`key_store.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/store/key_store.go)
- [`story/query_client.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/story/query_client.go), [`server/server.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/server/server.go), [`types/types.go`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/types/types.go)
- [`story-kernel.manifest.template`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/story-kernel.manifest.template), [`go.mod`](https://github.com/piplabs/story-kernel/blob/6eb8e9345407a58f1cd49ff61a1d7f193a72b178/go.mod)

### L1 Module Layer (35+ files)

- [`client/x/dkg/keeper/`](https://github.com/piplabs/story/tree/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper) — [`abci.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/abci.go), [`dkg_handler.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_handler.go), [`dkg_initialization.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_initialization.go), [`dkg_registration.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_registration.go), [`dkg_dealing.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_dealing.go), [`dkg_finalization.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_finalization.go), [`dkg_justification.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_justification.go), [`dkg_partial_decryption.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_partial_decryption.go), [`dkg_decrypt_registry.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_decrypt_registry.go), [`dkg_network.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_network.go), [`dkg_round.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_round.go), [`dkg_rewards.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_rewards.go), [`dkg_cdr_fees.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_cdr_fees.go), [`dkg_events.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_events.go), [`dkg_svc.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_svc.go), [`dkg_svc_registration.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_svc_registration.go), [`dkg_svc_dealing.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_svc_dealing.go), [`dkg_svc_finalization.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_svc_finalization.go), [`dkg_svc_complete.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/dkg_svc_complete.go), [`kernel_router.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/kernel_router.go), [`kernel_client.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/kernel_client.go), [`kernel_upgrade.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/kernel_upgrade.go), [`state_manager.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/state_manager.go), [`keeper.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/keeper.go), [`genesis.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/genesis.go), [`msg_server.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/msg_server.go), [`proposal_server.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/proposal_server.go), [`query.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/query.go), [`queue.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/queue.go), [`vote.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/vote.go), [`params.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/params.go), [`contract_client.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/dkg/keeper/contract_client.go)
- [`client/x/evmengine/keeper/dkg.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/evmengine/keeper/dkg.go), [`cdr.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/client/x/evmengine/keeper/cdr.go)
- [`lib/vss/verify.go`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/lib/vss/verify.go)

### Design Documents

- [`docs/design/CDR.md`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/docs/design/CDR.md), [`docs/design/DKG.md`](https://github.com/piplabs/story/blob/11307a50d2862d325898f2239e3c22bafd3a1565/docs/design/DKG.md)