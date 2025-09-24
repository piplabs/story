package events

import (
	"encoding/json"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"
)

func TestParseBeginInitializationEvent(t *testing.T) {
	eventData := map[string]interface{}{
		"mrenclave":         "initmrenclave",
		"round":             float64(1),
		"start_block":       float64(100),
		"active_validators": []interface{}{"val1", "val2"},
	}
	eventDataBytes, _ := json.Marshal(eventData)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginInitialization",
		Attributes: []abcitypes.EventAttribute{
			{Key: "data", Value: string(eventDataBytes)},
		},
	}
	height := int64(10)
	result := parseBeginInitializationEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_initialization", result.EventType)
	require.Equal(t, "initmrenclave", result.Mrenclave)
	require.Equal(t, uint32(1), result.Round)
	require.Equal(t, height, result.BlockHeight)
	require.ElementsMatch(t, []string{"val1", "val2"}, result.ActiveValidators)
}

func TestParseBeginNetworkSetEvent(t *testing.T) {
	eventData := map[string]interface{}{
		"mrenclave": "netmrenclave",
		"round":     float64(2),
		"total":     float64(5),
		"threshold": float64(3),
	}
	eventDataBytes, _ := json.Marshal(eventData)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginNetworkSet",
		Attributes: []abcitypes.EventAttribute{
			{Key: "data", Value: string(eventDataBytes)},
		},
	}
	height := int64(20)
	result := parseBeginNetworkSetEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_network_set", result.EventType)
	require.Equal(t, "netmrenclave", result.Mrenclave)
	require.Equal(t, uint32(2), result.Round)
	require.Equal(t, uint32(5), result.Total)
	require.Equal(t, uint32(3), result.Threshold)
	require.Equal(t, height, result.BlockHeight)
}

func TestParseBeginDealingEvent(t *testing.T) {
	eventData := map[string]interface{}{
		"mrenclave": "dealermrenclave",
		"round":     float64(3),
	}
	eventDataBytes, _ := json.Marshal(eventData)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginDealing",
		Attributes: []abcitypes.EventAttribute{
			{Key: "data", Value: string(eventDataBytes)},
		},
	}
	height := int64(30)
	result := parseBeginDealingEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_dealing", result.EventType)
	require.Equal(t, "dealermrenclave", result.Mrenclave)
	require.Equal(t, uint32(3), result.Round)
	require.Equal(t, height, result.BlockHeight)
}

func TestParseBeginFinalizationEvent(t *testing.T) {
	eventData := map[string]interface{}{
		"mrenclave": "finalizemrenclave",
		"round":     float64(4),
	}
	eventDataBytes, _ := json.Marshal(eventData)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginFinalization",
		Attributes: []abcitypes.EventAttribute{
			{Key: "data", Value: string(eventDataBytes)},
		},
	}
	height := int64(40)
	result := parseBeginFinalizationEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_finalization", result.EventType)
	require.Equal(t, "finalizemrenclave", result.Mrenclave)
	require.Equal(t, uint32(4), result.Round)
	require.Equal(t, height, result.BlockHeight)
}

func TestParseDKGFinalizedEvent(t *testing.T) {
	eventData := map[string]interface{}{
		"mrenclave": "finalizedmrenclave",
		"round":     float64(5),
	}
	eventDataBytes, _ := json.Marshal(eventData)
	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventDKGFinalized",
		Attributes: []abcitypes.EventAttribute{
			{Key: "data", Value: string(eventDataBytes)},
		},
	}
	height := int64(50)
	result := parseDKGFinalizedEvent(event, height)
	require.NotNil(t, result)
	require.Equal(t, "dkg_finalized", result.EventType)
	require.Equal(t, "finalizedmrenclave", result.Mrenclave)
	require.Equal(t, uint32(5), result.Round)
	require.Equal(t, height, result.BlockHeight)
}

func TestParseBeginProcessResponsesEvent(t *testing.T) {
	eventData := map[string]interface{}{
		"mrenclave": "testmrenclave",
		"round":     float64(5),
		"responses": []interface{}{
			map[string]interface{}{
				"index": float64(1),
				"vss_response": map[string]interface{}{
					"session_id": "sessid1",
					"index":      float64(2),
					"status":     true,
					"signature":  "sig1",
				},
			},
			map[string]interface{}{
				"index": float64(3),
				"vss_response": map[string]interface{}{
					"session_id": "sessid2",
					"index":      float64(4),
					"status":     false,
					"signature":  "sig2",
				},
			},
		},
	}
	eventDataBytes, _ := json.Marshal(eventData)

	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginProcessResponses",
		Attributes: []abcitypes.EventAttribute{
			{Key: "data", Value: string(eventDataBytes)},
		},
	}

	height := int64(456)
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
	require.Equal(t, true, resp1.VssResponse.Status)
	require.Equal(t, []byte("sig1"), resp1.VssResponse.Signature)

	resp2 := result.Responses[1]
	require.Equal(t, uint32(3), resp2.Index)
	require.NotNil(t, resp2.VssResponse)
	require.Equal(t, []byte("sessid2"), resp2.VssResponse.SessionId)
	require.Equal(t, uint32(4), resp2.VssResponse.Index)
	require.Equal(t, false, resp2.VssResponse.Status)
	require.Equal(t, []byte("sig2"), resp2.VssResponse.Signature)
}

func TestParseBeginProcessDealsEvent(t *testing.T) {
	eventData := map[string]interface{}{
		"mrenclave": "testmrenclave",
		"round":     float64(3),
		"deals": []interface{}{
			map[string]interface{}{
				"index":           float64(1),
				"recipient_index": float64(2),
				"deal": []interface{}{
					map[string]interface{}{"key": "dh_key", "value": "dhkeybytes"},
					map[string]interface{}{"key": "signature", "value": "sigbytes"},
					map[string]interface{}{"key": "nonce", "value": "noncebytes"},
					map[string]interface{}{"key": "cipher", "value": "cipherbytes"},
				},
				"signature": "dealsig",
			},
		},
	}
	eventDataBytes, _ := json.Marshal(eventData)

	event := abcitypes.Event{
		Type: "/story.dkg.v1.types.EventBeginProcessDeals",
		Attributes: []abcitypes.EventAttribute{
			{Key: "data", Value: string(eventDataBytes)},
		},
	}

	height := int64(123)
	result := parseBeginProcessDealsEvent(event, height)

	require.NotNil(t, result)
	require.Equal(t, "dkg_begin_process_deals", result.EventType)
	require.Equal(t, "testmrenclave", result.Mrenclave)
	require.Equal(t, uint32(3), result.Round)
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
