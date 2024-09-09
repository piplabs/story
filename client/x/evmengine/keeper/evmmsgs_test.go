package keeper

import (
	"bytes"
	"math/big"
	"slices"
	"testing"

	"github.com/cometbft/cometbft/crypto"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/piplabs/story/client/genutil/evm/predeploys"
	moduletestutil "github.com/piplabs/story/client/x/evmengine/testutil"
	"github.com/piplabs/story/client/x/evmengine/types"
	evmstakingtypes "github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/ethclient/mock"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/tutil"
)

var (
	// dummy block hashes for testing
	unstakeBlkHash = common.HexToHash("0x1398C32A45Bc409b6C652E25bb0a3e702492A4ab")
	unjailBlkHash  = common.HexToHash("0x1398C32A45Bc409b6C652E25bb0a3e702492A4ac")
	mixedBlkHash   = common.HexToHash("0x1398C32A45Bc409b6C652E25bb0a3e702492A4ad")
	failedBlock    = common.HexToHash("0x1398C32A45Bc409b6C652E25bb0a3e702492A4ae")

	// pre-deploy contract addresses
	stakingAddr  = common.HexToAddress(predeploys.IPTokenStaking)
	slashingAddr = common.HexToAddress(predeploys.IPTokenSlashing)
)

var (
	delPubKey   crypto.PubKey
	valPubKey   crypto.PubKey
	valEvmAddr  common.Address
	unstakeData []byte
	unjailData  []byte
)

// TestKeeper_evmEvents tests the evmEvents function of the keeper.
func TestKeeper_evmEvents(t *testing.T) {
	t.Parallel()
	keeper, ctx, _, _ := setupTestEnvironment(t)

	tcs := []struct {
		name           string
		blkHash        common.Hash
		expectedResult []*types.EVMEvent
		expectedError  string
	}{
		{
			name:          "fail: failed block hash",
			blkHash:       failedBlock,
			expectedError: "filter logs: failed to fetch logs",
		},
		{
			name: "pass: no logs",
			blkHash: common.HexToHash(
				"0x1398C32A45Bc409b6C652E25bb0a3e702492A4aa",
			),
			expectedResult: []*types.EVMEvent{},
		},
		{
			name:    "pass: only IPTokenStaking WithdrawEvent log",
			blkHash: unstakeBlkHash,
			expectedResult: []*types.EVMEvent{
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{evmstakingtypes.WithdrawEvent.ID.Bytes()},
					Data:    unstakeData,
				},
			},
		},
		{
			name:    "pass: only IPTokenSlashing UnjailEvent log",
			blkHash: unjailBlkHash,
			expectedResult: []*types.EVMEvent{
				{
					Address: slashingAddr.Bytes(),
					Topics:  [][]byte{evmstakingtypes.UnjailEvent.ID.Bytes(), common.BytesToHash(valEvmAddr.Bytes()).Bytes()},
					Data:    unjailData,
				},
			},
		},
		{
			name:    "pass: mixed logs",
			blkHash: mixedBlkHash,
			expectedResult: []*types.EVMEvent{
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{evmstakingtypes.WithdrawEvent.ID.Bytes()},
					Data:    unstakeData,
				},
				{
					Address: slashingAddr.Bytes(),
					Topics:  [][]byte{evmstakingtypes.UnjailEvent.ID.Bytes(), common.BytesToHash(valEvmAddr.Bytes()).Bytes()},
					Data:    unjailData,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			events, err := keeper.evmEvents(ctx, tc.blkHash)
			if tc.expectedError != "" {
				require.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, events)
				require.True(t, isSorted(events), "events are not sorted")
			}
		})
	}
}

func setupTestEnvironment(t *testing.T) (*Keeper, sdk.Context, *gomock.Controller, *moduletestutil.MockUpgradeKeeper) {
	t.Helper()
	cdc := getCodec(t)
	txConfig := authtx.NewTxConfig(cdc, nil)

	pubKeys, _, _ := createAddresses(3)
	delPubKey = pubKeys[0]
	valPubKey = pubKeys[1]
	valEvmAddr = common.BytesToAddress(valPubKey.Address())

	var err error
	// setup engine mock
	stakingAbi := mustGetABI(bindings.IPTokenStakingMetaData)
	slashingAbi := mustGetABI(bindings.IPTokenSlashingMetaData)
	amt := int64(100)
	unstakeData, err = stakingAbi.Events["Withdraw"].Inputs.NonIndexed().Pack(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).Mul(big.NewInt(amt), big.NewInt(params.Ether)))
	require.NoError(t, err)
	unjailData, err = slashingAbi.Events["Unjail"].Inputs.NonIndexed().Pack(valPubKey.Bytes())
	mockEngine, err := ethclient.NewEngineMock(
		ethclient.WithMockUnstake(
			unstakeBlkHash, stakingAddr, delPubKey.Bytes(), valPubKey.Bytes(), amt,
		),
		ethclient.WithMockUnjail(
			unjailBlkHash, slashingAddr, valEvmAddr, valPubKey.Bytes(),
		),
		ethclient.WithMockUnstakeAndUnjail(
			mixedBlkHash, stakingAddr, slashingAddr, valEvmAddr, delPubKey.Bytes(), valPubKey.Bytes(), amt,
		),
		ethclient.WithMockFailedOnBlockHashes([]common.Hash{failedBlock}),
	)
	require.NoError(t, err)

	cmtAPI := newMockCometAPI(t, nil)
	header := cmtproto.Header{Height: 1, AppHash: tutil.RandomHash().Bytes(), ProposerAddress: cmtAPI.validatorSet.Validators[0].Address}
	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)
	ak := moduletestutil.NewMockAccountKeeper(ctrl)
	esk := moduletestutil.NewMockEvmStakingKeeper(ctrl)
	uk := moduletestutil.NewMockUpgradeKeeper(ctrl)

	ctx, storeService := setupCtxStore(t, &header)

	keeper, err := NewKeeper(cdc, storeService, mockEngine, mockClient, txConfig, ak, esk, uk)
	require.NoError(t, err)
	keeper.SetCometAPI(cmtAPI)
	nxtAddr, err := k1util.PubKeyToAddress(cmtAPI.validatorSet.Validators[1].PubKey)
	require.NoError(t, err)
	keeper.SetValidatorAddress(nxtAddr)
	populateGenesisHead(ctx, t, keeper)

	return keeper, ctx, ctrl, uk
}

// mustGetABI returns the metadata's ABI as an abi.ABI type.
// It panics on error.
func mustGetABI(metadata *bind.MetaData) *abi.ABI {
	abi, err := metadata.GetAbi()
	if err != nil {
		panic(err)
	}

	return abi
}

// Helper function to check if the events are sorted by ascending order of address, topics, and data.
func isSorted(events []*types.EVMEvent) bool {
	for i := 1; i < len(events); i++ {
		// Compare addresses first
		addressComparison := bytes.Compare(events[i-1].Address, events[i].Address)
		if addressComparison > 0 {
			// it is not sorted by ascending order of address
			return false
		}

		if addressComparison == 0 {
			// If addresses are equal, compare by topics
			previousTopic := slices.Concat(events[i-1].Topics...)
			currentTopic := slices.Concat(events[i].Topics...)
			topicComparison := bytes.Compare(previousTopic, currentTopic)

			if topicComparison > 0 {
				// it is not sorted by ascending order of topics
				return false
			}

			if topicComparison == 0 {
				// If topics are also equal, compare by data
				if bytes.Compare(events[i-1].Data, events[i].Data) > 0 {
					// it is not sorted by ascending order of data
					return false
				}
			}
		}
	}
	return true
}
