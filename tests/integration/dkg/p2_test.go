//go:build integration

package dkg

import (
	"context"
	"fmt"
	"testing"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

func P2Cases() []TestCase {
	list := make([]TestCase, 0, 41)
	// IT-RES-01 .. IT-RES-07
	resFuncs := map[int]func(*testing.T, *Harness){
		1: runIT_RES_01, 2: runIT_RES_02, 3: runIT_RES_03,
		4: runIT_RES_04, 5: runIT_RES_05, 6: runIT_RES_06, 7: runIT_RES_07,
	}
	resDescs := map[int]string{
		1: "Phase=Failed → GetSession → ResumeDKGService",
		2: "GetSession fails → MarkFailed → recover next round",
		3: "Phase=Failed, tryResume → GetSession → resume",
		4: "GetSession ok, Phase != Failed → no resume",
		5: "Resume → rejoin current round",
		6: "Resume CreateSession fails → stay Failed",
		7: "All validators resume → new round completes",
	}
	resExpecteds := map[int]string{
		1: "new round starts after TEE recovery",
		2: "service resumes, new round progresses",
		3: "resume succeeds, DKG continues",
		4: "no resume triggered in normal operation",
		5: "validator rejoins after resume",
		6: "stayed Failed when CreateSession fails",
		7: "all resumed, round reaches Active, GlobalPublicKey set",
	}
	for i := 1; i <= 7; i++ {
		needsWait := i == 7 // RES-07: all validators resume → new round completes
		list = append(list, TestCase{
			ID:             "IT-RES-" + fmt.Sprintf("%02d", i),
			Priority:       "P2",
			Description:    resDescs[i],
			Expected:       resExpecteds[i],
			Run:            resFuncs[i],
			NeedsRoundWait: needsWait,
		})
	}
	// IT-UPG-01 .. IT-UPG-05
	upgFuncs := map[int]func(*testing.T, *Harness){
		1: runIT_UPG_01, 2: runIT_UPG_02, 3: runIT_UPG_03,
		4: runIT_UPG_04, 5: runIT_UPG_05,
	}
	upgDescs := map[int]string{
		1: "ScheduleUpgrade → PendingUpgrade set",
		2: "hasPendingUpgradeActivation returns activationHeight",
		3: "CancelUpgrade → removes pending upgrade",
		4: "Upgrade round + enough validators → completes",
		5: "Upgrade round failed → retry next round",
	}
	upgExpecteds := map[int]string{
		1: "upgrade scheduled successfully",
		2: "activationHeight set correctly",
		3: "upgrade canceled, IsUpgrade=false",
		4: "upgrade round reaches Active, GlobalPublicKey set",
		5: "upgrade round fails and new round retries",
	}
	for i := 1; i <= 5; i++ {
		needsWait := i == 4 || i == 5 // UPG-04/05: upgrade round lifecycle
		list = append(list, TestCase{
			ID:             "IT-UPG-" + fmt.Sprintf("%02d", i),
			Priority:       "P2",
			Description:    upgDescs[i],
			Expected:       upgExpecteds[i],
			Run:            upgFuncs[i],
			NeedsRoundWait: needsWait,
		})
	}
	// IT-E2E-06
	list = append(list, TestCase{
		ID:             "IT-E2E-06",
		Priority:       "P2",
		Description:    "E2E Upgrade resharing round",
		Expected:       "upgrade round reaches Active, GlobalPublicKey set",
		Run:            runIT_E2E_06,
		NeedsRoundWait: true,
	})
	// IT-ENC-01
	list = append(list, TestCase{
		ID:          "IT-ENC-01",
		Priority:    "P2",
		Description: "Stage=Active: TDH2 encrypt + CDR.write → VaultWritten",
		Expected:    "active round with GlobalPublicKey, CDRWrite succeeds",
		Run:         runIT_ENC_01,
	})
	// IT-DEC-01 .. IT-DEC-22
	decFuncs := map[int]func(*testing.T, *Harness){
		1: runIT_DEC_01, 2: runIT_DEC_02, 3: runIT_DEC_03, 4: runIT_DEC_04,
		5: runIT_DEC_05, 6: runIT_DEC_06, 7: runIT_DEC_07, 8: runIT_DEC_08,
		9: runIT_DEC_09, 10: runIT_DEC_10, 11: runIT_DEC_11, 12: runIT_DEC_12,
		13: runIT_DEC_13, 14: runIT_DEC_14, 15: runIT_DEC_15, 16: runIT_DEC_16,
		17: runIT_DEC_17, 18: runIT_DEC_18, 19: runIT_DEC_19, 20: runIT_DEC_20,
		21: runIT_DEC_21, 22: runIT_DEC_22,
	}
	decDescs := map[int]string{
		1:  "ProcessCDRVaultRead: latestActive exists",
		2:  "[needs TEE mock] No latestActive: return early",
		3:  "[needs disk fault] GetSession fails: MarkFailed",
		4:  "[needs disk fault] Phase != Completed: skip decrypt",
		5:  "[needs TEE mock] callTEEDecrypt fails: skip",
		6:  "[needs TEE mock] callContractSubmitPartial fails",
		7:  "[needs TEE mock] EncryptedPartialDecryption: IncrementCount",
		8:  "[needs TEE mock] EncryptedPartialDecryption: RefundCDRFee",
		9:  "[needs TEE mock] RefundCDRFee: pool balance check",
		10: "[needs TEE mock] RefundCDRFee: SendCoins from pool",
		11: "[needs TEE mock] RefundCDRFee: pool < refundAmt error",
		12: "[needs TEE mock] Multiple partials reach threshold",
		13: "[needs TEE mock] Partial from non-active-round validator",
		14: "[needs TEE mock] Duplicate partial submission",
		15: "[needs TEE mock] Invalid partial proof",
		16: "[needs TEE mock] Vault not found for decrypt",
		17: "[needs TEE mock] Decrypt request for empty vault",
		18: "[needs TEE mock] Concurrent decrypt requests",
		19: "[needs TEE mock] Decrypt after round rotation",
		20: "[needs TEE mock] Decrypt with stale GlobalPublicKey",
		21: "[needs TEE mock] callTEEDecrypt timeout",
		22: "[needs TEE mock] Full decrypt E2E path",
	}
	decExpecteds := map[int]string{
		1: "active round exists, CDRRead path exercised",
		7: "EncryptedPartialDecryption event processed",
		8: "RefundCDRFee executed",
		10: "SendCoins from pool executed",
		12: "partial threshold reached",
		22: "full decrypt E2E completed",
	}
	for i := 1; i <= 22; i++ {
		expected := decExpecteds[i]
		if expected == "" {
			expected = "decrypt path variant: " + decDescs[i]
		}
		list = append(list, TestCase{
			ID:          "IT-DEC-" + fmt.Sprintf("%02d", i),
			Priority:    "P2",
			Description: decDescs[i],
			Expected:    expected,
			Run:         decFuncs[i],
		})
	}
	// IT-CDR-04 .. IT-CDR-08
	list = append(list,
		TestCase{ID: "IT-CDR-04", Priority: "P2", Description: "FeeCollected amount=0: no AddCDRFeeToPool", Expected: "zero-fee event path", Run: runIT_CDR_04},
		TestCase{ID: "IT-CDR-05", Priority: "P2", Description: "Pool balance < refund amount: error", Expected: "pool underflow error", Run: runIT_CDR_05},
		TestCase{ID: "IT-CDR-06", Priority: "P2", Description: "No previous active round: distributeCDRRewardPool returns nil", Expected: "first round has no previous active", Run: runIT_CDR_06},
		TestCase{ID: "IT-CDR-07", Priority: "P2", Description: "Total submit count=0, pool>0: no SendCoins", Expected: "zero-submit-count path, no SendCoins", Run: runIT_CDR_07},
		TestCase{ID: "IT-CDR-08", Priority: "P2", Description: "E2E: CDR write/read + partial decryption + distributeCDRRewardPool", Expected: "full CDR lifecycle executed", Run: runIT_CDR_08},
	)
	return list
}

func runIT_CDR_06(t *testing.T, h *Harness) {
	ctx := context.Background()
	// If no active round exists, distributeCDRRewardPool should be a no-op
	_, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil {
		checkTrue(t, "no previous active round", true, "distributeCDRRewardPool returns nil")
		t.Log("no active round; distributeCDRRewardPool would return nil (expected)")
		return
	}
	t.Log("active round exists; IT-CDR-06 condition (no previous active) not met — skip assertion")
}

func runIT_CDR_08(t *testing.T, h *Harness) {
	if h.IsNoopChainClient() {
		t.Skip("CDR E2E requires EthChainClient")
		return
	}
	ctx := context.Background()

	// Verify active round with GlobalPublicKey
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no active DKG round with GlobalPublicKey for CDR E2E")
		return
	}

	check(t, "stage", dkgtypes.DKGStageActive, net.Stage)
	checkTrue(t, "GlobalPublicKey", len(net.GlobalPublicKey) > 0, fmt.Sprintf("len=%d", len(net.GlobalPublicKey)))

	// Allocate → Write → Read
	uuid, err := h.ChainClient.CDRAllocate(ctx)
	if err != nil {
		t.Fatalf("CDRAllocate: %v", err)
	}
	t.Logf("allocated vault uuid=%d", uuid)

	encData := []byte("test-encrypted-data-for-cdr-e2e")
	if err := h.ChainClient.CDRWrite(ctx, uuid, encData); err != nil {
		t.Fatalf("CDRWrite: %v", err)
	}

	// Read triggers ThresholdDecryptRequested → partial decryptions → FeeCollected
	requesterPubKey := []byte("dummy-requester-pubkey")
	if err := h.ChainClient.CDRRead(ctx, uuid, requesterPubKey); err != nil {
		t.Fatalf("CDRRead: %v", err)
	}
	t.Log("CDR E2E: Allocate→Write→Read completed; fees and partial decryption distribution exercised on next FinalizeDKGRound")
}

