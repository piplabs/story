# Blockchain Security Reviewer - Memory

## Project: Story Protocol DKG Module

### Architecture (confirmed 2026-02-18)
- DKG finalization signature verification moved from DKG.sol to Go consensus layer (dkg_handler.go)
- TEE (story-dkg-tee) signs with a sealed secp256k1 key; commPubKey stored as 64-byte uncompressed (no 0x04 prefix)
- TEE V-value convention: adds 27 before returning signature (sig[64] += 27 if < 27)
- Go handler convention: subtracts 27 before crypto.SigToPub (sig[64] -= 27 if >= 27)
- Address derivation: common.BytesToAddress(crypto.Keccak256(commPubKey64)) is equivalent to Ethereum standard

### Key Files
- `client/x/dkg/keeper/dkg_handler.go` — verifyFinalizationSignature (line 262), Finalized (line 96)
- `client/x/dkg/keeper/dkg_registration.go` — getDKGRegistration, updateDKGRegistrationStatus
- `contracts/src/protocol/DKG.sol` — signature verification removed from finalizeDKG; now just emits event
- `client/x/dkg/keeper/dkg_justification.go` — verifyJustificationSignature, verifyJustification, deduplicateJustifications, buildDealerPubKeyMap, MaxJustificationsPerBlock
- `client/x/dkg/keeper/dkg_dealing.go` — ProcessJustifications (3-line async pattern), ProcessDeals, ProcessResponses
- `client/x/dkg/keeper/dkg_svc_dealing.go` — handleDKGProcessJustifications (async pipeline: cap→ActiveValSet→session→sig verify→dedup→VSS→forward)
- `client/x/dkg/keeper/queue.go` — package-level queues (global state, no round scoping)
- `client/x/dkg/keeper/vote.go` — ExtendVote, aggregateVotes (no dedup)
- `lib/vss/verify.go` — VerifyPedersenVSS (Edwards25519, MaxCommitments=80 guard added)

### Encoding/Abi Notes
- Solidity abi.encodePacked for uint32 = 4 bytes big-endian (matches Go binary.BigEndian.PutUint32)
- participantsRoot computed identically in TEE (dkg_finalize.go) and Go handler (validateParticipantsRoot)
- kyber PriPoly.Eval(i): x-coordinate = 1+i (0-based i → x=1 means Eval(0)). PubPoly.Check(PriShare{I}) evaluates at x=1+I.
- verifyJustification passes recipientIndex = secShare.GetI()+1 (converts 0-based proto to 1-based VerifyPedersenVSS API)
- VerifyPedersenVSS internally does PriShare.I = recipientIndex-1 (back to 0-based for kyber)
- Net result: x-coordinate = 1 + (recipientIndex-1) = recipientIndex. Correct.

### Justification Processing Flow (2026-03-03 final state)
1. ExtendVote dequeues up to 10 justifications per validator
2. aggregateVotes concatenates all validators' justifications (N*10 total, no dedup)
3. ProcessJustifications: emit event (ALL justifications embedded), then spawn goroutine
4. handleDKGProcessJustifications (async): cap to 10 → ActiveValSet check → session/phase check
   → buildDealerPubKeyMap → Schnorr sig verify (filter) → dedup(dealerIdx,recipientIdx) → VSS verify → forward to TEE
5. invalidateDealerRegistration is NOT called from justification processing (off-chain only)

