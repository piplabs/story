package virgil

import (
	"context"
	"errors"
	"testing"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/piplabs/story/client/app/upgrades/virgil/testutil"
	"github.com/piplabs/story/lib/netconf"
)

func initBech32Config() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("story", "storypub")
	cfg.SetBech32PrefixForValidator("storyvaloper", "storyvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("storyvalcons", "storyvalconspub")
	cfg.Seal()
}

// defaultStakingParamsWithPeriods returns staking params that contain short (1), medium (2), and long (3) periods.
func defaultStakingParamsWithPeriods() stypes.Params {
	p := stypes.DefaultParams()
	p.Periods = []stypes.Period{
		{
			PeriodType:        1,
			Duration:          time.Second * 3600, // 1 hour (old short)
			RewardsMultiplier: math.LegacyNewDec(1),
		},
		{
			PeriodType:        2,
			Duration:          time.Second * 86400, // 1 day
			RewardsMultiplier: math.LegacyNewDec(1),
		},
		{
			PeriodType:        3,
			Duration:          time.Second * 604800, // 7 days
			RewardsMultiplier: math.LegacyNewDec(1),
		},
	}
	return p
}

func makePeriodDelegation(delBech, valBech string, periodType int32, pdID string, endTime time.Time) stypes.PeriodDelegation {
	return stypes.PeriodDelegation{
		DelegatorAddress:   delBech,
		ValidatorAddress:   valBech,
		Shares:             math.LegacyNewDec(100),
		RewardsShares:      math.LegacyNewDec(100),
		PeriodDelegationId: pdID,
		PeriodType:         periodType,
		EndTime:            endTime,
	}
}

