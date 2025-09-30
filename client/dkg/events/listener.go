package events

import (
	"context"
	"time"

	"github.com/cometbft/cometbft/rpc/client/http"

	"github.com/piplabs/story/client/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

type EventListener struct {
	client          *http.HTTP
	eventChan       chan *types.DKGEventData
	stopChan        chan struct{}
	pollingInterval time.Duration
	lastHeight      int64
}

func NewEventListener(rpcEndpoint string, pollingInterval time.Duration) (*EventListener, error) {
	client, err := http.New(rpcEndpoint, "/websocket")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tendermint client")
	}

	return &EventListener{
		client:          client,
		eventChan:       make(chan *types.DKGEventData, 100),
		stopChan:        make(chan struct{}),
		pollingInterval: pollingInterval,
		lastHeight:      0,
	}, nil
}

func (l *EventListener) Start(ctx context.Context) error {
	log.Info(ctx, "Starting DKG event listener")

	if err := l.client.Start(); err != nil {
		return errors.Wrap(err, "failed to start tendermint client")
	}

	go l.pollEvents(ctx)

	return nil
}

func (l *EventListener) Stop(ctx context.Context) error {
	log.Info(ctx, "Stopping DKG event listener")

	close(l.stopChan)

	err := l.client.Stop()
	if err != nil {
		return errors.Wrap(err, "stop dkg event listener")
	}

	return nil
}

func (l *EventListener) EventChannel() <-chan *types.DKGEventData {
	return l.eventChan
}

func (l *EventListener) pollEvents(ctx context.Context) {
	ticker := time.NewTicker(l.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-l.stopChan:
			return
		case <-ticker.C:
			if err := l.checkNewBlocks(ctx); err != nil {
				log.Error(ctx, "Error checking new blocks", err)
			}
		}
	}
}

// checkNewBlocks checks for new blocks since the last processed height.
func (l *EventListener) checkNewBlocks(ctx context.Context) error {
	status, err := l.client.Status(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get node status")
	}

	currentHeight := status.SyncInfo.LatestBlockHeight

	for height := l.lastHeight + 1; height <= currentHeight; height++ {
		if err := l.processBlock(ctx, height); err != nil {
			log.Error(ctx, "Error processing block", err, "height", height)
			continue
		}
	}

	l.lastHeight = currentHeight

	return nil
}

// processBlock processes a single block for DKG events.
func (l *EventListener) processBlock(ctx context.Context, height int64) error {
	blockResults, err := l.client.BlockResults(ctx, &height)
	if err != nil {
		return errors.Wrap(err, "failed to get block results")
	}

	//nolint:nestif // ignore nestedif linting
	if blockResults.FinalizeBlockEvents != nil {
		for _, event := range blockResults.FinalizeBlockEvents {
			if isDKGEvent(event.Type) {
				dkgEvent := l.parseEvent(event, height)
				if dkgEvent != nil {
					if err := l.sendEventWithRetry(ctx, dkgEvent, height); err != nil {
						return err
					}
				} else {
					log.Warn(ctx, "Failed to parse DKG event", nil,
						"event_type", event.Type,
						"height", height)
				}
			}
		}
	}

	for _, txResult := range blockResults.TxsResults {
		for _, event := range txResult.Events {
			if isDKGEvent(event.Type) {
				dkgEvent := l.parseEvent(event, height)
				if dkgEvent != nil {
					if err := l.sendEventWithRetry(ctx, dkgEvent, height); err != nil {
						return err
					}
				} else {
					log.Warn(ctx, "Failed to parse DKG event", nil,
						"event_type", event.Type,
						"height", height)
				}
			}
		}
	}

	return nil
}

// sendEventWithRetry attempts to send a DKG event to the channel with retry logic when the channel is full.
func (l *EventListener) sendEventWithRetry(ctx context.Context, dkgEvent *types.DKGEventData, height int64) error {
	const (
		maxRetries = 5
		baseDelay  = 100 * time.Millisecond
		maxDelay   = 2 * time.Second
	)

	for attempt := range maxRetries {
		select {
		case l.eventChan <- dkgEvent:
			log.Debug(ctx, "DKG event detected",
				"parsed_type", dkgEvent.EventType,
				"mrenclave", dkgEvent.Mrenclave,
				"round", dkgEvent.Round,
				"height", height,
			)

			return nil
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "context done")
		case <-l.stopChan:
			return nil
		default:
			if attempt == 0 {
				log.Warn(ctx, "DKG event channel is full, retrying with backoff", nil,
					"event_type", dkgEvent.EventType,
					"round", dkgEvent.Round,
					"height", height)
			}
			delay := baseDelay * time.Duration(1<<attempt)
			if delay > maxDelay {
				delay = maxDelay
			}

			select {
			case <-time.After(delay):
				// Continue to next retry attempt
			case <-ctx.Done():
				return errors.Wrap(ctx.Err(), "context done during retry")
			case <-l.stopChan:
				return nil
			}
		}
	}

	log.Error(ctx, "Failed to send DKG event after maximum retries, dropping event", nil,
		"event_type", dkgEvent.EventType,
		"round", dkgEvent.Round,
		"height", height,
		"max_retries", maxRetries)

	return nil // Don't return an error to avoid stopping the entire processing
}
