package keeper_test

import (
	"context"
	"math/big"

	"github.com/cometbft/cometbft/crypto"
	k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
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
	delPubKey := pubKeys[0]
	valSrcPubKey := pubKeys[1]
	valSrcAddr := valAddrs[1]
	valDstPubKey := pubKeys[2]
	valDstAddr := valAddrs[2]

	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	s.setupValidatorAndDelegation(ctx, valSrcPubKey, delPubKey, valSrcAddr, delAddr, valTokens)
	s.setupValidatorAndDelegation(ctx, valDstPubKey, delPubKey, valDstAddr, delAddr, valTokens)

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

	redelTokens := stakingKeeper.TokensFromConsensusPower(ctx, 5) // multiply power reduction of 1000000
	validInput := &bindings.IPTokenStakingRedelegate{
		DelegatorCmpPubkey: delPubKey.Bytes(),
		ValidatorSrcPubkey: valSrcPubKey.Bytes(),
		ValidatorDstPubkey: valDstPubKey.Bytes(),
		Amount:             big.NewInt(redelTokens.Int64()),
	}
	checkStateAfterRedelegation := func(c context.Context) {
		// check the amount of delegated tokens after redelegation
		delSrc, err = stakingKeeper.GetDelegatorValidator(c, delAddr, valSrcAddr)
		require.NoError(err)
		require.True(delSrc.Tokens.Equal(valTokens.Sub(redelTokens)))

		delDst, err = stakingKeeper.GetDelegatorValidator(c, delAddr, valDstAddr)
		require.NoError(err)
		require.True(delDst.Tokens.Equal(valTokens.Add(redelTokens)))

		// params
		params, err := s.StakingKeeper.GetParams(c)
		require.NoError(err)

		redelegation, err := stakingKeeper.GetRedelegation(c, delAddr, valSrcAddr, valDstAddr)
		require.NoError(err)
		require.Equal(delAddr.String(), redelegation.DelegatorAddress)
		require.Equal(valSrcAddr.String(), redelegation.ValidatorSrcAddress)
		require.Equal(valDstAddr.String(), redelegation.ValidatorDstAddress)
		require.Equal(redelTokens, redelegation.Entries[0].InitialBalance)
		sdkCtx := sdk.UnwrapSDKContext(c)
		require.Equal(sdkCtx.BlockTime().Add(params.UnbondingTime), redelegation.Entries[0].CompletionTime)
	}

	tcs := []struct {
		name          string
		input         func() bindings.IPTokenStakingRedelegate
		expectedError string
		// postCheck checks the state is changed after the successful operation
		postCheck func(c context.Context)
	}{
		{
			name: "pass: valid redelegation",
			input: func() bindings.IPTokenStakingRedelegate {
				return *validInput
			},
			postCheck: checkStateAfterRedelegation,
		},
		{
			name: "fail: zero amount",
			input: func() bindings.IPTokenStakingRedelegate {
				inputCpy := *validInput
				inputCpy.Amount = big.NewInt(0)

				return inputCpy
			},
			expectedError: "invalid shares amount",
		},
		{
			name: "fail: invalid delegator pubkey",
			input: func() bindings.IPTokenStakingRedelegate {
				inputCpy := *validInput
				inputCpy.DelegatorCmpPubkey = delPubKey.Bytes()[1:]

				return inputCpy
			},
			expectedError: "depositor pubkey to cosmos",
		},
		{
			name: "fail: invalid src validator pubkey",
			input: func() bindings.IPTokenStakingRedelegate {
				inputCpy := *validInput
				inputCpy.ValidatorSrcPubkey = valSrcPubKey.Bytes()[1:]

				return inputCpy
			},
			expectedError: "src validator pubkey to cosmos",
		},
		{
			name: "fail: invalid dst validator pubkey",
			input: func() bindings.IPTokenStakingRedelegate {
				inputCpy := *validInput
				inputCpy.ValidatorDstPubkey = valDstPubKey.Bytes()[1:]

				return inputCpy
			},
			expectedError: "dst validator pubkey to cosmos",
		},
		{
			name: "fail: corrupted delegator pubkey",
			input: func() bindings.IPTokenStakingRedelegate {
				inputCpy := *validInput
				inputCpy.DelegatorCmpPubkey = createCorruptedPubKey(delPubKey.Bytes())

				return inputCpy
			},
			expectedError: "deledator pubkey to evm address",
		},
		{
			name: "fail: corrupted src validator pubkey",
			input: func() bindings.IPTokenStakingRedelegate {
				inputCpy := *validInput
				inputCpy.ValidatorSrcPubkey = createCorruptedPubKey(valSrcPubKey.Bytes())

				return inputCpy
			},
			expectedError: "src validator pubkey to evm address",
		},
		{
			name: "fail: corrupted dst validator pubkey",
			input: func() bindings.IPTokenStakingRedelegate {
				inputCpy := *validInput
				inputCpy.ValidatorDstPubkey = createCorruptedPubKey(valDstPubKey.Bytes())

				return inputCpy
			},
			expectedError: "dst validator pubkey to evm address",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			input := tc.input()
			err := keeper.HandleRedelegateEvent(cachedCtx, &input)
			if tc.expectedError != "" {
				require.ErrorContains(err, tc.expectedError)
			} else {
				require.NoError(err, tc.expectedError)
				tc.postCheck(cachedCtx)
			}
		})
	}
}

func (s *TestSuite) TestParseRedelegationLog() {
	require := s.Require()
	keeper := s.EVMStakingKeeper

	testCases := []struct {
		name      string
		log       gethtypes.Log
		expectErr bool
	}{
		{
			name: "Unknown Topic",
			log: gethtypes.Log{
				Topics: []common.Hash{common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")},
			},
			expectErr: true,
		},
		{
			name: "Valid Topic",
			log: gethtypes.Log{
				Topics: []common.Hash{types.RedelegateEvent.ID},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := keeper.ParseRedelegateLog(tc.log)
			if tc.expectErr {
				require.Error(err, "should return error for %s", tc.name)
			} else {
				require.NoError(err, "should not return error for %s", tc.name)
			}
		})
	}
}
