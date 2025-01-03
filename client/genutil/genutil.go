package genutil

import (
	"encoding/json"
	"os"
	"time"

	"cosmossdk.io/math"
	"cosmossdk.io/x/tx/signing"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosstd "github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	atypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	btypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	gtypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	sltypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	sttypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"

	evmenginetypes "github.com/piplabs/story/client/x/evmengine/types"
	evmstakingtypes "github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/buildinfo"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/netconf"
)

// slashingWindows overrides the default slashing signed_blocks_window from 100 to 1000
// since Story block period (+-1s) is very fast, roughly 10x normal period of 10s.
const slashingBlocksWindow = 1000

func MakeGenesis(
	network netconf.ID,
	genesisTime time.Time,
	executionBlockHash common.Hash,
	valPubkeys ...crypto.PubKey,
) (*gtypes.AppGenesis, error) {
	cdc := getCodec()
	txConfig := authtx.NewTxConfig(cdc, nil)

	// Step 1: Create the default genesis app state for all modules.
	appState1 := defaultAppState(network.Static().MaxValidators, executionBlockHash, cdc.MustMarshalJSON)
	appState1Bz, err := json.MarshalIndent(appState1, "", " ")
	if err != nil {
		return nil, errors.Wrap(err, "marshal app state")
	}

	// Step 2: Create the app genesis object and store it to disk.
	appGen := &gtypes.AppGenesis{
		AppName:       "story",
		AppVersion:    buildinfo.Version(),
		GenesisTime:   genesisTime.UTC(),
		ChainID:       network.Static().StoryConsensusChainIDStr(),
		InitialHeight: 1,
		Consensus:     defaultConsensusGenesis(),
		AppState:      appState1Bz,
	}

	// Use this temp file as "disk cache", since the genutil functions require a file path
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		return nil, errors.Wrap(err, "create temp file")
	}
	if err := genutil.ExportGenesisFile(appGen, tempFile.Name()); err != nil {
		return nil, errors.Wrap(err, "export genesis file")
	}

	// Step 3: Create the genesis validators; genesis account and a MsgCreateValidator.
	valTxs := make([]sdk.Tx, 0, len(valPubkeys))
	for _, pubkey := range valPubkeys {
		tx, err := addValidator(txConfig, pubkey, cdc, tempFile.Name())
		if err != nil {
			return nil, errors.Wrap(err, "add validator")
		}
		valTxs = append(valTxs, tx)
	}

	// Step 4: Collect the MsgCreateValidator txs and update the app state (again).
	appState2, err := collectGenTxs(cdc, txConfig, tempFile.Name(), valTxs)
	if err != nil {
		return nil, errors.Wrap(err, "collect genesis transactions")
	}
	appGen.AppState, err = json.MarshalIndent(appState2, "", " ")
	if err != nil {
		return nil, errors.Wrap(err, "marshal app state")
	}

	// Step 5: Validate
	if err := appGen.ValidateAndComplete(); err != nil {
		return nil, errors.Wrap(err, "validate and complete genesis")
	}

	return appGen, validateGenesis(cdc, appState2)
}

func defaultConsensusGenesis() *gtypes.ConsensusGenesis {
	pb := DefaultConsensusParams().ToProto()
	resp := gtypes.NewConsensusGenesis(pb, nil)
	// NewConsensusGenesis has a bug, it doesn't set VoteExtensionsEnableHeight
	resp.Params.ABCI.VoteExtensionsEnableHeight = pb.Abci.VoteExtensionsEnableHeight

	return resp
}

func validateGenesis(cdc codec.Codec, appState map[string]json.RawMessage) error {
	// Staking module
	ststate := sttypes.GetGenesisStateFromAppState(cdc, appState)
	if err := staking.ValidateGenesis(ststate); err != nil {
		return errors.Wrap(err, "validate staking genesis")
	}

	// Slashing module
	var slstate sltypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[sltypes.ModuleName], &slstate); err != nil {
		return errors.Wrap(err, "unmarshal slashing genesis")
	}
	if err := sltypes.ValidateGenesis(slstate); err != nil {
		return errors.Wrap(err, "validate slashing genesis")
	}

	// Bank module
	bstate := btypes.GetGenesisStateFromAppState(cdc, appState)
	if err := bstate.Validate(); err != nil {
		return errors.Wrap(err, "validate bank genesis")
	}

	// Distribution module
	dstate := new(dtypes.GenesisState)
	if err := cdc.UnmarshalJSON(appState[dtypes.ModuleName], dstate); err != nil {
		return errors.Wrap(err, "unmarshal distribution genesis")
	}
	if err := dtypes.ValidateGenesis(dstate); err != nil {
		return errors.Wrap(err, "validate distribution genesis")
	}

	// Auth module
	astate := atypes.GetGenesisStateFromAppState(cdc, appState)
	if err := atypes.ValidateGenesis(astate); err != nil {
		return errors.Wrap(err, "validate auth genesis")
	}

	return nil
}

