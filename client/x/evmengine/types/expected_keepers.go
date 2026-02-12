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
	RegistrationInitialized(ctx context.Context, msgSender common.Address, codeCommitment [32]byte, round uint32, startBlockHeight uint64, startBlockHash [32]byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) error
	Finalized(ctx context.Context, round uint32, msgSender common.Address, codeCommitment, participantsRoot [32]byte, signature, globalPubKey []byte, publicCoeffs [][]byte) error

	// TODO: complete these functions
	UpgradeScheduled(ctx context.Context, activationHeight uint32, codeCommitment [32]byte) error
	RemoteAttestationProcessedOnChain(ctx context.Context, validator common.Address, chalStatus int, round uint32, codeCommitment [32]byte) error
	DealComplaintsSubmitted(ctx context.Context, index uint32, complainIndexes []uint32, round uint32, codeCommitment [32]byte) error
	DealVerified(ctx context.Context, index uint32, recipientIndex uint32, round uint32, codeCommitment [32]byte) error
	InvalidDeal(ctx context.Context, index uint32, round uint32, codeCommitment [32]byte) error
	ThresholdDecryptRequested(ctx context.Context, requester common.Address, round uint32, codeCommitment [32]byte, requesterPubKey []byte, ciphertext []byte, label []byte) error
}