func TestRunVirgilUpgrade(t *testing.T) {
	initBech32Config()

	ctx := context.Background()

	validDelAcc := sdk.AccAddress([]byte("delegator_address_20"))
	validValAcc := sdk.ValAddress([]byte("validator_address_20"))

	delBech := validDelAcc.String()
	valBech := validValAcc.String()

	chainID := netconf.TestChainID
	expectedMultipliers := GetRewardsMultipliers(chainID)
	oldShortDuration := time.Second * 3600

	baseEndTime := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	expectedNewEndTime := baseEndTime.Add(NewShortPeriodDuration - oldShortDuration)

	tests := []struct {
		name       string
		setupMocks func(sk *testutil.MockStakingKeeper)
		wantErr    bool
	}{
		{
			name: "fail: staking GetParams error",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				sk.EXPECT().GetParams(gomock.Any()).Return(stypes.Params{}, errors.New("boom"))
			},
			wantErr: true,
		},
		{
			name: "fail: staking SetParams error",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(errors.New("set params failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: staking reload GetParams error",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(stypes.Params{}, errors.New("reload failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: short period duration not updated on reload",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)

				// Return params where short period duration is still the old value.
				p2 := defaultStakingParamsWithPeriods()
				// Duration NOT updated — still old value.
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: short period rewards multiplier not updated on reload",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)

				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				// RewardsMultiplier NOT updated — still old value.
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: medium period rewards multiplier not updated on reload",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)

				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				// Medium NOT updated.
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: long period rewards multiplier not updated on reload",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)

				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				// Long NOT updated.
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: GetAllPeriodDelegations error",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return(nil, errors.New("get all period delegations failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: invalid delegator bech32 in period delegation",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				pd := makePeriodDelegation("bad-del-bech32", valBech, 1, "pd1", baseEndTime)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{pd}, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: invalid validator bech32 in period delegation",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				pd := makePeriodDelegation(delBech, "bad-val-bech32", 1, "pd1", baseEndTime)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{pd}, nil)
			},
			wantErr: true,
		},
		{
			name: "fail: SetPeriodDelegation error",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				pd := makePeriodDelegation(delBech, valBech, 1, "pd1", baseEndTime)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{pd}, nil)

				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("set period delegation failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: GetPeriodDelegation after update error",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				pd := makePeriodDelegation(delBech, valBech, 1, "pd1", baseEndTime)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{pd}, nil)

				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(stypes.PeriodDelegation{}, errors.New("get period delegation failed"))
			},
			wantErr: true,
		},
		{
			name: "fail: period delegation end time mismatch after update",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				pd := makePeriodDelegation(delBech, valBech, 1, "pd1", baseEndTime)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{pd}, nil)

				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				// Return period delegation with wrong end time.
				badPd := pd
				badPd.EndTime = time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(badPd, nil)
			},
			wantErr: true,
		},
		{
			name: "pass: no period delegations (empty list)",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{}, nil)
			},
			wantErr: false,
		},
		{
			name: "pass: skip non-short period delegations (only medium/long present)",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				// Only medium (2) and long (3) period delegations — should be skipped.
				pdMedium := makePeriodDelegation(delBech, valBech, 2, "pd2", baseEndTime)
				pdLong := makePeriodDelegation(delBech, valBech, 3, "pd3", baseEndTime)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{pdMedium, pdLong}, nil)
			},
			wantErr: false,
		},
		{
			name: "pass: happy path with single short period delegation",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				pd := makePeriodDelegation(delBech, valBech, 1, "pd1", baseEndTime)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{pd}, nil)

				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				updatedPd := pd
				updatedPd.EndTime = expectedNewEndTime
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(updatedPd, nil)
			},
			wantErr: false,
		},
		{
			name: "pass: happy path with mixed period delegations (short + medium + long)",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)

				pdShort := makePeriodDelegation(delBech, valBech, 1, "pd1", baseEndTime)
				pdMedium := makePeriodDelegation(delBech, valBech, 2, "pd2", baseEndTime)
				pdLong := makePeriodDelegation(delBech, valBech, 3, "pd3", baseEndTime)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).
					Return([]stypes.PeriodDelegation{pdShort, pdMedium, pdLong}, nil)

				// Only short period delegation should be modified.
				sk.EXPECT().SetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				updatedPd := pdShort
				updatedPd.EndTime = expectedNewEndTime
				sk.EXPECT().GetPeriodDelegation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(updatedPd, nil)
			},
			wantErr: false,
		},
		{
			name: "pass: default case in update switch — unknown period type is skipped",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := defaultStakingParamsWithPeriods()
				// Append an unknown period type (99) — hits default in the first switch.
				p.Periods = append(p.Periods, stypes.Period{
					PeriodType:        99,
					Duration:          time.Second * 999,
					RewardsMultiplier: math.LegacyNewDec(1),
				})

				p2 := defaultStakingParamsWithPeriods()
				p2.Periods[0].Duration = NewShortPeriodDuration
				p2.Periods[0].RewardsMultiplier = expectedMultipliers.Short
				p2.Periods[1].RewardsMultiplier = expectedMultipliers.Medium
				p2.Periods[2].RewardsMultiplier = expectedMultipliers.Long
				// Same unknown period type appears in reload — hits default in the second switch.
				p2.Periods = append(p2.Periods, stypes.Period{
					PeriodType:        99,
					Duration:          time.Second * 999,
					RewardsMultiplier: math.LegacyNewDec(1),
				})

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p2, nil)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{}, nil)
			},
			wantErr: false,
		},
		{
			name: "pass: default case only — params contain only unknown period types",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := stypes.DefaultParams()
				p.Periods = []stypes.Period{
					{
						PeriodType:        42,
						Duration:          time.Second * 100,
						RewardsMultiplier: math.LegacyNewDec(1),
					},
					{
						PeriodType:        99,
						Duration:          time.Second * 200,
						RewardsMultiplier: math.LegacyNewDec(1),
					},
				}

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				// Reload returns same unknown-only periods — both switches hit only default.
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{}, nil)
			},
			wantErr: false,
		},
		{
			name: "pass: params with no periods (no verification failures, no delegations to sweep)",
			setupMocks: func(sk *testutil.MockStakingKeeper) {
				p := stypes.DefaultParams()
				p.Periods = nil

				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil)
				sk.EXPECT().SetParams(gomock.Any(), gomock.Any()).Return(nil)
				sk.EXPECT().GetParams(gomock.Any()).Return(p, nil) // reload returns same (no periods)
				sk.EXPECT().GetAllPeriodDelegations(gomock.Any()).Return([]stypes.PeriodDelegation{}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sk := testutil.NewMockStakingKeeper(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(sk)
			}

			err := runVirgilUpgrade(ctx, sk, chainID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
