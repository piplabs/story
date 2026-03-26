# Findings Log: Cosmos SDK & Upgrade Handler Auditor

## Audit Session: 2026-03-22 — story@dkg/refactor-v200-upgrade, story-kernel@hans/fix-dealer-poly-persistence

### C-01: CDR Fee Non-Deterministic Map Iteration
- **Severity:** CRITICAL
- **Status:** TRUE_POSITIVE — Fixed in `origin/dkg/dev` via `71ded0d1`
- **Outcome:** Real consensus failure bug. Added sort.Strings() before iteration.

### H-02: Params Store Dual Registration
- **Severity:** HIGH
- **Status:** TRUE_POSITIVE
- **Outcome:** ParamsStore field unused, creates confusion. Should remove.

### H-04: Upgrade Activation Doesn't Flush Queues
- **Severity:** HIGH
- **Status:** TRUE_POSITIVE
- **Outcome:** InitiateDKGRound needs FlushAllQueues() at start.

### H-05: Package-Level Global State
- **Severity:** HIGH
- **Status:** TRUE_POSITIVE
- **Outcome:** Large refactor needed — move queue globals to Keeper struct.

### COSMOS-006, COSMOS-007, COSMOS-008, COSMOS-010
- **Severity:** various
- **Status:** FALSE_POSITIVE — Design-documented behaviors
- **Outcome:** Need to read design docs MORE carefully before flagging.

## MISSED by this agent (found by APEX):
- STOR-7: Vote dedup censorship (adversarial ordering)
- STOR-19: First dealing block rejects last registration block's events
- STOR-22: First active block rejects last finalization block's events
- STOR-23: Dealing-to-finalization drops last block's complaints
