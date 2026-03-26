# DKG/CDR Production Final Audit Report

**Date:** 2026-03-23
**Auditors:** 6 specialized agents (Cosmos SDK, DKG Service, Smart Contract, TEE/SGX, Cryptography, Cross-cutting Security)
**Design Docs:** `story/docs/design/DKG.md`, `story/docs/design/CDR.md`
**Scope:**
- story: branch `origin/dkg/dev`
- story-kernel: branch `origin/main`

**Overall Risk Rating: HIGH** — Multiple CRITICAL findings must be resolved before production deployment.

---

## Findings Summary

| Severity | Count |
|----------|-------|
| CRITICAL | 3 |
| HIGH | 6 |
| MEDIUM | 9 |
| LOW | 8 |

---

## CRITICAL Findings (3)

### C-01: CDR Fee Distribution — Non-Deterministic Map Iteration [CONSENSUS FAILURE]

**Location:** `client/x/dkg/keeper/dkg_cdr_fees.go:182`

**Description:** `distributeCDRRewardPool` iterates over `map[string]uint64` to send bank transfers. Go map iteration order is non-deterministic, meaning different validators iterate in different orders, producing different state transitions and different app hashes.

**Evaluation:** The DKG reward distribution in `dkg_rewards.go` correctly uses `sort.Strings(memberAddrs)` before iteration. The CDR fee distribution does not follow the same pattern. This guarantees a chain halt when the CDR fee pool has a balance and multiple validators contributed partial decryptions.

**Fix:** Sort map keys before iteration, matching the pattern already used in `dkg_rewards.go`.

---

### C-02: Decrypt Request Registry Key Collision

**Location:** `client/x/dkg/keeper/dkg_decrypt_registry.go:17-27`

**Description:** The registry key is formed using `fmt.Sprintf("%s_%s_%d_%s", requesterHash, label, round, ciphertextHash)` where `label` is raw bytes interpolated via `%s`. Since `label` can contain underscore characters or other delimiters, two different `(label, round)` pairs can produce the same key string. For example, `label="a_1_b", round=2` and `label="a", round=1` with subsequent field `"b_2_..."` can collide. Additionally, there is no nil check on the label input.

**Evaluation:** The design doc (Appendix A) specifies `"{requesterHash}_{labelHex}"` using hex encoding, but the implementation uses raw bytes. This creates both a key collision vulnerability and a deviation from the spec.

**Fix:** Use `hex.EncodeToString(label)` instead of raw `%s` interpolation. Add a nil/empty check for label at the top of the function.

---

### C-03: SGX Debug Mode Enabled in Production Manifest

**Location:** `story-kernel.manifest.template:51`

**Description:** The Gramine SGX manifest has `sgx.debug = true`, which allows the host OS to inspect enclave memory. All private keys, DKG polynomial coefficients, and threshold decryption shares become readable by any process with host-level access.

**Evaluation:** This negates the entire TEE security model that the DKG/CDR system relies on. Sealed storage, attestation, and confidential computation are all void in debug mode.

**Fix:** Set `sgx.debug = false` for production. Use a template variable or separate manifest for dev/prod. This will be addressed separately as part of production manifest hardening.

---

## HIGH Findings (6)

### H-01: Placeholder MRENCLAVE in Deploy Scripts

**Location:** `contracts/script/upgrades/DeploySGXValidationHook.s.sol:35`, `GenerateAlloc.s.sol:63-65`

**Description:** The deploy scripts contain `SGX_CODE_COMMITMENT = hex"0000...0001"` and `AUTOMATA_VALIDATION_ADDR = address(uint160(1000))`, which are development placeholders. `GenerateAlloc.s.sol` generates the actual genesis allocation files used for chain initialization.

**Evaluation:** If these placeholders are not replaced before mainnet deployment, SGX attestation verification will either always fail (wrong code commitment) or point to a non-existent validation contract. This is an operational risk rather than a code bug.

**Fix:** Replace before deployment. Add a CI check that rejects placeholder values in genesis-generating scripts.

---

### H-02: Params Store Dual Registration

