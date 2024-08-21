package keeper

import (
	"context"
	"sync"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/expbackoff"
)

// backoffFunc aliased for testing.
var (
	backoffFuncMu sync.RWMutex
	backoffFunc   = expbackoff.New
)

func retryForever(ctx context.Context, fn func(ctx context.Context) (bool, error)) error {
	backoffFuncMu.RLock()
	backoff := backoffFunc(ctx)
	backoffFuncMu.RUnlock()
	for {
		ok, err := fn(ctx)
		if ctx.Err() != nil {
			return errors.Wrap(ctx.Err(), "retry canceled")
		} else if err != nil {
			return err
		} else if !ok {
			backoff()
			continue
		}

		return nil
	}
}
