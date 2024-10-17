package keeper

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/piplabs/story/client/x/signal/types"
)

// Keeper implements the MsgServer and QueryServer interfaces.
var (
	_ types.MsgServiceServer = &Keeper{}
	_ types.QueryServer      = Keeper{}

	// defaultSignalThreshold is 5/6 or approximately 83.33%.
	defaultSignalThreshold = sdkmath.LegacyNewDec(5).Quo(sdkmath.LegacyNewDec(6))
)

// Threshold is the fraction of voting power that is required
// to signal for a version change. It is set to 5/6 as the middle point
// between 2/3 and 3/3 providing 1/6 fault tolerance to halting the
// network during an upgrade period. It can be modified through a
// hard fork change that modified the app version.
func Threshold(_ uint64) sdkmath.LegacyDec {
	return defaultSignalThreshold
}

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	authority    string

	authKeeper    types.AccountKeeper
	stakingKeeper types.StakingKeeper
}

// NewKeeper returns a signal keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	ak types.AccountKeeper,
	stk types.StakingKeeper,
	authority string,
) *Keeper {
	// ensure that authority is a valid AccAddress
	if _, err := ak.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	// ensure the module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(types.ModuleName + " module account has not been set")
	}

	return &Keeper{
		cdc:           cdc,
		storeService:  storeService,
		authKeeper:    ak,
		stakingKeeper: stk,
	}
}

// SetValidatorVersion saves a signaled version for a validator.
func (k Keeper) SetValidatorVersion(ctx sdk.Context, valAddress sdk.ValAddress, version uint64) {
	stores := k.storeService.OpenKVStore(ctx)
	if err := stores.Set(valAddress, VersionToBytes(version)); err != nil {
		panic(err)
	}
}

// DeleteValidatorVersion deletes a signaled version for a validator.
func (k Keeper) DeleteValidatorVersion(ctx sdk.Context, valAddress sdk.ValAddress) {
	stores := k.storeService.OpenKVStore(ctx)
	if err := stores.Delete(valAddress); err != nil {
		panic(err)
	}
}

// ShouldUpgrade returns whether the signaling mechanism has concluded that the
// network is ready to upgrade and the version to upgrade to. It returns false
// and 0 if no version has reached quorum.
func (k *Keeper) ShouldUpgrade(ctx sdk.Context) (isQuorumVersion bool, version uint64) {
	upgrade, ok := k.getUpgrade(ctx)
	if !ok {
		return false, 0
	}

	hasUpgradeHeightBeenReached := ctx.BlockHeight() >= upgrade.UpgradeHeight
	if hasUpgradeHeightBeenReached {
		return true, upgrade.AppVersion
	}
	return false, 0
}

func VersionToBytes(version uint64) []byte {
	return binary.BigEndian.AppendUint64(nil, version)
}

func VersionFromBytes(version []byte) uint64 {
	return binary.BigEndian.Uint64(version)
}

func (k Keeper) ScheduleUpgrade(ctx context.Context, msg *types.MsgScheduleUpgrade) (*types.MsgScheduleUpgradeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if k.IsUpgradePending(sdkCtx) {
		return nil, types.ErrUpgradePending
	}

	upgrade := msg.Upgrade
	if upgrade.AppVersion <= sdkCtx.BlockHeader().Version.App {
		return nil, types.ErrInvalidUpgradeVersion.Wrapf("can not upgrade to version %v because it is less than or equal to current version %v", upgrade.AppVersion, sdkCtx.BlockHeader().Version.App)
	}
	if err := k.setUpgrade(sdkCtx, *upgrade); err != nil {
		return nil, err
	}

	return &types.MsgScheduleUpgradeResponse{}, nil
}

// GetUpgrade returns the current upgrade information.
func (k Keeper) GetUpgrade(ctx context.Context, _ *types.QueryGetUpgradeRequest) (*types.QueryGetUpgradeResponse, error) {
	upgrade, ok := k.getUpgrade(sdk.UnwrapSDKContext(ctx))
	if !ok {
		return &types.QueryGetUpgradeResponse{}, nil
	}
	return &types.QueryGetUpgradeResponse{Upgrade: &upgrade}, nil
}

// IsUpgradePending returns true if an app version has reached quorum and the
// chain should upgrade to the app version at the upgrade height. While the
// keeper has an upgrade pending the SignalVersion and TryUpgrade messages will
// be rejected.
func (k *Keeper) IsUpgradePending(ctx sdk.Context) bool {
	_, ok := k.getUpgrade(ctx)
	return ok
}

// getUpgrade returns the current upgrade information from the store.
// If an upgrade is found, it returns the upgrade object and true.
// If no upgrade is found, it returns an empty upgrade object and false.
func (k *Keeper) getUpgrade(ctx sdk.Context) (upgrade types.Upgrade, ok bool) {
	stores := k.storeService.OpenKVStore(ctx)
	bz, err := stores.Get(types.UpgradeKey)
	if err != nil {
		return types.Upgrade{}, false
	}

	err = k.cdc.Unmarshal(bz, &upgrade)
	if err != nil {
		return upgrade, false
	}

	return upgrade, true
}

// setUpgrade sets the upgrade in the store.
func (k *Keeper) setUpgrade(ctx sdk.Context, upgrade types.Upgrade) error {
	stores := k.storeService.OpenKVStore(ctx)
	value := k.cdc.MustMarshal(&upgrade)

	if err := stores.Set(types.UpgradeKey, value); err != nil {
		return err
	}

	return nil
}
