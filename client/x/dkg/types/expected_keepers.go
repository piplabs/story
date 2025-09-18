package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type AccountKeeper interface {
	AddressCodec() address.Codec
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type StakingKeeper interface {
	GetAllValidators(ctx context.Context) (validators []stakingtypes.Validator, err error)
}
