# DKG/CDR Production Final Audit Report (Revised)

**Date:** 2026-03-22 (initial), 2026-03-22 (design doc re-validation)
**Auditors:** 6 specialized agents (Cosmos SDK, DKG Service, Smart Contract, TEE/SGX, Cryptography, Cross-cutting Security)
**Design Docs:** `story/docs/design/DKG.md`, `story/docs/design/CDR.md`
**Scope:**
- story: branch `dkg/refactor-v200-upgrade` (v1.5.3 → v2.0.0 upgrade, DKG module, DKG service, rewards/fees, contracts, deploy scripts)
- story-kernel: branch `hans/fix-dealer-poly-persistence` (full codebase, first production deployment)

**Overall Risk Rating: HIGH** — Multiple CRITICAL findings must be resolved before production deployment.

---

## Revision Notes

Initial audit findings were re-validated against the official design documents (DKG.md, CDR.md). This resulted in:
- **14 findings reclassified as FALSE POSITIVE** (behavior matches design intent)
- **1 finding upgraded to CRITICAL** (CRYPTO-001: sign/verify field mismatch)
- **Multiple severity adjustments** (mostly downgrades due to TEE mitigation and design intent)

---

## Findings Summary (Post Re-validation)

| Severity | Count | Change from Initial |
|----------|-------|---------------------|
| CRITICAL | 5 | was 10, -6 FP/downgrade, +1 new |
| HIGH | 8 | was 16, multiple downgrades |
| MEDIUM | 17 | adjusted |
| LOW | 12 | adjusted |
| FALSE POSITIVE | 14 | removed from active findings |

---

## CRITICAL Findings (5)

### C-01: CDR Fee Distribution — Non-Deterministic Map Iteration [CONSENSUS FAILURE]
- **IDs:** COSMOS-001 / DKG-SVC-009
- **Location:** `client/x/dkg/keeper/dkg_cdr_fees.go:182`
- **Description:** `distributeCDRRewardPool` iterates over `map[string]uint64` to send bank transfers. Go map iteration is non-deterministic → different validators produce different app hashes.
- **Impact:** **Guaranteed chain halt** when CDR fee pool has balance and multiple validators contributed.
- **Re-validation:** CONFIRMED. Design doc Section 3.1 shows `dkg_rewards.go` correctly uses `sort.Strings(memberAddrs)` but CDR fees do not.
- **Fix:** Sort map keys before iteration.

### C-02: SGX Debug Mode Enabled in Production Manifest
- **ID:** KERNEL-001
- **Location:** `story-kernel.manifest.template:51`
- **Description:** Gramine manifest has `sgx.debug = true`. Enclave memory is inspectable by host OS.
- **Impact:** **Entire TEE security model void.** All private keys, DKG shares readable by host.
- **Re-validation:** CONFIRMED. Design doc Section 5.5 assumes SGX sealing integrity which debug mode negates.
- **Fix:** Set `debug = false` for production. Template-variable for dev/prod separation.

### C-03: Partial Decryption Signature — Sign/Verify Field Mismatch [CDR BROKEN]
- **ID:** CRYPTO-001 (upgraded from MEDIUM)
- **Locations:**
  - Kernel sign: `story-kernel/service/dkg_service.go:711-718` — signs `codeCommitment || round || encryptedPartial || ephPubKey || pubShare`
  - CL verify: `story/client/x/dkg/keeper/dkg_handler.go:362-410` — verifies `round || ciphertext || encryptedPartial || ephPubKey || pubShare`
- **Description:** Kernel includes `codeCommitment` (32B) at front, CL includes `ciphertext` instead. Different hash → `SigToPub` recovers wrong address → **ALL partial decryption submissions fail signature verification**.
- **Impact:** **CDR threshold decryption is completely non-functional.** No partial decryption can pass CL verification.
- **Re-validation:** CONFIRMED as CRITICAL. CDR design doc Section 5 line 119 specifies kernel format (`CC || round || ...`), but CL code doesn't follow it.
- **Fix:** Align both sides to the design doc specification: `keccak256(CC || round || encryptedPartial || ephPubKey || pubShare)`.

### C-04: Decrypt Worker Context Leak — CDR Non-Functional
- **ID:** DKG-SVC-007
- **Location:** `client/x/dkg/keeper/dkg_svc.go:240-261`, `dkg_svc_complete.go:60`
- **Description:** `StartDecryptWorker` receives `dkgAsyncContext()` (1-min timeout). When `handleDKGComplete` finishes, context is canceled → decrypt worker goroutine terminates after ~1 minute.
- **Impact:** DecryptWorker stops after 1 minute. `ResumeDKGService` can restart it on next block, but creates periodic 1-min gaps in CDR processing.
- **Re-validation:** CONFIRMED. Design doc Section 2.9 says "StartDecryptWorker" without specifying context lifetime. Code uses dkgAsyncContext which is designed for short-lived operations.
- **Fix:** Pass an independent long-lived context bound to app lifecycle.

