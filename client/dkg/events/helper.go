package events

import (
	"encoding/json"

	abcitypes "github.com/cometbft/cometbft/abci/types"

	"github.com/piplabs/story/client/dkg/types"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

// isDKGEvent checks if an event is a DKG-related event.
func isDKGEvent(eventType string) bool {
	// Check for typed proto events from x/dkg module
	dkgEventTypes := []string{
		"/story.dkg.v1.types.EventBeginInitialization",
		"/story.dkg.v1.types.EventBeginNetworkSet",
		"/story.dkg.v1.types.EventBeginDealing",
		"/story.dkg.v1.types.EventBeginProcessDeals",
		"/story.dkg.v1.types.EventBeginProcessResponses",
		"/story.dkg.v1.types.EventBeginFinalization",
		"/story.dkg.v1.types.EventDKGFinalized",
	}

	for _, validType := range dkgEventTypes {
		if eventType == validType {
			return true
		}
	}

	return false
}

// parseEvent parses a blockchain event into DKG event data.
func (*EventListener) parseEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	// Handle typed proto events from x/dkg module
	switch event.Type {
	case "/story.dkg.v1.types.EventBeginInitialization":
		return parseBeginInitializationEvent(event, height)
	case "/story.dkg.v1.types.EventBeginNetworkSet":
		return parseBeginNetworkSetEvent(event, height)
	case "/story.dkg.v1.types.EventBeginDealing":
		return parseBeginDealingEvent(event, height)
	case "/story.dkg.v1.types.EventBeginProcessDeals":
		return parseBeginProcessDealsEvent(event, height)
	case "/story.dkg.v1.types.EventBeginProcessResponses":
		return parseBeginProcessResponsesEvent(event, height)
	case "/story.dkg.v1.types.EventBeginFinalization":
		return parseBeginFinalizationEvent(event, height)
	case "/story.dkg.v1.types.EventDKGFinalized":
		return parseDKGFinalizedEvent(event, height)
	default:
		return nil
	}
}

// parseBeginInitializationEvent parses EventBeginInitialization typed event.
func parseBeginInitializationEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventBeginInitialization

	for _, attr := range event.Attributes {
		if attr.Key == "data" {
			var eventData map[string]any
			if err := json.Unmarshal([]byte(attr.Value), &eventData); err != nil {
				break
			}

			if mrenclaveStr, ok := eventData["mrenclave"].(string); ok {
				protoEvent.Mrenclave = []byte(mrenclaveStr)
			}
			if roundFloat, ok := eventData["round"].(float64); ok {
				protoEvent.Round = uint32(roundFloat)
			}
			if startBlockFloat, ok := eventData["start_block"].(float64); ok {
				//nolint:govet // keep field
				protoEvent.StartBlock = uint32(startBlockFloat)
			}
			if activeValidators, ok := eventData["active_validators"].([]string); ok {
				protoEvent.ActiveValidators = append(protoEvent.ActiveValidators, activeValidators...)
			}

			break
		}
	}

	return &types.DKGEventData{
		EventType:        "dkg_begin_initialization",
		Mrenclave:        string(protoEvent.Mrenclave),
		Round:            protoEvent.Round,
		BlockHeight:      height,
		ActiveValidators: protoEvent.ActiveValidators,
		Attributes:       extractAttributes(event),
	}
}

// parseBeginNetworkSetEvent parses EventBeginNetworkSet typed event.
func parseBeginNetworkSetEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventBeginNetworkSet

	for _, attr := range event.Attributes {
		//nolint:nestif // ignore nestedif linting
		if attr.Key == "data" {
			var eventData map[string]any
			if err := json.Unmarshal([]byte(attr.Value), &eventData); err != nil {
				break
			}

			if mrenclaveStr, ok := eventData["mrenclave"].(string); ok {
				protoEvent.Mrenclave = []byte(mrenclaveStr)
			}
			if roundFloat, ok := eventData["round"].(float64); ok {
				protoEvent.Round = uint32(roundFloat)
			}
			if totalFloat, ok := eventData["total"].(float64); ok {
				protoEvent.Total = uint32(totalFloat)
			}
			if thresholdFloat, ok := eventData["threshold"].(float64); ok {
				protoEvent.Threshold = uint32(thresholdFloat)
			}
			break
		}
	}

	return &types.DKGEventData{
		EventType:   "dkg_begin_network_set",
		Mrenclave:   string(protoEvent.Mrenclave),
		Round:       protoEvent.Round,
		BlockHeight: height,
		Total:       protoEvent.Total,
		Threshold:   protoEvent.Threshold,
		Attributes:  extractAttributes(event),
	}
}

