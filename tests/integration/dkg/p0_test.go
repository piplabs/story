//go:build integration

package dkg

import (
	"context"
	"testing"
	"time"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

func P0Cases() []TestCase {
	return []TestCase{
		{
			ID:             "IT-E2E-01",
			Priority:       "P0",
			Description:    "Full first round happy path: Registration → Dealing → Finalization → Active; GlobalPublicKey set",
			Expected:       "stage=Active, GlobalPublicKey non-empty, round completes all 4 stages",
			Run:            runIT_E2E_01,
			NeedsRoundWait: true,
		},
	}
}

func runIT_E2E_01(t *testing.T, h *Harness) {
	ctx := context.Background()

	// Refresh the round — GetLatestDKGNetwork may return a stale round.
	// Poll until we get a round that's still in progress or active.
	net, err := h.GetLatestDKGNetwork(ctx)
	if err != nil {
		t.Fatalf("GetLatestDKGNetwork: %v", err)
	}
	if net == nil {
		t.Fatal("no DKG round (ensure height >= dkgStartBlock)")
	}

	round := net.Round
	stage := net.Stage
	t.Logf("DKG state: round=%d stage=%s isResharing=%v", round, stage, net.IsResharing)

	// 如果已经在 Active 或更后面的阶段，直接检查并返回
	if stage >= dkgtypes.DKGStageActive {
		t.Logf("round %d already at stage=%s, skipping wait", round, stage)
		assertActiveWithGlobalKey(t, h)
		return
	}

	// 若在 Registration，等待 Verified >= 2（refresh round each poll to avoid stale round)
	if stage == dkgtypes.DKGStageRegistration {
		deadline := time.Now().Add(h.MaxWait)
		found := false
		for time.Now().Before(deadline) {
			cur, _ := h.GetLatestDKGNetwork(ctx)
			if cur == nil {
				time.Sleep(h.PollInterval)
				continue
			}
			if cur.Round != round {
				round = cur.Round
				stage = cur.Stage
				t.Logf("round advanced to %d stage=%s", round, stage)
			}
			if cur.Stage > dkgtypes.DKGStageRegistration {
				found = true
				break
			}
			regs, _ := h.GetVerifiedRegistrations(ctx, cur.Round)
			if len(regs) >= 2 {
				found = true
				break
			}
			time.Sleep(h.PollInterval)
		}
		if !found {
			t.Fatal("timeout waiting for Verified >= 2")
		}
	}

	// 使用 WaitForRoundStage 追踪特定 round，避免被新 round 覆盖
	// 等待 Dealing
	if stage <= dkgtypes.DKGStageRegistration {
		if !h.WaitForRoundStage(ctx, round, dkgtypes.DKGStageDealing) {
			t.Fatal("timeout waiting for Stage=Dealing")
		}
	}

	// 等待 Finalization
	if stage <= dkgtypes.DKGStageDealing {
		if !h.WaitForRoundStage(ctx, round, dkgtypes.DKGStageFinalization) {
			t.Fatal("timeout waiting for Stage=Finalization")
		}
	}

	// 等待 Active
	if !h.WaitForRoundStage(ctx, round, dkgtypes.DKGStageActive) {
		t.Fatal("timeout waiting for Stage=Active")
	}

	assertActiveWithGlobalKey(t, h)
}


func TestDKG_P0(t *testing.T) {
	for _, tc := range P0Cases() {
		t.Run(tc.ID, func(t *testing.T) {
			if ScenarioNameForCase(tc.ID) == "" && tc.SkipIfLive != "" {
				t.Skip(tc.SkipIfLive)
			}
			globalHarness.RunCase(t, tc)
		})
	}
}
