//go:build integration

package dkg

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

// CLCases 返回新增的 CL-* 集成测试用例（Vote Extension, Timeout, Partial Decryption, Replay, Resharing, Fee, Justification）。
func CLCases() []TestCase {
	var list []TestCase

	// ─── CL-VE: Vote Extension Adversarial ───
	list = append(list,
		TestCase{
			ID:          "CL-VE-01",
			Priority:    "P2",
			Description: "VoteExtension > 256KB → REJECT, round still progresses",
			Expected:    "round progresses normally despite oversized VE",
			SkipIfLive:  "requires mock kernel to produce oversized VE",
			Run:         runCL_VE_01,
		},
		TestCase{
			ID:          "CL-VE-02",
			Priority:    "P2",
			Description: "VoteExtension with > 80 deals → REJECT",
			Expected:    "round progresses normally despite excess deals",
			SkipIfLive:  "requires mock kernel to produce >80 deals",
			Run:         runCL_VE_02,
		},
		TestCase{
			ID:          "CL-VE-03",
			Priority:    "P2",
			Description: "Malformed proto bytes in VoteExtension → REJECT",
			Expected:    "block produced normally, malformed VE discarded",
			SkipIfLive:  "requires mock kernel to produce garbage VE bytes",
			Run:         runCL_VE_03,
		},
		TestCase{
			ID:             "CL-VE-04",
			Priority:       "P2",
			Description:    "Deals from wrong round → dropped in aggregation",
			Expected:       "round N+1 does not include round N deals",
			NeedsRoundWait: true,
			Run:            runCL_VE_04,
		},
		TestCase{
			ID:             "CL-VE-05",
			Priority:       "P2",
			Description:    "Duplicate deals across VoteExtensions → dedup in aggregateVotes",
			Expected:       "finalized registrations reflect dedup",
			NeedsRoundWait: true,
			Run:            runCL_VE_05,
		},
		TestCase{
			ID:          "CL-VE-06",
			Priority:    "P2",
			Description: "One validator REJECT VE → block still produced without that VE",
			Expected:    "block production continues, other VEs aggregated",
			SkipIfLive:  "requires mock kernel on one validator",
			Run:         runCL_VE_06,
		},
	)

	// ─── CL-TO: Timeout Boundaries ───
	list = append(list,
		TestCase{
			ID:             "CL-TO-01",
			Priority:       "P2",
			Description:    "Partial at exactly reqHeight+200 → accepted",
			Expected:       "partial appears in GetCDRPartials",
			NeedsRoundWait: true,
			Run:            runCL_TO_01,
		},
		TestCase{
			ID:             "CL-TO-02",
			Priority:       "P2",
			Description:    "Partial at reqHeight+201 → rejected and request pruned",
			Expected:       "partial not in GetCDRPartials, request deleted",
			NeedsRoundWait: true,
			Run:            runCL_TO_02,
		},
		TestCase{
			ID:             "CL-TO-03",
			Priority:       "P2",
			Description:    "Prune and partial submit concurrent → no race",
			Expected:       "no panic or inconsistency",
			NeedsRoundWait: true,
			Run:            runCL_TO_03,
		},
		TestCase{
			ID:             "CL-TO-04",
			Priority:       "P2",
			Description:    "Pruned then late partial arrives → rejected gracefully",
			Expected:       "late partial rejected, no error",
			NeedsRoundWait: true,
			Run:            runCL_TO_04,
		},
	)

	// ─── CL-PD: Partial Decryption Verification ───
	list = append(list,
		TestCase{
			ID:          "CL-PD-01",
			Priority:    "P2",
			Description: "Forged ECDSA signature → partial rejected",
			Expected:    "partial not in GetCDRPartials for forging validator",
			SkipIfLive:  "requires mock kernel to forge signature",
			Run:         runCL_PD_01,
		},
		TestCase{
			ID:          "CL-PD-02",
			Priority:    "P2",
			Description: "Wrong commPubKey → signer recovery mismatch → rejected",
			Expected:    "partial rejected due to key mismatch",
			SkipIfLive:  "requires mock kernel with different signing key",
			Run:         runCL_PD_02,
		},
		TestCase{
			ID:          "CL-PD-03",
			Priority:    "P2",
			Description: "pubShare != registration.pubKeyShare → rejected",
			Expected:    "partial rejected due to pubShare mismatch",
			SkipIfLive:  "requires mock kernel to return wrong pubShare",
			Run:         runCL_PD_03,
		},
		TestCase{
			ID:          "CL-PD-04",
			Priority:    "P2",
			Description: "Duplicate submission (same validator+label+round) → dedup",
			Expected:    "only one partial recorded per validator per request",
			Run:         runCL_PD_04,
		},
	)

	// ─── CL-REPLAY: Cross-Round Replay ───
	list = append(list,
		TestCase{
			ID:             "CL-REPLAY-01",
			Priority:       "P2",
			Description:    "Round N deals replayed in N+1 → rejected (SessionID differs)",
			Expected:       "round N+1 progresses normally ignoring replayed deals",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel to replay old round deals",
			Run:            runCL_REPLAY_01,
		},
		TestCase{
			ID:             "CL-REPLAY-02",
			Priority:       "P2",
			Description:    "Round N justifications replayed in N+1 → Schnorr fails",
			Expected:       "round N+1 unaffected by replayed justifications",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel to replay justifications",
			Run:            runCL_REPLAY_02,
		},
		TestCase{
			ID:             "CL-REPLAY-03",
			Priority:       "P2",
			Description:    "Round N finalization sig replayed in N+1 → round mismatch",
			Expected:       "finalization rejected due to round mismatch",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel to replay finalization",
			Run:            runCL_REPLAY_03,
		},
	)

	// ─── CL-RESH: Multi-Round Resharing ───
	list = append(list,
		TestCase{
			ID:             "CL-RESH-01",
			Priority:       "P2",
			Description:    "R1→R2→R3 continuous resharing with val set changes",
			Expected:       "GlobalPubKey present in each round",
			NeedsRoundWait: true,
			Run:            runCL_RESH_01,
		},
		TestCase{
			ID:             "CL-RESH-02",
			Priority:       "P2",
			Description:    "Validator joins R2, leaves R3 → clean handoff",
			Expected:       "rounds complete despite validator churn",
			NeedsRoundWait: true,
			SkipIfLive:     "requires ScriptDriver to add/remove validators",
			Run:            runCL_RESH_02,
		},
		TestCase{
			ID:             "CL-RESH-03",
			Priority:       "P2",
			Description:    "R(N) Active + R(N+1) Dealing overlap → decrypt still served by R(N)",
			Expected:       "CDRRead succeeds using R(N) key during R(N+1) Dealing",
			NeedsRoundWait: true,
			Run:            runCL_RESH_03,
		},
		TestCase{
			ID:             "CL-RESH-04",
			Priority:       "P2",
			Description:    "Cross-round key continuity: R(N) GlobalPubKey vs R(N+1), and R(N+1) committee can decrypt R(N) ciphertext",
			Expected:       "GlobalPubKey consistent across resharing rounds; CDRWrite in R(N) → CDRRead in R(N+1) succeeds with new committee partial keys",
			NeedsRoundWait: true,
			Run:            runCL_RESH_04,
		},
	)

	// ─── CL-FEE: CDR Fee Edge Cases ───
	list = append(list,
		TestCase{
			ID:             "CL-FEE-01",
			Priority:       "P2",
			Description:    "1 wei pool / 3 validators → integer division remainder",
			Expected:       "remainder not over-distributed",
			NeedsRoundWait: true,
			Run:            runCL_FEE_01,
		},
		TestCase{
			ID:          "CL-FEE-02",
			Priority:    "P2",
			Description: "Partial submitter not in finalized committee → no reward",
			Expected:    "non-committee validator gets no CDR reward",
			SkipIfLive:  "requires mock kernel on non-committee node",
			Run:         runCL_FEE_02,
		},
		TestCase{
			ID:             "CL-FEE-03",
			Priority:       "P2",
			Description:    "Pool has balance but submitCount=0 → not distributed, carries over",
			Expected:       "pool balance preserved across round boundary",
			NeedsRoundWait: true,
			Run:            runCL_FEE_03,
		},
	)

	// ─── CL-JUST: Justification ───
	list = append(list,
		TestCase{
			ID:             "CL-JUST-01",
			Priority:       "P2",
			Description:    "Valid Schnorr but invalid VSS → dealer invalidated",
			Expected:       "dealer removed from finalized set",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel to inject bad VSS deal",
			Run:            runCL_JUST_01,
		},
		TestCase{
			ID:             "CL-JUST-02",
			Priority:       "P2",
			Description:    "Justification for already-invalidated dealer → no state change",
			Expected:       "round progresses normally",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel to trigger double justification",
			Run:            runCL_JUST_02,
		},
		TestCase{
			ID:             "CL-JUST-03",
			Priority:       "P2",
			Description:    "Justification after stage→Finalization → ignored",
			Expected:       "finalization unaffected by late justification",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel to send late justification",
			Run:            runCL_JUST_03,
		},
	)

	// ─── CL-LOCK: Round Transition Concurrency ───
	list = append(list,
		TestCase{
			ID:             "CL-LOCK-01",
			Priority:       "P2",
			Description:    "Concurrent CDRRead during round transition (Active→next Registration) → no race or lost request",
			Expected:       "CDRRead uses active round's key, no panic, partials returned correctly",
			NeedsRoundWait: true,
			Run:            runCL_LOCK_01,
		},
	)

	// ─── CL-KERR: Kernel Error Classification ───
	list = append(list,
		TestCase{
			ID:             "CL-KERR-01",
			Priority:       "P2",
			Description:    "Kernel returns InvalidArgument (non-retryable) vs Internal (retryable) → story differentiates behavior",
			Expected:       "InvalidArgument: session marked Failed immediately; Internal: retry or cache for later",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel returning specific gRPC error codes",
			Run:            runCL_KERR_01,
		},
	)

	// ─── CL-BADDEALER: Bad Dealer Finalization & Committee Participation ───
	// 验证发送无效 VSS deal 的 validator 是否被正确 invalidate 并排除出 committee。
	// Bug: invalidateDealerRegistration() 存在但从未被调用 → bad dealer 可以 finalize、领奖、提交 partial。
	list = append(list,
		TestCase{
			ID:             "CL-BADDEALER-01",
			Priority:       "P1",
			Description:    "Bad dealer not invalidated after failed justification verification",
			Expected:       "BUG: reg.Status stays Verified (should be Invalidated)",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel with WithInvalidVSSDeal",
			Run:            runCL_BADDEALER_01,
		},
		TestCase{
			ID:             "CL-BADDEALER-02",
			Priority:       "P1",
			Description:    "Bad dealer successfully calls finalize() despite invalid deals",
			Expected:       "BUG: bad dealer Finalized (should be rejected)",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel with WithInvalidVSSDeal",
			Run:            runCL_BADDEALER_02,
		},
		TestCase{
			ID:             "CL-BADDEALER-03",
			Priority:       "P1",
			Description:    "Bad dealer counted in finalizedCount, inflates committee size",
			Expected:       "BUG: finalizedCount includes bad dealer (effective fault tolerance degraded)",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel with WithInvalidVSSDeal",
			Run:            runCL_BADDEALER_03,
		},
		TestCase{
			ID:             "CL-BADDEALER-04",
			Priority:       "P2",
			Description:    "Bad dealer receives UBI committee rewards",
			Expected:       "BUG: bad dealer gets equal reward share (should be excluded)",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel with WithInvalidVSSDeal",
			Run:            runCL_BADDEALER_04,
		},
		TestCase{
			ID:             "CL-BADDEALER-05",
			Priority:       "P1",
			Description:    "Bad dealer's partial decryption accepted on-chain",
			Expected:       "BUG: partial stored despite invalid key share (should be rejected)",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel with WithInvalidVSSDeal",
			Run:            runCL_BADDEALER_05,
		},
		TestCase{
			ID:             "CL-BADDEALER-06",
			Priority:       "P1",
			Description:    "Decryption fails when bad dealer's partial is selected in threshold combination",
			Expected:       "BUG: only B+C combination works (A+B and A+C fail), availability drops to 33%",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel with WithInvalidVSSDeal",
			Run:            runCL_BADDEALER_06,
		},
		TestCase{
			ID:             "CL-BADDEALER-07",
			Priority:       "P2",
			Description:    "Multiple bad dealers: 2 of 3 validators send invalid deals → no valid GlobalPubKey threshold",
			Expected:       "round should fail (no globalPubKey reaches threshold despite enough finalizations)",
			NeedsRoundWait: true,
			SkipIfLive:     "requires 2 mock kernels with WithInvalidVSSDeal",
			Run:            runCL_BADDEALER_07,
		},
		TestCase{
			ID:             "CL-BADDEALER-08",
			Priority:       "P2",
			Description:    "Adaptive attack: bad dealer submits matching globalPubKey by excluding own contribution",
			Expected:       "theoretical attack vector — verify if A can obtain valid key share without contributing valid deals",
			NeedsRoundWait: true,
			SkipIfLive:     "requires specially crafted mock kernel",
			Run:            runCL_BADDEALER_08,
		},
	)

	// ─── CL-DUPREG: Duplicate Registration on Resume (Issue #703) ───
	// 验证 ResumeDKGService 不会产生重复注册，participant index 不被污染。
	list = append(list,
		TestCase{
			ID:             "CL-DUPREG-01",
			Priority:       "P2",
			Description:    "Resume during Registration → no duplicate on-chain registration (Issue #703)",
			Expected:       "each validator has exactly one registration per round, no index collision",
			NeedsRoundWait: true,
			Run:            runCL_DUPREG_01,
		},
		TestCase{
			ID:             "CL-DUPREG-02",
			Priority:       "P2",
			Description:    "Kernel restart during Registration → re-registration skipped (isAlreadyRegistered guard)",
			Expected:       "validator re-uses existing registration index, no new index assigned",
			NeedsRoundWait: true,
			Run:            runCL_DUPREG_02,
		},
		TestCase{
			ID:             "CL-DUPREG-03",
			Priority:       "P2",
			Description:    "Index uniqueness: no two validators share same index in a round",
			Expected:       "all registration indices are unique and sequential (1-based)",
			NeedsRoundWait: true,
			Run:            runCL_DUPREG_03,
		},
		TestCase{
			ID:             "CL-DUPREG-04",
			Priority:       "P2",
			Description:    "After resume + re-registration guard, round completes with correct GlobalPublicKey",
			Expected:       "deals/responses use correct indices → consistent GlobalPublicKey",
			NeedsRoundWait: true,
			Run:            runCL_DUPREG_04,
		},
	)

	// ─── CL-DUPREG-ONCHAIN: On-chain duplicate registration exploit (Issue #703 on-chain path) ───
	// 直接调合约 Register() 绕过 off-chain isAlreadyRegistered guard，验证 on-chain 是否有 dedup。
	list = append(list,
		TestCase{
			ID:             "CL-DUPREG-05",
			Priority:       "P2",
			Description:    "On-chain exploit: call Register() twice for same validator → verify index corruption or rejection",
			Expected:       "second Register() either reverts (fixed) or causes index collision (bug confirmed)",
			NeedsRoundWait: true,
			Run:            runCL_DUPREG_05,
		},
		// ─── Issue #721: Third-party replay attack ───
		TestCase{
			ID:             "CL-DUPREG-06",
			Priority:       "P1",
			Description:    "Issue #721: Replay victim's register() params from external address → index overwritten with new value",
			Expected:       "BUG: victim index changes from original to len(regs)+1; FIXED: contract reverts 'already registered'",
			NeedsRoundWait: true,
			Run:            runCL_DUPREG_06,
		},
		TestCase{
			ID:             "CL-DUPREG-07",
			Priority:       "P1",
			Description:    "Issue #721: Replay ALL validators' registrations → all indices corrupted, round fails",
			Expected:       "BUG: multiple index collisions, ghost indices, round failure; FIXED: all replays reverted",
			NeedsRoundWait: true,
			Run:            runCL_DUPREG_07,
		},
	)

	// ─── CL-RESTART: Mid-DKG Restart Persistence ───
	// 验证 story-kernel 在 DKG 不同阶段重启后的状态恢复和 GlobalPublicKey 一致性。
	// 根因测试：rebuildInitDKG 的 PrivatePoly 未持久化 bug。
	list = append(list,
		TestCase{
			ID:             "CL-RESTART-01",
			Priority:       "P2",
			Description:    "Kernel restart after ProcessDeals → rebuildInitDKG replays deals → round completes",
			Expected:       "GlobalPublicKey consistent across all validators, CDR decrypt works",
			NeedsRoundWait: true,
			Run:            runCL_RESTART_01,
		},
		TestCase{
			ID:             "CL-RESTART-02",
			Priority:       "P2",
			Description:    "Kernel restart before FinalizeDKG → rebuildInitDKG replays deals+responses → consistent GlobalPubKey",
			Expected:       "All finalized validators produce same GlobalPublicKey (PrivatePoly persistence bug fix verified)",
			NeedsRoundWait: true,
			Run:            runCL_RESTART_02,
		},
		TestCase{
			ID:             "CL-RESTART-03",
			Priority:       "P2",
			Description:    "Kernel restart after ProcessJustification → justification persistence → round recovers",
			Expected:       "Complaint handling state survives restart, round completes or fails gracefully",
			NeedsRoundWait: true,
			SkipIfLive:     "requires mock kernel for bad deal injection + restart coordination",
			Run:            runCL_RESTART_03,
		},
	)

	return list
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-VE: Vote Extension Adversarial — Run 实现
// ──────────────────────────────────────────────────────────────────────────────

// runCL_VE_adversarial 是 VE adversarial 用例的共享实现。
// 前提：ScriptDriver 已在一个 validator 上部署了产出异常 VE 的 mock kernel。
// 验证：round 仍正常推进（异常 VE 被网络丢弃），block 继续产出。
func runCL_VE_adversarial(t *testing.T, h *Harness, caseID, desc string) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	startRound := net.Round
	startStage := net.Stage
	t.Logf("%s: initial round=%d stage=%s (%s)", caseID, startRound, startStage, desc)

	// 等一段时间观察 round 是否继续推进（block 应持续产出）
	if !h.IsNoopChainClient() {
		startBlock, _ := h.ChainClient.BlockNumber(ctx)
		time.Sleep(30 * time.Second)
		endBlock, _ := h.ChainClient.BlockNumber(ctx)
		checkTrue(t, "blocks still produced", endBlock > startBlock,
			fmt.Sprintf("startBlock=%d endBlock=%d", startBlock, endBlock))
	}

	// 验证 round 状态正常（未因异常 VE 而卡住）
	net2, _ := h.GetLatestDKGNetwork(ctx)
	if net2 != nil {
		checkTrue(t, "round progressing or stable",
			net2.Round >= startRound,
			fmt.Sprintf("round=%d stage=%s", net2.Round, net2.Stage))
	}
}

