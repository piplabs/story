package keeper

import (
	"context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const (
	// maxVoteExtensionSize is the maximum allowed size of a raw vote extension in bytes (256 KB).
	maxVoteExtensionSize = 256 << 10

	// maxItemsPerVote is the maximum number of deals, responses, or justifications
	// allowed in a single vote extension. With 80 validators, a single validator
	// produces at most 79 deals and 79 responses per round.
	maxItemsPerVote = 80
)

func (k *Keeper) ExtendVote(_ sdk.Context, _ *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
	dequeuedDeals := k.DequeueDeals(maxItemsPerVote)
	dequeuedResponses := k.DequeueResponses(maxItemsPerVote)
	dequeuedJustifications := k.DequeueJustifications(maxItemsPerVote)

	bz, err := proto.Marshal(&types.Vote{
		Deals:          dequeuedDeals,
		Responses:      dequeuedResponses,
		Justifications: dequeuedJustifications,
	})
	if err != nil {
		return nil, errors.Wrap(err, "marshal vote")
	}

	return &abci.ResponseExtendVote{
		VoteExtension: bz,
	}, nil
}

func (k *Keeper) VerifyVoteExtension(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
	// Reject malformed vote extensions via ABCI status (not Go error),
	// as returning a Go error is treated as an application bug by CometBFT.
	_, _, err := k.parseAndVerifyVoteExtension(req.VoteExtension)
	if err != nil {
		log.Warn(ctx, "Rejecting malformed vote extension", err)
		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_REJECT}, nil
	}

	return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
}

//nolint:unparam // ignore unused param error
func (*Keeper) parseAndVerifyVoteExtension(voteExt []byte) ([]*types.Vote, bool, error) {
	if len(voteExt) > maxVoteExtensionSize {
		return nil, false, fmt.Errorf("vote extension too large: %d bytes exceeds max %d", len(voteExt), maxVoteExtensionSize)
	}

	vote, ok, err := votesFromExtension(voteExt)
	if err != nil {
		return nil, false, errors.Wrap(err, "parse vote extension")
	} else if !ok {
		return nil, true, nil // Empty vote extension is fine
	}

	if len(vote.Deals) > maxItemsPerVote {
		return nil, false, fmt.Errorf("too many deals in vote extension: %d exceeds max %d", len(vote.Deals), maxItemsPerVote)
	}

	if len(vote.Responses) > maxItemsPerVote {
		return nil, false, fmt.Errorf("too many responses in vote extension: %d exceeds max %d", len(vote.Responses), maxItemsPerVote)
	}

	if len(vote.Justifications) > maxItemsPerVote {
		return nil, false, fmt.Errorf("too many justifications in vote extension: %d exceeds max %d", len(vote.Justifications), maxItemsPerVote)
	}

	return []*types.Vote{vote}, true, nil
}

// PrepareVotes returns the cosmosSDK transaction MsgAddVotes that will include all the validator votes included
// in the previous block's vote extensions into the attest module.
//
// Note that the commit is assumed to be valid and only contains valid VEs from the previous block as
// provided by a trusted cometBFT. Some votes (contained inside VE) may however be invalid, they are discarded.
func (k *Keeper) PrepareVotes(ctx context.Context, commit abci.ExtendedCommitInfo, commitHeight uint64) (sdk.Msg, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Vote extensions become available in LocalLastCommit one block after the
	// VoteExtensionsEnableHeight. At the enable height itself, LocalLastCommit
	// contains votes from the previous height which have no VE data.
	// Return an empty MsgAddDkgVote until VEs are actually present.
	cp := sdkCtx.ConsensusParams()
	veHeight := int64(0)
	if cp.Abci != nil {
		veHeight = cp.Abci.VoteExtensionsEnableHeight
	}
	if veHeight == 0 || sdkCtx.BlockHeight() <= veHeight {
		return &types.MsgAddDkgVote{
			Authority: k.GetAuthority(),
			Vote:      &types.Vote{},
		}, nil
	}

	// The VEs in LastLocalCommit is expected to be valid
	if err := baseapp.ValidateVoteExtensions(sdkCtx, k.valStore, 0, "", commit); err != nil {
		return nil, errors.Wrap(err, "validate extensions [BUG]")
	}

	// Verify and discard invalid votes.
	// Votes inside the VEs are NOT guaranteed to be valid, since
	// VerifyVoteExtension isn't called after quorum is reached.
	var allVotes []*types.Vote

	log.Info(ctx, "Processing vote extensions", "height", commitHeight, "num_votes", len(commit.Votes))

	for _, vote := range commit.Votes {
		selected, _, err := k.parseAndVerifyVoteExtension(vote.VoteExtension)
		if err != nil {
			log.Warn(ctx, "Discarding invalid vote extension", err, log.Hex7("validator", vote.Validator.Address))
			continue
		}

		allVotes = append(allVotes, selected...)
	}

	votes := aggregateVotes(allVotes)

	return &types.MsgAddDkgVote{
		Authority: k.GetAuthority(),
		Vote:      votes,
	}, nil
}

func aggregateVotes(votes []*types.Vote) *types.Vote {
	allDeals := make([]types.Deal, 0)
	allResponses := make([]types.Response, 0)
	allJustifications := make([]types.Justification, 0)
	for _, vote := range votes {
		allDeals = append(allDeals, vote.Deals...)
		allResponses = append(allResponses, vote.Responses...)
		allJustifications = append(allJustifications, vote.Justifications...)
	}

	return &types.Vote{
		Deals:          deduplicateDeals(allDeals),
		Responses:      deduplicateResponses(allResponses),
		Justifications: deduplicateJustifications(allJustifications),
	}
}

// deduplicateDeals removes duplicate deals by (dealerIndex, recipientIndex).
func deduplicateDeals(deals []types.Deal) []types.Deal {
	type dedupKey struct {
		dealerIndex    uint32
		recipientIndex uint32
	}

	seen := make(map[dedupKey]struct{})
	result := make([]types.Deal, 0, len(deals))

	for _, d := range deals {
		key := dedupKey{dealerIndex: d.Index, recipientIndex: d.RecipientIndex}
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}

		result = append(result, d)
	}

	return result
}

// deduplicateResponses removes duplicate responses by (responderIndex, dealerIndex).
func deduplicateResponses(responses []types.Response) []types.Response {
	type dedupKey struct {
		responderIndex uint32
		dealerIndex    uint32
	}

	seen := make(map[dedupKey]struct{})
	result := make([]types.Response, 0, len(responses))

	for _, r := range responses {
		var dealerIdx uint32
		if r.VssResponse != nil {
			dealerIdx = r.VssResponse.Index
		}

		key := dedupKey{responderIndex: r.Index, dealerIndex: dealerIdx}
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}

		result = append(result, r)
	}

	return result
}

// deduplicateJustifications removes duplicate justifications by (dealerIndex, recipientIndex).
// When multiple validators broadcast the same justification, only the first is processed.
func deduplicateJustifications(justifications []types.Justification) []types.Justification {
	type dedupKey struct {
		dealerIndex    uint32
		recipientIndex uint32
	}

	seen := make(map[dedupKey]struct{})
	result := make([]types.Justification, 0, len(justifications))

	for _, j := range justifications {
		var recipientIdx uint32

		if vssJ := j.GetVssJustification(); vssJ != nil {
			if pd := vssJ.GetPlainDeal(); pd != nil {
				if ss := pd.GetSecShare(); ss != nil {
					recipientIdx = ss.GetI()
				}
			}
		}

		key := dedupKey{dealerIndex: j.Index, recipientIndex: recipientIdx}
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}

		result = append(result, j)
	}

	return result
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