### TEE Upgrade Mechanism — STATUS AS OF 2026-03-05 (RE-AUDIT)
- Key files: abci.go (checkPendingUpgradeActivation), dkg_handler.go (UpgradeScheduled/UpgradeCancelled), kernel_upgrade.go (CRUD), kernel_router.go, dkg_svc_registration.go (getRegistrationKernelClient), dkg_svc_dealing.go (dealerCC routing)
- DKG.sol scheduleUpgrade/cancelUpgrade are onlyOwner, event-only (no on-chain state)
- **FIXED (H-1)**: activationHeight overflow — IsInt64() + positivity check at evmengine/keeper/dkg.go:355; now int64 end-to-end. VERIFIED CLEAN.
- **FIXED (M-data-race)**: processDecryptQueue data race — GetDecryptRequests()/SetDecryptRequests() added in types/dkg.go, used in dkg_svc.go. VERIFIED CLEAN (race detector passes).
- **FIXED (L-empty-version)**: UpgradeScheduled now checks len(upgradeVersion)==0 at dkg_handler.go:205. VERIFIED CLEAN.
- **FIXED (L-half-commit)**: checkPendingUpgradeActivation rolls back IsActivated=false if InitiateDKGRound fails at abci.go:106-114. VERIFIED CLEAN.
- **FIXED (gRPC-leak)**: KernelRouter.closers map + CreateTEEClient returns io.Closer; ConnectAndDiscover closes on error. VERIFIED CLEAN.
- **FIXED (genesis)**: KernelUpgradeInfos in GenesisState proto; InitGenesis/ExportGenesis wired. VERIFIED CLEAN.
- **FIXED (M-1)**: GetClient() strict CC matching only. VERIFIED STILL INTACT.
- **FIXED (M-2)**: Duplicate upgrade check via GetPendingUpgrade() before scheduling. VERIFIED STILL INTACT.
- **FIXED (M-3)**: getRegistrationKernelClient returns error if no new CC. VERIFIED STILL INTACT.
- **FIXED (L-2)**: UpgradeCancelled returns error when no pending. VERIFIED STILL INTACT.
- **NEW LOW (re-audit): Stale log read of session.DecryptRequests after SetDecryptRequests** — dkg_svc.go:127 reads session.DecryptRequests directly in log.Error after SetDecryptRequests(remaining). Not concurrent (same goroutine), but inconsistent (shows pre-Set slice len). Cosmetic issue, not a race.
- **NEW LOW (re-audit): Log read in dkg_handler.go:392** — len(session.DecryptRequests) read without lock after AddDecryptRequest(). Concurrent decrypt worker goroutine could be accessing same field. Technically benign (int read), but go race detector may flag it. Fix: use len(session.GetDecryptRequests()).
- **NEW LOW (re-audit): ConnectAndDiscover no close-on-replace** — if same codeCommitmentHex already in r.closers, old closer is silently overwritten without calling Close(). Leak if endpoint re-registers same binary.
- **NEW LOW (re-audit): ValidateGenesis does not validate KernelUpgradeInfos** — only validates Params; malformed genesis with empty upgradeVersion in KernelUpgradeInfos would bypass UpgradeScheduled validation and be stored directly.
- **OPEN INFO**: ThresholdDecryptRequested — no size limits on ciphertext/label/requesterPubKey; unbounded memory via DecryptRequests append.
- DKGThresholdDecryptRequestedEvent uses getEventSafe (zero-value ID) — never matched until event added to ABI (safe by design).

### Known Security Gaps — STATUS AS OF 2026-03-03
1. **FIXED**: Invalidation + double-finalization guards in Finalized() ✓
2. **FIXED**: Deduplication now in deduplicateJustifications (signature-first before dedup) ✓
3. **FIXED**: Session.Phase == PhaseDealing guard in handleDKGProcessJustifications ✓
4. **FIXED**: MockTEEClient now has ProcessJustification; tests build and pass ✓
5. **FIXED**: vss package now has 100% coverage ✓
6. **MEDIUM (open 2026-03-03): Justification plain-deal (SecShare) in consensus event** — EventBeginProcessJustifications embeds full PlainDeal including SecShare.V (secret scalar). All validators see this in block events, but this is protocol-correct: justifications are public reveals by design (dealer proving share validity). No confidentiality risk.
7. **MEDIUM (open 2026-03-03): MaxJustificationsPerBlock cap misalignment** — Cap (10) applied only in async handler after event is emitted with uncapped list. A block can embed N_validators × 10 justifications in the event. With 80 validators, that is 800 justifications per event. The async handler silently drops 790. No consensus impact, but event bloat is observable.
8. **LOW (open 2026-03-03): aggregateVotes still lacks deduplication** — same justification from N validators → N entries in the slice passed to ProcessJustifications. Schnorr verification filters these, but wastes CPU in the async handler pre-dedup.
9. **INFO: handleDKGProcessJustifications has 0% unit test coverage** — async path not directly tested.