**Location:** `client/x/dkg/keeper/params.go:11-25`, `keeper.go:120`

**Description:** DKG params are registered as a `collections.Item[types.Params]` via the `ParamsStore` field in the Keeper, but all actual read/write operations use raw KV store access (`store.Set(types.ParamsKey, bz)` and `store.Get(types.ParamsKey)`). The collections registration is unused but creates a confusing dual-path that could break if the collections library changes its key encoding.

**Evaluation:** The `ParamsStore` field creates dead code and a false assumption that params are managed through the collections framework. If a future developer uses `ParamsStore` directly, the data would be written under a different key prefix, causing silent data inconsistency.

**Fix:** Remove the `ParamsStore` field from the Keeper struct and its initialization in `NewKeeper`. All params operations already use the raw KV store correctly.

---

### H-03: Replay Messages Silently Discards All Errors (Kernel)

**Location:** `service/dist_key_gen.go:586-599`

**Description:** `replayMessages` uses `_, _ =` for all `ProcessDeal`, `ProcessResponse`, and `ProcessJustification` calls during DKG state reconstruction after a kernel restart. If a persisted message is corrupted or the state is inconsistent, all errors are silently swallowed. The kernel then operates with a partially reconstructed DKG that will produce incorrect results.

**Evaluation:** While some "already existing response" errors are expected during replay (because `ProcessDeal` auto-adds the verifier's own response, which then collides with the stored response), critical errors like "different sessionIDs" or "invalid signature" indicate real data corruption and should not be silently ignored.

**Fix:** Log all errors at WARN level with deal/response/justification indices. Return an error from `replayMessages` if a critical (non-duplicate) error occurs, preventing the corrupted DKG instance from being cached and used.

---

### H-04: Upgrade Activation Doesn't Flush Queues

**Location:** `client/x/dkg/keeper/abci.go:55-70`

**Description:** When `InitiateDKGRound` is called (for both regular and upgrade rounds), it does not call `FlushAllQueues()` to clear stale deals, responses, and justifications from the previous round. The `SkipToNextRound` function correctly calls `FlushAllQueues()`, but `InitiateDKGRound` does not. This means stale vote extension data from a previous round can be broadcast into the new round.

**Evaluation:** This affects all DKG round initiations, not just upgrade-triggered ones. Any time a new round starts, the pending incoming queues should be cleared to prevent cross-round data contamination. The stale data would be rejected by session ID checks, but it wastes network bandwidth and creates confusing error logs.

**Fix:** Add `FlushAllQueues()` at the beginning of `InitiateDKGRound`, regardless of whether it is an upgrade round or a regular round.

---

### H-05: Package-Level Global State

**Location:** `client/x/dkg/keeper/keeper.go:21-51`

**Description:** DKG vote extension queues (`deals`, `responses`, `justifications`, and their `pendingIncoming*` counterparts) and the `dkgKernelMu` mutex are declared as package-level global variables rather than fields on the Keeper struct. Each variable has its own mutex, and there are 6 queue variables plus `dkgSvcRound` and `maxPendingIncoming` — a total of 8 package-level globals with 7 associated mutexes.

**Evaluation:** This creates two problems: (1) Test isolation — parallel tests sharing the same process will read/write the same global state, which is why DKG tests explicitly note "NOT parallel because they share package-level globals". (2) If the same binary ever runs multiple Keeper instances (e.g., in a test harness or future multi-chain scenario), the queues would be silently shared, causing data corruption. The pattern is functional for the current single-validator-per-process deployment but fragile.

**Fix:** Move all queue variables and their mutexes into the Keeper struct as fields. Update all queue accessor functions (AddDeals, DrainDeals, FlushAllQueues, etc.) to use receiver fields instead of package globals. This is a large refactor touching ~20 functions.

---

### H-06: cachedLastBlockHeight Race Condition (Kernel)

**Location:** `story-kernel/story/query_client.go:650-661`

**Description:** The `cachedLastBlockHeight` field (int64) is written by a background goroutine (`lastBlockCaching`) and read by query methods (`getQueryBlockHeight`) without synchronization. Under the Go memory model, concurrent reads and writes of non-atomic values constitute a data race, which can produce torn reads (partially updated values).

**Evaluation:** In practice, int64 reads/writes are atomic on 64-bit architectures, but this is not guaranteed by the Go specification and will be flagged by the race detector (`go test -race`).

**Fix:** Use `atomic.Int64` for `cachedLastBlockHeight`, or protect it with the existing `mutex`.

---

## MEDIUM Findings (9)

### M-01: CDR readConditionAddr=address(0) Permanent Vault Lock

**Location:** `contracts/src/protocol/CDR.sol`

**Description:** The `allocate()` function validates that at least one of `writeConditionAddr` or `readConditionAddr` is non-zero: `require(writeConditionAddr != address(0) || readConditionAddr != address(0))`. However, this allows `readConditionAddr` to be `address(0)` as long as `writeConditionAddr` is set. When `readConditionAddr == address(0)`, the `read()` function's access control check `require(msg.sender == vault.readConditionAddr)` requires `msg.sender == address(0)`, which is impossible — permanently locking the vault's read path.

**Evaluation:** A user who sets `writeConditionAddr` but leaves `readConditionAddr` as zero can write to the vault but never read from it. The vault's data becomes permanently inaccessible. However, this requires the user to explicitly pass `address(0)` for `readConditionAddr`, which is a user error. The contract could be more defensive.

**Fix:** Add an explicit check: `require(readConditionAddr != address(0), "read condition address cannot be zero")` in `allocate()`. Alternatively, if a "write-only vault" is a valid use case, document this behavior clearly.

---

### M-02: FinalizeDKGRound Missing GlobalPublicKey Check

**Location:** `client/x/dkg/keeper/dkg_finalization.go`

**Description:** `FinalizeDKGRound` checks that enough validators have finalized (threshold count), distributes rewards, and sets the latest active round. However, it does not verify that `latestRound.GlobalPublicKey` has been set before marking the round as finalized. The GlobalPublicKey is set during the finalization vote processing (when threshold votes are received), so in normal operation it is always set before `FinalizeDKGRound` is called. But there is no explicit guard.

**Evaluation:** If `FinalizeDKGRound` is called before the GlobalPublicKey is set (due to a logic error or race condition), the round would be marked as active without a usable public key. Subsequent CDR encryption operations would fail or produce unusable ciphertexts.

**Fix:** Add a guard at the beginning of `FinalizeDKGRound`: `if len(latestRound.GlobalPublicKey) == 0 { return error }`.

---

### M-03: No Rate Limiting on Kernel Reconnection

**Location:** `client/x/dkg/keeper/kernel_router.go:65`

**Description:** `TryReconnect` makes a single reconnection attempt per disconnected endpoint with a 10-second timeout, but does not implement exponential backoff or per-endpoint cooldown. Since `TryReconnect` is called from `BeginBlock` (every block), a permanently unreachable endpoint will be retried every block (~2 seconds), each attempt allocating a new gRPC connection and waiting up to 10 seconds before timing out.

**Evaluation:** With 3 disconnected endpoints and 10-second timeouts, `BeginBlock` could block for up to 30 seconds per block, significantly impacting block production latency. However, the 10-second context timeout bounds the worst case per call. The main concern is resource waste (connection attempts, goroutines) rather than a correctness issue.

**Fix:** Add per-endpoint cooldown (e.g., exponential backoff starting at 5 seconds, capped at 5 minutes) to avoid retrying unreachable endpoints every block.

---

### M-04: Unbounded In-Memory Caches in Kernel

**Location:** `story-kernel/service/dkg_server.go:30`

**Description:** The kernel maintains four in-memory caches: `InitDKGCache`, `ResharingPrevCache`, `ResharingNextCache`, and `DistKeyShareCache`. None of these caches have eviction policies or maximum size limits. Each cache is keyed by round number (uint32), and entries are `DistKeyGenerator` or `DistKeyShare` structs containing kyber points and scalars.

**Evaluation:** In normal operation, the caches hold at most 2 rounds of data (current + previous for resharing). A `DistKeyGenerator` for a 100-node committee is approximately 50-100KB (100 verifiers × ~500 bytes each). With 4 caches × 2 rounds × 100KB, total memory is under 1MB — well within SGX's 4GB enclave limit. The risk is theoretical: a malicious or buggy consensus layer could request DKG instances for many rounds, growing the cache unboundedly. However, the kernel only creates entries for rounds that have on-chain state, which naturally limits growth.

**Fix:** Add a maximum cache size (e.g., 10 entries per cache) with LRU eviction. This is a low-priority improvement given the natural bound from on-chain state.

---

### M-05: Gramine Manifest Broad File Access

**Location:** `story-kernel/story-kernel.manifest.template:56-69`

**Description:** The Gramine SGX manifest includes `"file:/home/ubuntu/"` in `sgx.allowed_files`, granting the enclave read/write access to the entire home directory. This includes SSH keys, shell history, other application data, and any files the host user has access to. While `allowed_files` means the data passes through without integrity verification (unlike `trusted_files`), it still represents an unnecessarily broad access surface.

**Evaluation:** The kernel only needs access to its own data directory (`~/.story-kernel/` or similar) and configuration files. Granting access to the entire home directory violates the principle of least privilege. A compromised enclave runtime could read sensitive host files. Additionally, `/etc/ssl/` grants access to all SSL certificates and potentially private keys if stored there.

**Fix:** Restrict `allowed_files` to the specific paths the kernel needs:
```
"file:/home/ubuntu/.story-kernel/",
"file:/etc/ssl/certs/",  # only certificates, not keys
```

---

### M-06: GetAllCodeCommitments Non-Deterministic Order

**Location:** `client/x/dkg/keeper/kernel_router.go:143-159`

**Description:** `GetAllCodeCommitments` iterates over a Go map (`r.clients`) to collect code commitments, returning them in non-deterministic order. The caller then selects the first element (`allCCs[0]`) when choosing which kernel client to use for a DKG round.

**Evaluation:** This function is used in the registration flow to select which kernel binary to register with. Since the selection happens off-chain (not in a consensus-critical path), different validators picking different kernel binaries is actually acceptable — each validator registers its own binary's code commitment. However, for upgrade rounds where a specific new binary should be selected, non-deterministic ordering could cause a validator to accidentally register with the old binary.

**Fix:** Sort the code commitments before returning. This ensures deterministic selection across all validators, which is especially important during upgrade rounds.

---

### M-07: DKG State File Not Written Atomically (Kernel)

**Location:** `story-kernel/store/dkg_state.go:88-103`

**Description:** `saveState` writes the DKG state JSON file using `os.WriteFile`, which truncates the file and writes new content in place. If the kernel process crashes mid-write (e.g., SGX enclave terminated, OOM kill, power failure), the file will contain partial JSON data. On the next startup, `loadState` will fail to parse the corrupted file, and `HasDKGState` will return false — causing the kernel to start a fresh DKG instead of recovering from the persisted state.

**Evaluation:** The DKG state file contains deals, responses, and justifications that are critical for restart recovery. Losing this file means the node must restart the DKG from scratch, which may not be possible if the round has progressed beyond the dealing phase. The state file is updated multiple times during a DKG round (once per AddDeals, AddResponses, AddJustifications call), so the window for corruption is not negligible.

**Fix:** Use atomic write pattern: write to a temporary file in the same directory, then `os.Rename` (which is atomic on POSIX filesystems). Same pattern should be applied to `StateManager.saveSession` in the story consensus client (`client/x/dkg/keeper/state_manager.go`).

---

### M-08: DKG.sol finalize() Missing msg.sender Check

**Location:** `contracts/src/protocol/DKG.sol`

**Description:** The `finalize()` function on the DKG smart contract does not validate `msg.sender`. While the consensus layer verifies the TEE attestation signature before calling the contract, the contract itself has no defense-in-depth check. Any externally owned account could call `finalize()` directly with crafted parameters.

**Evaluation:** In practice, the consensus layer's signature verification provides the primary security boundary. However, defense-in-depth is a standard smart contract practice. The risk is that a compromised EVM execution path could bypass the CL checks and call the contract directly.

**Fix:** Add `onlyOwner` or a whitelist check on `msg.sender` in `finalize()`. Alternatively, document why the current approach is acceptable (CL verification is sufficient because the contract is only called via system transactions).

---

### M-09: getDKGRegistrationsByRound Full Store Walk O(n)

**Location:** `client/x/dkg/keeper/dkg_registration.go`

**Description:** `getDKGRegistrationsByRound` iterates over all DKG registrations in the KV store using a prefix iterator to find registrations for a specific round. Each call walks O(n) entries where n is the total number of registrations across all rounds. This is called during registration processing, deal generation, and finalization.

**Evaluation:** For the current scale (tens of validators, few rounds), this is not a performance concern. As the number of historical rounds grows, performance will degrade linearly. However, registrations are naturally bounded by the number of validators per round, and old rounds can be pruned.

**Fix:** Consider indexing registrations by round (secondary index) or pruning old round data during `FinalizeDKGRound`. Low priority given current scale.

---

## LOW Findings (8)

### L-01: `retryAttemts` Typo

**Location:** `client/x/dkg/keeper/dkg_svc.go`

**Description:** The variable `retryAttemts` is a misspelling of "retryAttempts". While functionally harmless, this affects code readability and searchability.

**Fix:** Rename to `retryAttempts`.

---

### L-02: Event StartBlockHeight uint32 Truncation

**Location:** `client/x/dkg/keeper/dkg_initialization.go`

**Description:** The DKG initialization emits the start block height as a uint32 in events, while the actual block height is int64. For block heights exceeding 2^32 (~4.3 billion blocks), the value would be silently truncated. At 2-second block times, this limit would be reached after ~272 years.

**Evaluation:** This affects event data only, not consensus state. Event consumers parsing the truncated height could misinterpret the data, but the on-chain state uses the full int64 value.

**Fix:** Use uint64 or int64 in the event emission. Low priority given the timeline.

---

### L-03: ExportGenesis Only Exports Params

**Location:** `client/x/dkg/module/module.go`

**Description:** The DKG module's `ExportGenesis` function only exports the module params, not the full DKG state (registrations, networks, active rounds). This means a chain export/import cycle would lose all DKG state.

**Evaluation:** For a module with TEE-sealed keys and ongoing DKG rounds, full state export may not be meaningful (sealed keys can't be migrated). The params-only export is consistent with the module's design.

**Fix:** Document that DKG state is not exported by design. If full state export is needed in the future, add registration and network state to the genesis export.

---

### L-04: Deploy Script Salt Reuse Limitation

**Location:** `contracts/script/upgrades/DeploySGXValidationHook.s.sol`

**Description:** The deploy script uses `CREATE2` with a fixed salt for deterministic deployment addresses. If the contract needs to be redeployed (e.g., after a bug fix), the same salt would attempt to deploy to the same address, which would fail because the address is already occupied.

**Fix:** Use a versioned salt (e.g., `keccak256("SGXValidationHook_v2")`) for redeployments. Low priority — redeployment requires a full upgrade flow regardless.

---

### L-05: Session File Permission 0644

**Location:** `client/x/dkg/keeper/state_manager.go`

**Description:** DKG session files are written with permission `0644` (world-readable). While these files don't contain private key material (sealed keys are in the kernel), they do contain deal and response data that reveals which validators participated in which rounds.

**Fix:** Use `0600` (owner-readable only) for all DKG state files.

---

### L-06: StateManager File I/O Not Atomic

**Location:** `client/x/dkg/keeper/state_manager.go:249-262`

**Description:** `saveSession` uses `os.WriteFile` which truncates and rewrites the file in place. If the process crashes during the write, the file will contain partial JSON data. On the next start, `loadSessionFromFile` will fail to parse the corrupted file, returning an error that causes the session to be treated as non-existent.

**Evaluation:** This is the consensus client side (story), not the kernel. The session file stores DKG round state that helps the node recover after a crash. Losing this file is recoverable (the node can rejoin the round from on-chain state), but causes unnecessary round restarts.

**Fix:** Use atomic write pattern: `os.WriteFile(tmpPath, data, 0600)` followed by `os.Rename(tmpPath, path)`.

---

### L-07: retry() time.Sleep Ignores Context

**Location:** `client/x/dkg/keeper/dkg_svc.go`

**Description:** The `retry()` helper function uses `time.Sleep` for backoff delays between retry attempts. If the caller's context is cancelled (e.g., block timeout), the sleep continues for its full duration before checking the context. The maximum delay per attempt is 1 second with 6 attempts, so the worst case is a 6-second uninterruptible sleep.

**Fix:** Replace `time.Sleep(delay)` with a `select` on `time.After(delay)` and `ctx.Done()`:
```go
select {
case <-time.After(delay):
case <-ctx.Done():
    return ctx.Err()
}
```

---

### L-08: HKDF Nil Salt in Partial Decrypt

**Location:** `story-kernel/service/dkg_partial_decrypt.go:300`

**Description:** The HKDF call for deriving the AES encryption key uses `nil` as the salt parameter: `hkdf.New(sha256.New, sharedBytes, nil, []byte("dkg-tdh2-partial"))`. RFC 5869 Section 2.2 specifies that when salt is not provided, it defaults to a string of `HashLen` zeros (32 bytes for SHA-256).

**Evaluation:** This is technically RFC 5869 compliant. The Go `golang.org/x/crypto/hkdf` implementation correctly handles nil salt by using the zero-filled default. The input keying material (`sharedBytes`) comes from an ECDH key exchange, which provides sufficient entropy. Using an explicit non-zero salt would provide additional domain separation but is not required for security.

**Fix:** No action required. The current implementation is correct per RFC 5869. Optionally, pass an explicit salt for clarity: `[]byte("story-kernel-tdh2-v1")`.

---

## Remediation Priority

### P0 — Must Fix Before Production (3 items)

| # | ID | Issue | Effort |
|---|----|-------|--------|
| 1 | C-01 | CDR fee non-deterministic map iteration | Small |
| 2 | C-02 | Registry key collision (raw label) | Small |
| 3 | C-03 | SGX debug=true in manifest | Trivial (ops) |

### P1 — Should Fix Before Production (6 items)

| # | ID | Issue | Effort |
|---|----|-------|--------|
| 4 | H-01 | Placeholder MRENCLAVE in deploy scripts | Trivial (ops) |
| 5 | H-02 | Params dual registration | Small |
| 6 | H-03 | Silent replay error discard (kernel) | Small |
| 7 | H-04 | Upgrade activation doesn't flush queues | Small |
| 8 | H-05 | Package-level global state | Large |
| 9 | H-06 | cachedLastBlockHeight race | Small |

### P2 — Should Fix Shortly After Production

All MEDIUM items (M-01 through M-09).

### P3 — Backlog

All LOW items (L-01 through L-08).

---

## Positive Design Observations

- **Upgrade Handler:** Correct migration ordering, VE activation at upgradeHeight+1, idempotent
- **Vote Extension:** ABCI status returns per CometBFT spec, proper size limits
- **Reward Distribution:** `sort.Strings(memberAddrs)` ensures determinism in DKG rewards
- **CacheContext:** Per-event state rollback in evmengine prevents cross-event contamination
- **Async Context:** Pre-computation of SDK data before goroutine launch
- **dkgKernelMu:** CL-side serialization of kernel mutations
- **Binary-Swap Upgrade:** On-chain ScheduleUpgrade + disk fallback for dual-binary support
- **Module Permissions:** Separate minter for DKG, CDR fee pool, both in blockAccAddrs
- **Contract Security:** UUPS + _disableInitializers, ERC-7201, Ownable2Step, TimelockController
- **Polynomial Persistence:** Dealer polynomial correctly persisted and restored across restarts for both initial and resharing rounds

---

*Generated by 6-agent audit team, revised 2026-03-23*
*Repositories: story@origin/dkg/dev, story-kernel@origin/main*
