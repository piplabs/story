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
			if err := s.Keeper.ProcessDeals(ctx, latestRound, msg.Vote.Deals); err != nil {
				// Note: no need to return error since no state changes in processing deals
				log.Error(ctx, "Error occurred while processing deals", err)
			}
		}

		if len(msg.Vote.Responses) > 0 {
			if err := s.Keeper.ProcessResponses(ctx, latestRound, msg.Vote.Responses); err != nil {
				// Note: no need to return error since no state changes in processing responses
				log.Error(ctx, "Error occurred while processing responses", err)
			}
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
