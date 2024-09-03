package keeper_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"

	"github.com/cometbft/cometbft/crypto"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"

	evmenginetypes "github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/client/x/evmstaking/keeper"
	"github.com/piplabs/story/client/x/evmstaking/module"
	estestutil "github.com/piplabs/story/client/x/evmstaking/testutil"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

var (
	PKs = simtestutil.CreateTestPubKeys(3)
)

type TestSuite struct {
	suite.Suite

	Ctx sdk.Context

	AccountKeeper    *estestutil.MockAccountKeeper
	BankKeeper       *estestutil.MockBankKeeper
	DistrKeeper      *estestutil.MockDistributionKeeper
	StakingKeeper    *skeeper.Keeper
	EVMStakingKeeper *keeper.Keeper
	msgServer        types.MsgServiceServer

	encCfg moduletestutil.TestEncodingConfig
}

func (s *TestSuite) SetupTest() {
	s.encCfg = moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})
	evmstakingKey := storetypes.NewKVStoreKey(types.StoreKey)
	stakingKey := storetypes.NewKVStoreKey(stypes.StoreKey)
	storeService := runtime.NewKVStoreService(evmstakingKey)
	stakingStoreService := runtime.NewKVStoreService(stakingKey)

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	cms.MountStoreWithDB(evmstakingKey, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(stakingKey, storetypes.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	s.Require().NoError(err)

	s.Ctx = sdk.NewContext(cms, cmtproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	legacyAmino := codec.NewLegacyAmino()
	stypes.RegisterLegacyAminoCodec(legacyAmino)
	stypes.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("story", "storypub")
	cfg.SetBech32PrefixForValidator("storyvaloper", "storyvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("storyvalcons", "storyvalconspub")

	// gomock initializations
	ctrl := gomock.NewController(s.T())

	// mock keepers
	accountKeeper := estestutil.NewMockAccountKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName)).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(stypes.ModuleName).Return(authtypes.NewModuleAddress(stypes.ModuleName)).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(stypes.BondedPoolName).Return(authtypes.NewModuleAddress(stypes.BondedPoolName)).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(stypes.NotBondedPoolName).Return(authtypes.NewModuleAddress(stypes.NotBondedPoolName)).AnyTimes()
	accountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("story")).AnyTimes()
	s.AccountKeeper = accountKeeper
	bankKeeper := estestutil.NewMockBankKeeper(ctrl)
	s.BankKeeper = bankKeeper
	distrKeeper := estestutil.NewMockDistributionKeeper(ctrl)
	s.DistrKeeper = distrKeeper
	slashingKeeper := estestutil.NewMockSlashingKeeper(ctrl)

	// staking keeper
	stakingKeeper := skeeper.NewKeeper(
		marshaler,
		stakingStoreService,
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(stypes.ModuleName).String(),
		address.NewBech32Codec("storyvaloper"),
		address.NewBech32Codec("storyvalcons"),
	)
	s.StakingKeeper = stakingKeeper
	s.Require().NoError(s.StakingKeeper.SetParams(s.Ctx, stypes.DefaultParams()))

	// emvstaking keeper
	ethCl, err := ethclient.NewEngineMock()
	s.Require().NoError(err)
	evmstakingKeeper := keeper.NewKeeper(
		marshaler,
		storeService,
		accountKeeper,
		bankKeeper,
		slashingKeeper,
		stakingKeeper,
		distrKeeper,
		authtypes.NewModuleAddress(types.ModuleName).String(),
		ethCl,
		address.NewBech32Codec("storyvaloper"),
	)
	s.EVMStakingKeeper = evmstakingKeeper
	s.msgServer = keeper.NewMsgServerImpl(evmstakingKeeper)
}

func (s *TestSuite) TestLogger() {
	require := s.Require()
	logger := keeper.Logger(s.Ctx)
	require.NotNil(logger)
}

func (s *TestSuite) TestGetAuthority() {
	require := s.Require()
	require.Equal(authtypes.NewModuleAddress(types.ModuleName).String(), s.EVMStakingKeeper.GetAuthority())
}

func (s *TestSuite) TestValidatorAddressCodec() {
	require := s.Require()
	keeper := s.EVMStakingKeeper
	require.NotNil(keeper.ValidatorAddressCodec())
	_, err := keeper.ValidatorAddressCodec().StringToBytes("storyvaloper1hmjw3pvkjtndpg8wqppwdn8udd835qpaa6r6y0")
	require.NoError(err)
}

