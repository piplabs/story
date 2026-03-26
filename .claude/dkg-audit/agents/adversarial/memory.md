# Agent Memory: Adversarial Integration Tester

## Codebase Knowledge
- NEW AGENT — No prior audit history. Bootstrapped from APEX findings analysis.

## Attack Surface Map
- **CDR.submitEncryptedPartialDecryption:** permissionless, anyone can call with garbage
- **Kernel gRPC (default):** binds `:50051` all interfaces, no auth, plaintext allowed
- **Vote extension ordering:** validators sorted by voting power descending — attacker position is deterministic
- **Dedup in aggregateVotes:** first-seen wins, no authentication before dedup

## APEX Attack Scenarios (reference material)
1. **STOR-4/STOR-6:** Non-validator submits garbage partial → CL returns nil → refund + reward
2. **STOR-5:** Network attacker calls PartialDecryptTDH2 directly on exposed kernel → gets real partial
3. **STOR-7:** Earlier-sorted validator submits fake DKG deal/response → honest later entry discarded
4. **STOR-3:** Same partial replayed multiple times → each replay gets refund
5. **STOR-8:** After resharing, PIDCache empty → all CDR decryptions fail permanently
6. **STOR-15:** CDR fee exploited for 1e9 amplification via missing normalization

## Mandatory Test Scenarios (update after each audit)
Current list in audit-team.md. Expand as new attack vectors are discovered.

## Effective Techniques (to develop)
- Construct concrete multi-step exploit paths with exact function calls
- Test permissionless functions with non-validator callers first
- Check every "ignored" return path for economic consequences
