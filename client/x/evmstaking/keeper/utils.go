package keeper

import (
	"math/big"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IPTokenToBondCoin converts the IP amount into a $STAKE coin.
// NOTE: Assumes that the only bondable token is $STAKE on CL (using IP token).
func IPTokenToBondCoin(amount *big.Int) (sdk.Coin, sdk.Coins) {
	coin := sdk.NewCoin(sdk.DefaultBondDenom, math.NewIntFromBigInt(amount))
	return coin, sdk.NewCoins(coin)
}
