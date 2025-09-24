package types

import (
	"context"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/contracts/bindings"
)

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type EvmStakingKeeper interface {
	ParseDepositLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingDeposit, error)
	ParseWithdrawLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingWithdraw, error)
	ProcessStakingEvents(ctx context.Context, height uint64, logs []*ethtypes.Log) error
	MaxWithdrawalPerBlock(ctx context.Context) (uint32, error)
	DequeueEligibleWithdrawals(ctx context.Context, maxDequeue uint32) (withdrawals ethtypes.Withdrawals, err error)
	PeekEligibleWithdrawals(ctx context.Context, maxPeek uint32) (withdrawals ethtypes.Withdrawals, err error)
	DequeueEligibleRewardWithdrawals(ctx context.Context, maxDequeue uint32) (withdrawals ethtypes.Withdrawals, err error)
	PeekEligibleRewardWithdrawals(ctx context.Context, maxPeek uint32) (withdrawals ethtypes.Withdrawals, err error)
}

type UpgradeKeeper interface {
	ClearUpgradePlan(ctx context.Context) error
	ScheduleUpgrade(ctx context.Context, plan upgradetypes.Plan) error
	DumpUpgradeInfoToDisk(height int64, p upgradetypes.Plan) error
}

type DistrKeeper interface {
	SetUbi(ctx context.Context, newUbi math.LegacyDec) error
}

type DKGKeeper interface {
	// NOTE: completed
	RegistrationInitialized(ctx context.Context, msgSender common.Address, mrenclave []byte, round uint32, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) error
	NetworkSet(ctx context.Context, msgSender common.Address, mrenclave []byte, round uint32, total uint32, threshold uint32, signature []byte) error
	// TODO: complete these functions
	Finalized(ctx context.Context, round uint32, msgSender common.Address, mrenclave []byte, signature []byte) error
	UpgradeScheduled(ctx context.Context, activationHeight uint32, mrenclave []byte) error
	RemoteAttestationProcessedOnChain(ctx context.Context, validator common.Address, chalStatus int, round uint32, mrenclave []byte) error
	DealComplaintsSubmitted(ctx context.Context, index uint32, complainIndexes []uint32, round uint32, mrenclave []byte) error
	DealVerified(ctx context.Context, index uint32, recipientIndex uint32, round uint32, mrenclave []byte) error
	InvalidDeal(ctx context.Context, index uint32, round uint32, mrenclave []byte) error
}
