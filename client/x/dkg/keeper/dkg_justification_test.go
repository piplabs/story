package keeper

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v4/group/edwards25519"
	"go.dedis.ch/kyber/v4/share"

	"github.com/piplabs/story/client/x/dkg/types"
)

// generateShareAndCommitments generates a valid Pedersen VSS share and commitment
// bytes for a given recipientIndex (1-based) using the Edwards25519 suite.
// The underlying polynomial uses threshold=2 (linear), producing 2 commitments.
func generateShareAndCommitments(t *testing.T, recipientIndex int) (shareBytes []byte, commitmentBytes [][]byte) {
	t.Helper()

	suite := edwards25519.NewBlakeSHA256Ed25519()
	n := 3
	threshold := int(testThreshold)

	secret := suite.Scalar().Pick(suite.RandomStream())
	priPoly := share.NewPriPoly(suite, threshold, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.Point().Base())

	shares := priPoly.Shares(n)
	// Info() returns (base kyber.Point, commits []kyber.Point)
	_, commits := pubPoly.Info()

	idx := recipientIndex - 1
	if idx < 0 || idx >= n {
		idx = 0
	}

	var err error
	shareBytes, err = shares[idx].V.MarshalBinary()
	require.NoError(t, err)

	commitmentBytes = make([][]byte, len(commits))
	for i, c := range commits {
		bz, err := c.MarshalBinary()
		require.NoError(t, err)
		commitmentBytes[i] = bz
	}

	return shareBytes, commitmentBytes
}

// threshold is the number of commitments in these tests.
const testThreshold = uint32(2)

// TestVerifyJustificationVSS_ValidShare verifies that a correctly generated share
// passes the VSS wrapper function.
func TestVerifyJustificationVSS_ValidShare(t *testing.T) {
	t.Parallel()

	shareBytes, commitmentBytes := generateShareAndCommitments(t, 1)

	ok, err := verifyJustificationVSS(shareBytes, 1, commitmentBytes, testThreshold)
	require.NoError(t, err)
	require.True(t, ok, "valid share should pass VSS verification")
}

// TestVerifyJustificationVSS_InvalidShare verifies that a random scalar fails
// verification against commitments it was not derived from.
func TestVerifyJustificationVSS_InvalidShare(t *testing.T) {
	t.Parallel()

	suite := edwards25519.NewBlakeSHA256Ed25519()
	_, commitmentBytes := generateShareAndCommitments(t, 1)

	// Use a completely different random scalar
	wrongScalar := suite.Scalar().Pick(suite.RandomStream())
	wrongBytes, err := wrongScalar.MarshalBinary()
	require.NoError(t, err)

	ok, err := verifyJustificationVSS(wrongBytes, 1, commitmentBytes, testThreshold)
	require.NoError(t, err)
	require.False(t, ok, "wrong scalar should fail VSS verification")
}

// TestVerifyJustificationVSS_EmptyCommitments verifies that empty commitments
// cause verifyJustificationVSS to return an error.
func TestVerifyJustificationVSS_EmptyCommitments(t *testing.T) {
	t.Parallel()

	shareBytes, _ := generateShareAndCommitments(t, 1)

	ok, err := verifyJustificationVSS(shareBytes, 1, [][]byte{}, 0)
	require.Error(t, err, "empty commitments should return an error")
	require.False(t, ok)
}

// TestVerifyJustificationVSS_MalformedShareBytes verifies that malformed share bytes
// produce an error.
func TestVerifyJustificationVSS_MalformedShareBytes(t *testing.T) {
	t.Parallel()

	_, commitmentBytes := generateShareAndCommitments(t, 1)

	ok, err := verifyJustificationVSS([]byte("not-a-scalar"), 1, commitmentBytes, testThreshold)
	require.Error(t, err, "malformed share bytes should return an error")
	require.False(t, ok)
}

// TestVerifyJustificationVSS_MalformedCommitmentBytes verifies that malformed
// commitment bytes produce an error.
func TestVerifyJustificationVSS_MalformedCommitmentBytes(t *testing.T) {
	t.Parallel()

	shareBytes, _ := generateShareAndCommitments(t, 1)

	// threshold=1 but we pass a malformed single commitment
	ok, err := verifyJustificationVSS(shareBytes, 1, [][]byte{[]byte("not-a-point")}, 1)
	require.Error(t, err, "malformed commitment bytes should return an error")
	require.False(t, ok)
}

