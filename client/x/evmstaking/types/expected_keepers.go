package types

import (
	"context"
	"time"

	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/math"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper (noalias).
type AccountKeeper interface {
	AddressCodec() address.Codec
	HasAccount(ctx context.Context, addr sdk.AccAddress) bool
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
	GetModuleAddress(moduleName string) sdk.AccAddress
	// only used for simulation
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	IterateAccounts(ctx context.Context, process func(sdk.AccountI) (stop bool))
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	SetModuleAccount(ctx context.Context, modAcc sdk.ModuleAccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	// only used for simulation
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	LockedCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoin(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetSupply(ctx context.Context, denom string) sdk.Coin
	SendCoinsFromModuleToModule(ctx context.Context, senderPool, recipientPool string, amt sdk.Coins) error
}

// StakingKeeper defines the expected interface for the staking module.
type StakingKeeper interface {
	ValidatorAddressCodec() address.Codec

	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	GetAllValidators(ctx context.Context) (validators []stakingtypes.Validator, err error)
	BondDenom(ctx context.Context) (string, error)

	SetUnbondingDelegationEntry(
		ctx context.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
		creationHeight int64, minTime time.Time, balance math.Int,
	) (stakingtypes.UnbondingDelegation, error)
	UBDQueueIterator(ctx context.Context, endTime time.Time) (corestore.Iterator, error)
	// GetUnbondingDelegation returns a unbonding delegation.
	GetUnbondingDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubd stakingtypes.UnbondingDelegation, err error)
	DeleteUnbondingIndex(ctx context.Context, id uint64) error

	GetAllDelegations(ctx context.Context) (delegations []stakingtypes.Delegation, err error)
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, err error)
	GetValidatorDelegations(ctx context.Context, valAddr sdk.ValAddress) (delegations []stakingtypes.Delegation, err error)
	GetUnbondingDelegations(ctx context.Context, delegator sdk.AccAddress, maxRetrieve uint16) (unbondingDelegations []stakingtypes.UnbondingDelegation, err error)
	GetUnbondingDelegationsFromValidator(ctx context.Context, valAddr sdk.ValAddress) (ubds []stakingtypes.UnbondingDelegation, err error)

	EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error)
	EndBlockerWithUnbondedEntries(ctx context.Context) ([]abci.ValidatorUpdate, []stakingtypes.UnbondedEntry, error)

	MinDelegation(ctx context.Context) (math.Int, error)
	GetFlexiblePeriodType(ctx context.Context) (int32, error)
	GetPeriodInfo(ctx context.Context, periodType int32) (stakingtypes.Period, error)
	GetLockedTokenType(ctx context.Context) (int32, error)
	GetTokenTypeInfo(ctx context.Context, tokenType int32) (stakingtypes.TokenTypeInfo, error)

	GetSingularityHeight(ctx context.Context) (uint64, error)
}

// SlashingKeeper defines the expected interface for the slashing module.
type SlashingKeeper interface {
	Unjail(ctx context.Context, validatorAddr sdk.ValAddress) error
}

// DistributionKeeper defines the expected interface needed to calculate validator commission and delegator rewards.
type DistributionKeeper interface {
	GetValidatorCurrentRewards(ctx context.Context, val sdk.ValAddress) (rewards distributiontypes.ValidatorCurrentRewards, err error)
	GetValidatorAccumulatedCommission(ctx context.Context, val sdk.ValAddress) (commission distributiontypes.ValidatorAccumulatedCommission, err error)
	CalculateDelegationRewards(ctx context.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins, err error)

	WithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
	WithdrawValidatorCommission(ctx context.Context, valAddr sdk.ValAddress) (sdk.Coins, error)

	IncrementValidatorPeriod(ctx context.Context, val stakingtypes.ValidatorI) (uint64, error)

	GetUbiBalanceByDenom(ctx context.Context, denom string) (math.Int, error)
	WithdrawUbiByDenomToModule(ctx context.Context, denom string, recipientModule string) (sdk.Coin, error)
}
