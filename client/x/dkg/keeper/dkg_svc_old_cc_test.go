package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

func TestGetOldCodeCommitment_NoPreviousActive(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.validatorEVMAddr = testValidatorAddr

	cc, err := k.getOldCodeCommitment(ctx)
	require.NoError(t, err)
	require.Nil(t, cc, "should return nil when no previous active round exists")
}

func TestGetOldCodeCommitment_WithPreviousActive(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.validatorEVMAddr = testValidatorAddr

	// Set up previous active round
	network := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))
	require.NoError(t, k.setLatestActiveRound(ctx, network))

	// Set up registration for validator in previous round
	expectedCC := []byte("previous-code-commitment")
	require.NoError(t, k.setDKGRegistration(ctx, common.HexToAddress(testValidatorAddr), &types.DKGRegistration{
		Round:          1,
		ValidatorAddr:  testValidatorAddr,
		DkgPubKey:      []byte("pub-key"),
		CodeCommitment: expectedCC,
	}))

	cc, err := k.getOldCodeCommitment(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedCC, cc)
}

func TestGetOldCodeCommitment_NoRegistration(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.validatorEVMAddr = testValidatorAddr

	// Set up previous active round but no registration for this validator
	network := &types.DKGNetwork{
		Round: 2,
		Stage: types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))
	require.NoError(t, k.setLatestActiveRound(ctx, network))

	_, err := k.getOldCodeCommitment(ctx)
	require.Error(t, err, "should return error when no registration exists")
}