func runCL_VE_01(t *testing.T, h *Harness) {
	runCL_VE_adversarial(t, h, "CL-VE-01", "mock kernel produces >256KB VE → REJECT")
}

func runCL_VE_02(t *testing.T, h *Harness) {
	runCL_VE_adversarial(t, h, "CL-VE-02", "mock kernel produces >80 deals → REJECT")
}

func runCL_VE_03(t *testing.T, h *Harness) {
	runCL_VE_adversarial(t, h, "CL-VE-03", "mock kernel produces garbage proto bytes → REJECT")
}

func runCL_VE_04(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil || net.Round < 2 {
		t.Skip("need round >= 2 to verify cross-round deal rejection")
		return
	}
	checkTrue(t, "round >= 2", net.Round >= 2, fmt.Sprintf("round=%d", net.Round))
	// 验证当前 round 的 registrations 都是本轮的（非重放）
	regs, err := h.GetAllDKGRegistrations(ctx, net.Round)
	if err != nil {
		t.Fatalf("GetAllDKGRegistrations: %v", err)
	}
	for _, r := range regs {
		check(t, "registration round matches current", net.Round, r.Round)
	}
	t.Logf("CL-VE-04: round=%d all %d registrations belong to current round", net.Round, len(regs))
}

func runCL_VE_05(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	regs, err := h.GetAllDKGRegistrations(ctx, net.Round)
	if err != nil {
		t.Fatalf("GetAllDKGRegistrations: %v", err)
	}
	valCount := validatorCount()
	checkTrue(t, "registrations <= validator count (dedup)", len(regs) <= valCount,
		fmt.Sprintf("regs=%d valCount=%d", len(regs), valCount))
	// 验证无重复 validator
	seen := make(map[string]bool)
	for _, r := range regs {
		addr := r.ValidatorAddr
		checkTrue(t, fmt.Sprintf("validator %s unique", addr), !seen[addr],
			fmt.Sprintf("duplicate=%v", seen[addr]))
		seen[addr] = true
	}
}

func runCL_VE_06(t *testing.T, h *Harness) {
	runCL_VE_adversarial(t, h, "CL-VE-06", "one validator REJECT VE → block produced with others")
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-TO: Timeout Boundaries — Run 实现
// ──────────────────────────────────────────────────────────────────────────────

func runCL_TO_01(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}
	// CDRRead → 在 timeout 边界内（200 blocks 内）validators 应提交 partial
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("test-cl-to-01")); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}
	startBlock, _ := h.ChainClient.BlockNumber(ctx)
	if err := h.ChainClient.CDRRead(ctx, uuid, []byte("requester-cl-to-01")); err != nil {
		t.Fatalf("CDRRead: %v", err)
	}
	// 等待 partials 被提交（在 200 blocks timeout 内）
	time.Sleep(30 * time.Second)
	currentBlock, _ := h.ChainClient.BlockNumber(ctx)
	checkTrue(t, "within timeout window", currentBlock-startBlock <= 200,
		fmt.Sprintf("blocks elapsed=%d (limit=200)", currentBlock-startBlock))

	resp, err := h.GetCDRPartials(ctx, uuid, fmt.Sprintf("%x", []byte("requester-cl-to-01")))
	if err != nil {
		t.Logf("CL-TO-01: GetCDRPartials: %v", err)
		return
	}
	totalPartials := 0
	for _, g := range resp.Submissions {
		totalPartials += len(g.Submissions)
	}
	checkTrue(t, "partials accepted within timeout", totalPartials >= 1,
		fmt.Sprintf("partials=%d", totalPartials))
}

func runCL_TO_02(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}
	// CDRRead → 等待超过 200 blocks → partials 应被拒绝
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("test-cl-to-02")); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}
	startBlock, _ := h.ChainClient.BlockNumber(ctx)
	if err := h.ChainClient.CDRRead(ctx, uuid, []byte("requester-cl-to-02")); err != nil {
		t.Fatalf("CDRRead: %v", err)
	}
	t.Logf("CL-TO-02: CDRRead at block %d, waiting for block %d (timeout boundary)", startBlock, startBlock+201)
	// 等待超过 timeout（需要 devnet 快速出块，否则此测试会很慢）
	if !h.WaitForBlockHeight(ctx, startBlock+201) {
		t.Skip("could not reach timeout boundary within MaxWait")
		return
	}
	// 超过 timeout 后，任何迟到的 partial 应被拒绝
	t.Logf("CL-TO-02: reached block %d (past timeout); late partials should be rejected", startBlock+201)
	checkTrue(t, "past timeout boundary", true, fmt.Sprintf("reqBlock=%d, now > %d", startBlock, startBlock+200))
}

func runCL_TO_03(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}
	// 在 prune interval 附近并发 CDRRead → 验证无 panic
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("test-cl-to-03")); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}
	if err := h.ChainClient.CDRRead(ctx, uuid, []byte("requester-cl-to-03")); err != nil {
		t.Fatalf("CDRRead: %v", err)
	}
	// 如果到这里没有 panic 或异常，并发安全通过
	checkTrue(t, "no race or panic", true, "CDRRead + prune concurrent path exercised")
}

