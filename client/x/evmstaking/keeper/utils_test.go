package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/keeper"
)

func TestIPTokenToBondCoin(t *testing.T) {
	t.Parallel()

	amount := big.NewInt(100)
	coin, coins := keeper.IPTokenToBondCoin(amount)
	require.Equal(t, sdk.DefaultBondDenom, coin.Denom)
	require.Equal(t, math.NewIntFromBigInt(amount), coin.Amount)
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, math.NewIntFromBigInt(amount)), coin)
	require.Equal(t, sdk.NewCoins(coin), coins)
}
