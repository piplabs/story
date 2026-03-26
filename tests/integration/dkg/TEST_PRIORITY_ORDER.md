# DKG/CDR E2E Test Priority Order

**229 test cases · P0=44 · P1=108 · P2=77**

```bash
# P0 only:
go test -tags=integration -run TestDKG_PrioritySuite/P0 -timeout 60m
# P0+P1:
go test -tags=integration -run "TestDKG_PrioritySuite/(P0|P1)" -timeout 120m
# Full:
go test -tags=integration -run TestDKG_PrioritySuite -timeout 180m
```

---


## P0 — Happy Path + Must-Work

Production daily paths. Failure = system unusable.

| # | Case ID | Description | Reason |
|---|---------|-------------|--------|
| 1 | CL-AUDIT-01 | [STOR-28] CDR decrypt works >1min after Active (worker timeout) | [STOR-28] Worker timeout fix regression |
| 2 | CL-AUDIT-04 | [STOR-9] CDR read path alive during Active stage | [STOR-9] CDR read path fix regression |
| 3 | CL-CDR-01 | CDR happy path: allocate → write → read → partials → threshold met | CDR core: allocate→write→read→decrypt |
| 4 | CL-CDR-02 | Cross-round CDR: write in R(N) → read in R(N+1) → decrypt succeeds | Cross-round: old data must decrypt after resharing |
| 5 | CL-CDR-03 | Multiple CDR read/write cycles within same round | Multiple vaults in parallel |
| 6 | CL-CRASH-01 | Single story validator crash: stop val2 story → chain continues (2/3) → restart → catches up → DKG resumes | Single validator crash (production certainty) |
| 7 | CL-CRASH-02 | All nodes crash: stop all story+kernel → restart → consensus recovers → DKG resumes | All nodes crash (power failure scenario) |
| 8 | CL-DECRYPT-RESUME-01 | Decrypt worker resume: Active CDR Read → stop kernel → restart → partials submitted | Decrypt worker restart (STOR-28 fixed) |
| 9 | CL-FEE-05 | [STOR-15] CDR fee wei/gwei 1e9 mismatch verification | [STOR-15] HIGH: 1e9 fee amplification bug |
| 10 | CL-RECONNECT-01 | Kernel auto-reconnect: stop kernel → restart (not story) → 3/3 registration | Kernel ops restart → auto-reconnect |
| 11 | CL-RESH-01 | R1→R2→R3 continuous resharing | Continuous key rotation R1→R2→R3 |
| 12 | CL-RESH-04 | Cross-round key continuity: R(N) encrypt → R(N+1) decrypt | R(N) encrypt → R(N+1) decrypt continuity |
| 13 | CL-RESTART-01 | Kernel restart after ProcessDeals → round recovers | Mid-DKG kernel restart recovery |
| 14 | CL-RESTART-02 | Kernel restart before FinalizeDKG → GlobalPubKey consistent | PrivatePoly persistence (kernel #25) |
| 15 | IT-ACT-01 | FinalizeDKGRound: setLatestActiveRound, handleDKGComplete | Active round setup |
| 16 | IT-ACT-05 | IsUpgrade round: complete normally | Upgrade round completion |
| 17 | IT-ACT-06 | First round: settleRewards nil (no prevActive) | First round no prev rewards |
| 18 | IT-ACT-07 | DkgCommitteeRewardPortion=0 | Reward portion=0 boundary |
| 19 | IT-ACT-08 | UBI balance=0: rewards settled | UBI balance=0 settlement |
| 20 | IT-BB-02 | latestRound == nil: InitiateDKGRound(ctx, false) | First DKG round entry point |
| 21 | IT-BB-03 | Pending upgrade + height >= ActivationHeight | Upgrade activation trigger |
| 22 | IT-BB-04 | Stage=Registration, elapsed >= RegistrationPeriod → Dealing | Core stage transition: Reg→Deal |
| 23 | IT-BB-05 | Stage=Dealing → Finalization | Core stage transition: Deal→Fin |
| 24 | IT-BB-06 | Stage=Finalization → Active | Core stage transition: Fin→Active |
| 25 | IT-CDR-01 | FeeCollected → AddCDRFeeToPool, pool balance increases | Fee collection into pool |
| 26 | IT-CDR-02 | EncryptedPartialDecryption → IncrementCDRPartialSubmitCount, RefundCDRFee | Partial decrypt count + refund |
| 27 | IT-CDR-06 | No previous active round: distributeCDRRewardPool returns nil | First round no distribution |
| 28 | IT-DEC-01 | ProcessCDRVaultRead: latestActive exists | CDR vault read entry |
| 29 | IT-DEC-16 | Vault not found for decrypt | Vault not found handling |
| 30 | IT-DL-01 | Verified >= MinReq: BeginDealing | Dealing start condition |
| 31 | IT-E2E-01 | Full happy path: Registration → Active with GlobalPubKey | DKG core flow, no key = no CDR |
| 32 | IT-E2E-02 | Resharing round (non-upgrade) | Auto-triggered every 21 days |
| 33 | IT-E2E-03 | Failed round: insufficient Verified | Graceful degradation when validators insufficient |
| 34 | IT-E2E-04 | Failed round: insufficient Finalized | Graceful degradation when finalization insufficient |
| 35 | IT-E2E-06 | E2E Upgrade resharing round | Upgrade is mandatory ops path |
| 36 | IT-ENC-01 | TDH2 encrypt + CDR.write → VaultWritten | Encrypt + write entry |
| 37 | IT-FN-01 | BeginFinalization: emit, go handleDKGFinalization | Finalization start |
| 38 | IT-FN-04 | callTEEFinalizeDKG succeeds | TEE finalize success path |
| 39 | IT-REG-01 | First round: round 1, Stage=Registration, IsResharing=false | First round registration |
| 40 | IT-REG-02 | Previous round ended: new round IsResharing=true | Resharing registration |
| 41 | IT-RES-01 | Phase=Failed → ResumeDKGService | Failure recovery entry point |
| 42 | IT-SKIP-03 | FlushAllQueues | Round cleanup prevents cross-round pollution |
| 43 | IT-UPG-01 | ScheduleUpgrade → PendingUpgrade set | Upgrade scheduling |
| 44 | IT-UPG-04 | Upgrade round completes normally | Upgrade round completion |

## P1 — Adversarial + Bug Regression

Contract-layer attacks + known issue regression. No kernel control needed.

| # | Case ID | Description | Reason |
|---|---------|-------------|--------|
| 45 | CL-AUDIT-02 | [STOR-13] pubKeyShare finalization vs partial submission consistency | [STOR-13] pubKeyShare prefix |
| 46 | CL-AUDIT-03 | [STOR-8] CDR decrypt works after resharing (PIDCache populated) | [STOR-8] PIDCache |
| 47 | CL-AUDIT-05 | [STOR-16] LatestActiveRound updated correctly at round transition | [STOR-16] Round rollover |
| 48 | CL-AUDIT-06 | [STOR-22] Finalization events at stage boundary not dropped | [STOR-22] Boundary events |
| 49 | CL-COND-01 | [CDR-006] Condition contract as msg.sender bypasses access control | [CDR-006] Condition bypass |
| 54 | CL-COND-02 | [CDR-015/M-01] allocate with readConditionAddr=0 → vault unreadable | [CDR-015] Unreadable vault |
| 55 | CL-DUPREG-01 | Resume → no duplicate registration | #703: off-chain guard |
| 56 | CL-DUPREG-02 | Kernel restart → index unchanged | #703: restart index stable |
| 57 | CL-DUPREG-03 | Index uniqueness: no collision | #703: index uniqueness |
| 58 | CL-DUPREG-04 | Correct indices → round completes | #703: round completes |
| 59 | CL-DUPREG-05 | On-chain: Register() twice same validator | #703: on-chain dedup |
| 60 | CL-DUPREG-06 | [#721] Replay victim's register() from external address | #721: third-party replay |
| 61 | CL-DUPREG-07 | [#721] Replay ALL validators' registrations | #721: mass replay |
| 62 | CL-FEE-01 | 1 wei pool / 3 validators → remainder | Integer division remainder |
| 63 | CL-FEE-03 | Pool>0, submitCount=0 → not distributed | Fee pool carryover |
| 65 | CL-FEE-04 | [CDR-003] CDR reward distribution consistency across validators | [CDR-003] Sorted map iteration |
| 66 | CL-FIN-01 | [CDR-005/M-08] finalize() from non-validator address → rejected by CL | [CDR-005] Permissionless finalize() |
| 67 | CL-FIN-02 | [M-02] Active round always has non-empty GlobalPublicKey | [M-02] GPK guard |
| 68 | CL-FLUSH-01 | [H-04] No cross-round data contamination at round boundary | [H-04] Cross-round cleanup |
| 69 | CL-PS-01 | [STOR-4] Submit partial for unknown/expired request → not rewarded | [STOR-4] Permissionless submitPartial |
| 70 | CL-PS-02 | Submit partial with wrong round → rejected | Contract round validation |
| 71 | CL-PS-03 | Submit partial with mismatched ciphertext → rejected | Contract ciphertext validation |
| 72 | CL-PS-04 | Submit partial from unregistered address → rejected | Contract registration check |
| 73 | CL-PS-05 | [STOR-3] Duplicate partial replay → no double reward | [STOR-3] Duplicate replay |
| 74 | CL-PS-06 | [L1-02] Oversized ciphertext in CDR write | [L1-02] State bloat |
| 76 | CL-TO-01 | Partial at reqHeight+200 → accepted | Timeout boundary exact |
| 77 | CL-TO-02 | Partial at reqHeight+201 → rejected | Timeout boundary +1 |
| 78 | CL-TO-03 | Prune + partial concurrent → no race | Concurrent prune safety |
| 79 | CL-TO-04 | Pruned → late partial → rejected gracefully | Expired request handling |
| 80 | IT-ACT-02 | finalizedCount < MinReq: SkipToNextRound | Chain validation logic |
| 81 | IT-ACT-03 | finalizedCount < Threshold: SkipToNextRound | Chain validation logic |
| 82 | IT-ACT-04 | isDKGSvcEnabled=false: no handleDKGComplete | Chain validation logic |
| 83 | IT-ACT-09 | settleRewards: WithdrawUbiToModule, distribute | Chain validation logic |
| 84 | IT-ACT-10 | handleDKGComplete: PhaseCompleted | Chain validation logic |
| 85 | IT-ACT-12 | Phase != PhaseFinalized: MarkFailed | Chain validation logic |
| 86 | IT-CDR-03 | distributeCDRRewardPool: validators receive share by count | Chain validation logic |
| 87 | IT-CDR-04 | FeeCollected amount=0: no AddCDRFeeToPool | Chain validation logic |
| 88 | IT-CDR-05 | Pool balance < refund amount: error | Chain validation logic |
| 89 | IT-CDR-07 | Total submit count=0, pool>0: no SendCoins | Chain validation logic |
| 90 | IT-CDR-08 | E2E: CDR write/read + partial decryption + distributeCDRRewardPool | Chain validation logic |
| 91 | IT-DEC-02 | No latestActive: return early | Chain validation logic |
| 92 | IT-DEC-03 | GetSession fails: MarkFailed | Chain validation logic |
| 93 | IT-DEC-04 | Phase != Completed: skip decrypt | Chain validation logic |
| 94 | IT-DEC-05 | callTEEDecrypt fails: skip | Chain validation logic |
| 95 | IT-DEC-06 | callContractSubmitPartial fails | Chain validation logic |
| 96 | IT-DEC-07 | EncryptedPartialDecryption: IncrementCount | Chain validation logic |
| 97 | IT-DEC-08 | EncryptedPartialDecryption: RefundCDRFee | Chain validation logic |
| 98 | IT-DEC-09 | RefundCDRFee: pool balance check | Chain validation logic |
| 99 | IT-DEC-10 | RefundCDRFee: SendCoins from pool | Chain validation logic |
| 100 | IT-DEC-11 | RefundCDRFee: pool < refundAmt error | Chain validation logic |
| 101 | IT-DEC-12 | Multiple partials reach threshold | Chain validation logic |
| 102 | IT-DEC-13 | Partial from non-active-round validator | Chain validation logic |
| 103 | IT-DEC-14 | Duplicate partial submission | Chain validation logic |
| 104 | IT-DEC-15 | Invalid partial proof | Chain validation logic |
| 105 | IT-DEC-17 | Decrypt request for empty vault | Chain validation logic |
| 106 | IT-DEC-18 | Concurrent decrypt requests | Chain validation logic |
| 107 | IT-DEC-19 | Decrypt after round rotation | Chain validation logic |
| 108 | IT-DEC-20 | Decrypt with stale GlobalPublicKey | Chain validation logic |
| 109 | IT-DEC-21 | callTEEDecrypt timeout | Chain validation logic |
| 110 | IT-DEC-22 | Full decrypt E2E path | Chain validation logic |
| 111 | IT-DL-02 | Verified < MinReq: SkipToNextRound | Chain validation logic |
| 112 | IT-DL-03 | isDKGSvcEnabled=false: no handleDKGDealing | Chain validation logic |
| 113 | IT-DL-06 | Old member (not in CurRoundSet) still participates in dealing | Chain validation logic |
| 114 | IT-DL-07 | callTEEGenerateDeals succeeds | Chain validation logic |
| 115 | IT-DL-08 | TEE down: callTEEGenerateDeals fails | Chain validation logic |
| 116 | IT-DL-09 | Upgrade round dealing | Chain validation logic |
| 117 | IT-DL-10 | TEE down: callTEEVerifyDeals fails | Chain validation logic |
| 118 | IT-DL-11 | callTEEVerifyDeals succeeds | Chain validation logic |
| 119 | IT-DL-12 | Deal broadcast via VoteExtension | Chain validation logic |
| 120 | IT-DL-13 | Response broadcast via VoteExtension | Chain validation logic |
| 121 | IT-DL-14 | VerifyVoteExtension: valid deal accepted | Chain validation logic |
| 122 | IT-DL-15 | Invalid deal detected: complaint generated | Chain validation logic |
| 123 | IT-DL-16 | Justification for complaint resolves deal | Chain validation logic |
| 124 | IT-DL-17 | All deals valid: no complaints | Chain validation logic |
| 125 | IT-FN-02 | isDKGSvcEnabled=false: Emit only, no goroutine | Chain validation logic |
| 126 | IT-FN-03 | Finalization signature verification failed | Chain validation logic |
| 127 | IT-FN-05 | callContractFinalizeDKG fails with invalid params | Chain validation logic |
| 128 | IT-FN-06 | Double finalization rejected | Chain validation logic |
| 129 | IT-FN-07 | TEE down: handleDKGFinalization callTEE fails | Chain validation logic |
| 130 | IT-FN-08 | TEE down: node does not submit finalize | Chain validation logic |
| 131 | IT-FN-09 | Invalidated dealer cannot finalize | Chain validation logic |
| 132 | IT-FN-10 | Finalized count >= Threshold: GlobalPublicKey set | Chain validation logic |
| 133 | IT-FN-11 | GlobalPubKey vote counting | Chain validation logic |
| 134 | IT-FN-12 | GlobalPubKey set when votes >= threshold | Chain validation logic |
| 135 | IT-FN-13 | Finalization with minimum threshold | Chain validation logic |
| 136 | IT-FN-14 | All validators finalize successfully | Chain validation logic |
| 137 | IT-REG-03 | Upgrade activation: IsUpgrade=true | Chain validation logic |
| 138 | IT-REG-04 | isDKGSvcEnabled=false: no handleDKGRegistration | Chain validation logic |
| 139 | IT-REG-05 | GetActiveValidators empty | Chain validation logic |
| 140 | IT-REG-06 | Stage=Registration, in CurRoundSet: CreateSession→Register | Chain validation logic |
| 141 | IT-REG-07 | NOT in CurRoundSet: no Register | Chain validation logic |
| 142 | IT-REG-10 | CreateSession fails: MarkFailed | Chain validation logic |
| 143 | IT-REG-11 | callTEEGenerateAndSealKey fails | Chain validation logic |
| 144 | IT-REG-12 | callContractRegister fails | Chain validation logic |
| 145 | IT-REG-13 | IsUpgrade, in CurRoundSet | Chain validation logic |
| 146 | IT-REG-14 | IsUpgrade, NOT in CurRoundSet | Chain validation logic |
| 147 | IT-REG-15 | Valid register: Status=Verified | Chain validation logic |
| 148 | IT-REG-16 | Round mismatch: return | Chain validation logic |
| 149 | IT-REG-17 | StartBlockHeight mismatch | Chain validation logic |
| 150 | IT-REG-18 | StartBlockHash mismatch | Chain validation logic |
| 151 | IT-REG-19 | Stage != Registration | Chain validation logic |
| 152 | IT-REG-20 | Validator not in ActiveValSet | Chain validation logic |
| 153 | IT-RES-02 | GetSession fails → MarkFailed → recover | Chain validation logic |
| 154 | IT-RES-03 | Phase=Failed, tryResume → resume | Chain validation logic |
| 155 | IT-RES-04 | Phase != Failed → no resume | Chain validation logic |
| 156 | IT-RES-05 | Resume → rejoin current round | Chain validation logic |
| 157 | IT-RES-06 | Resume CreateSession fails → stay Failed | Chain validation logic |
| 158 | IT-RES-07 | All validators resume → new round completes | Chain validation logic |

## P2 — Edge/Defensive + Requires Mock Kernel

SGX prevents kernel control in production. Defensive guards + boundary conditions.

| # | Case ID | Description | Reason |
|---|---------|-------------|--------|
| 159 | CL-BADDEALER-01 | Bad dealer not invalidated after justification | Needs precise VSS-invalid deal (mock too coarse) |
| 160 | CL-BADDEALER-02 | Bad dealer successfully finalizes | Needs precise VSS-invalid deal (mock too coarse) |
| 161 | CL-BADDEALER-03 | Bad dealer inflates finalizedCount | Needs precise VSS-invalid deal (mock too coarse) |
| 162 | CL-BADDEALER-04 | Bad dealer receives UBI rewards | Needs precise VSS-invalid deal (mock too coarse) |
| 163 | CL-BADDEALER-05 | Bad dealer's partial decryption accepted | Mock kernel + fake DCAP |
| 160 | CL-BADDEALER-06 | Decryption fails with bad dealer's partial | Mock kernel + fake DCAP |
| 161 | CL-BADDEALER-07 | Multiple bad dealers vs threshold | Mock kernel + fake DCAP |
| 162 | CL-BADDEALER-08 | Adaptive attack: match honest globalPubKey | Mock kernel + fake DCAP |
| 163 | CL-JUST-01 | Valid Schnorr + invalid VSS → dealer invalidated | Needs kernel bad deal injection |
| 164 | CL-JUST-02 | Already-invalidated dealer → no state change | Needs kernel bad deal injection |
| 165 | CL-JUST-03 | Justification after Finalization → ignored | Needs kernel bad deal injection |
| 166 | CL-KERR-01 | Kernel InvalidArgument vs Internal error handling | Needs kernel error control |
| 167 | CL-LOCK-01 | Concurrent CDRRead during round transition | Extremely low probability race |
| 168 | CL-PD-01 | Forged ECDSA sig → partial rejected | Needs kernel signing control |
| 169 | CL-PD-02 | Wrong commPubKey → signer mismatch | Needs kernel key control |
| 170 | CL-PD-03 | pubShare mismatch → rejected | Needs kernel pubShare control |
| 171 | CL-PD-04 | Duplicate partial submission → dedup | Needs mock kernel |
| 172 | CL-REPLAY-01 | Round N deals replayed in N+1 → rejected | Needs kernel replay control |
| 173 | CL-REPLAY-02 | Round N justifications replayed → Schnorr fails | Needs kernel replay control |
| 174 | CL-REPLAY-03 | Round N finalization sig replayed → round mismatch | Needs kernel replay control |
| 175 | CL-RESH-02 | Validator joins R2, leaves R3 → clean handoff | Needs dynamic validator changes |
| 176 | CL-RESH-03 | R(N) Active + R(N+1) Dealing → decrypt uses R(N) key | Needs precise timing |
| 177 | CL-VE-01 | VE > 256KB → REJECT, round progresses | Needs kernel control (SGX) |
| 178 | CL-VE-02 | VE > 80 deals → REJECT | Needs kernel control (SGX) |
| 179 | CL-VE-03 | Malformed proto VE → REJECT | Needs kernel control (SGX) |
| 180 | CL-VE-04 | Wrong round deals → dropped in aggregation | Needs kernel control (SGX) |
| 181 | CL-VE-05 | Duplicate deals → dedup in aggregateVotes | Needs kernel control (SGX) |
| 182 | CL-VE-06 | One REJECT VE → block still produced | Needs kernel control (SGX) |
| 183 | IT-ACT-11 | [defensive] PhaseCompleted and IsFinalized: return | Defensive guard |
| 184 | IT-ACT-13 | [defensive] dkgSvcRunning true: return | Defensive guard |
| 185 | IT-BB-01 | height < dkgStartBlock: BeginBlocker return nil | Only at chain start |
| 186 | IT-BB-07 | Stage=Active, elapsed >= ActiveEnd → new round | 21-day trigger |
| 187 | IT-BB-08 | Stage transition time not reached: no transition | Defensive: time not reached |
| 188 | IT-BB-09 | Pending upgrade, height < ActivationHeight: no activation | Defensive: height not reached |
| 189 | IT-BB-10 | No pending upgrade: hasPendingUpgradeActivation returns nil | Defensive: no upgrade returns nil |
| 190 | IT-DL-04 | [defensive] Stage != Dealing: handleDKGDealing returns early | Defensive guard |
| 191 | IT-DL-05 | [defensive] dkgSvcRunning true: concurrent call ignored | Defensive guard |
| 192 | IT-DL-18 | [implicit] VE internal logic case 18 | Implicit/boundary/defensive |
| 193 | IT-DL-19 | [implicit] VE internal logic case 19 | Implicit/boundary/defensive |
| 194 | IT-DL-20 | [implicit] VE internal logic case 20 | Implicit/boundary/defensive |
| 195 | IT-DL-21 | [implicit] VE internal logic case 21 | Implicit/boundary/defensive |
| 196 | IT-DL-22 | [implicit] VE internal logic case 22 | Implicit/boundary/defensive |
| 197 | IT-DL-23 | [implicit] VE internal logic case 23 | Implicit/boundary/defensive |
| 198 | IT-DL-24 | [implicit] VE internal logic case 24 | Implicit/boundary/defensive |
| 199 | IT-DL-25 | [implicit] VE internal logic case 25 | Implicit/boundary/defensive |
| 200 | IT-DL-26 | [implicit] VE internal logic case 26 | Implicit/boundary/defensive |
| 201 | IT-DL-27 | [implicit] VE internal logic case 27 | Implicit/boundary/defensive |
| 202 | IT-DL-28 | [implicit] VE internal logic case 28 | Implicit/boundary/defensive |
| 203 | IT-DL-29 | [implicit] VE internal logic case 29 | Implicit/boundary/defensive |
| 204 | IT-DL-30 | [implicit] VE internal logic case 30 | Implicit/boundary/defensive |
| 205 | IT-DL-31 | [implicit] VE internal logic case 31 | Implicit/boundary/defensive |
| 206 | IT-DL-32 | [implicit] VE internal logic case 32 | Implicit/boundary/defensive |
| 207 | IT-DL-33 | [implicit] VE internal logic case 33 | Implicit/boundary/defensive |
| 208 | IT-DL-34 | [implicit] VE internal logic case 34 | Implicit/boundary/defensive |
| 209 | IT-DL-35 | [implicit] VE internal logic case 35 | Implicit/boundary/defensive |
| 210 | IT-E2E-05 | Complaint / Justification path | Needs mock kernel bad deal |
| 211 | IT-EDGE-01 | Minimum validators boundary | Implicit/boundary/defensive |
| 212 | IT-EDGE-02 | MinReq params boundary | Implicit/boundary/defensive |
| 213 | IT-EDGE-03 | Max validators boundary | Implicit/boundary/defensive |
| 214 | IT-EDGE-04 | Failed round recovery | Implicit/boundary/defensive |
| 215 | IT-EDGE-05 | Period values boundary | Implicit/boundary/defensive |
| 216 | IT-EDGE-06 | DKG completes with one TEE down | Implicit/boundary/defensive |
| 217 | IT-FN-15 | [implicit] Finalization VE case 15 | Implicit/boundary/defensive |
| 218 | IT-FN-16 | [implicit] Finalization VE case 16 | Implicit/boundary/defensive |
| 219 | IT-FN-17 | [implicit] Finalization VE case 17 | Implicit/boundary/defensive |
| 220 | IT-FN-18 | [implicit] Finalization VE case 18 | Implicit/boundary/defensive |
| 221 | IT-FN-19 | [implicit] Finalization VE case 19 | Implicit/boundary/defensive |
| 222 | IT-REG-08 | [defensive] Stage != Registration: return | Defensive guard |
| 223 | IT-REG-09 | [defensive] dkgSvcRunning true: return | Defensive guard |
| 224 | IT-SKIP-01 | Verified < MinReq: SkipToNextRound | Overlaps IT-DL-02 |
| 225 | IT-SKIP-02 | Finalized < MinReq or < Threshold: SkipToNextRound | Overlaps IT-ACT-02 |
| 226 | IT-SKIP-04 | [defensive] setDKGNetwork(Stage=Failed) fails | Defensive guard |
| 227 | IT-UPG-02 | hasPendingUpgradeActivation returns activationHeight | Covered by IT-BB-03 |
| 228 | IT-UPG-03 | CancelUpgrade → removes pending | LIMITATION |
| 229 | IT-UPG-05 | Upgrade round failed → retry next round | Low probability upgrade failure |