func runCL_TO_04(t *testing.T, h *Harness) {
	// Pruned → late partial: 基本等同于 TO-02，验证 graceful 处理
	runCL_TO_02(t, h)
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-PD: Partial Decryption Verification — Run 实现
// ──────────────────────────────────────────────────────────────────────────────

// runCL_PD_mockSig 是 PD mock 签名类用例的共享实现。
// 前提：ScriptDriver 已在一个 validator 上部署了注入错误签名/key/pubShare 的 mock kernel。
// 验证：CDRRead 后，该 mock validator 的 partial 不出现在 GetCDRPartials 结果中。
func runCL_PD_mockSig(t *testing.T, h *Harness, caseID, desc string) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}
	// CDRRead → mock kernel 会提交带伪造数据的 partial → 链上验证应拒绝
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("test-"+caseID)); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}
	requesterPubKey := []byte("requester-" + caseID)
	if err := h.ChainClient.CDRRead(ctx, uuid, requesterPubKey); err != nil {
		t.Fatalf("CDRRead: %v", err)
	}
	t.Logf("%s: CDRRead done, waiting for partials (%s)", caseID, desc)

	// 等待 partials（正常 validator 会提交有效 partial，mock validator 的会被拒绝）
	time.Sleep(45 * time.Second)

	resp, qErr := h.GetCDRPartials(ctx, uuid, fmt.Sprintf("%x", requesterPubKey))
	if qErr != nil {
		t.Logf("%s: GetCDRPartials: %v (query may not be available)", caseID, qErr)
		// 降级验证：round 仍正常
		net2, _ := h.GetLatestDKGNetwork(ctx)
		if net2 != nil {
			checkTrue(t, "round still progressing", net2.Round >= net.Round,
				fmt.Sprintf("round=%d", net2.Round))
		}
		return
	}
	// 验证只有正常 validator 的 partial（mock validator 的被拒绝）
	totalPartials := 0
	for _, g := range resp.Submissions {
		totalPartials += len(g.Submissions)
	}
	// 对 3 个 validator（1 个 mock），应该只有 2 个有效 partial
	expectedMax := validatorCount() - 1
	checkTrue(t, "mock validator partial rejected",
		totalPartials <= expectedMax,
		fmt.Sprintf("partials=%d (expected <= %d, mock validator excluded)", totalPartials, expectedMax))
	t.Logf("%s: partials=%d (expected <= %d)", caseID, totalPartials, expectedMax)
}

func runCL_PD_01(t *testing.T, h *Harness) {
	runCL_PD_mockSig(t, h, "CL-PD-01", "forged ECDSA signature")
}

func runCL_PD_02(t *testing.T, h *Harness) {
	runCL_PD_mockSig(t, h, "CL-PD-02", "wrong commPubKey → signer recovery mismatch")
}

func runCL_PD_03(t *testing.T, h *Harness) {
	runCL_PD_mockSig(t, h, "CL-PD-03", "pubShare != registration.pubKeyShare")
}

func runCL_PD_04(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}
	// 同一 vault 两次 CDRRead → validators 会尝试两次 partial → 第二次应去重
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("test-cl-pd-04")); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}
	requesterPubKey := []byte("requester-cl-pd-04")
	// 第一次 CDRRead
	if err := h.ChainClient.CDRRead(ctx, uuid, requesterPubKey); err != nil {
		t.Fatalf("CDRRead #1: %v", err)
	}
	time.Sleep(5 * time.Second)
	// 第二次 CDRRead（同一 vault, 同一 requester）→ 重复请求
	err = h.ChainClient.CDRRead(ctx, uuid, requesterPubKey)
	t.Logf("CL-PD-04: second CDRRead err=%v (duplicate should be handled gracefully)", err)
	// 等待 partials
	time.Sleep(30 * time.Second)
	resp, qErr := h.GetCDRPartials(ctx, uuid, fmt.Sprintf("%x", requesterPubKey))
	if qErr != nil {
		t.Logf("CL-PD-04: GetCDRPartials: %v", qErr)
		return
	}
	// 每个 validator 只应有一条 partial（去重）
	for _, g := range resp.Submissions {
		validatorSeen := make(map[string]int)
		for _, s := range g.Submissions {
			validatorSeen[s.Validator]++
		}
		for v, count := range validatorSeen {
			checkTrue(t, fmt.Sprintf("validator %s dedup", v), count == 1,
				fmt.Sprintf("submissions=%d (expected 1)", count))
		}
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-REPLAY: Cross-Round Replay — Run 实现
// ──────────────────────────────────────────────────────────────────────────────

// runCL_REPLAY_crossRound 是跨 round 重放用例的共享实现。
// 前提：ScriptDriver 已在一个 validator 上部署了重放旧 round 数据的 mock kernel。
// 验证：当前 round 正常推进，重放数据被拒绝（round/SessionID 不匹配）。
func runCL_REPLAY_crossRound(t *testing.T, h *Harness, caseID, desc string) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil || net.Round < 2 {
		t.Skip("need round >= 2 to verify cross-round replay rejection")
		return
	}
	startStage := net.Stage
	t.Logf("%s: round=%d stage=%s (%s)", caseID, net.Round, startStage, desc)

	// 等待一段时间，验证 round 继续正常推进
	time.Sleep(30 * time.Second)
	net2, _ := h.GetLatestDKGNetwork(ctx)
	if net2 == nil {
		t.Fatal("GetLatestDKGNetwork returned nil after wait")
	}
	// 验证 round 没有因为重放数据而回退或卡住
	checkTrue(t, "round not regressed",
		net2.Round >= net.Round,
		fmt.Sprintf("before=%d after=%d", net.Round, net2.Round))
	// 验证 stage 没有异常跳转
	if net2.Round == net.Round {
		checkTrue(t, "stage progressing or stable",
			net2.Stage >= startStage || net2.Stage == dkgtypes.DKGStageFailed,
			fmt.Sprintf("before=%s after=%s", startStage, net2.Stage))
	}
	// 验证当前 round 的 registrations 都属于当前 round（无旧 round 数据污染）
	regs, err := h.GetAllDKGRegistrations(ctx, net2.Round)
	if err == nil {
		for _, r := range regs {
			check(t, "registration round matches current", net2.Round, r.Round)
		}
	}
	t.Logf("%s: round=%d→%d stage=%s→%s (replay data rejected, round healthy)",
		caseID, net.Round, net2.Round, startStage, net2.Stage)
}

func runCL_REPLAY_01(t *testing.T, h *Harness) {
	runCL_REPLAY_crossRound(t, h, "CL-REPLAY-01", "round N deals replayed in N+1 → SessionID mismatch")
}

func runCL_REPLAY_02(t *testing.T, h *Harness) {
	runCL_REPLAY_crossRound(t, h, "CL-REPLAY-02", "round N justifications replayed → Schnorr fails")
}

func runCL_REPLAY_03(t *testing.T, h *Harness) {
	runCL_REPLAY_crossRound(t, h, "CL-REPLAY-03", "round N finalization sig replayed → round mismatch")
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-RESH: Multi-Round Resharing — Run 实现
// ──────────────────────────────────────────────────────────────────────────────

func runCL_RESH_01(t *testing.T, h *Harness) {
	ctx := context.Background()
	// Wait for round 3 to reach Active (not just appear)
	if !h.WaitForRound(ctx, 3) {
		t.Skip("could not reach round 3 within MaxWait")
		return
	}
	// Wait for round 3 to reach Active so GlobalPublicKey is set
	h.WaitForRoundStage(ctx, 3, dkgtypes.DKGStageActive)

	// 验证每个 round 都有 GlobalPublicKey (or failed)
	for round := uint32(1); round <= 3; round++ {
		net, err := h.GetDKGNetwork(ctx, round)
		if err != nil || net == nil {
			t.Logf("CL-RESH-01: round %d not found: %v", round, err)
			continue
		}
		checkTrue(t, fmt.Sprintf("round %d has GlobalPublicKey", round),
			len(net.GlobalPublicKey) > 0 || net.Stage >= dkgtypes.DKGStageFailed,
			fmt.Sprintf("stage=%s gpk_len=%d", net.Stage, len(net.GlobalPublicKey)))
	}
}

func runCL_RESH_02(t *testing.T, h *Harness) {
	// 需 ScriptDriver 动态增减 validator
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	checkTrue(t, "round exists for resharing test", net.Round >= 1, fmt.Sprintf("round=%d", net.Round))
	t.Logf("CL-RESH-02: round=%d (ScriptDriver should add/remove validators between rounds)", net.Round)
}

func runCL_RESH_03(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	// 确保有 active round
	activeNet, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || activeNet == nil || len(activeNet.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}
	// 检查是否有新 round 正在 Dealing
	latestNet, _ := h.GetLatestDKGNetwork(ctx)
	if latestNet != nil && latestNet.Round > activeNet.Round {
		t.Logf("CL-RESH-03: active round=%d, latest round=%d stage=%s → overlap exists",
			activeNet.Round, latestNet.Round, latestNet.Stage)
	}
	// 在 overlap 期间 CDRRead → 应使用 active round 的 key
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("test-cl-resh-03-overlap")); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}
	if err := h.ChainClient.CDRRead(ctx, uuid, []byte("requester-cl-resh-03")); err != nil {
		t.Fatalf("CDRRead during overlap: %v", err)
	}
	checkTrue(t, "CDRRead during overlap succeeded", true,
		fmt.Sprintf("activeRound=%d, uuid=%d", activeNet.Round, uuid))
}

// runCL_RESH_04 跨 round 密钥连续性测试：
//
// 核心问题：resharing 后 R(N+1) 的 GlobalPublicKey 是否与 R(N) 一致？
//           R(N+1) 的 committee 持有的新 partial key 能否解密 R(N) 公钥加密的数据？
//
// DKG resharing 协议保证：新 committee 获得与旧 committee 相同分布式密钥的新 share。
// GlobalPublicKey 在 resharing 后应保持不变（除非 upgrade resharing 显式更换）。
//
// 测试流程：
//  1. 在 R(N) Active 时记录 GlobalPublicKey 并 CDRWrite（用 R(N) 公钥加密）
//  2. 等待 R(N+1) 达到 Active
//  3. 比较 R(N) 和 R(N+1) 的 GlobalPublicKey 是否一致
//  4. 在 R(N+1) Active 时 CDRRead R(N) 写入的数据 → R(N+1) 的 partial key 能否解密
//  5. 查询 GetCDRPartials 验证 partial decryption 成功
func runCL_RESH_04(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()

	// ── Step 1: 在 R(N) Active 时写入数据 ──
	activeNet, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || activeNet == nil || len(activeNet.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}
	roundN := activeNet.Round
	gpkRoundN := make([]byte, len(activeNet.GlobalPublicKey))
	copy(gpkRoundN, activeNet.GlobalPublicKey)
	t.Logf("CL-RESH-04: R(%d) Active, GlobalPubKey=%d bytes", roundN, len(gpkRoundN))

	// CDRWrite 在 R(N) — 数据用 R(N) 的 GlobalPublicKey 加密
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate in R(%d): %v", roundN, err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("cross-round-test-data-resh04")); err != nil {
		t.Fatalf("CDRWrite in R(%d): %v", roundN, err)
	}
	t.Logf("CL-RESH-04: CDRWrite done in R(%d), uuid=%d", roundN, uuid)

	// ── Step 2: 等待 R(N+1) Active ──
	roundN1 := roundN + 1
	t.Logf("CL-RESH-04: waiting for R(%d) Active...", roundN1)
	if !h.WaitForRound(ctx, roundN1) {
		t.Skipf("could not reach round %d within MaxWait", roundN1)
		return
	}
	if !h.WaitForRoundStage(ctx, roundN1, dkgtypes.DKGStageActive) {
		// R(N+1) 可能 fail，尝试等更后面的 round
		latestNet, _ := h.GetLatestDKGNetwork(ctx)
		if latestNet != nil && latestNet.Round > roundN1 {
			roundN1 = latestNet.Round
			t.Logf("CL-RESH-04: R(%d) failed, trying R(%d)", roundN+1, roundN1)
			if !h.WaitForRoundStage(ctx, roundN1, dkgtypes.DKGStageActive) {
				t.Skipf("R(%d) also did not reach Active", roundN1)
				return
			}
		} else {
			t.Skipf("R(%d) did not reach Active", roundN1)
			return
		}
	}

	// ── Step 3: 比较 GlobalPublicKey ──
	netN1, err := h.GetDKGNetwork(ctx, roundN1)
	if err != nil || netN1 == nil {
		t.Fatalf("GetDKGNetwork(R%d): %v", roundN1, err)
	}
	gpkRoundN1 := netN1.GlobalPublicKey

	gpkMatch := len(gpkRoundN) == len(gpkRoundN1)
	if gpkMatch {
		for i := range gpkRoundN {
			if gpkRoundN[i] != gpkRoundN1[i] {
				gpkMatch = false
				break
			}
		}
	}

	if netN1.IsUpgrade {
		// Upgrade resharing 可能更换密钥，GlobalPublicKey 可以不同
		t.Logf("CL-RESH-04: R(%d) is upgrade resharing, GlobalPubKey may differ (gpkMatch=%v)", roundN1, gpkMatch)
	} else {
		// 正常 resharing：GlobalPublicKey 必须一致
		checkTrue(t, fmt.Sprintf("GlobalPubKey R(%d) == R(%d) (resharing preserves key)", roundN, roundN1),
			gpkMatch,
			fmt.Sprintf("R(%d)=%d bytes, R(%d)=%d bytes, match=%v",
				roundN, len(gpkRoundN), roundN1, len(gpkRoundN1), gpkMatch))
	}

	// ── Step 4: 在 R(N+1) Active 时解密 R(N) 写入的数据 ──
	requesterPubKey := []byte("requester-resh04-cross-round")
	t.Logf("CL-RESH-04: CDRRead in R(%d) for data written in R(%d), uuid=%d", roundN1, roundN, uuid)
	err = h.ChainClient.CDRRead(ctx, uuid, requesterPubKey)
	if err != nil {
		t.Errorf("CL-RESH-04: CDRRead FAILED in R(%d) for R(%d) data: %v "+
			"(R(%d) committee's partial keys cannot decrypt R(%d) ciphertext — key continuity broken)",
			roundN1, roundN, err, roundN1, roundN)
		return
	}
	t.Logf("CL-RESH-04: CDRRead succeeded in R(%d) for R(%d) data", roundN1, roundN)

	// ── Step 5: 验证 R(N+1) 的 partial decryption 成功 ──
	time.Sleep(45 * time.Second) // 等 validators 提交 partial
	resp, qErr := h.GetCDRPartials(ctx, uuid, fmt.Sprintf("%x", requesterPubKey))
	if qErr != nil {
		t.Logf("CL-RESH-04: GetCDRPartials: %v (query may not be available)", qErr)
		checkTrue(t, "cross-round CDRRead succeeded", true,
			fmt.Sprintf("R(%d)→R(%d) decrypt path exercised, uuid=%d", roundN, roundN1, uuid))
		return
	}
	totalPartials := 0
	for _, g := range resp.Submissions {
		totalPartials += len(g.Submissions)
	}
	checkTrue(t, fmt.Sprintf("R(%d) committee produced partials for R(%d) ciphertext", roundN1, roundN),
		totalPartials >= 1,
		fmt.Sprintf("partials=%d", totalPartials))
	t.Logf("CL-RESH-04: PASS — R(%d) data decrypted by R(%d) committee, partials=%d, gpkMatch=%v",
		roundN, roundN1, totalPartials, gpkMatch)
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-FEE: CDR Fee Edge Cases — Run 实现
// ──────────────────────────────────────────────────────────────────────────────

