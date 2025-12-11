package keeper

import (
	"context"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

type proposalServer struct {
	*Keeper
	types.UnimplementedMsgServiceServer
}

// AddVotes verifies all aggregated votes included in a proposed block.
func (s proposalServer) AddVote(ctx context.Context, msg *types.MsgAddDkgVote,
) (*types.AddDkgVoteResponse, error) {
	// TODO: add verification of deals and responses

	if s.isDKGSvcEnabled {
		latestRound, err := s.GetLatestDKGRound(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get latest DKG round")
		}

		if latestRound != nil {
			s.ResumeDKGService(ctx, latestRound)
		}
	}

	return &types.AddDkgVoteResponse{}, nil
}

// NewProposalServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewProposalServer(keeper *Keeper) types.MsgServiceServer {
	return &proposalServer{Keeper: keeper}
}

var _ types.MsgServiceServer = proposalServer{}