// TestVerifyJustificationVSS_WrongIndex verifies that a valid share fails
// when checked against a different recipient index.
func TestVerifyJustificationVSS_WrongIndex(t *testing.T) {
	t.Parallel()

	// Generate share for index 1
	shareBytes, commitmentBytes := generateShareAndCommitments(t, 1)

	// Check at index 2 - should fail
	ok, err := verifyJustificationVSS(shareBytes, 2, commitmentBytes, testThreshold)
	require.NoError(t, err)
	require.False(t, ok, "share for index 1 should not verify at index 2")
}

// TestInvalidateDealerRegistration_Success verifies that a Verified dealer is
// successfully transitioned to Invalidated status.
func TestInvalidateDealerRegistration_Success(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	round := uint32(5)
	dealerAddr := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	dealerIndex := uint32(3)

	setupDealerRegistrationForInvalidation(t, k, ctx, round, dealerAddr, dealerIndex, types.DKGRegStatusVerified)

	network := &types.DKGNetwork{
		Round: round,
	}

	err := k.invalidateDealerRegistration(ctx, network, dealerIndex)
	require.NoError(t, err)

	reg, err := k.getDKGRegistration(ctx, round, dealerAddr)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusInvalidated, reg.Status)
}

// TestInvalidateDealerRegistration_AlreadyFinalized verifies that a dealer that
// has already finalized is not re-invalidated (the function returns nil without change).
func TestInvalidateDealerRegistration_AlreadyFinalized(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	round := uint32(1)
	dealerAddr := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	dealerIndex := uint32(2)

	// Set up as already Finalized
	setupDealerRegistrationForInvalidation(t, k, ctx, round, dealerAddr, dealerIndex, types.DKGRegStatusFinalized)

	network := &types.DKGNetwork{
		Round: round,
	}

	err := k.invalidateDealerRegistration(ctx, network, dealerIndex)
	require.Error(t, err, "should error when trying to invalidate a finalized dealer")
	require.Contains(t, err.Error(), "already finalized")
	require.Contains(t, err.Error(), "possible bug")

	// Status should remain Finalized
	reg, err := k.getDKGRegistration(ctx, round, dealerAddr)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusFinalized, reg.Status, "status should not change for finalized dealers")
}

// TestInvalidateDealerRegistration_NotFound verifies that an error is returned
// when no registration with the given dealer index exists.
func TestInvalidateDealerRegistration_NotFound(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	network := &types.DKGNetwork{
		Round: 1,
	}

	// No registration has been created; index 99 does not exist
	err := k.invalidateDealerRegistration(ctx, network, 99)
	require.Error(t, err, "should return error when dealer index is not found")
	require.Contains(t, err.Error(), "no registration found with dealer index 99")
}

// TestInvalidateDealerRegistration_MultipleRegistrations verifies that when
// multiple dealers exist only the correct one gets invalidated.
func TestInvalidateDealerRegistration_MultipleRegistrations(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	round := uint32(3)

	dealerA := common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")
	dealerB := common.HexToAddress("0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	indexA := uint32(1)
	indexB := uint32(2)

	setupDealerRegistrationForInvalidation(t, k, ctx, round, dealerA, indexA, types.DKGRegStatusVerified)
	setupDealerRegistrationForInvalidation(t, k, ctx, round, dealerB, indexB, types.DKGRegStatusVerified)

	network := &types.DKGNetwork{
		Round: round,
	}

	// Invalidate only dealer A
	err := k.invalidateDealerRegistration(ctx, network, indexA)
	require.NoError(t, err)

	regA, err := k.getDKGRegistration(ctx, round, dealerA)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusInvalidated, regA.Status, "dealer A should be invalidated")

	regB, err := k.getDKGRegistration(ctx, round, dealerB)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusVerified, regB.Status, "dealer B should remain Verified")
}

