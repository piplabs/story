package keeper

import (
	"context"
	"github.com/piplabs/story/lib/log"

	"github.com/piplabs/story/client/x/dkg/types"
)

type msgServer struct {
	*Keeper
	types.UnimplementedMsgServiceServer
}

// AddVotes is called with all aggregated votes included in a new finalized block.
func (s msgServer) AddVote(ctx context.Context, msg *types.MsgAddDkgVote,
) (*types.AddDkgVoteResponse, error) {
	latestRound, err := s.Keeper.GetLatestDKGRound(ctx)
	if err != nil {
		return nil, err
	}

	if latestRound != nil && latestRound.Stage == types.DKGStageDealing {
		if len(msg.Vote.Deals) > 0 {
			log.Info(ctx, "there are deals to process", "len_deals", len(msg.Vote.Deals))
			_ = s.Keeper.emitBeginProcessDeals(ctx, latestRound, msg.Vote.Deals)
		}

		if len(msg.Vote.Responses) > 0 {
			log.Info(ctx, "there are responses to process", "len_responses", len(msg.Vote.Responses))
			_ = s.Keeper.emitBeginProcessResponses(ctx, latestRound, msg.Vote.Responses)
		}
	}

	return &types.AddDkgVoteResponse{}, nil
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServiceServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServiceServer = msgServer{}
