package events

import (
	"encoding/json"
	"strconv"

	abcitypes "github.com/cometbft/cometbft/abci/types"

	"github.com/piplabs/story/client/dkg/types"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

// isDKGEvent checks if an event is a DKG-related event.
func isDKGEvent(eventType string) bool {
	// Check for typed proto events from x/dkg module
	dkgEventTypes := []string{
		"story.dkg.v1.types.EventBeginInitialization",
		"story.dkg.v1.types.EventBeginNetworkSet",
		"story.dkg.v1.types.EventBeginDealing",
		"story.dkg.v1.types.EventBeginProcessDeals",
		"story.dkg.v1.types.EventBeginProcessResponses",
		"story.dkg.v1.types.EventBeginFinalization",
		"story.dkg.v1.types.EventDKGFinalized",
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
	case "story.dkg.v1.types.EventBeginInitialization":
		return parseBeginInitializationEvent(event, height)
	case "story.dkg.v1.types.EventBeginNetworkSet":
		return parseBeginNetworkSetEvent(event, height)
	case "story.dkg.v1.types.EventBeginDealing":
		return parseBeginDealingEvent(event, height)
	case "story.dkg.v1.types.EventBeginProcessDeals":
		return parseBeginProcessDealsEvent(event, height)
	case "story.dkg.v1.types.EventBeginProcessResponses":
		return parseBeginProcessResponsesEvent(event, height)
	case "story.dkg.v1.types.EventBeginFinalization":
		return parseBeginFinalizationEvent(event, height)
	case "story.dkg.v1.types.EventDKGFinalized":
		return parseDKGFinalizedEvent(event, height)
	default:
		return nil
	}
}

// parseBeginInitializationEvent parses EventBeginInitialization typed event.
func parseBeginInitializationEvent(event abcitypes.Event, height int64) *types.DKGEventData {
	var protoEvent dkgtypes.EventBeginInitialization

	for _, attr := range event.Attributes {
		key, val := attr.Key, attr.Value

		switch key {
		case "mrenclave":
			unquoted, err := strconv.Unquote(val)
			if err != nil {
				protoEvent.Mrenclave = []byte(val)
			} else {
				protoEvent.Mrenclave = []byte(unquoted)
			}
		case "round":
			round, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			protoEvent.Round = uint32(round)
		case "start_block":
			h, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				continue
			}
			protoEvent.StartBlock = uint32(h)
		case "active_validators":
			var validators []string
			if err := json.Unmarshal([]byte(val), &validators); err != nil {
				continue
			}
			protoEvent.ActiveValidators = validators
		default:
			continue
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
		key, val := attr.Key, attr.Value

		switch key {
		case "mrenclave":
			unquoted, err := strconv.Unquote(val)
			if err != nil {
				protoEvent.Mrenclave = []byte(val)
			} else {
				protoEvent.Mrenclave = []byte(unquoted)
			}
		case "round":
			round, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			protoEvent.Round = uint32(round)
		case "total":
			total, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			protoEvent.Total = uint32(total)
		case "threshold":
			threshold, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			protoEvent.Threshold = uint32(threshold)
		default:
			continue
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
		key, val := attr.Key, attr.Value

		switch key {
		case "mrenclave":
			unquoted, err := strconv.Unquote(val)
			if err != nil {
				protoEvent.Mrenclave = []byte(val)
			} else {
				protoEvent.Mrenclave = []byte(unquoted)
			}
		case "round":
			round, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			protoEvent.Round = uint32(round)
		default:
			continue
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
		//nolint:nestif // ignore nestedif linting
		key, val := attr.Key, attr.Value

		switch key {
		case "mrenclave":
			unquoted, err := strconv.Unquote(val)
			if err != nil {
				protoEvent.Mrenclave = []byte(val)
			} else {
				protoEvent.Mrenclave = []byte(unquoted)
			}
		case "round":
			round, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			protoEvent.Round = uint32(round)
		case "deals":
			var deals []*dkgtypes.Deal
			if err := json.Unmarshal([]byte(val), &deals); err != nil {
				continue
			}
			protoEvent.Deals = deals
		default:
			continue
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
		//nolint:nestif // ignore nestedif linting
		key, val := attr.Key, attr.Value

		switch key {
		case "mrenclave":
			unquoted, err := strconv.Unquote(val)
			if err != nil {
				protoEvent.Mrenclave = []byte(val)
			} else {
				protoEvent.Mrenclave = []byte(unquoted)
			}
		case "round":
			round, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			protoEvent.Round = uint32(round)
		case "responses":
			var responses []*dkgtypes.Response
			if err := json.Unmarshal([]byte(val), &responses); err != nil {
				continue
			}
			protoEvent.Responses = responses
		default:
			continue
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
		key, val := attr.Key, attr.Value

		switch key {
		case "mrenclave":
			unquoted, err := strconv.Unquote(val)
			if err != nil {
				protoEvent.Mrenclave = []byte(val)
			} else {
				protoEvent.Mrenclave = []byte(unquoted)
			}
		case "round":
			round, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			protoEvent.Round = uint32(round)
		default:
			continue
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
		key, val := attr.Key, attr.Value

		switch key {
		case "mrenclave":
			unquoted, err := strconv.Unquote(val)
			if err != nil {
				protoEvent.Mrenclave = []byte(val)
			} else {
				protoEvent.Mrenclave = []byte(unquoted)
			}
		case "round":
			round, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			protoEvent.Round = uint32(round)
		default:
			continue
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