### Test Coverage (2026-03-05 — re-audit)
- keeper package: 26.0% overall — 80% threshold STILL NOT MET (was 26.3%, effectively unchanged)
- BeginBlocker: 0.0%, checkPendingUpgradeActivation: 0.0%, UpgradeCancelled: 0.0%
- kernel_upgrade.go: SetKernelUpgradeInfo 75%, GetKernelUpgradeInfo 0%, GetAllKernelUpgradeInfos 0%, DeleteKernelUpgradeInfo 0%, SetKernelUpgradeInfos 0%, GetPendingUpgrade 44.4%
- getRegistrationKernelClient: 0.0%, getOldCodeCommitment: 0.0%
- handleDKGDealing/Finalization/Registration (svc_*): 0.0%
- processDecryptQueue/handleDecryptRequest: 0.0%, StartDecryptWorker: 0.0%
- ThresholdDecryptRequested: 0.0%
- KernelRouter: GetClient 100%, GetDefaultClient 100%, HasClients 100%, Disconnect 71.4%, ConnectAndDiscover 0%, RegisterClientForEndpoint 0%
- UpgradeScheduled (keeper): 66.7% — empty-version and duplicate-pending paths NOT tested
- UpgradeCancelled (keeper): 0.0%
- InitGenesis: 0.0%, ExportGenesis: 0.0%
- verifyFinalizationSignature: 100.0%, Finalized: 82.1%, Registered: 77.3%

### Security Patterns in this Codebase
- Uses go-ethereum crypto package (not standard lib crypto/ecdsa) for all secp256k1 operations
- Error wrapping via `github.com/piplabs/story/lib/errors`
- No reentrancy concern (Go consensus layer, not EVM)
- No flash loan vectors (DKG is a ceremony, not a financial primitive)
- Async handler pattern: dkgAsyncContext() gives 1-minute timeout; goroutines detach from consensus ctx
- Signature-before-dedup ordering is intentional security invariant (prevents shadow-entry attack)

### Final DKG Audit Fixes — STATUS AS OF 2026-03-10 (VERIFIED)
- H-03 Authority: depinject.go sets authority = govtypes.ModuleName by default; evmengine calls dkgKeeper.UpgradeScheduled etc. directly (no authority check needed, already trusted). msg_server.go and proposal_server.go both check msg.Authority == GetAuthority(). SECURE.
- H-06 VE Size: parseAndVerifyVoteExtension checks len>256KB BEFORE proto.Unmarshal, then count checks per category (max 80). Called by both VerifyVoteExtension and PrepareVotes. SECURE. NOTE: VerifyVoteExtension still returns error (not REJECT) on parse failure — OPEN HIGH from v2.0.0 audit.
- H-08 Upgrade deletion timing: abci.go sets IsActivated=true (not delete) on activation; deleteActivatedUpgradeInfo called only after FinalizeDKGRound succeeds. hasPendingUpgradeActivation filters IsActivated=true. SECURE.
- C-01 Phase transition: handleDKGDealing checks session.Phase != PhaseInitialized and returns early; sets PhaseDealing at end. ResumeDKGService sets PhaseDealing before spawning goroutine. SECURE.
- H-02 Dedup: deduplicateDeals/Responses/Justifications all in vote.go only (not duplicated in dkg_justification.go). Keys correct. First-wins semantics. SECURE.
- M-13 DKGStageEnded: endPreviousActiveRound only transitions Active→Ended (not other stages). shouldTransitionStage case DKGStageEnded returns false. SECURE.
- story-kernel crypto_wipe.go: zeroPrivateKey uses SetInt64(0) + new(big.Int) — NOTE: original big.Int backing array not zeroed, but new(big.Int) clears the D pointer (GC eligible). Acceptable for TEE context.
- story-kernel dkg_service.go: defer zeroPrivateKey used in FinalizeDKG and signPartialDecryptResponse. defer zeroBytes used for privShare.Bytes, sharedBytes, aesKey, ephemeral. SECURE.
- story-kernel dkg_state.go: FromRound field correctly serialized in dkgStateDisk as json:"from_round,omitempty". SECURE.
- story-kernel dkg_generate_deals.go: CachePID returns error if ownPID==0 (not found). GenerateDeals propagates this error. SECURE.
- story-kernel server.go: reflection.Register gated by cfg.GRPC.DebugMode. config.go DebugMode defaults false. SECURE.
- OPEN TODO (story-kernel): PartialDecryptTDH2 has comment "TEE should verify if the request transaction was indeed submitted to the canonical chain" — not yet implemented. Medium risk.
- deduplicateJustifications called in BOTH vote.go (aggregateVotes) and dkg_svc_dealing.go (handleDKGProcessJustifications) — correct, two separate dedup points at different layers.

