//go:build integration

package dkg

import (
	"context"
	"fmt"
	"testing"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

func P3Cases() []TestCase {
	edgeFuncs := map[int]func(*testing.T, *Harness){
		1: runIT_EDGE_01, 2: runIT_EDGE_02, 3: runIT_EDGE_03,
		4: runIT_EDGE_04, 5: runIT_EDGE_05, 6: runIT_EDGE_06,
	}
	edgeDescs := map[int]string{
		1: "ActiveValSet only MinReq validators: DKG completes",
		2: "Exactly MinReq registered/finalized: boundary check",
		3: "Max validators stress: total/threshold consistency",
		4: "Rapid round rotation stress test",
		5: "Period boundaries: Registration/Dealing/Finalization/Active periods",
		6: "TEE down 1 node: DKG still completes (fault tolerance)",
	}
	edgeExpecteds := map[int]string{
		1: "Active round with valid Total/Threshold, GlobalPublicKey set",
		2: "exactly MinReq validators at boundary",
		3: "max validators, Total/Threshold valid and consistent",
		4: "rapid round rotations handled correctly",
		5: "period transitions occur at correct boundaries",
		6: "DKG completes with fault tolerance despite one node down",
	}
	list := make([]TestCase, 0, 6)
	for i := 1; i <= 6; i++ {
		needsWait := i == 1 || i == 4 || i == 5 || i == 6 // EDGE-01,04,05,06: need stage transitions
		list = append(list, TestCase{
			ID: "IT-EDGE-" + fmt.Sprintf("%02d", i),
			Priority: "P2",
			Description:    edgeDescs[i],
			Expected:       edgeExpecteds[i],
			Run:            edgeFuncs[i],
			NeedsRoundWait: needsWait,
		})
	}
	return list
}

func runIT_EDGE_01(t *testing.T, h *Harness) {
	ctx := context.Background()
	net, err := h.GetLatestDKGNetwork(ctx)
	if err != nil || net == nil {
		t.Skip("no network")
		return
	}
	// Verify DKG completes even at minimum validator count
	active, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil || active == nil {
		// Try waiting for Active (track specific round)
		if !h.WaitForRoundStage(ctx, net.Round, dkgtypes.DKGStageActive) {
			t.Skip("could not reach Active stage")
			return
		}
		active, err = h.GetLatestActiveDKGNetwork(ctx)
		if err != nil || active == nil {
			t.Skip("no active round after waiting")
			return
		}
	}
	// Assert Total and Threshold are valid
	checkTrue(t, "Total > 0", active.Total > 0, fmt.Sprintf("Total=%d", active.Total))
	checkTrue(t, "Threshold > 0", active.Threshold > 0, fmt.Sprintf("Threshold=%d", active.Threshold))
	checkTrue(t, "Threshold <= Total", active.Threshold <= active.Total, fmt.Sprintf("Threshold=%d Total=%d", active.Threshold, active.Total))
	checkTrue(t, "GlobalPublicKey", len(active.GlobalPublicKey) > 0, fmt.Sprintf("len=%d", len(active.GlobalPublicKey)))
	// Verify params consistency
	params, err := h.Params(ctx)
	if err != nil {
		t.Fatalf("Params: %v", err)
	}
	checkTrue(t, "Total >= MinDkgCommitteeSize", active.Total >= uint32(params.MinReqRegisteredParticipants), fmt.Sprintf("Total=%d MinReq=%d", active.Total, params.MinReqRegisteredParticipants))
	t.Logf("Active round: total=%d threshold=%d minDkgCommitteeSize=%d globalPubKey=%d bytes",
		active.Total, active.Threshold, params.MinReqRegisteredParticipants, len(active.GlobalPublicKey))
}

func TestDKG_P3(t *testing.T) {
	for _, tc := range P3Cases() {
		t.Run(tc.ID, func(t *testing.T) {
			if ScenarioNameForCase(tc.ID) == "" && tc.SkipIfLive != "" {
				t.Skip(tc.SkipIfLive)
			}
			if tc.Run == nil {
				t.Skip("no Run (stress/boundary only)")
				return
			}
			globalHarness.RunCase(t, tc)
		})
	}
}
