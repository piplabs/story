package keeper

import (
	"github.com/cosmos/gogoproto/proto"

	"github.com/piplabs/story/client/x/dkg/types"
)

// EnqueueDeals adds multiple deals to the queue in a thread-safe manner.
func (*Keeper) EnqueueDeals(newDeals []types.Deal) {
	dealsMu.Lock()
	defer dealsMu.Unlock()

	deals = append(deals, newDeals...)
}

// DequeueDeals dequeues up to count deals in a thread-safe manner.
func (*Keeper) DequeueDeals(count int) []types.Deal {
	dealsMu.Lock()
	defer dealsMu.Unlock()

	if len(deals) == 0 {
		return nil
	}

	if count > len(deals) {
		count = len(deals)
	}

	out := make([]types.Deal, count)
	copy(out, deals[:count])
	deals = deals[count:]

	return out
}

// PeekDeals returns up to count deals without removing them.
func (*Keeper) PeekDeals(count int) []types.Deal {
	dealsMu.Lock()
	defer dealsMu.Unlock()

	if len(deals) == 0 {
		return nil
	}

	if count > len(deals) {
		count = len(deals)
	}

	out := make([]types.Deal, count)
	copy(out, deals[:count])

	return out
}

// EnqueueResponses adds multiple responses to the queue in a thread-safe manner.
func (*Keeper) EnqueueResponses(newResponses []types.Response) {
	responsesMu.Lock()
	defer responsesMu.Unlock()

	responses = append(responses, newResponses...)
}

// DequeueResponses dequeues up to count responses in a thread-safe manner.
func (*Keeper) DequeueResponses(count int) []types.Response {
	responsesMu.Lock()
	defer responsesMu.Unlock()

	if len(responses) == 0 {
		return nil
	}

	if count > len(responses) {
		count = len(responses)
	}

	out := make([]types.Response, count)
	copy(out, responses[:count])
	responses = responses[count:]

	return out
}

// PeekResponses returns up to count responses without removing them.
func (*Keeper) PeekResponses(count int) []types.Response {
	responsesMu.Lock()
	defer responsesMu.Unlock()

	if len(responses) == 0 {
		return nil
	}

	if count > len(responses) {
		count = len(responses)
	}

	out := make([]types.Response, count)
	copy(out, responses[:count])

	return out
}

// EnqueueJustifications adds multiple justifications to the queue in a thread-safe manner.
func (*Keeper) EnqueueJustifications(newJustifications []*types.Justification) {
	justificationsMu.Lock()
	defer justificationsMu.Unlock()

	for _, j := range newJustifications {
		if j != nil {
			justifications = append(justifications, *j)
		}
	}
}

// DequeueJustifications dequeues up to count justifications in a thread-safe manner.
func (*Keeper) DequeueJustifications(count int) []types.Justification {
	justificationsMu.Lock()
	defer justificationsMu.Unlock()

	if len(justifications) == 0 {
		return nil
	}

	if count > len(justifications) {
		count = len(justifications)
	}

	out := make([]types.Justification, count)
	copy(out, justifications[:count])
	justifications = justifications[count:]

	return out
}

// PeekJustifications returns up to count justifications without removing them.
func (*Keeper) PeekJustifications(count int) []types.Justification {
	justificationsMu.Lock()
	defer justificationsMu.Unlock()

	if len(justifications) == 0 {
		return nil
	}

	if count > len(justifications) {
		count = len(justifications)
	}

	out := make([]types.Justification, count)
	copy(out, justifications[:count])

	return out
}

// RemoveBroadcastedVotes removes any queued items that were included in the block's vote aggregation.
func (*Keeper) RemoveBroadcastedVotes(vote *types.Vote) {
	if vote == nil {
		return
	}

	if len(vote.Deals) > 0 {
		removeDeals(vote.Deals)
	}

	if len(vote.Responses) > 0 {
		removeResponses(vote.Responses)
	}

	if len(vote.Justifications) > 0 {
		removeJustifications(vote.Justifications)
	}
}

func removeDeals(included []types.Deal) {
	dealsMu.Lock()
	defer dealsMu.Unlock()

	if len(deals) == 0 {
		return
	}

	keys := make(map[string]struct{}, len(included))
	for _, d := range included {
		if key, ok := voteItemKey(&d); ok {
			keys[key] = struct{}{}
		}
	}

	filtered := make([]types.Deal, 0, len(deals))
	for _, d := range deals {
		if key, ok := voteItemKey(&d); ok {
			if _, exists := keys[key]; exists {
				continue
			}
		}

		filtered = append(filtered, d)
	}

	deals = filtered
}

func removeResponses(included []types.Response) {
	responsesMu.Lock()
	defer responsesMu.Unlock()

	if len(responses) == 0 {
		return
	}
	keys := make(map[string]struct{}, len(included))
	for _, r := range included {
		if key, ok := voteItemKey(&r); ok {
			keys[key] = struct{}{}
		}
	}

	filtered := make([]types.Response, 0, len(responses))
	for _, r := range responses {
		if key, ok := voteItemKey(&r); ok {
			if _, exists := keys[key]; exists {
				continue
			}
		}

		filtered = append(filtered, r)
	}

	responses = filtered
}

func removeJustifications(included []types.Justification) {
	justificationsMu.Lock()
	defer justificationsMu.Unlock()

	if len(justifications) == 0 {
		return
	}
	keys := make(map[string]struct{}, len(included))
	for _, j := range included {
		if key, ok := voteItemKey(&j); ok {
			keys[key] = struct{}{}
		}
	}

	filtered := make([]types.Justification, 0, len(justifications))
	for _, j := range justifications {
		if key, ok := voteItemKey(&j); ok {
			if _, exists := keys[key]; exists {
				continue
			}
		}

		filtered = append(filtered, j)
	}

	justifications = filtered
}

func voteItemKey(msg proto.Message) (string, bool) {
	bz, err := proto.Marshal(msg)
	if err != nil {
		return "", false
	}

	return string(bz), true
}

// FlushAllQueues clears all deal, response, justification, and pending incoming queues.
// This should be called when a DKG round transitions to prevent stale data
// from a previous round from being broadcast in the new round.
func (*Keeper) FlushAllQueues() {
	dealsMu.Lock()

	deals = nil

	dealsMu.Unlock()

	responsesMu.Lock()

	responses = nil

	responsesMu.Unlock()

	justificationsMu.Lock()

	justifications = nil

	justificationsMu.Unlock()

	// Also flush pending incoming data from failed kernel calls.
	flushPendingIncoming()
}