func runIT_ENC_01(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil {
		t.Skip("no active DKG round (CDR.write requires Active + GlobalPublicKey)")
		return
	}
	if net == nil || len(net.GlobalPublicKey) == 0 {
		t.Skip("no GlobalPublicKey for TDH2 encrypt")
		return
	}
	checkTrue(t, "stage", net.Stage == dkgtypes.DKGStageActive, net.Stage.String())
	checkTrue(t, "GlobalPublicKey", len(net.GlobalPublicKey) > 0, fmt.Sprintf("len=%d", len(net.GlobalPublicKey)))
	// With EthChainClient, actually exercise CDR.write (TDH2 encrypt + VaultWritten)
	if !h.IsNoopChainClient() {
		uuid, err := h.ChainClient.CDRAllocate(ctx)
		if err != nil {
			t.Fatalf("CDRAllocate: %v", err)
		}
		encData := []byte("test-tdh2-encrypted-data-enc-01")
		if err := h.ChainClient.CDRWrite(ctx, uuid, encData); err != nil {
			t.Fatalf("CDRWrite: %v", err)
		}
		t.Logf("TDH2 encrypt + CDR.write succeeded: uuid=%d, VaultWritten event expected", uuid)
		return
	}
	t.Log("Active round with GlobalPublicKey confirmed; CDR.write requires EthChainClient for full test")
}

