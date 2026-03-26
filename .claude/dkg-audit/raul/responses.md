# External Audit Response: CDR AI Security Audit + Apex Report

**Date:** 2026-03-24
**Reviewed by:** Story Blockchain Architect
**External Reports:**
- CDR AI Security Audit Report (commits: story@`11307a5`, story-kernel@`6eb8e93`)
- Apex Report - Story network / Scan #1 (same commits)

**Internal Reference:** `story/.claude/dkg-audit/AUDIT_REPORT.md` and `story/.claude/FINAL_AUDIT_1.md`

**Fix branches:**
- story: `dkg/fix-audit-findings` (14 commits on top of dkg/dev)
- story-kernel: `fix/audit-findings` (4 commits on top of main)

**Note:** Several fixes have already been merged into `origin/dkg/dev` (commits `a56fbc5f`, `71ded0d1`, `0382ec7f`, `efd59a56`) and `origin/main` (commit `600bf54`) since the audited commits.

---

## CDR AI Security Audit Report Findings

---

### [CDR Report] Finding CDR-001: Gramine SGX Debug Mode Enabled in Manifest Template
- **Severity**: P0 (Key Compromise)
- **Status**: DUPLICATE / ALREADY FIXED
- **Duplicate of**: Our C-02 (KERNEL-001) in FINAL_AUDIT_1.md, C-03 in AUDIT_REPORT.md
- **Fixed in**: Not yet fixed in code (operational fix, needs `sgx.debug = false` for production manifest)
- **Analysis**: Valid finding. Exact same issue identified in our audit. `sgx.debug = true` at `story-kernel.manifest.template:51` disables SGX memory encryption. Our report rates this as CRITICAL. The fix is an operational change for the production manifest — no code branch fix needed, as this is a deployment-time configuration.

---

### [CDR Report] Finding CDR-002: TEE Does Not Verify On-Chain Read Request Before Partial Decryption
- **Severity**: P1 (Confidentiality Breach)
- **Status**: DUPLICATE / PARTIALLY FIXED
- **Duplicate of**: Our M-13 (KERNEL-015) in FINAL_AUDIT_1.md, H-06 (SECURITY-002) in AUDIT_REPORT.md
- **Fixed in**: `origin/main` commit `600bf54` (story-kernel) adds `QueryClient` RPC for validating decrypt request existence canonically; story commit `6a4c4174` adds the new QueryClient RPC
- **Analysis**: Valid finding. Our audit identified this as M-13/H-06 — the known TODO in the code acknowledging the missing on-chain verification. Since the audited commit, a partial fix has been merged: story-kernel now has a `ValidateDecryptRequest` RPC that queries canonical chain state. However, full light-client-verified proof validation inside the enclave may still be incomplete. The CDR report recommends verifying a `VaultRead` event which is the same direction as the implemented fix.

---

