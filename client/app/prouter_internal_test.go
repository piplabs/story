package app

import (
	"context"
	"reflect"
	"testing"
	"unsafe"

	"cosmossdk.io/depinject"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmengine/module"
	etypes "github.com/piplabs/story/client/x/evmengine/types"
	esmodule "github.com/piplabs/story/client/x/evmstaking/module"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"

	protov2 "google.golang.org/protobuf/proto"
)

func createRequest(t *testing.T, txConfig client.TxConfig, msg []types.Msg, isFirst, multipleTx bool) *abci.RequestProcessProposal {
	t.Helper()

	b := txConfig.NewTxBuilder()
	require.NoError(t, b.SetMsgs(msg...))

	tx := b.GetTx()

	txBz, err := txConfig.TxEncoder()(tx)
	require.NoError(t, err)

	txs := [][]byte{txBz}
	if msg == nil {
		txs = nil
	}

	if multipleTx {
		txs = append(txs, txBz)
	}

	height := int64(99)
	if isFirst {
		height = 1
	}

	return &abci.RequestProcessProposal{
		Height: height,
		Txs:    txs,
		ProposedLastCommit: abci.CommitInfo{
			Votes: []abci.VoteInfo{
				{BlockIdFlag: cmttypes.BlockIDFlagCommit, Validator: abci.Validator{Power: 1}},
			},
		},
	}
}

func TestProcessProposalRouter(t *testing.T) {
	executionPayloadMsg := &etypes.MsgExecutionPayload{
		Authority: authtypes.NewModuleAddress(etypes.ModuleName).String(),
	}

	tcs := []struct {
		name            string
		first           bool
		payloadMsgs     []types.Msg
		accept          bool
		multipleTx      bool
		expectedSrvCall int
	}{
		{
			name:            "first empty",
			first:           true,
			accept:          true,
			expectedSrvCall: 0,
		},
		{
			name:            "first not empty",
			payloadMsgs:     []types.Msg{executionPayloadMsg},
			first:           true,
			accept:          false,
			expectedSrvCall: 0,
		},
		{
			name:            "too many txs",
			payloadMsgs:     []types.Msg{executionPayloadMsg},
			multipleTx:      true,
			accept:          false,
			expectedSrvCall: 0,
		},
		{
			name:            "one payload message",
			payloadMsgs:     []types.Msg{executionPayloadMsg},
			accept:          true,
			expectedSrvCall: 1,
		},
		{
			name:            "two payload messages",
			payloadMsgs:     []types.Msg{executionPayloadMsg, executionPayloadMsg},
			accept:          false,
			expectedSrvCall: 1,
		},
		{
			name:            "unexpected msg",
			payloadMsgs:     []types.Msg{&stypes.Delegation{}},
			accept:          false,
			expectedSrvCall: 0,
		},
		{
			name: "invalid tx - authority",
			payloadMsgs: []types.Msg{
				&etypes.MsgExecutionPayload{
					Authority: authtypes.NewModuleAddress("test").String(),
				},
			},
			accept:          false,
			expectedSrvCall: 0,
		},
		{
			name: "invalid payload",
			payloadMsgs: []types.Msg{
				&etypes.MsgExecutionPayload{
					Authority:        authtypes.NewModuleAddress(etypes.ModuleName).String(),
					ExecutionPayload: []byte("invalid payload"),
				},
			},
			accept:          false,
			expectedSrvCall: 1,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			key := storetypes.NewKVStoreKey("test")
			ctx := sdktestutil.DefaultContext(key, storetypes.NewTransientStoreKey("test_key"))

			srv := &mockServer{}
			encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{}, esmodule.AppModuleBasic{})

			engineCl := struct {
				ethclient.EngineClient
			}{}
			depCfg := depinject.Configs(DepConfig(), depinject.Supply(newSDKLogger(ctx), engineCl))
			require.NoError(t, depinject.Inject(depCfg, []any{&encCfg.InterfaceRegistry, &encCfg.Codec, &encCfg.TxConfig}...))

			txConfig := encCfg.TxConfig

			router := baseapp.NewMsgServiceRouter()
			router.SetInterfaceRegistry(encCfg.InterfaceRegistry)
			etypes.RegisterMsgServiceServer(router, srv)

			handler := makeProcessProposalHandler(router, txConfig)

			newReq := createRequest(t, txConfig, tc.payloadMsgs, tc.first, tc.multipleTx)

			res, err := handler(ctx, newReq)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSrvCall, srv.payload)
			if tc.accept {
				require.Equal(t, abci.ResponseProcessProposal_ACCEPT, res.Status)
			} else {
				require.Equal(t, abci.ResponseProcessProposal_REJECT, res.Status)
			}
		})
	}
}

var _ types.Tx = &mockTx{}

type mockTx struct{}

func NewMockTx() *mockTx {
	return &mockTx{}
}

func (m *mockTx) GetMsgs() []types.Msg {
	return nil
}

func (m *mockTx) GetMsgsV2() ([]protov2.Message, error) {
	return nil, nil
}

