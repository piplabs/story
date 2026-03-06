package types

import (
	"context"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"math/big"

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
	GetLatestActiveRound(ctx context.Context) (*dkgtypes.DKGNetwork, error)

	Registered(ctx context.Context, msgSender common.Address, codeCommitment [32]byte, round uint32, startBlockHeight *big.Int, startBlockHash, enclaveType [32]byte, dkgPubKey []byte, commPubKey []byte, enclaveReport []byte) error
	Finalized(ctx context.Context, round uint32, msgSender common.Address, codeCommitment, participantsRoot [32]byte, signature, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte) error
	UpgradeScheduled(ctx context.Context, activationHeight int64, upgradeVersion string) error
	UpgradeCancelled(ctx context.Context, upgradeVersion string) error
	ThresholdDecryptRequested(ctx context.Context, round uint32, requesterPubKey []byte, ciphertext []byte, label [32]byte) error

	// Parameter setters (driven by DKG.sol contract events)
	SetMinReqRegisteredParticipants(ctx context.Context, value uint32) error
	SetMinReqFinalizedParticipants(ctx context.Context, value uint32) error
	SetOperationalThreshold(ctx context.Context, value uint32) error
}
