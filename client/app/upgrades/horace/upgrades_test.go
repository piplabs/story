package horace

import (
	"context"
	"cosmossdk.io/core/address"
	"errors"
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/piplabs/story/client/app/upgrades/horace/testutil"
	minttypes "github.com/piplabs/story/client/x/mint/types"
)

// initBech32Config ensures address conversion helpers work in tests.
func initBech32Config() {
	cfg := sdk.GetConfig()
	// These values must match your chain config if you rely on specific prefixes.
	cfg.SetBech32PrefixForAccount("story", "storypub")
	cfg.SetBech32PrefixForValidator("storyvaloper", "storyvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("storyvalcons", "storyvalconspub")
	cfg.Seal()
}

// stubAddressCodec is a small test helper implementing address.Codec.
// It allows injecting errors for StringToBytes/BytesToString.
type stubAddressCodec struct {
	stringToBytes func(string) ([]byte, error)
	bytesToString func([]byte) (string, error)
}

var _ address.Codec = (*stubAddressCodec)(nil)

func (s *stubAddressCodec) StringToBytes(str string) ([]byte, error) {
	if s.stringToBytes != nil {
		return s.stringToBytes(str)
	}
	return nil, errors.New("StringToBytes not configured")
}

func (s *stubAddressCodec) BytesToString(bz []byte) (string, error) {
	if s.bytesToString != nil {
		return s.bytesToString(bz)
	}
	return "", errors.New("BytesToString not configured")
}

func mustValBech32(t *testing.T, bz []byte) string {
	t.Helper()
	return sdk.ValAddress(bz).String()
}

func defaultStakingParamsLocked() stypes.Params {
	p := stypes.DefaultParams()

	// Deep copy TokenTypes to avoid sharing the backing array returned by DefaultParams.
	// This prevents mutations in one test case from leaking into other tests.
	p.TokenTypes = append([]stypes.TokenTypeInfo(nil), p.TokenTypes...)
	return p
}

func makeValidator(operatorBech32 string, supportTokenType int32, tokens math.Int, rewardsTokens math.LegacyDec) stypes.Validator {
	return stypes.Validator{
		OperatorAddress:        operatorBech32,
		SupportTokenType:       supportTokenType,
		Tokens:                 tokens,
		DelegatorRewardsShares: math.LegacyZeroDec(),
		RewardsTokens:          rewardsTokens,
	}
}

func makeDelegation(delegatorBech32 string, validatorBech32 string, shares math.LegacyDec) stypes.Delegation {
	return stypes.Delegation{
		DelegatorAddress: delegatorBech32,
		ValidatorAddress: validatorBech32,
		Shares:           shares,
		RewardsShares:    shares,
	}
}

func TestRunHoraceUpgrade(t *testing.T) {
	// Do not run in parallel because sdk.GetConfig().Seal() is global.
	initBech32Config()

	ctx := context.Background()

	validDelAcc := sdk.AccAddress([]byte("delegator_address_20")) // 20 bytes-ish
	validValAcc := sdk.ValAddress([]byte("validator_address_20"))

	delBech := validDelAcc.String()
	valBech := validValAcc.String()

	tests := []struct {
		name       string
		setupMocks func(
			t *testing.T,
			ctrl *gomock.Controller,
			ak *testutil.MockAccountKeeper,
			sk *testutil.MockStakingKeeper,
			dk *testutil.MockDistributionKeeper,
			mk *testutil.MockMintKeeper,
		)
		wantErr bool
	}{
		{
			name: "fail: staking GetParams error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				sk.EXPECT().GetParams(gomock.Any()).Return(stypes.Params{}, errors.New("boom"))
			},
			wantErr: true,
		},
		{
			name: "fail: locked token type mismatch in TokenTypes[0]",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				p.TokenTypes[0].TokenType = 1
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: staking SetParams error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(errors.New("set params failed")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "fail: staking reload GetParams error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(stypes.Params{}, errors.New("reload failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: locked multiplier not updated on reload",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				// p2 keeps the old multiplier to trigger the verification failure.
				p2.TokenTypes[0].RewardsMultiplier = math.LegacyNewDec(1)

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: staking GetAllValidators error (initial)",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, errors.New("get validators failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: oldTotalRewardsTokens is zero (edge case: all rewards tokens zero)",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				v1 := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyZeroDec())
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{v1}, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: invalid validator operator address bech32",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				v1 := makeValidator("not-a-bech32", 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{v1}, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: GetValidatorDelegations error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("get dels failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: WithdrawValidatorCommission error (not ErrNoValidatorCommission)",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("commission withdraw failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: delegation bech32 invalid during WithdrawDelegationRewards",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				del := makeDelegation("bad-del-bech32", valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{del}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
			},
			wantErr: true,
		},
		{
			name: "fail: WithdrawDelegationRewards error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				del := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{del}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("withdraw del rewards failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: SetDelegation error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return sdk.AccAddress([]byte("delegator_address_20")).Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				del := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{del}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).
					Return(errors.New("set delegation failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: AddressCodec StringToBytes error (delegator address)",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return nil, errors.New("codec error") },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				del := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{del}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "fail: GetDelegation error during verification",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				del := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{del}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(stypes.Delegation{}, errors.New("get delegation failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: delegation rewards shares mismatch after SetDelegation",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)

				// Return a delegation with an incorrect RewardsShares to trigger the mismatch.
				bad := orig
				bad.RewardsShares = math.LegacyNewDec(9999)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(bad, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: GetPeriodDelegation error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)

				// Return correct updated delegation so it passes the previous check.
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(updated, nil)

				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(stypes.PeriodDelegation{}, errors.New("get period del failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: SetPeriodDelegation error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().
					SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("set period delegation failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: GetPeriodDelegation after update error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().
					SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				sk.EXPECT().
					GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(stypes.PeriodDelegation{}, errors.New("reload period delegation failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: period delegation rewards shares mismatch",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().
					SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				pd2 := stypes.PeriodDelegation{
					RewardsShares: math.LegacyNewDec(9999),
				}
				sk.EXPECT().
					GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(pd2, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: SetValidator error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation round trip succeeds (minimal).
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				// Fail here.
				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(errors.New("set validator failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: GetValidator error after update",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation round trip succeeds (minimal).
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().
					GetValidator(gomock.Any(), gomock.Any()).
					Return(stypes.Validator{}, errors.New("get validator failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: delegator rewards shares mismatch",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation round trip succeeds (minimal).
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				v := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				v.DelegatorRewardsShares = math.LegacyNewDec(999)

				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(v, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: validator rewards tokens mismatch",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation round trip succeeds (minimal).
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDec(999)

				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: GetValidatorDelegations during starting info update",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation round trip succeeds (minimal).
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				sk.EXPECT().
					GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("get delegations failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: GetDelegatorStartingInfo error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).
					Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation round trip succeeds (minimal).
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)

				dk.EXPECT().
					GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{}, errors.New("get starting info failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: SetDelegatorStartingInfo error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				// Update validator.
				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				// Updated delegations for starting info scaling.
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)

				// Starting info scaling.
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)

				dk.EXPECT().
					SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("set starting info failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: GetDelegatorStartingInfo after updating the starting info",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				// Update validator.
				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				// Updated delegations for starting info scaling.
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)

				// Starting info scaling.
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)

				dk.EXPECT().
					SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, errors.New("failed to get starting info"))
			},
			wantErr: true,
		},
		{
			name: "fail: scaled starting info rewards stake mismatch",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				// Update validator.
				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				// Updated delegations for starting info scaling.
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)

				// Starting info scaling.
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)

				dk.EXPECT().
					SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				dk.EXPECT().
					GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: math.LegacyNewDec(999)}, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: unlocked rewards tokens changed (edge case across GetAllValidators calls)",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				vUnlocked := makeValidator(mustValBech32(t, []byte("unlocked_val_addr20")), 1, math.NewInt(2000), math.LegacyNewDec(20))

				// First GetAllValidators returns both.
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked, vUnlocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				// Validator set/update/validation should pass.
				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)
				// Return a validator with matching expected fields.
				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				// Updated delegations list for starting info scaling.
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)

				// Starting info round trip.
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)
				dk.EXPECT().SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				expectedStake := newV.RewardsTokensFromRewardsSharesTruncated(updated.RewardsShares)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: expectedStake}, nil)

				// Second GetAllValidators returns the unlocked validator with changed rewards tokens to trigger the check.
				vUnlockedChanged := vUnlocked
				vUnlockedChanged.RewardsTokens = math.LegacyNewDec(999)
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{newV, vUnlockedChanged}, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: GetAllValidators after update error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				// Update validator.
				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				// Updated delegations for starting info scaling.
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)

				// Starting info scaling.
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)
				dk.EXPECT().SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				expectedStake := newV.RewardsTokensFromRewardsSharesTruncated(updated.RewardsShares)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: expectedStake}, nil)

				sk.EXPECT().
					GetAllValidators(gomock.Any()).
					Return(nil, errors.New("get validators failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: new total rewards tokens is zero",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				// Update validator.
				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				// Updated delegations for starting info scaling.
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)

				// Starting info scaling.
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)
				dk.EXPECT().SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				expectedStake := newV.RewardsTokensFromRewardsSharesTruncated(updated.RewardsShares)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: expectedStake}, nil)

				v := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyZeroDec())
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{v}, nil)
			},
			wantErr: true,
		},
		{
			name: "pass: happy path (minimal single locked validator, single delegation, no unlocked validators)",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				// Update validator.
				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				// Updated delegations for starting info scaling.
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)

				// Starting info scaling.
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)
				dk.EXPECT().SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				expectedStake := newV.RewardsTokensFromRewardsSharesTruncated(updated.RewardsShares)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: expectedStake}, nil)

				// Second GetAllValidators for post-check.
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{newV}, nil)

				// Mint params update and verification.
				initialMint := minttypes.Params{
					MintDenom:         "stake",
					InflationsPerYear: math.LegacyNewDec(1),
					BlocksPerYear:     uint64(1),
				}
				// First GetParams.
				mk.EXPECT().GetParams(gomock.Any()).Return(initialMint, nil)
				// SetParams is successful.
				mk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				// Reload GetParams returns updated values.
				updatedMint := initialMint
				updatedMint.InflationsPerYear = NewAnnualInflationsPerYear
				updatedMint.BlocksPerYear = NewBlocksPerYear
				mk.EXPECT().GetParams(gomock.Any()).Return(updatedMint, nil)
			},
			wantErr: false,
		},
		{
			name: "fail: mint GetParams error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				// Update validator.
				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				// Updated delegations for starting info scaling.
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)

				// Starting info scaling.
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)
				dk.EXPECT().SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				expectedStake := newV.RewardsTokensFromRewardsSharesTruncated(updated.RewardsShares)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: expectedStake}, nil)

				// Second GetAllValidators for post-check.
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{newV}, nil)

				// Fail here.
				mk.EXPECT().GetParams(gomock.Any()).Return(minttypes.Params{}, errors.New("mint get params failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: mint params validation error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				// Staking params update and verification.
				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// One locked validator with non-zero rewards tokens.
				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				// One delegation.
				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				// Withdraw commission: no commission is fine.
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)

				// Withdraw delegation rewards.
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				// Set delegation and verify it is updated.
				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				// Period delegation updated and verified.
				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				// Update validator.
				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)

				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				// Updated delegations for starting info scaling.
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)

				// Starting info scaling.
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)
				dk.EXPECT().SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				expectedStake := newV.RewardsTokensFromRewardsSharesTruncated(updated.RewardsShares)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: expectedStake}, nil)

				// Second GetAllValidators for post-check.
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{newV}, nil)

				// Mint params update and verification.
				initialMint := minttypes.Params{
					InflationsPerYear: math.LegacyNewDec(1),
					BlocksPerYear:     uint64(1),
				}
				// First GetParams.
				mk.EXPECT().GetParams(gomock.Any()).Return(initialMint, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: mint SetParams error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				// Same as happy path until mint section.
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)
				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)
				dk.EXPECT().SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				expectedStake := newV.RewardsTokensFromRewardsSharesTruncated(updated.RewardsShares)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: expectedStake}, nil)

				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{newV}, nil)

				initialMint := minttypes.Params{
					MintDenom:         "stake",
					InflationsPerYear: math.LegacyNewDec(1),
					BlocksPerYear:     uint64(1),
				}
				mk.EXPECT().GetParams(gomock.Any()).Return(initialMint, nil)
				mk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(errors.New("mint set params failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: mint reload GetParams error",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				// Same as happy path until mint section.
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)
				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)
				dk.EXPECT().SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				expectedStake := newV.RewardsTokensFromRewardsSharesTruncated(updated.RewardsShares)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: expectedStake}, nil)

				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{newV}, nil)

				initialMint := minttypes.Params{
					MintDenom:         "stake",
					InflationsPerYear: math.LegacyNewDec(1),
					BlocksPerYear:     uint64(1),
				}
				mk.EXPECT().GetParams(gomock.Any()).Return(initialMint, nil)
				mk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				mk.EXPECT().GetParams(gomock.Any()).Return(minttypes.Params{}, errors.New("mint reload failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: mint inflations_per_year not updated on reload",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				// Same as happy path until mint section.
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)
				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)
				dk.EXPECT().SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				expectedStake := newV.RewardsTokensFromRewardsSharesTruncated(updated.RewardsShares)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: expectedStake}, nil)

				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{newV}, nil)

				initialMint := minttypes.Params{
					MintDenom:         "stake",
					InflationsPerYear: math.LegacyNewDec(1),
					BlocksPerYear:     uint64(1),
				}
				mk.EXPECT().GetParams(gomock.Any()).Return(initialMint, nil)
				mk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)

				// Return mint params without updated inflation to trigger the check.
				badReload := initialMint
				badReload.BlocksPerYear = NewBlocksPerYear
				// Inflation remains old (1)
				mk.EXPECT().GetParams(gomock.Any()).Return(badReload, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: mint blocks per year mismatch after reload",
			setupMocks: func(t *testing.T, ctrl *gomock.Controller, ak *testutil.MockAccountKeeper, sk *testutil.MockStakingKeeper, dk *testutil.MockDistributionKeeper, mk *testutil.MockMintKeeper) {
				// Same as happy path until mint section.
				codec := &stubAddressCodec{
					stringToBytes: func(s string) ([]byte, error) { return validDelAcc.Bytes(), nil },
				}
				ak.EXPECT().AddressCodec().Return(codec).AnyTimes()

				p := defaultStakingParamsLocked()
				p2 := defaultStakingParamsLocked()
				p2.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				vLocked := makeValidator(valBech, 0, math.NewInt(1000), math.LegacyNewDec(10))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{vLocked}, nil)

				orig := makeDelegation(delBech, valBech, math.LegacyNewDec(100))
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{orig}, nil)

				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.Coins{}, nil)

				sk.EXPECT().SetDelegation(gomock.Any(), gomock.Any()).Return(nil)
				updated := orig
				updated.RewardsShares = orig.Shares.MulTruncate(NewLockedTokenMultiplier)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(updated, nil)

				pd := stypes.PeriodDelegation{RewardsShares: updated.RewardsShares}
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pd, nil)

				sk.EXPECT().SetValidator(gomock.Any(), gomock.Any()).Return(nil)
				newV := vLocked
				newV.DelegatorRewardsShares = updated.RewardsShares
				newV.RewardsTokens = math.LegacyNewDecFromInt(vLocked.Tokens).Mul(NewLockedTokenMultiplier)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(newV, nil)

				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{updated}, nil)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(dtypes.DelegatorStartingInfo{}, nil)
				dk.EXPECT().SetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				expectedStake := newV.RewardsTokensFromRewardsSharesTruncated(updated.RewardsShares)
				dk.EXPECT().GetDelegatorStartingInfo(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dtypes.DelegatorStartingInfo{RewardsStake: expectedStake}, nil)

				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{newV}, nil)

				initialMint := minttypes.Params{
					MintDenom:         "stake",
					InflationsPerYear: math.LegacyNewDec(1),
					BlocksPerYear:     uint64(1),
				}
				mk.EXPECT().GetParams(gomock.Any()).Return(initialMint, nil)
				mk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)

				// Return mint params without updated inflation to trigger the check.
				badReload := initialMint
				badReload.InflationsPerYear = NewAnnualInflationsPerYear
				badReload.BlocksPerYear = 999

				mk.EXPECT().GetParams(gomock.Any()).Return(badReload, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ak := testutil.NewMockAccountKeeper(ctrl)
			sk := testutil.NewMockStakingKeeper(ctrl)
			dk := testutil.NewMockDistributionKeeper(ctrl)
			mk := testutil.NewMockMintKeeper(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(t, ctrl, ak, sk, dk, mk)
			}

			err := runHoraceUpgrade(ctx, ak, sk, dk, mk)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