// TestDeduplicateJustifications_RemovesDuplicates verifies that justifications with
// the same (dealerIndex, recipientIndex) pair are deduplicated.
func TestDeduplicateJustifications_RemovesDuplicates(t *testing.T) {
	t.Parallel()

	j1 := types.Justification{Index: 1, VssJustification: &types.VSSJustification{
		PlainDeal: &types.PlainDeal{SecShare: &types.SecShare{I: 10}},
	}}
	j2 := types.Justification{Index: 1, VssJustification: &types.VSSJustification{
		PlainDeal: &types.PlainDeal{SecShare: &types.SecShare{I: 10}},
	}}
	j3 := types.Justification{Index: 2, VssJustification: &types.VSSJustification{
		PlainDeal: &types.PlainDeal{SecShare: &types.SecShare{I: 10}},
	}}

	result := deduplicateJustifications([]types.Justification{j1, j2, j3})
	require.Len(t, result, 2, "duplicate (1,10) should be removed")
	require.Equal(t, uint32(1), result[0].Index)
	require.Equal(t, uint32(2), result[1].Index)
}

// TestDeduplicateJustifications_Empty verifies that an empty slice returns empty.
func TestDeduplicateJustifications_Empty(t *testing.T) {
	t.Parallel()

	result := deduplicateJustifications(nil)
	require.Empty(t, result)
}

// TestDeduplicateJustifications_NilVSSJustification verifies that justifications
// with nil VSSJustification are handled gracefully (using zero recipient index).
func TestDeduplicateJustifications_NilVSSJustification(t *testing.T) {
	t.Parallel()

	j1 := types.Justification{Index: 1}
	j2 := types.Justification{Index: 1}

	result := deduplicateJustifications([]types.Justification{j1, j2})
	require.Len(t, result, 1, "both have (1,0) key so second should be deduplicated")
}

// TestDeduplicateJustifications_PreservesOrder verifies that the first occurrence
// is kept and order is preserved.
func TestDeduplicateJustifications_PreservesOrder(t *testing.T) {
	t.Parallel()

	j1 := types.Justification{Index: 3, VssJustification: &types.VSSJustification{
		PlainDeal: &types.PlainDeal{SecShare: &types.SecShare{I: 1}},
	}}
	j2 := types.Justification{Index: 1, VssJustification: &types.VSSJustification{
		PlainDeal: &types.PlainDeal{SecShare: &types.SecShare{I: 2}},
	}}
	j3 := types.Justification{Index: 2, VssJustification: &types.VSSJustification{
		PlainDeal: &types.PlainDeal{SecShare: &types.SecShare{I: 3}},
	}}

	result := deduplicateJustifications([]types.Justification{j1, j2, j3})
	require.Len(t, result, 3, "all unique keys, nothing deduplicated")
	require.Equal(t, uint32(3), result[0].Index)
	require.Equal(t, uint32(1), result[1].Index)
	require.Equal(t, uint32(2), result[2].Index)
}

// TestInvalidateDealerRegistration_AlreadyInvalidated verifies idempotency:
// re-invalidating an already-invalidated dealer is a no-op.
func TestInvalidateDealerRegistration_AlreadyInvalidated(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	round := uint32(2)
	dealerAddr := common.HexToAddress("0xEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
	dealerIndex := uint32(5)

	setupDealerRegistrationForInvalidation(t, k, ctx, round, dealerAddr, dealerIndex, types.DKGRegStatusInvalidated)

	network := &types.DKGNetwork{
		Round: round,
	}

	err := k.invalidateDealerRegistration(ctx, network, dealerIndex)
	require.NoError(t, err, "re-invalidating should be a no-op")

	reg, err := k.getDKGRegistration(ctx, round, dealerAddr)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusInvalidated, reg.Status)
}

// setupDealerRegistrationForInvalidation is a helper that creates a registration
// with a specific status for use in invalidation tests.
func setupDealerRegistrationForInvalidation(t *testing.T, k *Keeper, ctx context.Context, round uint32, addr common.Address, index uint32, status types.DKGRegStatus) {
	t.Helper()

	reg := &types.DKGRegistration{
		Round:         round,
		ValidatorAddr: addr.Hex(),
		Index:         index,
		DkgPubKey:     []byte("dkg-pub"),
		CommPubKey:    []byte("comm-pub"),
		Status:        status,
	}
	require.NoError(t, k.setDKGRegistration(ctx, addr, reg))
}
