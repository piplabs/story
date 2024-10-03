//nolint:revive // just use interface{}
package keeper

import (
	"context"
	"fmt"
	"math/big"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	evmenginetypes "github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/client/x/mint/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	clog "github.com/piplabs/story/lib/log"
)

// Keeper of the mint store.
type Keeper struct {
	cdc              codec.BinaryCodec
	storeService     storetypes.KVStoreService
	stakingKeeper    types.StakingKeeper
	bankKeeper       types.BankKeeper
	feeCollectorName string

	inflationUpdateContract *bindings.IPTokenStaking // (rayden) TODO

	Schema collections.Schema
	Params collections.Item[types.Params]
}

// NewKeeper creates a new mint Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	sk types.StakingKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	feeCollectorName string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the x/%s module account has not been set", types.ModuleName))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:              cdc,
		storeService:     storeService,
		stakingKeeper:    sk,
		bankKeeper:       bk,
		feeCollectorName: feeCollectorName,
		Params:           collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx context.Context) (math.Int, error) {
	return k.stakingKeeper.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx context.Context) (math.LegacyDec, error) {
	return k.stakingKeeper.BondedRatio(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx context.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx context.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}

func (k Keeper) ProcessInflationEvents(ctx context.Context, height uint64, logs []*evmenginetypes.EVMEvent) error {
	gwei, exp := big.NewInt(10), big.NewInt(9)
	gwei.Exp(gwei, exp, nil)

	for _, evmLog := range logs {
		if err := evmLog.Verify(); err != nil {
			return errors.Wrap(err, "verify log [BUG]") // This shouldn't happen
		}
		ethlog := evmLog.ToEthLog()

		// (rayden) TODO: handle when each event processing fails.
		if ethlog.Topics[0] == types.SetInflationParameters.ID {
			ev, err := k.inflationUpdateContract.ParseDeposit(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse SetInflationParameters log", err)
				continue
			}
			if err = k.ProcessSetInflationParameters(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process update inflation parameters", err)
				continue
			}
		}
	}

	clog.Debug(ctx, "Processed inflation events", "height", height, "count", len(logs))

	return nil
}
