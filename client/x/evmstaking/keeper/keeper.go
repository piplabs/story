package keeper

import (
	"context"
	"fmt"
	"math/big"

	"cosmossdk.io/collections"
	addresscodec "cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	addcollections "github.com/piplabs/story/client/collections"
	"github.com/piplabs/story/client/genutil/evm/predeploys"
	evmenginetypes "github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"
	clog "github.com/piplabs/story/lib/log"
)

// Keeper of the x/evmstaking store.
type Keeper struct {
	cdc                   codec.BinaryCodec
	storeService          store.KVStoreService
	validatorAddressCodec addresscodec.Codec
	authority             string

	authKeeper         types.AccountKeeper
	bankKeeper         types.BankKeeper
	slashingKeeper     types.SlashingKeeper
	stakingKeeper      types.StakingKeeper
	distributionKeeper types.DistributionKeeper

	ipTokenStakingContract *bindings.IPTokenStaking

	WithdrawalQueue          addcollections.Queue[types.Withdrawal]
	DelegatorWithdrawAddress collections.Map[string, string]
	DelegatorRewardAddress   collections.Map[string, string]
	DelegatorOperatorAddress collections.Map[string, string]
}

// NewKeeper creates a new evmstaking Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	slk types.SlashingKeeper,
	stk types.StakingKeeper,
	dk types.DistributionKeeper,
	authority string,
	ethCl ethclient.Client,
	validatorAddressCodec addresscodec.Codec,
) *Keeper {
	// ensure that authority is a valid AccAddress
	if _, err := ak.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	// ensure the module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(types.ModuleName + " module account has not been set")
	}

	sb := collections.NewSchemaBuilder(storeService)

	ipTokenStakingContract, err := bindings.NewIPTokenStaking(common.HexToAddress(predeploys.IPTokenStaking), ethCl)
	if err != nil {
		panic(fmt.Sprintf("failed to bind to the IPTokenStaking contract: %s", err))
	}

	return &Keeper{
		cdc:                      cdc,
		storeService:             storeService,
		authKeeper:               ak,
		bankKeeper:               bk,
		slashingKeeper:           slk,
		stakingKeeper:            stk,
		distributionKeeper:       dk,
		authority:                authority,
		validatorAddressCodec:    validatorAddressCodec,
		ipTokenStakingContract:   ipTokenStakingContract,
		WithdrawalQueue:          addcollections.NewQueue(sb, types.WithdrawalQueueKey, "withdrawal_queue", codec.CollValue[types.Withdrawal](cdc)),
		DelegatorWithdrawAddress: collections.NewMap(sb, types.DelegatorWithdrawAddressMapKey, "delegator_withdraw_address_map", collections.StringKey, collections.StringValue),
		DelegatorRewardAddress:   collections.NewMap(sb, types.DelegatorRewardAddressMapKey, "delegator_reward_address_map", collections.StringKey, collections.StringValue),
		DelegatorOperatorAddress: collections.NewMap(sb, types.DelegatorOperatorAddressMapKey, "delegator_operator_address_map", collections.StringKey, collections.StringValue),
	}
}

func Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetAuthority returns the x/evmstaking module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// ValidatorAddressCodec returns the app validator address codec.
func (k Keeper) ValidatorAddressCodec() addresscodec.Codec {
	return k.validatorAddressCodec
}

//nolint:gocyclo // TODO
func (k Keeper) ProcessStakingEvents(ctx context.Context, height uint64, logs []*evmenginetypes.EVMEvent) error {
	gwei, exp := big.NewInt(10), big.NewInt(9)
	gwei.Exp(gwei, exp, nil)

	for _, evmLog := range logs {
		if err := evmLog.Verify(); err != nil {
			return errors.Wrap(err, "verify log [BUG]") // This shouldn't happen
		}
		ethlog, err := evmLog.ToEthLog()
		if err != nil {
			return err
		}

		// TODO: handle when each event processing fails.

		// Convert the amount from wei to gwei (Eth2 spec withdrawal is specified in gwei) by dividing by 10^9.
		// TODO: consider rounding and decimal precision when dividing bigint.

		switch ethlog.Topics[0] {
		case types.UpdateValidatorCommission.ID:
			ev, err := k.ipTokenStakingContract.ParseUpdateValidatorCommssion(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse UpdateValidatorCommission log", err)
				continue
			}
			if err = k.ProcessUpdateValidatorCommission(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process update validator commission", err)
				continue
			}
		case types.SetWithdrawalAddress.ID:
			ev, err := k.ipTokenStakingContract.ParseSetWithdrawalAddress(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse SetWithdrawalAddress log", err)
				continue
			}
			if err = k.ProcessSetWithdrawalAddress(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process set withdrawal address", err)
				continue
			}
		case types.SetRewardAddress.ID:
			ev, err := k.ipTokenStakingContract.ParseSetRewardAddress(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse SetRewardAddress log", err)
				continue
			}
			if err = k.ProcessSetRewardAddress(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process set reward address", err)
				continue
			}
		case types.AddOperator.ID:
			ev, err := k.ipTokenStakingContract.ParseAddOperator(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse SetRewardAddress log", err)
				continue
			}
			if err = k.ProcessAddOperator(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process add operator", err)
				continue
			}
		case types.RemoveOperator.ID:
			ev, err := k.ipTokenStakingContract.ParseRemoveOperator(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse SetRewardAddress log", err)
				continue
			}
			if err = k.ProcessRemoveOperator(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process add operator", err)
				continue
			}
		case types.CreateValidatorEvent.ID:
			ev, err := k.ParseCreateValidatorLog(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse CreateValidator log", err)
				continue
			}
			ev.StakeAmount.Div(ev.StakeAmount, gwei)
			if err = k.ProcessCreateValidator(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process create validator", err)
				continue
			}
		case types.DepositEvent.ID:
			ev, err := k.ParseDepositLog(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse Deposit log", err)
				continue
			}
			ev.StakeAmount.Div(ev.StakeAmount, gwei)
			if err = k.ProcessDeposit(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process deposit", err)
				continue
			}
		case types.RedelegateEvent.ID:
			ev, err := k.ParseRedelegateLog(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse Redelegate log", err)
				continue
			}
			ev.Amount.Div(ev.Amount, gwei)
			if err = k.ProcessRedelegate(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process redelegate", err)
				continue
			}
		case types.WithdrawEvent.ID:
			ev, err := k.ParseWithdrawLog(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse Withdraw log", err)
				continue
			}
			ev.StakeAmount.Div(ev.StakeAmount, gwei)
			if err = k.ProcessWithdraw(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process withdraw", err)
				continue
			}
		case types.UnjailEvent.ID:
			ev, err := k.ParseUnjailLog(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse Unjail log", err)
				continue
			}
			if err = k.ProcessUnjail(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process unjail", err)
				continue
			}
		}
	}

	clog.Debug(ctx, "Processed staking events", "height", height, "count", len(logs))

	return nil
}