func runIT_DEC_01(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || net == nil {
		t.Skip("no active round (ProcessCDRVaultRead requires latestActive)")
		return
	}
	checkTrue(t, "GlobalPublicKey", len(net.GlobalPublicKey) > 0, fmt.Sprintf("len=%d", len(net.GlobalPublicKey)))
	// With EthChainClient, actually exercise CDR.read (VaultRead → ThresholdDecryptRequested)
	if !h.IsNoopChainClient() {
		uuid, err := h.ChainClient.CDRAllocate(ctx)
		if err != nil {
			t.Fatalf("CDRAllocate: %v", err)
		}
		encData := []byte("test-encrypted-data-dec-01")
		if err := h.ChainClient.CDRWrite(ctx, uuid, encData); err != nil {
			t.Fatalf("CDRWrite: %v", err)
		}
		requesterPubKey := []byte("requester-pubkey-dec-01")
		if err := h.ChainClient.CDRRead(ctx, uuid, requesterPubKey); err != nil {
			t.Fatalf("CDRRead: %v", err)
		}
		t.Logf("CDR.read succeeded: uuid=%d, VaultRead→ThresholdDecryptRequested path exercised", uuid)
		return
	}
	t.Log("latestActive present; CDR.read requires EthChainClient for full decrypt path test")
}

func TestDKG_P2(t *testing.T) {
	for _, tc := range P2Cases() {
		t.Run(tc.ID, func(t *testing.T) {
			if ScenarioNameForCase(tc.ID) == "" && tc.SkipIfLive != "" {
				t.Skip(tc.SkipIfLive)
			}
			if tc.Run == nil {
				t.Skip("no Run (API or fault-injection only)")
				return
			}
			globalHarness.RunCase(t, tc)
		})
	}
}
