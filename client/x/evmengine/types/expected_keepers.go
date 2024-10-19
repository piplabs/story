package types

import (
	"context"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/contracts/bindings"
)

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type EvmStakingKeeper interface {
	ProcessDeposit(ctx context.Context, ev *bindings.IPTokenStakingDeposit) error
	ProcessWithdraw(ctx context.Context, ev *bindings.IPTokenStakingWithdraw) error
	DequeueEligibleWithdrawals(ctx context.Context) (ethtypes.Withdrawals, error)
	ParseDepositLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingDeposit, error)
	ParseWithdrawLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingWithdraw, error)
	ProcessStakingEvents(ctx context.Context, height uint64, logs []*EVMEvent) error
	PeekEligibleWithdrawals(ctx context.Context) (withdrawals ethtypes.Withdrawals, err error)
}

type UpgradeKeeper interface {
	ScheduleUpgrade(ctx context.Context, plan upgradetypes.Plan) error
}

type DistrKeeper interface {
	SetUbi(ctx context.Context, newUbi math.LegacyDec) error
}
