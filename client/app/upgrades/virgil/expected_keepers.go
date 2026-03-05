package virgil

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper is a test-only interface for the virgil upgrade handler.
type StakingKeeper interface {
	GetParams(ctx context.Context) (stypes.Params, error)
	SetParams(ctx context.Context, params stypes.Params) error

	GetAllPeriodDelegations(ctx context.Context) ([]stypes.PeriodDelegation, error)

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