### v2.0.0 Upgrade Handler — STATUS AS OF 2026-03-07 (AUDIT)
Key files: lib/netconf/upgrades.go, client/app/upgrades/v_2_0_0/{upgrades.go,constants.go},
  client/app/upgrades.go, client/app/prouter.go, client/x/dkg/keeper/{vote.go,abci.go},
  client/x/evmengine/keeper/abci.go, client/x/evmstaking/keeper/ubi.go

**CRITICAL (open): V200=0 for aeneid/mainnet — IsV200 always true before height is set**
- netconf/upgrades.go: AeneidChainID and StoryChainID have V200=0 (TBD)
- IsV200(chainID, blockNumber) = blockNumber >= 0 → always true for any block
- DKG BeginBlocker, vote extensions, UBI reward logic all fire from block 1 without upgrade handler running
- No startup validation guards against deploying with V200=0

**HIGH (open): VerifyVoteExtension returns Go error instead of REJECT status**
- vote.go:65: `return nil, errors.Wrap(err, "failed to parse vote extension")`
- Per ABCI++ spec: returning error != returning REJECT; error signals app failure
- A malformed vote extension from a peer causes abci.go to log "[BUG]" and propagate error
- Should return `ResponseVerifyVoteExtension{Status: REJECT}` for invalid extensions

**MEDIUM (open): enableVoteExtensions idempotency check is overbroad**
- upgrades.go:66: `if currentParams.Abci.VoteExtensionsEnableHeight > 0 { return nil }`
- If VE was set to a wrong height previously, upgrade handler silently skips correction
- Should check `== upgradeHeight+1`, not just `> 0`

**MEDIUM (open): UpgradeStoreLoader pre-adds all upgrades without dedup guard**
- If a future upgrade also has "dkg" in Added, both pre-add → duplicate entries in StoreUpgrades
- SDK behavior with duplicate Added keys is not hardened; risk increases as upgrade history grows

**LOW (open): No tests for v_2_0_0 upgrade package (0.0% coverage)**
- CreateUpgradeHandler, enableVoteExtensions, GetUpgradeHeight: zero test coverage
- Horace upgrade has 92.4% coverage as the gold standard

**INFO (open): Proposal compatibility during transition (blocks H+1 to H+2)**
- PrepareVotes/ProcessProposal both correctly use v200Height+1 threshold — timing is correct
- Hard fork: all validators must upgrade simultaneously or proposal rejection chain occurs

**Test Coverage (2026-03-07 — v2.0.0 audit)**
- client/app/upgrades/v_2_0_0: 0.0% (no test files)
- client/app: 24.9% (UpgradeStoreLoader 0%, GetUpgradeHeight 0%)
- client/x/dkg/keeper: 31.6% (VerifyVoteExtension 0%, PrepareVotes 0%, BeginBlocker 0%)
- client/x/evmstaking/keeper/ubi.go: 91.9% (IsV200 branch tested)
- lib/netconf/upgrades.go IsV200: 0.0%

