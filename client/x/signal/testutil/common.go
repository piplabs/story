package testutil

import (
	"bytes"
	cosmosmath "cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	tmed "github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ccrypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	// TestingStakeParams is a set of staking params for testing
	TestingStakeParams = stakingtypes.Params{
		UnbondingTime:     100,
		MaxValidators:     10,
		MaxEntries:        10,
		HistoricalEntries: 10000,
		BondDenom:         "stake",
		MinCommissionRate: cosmosmath.LegacyNewDec(0),
	}

	// HardcodedConsensusPrivKeys
	FixedConsensusPrivKeys = []tmed.PrivKey{
		tmed.GenPrivKeyFromSecret([]byte("12345678901234567890123456389012")),
		tmed.GenPrivKeyFromSecret([]byte("12345678901234567890123456389013")),
		tmed.GenPrivKeyFromSecret([]byte("12345678901234567890123456389014")),
		tmed.GenPrivKeyFromSecret([]byte("12345678901234567890123456389015")),
		tmed.GenPrivKeyFromSecret([]byte("12345678901234567890123456389016")),
	}

	FixedNetworkPrivKeys = []tmed.PrivKey{
		tmed.GenPrivKeyFromSecret([]byte("12345678901234567890123456786012")),
		tmed.GenPrivKeyFromSecret([]byte("12345678901234567890123456786013")),
		tmed.GenPrivKeyFromSecret([]byte("12345678901234567890123456786014")),
		tmed.GenPrivKeyFromSecret([]byte("12345678901234567890123456786015")),
		tmed.GenPrivKeyFromSecret([]byte("12345678901234567890123456786016")),
	}

	// FixedMnemonics is a set of fixed mnemonics for testing.
	// Account names are: validator1, validator2, validator3, validator4, validator5
	FixedMnemonics = []string{
		"body world north giggle crop reduce height copper damp next verify orphan lens loan adjust inform utility theory now ranch motion opinion crowd fun",
		"body champion street fat bone above office guess waste vivid gift around approve elevator depth fiber alarm usual skirt like organ space antique silk",
		"cheap alpha render punch clap prize duty drive steel situate person radar smooth elegant over chronic wait danger thumb soft letter spatial acquire rough",
		"outdoor ramp suspect office disagree world attend vanish small wish capable fall wall soon damp session emotion chest toss viable meat host clerk truth",
		"ability evidence casino cram weasel chest brush bridge sister blur onion found glad own mansion amateur expect force fun dragon famous alien appear open",
	}

	// ConsPrivKeys generate ed25519 ConsPrivKeys to be used for validator operator keys
	ConsPrivKeys = []ccrypto.PrivKey{
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
	}

	// ConsPubKeys holds the consensus public keys to be used for validator operator keys
	ConsPubKeys = []ccrypto.PubKey{
		ConsPrivKeys[0].PubKey(),
		ConsPrivKeys[1].PubKey(),
		ConsPrivKeys[2].PubKey(),
		ConsPrivKeys[3].PubKey(),
		ConsPrivKeys[4].PubKey(),
	}

	// AccPrivKeys generate secp256k1 pubkeys to be used for account pub keys
	AccPrivKeys = []ccrypto.PrivKey{
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
	}

	// AccPubKeys holds the pub keys for the account keys
	AccPubKeys = []ccrypto.PubKey{
		AccPrivKeys[0].PubKey(),
		AccPrivKeys[1].PubKey(),
		AccPrivKeys[2].PubKey(),
		AccPrivKeys[3].PubKey(),
		AccPrivKeys[4].PubKey(),
	}

	// AccAddrs holds the sdk.AccAddresses
	AccAddrs = []sdk.AccAddress{
		sdk.AccAddress(AccPubKeys[0].Address()),
		sdk.AccAddress(AccPubKeys[1].Address()),
		sdk.AccAddress(AccPubKeys[2].Address()),
		sdk.AccAddress(AccPubKeys[3].Address()),
		sdk.AccAddress(AccPubKeys[4].Address()),
	}

	// ValAddrs holds the sdk.ValAddresses
	ValAddrs = []sdk.ValAddress{
		sdk.ValAddress(AccPubKeys[0].Address()),
		sdk.ValAddress(AccPubKeys[1].Address()),
		sdk.ValAddress(AccPubKeys[2].Address()),
		sdk.ValAddress(AccPubKeys[3].Address()),
		sdk.ValAddress(AccPubKeys[4].Address()),
	}

	// EVMAddrs holds etheruem addresses
	EVMAddrs = initEVMAddrs(100)

	// InitTokens holds the number of tokens to initialize an account with
	InitTokens = sdk.TokensFromConsensusPower(110, sdk.DefaultPowerReduction)

	// InitCoins holds the number of coins to initialize an account with
	InitCoins = sdk.NewCoins(sdk.NewCoin(TestingStakeParams.BondDenom, InitTokens))

	// StakingAmount holds the staking power to start a validator with
	StakingAmount = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
)

func initEVMAddrs(count int) []gethcommon.Address {
	addresses := make([]gethcommon.Address, count)
	for i := 0; i < count; i++ {
		evmAddr := gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(i + 1)}, gethcommon.AddressLength))
		addresses[i] = evmAddr
	}
	return addresses
}

// TestInput stores the various keepers required to test Blobstream
type TestInput struct {
	AccountKeeper authkeeper.AccountKeeper
	StakingKeeper stakingkeeper.Keeper
	Context       sdk.Context
	Marshaler     codec.Codec
	LegacyAmino   *codec.LegacyAmino
}

func CreateValidator(
	t *testing.T,
	input TestInput,
	accAddr sdk.AccAddress,
	accPubKey ccrypto.PubKey,
	accountNumber uint64,
	valAddr sdk.ValAddress,
	consPubKey ccrypto.PubKey,
	stakingAmount cosmosmath.Int,
) {
	// Initialize the account for the key
	acc := input.AccountKeeper.NewAccount(
		input.Context,
		authtypes.NewBaseAccount(accAddr, accPubKey, accountNumber, 0),
	)

	// Set the account in state
	input.AccountKeeper.SetAccount(input.Context, acc)

	// Create a validator for that account using some tokens in the account
	// and the staking handler
	msgServer := stakingkeeper.NewMsgServerImpl(&input.StakingKeeper)
	_, err := msgServer.CreateValidator(input.Context, NewTestMsgCreateValidator(valAddr, consPubKey, stakingAmount))
	require.NoError(t, err)
}

func NewTestMsgCreateValidator(
	address sdk.ValAddress,
	pubKey ccrypto.PubKey,
	amt cosmosmath.Int,
) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec())
	out, err := stakingtypes.NewMsgCreateValidator(
		address.String(), pubKey, sdk.NewCoin("stake", amt),
		stakingtypes.Description{
			Moniker:         "",
			Identity:        "",
			Website:         "",
			SecurityContact: "",
			Details:         "",
		}, commission, sdkmath.OneInt(), 0,
	)
	if err != nil {
		panic(err)
	}
	return out
}
