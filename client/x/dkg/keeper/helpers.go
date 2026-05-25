package keeper

import (
	"context"
	"time"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const (
	retryAttempts = 3
	retryDelay    = 2 * time.Second
)

func retry(ctx context.Context, fn func(ctx context.Context) error) error {
	var lastErr error
	for i := range retryAttempts {
		if err := fn(ctx); err != nil {
			lastErr = err
			log.Warn(context.Background(), "retry failed", err, "attempt", i+1)

			// Use context-aware sleep so that cancellation can interrupt the delay.
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay):
			}

			continue
		}

		return nil
	}

	if lastErr != nil {
		return lastErr
	}
	return errors.New("all retries failed")
}
