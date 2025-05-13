package keeper_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cometbft/cometbft/crypto"
	k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/genutil"
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
	dummyHash = common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
)

func TestLogger(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, _ := createKeeperWithMockStaking(t)
	logger := keeper.Logger(ctx)
	require.NotNil(t, logger)
}

func TestGetAuthority(t *testing.T) {
	//nolint:dogsled // This is common helper function
	_, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

	require.Equal(t, authtypes.NewModuleAddress(types.ModuleName).String(), esk.GetAuthority())
}

func TestValidatorAddressCodec(t *testing.T) {
	//nolint:dogsled // This is common helper function
	_, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

	require.NotNil(t, esk.ValidatorAddressCodec())
	_, err := esk.ValidatorAddressCodec().StringToBytes("storyvaloper1hmjw3pvkjtndpg8wqppwdn8udd835qpaa6r6y0")
	require.NoError(t, err)
}

func ConvertAddressToBytes32(addr common.Address) [32]byte {
	uint160 := new(big.Int).SetBytes(addr.Bytes())
	var bytes32 [32]byte
	copy(bytes32[32-len(uint160.Bytes()):], uint160.Bytes())

	return bytes32
}

func TestProcessStakingEvents(t *testing.T) {
	pubKeys, _, valAddrs := createAddresses(3)

	// delegator
	delPubKey := pubKeys[0]
	delEVMAddr, err := keeper.CmpPubKeyToEVMAddress(delPubKey.Bytes())
	require.NoError(t, err)
	delAccAddrFromEVM := sdk.AccAddress(delEVMAddr.Bytes())

	// validator
	val1PubKey := pubKeys[1]
	val1ValAddr := valAddrs[1]
	val1EVMAddr, err := k1util.CosmosPubkeyToEVMAddress(val1PubKey.Bytes())
	require.NoError(t, err)
	val2PubKey := pubKeys[2]
	val2ValAddr := valAddrs[2]
	invalidPubKey := append([]byte{0x04}, val1PubKey.Bytes()[1:]...)

	// ABI of IPTokenStaking contract
	stakingAbi, err := bindings.IPTokenStakingMetaData.GetAbi()
	require.NoError(t, err)

	// gwei multiplication for evm input
	gwei, exp := big.NewInt(10), big.NewInt(9)
	gwei.Exp(gwei, exp, nil)

	// delegation amount
	delCoin := sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))
	delAmtGwei := new(big.Int).Mul(gwei, new(big.Int).SetUint64(delCoin.Amount.Uint64()))

	tcs := []struct {
		name           string
		setupMocks     func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper)
		evmEvents      func() []*evmenginetypes.EVMEvent
		setup          func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context
		expectedError  string
		postStateCheck func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper)
	}{
		{
			name: "fail: invalid evm event",
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.SetWithdrawalAddress.ID, dummyHash}, TxHash: dummyHash}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				evmEvents[0].Address = nil

				return evmEvents
			},
			expectedError: "verify log [BUG]",
		},

		// ********** UpdateValidatorCommission **********
		{
			name: "fail(continue): invalid UpdateValidatorCommission log",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// create validator
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.UpdateValidatorCommission.ID, dummyHash}, TxHash: dummyHash}}
				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// commission should not be changed
				validator, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				require.Equal(t, stypes.CommissionRates{
					Rate:          math.LegacyNewDecWithPrec(int64(1000), 4),
					MaxRate:       math.LegacyNewDecWithPrec(int64(5000), 4),
					MaxChangeRate: math.LegacyNewDecWithPrec(int64(1000), 4),
				}, validator.Commission.CommissionRates)
			},
		},
		{
			name: "fail(continue): process UpdateValidatorCommission event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// create validator
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["UpdateValidatorCommission"].Inputs.NonIndexed().Pack(
					invalidPubKey,
					uint32(1000),
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.UpdateValidatorCommission.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// commission should not be changed
				validator, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				require.Equal(t, stypes.CommissionRates{
					Rate:          math.LegacyNewDecWithPrec(int64(1000), 4),
					MaxRate:       math.LegacyNewDecWithPrec(int64(5000), 4),
					MaxChangeRate: math.LegacyNewDecWithPrec(int64(1000), 4),
				}, validator.Commission.CommissionRates)
			},
		},
		{
			name: "pass: process UpdateValidatorCommission event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// create validator
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["UpdateValidatorCommission"].Inputs.NonIndexed().Pack(
					val1PubKey.Bytes(),
					uint32(1500), // change to 15%
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.UpdateValidatorCommission.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check the new commission rate
				validator, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				require.Equal(t, stypes.CommissionRates{
					Rate:          math.LegacyNewDecWithPrec(int64(1500), 4),
					MaxRate:       math.LegacyNewDecWithPrec(int64(5000), 4),
					MaxChangeRate: math.LegacyNewDecWithPrec(int64(1000), 4),
				}, validator.Commission.CommissionRates)
			},
		},

		// ********** SetWithdrawalAddress **********
		{
			name: "fail(continue): invalid SetWithdrawalAddress event",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorWithdrawAddress.Set(ctx, delAccAddrFromEVM.String(), delEVMAddr.String()))

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.SetWithdrawalAddress.ID, dummyHash}, TxHash: dummyHash}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				withdrawAddr, err := esk.DelegatorWithdrawAddress.Get(ctx, delAccAddrFromEVM.String())
				require.NoError(t, err)
				require.Equal(t, delEVMAddr.String(), withdrawAddr)
			},
		},
		{
			name: "pass: process SetWithdrawalAddress event",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorWithdrawAddress.Set(ctx, delAccAddrFromEVM.String(), delEVMAddr.String()))

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["SetWithdrawalAddress"].Inputs.NonIndexed().Pack(
					delEVMAddr,
					ConvertAddressToBytes32(common.MaxAddress),
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.SetWithdrawalAddress.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				withdrawAddr, err := esk.DelegatorWithdrawAddress.Get(ctx, delAccAddrFromEVM.String())
				require.NoError(t, err)
				require.Equal(t, common.MaxAddress.String(), withdrawAddr)
			},
		},

		// ********** SetRewardAddress **********
		{
			name: "fail(continue): invalid SetRewardAddress event",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, delAccAddrFromEVM.String(), delEVMAddr.String()))

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.SetRewardAddress.ID, dummyHash}, TxHash: dummyHash}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				rewardAddr, err := esk.DelegatorRewardAddress.Get(ctx, delAccAddrFromEVM.String())
				require.NoError(t, err)
				require.Equal(t, delEVMAddr.String(), rewardAddr)
			},
		},
		{
			name: "pass: process SetRewardAddress event",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, delAccAddrFromEVM.String(), delEVMAddr.String()))

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["SetRewardAddress"].Inputs.NonIndexed().Pack(
					delEVMAddr,
					ConvertAddressToBytes32(common.MaxAddress),
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.SetRewardAddress.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				rewardAddr, err := esk.DelegatorRewardAddress.Get(ctx, delAccAddrFromEVM.String())
				require.NoError(t, err)
				require.Equal(t, common.MaxAddress.String(), rewardAddr)
			},
		},

		// ********** SetOperator **********
		{
			name: "fail(continue): invalid SetOperator event",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, delAccAddrFromEVM.String(), delEVMAddr.String()))

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.SetOperator.ID, dummyHash}, TxHash: dummyHash}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				operatorAddr, err := esk.DelegatorOperatorAddress.Get(ctx, delAccAddrFromEVM.String())
				require.NoError(t, err)
				require.Equal(t, delEVMAddr.String(), operatorAddr)
			},
		},
		{
			name: "pass: process SetOperator event",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, delAccAddrFromEVM.String(), delEVMAddr.String()))

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["SetOperator"].Inputs.NonIndexed().Pack(
					delEVMAddr,
					ConvertAddressToBytes32(common.MaxAddress),
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.SetOperator.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				operatorAddr, err := esk.DelegatorOperatorAddress.Get(ctx, delAccAddrFromEVM.String())
				require.NoError(t, err)
				require.Equal(t, common.MaxAddress.String(), operatorAddr)
			},
		},

		// ********** UnsetOperator **********
		{
			name: "fail(continue): invalid UnsetOperator event",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, delAccAddrFromEVM.String(), delEVMAddr.String()))

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.UnsetOperator.ID, dummyHash}, TxHash: dummyHash}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				operatorAddr, err := esk.DelegatorOperatorAddress.Get(ctx, delAccAddrFromEVM.String())
				require.NoError(t, err)
				require.Equal(t, delEVMAddr.String(), operatorAddr)
			},
		},
		{
			name: "pass: process UnsetOperator event",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, delAccAddrFromEVM.String(), delEVMAddr.String()))

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["UnsetOperator"].Inputs.NonIndexed().Pack(
					delEVMAddr,
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.UnsetOperator.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				_, err := esk.DelegatorOperatorAddress.Get(ctx, delAccAddrFromEVM.String())
				require.ErrorContains(t, err, "not found")
			},
		},

		// ********** CreateValidator **********
		{
			name: "fail(continue): invalid CreateValidator event",
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.CreateValidatorEvent.ID, dummyHash}, TxHash: dummyHash}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				_, err := sk.GetValidator(ctx, val1ValAddr)
				require.ErrorContains(t, err, stypes.ErrNoValidatorFound.Error())
			},
		},
		{
			name: "fail(continue): process CreateValidator event",
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["CreateValidator"].Inputs.NonIndexed().Pack(
					invalidPubKey,
					"moniker",
					delAmtGwei,
					uint32(1000),
					uint32(5000),
					uint32(500),
					uint8(0),
					cmpToEVM(delPubKey.Bytes()),
					[]byte{},
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.CreateValidatorEvent.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				_, err := sk.GetValidator(ctx, val1ValAddr)
				require.ErrorContains(t, err, stypes.ErrNoValidatorFound.Error())
			},
		},
		{
			name: "pass: process CreateValidator event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["CreateValidator"].Inputs.NonIndexed().Pack(
					val1PubKey.Bytes(),
					"test",
					delAmtGwei,
					uint32(1000),
					uint32(5000),
					uint32(500),
					uint8(0),
					cmpToEVM(val1PubKey.Bytes()),
					[]byte{},
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.CreateValidatorEvent.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				validator, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				require.NotNil(t, validator)
			},
		},

		// ********** Deposit **********
		{
			name: "fail(continue): invalid Deposit event",
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.DepositEvent.ID, dummyHash}, TxHash: dummyHash}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				_, err := sk.GetDelegation(ctx, delAccAddrFromEVM, val1ValAddr)
				require.ErrorContains(t, err, stypes.ErrNoDelegation.Error())
			},
		},
		{
			name: "fail(continue): process Deposit event",
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["Deposit"].Inputs.NonIndexed().Pack(
					delEVMAddr,
					invalidPubKey,
					delAmtGwei,
					new(big.Int).SetUint64(0),
					new(big.Int).SetUint64(0),
					cmpToEVM(delPubKey.Bytes()),
					[]byte{},
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.DepositEvent.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				_, err := sk.GetDelegation(ctx, delAccAddrFromEVM, val1ValAddr)
				require.ErrorContains(t, err, stypes.ErrNoDelegation.Error())
			},
		},
		{
			name: "pass: process Deposit event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["Deposit"].Inputs.NonIndexed().Pack(
					delEVMAddr,
					val1PubKey.Bytes(),
					delAmtGwei,
					new(big.Int).SetUint64(0),
					new(big.Int).SetUint64(0),
					cmpToEVM(delPubKey.Bytes()),
					[]byte{},
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.DepositEvent.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check delegation
				delegation, err := sk.GetDelegation(ctx, delAccAddrFromEVM, val1ValAddr)
				require.NoError(t, err)
				require.NotNil(t, delegation)

				// check delegated token
				validator, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				require.Equal(t, uint64(10000100), validator.Tokens.Uint64())
			},
		},

		// ********** Redelegate **********
		{
			name: "fail(continue): invalid Redelegate event",
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.RedelegateEvent.ID, dummyHash}, TxHash: dummyHash}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check delegation from source validator
				_, err := sk.GetDelegation(ctx, delAccAddrFromEVM, val1ValAddr)
				require.ErrorContains(t, err, stypes.ErrNoDelegation.Error())

				// check delegatio from destination validator
				_, err = sk.GetDelegation(ctx, delAccAddrFromEVM, val2ValAddr)
				require.ErrorContains(t, err, stypes.ErrNoDelegation.Error())
			},
		},
		{
			name: "fail(continue): process Relegate event",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["Redelegate"].Inputs.NonIndexed().Pack(
					delEVMAddr,
					invalidPubKey,
					val2PubKey.Bytes(),
					new(big.Int).SetUint64(0),
					cmpToEVM(delPubKey.Bytes()),
					delAmtGwei,
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.RedelegateEvent.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check delegation from source validator
				_, err := sk.GetDelegation(ctx, delAccAddrFromEVM, val1ValAddr)
				require.ErrorContains(t, err, stypes.ErrNoDelegation.Error())

				// check delegatio from destination validator
				_, err = sk.GetDelegation(ctx, delAccAddrFromEVM, val2ValAddr)
				require.ErrorContains(t, err, stypes.ErrNoDelegation.Error())
			},
		},
		{
			name: "pass: process Redelegate event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				// create destination validator
				createValidator(t, ctx, sk, val2PubKey, val2ValAddr, 0)

				// create delegation
				createDelegation(t, ctx, sk, delAccAddrFromEVM, val1ValAddr, delCoin.Amount.Int64(), true)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["Redelegate"].Inputs.NonIndexed().Pack(
					delEVMAddr,
					val1PubKey.Bytes(),
					val2PubKey.Bytes(),
					new(big.Int).SetUint64(0),
					cmpToEVM(delPubKey.Bytes()),
					delAmtGwei,
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.RedelegateEvent.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check delegation from source validator
				_, err := sk.GetDelegation(ctx, delAccAddrFromEVM, val1ValAddr)
				require.ErrorContains(t, err, stypes.ErrNoDelegation.Error())

				// check delegatio from destination validator
				dstDelegation, err := sk.GetDelegation(ctx, delAccAddrFromEVM, val2ValAddr)
				require.NoError(t, err)
				require.NotNil(t, dstDelegation)

				// check validators
				srcVal, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				require.Equal(t, uint64(10000000), srcVal.Tokens.Uint64())

				dstVal, err := sk.GetValidator(ctx, val2ValAddr)
				require.NoError(t, err)
				require.Equal(t, uint64(10000100), dstVal.Tokens.Uint64())
			},
		},

		// ********** Withdraw **********
		{
			name: "fail(continue): invalid Withdraw event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				// create delegation
				createDelegation(t, ctx, sk, delAccAddrFromEVM, val1ValAddr, delCoin.Amount.Int64(), true)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.WithdrawEvent.ID, dummyHash}, TxHash: dummyHash}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check delegation from validator
				del, err := sk.GetDelegation(ctx, delAccAddrFromEVM, val1ValAddr)
				require.NoError(t, err)
				require.NotNil(t, del)
			},
		},
		{
			name: "fail(continue): process Withdraw event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				// create delegation
				createDelegation(t, ctx, sk, delAccAddrFromEVM, val1ValAddr, delCoin.Amount.Int64(), true)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["Withdraw"].Inputs.NonIndexed().Pack(
					delEVMAddr,
					invalidPubKey,
					delAmtGwei,
					new(big.Int).SetUint64(0),
					cmpToEVM(delPubKey.Bytes()),
					[]byte{},
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.WithdrawEvent.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check delegation from validator
				del, err := sk.GetDelegation(ctx, delAccAddrFromEVM, val1ValAddr)
				require.NoError(t, err)
				require.NotNil(t, del)
			},
		},
		{
			name: "pass: process Withdraw event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				// create delegation
				createDelegation(t, ctx, sk, delAccAddrFromEVM, val1ValAddr, delCoin.Amount.Int64(), true)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["Withdraw"].Inputs.NonIndexed().Pack(
					delEVMAddr,
					val1PubKey.Bytes(),
					delAmtGwei,
					new(big.Int).SetUint64(0),
					cmpToEVM(delPubKey.Bytes()),
					[]byte{},
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.WithdrawEvent.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check delegation from validator
				_, err := sk.GetDelegation(ctx, delAccAddrFromEVM, val1ValAddr)
				require.ErrorContains(t, err, stypes.ErrNoDelegation.Error())

				// check unbonding from the delegator
				ubd, err := sk.GetUnbondingDelegation(ctx, delAccAddrFromEVM, val1ValAddr)
				require.NoError(t, err)
				require.Len(t, ubd.Entries, 1)
				require.Equal(t, uint64(100), ubd.Entries[0].Balance.Uint64())

				// check validator
				val, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				require.Equal(t, uint64(10000000), val.Tokens.Uint64())
			},
		},

		// ********** Unjail **********
		{
			name: "fail(continue): invalid Unjail event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// create source validator
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				// set validator jailed
				val, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				val.Jailed = true
				_ = skeeper.TestingUpdateValidator(sk, ctx, val, true)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				logs := []ethtypes.Log{{Topics: []common.Hash{types.UnjailEvent.ID, dummyHash}, TxHash: dummyHash}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check validator status
				val, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				require.True(t, val.Jailed)
			},
		},
		{
			name: "fail(continue): process Unjail event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// create source validator
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				// set validator jailed
				val, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				val.Jailed = true
				_ = skeeper.TestingUpdateValidator(sk, ctx, val, true)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["Unjail"].Inputs.NonIndexed().Pack(
					delEVMAddr,
					invalidPubKey,
					[]byte{},
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.UnjailEvent.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check validator status
				val, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				require.True(t, val.Jailed)
			},
		},
		{
			name: "pass: process Unjail event",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, slk *estestutil.MockSlashingKeeper, sk *skeeper.Keeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
				slk.EXPECT().Unjail(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx sdk.Context, valValAddr sdk.ValAddress) error {
					val, err := sk.GetValidator(ctx, val1ValAddr)
					require.NoError(t, err)

					consAddr, err := val.GetConsAddr()
					require.NoError(t, err)

					return sk.Unjail(ctx, consAddr)
				})
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// create source validator
				createValidator(t, ctx, sk, val1PubKey, val1ValAddr, 0)

				// set validator jailed
				val, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				val.Jailed = true
				_ = skeeper.TestingUpdateValidator(sk, ctx, val, true)

				return ctx
			},
			evmEvents: func() []*evmenginetypes.EVMEvent {
				data, err := stakingAbi.Events["Unjail"].Inputs.NonIndexed().Pack(
					val1EVMAddr,
					val1PubKey.Bytes(),
					[]byte{},
				)
				require.NoError(t, err)

				logs := []ethtypes.Log{{
					Topics: []common.Hash{types.UnjailEvent.ID},
					Data:   data,
					TxHash: dummyHash,
				}}

				evmEvents, err := ethLogsToEvmEvents(logs)
				require.NoError(t, err)

				return evmEvents
			},
			postStateCheck: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) {
				// check validator status
				val, err := sk.GetValidator(ctx, val1ValAddr)
				require.NoError(t, err)
				require.False(t, val.Jailed)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, ak, bk, sk, slk, esk := createKeeperWithRealStaking(t)

			if tc.setupMocks != nil {
				tc.setupMocks(ak, bk, slk, sk)
			}

			cachedCtx, _ := ctx.CacheContext()
			if tc.setup != nil {
				cachedCtx = tc.setup(cachedCtx, sk, esk)
			}

			err = esk.ProcessStakingEvents(cachedCtx, 1, tc.evmEvents())
			if tc.expectedError != "" {
				require.ErrorContains(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				if tc.postStateCheck != nil {
					tc.postStateCheck(cachedCtx, sk, esk)
				}
			}
		})
	}
}

