package genutil

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	etypes "github.com/piplabs/story/client/x/evmengine/types"
)

func TestDefaultConsensusParams(t *testing.T) {
	t.Parallel()
	cons := defaultConsensusGenesis()
	require.EqualValues(t, 0, cons.Params.ABCI.VoteExtensionsEnableHeight)
	require.EqualValues(t, types.ABCIPubKeyTypeSecp256k1, cons.Params.Validator.PubKeyTypes[0])
	require.EqualValues(t, -1, cons.Params.Block.MaxBytes)
	require.EqualValues(t, -1, cons.Params.Block.MaxGas)
}

func TestEncodeTXs(t *testing.T) {
	t.Parallel()
	msgs := []sdk.Msg{
		&etypes.MsgExecutionPayload{
			Authority: authtypes.NewModuleAddress("evm").String(),
		},
	}

	cdc := getCodec()
	txConfig := authtx.NewTxConfig(cdc, nil)

	b := txConfig.NewTxBuilder()
	err := b.SetMsgs(msgs...)
	require.NoError(t, err)

	tx := b.GetTx()

	require.Len(t, tx.GetMsgs(), 1)
	msgsV2, err := tx.GetMsgsV2()
	require.NoError(t, err)
	require.Len(t, msgsV2, 1)
}
