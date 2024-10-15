package types

import (
	"context"
	signaltypes "github.com/piplabs/story/client/x/signal/types"

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

type SignalKeeper interface {
	ScheduleUpgrade(ctx context.Context, upgrade signaltypes.Upgrade)
}

type MintKeeper interface {
	ProcessInflationEvents(ctx context.Context, height uint64, logs []*EVMEvent) error
}
