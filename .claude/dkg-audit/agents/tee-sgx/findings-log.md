# Findings Log: TEE/SGX Full Audit

## Audit Session: 2026-03-22

### C-03: SGX Debug Mode in Production Manifest
- **Severity:** CRITICAL
- **Status:** TRUE_POSITIVE
- **Outcome:** sgx.debug = true negates entire TEE security model.

### H-03: Replay Messages Silently Discards Errors
- **Severity:** HIGH
- **Status:** TRUE_POSITIVE
- **Outcome:** replayMessages swallows all errors with `_, _ =`. Critical errors invisible.

### H-06: cachedLastBlockHeight Race Condition
- **Severity:** HIGH
- **Status:** TRUE_POSITIVE
- **Outcome:** Use atomic.Int64.

### M-04: Unbounded In-Memory Caches
- **Severity:** MEDIUM
- **Status:** TRUE_POSITIVE
- **Outcome:** Add LRU eviction with max size.

### M-05: Gramine Manifest Broad File Access
- **Severity:** MEDIUM
- **Status:** TRUE_POSITIVE
- **Outcome:** Restrict allowed_files to ~/.story-kernel/ only.

## MISSED by this agent (found by APEX):
- STOR-5: Unauthenticated PartialDecryptTDH2 = remote decryption oracle
- STOR-17: Sealed light-client state rollbackable by hostile host
- STOR-18: VerifiedQueryClient accepts proofs for arbitrary keys/stores