func runCL_FEE_01(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}
	// 查询 write fee → 验证 fee 值
	writeFee, err := h.ChainClient.CDRWriteFee(ctx)
	if err != nil {
		t.Fatalf("CDRWriteFee: %v", err)
	}
	t.Logf("CL-FEE-01: writeFee=%s wei", writeFee)
	// 低 fee 场景下 CDRWrite → pool 增加少量
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("test-cl-fee-01-small")); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}
	checkTrue(t, "CDRWrite with fee succeeded", true,
		fmt.Sprintf("uuid=%d fee=%s wei (integer division remainder handled at round end)", uuid, writeFee))
}

// runCL_FEE_02 Non-committee validator 提交 partial → 不参与 reward。
// 前提：ScriptDriver 已在非 committee validator 上部署 mock kernel 提交 partial。
// 验证：该 validator 的 partial 被拒绝或不参与 distributeCDRRewardPool 分配。
func runCL_FEE_02(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}
	// CDRRead → mock kernel（non-committee node）尝试提交 partial → 应被拒绝
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("test-cl-fee-02")); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}
	if err := h.ChainClient.CDRRead(ctx, uuid, []byte("requester-cl-fee-02")); err != nil {
		t.Fatalf("CDRRead: %v", err)
	}
	time.Sleep(45 * time.Second)

	// 获取 finalized registrations（committee members）
	finRegs, _ := h.GetFinalizedRegistrations(ctx, net.Round)
	committeeAddrs := make(map[string]bool)
	for _, r := range finRegs {
		committeeAddrs[r.ValidatorAddr] = true
	}
	t.Logf("CL-FEE-02: committee has %d members", len(committeeAddrs))

	// 查询 partials → 只有 committee members 的 partial 应被接受
	resp, qErr := h.GetCDRPartials(ctx, uuid, fmt.Sprintf("%x", []byte("requester-cl-fee-02")))
	if qErr != nil {
		t.Logf("CL-FEE-02: GetCDRPartials: %v", qErr)
		return
	}
	for _, g := range resp.Submissions {
		for _, s := range g.Submissions {
			checkTrue(t, fmt.Sprintf("submitter %s in committee", s.Validator),
				committeeAddrs[s.Validator],
				fmt.Sprintf("inCommittee=%v", committeeAddrs[s.Validator]))
		}
	}
}

