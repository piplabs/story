package events

import (
	"encoding/json"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/dkg/types"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

// isDKGEvent checks if an event is a DKG-related event.
func isDKGEvent(eventType string) bool {
	// Check for typed proto events from x/dkg module
	dkgEventTypes := []string{
		"/story.dkg.v1.types.EventBeginInitialization",
		"/story.dkg.v1.types.EventBeginChallengePeriod",
		"/story.dkg.v1.types.EventBeginNetworkSet",
		"/story.dkg.v1.types.EventBeginDealing",
		"/story.dkg.v1.types.EventBeginFinalization",
		"/story.dkg.v1.types.EventDKGFinalized",
		"/story.dkg.v1.types.EventDKGRegistrationInitialized",
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
	case "/story.dkg.v1.types.EventBeginChallengePeriod":
		return parseBeginChallengePeriodEvent(event, height)
	case "/story.dkg.v1.types.EventBeginNetworkSet":
		return parseBeginNetworkSetEvent(event, height)
	case "/story.dkg.v1.types.EventBeginDealing":
		return parseBeginDealingEvent(event, height)
	case "/story.dkg.v1.types.EventBeginFinalization":
		return parseBeginFinalizationEvent(event, height)
	case "/story.dkg.v1.types.EventDKGFinalized":
		return parseDKGFinalizedEvent(event, height)
	case "/story.dkg.v1.types.EventDKGRegistrationInitialized":
		return parseDKGRegistrationInitializedEvent(event, height)
	case "/story.dkg.v1.types.EventDKGRegistrationCommitmentsUpdated":
		return parseDKGRegistrationCommitmentsUpdatedEvent(event, height)
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
			if totalFloat, ok := eventData["total"].(float64); ok {
				protoEvent.Total = uint32(totalFloat)
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
		Total:            protoEvent.Total,
		Attributes:       extractAttributes(event),
	}
}

// parseBeginChallengePeriodEvent parses EventBeginChallengePeriod typed event.
func parseBeginChallengePeriodEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventBeginChallengePeriod

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
		EventType:   "dkg_begin_challenge_period",
		Mrenclave:   string(protoEvent.Mrenclave),
		Round:       protoEvent.Round,
		BlockHeight: height,
		Attributes:  extractAttributes(event),
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

// parseDKGRegistrationInitializedEvent parses EventDKGRegistrationInitialized typed event.
func parseDKGRegistrationInitializedEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventDKGRegistrationInitialized
	var msgSenderStr string

	for _, attr := range event.Attributes {
		if attr.Key == "data" {
			var eventData map[string]any
			if err := json.Unmarshal([]byte(attr.Value), &eventData); err != nil {
				break
			}

			if msgSender, ok := eventData["msgSender"].(common.Address); ok {
				msgSenderStr = msgSender.Hex()
			}
			if mrenclaveStr, ok := eventData["mrenclave"].(string); ok {
				protoEvent.Mrenclave = []byte(mrenclaveStr)
			}
			if roundFloat, ok := eventData["round"].(float64); ok {
				protoEvent.Round = uint32(roundFloat)
			}
			if indexFloat, ok := eventData["index"].(float64); ok {
				protoEvent.Index = uint32(indexFloat)
			}
			if dkgPubKeyStr, ok := eventData["dkgPubKey"].(string); ok {
				protoEvent.DkgPubKey = []byte(dkgPubKeyStr)
			}
			if commPubKeyStr, ok := eventData["commPubKey"].(string); ok {
				protoEvent.CommPubKey = []byte(commPubKeyStr)
			}
			if rawQuoteStr, ok := eventData["rawQuote"].(string); ok {
				protoEvent.RawQuote = []byte(rawQuoteStr)
			}

			break
		}
	}

	return &types.DKGEventData{
		EventType:     "dkg_registration_initialized",
		Mrenclave:     string(protoEvent.Mrenclave),
		Round:         protoEvent.Round,
		ValidatorAddr: msgSenderStr,
		Index:         protoEvent.Index,
		BlockHeight:   height,
		DkgPubKey:     protoEvent.DkgPubKey,
		CommPubKey:    protoEvent.CommPubKey,
		RawQuote:      protoEvent.RawQuote,
		Attributes:    extractAttributes(event),
	}
}

func parseDKGRegistrationCommitmentsUpdatedEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventDKGRegistrationCommitmentsUpdated
	var msgSenderStr string

	for _, attr := range event.Attributes {
		if attr.Key == "data" {
			var eventData map[string]any
			if err := json.Unmarshal([]byte(attr.Value), &eventData); err != nil {
				break
			}

			if msgSender, ok := eventData["msgSender"].(common.Address); ok {
				msgSenderStr = msgSender.Hex()
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
			if indexFloat, ok := eventData["index"].(float64); ok {
				protoEvent.Index = uint32(indexFloat)
			}
			if commitmentsStr, ok := eventData["commitments"].(string); ok {
				protoEvent.Commitments = []byte(commitmentsStr)
			}
			if signatureStr, ok := eventData["signature"].(string); ok {
				protoEvent.Signature = []byte(signatureStr)
			}

			break
		}
	}

	return &types.DKGEventData{
		EventType:     "dkg_registration_commitments_updated",
		Mrenclave:     string(protoEvent.Mrenclave),
		Round:         protoEvent.Round,
		ValidatorAddr: msgSenderStr,
		Total:         protoEvent.Total,
		Threshold:     protoEvent.Threshold,
		Index:         protoEvent.Index,
		BlockHeight:   height,
		Commitments:   protoEvent.Commitments,
		Signature:     protoEvent.Signature,
		Attributes:    extractAttributes(event),
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
