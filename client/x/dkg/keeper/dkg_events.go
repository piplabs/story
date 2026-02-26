package keeper

import (
	"context"
	"encoding/hex"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

//
// events for stages (ie. events that drive the stage transitions)
//

func (*Keeper) emitBeginDKGInitialization(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginInitialization{
		CodeCommitment:   dkgNetwork.CodeCommitment,
		Round:            dkgNetwork.Round,
		ActiveValidators: dkgNetwork.ActiveValSet,
		StartBlockHeight: uint32(dkgNetwork.StartBlockHeight),
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_initialization event")
	}

	return nil
}

func (*Keeper) emitBeginDKGDealing(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginDealing{
		CodeCommitment: dkgNetwork.CodeCommitment,
		Round:          dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_dealing event")
	}

	log.Info(ctx, "Emitted BeginDKGDealing event", "round", dkgNetwork.Round)

	return nil
}

func (*Keeper) emitBeginProcessDeals(ctx context.Context, dkgNetwork *types.DKGNetwork, deals []types.Deal) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginProcessDeals{
		CodeCommitment: dkgNetwork.CodeCommitment,
		Round:          dkgNetwork.Round,
		Deals:          deals,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_process_deals event")
	}

	log.Info(ctx, "Emitted BeginProcessDeals event", "round", dkgNetwork.Round, "code_commitment", hex.EncodeToString(dkgNetwork.CodeCommitment), "num_deals", len(deals))

	return nil
}

func (*Keeper) emitBeginProcessResponses(ctx context.Context, dkgNetwork *types.DKGNetwork, responses []types.Response) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginProcessResponses{
		CodeCommitment: dkgNetwork.CodeCommitment,
		Round:          dkgNetwork.Round,
		Responses:      responses,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_process_responses event")
	}

	log.Info(ctx, "Emitted BeginProcessResponses event", "round", dkgNetwork.Round, "code_commitment", hex.EncodeToString(dkgNetwork.CodeCommitment), "num_responses", len(responses))

	return nil
}

func (*Keeper) emitBeginDKGFinalization(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginFinalization{
		CodeCommitment: dkgNetwork.CodeCommitment,
		Round:          dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_finalization event")
	}

	log.Info(ctx, "Emitted BeginDKGFinalization event", "round", dkgNetwork.Round, "code_commitment", hex.EncodeToString(dkgNetwork.CodeCommitment))

	return nil
}

func (*Keeper) emitDKGFinalized(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventDKGFinalized{
		CodeCommitment: dkgNetwork.CodeCommitment,
		Round:          dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_finalized_event")
	}

	log.Info(ctx, "Emitted DKGFinalized event", "round", dkgNetwork.Round, "code_commitment", hex.EncodeToString(dkgNetwork.CodeCommitment))

	return nil
}

func (*Keeper) emitDKGCommitteeRewarded(ctx context.Context, dkgNetwork *types.DKGNetwork, memberCount uint32, totalReward, perMemberReward, remainingUbi math.Int) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventDKGCommitteeRewarded{
		CodeCommitment:  dkgNetwork.CodeCommitment,
		Round:           dkgNetwork.Round,
		MemberCount:     memberCount,
		TotalReward:     totalReward.String(),
		PerMemberReward: perMemberReward.String(),
		RemainingUbi:    remainingUbi.String(),
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_committee_rewarded event")
	}

	log.Info(ctx, "Emitted DKGCommitteeRewarded event",
		"round", dkgNetwork.Round,
		"member_count", memberCount,
		"total_reward", totalReward.String(),
	)

	return nil
}
