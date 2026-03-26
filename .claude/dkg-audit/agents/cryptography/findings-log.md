# Findings Log: Cryptography Auditor

## Audit Session: 2026-03-22

### C-03 (CRYPTO-001): Sign/Verify Field Mismatch
- **Severity:** CRITICAL (upgraded from HIGH during re-validation)
- **Status:** TRUE_POSITIVE
- **Outcome:** CDR completely broken. Kernel signs CC||round||..., CL verifies round||ciphertext||...

### L-08: HKDF Nil Salt
- **Severity:** LOW
- **Status:** TRUE_POSITIVE (but not a security issue — RFC compliant)
- **Outcome:** Optionally add explicit salt for clarity.

## MISSED by this agent (found by APEX):
- STOR-1: Partial decrypt signature omits requesterPubKey and label binding
- STOR-13: pubKeyShare raw vs prefixed format mismatch across layers
- CDR-004: EIP-191 prefix inconsistency between finalization and partial decrypt paths
- CDR-007: Signature missing codeCommitment field
