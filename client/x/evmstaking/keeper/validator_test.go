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
	ctx, eskeeper := s.Ctx, s.EVMStakingKeeper

	pubKeys, _, _ := createAddresses(3)

	corruptedPubKey := append([]byte{}, pubKeys[0].Bytes()...)
	corruptedPubKey[0] = 0x09
	corruptedPubKey[1] = 0xFF

	tcs := []struct {
		name           string
		valDelAddr     sdk.AccAddress
		valAddr        sdk.ValAddress
		valPubKey      crypto.PubKey
		valPubKeyBytes []byte
		valCmpPubKey   []byte
		valTokens      math.Int
		moniker        string
		preRun         func(t *testing.T, c sdk.Context, valDelAddr sdk.AccAddress, valPubKey crypto.PubKey, valTokens math.Int)
		postCheck      func(c sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, delEvmAddr common.Address, previousTokens math.Int)
		expectedError  string
	}{
		{
			name:          "fail: nil validator pubkey",
			valCmpPubKey:  nil,
			expectedError: "compress validator pubkey: invalid uncompressed public key length or format",
		},
		{
			name:          "fail: invalid validator pubkey",
			valCmpPubKey:  pubKeys[0].Bytes()[1:],
			expectedError: "compress validator pubkey: invalid uncompressed public key length or format",
		},
		{
			name:          "fail: corrupted validator pubkey",
			valCmpPubKey:  corruptedPubKey,
			expectedError: "validator pubkey to evm address: invalid public key",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			// create valPubKey using tc.valPubKeyBytes
			if tc.preRun != nil {
				tc.preRun(s.T(), cachedCtx, tc.valDelAddr, tc.valPubKey, tc.valTokens)
			}
			moniker := tc.moniker
			if moniker == "" {
				moniker = "testing"
			}
			err := eskeeper.ProcessCreateValidator(cachedCtx, &bindings.IPTokenStakingCreateValidator{
				ValidatorCmpPubkey:      tc.valCmpPubKey,
				Moniker:                 moniker,
				StakeAmount:             new(big.Int).SetUint64(100),
				CommissionRate:          1000, // 10%
				MaxCommissionRate:       5000, // 50%
				MaxCommissionChangeRate: 500,  // 5%
				SupportsUnlocked:        uint8(0),
				Raw:                     gethtypes.Log{},
			})
			if tc.expectedError != "" {
				require.Error(err, tc.expectedError)
			} else {
				require.NoError(err)
				delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(tc.valPubKeyBytes)
				require.NoError(err)
				tc.postCheck(cachedCtx, tc.valDelAddr, tc.valAddr, delEvmAddr, tc.valTokens)
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
