//nolint:testpackage // fix this later to events_test
package events

import (
	"encoding/json"
	"strconv"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

func TestParseBeginInitializationEvent(t *testing.T) {
	mrenclave := "initmrenclave"
	activeValidators := []string{"val1", "val2"}
	activeValidatorsBytes, err := json.Marshal(activeValidators)
	require.NoError(t, err)
	height := int64(10)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginInitialization",
		Attributes: []abcitypes.EventAttribute{
			{Key: "mrenclave", Value: mrenclave},
			{Key: "round", Value: "1"},
			{Key: "start_block", Value: strconv.FormatInt(height, 10)},
			{Key: "active_validators", Value: string(activeValidatorsBytes)},
		},
	}

	result := parseBeginInitializationEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_initialization", result.EventType)
	require.Equal(t, mrenclave, result.Mrenclave)
	require.Equal(t, uint32(1), result.Round)
	require.Equal(t, height, result.BlockHeight)
	require.ElementsMatch(t, []string{"val1", "val2"}, result.ActiveValidators)
}

func TestParseBeginNetworkSetEvent(t *testing.T) {
	height := int64(20)
	mrenclave := "netmrenclave"
	total := uint32(5)
	threshold := uint32(3)
	round := uint32(2)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginNetworkSet",
		Attributes: []abcitypes.EventAttribute{
			{Key: "mrenclave", Value: mrenclave},
			{Key: "round", Value: strconv.FormatUint(uint64(round), 10)},
			{Key: "total", Value: strconv.FormatUint(uint64(total), 10)},
			{Key: "threshold", Value: strconv.FormatUint(uint64(threshold), 10)},
		},
	}

	result := parseBeginNetworkSetEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_network_set", result.EventType)
	require.Equal(t, mrenclave, result.Mrenclave)
	require.Equal(t, round, result.Round)
	require.Equal(t, total, result.Total)
	require.Equal(t, threshold, result.Threshold)
	require.Equal(t, height, result.BlockHeight)
}

func TestParseBeginDealingEvent(t *testing.T) {
	height := int64(30)
	mrenclave := "dealermrenclave"
	round := uint32(3)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginDealing",
		Attributes: []abcitypes.EventAttribute{
			{Key: "mrenclave", Value: mrenclave},
			{Key: "round", Value: strconv.FormatUint(uint64(round), 10)},
		},
	}

	result := parseBeginDealingEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_dealing", result.EventType)
	require.Equal(t, mrenclave, result.Mrenclave)
	require.Equal(t, round, result.Round)
	require.Equal(t, height, result.BlockHeight)
}

func TestParseBeginFinalizationEvent(t *testing.T) {
	height := int64(40)
	mrenclave := "finalizemrenclave"
	round := uint32(4)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginFinalization",
		Attributes: []abcitypes.EventAttribute{
			{Key: "mrenclave", Value: mrenclave},
			{Key: "round", Value: strconv.FormatUint(uint64(round), 10)},
		},
	}

	result := parseBeginFinalizationEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_finalization", result.EventType)
	require.Equal(t, mrenclave, result.Mrenclave)
	require.Equal(t, round, result.Round)
	require.Equal(t, height, result.BlockHeight)
}

func TestParseDKGFinalizedEvent(t *testing.T) {
	height := int64(50)
	mrenclave := "finalizedmrenclave"
	round := uint32(5)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventDKGFinalized",
		Attributes: []abcitypes.EventAttribute{
			{Key: "mrenclave", Value: mrenclave},
			{Key: "round", Value: strconv.FormatUint(uint64(round), 10)},
		},
	}

	result := parseDKGFinalizedEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_finalized", result.EventType)
	require.Equal(t, mrenclave, result.Mrenclave)
	require.Equal(t, round, result.Round)
	require.Equal(t, height, result.BlockHeight)
}