### [CDR Report] Finding CDR-003: Non-Deterministic Map Iteration in CDR Reward Distribution
- **Severity**: P1 (Consensus Safety)
- **Status**: DUPLICATE / ALREADY FIXED
- **Duplicate of**: Our C-01 (COSMOS-001 / DKG-SVC-009) in both audit reports
- **Fixed in**: `dkg/fix-audit-findings` commit `930a1ab0`; also merged to `origin/dkg/dev` via `71ded0d1` (fix(cdr): fix fee distribution #737)
- **Analysis**: Valid finding. Exact match with our highest-priority finding. `distributeCDRRewardPool` iterates `map[string]uint64` without sorting keys. Fixed by adding `sort.Strings()` before iteration, matching the pattern in `dkg_rewards.go`.

---

### [CDR Report] Finding CDR-004: Partial Decryption Signature Missing Ethereum Signed Message Prefix
- **Severity**: P1 (Cryptographic Integrity)
- **Status**: VALID (partially overlaps with C-03)
- **Duplicate of**: Partially overlaps with our C-03 (CRYPTO-001) in FINAL_AUDIT_1.md regarding sign/verify field mismatch
- **Fixed in**: Not specifically fixed
- **Analysis**: Valid observation about inconsistency between `verifyPartialDecryptionSignature` (no EIP-191 prefix) and `verifyFinalizationSignature` (has EIP-191 prefix). Our C-03 focuses on a different aspect of the same signature path (field mismatch between kernel and CL). The CDR report's finding about prefix inconsistency is a separate valid concern. However, since both the kernel `signPartialDecryptResponse` and CL `verifyPartialDecryptionSignature` consistently omit the prefix, this is functionally correct — just inconsistent with the finalization path. Severity should be downgraded to MEDIUM (documentation/consistency issue rather than a functional break).

---

### [CDR Report] Finding CDR-005: DKG.sol `finalize()` Has No On-Chain Validation
- **Severity**: P2 (Integrity)
- **Status**: DUPLICATE
- **Duplicate of**: Our M-02 (CONTRACT-001) in FINAL_AUDIT_1.md, M-08 in AUDIT_REPORT.md
- **Fixed in**: Not yet fixed in fix branches
- **Analysis**: Valid finding. Our audit identified this as MEDIUM severity (downgraded from HIGH) because the CL verifies TEE attestation signatures before the contract call. The contract provides no defense-in-depth `validatorAddr == msg.sender` check. The CDR report's P2 rating aligns with our MEDIUM assessment.

---

### [CDR Report] Finding CDR-006: CDR.sol Condition Bypass When msg.sender IS the Condition Contract
- **Severity**: P2 (Integrity)
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid finding not covered in our audit. When `msg.sender == vault.writeConditionAddr` or `msg.sender == vault.readConditionAddr`, the condition check is bypassed. This is likely an intentional design choice (the condition contract itself is trusted), but it should be explicitly documented. If a condition contract is compromised or upgradeable, it gains unconditional vault access. LOW priority — design documentation issue rather than a code bug.

---

### [CDR Report] Finding CDR-007: Partial Decryption Signature Missing Code Commitment
- **Severity**: P2 (Integrity)
- **Status**: VALID — NEW FINDING (partially overlaps with STOR-1)
- **Duplicate of**: Overlaps with Apex STOR-1 regarding incomplete signature domain
- **Fixed in**: Not fixed
- **Analysis**: Valid finding. The partial decryption signature covers `round || ciphertext || encryptedPartial || ephPubKey || pubShare` but not `codeCommitment`. During rolling upgrades with two kernel binaries, signatures from different enclave builds would be indistinguishable. This is a real concern for upgrade safety and should be addressed alongside the broader signature domain expansion recommended in STOR-1.

---

### [CDR Report] Finding CDR-008: Pre-Release kyber v4 and Forked cb-mpc-go Dependencies
- **Severity**: P2 (Integrity)
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid supply chain concern. `kyber v4.0.0-pre2` is a pre-release version and `cb-mpc-go` is a Piplabs fork of a Coinbase library. Neither is covered by standard advisory databases. This is an operational/supply-chain risk rather than a code bug. Should be tracked as a long-term item: audit the specific commits used, diff the fork against upstream, and plan migration to stable kyber when available.

---

### [CDR Report] Finding CDR-009: No Witness Count Enforcement in Light Client Configuration
- **Severity**: P2 (Integrity)
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid finding. The light client at `story/query_client.go:87-94` accepts 0 witnesses with `SkippingVerification`. With 0 witnesses, a compromised primary RPC can serve forged state. However, the practical impact depends on deployment: if operators configure witnesses (as recommended in deployment docs), this is mitigated. Should enforce `len(witness_addrs) >= 1` at startup as defense-in-depth.

---

### [CDR Report] Finding CDR-010: `replayMessages()` Silently Ignores All Errors
- **Severity**: P2 (Integrity)
- **Status**: DUPLICATE / ALREADY FIXED
- **Duplicate of**: Our H-03 (KERNEL-005 / CRYPTO-002) in both audit reports
- **Fixed in**: `fix/audit-findings` commit `62b11f4` (story-kernel)
- **Analysis**: Valid finding. Exact match with our H-03. Fixed by logging errors at WARN level and propagating critical (non-duplicate) errors from replay.

---

### [CDR Report] Finding CDR-011: CDR Fee Pool Minting Without EVM-Side Burn Verification
- **Severity**: P2 (Economic)
- **Status**: DUPLICATE
- **Duplicate of**: Our M-14 (SECURITY-009/CONTRACT-011) in FINAL_AUDIT_1.md
- **Fixed in**: Not specifically fixed (design-level concern)
- **Analysis**: Valid concern. Our audit identified this as M-14 — the CDR fee pool dual tracking issue where CL mints based on EVM event amounts without verifying the EVM side actually burned equivalent value. The CDR report's P2 rating aligns with our MEDIUM assessment. Note: the gwei normalization issue (STOR-15) is the more severe manifestation of this same fee-bridging concern, and that HAS been fixed on `origin/dkg/dev`.

---

### [CDR Report] Finding CDR-012: Raw Label Bytes in Registry Key Causing Potential Collisions
- **Severity**: P2 (Integrity)
- **Status**: DUPLICATE / ALREADY FIXED
- **Duplicate of**: Our C-05 (COSMOS-002 / SECURITY-004) in FINAL_AUDIT_1.md, C-02 in AUDIT_REPORT.md
- **Fixed in**: `dkg/fix-audit-findings` commit `0e2b68e6`
- **Analysis**: Valid finding. Exact match with our CRITICAL finding. Fixed by using `hex.EncodeToString(label)` instead of raw `%s` interpolation, plus adding a nil/empty check.

---

### [CDR Report] Finding CDR-013: Unbounded In-Memory Cache Growth in TEE Kernel
- **Severity**: P2 (Availability)
- **Status**: DUPLICATE / ALREADY FIXED
- **Duplicate of**: Our M-10 (KERNEL-012) in FINAL_AUDIT_1.md, M-04 in AUDIT_REPORT.md
- **Fixed in**: `fix/audit-findings` commit `dda1c62` (story-kernel)
- **Analysis**: Valid finding. Our audit rated this as MEDIUM. Fixed by adding max size limits to DKG caches with LRU eviction.

---

### [CDR Report] Finding CDR-014: Broad `allowed_files` and `insecure__` Flags in Gramine Manifest
- **Severity**: P2 (Integrity)
- **Status**: DUPLICATE
- **Duplicate of**: Our M-11 (KERNEL-013) in FINAL_AUDIT_1.md, M-05 in AUDIT_REPORT.md
- **Fixed in**: Not fixed (operational/manifest change)
- **Analysis**: Valid finding. Our audit identified the same broad file access in the Gramine manifest. The CDR report adds detail about `insecure__use_cmdline_argv = true` which we also noted. This is a production hardening item.

---

### [CDR Report] Finding CDR-015: CDR.sol `read()` Reverts When readConditionAddr Is address(0)
- **Severity**: P3 (Liveness)
- **Status**: DUPLICATE / ALREADY FIXED
- **Duplicate of**: Our M-01 (CONTRACT-003) in FINAL_AUDIT_1.md
- **Fixed in**: `dkg/fix-audit-findings` commit `575442f4` (contracts/src/protocol/CDR.sol)
- **Analysis**: Valid finding. Our audit identified this as MEDIUM. Fixed by adding `require(readConditionAddr != address(0))` in `allocate()`.

---

### [CDR Report] Finding CDR-016: Session Files Written World-Readable (0644)
- **Severity**: P3 (Information Disclosure)
- **Status**: DUPLICATE / ALREADY FIXED
- **Duplicate of**: Our L-09 (DKG-SVC-015) in FINAL_AUDIT_1.md, L-05 in AUDIT_REPORT.md
- **Fixed in**: `dkg/fix-audit-findings` commit `27faae89`
- **Analysis**: Valid finding. Our audit rated this as LOW because the session files don't contain private key material. Fixed by changing to 0600 permissions.

---

### [CDR Report] Finding CDR-017: Integer Overflow in CDR Share Calculation (uint64 -> int64)
- **Severity**: P3 (Economic)
- **Status**: DUPLICATE
- **Duplicate of**: Our M-04 (COSMOS-009) in FINAL_AUDIT_1.md
- **Fixed in**: Not specifically fixed in fix branches (practically unreachable given validator counts)
- **Analysis**: Valid finding. Our audit rated this as MEDIUM. The `uint64` cast to `int64` via `math.NewInt(int64(count))` is technically unsafe for values > 2^63, but practically unreachable since `count` represents CDR partial submissions per validator per round (bounded by number of decrypt requests).

---

### [CDR Report] Finding CDR-018: CDR.sol Vault UUID Overflow (uint32)
- **Severity**: P4 (Economic)
- **Status**: DUPLICATE
- **Duplicate of**: Our L-08 (CONTRACT-005) in FINAL_AUDIT_1.md
- **Fixed in**: Not fixed
- **Analysis**: Valid but extremely low priority. Solidity 0.8.x reverts on overflow, so this would brick `allocate()` after ~4.3 billion vaults. At any realistic usage rate, this limit won't be reached for decades.

---

### [CDR Report] Finding CDR-019: DKG State Persisted as Plaintext JSON
- **Severity**: P5 (Information Disclosure)
- **Status**: DUPLICATE (FALSE POSITIVE in our assessment)
- **Duplicate of**: Our KERNEL-002 in FINAL_AUDIT_1.md (reclassified as FALSE POSITIVE)
- **Fixed in**: N/A
- **Analysis**: Our audit reclassified this as FALSE POSITIVE after detailed analysis. The DKG state file does NOT contain private keys — deals are encrypted, and justification SecShares are publicly broadcast data (they appear in on-chain vote extensions). The CDR report's concern about "plaintext SecShare scalars" is technically correct for justification data, but those values are already public information once broadcast.

---

### [CDR Report] Finding CDR-020: Go GC May Copy Secrets Before zeroBytes() Wipe
- **Severity**: P5 (Information Disclosure)
- **Status**: DUPLICATE
- **Duplicate of**: Our M-17 (CRYPTO-006) in FINAL_AUDIT_1.md
- **Fixed in**: Not fixed
- **Analysis**: Valid theoretical concern, mitigated by SGX. Our audit rated this as MEDIUM (mitigated by SGX enclave isolation). The Go GC may copy heap objects before `defer zeroBytes()` executes, and kyber scalar internals are never zeroed. However, within an SGX enclave, memory is encrypted and isolated from the host, so this is only exploitable after an SGX breach.

---

## CDR Supplementary Findings — Smart Contract

---

### [CDR Report] Finding SC-01: DKG.sol `register()` Duplicate Registrations
- **Severity**: Medium
- **Status**: VALID — PARTIALLY ADDRESSED
- **Duplicate of**: None directly, but related to STOR-11 (duplicate key reuse)
- **Fixed in**: `origin/dkg/dev` commit `62ce0fd1` ("add defense in depth guard against duplicate registration")
- **Analysis**: Valid finding. The EVM contract has no on-chain state tracking which validators registered. A defense-in-depth guard was added on `dkg/dev` to reject duplicate registrations.

---

### [CDR Report] Finding SC-02: Fee Burning via `transfer()` to address(0) May Brick
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid concern. `payable(address(0x0)).transfer(amount)` forwards only 2300 gas. If Story ever deploys a precompile at address(0), fee operations break. LOW practical risk since address(0) precompile deployment is not planned, but using `.call{value: amount}("")` or a dedicated burn address is better practice.

---

### [CDR Report] Finding SC-03: SGXValidationHook Inline Assembly Lacks Bounds Check
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid defense-in-depth concern. `_extractReportInstanceDataCommitment` uses `mload` at offset 368 without checking `enclaveReport.length >= 432`. Unlike `_extractReportCodeCommitment` which uses `BytesUtils.substring` with bounds checking. Automata DCAP should reject short reports upstream, but this is inconsistent and should have the same bounds check.

---

### [CDR Report] Finding SC-04: Vault Conditions Immutable But Condition Contracts May Be Upgradeable
- **Severity**: Medium (Design)
- **Status**: VALID — NEW FINDING (design documentation)
- **Duplicate of**: Related to CDR-006 (condition bypass)
- **Fixed in**: Not fixed
- **Analysis**: Valid design observation. If a condition contract is a proxy, its logic can change without the vault's knowledge. Should be documented as a known design characteristic, with warnings against using upgradeable proxies as condition addresses.

---

### [CDR Report] Finding SC-05: DKG.sol `whitelistEnclaveType()` Requires Valid Data When De-Whitelisting
- **Severity**: Low
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid minor UX issue. Owner must provide dummy valid data to de-whitelist. Easy fix: only validate when `isWhitelisted = true`.

---

### [CDR Report] Finding SC-06: `submitEncryptedPartialDecryption()` Has No On-Chain Validation
- **Severity**: Low
- **Status**: DUPLICATE (partially)
- **Duplicate of**: Related to Apex STOR-4 and STOR-6 (permissionless partial submissions)
- **Fixed in**: Not fixed at contract level
- **Analysis**: Valid finding. The contract accepts arbitrary payloads with fee as only spam deterrent. CL-side validation provides the security boundary, but on-chain validation (e.g., `require(vaults[uuid].encryptedData.length > 0)`) would be better defense-in-depth.

---

### [CDR Report] Finding SC-07: SGXValidationHook Extraction Functions Not Marked `pure`
- **Severity**: Low
- **Status**: DUPLICATE
- **Duplicate of**: Our L-06 (CONTRACT-009) in FINAL_AUDIT_1.md
- **Fixed in**: Not fixed
- **Analysis**: Valid minor issue. Functions should be marked `internal pure`.

---

### [CDR Report] Finding SC-08: DKG.sol Missing ReentrancyGuard
- **Severity**: Low
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid defense-in-depth concern. CDR.sol has `ReentrancyGuardUpgradeable` but DKG.sol does not, despite making external calls to `IAttestationReportValidator.validateReport()`. LOW risk since the validation hook is controlled, but inconsistent with CDR.sol's approach.

---

### [CDR Report] Finding SC-09: CDR.sol `allocate()` Missing `nonReentrant`
- **Severity**: Low
- **Status**: DUPLICATE
- **Duplicate of**: Our L-07 (CONTRACT-004) in FINAL_AUDIT_1.md
- **Fixed in**: Not fixed
- **Analysis**: Valid consistency issue. Our audit rated this as LOW since `allocate()` doesn't transfer to `address(0)` in a way that could callback.

---

### [CDR Report] Finding SC-10: No Timelock Enforced in `_authorizeUpgrade`
- **Severity**: Info
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid governance concern. Tests show owner is a `TimelockController`, but the contract doesn't enforce it. Documentation issue.

---

### [CDR Report] Finding SC-11: Test Coverage Critically Insufficient
- **Severity**: Info
- **Status**: DUPLICATE
- **Duplicate of**: General observation in our Test Coverage section
- **Fixed in**: Partially addressed in `dkg/fix-audit-findings` (some new tests added, but many old tests removed)
- **Analysis**: Valid observation.

---

### [CDR Report] Finding SC-12: Gas Inefficiency — Full Vault Struct Copied to Memory
- **Severity**: Info (Gas)
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid optimization suggestion. Using storage pointers instead of memory copies would save gas for vaults with large `encryptedData`. LOW priority.

---

## CDR Supplementary Findings — TEE Kernel

---

### [CDR Report] Finding TEE-05: HKDF Used Without Salt in Partial Decryption Encryption
- **Severity**: Medium
- **Status**: DUPLICATE
- **Duplicate of**: Our L-12 (CRYPTO-008) in FINAL_AUDIT_1.md (our KERNEL-009 was reclassified as FP)
- **Fixed in**: Not fixed (RFC 5869 compliant, no action needed)
- **Analysis**: Our audit concluded this is RFC 5869 compliant. The CDR report acknowledges it's "secure when IKM has sufficient entropy (ECDH does)" and recommends salt as an optional improvement. We agree — no action required.

---

### [CDR Report] Finding TEE-06: TLS Not Enforced by Default
- **Severity**: Low
- **Status**: DUPLICATE / ALREADY FIXED
- **Duplicate of**: Our M-12 (KERNEL-014) in FINAL_AUDIT_1.md
- **Fixed in**: `origin/dkg/dev` commit `f1fcc745` ("add tls and mtls support for grpc connections to the kernel")
- **Analysis**: Valid finding. TLS/mTLS support has been added since the audited commit.

---

### [CDR Report] Finding TEE-07: Panic Recovery Interceptor May Mask Crypto State Corruption
- **Severity**: Low
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid concern. The panic recovery interceptor at `server.go:192-208` catches all panics including those from kyber/cb-mpc. For crypto panics specifically, allowing the enclave to crash may be safer than recovering with potentially corrupted state. LOW priority — the practical risk is small since crypto panics in production are rare.

---

### [CDR Report] Finding TEE-08: `sync.Map` Per-Round Mutexes Never Pruned
- **Severity**: Low
- **Status**: DUPLICATE
- **Duplicate of**: Included in our M-10 (unbounded caches) assessment
- **Fixed in**: Partially addressed in `fix/audit-findings` commit `dda1c62` (cache size limits added)
- **Analysis**: Valid minor concern. Each entry is ~16 bytes so practical impact is minimal. The cache size limit fix partially addresses this.

---

## CDR Supplementary Findings — L1 Module

---

### [CDR Report] Finding L1-01: `getDKGRegistrationsByRound` Full Table Scan
- **Severity**: Medium (Performance)
- **Status**: DUPLICATE / ALREADY FIXED
- **Duplicate of**: Our M-03 (COSMOS-005) in FINAL_AUDIT_1.md, M-09 in AUDIT_REPORT.md
- **Fixed in**: `dkg/fix-audit-findings` commit `4183a3be`
- **Analysis**: Valid finding. Fixed by using prefix range for round-scoped registration queries.

---

### [CDR Report] Finding L1-02: No Size Limit on Ciphertext in Decrypt Registry
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid concern. `ThresholdDecryptRequested` stores full `ciphertext` in state without size limits. Large ciphertexts could bloat consensus state. Should enforce a maximum ciphertext size (e.g., 1MB). MEDIUM priority.

---

### [CDR Report] Finding L1-03: `GetAllCodeCommitments` Non-Deterministic Map Iteration
- **Severity**: Medium (Off-Chain)
- **Status**: DUPLICATE / ALREADY FIXED
- **Duplicate of**: Our M-15 (SECURITY-010) in FINAL_AUDIT_1.md, M-06 in AUDIT_REPORT.md
- **Fixed in**: `dkg/fix-audit-findings` commit `b0d0d68b`
- **Analysis**: Valid finding. Fixed by sorting code commitments for deterministic selection.

---

### [CDR Report] Finding L1-04: `participantsRoot` Validation With Invalidated Dealers
- **Severity**: Medium (Documentation)
- **Status**: VALID — NEW FINDING (documentation)
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid documentation concern. The behavior (excluding invalidated dealers by virtue of their status) is correct but implicit. A comment should be added.

---

### [CDR Report] Finding L1-05: `UpgradeScheduled` Does Not Validate Activation Height vs Current DKG Round End
- **Severity**: Low
- **Status**: VALID — NEW FINDING
- **Duplicate of**: Our SECURITY-007 was declared FALSE POSITIVE (height validation exists), but this is a different concern
- **Fixed in**: Not fixed
- **Analysis**: Valid concern. The activation height is validated against current block height, but not against whether it's far enough in the future to complete the current DKG round. Activating mid-round forces a `SkipToNextRound`. LOW priority — operationally, upgrade scheduling should account for round timing.

---

### [CDR Report] Finding L1-06: `proposalServer.AddVote` Does Not Validate Vote Content
- **Severity**: Low
- **Status**: VALID — NEW FINDING (but intentional design)
- **Duplicate of**: Related to Apex STOR-25 (proposer censorship)
- **Fixed in**: Not fixed
- **Analysis**: This is intentional per ABCI++ pattern — `PrepareProposal` is permissive, `FinalizeBlock` validates. However, as Apex STOR-25 points out, this creates a censorship vector. A comment documenting the intentional design would help.

---

### [CDR Report] Finding L1-07: `dkgSvcRound` Global Atomic Shared Across Test Instances
- **Severity**: Low (Testing)
- **Status**: DUPLICATE
- **Duplicate of**: Part of our H-07 (DKG-SVC-002 / SECURITY-003) — package-level global state
- **Fixed in**: Not yet fully fixed (part of the larger refactor)
- **Analysis**: Valid testing concern, subset of our H-07 finding about package-level globals.

---

## Apex Report Findings

---

### [Apex Report] Finding STOR-15: CDR Fee Bridging 1e9 Reward Over-Credit (wei vs gwei)
- **Severity**: High
- **Status**: VALID — NEW FINDING / ALREADY FIXED
- **Duplicate of**: Tangentially related to our M-14 (CDR fee pool dual tracking), but this is a distinct and much more severe issue
- **Fixed in**: `origin/dkg/dev` — the gwei normalization (`ev.Amount.Div(ev.Amount, gwei)`) is present in current `cdr.go` on dkg/dev. Was NOT present at audited commit `11307a5`.
- **Analysis**: **This is a genuinely new HIGH-severity finding not identified in our audit.** At the audited commit, `ProcessCDRFeeCollected` forwarded raw wei-denominated amounts from EVM events into CL stake minting without the gwei division that the staking bridge applies. This created a 1e9 amplification: 1 wei burned on EL became 1 gwei-worth of CL stake. The fix was merged to dkg/dev (visible in current code). Our M-14 noted the "dual tracking" concern but missed the critical unit mismatch. **Credit to Apex for finding this.**

---

### [Apex Report] Finding STOR-6: Validators Can Submit Arbitrary or Replayed Partial Decryptions
- **Severity**: High
- **Status**: DUPLICATE / PARTIALLY FIXED
- **Duplicate of**: Our H-06 (SECURITY-002) in FINAL_AUDIT_1.md
- **Fixed in**: The `accepted` return value fix on `origin/dkg/dev` (`a56fbc5f`) prevents rewarding non-validated submissions
- **Analysis**: Valid finding. Our audit identified the lack of on-chain TDH2 correctness verification as H-06. The Apex report provides much more detailed exploitation analysis showing how arbitrary or replayed partials could collect refunds/rewards. The core issue (CL doesn't verify cryptographic correctness of partials) remains unfixed by design (expensive on-chain), but the immediate reward-accounting vulnerability (nil return = success) has been fixed.

---

### [Apex Report] Finding STOR-5: Unauthenticated PartialDecryptTDH2 — Remote Threshold-Decryption Oracle
- **Severity**: High
- **Status**: DUPLICATE / PARTIALLY FIXED
- **Duplicate of**: Our M-12/M-13 (KERNEL-014/KERNEL-015) in FINAL_AUDIT_1.md, and CDR-002 above
- **Fixed in**: mTLS support added (`f1fcc745`), canonical chain validation RPC added (`600bf54`, `6a4c4174`)
- **Analysis**: Valid finding. The Apex report provides a detailed attack path showing how unauthenticated gRPC access to kernel enables a remote decryption oracle. Our audit identified the components (no auth M-12, no on-chain verification M-13) but Apex ties them together into a complete attack. The mTLS addition and on-chain verification RPC partially address this.

---

### [Apex Report] Finding STOR-4: Unknown or Expired Partial Submissions Rewarded as Valid Work
- **Severity**: High
- **Status**: VALID / ALREADY FIXED
- **Duplicate of**: Closely related to our H-06 but identifies a distinct attack vector (not-found/timed-out returns nil = success)
- **Fixed in**: `origin/dkg/dev` commit `a56fbc5f` ("avoid adding submission count when partials are not accepted"). `PartialDecryptionSubmitted` now returns `(bool, error)` and the caller checks `accepted && partialErr == nil`.
- **Analysis**: **This is a genuinely important finding.** At the audited commit, `PartialDecryptionSubmitted` returned `nil` for not-found, timed-out, and duplicate submissions, and the caller treated `nil` as success. This allowed any external account to submit garbage partials and collect refunds/rewards. The fix changes the function signature to return `(bool, error)` where `accepted=true` only for genuinely valid submissions. **Well-caught by Apex.**

---

### [Apex Report] Finding STOR-3: Duplicate Partial Replays Accrue Fresh Refunds
- **Severity**: High
- **Status**: VALID / ALREADY FIXED
- **Duplicate of**: Subset of STOR-4 (same root cause: nil return = success)
- **Fixed in**: Same fix as STOR-4 (`a56fbc5f`)
- **Analysis**: Valid finding, same root cause as STOR-4. Duplicate detection exists in storage but returned `nil` to the caller, which was treated as success. Fixed by the `(bool, error)` return value change.

---

### [Apex Report] Finding STOR-29: Rebooting During Active Round Permanently Disables Threshold Decryption
- **Severity**: Medium
- **Status**: DUPLICATE
- **Duplicate of**: Closely related to our C-04 (DKG-SVC-007) — Decrypt Worker Context Leak
- **Fixed in**: Not fully fixed
- **Analysis**: Valid finding. Our C-04 identified the decrypt worker context leak (1-min timeout kills worker). STOR-29 focuses on the recovery side: `ResumeDKGService` doesn't restart the worker for `PhaseCompleted` sessions during `DKGStageActive`. These are two aspects of the same problem. Both need to be addressed for CDR to be functional.

---

### [Apex Report] Finding STOR-28: Decrypt Worker Dies After One Minute
- **Severity**: Medium
- **Status**: DUPLICATE
- **Duplicate of**: Our C-04 (DKG-SVC-007) in FINAL_AUDIT_1.md
- **Fixed in**: Not fixed
- **Analysis**: Exact same finding as our C-04. The `dkgAsyncContext()` provides a 1-minute timeout, but the decrypt worker needs to run for the entire 21-day active period. Apex provides more detailed impact analysis showing the economic harm (fee-for-no-service).

---

### [Apex Report] Finding STOR-27: Final EL-Block Partial Submissions Settled in Wrong Committee Epoch
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None (related pattern to STOR-22/STOR-19/STOR-23 but distinct)
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** CDR fee-pool settlement runs in `BeginBlocker` before `MsgExecutionPayload` processes the previous EL block's partial submission events. Late submissions from committee R are cleared during settlement, then re-populate `CDRPartialSubmitCount` from the next `DeliverTx`, causing them to be paid in the next epoch. This is a timing/ordering issue inherent to the BeginBlock-before-DeliverTx execution model. MEDIUM priority — the economic impact is bounded to the final EL block's submissions per epoch.

---

### [Apex Report] Finding STOR-26: Permissionless Pre-v2 CDR Fees Mutate Future DKG Store
- **Severity**: Medium
- **Status**: VALID / ALREADY FIXED
- **Duplicate of**: None
- **Fixed in**: `origin/dkg/dev` commit `efd59a56` ("remove mounted stores filter from upgrade store loader")
- **Analysis**: **Valid new finding at time of audit.** The upgrade store loader pre-mounted future stores, and CDR fee events were not gated on activation height. The fix removes this pre-mounting behavior. Additionally, `ProcessDKGEvents`/`ProcessCDREvents` should ideally be gated on DKG activation, which commit `0382ec7f` ("remove redundant checks and support on chain upgrade scheduling") may address.

---

### [Apex Report] Finding STOR-25: Single Proposer Can Censor DKG Complaints
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None (related to STOR-7 and STOR-10 but distinct attack vector)
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** `ProcessProposal` only checks that one `MsgAddDkgVote` exists with correct authority — it never verifies the vote body matches `ProposedLastCommit`. A Byzantine proposer can omit complaints from the aggregated vote. This is a fundamental design concern with the current ABCI++ integration. Combined with the destructive dequeue pattern (items removed from queues during `ExtendVote`), omitted items are lost forever. MEDIUM-HIGH priority — requires architectural change to either verify proposal-time DKG aggregation against committed vote extensions, or make queue removal acknowledgement-based.

---

### [Apex Report] Finding STOR-24: Global First-10 Truncation Discards Valid Justifications
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** `ProcessJustifications` silently truncates to `MaxJustificationsPerBlock=10`, but vote extensions allow up to 80 per validator. Combined with the destructive dequeue, overflow justifications are permanently lost. A colluding set can flood the first 10 slots to protect malicious dealers from invalidation. MEDIUM priority — should either reject proposals exceeding the budget or carry overflow into subsequent blocks.

---

### [Apex Report] Finding STOR-23: Dealing-to-Finalization Transition Drops Last Block's DKG Complaints
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None (same class as STOR-22/STOR-19 — stage-boundary timing bugs)
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** On the first block of `DKGStageFinalization`, `BeginBlocker` has already advanced the stage, but `MsgAddDkgVote` still carries the previous dealing block's vote extensions. `msgServer.AddVote` checks the current stage and silently skips the payload because it's no longer `DKGStageDealing`. This deterministically drops the final dealing block's complaints. This is a systematic issue — the stage transition happens before the previous block's vote extensions are processed. MEDIUM-HIGH priority.

---

### [Apex Report] Finding STOR-22: First Active Block Rejects Last Finalization Block's DKG Finalize Events
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None (same class as STOR-23/STOR-19)
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** Same pattern as STOR-23 but for the Finalization-to-Active transition. `BeginBlocker` activates the round, but `MsgExecutionPayload` still carries the last finalization block's `DKG.Finalized` events. `dkgKeeper.Finalized` rejects them because stage is already `Active`. Honest validators who finalized in the last EL block are excluded from the committee. MEDIUM-HIGH priority — this is a systematic issue with the BeginBlocker-before-DeliverTx execution model.

---

### [Apex Report] Finding STOR-21: Resharing Mutates Serving-Round Pointer
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None (closely related to STOR-16, same root cause)
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** When active period expires, `BeginBlocker` mutates the current round's stage to `Registration` but `LatestActiveRound` still points to it. CDR reads during resharing are accepted and charged but the honest auto-queue skips them because the round is no longer `Active`. MEDIUM priority — same root cause as STOR-16 (stale `LatestActiveRound` pointer).

---

### [Apex Report] Finding STOR-20: PartialDecryptTDH2 Protobuf Mismatch Between CL and Kernel
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: Related to our C-03 (sign/verify field mismatch) but a distinct protobuf schema-level incompatibility
- **Fixed in**: Not fixed
- **Analysis**: **Critical new finding.** The CL proto has `pid` at field 5, `global_pub_key` at field 6, `requester_pub_key` at field 8. The kernel proto has `global_pub_key` at field 5, `requester_pub_key` at field 6, no field 8. This is a wire-type mismatch (varint vs bytes for field 5) that makes every CL-to-kernel decrypt RPC malformed. This should be HIGH severity — it means CDR decryption is broken at the protobuf layer. The proto schemas must be unified.

---

### [Apex Report] Finding STOR-19: First Dealing Block Rejects Last Registration Block's DKG Register Events
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None (same class as STOR-22/STOR-23)
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** Same systematic pattern: Registration-to-Dealing transition in `BeginBlocker` happens before previous block's `DKG.Registered` events are processed. Valid registrations from the last EL registration block are rejected. MEDIUM-HIGH priority — part of the systematic stage-boundary timing issue.

---

### [Apex Report] Finding STOR-18: VerifiedQueryClient Accepts Proofs for Arbitrary Keys
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** The kernel's `queryWithProof()` checks that `resp.Key` is non-empty but never verifies it matches the requested key. `verifyProof()` trusts proof-supplied keys rather than the caller's query key. A hostile RPC can therefore answer with a valid proof for a different key. This breaks the core verified-query trust boundary. MEDIUM-HIGH priority — requires binding the proof chain to the caller's expected key at every layer.

---

### [Apex Report] Finding STOR-16: Expired DKG Rounds Keep Collecting UBI
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None (related to STOR-21 — same stale `LatestActiveRound` root cause)
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** When active period expires, `LatestActiveRound` is not cleared until the successor round finalizes. UBI rewards are sent to whatever `LatestActiveRound` returns without checking stage. Self-unstake is also blocked. A coalition that keeps resharing from finalizing can farm UBI indefinitely. MEDIUM priority — requires clearing/updating `LatestActiveRound` at rollover.

---

### [Apex Report] Finding STOR-13: Finalization Stores Raw pubKeyShare While Partial Decryptions Submit Prefixed pubShare
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None (different from our C-03 which is about sign/verify field mismatch)
- **Fixed in**: Not fixed
- **Analysis**: **Critical new finding that should be HIGH.** During finalization, kernel returns raw Edwards25519 point bytes stored as `PubKeyShare`. During partial decryption, kernel prepends `[0x04, 0x3f]` TDH2 prefix. CL verifier enforces `bytes.Equal(pubShare, reg.PubKeyShare)` which always fails for honest validators. This means the CDR read path is broken even if the protobuf mismatch (STOR-20) is fixed. Must normalize to the same representation on both sides.

---

### [Apex Report] Finding STOR-11: Reusing One TEE Keypair Across Validator Addresses
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** `story-kernel` seals keys only by `(round, codeCommitment)`, not by validator address. A multi-validator operator can register the same TEE keypair for multiple validator addresses, inflating the finalized count beyond the actual number of unique key shares. The committee appears to have met threshold but may not be able to decrypt if unique shares < threshold. MEDIUM priority — requires namespacing key paths by validator address.

---

### [Apex Report] Finding STOR-10: Early-Sorted Validator Can Suppress Later DKG Complaints
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: Same class as STOR-7 but specifically for responses
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** Vote-extension deduplication uses first-wins on `(dealer, responder)` pairs. An earlier-sorted Byzantine validator can include forged response entries in its vote extension that consume the dedup slots for later honest validators' real complaints. The forged entries fail kernel signature checks, but by then the real complaints are already gone. Combined with STOR-25 (proposer censorship), this creates a powerful complaint suppression mechanism.

---

### [Apex Report] Finding STOR-9: CDR Read Path Dead in Every DKG Phase
- **Severity**: Medium
- **Status**: DUPLICATE
- **Duplicate of**: Combination of our C-04 (decrypt worker context leak) and STOR-21/STOR-28/STOR-29
- **Fixed in**: Not fixed
- **Analysis**: Valid finding that ties together multiple issues: (1) active rounds kill the decrypt worker via context leak (our C-04), and (2) resharing rounds reject request queueing because the round stage was overwritten. This is a comprehensive statement of the CDR read path being broken in all phases.

---

### [Apex Report] Finding STOR-8: Resharing Rounds Never Populate PIDCache
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: **Critical new finding.** `GenerateDeals()` explicitly skips `CachePID()` when `isResharing == true`. Every round after genesis is a resharing round. `PartialDecryptTDH2()` requires `PIDCache` and ignores the `pid` field in the request. Result: after the first resharing, CDR decryption permanently fails. This should be HIGH severity — it means CDR is broken after the first round rotation.

---

### [Apex Report] Finding STOR-2: Upgrade Resharing Rejects New Validators
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: **Valid new finding.** `handleDKGRegistration` requires `oldCC` (previous round's code commitment) for upgrade rounds, but new validators have no previous registration. `getOldCodeCommitment` fails for them, and `getRegistrationKernelClient` rejects them with "old code commitment required". The round is initialized from the current bonded set but new validators can't participate. MEDIUM priority — should allow new validators to register directly with the new binary.

---

### [Apex Report] Finding STOR-1: Partial-Decrypt Signatures Omit Requester and UUID Binding
- **Severity**: Medium
- **Status**: VALID — NEW FINDING
- **Duplicate of**: Related to CDR-004/CDR-007 (incomplete signature domain) but distinct attack
- **Fixed in**: Not fixed
- **Analysis**: **Valid HIGH-severity finding.** The kernel signature covers only `round || ciphertext || encryptedPartial || ephPubKey || pubShare`, omitting `requesterPubKey`, `label/UUID`, and `pid`. A hostile keeper can swap `requesterPubKey` when calling the enclave, receive a partial encrypted to the attacker's key, and submit it to the contract with the victim's key. CL accepts it because the forged context is outside the signed domain. This enables plaintext disclosure across threshold-many malicious validators.

---

## Low-Severity Apex Findings

---

### [Apex Report] Finding STOR-17: Sealed Light-Client State Is Rollbackable
- **Severity**: Low
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid concern about rollback protection. Sealing provides confidentiality/integrity but not freshness. A hostile host can snapshot and restore older sealed state. Requires a monotonic counter or append-only checkpoint for full rollback protection. LOW priority given the threat model already assumes host control requires additional RPC compromise.

---

### [Apex Report] Finding STOR-14: Honest Signed DKG Deals Can Be Turned Into Undeliverable Shares
- **Severity**: Low
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid finding about kyber's narrow signature domain. `dkg.Deal.MarshalBinary()` signs only `(dealerIndex, cipher)`, not `recipientIndex` or `nonce`. A malicious proposer can alter these unsigned fields to misroute honest deals. Combined with vote-extension one-shot delivery, the honest dealer is permanently excluded from `QUAL()`. MEDIUM priority (Apex rates LOW but combined with STOR-25, the practical impact is significant).

---

### [Apex Report] Finding STOR-12: Whitelisting New Enclave Type Bypasses scheduleUpgrade
- **Severity**: Low
- **Status**: VALID — NEW FINDING
- **Duplicate of**: None
- **Fixed in**: Not fixed
- **Analysis**: Valid finding. CL never enforces that a registration's code commitment matches the expected one for the current round type (upgrade vs ordinary). Once a second enclave type is whitelisted, validators can register with it in an ordinary round without going through `scheduleUpgrade`. LOW priority — requires the admin to whitelist the new type first, which is already an authorized action.

---

### [Apex Report] Finding STOR-7: First-Seen Vote Dedup Lets Earlier Validator Erase Later DKG Messages
- **Severity**: Low
- **Status**: VALID — NEW FINDING
- **Duplicate of**: Same class as STOR-10 but for deals
- **Fixed in**: Not fixed
- **Analysis**: Valid finding. Same first-wins deduplication vulnerability as STOR-10, applied to both deals and responses. Does not require proposer control — any earlier-sorted validator can win the dedup slot. MEDIUM priority (Apex rates LOW but the DKG transcript integrity impact is significant).

---

## Summary of New Findings Not In Our Audit

| Finding | Source | Severity | Priority |
|---------|--------|----------|----------|
| STOR-15: CDR fee wei/gwei 1e9 amplification | Apex | HIGH | **ALREADY FIXED** on dkg/dev |
| STOR-4: Not-found/timed-out partials rewarded | Apex | HIGH | **ALREADY FIXED** on dkg/dev |
| STOR-20: Protobuf field mismatch CL↔kernel | Apex | MEDIUM (should be HIGH) | **UNFIXED — CRITICAL** |
| STOR-13: pubKeyShare prefix mismatch | Apex | MEDIUM (should be HIGH) | **UNFIXED — CRITICAL** |
| STOR-8: PIDCache not populated for resharing | Apex | MEDIUM (should be HIGH) | **UNFIXED — CRITICAL** |
| STOR-25: Proposer can censor DKG complaints | Apex | MEDIUM | UNFIXED |
| STOR-23: Stage boundary drops last block complaints | Apex | MEDIUM | UNFIXED |
| STOR-22: Stage boundary drops last finalization events | Apex | MEDIUM | UNFIXED |
| STOR-19: Stage boundary drops last registration events | Apex | MEDIUM | UNFIXED |
| STOR-27: Final EL-block partials in wrong epoch | Apex | MEDIUM | UNFIXED |
| STOR-26: Pre-v2 CDR fees mutate DKG store | Apex | MEDIUM | **ALREADY FIXED** on dkg/dev |
| STOR-21: Resharing mutates serving-round pointer | Apex | MEDIUM | UNFIXED |
| STOR-16: Expired rounds keep collecting UBI | Apex | MEDIUM | UNFIXED |
| STOR-18: VerifiedQueryClient key/store binding | Apex | MEDIUM | UNFIXED |
| STOR-11: TEE keypair reuse across validators | Apex | MEDIUM | UNFIXED |
| STOR-10: Vote dedup complaint suppression | Apex | MEDIUM | UNFIXED |
| STOR-2: Upgrade resharing rejects new validators | Apex | MEDIUM | UNFIXED |
| STOR-1: Partial-decrypt sig omits requester binding | Apex | MEDIUM (should be HIGH) | UNFIXED |
| STOR-14: Unsigned deal routing fields | Apex | LOW | UNFIXED |
| STOR-17: Sealed state rollbackable | Apex | LOW | UNFIXED |
| STOR-12: Whitelist bypasses scheduleUpgrade | Apex | LOW | UNFIXED |
| STOR-7: First-seen vote dedup censorship | Apex | LOW | UNFIXED |
| CDR-006: Condition contract bypass | CDR | P2 | UNFIXED (design doc) |
| CDR-007: Partial decrypt sig missing CC | CDR | P2 | UNFIXED |
| CDR-008: Pre-release kyber/forked deps | CDR | P2 | UNFIXED (supply chain) |
| CDR-009: No witness enforcement in light client | CDR | P2 | UNFIXED |
| SC-02: Fee burn via transfer to address(0) | CDR | Medium | UNFIXED |
| SC-03: SGXValidationHook missing bounds check | CDR | Medium | UNFIXED |
| SC-04: Upgradeable condition contracts | CDR | Medium | UNFIXED (design doc) |
| SC-08: DKG.sol missing ReentrancyGuard | CDR | Low | UNFIXED |
| SC-10: No timelock enforced in _authorizeUpgrade | CDR | Info | UNFIXED (design doc) |
| L1-02: No size limit on ciphertext | CDR | Medium | UNFIXED |
| L1-05: UpgradeScheduled height validation | CDR | Low | UNFIXED |
| TEE-07: Panic recovery masks crypto corruption | CDR | Low | UNFIXED |

---

## Key Observations

### 1. CDR Read Path Is Fundamentally Broken (Multiple Independent Bugs)
The Apex report identifies at least 4 independent bugs that each individually break CDR threshold decryption:
- **STOR-20**: Protobuf schema mismatch between CL and kernel
- **STOR-13**: pubKeyShare prefix mismatch (finalization vs partial decrypt)
- **STOR-8**: PIDCache not populated for resharing rounds
- **C-04/STOR-28**: Decrypt worker dies after 1 minute

These are separate from the sign/verify field mismatch (our C-03). Even fixing C-03 alone would not make CDR work.

### 2. Systematic Stage-Boundary Timing Bug
STOR-19, STOR-22, STOR-23, and STOR-27 all share the same root cause: `BeginBlocker` advances the DKG stage before `MsgExecutionPayload` processes the previous block's EVM events. This is a fundamental ordering issue in the ABCI lifecycle that affects all stage transitions (Registration→Dealing, Dealing→Finalization, Finalization→Active). The fix requires either delaying stage transitions until after previous-block events are processed, or accepting previous-stage events during the transition block.

### 3. Vote Extension Security Model Has Gaps
STOR-7, STOR-10, STOR-24, and STOR-25 collectively show that the vote extension transport is vulnerable to censorship and manipulation by both earlier-sorted validators (dedup exploitation) and Byzantine proposers (vote body substitution). The core issue is that DKG message authenticity is verified too late — after deduplication has already discarded honest entries.

### 4. `LatestActiveRound` Stale Pointer Problem
STOR-16 and STOR-21 identify the same root cause: `LatestActiveRound` is not cleared or updated when the active period expires. Downstream consumers (rewards, self-unstake, CDR reads) trust this stale pointer, creating economic misalignment during resharing.

### 5. Fixes Already Merged Address Critical Economic Bugs
The most immediately exploitable economic vulnerabilities have already been fixed:
- Wei/gwei amplification (STOR-15) — fixed on dkg/dev
- Rewarding invalid/duplicate/expired partials (STOR-3/STOR-4) — fixed on dkg/dev
- Pre-v2 DKG store mutation (STOR-26) — fixed on dkg/dev

### 6. Quality of Apex Report
The Apex report is exceptionally thorough and identifies several genuinely new critical findings that our 6-agent audit missed. The analysis of stage-boundary timing bugs, vote-extension security gaps, and CDR read-path failures demonstrates deep understanding of the ABCI lifecycle and cross-layer interactions. Several findings rated as "Medium" by Apex should be considered HIGH given their impact on CDR functionality (STOR-8, STOR-13, STOR-20).
