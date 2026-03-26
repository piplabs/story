# Findings Log: Smart Contract & Deploy Script Auditor

## Audit Session: 2026-03-22

### H-01: Placeholder MRENCLAVE in Deploy Scripts
- **Severity:** HIGH
- **Status:** TRUE_POSITIVE (operational)
- **Outcome:** Must replace before mainnet.

### M-01: CDR readConditionAddr=address(0) Permanent Lock
- **Severity:** MEDIUM
- **Status:** TRUE_POSITIVE
- **Outcome:** User can make vault permanently unreadable.

### M-08: DKG.sol finalize() Missing msg.sender Check
- **Severity:** MEDIUM
- **Status:** TRUE_POSITIVE
- **Outcome:** Defense-in-depth gap.

## MISSED by this agent (found by APEX):
- STOR-15: CDR fee wei/gwei 1e9 amplification (cross-layer value flow)
- STOR-4: Permissionless submitEncryptedPartialDecryption fee drain
- STOR-6: Arbitrary partial submissions still rewarded
- CDR-005: finalize() no sender validation (found by CDR report, not us)
- CDR-006: Condition bypass when msg.sender is condition contract
