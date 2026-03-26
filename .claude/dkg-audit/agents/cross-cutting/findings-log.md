# Findings Log: Cross-Cutting Security & Workflow Auditor

## Audit Session: 2026-03-22

### C-02: Registry Key Collision
- **Severity:** CRITICAL
- **Status:** TRUE_POSITIVE
- **Outcome:** Raw label bytes in key string. Use hex encoding.

### M-06: GetAllCodeCommitments Non-Deterministic Order
- **Severity:** MEDIUM
- **Status:** TRUE_POSITIVE
- **Outcome:** Sort code commitments for deterministic selection.

## MISSED by this agent (found by APEX):
- STOR-7: Vote dedup censorship by adversarial ordering
- STOR-9: CDR read dead during all non-active DKG phases
- STOR-8: PIDCache empty after resharing
- STOR-21: Resharing mutates serving-round pointer
- STOR-19, STOR-22, STOR-23: Stage transition boundary bugs
- STOR-2: Upgrade resharing rejects non-committee validators
