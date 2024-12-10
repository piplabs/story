package keeper

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/comet"
	moduletestutil "github.com/piplabs/story/client/x/evmengine/testutil"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/ethclient/mock"
	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

type args struct {
	height         int64
	validatorsFunc func(context.Context, int64) (*cmttypes.ValidatorSet, error)
	isNextProposer bool
	header         func(height int64) cmtproto.Header
}

func createTestKeeper(t *testing.T) (context.Context, *Keeper) {
	t.Helper()

	cdc := getCodec(t)
	txConfig := authtx.NewTxConfig(cdc, nil)

	cmtAPI := newMockCometAPI(t, nil)
	header := cmtproto.Header{Height: 1}

	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)
	ak := moduletestutil.NewMockAccountKeeper(ctrl)
	esk := moduletestutil.NewMockEvmStakingKeeper(ctrl)
	uk := moduletestutil.NewMockUpgradeKeeper(ctrl)
	dk := moduletestutil.NewMockDistrKeeper(ctrl)

	ctx, storeKey, storeService := setupCtxStore(t, &header)
	mockEngine, err := newMockEngineAPI(storeKey, 0)
	require.NoError(t, err)

	keeper, err := NewKeeper(cdc, storeService, &mockEngine, mockClient, txConfig, ak, esk, uk, dk)
	require.NoError(t, err)
	keeper.SetCometAPI(cmtAPI)

	return ctx, keeper
}

func createKeeper(t *testing.T, args args) (sdk.Context, *mockCometAPI, *Keeper) {
	t.Helper()

	cdc := getCodec(t)
	txConfig := authtx.NewTxConfig(cdc, nil)

	cmtAPI := newMockCometAPI(t, args.validatorsFunc)
	header := args.header(args.height)

	var nxtAddr common.Address
	var err error
	if args.isNextProposer {
		nxtAddr, err = k1util.PubKeyToAddress(cmtAPI.validatorSet.CopyIncrementProposerPriority(1).Proposer.PubKey)
	} else {
		nxtAddr = common.HexToAddress("0x0000000000000000000000000000000000000000")
	}
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)
	ak := moduletestutil.NewMockAccountKeeper(ctrl)
	esk := moduletestutil.NewMockEvmStakingKeeper(ctrl)
	uk := moduletestutil.NewMockUpgradeKeeper(ctrl)
	dk := moduletestutil.NewMockDistrKeeper(ctrl)

	ctx, storeKey, storeService := setupCtxStore(t, &header)
	mockEngine, err := newMockEngineAPI(storeKey, 0)
	require.NoError(t, err)
	keeper, err := NewKeeper(cdc, storeService, &mockEngine, mockClient, txConfig, ak, esk, uk, dk)
	require.NoError(t, err)
	keeper.SetCometAPI(cmtAPI)
	keeper.SetValidatorAddress(nxtAddr)
	populateGenesisHead(ctx, t, keeper)

	return ctx, cmtAPI, keeper
}

func TestKeeper_SetBuildDelay(t *testing.T) {
	t.Parallel()
	keeper := new(Keeper)
	// check existing value
	require.Equal(t, 0*time.Second, keeper.buildDelay)
	// set new value
	keeper.SetBuildDelay(10 * time.Second)
	require.Equal(t, 10*time.Second, keeper.buildDelay)
}

func TestKeeper_SetBuildOptimistic(t *testing.T) {
	t.Parallel()
	keeper := new(Keeper)
	// check existing value
	require.False(t, keeper.buildOptimistic)
	// set new value
	keeper.SetBuildOptimistic(true)
	require.True(t, keeper.buildOptimistic)
}

