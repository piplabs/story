package horace

import (
	"context"
	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	minttypes "github.com/piplabs/story/client/x/mint/types"
)

// StakingKeeper is a test-only interface for upgrade handlers.
type StakingKeeper interface {
	GetParams(ctx context.Context) (stypes.Params, error)
	SetParams(ctx context.Context, params stypes.Params) error

	GetAllValidators(ctx context.Context) ([]stypes.Validator, error)
	GetValidator(ctx context.Context, valAddr sdk.ValAddress) (stypes.Validator, error)
	SetValidator(ctx context.Context, val stypes.Validator) error

	GetValidatorDelegations(ctx context.Context, valAddr sdk.ValAddress) ([]stypes.Delegation, error)
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stypes.Delegation, error)
	SetDelegation(ctx context.Context, del stypes.Delegation) error

	GetPeriodDelegation(
		ctx context.Context,
		delAddr sdk.AccAddress,
		valAddr sdk.ValAddress,
		periodDelegationID string,
	) (stypes.PeriodDelegation, error)

	SetPeriodDelegation(
		ctx context.Context,
		delAddr sdk.AccAddress,
		valAddr sdk.ValAddress,
		pd stypes.PeriodDelegation,
	) error
}

// DistributionKeeper is a test-only interface for upgrade handlers.
type DistributionKeeper interface {
	WithdrawValidatorCommission(
		ctx context.Context,
		valAddr sdk.ValAddress,
	) (sdk.Coins, error)

	WithdrawDelegationRewards(
		ctx context.Context,
		delAddr sdk.AccAddress,
		valAddr sdk.ValAddress,
	) (sdk.Coins, error)

	GetDelegatorStartingInfo(
		ctx context.Context,
		valAddr sdk.ValAddress,
		delAddr sdk.AccAddress,
	) (dtypes.DelegatorStartingInfo, error)

	SetDelegatorStartingInfo(
		ctx context.Context,
		valAddr sdk.ValAddress,
		delAddr sdk.AccAddress,
		info dtypes.DelegatorStartingInfo,
	) error
}

// MintKeeper is a test-only interface for upgrade handlers.
type MintKeeper interface {
	GetParams(ctx context.Context) (minttypes.Params, error)
	SetParams(ctx context.Context, params minttypes.Params) error
}

// AccountKeeper is a test-only interface for upgrade handlers.
type AccountKeeper interface {
	AddressCodec() address.Codec
}
