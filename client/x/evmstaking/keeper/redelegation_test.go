package keeper_test

import (
	"math/big"

	"github.com/cometbft/cometbft/crypto"
	k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

func createAddresses(count int) ([]crypto.PubKey, []sdk.AccAddress, []sdk.ValAddress) {
	var pubKeys []crypto.PubKey
	var accAddrs []sdk.AccAddress
	var valAddrs []sdk.ValAddress
	for range count {
		pubKey := k1.GenPrivKey().PubKey()
		accAddr := sdk.AccAddress(pubKey.Address().Bytes())
		valAddr := sdk.ValAddress(pubKey.Address().Bytes())
		pubKeys = append(pubKeys, pubKey)
		accAddrs = append(accAddrs, accAddr)
		valAddrs = append(valAddrs, valAddr)
	}

	return pubKeys, accAddrs, valAddrs
}

func (s *TestSuite) TestRedelegation() {
	ctx, keeper, stakingKeeper := s.Ctx, s.EVMStakingKeeper, s.StakingKeeper
	require := s.Require()

	// create addresses
	pubKeys, accAddrs, valAddrs := createAddresses(3)
	delAddr := accAddrs[0]
	valSrcAddr := valAddrs[1]
	valDstAddr := valAddrs[2]

	// create a validator (src)
	validatorSrc := testutil.NewValidator(s.T(), valSrcAddr, PKs[1])
	require.NoError(stakingKeeper.SetValidatorByConsAddr(ctx, validatorSrc))

	// delegate to the validator
	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	validator, issuedShares := validatorSrc.AddTokensFromDel(valTokens)
	require.Equal(valTokens, issuedShares.RoundInt())
	s.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
	_ = skeeper.TestingUpdateValidator(stakingKeeper, ctx, validator, true)
	delegation := stypes.NewDelegation(delAddr.String(), valSrcAddr.String(), issuedShares)
	require.NoError(stakingKeeper.SetDelegation(ctx, delegation))
	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(pubKeys[0].Bytes())
	require.NoError(err)
	require.NoError(keeper.DelegatorMap.Set(ctx, delAddr.String(), delEvmAddr.String()))

	// create a second validator(dst) and delegate the same amount of token
	validatorDst := testutil.NewValidator(s.T(), valDstAddr, PKs[1])
	validatorDst, issuedShares = validatorDst.AddTokensFromDel(valTokens)
	require.Equal(valTokens, issuedShares.RoundInt())
	s.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
	_ = skeeper.TestingUpdateValidator(stakingKeeper, ctx, validatorDst, true)
	delegation = stypes.NewDelegation(delAddr.String(), valDstAddr.String(), issuedShares)
	require.NoError(stakingKeeper.SetDelegation(ctx, delegation))

	// check the amount of delegated tokens
	delSrc, err := stakingKeeper.GetDelegatorValidator(ctx, delAddr, valSrcAddr)
	require.NoError(err)
	require.True(delSrc.Tokens.Equal(valTokens))

	delDst, err := stakingKeeper.GetDelegatorValidator(ctx, delAddr, valDstAddr)
	require.NoError(err)
	require.True(delDst.Tokens.Equal(valTokens))

	// test shouldn't have and redelegations
	has, err := stakingKeeper.HasReceivingRedelegation(ctx, delAddr, valDstAddr)
	require.NoError(err)
	require.False(has)

	redelTokens := stakingKeeper.TokensFromConsensusPower(ctx, 5)

	ipTokenRedelegate := &bindings.IPTokenStakingRedelegate{
		DepositorPubkey:    pubKeys[0].Bytes(),
		ValidatorSrcPubkey: pubKeys[1].Bytes(),
		ValidatorDstPubkey: pubKeys[2].Bytes(),
		Amount:             big.NewInt(redelTokens.Int64()), // multiply power reduction of 1000000
		Raw:                gethtypes.Log{},
	}

	// redelegation
	require.NoError(keeper.ProcessRedelegate(ctx, ipTokenRedelegate))

	// check the amount of delegated tokens after redelegation
	delSrc, err = stakingKeeper.GetDelegatorValidator(ctx, delAddr, valSrcAddr)
	require.NoError(err)
	require.True(delSrc.Tokens.Equal(valTokens.Sub(redelTokens)))

	delDst, err = stakingKeeper.GetDelegatorValidator(ctx, delAddr, valDstAddr)
	require.NoError(err)
	require.True(delDst.Tokens.Equal(valTokens.Add(redelTokens)))

	// params
	params, err := s.StakingKeeper.GetParams(ctx)
	require.NoError(err)

	redelegation, err := stakingKeeper.GetRedelegation(ctx, delAddr, valSrcAddr, valDstAddr)
	require.NoError(err)
	require.Equal(delAddr.String(), redelegation.DelegatorAddress)
	require.Equal(valSrcAddr.String(), redelegation.ValidatorSrcAddress)
	require.Equal(valDstAddr.String(), redelegation.ValidatorDstAddress)
	require.Equal(redelTokens, redelegation.Entries[0].InitialBalance)
	require.Equal(ctx.BlockTime().Add(params.UnbondingTime), redelegation.Entries[0].CompletionTime)

	// TODO: test EndBlock
}