func createAddresses(count int) ([]crypto.PubKey, []sdk.AccAddress, []sdk.ValAddress) {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("story", "storypub")
	cfg.SetBech32PrefixForValidator("storyvaloper", "storyvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("storyvalcons", "storyvalconspub")

	var pubKeys []crypto.PubKey
	var accAddrs []sdk.AccAddress
	var valAddrs []sdk.ValAddress
	for range count {
		pubKey := k1.GenPrivKey().PubKey()
		evmAddr, _ := k1util.CosmosPubkeyToEVMAddress(pubKey.Bytes())
		accAddr := sdk.AccAddress(pubKey.Address().Bytes())
		valAddr := sdk.ValAddress(evmAddr.Bytes())
		pubKeys = append(pubKeys, pubKey)
		accAddrs = append(accAddrs, accAddr)
		valAddrs = append(valAddrs, valAddr)
	}

	return pubKeys, accAddrs, valAddrs
}

func createKeeperWithRealStaking(t *testing.T) (sdk.Context, *estestutil.MockAccountKeeper, *estestutil.MockBankKeeper, *skeeper.Keeper, *estestutil.MockSlashingKeeper, *keeper.Keeper) {
	t.Helper()

	ctx, storeKey, storeService := setupCtxStore(t, nil)
	cdc := getCodec()

	ctrl := gomock.NewController(t)

	ak := estestutil.NewMockAccountKeeper(ctrl)
	bk := estestutil.NewMockBankKeeper(ctrl)
	dk := estestutil.NewMockDistributionKeeper(ctrl)
	slk := estestutil.NewMockSlashingKeeper(ctrl)

	// mock keeper funcs
	ak.EXPECT().AddressCodec().Return(authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())).AnyTimes()
	ak.EXPECT().GetModuleAddress(gomock.Any()).Return(authtypes.NewModuleAddress(types.ModuleName)).Times(3)

	ethCl, err := ethclient.NewEngineMock(storeKey)
	require.NoError(t, err)

	sk := skeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		ak,
		bk,
		authtypes.NewModuleAddress(stypes.ModuleName).String(),
		address.NewBech32Codec("storyvaloper"),
		address.NewBech32Codec("storyvalcons"),
	)

	// set min delegation amount to 2
	var stakingParams = stypes.NewParams(
		stypes.DefaultUnbondingTime,
		stypes.DefaultMaxValidators,
		stypes.DefaultMaxEntries,
		stypes.DefaultHistoricalEntries,
		sdk.DefaultBondDenom,
		stypes.DefaultMinCommissionRate,
		math.NewInt(2),
		stypes.DefaultFlexiblePeriodType,
		stypes.DefaultPeriods,
		stypes.DefaultLockedTokenType,
		stypes.DefaultTokenTypes,
		stypes.DefaultSingularityHeight,
	)
	require.NoError(t, sk.SetParams(ctx, stakingParams))

	esk := keeper.NewKeeper(
		cdc,
		storeService,
		ak,
		bk,
		slk,
		sk,
		dk,
		authtypes.NewModuleAddress(types.ModuleName).String(),
		ethCl,
		address.NewBech32Codec("storyvaloper"),
	)
	require.NoError(t, esk.SetParams(ctx, types.DefaultParams()))

	return ctx, ak, bk, sk, slk, esk
}