### C-05: Decrypt Request Registry Key Collision
- **ID:** COSMOS-002 / SECURITY-004
- **Location:** `client/x/dkg/keeper/dkg_decrypt_registry.go:17-27`
- **Description:** Code uses raw `label` bytes via `%s` format. Label is 32 bytes with null padding — null bytes in KV key cause parsing issues. Design doc Appendix A specifies `"{requesterHash}_{labelHex}"`.
- **Re-validation:** CONFIRMED. Design doc explicitly says `hexLabel` but code uses raw bytes.
- **Fix:** Use `hex.EncodeToString(label)`.

---

## HIGH Findings (8)

### H-01: Placeholder MRENCLAVE in Deploy Scripts
- **ID:** CONTRACT-002
- **Location:** `contracts/script/upgrades/DeploySGXValidationHook.s.sol:35`, `GenerateAlloc.s.sol:63-65`
- **Description:** `SGX_CODE_COMMITMENT = hex"0000...0001"` and `AUTOMATA_VALIDATION_ADDR = address(uint160(1000))` are dev placeholders.
- **Re-validation:** CONFIRMED. GenerateAlloc.s.sol generates actual genesis alloc files.
- **Fix:** Replace before deployment. CI check for placeholder values.

### H-02: Params Store Dual Registration
- **ID:** COSMOS-003
- **Location:** `client/x/dkg/keeper/params.go:11-25`, `keeper.go:120`
- **Description:** Params registered as `collections.Item` but accessed via raw KV store.
- **Re-validation:** CONFIRMED. Fragile pattern that breaks on collections library changes.
- **Fix:** Unify to one approach.

### H-03: Replay Messages Silently Discards All Errors (Kernel)
- **ID:** KERNEL-005 / CRYPTO-002
- **Location:** `service/dist_key_gen.go:454-467`
- **Description:** `_, _ =` for all ProcessDeal/ProcessResponse/ProcessJustification during replay. Corrupted state silently accepted.
- **Re-validation:** CONFIRMED. Not documented in design.
- **Fix:** Log errors at WARN. Return error if critical messages fail.

### H-04: cachedLastBlockHeight Race Condition (Kernel)
- **ID:** KERNEL-007
- **Location:** `story/query_client.go:650-661`
- **Description:** int64 field written by background goroutine, read by query methods without synchronization.
- **Re-validation:** CONFIRMED. Data race under Go memory model.
- **Fix:** Use `atomic.Int64`.

### H-05: Upgrade Activation Doesn't Terminate Current Round
- **ID:** SECURITY-001
- **Location:** `client/x/dkg/keeper/abci.go:55-70`
- **Description:** When upgrade activates, `InitiateDKGRound(isUpgrade=true)` is called without terminating the current round or flushing queues. Stale deals/responses may contaminate new round.
- **Re-validation:** CONFIRMED. Design doc Section 4.3 doesn't address this scenario.
- **Fix:** Call `FlushAllQueues()` before initiating upgrade round.

### H-06: CDR Partial Decryption Correctness Not Verified On-Chain
- **ID:** SECURITY-002
- **Location:** `client/x/dkg/keeper/dkg_handler.go`
- **Description:** CL verifies signature + pubShare but not cryptographic correctness. Malicious validator can submit garbage partials, collect fees, prevent decryption.
- **Re-validation:** CONFIRMED as intentional trade-off (on-chain TDH2 verification is expensive). Economic attack vector is real.
- **Fix:** Consider post-hoc slashing or requester feedback mechanism.

### H-07: Package-Level Global State
- **ID:** DKG-SVC-002 / SECURITY-003
- **Location:** `client/x/dkg/keeper/keeper.go:21-51`
- **Description:** deals, responses, justifications queues, dkgSvcRound as package-level globals.
- **Re-validation:** CONFIRMED as TRUE ISSUE despite being documented in design doc Section 5.3. Design documents the pattern but doesn't justify it — tests already acknowledge the problem ("NOT parallel because they share package-level globals").
- **Fix:** Move to Keeper struct fields.

### H-08: No Mutual Exclusion on Kernel DKG RPCs After Cache Hit
- **ID:** KERNEL-008
- **Location:** `service/dkg_server.go:31-50`
- **Description:** CL-side `dkgKernelMu` provides serialization, but kernel has no internal mutex. Without mTLS, unauthorized callers bypass CL serialization.
- **Re-validation:** SEVERITY ADJUSTED (CRITICAL→HIGH→MEDIUM→HIGH). Risk depends on network isolation. HIGH if gRPC is accessible beyond localhost.
- **Fix:** Add per-round mutex inside kernel for all DKG operations.

