package keeper

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
	dkgservice "github.com/piplabs/story/client/dkg/service"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func (k *Keeper) ExtendVote(ctx sdk.Context, _ *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
	deals := dkgservice.Deals
	bz, err := proto.Marshal(&types.Vote{
		Deals: deals,
	})
	if err != nil {
		return nil, errors.Wrap(err, "marshal vote")
	}
	return &abci.ResponseExtendVote{
		VoteExtension: bz,
	}, nil
}

func (k *Keeper) VerifyVoteExtension(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
	// todo: consider adding more checks here
	_, _, err := k.parseAndVerifyVoteExtension(ctx, req.VoteExtension)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse vote extension")
	}
	return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
}

func (k *Keeper) parseAndVerifyVoteExtension(
	ctx sdk.Context,
	voteExt []byte,
) ([]*types.Vote, bool, error) {
	vote, ok, err := votesFromExtension(voteExt)
	if err != nil {
		return nil, false, errors.Wrap(err, "parse vote extension")
	} else if !ok {
		return nil, true, nil // Empty vote extension is fine
	} else {
		return []*types.Vote{vote}, true, nil
	}
}

// PrepareVotes returns the cosmosSDK transaction MsgAddVotes that will include all the validator votes included
// in the previous block's vote extensions into the attest module.
//
// Note that the commit is assumed to be valid and only contains valid VEs from the previous block as
// provided by a trusted cometBFT. Some votes (contained inside VE) may however be invalid, they are discarded.
func (k *Keeper) PrepareVotes(ctx context.Context, commit abci.ExtendedCommitInfo, commitHeight uint64) ([]sdk.Msg, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// The VEs in LastLocalCommit is expected to be valid
	if err := baseapp.ValidateVoteExtensions(sdkCtx, k.skeeper, 0, "", commit); err != nil {
		return nil, errors.Wrap(err, "validate extensions [BUG]")
	}

	// Verify and discard invalid votes.
	// Votes inside the VEs are NOT guaranteed to be valid, since
	// VerifyVoteExtension isn't called after quorum is reached.
	var allVotes []*types.Vote
	log.Info(ctx, "Processing vote extensions", "height", commitHeight, "num_votes", len(commit.Votes))
	for _, vote := range commit.Votes {
		selected, _, err := k.parseAndVerifyVoteExtension(sdkCtx, vote.VoteExtension) //nolint:contextcheck // sdkCtx passed
		if err != nil {
			log.Warn(ctx, "Discarding invalid vote extension", err, log.Hex7("validator", vote.Validator.Address))
			continue
		}

		allVotes = append(allVotes, selected...)
	}

	votes, err := aggregateVotes(allVotes)
	if err != nil {
		return nil, err
	}

	return []sdk.Msg{&types.MsgAddDkgVote{
		Authority: authtypes.NewModuleAddress(types.ModuleName).String(),
		Vote:      votes,
	}}, nil
}

func aggregateVotes(votes []*types.Vote) (*types.Vote, error) {
	dealMap := make([]*types.Deal, 0)
	for _, vote := range votes {
		dealMap = append(dealMap, vote.Deals...)
	}
	return &types.Vote{
		Deals: dealMap,
	}, nil
}

// votesFromExtension returns the attestations contained in the vote extension, or false if none or an error.
func votesFromExtension(voteExtension []byte) (*types.Vote, bool, error) {
	if len(voteExtension) == 0 {
		return nil, false, nil
	}

	resp := new(types.Vote)
	if err := proto.Unmarshal(voteExtension, resp); err != nil {
		return nil, false, errors.Wrap(err, "decode vote extension")
	}

	return resp, true, nil
}
