# DKG/CDR Audit Context

## First Audit — 2026-03-22

### Audit Request
Production-level final audit before live network deployment. Covers security, bugs, code quality, and corner cases. This is the last gate before production.

### Target Branches
- **story:** `dkg/refactor-v200-upgrade`
- **story-kernel:** `hans/fix-dealer-poly-persistence`

### Design Documents (ALWAYS reference in audits)
- **DKG:** `story/docs/design/DKG.md` — system architecture, DKG lifecycle, reward distribution, kernel upgrade workflow, key technical details, story-kernel TEE architecture
- **CDR:** `story/docs/design/CDR.md` — allocate, encrypt/write, decryption request, partial decryption, threshold combination, decrypt request lifecycle

### Audit Scope
1. v2.0.0 upgrade handler (upgrading from v1.5.3 live network)
2. DKG module (Cosmos SDK module)
3. DKG service (communication with story-kernel)
4. DKG rewards and fee distribution
5. DKG/CDR related contracts and deploy scripts
6. story-kernel full codebase (first production release)
7. End-to-end DKG/CDR workflow corner cases

### Key Context
- Live network is running v1.5.3, upgrade to v2.0.0 adds DKG module
- story-kernel is going to production for the first time
- Unit test coverage is low (being worked on in parallel in separate sessions)
- DKG uses Pedersen DKG (kyber v4) with Edwards25519; TDH2 via cb-mpc for threshold decryption
- story-kernel runs inside Intel SGX enclave via Gramine

### Results (Post Design-Doc Re-validation)
- **Report:** `story/.claude/FINAL_AUDIT.md`
- **Initial findings:** 10 CRITICAL, 16 HIGH, 19 MEDIUM, 12 LOW
- **After re-validation:** 5 CRITICAL, 8 HIGH, 17 MEDIUM, 12 LOW (14 FALSE POSITIVES removed)
- **Key change:** CRYPTO-001 (sign/verify field mismatch) upgraded to CRITICAL — CDR completely broken

### P0 Blockers (5)
1. C-01: CDR fee non-deterministic map iteration → consensus failure
2. C-02: SGX debug=true in production manifest → TEE security void
3. C-03: Partial decrypt sign/verify field mismatch → CDR completely broken
4. C-04: Decrypt worker context leak → CDR stops after 1 min
5. C-05: Registry key collision (raw label bytes)

### Remediation Checklist
**CRITICAL:**
- [ ] C-01: CDR fee `distributeCDRRewardPool` — sort map keys before iteration
- [ ] C-02: SGX manifest `debug = false` for production
- [ ] C-03: Align sign/verify fields (kernel signs CC||round||..., CL verifies round||ciphertext||...)
- [ ] C-04: Decrypt worker — pass long-lived context instead of dkgAsyncContext
- [ ] C-05: Registry key — use `hex.EncodeToString(label)`

**HIGH:**
- [ ] H-01: Replace placeholder MRENCLAVE before deployment
- [ ] H-02: Params — unify raw KV vs collections.Item
- [ ] H-03: Kernel replay — log errors, fail on critical message errors
- [ ] H-04: cachedLastBlockHeight — use atomic.Int64
- [ ] H-05: Upgrade activation — flush queues before new round
- [ ] H-06: CDR partial decrypt on-chain verification (or slashing)
- [ ] H-07: Move global state to Keeper struct
- [ ] H-08: Kernel per-round mutex for all DKG RPCs

### False Positives Identified (14)
COSMOS-006, COSMOS-007, COSMOS-008, COSMOS-010, DKG-SVC-002 (design-documented), DKG-SVC-012 (correct dedup key), KERNEL-002 (no private keys), KERNEL-009 (RFC compliant), KERNEL-010 (no reflection), KERNEL-016 (nonexistent function), CRYPTO-004, CRYPTO-007, SECURITY-007, SECURITY-011

---

## External Audit Comparison — 2026-03-24

### External Reports
- **APEX Report:** `raul/_extracted/Apex_Report_-_Story_network_Scan_#1_...md` — 5 HIGH, 20 MEDIUM, 4 LOW
- **CDR AI Security Audit:** `raul/_extracted/CDR_-_AI_Security_Audit_Report_...md` — 1 P0, 3 P1, 10 P2, 3 P3, 1 P4, 2 P5

### Critical Findings We Missed (APEX-exclusive)

| ID | Issue | Impact | Root Cause of Gap |
|---|---|---|---|
| STOR-15 | CDR fee wei→gwei 1e9 amplification | Unbacked stake minting, bridge drain | No agent traced EL↔CL unit conversion |
| STOR-6 | Arbitrary/replayed partials rewarded | Fee pool theft | No agent checked non-validator callers |
| STOR-5 | Kernel = unauthenticated decrypt oracle | Vault confidentiality breach | TEE agent didn't model network attacker |
| STOR-4 | nil-return on not-found = rewarded | Fee pool drain | Error-path economics not analyzed |
| STOR-3 | Duplicate partial replays accrue rewards | Reward inflation | Replay analysis incomplete |
| STOR-7 | Vote dedup censors honest DKG messages | DKG transcript corruption | Adversarial ordering not modeled |
| STOR-8 | PIDCache empty after resharing | CDR permanently broken | State lifecycle not traced across transitions |
| STOR-13 | pubKeyShare format mismatch | CDR reads fail | Cross-layer serialization not compared |

### Lessons Learned → Applied to audit-team.md v2
1. Added **Economic & Value Flow Auditor** (Agent #7)
2. Added **Adversarial Integration Tester** (Agent #8)
3. Added **Phase 2: Cross-Layer Trace** between audit and re-validation
4. Enhanced every agent with APEX-derived mandatory checks
5. Upgraded finding format to APEX-quality standard (executive summary, code trace, impact cascade, anticipated response, test gap analysis)
6. Added mandatory adversarial scenarios for Agent #8

### Response Document
- **File:** `raul/responses.md` — Our response to each external finding with fix status
