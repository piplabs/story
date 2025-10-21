package keeper

import (
	"context"
	"fmt"
	"github.com/piplabs/story/lib/log"
	"time"
)

const (
	retryAttemts = 3
	retryDelay   = 2 * time.Second
)

func retry(ctx context.Context, fn func(ctx context.Context) error) error {
	for i := 0; i < retryAttemts; i++ {
		if err := fn(ctx); err != nil {
			log.Warn(context.Background(), "retry failed", err, "attempt", i+1)
			time.Sleep(retryDelay)

			continue
		}

		return nil
	}

	return fmt.Errorf("all retries failed")
}