---

## MEDIUM Findings (17)

| ID | Finding | Re-validation |
|----|---------|---------------|
| M-01 | CDR readConditionAddr=address(0) permanent vault lock | CONFIRMED (CONTRACT-003) |
| M-02 | DKG.sol finalize() missing msg.sender check | SEVERITY ADJUSTED H→M (CONTRACT-001). CL verifies TEE signature, but defense-in-depth missing |
| M-03 | getDKGRegistrationsByRound full store walk O(n) | SEVERITY ADJUSTED H→M (COSMOS-005). Performance issue, not immediate danger |
| M-04 | totalCount uint64→int64 overflow in CDR fees | CONFIRMED (COSMOS-009) |
| M-05 | ProcessDKGEvents silent continue on error | CONFIRMED (DKG-SVC-014). Intentional per-event isolation but missing retry |
| M-06 | FinalizeDKGRound missing GlobalPublicKey check | CONFIRMED (SECURITY-005). Could finalize round without agreed key |
| M-07 | No rate limiting on kernel reconnection | SEVERITY ADJUSTED H→M (DKG-SVC-011→LOW by Cosmos, keeping M for completeness) |
| M-08 | Unbounded vote extension aggregation | SEVERITY ADJUSTED H→M (DKG-SVC-004). Dedup provides natural bound |
| M-09 | DKG state file not written atomically | CONFIRMED (KERNEL-011). No fsync/rename pattern |
| M-10 | Unbounded in-memory caches in kernel | CONFIRMED (KERNEL-012). SGX 4GB enclave limit risk |
| M-11 | Gramine manifest broad file access | CONFIRMED (KERNEL-013). /home/ubuntu/ too permissive |
| M-12 | No auth on kernel gRPC + reflection enabled | CONFIRMED (KERNEL-014). TODO to remove reflection |
| M-13 | No on-chain verification for decrypt requests | CONFIRMED (KERNEL-015). Known TODO in code |
| M-14 | CDR fee pool dual tracking (EL burn + CL re-mint) | CONFIRMED (SECURITY-009/CONTRACT-011) |
| M-15 | GetAllCodeCommitments non-deterministic order | CONFIRMED (SECURITY-010). Upgrade round client selection |
| M-16 | Resharing polynomial coefficients not persisted | SEVERITY ADJUSTED CRITICAL→M (CRYPTO-009). Limited exploit window, complaint-only path |
| M-17 | Private key zeroing best-effort | CONFIRMED (CRYPTO-006). Mitigated by SGX |

---

## LOW Findings (12)

| ID | Finding | Re-validation |
|----|---------|---------------|
| L-01 | `retryAttemts` typo | CONFIRMED |
| L-02 | Event StartBlockHeight uint32 truncation | SEVERITY ADJUSTED H→L (COSMOS-004). Event-only, not consensus |
| L-03 | ExportGenesis only exports Params | CONFIRMED |
| L-04 | Deploy script salt reuse limitation | CONFIRMED (CONTRACT-007) |
| L-05 | BytesUtils.substring missing error message | CONFIRMED (CONTRACT-008) |
| L-06 | SGXValidationHook functions should be pure | CONFIRMED (CONTRACT-009) |
| L-07 | submitEncryptedPartialDecryption missing nonReentrant | SEVERITY ADJUSTED M→L (CONTRACT-004). address(0) can't callback |
| L-08 | uint32 vault UUID | SEVERITY ADJUSTED M→L (CONTRACT-005). Solidity 0.8 reverts on overflow |
| L-09 | Session file permission 0644 | SEVERITY ADJUSTED M→L (DKG-SVC-015). No key material in file |
| L-10 | retry() time.Sleep ignores context | SEVERITY ADJUSTED H→M→L. Max 6s delay |
| L-11 | StateManager file I/O not atomic (story side) | SEVERITY ADJUSTED H→M→L. Off-chain, recoverable |
| L-12 | HKDF nil salt in partial decrypt | SEVERITY ADJUSTED M→L (CRYPTO-008). RFC 5869 compliant |

---

## FALSE POSITIVES Removed (14)

