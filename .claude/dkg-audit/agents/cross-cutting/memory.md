# Agent Memory: Cross-Cutting Security & Workflow Auditor

## Codebase Knowledge
- Three-layer architecture: EVM contracts → CL (Cosmos SDK) → TEE kernel (SGX)
- Trust flow: Contracts emit events → CL processes deterministically → Kernel does crypto
- DKG state machine: registration → dealing → response → finalization → active → resharing
- CDR flow: allocate → write → read (triggers decrypt request) → partial decrypt → combine

## Past False Positives
- SECURITY-007, SECURITY-011: flagged issues that were design-documented

## Blind Spots Discovered
- **CDR availability mapping (STOR-9):** Never systematically mapped which CDR operations work in which DKG states. CDR read is dead during any non-active phase. This is a fundamental availability gap I should have caught
- **Resharing → active transition gaps (STOR-8, STOR-21):** Resharing doesn't populate PIDCache, doesn't carry over serving round pointer correctly. State transition analysis was too shallow
- **Adversarial validator ordering (STOR-7):** Never modeled what happens when a malicious validator sorts BEFORE honest ones in ExtendedCommitInfo. The first-seen dedup becomes a censorship tool
- **DKG stage boundary off-by-one (STOR-19, STOR-22, STOR-23):** Events from the last block of one stage are rejected in the first block of the next stage. I should have checked every transition boundary for off-by-one handling
- **Cross-layer trust assumptions:** I assumed CL validates event data from contracts. For CDR submits, it doesn't — it trusts the event fields and rewards based on nil-return

## Effective Techniques
- Analyzing upgrade workflow corner cases
- Identifying trust boundaries between components

## Ineffective Techniques
- High-level system analysis without concrete attack scenario construction
- Not mapping state machine transitions systematically against functionality availability
- Not modeling adversarial control of input ordering
