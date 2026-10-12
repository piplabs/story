//nolint:staticcheck,revive,nolintlint // versioning
package v_1_8_1

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper is the slice of the bank keeper the sweep needs.
type BankKeeper interface {
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

// DistributionKeeper is the slice of the distribution keeper the sweep needs.
// FundUbi moves the coins into the distribution module and credits the fee
// pool, which is the accounted way to hand it protocol revenue.
type DistributionKeeper interface {
	FundUbi(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