func createKeeperWithMockStaking(t *testing.T) (sdk.Context, *estestutil.MockAccountKeeper, *estestutil.MockBankKeeper, *estestutil.MockDistributionKeeper, *estestutil.MockStakingKeeper, *estestutil.MockSlashingKeeper, *keeper.Keeper) {
	t.Helper()

	ctx, storeKey, storeService := setupCtxStore(t, nil)
	cdc := getCodec()

	ctrl := gomock.NewController(t)

	ak := estestutil.NewMockAccountKeeper(ctrl)
	bk := estestutil.NewMockBankKeeper(ctrl)
	dk := estestutil.NewMockDistributionKeeper(ctrl)
	slk := estestutil.NewMockSlashingKeeper(ctrl)
	stk := estestutil.NewMockStakingKeeper(ctrl)

	// mock keeper funcs
	ak.EXPECT().AddressCodec().Return(authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())).AnyTimes()
	ak.EXPECT().GetModuleAddress(gomock.Any()).Return(authtypes.NewModuleAddress(types.ModuleName))

	ethCl, err := ethclient.NewEngineMock(storeKey)
	require.NoError(t, err)

	esk := keeper.NewKeeper(
		cdc,
		storeService,
		ak,
		bk,
		slk,
		stk,
		dk,
		authtypes.NewModuleAddress(types.ModuleName).String(),
		ethCl,
		address.NewBech32Codec("storyvaloper"),
	)
	require.NoError(t, esk.SetParams(ctx, types.DefaultParams()))

	return ctx, ak, bk, dk, stk, slk, esk
}