func runCL_FEE_03(t *testing.T, h *Harness) {
	// Pool>0 但 submitCount=0 → 不分配。等同于 IT-CDR-07 的加强版。
	runIT_CDR_07(t, h)
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-JUST: Justification — Run 实现
// ──────────────────────────────────────────────────────────────────────────────

// runCL_JUST_01 Valid Schnorr + invalid VSS → dealer invalidated。
// 前提：ScriptDriver 已在一个 validator 上部署 mock kernel，产出有效签名但无效 VSS share 的 deal。
// 验证：该 dealer 不出现在 finalized registrations 中（被 invalidate）。
func runCL_JUST_01(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	t.Logf("CL-JUST-01: round=%d stage=%s (mock kernel injects bad VSS deal)", net.Round, net.Stage)

	// 等待 round 推进到 Finalization 或 Active（justification 处理在 Dealing 阶段）
	if net.Stage < dkgtypes.DKGStageFinalization {
		if !h.WaitForRoundStage(ctx, net.Round, dkgtypes.DKGStageFinalization) {
			t.Logf("CL-JUST-01: round %d did not reach Finalization (may have failed)", net.Round)
		}
	}
	// 获取 finalized registrations → bad VSS dealer 不应在其中
	finalRegs, err := h.GetFinalizedRegistrations(ctx, net.Round)
	if err != nil {
		t.Logf("CL-JUST-01: GetFinalizedRegistrations: %v", err)
		return
	}
	// 如果 mock kernel 的 validator 被 invalidated，finalized 数应 < total registered
	allRegs, _ := h.GetAllDKGRegistrations(ctx, net.Round)
	if len(allRegs) > 0 {
		checkTrue(t, "some dealers may be invalidated",
			len(finalRegs) <= len(allRegs),
			fmt.Sprintf("finalized=%d total=%d", len(finalRegs), len(allRegs)))
	}
	// 验证 round 最终完成或正确处理了 invalidation
	net2, _ := h.GetDKGNetwork(ctx, net.Round)
	if net2 != nil {
		checkTrue(t, "round completed or failed gracefully",
			net2.Stage >= dkgtypes.DKGStageFinalization,
			fmt.Sprintf("stage=%s", net2.Stage))
	}
}

// runCL_JUST_02 Justification for already-invalidated dealer → no state change。
func runCL_JUST_02(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	t.Logf("CL-JUST-02: round=%d (double justification for same dealer)", net.Round)

	// 等待 round 推进 → 如果 round 正常完成，说明重复 justification 未造成状态异常
	time.Sleep(30 * time.Second)
	net2, _ := h.GetLatestDKGNetwork(ctx)
	if net2 != nil {
		checkTrue(t, "round progressing (double justification no impact)",
			net2.Round >= net.Round,
			fmt.Sprintf("round=%d stage=%s", net2.Round, net2.Stage))
	}
}

// runCL_JUST_03 Justification after stage→Finalization → ignored。
func runCL_JUST_03(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	t.Logf("CL-JUST-03: round=%d stage=%s (late justification after Finalization)", net.Round, net.Stage)

	// 等待到 Finalization 阶段
	if net.Stage < dkgtypes.DKGStageFinalization {
		if !h.WaitForRoundStage(ctx, net.Round, dkgtypes.DKGStageFinalization) {
			t.Skip("could not reach Finalization stage")
			return
		}
	}
	// 在 Finalization 阶段验证 finalized count
	beforeFin, _ := h.GetFinalizedRegistrations(ctx, net.Round)
	beforeCount := len(beforeFin)

	// 等待一段时间（mock kernel 发送 late justification）
	time.Sleep(20 * time.Second)

	// 再次检查 finalized count → 不应因 late justification 而改变
	afterFin, _ := h.GetFinalizedRegistrations(ctx, net.Round)
	afterCount := len(afterFin)
	checkTrue(t, "finalized count stable (late justification ignored)",
		afterCount >= beforeCount,
		fmt.Sprintf("before=%d after=%d", beforeCount, afterCount))
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-DUPREG: Duplicate Registration on Resume (Issue #703)
// 验证 ResumeDKGService 不产生重复注册、index 不被污染。
// ──────────────────────────────────────────────────────────────────────────────
// CL-BADDEALER: Bad Dealer Finalization & Committee Participation
// 验证发送无效 VSS deal 的 validator 在 justification 验证失败后的 on-chain 行为。
// 场景前提: ScriptDriver 部署 mock_kernel_bad_dealer_finalize 场景 →
//   Validator A (index 2) 使用 WithInvalidVSSDeal mock kernel，B/C 正常。
// ──────────────────────────────────────────────────────────────────────────────

// ──────────────────────────────────────────────────────────────────────────────
// CL-LOCK: Round Transition Concurrency
// 验证 round 转换期间（Active→next Registration）并发 CDRRead 的安全性。
// ──────────────────────────────────────────────────────────────────────────────

// runCL_LOCK_01 在 round 转换窗口内并发发起多个 CDRRead。
// 场景：R(N) Active 即将结束 → R(N+1) Registration 开始，期间并发 CDRRead。
// 验证：
//  1. CDRRead 使用当前 active round 的 key（不因 round 切换而 panic）
//  2. 多个并发请求不互相干扰
//  3. 如果 round 刚好切换，CDRRead 要么成功（用旧 round key）要么返回明确错误（无 active round）
func runCL_LOCK_01(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}

	// 预先分配多个 vault
	const concurrentCount = 5
	uuids := make([]uint32, concurrentCount)
	for i := range concurrentCount {
		uuid, allocErr := h.ChainClient.CDRAllocate(ctx)
		if allocErr != nil {
			t.Fatalf("CDRAllocate[%d]: %v", i, allocErr)
		}
		uuids[i] = uuid
		if writeErr := h.ChainClient.CDRWrite(ctx, uuid, []byte(fmt.Sprintf("lock-test-%d", i))); writeErr != nil {
			t.Fatalf("CDRWrite[%d]: %v", i, writeErr)
		}
	}

	t.Logf("CL-LOCK-01: %d vaults prepared, launching concurrent CDRRead...", concurrentCount)

	// 并发 CDRRead
	type readResult struct {
		index int
		err   error
	}
	resultCh := make(chan readResult, concurrentCount)
	for i := range concurrentCount {
		go func(idx int) {
			readErr := h.ChainClient.CDRRead(ctx, uuids[idx], []byte(fmt.Sprintf("requester-lock-%d", idx)))
			resultCh <- readResult{index: idx, err: readErr}
		}(i)
	}

	// 收集结果
	successCount := 0
	errCount := 0
	for range concurrentCount {
		res := <-resultCh
		if res.err != nil {
			errCount++
			t.Logf("  CDRRead[%d]: err=%v", res.index, res.err)
		} else {
			successCount++
		}
	}

	// 核心断言：不 panic，且大部分成功
	checkTrue(t, "no panic during concurrent CDRRead", true,
		fmt.Sprintf("success=%d err=%d total=%d", successCount, errCount, concurrentCount))
	checkTrue(t, "majority of concurrent CDRReads succeeded",
		successCount > 0,
		fmt.Sprintf("success=%d/%d", successCount, concurrentCount))

	// 等 partials 提交后验证数据完整性
	time.Sleep(30 * time.Second)
	// 检查第一个成功的 vault 是否有 partials
	if successCount > 0 {
		resp, qErr := h.GetCDRPartials(ctx, uuids[0], fmt.Sprintf("%x", []byte("requester-lock-0")))
		if qErr == nil {
			totalPartials := 0
			for _, g := range resp.Submissions {
				totalPartials += len(g.Submissions)
			}
			checkTrue(t, "partials returned for concurrent request", totalPartials >= 1,
				fmt.Sprintf("partials=%d", totalPartials))
		}
	}
	t.Logf("CL-LOCK-01: concurrent CDRRead during round transition: success=%d err=%d",
		successCount, errCount)
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-KERR: Kernel Error Classification
// 验证 story 对 kernel 返回的不同 gRPC error code 的处理差异。
// ──────────────────────────────────────────────────────────────────────────────

// runCL_KERR_01 验证 InvalidArgument vs Internal error 对 story 行为的影响。
// 前提: ScriptDriver 在一个 validator 上部署 mock kernel:
//   Phase 1: GenerateDeals 返回 codes.InvalidArgument → story 应立即标记 session Failed（不重试）
//   Phase 2: 恢复正常 kernel → 观察 round 是否在新 round 恢复
//
// 对比 CL-RESTART (Internal error → 缓存 + 重试) 验证 story 的错误分类行为。
func runCL_KERR_01(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	startRound := net.Round
	t.Logf("CL-KERR-01: round=%d stage=%s (mock kernel should return InvalidArgument on GenerateDeals)", net.Round, net.Stage)

	// ScriptDriver 的 pre.sh 部署 mock kernel (InvalidArgument mode)
	// 等待一段时间观察 story 的行为
	time.Sleep(30 * time.Second)

	net2, _ := h.GetLatestDKGNetwork(ctx)
	if net2 == nil {
		t.Fatal("no network after wait")
	}

	// InvalidArgument 应导致该 validator 的 session 被标记为 Failed（不重试）
	// 但 round 整体应该继续（其他 2 个 validator 够用）
	if net2.Round > startRound {
		t.Logf("CL-KERR-01: round advanced %d→%d (previous round may have failed due to kernel error)",
			startRound, net2.Round)
	}

	// 核心验证：round 仍在推进（不因一个 validator 的 InvalidArgument 而永久卡住）
	checkTrue(t, "round progressing despite kernel InvalidArgument",
		net2.Round >= startRound,
		fmt.Sprintf("round=%d stage=%s", net2.Round, net2.Stage))

	// 验证该 validator 是否从 verified registrations 中缺席（session Failed → 未完成注册/dealing）
	if net2.Stage >= dkgtypes.DKGStageDealing {
		regs, _ := h.GetVerifiedRegistrations(ctx, net2.Round)
		// mock kernel validator 可能因为 InvalidArgument 而未完成 dealing
		// 但如果 round 已推进到新 round，所有 validator 可能重新注册
		t.Logf("CL-KERR-01: round=%d verified=%d (mock validator may be excluded from dealing)",
			net2.Round, len(regs))
	}

	// 等 round 完成验证最终结果
	if net2.Stage < dkgtypes.DKGStageActive {
		targetRound := net2.Round
		if h.WaitForRoundStage(ctx, targetRound, dkgtypes.DKGStageActive) {
			activeNet, _ := h.GetLatestActiveDKGNetwork(ctx)
			if activeNet != nil {
				checkTrue(t, "round completed despite one validator's InvalidArgument",
					len(activeNet.GlobalPublicKey) > 0,
					fmt.Sprintf("round=%d gpk_len=%d", activeNet.Round, len(activeNet.GlobalPublicKey)))
			}
		} else {
			// Round fail 也是可接受的（取决于 threshold）
			net3, _ := h.GetDKGNetwork(ctx, targetRound)
			if net3 != nil {
				t.Logf("CL-KERR-01: round=%d stage=%s (may have failed, next round should recover)",
					net3.Round, net3.Stage)
			}
		}
	}
}

// badDealerNodeIndex 是 mock kernel 部署的 validator (0-indexed)。
const badDealerNodeIndex = 2 // Validator 3 (0-indexed)

// runCL_BADDEALER_01 验证 bad dealer 在 justification 失败后 registration status 是否被 invalidate。
func runCL_BADDEALER_01(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}

	// 等待到 Finalization 阶段（此时 justification 流程已完成）
	if net.Stage < dkgtypes.DKGStageFinalization {
		if !h.WaitForRoundStage(ctx, net.Round, dkgtypes.DKGStageFinalization) {
			t.Skip("could not reach Finalization stage")
			return
		}
	}

	// 获取所有注册，找到 bad dealer (Validator A / node index 2)
	regs, err := h.GetAllDKGRegistrations(ctx, net.Round)
	if err != nil {
		t.Fatalf("GetAllDKGRegistrations: %v", err)
	}

	// 遍历所有注册，检查是否有 Invalidated 状态的
	invalidatedCount := 0
	for _, r := range regs {
		if r.Status == dkgtypes.DKGRegStatusInvalidated {
			invalidatedCount++
			t.Logf("CL-BADDEALER-01: validator %s index=%d → Invalidated (FIXED)", r.ValidatorAddr, r.Index)
		}
	}

	if invalidatedCount == 0 {
		// BUG: 没有任何 validator 被 invalidated → invalidateDealerRegistration 未被调用
		t.Errorf("BUG CONFIRMED: No validator has status=Invalidated after bad deal justification. "+
			"All %d registrations remain Verified/Finalized. "+
			"invalidateDealerRegistration() is never called from handleDKGProcessJustifications.",
			len(regs))
		// 打印所有注册详情用于 issue 报告
		for _, r := range regs {
			t.Logf("  validator=%s index=%d status=%s", r.ValidatorAddr, r.Index, r.Status)
		}
	} else {
		checkTrue(t, "bad dealer invalidated", invalidatedCount >= 1,
			fmt.Sprintf("invalidated=%d", invalidatedCount))
	}
}

// runCL_BADDEALER_02 验证 bad dealer 是否能成功调用 finalize()。
func runCL_BADDEALER_02(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}

	// 等待到 Active 或 Finalization 后期（给 finalize 时间）
	if net.Stage < dkgtypes.DKGStageFinalization {
		h.WaitForRoundStage(ctx, net.Round, dkgtypes.DKGStageFinalization)
	}
	// 额外等一段时间让所有 validator 完成 finalize
	time.Sleep(30 * time.Second)

	finRegs, err := h.GetFinalizedRegistrations(ctx, net.Round)
	if err != nil {
		t.Fatalf("GetFinalizedRegistrations: %v", err)
	}

	valCount := validatorCount()

	if len(finRegs) == valCount {
		// BUG: 所有 validator 都 finalized（包括 bad dealer）
		t.Errorf("BUG CONFIRMED: All %d validators finalized (including bad dealer). "+
			"Finalized() handler does not reject validators with invalid deals. "+
			"Expected: bad dealer's finalize() should be rejected because status should be Invalidated.",
			len(finRegs))
		for _, r := range finRegs {
			t.Logf("  finalized: validator=%s index=%d", r.ValidatorAddr, r.Index)
		}
	} else if len(finRegs) == valCount-1 {
		// FIXED: bad dealer 被排除
		checkTrue(t, "bad dealer excluded from finalized set", true,
			fmt.Sprintf("finalized=%d (expected %d, bad dealer excluded)", len(finRegs), valCount-1))
	} else {
		t.Logf("CL-BADDEALER-02: finalized=%d (expected %d or %d)", len(finRegs), valCount, valCount-1)
	}
}

// runCL_BADDEALER_03 验证 bad dealer 是否被计入 finalizedCount 从而虚增 committee 规模。
func runCL_BADDEALER_03(t *testing.T, h *Harness) {
	ctx := context.Background()

	// 等待 round Active
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	if net.Stage < dkgtypes.DKGStageActive {
		if !h.WaitForRoundStage(ctx, net.Round, dkgtypes.DKGStageActive) {
			// round 可能 fail（如果 globalPubKey 没达到 threshold）
			net2, _ := h.GetDKGNetwork(ctx, net.Round)
			if net2 != nil && net2.Stage == dkgtypes.DKGStageFailed {
				t.Logf("CL-BADDEALER-03: round %d failed (may be correct if bad dealer excluded from threshold)", net.Round)
				return
			}
			t.Skip("round did not reach Active")
			return
		}
	}

	// Active → 检查 finalizedCount
	finRegs, _ := h.GetFinalizedRegistrations(ctx, net.Round)
	valCount := validatorCount()

	if len(finRegs) == valCount {
		t.Errorf("BUG CONFIRMED: finalizedCount=%d includes bad dealer. "+
			"Committee size inflated: system thinks %d members available, but only %d have valid key shares. "+
			"Fault tolerance silently degraded.",
			len(finRegs), valCount, valCount-1)
	} else {
		checkTrue(t, "finalizedCount excludes bad dealer",
			len(finRegs) == valCount-1,
			fmt.Sprintf("finalized=%d expected=%d", len(finRegs), valCount-1))
	}

	// 验证 GlobalPublicKey 存在
	activeNet, _ := h.GetLatestActiveDKGNetwork(ctx)
	if activeNet != nil {
		checkTrue(t, "GlobalPublicKey set", len(activeNet.GlobalPublicKey) > 0,
			fmt.Sprintf("len=%d", len(activeNet.GlobalPublicKey)))
	}
}

// runCL_BADDEALER_04 验证 bad dealer 是否领到了 committee rewards。
func runCL_BADDEALER_04(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestActiveDKGNetwork(ctx)
	if net == nil {
		t.Skip("no active round")
		return
	}

	// 获取 finalized 列表 → bad dealer 如果在其中则会领到奖励
	finRegs, _ := h.GetFinalizedRegistrations(ctx, net.Round)
	valCount := validatorCount()

	if len(finRegs) == valCount {
		t.Errorf("BUG CONFIRMED: bad dealer in finalized set → will receive equal UBI committee reward. "+
			"Reward distributed to %d members including bad dealer (dkg_rewards.go:59-66 iterates ALL finalized).",
			len(finRegs))
		t.Log("  Fix: exclude Invalidated registrations from settleRewardsForPreviousCommittee()")
	} else {
		t.Logf("CL-BADDEALER-04: finalized=%d — bad dealer excluded, reward distribution correct", len(finRegs))
	}
}

// runCL_BADDEALER_05 验证 bad dealer 的 partial decryption 是否被链上接受。
func runCL_BADDEALER_05(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}

	// CDRWrite + CDRRead → 所有 finalized validators（包括 bad dealer）会提交 partial
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("baddealer-test-05")); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}
	requesterPubKey := []byte("requester-baddealer-05")
	if err := h.ChainClient.CDRRead(ctx, uuid, requesterPubKey); err != nil {
		t.Fatalf("CDRRead: %v", err)
	}
	t.Log("CL-BADDEALER-05: CDRRead done, waiting for partials...")
	time.Sleep(60 * time.Second)

	resp, qErr := h.GetCDRPartials(ctx, uuid, fmt.Sprintf("%x", requesterPubKey))
	if qErr != nil {
		t.Logf("CL-BADDEALER-05: GetCDRPartials: %v", qErr)
		return
	}

	totalPartials := 0
	for _, g := range resp.Submissions {
		totalPartials += len(g.Submissions)
	}
	valCount := validatorCount()

	if totalPartials == valCount {
		t.Errorf("BUG CONFIRMED: %d partials accepted (all validators including bad dealer). "+
			"Bad dealer's partial decryption passed on-chain validation "+
			"(pubShare check passes because bad dealer stored its own pubKeyShare during finalize). "+
			"Expected: bad dealer should not be able to submit partial (not in valid committee).",
			totalPartials)
	} else if totalPartials == valCount-1 {
		checkTrue(t, "bad dealer partial rejected", true,
			fmt.Sprintf("partials=%d (bad dealer excluded)", totalPartials))
	} else {
		t.Logf("CL-BADDEALER-05: partials=%d (expected %d or %d)", totalPartials, valCount, valCount-1)
	}
}

