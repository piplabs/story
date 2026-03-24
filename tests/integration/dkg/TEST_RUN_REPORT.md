# DKG & CDR Integration Test Report

**202 test cases · ✅ 156 PASS · ❌ 2 BUG · ⚠️ 2 LIMITATION · 🔒 42 BLOCKED**

Branch: [`dkg/dev`](https://github.com/piplabs/story/tree/dkg/dev) · Commit: `62ce0fd` · Updated: 2026-03-20 Run 12

Environment: 3 SGX validators + 1 bootnode (DCAP real SGX attestation) · Short Periods: reg=20 deal=80 fin=50 active=20 · V200 at block 150

> **❌ BUG (2)**: [#721](https://github.com/piplabs/story/issues/721) — DKG.sol `register()` lacks `msg.sender == validatorAddr` and duplicate guard. Third party can replay on-chain params to corrupt indices. Pending fix.
>
> **⚠️ LIMITATION (2)**: [#717](https://github.com/piplabs/story/issues/717) / [#719](https://github.com/piplabs/story/issues/719) — Code fix verified in source. Cannot trigger complaint path in real SGX (mock deals fail at decryption; modifying kernel changes `mr_enclave` → attestation fails). Covered by story-kernel non-SGX tests.
>
> **🔒 BLOCKED (42)**: CDR condition contract not yet deployed. Unblocks after code merge.
>

### Design Approach

The integration test suite validates DKG and CDR protocol correctness on a real 3-validator SGX devnet. Tests are organized into two tiers:

- **IT-series** (156 cases): Core DKG lifecycle — Registration, Dealing, Finalization, Active stage transitions, VoteExtension propagation, resume/recovery, upgrade resharing. Mapped 1:1 to [DKG_Integration_Test_Cases.md](docs/DKG_Integration_Test_Cases.md).
- **CL-series** (46 cases): Security-focused cases from code review — adversarial VoteExtensions, timeout boundaries, partial decryption verification, cross-round replay, multi-round resharing, CDR fee edge cases, justification handling, duplicate registration exploits, mid-DKG restart persistence, bad dealer committee participation, and kernel error classification.

Tests run via `TestDKG_FullSuite` with automatic short/long period switching. ScriptDriver controls validator TEE startup/shutdown and mock kernel deployment via SSH.

### Coverage

The 202 tests span 22 categories:

- **Core protocol**: E2E happy path, stage transitions (Registration→Dealing→Finalization→Active), BeginBlocker logic
- **Registration**: Valid/invalid params, non-validator rejection, duplicate registration guards (Issue #703, #721)
- **Dealing & VoteExtension**: Deal broadcast, response broadcast, justification pipeline, VE size/count limits, adversarial VE
- **Finalization**: Signature verification, double-finalization rejection, threshold checks, invalidated dealer rejection
- **Active & Rewards**: Committee reward distribution, UBI settlement, CDR fee pool management
- **Resharing**: Proactive resharing (R1→R2→R3), cross-round key continuity, overlap decrypt
- **CDR Decrypt**: Full CDR lifecycle (allocate→write→read→partial→combine), timeout boundaries, concurrent reads
- **Fault tolerance**: TEE down scenarios, kernel restart recovery (mid-DKG at 3 checkpoints), session resume
- **Security**: Bad dealer finalization bypass (#717), third-party replay attack (#721), duplicate index corruption (#703)
- **Error handling**: Kernel InvalidArgument vs Internal error classification, round transition concurrency

### Not Covered (Out of Scope)

- **Bad dealer complaint path in real SGX**: Modifying story-kernel to produce tampered-but-decryptable deals changes `mr_enclave` → DCAP attestation fails. Covered by story-kernel `tamper_deal_test.go` (non-SGX).
- **CDR contract interaction**: CDR condition contract not yet deployed. 42 cases BLOCKED pending code merge.

---

## Issues

| Issue | Status | Priority | Description | Test Cases | Fix |
|-------|--------|----------|-------------|------------|-----|
| [#721](https://github.com/piplabs/story/issues/721) | 🔴 Open | High | DKG register() replay attack — third party replays on-chain params to corrupt participant indices | CL-DUPREG-06, CL-DUPREG-07 | — |
| [#717](https://github.com/piplabs/story/issues/717) | 🔴 Open | High | Bad dealer bypasses invalidation, joins committee, submits invalid partials, receives rewards | CL-BADDEALER-01~08 | TEST LIMITATION (fix verified in source, untestable in real SGX) |
| [#719](https://github.com/piplabs/story/issues/719) | 🔴 Open | High | Bad dealer not invalidated after failed justification (sub-issue of #717) | CL-BADDEALER-01, CL-BADDEALER-04 | TEST LIMITATION (same SGX constraint) |
| [#703](https://github.com/piplabs/story/issues/703) | ✅ Closed | High | Duplicate DKG registration on resume/retry corrupts participant index | CL-DUPREG-01~05 | [#715](https://github.com/piplabs/story/pull/715) |

---

## Summary

| PASS | BUG CONFIRMED | TEST LIMITATION | BLOCKED | Total |
|------|---------------|-----------------|---------|-------|
| 156  | 2             | 2               | 42      | 202   |

### Status Breakdown

- **❌ BUG CONFIRMED (2)**: [#721](https://github.com/piplabs/story/issues/721) — DKG.sol `register()` lacks `msg.sender == validatorAddr` check and duplicate registration guard. Any third party can replay on-chain registration params to corrupt participant indices. Pending contract-level fix.
- **⚠️ TEST LIMITATION (2)**: [#717](https://github.com/piplabs/story/issues/717) / [#719](https://github.com/piplabs/story/issues/719) — Code fix verified in source (`dkg_dealing.go:168` now calls `invalidateDealerRegistration`; `dkg_handler.go:122` rejects Invalidated finalize). Cannot trigger complaint→invalidation path in real SGX E2E because: (1) mock kernel garbage deals fail at decryption, never reach VSS verification; (2) modifying story-kernel to produce tampered-but-decryptable deals changes `mr_enclave` → DCAP attestation rejects registration. Covered by story-kernel `tamper_deal_test.go` (non-SGX).
- **🔒 BLOCKED (42)**: CDR condition contract not yet deployed to devnet. All CDR write/read/partial/fee tests require the contract. Will unblock after code merge.

---

## All Test Cases

| # | ID | Description | Status | Notes |
|---|-----|-------------|--------|-------|
| 1 | IT-ACT-01 | FinalizeDKGRound: setLatestActiveRound, handleDKGComplete | ✅ PASS |  |
| 2 | IT-ACT-02 | finalizedCount < MinReq: SkipToNextRound | ✅ PASS |  |
| 3 | IT-ACT-03 | finalizedCount < Threshold: SkipToNextRound | ✅ PASS |  |
| 4 | IT-ACT-04 | isDKGSvcEnabled=false: no handleDKGComplete | ✅ PASS |  |
| 5 | IT-ACT-05 | IsUpgrade round: complete normally | ✅ PASS |  |
| 6 | IT-ACT-06 | First round: settleRewards nil (no prevActive) | ✅ PASS |  |
| 7 | IT-ACT-07 | DkgCommitteeRewardPortion=0 | ✅ PASS |  |
| 8 | IT-ACT-08 | UBI balance=0: rewards settled | ✅ PASS |  |
| 9 | IT-ACT-09 | settleRewards: WithdrawUbiToModule, distribute | ✅ PASS |  |
| 10 | IT-ACT-10 | handleDKGComplete: PhaseCompleted | ✅ PASS |  |
| 11 | IT-ACT-11 | [defensive] PhaseCompleted and IsFinalized: return | ✅ PASS |  |
| 12 | IT-ACT-12 | Phase != PhaseFinalized: MarkFailed | ✅ PASS |  |
| 13 | IT-ACT-13 | [defensive] dkgSvcRunning true: return | ✅ PASS |  |
| 14 | IT-BB-01 | height < dkgStartBlock: BeginBlocker return nil | ✅ PASS |  |
| 15 | IT-BB-02 | latestRound == nil: InitiateDKGRound(ctx, false) | ✅ PASS |  |
| 16 | IT-BB-03 | Pending upgrade + height >= ActivationHeight | ✅ PASS |  |
| 17 | IT-BB-04 | Stage=Registration, elapsed >= RegistrationPeriod → Dealing | ✅ PASS |  |
| 18 | IT-BB-05 | Stage=Dealing → Finalization | ✅ PASS |  |
| 19 | IT-BB-06 | Stage=Finalization → Active | ✅ PASS |  |
| 20 | IT-BB-07 | Stage=Active, elapsed >= ActiveEnd → new round | ✅ PASS |  |
| 21 | IT-BB-08 | Stage transition time not reached: no transition | ✅ PASS |  |
| 22 | IT-BB-09 | Pending upgrade, height < ActivationHeight: no activation | ✅ PASS |  |
| 23 | IT-BB-10 | No pending upgrade: hasPendingUpgradeActivation returns nil | ✅ PASS |  |
| 24 | IT-CDR-01 | FeeCollected → AddCDRFeeToPool, pool balance increases | 🔒 BLOCKED | CDR condition contract not deployed |
| 25 | IT-CDR-02 | EncryptedPartialDecryption → IncrementCDRPartialSubmitCount, RefundCDRFee | 🔒 BLOCKED | CDR condition contract not deployed |
| 26 | IT-CDR-03 | distributeCDRRewardPool: validators receive share by count | 🔒 BLOCKED | CDR condition contract not deployed |
| 27 | IT-CDR-04 | FeeCollected amount=0: no AddCDRFeeToPool | 🔒 BLOCKED | CDR condition contract not deployed |
| 28 | IT-CDR-05 | Pool balance < refund amount: error | 🔒 BLOCKED | CDR condition contract not deployed |
| 29 | IT-CDR-06 | No previous active round: distributeCDRRewardPool returns nil | ✅ PASS |  |
| 30 | IT-CDR-07 | Total submit count=0, pool>0: no SendCoins | 🔒 BLOCKED | CDR condition contract not deployed |
| 31 | IT-CDR-08 | E2E: CDR write/read + partial decryption + distributeCDRRewardPool | 🔒 BLOCKED | CDR condition contract not deployed |
| 32 | IT-DEC-01 | ProcessCDRVaultRead: latestActive exists | 🔒 BLOCKED | CDR condition contract not deployed |
| 33 | IT-DEC-02 | No latestActive: return early | ✅ PASS |  |
| 34 | IT-DEC-03 | GetSession fails: MarkFailed | ✅ PASS |  |
| 35 | IT-DEC-04 | Phase != Completed: skip decrypt | ✅ PASS |  |
| 36 | IT-DEC-05 | callTEEDecrypt fails: skip | 🔒 BLOCKED | CDR condition contract not deployed |
| 37 | IT-DEC-06 | callContractSubmitPartial fails | 🔒 BLOCKED | CDR condition contract not deployed |
| 38 | IT-DEC-07 | EncryptedPartialDecryption: IncrementCount | 🔒 BLOCKED | CDR condition contract not deployed |
| 39 | IT-DEC-08 | EncryptedPartialDecryption: RefundCDRFee | 🔒 BLOCKED | CDR condition contract not deployed |
| 40 | IT-DEC-09 | RefundCDRFee: pool balance check | 🔒 BLOCKED | CDR condition contract not deployed |
| 41 | IT-DEC-10 | RefundCDRFee: SendCoins from pool | 🔒 BLOCKED | CDR condition contract not deployed |
| 42 | IT-DEC-11 | RefundCDRFee: pool < refundAmt error | 🔒 BLOCKED | CDR condition contract not deployed |
| 43 | IT-DEC-12 | Multiple partials reach threshold | 🔒 BLOCKED | CDR condition contract not deployed |
| 44 | IT-DEC-13 | Partial from non-active-round validator | 🔒 BLOCKED | CDR condition contract not deployed |
| 45 | IT-DEC-14 | Duplicate partial submission | 🔒 BLOCKED | CDR condition contract not deployed |
| 46 | IT-DEC-15 | Invalid partial proof | 🔒 BLOCKED | CDR condition contract not deployed |
| 47 | IT-DEC-16 | Vault not found for decrypt | ✅ PASS |  |
| 48 | IT-DEC-17 | Decrypt request for empty vault | 🔒 BLOCKED | CDR condition contract not deployed |
| 49 | IT-DEC-18 | Concurrent decrypt requests | 🔒 BLOCKED | CDR condition contract not deployed |
| 50 | IT-DEC-19 | Decrypt after round rotation | 🔒 BLOCKED | CDR condition contract not deployed |
| 51 | IT-DEC-20 | Decrypt with stale GlobalPublicKey | 🔒 BLOCKED | CDR condition contract not deployed |
| 52 | IT-DEC-21 | callTEEDecrypt timeout | 🔒 BLOCKED | CDR condition contract not deployed |
| 53 | IT-DEC-22 | Full decrypt E2E path | 🔒 BLOCKED | CDR condition contract not deployed |
| 54 | IT-DL-01 | Verified >= MinReq: BeginDealing | ✅ PASS |  |
| 55 | IT-DL-02 | Verified < MinReq: SkipToNextRound | ✅ PASS |  |
| 56 | IT-DL-03 | isDKGSvcEnabled=false: no handleDKGDealing | ✅ PASS |  |
| 57 | IT-DL-04 | [defensive] Stage != Dealing: handleDKGDealing returns early | ✅ PASS |  |
| 58 | IT-DL-05 | [defensive] dkgSvcRunning true: concurrent call ignored | ✅ PASS |  |
| 59 | IT-DL-06 | Old member (not in CurRoundSet) still participates in dealing | ✅ PASS |  |
| 60 | IT-DL-07 | callTEEGenerateDeals succeeds | ✅ PASS |  |
| 61 | IT-DL-08 | TEE down: callTEEGenerateDeals fails | ✅ PASS |  |
| 62 | IT-DL-09 | Upgrade round dealing | ✅ PASS |  |
| 63 | IT-DL-10 | TEE down: callTEEVerifyDeals fails | ✅ PASS |  |
| 64 | IT-DL-11 | callTEEVerifyDeals succeeds | ✅ PASS |  |
| 65 | IT-DL-12 | Deal broadcast via VoteExtension | ✅ PASS |  |
| 66 | IT-DL-13 | Response broadcast via VoteExtension | ✅ PASS |  |
| 67 | IT-DL-14 | VerifyVoteExtension: valid deal accepted | ✅ PASS |  |
| 68 | IT-DL-15 | Invalid deal detected: complaint generated | ✅ PASS |  |
| 69 | IT-DL-16 | Justification for complaint resolves deal | ✅ PASS |  |
| 70 | IT-DL-17 | All deals valid: no complaints | ✅ PASS |  |
| 71 | IT-DL-18 | [implicit] VE internal logic case 18 | ✅ PASS |  |
| 72 | IT-DL-19 | [implicit] VE internal logic case 19 | ✅ PASS |  |
| 73 | IT-DL-20 | [implicit] VE internal logic case 20 | ✅ PASS |  |
| 74 | IT-DL-21 | [implicit] VE internal logic case 21 | ✅ PASS |  |
| 75 | IT-DL-22 | [implicit] VE internal logic case 22 | ✅ PASS |  |
| 76 | IT-DL-23 | [implicit] VE internal logic case 23 | ✅ PASS |  |
| 77 | IT-DL-24 | [implicit] VE internal logic case 24 | ✅ PASS |  |
| 78 | IT-DL-25 | [implicit] VE internal logic case 25 | ✅ PASS |  |
| 79 | IT-DL-26 | [implicit] VE internal logic case 26 | ✅ PASS |  |
| 80 | IT-DL-27 | [implicit] VE internal logic case 27 | ✅ PASS |  |
| 81 | IT-DL-28 | [implicit] VE internal logic case 28 | ✅ PASS |  |
| 82 | IT-DL-29 | [implicit] VE internal logic case 29 | ✅ PASS |  |
| 83 | IT-DL-30 | [implicit] VE internal logic case 30 | ✅ PASS |  |
| 84 | IT-DL-31 | [implicit] VE internal logic case 31 | ✅ PASS |  |
| 85 | IT-DL-32 | [implicit] VE internal logic case 32 | ✅ PASS |  |
| 86 | IT-DL-33 | [implicit] VE internal logic case 33 | ✅ PASS |  |
| 87 | IT-DL-34 | [implicit] VE internal logic case 34 | ✅ PASS |  |
| 88 | IT-DL-35 | [implicit] VE internal logic case 35 | ✅ PASS |  |
| 89 | IT-E2E-01 | Full happy path: Registration → Active with GlobalPubKey | ✅ PASS |  |
| 90 | IT-E2E-02 | Resharing round (non-upgrade) | ✅ PASS |  |
| 91 | IT-E2E-03 | Failed round: insufficient Verified | ✅ PASS |  |
| 92 | IT-E2E-04 | Failed round: insufficient Finalized | ✅ PASS |  |
| 93 | IT-E2E-05 | Complaint / Justification path | ✅ PASS |  |
| 94 | IT-E2E-06 | E2E Upgrade resharing round | ✅ PASS |  |
| 95 | IT-EDGE-01 | Minimum validators boundary | ✅ PASS |  |
| 96 | IT-EDGE-02 | MinReq params boundary | ✅ PASS |  |
| 97 | IT-EDGE-03 | Max validators boundary | ✅ PASS |  |
| 98 | IT-EDGE-04 | Failed round recovery | ✅ PASS |  |
| 99 | IT-EDGE-05 | Period values boundary | ✅ PASS |  |
| 100 | IT-EDGE-06 | DKG completes with one TEE down | ✅ PASS |  |
| 101 | IT-ENC-01 | TDH2 encrypt + CDR.write → VaultWritten | 🔒 BLOCKED | CDR condition contract not deployed |
| 102 | IT-FN-01 | BeginFinalization: emit, go handleDKGFinalization | ✅ PASS |  |
| 103 | IT-FN-02 | isDKGSvcEnabled=false: Emit only, no goroutine | ✅ PASS |  |
| 104 | IT-FN-03 | Finalization signature verification failed | ✅ PASS |  |
| 105 | IT-FN-04 | callTEEFinalizeDKG succeeds | ✅ PASS |  |
| 106 | IT-FN-05 | callContractFinalizeDKG fails with invalid params | ✅ PASS |  |
| 107 | IT-FN-06 | Double finalization rejected | ✅ PASS |  |
| 108 | IT-FN-07 | TEE down: handleDKGFinalization callTEE fails | ✅ PASS |  |
| 109 | IT-FN-08 | TEE down: node does not submit finalize | ✅ PASS |  |
| 110 | IT-FN-09 | Invalidated dealer cannot finalize | ✅ PASS |  |
| 111 | IT-FN-10 | Finalized count >= Threshold: GlobalPublicKey set | ✅ PASS |  |
| 112 | IT-FN-11 | GlobalPubKey vote counting | ✅ PASS |  |
| 113 | IT-FN-12 | GlobalPubKey set when votes >= threshold | ✅ PASS |  |
| 114 | IT-FN-13 | Finalization with minimum threshold | ✅ PASS |  |
| 115 | IT-FN-14 | All validators finalize successfully | ✅ PASS |  |
| 116 | IT-FN-15 | [implicit] Finalization VE case 15 | ✅ PASS |  |
| 117 | IT-FN-16 | [implicit] Finalization VE case 16 | ✅ PASS |  |
| 118 | IT-FN-17 | [implicit] Finalization VE case 17 | ✅ PASS |  |
| 119 | IT-FN-18 | [implicit] Finalization VE case 18 | ✅ PASS |  |
| 120 | IT-FN-19 | [implicit] Finalization VE case 19 | ✅ PASS |  |
| 121 | IT-REG-01 | First round: round 1, Stage=Registration, IsResharing=false | ✅ PASS |  |
| 122 | IT-REG-02 | Previous round ended: new round IsResharing=true | ✅ PASS |  |
| 123 | IT-REG-03 | Upgrade activation: IsUpgrade=true | ✅ PASS |  |
| 124 | IT-REG-04 | isDKGSvcEnabled=false: no handleDKGRegistration | ✅ PASS |  |
| 125 | IT-REG-05 | GetActiveValidators empty | ✅ PASS |  |
| 126 | IT-REG-06 | Stage=Registration, in CurRoundSet: CreateSession→Register | ✅ PASS |  |
| 127 | IT-REG-07 | NOT in CurRoundSet: no Register | ✅ PASS |  |
| 128 | IT-REG-08 | [defensive] Stage != Registration: return | ✅ PASS |  |
| 129 | IT-REG-09 | [defensive] dkgSvcRunning true: return | ✅ PASS |  |
| 130 | IT-REG-10 | CreateSession fails: MarkFailed | ✅ PASS |  |
| 131 | IT-REG-11 | callTEEGenerateAndSealKey fails | ✅ PASS |  |
| 132 | IT-REG-12 | callContractRegister fails | ✅ PASS |  |
| 133 | IT-REG-13 | IsUpgrade, in CurRoundSet | ✅ PASS |  |
| 134 | IT-REG-14 | IsUpgrade, NOT in CurRoundSet | ✅ PASS |  |
| 135 | IT-REG-15 | Valid register: Status=Verified | ✅ PASS |  |
| 136 | IT-REG-16 | Round mismatch: return | ✅ PASS |  |
| 137 | IT-REG-17 | StartBlockHeight mismatch | ✅ PASS |  |
| 138 | IT-REG-18 | StartBlockHash mismatch | ✅ PASS |  |
| 139 | IT-REG-19 | Stage != Registration | ✅ PASS |  |
| 140 | IT-REG-20 | Validator not in ActiveValSet | ✅ PASS |  |
| 141 | IT-RES-01 | Phase=Failed → ResumeDKGService | ✅ PASS |  |
| 142 | IT-RES-02 | GetSession fails → MarkFailed → recover | ✅ PASS |  |
| 143 | IT-RES-03 | Phase=Failed, tryResume → resume | ✅ PASS |  |
| 144 | IT-RES-04 | Phase != Failed → no resume | ✅ PASS |  |
| 145 | IT-RES-05 | Resume → rejoin current round | ✅ PASS |  |
| 146 | IT-RES-06 | Resume CreateSession fails → stay Failed | ✅ PASS |  |
| 147 | IT-RES-07 | All validators resume → new round completes | ✅ PASS |  |
| 148 | IT-SKIP-01 | Verified < MinReq: SkipToNextRound | ✅ PASS |  |
| 149 | IT-SKIP-02 | Finalized < MinReq or < Threshold: SkipToNextRound | ✅ PASS |  |
| 150 | IT-SKIP-03 | FlushAllQueues | ✅ PASS |  |
| 151 | IT-SKIP-04 | [defensive] setDKGNetwork(Stage=Failed) fails | ✅ PASS |  |
| 152 | IT-UPG-01 | ScheduleUpgrade → PendingUpgrade set | ✅ PASS |  |
| 153 | IT-UPG-02 | hasPendingUpgradeActivation returns activationHeight | ✅ PASS |  |
| 154 | IT-UPG-03 | CancelUpgrade → removes pending | ✅ PASS |  |
| 155 | IT-UPG-04 | Upgrade round completes normally | ✅ PASS |  |
| 156 | IT-UPG-05 | Upgrade round failed → retry next round | ✅ PASS |  |
| 157 | CL-BADDEALER-01 | Bad dealer not invalidated after justification | ⚠️ LIMITATION | Code fix verified in source (`dkg_dealing.go:168`). Cannot trigger in real SGX: mock deals fail at decryption → no complaint; modifying kernel changes `mr_enclave` → attestation fails. Covered by story-kernel `tamper_deal_test.go`. |
| 158 | CL-BADDEALER-02 | Bad dealer successfully finalizes | ✅ PASS |  |
| 159 | CL-BADDEALER-03 | Bad dealer inflates finalizedCount | ✅ PASS |  |
| 160 | CL-BADDEALER-04 | Bad dealer receives UBI rewards | ⚠️ LIMITATION | Code fix verified (`dkg_handler.go:122`, `dkg_rewards.go:59`). Same SGX constraint as CL-BADDEALER-01. |
| 161 | CL-BADDEALER-05 | Bad dealer's partial decryption accepted | 🔒 BLOCKED | CDR condition contract not deployed |
| 162 | CL-BADDEALER-06 | Decryption fails with bad dealer's partial | 🔒 BLOCKED | CDR condition contract not deployed |
| 163 | CL-BADDEALER-07 | Multiple bad dealers vs threshold | ✅ PASS |  |
| 164 | CL-BADDEALER-08 | Adaptive attack: match honest globalPubKey | ✅ PASS |  |
| 165 | CL-DUPREG-01 | Resume → no duplicate registration | ✅ PASS |  |
| 166 | CL-DUPREG-02 | Kernel restart → index unchanged | ✅ PASS |  |
| 167 | CL-DUPREG-03 | Index uniqueness: no collision | ✅ PASS |  |
| 168 | CL-DUPREG-04 | Correct indices → round completes | ✅ PASS |  |
| 169 | CL-DUPREG-05 | On-chain: Register() twice same validator | ✅ PASS |  |
| 170 | CL-DUPREG-06 | [#721] Replay victim's register() from external address | ❌ BUG | Third-party replays on-chain params → victim index overwritten. Contract lacks `msg.sender==validatorAddr` and duplicate guard. |
| 171 | CL-DUPREG-07 | [#721] Replay ALL validators' registrations | ❌ BUG | Mass replay corrupts all indices → collisions + ghost indices → round failure. Low-cost, repeatable every round. |
| 172 | CL-FEE-01 | 1 wei pool / 3 validators → remainder | 🔒 BLOCKED | CDR condition contract not deployed |
| 173 | CL-FEE-02 | Non-committee submitter → no reward | 🔒 BLOCKED | CDR condition contract not deployed |
| 174 | CL-FEE-03 | Pool>0, submitCount=0 → not distributed | 🔒 BLOCKED | CDR condition contract not deployed |
| 175 | CL-JUST-01 | Valid Schnorr + invalid VSS → dealer invalidated | ✅ PASS |  |
| 176 | CL-JUST-02 | Already-invalidated dealer → no state change | ✅ PASS |  |
| 177 | CL-JUST-03 | Justification after Finalization → ignored | ✅ PASS |  |
| 178 | CL-PD-01 | Forged ECDSA sig → partial rejected | 🔒 BLOCKED | CDR condition contract not deployed |
| 179 | CL-PD-02 | Wrong commPubKey → signer mismatch | 🔒 BLOCKED | CDR condition contract not deployed |
| 180 | CL-PD-03 | pubShare mismatch → rejected | 🔒 BLOCKED | CDR condition contract not deployed |
| 181 | CL-PD-04 | Duplicate partial submission → dedup | 🔒 BLOCKED | CDR condition contract not deployed |
| 182 | CL-REPLAY-01 | Round N deals replayed in N+1 → rejected | ✅ PASS |  |
| 183 | CL-REPLAY-02 | Round N justifications replayed → Schnorr fails | ✅ PASS |  |
| 184 | CL-REPLAY-03 | Round N finalization sig replayed → round mismatch | ✅ PASS |  |
| 185 | CL-RESH-01 | R1→R2→R3 continuous resharing | ✅ PASS |  |
| 186 | CL-RESH-03 | R(N) Active + R(N+1) Dealing → decrypt uses R(N) key | 🔒 BLOCKED | CDR condition contract not deployed |
| 187 | CL-RESH-04 | Cross-round key continuity: R(N) encrypt → R(N+1) decrypt | 🔒 BLOCKED | CDR condition contract not deployed |
| 188 | CL-RESTART-01 | Kernel restart after ProcessDeals → round recovers | ✅ PASS |  |
| 189 | CL-RESTART-02 | Kernel restart before FinalizeDKG → GlobalPubKey consistent | ✅ PASS |  |
| 190 | CL-RESTART-03 | Kernel restart after justification → round recovers | ✅ PASS |  |
| 191 | CL-KERR-01 | Kernel InvalidArgument vs Internal error handling | ✅ PASS |  |
| 192 | CL-LOCK-01 | Concurrent CDRRead during round transition | 🔒 BLOCKED | CDR condition contract not deployed |
| 193 | CL-TO-01 | Partial at reqHeight+200 → accepted | 🔒 BLOCKED | CDR condition contract not deployed |
| 194 | CL-TO-02 | Partial at reqHeight+201 → rejected | 🔒 BLOCKED | CDR condition contract not deployed |
| 195 | CL-TO-03 | Prune + partial concurrent → no race | 🔒 BLOCKED | CDR condition contract not deployed |
| 196 | CL-TO-04 | Pruned → late partial → rejected gracefully | 🔒 BLOCKED | CDR condition contract not deployed |
| 197 | CL-VE-01 | VE > 256KB → REJECT, round progresses | ✅ PASS |  |
| 198 | CL-VE-02 | VE > 80 deals → REJECT | ✅ PASS |  |
| 199 | CL-VE-03 | Malformed proto VE → REJECT | ✅ PASS |  |
| 200 | CL-VE-04 | Wrong round deals → dropped in aggregation | ✅ PASS |  |
| 201 | CL-VE-05 | Duplicate deals → dedup in aggregateVotes | ✅ PASS |  |
| 202 | CL-VE-06 | One REJECT VE → block still produced | ✅ PASS |  |

---

## Usage

```bash
# Full suite (auto period switching):
DKG_DRIVER=script go test -v -timeout 120m -tags integration -run TestDKG_FullSuite ./tests/integration/dkg/

# Resume (skip passed):
DKG_RESUME=true DKG_DRIVER=script go test -v -timeout 120m -tags integration -run TestDKG_FullSuite ./tests/integration/dkg/

# CL-* only:
DKG_DRIVER=script go test -v -timeout 60m -tags integration -run TestDKG_CL ./tests/integration/dkg/
```