func TestKeeper_parseAndVerifyProposedPayload(t *testing.T) {
	t.Parallel()
	now := time.Now()
	fuzzer := ethclient.NewFuzzer(now.Unix())
	ctx, _, keeper := createKeeper(t, args{
		height: 0,
		header: func(height int64) cmtproto.Header {
			return cmtproto.Header{Height: height}
		},
	})

	tcs := []struct {
		name        string
		setup       func(context.Context) sdk.Context
		msg         func(context.Context) *types.MsgExecutionPayload
		expectedErr string
	}{
		{
			name: "fail: unmarshal payload because of invalid json",
			msg: func(_ context.Context) *types.MsgExecutionPayload {
				return &types.MsgExecutionPayload{ExecutionPayload: []byte("invalid")}
			},
			expectedErr: "unmarshal payload",
		},
		{
			name: "fail: payload number is not equal to head block height + 1",
			msg: func(_ context.Context) *types.MsgExecutionPayload {
				payload, err := ethclient.MakePayload(fuzzer, 100, uint64(now.Unix()), common.Hash{}, common.Address{}, common.Hash{}, &common.Hash{})
				require.NoError(t, err)

				marshaled, err := json.Marshal(payload)
				require.NoError(t, err)

				return &types.MsgExecutionPayload{ExecutionPayload: marshaled}
			},
			expectedErr: "invalid proposed payload number",
		},
		{
			name: "fail: payload parent hash is not equal to head hash",
			msg: func(c context.Context) *types.MsgExecutionPayload {
				execHead, err := keeper.getExecutionHead(c)
				require.NoError(t, err)

				payload, err := ethclient.MakePayload(fuzzer, execHead.GetBlockHeight()+1, uint64(now.Unix()), common.Hash{}, common.Address{}, common.Hash{}, &common.Hash{})
				require.NoError(t, err)

				marshaled, err := json.Marshal(payload)
				require.NoError(t, err)

				return &types.MsgExecutionPayload{ExecutionPayload: marshaled}
			},
			expectedErr: "invalid proposed payload parent hash",
		},
		{
			name: "fail: invalid payload timestamp",
			msg: func(c context.Context) *types.MsgExecutionPayload {
				execHead, err := keeper.getExecutionHead(c)
				require.NoError(t, err)
				weekAgo := execHead.GetBlockTime() - 604800

				payload, err := ethclient.MakePayload(fuzzer, 1, weekAgo, execHead.Hash(), common.Address{}, common.Hash{}, &common.Hash{})
				require.NoError(t, err)

				marshaled, err := json.Marshal(payload)
				require.NoError(t, err)

				return &types.MsgExecutionPayload{ExecutionPayload: marshaled}
			},
			expectedErr: "invalid payload timestamp",
		},
		{
			name: "fail: invalid payload random",
			msg: func(c context.Context) *types.MsgExecutionPayload {
				execHead, err := keeper.getExecutionHead(c)
				require.NoError(t, err)

				payload, err := ethclient.MakePayload(fuzzer, execHead.GetBlockHeight()+1, uint64(now.Unix()), execHead.Hash(), common.Address{}, common.Hash{}, &common.Hash{})
				require.NoError(t, err)

				marshaled, err := json.Marshal(payload)
				require.NoError(t, err)

				return &types.MsgExecutionPayload{ExecutionPayload: marshaled}
			},
			expectedErr: "invalid payload random",
		},
		{
			name: "pass: valid payload",
			msg: func(c context.Context) *types.MsgExecutionPayload {
				execHead, err := keeper.getExecutionHead(c)
				require.NoError(t, err)

				payload, err := ethclient.MakePayload(fuzzer, execHead.GetBlockHeight()+1, uint64(now.Unix()), execHead.Hash(), common.Address{}, execHead.Hash(), &common.Hash{})
				require.NoError(t, err)

				marshaled, err := json.Marshal(payload)
				require.NoError(t, err)

				return &types.MsgExecutionPayload{ExecutionPayload: marshaled}
			},
		},
		{
			name: "pass: valid payload when consensus block time is less than execution block time",
			setup: func(c context.Context) sdk.Context {
				execHead, err := keeper.getExecutionHead(c)
				require.NoError(t, err)
				// update execution head with current block time
				err = keeper.updateExecutionHead(c, engine.ExecutableData{
					Number:    execHead.GetBlockHeight(),
					BlockHash: common.BytesToHash(execHead.GetBlockHash()),
					Timestamp: uint64(now.Unix()),
				})
				require.NoError(t, err)

				// set block time to be less than execution block time
				sdkCtx := sdk.UnwrapSDKContext(c)
				sdkCtx = sdkCtx.WithBlockTime(now.Add(-24 * time.Hour))

				return sdkCtx
			},
			msg: func(c context.Context) *types.MsgExecutionPayload {
				execHead, err := keeper.getExecutionHead(c)
				require.NoError(t, err)

				payload, err := ethclient.MakePayload(fuzzer, execHead.GetBlockHeight()+1, execHead.GetBlockTime()+1, execHead.Hash(), common.Address{}, execHead.Hash(), &common.Hash{})
				require.NoError(t, err)

				marshaled, err := json.Marshal(payload)
				require.NoError(t, err)

				return &types.MsgExecutionPayload{ExecutionPayload: marshaled}
			},
		},
	}

	for _, tc := range tcs {
		//nolint:tparallel // cannot run parallel because of data race on execution head table
		t.Run(tc.name, func(t *testing.T) {
			cachedCtx, _ := ctx.CacheContext()
			if tc.setup != nil {
				cachedCtx = tc.setup(cachedCtx)
			}
			_, err := keeper.parseAndVerifyProposedPayload(cachedCtx, tc.msg(cachedCtx))
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_setOptimisticPayload(t *testing.T) {
	t.Parallel()
	_, _, keeper := createKeeper(t, args{
		height: 0,
		header: func(height int64) cmtproto.Header {
			return cmtproto.Header{Height: height}
		},
	})

	// check existing values
	require.Equal(t, engine.PayloadID{}, keeper.mutablePayload.ID)
	require.Zero(t, keeper.mutablePayload.Height)

	// set new values
	keeper.setOptimisticPayload(engine.PayloadID{1}, 1)
	require.Equal(t, uint64(1), keeper.mutablePayload.Height)
	require.Equal(t, engine.PayloadID{1}, keeper.mutablePayload.ID)
}

func TestKeeper_isNextProposer(t *testing.T) {
	t.Parallel()
	height := int64(1)
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "not proposer",
			args: args{
				height:         height,
				isNextProposer: false,
				header: func(height int64) cmtproto.Header {
					return cmtproto.Header{Height: height}
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "next proposer",
			args: args{
				height:         height,
				isNextProposer: true,
				header: func(height int64) cmtproto.Header {
					return cmtproto.Header{Height: height}
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "validatorsFunc error",
			args: args{
				height: height,
				validatorsFunc: func(ctx context.Context, i int64) (*cmttypes.ValidatorSet, error) {
					return nil, errors.New("error")
				},
				header: func(height int64) cmtproto.Header {
					return cmtproto.Header{Height: height}
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cmtAPI, keeper := createKeeper(t, tt.args)

			got, err := keeper.isNextProposer(ctx, ctx.BlockHeader().Height)
			if (err != nil) != tt.wantErr {
				t.Errorf("isNextProposer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isNextProposer() got = %v, want %v", got, tt.want)
			}
			// make sure that height passed into Validators is correct
			require.Equal(t, tt.args.height, cmtAPI.height)
		})
	}
}

var _ comet.API = (*mockCometAPI)(nil)

type mockCometAPI struct {
	comet.API
	fuzzer         *fuzz.Fuzzer
	validatorSet   *cmttypes.ValidatorSet
	validatorsFunc func(context.Context, int64) (*cmttypes.ValidatorSet, error)
	height         int64
}

func newMockCometAPI(t *testing.T, valFun func(context.Context, int64) (*cmttypes.ValidatorSet, error)) *mockCometAPI {
	t.Helper()
	fuzzer := newFuzzer(0)
	valSet := fuzzValidators(t, fuzzer)

	return &mockCometAPI{
		fuzzer:         fuzzer,
		validatorSet:   valSet,
		validatorsFunc: valFun,
	}
}

func fuzzValidators(t *testing.T, fuzzer *fuzz.Fuzzer) *cmttypes.ValidatorSet {
	t.Helper()
	var validators []*cmttypes.Validator

	fuzzer.NilChance(0).NumElements(3, 7).Fuzz(&validators)

	valSet := new(cmttypes.ValidatorSet)
	err := valSet.UpdateWithChangeSet(validators)
	require.NoError(t, err)

	return valSet
}

func (m *mockCometAPI) Validators(ctx context.Context, height int64) (*cmttypes.ValidatorSet, error) {
	m.height = height
	if m.validatorsFunc != nil {
		return m.validatorsFunc(ctx, height)
	}

	return m.validatorSet, nil
}

// newFuzzer - create a new custom cmttypes.Validator fuzzer.
func newFuzzer(seed int64) *fuzz.Fuzzer {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	f := fuzz.NewWithSeed(seed).NilChance(0)
	f.Funcs(
		func(v *cmttypes.Validator, c fuzz.Continue) {
			privKey := k1.GenPrivKey()
			v.PubKey = privKey.PubKey()
			v.VotingPower = 200
			val := cmttypes.NewValidator(v.PubKey, v.VotingPower)

			*v = *val
		},
	)

	return f
}
