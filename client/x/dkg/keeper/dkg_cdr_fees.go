package keeper

import (
	"context"
	"math/big"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

func cdrSubmitCountKey(validator common.Address) string {
	return strings.ToLower(validator.Hex())
}

// AddCDRFeeToPool mints tokens into the DKG module and transfers them to the CDR fee pool.
func (k *Keeper) AddCDRFeeToPool(ctx context.Context, amount *big.Int) error {
	feeAmount, ok, err := parseCDRFeeAmount(amount)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, feeAmount))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return errors.Wrap(err, "mint CDR fee coins")
	}
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.CDRFeePoolName, coins); err != nil {
		return errors.Wrap(err, "transfer CDR fee coins to pool")
	}

	current, _, err := k.getCDRFeePoolBalance(ctx)
	if err != nil {
		return err
	}

	newBalance := current.Add(feeAmount)
	if err := k.CDRFeePoolBalance.Set(ctx, newBalance.String()); err != nil {
		return errors.Wrap(err, "set CDR fee pool balance")
	}

	return nil
}

// RefundCDRFee sends tokens from the CDR fee pool to the validator.
func (k *Keeper) RefundCDRFee(ctx context.Context, validator common.Address, amount *big.Int) error {
	feeAmount, ok, err := parseCDRFeeAmount(amount)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	recipient, err := utils.EvmAddressToBech32AccAddress(validator.Hex())
	if err != nil {
		return errors.Wrap(err, "convert validator address")
	}

	current, found, err := k.getCDRFeePoolBalance(ctx)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("cdr fee pool balance not found")
	}

	newBalance := current.Sub(feeAmount)
	if newBalance.IsNegative() {
		return errors.New("cdr fee pool balance underflow",
			"balance", current.String(),
			"refund", feeAmount.String(),
		)
	}

	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, feeAmount))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.CDRFeePoolName, recipient, coins); err != nil {
		return errors.Wrap(err, "refund CDR fee coins")
	}

	if newBalance.IsZero() {
		if err := k.CDRFeePoolBalance.Remove(ctx); err != nil {
			return errors.Wrap(err, "remove CDR fee pool balance")
		}
		return nil
	}

	if err := k.CDRFeePoolBalance.Set(ctx, newBalance.String()); err != nil {
		return errors.Wrap(err, "set CDR fee pool balance")
	}

	return nil
}

// IncrementCDRPartialSubmitCount increments the valid partial submission count for a validator.
func (k *Keeper) IncrementCDRPartialSubmitCount(ctx context.Context, validator common.Address) error {
	key := cdrSubmitCountKey(validator)
	count, err := k.CDRPartialSubmitCount.Get(ctx, key)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return errors.Wrap(err, "get CDR submit count")
		}

		count = 0
	}

	if err := k.CDRPartialSubmitCount.Set(ctx, key, count+1); err != nil {
		return errors.Wrap(err, "set CDR submit count")
	}

	return nil
}

func (k *Keeper) distributeCDRRewardPool(ctx context.Context) error {
	iter, err := k.CDRPartialSubmitCount.Iterate(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "iterate CDR submit counts")
	}
	defer iter.Close()

	counts := map[string]uint64{}
	var totalCount uint64

	for ; iter.Valid(); iter.Next() {
		key, err := iter.Key()
		if err != nil {
			return errors.Wrap(err, "iterate CDR submit count key")
		}
		count, err := iter.Value()
		if err != nil {
			return errors.Wrap(err, "iterate CDR submit count value")
		}

		if key == "" {
			continue
		}

		counts[key] += count
		totalCount += count
	}

	if totalCount == 0 {
		return nil
	}

	poolBalance, found, err := k.getCDRFeePoolBalance(ctx)
	if err != nil {
		return err
	}
	if !found || poolBalance.IsZero() {
		if err := k.CDRPartialSubmitCount.Clear(ctx, nil); err != nil {
			return errors.Wrap(err, "clear CDR submit count")
		}
		if found && poolBalance.IsZero() {
			if err := k.CDRFeePoolBalance.Remove(ctx); err != nil {
				return errors.Wrap(err, "remove CDR fee pool balance")
			}
		}
		return nil
	}

	totalCountInt := math.NewInt(int64(totalCount))
	distributed := math.ZeroInt()

	for addr, count := range counts {
		share := poolBalance.Mul(math.NewInt(int64(count))).Quo(totalCountInt)
		if share.IsZero() {
			continue
		}

		recipient, err := utils.EvmAddressToBech32AccAddress(addr)
		if err != nil {
			return errors.Wrap(err, "convert validator address", "address", addr)
		}

		coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, share))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.CDRFeePoolName, recipient, coins); err != nil {
			return errors.Wrap(err, "distribute CDR fee pool", "address", addr)
		}

		distributed = distributed.Add(share)
	}

	if err := k.CDRPartialSubmitCount.Clear(ctx, nil); err != nil {
		return errors.Wrap(err, "clear CDR submit count")
	}

	remaining := poolBalance.Sub(distributed)
	if remaining.IsZero() {
		if err := k.CDRFeePoolBalance.Remove(ctx); err != nil {
			return errors.Wrap(err, "remove CDR fee pool balance")
		}
		return nil
	}

	if err := k.CDRFeePoolBalance.Set(ctx, remaining.String()); err != nil {
		return errors.Wrap(err, "set CDR fee pool balance")
	}

	return nil
}

func (k *Keeper) getCDRFeePoolBalance(ctx context.Context) (math.Int, bool, error) {
	balStr, err := k.CDRFeePoolBalance.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return math.ZeroInt(), false, nil
		}
		return math.ZeroInt(), false, errors.Wrap(err, "get CDR fee pool balance")
	}

	balance, ok := math.NewIntFromString(balStr)
	if !ok {
		return math.ZeroInt(), false, errors.New("invalid CDR fee pool balance", "value", balStr)
	}

	return balance, true, nil
}

func parseCDRFeeAmount(amount *big.Int) (math.Int, bool, error) {
	if amount == nil || amount.Sign() == 0 {
		return math.ZeroInt(), false, nil
	}

	feeAmount, ok := math.NewIntFromString(amount.String())
	if !ok {
		return math.ZeroInt(), false, errors.New("invalid CDR fee amount", "value", amount.String())
	}

	return feeAmount, true, nil
}