func TestValidateTx(t *testing.T) {
	authority := authtypes.NewModuleAddress(etypes.ModuleName).String()

	tcs := []struct {
		name           string
		msgs           []types.Msg
		isNotSigningTx bool
		callback       func(client.TxBuilder)
		expectedErr    string
	}{
		{
			name:           "invalid proto tx",
			isNotSigningTx: true,
			msgs:           []types.Msg{&etypes.MsgExecutionPayload{Authority: authority}},
			expectedErr:    "invalid proto tx",
		},
		{
			name: "valid payload message",
			msgs: []types.Msg{&etypes.MsgExecutionPayload{Authority: authority}},
		},
		{
			name: "memo not empty",
			msgs: []types.Msg{&etypes.MsgExecutionPayload{Authority: authority}},
			callback: func(b client.TxBuilder) {
				b.SetMemo("memo")
			},
			expectedErr: "disallowed memo in tx",
		},
		{
			name: "nil fee",
			msgs: []types.Msg{&etypes.MsgExecutionPayload{Authority: authority}},
			callback: func(b client.TxBuilder) {
				wrappedTx := b.GetTx()

				wrappedTxField := reflect.ValueOf(wrappedTx).Elem()
				txField := wrappedTxField.FieldByName("tx").Elem()
				authInfoField := txField.FieldByName("AuthInfo").Elem()
				feeField := authInfoField.FieldByName("Fee")
				fieldPtr := unsafe.Pointer(feeField.UnsafeAddr())
				fieldVal := reflect.NewAt(feeField.Type(), fieldPtr).Elem()
				fieldVal.Set(reflect.Zero(feeField.Type()))
			},
			expectedErr: "invalid fee in tx",
		},
		{
			name: "fee not empty",
			msgs: []types.Msg{&etypes.MsgExecutionPayload{Authority: authority}},
			callback: func(b client.TxBuilder) {
				b.SetFeeAmount(types.Coins{types.NewCoin(types.DefaultBondDenom, math.NewInt(1))})
			},
			expectedErr: "invalid fee in tx",
		},
		{
			name: "signatures v2 not empty",
			msgs: []types.Msg{&etypes.MsgExecutionPayload{Authority: authority}},
			callback: func(b client.TxBuilder) {
				sig := signing.SignatureV2{
					PubKey: &secp256k1.PubKey{
						Key: []byte("key"),
					},
					Data: &signing.SingleSignatureData{
						SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
						Signature: []byte("sig"),
					},
					Sequence: 0,
				}

				_ = b.SetSignatures(sig)
			},
			expectedErr: "disallowed signatures in tx",
		},
		{
			name: "fee granter not empty",
			msgs: []types.Msg{&etypes.MsgExecutionPayload{Authority: authority}},
			callback: func(b client.TxBuilder) {
				b.SetFeeGranter(authtypes.NewModuleAddress("granter"))
			},
			expectedErr: "disallowed fee granter in tx",
		},
		{
			name: "tip not empty",
			msgs: []types.Msg{&etypes.MsgExecutionPayload{Authority: authority}},
			callback: func(b client.TxBuilder) {
				var tip = &txtypes.Tip{
					Amount: types.NewCoins(),
					Tipper: "invalid tip",
				}
				wrappedTx := b.GetTx()

				wrappedTxField := reflect.ValueOf(wrappedTx).Elem()
				txField := wrappedTxField.FieldByName("tx").Elem()
				authInfoField := txField.FieldByName("AuthInfo").Elem()
				tipField := authInfoField.FieldByName("Tip")
				fieldPtr := unsafe.Pointer(tipField.UnsafeAddr())
				fieldVal := reflect.NewAt(tipField.Type(), fieldPtr).Elem()
				fieldVal.Set(reflect.ValueOf(tip))
			},
			expectedErr: "disallowed tip in tx",
		},
		{
			name: "signatures not empty",
			msgs: []types.Msg{&etypes.MsgExecutionPayload{Authority: authority}},
			callback: func(b client.TxBuilder) {
				signatures := [][]byte{
					[]byte("invalid signatures"),
				}
				wrappedTx := b.GetTx()

				wrappedTxField := reflect.ValueOf(wrappedTx).Elem()
				txField := wrappedTxField.FieldByName("tx").Elem()
				signaturesField := txField.FieldByName("Signatures")
				fieldPtr := unsafe.Pointer(signaturesField.UnsafeAddr())
				fieldVal := reflect.NewAt(signaturesField.Type(), fieldPtr).Elem()
				fieldVal.Set(reflect.ValueOf(signatures))
			},
			expectedErr: "disallowed signatures in tx",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			key := storetypes.NewKVStoreKey("test")
			ctx := sdktestutil.DefaultContext(key, storetypes.NewTransientStoreKey("test_key"))

			srv := &mockServer{}
			encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{}, esmodule.AppModuleBasic{})

			engineCl := struct {
				ethclient.EngineClient
			}{}
			depCfg := depinject.Configs(DepConfig(), depinject.Supply(newSDKLogger(ctx), engineCl))
			require.NoError(t, depinject.Inject(depCfg, []any{&encCfg.InterfaceRegistry, &encCfg.Codec, &encCfg.TxConfig}...))

			txConfig := encCfg.TxConfig

			router := baseapp.NewMsgServiceRouter()
			router.SetInterfaceRegistry(encCfg.InterfaceRegistry)
			etypes.RegisterMsgServiceServer(router, srv)

			b := txConfig.NewTxBuilder()
			if tc.callback != nil {
				tc.callback(b)
			}

			require.NoError(t, b.SetMsgs(tc.msgs...))

			var tx types.Tx
			if tc.isNotSigningTx {
				tx = NewMockTx()
			} else {
				tx = b.GetTx()
			}
			err := validateTx(tx)

			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

var _ etypes.MsgServiceServer = &mockServer{}

type mockServer struct {
	etypes.MsgServiceServer
	payload int
}

func (s *mockServer) ExecutionPayload(_ context.Context, payload *etypes.MsgExecutionPayload) (*etypes.ExecutionPayloadResponse, error) {
	s.payload++

	if payload.ExecutionPayload == nil {
		return &etypes.ExecutionPayloadResponse{}, nil
	}

	if err := etypes.ValidateExecPayload(payload); err != nil {
		return nil, errors.New("invalid execution payload")
	}

	return &etypes.ExecutionPayloadResponse{}, nil
}