func createQueryClient(ctx sdk.Context, esk *keeper.Keeper) types.QueryClient {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("story", "storypub")
	cfg.SetBech32PrefixForValidator("storyvaloper", "storyvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("storyvalcons", "storyvalconspub")

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, esk)

	return types.NewQueryClient(queryHelper)
}

func setupCtxStore(t *testing.T, header *cmtproto.Header) (sdk.Context, *storetypes.KVStoreKey, corestore.KVStoreService) {
	t.Helper()

	key := storetypes.NewKVStoreKey("test")
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	if header == nil {
		header = &cmtproto.Header{Time: cmttime.Now()}
	}
	ctx := testCtx.Ctx.WithBlockHeader(*header)
	defaultConsensusParams := genutil.DefaultConsensusParams()
	ctx = ctx.WithConsensusParams(defaultConsensusParams.ToProto())

	return ctx, key, storeService
}

func getCodec() codec.Codec {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	legacyAmino := codec.NewLegacyAmino()
	authtypes.RegisterLegacyAminoCodec(legacyAmino)
	authtypes.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterLegacyAminoCodec(legacyAmino)
	banktypes.RegisterInterfaces(interfaceRegistry)
	stypes.RegisterLegacyAminoCodec(legacyAmino)
	stypes.RegisterInterfaces(interfaceRegistry)

	return codec.NewProtoCodec(interfaceRegistry)
}