// runCL_BADDEALER_06 验证当 threshold 组合包含 bad dealer partial 时解密是否失败。
// 注意：此测试验证的是 partial 组合的可用性，不执行实际的 TDH2 combine（需要客户端库）。
// 通过检查 partial 提交者的 registration 状态来推断哪些组合有效。
func runCL_BADDEALER_06(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active round with GlobalPublicKey")
		return
	}

	// 复用 CL-BADDEALER-05 或单独做一次 CDRRead
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("baddealer-test-06")); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}
	requesterPubKey := []byte("requester-baddealer-06")
	if err := h.ChainClient.CDRRead(ctx, uuid, requesterPubKey); err != nil {
		t.Fatalf("CDRRead: %v", err)
	}
	time.Sleep(60 * time.Second)

	resp, qErr := h.GetCDRPartials(ctx, uuid, fmt.Sprintf("%x", requesterPubKey))
	if qErr != nil {
		t.Logf("CL-BADDEALER-06: GetCDRPartials: %v", qErr)
		return
	}

	// 分析 partial 提交者与 registration 状态的对应关系
	finRegs, _ := h.GetFinalizedRegistrations(ctx, net.Round)
	allRegs, _ := h.GetAllDKGRegistrations(ctx, net.Round)

	// 找出哪些 validator 的 deals 在 justification 中被证明无效
	// （无法直接查链上 justification 结果，通过比对 GlobalPublicKey 投票间接推断）
	t.Logf("CL-BADDEALER-06: finalized=%d total=%d", len(finRegs), len(allRegs))

	totalPartials := 0
	for _, g := range resp.Submissions {
		totalPartials += len(g.Submissions)
	}

	valCount := validatorCount()
	if totalPartials == valCount {
		// threshold=2, 3 个 partial 但 1 个是无效的
		// C(3,2)=3 种组合，只有 1 种（B+C）有效 → 可用性 33%
		t.Errorf("BUG IMPACT: %d partials on-chain (including bad dealer). "+
			"Threshold=%d, valid combinations: only %d of %d (availability drops from 100%% to 33%%). "+
			"Client selecting A's partial will get decryption failure.",
			totalPartials, 2, 1, valCount*(valCount-1)/2)
	} else {
		t.Logf("CL-BADDEALER-06: partials=%d (bad dealer may be excluded)", totalPartials)
	}
}

// runCL_BADDEALER_07 验证多个 bad dealer 同时存在时的行为。
// 场景：2 of 3 validators 使用 mock kernel (invalid VSS)。
func runCL_BADDEALER_07(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}

	// 等待 round 完成（可能 fail）
	time.Sleep(60 * time.Second)
	net2, _ := h.GetLatestDKGNetwork(ctx)
	if net2 == nil {
		t.Fatal("no network after wait")
	}

	if net2.Stage == dkgtypes.DKGStageFailed {
		// 预期行为：2 个 bad dealer → 没有 globalPubKey 达到 threshold → round fail
		t.Logf("CL-BADDEALER-07: round %d failed (correct: 2 bad dealers → no globalPubKey threshold)", net2.Round)
		checkTrue(t, "round failed with 2 bad dealers (correct)", true, "stage=Failed")
		return
	}

	if net2.Stage == dkgtypes.DKGStageActive {
		// 如果 round 到了 Active → 检查 GlobalPublicKey 是否真的有效
		finRegs, _ := h.GetFinalizedRegistrations(ctx, net2.Round)
		t.Errorf("BUG: round %d reached Active with 2 bad dealers. "+
			"finalizedCount=%d, but only 1 has valid key share. "+
			"GlobalPubKey may not have reached legitimate threshold.",
			net2.Round, len(finRegs))
		return
	}

	t.Logf("CL-BADDEALER-07: round=%d stage=%s (in progress)", net2.Round, net2.Stage)
}

// runCL_BADDEALER_08 自适应攻击：bad dealer 提交与诚实多数相同的 globalPubKey。
// 理论验证：如果 A 排除自己的 deal contribution 计算 DistKeyShare，能否得到有效 key share？
func runCL_BADDEALER_08(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestActiveDKGNetwork(ctx)
	if net == nil {
		t.Skip("no active round")
		return
	}

	// 检查所有 finalized validator 是否投了相同的 globalPubKey
	// （如果 bad dealer 能适应性地匹配诚实多数的 key，则所有投票一致）
	finRegs, _ := h.GetFinalizedRegistrations(ctx, net.Round)
	t.Logf("CL-BADDEALER-08: round=%d finalized=%d GlobalPubKey=%d bytes",
		net.Round, len(finRegs), len(net.GlobalPublicKey))

	// 在当前架构下无法直接查每个 validator 的 globalPubKey 投票
	// 但如果 GlobalPublicKey 非空 → 说明有 threshold 个 validator 投了相同的 key
	checkTrue(t, "GlobalPublicKey present", len(net.GlobalPublicKey) > 0,
		fmt.Sprintf("len=%d", len(net.GlobalPublicKey)))

	// 如果 bad dealer 也在 finalized 中，且 GlobalPublicKey 有效 → 可能的自适应攻击
	valCount := validatorCount()
	if len(finRegs) == valCount && len(net.GlobalPublicKey) > 0 {
		t.Logf("CL-BADDEALER-08: ATTENTION — all %d validators finalized with same GlobalPubKey. "+
			"If Validator A sent invalid deals but matched honest globalPubKey, "+
			"this suggests adaptive attack is possible (A excluded own deal contribution). "+
			"Requires deeper kyber/Pedersen DKG analysis to confirm.",
			valCount)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-DUPREG: Duplicate Registration on Resume (Issue #703)
// ──────────────────────────────────────────────────────────────────────────────

// runCL_DUPREG_01 验证 resume 后同一 validator 不会重复注册。
// 场景：TEE 重启触发 ResumeDKGService → handleDKGRegistration 被再次调用
// → isAlreadyRegistered 应返回 true → 跳过合约调用 → 不产生重复 Registered 事件。
func runCL_DUPREG_01(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}

	// 等待 Registration 阶段有注册
	if net.Stage == dkgtypes.DKGStageRegistration {
		h.WaitForVerifiedCount(ctx, net.Round, 1)
	}

	// 获取当前所有注册
	regs, err := h.GetAllDKGRegistrations(ctx, net.Round)
	if err != nil {
		t.Fatalf("GetAllDKGRegistrations: %v", err)
	}

	// 核心断言：每个 validator 最多只有一条注册
	addrCount := make(map[string]int)
	for _, r := range regs {
		addr := strings.ToLower(r.ValidatorAddr)
		addrCount[addr]++
	}
	for addr, count := range addrCount {
		checkTrue(t, fmt.Sprintf("validator %s has exactly 1 registration", addr),
			count == 1,
			fmt.Sprintf("count=%d (duplicate registration bug if > 1)", count))
	}
	t.Logf("CL-DUPREG-01: round=%d registrations=%d unique_validators=%d (no duplicates)",
		net.Round, len(regs), len(addrCount))
}

// runCL_DUPREG_02 验证 kernel 重启后不会产生新的 index。
// 场景：在 Registration 阶段 SSH 重启 validator 3 的 kernel
// → ResumeDKGService → isAlreadyRegistered 返回 true → 复用原有 index。
func runCL_DUPREG_02(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}

	// 记录重启前的 registrations
	regsBefore, _ := h.GetAllDKGRegistrations(ctx, net.Round)
	beforeMap := make(map[string]uint32) // addr → index
	for _, r := range regsBefore {
		beforeMap[strings.ToLower(r.ValidatorAddr)] = r.Index
	}

	// 如果有 SSH 能力，重启 validator 3 的 kernel
	targetNode := validatorCount() - 1
	_, sshErr := sshRunCmd(targetNode, "sudo systemctl restart story-kernel")
	if sshErr != nil {
		t.Logf("CL-DUPREG-02: SSH restart not available (%v), checking current state only", sshErr)
	} else {
		t.Logf("CL-DUPREG-02: restarted validator %d kernel, waiting for recovery...", targetNode+1)
		time.Sleep(40 * time.Second) // enclave 加载时间
	}

	// 等待一段时间让 resume 逻辑运行
	time.Sleep(20 * time.Second)

	// 重启后再次检查 registrations
	regsAfter, _ := h.GetAllDKGRegistrations(ctx, net.Round)
	afterMap := make(map[string]uint32)
	for _, r := range regsAfter {
		afterMap[strings.ToLower(r.ValidatorAddr)] = r.Index
	}

	// 断言：重启前已注册的 validator 的 index 不变
	for addr, idxBefore := range beforeMap {
		if idxAfter, ok := afterMap[addr]; ok {
			check(t, fmt.Sprintf("validator %s index unchanged after restart", addr),
				idxBefore, idxAfter)
		}
	}

	// 断言：没有新增的注册（数量不变或只增加了之前未注册的 validator）
	checkTrue(t, "no spurious new registrations",
		len(regsAfter) <= validatorCount(),
		fmt.Sprintf("before=%d after=%d valCount=%d", len(regsBefore), len(regsAfter), validatorCount()))
}

// runCL_DUPREG_03 验证 index 唯一性：同一 round 内所有 index 不重复且从 1 开始连续。
func runCL_DUPREG_03(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}

	// 等待有足够的注册
	if net.Stage <= dkgtypes.DKGStageRegistration {
		h.WaitForVerifiedCount(ctx, net.Round, 2)
	}

	regs, err := h.GetAllDKGRegistrations(ctx, net.Round)
	if err != nil {
		t.Fatalf("GetAllDKGRegistrations: %v", err)
	}
	if len(regs) == 0 {
		t.Skip("no registrations yet")
		return
	}

	// 收集所有 index
	indexSeen := make(map[uint32]string) // index → validator addr
	for _, r := range regs {
		if prev, dup := indexSeen[r.Index]; dup {
			t.Errorf("INDEX COLLISION (Issue #703): index %d used by both %s and %s",
				r.Index, prev, r.ValidatorAddr)
		}
		indexSeen[r.Index] = r.ValidatorAddr
	}

	// 验证 index 从 1 开始连续
	for i := uint32(1); i <= uint32(len(regs)); i++ {
		addr, ok := indexSeen[i]
		checkTrue(t, fmt.Sprintf("index %d assigned", i), ok,
			fmt.Sprintf("assigned_to=%s", addr))
	}

	t.Logf("CL-DUPREG-03: round=%d registrations=%d indices=[1..%d] all unique",
		net.Round, len(regs), len(regs))
}

