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
//
// This handler is intentionally permissive — it only checks Authority
// and returns success without validating vote content. This follows the
// ABCI++ pattern where ProcessProposal is permissive (accepts all well-formed
// proposals) and FinalizeBlock performs the actual validation via msg_server.AddVote.
func (s proposalServer) AddVote(ctx context.Context, msg *types.MsgAddDkgVote,
) (*types.AddDkgVoteResponse, error) {
	if msg.Authority != s.Keeper.GetAuthority() {
		return nil, errors.New("unauthorized")
	}

	return &types.AddDkgVoteResponse{}, nil
}

// NewProposalServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewProposalServer(keeper *Keeper) types.MsgServiceServer {
	return &proposalServer{Keeper: keeper}
}

var _ types.MsgServiceServer = proposalServer{}