func (s *TestSuite) TestProcessStakingEvents() {
	require := s.Require()
	ctx, evmstakingKeeper := s.Ctx, s.EVMStakingKeeper
	pubKeys, _, _ := createAddresses(3)
	delPubKey := pubKeys[0]
	delEvmAddr := common.BytesToAddress(pubKeys[0].Address().Bytes())
	var evmAddrBytes [32]byte
	copy(evmAddrBytes[:], delEvmAddr.Bytes())
	delSecp256k1PubKey, err := secp256k1.ParsePubKey(delPubKey.Bytes())
	require.NoError(err)
	uncompressedDelPubKeyBytes := delSecp256k1PubKey.SerializeUncompressed()
	valPubKey := pubKeys[1]
	valPubKey2 := pubKeys[2]
	dummyHash := common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	stakingAbi, err := bindings.IPTokenStakingMetaData.GetAbi()
	require.NoError(err)
	slashingAbi, err := bindings.IPTokenSlashingMetaData.GetAbi()
	require.NoError(err)
	tcs := []struct {
		name          string
		evmEvents     func() ([]*evmenginetypes.EVMEvent, error)
		setup         func(c context.Context)
		expectedError string
	}{
		{
			name: "fail: invalid evm event",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.SetWithdrawalAddress.ID, dummyHash}}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}
				evmEvents[0].Address = nil

				return evmEvents, nil
			},
			expectedError: "verify log [BUG]",
		},
		// INVALID LOGS but PASS Cases because currently we are handling it as a continued
		{
			name: "pass(continue): invalid SetWithdrawalEvent log",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.SetWithdrawalAddress.ID, dummyHash}}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		{
			name: "pass(continue): invalid CreateValidatorEvent log",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.CreateValidatorEvent.ID, dummyHash}}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		{
			name: "pass(continue): invalid DepositEvent log",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.DepositEvent.ID, dummyHash}}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		{
			name: "pass(continue): invalid RedelegateEvent log",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.RedelegateEvent.ID, dummyHash}}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		{
			name: "pass(continue): invalid WithdrawEvent log",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.WithdrawEvent.ID, dummyHash}}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		{
			name: "pass(continue): invalid UnjailEvent log",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.UnjailEvent.ID, dummyHash, dummyHash}}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		// FAIL TO PROCESS but PASS Cases because currently we are handling it as a continued.
		// Only basic failure cases are validated. Various failure and success scenarios that may occur during the actual process
		// are tested separately with unit tests in the files where each processing logic is defined.
		{
			name: "pass(continue): fail to process SetWithdrawalAddressEvent - invalid delegator pubkey",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				invalidDelPubKey := delPubKey.Bytes()[1:]
				data, err := stakingAbi.Events["SetWithdrawalAddress"].Inputs.NonIndexed().Pack(
					invalidDelPubKey,
					evmAddrBytes,
				)
				require.NoError(err)
				logs := []ethtypes.Log{{Topics: []common.Hash{types.SetWithdrawalAddress.ID}, Data: data}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		{
			name: "pass(continue): fail to process CreateValidatorEvent - corrupted pubkey",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				data, err := stakingAbi.Events["CreateValidator"].Inputs.NonIndexed().Pack(
					uncompressedDelPubKeyBytes,
					createCorruptedPubKey(delPubKey.Bytes()),
					"moniker",
					new(big.Int).SetUint64(100),
					uint32(1000),
					uint32(5000),
					uint32(500),
				)
				require.NoError(err)
				logs := []ethtypes.Log{{Topics: []common.Hash{types.CreateValidatorEvent.ID}, Data: data}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		{
			name: "pass(continue): fail to process DepositEvent - corrupted delegator pubkey",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				data, err := stakingAbi.Events["Deposit"].Inputs.NonIndexed().Pack(
					uncompressedDelPubKeyBytes,
					createCorruptedPubKey(delPubKey.Bytes()),
					valPubKey.Bytes(),
					new(big.Int).SetUint64(100),
				)
				require.NoError(err)
				logs := []ethtypes.Log{{Topics: []common.Hash{types.DepositEvent.ID}, Data: data}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		{
			name: "pass(continue): fail to process RedelegateEvent - corrupted delegator pubkey",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				data, err := stakingAbi.Events["Redelegate"].Inputs.NonIndexed().Pack(
					createCorruptedPubKey(delPubKey.Bytes()),
					valPubKey.Bytes(),
					valPubKey2.Bytes(),
					new(big.Int).SetUint64(100),
				)
				require.NoError(err)
				logs := []ethtypes.Log{{Topics: []common.Hash{types.RedelegateEvent.ID}, Data: data}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		{
			name: "pass(continue): fail to process WithdrawEvent - corrupted delegator pubkey",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				data, err := stakingAbi.Events["Withdraw"].Inputs.NonIndexed().Pack(
					createCorruptedPubKey(delPubKey.Bytes()),
					valPubKey.Bytes(),
					new(big.Int).SetUint64(100),
				)
				require.NoError(err)
				logs := []ethtypes.Log{{Topics: []common.Hash{types.WithdrawEvent.ID}, Data: data}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		{
			name: "pass(continue): fail to process UnjailEvent - invalid validator pubkey",
			evmEvents: func() ([]*evmenginetypes.EVMEvent, error) {
				invalidValPubKey := valPubKey.Bytes()[1:]
				data, err := slashingAbi.Events["Unjail"].Inputs.NonIndexed().Pack(invalidValPubKey)
				require.NoError(err)
				logs := []ethtypes.Log{{Topics: []common.Hash{types.UnjailEvent.ID, common.BytesToHash(delEvmAddr.Bytes())}, Data: data}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				if err != nil {
					return nil, err
				}

				return evmEvents, nil
			},
		},
		// SUCCESS Cases should be tested separately with unit tests in the files where each processing logic is defined.
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			if tc.setup != nil {
				tc.setup(cachedCtx)
			}
			evmLogs, err := tc.evmEvents()
			require.NoError(err)
			err = evmstakingKeeper.ProcessStakingEvents(cachedCtx, 1, evmLogs)
			if tc.expectedError != "" {
				require.Error(err)
				require.Contains(err.Error(), tc.expectedError)
			} else {
				require.NoError(err)
			}
		})
	}
}

func TestTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TestSuite))
}

// setupValidatorAndDelegation creates a validator and delegation for testing.
func (s *TestSuite) setupValidatorAndDelegation(ctx context.Context, valPubKey, delPubKey crypto.PubKey, valAddr sdk.ValAddress, delAddr sdk.AccAddress) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	require := s.Require()
	bankKeeper, stakingKeeper, keeper := s.BankKeeper, s.StakingKeeper, s.EVMStakingKeeper

	// Convert public key to cosmos format
	valCosmosPubKey, err := k1util.PubKeyToCosmos(valPubKey)
	require.NoError(err)

	// Create and update validator
	val := testutil.NewValidator(s.T(), valAddr, valCosmosPubKey)
	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	validator, _ := val.AddTokensFromDel(valTokens)
	bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
	_ = skeeper.TestingUpdateValidator(stakingKeeper, sdkCtx, validator, true)

	// Create and set delegation
	delAmt := stakingKeeper.TokensFromConsensusPower(ctx, 100).ToLegacyDec()
	delegation := stypes.NewDelegation(delAddr.String(), valAddr.String(), delAmt)
	require.NoError(stakingKeeper.SetDelegation(ctx, delegation))

	// Map delegator to EVM address
	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(err)
	require.NoError(keeper.DelegatorMap.Set(ctx, delAddr.String(), delEvmAddr.String()))

	// Ensure delegation is set correctly
	delegation, err = stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	require.NoError(err)
	require.Equal(delAmt, delegation.GetShares())
}

func createCorruptedPubKey(pubKey []byte) []byte {
	corruptedPubKey := append([]byte(nil), pubKey...)
	corruptedPubKey[0] = 0x04
	corruptedPubKey[1] = 0xFF

	return corruptedPubKey
}

// ethLogsToEvmEvents converts Ethereum logs to a slice of EVM events.
func ethLogsToEvmEvents(logs []ethtypes.Log) ([]*evmenginetypes.EVMEvent, error) {
	events := make([]*evmenginetypes.EVMEvent, 0, len(logs))
	for _, l := range logs {
		topics := make([][]byte, 0, len(l.Topics))
		for _, t := range l.Topics {
			topics = append(topics, t.Bytes())
		}
		events = append(events, &evmenginetypes.EVMEvent{
			Address: l.Address.Bytes(),
			Topics:  topics,
			Data:    l.Data,
		})
	}

	for _, log := range events {
		if err := log.Verify(); err != nil {
			return nil, errors.Wrap(err, "verify log")
		}
	}

	return events, nil
}

