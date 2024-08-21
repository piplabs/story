package keeper

import (
	"math/big"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IPTokenToBondCoin converts the IP amount into a $STAKE coin.
// TODO: At this point, it is 1-to1, but this might change in the future.
// TODO: parameterized bondDenom.
func IPTokenToBondCoin(amount *big.Int) (sdk.Coin, sdk.Coins) {
	coin := sdk.NewCoin(sdk.DefaultBondDenom, math.NewIntFromBigInt(amount))
	return coin, sdk.NewCoins(coin)
}