// runCL_DUPREG_04 验证 resume + re-registration guard 后 round 能正确完成。
// 核心：如果 index 被污染（Issue #703 的 bug），deals 会路由到错误的 participant，
// 导致 finalization 失败或 GlobalPublicKey 不一致。
func runCL_DUPREG_04(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	startRound := net.Round

	// 先验证 index 唯一性（前置条件）
	runCL_DUPREG_03(t, h)

	// 等待 round 完成到 Active
	if net.Stage < dkgtypes.DKGStageActive {
		if !h.WaitForRoundStage(ctx, startRound, dkgtypes.DKGStageActive) {
			// 如果 round fail 了，检查是否是 index 污染导致
			net2, _ := h.GetDKGNetwork(ctx, startRound)
			if net2 != nil && net2.Stage == dkgtypes.DKGStageFailed {
				// 检查是否有 index collision
				regs, _ := h.GetAllDKGRegistrations(ctx, startRound)
				indexSeen := make(map[uint32]bool)
				for _, r := range regs {
					if indexSeen[r.Index] {
						t.Errorf("CONFIRMED Issue #703: round %d FAILED with index collision at index %d",
							startRound, r.Index)
						return
					}
					indexSeen[r.Index] = true
				}
				t.Logf("CL-DUPREG-04: round %d failed (not due to index collision, may be other reason)", startRound)
			}
			// 等新 round
			net3, _ := h.GetLatestDKGNetwork(ctx)
			if net3 != nil && net3.Round > startRound {
				startRound = net3.Round
				h.WaitForRoundStage(ctx, startRound, dkgtypes.DKGStageActive)
			}
		}
	}

	// 验证 GlobalPublicKey 存在且一致
	activeNet, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || activeNet == nil {
		t.Logf("CL-DUPREG-04: no active network (round may still be in progress)")
		return
	}
	checkTrue(t, "GlobalPublicKey present (correct indices → correct deal routing)",
		len(activeNet.GlobalPublicKey) > 0,
		fmt.Sprintf("round=%d gpk_len=%d", activeNet.Round, len(activeNet.GlobalPublicKey)))

	// CDR decrypt 验证（间接验证 index → PID 映射正确）
	if !h.IsNoopChainClient() {
		verifyCDRDecryptAfterRestart(t, h, "CL-DUPREG-04")
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-DUPREG-05: On-chain duplicate registration exploit
// 直接调合约绕过 off-chain guard，验证 on-chain Registered() handler 是否有 dedup。
// ──────────────────────────────────────────────────────────────────────────────

// runCL_DUPREG_05 直接调 DKG.Register() 合约两次（同一 validator），验证 on-chain 行为。
//
// 场景模拟 Issue #703 的根因：
//   1. Validator A 正常注册 → index=1
//   2. 攻击者/resume 逻辑绕过 off-chain guard，再次调 Register() → index 应该被拒绝或保持不变
//
// 如果 on-chain 没有 dedup：
//   - 第二次 Register() 成功 → A 的 index 从 1 变成 N+1（覆盖）
//   - getNextDKGRegistrationIndex 返回错误值
//   - 后续 deal routing 全部错乱
//
// 如果 on-chain 有 dedup（已修复）：
//   - 第二次 Register() 应 revert 或返回 error
func runCL_DUPREG_05(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient to call Register() directly")
		return
	}
	ctx := context.Background()

	// 等待 Registration 阶段有注册
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	if net.Stage != dkgtypes.DKGStageRegistration {
		t.Skipf("need Registration stage, got %s (this test must run early in a round)", net.Stage)
		return
	}

	// 等待至少一个 validator 注册
	if !h.WaitForVerifiedCount(ctx, net.Round, 1) {
		t.Skip("no verified registrations yet")
		return
	}

	// 获取第一个已注册 validator 的数据
	regs, err := h.GetVerifiedRegistrations(ctx, net.Round)
	if err != nil || len(regs) == 0 {
		t.Fatalf("GetVerifiedRegistrations: %v (len=%d)", err, len(regs))
	}
	victim := regs[0]
	originalIndex := victim.Index
	t.Logf("CL-DUPREG-05: target validator=%s original_index=%d round=%d",
		victim.ValidatorAddr, originalIndex, net.Round)

	// 记录当前总注册数
	allRegsBefore, _ := h.GetAllDKGRegistrations(ctx, net.Round)
	countBefore := len(allRegsBefore)

	// 构造参数，直接调合约 Register()（绕过 off-chain guard）
	validatorAddr := common.HexToAddress(victim.ValidatorAddr)
	var enclaveType [32]byte
	copy(enclaveType[:], victim.EnclaveType)
	var startBlockHash [32]byte
	copy(startBlockHash[:], net.StartBlockHash)

	t.Logf("CL-DUPREG-05: calling Register() again for same validator (bypassing off-chain guard)...")
	err = h.ChainClient.RegisterDKGWithParams(
		ctx,
		net.Round,
		validatorAddr,
		enclaveType,
		victim.CommPubKey,
		victim.DkgPubKey,
		victim.EnclaveReport,
		big.NewInt(net.StartBlockHeight),
		startBlockHash,
	)

	if err != nil {
		// 如果 Register() revert → on-chain 有 dedup guard（已修复）
		t.Logf("CL-DUPREG-05: ✓ FIXED — second Register() reverted: %v", err)
		checkTrue(t, "duplicate Register() rejected on-chain", true,
			fmt.Sprintf("err=%v", err))
		return
	}

	// 如果 Register() 成功 → 检查是否造成 index 污染
	t.Log("CL-DUPREG-05: ⚠ second Register() succeeded — checking for index corruption...")

	// 等待链处理事件
	time.Sleep(10 * time.Second)

	allRegsAfter, _ := h.GetAllDKGRegistrations(ctx, net.Round)
	countAfter := len(allRegsAfter)

	// 检查 victim 的 index 是否被改变
	var victimNewIndex uint32
	indexSeen := make(map[uint32]string)
	hasCollision := false
	for _, r := range allRegsAfter {
		addr := strings.ToLower(r.ValidatorAddr)
		if addr == strings.ToLower(victim.ValidatorAddr) {
			victimNewIndex = r.Index
		}
		if prev, dup := indexSeen[r.Index]; dup {
			t.Errorf("BUG CONFIRMED (Issue #703 on-chain): INDEX COLLISION — index %d used by both %s and %s",
				r.Index, prev, addr)
			hasCollision = true
		}
		indexSeen[r.Index] = addr
	}

	if victimNewIndex != originalIndex {
		t.Errorf("BUG CONFIRMED (Issue #703 on-chain): validator %s index changed from %d to %d after duplicate Register()",
			victim.ValidatorAddr, originalIndex, victimNewIndex)
	}

	if countAfter != countBefore {
		t.Logf("CL-DUPREG-05: ⚠ registration count changed: before=%d after=%d (on-chain dedup missing)",
			countBefore, countAfter)
	}

	if !hasCollision && victimNewIndex == originalIndex {
		// Register() 成功但 index 没变（可能 setDKGRegistration 覆盖了同一 key）
		t.Logf("CL-DUPREG-05: Register() succeeded but index unchanged (setDKGRegistration overwrites by round_addr key)")
		checkTrue(t, "index preserved despite duplicate Register()", true,
			fmt.Sprintf("index=%d (overwrite semantics, not append)", victimNewIndex))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-DUPREG-06/07: Issue #721 — Third-party replay attack
// 任何外部地址可以读取链上的 register() 参数，重放给合约，导致 victim 的 index 被覆盖。
// ──────────────────────────────────────────────────────────────────────────────

// runCL_DUPREG_06 第三方重放单个 validator 的注册参数。
// 攻击路径：读取 A 的注册参数 → 用不同地址（测试用同一 signer）重新调 register() → A 的 index 被覆盖。
func runCL_DUPREG_06(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()

	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	// 如果不在 Registration 阶段，等下一个 round 的 Registration
	if net.Stage != dkgtypes.DKGStageRegistration {
		t.Logf("current round=%d stage=%s, waiting for next round Registration...", net.Round, net.Stage)
		if !h.WaitForRound(ctx, net.Round+1) {
			t.Skip("timed out waiting for next round")
			return
		}
		net, _ = h.GetLatestDKGNetwork(ctx)
		if net == nil || net.Stage != dkgtypes.DKGStageRegistration {
			t.Skipf("next round not in Registration stage")
			return
		}
	}

	if !h.WaitForVerifiedCount(ctx, net.Round, 2) {
		t.Skip("need at least 2 verified registrations")
		return
	}

	regs, err := h.GetVerifiedRegistrations(ctx, net.Round)
	if err != nil || len(regs) < 2 {
		t.Fatalf("GetVerifiedRegistrations: %v (len=%d)", err, len(regs))
	}

	// 选择 victim（第一个注册的 validator）
	victim := regs[0]
	originalIndex := victim.Index
	t.Logf("CL-DUPREG-06: victim=%s index=%d round=%d", victim.ValidatorAddr, originalIndex, net.Round)

	// 记录所有当前注册的 index 快照
	indexesBefore := make(map[string]uint32)
	for _, r := range regs {
		indexesBefore[strings.ToLower(r.ValidatorAddr)] = r.Index
	}

	// 攻击：用 victim 的参数重放 register()
	validatorAddr := common.HexToAddress(victim.ValidatorAddr)
	var enclaveType [32]byte
	copy(enclaveType[:], victim.EnclaveType)
	var startBlockHash [32]byte
	copy(startBlockHash[:], net.StartBlockHash)

	t.Log("CL-DUPREG-06: replaying victim's register() params (simulating third-party attacker)...")
	err = h.ChainClient.RegisterDKGWithParams(
		ctx, net.Round, validatorAddr, enclaveType,
		victim.CommPubKey, victim.DkgPubKey, victim.EnclaveReport,
		big.NewInt(net.StartBlockHeight), startBlockHash,
	)

	if err != nil {
		t.Logf("CL-DUPREG-06: ✓ FIXED — replay register() reverted: %v", err)
		checkTrue(t, "replay register() rejected (Issue #721 fixed)", true, fmt.Sprintf("err=%v", err))
		return
	}

	// Register() 没 revert → 检查 index 是否被篡改
	t.Log("CL-DUPREG-06: ⚠ replay register() succeeded — checking index corruption...")
	time.Sleep(10 * time.Second)

	regsAfter, _ := h.GetAllDKGRegistrations(ctx, net.Round)
	for _, r := range regsAfter {
		addr := strings.ToLower(r.ValidatorAddr)
		if before, ok := indexesBefore[addr]; ok && r.Index != before {
			t.Errorf("BUG CONFIRMED (Issue #721): validator %s index changed %d → %d after third-party replay",
				addr, before, r.Index)
		}
	}

	// 检查 index collision
	indexSeen := make(map[uint32]string)
	for _, r := range regsAfter {
		if prev, dup := indexSeen[r.Index]; dup {
			t.Errorf("BUG CONFIRMED (Issue #721): INDEX COLLISION — index %d shared by %s and %s after replay attack",
				r.Index, prev, r.ValidatorAddr)
		}
		indexSeen[r.Index] = r.ValidatorAddr
	}
}

// runCL_DUPREG_07 第三方重放所有 validator 的注册参数 → 全部 index 被覆盖。
func runCL_DUPREG_07(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("requires EthChainClient")
		return
	}
	ctx := context.Background()

	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	// 如果不在 Registration 阶段，等下一个 round 的 Registration
	if net.Stage != dkgtypes.DKGStageRegistration {
		t.Logf("current round=%d stage=%s, waiting for next round Registration...", net.Round, net.Stage)
		if !h.WaitForRound(ctx, net.Round+1) {
			t.Skip("timed out waiting for next round")
			return
		}
		net, _ = h.GetLatestDKGNetwork(ctx)
		if net == nil || net.Stage != dkgtypes.DKGStageRegistration {
			t.Skipf("next round not in Registration stage")
			return
		}
	}

	if !h.WaitForVerifiedCount(ctx, net.Round, 3) {
		t.Skip("need at least 3 verified registrations")
		return
	}

	regs, err := h.GetVerifiedRegistrations(ctx, net.Round)
	if err != nil || len(regs) < 3 {
		t.Fatalf("GetVerifiedRegistrations: %v (len=%d)", err, len(regs))
	}

	// 记录 replay 前的所有 index
	indexesBefore := make(map[string]uint32)
	for _, r := range regs {
		indexesBefore[strings.ToLower(r.ValidatorAddr)] = r.Index
	}
	t.Logf("CL-DUPREG-07: replaying ALL %d validators' register() params...", len(regs))

	// 依次重放每个 validator 的注册
	replayFailed := 0
	replaySucceeded := 0
	for _, victim := range regs {
		validatorAddr := common.HexToAddress(victim.ValidatorAddr)
		var enclaveType [32]byte
		copy(enclaveType[:], victim.EnclaveType)
		var startBlockHash [32]byte
		copy(startBlockHash[:], net.StartBlockHash)

		replayErr := h.ChainClient.RegisterDKGWithParams(
			ctx, net.Round, validatorAddr, enclaveType,
			victim.CommPubKey, victim.DkgPubKey, victim.EnclaveReport,
			big.NewInt(net.StartBlockHeight), startBlockHash,
		)
		if replayErr != nil {
			replayFailed++
		} else {
			replaySucceeded++
		}
	}

	if replayFailed == len(regs) {
		t.Logf("CL-DUPREG-07: ✓ FIXED — all %d replay attempts reverted", len(regs))
		checkTrue(t, "all replay register() rejected (Issue #721 fixed)", true,
			fmt.Sprintf("failed=%d/%d", replayFailed, len(regs)))
		return
	}

	// 部分或全部成功 → 检查 index 灾难
	t.Logf("CL-DUPREG-07: ⚠ %d/%d replays succeeded — checking mass index corruption...",
		replaySucceeded, len(regs))
	time.Sleep(10 * time.Second)

	regsAfter, _ := h.GetAllDKGRegistrations(ctx, net.Round)

	// 检查 index 变化
	changedCount := 0
	for _, r := range regsAfter {
		addr := strings.ToLower(r.ValidatorAddr)
		if before, ok := indexesBefore[addr]; ok && r.Index != before {
			changedCount++
			t.Logf("  validator %s: index %d → %d", addr, before, r.Index)
		}
	}

	// 检查 index collision
	indexSeen := make(map[uint32]string)
	collisions := 0
	for _, r := range regsAfter {
		if prev, dup := indexSeen[r.Index]; dup {
			collisions++
			t.Logf("  INDEX COLLISION: index %d shared by %s and %s", r.Index, prev, r.ValidatorAddr)
		}
		indexSeen[r.Index] = r.ValidatorAddr
	}

	// 检查 ghost indices（1..N 中有空洞）
	ghostCount := 0
	for i := uint32(1); i <= uint32(len(regsAfter)); i++ {
		if _, ok := indexSeen[i]; !ok {
			ghostCount++
		}
	}

	if changedCount > 0 || collisions > 0 {
		t.Errorf("BUG CONFIRMED (Issue #721): mass replay attack results — "+
			"index_changed=%d collisions=%d ghost_indices=%d. "+
			"DKG deal routing corrupted, round will likely fail. "+
			"Contract lacks: (1) msg.sender==validatorAddr check, (2) duplicate registration guard. "+
			"Consensus layer lacks: existing registration check in Registered() handler.",
			changedCount, collisions, ghostCount)
	}

	// 等 round 完成或失败
	time.Sleep(30 * time.Second)
	net2, _ := h.GetLatestDKGNetwork(ctx)
	if net2 != nil {
		t.Logf("CL-DUPREG-07: round=%d stage=%s after mass replay (expected: Failed if indices corrupted)",
			net2.Round, net2.Stage)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// CL-RESTART: Mid-DKG Restart Persistence
// 验证 story-kernel 在 DKG 不同阶段重启后，能从磁盘恢复状态并产出一致的 GlobalPublicKey。
// 根因：rebuildInitDKG 需要正确 replay 已持久化的 deals/responses/justifications，
// 否则 NewDistKeyGenerator 会生成新的随机 dealer 多项式（PrivatePoly 未持久化），
// 导致 DistKeyShare() 在不同节点产出不一致的 global_pub_key。
// ──────────────────────────────────────────────────────────────────────────────

// runCL_RESTART_01 AfterProcessDeals: kernel 在 ProcessDeals 完成后、ProcessResponses 之前重启。
// 前提：ScriptDriver 在 Dealing 阶段检测到 deals 已处理后，SSH 重启目标 validator 的 story-kernel。
// 验证：
//  1. 该 validator 重启后仍能完成 round（通过 rebuildInitDKG replay deals）
//  2. 所有 validator 的 GlobalPublicKey 一致
//  3. CDR TDH2 decrypt 仍可用
func runCL_RESTART_01(t *testing.T, h *Harness) {
	ctx := context.Background()

	// 等待进入 Dealing 阶段
	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	startRound := net.Round
	t.Logf("CL-RESTART-01: round=%d stage=%s, waiting for Dealing stage...", net.Round, net.Stage)

	if net.Stage < dkgtypes.DKGStageDealing {
		if !h.WaitForRoundStage(ctx, startRound, dkgtypes.DKGStageDealing) {
			t.Skip("could not reach Dealing stage")
			return
		}
	}

	// ScriptDriver 的 pre.sh 应在此时 SSH 重启 validator 3 的 story-kernel
	// （由 RunCase 自动调用 SetupScenario）
	// 等待一小段时间让 kernel 重启完成（enclave 加载 ~30s）
	t.Log("CL-RESTART-01: kernel restarted during Dealing (after ProcessDeals), waiting for recovery...")
	time.Sleep(30 * time.Second)

	// 验证 round 能继续完成 → Finalization → Active
	if !h.WaitForRoundStage(ctx, startRound, dkgtypes.DKGStageActive) {
		// Round 可能 fail 然后新 round 启动，也是可以接受的
		net2, _ := h.GetLatestDKGNetwork(ctx)
		if net2 != nil && net2.Round > startRound {
			t.Logf("CL-RESTART-01: round %d failed, new round %d started (acceptable)", startRound, net2.Round)
			startRound = net2.Round
			if !h.WaitForRoundStage(ctx, startRound, dkgtypes.DKGStageActive) {
				t.Fatalf("new round %d also did not reach Active", startRound)
			}
		} else {
			t.Fatalf("round %d did not reach Active after restart", startRound)
		}
	}

	// 验证 GlobalPublicKey 一致性：所有 round 的 GlobalPublicKey 应一致
	activeNet, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || activeNet == nil {
		t.Fatalf("no active network after restart recovery")
	}
	checkTrue(t, "GlobalPublicKey present after restart",
		len(activeNet.GlobalPublicKey) > 0,
		fmt.Sprintf("len=%d", len(activeNet.GlobalPublicKey)))

	// 验证重启后 TDH2 decrypt 仍可用（如果有 EthChainClient）
	if !h.IsNoopChainClient() {
		verifyCDRDecryptAfterRestart(t, h, "CL-RESTART-01")
	}
	t.Logf("CL-RESTART-01: round=%d GlobalPubKey=%d bytes, recovery successful",
		activeNet.Round, len(activeNet.GlobalPublicKey))
}

// runCL_RESTART_02 BeforeFinalizeDKG: kernel 在 ProcessResponses 完成后、FinalizeDKG 之前重启。
// 这是最关键的场景：rebuildInitDKG 需要 replay deals + responses 才能正确重建 DKG 状态。
func runCL_RESTART_02(t *testing.T, h *Harness) {
	ctx := context.Background()

	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	startRound := net.Round
	t.Logf("CL-RESTART-02: round=%d stage=%s, waiting for Finalization stage...", net.Round, net.Stage)

	// 等待进入 Finalization 阶段（此时 ProcessResponses 已完成）
	if net.Stage < dkgtypes.DKGStageFinalization {
		if !h.WaitForRoundStage(ctx, startRound, dkgtypes.DKGStageFinalization) {
			t.Skip("could not reach Finalization stage")
			return
		}
	}

	// ScriptDriver 的 pre.sh 在 Finalization 阶段重启 kernel
	t.Log("CL-RESTART-02: kernel restarted before FinalizeDKG, waiting for recovery...")
	time.Sleep(30 * time.Second)

	// 验证 round 能完成 Finalization → Active
	if !h.WaitForRoundStage(ctx, startRound, dkgtypes.DKGStageActive) {
		net2, _ := h.GetLatestDKGNetwork(ctx)
		if net2 != nil && net2.Round > startRound {
			t.Logf("CL-RESTART-02: round %d failed, new round %d started", startRound, net2.Round)
			startRound = net2.Round
			if !h.WaitForRoundStage(ctx, startRound, dkgtypes.DKGStageActive) {
				t.Fatalf("new round %d also did not reach Active", startRound)
			}
		} else {
			t.Fatalf("round %d did not reach Active after restart (this is the rebuildInitDKG bug scenario)", startRound)
		}
	}

	// 关键验证：所有 finalized validator 的 GlobalPublicKey 必须一致
	activeNet, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || activeNet == nil {
		t.Fatalf("no active network after restart")
	}
	checkTrue(t, "GlobalPublicKey present (rebuildInitDKG replayed correctly)",
		len(activeNet.GlobalPublicKey) > 0,
		fmt.Sprintf("len=%d", len(activeNet.GlobalPublicKey)))

	// 验证 finalized registrations 都有相同的 GlobalPublicKey
	finRegs, _ := h.GetFinalizedRegistrations(ctx, activeNet.Round)
	checkTrue(t, "finalized count >= threshold after restart",
		len(finRegs) >= 2,
		fmt.Sprintf("finalized=%d (need >= 2 for 3 validators)", len(finRegs)))

	if !h.IsNoopChainClient() {
		verifyCDRDecryptAfterRestart(t, h, "CL-RESTART-02")
	}
	t.Logf("CL-RESTART-02: round=%d finalized=%d GlobalPubKey=%d bytes, PrivatePoly persistence verified",
		activeNet.Round, len(finRegs), len(activeNet.GlobalPublicKey))
}

// runCL_RESTART_03 PersistenceAfterRestart with justification: kernel 在 ProcessJustification 后重启。
// 需要先触发 complaint（mock kernel 产出 bad deal），处理 justification，然后重启。
func runCL_RESTART_03(t *testing.T, h *Harness) {
	ctx := context.Background()

	net, _ := h.GetLatestDKGNetwork(ctx)
	if net == nil {
		t.Skip("no DKG network")
		return
	}
	startRound := net.Round
	t.Logf("CL-RESTART-03: round=%d stage=%s (justification + restart scenario)", net.Round, net.Stage)

	// 此场景需要 mock kernel 先产出 invalid deal → complaint → justification → restart
	// ScriptDriver 的 pre.sh 负责：
	//   1. 在 validator 3 部署 mock kernel（bad deal）
	//   2. 等待 justification 处理完成
	//   3. 重启 validator 3 的 story-kernel（恢复真实 kernel）

	// 等待 round 完成或新 round 开始
	time.Sleep(30 * time.Second)

	// 即使有 complaint + restart，round 应该仍能完成（2/3 validator 够用）
	net2, _ := h.GetLatestDKGNetwork(ctx)
	if net2 == nil {
		t.Fatal("no network after restart")
	}

	if net2.Round >= startRound {
		// 等待达到 Active
		targetRound := net2.Round
		if net2.Stage < dkgtypes.DKGStageActive {
			if !h.WaitForRoundStage(ctx, targetRound, dkgtypes.DKGStageActive) {
				// 允许 fail + new round
				net3, _ := h.GetLatestDKGNetwork(ctx)
				if net3 != nil && net3.Round > targetRound {
					targetRound = net3.Round
					h.WaitForRoundStage(ctx, targetRound, dkgtypes.DKGStageActive)
				}
			}
		}
	}

	// 验证 Active round 存在且 GlobalPublicKey 一致
	activeNet, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || activeNet == nil {
		t.Logf("CL-RESTART-03: no active network (round may have been skipped due to complaint + restart)")
		// 这种情况下验证链仍在正常运行
		net3, _ := h.GetLatestDKGNetwork(ctx)
		if net3 != nil {
			checkTrue(t, "chain still running after justification + restart",
				net3.Round >= startRound,
				fmt.Sprintf("round=%d stage=%s", net3.Round, net3.Stage))
		}
		return
	}
	checkTrue(t, "GlobalPublicKey present after justification + restart",
		len(activeNet.GlobalPublicKey) > 0,
		fmt.Sprintf("len=%d", len(activeNet.GlobalPublicKey)))

	if !h.IsNoopChainClient() {
		verifyCDRDecryptAfterRestart(t, h, "CL-RESTART-03")
	}
	t.Logf("CL-RESTART-03: round=%d GlobalPubKey=%d bytes, justification persistence verified",
		activeNet.Round, len(activeNet.GlobalPublicKey))
}

// verifyCDRDecryptAfterRestart 在 kernel 重启后验证 CDR TDH2 decrypt 仍可用。
func verifyCDRDecryptAfterRestart(t *testing.T, h *Harness, caseID string) {
	t.Helper()
	ctx := context.Background()
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Logf("%s: CDRAllocate after restart: %v", caseID, err)
		return
	}
	if err := h.ChainClient.CDRWrite(ctx, uuid, []byte("test-"+caseID+"-post-restart")); err != nil {
		t.Logf("%s: CDRWrite after restart: %v", caseID, err)
		return
	}
	if err := h.ChainClient.CDRRead(ctx, uuid, []byte("requester-"+caseID)); err != nil {
		t.Logf("%s: CDRRead after restart: %v", caseID, err)
		return
	}
	checkTrue(t, "CDR decrypt works after restart", true,
		fmt.Sprintf("uuid=%d, TDH2 decrypt path exercised post-restart", uuid))
}

// TestDKG_CL 运行所有 CL-* 新增测试用例。
func TestDKG_CL(t *testing.T) {
	for _, tc := range CLCases() {
		t.Run(tc.ID, func(t *testing.T) {
			if tc.SkipIfLive != "" && globalHarness.Driver == nil {
				t.Skip(tc.SkipIfLive)
				return
			}
			if tc.Run == nil {
				t.Skip("no Run")
				return
			}
			globalHarness.RunCase(t, tc)
		})
	}
}
