package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/dkg/types"
)

type msgServer struct {
	*Keeper
	types.UnimplementedMsgServiceServer
}

// AddVotes is called with all aggregated votes included in a new finalized block.
func (msgServer) AddVote(_ context.Context, _ *types.MsgAddDkgVote,
) (*types.AddDkgVoteResponse, error) {
	return &types.AddDkgVoteResponse{}, nil
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServiceServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServiceServer = msgServer{}
