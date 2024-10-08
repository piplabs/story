package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/epochs/types"
)

// AddEpochInfo adds a new epoch info. Will return an error if the epoch fails validation,
// or re-uses an existing identifier.
// This method also sets the start time if left unset, and sets the epoch start height.
func (k Keeper) AddEpochInfo(ctx context.Context, epoch types.EpochInfo) error {
	err := epoch.Validate()
	if err != nil {
		return err
	}
	// Check if identifier already exists
	isExist, err := k.EpochInfo.Has(ctx, epoch.Identifier)
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("epoch with identifier %s already exists", epoch.Identifier)
	}

	// Initialize empty and default epoch values
	if epoch.StartTime.IsZero() {
		epoch.StartTime = sdk.UnwrapSDKContext(ctx).HeaderInfo().Time
	}
	if epoch.CurrentEpochStartHeight == 0 {
		epoch.CurrentEpochStartHeight = sdk.UnwrapSDKContext(ctx).HeaderInfo().Height
	}

	return k.EpochInfo.Set(ctx, epoch.Identifier, epoch)
}

// AllEpochInfos iterate through epochs to return all epochs info.
func (k Keeper) AllEpochInfos(ctx context.Context) ([]types.EpochInfo, error) {
	epochs := []types.EpochInfo{}
	err := k.EpochInfo.Walk(
		ctx,
		nil,
		func(_ string, value types.EpochInfo) (stop bool, err error) {
			epochs = append(epochs, value)
			return false, nil
		},
	)

	return epochs, err
}

// NumBlocksSinceEpochStart returns the number of blocks since the epoch started.
// if the epoch started on block N, then calling this during block N (after BeforeEpochStart)
// would return 0.
// Calling it any point in block N+1 (assuming the epoch doesn't increment) would return 1.
func (k Keeper) NumBlocksSinceEpochStart(ctx context.Context, identifier string) (int64, error) {
	epoch, err := k.EpochInfo.Get(ctx, identifier)
	if err != nil {
		return 0, fmt.Errorf("epoch with identifier %s not found", identifier)
	}

	return sdk.UnwrapSDKContext(ctx).HeaderInfo().Height - epoch.CurrentEpochStartHeight, nil
}

// GetEpochInfo gets current epoch info by identifier.
// NOTE(Narangde): add this func which can be a diff from cosmos-sdk.
func (k Keeper) GetEpochInfo(ctx context.Context, identifier string) (types.EpochInfo, error) {
	return k.EpochInfo.Get(ctx, identifier)
}
