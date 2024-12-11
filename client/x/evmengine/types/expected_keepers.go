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
	ParseDepositLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingDeposit, error)
	ParseWithdrawLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingWithdraw, error)
	ProcessStakingEvents(ctx context.Context, height uint64, logs []*EVMEvent) error
	MaxWithdrawalPerBlock(ctx context.Context) (uint32, error)
	DequeueEligibleWithdrawals(ctx context.Context, maxDequeue uint32) (withdrawals ethtypes.Withdrawals, err error)
	PeekEligibleWithdrawals(ctx context.Context, maxPeek uint32) (withdrawals ethtypes.Withdrawals, err error)
	DequeueEligibleRewardWithdrawals(ctx context.Context, maxDequeue uint32) (withdrawals ethtypes.Withdrawals, err error)
	PeekEligibleRewardWithdrawals(ctx context.Context, maxPeek uint32) (withdrawals ethtypes.Withdrawals, err error)
}

type UpgradeKeeper interface {
	ClearUpgradePlan(ctx context.Context) error
	ScheduleUpgrade(ctx context.Context, plan upgradetypes.Plan) error
}

type DistrKeeper interface {
	SetUbi(ctx context.Context, newUbi math.LegacyDec) error
}
