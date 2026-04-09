package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

type msgServer struct {
	*Keeper
	types.UnimplementedMsgServiceServer
}

// AddVotes is called with all aggregated votes included in a new finalized block.
func (s msgServer) AddVote(ctx context.Context, msg *types.MsgAddDkgVote,
) (*types.AddDkgVoteResponse, error) {
	if msg.Authority != s.Keeper.GetAuthority() {
		return nil, errors.New("unauthorized")
	}

	s.RemoveBroadcastedVotes(msg.Vote)

	latestRound, err := s.GetLatestDKGRound(ctx)
	if err != nil {
		return nil, err
	}

	if latestRound != nil && latestRound.Stage == types.DKGStageDealing {
		if len(msg.Vote.Deals) > 0 {
			if err := s.ProcessDeals(ctx, latestRound, msg.Vote.Deals); err != nil {
				// Note: no need to return error since no state changes in processing deals
				log.Error(ctx, "Error occurred while processing deals", err)
			}
		}

		if len(msg.Vote.Responses) > 0 {
			if err := s.ProcessResponses(ctx, latestRound, msg.Vote.Responses); err != nil {
				// Note: no need to return error since no state changes in processing responses
				log.Error(ctx, "Error occurred while processing responses", err)
			}
		}

		if len(msg.Vote.Justifications) > 0 {
			if err := s.ProcessJustifications(ctx, latestRound, msg.Vote.Justifications); err != nil {
				// Note: no need to return error since no state changes in processing justifications
				log.Error(ctx, "Error occurred while processing justifications", err)
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