### Deal/Response Caching Proposal — STATUS AS OF 2026-03-15 (SECURITY REVIEW)
- Proposal: cache unprocessed deals/responses in session file (PendingDeals/PendingResponses) when kernel unavailable, replay in ResumeDKGService
- **CRITICAL (open): Consensus divergence by design** — deals from block B replayed at block B+K violates consensus ordering; different validators process deals at different heights → TEE DKG state diverges; VssResponses generated at different times become incoherent
- **HIGH (open): No disk integrity protection** — session files at 0644 allow admin-level attacker to tamper Cipher/DhKey in PendingDeals; replayed corrupted deals trigger false complaints against honest dealers
- **HIGH (open): Stage-transition replay** — ResumeDKGService replay can run after session reaches PhaseFinalized/PhaseCompleted; must gate on Phase in {PhaseInitialized, PhaseDealing}
- **MEDIUM (open): Per-block replay rate** — ResumeDKGService called every block (BeginBlocker + AddVote); replay adds 3×retry×2s goroutines per call during kernel downtime
- **MEDIUM (open): No cap on PendingDeals size** — deduplication applied at vote extension aggregation but not at cache append; unbounded growth possible
- **MEDIUM (open): dkgKernelMu deadlock/race** — replay without mutex exposes DistKeyGenerator to concurrent mutation; replay with mutex risks timeout under contention
- **LOW (open): Non-atomic saveSession** — os.WriteFile not atomic; partial write on crash → json.Unmarshal failure → silent deal loss
- **Recommended alternative**: in-memory prepend to existing deals/responses queues (already have FlushAllQueues for round-transition safety); avoids all disk attack surface and stage-mismatch issues; loses crash persistence (acceptable if kernel SLA enforced)
- **Test coverage (2026-03-15)**: keeper 37.8% — handleDKGProcessDeals/Responses 0.0%, ResumeDKGService 0.0%, state_manager saveSession/loadSession 0.0%

### Final DKG Full-Flow Audit — STATUS AS OF 2026-03-10
- keeper package: 32.9% (80% threshold STILL NOT MET)
- **OPEN MEDIUM: GetPendingUpgrade does NOT filter IsActivated=true entries** — kernel_upgrade.go:68 returns first entry unconditionally. hasPendingUpgradeActivation filters at call site, but UpgradeScheduled duplicate-check uses GetPendingUpgrade and will false-positive reject new upgrades while resharing round is in progress (activated entry still in store).
- **OPEN MEDIUM: nil session MarkFailed panic** — handleDKGDealing:53, handleDKGFinalization:44, handleDKGComplete:27 all call MarkFailed(ctx, nil) when GetSession fails. session=nil → UpdatePhase() panic. Fix: remove MarkFailed call when session is nil.
- **OPEN MEDIUM: ResumeDKGService called every block from proposal_server.go** — PhaseFailed state spawns new goroutine every block; goroutine terminates via dkgSvcRunning CAS but creation overhead accumulates during TEE downtime.
- **OPEN MEDIUM: old-only member session stuck at PhaseInitializing** — handleDKGRegistration returns early without updating phase for old-only members; subsequent deals/responses processing requires PhaseDealing which is never reached.
- **OPEN LOW: nil VssResponse dereference in handleDKGProcessResponses** — dkg_svc_dealing.go:244 resp.VssResponse.Index without nil check.
- **OPEN LOW: ThresholdDecryptRequested has no size limits** — ciphertext and requesterPubKey have no max size validation; unbounded DecryptRequests append.
- **OPEN LOW: ResumeDKGService switch missing DKGStageEnded/DKGStageFailed cases** — implicit fall-through is safe but violates defensive coding.
- **OPEN INFO: parseAndVerifyVoteExtension returns []*types.Vote but always 0 or 1 element** — nolint:unparam needed; simplify return type.
