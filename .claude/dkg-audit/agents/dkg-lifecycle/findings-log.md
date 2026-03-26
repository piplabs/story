# Findings Log: DKG Service & Lifecycle Auditor

## Audit Session: 2026-03-22

### C-04: Decrypt Worker Context Leak
- **Severity:** CRITICAL
- **Status:** TRUE_POSITIVE
- **Outcome:** dkgAsyncContext expires after 1 min, killing decrypt worker. Need long-lived context.

### M-03: No Rate Limiting on Kernel Reconnection
- **Severity:** MEDIUM
- **Status:** TRUE_POSITIVE
- **Outcome:** Added exponential backoff for kernel reconnection.

## MISSED by this agent (found by APEX):
- STOR-4: nil-return on not-found = rewarded (error-path economics)
- STOR-3: Duplicate partial replays accrue rewards
- STOR-8: PIDCache empty after resharing breaks CDR
- STOR-9: CDR read dead during all non-active DKG phases
- STOR-28: Decrypt worker dies after 1 min while CDR charges 21-day fees
