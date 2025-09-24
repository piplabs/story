package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/dkg/types"
)

type proposalServer struct {
	*Keeper
	types.UnimplementedMsgServiceServer
}

// AddVotes verifies all aggregated votes included in a proposed block.
func (s proposalServer) AddVote(ctx context.Context, msg *types.MsgAddDkgVote,
) (*types.AddDkgVoteResponse, error) {
	latestRound, err := s.Keeper.GetLatestDKGRound(ctx)
	if err != nil {
		return nil, err
	}

	if latestRound != nil && latestRound.Stage == types.DKGStageDealing {
		_ = s.Keeper.emitBeginProcessDeal(ctx, latestRound, msg.Vote.Deals)
		_ = s.Keeper.emitBeginProcessResponses(ctx, latestRound, msg.Vote.Responses)
	}

	return &types.AddDkgVoteResponse{}, nil
}

// NewProposalServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewProposalServer(keeper *Keeper) types.MsgServiceServer {
	return &proposalServer{Keeper: keeper}
}

var _ types.MsgServiceServer = proposalServer{}