// parseBeginDealingEvent parses EventBeginDealing typed event.

func parseBeginDealingEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventBeginDealing

	for _, attr := range event.Attributes {
		if attr.Key == "data" {
			var eventData map[string]any
			if err := json.Unmarshal([]byte(attr.Value), &eventData); err != nil {
				break
			}

			if mrenclaveStr, ok := eventData["mrenclave"].(string); ok {
				protoEvent.Mrenclave = []byte(mrenclaveStr)
			}
			if roundFloat, ok := eventData["round"].(float64); ok {
				protoEvent.Round = uint32(roundFloat)
			}

			break
		}
	}

	return &types.DKGEventData{
		EventType:   "dkg_begin_dealing",
		Mrenclave:   string(protoEvent.Mrenclave),
		Round:       protoEvent.Round,
		BlockHeight: height,
		Attributes:  extractAttributes(event),
	}
}

func parseBeginProcessDealsEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventBeginProcessDeals

	for _, attr := range event.Attributes {
		if attr.Key == "data" {
			var eventData map[string]any
			if err := json.Unmarshal([]byte(attr.Value), &eventData); err != nil {
				break
			}
			if mrenclaveStr, ok := eventData["mrenclave"].(string); ok {
				protoEvent.Mrenclave = []byte(mrenclaveStr)
			}
			if roundFloat, ok := eventData["round"].(float64); ok {
				protoEvent.Round = uint32(roundFloat)
			}
			if deals, ok := eventData["deals"].([]map[string]any); ok {
				for _, deal := range deals {
					var dkgDeal dkgtypes.Deal
					if indexFloat, ok := deal["index"].(float64); ok {
						dkgDeal.Index = uint32(indexFloat)
					}
					if recipientIndex, ok := deal["recipient_index"].(float64); ok {
						dkgDeal.RecipientIndex = uint32(recipientIndex)
					}
					if encryptedDeal, ok := deal["deal"].([]map[string]any); ok {
						dkgDeal.Deal = dkgtypes.EncryptedDeal{}
						for _, edAttr := range encryptedDeal {
							if edKey, ok := edAttr["key"].(string); ok {
								switch edKey {
								case "dh_key":
									if dhKey, ok := edAttr["value"].(string); ok {
										dkgDeal.Deal.DhKey = []byte(dhKey)
									}
								case "signature":
									if signature, ok := edAttr["value"].(string); ok {
										dkgDeal.Deal.Signature = []byte(signature)
									}
								case "nonce":
									if nonce, ok := edAttr["value"].(string); ok {
										dkgDeal.Deal.Nonce = []byte(nonce)
									}
								case "cipher":
									if cipher, ok := edAttr["value"].(string); ok {
										dkgDeal.Deal.Cipher = []byte(cipher)
									}
								}
							}
						}
					}

					if signature, ok := deal["signature"].(string); ok {
						dkgDeal.Signature = []byte(signature)
					}
					protoEvent.Deals = append(protoEvent.Deals, &dkgDeal)
				}
			}

			break
		}
	}

	return &types.DKGEventData{
		EventType:   "dkg_begin_process_deals",
		Mrenclave:   string(protoEvent.Mrenclave),
		Round:       protoEvent.Round,
		BlockHeight: height,
		Deals:       protoEvent.Deals,
		Attributes:  extractAttributes(event),
	}
}

func parseBeginProcessResponsesEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventBeginProcessResponses

	for _, attr := range event.Attributes {
		if attr.Key == "data" {
			var eventData map[string]any
			if err := json.Unmarshal([]byte(attr.Value), &eventData); err != nil {
				break
			}

			if mrenclaveStr, ok := eventData["mrenclave"].(string); ok {
				protoEvent.Mrenclave = []byte(mrenclaveStr)
			}
			if roundFloat, ok := eventData["round"].(float64); ok {
				protoEvent.Round = uint32(roundFloat)
			}
			if responses, ok := eventData["responses"].([]map[string]any); ok {
				for _, response := range responses {
					var dkgResponse dkgtypes.Response
					if indexFloat, ok := response["index"].(float64); ok {
						dkgResponse.Index = uint32(indexFloat)
					}
					if vssResponse, ok := response["vss_response"].(map[string]any); ok {
						dkgResponse.VssResponse = &dkgtypes.VSSResponse{}
						if sessionID, ok := vssResponse["session_id"].(string); ok {
							dkgResponse.VssResponse.SessionId = []byte(sessionID)
						}
						if index, ok := vssResponse["index"].(float64); ok {
							dkgResponse.VssResponse.Index = uint32(index)
						}
						if status, ok := vssResponse["status"].(bool); ok {
							dkgResponse.VssResponse.Status = status
						}
						if signature, ok := vssResponse["signature"].(string); ok {
							dkgResponse.VssResponse.Signature = []byte(signature)
						}
					}
					protoEvent.Responses = append(protoEvent.Responses, &dkgResponse)
				}
			}
			break
		}
	}

	return &types.DKGEventData{
		EventType:   "dkg_begin_process_responses",
		Mrenclave:   string(protoEvent.Mrenclave),
		Round:       protoEvent.Round,
		BlockHeight: height,
		Responses:   protoEvent.Responses,
		Attributes:  extractAttributes(event),
	}
}

// parseBeginFinalizationEvent parses EventBeginFinalization typed event.
func parseBeginFinalizationEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventBeginFinalization

	for _, attr := range event.Attributes {
		if attr.Key == "data" {
			var eventData map[string]any
			if err := json.Unmarshal([]byte(attr.Value), &eventData); err != nil {
				break
			}

			if mrenclaveStr, ok := eventData["mrenclave"].(string); ok {
				protoEvent.Mrenclave = []byte(mrenclaveStr)
			}
			if roundFloat, ok := eventData["round"].(float64); ok {
				protoEvent.Round = uint32(roundFloat)
			}

			break
		}
	}

	return &types.DKGEventData{
		EventType:   "dkg_begin_finalization",
		Mrenclave:   string(protoEvent.Mrenclave),
		Round:       protoEvent.Round,
		BlockHeight: height,
		Attributes:  extractAttributes(event),
	}
}

// parseDKGFinalizedEvent parses EventDKGFinalized typed event.
func parseDKGFinalizedEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventDKGFinalized

	for _, attr := range event.Attributes {
		if attr.Key == "data" {
			var eventData map[string]any
			if err := json.Unmarshal([]byte(attr.Value), &eventData); err != nil {
				break
			}

			if mrenclaveStr, ok := eventData["mrenclave"].(string); ok {
				protoEvent.Mrenclave = []byte(mrenclaveStr)
			}
			if roundFloat, ok := eventData["round"].(float64); ok {
				protoEvent.Round = uint32(roundFloat)
			}

			break
		}
	}

	return &types.DKGEventData{
		EventType:   "dkg_finalized",
		Mrenclave:   string(protoEvent.Mrenclave),
		Round:       protoEvent.Round,
		BlockHeight: height,
		Attributes:  extractAttributes(event),
	}
}

// extractAttributes extracts all attributes from an event as a map.
func extractAttributes(event abcitypes.Event) map[string]string {
	attributes := make(map[string]string)
	for _, attr := range event.Attributes {
		attributes[attr.Key] = attr.Value
	}

	return attributes
}
