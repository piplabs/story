package keeper

import (
	"context"
	"time"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const (
	retryAttemts = 3
	retryDelay   = 2 * time.Second
)

func retry(ctx context.Context, fn func(ctx context.Context) error) error {
	for i := range retryAttemts {
		if err := fn(ctx); err != nil {
			log.Warn(context.Background(), "retry failed", err, "attempt", i+1)
			time.Sleep(retryDelay)

			continue
		}

		return nil
	}

	return errors.New("all retries failed")
}