| ID | Original Finding | Reason |
|----|-----------------|--------|
| COSMOS-006 | SettlementBalance stored as string | math.Int standard string encoding |
| COSMOS-007 | Fork.UpgradeHeight=0 | Fork struct has no UpgradeHeight field in code |
| COSMOS-008 | AddVote ignores ProcessDeals errors | Intentional: off-chain async, no state changes in error path |
| COSMOS-010 | Partial decryption uses JSON | Design doc Appendix A explicitly specifies "JSON bytes" |
| DKG-SVC-002 | Global state package variables | Design doc Section 5.3 explicitly documents this as the concurrency model |
| DKG-SVC-012 | Justification dedup key semantic | SecShare.I = recipient index, matches design doc dedup key spec |
| KERNEL-002 | state.json plaintext (CRITICAL→FP) | No private keys in state.json. Deals are encrypted. Justification SecShares are publicly broadcast data |
| KERNEL-009 | HKDF nil salt (HIGH→FP) | RFC 5869 compliant, design doc matches, high-entropy ECDH input |
| KERNEL-010 | Unsafe reflection on kyber | No reflection in business logic — only in protobuf auto-generated code |
| KERNEL-016 | zeroPrivateKey function | Function does not exist in story-kernel codebase |
| CRYPTO-004 | Unsafe reflection for poly restore | Same as KERNEL-010 — no reflection on kyber internals |
| CRYPTO-007 | No threshold boundary validation | Validated at CL params, CalculateThreshold, AND kyber library |
| SECURITY-007 | UpgradeScheduled height not validated | Both contract AND CL validate activationHeight > current block |
| SECURITY-011 | Activated upgrades accumulate | deleteActivatedUpgradeInfo confirmed in FinalizeDKGRound |

---

## Remediation Priority (Revised)

### P0 — Must Fix Before Production (5 items)

| # | ID | Issue | Effort |
|---|----|-------|--------|
| 1 | C-01 | CDR fee non-deterministic map iteration | Small |
| 2 | C-02 | SGX debug=true in manifest | Trivial |
| 3 | C-03 | **Partial decrypt sign/verify field mismatch** | Small |
| 4 | C-04 | Decrypt worker context leak | Small |
| 5 | C-05 | Registry key collision (raw label) | Small |

### P1 — Should Fix Before Production (8 items)

| # | ID | Issue | Effort |
|---|----|-------|--------|
| 6 | H-01 | Placeholder MRENCLAVE in deploy scripts | Trivial (ops) |
| 7 | H-02 | Params dual registration | Small |
| 8 | H-03 | Silent replay error discard (kernel) | Small |
| 9 | H-04 | cachedLastBlockHeight race | Small |
| 10 | H-05 | Upgrade activation doesn't flush queues | Small |
| 11 | H-06 | CDR partial decrypt not verified on-chain | Large |
| 12 | H-07 | Package-level global state | Large |
| 13 | H-08 | Kernel DKG RPC mutex | Medium |

### P2 — Should Fix Shortly After Production

All MEDIUM items (M-01 through M-17).

### P3 — Backlog

All LOW items.

---

## Test Coverage

| Package | Coverage | Target | Status |
|---------|----------|--------|--------|
| story/x/dkg/keeper | 38.7% | 80% | FAIL |
| story/x/dkg/types | 1.1% | 80% | FAIL |
| story/app/upgrades/v_2_0_0 | 58.1% | 80% | FAIL |
| story-kernel/config | 95.0% | 80% | PASS |
| story-kernel/crypto | 81.8% | 80% | PASS |
| story-kernel/dkgutil | 78.0% | 80% | WARN |
| story-kernel/store | 76.6% | 80% | WARN |
| story-kernel/story | 84.4% | 80% | PASS |
| story-kernel/types | 84.8% | 80% | PASS |
| story-kernel/service | N/A (CGO) | 80% | FAIL |
| story/lib/vss | 100.0% | 80% | PASS |

---

## Positive Design Observations

- Upgrade Handler: correct migration ordering, VE height = upgradeHeight+1, idempotent
- Vote Extension: ABCI status returns per CometBFT spec, proper size limits
- Reward Distribution: `sort.Strings(memberAddrs)` ensures determinism (CDR fees should follow same pattern)
- CacheContext: per-event state rollback in evmengine
- Async Context: pre-computation of SDK data before goroutine launch
- dkgKernelMu: CL-side serialization of kernel mutations
- Binary-Swap Upgrade: on-chain ScheduleUpgrade + disk fallback
- Module Permissions: Minter for DKG, separate CDR fee pool, both in blockAccAddrs
- Contract Security: UUPS + _disableInitializers, ERC-7201, Ownable2Step, TimelockController

---

*Generated by 6-agent audit team on 2026-03-22*
*Re-validated against design docs: story/docs/design/DKG.md, story/docs/design/CDR.md*
*Repositories: story@dkg/refactor-v200-upgrade, story-kernel@hans/fix-dealer-poly-persistence*