// createValidator creates a validator.
func createValidator(t *testing.T, ctx context.Context, sKeeper *skeeper.Keeper, valPubKey crypto.PubKey, valAddr sdk.ValAddress, supportTokenType int32) {
	t.Helper()

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Convert public key to cosmos format
	valCosmosPubKey, err := k1util.PubKeyToCosmos(valPubKey)
	require.NoError(t, err)

	// Create and update validator
	val, err := stypes.NewValidator(valAddr.String(), valCosmosPubKey, stypes.Description{Moniker: "test"}, supportTokenType)
	require.NoError(t, err)

	// Set commission
	val.Commission = stypes.NewCommission(
		math.LegacyNewDecWithPrec(int64(1000), 4),
		math.LegacyNewDecWithPrec(int64(5000), 4),
		math.LegacyNewDecWithPrec(int64(1000), 4),
	)

	valTokens := sKeeper.TokensFromConsensusPower(ctx, 10)
	validator, _, _ := val.AddTokensFromDel(valTokens, math.LegacyNewDecFromInt(valTokens).Quo(math.LegacyNewDec(2)))
	_ = skeeper.TestingUpdateValidator(sKeeper, sdkCtx, validator, true)
	require.NoError(t, sKeeper.SetValidatorByConsAddr(ctx, validator))
}