func TestParseBeginProcessDealsEvent(t *testing.T) {
	deals := []*dkgtypes.Deal{
		{
			Index:          1,
			RecipientIndex: 2,
			Deal: dkgtypes.EncryptedDeal{
				DhKey:     []byte("dhkeybytes"),
				Signature: []byte("sigbytes"),
				Nonce:     []byte("noncebytes"),
				Cipher:    []byte("cipherbytes"),
			},
			Signature: []byte("dealsig"),
		},
	}
	dealsBytes, err := json.Marshal(deals)
	require.NoError(t, err)

	height := int64(123)
	mrenclave := "testmrenclave"
	round := uint32(3)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginProcessDeals",
		Attributes: []abcitypes.EventAttribute{
			{Key: "mrenclave", Value: mrenclave},
			{Key: "round", Value: strconv.FormatUint(uint64(round), 10)},
			{Key: "deals", Value: string(dealsBytes)},
		},
	}

	result := parseBeginProcessDealsEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_process_deals", result.EventType)
	require.Equal(t, mrenclave, result.Mrenclave)
	require.Equal(t, round, result.Round)
	require.Equal(t, height, result.BlockHeight)
	require.Len(t, result.Deals, 1)

	deal := result.Deals[0]
	require.Equal(t, uint32(1), deal.Index)
	require.Equal(t, uint32(2), deal.RecipientIndex)
	require.Equal(t, []byte("dhkeybytes"), deal.Deal.DhKey)
	require.Equal(t, []byte("sigbytes"), deal.Deal.Signature)
	require.Equal(t, []byte("noncebytes"), deal.Deal.Nonce)
	require.Equal(t, []byte("cipherbytes"), deal.Deal.Cipher)
	require.Equal(t, []byte("dealsig"), deal.Signature)
}

func TestParseBeginProcessResponsesEvent(t *testing.T) {
	responses := []*dkgtypes.Response{
		{
			Index: 1,
			VssResponse: &dkgtypes.VSSResponse{
				SessionId: []byte("sessid1"),
				Index:     2,
				Status:    true,
				Signature: []byte("sig1"),
			},
		},
		{
			Index: 3,
			VssResponse: &dkgtypes.VSSResponse{
				SessionId: []byte("sessid2"),
				Index:     4,
				Status:    false,
				Signature: []byte("sig2"),
			},
		},
	}
	responsesBytes, err := json.Marshal(responses)
	require.NoError(t, err)

	height := int64(456)
	mrenclave := "testmrenclave"
	round := uint32(5)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginProcessResponses",
		Attributes: []abcitypes.EventAttribute{
			{Key: "mrenclave", Value: mrenclave},
			{Key: "round", Value: strconv.FormatUint(uint64(round), 10)},
			{Key: "responses", Value: string(responsesBytes)},
		},
	}

	result := parseBeginProcessResponsesEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_process_responses", result.EventType)
	require.Equal(t, "testmrenclave", result.Mrenclave)
	require.Equal(t, uint32(5), result.Round)
	require.Equal(t, height, result.BlockHeight)
	require.Len(t, result.Responses, 2)

	resp1 := result.Responses[0]
	require.Equal(t, uint32(1), resp1.Index)
	require.NotNil(t, resp1.VssResponse)
	require.Equal(t, []byte("sessid1"), resp1.VssResponse.SessionId)
	require.Equal(t, uint32(2), resp1.VssResponse.Index)
	require.True(t, resp1.VssResponse.Status)
	require.Equal(t, []byte("sig1"), resp1.VssResponse.Signature)

	resp2 := result.Responses[1]
	require.Equal(t, uint32(3), resp2.Index)
	require.NotNil(t, resp2.VssResponse)
	require.Equal(t, []byte("sessid2"), resp2.VssResponse.SessionId)
	require.Equal(t, uint32(4), resp2.VssResponse.Index)
	require.False(t, resp2.VssResponse.Status)
	require.Equal(t, []byte("sig2"), resp2.VssResponse.Signature)
}
