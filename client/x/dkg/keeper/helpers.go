package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/netconf"
)

const (
	retryAttempts = 3
	retryDelay    = 2 * time.Second
)

// isV190Round reports whether a round is governed by the v1.9.0 DKG validation rules.
// It checks netconf.IsV190 against the round's start height (not the current height)
// so a round spanning the upgrade boundary keeps a single rule set, and treats an
// unregistered chain as not active.
func (k *Keeper) isV190Round(ctx context.Context, round *types.DKGNetwork) bool {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	active, err := netconf.IsV190(sdkCtx.ChainID(), round.StartBlockHeight)

	return err == nil && active
}

// dealtDealerKey is the DealtDealers state key for a dealer in a round, keyed by the
// dealer's validator address. Addresses (not indices) are used because in resharing
// rounds the dealer committee's index space differs from the current round's.
func dealtDealerKey(round uint32, dealerAddr string) string {
	return fmt.Sprintf("%d_%s", round, common.HexToAddress(dealerAddr).Hex())
}

func retry(ctx context.Context, fn func(ctx context.Context) error) error {
	for i := range retryAttempts {
		if err := fn(ctx); err != nil {
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

	return errors.New("all retries failed")
}
