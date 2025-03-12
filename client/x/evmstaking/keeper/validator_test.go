package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/k1util"
)

func (s *TestSuite) TestProcessCreateValidator() {
	require := s.Require()
	ctx, skeeper, eskeeper := s.Ctx, s.StakingKeeper, s.EVMStakingKeeper

	pubKeys, _, valAddrs := createAddresses(3)

	corruptedPubKey := append([]byte{}, pubKeys[0].Bytes()...)
	corruptedPubKey[0] = 0x09
	corruptedPubKey[1] = 0xFF

	valPubkey := pubKeys[0]
	valPubkeyBytes := valPubkey.Bytes()
	valAddr := valAddrs[0]
	valDelAddr := sdk.AccAddress(valAddr.Bytes())

	stakingParams, err := skeeper.GetParams(ctx)
	require.NoError(err)

	// set min delegation to 1024 for all test cases
	stakingParams.MinDelegation = math.NewInt(1024)
	require.NoError(skeeper.SetParams(ctx, stakingParams))

	refundFeeBps, err := eskeeper.RefundFeeBps(ctx)
	require.NoError(err)
	refundPeriod, err := eskeeper.RefundPeriod(ctx)
	require.NoError(err)

	tcs := []struct {
		name           string
		valDelAddr     sdk.AccAddress
		valAddr        sdk.ValAddress
		valPubKey      crypto.PubKey
		valPubKeyBytes []byte
		valCmpPubKey   []byte // In fail cases, use this to test the compressed pubkey. For pass cases, don't set (i.e. nil)
		valTokens      math.Int
		stakeAmount    *big.Int
		moniker        string
		preRun         func(t *testing.T, c sdk.Context, valDelAddr sdk.AccAddress, valPubKey crypto.PubKey, valTokens math.Int) sdk.Context
		postCheck      func(t *testing.T, c sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, previousTokens math.Int, stakeAmount *big.Int)
		expectedError  string
	}{
		{
			name:           "fail: nil validator pubkey",
			valPubKeyBytes: valPubkeyBytes, // must be valid value but not used in ProcessCreateValidator
			valCmpPubKey:   nil,
			expectedError:  "compress validator pubkey: invalid uncompressed public key length or format",
		},
		{
			name:           "fail: invalid validator pubkey",
			valPubKeyBytes: valPubkeyBytes, // must be valid value but not used in ProcessCreateValidator
			valCmpPubKey:   pubKeys[0].Bytes()[1:],
			expectedError:  "compress validator pubkey: invalid uncompressed public key length or format",
		},
		{
			name:           "fail: corrupted validator pubkey",
			valPubKeyBytes: valPubkeyBytes, // must be valid value but not used in ProcessCreateValidator
			valCmpPubKey:   corruptedPubKey,
			expectedError:  "validator pubkey to evm address: invalid public key",
		},
		{
			name:           "fail: minimum stake amount",
			valPubKey:      valPubkey,
			valPubKeyBytes: valPubkeyBytes,
			valTokens:      math.NewInt(1),
			stakeAmount:    big.NewInt(1),
			expectedError:  "IPTokenStaking: Stake amount under min",
			preRun: func(t *testing.T, c sdk.Context, valDelAddr sdk.AccAddress, valPubKey crypto.PubKey, valTokens math.Int) sdk.Context {
				t.Helper()
				p, err := skeeper.GetParams(ctx)
				require.NoError(err)
				require.True(p.MinDelegation.GT(math.NewInt(1)))
				// skip to first block since the "min self delegation < min delegation" check is skipped in the genesis block
				return c.WithBlockHeight(1)
			},
		},
		{
			name:           "pass: create new validator",
			valDelAddr:     valDelAddr,
			valAddr:        valAddr,
			valPubKey:      valPubkey,
			valPubKeyBytes: valPubkeyBytes,
			valTokens:      math.NewInt(1024),
			stakeAmount:    big.NewInt(1024),
			postCheck: func(t *testing.T, c sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, previousTokens math.Int, stakeAmount *big.Int) {
				t.Helper()
				val, err := skeeper.GetValidator(c, valAddr)
				require.NoError(err)
				require.Equal(math.NewInt(1024), val.MinSelfDelegation)
				require.Equal(big.NewInt(1024), val.Tokens.BigInt())

				del, err := skeeper.GetDelegation(c, delAddr, valAddr)
				require.NoError(err)
				require.Equal(valAddr.String(), del.ValidatorAddress)
				require.Equal(math.LegacyNewDecFromBigInt(stakeAmount), del.Shares)

				delVal, err := skeeper.GetDelegatorValidator(c, delAddr, valAddr)
				require.NoError(err)
				valConstPubKey, err := val.ConsPubKey()
				require.NoError(err)
				delValConstPubKey, err := delVal.ConsPubKey()
				require.NoError(err)
				require.Equal(valConstPubKey, delValConstPubKey)
			},
		},
		{
			name:           "pass: create existing validator, convert to self-delegation",
			valDelAddr:     valDelAddr,
			valAddr:        valAddr,
			valPubKey:      valPubkey,
			valPubKeyBytes: valPubkeyBytes,
			valTokens:      math.NewInt(1500),
			stakeAmount:    big.NewInt(1500),
			postCheck: func(t *testing.T, c sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, previousTokens math.Int, stakeAmount *big.Int) {
				t.Helper()
				_, err := skeeper.GetValidator(c, valAddr)
				require.NoError(err)

				del, err := skeeper.GetDelegation(c, delAddr, valAddr)
				require.NoError(err)
				require.Equal(valAddr.String(), del.ValidatorAddress)
				require.Equal(math.LegacyNewDecFromBigInt(stakeAmount), del.Shares)

				// create validator again with the same address
				err = eskeeper.ProcessCreateValidator(c, &bindings.IPTokenStakingCreateValidator{
					ValidatorCmpPubkey:      valPubkeyBytes,
					Moniker:                 "testing",
					StakeAmount:             stakeAmount,
					CommissionRate:          1000,
					MaxCommissionRate:       5000,
					MaxCommissionChangeRate: 500,
					SupportsUnlocked:        uint8(0), // unlocked token
					OperatorAddress:         common.Address{},
					Raw:                     gethtypes.Log{},
				})
				require.NoError(err)

				// expect self-delegation to increase by stake amount
				delAfter, err := skeeper.GetDelegation(c, delAddr, valAddr)
				require.NoError(err)
				require.Equal(valAddr.String(), delAfter.ValidatorAddress)
				require.Equal(math.LegacyNewDecFromBigInt(stakeAmount).Add(del.Shares), delAfter.Shares)
			},
		},
		{
			name:           "pass: create existing validator, mismatch in stake token type, convert to refund",
			valDelAddr:     valDelAddr,
			valAddr:        valAddr,
			valPubKey:      valPubkey,
			valPubKeyBytes: valPubkeyBytes,
			valTokens:      math.NewInt(2000),
			stakeAmount:    big.NewInt(2000),
			postCheck: func(t *testing.T, c sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, previousTokens math.Int, stakeAmount *big.Int) {
				t.Helper()
				_, err := skeeper.GetValidator(c, valAddr)
				require.NoError(err)

				del, err := skeeper.GetDelegation(c, delAddr, valAddr)
				require.NoError(err)
				require.Equal(valAddr.String(), del.ValidatorAddress)
				require.Equal(math.LegacyNewDecFromBigInt(stakeAmount), del.Shares)

				// create validator again with the same address, but locked token
				err = eskeeper.ProcessCreateValidator(c, &bindings.IPTokenStakingCreateValidator{
					ValidatorCmpPubkey:      valPubkeyBytes,
					Moniker:                 "testing",
					StakeAmount:             stakeAmount,
					CommissionRate:          1000,
					MaxCommissionRate:       5000,
					MaxCommissionChangeRate: 500,
					SupportsUnlocked:        uint8(1), // locked token
					OperatorAddress:         common.Address{},
					Raw:                     gethtypes.Log{},
				})
				require.NoError(err)

				// expect a refund in queue
				ubd, err := skeeper.GetUnbondingDelegation(c, delAddr, valAddr)
				require.NoError(err)
				require.Equal(valAddr.String(), ubd.ValidatorAddress)
				require.Equal(delAddr.String(), ubd.DelegatorAddress)

				completionTime := ctx.BlockTime().Add(refundPeriod)
				expectedRefundAmount := math.NewInt(2000).Sub(math.NewInt(2000).Mul(math.NewInt(int64(refundFeeBps))).Quo(math.NewInt(10_000)))
				for _, entry := range ubd.Entries {
					require.True(expectedRefundAmount.Equal(math.NewIntFromBigInt(entry.Balance.BigInt())))
					require.Equal(completionTime, entry.CompletionTime)
				}
			},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			// cache ctx after preRun because preRun may change the block height
			cachedCtx, _ := ctx.CacheContext()
			if tc.preRun != nil {
				cachedCtx = tc.preRun(s.T(), cachedCtx, tc.valDelAddr, tc.valPubKey, tc.valTokens)
			}
			moniker := tc.moniker
			if moniker == "" {
				moniker = "testing"
			}
			delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(tc.valPubKeyBytes)
			require.NoError(err)
			valCmpPubKey := tc.valPubKeyBytes
			if tc.valCmpPubKey != nil { // used for testing fail cases
				valCmpPubKey = tc.valCmpPubKey
			}
			err = eskeeper.ProcessCreateValidator(cachedCtx, &bindings.IPTokenStakingCreateValidator{
				ValidatorCmpPubkey:      valCmpPubKey,
				Moniker:                 moniker,
				StakeAmount:             tc.stakeAmount,
				CommissionRate:          1000, // 10%
				MaxCommissionRate:       5000, // 50%
				MaxCommissionChangeRate: 500,  // 5%
				SupportsUnlocked:        uint8(0),
				OperatorAddress:         delEvmAddr,
				Raw:                     gethtypes.Log{},
			})
			if tc.expectedError != "" {
				require.Error(err, tc.expectedError)
			} else {
				require.NoError(err)
				if tc.postCheck != nil {
					tc.postCheck(s.T(), cachedCtx, tc.valDelAddr, tc.valAddr, tc.valTokens, tc.stakeAmount)
				}
			}
		})
	}
}

func (s *TestSuite) TestParseCreateValidatorLog() {
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
				Topics: []common.Hash{types.CreateValidatorEvent.ID},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := keeper.ParseCreateValidatorLog(tc.log)
			if tc.expectErr {
				require.Error(err, "should return error for %s", tc.name)
			} else {
				require.NoError(err, "should not return error for %s", tc.name)
			}
		})
	}
}