func collectGenTxs(cdc codec.Codec, txConfig client.TxConfig, genFile string, genTXs []sdk.Tx,
) (map[string]json.RawMessage, error) {
	appState, _, err := gtypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal genesis state")
	}

	appState, err = genutil.SetGenTxsInAppGenesisState(cdc, txConfig.TxJSONEncoder(), appState, genTXs)
	if err != nil {
		return nil, errors.Wrap(err, "set genesis transactions")
	}

	return appState, nil
}

func addValidator(txConfig client.TxConfig, pubkey crypto.PubKey, cdc codec.Codec, genFile string) (sdk.Tx, error) {
	// We use the validator pubkey as the account address
	addr, err := k1util.PubKeyToAddress(pubkey)
	if err != nil {
		return nil, err
	}

	// Add validator with 1 power (1e18 $STAKE ~= 1 ether $STAKE)
	amount := sdk.NewCoin(sdk.DefaultBondDenom, sdk.DefaultPowerReduction)

	err = genutil.AddGenesisAccount(cdc, addr.Bytes(), false, genFile, amount.String(), "", 0, 0, "")
	if err != nil {
		return nil, errors.Wrap(err, "add genesis account")
	}

	pub, err := k1util.PubKeyToCosmos(pubkey)
	if err != nil {
		return nil, err
	}

	var zero = math.LegacyZeroDec()

	msg, err := sttypes.NewMsgCreateValidator(
		sdk.ValAddress(addr.Bytes()).String(),
		pub,
		amount,
		sttypes.Description{Moniker: addr.Hex()},
		sttypes.NewCommissionRates(zero, zero, zero),
		sdk.DefaultPowerReduction,
		sttypes.DefaultLockedTokenType,
	)
	if err != nil {
		return nil, errors.Wrap(err, "create validator message")
	}

	builder := txConfig.NewTxBuilder()

	if err := builder.SetMsgs(msg); err != nil {
		return nil, errors.Wrap(err, "set message")
	}

	return builder.GetTx(), nil
}

// defaultAppState returns the default genesis application state.
func defaultAppState(
	maxVals uint32,
	executionBlockHash common.Hash,
	marshal func(proto.Message) []byte,
) map[string]json.RawMessage {
	stakingGenesis := sttypes.DefaultGenesisState()
	stakingGenesis.Params.MaxValidators = maxVals

	slashingGenesis := sltypes.DefaultGenesisState()
	slashingGenesis.Params.SignedBlocksWindow = slashingBlocksWindow

	evmengGenesis := evmenginetypes.NewGenesisState(evmenginetypes.Params{ExecutionBlockHash: executionBlockHash.Bytes()})

	return map[string]json.RawMessage{
		sttypes.ModuleName:         marshal(stakingGenesis),
		sltypes.ModuleName:         marshal(slashingGenesis),
		atypes.ModuleName:          marshal(atypes.DefaultGenesisState()),
		btypes.ModuleName:          marshal(btypes.DefaultGenesisState()),
		dtypes.ModuleName:          marshal(dtypes.DefaultGenesisState()),
		evmenginetypes.ModuleName:  marshal(evmengGenesis),
		evmstakingtypes.ModuleName: marshal(evmstakingtypes.DefaultGenesisState()),
	}
}

func getCodec() *codec.ProtoCodec {
	sdkConfig := sdk.GetConfig()
	reg, err := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec:          authcodec.NewBech32Codec(sdkConfig.GetBech32AccountAddrPrefix()),
			ValidatorAddressCodec: authcodec.NewBech32Codec(sdkConfig.GetBech32ValidatorAddrPrefix()),
		},
	})
	if err != nil {
		panic(err)
	}

	cosmosstd.RegisterInterfaces(reg)
	atypes.RegisterInterfaces(reg)
	sttypes.RegisterInterfaces(reg)
	sltypes.RegisterInterfaces(reg)
	btypes.RegisterInterfaces(reg)
	dtypes.RegisterInterfaces(reg)
	evmstakingtypes.RegisterInterfaces(reg)
	evmenginetypes.RegisterInterfaces(reg)

	return codec.NewProtoCodec(reg)
}
