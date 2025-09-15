package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/epochs/types"
)

func TestEpochsExportGenesis(t *testing.T) {
	ctx, epochsKeeper := Setup(t)

	chainStartTime := ctx.BlockTime()
	chainStartHeight := ctx.BlockHeight()

	genesis, err := epochsKeeper.ExportGenesis(ctx)
	require.NoError(t, err)
	require.Len(t, genesis.Epochs, 4)

	expectedEpochs := types.DefaultGenesis().Epochs
	for i := range expectedEpochs {
		expectedEpochs[i].CurrentEpochStartHeight = chainStartHeight
		expectedEpochs[i].StartTime = chainStartTime
	}
	require.Equal(t, expectedEpochs, genesis.Epochs)
}

func TestEpochsInitGenesis(t *testing.T) {
	ctx, epochsKeeper := Setup(t)

	// On init genesis, default epochs information is set
	// To check init genesis again, should make it fresh status
	epochInfos, err := epochsKeeper.AllEpochInfos(ctx)
	require.NoError(t, err)
	for _, epochInfo := range epochInfos {
		err := epochsKeeper.EpochInfo.Remove(ctx, epochInfo.Identifier)
		require.NoError(t, err)
	}

	// now := time.Now()
	ctx = ctx.WithBlockHeight(1).WithBlockTime(time.Now().UTC())

	// test genesisState validation
	genesisState := types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
		},
	}
	require.EqualError(t, genesisState.Validate(), "epoch identifier should be unique")

	genesisState = types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
		},
	}

	err = epochsKeeper.InitGenesis(ctx, genesisState)
	require.NoError(t, err)
	epochInfo, err := epochsKeeper.EpochInfo.Get(ctx, "monthly")
	require.NoError(t, err)
	require.Equal(t, "monthly", epochInfo.Identifier)
	require.Equal(t, ctx.BlockTime().UTC().String(), epochInfo.StartTime.UTC().String())
	require.Equal(t, time.Hour*24, epochInfo.Duration)
	require.Equal(t, int64(0), epochInfo.CurrentEpoch)
	require.Equal(t, ctx.BlockHeight(), epochInfo.CurrentEpochStartHeight)
	require.Equal(t, time.Time{}.UTC().String(), epochInfo.CurrentEpochStartTime.UTC().String())
	require.True(t, epochInfo.EpochCountingStarted)
}
