package keeper

import (
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