func createDelegation(t *testing.T, ctx context.Context, sKeeper *skeeper.Keeper, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount int64, withPeriodDelegation bool) {
	t.Helper()

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// set delegation
	delegation := stypes.NewDelegation(delAddr.String(), valAddr.String(), math.LegacyNewDec(amount), math.LegacyNewDec(amount).Quo(math.LegacyNewDec(2)))
	require.NoError(t, sKeeper.SetDelegation(sdkCtx, delegation))

	// set period delegation
	if withPeriodDelegation {
		periodDelegation := stypes.NewPeriodDelegation(delAddr.String(), valAddr.String(), stypes.FlexiblePeriodDelegationID, math.LegacyNewDec(amount), math.LegacyNewDec(amount).Quo(math.LegacyNewDec(2)), stypes.DefaultFlexiblePeriodType, time.Time{})
		require.NoError(t, sKeeper.SetPeriodDelegation(sdkCtx, delAddr, valAddr, periodDelegation))
	}

	// add to validator
	validator, err := sKeeper.GetValidator(ctx, valAddr)
	require.NoError(t, err)
	validator, _, _ = validator.AddTokensFromDel(math.NewInt(amount), math.LegacyNewDecFromInt(math.NewInt(amount)).Quo(math.LegacyNewDec(2)))
	_ = skeeper.TestingUpdateValidator(sKeeper, sdkCtx, validator, true)
}

func cmpToEVM(cmpPubKey []byte) common.Address {
	evmAddr, err := keeper.CmpPubKeyToEVMAddress(cmpPubKey)
	if err != nil {
		panic(err)
	}

	return evmAddr
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
			TxHash:  l.TxHash.Bytes(),
		})
	}

	for _, log := range events {
		if err := log.Verify(); err != nil {
			return nil, errors.Wrap(err, "verify log")
		}
	}

	return events, nil
}
